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
    } else if ( $f->go_type == "*DateTime" ) {
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
$j = 0;
foreach ($t->fields as $f) {
    $k = "\"{$f->Field}\"";
    $fmname = maybeLC(convertFieldName($f->Field));
    if ( $f->go_type == 'string' ) {
        $v = "\"AString\"";
    } else if ( preg_match("/^int/",$f->go_type) ) {
        $v = 999;
    } else if ( $f->go_type == "*DateTime" ) {
        $v = "\"2016-01-01 10:50:23.5Z\"";
    } else {
        die("What? no go_type support for {$f->go_type}\n");
    }
    if ($f->go_type != "*DateTime" ) {
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
    if (o.{$fmname}.Year != 2016 || 
        o.{$fmname}.Month != 1 ||
        o.{$fmname}.Day != 1 ||
        o.{$fmname}.Hours != 10 ||
        o.{$fmname}.Minutes != 50 ||
        o.{$fmname}.Seconds != 23 ||
        o.{$fmname}.Offset != 5 ||
        o.{$fmname}.Zone != `Z`) {
        $fail(`fields don't match up for %+v`,o.{$fmname})
    }
    r{$j},_ := m[$k].AsString()
    if o.{$fmname}.ToString() != r{$j} {
        $fail(`restring of o.{$fmname} failed %s`,o.{$fmname}.ToString())
    }
";
    }
$j++;

}
$txt .= "}";
puts($txt);
