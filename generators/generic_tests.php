<?php
$fail = "t.Errorf";
$txt = "
func TestMysqlAdapterFromYAML(t *testing.T) {
    a := NewMysqlAdapter(`pw_`)
    y,err := fileGetContents(`test_data/adapter.yml`)
    if err != nil {
        $fail(`failed to load yaml %s`,err)
        return
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
    array(67859.58686,'float64'),
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
    dvar := a.NewDBValue()
    dvar.SetInternalValue(`x`,`2016-01-09 23:24:50.7Z`)
    dc,err := dvar.AsDateTime()
    if err != nil {
        $fail(`failed to convert datetime %+v`,dc)
    }

    if (dc.Year != 2016 || 
        dc.Month != 1 ||
        dc.Day != 9 ||
        dc.Hours != 23 ||
        dc.Minutes != 24 ||
        dc.Seconds != 50 ||
        dc.Offset != 7 ||
        dc.Zone != `Z`) {
        $fail(`fields don't match up for %+v`,dc)
    }
    r,_ := dvar.AsString()
    if dc.ToString() != r {
        $fail(`restring of dvar failed %s`,dc.ToString())
    }

}
";
puts($txt);