<?php
$defs = "
type ArgMap map[string]string
";
$wp_get_object_terms = "
func (a *MysqlAdapter) WpGetObjectTerms(object_ids []int64, taxonomy_array []string, args ArgMap) ([]*Term,err) {
    if len(tids) == 0 || len(taxonomies) == 0 {
        return errors.New(`both tids and taxonomies must be non-zero`)
    }
    for _,t := range {
        if TaxonomyExists(t) != true {
            return errors.New(fmt.Sprintf(`%s does not exist`,t))
        }
    }
    defaults := make(ArgMap)
    var terms []*TermTaxonomy

    orderby := args[\"orderby\"]
    order := args[\"order\"]

    if InStringSlice(orderby,[]string{`term_id`,`slug`,`term_group`}) {
        orderby = fmt.Sprintf(`t.%s`,orderby)
    }
    if InStringSlice(orderby,[]string{`count`,`parent`,`taxonomy`,`term_taxonomy_id`}) {
        orderby = fmt.Sprintf(`tt.%s`,orderby)
    }
    if \"term_order\" == orderby {
        orderby = fmt.Sprintf(`tr.%s`,orderby)
    }

    if \"tt_ids\" == fields && orderby != \"\" {
        orderby = `tr.term_taxonomy_id`
    }

    if orderby != `` {
        orderby = fmt.Sprintf(`ORDER BY %s`,orderby)
    }

    if order != `` && InStringSlice(order,[]string{`ASC`,`DSC`}) {
        order = `ASC`
    }

    taxonomies := strings.Join(`,`,taxonomy_array)
    oids := strings.join(`,`,object_ids)

    var where []string
    a_term := NewTerm(a)
    a_term_tax := NewTermTaxonomy(a)
    a_term_rel := NewTermRelationship(a)
    where = append(where,fmt.Sprintf(`tt.taxonomy IN (%s)`,taxonomies))
    where = append(where,fmt.Sprintf(`tr.object_id IN (%s)`,oids))

    where_str := strings.Join(` AND `, where)
    query := fmt.Sprintf(`SELECT 
        t.term_id 
    FROM 
        %s 
    AS 
        t 
    INNER JOIN 
        %s 
    AS 
        tt 
    ON 
        tt.term_id = t.term_id 
    INNER JOIN 
        %s 
    AS 
        tr 
    ON 
        tr.term_taxonomy_id = tt.term_taxonomy_id 
    WHERE 
        %s %s %s`, 
    a_term._table, 
    a_term_tax._table, 
    a_term_rel._table,
    where_str, 
    orderby, 
    order)

    results,err := a.Query(query)
    if err != nil {
        return nil,err
    }
}
";
$tax_exists = "
func (a *MysqlAdapter) WpTaxonomyExists(t string) bool {
    tx := NewTermTaxonomy(a)
    found, err := tx.FindByTaxonomyOnly(t)
    if err != nil {
        a.Oops(err)
        return false
    }
    return found
}
";