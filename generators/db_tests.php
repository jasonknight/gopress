<?php
// in order for this to work user'll need to
// have a special yaml file
$fail = "t.Errorf";
$cnf = "../gopress.db.yml";
$fields = array();
foreach ($t->fields as $f) {
    if ( $f->Key == 'PRI') {
        // skip primary key
        continue;
    }
    $goname = convertFieldName($f->Field);
    $dbname = $f->Field;
    $gotype = $f->go_type;
    if ( $gotype == "string") {
        $tval = "the rain in spain";
    } else if (preg_match("/^int/",$gotype) ) {
        $tval = 999;
    } else if (preg_match("/^float/",$gotype) ) {
        $tval = 3.145;
    } else if ( $gotype == "DateTime") {
        $tval = "2016-10-11 03:05:21.4Z";
    }
    $fields[] = array(
        'type' => $gotype,
        'name' => $goname,
        'dbname' => $dbname,
        'value' => $tval
    );
}
$txt = "
func Test{$t->model_name}Create(t *testing.T) {
    if fileExists(`$cnf`) {
    a,err := NewMysqlAdapterEx(`$cnf`)
    model := New{$t->model_name}(a)
";
$i = 0;
foreach ($fields as $f) {
    if ($f['type'] == '*DateTime') {
        $txt .= "
    d{$i} := NewDateTime()
    d{$i}.FromString(`{$f['value']}`)
    model.{$f['name']} = d{$i}
        ";
    } else if ($f['type'] == 'string') {
        $txt .= "
    model.{$f['name']} = `{$f['value']}`
        ";
    } else {
        $txt .= "
    model.{$f['name']} = {$f['value']}
        ";        
    }
    $i++;
}
$txt .= "
    i,err := model.Create()
    if err != nil {
        $fail(`failed to create model %+v error: %s`,model,err)
        return
    }
    if i == 0 {
        $fail(`zero affected rows`)
        return
    }
";

$txt .= "}
};\n";

puts($txt);