<?php
$fail = "t.Errorf";
$txt = "
func TestMysqlAdapterFromYAML(t *testing.T) {
    a := NewMysqlAdapter(`pw_`)
    y,err := fileGetContents(`test_data/adapter.yml`)
    if err != nil {
        $fail(`failed to load yaml %s`,err)
    }
    err = a.FromYAML(y)
    if err != nil {
        $fail(`failed to apply yaml %s`,err)
        return
    }

    if (a.User != `root` ||
        a.Pass != `rootpass` ||
        a.Host != `localhost` ||
        a.Database != `my_db` ||
        a.DBPrefix != `wp_`) {
        $fail(`did not fully apply yaml file %+v`,a)
    }
}
func TestDBValue(t *testing.T) {
    a := NewMysqlAdapter(`wp_`)
";
$tbl = array(
    array(999,'int32'),
    array(666,'int'),
    array('hello world','string'),
    array(3.14,'float32'),
    array(67859.58686,'float32'),
);
for ($i = 0; $i < count($tbl); $i++) {
    $var = "v{$i}";
    $c = "c{$i}";
    $fn = "As" . ucfirst($tbl[$i][1]) . "()";
    if ($tbl[$i][1] == 'string') {
        $tval = "\"{$tbl[$i][0]}\"";
    } else {
        $tval = $tbl[$i][0];
    }
$txt .= "
    $var := a.NewDBValue()
    $var.SetInternalValue(`x`,`{$tbl[$i][0]}`)
    $c,err := $var.$fn
    if err != nil {
        $fail(`failed to convert with $fn %+v`,$var)
        return
    }
    if $c != $tval {
        $fail(`values don't match `)
        return
    }
";
}
$txt .= "
}
";
puts($txt);