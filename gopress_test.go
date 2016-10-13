package gopress

import (
	"strconv"
	"testing"
)

func TestNewCommentMeta(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewCommentMeta(a)
	if o._table != "wp_commentmeta" {
		t.Errorf("failed creating %+v", o)
		return
	}
}
func TestCommentMetaFromDBValueMap(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewCommentMeta(a)
	m := make(map[string]DBValue)
	m["meta_id"] = a.NewDBValue()
	m["meta_id"].SetInternalValue("meta_id", strconv.Itoa(999))
	m["comment_id"] = a.NewDBValue()
	m["comment_id"].SetInternalValue("comment_id", strconv.Itoa(999))
	m["meta_key"] = a.NewDBValue()
	m["meta_key"].SetInternalValue("meta_key", "AString")
	m["meta_value"] = a.NewDBValue()
	m["meta_value"].SetInternalValue("meta_value", "AString")

	err := o.FromDBValueMap(m)
	if err != nil {
		t.Errorf("FromDBValueMap failed %s", err)
	}

	if o.MetaId != 999 {
		t.Errorf("o.MetaId test failed %+v", o)
		return
	}

	if o.CommentId != 999 {
		t.Errorf("o.CommentId test failed %+v", o)
		return
	}

	if o.MetaKey != "AString" {
		t.Errorf("o.MetaKey test failed %+v", o)
		return
	}

	if o.MetaValue != "AString" {
		t.Errorf("o.MetaValue test failed %+v", o)
		return
	}
}

func TestNewComment(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewComment(a)
	if o._table != "wp_comments" {
		t.Errorf("failed creating %+v", o)
		return
	}
}
func TestCommentFromDBValueMap(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewComment(a)
	m := make(map[string]DBValue)
	m["comment_ID"] = a.NewDBValue()
	m["comment_ID"].SetInternalValue("comment_ID", strconv.Itoa(999))
	m["comment_post_ID"] = a.NewDBValue()
	m["comment_post_ID"].SetInternalValue("comment_post_ID", strconv.Itoa(999))
	m["comment_author"] = a.NewDBValue()
	m["comment_author"].SetInternalValue("comment_author", "AString")
	m["comment_author_email"] = a.NewDBValue()
	m["comment_author_email"].SetInternalValue("comment_author_email", "AString")
	m["comment_author_url"] = a.NewDBValue()
	m["comment_author_url"].SetInternalValue("comment_author_url", "AString")
	m["comment_author_IP"] = a.NewDBValue()
	m["comment_author_IP"].SetInternalValue("comment_author_IP", "AString")
	m["comment_date"] = a.NewDBValue()
	m["comment_date"].SetInternalValue("comment_date", "2016-01-01 10:50:23.5Z")
	m["comment_date_gmt"] = a.NewDBValue()
	m["comment_date_gmt"].SetInternalValue("comment_date_gmt", "2016-01-01 10:50:23.5Z")
	m["comment_content"] = a.NewDBValue()
	m["comment_content"].SetInternalValue("comment_content", "AString")
	m["comment_karma"] = a.NewDBValue()
	m["comment_karma"].SetInternalValue("comment_karma", strconv.Itoa(999))
	m["comment_approved"] = a.NewDBValue()
	m["comment_approved"].SetInternalValue("comment_approved", "AString")
	m["comment_agent"] = a.NewDBValue()
	m["comment_agent"].SetInternalValue("comment_agent", "AString")
	m["comment_type"] = a.NewDBValue()
	m["comment_type"].SetInternalValue("comment_type", "AString")
	m["comment_parent"] = a.NewDBValue()
	m["comment_parent"].SetInternalValue("comment_parent", strconv.Itoa(999))
	m["user_id"] = a.NewDBValue()
	m["user_id"].SetInternalValue("user_id", strconv.Itoa(999))

	err := o.FromDBValueMap(m)
	if err != nil {
		t.Errorf("FromDBValueMap failed %s", err)
	}

	if o.CommentID != 999 {
		t.Errorf("o.CommentID test failed %+v", o)
		return
	}

	if o.CommentPostID != 999 {
		t.Errorf("o.CommentPostID test failed %+v", o)
		return
	}

	if o.CommentAuthor != "AString" {
		t.Errorf("o.CommentAuthor test failed %+v", o)
		return
	}

	if o.CommentAuthorEmail != "AString" {
		t.Errorf("o.CommentAuthorEmail test failed %+v", o)
		return
	}

	if o.CommentAuthorUrl != "AString" {
		t.Errorf("o.CommentAuthorUrl test failed %+v", o)
		return
	}

	if o.CommentAuthorIP != "AString" {
		t.Errorf("o.CommentAuthorIP test failed %+v", o)
		return
	}

	if o.CommentDate.Year != 2016 {
		t.Errorf("year not set for %+v", o.CommentDate)
		return
	}

	if o.CommentDateGmt.Year != 2016 {
		t.Errorf("year not set for %+v", o.CommentDateGmt)
		return
	}

	if o.CommentContent != "AString" {
		t.Errorf("o.CommentContent test failed %+v", o)
		return
	}

	if o.CommentKarma != 999 {
		t.Errorf("o.CommentKarma test failed %+v", o)
		return
	}

	if o.CommentApproved != "AString" {
		t.Errorf("o.CommentApproved test failed %+v", o)
		return
	}

	if o.CommentAgent != "AString" {
		t.Errorf("o.CommentAgent test failed %+v", o)
		return
	}

	if o.CommentType != "AString" {
		t.Errorf("o.CommentType test failed %+v", o)
		return
	}

	if o.CommentParent != 999 {
		t.Errorf("o.CommentParent test failed %+v", o)
		return
	}

	if o.UserId != 999 {
		t.Errorf("o.UserId test failed %+v", o)
		return
	}
}

func TestNewLink(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewLink(a)
	if o._table != "wp_links" {
		t.Errorf("failed creating %+v", o)
		return
	}
}
func TestLinkFromDBValueMap(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewLink(a)
	m := make(map[string]DBValue)
	m["link_id"] = a.NewDBValue()
	m["link_id"].SetInternalValue("link_id", strconv.Itoa(999))
	m["link_url"] = a.NewDBValue()
	m["link_url"].SetInternalValue("link_url", "AString")
	m["link_name"] = a.NewDBValue()
	m["link_name"].SetInternalValue("link_name", "AString")
	m["link_image"] = a.NewDBValue()
	m["link_image"].SetInternalValue("link_image", "AString")
	m["link_target"] = a.NewDBValue()
	m["link_target"].SetInternalValue("link_target", "AString")
	m["link_description"] = a.NewDBValue()
	m["link_description"].SetInternalValue("link_description", "AString")
	m["link_visible"] = a.NewDBValue()
	m["link_visible"].SetInternalValue("link_visible", "AString")
	m["link_owner"] = a.NewDBValue()
	m["link_owner"].SetInternalValue("link_owner", strconv.Itoa(999))
	m["link_rating"] = a.NewDBValue()
	m["link_rating"].SetInternalValue("link_rating", strconv.Itoa(999))
	m["link_updated"] = a.NewDBValue()
	m["link_updated"].SetInternalValue("link_updated", "2016-01-01 10:50:23.5Z")
	m["link_rel"] = a.NewDBValue()
	m["link_rel"].SetInternalValue("link_rel", "AString")
	m["link_notes"] = a.NewDBValue()
	m["link_notes"].SetInternalValue("link_notes", "AString")
	m["link_rss"] = a.NewDBValue()
	m["link_rss"].SetInternalValue("link_rss", "AString")

	err := o.FromDBValueMap(m)
	if err != nil {
		t.Errorf("FromDBValueMap failed %s", err)
	}

	if o.LinkId != 999 {
		t.Errorf("o.LinkId test failed %+v", o)
		return
	}

	if o.LinkUrl != "AString" {
		t.Errorf("o.LinkUrl test failed %+v", o)
		return
	}

	if o.LinkName != "AString" {
		t.Errorf("o.LinkName test failed %+v", o)
		return
	}

	if o.LinkImage != "AString" {
		t.Errorf("o.LinkImage test failed %+v", o)
		return
	}

	if o.LinkTarget != "AString" {
		t.Errorf("o.LinkTarget test failed %+v", o)
		return
	}

	if o.LinkDescription != "AString" {
		t.Errorf("o.LinkDescription test failed %+v", o)
		return
	}

	if o.LinkVisible != "AString" {
		t.Errorf("o.LinkVisible test failed %+v", o)
		return
	}

	if o.LinkOwner != 999 {
		t.Errorf("o.LinkOwner test failed %+v", o)
		return
	}

	if o.LinkRating != 999 {
		t.Errorf("o.LinkRating test failed %+v", o)
		return
	}

	if o.LinkUpdated.Year != 2016 {
		t.Errorf("year not set for %+v", o.LinkUpdated)
		return
	}

	if o.LinkRel != "AString" {
		t.Errorf("o.LinkRel test failed %+v", o)
		return
	}

	if o.LinkNotes != "AString" {
		t.Errorf("o.LinkNotes test failed %+v", o)
		return
	}

	if o.LinkRss != "AString" {
		t.Errorf("o.LinkRss test failed %+v", o)
		return
	}
}

func TestNewOption(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewOption(a)
	if o._table != "wp_options" {
		t.Errorf("failed creating %+v", o)
		return
	}
}
func TestOptionFromDBValueMap(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewOption(a)
	m := make(map[string]DBValue)
	m["option_id"] = a.NewDBValue()
	m["option_id"].SetInternalValue("option_id", strconv.Itoa(999))
	m["option_name"] = a.NewDBValue()
	m["option_name"].SetInternalValue("option_name", "AString")
	m["option_value"] = a.NewDBValue()
	m["option_value"].SetInternalValue("option_value", "AString")
	m["autoload"] = a.NewDBValue()
	m["autoload"].SetInternalValue("autoload", "AString")

	err := o.FromDBValueMap(m)
	if err != nil {
		t.Errorf("FromDBValueMap failed %s", err)
	}

	if o.OptionId != 999 {
		t.Errorf("o.OptionId test failed %+v", o)
		return
	}

	if o.OptionName != "AString" {
		t.Errorf("o.OptionName test failed %+v", o)
		return
	}

	if o.OptionValue != "AString" {
		t.Errorf("o.OptionValue test failed %+v", o)
		return
	}

	if o.Autoload != "AString" {
		t.Errorf("o.Autoload test failed %+v", o)
		return
	}
}

func TestNewPostMeta(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewPostMeta(a)
	if o._table != "wp_postmeta" {
		t.Errorf("failed creating %+v", o)
		return
	}
}
func TestPostMetaFromDBValueMap(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewPostMeta(a)
	m := make(map[string]DBValue)
	m["meta_id"] = a.NewDBValue()
	m["meta_id"].SetInternalValue("meta_id", strconv.Itoa(999))
	m["post_id"] = a.NewDBValue()
	m["post_id"].SetInternalValue("post_id", strconv.Itoa(999))
	m["meta_key"] = a.NewDBValue()
	m["meta_key"].SetInternalValue("meta_key", "AString")
	m["meta_value"] = a.NewDBValue()
	m["meta_value"].SetInternalValue("meta_value", "AString")

	err := o.FromDBValueMap(m)
	if err != nil {
		t.Errorf("FromDBValueMap failed %s", err)
	}

	if o.MetaId != 999 {
		t.Errorf("o.MetaId test failed %+v", o)
		return
	}

	if o.Id != 999 {
		t.Errorf("o.Id test failed %+v", o)
		return
	}

	if o.MetaKey != "AString" {
		t.Errorf("o.MetaKey test failed %+v", o)
		return
	}

	if o.MetaValue != "AString" {
		t.Errorf("o.MetaValue test failed %+v", o)
		return
	}
}

func TestNewPost(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewPost(a)
	if o._table != "wp_posts" {
		t.Errorf("failed creating %+v", o)
		return
	}
}
func TestPostFromDBValueMap(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewPost(a)
	m := make(map[string]DBValue)
	m["ID"] = a.NewDBValue()
	m["ID"].SetInternalValue("ID", strconv.Itoa(999))
	m["post_author"] = a.NewDBValue()
	m["post_author"].SetInternalValue("post_author", strconv.Itoa(999))
	m["post_date"] = a.NewDBValue()
	m["post_date"].SetInternalValue("post_date", "2016-01-01 10:50:23.5Z")
	m["post_date_gmt"] = a.NewDBValue()
	m["post_date_gmt"].SetInternalValue("post_date_gmt", "2016-01-01 10:50:23.5Z")
	m["post_content"] = a.NewDBValue()
	m["post_content"].SetInternalValue("post_content", "AString")
	m["post_title"] = a.NewDBValue()
	m["post_title"].SetInternalValue("post_title", "AString")
	m["post_excerpt"] = a.NewDBValue()
	m["post_excerpt"].SetInternalValue("post_excerpt", "AString")
	m["post_status"] = a.NewDBValue()
	m["post_status"].SetInternalValue("post_status", "AString")
	m["comment_status"] = a.NewDBValue()
	m["comment_status"].SetInternalValue("comment_status", "AString")
	m["ping_status"] = a.NewDBValue()
	m["ping_status"].SetInternalValue("ping_status", "AString")
	m["post_password"] = a.NewDBValue()
	m["post_password"].SetInternalValue("post_password", "AString")
	m["post_name"] = a.NewDBValue()
	m["post_name"].SetInternalValue("post_name", "AString")
	m["to_ping"] = a.NewDBValue()
	m["to_ping"].SetInternalValue("to_ping", "AString")
	m["pinged"] = a.NewDBValue()
	m["pinged"].SetInternalValue("pinged", "AString")
	m["post_modified"] = a.NewDBValue()
	m["post_modified"].SetInternalValue("post_modified", "2016-01-01 10:50:23.5Z")
	m["post_modified_gmt"] = a.NewDBValue()
	m["post_modified_gmt"].SetInternalValue("post_modified_gmt", "2016-01-01 10:50:23.5Z")
	m["post_content_filtered"] = a.NewDBValue()
	m["post_content_filtered"].SetInternalValue("post_content_filtered", "AString")
	m["post_parent"] = a.NewDBValue()
	m["post_parent"].SetInternalValue("post_parent", strconv.Itoa(999))
	m["guid"] = a.NewDBValue()
	m["guid"].SetInternalValue("guid", "AString")
	m["menu_order"] = a.NewDBValue()
	m["menu_order"].SetInternalValue("menu_order", strconv.Itoa(999))
	m["post_type"] = a.NewDBValue()
	m["post_type"].SetInternalValue("post_type", "AString")
	m["post_mime_type"] = a.NewDBValue()
	m["post_mime_type"].SetInternalValue("post_mime_type", "AString")
	m["comment_count"] = a.NewDBValue()
	m["comment_count"].SetInternalValue("comment_count", strconv.Itoa(999))

	err := o.FromDBValueMap(m)
	if err != nil {
		t.Errorf("FromDBValueMap failed %s", err)
	}

	if o.ID != 999 {
		t.Errorf("o.ID test failed %+v", o)
		return
	}

	if o.Author != 999 {
		t.Errorf("o.Author test failed %+v", o)
		return
	}

	if o.Date.Year != 2016 {
		t.Errorf("year not set for %+v", o.Date)
		return
	}

	if o.DateGmt.Year != 2016 {
		t.Errorf("year not set for %+v", o.DateGmt)
		return
	}

	if o.Content != "AString" {
		t.Errorf("o.Content test failed %+v", o)
		return
	}

	if o.Title != "AString" {
		t.Errorf("o.Title test failed %+v", o)
		return
	}

	if o.Excerpt != "AString" {
		t.Errorf("o.Excerpt test failed %+v", o)
		return
	}

	if o.Status != "AString" {
		t.Errorf("o.Status test failed %+v", o)
		return
	}

	if o.CommentStatus != "AString" {
		t.Errorf("o.CommentStatus test failed %+v", o)
		return
	}

	if o.PingStatus != "AString" {
		t.Errorf("o.PingStatus test failed %+v", o)
		return
	}

	if o.Password != "AString" {
		t.Errorf("o.Password test failed %+v", o)
		return
	}

	if o.Name != "AString" {
		t.Errorf("o.Name test failed %+v", o)
		return
	}

	if o.ToPing != "AString" {
		t.Errorf("o.ToPing test failed %+v", o)
		return
	}

	if o.Pinged != "AString" {
		t.Errorf("o.Pinged test failed %+v", o)
		return
	}

	if o.Modified.Year != 2016 {
		t.Errorf("year not set for %+v", o.Modified)
		return
	}

	if o.ModifiedGmt.Year != 2016 {
		t.Errorf("year not set for %+v", o.ModifiedGmt)
		return
	}

	if o.ContentFiltered != "AString" {
		t.Errorf("o.ContentFiltered test failed %+v", o)
		return
	}

	if o.Parent != 999 {
		t.Errorf("o.Parent test failed %+v", o)
		return
	}

	if o.Guid != "AString" {
		t.Errorf("o.Guid test failed %+v", o)
		return
	}

	if o.MenuOrder != 999 {
		t.Errorf("o.MenuOrder test failed %+v", o)
		return
	}

	if o.Type != "AString" {
		t.Errorf("o.Type test failed %+v", o)
		return
	}

	if o.MimeType != "AString" {
		t.Errorf("o.MimeType test failed %+v", o)
		return
	}

	if o.CommentCount != 999 {
		t.Errorf("o.CommentCount test failed %+v", o)
		return
	}
}

func TestNewTermRelationship(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewTermRelationship(a)
	if o._table != "wp_term_relationships" {
		t.Errorf("failed creating %+v", o)
		return
	}
}
func TestTermRelationshipFromDBValueMap(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewTermRelationship(a)
	m := make(map[string]DBValue)
	m["object_id"] = a.NewDBValue()
	m["object_id"].SetInternalValue("object_id", strconv.Itoa(999))
	m["term_taxonomy_id"] = a.NewDBValue()
	m["term_taxonomy_id"].SetInternalValue("term_taxonomy_id", strconv.Itoa(999))
	m["term_order"] = a.NewDBValue()
	m["term_order"].SetInternalValue("term_order", strconv.Itoa(999))

	err := o.FromDBValueMap(m)
	if err != nil {
		t.Errorf("FromDBValueMap failed %s", err)
	}

	if o.ObjectId != 999 {
		t.Errorf("o.ObjectId test failed %+v", o)
		return
	}

	if o.TermTaxonomyId != 999 {
		t.Errorf("o.TermTaxonomyId test failed %+v", o)
		return
	}

	if o.TermOrder != 999 {
		t.Errorf("o.TermOrder test failed %+v", o)
		return
	}
}

func TestNewTermTaxonomy(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewTermTaxonomy(a)
	if o._table != "wp_term_taxonomy" {
		t.Errorf("failed creating %+v", o)
		return
	}
}
func TestTermTaxonomyFromDBValueMap(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewTermTaxonomy(a)
	m := make(map[string]DBValue)
	m["term_taxonomy_id"] = a.NewDBValue()
	m["term_taxonomy_id"].SetInternalValue("term_taxonomy_id", strconv.Itoa(999))
	m["term_id"] = a.NewDBValue()
	m["term_id"].SetInternalValue("term_id", strconv.Itoa(999))
	m["taxonomy"] = a.NewDBValue()
	m["taxonomy"].SetInternalValue("taxonomy", "AString")
	m["description"] = a.NewDBValue()
	m["description"].SetInternalValue("description", "AString")
	m["parent"] = a.NewDBValue()
	m["parent"].SetInternalValue("parent", strconv.Itoa(999))
	m["count"] = a.NewDBValue()
	m["count"].SetInternalValue("count", strconv.Itoa(999))

	err := o.FromDBValueMap(m)
	if err != nil {
		t.Errorf("FromDBValueMap failed %s", err)
	}

	if o.TermTaxonomyId != 999 {
		t.Errorf("o.TermTaxonomyId test failed %+v", o)
		return
	}

	if o.TermId != 999 {
		t.Errorf("o.TermId test failed %+v", o)
		return
	}

	if o.Taxonomy != "AString" {
		t.Errorf("o.Taxonomy test failed %+v", o)
		return
	}

	if o.Description != "AString" {
		t.Errorf("o.Description test failed %+v", o)
		return
	}

	if o.Parent != 999 {
		t.Errorf("o.Parent test failed %+v", o)
		return
	}

	if o.Count != 999 {
		t.Errorf("o.Count test failed %+v", o)
		return
	}
}

func TestNewTerm(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewTerm(a)
	if o._table != "wp_terms" {
		t.Errorf("failed creating %+v", o)
		return
	}
}
func TestTermFromDBValueMap(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewTerm(a)
	m := make(map[string]DBValue)
	m["term_id"] = a.NewDBValue()
	m["term_id"].SetInternalValue("term_id", strconv.Itoa(999))
	m["name"] = a.NewDBValue()
	m["name"].SetInternalValue("name", "AString")
	m["slug"] = a.NewDBValue()
	m["slug"].SetInternalValue("slug", "AString")
	m["term_group"] = a.NewDBValue()
	m["term_group"].SetInternalValue("term_group", strconv.Itoa(999))

	err := o.FromDBValueMap(m)
	if err != nil {
		t.Errorf("FromDBValueMap failed %s", err)
	}

	if o.TermId != 999 {
		t.Errorf("o.TermId test failed %+v", o)
		return
	}

	if o.Name != "AString" {
		t.Errorf("o.Name test failed %+v", o)
		return
	}

	if o.Slug != "AString" {
		t.Errorf("o.Slug test failed %+v", o)
		return
	}

	if o.TermGroup != 999 {
		t.Errorf("o.TermGroup test failed %+v", o)
		return
	}
}

func TestNewUserMeta(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewUserMeta(a)
	if o._table != "wp_usermeta" {
		t.Errorf("failed creating %+v", o)
		return
	}
}
func TestUserMetaFromDBValueMap(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewUserMeta(a)
	m := make(map[string]DBValue)
	m["umeta_id"] = a.NewDBValue()
	m["umeta_id"].SetInternalValue("umeta_id", strconv.Itoa(999))
	m["user_id"] = a.NewDBValue()
	m["user_id"].SetInternalValue("user_id", strconv.Itoa(999))
	m["meta_key"] = a.NewDBValue()
	m["meta_key"].SetInternalValue("meta_key", "AString")
	m["meta_value"] = a.NewDBValue()
	m["meta_value"].SetInternalValue("meta_value", "AString")

	err := o.FromDBValueMap(m)
	if err != nil {
		t.Errorf("FromDBValueMap failed %s", err)
	}

	if o.UMetaId != 999 {
		t.Errorf("o.UMetaId test failed %+v", o)
		return
	}

	if o.UserId != 999 {
		t.Errorf("o.UserId test failed %+v", o)
		return
	}

	if o.MetaKey != "AString" {
		t.Errorf("o.MetaKey test failed %+v", o)
		return
	}

	if o.MetaValue != "AString" {
		t.Errorf("o.MetaValue test failed %+v", o)
		return
	}
}

func TestNewUser(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewUser(a)
	if o._table != "wp_users" {
		t.Errorf("failed creating %+v", o)
		return
	}
}
func TestUserFromDBValueMap(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewUser(a)
	m := make(map[string]DBValue)
	m["ID"] = a.NewDBValue()
	m["ID"].SetInternalValue("ID", strconv.Itoa(999))
	m["user_login"] = a.NewDBValue()
	m["user_login"].SetInternalValue("user_login", "AString")
	m["user_pass"] = a.NewDBValue()
	m["user_pass"].SetInternalValue("user_pass", "AString")
	m["user_nicename"] = a.NewDBValue()
	m["user_nicename"].SetInternalValue("user_nicename", "AString")
	m["user_email"] = a.NewDBValue()
	m["user_email"].SetInternalValue("user_email", "AString")
	m["user_url"] = a.NewDBValue()
	m["user_url"].SetInternalValue("user_url", "AString")
	m["user_registered"] = a.NewDBValue()
	m["user_registered"].SetInternalValue("user_registered", "2016-01-01 10:50:23.5Z")
	m["user_activation_key"] = a.NewDBValue()
	m["user_activation_key"].SetInternalValue("user_activation_key", "AString")
	m["user_status"] = a.NewDBValue()
	m["user_status"].SetInternalValue("user_status", strconv.Itoa(999))
	m["display_name"] = a.NewDBValue()
	m["display_name"].SetInternalValue("display_name", "AString")

	err := o.FromDBValueMap(m)
	if err != nil {
		t.Errorf("FromDBValueMap failed %s", err)
	}

	if o.ID != 999 {
		t.Errorf("o.ID test failed %+v", o)
		return
	}

	if o.UserLogin != "AString" {
		t.Errorf("o.UserLogin test failed %+v", o)
		return
	}

	if o.UserPass != "AString" {
		t.Errorf("o.UserPass test failed %+v", o)
		return
	}

	if o.UserNicename != "AString" {
		t.Errorf("o.UserNicename test failed %+v", o)
		return
	}

	if o.UserEmail != "AString" {
		t.Errorf("o.UserEmail test failed %+v", o)
		return
	}

	if o.UserUrl != "AString" {
		t.Errorf("o.UserUrl test failed %+v", o)
		return
	}

	if o.UserRegistered.Year != 2016 {
		t.Errorf("year not set for %+v", o.UserRegistered)
		return
	}

	if o.UserActivationKey != "AString" {
		t.Errorf("o.UserActivationKey test failed %+v", o)
		return
	}

	if o.UserStatus != 999 {
		t.Errorf("o.UserStatus test failed %+v", o)
		return
	}

	if o.DisplayName != "AString" {
		t.Errorf("o.DisplayName test failed %+v", o)
		return
	}
}

func TestNewWooAttrTaxonomie(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewWooAttrTaxonomie(a)
	if o._table != "wp_woocommerce_attribute_taxonomies" {
		t.Errorf("failed creating %+v", o)
		return
	}
}
func TestWooAttrTaxonomieFromDBValueMap(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewWooAttrTaxonomie(a)
	m := make(map[string]DBValue)
	m["attribute_id"] = a.NewDBValue()
	m["attribute_id"].SetInternalValue("attribute_id", strconv.Itoa(999))
	m["attribute_name"] = a.NewDBValue()
	m["attribute_name"].SetInternalValue("attribute_name", "AString")
	m["attribute_label"] = a.NewDBValue()
	m["attribute_label"].SetInternalValue("attribute_label", "AString")
	m["attribute_type"] = a.NewDBValue()
	m["attribute_type"].SetInternalValue("attribute_type", "AString")
	m["attribute_orderby"] = a.NewDBValue()
	m["attribute_orderby"].SetInternalValue("attribute_orderby", "AString")

	err := o.FromDBValueMap(m)
	if err != nil {
		t.Errorf("FromDBValueMap failed %s", err)
	}

	if o.AttrId != 999 {
		t.Errorf("o.AttrId test failed %+v", o)
		return
	}

	if o.AttrName != "AString" {
		t.Errorf("o.AttrName test failed %+v", o)
		return
	}

	if o.AttrLabel != "AString" {
		t.Errorf("o.AttrLabel test failed %+v", o)
		return
	}

	if o.AttrType != "AString" {
		t.Errorf("o.AttrType test failed %+v", o)
		return
	}

	if o.AttrOrderby != "AString" {
		t.Errorf("o.AttrOrderby test failed %+v", o)
		return
	}
}

func TestNewWooDownloadableProductPerm(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewWooDownloadableProductPerm(a)
	if o._table != "wp_woocommerce_downloadable_product_permissions" {
		t.Errorf("failed creating %+v", o)
		return
	}
}
func TestWooDownloadableProductPermFromDBValueMap(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewWooDownloadableProductPerm(a)
	m := make(map[string]DBValue)
	m["permission_id"] = a.NewDBValue()
	m["permission_id"].SetInternalValue("permission_id", strconv.Itoa(999))
	m["download_id"] = a.NewDBValue()
	m["download_id"].SetInternalValue("download_id", "AString")
	m["product_id"] = a.NewDBValue()
	m["product_id"].SetInternalValue("product_id", strconv.Itoa(999))
	m["order_id"] = a.NewDBValue()
	m["order_id"].SetInternalValue("order_id", strconv.Itoa(999))
	m["order_key"] = a.NewDBValue()
	m["order_key"].SetInternalValue("order_key", "AString")
	m["user_email"] = a.NewDBValue()
	m["user_email"].SetInternalValue("user_email", "AString")
	m["user_id"] = a.NewDBValue()
	m["user_id"].SetInternalValue("user_id", strconv.Itoa(999))
	m["downloads_remaining"] = a.NewDBValue()
	m["downloads_remaining"].SetInternalValue("downloads_remaining", "AString")
	m["access_granted"] = a.NewDBValue()
	m["access_granted"].SetInternalValue("access_granted", "2016-01-01 10:50:23.5Z")
	m["access_expires"] = a.NewDBValue()
	m["access_expires"].SetInternalValue("access_expires", "2016-01-01 10:50:23.5Z")
	m["download_count"] = a.NewDBValue()
	m["download_count"].SetInternalValue("download_count", strconv.Itoa(999))

	err := o.FromDBValueMap(m)
	if err != nil {
		t.Errorf("FromDBValueMap failed %s", err)
	}

	if o.PermissionId != 999 {
		t.Errorf("o.PermissionId test failed %+v", o)
		return
	}

	if o.DownloadId != "AString" {
		t.Errorf("o.DownloadId test failed %+v", o)
		return
	}

	if o.ProductId != 999 {
		t.Errorf("o.ProductId test failed %+v", o)
		return
	}

	if o.OrderId != 999 {
		t.Errorf("o.OrderId test failed %+v", o)
		return
	}

	if o.OrderKey != "AString" {
		t.Errorf("o.OrderKey test failed %+v", o)
		return
	}

	if o.UserEmail != "AString" {
		t.Errorf("o.UserEmail test failed %+v", o)
		return
	}

	if o.UserId != 999 {
		t.Errorf("o.UserId test failed %+v", o)
		return
	}

	if o.DownloadsRemaining != "AString" {
		t.Errorf("o.DownloadsRemaining test failed %+v", o)
		return
	}

	if o.AccessGranted.Year != 2016 {
		t.Errorf("year not set for %+v", o.AccessGranted)
		return
	}

	if o.AccessExpires.Year != 2016 {
		t.Errorf("year not set for %+v", o.AccessExpires)
		return
	}

	if o.DownloadCount != 999 {
		t.Errorf("o.DownloadCount test failed %+v", o)
		return
	}
}

func TestNewWooOrderItemMeta(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewWooOrderItemMeta(a)
	if o._table != "wp_woocommerce_order_itemmeta" {
		t.Errorf("failed creating %+v", o)
		return
	}
}
func TestWooOrderItemMetaFromDBValueMap(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewWooOrderItemMeta(a)
	m := make(map[string]DBValue)
	m["meta_id"] = a.NewDBValue()
	m["meta_id"].SetInternalValue("meta_id", strconv.Itoa(999))
	m["order_item_id"] = a.NewDBValue()
	m["order_item_id"].SetInternalValue("order_item_id", strconv.Itoa(999))
	m["meta_key"] = a.NewDBValue()
	m["meta_key"].SetInternalValue("meta_key", "AString")
	m["meta_value"] = a.NewDBValue()
	m["meta_value"].SetInternalValue("meta_value", "AString")

	err := o.FromDBValueMap(m)
	if err != nil {
		t.Errorf("FromDBValueMap failed %s", err)
	}

	if o.MetaId != 999 {
		t.Errorf("o.MetaId test failed %+v", o)
		return
	}

	if o.OrderItemId != 999 {
		t.Errorf("o.OrderItemId test failed %+v", o)
		return
	}

	if o.MetaKey != "AString" {
		t.Errorf("o.MetaKey test failed %+v", o)
		return
	}

	if o.MetaValue != "AString" {
		t.Errorf("o.MetaValue test failed %+v", o)
		return
	}
}

func TestNewWooOrderItem(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewWooOrderItem(a)
	if o._table != "wp_woocommerce_order_items" {
		t.Errorf("failed creating %+v", o)
		return
	}
}
func TestWooOrderItemFromDBValueMap(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewWooOrderItem(a)
	m := make(map[string]DBValue)
	m["order_item_id"] = a.NewDBValue()
	m["order_item_id"].SetInternalValue("order_item_id", strconv.Itoa(999))
	m["order_item_name"] = a.NewDBValue()
	m["order_item_name"].SetInternalValue("order_item_name", "AString")
	m["order_item_type"] = a.NewDBValue()
	m["order_item_type"].SetInternalValue("order_item_type", "AString")
	m["order_id"] = a.NewDBValue()
	m["order_id"].SetInternalValue("order_id", strconv.Itoa(999))

	err := o.FromDBValueMap(m)
	if err != nil {
		t.Errorf("FromDBValueMap failed %s", err)
	}

	if o.OrderItemId != 999 {
		t.Errorf("o.OrderItemId test failed %+v", o)
		return
	}

	if o.OrderItemName != "AString" {
		t.Errorf("o.OrderItemName test failed %+v", o)
		return
	}

	if o.OrderItemType != "AString" {
		t.Errorf("o.OrderItemType test failed %+v", o)
		return
	}

	if o.OrderId != 999 {
		t.Errorf("o.OrderId test failed %+v", o)
		return
	}
}

func TestNewWooTaxRateLocation(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewWooTaxRateLocation(a)
	if o._table != "wp_woocommerce_tax_rate_locations" {
		t.Errorf("failed creating %+v", o)
		return
	}
}
func TestWooTaxRateLocationFromDBValueMap(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewWooTaxRateLocation(a)
	m := make(map[string]DBValue)
	m["location_id"] = a.NewDBValue()
	m["location_id"].SetInternalValue("location_id", strconv.Itoa(999))
	m["location_code"] = a.NewDBValue()
	m["location_code"].SetInternalValue("location_code", "AString")
	m["tax_rate_id"] = a.NewDBValue()
	m["tax_rate_id"].SetInternalValue("tax_rate_id", strconv.Itoa(999))
	m["location_type"] = a.NewDBValue()
	m["location_type"].SetInternalValue("location_type", "AString")

	err := o.FromDBValueMap(m)
	if err != nil {
		t.Errorf("FromDBValueMap failed %s", err)
	}

	if o.LocationId != 999 {
		t.Errorf("o.LocationId test failed %+v", o)
		return
	}

	if o.LocationCode != "AString" {
		t.Errorf("o.LocationCode test failed %+v", o)
		return
	}

	if o.TaxRateId != 999 {
		t.Errorf("o.TaxRateId test failed %+v", o)
		return
	}

	if o.LocationType != "AString" {
		t.Errorf("o.LocationType test failed %+v", o)
		return
	}
}

func TestNewWooTaxRate(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewWooTaxRate(a)
	if o._table != "wp_woocommerce_tax_rates" {
		t.Errorf("failed creating %+v", o)
		return
	}
}
func TestWooTaxRateFromDBValueMap(t *testing.T) {
	a := NewMysqlAdapter("wp_")
	o := NewWooTaxRate(a)
	m := make(map[string]DBValue)
	m["tax_rate_id"] = a.NewDBValue()
	m["tax_rate_id"].SetInternalValue("tax_rate_id", strconv.Itoa(999))
	m["tax_rate_country"] = a.NewDBValue()
	m["tax_rate_country"].SetInternalValue("tax_rate_country", "AString")
	m["tax_rate_state"] = a.NewDBValue()
	m["tax_rate_state"].SetInternalValue("tax_rate_state", "AString")
	m["tax_rate"] = a.NewDBValue()
	m["tax_rate"].SetInternalValue("tax_rate", "AString")
	m["tax_rate_name"] = a.NewDBValue()
	m["tax_rate_name"].SetInternalValue("tax_rate_name", "AString")
	m["tax_rate_priority"] = a.NewDBValue()
	m["tax_rate_priority"].SetInternalValue("tax_rate_priority", strconv.Itoa(999))
	m["tax_rate_compound"] = a.NewDBValue()
	m["tax_rate_compound"].SetInternalValue("tax_rate_compound", strconv.Itoa(999))
	m["tax_rate_shipping"] = a.NewDBValue()
	m["tax_rate_shipping"].SetInternalValue("tax_rate_shipping", strconv.Itoa(999))
	m["tax_rate_order"] = a.NewDBValue()
	m["tax_rate_order"].SetInternalValue("tax_rate_order", strconv.Itoa(999))
	m["tax_rate_class"] = a.NewDBValue()
	m["tax_rate_class"].SetInternalValue("tax_rate_class", "AString")

	err := o.FromDBValueMap(m)
	if err != nil {
		t.Errorf("FromDBValueMap failed %s", err)
	}

	if o.TaxRateId != 999 {
		t.Errorf("o.TaxRateId test failed %+v", o)
		return
	}

	if o.TaxRateCountry != "AString" {
		t.Errorf("o.TaxRateCountry test failed %+v", o)
		return
	}

	if o.TaxRateState != "AString" {
		t.Errorf("o.TaxRateState test failed %+v", o)
		return
	}

	if o.TaxRate != "AString" {
		t.Errorf("o.TaxRate test failed %+v", o)
		return
	}

	if o.TaxRateName != "AString" {
		t.Errorf("o.TaxRateName test failed %+v", o)
		return
	}

	if o.TaxRatePriority != 999 {
		t.Errorf("o.TaxRatePriority test failed %+v", o)
		return
	}

	if o.TaxRateCompound != 999 {
		t.Errorf("o.TaxRateCompound test failed %+v", o)
		return
	}

	if o.TaxRateShipping != 999 {
		t.Errorf("o.TaxRateShipping test failed %+v", o)
		return
	}

	if o.TaxRateOrder != 999 {
		t.Errorf("o.TaxRateOrder test failed %+v", o)
		return
	}

	if o.TaxRateClass != "AString" {
		t.Errorf("o.TaxRateClass test failed %+v", o)
		return
	}
}

func TestMysqlAdapterFromYAML(t *testing.T) {
	a := NewMysqlAdapter(`pw_`)
	y, err := fileGetContents(`test_data/adapter.yml`)
	if err != nil {
		t.Errorf(`failed to load yaml %s`, err)
	}
	err = a.FromYAML(y)
	if err != nil {
		t.Errorf(`failed to apply yaml %s`, err)
		return
	}

	if a.User != `root` ||
		a.Pass != `rootpass` ||
		a.Host != `localhost` ||
		a.Database != `my_db` ||
		a.DBPrefix != `wp_` {
		t.Errorf(`did not fully apply yaml file %+v`, a)
	}
}
func TestDBValue(t *testing.T) {
	a := NewMysqlAdapter(`wp_`)

	v0 := a.NewDBValue()
	v0.SetInternalValue(`x`, `999`)
	c0, err := v0.AsInt32()
	if err != nil {
		t.Errorf(`failed to convert with AsInt32() %+v`, v0)
		return
	}
	if c0 != 999 {
		t.Errorf(`values don't match `)
		return
	}

	v1 := a.NewDBValue()
	v1.SetInternalValue(`x`, `666`)
	c1, err := v1.AsInt()
	if err != nil {
		t.Errorf(`failed to convert with AsInt() %+v`, v1)
		return
	}
	if c1 != 666 {
		t.Errorf(`values don't match `)
		return
	}

	v2 := a.NewDBValue()
	v2.SetInternalValue(`x`, `hello world`)
	c2, err := v2.AsString()
	if err != nil {
		t.Errorf(`failed to convert with AsString() %+v`, v2)
		return
	}
	if c2 != "hello world" {
		t.Errorf(`values don't match `)
		return
	}

	v3 := a.NewDBValue()
	v3.SetInternalValue(`x`, `3.14`)
	c3, err := v3.AsFloat32()
	if err != nil {
		t.Errorf(`failed to convert with AsFloat32() %+v`, v3)
		return
	}
	if c3 != 3.14 {
		t.Errorf(`values don't match `)
		return
	}

	v4 := a.NewDBValue()
	v4.SetInternalValue(`x`, `67859.58686`)
	c4, err := v4.AsFloat32()
	if err != nil {
		t.Errorf(`failed to convert with AsFloat32() %+v`, v4)
		return
	}
	if c4 != 67859.58686 {
		t.Errorf(`values don't match `)
		return
	}

}
