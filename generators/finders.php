<?php
$from_map_body = "";
$from_model_body = "";
foreach($t->fields as $f) {
    $fname = convertFieldName($f->Field);
    $field_name = convertFieldName($f->Field);
    $fname = "FindBy{$fname}";
    $arg = "_find_by_" . $f->model_field_name;
    $argtype = $f->go_type;
    $fmt_type = mysqlToFmtType($f->Type);
    if (isPrimaryKey($f)) {
        $fname = "Find";
        $rtype = "bool"; //i.e. we set the current model
    } else {
        $rtype = "[]*{$t->model_name}";
    }
    $scol = $f->Field;
    // these are here just to save having to loop
    if ( $f->go_type == "*DateTime" ) {
        $from_map_body .= "\t_" . $f->model_field_name . ",err := m[\"{$f->Field}\"].As" . ucfirst(substr($f->go_type,1)). "()\n";
    } else {
        $from_map_body .= "\t_" . $f->model_field_name . ",err := m[\"{$f->Field}\"].As" . ucfirst($f->go_type). "()\n";
    }
    $from_model_body .= "\to.{$f->model_field_name} = m.{$f->model_field_name}\n";
    
    $from_map_body .= "\tif err != nil {\n \t\treturn err\n\t}\n";
    $from_map_body .= "\to." . $f->model_field_name . " = _" . $f->model_field_name . "\n";
    
    if ( $fname == "Find" ) {
        $failure_return = "return false,err";
    } else {
        $failure_return = "return model_slice,err";
    }
    $sig = "func (o *{$t->model_name}) $fname($arg $argtype) ($rtype,error) {";
$body = "
    var model_slice []*{$t->model_name}
    q := fmt.Sprintf(\"SELECT * FROM %s WHERE `%s` = '$fmt_type'\",o._table, \"{$f->Field}\", $arg)
    results, err := o._adapter.Query(q)
    if err != nil {
        $failure_return
    }
    
    for _,result := range results {
        ro := {$t->model_name}{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            $failure_return
        }
        model_slice = append(model_slice,&ro)
    }
";
if ( $fname == "Find" ) {
    // we return the 0th element
    $body .= "
    if len(model_slice) == 0 {
        // there was an error!
        return false, errors.New(\"not found\")
    }
    o.From{$t->model_name}(model_slice[0])
    return true,nil
";
} else {
    $body .= "
    if len(model_slice) == 0 {
        // there was an error!
        return nil, errors.New(\"no results\")
    }
    return model_slice,nil
";
}
    puts($sig);
    puts($body);
    puts("}");
}
puts("
func (o *{$t->model_name}) FromDBValueMap(m map[string]DBValue) error {
$from_map_body
 \treturn nil
}
func (o *{$t->model_name}) From{$t->model_name}(m *{$t->model_name}) {
$from_model_body
}
");
