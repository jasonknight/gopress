<?php
$from_map_body = "";
$from_model_body = "";
$pkfmt_type = mysqlToFmtType($t->pfield->Type);
foreach($t->fields as $f) {

    $fname = convertFieldName($f->Field);
    $field_name = convertFieldName($f->Field);
    $fname = "FindBy{$fname}";
    $arg = "_findBy" . $f->model_field_name;
    $argtype = $f->go_type;
    $fmt_type = mysqlToFmtType($f->Type);
    if (isPrimaryKey($f) && $t->model_name != "TermRelationship") {

        $fname = "Find";
        $rtype = "bool"; //i.e. we set the current model
        
    } else {
        $rtype = "[]*{$t->model_name}";
    }
    if (isPrimaryKey($f) && $t->model_name == "TermRelationship") {
        $fname = "Find";
        $rtype = "bool"; //i.e. we set the current model
        $from_map_body .= "\t_" . $f->model_field_name . ",err := m[\"{$f->Field}\"].As" . ucfirst($f->go_type). "()\n";
        $from_map_body .= "\tif err != nil {\n \t\treturn o._adapter.Oops(fmt.Sprintf(`%s`,err))\n\t}\n";
        $from_map_body .= "\to." . $f->model_field_name . " = _" . $f->model_field_name . "\n";
        $from_model_body .= "\to.{$f->model_field_name} = m.{$f->model_field_name}\n";
        include "term_relationship_finder.php";
        continue;
    }
    $scol = $f->Field;
    // these are here just to save having to loop
    if ( $f->go_type == "*DateTime" ) {
        $from_map_body .= "\t_" . $f->model_field_name . ",err := m[\"{$f->Field}\"].As" . ucfirst(substr($f->go_type,1)). "()\n";
    } else {
        $from_map_body .= "\t_" . $f->model_field_name . ",err := m[\"{$f->Field}\"].As" . ucfirst($f->go_type). "()\n";
    }
    $from_model_body .= "\to.{$f->model_field_name} = m.{$f->model_field_name}\n";
    
    $from_map_body .= "\tif err != nil {\n \t\treturn o._adapter.Oops(fmt.Sprintf(`%s`,err))\n\t}\n";
    $from_map_body .= "\to." . $f->model_field_name . " = _" . $f->model_field_name . "\n";
    
    if ( $fname == "Find" ) {
        $failure_return = "return false,o._adapter.Oops(fmt.Sprintf(`%s`,err))";
    } else {
        $failure_return = "return _modelSlice,err";
    }
    $sig = "// {$fname} searchs against the database table field {$f->Field} and will return $rtype,error
// This method is a programatically generated finder for {$t->model_name}
";
    if ($fname == "Find") {
        $sig .= "//  
// Note that Find returns a bool of true|false if found or not, not err, in the case of
// found == true, the instance data will be filled out!
//
// A call to find ALWAYS overwrites the model you call Find on
// i.e. receiver is a pointer!
//
//```go
//      m := New{$t->model_name}(a)
//      found,err := m.Find(23)
//      .. handle err
//      if found == false {
//          // handle found
//      }
//      ... do what you want with m here
//```
//
";
    } else {
        $sig .= "//
//```go  
//    m := New{$t->model_name}(a)
//    results,err := m.{$fname}(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of {$t->model_name}
//    }
//```  
//
";
    }
    $sig .= "func (o *{$t->model_name}) $fname($arg $argtype) ($rtype,error) {";
$body = "
    var _modelSlice []*{$t->model_name}
    q := fmt.Sprintf(\"SELECT * FROM %s WHERE `%s` = '$fmt_type'\",o._table, \"{$f->Field}\", $arg)
    results, err := o._adapter.Query(q)
    if err != nil {
        $failure_return
    }
    
    for _,result := range results {
        ro := New{$t->model_name}(o._adapter)
        err = ro.FromDBValueMap(result)
        if err != nil {
            $failure_return
        }
        _modelSlice = append(_modelSlice,ro)
    }
";
if ( $fname == "Find" ) {
    // we return the 0th element
    $body .= "
    if len(_modelSlice) == 0 {
        // there was an error!
        return false, o._adapter.Oops(`not found`)
    }
    o.From{$t->model_name}(_modelSlice[0])
    return true,nil
";
} else {
    $body .= "
    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil
";
}
    puts($sig);
    puts($body);
    puts("}");
}
$find_line = "o.Find(o.GetPrimaryKeyValue())";
$destroy_sql = "fmt.Sprintf(\"DELETE FROM %s WHERE %s = '$pkfmt_type' LIMIT 1\",o._table,o.GetPrimaryKeyName(),o.GetPrimaryKeyValue())";
if ($t->model_name == "TermRelationship") {
    $find_line = "o.Find(o.TermTaxonomyId ,o.ObjectId)";
    $destroy_sql = "fmt.Sprintf(\"DELETE FROM %s WHERE object_id = '%d' AND term_taxonomy_id = %d LIMIT 1\",o._table,o.ObjectId,o.TermTaxonomyId)";
}
puts("
// FromDBValueMap Converts a DBValueMap returned from Adapter.Query to a {$t->model_name}
func (o *{$t->model_name}) FromDBValueMap(m map[string]DBValue) error {
$from_map_body
 \treturn nil
}
// From{$t->model_name} A kind of Clone function for {$t->model_name}
func (o *{$t->model_name}) From{$t->model_name}(m *{$t->model_name}) {
$from_model_body
}
// Reload A function to forcibly reload {$t->model_name}
func (o *{$t->model_name}) Reload() error {
    _,err := $find_line
    return err
}
// Destroy deletes the model
func (o *{$t->model_name}) Destroy() error {
    frmt := $destroy_sql
    err := o._adapter.Execute(frmt)
    if err != nil {
        return o._adapter.Oops(fmt.Sprintf(`%s led to %s`,frmt,err))
    }
    return nil
}
// FindBySQL allows you to search using a complete SQL string
func (o *{$t->model_name}) FindBySQL(s string) ([]*{$t->model_name},error) {
    var _modelSlice []*{$t->model_name}
    
    results, err := o._adapter.Query(s)
    if err != nil {
        $failure_return
    }
    
    for _,result := range results {
        ro := New{$t->model_name}(o._adapter)
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,o._adapter.Oops(fmt.Sprintf(`%s`,err))
        }
        _modelSlice = append(_modelSlice,ro)
    }
    return _modelSlice,nil
}
// Where is a shortcut to FindBySql, in this case you only 
// specify the WHERE clause, such as m.Where(`ID IN (23,25)`)
func (o *{$t->model_name}) Where(s string) ([]*{$t->model_name},error) {
    var _modelSlice []*{$t->model_name}
    
    results, err := o._adapter.Query(fmt.Sprintf(`SELECT * FROM %s WHERE %s`,o._table,s))
    if err != nil {
        $failure_return
    }
    
    for _,result := range results {
        ro := New{$t->model_name}(o._adapter)
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,o._adapter.Oops(fmt.Sprintf(`%s`,err))
        }
        _modelSlice = append(_modelSlice,ro)
    }
    return _modelSlice,nil
}
");
