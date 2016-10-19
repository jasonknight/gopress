<?php
$txt = "";
$txt .= "
func Test{$t->model_name}Create(t *testing.T) {
    if fileExists(`$cnf`) {
    a,err := NewMysqlAdapterEx(`$cnf`)
    defer a.Close()
    if err != nil {
        $fail(`could not load $cnf %s`,err)
        return
    }
    file, err := os.OpenFile(\"adapter.log\", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        $fail(\" Failed to open log file %s\", err)
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
$find_fail_line = "$fail(` did not find record for %s = {$t->pfield->mysql_fmt_type} because of %s`,model.GetPrimaryKeyName(),model.GetPrimaryKeyValue(),err)";
if ($t->model_name == "TermRelationship") {
    $find_line = "model2.Find(model.TermTaxonomyId,model.ObjectId)\nif model.TermTaxonomyId == 0 {\n$fail(`it's 0`)\n} \n";
    $find_fail_line = "$fail(` did not find record for term_taxonomy_id = %d AND object_id = %d because of %s`,model.TermTaxonomyId,model.ObjectId,err)";
}
$txt .= "
    err = model.Create()
    if err != nil {
        $fail(` failed to create model %s`,err)
        return
    }

    model2 := New{$t->model_name}(a)
    found,err := $find_line
    if err != nil {
        $find_fail_line
        return
    }
    if found == false {
        $find_fail_line
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
        $fail(`2: model.{$f['name']} != model2.{$f['name']} %+v --- %+v`,model.{$f['name']},model2.{$f['name']})
        return
    }
";
    continue;
   }
   $txt .= "
    if model.{$f['name']} != model2.{$f['name']} {
        $fail(` model.{$f['name']}[{$f['fmt']}] != model2.{$f['name']}[{$f['fmt']}]`,model.{$f['name']},model2.{$f['name']})
        return
    }
";
    $i++;
}
foreach ($fields as $f) {
    $txt .= "model2.Set{$f['name']}({$f['value']})\n";
    $i++;
}
$txt .= "
    err = model2.Save()
    if err != nil {
        $fail(`failed to save model2 %s`,err)
    }
";
foreach ($fields as $f) {
   if ( $f['type'] == "*DateTime") {
   $txt .= "
    if (model.{$f['name']}.Year == model2.{$f['name']}.Year) {
        $fail(` model.{$f['name']}.Year == model2.{$f['name']} but should not!`)
        return
    }
";
    continue;
   }
   $txt .= "
    if model.{$f['name']} == model2.{$f['name']} {
        $fail(`1: model.{$f['name']}[{$f['fmt']}] != model2.{$f['name']}[{$f['fmt']}]`,model.{$f['name']},model2.{$f['name']})
        return
    }
";
    $i++;
}

// now we need to test all the findBy methods
foreach ($fields as $f) {
    if ($t->model_name == "TermRelationship") {
        if ($f['name'] == 'TermTaxonomyId' || $f['name'] == 'ObjectId') {
            continue;
        } 
        
    }
    if (preg_match("/Order/",$f['name'])) {
        continue;
    }
    echo "Adding in FindBys\n";
$txt .= "
    res{$i},err := model.FindBy{$f['name']}(model2.Get{$f['name']}())
    if err != nil {
        $fail(`failed model.FindBy{$f['name']}(model2.Get{$f['name']}())`)
    }
    if len(res{$i}) == 0 {
        $fail(`failed to find any {$t->model_name}`)
    }
"; 
$i++;
}


$txt .= "} // end of if fileExists
};\n";
puts($txt);