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
    if ( preg_match("/woocommerce_/",$t->database_name)) {
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
function mysqlToGoRandom($t) {
    if (preg_match("/varchar|text/", $t) ) {
        preg_match("/varchar\((\d+)\)/",$t,$m);
        if (!empty($m)) {
            $m[1] = intval($m[1]);
            if ( $m[1] > 20 ) {
                $m[1] = 20;
            }
            $m[1] -= 1;
            $tval = "randomString({$m[1]})";
        } else {
            $tval = "randomString(25)";
        }
        return $tval;
    }
    if (preg_match("/bigint/", $t) ) {
        return "int64(randomInteger())";
    }
    if (preg_match("/int/", $t) ) {
        return "int(randomInteger())";
    }
    if ( $t == "longtext" || $t == "tinytext" || $t == "mediumtext") {
        return "randomString(35)";
    }

    if ( $t == "datetime" ) {
        return "randomDateTime(a)";
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
function getFieldDef($needle,$index=false) {
    global $tables;
    for ($i = 0; $i < count($tables); $i++ ) {
        for ($j = 0; $j < count($tables[$i]->fields);$j++) {
            if ( $tables[$i]->fields[$j]->Field == $needle) {
                if ( $index ) {
                    return array($i,$j);
                }
                return $tables[$i]->fields[$j];
            }
        }
    }
}
function pluralize($s) {
    if ($s[strlen($s)-1] == 'y') {
        return substr($s, 0,strlen($s)-1) . "ies";
    }
    return $s . "s";
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
        $row->go_random = mysqlToGoRandom($row->Type);
        $row->model_field_name = convertFieldName($row->Field);
        $row->mysql_fmt_type = mysqlToFmtType($row->Type);
        $row->dirty_marker = "Is" . convertFieldName($row->Field) . "Dirty"; 
        if ( isPrimaryKey($row) ) {
            $t->pfield = $row;
        }
        $t->fields[] = $row;
    }
}
// foreach ($tables as &$t) {
//     foreach ($t->fields as &$field) {
//         if (preg_match("/([\w_]+)_[idID]{2}$/",$field->Field) && isPrimaryKey($field) == false) {
//         //if ( $field->Field == "comment_post_ID") {
//             include "setup_habtm.php";
//         } else {
//            // echo "Failed for: {$field->Field}\n";
//         }
//     }
// }
// print_r($tables);
// die("Done\n");
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
    _ "github.com/go-sql-driver/mysql" // This is standard for this library.
    "strconv"
    "gopkg.in/yaml.v2"
    "regexp"
    "errors"
    "os"
    "io"
    "io/ioutil"
    "bufio"
    "log"
    "strings"
)

');

include "generic_bits.php";
//include "model_interface.php";
$_ii = 0;
$seen = array();
foreach ($tables as $t) {
    if (in_array($t->model_name, $seen) ) {
        continue;
    }
    $seen[] = $t->model_name;
puts("// {$t->model_name} is a Object Relational Mapping to
// the database table that represents it. In this case it is
// {$t->dname}. The table name will be Sprintf'd to include
// the prefix you define in your YAML configuration for the
// Adapter.
type {$t->model_name} struct {
    _table string
    _adapter Adapter
    _pkey string // $_ii The name of the primary key in this table
    _conds []string
    _new bool");
    include "models/interface_fields.php";
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
// New{$t->model_name} binds an Adapter to a new instance
// of {$t->model_name} and sets up the _table and primary keys
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
//include "models/interface_funcs.php";
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
    \"math/rand\"
    \"os\"
    \"time\"
    \"bytes\"
    \"regexp\"
    \"bufio\"
    \"errors\"
)
var letters = []rune(\"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ\")

func randomString(n int) string {
    rand.Seed(time.Now().UnixNano())
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}
func randomInteger() int {
    rand.Seed(time.Now().UnixNano())
    x := rand.Intn(10000) + 100
    if x == 0 {
        return randomInteger();
    }
    return x + 100
}
func randomFloat() float32 {
    rand.Seed(time.Now().UnixNano())
    return rand.Float32() * 100
}
func randomDateTime(a Adapter) *DateTime {
    rand.Seed(time.Now().UnixNano())
    d := NewDateTime(a)
    d.Year = rand.Intn(2017)
    d.Month = rand.Intn(11)
    d.Day = rand.Intn(28)
    d.Hours = rand.Intn(23)
    d.Minutes = rand.Intn(59)
    d.Seconds = rand.Intn(56)
    if d.Year < 1000 {
        d.Year = d.Year + 1000
    }
    return d
}
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
