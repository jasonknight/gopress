<?php
$fail = "t.Errorf";
$txt = "
func TestNew{$t->model_name}(t *testing.T) {
    a := NewMysqlAdapter(\"wp_\")
    o := New{$t->model_name}(a)
    if o._table != \"{$t->database_name}\" {
        $fail(\"failed creating %+v\",o);
        return
    }
}
func Test{$t->model_name}FromDBValueMap(t *testing.T) {
    a := NewMysqlAdapter(\"wp_\")
    o := New{$t->model_name}(a)
    m := make(map[string]DBValue)
";
// $fields = array();
// $ftypes = array();
// foreach ($t->fields as $f) {
//     $fields[] = $f->Field;
//     $ftypes[] = $f->go_type;
// }
//$quoted_fields = array_map(function ($x) {return "\"$x\"";},$fields);
foreach ($t->fields as $f) {
    $k = "\"{$f->Field}\"";
    if ( $f->go_type == 'string' ) {
        $v = "\"AString\"";
    } else if ( preg_match("/^int/",$f->go_type) ) {
        $v = "strconv.Itoa(999)";
    } else if ( $f->go_type == "DateTime" ) {
        $v = "\"2016-01-01 10:50:23.5Z\"";
    } else {
        die("What? no go_type support for {$f->go_type}\n");
    }
    $txt .= "\tm[$k] = a.NewDBValue()\n";
    $txt .= "\tm[$k].SetInternalValue($k,$v)\n";
}
$txt .= "
    err := o.FromDBValueMap(m)
    if err != nil {
        $fail(\"FromDBValueMap failed %s\",err)
    }
";
foreach ($t->fields as $f) {
    $k = "\"{$f->Field}\"";
    $fmname = lcfirst(convertFieldName($f->Field));
    if ( $f->go_type == 'string' ) {
        $v = "\"AString\"";
    } else if ( preg_match("/^int/",$f->go_type) ) {
        $v = 999;
    } else if ( $f->go_type == "DateTime" ) {
        $v = "\"2016-01-01 10:50:23.5Z\"";
    } else {
        die("What? no go_type support for {$f->go_type}\n");
    }
    if ($f->go_type != "DateTime" ) {
$txt .= "
    if o.{$fmname} != $v {
        $fail(\"o.{$fmname} test failed %+v\",o)
        return
    }    
";
    } else {
        // This is a DateTime, so let's
$txt .= "
    if o.{$fmname}.Year != 2016 {
        $fail(\"year not set for %+v\",o.{$fmname})
        return
    }
";
    }


}
$txt .= "}";
puts($txt);
