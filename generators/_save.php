<?php
if ( ! function_exists('_save_create') ) {
function _save_create($t) {
    $go_fnames = array();
    $mysql_fnames = array();
    $fmts = array();
    $i = 0;
    $pkeyfmt = "";
    $pkeyname = "";
    foreach ( $t->fields as $tf) {
        if (isPrimaryKey($tf)) {
            $pkeyfmt = mysqlToFmtType($tf->Type);
            $pkeyname =  maybeLC(convertFieldName($tf->Field));
            continue;
        }
        if ($tf->go_type == "*DateTime") {
             $go_fnames[$i] = maybeLC(convertFieldName($tf->Field)) . ".ToString()";
        } else {
             $go_fnames[$i] = maybeLC(convertFieldName($tf->Field));
        }
       
        $mysql_fnames[$i] = $tf->Field;
        $fmts[$i] = mysqlToFmtType($tf->Type);
        $i++;
    }
    $fmts = array_map(function ($x) { return "'$x'";},$fmts);
    $mysql_fnames = array_map(function ($x) { return "`$x`";},$mysql_fnames);
    $go_fnames = array_map(function ($x) { return "o.$x";},$go_fnames);

    $update_entries = array();
    $cr_cols = array();
    $cr_vals = array();
    for ($i = 0; $i < count($go_fnames); $i++) {
        $mf = $mysql_fnames[$i];
        $fm = $fmts[$i];
        $update_entries[] = "$mf = $fm";
        
        $cr_cols[] = $mf;
        $cr_vals[] = $fm;
    }
    $up_fmt_line = join(", ",$update_entries);
    $up_gn_line = join(", ",$go_fnames);

    $cr_gn_line = $up_gn_line;
    $cr_col_line = join(", ",$cr_cols);
    $cr_val_line = join(", ",$cr_vals);
$txt = "
func (o *{$t->model_name}) Save() (int64,error) {
    if o._new == true {
        return o.Create()
    }
    frmt := fmt.Sprintf(\"UPDATE %s SET $up_fmt_line WHERE %s = '$pkeyfmt' LIMIT 1\",o._table,$up_gn_line,o._pkey, o.$pkeyname)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }

    return o._adapter.AffectedRows(),nil
}
func (o *{$t->model_name}) Create() (int64,error) {
    frmt := fmt.Sprintf(\"INSERT INTO %s ($cr_col_line) VALUES ($cr_val_line)\",o._table,$cr_gn_line)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,o._adapter.Oops(fmt.Sprintf(`%s led to %s`,frmt,err))
    }
    o.{$t->pfield->model_field_name} = o._adapter.LastInsertedId()

    return o._adapter.AffectedRows(),nil
}
";   
    return $txt; 
}
}
puts(_save_create($t));
