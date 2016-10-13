<?php
$p = $t->pfield;
$pname = convertFieldName($p->Field);
$ptype = $p->go_type;
$txt = "
func (m *{$t->model_name}) GetPrimaryKeyValue() $ptype {
    return m.$pname
}
func (m *{$t->model_name}) GetPrimaryKeyName() string {
    return `{$p->Field}`
}
";
foreach ($t->fields as $f) {
$txt .= "
func (m *{$t->model_name}) Get{$f->model_field_name} () {$f->go_type} {
    return m.{$f->model_field_name}
}
func (m *{$t->model_name}) Set{$f->model_field_name} (arg {$f->go_type}) {
    m.{$f->model_field_name} = arg
    m.{$f->dirty_marker} = true
}
";
}

puts($txt);