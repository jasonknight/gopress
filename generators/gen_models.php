<?php 
include "../../db.php";
$con =  mysql_connect($host,$user,$pass);
        mysql_select_db($database);
$_contents = "";
function maybeLC($str) {
    return $str;
}
function convertTableName($n) {
    if ( preg_match("/^wp_/",$n) ) {
        $n = substr($n,2);
    }
    $camel = str_replace(' ', '', ucwords(str_replace('_', ' ', $n)));
    $camel = str_replace("meta","Meta",$camel);
    $camel = str_replace("Woocommerce","Woo",$camel);
    $camel = str_replace("Attribute","Attr",$camel);
    $camel = str_replace("Permissions","Perms",$camel);
    if ( $camel[strlen($camel) -1 ] == "s" ) {
        $camel = substr($camel,0,strlen($camel) - 1);
    }
    return $camel;
}
function convertFieldName($n) {
    $camel = str_replace(' ', '', ucwords(str_replace('_', ' ', $n)));
    $camel = str_replace("meta","Meta",$camel);
    $camel = str_replace("Woocommerce","Woo",$camel);
    $camel = str_replace("Attribute","Attr",$camel);
    $camel = str_replace("Permissions","Perms",$camel);
    return $camel;
}
// First we get a list of all the tables
$tables = array();
$res = mysql_query("SHOW TABLES;");
while($row = mysql_fetch_object($res)) {
    $t = new stdClass();
    $t->database_name = $row->Tables_in_philatelic_wp2;
    $t->dname = preg_replace("/^wp_/","",$t->database_name);
    if ( preg_match("/ign_/",$t->database_name)) {
        continue;
    }
    $t->model_name = convertTableName($t->database_name);
    $tables[] = $t;
}

// Second, we iterate over tables and describe each
function mysqlToGoType($t) {
    if (preg_match("/varchar|text/", $t) ) {
        return "string";
    }
    if (preg_match("/bigint/", $t) ) {
        return "int64";
    }
    if (preg_match("/int/", $t) ) {
        return "int";
    }
    if ( $t == "longtext" || $t == "tinytext" || $t == "mediumtext") {
        return "string";
    }

    if ( $t == "datetime" ) {
        return "*DateTime";
    }
}
function mysqlToFmtType($t) {
    if (preg_match("/varchar|text/", $t) ) {
        return "%s";
    }
    if (preg_match("/bigint/", $t) ) {
        return "%d";
    }
    if (preg_match("/int/", $t) ) {
        return "%d";
    }
    if ( $t == "longtext" || $t == "tinytext" || $t == "mediumtext") {
        return "%s";
    }

    if ( $t == "datetime" ) {
        return "%s";
    }
}
function isPrimaryKey($p) {
    if ( $p->Field == "object_id" ) {
        return false;
    }
    if ( $p->Key != "PRI" ) {
        return false;
    }
    return true;
}
foreach ($tables as &$t) {
    $t->fields = array();
    $t->belongs_to = array();
    $t->has_many = array();
    $res = mysql_query("DESCRIBE {$t->database_name};");
    while ( $row = mysql_fetch_object($res)) {
        $row->go_type = mysqlToGoType($row->Type);
        $row->model_field_name = convertFieldName($row->Field);
        $row->mysql_fmt_type = mysqlToFmtType($row->Type);
        $row->dirty_marker = "Is" . convertFieldName($row->Field) . "Dirty"; 
        if (preg_match("/\w+_id/",$row->Field) &&  isPrimaryKey($row) == false) {
            $bt = belongsTo($row);
            
            $bt->model_name = $t->model_name;
            foreach ($tables as &$tt) {
                if ( $bt->model == $tt->model_name) {
                    // i.e. only add it if there
                    // is a corresponding definition
                    $bt->ftable = $t->dname;
                    $bt->fmodel = $t->model_name;
                    $t->belongs_to[] = $bt;

                    break;
                }
            }
        }
        if ( isPrimaryKey($row) ) {
            $t->pfield = $row;
        }
        $t->fields[] = $row;
    }
}
foreach ($tables as &$t) {

    echo "Considering {$t->model_name}\n";
    if ( !empty($t->belongs_to)) {
        foreach ($t->belongs_to as $bt) {
            echo "\tBelongsTo: {$bt->model}\n";
            $i = getTableDef($bt->model,true);
            if ( $i != -1 ) {
                echo "\tAdding has_many: {$tables[$i]->model_name}\n";
                $tables[$i]->has_many[] = hasMany($bt,$tables[$i]);
            }
            
        }
    }
}
function getTableDef($model_name,$index=false) {
    global $tables;
    $i = 0;
    foreach ($tables as $t) {
        if ($t->model_name == $model_name) {
            if ( $index ) {
                return $i;
            }
            return $t;
        }
        $i++;
    }
    return -1;
}
function belongsTo($f) {
    preg_match("/([\w_]+)_id/",$f->Field,$m);
    $nf = new stdClass();
    $nf->model = convertTableName($m[1]);
    $nf->go_type = "*{$nf->model}";
    $nf->model_field_name = convertFieldName($f->Field);
    $nf->fkey = $f->Field;
    $nf->fkey_type = $f->go_type;
    $nf->fkey_myfmt = $f->mysql_fmt_type;
    return $nf;
}

function hasMany($bt,$t) {
    $f = $bt;
     $has_name = $f->model_name . "s";
     $has_type = "[]*{$f->model_name}";
    // $ft = getTableDef($f->model_name);
    // $has_fkey = $ft->pfield->Field;
    // $has_fkey_type = $ft->pfield->go_type;
    // $has_fkey_myfmt = $ft->pfield->mysql_fmt_type;

    $hm = new stdClass();
    $hm->model_name = $f->model_name;
    $hm->model_field_name = convertFieldName($f->fkey);
    $hm->name = $has_name;
    $hm->table = $f->ftable;
    $hm->type = $has_type;
    $hm->fkey = $f->fkey;
    $hm->fkey_type = $f->fkey_type;
    $hm->fkey_myfmt = $f->fkey_myfmt;
    return $hm;
}
// Now we loop over and create all the models
foreach ($tables as $t) {
    echo "{$t->model_name}: \n";
    foreach ($t->has_many as $hm) {
        echo "\t Has Many: {$hm->name}\n";
    }
    foreach ($t->belongs_to as $bt) {
        echo "\t Belongs_to: {$bt->model}\n";
    }
}
function puts($s) {
    global $_contents;
    $_contents .= $s . "\n";
}

puts("package gopress");
puts('import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "strconv"
    "gopkg.in/yaml.v2"
    "regexp"
    "errors"
    "os"
    "io/ioutil"
    "bufio"
)
');

include "generic_bits.php";

$_ii = 0;
$seen = array();
foreach ($tables as $t) {
    if (in_array($t->model_name, $seen) ) {
        continue;
    }
    $seen[] = $t->model_name;
puts("
type {$t->model_name} struct {
    _table string
    _adapter Adapter
    _pkey string // $_ii The name of the primary key in this table
    _conds []string
    _new bool");
    
    foreach($t->fields as $f) {
        $fname = $f->model_field_name;
        puts("    {$fname} {$f->go_type}");
    }
    puts("\t// Dirty markers for smart updates");
    foreach($t->fields as $f) {
        puts("    {$f->dirty_marker} bool");
    }
    puts("\t// Relationships");
    foreach ($t->belongs_to as $bt) {
        puts("\t{$bt->model} {$bt->go_type}");
        puts("\tIs{$bt->model}Loaded bool");
    }
    foreach ($t->has_many as $hm) {
        puts("\t{$hm->name} {$hm->type}");
        puts("\tIs{$hm->name}Loaded bool");
    }
    puts("}");
$newfunc = "
func New{$t->model_name}(a Adapter) *{$t->model_name} {
    var o {$t->model_name}
    o._table = fmt.Sprintf(\"%s{$t->dname}\",a.DatabasePrefix())
    o._adapter = a
    o._pkey = \"{$t->pfield->Field}\"
    o._new = false
    return &o
}
";
puts($newfunc);
include "habtm.php";
include "getset.php";
include "finders.php";
include "crud.php";

}



file_put_contents($_SERVER['argv'][1] . ".go",$_contents);
$_contents = "";

// Now we need to generate tests
puts("package gopress
import (
    \"testing\"
    \"strconv\"
)
");
$seen = array();
foreach ($tables as $t) {
    if (in_array($t->model_name, $seen) ) {
        continue;
    }
    $seen[] = $t->model_name;
    include "model_test.php";
    include "db_tests.php";
}
include "generic_tests.php";
file_put_contents($_SERVER['argv'][1]."_test.go",$_contents);
$_contents = "";
