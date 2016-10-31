<?php
$txt = "
func (o *PostMeta) FindByKeyValue(k string, v string) ([]*PostMeta,error) {
    var _modelSlice []*PostMeta
    q := fmt.Sprintf(\"SELECT * FROM %s WHERE `meta_key` = '%s' AND meta_value = '%s'\", o._table, o._adapter.SafeString(k),o._adapter.SafeString(v))
    results, err := o._adapter.Query(q)
    if err != nil {
        return nil, o._adapter.Oops(fmt.Sprintf(`%s`, err))
    }

    for _, result := range results {
        ro := NewPostMeta(o._adapter)
        err = ro.FromDBValueMap(result)
        if err != nil {
            return nil, o._adapter.Oops(fmt.Sprintf(`%s`, err))
        }
        _modelSlice = append(_modelSlice, ro)
    }
    return _modelSlice,nil
}
func (o *PostMeta) FindByKeyValueWithPostId(k string, v string, pid int64) ([]*PostMeta,error) {
    var _modelSlice []*PostMeta
    q := fmt.Sprintf(\"SELECT * FROM %s WHERE `post_id` = '%d' AND `meta_key` = '%s' AND meta_value = '%s'\", o._table, pid,o._adapter.SafeString(k),o._adapter.SafeString(v))
    results, err := o._adapter.Query(q)
    if err != nil {
        return nil, o._adapter.Oops(fmt.Sprintf(`%s`, err))
    }

    for _, result := range results {
        ro := NewPostMeta(o._adapter)
        err = ro.FromDBValueMap(result)
        if err != nil {
            return nil, o._adapter.Oops(fmt.Sprintf(`%s`, err))
        }
        _modelSlice = append(_modelSlice, ro)
    }
    return _modelSlice,nil
}
func (o *Post) FindByPostMetaKeyValue(k string, v string) ([]*Post,error) {
    var _modelSlice []*Post
    m := NewPostMeta(o._adapter)
    metas,err := m.FindByKeyValueWithPostId(k,v,o.ID)
    if err != nil {
        return nil, o._adapter.Oops(fmt.Sprintf(`%s`, err))
    }
    if len(metas) == 0 {
        return _modelSlice,nil
    }
    for _,meta := range metas {
        p := NewPost(o._adapter)
        found,err := p.Find(meta.PostId)
        if err != nil {
            return nil, o._adapter.Oops(fmt.Sprintf(`%s`, err))
        }
        if found == true {
            _modelSlice = append(_modelSlice,p)
        }
    }
    return _modelSlice,nil
}
";
puts($txt);