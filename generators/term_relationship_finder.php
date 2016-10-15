<?php
$arg = "termId";
$arg2 = 'objectId';
$fname = "Find";
$sig = "func (o *{$t->model_name}) $fname($arg $argtype,$arg2 $argtype) ($rtype,error) {";
if ( $fname == "Find" ) {
    $failure_return = "return false,err";
} else {
    $failure_return = "return model_slice,err";
}

$body = "
    var model_slice []*{$t->model_name}
    q := fmt.Sprintf(\"SELECT * FROM %s WHERE `term_taxonomy_id` = '%d' AND `object_id` = '%d'\",o._table, $arg,$arg2)
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