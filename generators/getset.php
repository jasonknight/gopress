<?php
$p = $t->pfield;
$pname = convertFieldName($p->Field);
$ptype = $p->go_type;
$txt = "
// GetPrimaryKeyValue returns the value, usually int64 of
// the PrimaryKey
func (o *{$t->model_name}) GetPrimaryKeyValue() $ptype {
    return o.$pname
}
// GetPrimaryKeyName returns the DB field name
func (o *{$t->model_name}) GetPrimaryKeyName() string {
    return `{$p->Field}`
}
";
foreach ($t->fields as $f) {
$txt .= "
// Get{$f->model_field_name} returns the value of 
// {$t->model_name}.{$f->model_field_name}
func (o *{$t->model_name}) Get{$f->model_field_name}() {$f->go_type} {
    return o.{$f->model_field_name}
}
// Set{$f->model_field_name} sets and marks as dirty the value of
// {$t->model_name}.{$f->model_field_name}
func (o *{$t->model_name}) Set{$f->model_field_name}(arg {$f->go_type}) {
    o.{$f->model_field_name} = arg
    o.{$f->dirty_marker} = true
}
";
}

puts($txt);