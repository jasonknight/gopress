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

		if model.CommentId != model2.CommentId {
			t.Errorf(`model.CommentId != model2.CommentId`)
			return
		}

		if model.MetaKey != model2.MetaKey {
			t.Errorf(`model.MetaKey != model2.MetaKey`)
			return
		}

		if model.MetaValue != model2.MetaValue {
			t.Errorf(`model.MetaValue != model2.MetaValue`)
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

		if model.CommentPostID != model2.CommentPostID {
			t.Errorf(`model.CommentPostID != model2.CommentPostID`)
			return
		}

		if model.CommentAuthor != model2.CommentAuthor {
			t.Errorf(`model.CommentAuthor != model2.CommentAuthor`)
			return
		}

		if model.CommentAuthorEmail != model2.CommentAuthorEmail {
			t.Errorf(`model.CommentAuthorEmail != model2.CommentAuthorEmail`)
			return
		}

		if model.CommentAuthorUrl != model2.CommentAuthorUrl {
			t.Errorf(`model.CommentAuthorUrl != model2.CommentAuthorUrl`)
			return
		}

		if model.CommentAuthorIP != model2.CommentAuthorIP {
			t.Errorf(`model.CommentAuthorIP != model2.CommentAuthorIP`)
			return
		}

		if model.CommentDate != model2.CommentDate {
			t.Errorf(`model.CommentDate != model2.CommentDate`)
			return
		}

		if model.CommentDateGmt != model2.CommentDateGmt {
			t.Errorf(`model.CommentDateGmt != model2.CommentDateGmt`)
			return
		}

		if model.CommentContent != model2.CommentContent {
			t.Errorf(`model.CommentContent != model2.CommentContent`)
			return
		}

		if model.CommentKarma != model2.CommentKarma {
			t.Errorf(`model.CommentKarma != model2.CommentKarma`)
			return
		}

		if model.CommentApproved != model2.CommentApproved {
			t.Errorf(`model.CommentApproved != model2.CommentApproved`)
			return
		}

		if model.CommentAgent != model2.CommentAgent {
			t.Errorf(`model.CommentAgent != model2.CommentAgent`)
			return
		}

		if model.CommentType != model2.CommentType {
			t.Errorf(`model.CommentType != model2.CommentType`)
			return
		}

		if model.CommentParent != model2.CommentParent {
			t.Errorf(`model.CommentParent != model2.CommentParent`)
			return
		}

		if model.UserId != model2.UserId {
			t.Errorf(`model.UserId != model2.UserId`)
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

		if model.LinkUrl != model2.LinkUrl {
			t.Errorf(`model.LinkUrl != model2.LinkUrl`)
			return
		}

		if model.LinkName != model2.LinkName {
			t.Errorf(`model.LinkName != model2.LinkName`)
			return
		}

		if model.LinkImage != model2.LinkImage {
			t.Errorf(`model.LinkImage != model2.LinkImage`)
			return
		}

		if model.LinkTarget != model2.LinkTarget {
			t.Errorf(`model.LinkTarget != model2.LinkTarget`)
			return
		}

		if model.LinkDescription != model2.LinkDescription {
			t.Errorf(`model.LinkDescription != model2.LinkDescription`)
			return
		}

		if model.LinkVisible != model2.LinkVisible {
			t.Errorf(`model.LinkVisible != model2.LinkVisible`)
			return
		}

		if model.LinkOwner != model2.LinkOwner {
			t.Errorf(`model.LinkOwner != model2.LinkOwner`)
			return
		}

		if model.LinkRating != model2.LinkRating {
			t.Errorf(`model.LinkRating != model2.LinkRating`)
			return
		}

		if model.LinkUpdated != model2.LinkUpdated {
			t.Errorf(`model.LinkUpdated != model2.LinkUpdated`)
			return
		}

		if model.LinkRel != model2.LinkRel {
			t.Errorf(`model.LinkRel != model2.LinkRel`)
			return
		}

		if model.LinkNotes != model2.LinkNotes {
			t.Errorf(`model.LinkNotes != model2.LinkNotes`)
			return
		}

		if model.LinkRss != model2.LinkRss {
			t.Errorf(`model.LinkRss != model2.LinkRss`)
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

		if model.OptionName != model2.OptionName {
			t.Errorf(`model.OptionName != model2.OptionName`)
			return
		}

		if model.OptionValue != model2.OptionValue {
			t.Errorf(`model.OptionValue != model2.OptionValue`)
			return
		}

		if model.Autoload != model2.Autoload {
			t.Errorf(`model.Autoload != model2.Autoload`)
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

	if o.PostId != 999 {
		t.Errorf("o.PostId test failed %+v", o)
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

		model.PostId = 999

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

		if model.PostId != model2.PostId {
			t.Errorf(`model.PostId != model2.PostId`)
			return
		}

		if model.MetaKey != model2.MetaKey {
			t.Errorf(`model.MetaKey != model2.MetaKey`)
			return
		}

		if model.MetaValue != model2.MetaValue {
			t.Errorf(`model.MetaValue != model2.MetaValue`)
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

	if o.PostAuthor != 999 {
		t.Errorf("o.PostAuthor test failed %+v", o)
		return
	}

	if o.PostDate.Year != 2016 {
		t.Errorf("year not set for %+v", o.PostDate)
		return
	}
	if o.PostDate.Year != 2016 ||
		o.PostDate.Month != 1 ||
		o.PostDate.Day != 1 ||
		o.PostDate.Hours != 10 ||
		o.PostDate.Minutes != 50 ||
		o.PostDate.Seconds != 23 ||
		o.PostDate.Offset != 5 ||
		o.PostDate.Zone != `Z` {
		t.Errorf(`fields don't match up for %+v`, o.PostDate)
	}
	r2, _ := m["post_date"].AsString()
	if o.PostDate.ToString() != r2 {
		t.Errorf(`restring of o.PostDate failed %s`, o.PostDate.ToString())
	}

	if o.PostDateGmt.Year != 2016 {
		t.Errorf("year not set for %+v", o.PostDateGmt)
		return
	}
	if o.PostDateGmt.Year != 2016 ||
		o.PostDateGmt.Month != 1 ||
		o.PostDateGmt.Day != 1 ||
		o.PostDateGmt.Hours != 10 ||
		o.PostDateGmt.Minutes != 50 ||
		o.PostDateGmt.Seconds != 23 ||
		o.PostDateGmt.Offset != 5 ||
		o.PostDateGmt.Zone != `Z` {
		t.Errorf(`fields don't match up for %+v`, o.PostDateGmt)
	}
	r3, _ := m["post_date_gmt"].AsString()
	if o.PostDateGmt.ToString() != r3 {
		t.Errorf(`restring of o.PostDateGmt failed %s`, o.PostDateGmt.ToString())
	}

	if o.PostContent != "AString" {
		t.Errorf("o.PostContent test failed %+v", o)
		return
	}

	if o.PostTitle != "AString" {
		t.Errorf("o.PostTitle test failed %+v", o)
		return
	}

	if o.PostExcerpt != "AString" {
		t.Errorf("o.PostExcerpt test failed %+v", o)
		return
	}

	if o.PostStatus != "AString" {
		t.Errorf("o.PostStatus test failed %+v", o)
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

	if o.PostPassword != "AString" {
		t.Errorf("o.PostPassword test failed %+v", o)
		return
	}

	if o.PostName != "AString" {
		t.Errorf("o.PostName test failed %+v", o)
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

	if o.PostModified.Year != 2016 {
		t.Errorf("year not set for %+v", o.PostModified)
		return
	}
	if o.PostModified.Year != 2016 ||
		o.PostModified.Month != 1 ||
		o.PostModified.Day != 1 ||
		o.PostModified.Hours != 10 ||
		o.PostModified.Minutes != 50 ||
		o.PostModified.Seconds != 23 ||
		o.PostModified.Offset != 5 ||
		o.PostModified.Zone != `Z` {
		t.Errorf(`fields don't match up for %+v`, o.PostModified)
	}
	r14, _ := m["post_modified"].AsString()
	if o.PostModified.ToString() != r14 {
		t.Errorf(`restring of o.PostModified failed %s`, o.PostModified.ToString())
	}

	if o.PostModifiedGmt.Year != 2016 {
		t.Errorf("year not set for %+v", o.PostModifiedGmt)
		return
	}
	if o.PostModifiedGmt.Year != 2016 ||
		o.PostModifiedGmt.Month != 1 ||
		o.PostModifiedGmt.Day != 1 ||
		o.PostModifiedGmt.Hours != 10 ||
		o.PostModifiedGmt.Minutes != 50 ||
		o.PostModifiedGmt.Seconds != 23 ||
		o.PostModifiedGmt.Offset != 5 ||
		o.PostModifiedGmt.Zone != `Z` {
		t.Errorf(`fields don't match up for %+v`, o.PostModifiedGmt)
	}
	r15, _ := m["post_modified_gmt"].AsString()
	if o.PostModifiedGmt.ToString() != r15 {
		t.Errorf(`restring of o.PostModifiedGmt failed %s`, o.PostModifiedGmt.ToString())
	}

	if o.PostContentFiltered != "AString" {
		t.Errorf("o.PostContentFiltered test failed %+v", o)
		return
	}

	if o.PostParent != 999 {
		t.Errorf("o.PostParent test failed %+v", o)
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

	if o.PostType != "AString" {
		t.Errorf("o.PostType test failed %+v", o)
		return
	}

	if o.PostMimeType != "AString" {
		t.Errorf("o.PostMimeType test failed %+v", o)
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

		model.PostAuthor = 999

		d1 := NewDateTime()
		d1.FromString(`999`)
		model.PostDate = d1

		d2 := NewDateTime()
		d2.FromString(`999`)
		model.PostDateGmt = d2

		model.PostContent = `the rain in spain`

		model.PostTitle = `the rain in spain`

		model.PostExcerpt = `the rain in spain`

		model.PostStatus = `the rain in spain`

		model.CommentStatus = `the rain in spain`

		model.PingStatus = `the rain in spain`

		model.PostPassword = `the rain in spain`

		model.PostName = `the rain in spain`

		model.ToPing = `the rain in spain`

		model.Pinged = `the rain in spain`

		d13 := NewDateTime()
		d13.FromString(`the rain in spain`)
		model.PostModified = d13

		d14 := NewDateTime()
		d14.FromString(`the rain in spain`)
		model.PostModifiedGmt = d14

		model.PostContentFiltered = `the rain in spain`

		model.PostParent = 999

		model.Guid = `the rain in spain`

		model.MenuOrder = 999

		model.PostType = `the rain in spain`

		model.PostMimeType = `the rain in spain`

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

		if model.PostAuthor != model2.PostAuthor {
			t.Errorf(`model.PostAuthor != model2.PostAuthor`)
			return
		}

		if model.PostDate != model2.PostDate {
			t.Errorf(`model.PostDate != model2.PostDate`)
			return
		}

		if model.PostDateGmt != model2.PostDateGmt {
			t.Errorf(`model.PostDateGmt != model2.PostDateGmt`)
			return
		}

		if model.PostContent != model2.PostContent {
			t.Errorf(`model.PostContent != model2.PostContent`)
			return
		}

		if model.PostTitle != model2.PostTitle {
			t.Errorf(`model.PostTitle != model2.PostTitle`)
			return
		}

		if model.PostExcerpt != model2.PostExcerpt {
			t.Errorf(`model.PostExcerpt != model2.PostExcerpt`)
			return
		}

		if model.PostStatus != model2.PostStatus {
			t.Errorf(`model.PostStatus != model2.PostStatus`)
			return
		}

		if model.CommentStatus != model2.CommentStatus {
			t.Errorf(`model.CommentStatus != model2.CommentStatus`)
			return
		}

		if model.PingStatus != model2.PingStatus {
			t.Errorf(`model.PingStatus != model2.PingStatus`)
			return
		}

		if model.PostPassword != model2.PostPassword {
			t.Errorf(`model.PostPassword != model2.PostPassword`)
			return
		}

		if model.PostName != model2.PostName {
			t.Errorf(`model.PostName != model2.PostName`)
			return
		}

		if model.ToPing != model2.ToPing {
			t.Errorf(`model.ToPing != model2.ToPing`)
			return
		}

		if model.Pinged != model2.Pinged {
			t.Errorf(`model.Pinged != model2.Pinged`)
			return
		}

		if model.PostModified != model2.PostModified {
			t.Errorf(`model.PostModified != model2.PostModified`)
			return
		}

		if model.PostModifiedGmt != model2.PostModifiedGmt {
			t.Errorf(`model.PostModifiedGmt != model2.PostModifiedGmt`)
			return
		}

		if model.PostContentFiltered != model2.PostContentFiltered {
			t.Errorf(`model.PostContentFiltered != model2.PostContentFiltered`)
			return
		}

		if model.PostParent != model2.PostParent {
			t.Errorf(`model.PostParent != model2.PostParent`)
			return
		}

		if model.Guid != model2.Guid {
			t.Errorf(`model.Guid != model2.Guid`)
			return
		}

		if model.MenuOrder != model2.MenuOrder {
			t.Errorf(`model.MenuOrder != model2.MenuOrder`)
			return
		}

		if model.PostType != model2.PostType {
			t.Errorf(`model.PostType != model2.PostType`)
			return
		}

		if model.PostMimeType != model2.PostMimeType {
			t.Errorf(`model.PostMimeType != model2.PostMimeType`)
			return
		}

		if model.CommentCount != model2.CommentCount {
			t.Errorf(`model.CommentCount != model2.CommentCount`)
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

		if model.TermOrder != model2.TermOrder {
			t.Errorf(`model.TermOrder != model2.TermOrder`)
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

		if model.TermId != model2.TermId {
			t.Errorf(`model.TermId != model2.TermId`)
			return
		}

		if model.Taxonomy != model2.Taxonomy {
			t.Errorf(`model.Taxonomy != model2.Taxonomy`)
			return
		}

		if model.Description != model2.Description {
			t.Errorf(`model.Description != model2.Description`)
			return
		}

		if model.Parent != model2.Parent {
			t.Errorf(`model.Parent != model2.Parent`)
			return
		}

		if model.Count != model2.Count {
			t.Errorf(`model.Count != model2.Count`)
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

		if model.Name != model2.Name {
			t.Errorf(`model.Name != model2.Name`)
			return
		}

		if model.Slug != model2.Slug {
			t.Errorf(`model.Slug != model2.Slug`)
			return
		}

		if model.TermGroup != model2.TermGroup {
			t.Errorf(`model.TermGroup != model2.TermGroup`)
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

		if model.UserId != model2.UserId {
			t.Errorf(`model.UserId != model2.UserId`)
			return
		}

		if model.MetaKey != model2.MetaKey {
			t.Errorf(`model.MetaKey != model2.MetaKey`)
			return
		}

		if model.MetaValue != model2.MetaValue {
			t.Errorf(`model.MetaValue != model2.MetaValue`)
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

		if model.UserLogin != model2.UserLogin {
			t.Errorf(`model.UserLogin != model2.UserLogin`)
			return
		}

		if model.UserPass != model2.UserPass {
			t.Errorf(`model.UserPass != model2.UserPass`)
			return
		}

		if model.UserNicename != model2.UserNicename {
			t.Errorf(`model.UserNicename != model2.UserNicename`)
			return
		}

		if model.UserEmail != model2.UserEmail {
			t.Errorf(`model.UserEmail != model2.UserEmail`)
			return
		}

		if model.UserUrl != model2.UserUrl {
			t.Errorf(`model.UserUrl != model2.UserUrl`)
			return
		}

		if model.UserRegistered != model2.UserRegistered {
			t.Errorf(`model.UserRegistered != model2.UserRegistered`)
			return
		}

		if model.UserActivationKey != model2.UserActivationKey {
			t.Errorf(`model.UserActivationKey != model2.UserActivationKey`)
			return
		}

		if model.UserStatus != model2.UserStatus {
			t.Errorf(`model.UserStatus != model2.UserStatus`)
			return
		}

		if model.DisplayName != model2.DisplayName {
			t.Errorf(`model.DisplayName != model2.DisplayName`)
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

		if model.AttrName != model2.AttrName {
			t.Errorf(`model.AttrName != model2.AttrName`)
			return
		}

		if model.AttrLabel != model2.AttrLabel {
			t.Errorf(`model.AttrLabel != model2.AttrLabel`)
			return
		}

		if model.AttrType != model2.AttrType {
			t.Errorf(`model.AttrType != model2.AttrType`)
			return
		}

		if model.AttrOrderby != model2.AttrOrderby {
			t.Errorf(`model.AttrOrderby != model2.AttrOrderby`)
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

		if model.DownloadId != model2.DownloadId {
			t.Errorf(`model.DownloadId != model2.DownloadId`)
			return
		}

		if model.ProductId != model2.ProductId {
			t.Errorf(`model.ProductId != model2.ProductId`)
			return
		}

		if model.OrderId != model2.OrderId {
			t.Errorf(`model.OrderId != model2.OrderId`)
			return
		}

		if model.OrderKey != model2.OrderKey {
			t.Errorf(`model.OrderKey != model2.OrderKey`)
			return
		}

		if model.UserEmail != model2.UserEmail {
			t.Errorf(`model.UserEmail != model2.UserEmail`)
			return
		}

		if model.UserId != model2.UserId {
			t.Errorf(`model.UserId != model2.UserId`)
			return
		}

		if model.DownloadsRemaining != model2.DownloadsRemaining {
			t.Errorf(`model.DownloadsRemaining != model2.DownloadsRemaining`)
			return
		}

		if model.AccessGranted != model2.AccessGranted {
			t.Errorf(`model.AccessGranted != model2.AccessGranted`)
			return
		}

		if model.AccessExpires != model2.AccessExpires {
			t.Errorf(`model.AccessExpires != model2.AccessExpires`)
			return
		}

		if model.DownloadCount != model2.DownloadCount {
			t.Errorf(`model.DownloadCount != model2.DownloadCount`)
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

		if model.OrderItemId != model2.OrderItemId {
			t.Errorf(`model.OrderItemId != model2.OrderItemId`)
			return
		}

		if model.MetaKey != model2.MetaKey {
			t.Errorf(`model.MetaKey != model2.MetaKey`)
			return
		}

		if model.MetaValue != model2.MetaValue {
			t.Errorf(`model.MetaValue != model2.MetaValue`)
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

		if model.OrderItemName != model2.OrderItemName {
			t.Errorf(`model.OrderItemName != model2.OrderItemName`)
			return
		}

		if model.OrderItemType != model2.OrderItemType {
			t.Errorf(`model.OrderItemType != model2.OrderItemType`)
			return
		}

		if model.OrderId != model2.OrderId {
			t.Errorf(`model.OrderId != model2.OrderId`)
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

		if model.LocationCode != model2.LocationCode {
			t.Errorf(`model.LocationCode != model2.LocationCode`)
			return
		}

		if model.TaxRateId != model2.TaxRateId {
			t.Errorf(`model.TaxRateId != model2.TaxRateId`)
			return
		}

		if model.LocationType != model2.LocationType {
			t.Errorf(`model.LocationType != model2.LocationType`)
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

		if model.TaxRateCountry != model2.TaxRateCountry {
			t.Errorf(`model.TaxRateCountry != model2.TaxRateCountry`)
			return
		}

		if model.TaxRateState != model2.TaxRateState {
			t.Errorf(`model.TaxRateState != model2.TaxRateState`)
			return
		}

		if model.TaxRate != model2.TaxRate {
			t.Errorf(`model.TaxRate != model2.TaxRate`)
			return
		}

		if model.TaxRateName != model2.TaxRateName {
			t.Errorf(`model.TaxRateName != model2.TaxRateName`)
			return
		}

		if model.TaxRatePriority != model2.TaxRatePriority {
			t.Errorf(`model.TaxRatePriority != model2.TaxRatePriority`)
			return
		}

		if model.TaxRateCompound != model2.TaxRateCompound {
			t.Errorf(`model.TaxRateCompound != model2.TaxRateCompound`)
			return
		}

		if model.TaxRateShipping != model2.TaxRateShipping {
			t.Errorf(`model.TaxRateShipping != model2.TaxRateShipping`)
			return
		}

		if model.TaxRateOrder != model2.TaxRateOrder {
			t.Errorf(`model.TaxRateOrder != model2.TaxRateOrder`)
			return
		}

		if model.TaxRateClass != model2.TaxRateClass {
			t.Errorf(`model.TaxRateClass != model2.TaxRateClass`)
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
