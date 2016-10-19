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
func TestAdapterFailures(t *testing.T) {
    _,err := NewMysqlAdapterEx(`file_that_does_not_exist123323`)
    if err == nil {
        $fail(`Did not receive an error when file should not exist!`)
        return
    }
    // Load a nonsense yaml file
    _,err = NewMysqlAdapterEx(`test_data/nonsenseyaml.yml`)
    if err == nil {
        $fail(`this should fail to load a nonsense yaml file`)
        return
    }
    // Load a test yaml with wrong Open
    _, err = NewMysqlAdapterEx(`test_data/adapter.yml`)
    if err == nil {
        $fail(`this should fail with wrong login info`)
        return
    }
    // Load a silly yml file with wrong data
    _, err = NewMysqlAdapterEx(`test_data/silly.yml`)
    if err == nil {
        $fail(`this should fail with wrong login info`)
        return
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
    dvar.SetInternalValue(`x`,`2016-01-09 23:24:50`)
    dc,err := dvar.AsDateTime()
    if err != nil {
        $fail(`failed to convert datetime %+v`,dc)
    }

    if (dc.Year != 2016 || 
        dc.Month != 1 ||
        dc.Day != 9 ||
        dc.Hours != 23 ||
        dc.Minutes != 24 ||
        dc.Seconds != 50 ) {
        $fail(`fields don't match up for %+v`,dc)
    }
    r,_ := dvar.AsString()
    if dc.ToString() != r {
        $fail(`restring of dvar failed %s`,dc.ToString())
    }

}
";
foreach (array("Info","Debug") as $ltype) {
    $ltag = strtoupper($ltype);
$txt .= "
func TestAdapter{$ltype}Logging(t *testing.T) {
    a := NewMysqlAdapter(`wp_`)
    var b bytes.Buffer
    r, err := regexp.Compile(`\\[{$ltag}\\]:.+Hello World`)
    if err != nil {
        $fail(`could not compile regex`)
        return
    }
    wr := bufio.NewWriter(&b)
    a.SetLogs(wr)
    a.Log{$ltype}(`Hello World`)
    wr.Flush()
    if r.MatchString(b.String()) == false {
        $fail(`failed to match info line`)
        return
    }
}
func TestAdapterEmpty{$ltype}Logging(t *testing.T) {
    a := NewMysqlAdapter(`wp_`)
    var b bytes.Buffer
    wr := bufio.NewWriter(&b)
    a.SetLogs(wr)
    a.Log{$ltype}(``)
    wr.Flush()
    if b.String() != `` {
        $fail(`Info should not occur in this case`)
        return
    }
    a.SetLogFilter(func (tag string,val string) string {
        return ``
    })
    a.Log{$ltype}(`Hello World`)
    wr.Flush()
    if b.String() != `` {
        $fail(`Info should not occur due to filter in this case`)
        return
    }
}
";
}
$txt .= "
func TestAdapterErrorLogging(t *testing.T) {
    a := NewMysqlAdapter(`wp_`)
    
    r, err := regexp.Compile(`\\[ERROR\\]:.+Hello World`)
    if err != nil {
        $fail(`could not compile regex`)
        return
    }
    var b bytes.Buffer
    wr := bufio.NewWriter(&b)
    a.SetLogs(wr)
    a.LogError(errors.New(`Hello World`))
    wr.Flush()
    if r.MatchString(b.String()) == false {
        $fail(`failed to match info line`)
        return
    }

    var b2 bytes.Buffer
    wr2 := bufio.NewWriter(&b2)
    a.SetLogs(wr2)
    a.SetLogFilter(func (tag string,val string) string {
        return ``
    })
    a.LogError(errors.New(`Hello World`))
    wr2.Flush()
    if b2.String() != `` {
        $fail(`Info should not occur due to filter in this case but equals %s`,b2.String())
        return
    }
}

";
puts($txt);