<?php
// in order for this to work user'll need to
// have a special yaml file
$fail = "t.Errorf";
$cnf = "../gopress.db.yml";
$fields = array();
if (! function_exists("randomString")) {
function randomString() {
    $c = "bcdfghjklmnpqrstvwxyz";
    $v = "aeiou";
    $str = "";
    for ($i = 0; $i < 10;$i++) {
        $t = $c[rand(0,strlen($c)-1)] . $v[rand(0,strlen($v)-1)];
        if (rand(0,10) < 5) {
            ucfirst($t);
        }
        $str .= $t;
    }
    return $str;
}
function randomInt() {
    return rand(1000,50000) + rand(50,rand(1000,9000)) - rand(0,10000);
}
}
foreach ($t->fields as $f) {
    if ( isPrimaryKey($f) && $t->model_name != "TermRelationship") {
        // skip primary key
        continue;
    }
    $goname = convertFieldName($f->Field);
    $dbname = $f->Field;
    $gotype = $f->go_type;
    $tval = $f->go_random;
    $fields[] = array(
        'type' => $gotype,
        'name' => $goname,
        'dbname' => $dbname,
        'value' => $tval,
        'fmt' => $f->mysql_fmt_type
    );
}
include "_create_model_test.php";
include "_model_update_test.php";

puts($txt);