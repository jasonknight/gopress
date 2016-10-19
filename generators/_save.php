<?php
if ( ! function_exists('_save_create') ) {
function _save_create($t) {
    $go_fnames = array();
    $mysql_fnames = array();
    $fmts = array();
    $i = 0;
    $pkeyfmt = "";
    $pkeyname = "";
    $sets = "";
    $pkeyfmt = mysqlToFmtType($t->pfield->Type);
    $pkeyname =  maybeLC(convertFieldName($t->pfield->Field));
    foreach ( $t->fields as $tf) {
        if (isPrimaryKey($tf) && $t->model_name != "TermRelationship") {
            continue;
        }
        if ($tf->go_type == "*DateTime") {
             $go_fnames[$i] = maybeLC(convertFieldName($tf->Field)) . ".ToString()";
        } else {
             $go_fnames[$i] = maybeLC(convertFieldName($tf->Field));
        }
        $gfn = convertFieldName($tf->Field);
        $mysql_fnames[$i] = $tf->Field;
        $fmts[$i] = mysqlToFmtType($tf->Type);
        if ($tf->go_type == 'string') {
            $the_gfn = "o._adapter.SafeString(o.$gfn)";
        } else {
            $the_gfn = "o.$gfn";
        }
        $sets .= "
    if o.Is{$gfn}Dirty == true {
        sets = append(sets,fmt.Sprintf(`{$mysql_fnames[$i]} = '{$fmts[$i]}'`,$the_gfn))
    }
";
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
    if ( $t->model_name == "TermRelationship") {
        $where = "`term_taxonomy_id` = '%d' AND object_id = '%d'";
        $up_gn_line = "o.TermTaxonomyId, o.ObjectId";
    } else {
        $where = "%s = '$pkeyfmt'";
        $up_gn_line = "o._pkey, o.$pkeyname";  
    }
    $set_primary_key_field = "o.{$t->pfield->model_field_name} = o._adapter.LastInsertedId()";
    if ($t->model_name == "TermRelationship") {
        $set_primary_key_field = "";
    }
$txt = "// Save is a dynamic saver 'inherited' by all models
func (o *{$t->model_name}) Save() error {
    if o._new == true {
        return o.Create()
    }
    var sets []string
    $sets
    frmt := fmt.Sprintf(\"UPDATE %s SET %s WHERE $where\",o._table,strings.Join(sets,`,`),$up_gn_line)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return err
    }
    return nil
}
// Update is a dynamic updater, it considers whether or not
// a field is 'dirty' and needs to be updated. Will only work
// if you use the Getters and Setters
func (o *{$t->model_name}) Update() error {
    var sets []string
    $sets
    frmt := fmt.Sprintf(\"UPDATE %s SET %s WHERE $where\",o._table,strings.Join(sets,`,`),$up_gn_line)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return err
    }
    return nil
}
// Create inserts the model. Calling Save will call this function
// automatically for new models
func (o *{$t->model_name}) Create() error {
    frmt := fmt.Sprintf(\"INSERT INTO %s ($cr_col_line) VALUES ($cr_val_line)\",o._table,$cr_gn_line)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return o._adapter.Oops(fmt.Sprintf(`%s led to %s`,frmt,err))
    }
    $set_primary_key_field
    o._new = false
    return nil
}
";   
    return $txt; 
}
}
puts(_save_create($t));
