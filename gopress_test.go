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

	if o.metaId != 999 {
		t.Errorf("o.metaId test failed %+v", o)
		return
	}

	if o.commentId != 999 {
		t.Errorf("o.commentId test failed %+v", o)
		return
	}

	if o.metaKey != "AString" {
		t.Errorf("o.metaKey test failed %+v", o)
		return
	}

	if o.metaValue != "AString" {
		t.Errorf("o.metaValue test failed %+v", o)
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

	if o.commentID != 999 {
		t.Errorf("o.commentID test failed %+v", o)
		return
	}

	if o.commentPostID != 999 {
		t.Errorf("o.commentPostID test failed %+v", o)
		return
	}

	if o.commentAuthor != "AString" {
		t.Errorf("o.commentAuthor test failed %+v", o)
		return
	}

	if o.commentAuthorEmail != "AString" {
		t.Errorf("o.commentAuthorEmail test failed %+v", o)
		return
	}

	if o.commentAuthorUrl != "AString" {
		t.Errorf("o.commentAuthorUrl test failed %+v", o)
		return
	}

	if o.commentAuthorIP != "AString" {
		t.Errorf("o.commentAuthorIP test failed %+v", o)
		return
	}

	if o.commentDate.Year != 2016 {
		t.Errorf("year not set for %+v", o.commentDate)
		return
	}

	if o.commentDateGmt.Year != 2016 {
		t.Errorf("year not set for %+v", o.commentDateGmt)
		return
	}

	if o.commentContent != "AString" {
		t.Errorf("o.commentContent test failed %+v", o)
		return
	}

	if o.commentKarma != 999 {
		t.Errorf("o.commentKarma test failed %+v", o)
		return
	}

	if o.commentApproved != "AString" {
		t.Errorf("o.commentApproved test failed %+v", o)
		return
	}

	if o.commentAgent != "AString" {
		t.Errorf("o.commentAgent test failed %+v", o)
		return
	}

	if o.commentType != "AString" {
		t.Errorf("o.commentType test failed %+v", o)
		return
	}

	if o.commentParent != 999 {
		t.Errorf("o.commentParent test failed %+v", o)
		return
	}

	if o.userId != 999 {
		t.Errorf("o.userId test failed %+v", o)
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

	if o.linkId != 999 {
		t.Errorf("o.linkId test failed %+v", o)
		return
	}

	if o.linkUrl != "AString" {
		t.Errorf("o.linkUrl test failed %+v", o)
		return
	}

	if o.linkName != "AString" {
		t.Errorf("o.linkName test failed %+v", o)
		return
	}

	if o.linkImage != "AString" {
		t.Errorf("o.linkImage test failed %+v", o)
		return
	}

	if o.linkTarget != "AString" {
		t.Errorf("o.linkTarget test failed %+v", o)
		return
	}

	if o.linkDescription != "AString" {
		t.Errorf("o.linkDescription test failed %+v", o)
		return
	}

	if o.linkVisible != "AString" {
		t.Errorf("o.linkVisible test failed %+v", o)
		return
	}

	if o.linkOwner != 999 {
		t.Errorf("o.linkOwner test failed %+v", o)
		return
	}

	if o.linkRating != 999 {
		t.Errorf("o.linkRating test failed %+v", o)
		return
	}

	if o.linkUpdated.Year != 2016 {
		t.Errorf("year not set for %+v", o.linkUpdated)
		return
	}

	if o.linkRel != "AString" {
		t.Errorf("o.linkRel test failed %+v", o)
		return
	}

	if o.linkNotes != "AString" {
		t.Errorf("o.linkNotes test failed %+v", o)
		return
	}

	if o.linkRss != "AString" {
		t.Errorf("o.linkRss test failed %+v", o)
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

	if o.optionId != 999 {
		t.Errorf("o.optionId test failed %+v", o)
		return
	}

	if o.optionName != "AString" {
		t.Errorf("o.optionName test failed %+v", o)
		return
	}

	if o.optionValue != "AString" {
		t.Errorf("o.optionValue test failed %+v", o)
		return
	}

	if o.autoload != "AString" {
		t.Errorf("o.autoload test failed %+v", o)
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

	if o.metaId != 999 {
		t.Errorf("o.metaId test failed %+v", o)
		return
	}

	if o.stId != 999 {
		t.Errorf("o.stId test failed %+v", o)
		return
	}

	if o.metaKey != "AString" {
		t.Errorf("o.metaKey test failed %+v", o)
		return
	}

	if o.metaValue != "AString" {
		t.Errorf("o.metaValue test failed %+v", o)
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

	if o.iD != 999 {
		t.Errorf("o.iD test failed %+v", o)
		return
	}

	if o.stAuthor != 999 {
		t.Errorf("o.stAuthor test failed %+v", o)
		return
	}

	if o.stDate.Year != 2016 {
		t.Errorf("year not set for %+v", o.stDate)
		return
	}

	if o.stDateGmt.Year != 2016 {
		t.Errorf("year not set for %+v", o.stDateGmt)
		return
	}

	if o.stContent != "AString" {
		t.Errorf("o.stContent test failed %+v", o)
		return
	}

	if o.stTitle != "AString" {
		t.Errorf("o.stTitle test failed %+v", o)
		return
	}

	if o.stExcerpt != "AString" {
		t.Errorf("o.stExcerpt test failed %+v", o)
		return
	}

	if o.stStatus != "AString" {
		t.Errorf("o.stStatus test failed %+v", o)
		return
	}

	if o.commentStatus != "AString" {
		t.Errorf("o.commentStatus test failed %+v", o)
		return
	}

	if o.pingStatus != "AString" {
		t.Errorf("o.pingStatus test failed %+v", o)
		return
	}

	if o.stPassword != "AString" {
		t.Errorf("o.stPassword test failed %+v", o)
		return
	}

	if o.stName != "AString" {
		t.Errorf("o.stName test failed %+v", o)
		return
	}

	if o.toPing != "AString" {
		t.Errorf("o.toPing test failed %+v", o)
		return
	}

	if o.pinged != "AString" {
		t.Errorf("o.pinged test failed %+v", o)
		return
	}

	if o.stModified.Year != 2016 {
		t.Errorf("year not set for %+v", o.stModified)
		return
	}

	if o.stModifiedGmt.Year != 2016 {
		t.Errorf("year not set for %+v", o.stModifiedGmt)
		return
	}

	if o.stContentFiltered != "AString" {
		t.Errorf("o.stContentFiltered test failed %+v", o)
		return
	}

	if o.stParent != 999 {
		t.Errorf("o.stParent test failed %+v", o)
		return
	}

	if o.guid != "AString" {
		t.Errorf("o.guid test failed %+v", o)
		return
	}

	if o.menuOrder != 999 {
		t.Errorf("o.menuOrder test failed %+v", o)
		return
	}

	if o.stType != "AString" {
		t.Errorf("o.stType test failed %+v", o)
		return
	}

	if o.stMimeType != "AString" {
		t.Errorf("o.stMimeType test failed %+v", o)
		return
	}

	if o.commentCount != 999 {
		t.Errorf("o.commentCount test failed %+v", o)
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

	if o.objectId != 999 {
		t.Errorf("o.objectId test failed %+v", o)
		return
	}

	if o.termTaxonomyId != 999 {
		t.Errorf("o.termTaxonomyId test failed %+v", o)
		return
	}

	if o.termOrder != 999 {
		t.Errorf("o.termOrder test failed %+v", o)
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

	if o.termTaxonomyId != 999 {
		t.Errorf("o.termTaxonomyId test failed %+v", o)
		return
	}

	if o.termId != 999 {
		t.Errorf("o.termId test failed %+v", o)
		return
	}

	if o.taxonomy != "AString" {
		t.Errorf("o.taxonomy test failed %+v", o)
		return
	}

	if o.description != "AString" {
		t.Errorf("o.description test failed %+v", o)
		return
	}

	if o.parent != 999 {
		t.Errorf("o.parent test failed %+v", o)
		return
	}

	if o.count != 999 {
		t.Errorf("o.count test failed %+v", o)
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

	if o.termId != 999 {
		t.Errorf("o.termId test failed %+v", o)
		return
	}

	if o.name != "AString" {
		t.Errorf("o.name test failed %+v", o)
		return
	}

	if o.slug != "AString" {
		t.Errorf("o.slug test failed %+v", o)
		return
	}

	if o.termGroup != 999 {
		t.Errorf("o.termGroup test failed %+v", o)
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

	if o.uMetaId != 999 {
		t.Errorf("o.uMetaId test failed %+v", o)
		return
	}

	if o.userId != 999 {
		t.Errorf("o.userId test failed %+v", o)
		return
	}

	if o.metaKey != "AString" {
		t.Errorf("o.metaKey test failed %+v", o)
		return
	}

	if o.metaValue != "AString" {
		t.Errorf("o.metaValue test failed %+v", o)
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

	if o.iD != 999 {
		t.Errorf("o.iD test failed %+v", o)
		return
	}

	if o.userLogin != "AString" {
		t.Errorf("o.userLogin test failed %+v", o)
		return
	}

	if o.userPass != "AString" {
		t.Errorf("o.userPass test failed %+v", o)
		return
	}

	if o.userNicename != "AString" {
		t.Errorf("o.userNicename test failed %+v", o)
		return
	}

	if o.userEmail != "AString" {
		t.Errorf("o.userEmail test failed %+v", o)
		return
	}

	if o.userUrl != "AString" {
		t.Errorf("o.userUrl test failed %+v", o)
		return
	}

	if o.userRegistered.Year != 2016 {
		t.Errorf("year not set for %+v", o.userRegistered)
		return
	}

	if o.userActivationKey != "AString" {
		t.Errorf("o.userActivationKey test failed %+v", o)
		return
	}

	if o.userStatus != 999 {
		t.Errorf("o.userStatus test failed %+v", o)
		return
	}

	if o.displayName != "AString" {
		t.Errorf("o.displayName test failed %+v", o)
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

	if o.attrId != 999 {
		t.Errorf("o.attrId test failed %+v", o)
		return
	}

	if o.attrName != "AString" {
		t.Errorf("o.attrName test failed %+v", o)
		return
	}

	if o.attrLabel != "AString" {
		t.Errorf("o.attrLabel test failed %+v", o)
		return
	}

	if o.attrType != "AString" {
		t.Errorf("o.attrType test failed %+v", o)
		return
	}

	if o.attrOrderby != "AString" {
		t.Errorf("o.attrOrderby test failed %+v", o)
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

	if o.permissionId != 999 {
		t.Errorf("o.permissionId test failed %+v", o)
		return
	}

	if o.downloadId != "AString" {
		t.Errorf("o.downloadId test failed %+v", o)
		return
	}

	if o.productId != 999 {
		t.Errorf("o.productId test failed %+v", o)
		return
	}

	if o.orderId != 999 {
		t.Errorf("o.orderId test failed %+v", o)
		return
	}

	if o.orderKey != "AString" {
		t.Errorf("o.orderKey test failed %+v", o)
		return
	}

	if o.userEmail != "AString" {
		t.Errorf("o.userEmail test failed %+v", o)
		return
	}

	if o.userId != 999 {
		t.Errorf("o.userId test failed %+v", o)
		return
	}

	if o.downloadsRemaining != "AString" {
		t.Errorf("o.downloadsRemaining test failed %+v", o)
		return
	}

	if o.accessGranted.Year != 2016 {
		t.Errorf("year not set for %+v", o.accessGranted)
		return
	}

	if o.accessExpires.Year != 2016 {
		t.Errorf("year not set for %+v", o.accessExpires)
		return
	}

	if o.downloadCount != 999 {
		t.Errorf("o.downloadCount test failed %+v", o)
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

	if o.metaId != 999 {
		t.Errorf("o.metaId test failed %+v", o)
		return
	}

	if o.orderItemId != 999 {
		t.Errorf("o.orderItemId test failed %+v", o)
		return
	}

	if o.metaKey != "AString" {
		t.Errorf("o.metaKey test failed %+v", o)
		return
	}

	if o.metaValue != "AString" {
		t.Errorf("o.metaValue test failed %+v", o)
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

	if o.orderItemId != 999 {
		t.Errorf("o.orderItemId test failed %+v", o)
		return
	}

	if o.orderItemName != "AString" {
		t.Errorf("o.orderItemName test failed %+v", o)
		return
	}

	if o.orderItemType != "AString" {
		t.Errorf("o.orderItemType test failed %+v", o)
		return
	}

	if o.orderId != 999 {
		t.Errorf("o.orderId test failed %+v", o)
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

	if o.locationId != 999 {
		t.Errorf("o.locationId test failed %+v", o)
		return
	}

	if o.locationCode != "AString" {
		t.Errorf("o.locationCode test failed %+v", o)
		return
	}

	if o.taxRateId != 999 {
		t.Errorf("o.taxRateId test failed %+v", o)
		return
	}

	if o.locationType != "AString" {
		t.Errorf("o.locationType test failed %+v", o)
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

	if o.taxRateId != 999 {
		t.Errorf("o.taxRateId test failed %+v", o)
		return
	}

	if o.taxRateCountry != "AString" {
		t.Errorf("o.taxRateCountry test failed %+v", o)
		return
	}

	if o.taxRateState != "AString" {
		t.Errorf("o.taxRateState test failed %+v", o)
		return
	}

	if o.taxRate != "AString" {
		t.Errorf("o.taxRate test failed %+v", o)
		return
	}

	if o.taxRateName != "AString" {
		t.Errorf("o.taxRateName test failed %+v", o)
		return
	}

	if o.taxRatePriority != 999 {
		t.Errorf("o.taxRatePriority test failed %+v", o)
		return
	}

	if o.taxRateCompound != 999 {
		t.Errorf("o.taxRateCompound test failed %+v", o)
		return
	}

	if o.taxRateShipping != 999 {
		t.Errorf("o.taxRateShipping test failed %+v", o)
		return
	}

	if o.taxRateOrder != 999 {
		t.Errorf("o.taxRateOrder test failed %+v", o)
		return
	}

	if o.taxRateClass != "AString" {
		t.Errorf("o.taxRateClass test failed %+v", o)
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
