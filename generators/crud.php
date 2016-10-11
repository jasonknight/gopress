<?php
$pkeyname = convertFieldName($t->pfield->Field);
$pkey_fmt_type = $fmt_type = mysqlToFmtType($t->pfield->Type);
foreach($t->fields as $f) {
    $fname = convertFieldName($f->Field);
    $field_name = convertFieldName($f->Field);
    $fname = "Update{$fname}";
    $arg = "_update_" . lcfirst(convertFieldName($f->Field));
    $argtype = $f->go_type;
    $fmt_type = mysqlToFmtType($f->Type);
    if ($f->Key == "PRI") {
        continue;
    }
    $scol = $f->Field;
$sig = "func (o *{$t->model_name}) $fname($arg $argtype) (error) {";
$body = "
    q := fmt.Sprintf(\"UPDATE %s WHERE SET %s = '$fmt_type' WHERE `%s` = '$pkey_fmt_type'\",o._table_, \"{$f->Field}\", $arg, o._pkey_,o.${pkeyname})
    err := o._adapter_.Execute(q)
    if err != nil {
        return err
    }
    return nil
";
puts($sig);
    puts($body);
    puts("}");
}