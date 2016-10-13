<?php
if ( !function_exists("gen_updaters")) {
function gen_updaters($t) {
    $txt = "";
    foreach ($t->fields as $f) {
        if ( $f->Key == 'PRI' ) {
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
$txt .= "
func (o *{$t->model_name}) Update{$fname}($arg $argtype) (int64,error) {
    frmt := fmt.Sprintf(\"UPDATE %s SET `{$f->Field}` = '$fmt_type' WHERE `$pkfname` = '$pkfmttype'\",o._table,$arg,o.{$pkmname})
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

