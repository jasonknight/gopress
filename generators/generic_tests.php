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
";
puts($txt);