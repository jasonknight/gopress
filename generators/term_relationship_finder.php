<?php
$arg = "termId";
$arg2 = 'objectId';
$fname = "Find";
$sig = "// Find for the TermRelationship is a bit tricky, as it has no
// primary key as such, but a composite key.
";
$sig .= "func (o *{$t->model_name}) $fname($arg $argtype,$arg2 $argtype) ($rtype,error) {";
if ( $fname == "Find" ) {
    $failure_return = "return false,err";
} else {
    $failure_return = "return _modelSlice,err";
}

$body = "
    var _modelSlice []*{$t->model_name}
    q := fmt.Sprintf(\"SELECT * FROM %s WHERE `term_taxonomy_id` = '%d' AND `object_id` = '%d'\",o._table, $arg,$arg2)
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
        return false, errors.New(\"not found\")
    }
    o.From{$t->model_name}(_modelSlice[0])
    return true,nil
";
} else {
    $body .= "
    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, errors.New(\"no results\")
    }
    return _modelSlice,nil
";
}
    puts($sig);
    puts($body);
    puts("}");