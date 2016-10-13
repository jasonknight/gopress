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
    if ( preg_match("/^post_/",$n) ) {
        $n = substr($n,5);
    }
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
foreach ($tables as &$t) {
    $t->fields = array();
    $res = mysql_query("DESCRIBE {$t->database_name};");
    while ( $row = mysql_fetch_object($res)) {
        $row->go_type = mysqlToGoType($row->Type);
        if ( $row->Key == "PRI" ) {
            $t->pfield = $row;
        }
        $row->model_field_name = convertFieldName($row->Field);
        $row->mysql_fmt_type = mysqlToFmtType($row->Type);
        $row->dirty_marker = "Is" . convertFieldName($row->Field) . "Dirty"; 
        $t->fields[] = $row;
    }
}

// Now we loop over and create all the models

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
