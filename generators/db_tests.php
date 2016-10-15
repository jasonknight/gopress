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
    if ( isPrimaryKey($f)) {
        // skip primary key
        continue;
    }
    $goname = convertFieldName($f->Field);
    $dbname = $f->Field;
    $gotype = $f->go_type;
    if ( $gotype == "string") {
        $tval = randomString();
    } else if (preg_match("/^int/",$gotype) ) {
        $tval = randomInt();
    } else if (preg_match("/^float/",$gotype) ) {
        $tval = randomInt() / 3.14;
    } else if ( $gotype == "DateTime") {
        $tval = "2016-10-11 03:05:21.4Z";
    }
    $fields[] = array(
        'type' => $gotype,
        'name' => $goname,
        'dbname' => $dbname,
        'value' => $tval,
        'fmt' => $f->mysql_fmt_type
    );
}
$txt = "
func Test{$t->model_name}Create(t *testing.T) {
    if fileExists(`$cnf`) {
    a,err := NewMysqlAdapterEx(`$cnf`)
    if err != nil {
        $fail(`could not load $cnf %s`,err)
        return
    }
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
        $fail(`failed to create model %s`,err)
        return
    }
    if i == 0 {
        $fail(`zero affected rows`)
        return
    }
    model2 := New{$t->model_name}(a)
    found,err := model2.Find(model.GetPrimaryKeyValue())
    if err != nil {
        $fail(`did not find record for %s = {$t->pfield->mysql_fmt_type} because of %s`,model.GetPrimaryKeyName(),model.GetPrimaryKeyValue(),err)
        return
    }
    if found == false {
        $fail(`did not find record for %s = {$t->pfield->mysql_fmt_type}`,model.GetPrimaryKeyName(),model.GetPrimaryKeyValue())
        return
    }

";
$i = 0;
foreach ($fields as $f) {
   if ( $f['type'] == "*DateTime") {
   $txt .= "
    if (model.{$f['name']}.Year != model2.{$f['name']}.Year ||
        model.{$f['name']}.Month != model2.{$f['name']}.Month ||
        model.{$f['name']}.Day != model2.{$f['name']}.Day ||
        model.{$f['name']}.Hours != model2.{$f['name']}.Hours ||
        model.{$f['name']}.Minutes != model2.{$f['name']}.Minutes ||
        model.{$f['name']}.Seconds != model2.{$f['name']}.Seconds ) {
        $fail(`model.{$f['name']} != model2.{$f['name']}`)
        return
    }
";
    continue;
   }
   $txt .= "
    if model.{$f['name']} != model2.{$f['name']} {
        $fail(`model.{$f['name']}[{$f['fmt']}] != model2.{$f['name']}[{$f['fmt']}]`,model.{$f['name']},model2.{$f['name']})
        return
    }
";
    $i++;
}
$txt .= "} // end of if fileExists
};\n";

puts($txt);