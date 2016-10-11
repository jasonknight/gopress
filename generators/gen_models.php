<?php 
include "../../db.php";
$con =  mysql_connect($host,$user,$pass);
        mysql_select_db($database);

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
        $n = substr($n,2);
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
        return "DateTime";
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
        return "string";
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
        $t->fields[] = $row;
    }
}

// Now we loop over and create all the models

function puts($s) {
    echo $s . "\n";
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
)');
$_ii = 0;
$seen = array();
foreach ($tables as $t) {
    if (in_array($t->model_name, $seen) ) {
        continue;
    }
    $seen[] = $t->model_name;
puts("
type {$t->model_name} struct {
    _table_ string
    _adapter_ Adapter
    _pkey_ string // $_ii The name of the primary key in this table
    _conds_ []string");
    
    foreach($t->fields as $f) {
        $fname = lcfirst(convertFieldName($f->Field));
    puts("    {$fname} {$f->go_type}");
    }
    puts("}");
$newfunc = "
func New{$t->model_name}(a Adapter) *{$t->model_name} {
    var o {$t->model_name}
    o._table_ = \"{$t->database_name}\"
    o._adapter_ = a
    o._pkey_ = \"{$t->pfield->Field}\"
    return &o
}
";
puts($newfunc);
include "finders.php";
//include "crud.php";
}

puts("
type Adapter interface {
    Open(string,string,string,string) error
    Close()
    Query(string) ([]map[string]DBValue,error)
    Execute(string) error
    LastInsertedId() int64
    AffectedRows() int64
}

type DBValue interface {
    AsInt() (int,error)
    AsInt64() (int64,error)
    AsFloat32() (float32,error)
    AsString() (string,error)
    AsDateTime() (DateTime,error)
    SetInternalValue(string,string)
}

type MysqlValue struct {
    _v string
    _k string
}
func (v *MysqlValue) SetInternalValue(key,value string) {
    v._v = value
    v._k = key

}
func (v *MysqlValue) AsString() (string,error) {
    return v._v,nil
}
func (v *MysqlValue) AsInt() (int,error) {
    i,err := strconv.ParseInt(v._v,10,32)
    return int(i),err
}
func (v *MysqlValue) AsInt64() (int64,error) {
    i,err := strconv.ParseInt(v._v,10,64)
    return i,err
}
func (v *MysqlValue) AsFloat32() (float32,error) {
    i,err := strconv.ParseFloat(v._v,32)
    if err != nil {
        return 0.0,err
    }
    return float32(i),err
}
func (v *MysqlValue) AsFloat64() (float64,error) {
    i,err := strconv.ParseFloat(v._v,64)
    if err != nil {
        return 0.0,err
    }
    return i,err
}

func (v *MysqlValue) AsDateTime() (DateTime,error) {
    var dt DateTime
    err := dt.FromString(v._v)
    if err != nil {
        return DateTime{}, err
    }
    return dt,nil
}

type MysqlAdapter struct {
    Host string `yaml:\"host\"`
    User string `yaml:\"user\"`
    Pass string `yaml: \"pass\"`
    Database string `yaml:\"database\"`
    _conn_ *sql.DB
    _lid int64
    _cnt int64
}

func (a *MysqlAdapter) FromYAML(b []byte) error {
    return yaml.Unmarshal(b,a)
}

func (a *MysqlAdapter) Open(h,u,p,d string) error {
    if ( h != \"localhost\") {
        tc, err := sql.Open(\"mysql\",fmt.Sprintf(\"%s:%s@tcp(%s)/%s\",u,p,h,d))
        if err != nil {
            return err
        }
        a._conn_ = tc
    } else {
        tc, err := sql.Open(\"mysql\",fmt.Sprintf(\"%s:%s@/%s\",u,p,d))
        if err != nil {
            return err
        }
        a._conn_ = tc
    }
    return nil

}
func (a *MysqlAdapter) Close() {
    a._conn_.Close()
}

func (a *MysqlAdapter) Query(q string) ([]map[string]DBValue,error) {
    results := new([]map[string]DBValue)
    rows, err := a._conn_.Query(q)
    if err != nil {
        return nil,err
    }
    defer rows.Close()
    columns, err := rows.Columns()
    if err != nil {
        return nil, err
    }
    values := make([]sql.RawBytes, len(columns))
    scanArgs := make([]interface{},len(values))
    for i := range values {
        scanArgs[i] = &values[i]
    }
    for rows.Next() {
        err = rows.Scan(scanArgs...)
        if err != nil {
            return nil,err
        }
        res := make(map[string]DBValue)
        for i,col := range values {
            k := columns[i]
            res[k].SetInternalValue(k,string(col))
        }
        *results = append(*results,res)
    }
    return *results,nil
}
func (a *MysqlAdapter) Execute(q string) error {
    tx, err := a._conn_.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback();
    stmt, err := tx.Prepare(q)
    if err != nil {
        return err
    }
    defer stmt.Close()
    res,err := stmt.Exec(q)
    if err != nil {
        return err
    }
    a._lid,err = res.LastInsertId()
    if err != nil {
        return err
    }
    a._cnt,err = res.RowsAffected()
    if err != nil {
        return err
    }
    err = tx.Commit()
    return err
}
func (a *MysqlAdapter) LastInsertedId() int64 {
    return a._lid
}
func (a *MysqlAdapter) AffectedRows() int64 {
    return a._cnt
}
type DateTime struct {
    Day int
    Month int
    Year int
    Hours int
    Minutes int
    Seconds int
    Zone string
    Offset int
}
func (d *DateTime) FromString(s string) error {
    es := s
    re := regexp.MustCompile(\"(?P<year>[\\\d]{4})-(?P<month>[\\\d]{2})-(?P<day>[\\\d]{2}) (?P<hours>[\\\d]{2}):(?P<minutes>[\\\d]{2}):(?P<seconds>[\\\d]{2})\\\.(?P<offset>[\\\d]+)(?P<zone>[\\\w]+)\")
    n1 := re.SubexpNames()
    ir2 := re.FindAllStringSubmatch(es, -1)
    if len(ir2) == 0 {
        return errors.New(fmt.Sprintf(\"found now data to capture in %s\",es))
    }
    r2 := ir2[0]
    for i, n := range r2 {
        if n1[i] == \"year\" {
            _Year,err := strconv.ParseInt(n,10,32)
            d.Year = int(_Year)
            if err != nil {
                return errors.New(fmt.Sprintf(\"failed to convert %s in %s received %s\",n[i],es,err))
            }
        }
        if n1[i] == \"month\" {
            _Month,err := strconv.ParseInt(n,10,32)
            d.Month = int(_Month)
            if err != nil {
                return errors.New(fmt.Sprintf(\"failed to convert %s in %s received %s\",n[i],es,err))
            }
        }
        if n1[i] == \"day\" {
            _Day,err := strconv.ParseInt(n,10,32)
            d.Day = int(_Day)
            if err != nil {
                return errors.New(fmt.Sprintf(\"failed to convert %s in %s received %s\",n[i],es,err))
            }
        }
        if n1[i] == \"hours\" {
            _Hours,err := strconv.ParseInt(n,10,32)
            d.Hours = int(_Hours)
            if err != nil {
                return errors.New(fmt.Sprintf(\"failed to convert %s in %s received %s\",n[i],es,err))
            }
        }
        if n1[i] == \"minutes\" {
            _Minutes,err := strconv.ParseInt(n,10,32)
            d.Minutes = int(_Minutes)
            if err != nil {
                return errors.New(fmt.Sprintf(\"failed to convert %s in %s received %s\",n[i],es,err))
            }
        }
        if n1[i] == \"seconds\" {
            _Seconds,err := strconv.ParseInt(n,10,32)
            d.Seconds = int(_Seconds)
            if err != nil {
                return errors.New(fmt.Sprintf(\"failed to convert %s in %s received %s\",n[i],es,err))
            }
        }
        if n1[i] == \"offset\" {
            _Offset,err := strconv.ParseInt(n,10,32)
            d.Offset = int(_Offset)
            if err != nil {
                return errors.New(fmt.Sprintf(\"failed to convert %s in %s received %s\",n[i],es,err))
            }
        }
        if n1[i] == \"zone\" {
            d.Zone = n
        }
    }
    return nil
}
func (d *DateTime) ToString() string {
    return fmt.Sprintf(\"%d-%d-%d %d:%d:%d.%d.%s\",d.Year,d.Month,d.Day,d.Hours,d.Minutes,d.Seconds,d.Offset,d.Zone)
}
");
//print_r($tables);
