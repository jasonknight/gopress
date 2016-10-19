<?php
$txt = "
func Test{$t->model_name}Updaters(t *testing.T) {
    if fileExists(`../gopress.db.yml`) == false {
        return
    }
    a,err := NewMysqlAdapterEx(`../gopress.db.yml`)
    if err != nil {
        $fail(`could not load $cnf %s`,err)
        return
    }
    file, err := os.OpenFile(\"adapter.log\", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        $fail(\"Failed to open log file %s\", err)
        return
    }
    a.SetLogs(file)
    model := New{$t->model_name}(a)
";

$tests = "";
$kls = $t->model_name;
$i = 0;
foreach ($t->fields as $f) {
    if ( isPrimaryKey($f) ) {
        continue;
    }
    
    $fname = convertFieldName($f->Field);

// so we need to set all of them
    $txt .= "
    model.Set{$fname}({$f->go_random})
    if model.Get{$fname}() != model.{$fname} {
        $fail(`$kls.Get{$fname}() != $kls.{$fname}`)
    }
    if model.Is{$fname}Dirty != true {
        $fail(`$kls.Is{$fname}Dirty != true`)
        return
    }
    
    u{$i} := {$f->go_random}
    _,err = model.Update{$fname}(u{$i})
    if err != nil {
        $fail(`failed Update{$fname}(u{$i}) %s`,err)
        return
    }

    if model.Get{$fname}() != u{$i} {
        $fail(`$kls.Get{$fname}() != u{$i} after Update{$fname}`)
        return
    }
    model.Reload()
    if model.Get{$fname}() != u{$i} {
        $fail(`$kls.Get{$fname}() != u{$i} after Reload`)
        return
    }
";
    $i++;
} 

$txt .= "
};\n";
puts($txt);