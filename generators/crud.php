<?php
// $pkeyname = lcfirst(convertFieldName($t->pfield->Field));
// $pkey_fmt_type = $fmt_type = mysqlToFmtType($t->pfield->Type);
// foreach($t->fields as $f) {
//     $fname = convertFieldName($f->Field);
//     $fname = "Update{$fname}";
//     $arg = "_update_" . lcfirst(convertFieldName($f->Field));
//     $argtype = $f->go_type;
//     $fmt_type = mysqlToFmtType($f->Type);
//     if ($f->Key == "PRI") {
//         continue;
//     }
//     $scol = $f->Field;
// $sig = "func (o *{$t->model_name}) $fname($arg $argtype) (error) {";
// $body = "
//     q := fmt.Sprintf(\"UPDATE %s WHERE SET %s = '$fmt_type' WHERE `%s` = '$pkey_fmt_type' LIMIT 1\",o._table, \"{$f->Field}\", $arg, o._pkey,o.${pkeyname})
//     err := o._adapter.Execute(q)
//     if err != nil {
//         return err
//     }
//     return nil
// ";
// puts($sig);
//     puts($body);
//     puts("}");

// }
include "_save.php";
include "updaters.php";