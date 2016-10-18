package gopress

import (
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func randomInteger() int {
	rand.Seed(time.Now().UnixNano())
	x := rand.Intn(10000) + 100
	if x == 0 {
		return randomInteger()
	}
	return x + 100
}
func randomFloat() float32 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float32() * 100
}
func randomDateTime(a Adapter) *DateTime {
	rand.Seed(time.Now().UnixNano())
	d := NewDateTime(a)
	d.Year = rand.Intn(2017)
	d.Month = rand.Intn(11)
	d.Day = rand.Intn(28)
	d.Hours = rand.Intn(23)
	d.Minutes = rand.Intn(59)
	d.Seconds = rand.Intn(56)
	if d.Year < 1000 {
		d.Year = d.Year + 1000
	}
	return d
}

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

func TestCommentMetaUpdaters(t *testing.T) {
	a, err := NewMysqlAdapterEx(`test_data/adapter.yml`)
	if err != nil {
		t.Errorf(`could not load ../gopress.db.yml %s`, err)
		return
	}
	file, err := os.OpenFile("adapter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Errorf("Failed to open log file %s", err)
	}
	a.SetLogs(file)
	model := NewCommentMeta(a)

	model.SetCommentId(int64(randomInteger()))

	model.SetMetaKey(randomString(19))

	model.SetMetaValue(randomString(25))

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
	m["comment_date"].SetInternalValue("comment_date", "2016-01-01 10:50:23")
	m["comment_date_gmt"] = a.NewDBValue()
	m["comment_date_gmt"].SetInternalValue("comment_date_gmt", "2016-01-01 10:50:23")
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
		o.CommentDate.Seconds != 23 {
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
		o.CommentDateGmt.Seconds != 23 {
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

func TestCommentUpdaters(t *testing.T) {
	a, err := NewMysqlAdapterEx(`test_data/adapter.yml`)
	if err != nil {
		t.Errorf(`could not load ../gopress.db.yml %s`, err)
		return
	}
	file, err := os.OpenFile("adapter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Errorf("Failed to open log file %s", err)
	}
	a.SetLogs(file)
	model := NewComment(a)

	model.SetCommentPostID(int64(randomInteger()))

	model.SetCommentAuthor(randomString(25))

	model.SetCommentAuthorEmail(randomString(19))

	model.SetCommentAuthorUrl(randomString(19))

	model.SetCommentAuthorIP(randomString(19))

	model.SetCommentDate(randomDateTime(a))

	model.SetCommentDateGmt(randomDateTime(a))

	model.SetCommentContent(randomString(25))

	model.SetCommentKarma(int(randomInteger()))

	model.SetCommentApproved(randomString(19))

	model.SetCommentAgent(randomString(19))

	model.SetCommentType(randomString(19))

	model.SetCommentParent(int64(randomInteger()))

	model.SetUserId(int64(randomInteger()))

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
	m["link_updated"].SetInternalValue("link_updated", "2016-01-01 10:50:23")
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
		o.LinkUpdated.Seconds != 23 {
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

func TestLinkUpdaters(t *testing.T) {
	a, err := NewMysqlAdapterEx(`test_data/adapter.yml`)
	if err != nil {
		t.Errorf(`could not load ../gopress.db.yml %s`, err)
		return
	}
	file, err := os.OpenFile("adapter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Errorf("Failed to open log file %s", err)
	}
	a.SetLogs(file)
	model := NewLink(a)

	model.SetLinkUrl(randomString(19))

	model.SetLinkName(randomString(19))

	model.SetLinkImage(randomString(19))

	model.SetLinkTarget(randomString(19))

	model.SetLinkDescription(randomString(19))

	model.SetLinkVisible(randomString(19))

	model.SetLinkOwner(int64(randomInteger()))

	model.SetLinkRating(int(randomInteger()))

	model.SetLinkUpdated(randomDateTime(a))

	model.SetLinkRel(randomString(19))

	model.SetLinkNotes(randomString(25))

	model.SetLinkRss(randomString(19))

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

func TestOptionUpdaters(t *testing.T) {
	a, err := NewMysqlAdapterEx(`test_data/adapter.yml`)
	if err != nil {
		t.Errorf(`could not load ../gopress.db.yml %s`, err)
		return
	}
	file, err := os.OpenFile("adapter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Errorf("Failed to open log file %s", err)
	}
	a.SetLogs(file)
	model := NewOption(a)

	model.SetOptionName(randomString(19))

	model.SetOptionValue(randomString(25))

	model.SetAutoload(randomString(19))

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

func TestPostMetaUpdaters(t *testing.T) {
	a, err := NewMysqlAdapterEx(`test_data/adapter.yml`)
	if err != nil {
		t.Errorf(`could not load ../gopress.db.yml %s`, err)
		return
	}
	file, err := os.OpenFile("adapter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Errorf("Failed to open log file %s", err)
	}
	a.SetLogs(file)
	model := NewPostMeta(a)

	model.SetPostId(int64(randomInteger()))

	model.SetMetaKey(randomString(19))

	model.SetMetaValue(randomString(25))

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
	m["post_date"].SetInternalValue("post_date", "2016-01-01 10:50:23")
	m["post_date_gmt"] = a.NewDBValue()
	m["post_date_gmt"].SetInternalValue("post_date_gmt", "2016-01-01 10:50:23")
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
	m["post_modified"].SetInternalValue("post_modified", "2016-01-01 10:50:23")
	m["post_modified_gmt"] = a.NewDBValue()
	m["post_modified_gmt"].SetInternalValue("post_modified_gmt", "2016-01-01 10:50:23")
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
		o.PostDate.Seconds != 23 {
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
		o.PostDateGmt.Seconds != 23 {
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
		o.PostModified.Seconds != 23 {
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
		o.PostModifiedGmt.Seconds != 23 {
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

func TestPostUpdaters(t *testing.T) {
	a, err := NewMysqlAdapterEx(`test_data/adapter.yml`)
	if err != nil {
		t.Errorf(`could not load ../gopress.db.yml %s`, err)
		return
	}
	file, err := os.OpenFile("adapter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Errorf("Failed to open log file %s", err)
	}
	a.SetLogs(file)
	model := NewPost(a)

	model.SetPostAuthor(int64(randomInteger()))

	model.SetPostDate(randomDateTime(a))

	model.SetPostDateGmt(randomDateTime(a))

	model.SetPostContent(randomString(25))

	model.SetPostTitle(randomString(25))

	model.SetPostExcerpt(randomString(25))

	model.SetPostStatus(randomString(19))

	model.SetCommentStatus(randomString(19))

	model.SetPingStatus(randomString(19))

	model.SetPostPassword(randomString(19))

	model.SetPostName(randomString(19))

	model.SetToPing(randomString(25))

	model.SetPinged(randomString(25))

	model.SetPostModified(randomDateTime(a))

	model.SetPostModifiedGmt(randomDateTime(a))

	model.SetPostContentFiltered(randomString(25))

	model.SetPostParent(int64(randomInteger()))

	model.SetGuid(randomString(19))

	model.SetMenuOrder(int(randomInteger()))

	model.SetPostType(randomString(19))

	model.SetPostMimeType(randomString(19))

	model.SetCommentCount(int64(randomInteger()))

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

func TestTermRelationshipUpdaters(t *testing.T) {
	a, err := NewMysqlAdapterEx(`test_data/adapter.yml`)
	if err != nil {
		t.Errorf(`could not load ../gopress.db.yml %s`, err)
		return
	}
	file, err := os.OpenFile("adapter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Errorf("Failed to open log file %s", err)
	}
	a.SetLogs(file)
	model := NewTermRelationship(a)

	model.SetObjectId(int64(randomInteger()))

	model.SetTermOrder(int(randomInteger()))

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

func TestTermTaxonomyUpdaters(t *testing.T) {
	a, err := NewMysqlAdapterEx(`test_data/adapter.yml`)
	if err != nil {
		t.Errorf(`could not load ../gopress.db.yml %s`, err)
		return
	}
	file, err := os.OpenFile("adapter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Errorf("Failed to open log file %s", err)
	}
	a.SetLogs(file)
	model := NewTermTaxonomy(a)

	model.SetTermId(int64(randomInteger()))

	model.SetTaxonomy(randomString(19))

	model.SetDescription(randomString(25))

	model.SetParent(int64(randomInteger()))

	model.SetCount(int64(randomInteger()))

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

func TestTermUpdaters(t *testing.T) {
	a, err := NewMysqlAdapterEx(`test_data/adapter.yml`)
	if err != nil {
		t.Errorf(`could not load ../gopress.db.yml %s`, err)
		return
	}
	file, err := os.OpenFile("adapter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Errorf("Failed to open log file %s", err)
	}
	a.SetLogs(file)
	model := NewTerm(a)

	model.SetName(randomString(19))

	model.SetSlug(randomString(19))

	model.SetTermGroup(int64(randomInteger()))

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

func TestUserMetaUpdaters(t *testing.T) {
	a, err := NewMysqlAdapterEx(`test_data/adapter.yml`)
	if err != nil {
		t.Errorf(`could not load ../gopress.db.yml %s`, err)
		return
	}
	file, err := os.OpenFile("adapter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Errorf("Failed to open log file %s", err)
	}
	a.SetLogs(file)
	model := NewUserMeta(a)

	model.SetUserId(int64(randomInteger()))

	model.SetMetaKey(randomString(19))

	model.SetMetaValue(randomString(25))

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
	dvar.SetInternalValue(`x`, `2016-01-09 23:24:50`)
	dc, err := dvar.AsDateTime()
	if err != nil {
		t.Errorf(`failed to convert datetime %+v`, dc)
	}

	if dc.Year != 2016 ||
		dc.Month != 1 ||
		dc.Day != 9 ||
		dc.Hours != 23 ||
		dc.Minutes != 24 ||
		dc.Seconds != 50 {
		t.Errorf(`fields don't match up for %+v`, dc)
	}
	r, _ := dvar.AsString()
	if dc.ToString() != r {
		t.Errorf(`restring of dvar failed %s`, dc.ToString())
	}

}
