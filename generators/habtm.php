<?php
if (!function_exists("habtm")) {
function habtm($t) {
    global $tables;
    $txt = "";

    foreach ($t->belongs_to as $bt) {
        $owner = "";
        foreach ($tables as $tt) {
            if ( $bt->model == $tt->model_name) {
                $owner = $tt;
                break;
            }
        }
        if ( $owner != "" ) {
$txt .= "
func (o *{$bt->model_name}) Load{$bt->model}() ({$bt->go_type},error) {
    if o.Is{$bt->model}Loaded == true {
        return o.{$bt->model},nil
    }
    m := New{$bt->model}(o._adapter)
    found,err := m.Find(o.Get{$bt->model_field_name}())
    if err != nil {
        return nil,err
    }
    if found == false {
        return nil,errors.New(fmt.Sprintf(`could not find {$bt->model} with {$owner->pfield->Field} of {$owner->pfield->mysql_fmt_type}`,o.Get{$bt->model_field_name}()))
    }
    o.Is{$bt->model}Loaded = true
    o.{$bt->model} = m
    return m,nil
}

";
        }
    }
    // stdClass Object
    // (
    //     [name] => CommentMetas
    //     [table] => commentmeta
    //     [type] => []CommentMeta
    //     [fkey] => comment_id
    //     [fkey_type] => int64
    //     [fkey_myfmt] => %d
    // )

    foreach ($t->has_many as $hm) {
$txt .= "
func (o *{$t->model_name}) Load{$hm->name}() ({$hm->type},error) {
    if o.Is{$hm->name}Loaded == true {
        return o.{$hm->name},nil
    }
    var finder {$hm->model_name}
    results, err := finder.FindBy{$hm->model_field_name}(o.{$t->pfield->model_field_name})
    if err != nil {
        return nil,err
    }
    o.Is{$hm->name}Loaded = true
    o.{$hm->name} = results
    return results,nil
}
";
    }
    if (!empty($t->has_many) || !empty($t->belongs_to) ) {
        $txt .= "
func (o *{$t->model_name}) SaveRelations() (error) {
";
    $i = 0;
        foreach ($t->has_many as $hm) {
$txt .= "
        for _,m := range o.{$hm->name} {
            _,err :=  m.Save()
            if err != nil {
                return errors.New(fmt.Sprintf(`while saving relation {$hm->name} %s`,err))
            }
        }
";
            $i++;
        }
        foreach ($t->belongs_to as $bt) {
$txt .= "
        _,err :=  o.{$bt->model}.Save()
        if err != nil {
            return errors.New(fmt.Sprintf(`while saving relation {$bt->model} %s`,err))
        }
";
            $i++;
        }
$txt .= "
    return nil
}
";
    }
    return $txt;
}
}
puts(habtm($t));