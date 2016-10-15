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
    if ( $gotype == "string") {
        preg_match("/varchar\((\d+)\)/",$f->Type,$m);
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
        
    } else if (preg_match("/^int/",$gotype) ) {
        if ($gotype == "int64") {
            $tval = "int64(randomInteger())";
        } else {
            $tval = "randomInteger()";
        }
    } else if (preg_match("/^float/",$gotype) ) {
        $tval = "randomFloat()";
    } else if ( $gotype == "*DateTime") {
        $tval = "randomDateTime(a)";
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
    file, err := os.OpenFile(\"adapter.log\", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        $fail(\"Failed to open log file %s\", err)
    }
    a.SetLogs(file)
    model := New{$t->model_name}(a)
";
$i = 0;
foreach ($fields as $f) {
    $txt .= "model.{$f['name']} = {$f['value']}\n";
    $i++;
}
$find_line = "model2.Find(model.GetPrimaryKeyValue())";
if ($t->model_name == "TermRelationship") {
    $find_line = "model2.Find(model.GetTermTaxonomyId(),model.GetObjectId())";
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
    found,err := $find_line
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
        $fail(`model.{$f['name']} != model2.{$f['name']} %+v --- %+v`,model.{$f['name']},model2.{$f['name']})
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