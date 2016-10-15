<?php
if ( !function_exists("gen_updaters")) {
function gen_updaters($t) {
    $txt = "";
    foreach ($t->fields as $f) {
        if ( isPrimaryKey($f) ) {
            continue;
        }
        $fname = convertFieldName($f->Field);
        $mname = maybeLC($fname);
        $arg = "_upd_" . maybeLC(convertFieldName($f->Field));
        $argtype = $f->go_type;
        $fmt_type = mysqlToFmtType($f->Type);
        $pkmname = maybeLC(convertFieldName($f->Field));
        $pkfname = $t->pfield->Field;
        $pkfmttype = mysqlToFmtType($t->pfield->go_type);
        $update_line = "\"UPDATE %s SET `{$f->Field}` = '$fmt_type' WHERE `$pkfname` = '$pkfmttype'\",o._table,$arg,o.{$pkmname}";
        if ( $t->model_name == "TermRelationship") {
           $update_line = "\"UPDATE %s SET `{$f->Field}` = '$fmt_type' WHERE term_taxonomy_id = '%d' AND object_id = '%d'\",o._table,$arg,o.TermTaxonomyId,o.ObjectId"; 
        }
$txt .= "
func (o *{$t->model_name}) Update{$fname}($arg $argtype) (int64,error) {
    frmt := fmt.Sprintf($update_line)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.{$mname} = $arg
    return o._adapter.AffectedRows(),nil
}
";
    } 
    return $txt;
}    
}
puts(gen_updaters($t));

