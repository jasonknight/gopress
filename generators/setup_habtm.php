<?php
// In this context, we are looping through the 
// already processed $t->fields , they are save as $field

// We need to decide if this "Field" specifies a foreign
// key. It already matched the preg_match, but we could
// have a false positive
$fail = false;
$owner = new stdClass();
$child = new stdClass();
preg_match("/(([\w_]+)_[IiDd]{2})/",$field->Field,$m);
// Here are the two cases we're looking for.
// $field is NOT a pfield, yet has comment_id, therefore it
// is a foreign key on this table.
// Array
// (
//     [0] => comment_id
//     [1] => comment_id
//     [2] => comment
// )

// Here is the second case, the Wordpress schema is badly
// inconsistent, so here we have a foreign key where the
// owner is Post sandwiched between comment_ and _ID
// Array
// (
//     [0] => comment_post_ID
//     [1] => comment_post_ID
//     [2] => comment_post
// )
if ( empty($m) ) {
    $fail = true;
}

// Our next step is to figure out who the owner is:

if ( $fail == false ) {
    preg_match("/([\w_]+)_id$/",$field->Field,$m2);
    if ( !empty($m2) ) {
        // we can be reasonably sure this is a later key
        // and so we don't need any magic
        $owner->model = convertTableName($m2[1]);
    } else {
        // let's try the other crazy way
        preg_match("/([\w]+)_([\w_]+)_ID$/",$field->Field,$m2);
        if ( !empty($m2) ) {
            $owner->model = convertTableName($m2[2]);
        }
    }
}
if ($fail == true ) {
    goto done;
}
$child->model = convertTableName($t->model_name);
$child->member_name = convertFieldName($field->Field);
$owner_table = getTableDef($owner->model);  
// Now we know which models are in play, i.e. which Go Structs
$child->parent_member_field_name = "{$owner->model}";
$child->parent_member_field_type = "*{$owner->model}";

$owner->children_member_field_name = pluralize("{$child->model}");
$owner->children_member_field_type = "{$child->model}";
$owner->children_find_by = "FindBy" . convertFieldName($field->Field);
$child->is_parent_loaded_name = "Is{$owner->model}Loaded";
$owner->are_children_loaded_name = "Are{$child->model}sLoaded";



// Now we actually need to fetch the child's parent

$child->load_function = "
func (o *{$child->model}) Load{$owner->model}() ({$child->parent_member_field_type}, error) {
    if o.{$child->is_parent_loaded_name} == true {
        return o.{$child->parent_member_field_name},nil 
    }

    p := New{$owner->model}(o._adapter)
    b,err := p.Find(o.{$child->member_name})
    if err != nil {
        return nil,o._adapter.Oops(errors.New(fmt.Sprintf(`failed {$child->model}.Load{$owner->model}(%d) because %s`,o.{$child->member_name},err)))
    }
    if b != true {
        return nil, o._adapter.Oops(errors.New(fmt.Sprintf(`failed {$child->model}.Load{$owner->model}(%d) because not found`,o.{$child->member_name})))
    }
    o.{$child->parent_member_field_name} = p
    o.{$child->is_parent_loaded_name} = true
    return p,nil
}
func (o *{$child->model}) Reload{$owner->model}() ({$child->parent_member_field_type}, error) {
    o.{$child->is_parent_loaded_name} = false
    return o.Load{$owner->model}()
}
";

// Now we need to fetch the children

$owner->load_function = "
func (o *{$owner->model}) Load{$owner->children_member_field_name}() ([]*{$child->model},error) {
    if o.{$owner->are_children_loaded_name} == true {
        return o.{$owner->children_member_field_name},nil
    }

    c := New{$child->model}(o._adapter)
    r,err := c.{$owner->children_find_by}(o.{$owner_table->pfield->model_field_name})
    if err != nil {
         return nil,o._adapter.Oops(errors.New(fmt.Sprintf(`failed {$owner->model}.{$owner->children_find_by}(%d) because %s`,o.{$owner_table->pfield->model_field_name},err)))
    }
    o.{$owner->children_member_field_name} = r
    o.{$owner->are_children_loaded_name} = true;
    return r,nil
}
func (o *{$owner->model}) Reload{$owner->children_member_field_name}() ([]*{$child->model},error) {
    o.{$owner->are_children_loaded_name} = false
    return o.Load{$owner->children_member_field_name}()
}
";

echo "Currnet table is: {$t->dname} field is {$field->Field}\n";
echo "Owner: \n";
print_r($owner);
echo "Child: \n";
print_r($child);
$i = getTableDef($owner->model,true);
if ($i != -1) {
    $tables[$i]->has_many[] = $owner;
}
$i = getTableDef($child->model,true);
if ($i != -1) {
    $tables[$i]->belongs_to[] = $child;
}
done:
echo "setup_habtm.php done\n";



