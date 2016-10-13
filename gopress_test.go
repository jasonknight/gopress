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

func TestCommentMetaCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		model := NewCommentMeta(a)

		model.CommentId = 999

		model.MetaKey = `the rain in spain`

		model.MetaValue = `the rain in spain`

		i, err := model.Create()
		if err != nil {
			t.Errorf(`failed to create model %+v error: %s`, model, err)
			return
		}
		if i == 0 {
			t.Errorf(`zero affected rows`)
			return
		}
		model2 := NewCommentMeta(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue())
			return
		}

	} // end of if fileExists
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
	if o.CommentDate.Year != 2016 ||
		o.CommentDate.Month != 1 ||
		o.CommentDate.Day != 1 ||
		o.CommentDate.Hours != 10 ||
		o.CommentDate.Minutes != 50 ||
		o.CommentDate.Seconds != 23 ||
		o.CommentDate.Offset != 5 ||
		o.CommentDate.Zone != `Z` {
		t.Errorf(`fields don't match up for %+v`, o.CommentDate)
	}
	r6, _ := m["comment_date"].AsString()
	if o.CommentDate.ToString() != r6 {
		t.Errorf(`restring of o.CommentDate failed %s`, o.CommentDate.ToString())
	}

	if o.CommentDateGmt.Year != 2016 {
		t.Errorf("year not set for %+v", o.CommentDateGmt)
		return
	}
	if o.CommentDateGmt.Year != 2016 ||
		o.CommentDateGmt.Month != 1 ||
		o.CommentDateGmt.Day != 1 ||
		o.CommentDateGmt.Hours != 10 ||
		o.CommentDateGmt.Minutes != 50 ||
		o.CommentDateGmt.Seconds != 23 ||
		o.CommentDateGmt.Offset != 5 ||
		o.CommentDateGmt.Zone != `Z` {
		t.Errorf(`fields don't match up for %+v`, o.CommentDateGmt)
	}
	r7, _ := m["comment_date_gmt"].AsString()
	if o.CommentDateGmt.ToString() != r7 {
		t.Errorf(`restring of o.CommentDateGmt failed %s`, o.CommentDateGmt.ToString())
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

func TestCommentCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		model := NewComment(a)

		model.CommentPostID = 999

		model.CommentAuthor = `the rain in spain`

		model.CommentAuthorEmail = `the rain in spain`

		model.CommentAuthorUrl = `the rain in spain`

		model.CommentAuthorIP = `the rain in spain`

		d5 := NewDateTime()
		d5.FromString(`the rain in spain`)
		model.CommentDate = d5

		d6 := NewDateTime()
		d6.FromString(`the rain in spain`)
		model.CommentDateGmt = d6

		model.CommentContent = `the rain in spain`

		model.CommentKarma = 999

		model.CommentApproved = `the rain in spain`

		model.CommentAgent = `the rain in spain`

		model.CommentType = `the rain in spain`

		model.CommentParent = 999

		model.UserId = 999

		i, err := model.Create()
		if err != nil {
			t.Errorf(`failed to create model %+v error: %s`, model, err)
			return
		}
		if i == 0 {
			t.Errorf(`zero affected rows`)
			return
		}
		model2 := NewComment(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue())
			return
		}

	} // end of if fileExists
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
	if o.LinkUpdated.Year != 2016 ||
		o.LinkUpdated.Month != 1 ||
		o.LinkUpdated.Day != 1 ||
		o.LinkUpdated.Hours != 10 ||
		o.LinkUpdated.Minutes != 50 ||
		o.LinkUpdated.Seconds != 23 ||
		o.LinkUpdated.Offset != 5 ||
		o.LinkUpdated.Zone != `Z` {
		t.Errorf(`fields don't match up for %+v`, o.LinkUpdated)
	}
	r9, _ := m["link_updated"].AsString()
	if o.LinkUpdated.ToString() != r9 {
		t.Errorf(`restring of o.LinkUpdated failed %s`, o.LinkUpdated.ToString())
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

func TestLinkCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		model := NewLink(a)

		model.LinkUrl = `the rain in spain`

		model.LinkName = `the rain in spain`

		model.LinkImage = `the rain in spain`

		model.LinkTarget = `the rain in spain`

		model.LinkDescription = `the rain in spain`

		model.LinkVisible = `the rain in spain`

		model.LinkOwner = 999

		model.LinkRating = 999

		d8 := NewDateTime()
		d8.FromString(`999`)
		model.LinkUpdated = d8

		model.LinkRel = `the rain in spain`

		model.LinkNotes = `the rain in spain`

		model.LinkRss = `the rain in spain`

		i, err := model.Create()
		if err != nil {
			t.Errorf(`failed to create model %+v error: %s`, model, err)
			return
		}
		if i == 0 {
			t.Errorf(`zero affected rows`)
			return
		}
		model2 := NewLink(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue())
			return
		}

	} // end of if fileExists
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

func TestOptionCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		model := NewOption(a)

		model.OptionName = `the rain in spain`

		model.OptionValue = `the rain in spain`

		model.Autoload = `the rain in spain`

		i, err := model.Create()
		if err != nil {
			t.Errorf(`failed to create model %+v error: %s`, model, err)
			return
		}
		if i == 0 {
			t.Errorf(`zero affected rows`)
			return
		}
		model2 := NewOption(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue())
			return
		}

	} // end of if fileExists
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

func TestPostMetaCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		model := NewPostMeta(a)

		model.Id = 999

		model.MetaKey = `the rain in spain`

		model.MetaValue = `the rain in spain`

		i, err := model.Create()
		if err != nil {
			t.Errorf(`failed to create model %+v error: %s`, model, err)
			return
		}
		if i == 0 {
			t.Errorf(`zero affected rows`)
			return
		}
		model2 := NewPostMeta(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue())
			return
		}

	} // end of if fileExists
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
	if o.Date.Year != 2016 ||
		o.Date.Month != 1 ||
		o.Date.Day != 1 ||
		o.Date.Hours != 10 ||
		o.Date.Minutes != 50 ||
		o.Date.Seconds != 23 ||
		o.Date.Offset != 5 ||
		o.Date.Zone != `Z` {
		t.Errorf(`fields don't match up for %+v`, o.Date)
	}
	r2, _ := m["post_date"].AsString()
	if o.Date.ToString() != r2 {
		t.Errorf(`restring of o.Date failed %s`, o.Date.ToString())
	}

	if o.DateGmt.Year != 2016 {
		t.Errorf("year not set for %+v", o.DateGmt)
		return
	}
	if o.DateGmt.Year != 2016 ||
		o.DateGmt.Month != 1 ||
		o.DateGmt.Day != 1 ||
		o.DateGmt.Hours != 10 ||
		o.DateGmt.Minutes != 50 ||
		o.DateGmt.Seconds != 23 ||
		o.DateGmt.Offset != 5 ||
		o.DateGmt.Zone != `Z` {
		t.Errorf(`fields don't match up for %+v`, o.DateGmt)
	}
	r3, _ := m["post_date_gmt"].AsString()
	if o.DateGmt.ToString() != r3 {
		t.Errorf(`restring of o.DateGmt failed %s`, o.DateGmt.ToString())
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
	if o.Modified.Year != 2016 ||
		o.Modified.Month != 1 ||
		o.Modified.Day != 1 ||
		o.Modified.Hours != 10 ||
		o.Modified.Minutes != 50 ||
		o.Modified.Seconds != 23 ||
		o.Modified.Offset != 5 ||
		o.Modified.Zone != `Z` {
		t.Errorf(`fields don't match up for %+v`, o.Modified)
	}
	r14, _ := m["post_modified"].AsString()
	if o.Modified.ToString() != r14 {
		t.Errorf(`restring of o.Modified failed %s`, o.Modified.ToString())
	}

	if o.ModifiedGmt.Year != 2016 {
		t.Errorf("year not set for %+v", o.ModifiedGmt)
		return
	}
	if o.ModifiedGmt.Year != 2016 ||
		o.ModifiedGmt.Month != 1 ||
		o.ModifiedGmt.Day != 1 ||
		o.ModifiedGmt.Hours != 10 ||
		o.ModifiedGmt.Minutes != 50 ||
		o.ModifiedGmt.Seconds != 23 ||
		o.ModifiedGmt.Offset != 5 ||
		o.ModifiedGmt.Zone != `Z` {
		t.Errorf(`fields don't match up for %+v`, o.ModifiedGmt)
	}
	r15, _ := m["post_modified_gmt"].AsString()
	if o.ModifiedGmt.ToString() != r15 {
		t.Errorf(`restring of o.ModifiedGmt failed %s`, o.ModifiedGmt.ToString())
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

func TestPostCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		model := NewPost(a)

		model.Author = 999

		d1 := NewDateTime()
		d1.FromString(`999`)
		model.Date = d1

		d2 := NewDateTime()
		d2.FromString(`999`)
		model.DateGmt = d2

		model.Content = `the rain in spain`

		model.Title = `the rain in spain`

		model.Excerpt = `the rain in spain`

		model.Status = `the rain in spain`

		model.CommentStatus = `the rain in spain`

		model.PingStatus = `the rain in spain`

		model.Password = `the rain in spain`

		model.Name = `the rain in spain`

		model.ToPing = `the rain in spain`

		model.Pinged = `the rain in spain`

		d13 := NewDateTime()
		d13.FromString(`the rain in spain`)
		model.Modified = d13

		d14 := NewDateTime()
		d14.FromString(`the rain in spain`)
		model.ModifiedGmt = d14

		model.ContentFiltered = `the rain in spain`

		model.Parent = 999

		model.Guid = `the rain in spain`

		model.MenuOrder = 999

		model.Type = `the rain in spain`

		model.MimeType = `the rain in spain`

		model.CommentCount = 999

		i, err := model.Create()
		if err != nil {
			t.Errorf(`failed to create model %+v error: %s`, model, err)
			return
		}
		if i == 0 {
			t.Errorf(`zero affected rows`)
			return
		}
		model2 := NewPost(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue())
			return
		}

	} // end of if fileExists
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

func TestTermRelationshipCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		model := NewTermRelationship(a)

		model.TermOrder = 999

		i, err := model.Create()
		if err != nil {
			t.Errorf(`failed to create model %+v error: %s`, model, err)
			return
		}
		if i == 0 {
			t.Errorf(`zero affected rows`)
			return
		}
		model2 := NewTermRelationship(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue())
			return
		}

	} // end of if fileExists
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

func TestTermTaxonomyCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		model := NewTermTaxonomy(a)

		model.TermId = 999

		model.Taxonomy = `the rain in spain`

		model.Description = `the rain in spain`

		model.Parent = 999

		model.Count = 999

		i, err := model.Create()
		if err != nil {
			t.Errorf(`failed to create model %+v error: %s`, model, err)
			return
		}
		if i == 0 {
			t.Errorf(`zero affected rows`)
			return
		}
		model2 := NewTermTaxonomy(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue())
			return
		}

	} // end of if fileExists
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

func TestTermCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		model := NewTerm(a)

		model.Name = `the rain in spain`

		model.Slug = `the rain in spain`

		model.TermGroup = 999

		i, err := model.Create()
		if err != nil {
			t.Errorf(`failed to create model %+v error: %s`, model, err)
			return
		}
		if i == 0 {
			t.Errorf(`zero affected rows`)
			return
		}
		model2 := NewTerm(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue())
			return
		}

	} // end of if fileExists
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

func TestUserMetaCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		model := NewUserMeta(a)

		model.UserId = 999

		model.MetaKey = `the rain in spain`

		model.MetaValue = `the rain in spain`

		i, err := model.Create()
		if err != nil {
			t.Errorf(`failed to create model %+v error: %s`, model, err)
			return
		}
		if i == 0 {
			t.Errorf(`zero affected rows`)
			return
		}
		model2 := NewUserMeta(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue())
			return
		}

	} // end of if fileExists
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
	if o.UserRegistered.Year != 2016 ||
		o.UserRegistered.Month != 1 ||
		o.UserRegistered.Day != 1 ||
		o.UserRegistered.Hours != 10 ||
		o.UserRegistered.Minutes != 50 ||
		o.UserRegistered.Seconds != 23 ||
		o.UserRegistered.Offset != 5 ||
		o.UserRegistered.Zone != `Z` {
		t.Errorf(`fields don't match up for %+v`, o.UserRegistered)
	}
	r6, _ := m["user_registered"].AsString()
	if o.UserRegistered.ToString() != r6 {
		t.Errorf(`restring of o.UserRegistered failed %s`, o.UserRegistered.ToString())
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

func TestUserCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		model := NewUser(a)

		model.UserLogin = `the rain in spain`

		model.UserPass = `the rain in spain`

		model.UserNicename = `the rain in spain`

		model.UserEmail = `the rain in spain`

		model.UserUrl = `the rain in spain`

		d5 := NewDateTime()
		d5.FromString(`the rain in spain`)
		model.UserRegistered = d5

		model.UserActivationKey = `the rain in spain`

		model.UserStatus = 999

		model.DisplayName = `the rain in spain`

		i, err := model.Create()
		if err != nil {
			t.Errorf(`failed to create model %+v error: %s`, model, err)
			return
		}
		if i == 0 {
			t.Errorf(`zero affected rows`)
			return
		}
		model2 := NewUser(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue())
			return
		}

	} // end of if fileExists
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

func TestWooAttrTaxonomieCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		model := NewWooAttrTaxonomie(a)

		model.AttrName = `the rain in spain`

		model.AttrLabel = `the rain in spain`

		model.AttrType = `the rain in spain`

		model.AttrOrderby = `the rain in spain`

		i, err := model.Create()
		if err != nil {
			t.Errorf(`failed to create model %+v error: %s`, model, err)
			return
		}
		if i == 0 {
			t.Errorf(`zero affected rows`)
			return
		}
		model2 := NewWooAttrTaxonomie(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue())
			return
		}

	} // end of if fileExists
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
	if o.AccessGranted.Year != 2016 ||
		o.AccessGranted.Month != 1 ||
		o.AccessGranted.Day != 1 ||
		o.AccessGranted.Hours != 10 ||
		o.AccessGranted.Minutes != 50 ||
		o.AccessGranted.Seconds != 23 ||
		o.AccessGranted.Offset != 5 ||
		o.AccessGranted.Zone != `Z` {
		t.Errorf(`fields don't match up for %+v`, o.AccessGranted)
	}
	r8, _ := m["access_granted"].AsString()
	if o.AccessGranted.ToString() != r8 {
		t.Errorf(`restring of o.AccessGranted failed %s`, o.AccessGranted.ToString())
	}

	if o.AccessExpires.Year != 2016 {
		t.Errorf("year not set for %+v", o.AccessExpires)
		return
	}
	if o.AccessExpires.Year != 2016 ||
		o.AccessExpires.Month != 1 ||
		o.AccessExpires.Day != 1 ||
		o.AccessExpires.Hours != 10 ||
		o.AccessExpires.Minutes != 50 ||
		o.AccessExpires.Seconds != 23 ||
		o.AccessExpires.Offset != 5 ||
		o.AccessExpires.Zone != `Z` {
		t.Errorf(`fields don't match up for %+v`, o.AccessExpires)
	}
	r9, _ := m["access_expires"].AsString()
	if o.AccessExpires.ToString() != r9 {
		t.Errorf(`restring of o.AccessExpires failed %s`, o.AccessExpires.ToString())
	}

	if o.DownloadCount != 999 {
		t.Errorf("o.DownloadCount test failed %+v", o)
		return
	}
}

func TestWooDownloadableProductPermCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		model := NewWooDownloadableProductPerm(a)

		model.DownloadId = `the rain in spain`

		model.ProductId = 999

		model.OrderId = 999

		model.OrderKey = `the rain in spain`

		model.UserEmail = `the rain in spain`

		model.UserId = 999

		model.DownloadsRemaining = `the rain in spain`

		d7 := NewDateTime()
		d7.FromString(`the rain in spain`)
		model.AccessGranted = d7

		d8 := NewDateTime()
		d8.FromString(`the rain in spain`)
		model.AccessExpires = d8

		model.DownloadCount = 999

		i, err := model.Create()
		if err != nil {
			t.Errorf(`failed to create model %+v error: %s`, model, err)
			return
		}
		if i == 0 {
			t.Errorf(`zero affected rows`)
			return
		}
		model2 := NewWooDownloadableProductPerm(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue())
			return
		}

	} // end of if fileExists
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

func TestWooOrderItemMetaCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		model := NewWooOrderItemMeta(a)

		model.OrderItemId = 999

		model.MetaKey = `the rain in spain`

		model.MetaValue = `the rain in spain`

		i, err := model.Create()
		if err != nil {
			t.Errorf(`failed to create model %+v error: %s`, model, err)
			return
		}
		if i == 0 {
			t.Errorf(`zero affected rows`)
			return
		}
		model2 := NewWooOrderItemMeta(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue())
			return
		}

	} // end of if fileExists
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

func TestWooOrderItemCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		model := NewWooOrderItem(a)

		model.OrderItemName = `the rain in spain`

		model.OrderItemType = `the rain in spain`

		model.OrderId = 999

		i, err := model.Create()
		if err != nil {
			t.Errorf(`failed to create model %+v error: %s`, model, err)
			return
		}
		if i == 0 {
			t.Errorf(`zero affected rows`)
			return
		}
		model2 := NewWooOrderItem(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue())
			return
		}

	} // end of if fileExists
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

func TestWooTaxRateLocationCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		model := NewWooTaxRateLocation(a)

		model.LocationCode = `the rain in spain`

		model.TaxRateId = 999

		model.LocationType = `the rain in spain`

		i, err := model.Create()
		if err != nil {
			t.Errorf(`failed to create model %+v error: %s`, model, err)
			return
		}
		if i == 0 {
			t.Errorf(`zero affected rows`)
			return
		}
		model2 := NewWooTaxRateLocation(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue())
			return
		}

	} // end of if fileExists
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

func TestWooTaxRateCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		model := NewWooTaxRate(a)

		model.TaxRateCountry = `the rain in spain`

		model.TaxRateState = `the rain in spain`

		model.TaxRate = `the rain in spain`

		model.TaxRateName = `the rain in spain`

		model.TaxRatePriority = 999

		model.TaxRateCompound = 999

		model.TaxRateShipping = 999

		model.TaxRateOrder = 999

		model.TaxRateClass = `the rain in spain`

		i, err := model.Create()
		if err != nil {
			t.Errorf(`failed to create model %+v error: %s`, model, err)
			return
		}
		if i == 0 {
			t.Errorf(`zero affected rows`)
			return
		}
		model2 := NewWooTaxRate(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue())
			return
		}

	} // end of if fileExists
}

func TestMysqlAdapterFromYAML(t *testing.T) {
	a := NewMysqlAdapter(`pw_`)
	y, err := fileGetContents(`test_data/adapter.yml`)
	if err != nil {
		t.Errorf(`failed to load yaml %s`, err)
		return
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
	c4, err := v4.AsFloat64()
	if err != nil {
		t.Errorf(`failed to convert with AsFloat64() %+v`, v4)
		return
	}
	if c4 != 67859.58686 {
		t.Errorf(`values don't match `)
		return
	}

	dvar := a.NewDBValue()
	dvar.SetInternalValue(`x`, `2016-01-09 23:24:50.7Z`)
	dc, err := dvar.AsDateTime()
	if err != nil {
		t.Errorf(`failed to convert datetime %+v`, dc)
	}

	if dc.Year != 2016 ||
		dc.Month != 1 ||
		dc.Day != 9 ||
		dc.Hours != 23 ||
		dc.Minutes != 24 ||
		dc.Seconds != 50 ||
		dc.Offset != 7 ||
		dc.Zone != `Z` {
		t.Errorf(`fields don't match up for %+v`, dc)
	}
	r, _ := dvar.AsString()
	if dc.ToString() != r {
		t.Errorf(`restring of dvar failed %s`, dc.ToString())
	}

}
