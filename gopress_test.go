package gopress

import (
	"bufio"
	"bytes"
	"errors"
	"math/rand"
	"os"
	"regexp"
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

func TestCommentMetaCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		defer a.Close()
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
		model.CommentId = int64(randomInteger())
		model.MetaKey = randomString(19)
		model.MetaValue = randomString(25)

		err = model.Create()
		if err != nil {
			t.Errorf(`failed to create model %s`, err)
			return
		}

		model2 := NewCommentMeta(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}

		if model.CommentId != model2.CommentId {
			t.Errorf(`model.CommentId[%d] != model2.CommentId[%d]`, model.CommentId, model2.CommentId)
			return
		}

		if model.MetaKey != model2.MetaKey {
			t.Errorf(`model.MetaKey[%s] != model2.MetaKey[%s]`, model.MetaKey, model2.MetaKey)
			return
		}

		if model.MetaValue != model2.MetaValue {
			t.Errorf(`model.MetaValue[%s] != model2.MetaValue[%s]`, model.MetaValue, model2.MetaValue)
			return
		}
		model2.SetCommentId(int64(randomInteger()))
		model2.SetMetaKey(randomString(19))
		model2.SetMetaValue(randomString(25))

		err = model2.Save()
		if err != nil {
			t.Errorf(`failed to save model2 %s`, err)
		}

		if model.CommentId == model2.CommentId {
			t.Errorf(`model.CommentId[%d] != model2.CommentId[%d]`, model.CommentId, model2.CommentId)
			return
		}

		if model.MetaKey == model2.MetaKey {
			t.Errorf(`model.MetaKey[%s] != model2.MetaKey[%s]`, model.MetaKey, model2.MetaKey)
			return
		}

		if model.MetaValue == model2.MetaValue {
			t.Errorf(`model.MetaValue[%s] != model2.MetaValue[%s]`, model.MetaValue, model2.MetaValue)
			return
		}

		res9, err := model.FindByCommentId(model2.GetCommentId())
		if err != nil {
			t.Errorf(`failed model.FindByCommentId(model2.GetCommentId())`)
		}
		if len(res9) == 0 {
			t.Errorf(`failed to find any CommentMeta`)
		}

		res10, err := model.FindByMetaKey(model2.GetMetaKey())
		if err != nil {
			t.Errorf(`failed model.FindByMetaKey(model2.GetMetaKey())`)
		}
		if len(res10) == 0 {
			t.Errorf(`failed to find any CommentMeta`)
		}

		res11, err := model.FindByMetaValue(model2.GetMetaValue())
		if err != nil {
			t.Errorf(`failed model.FindByMetaValue(model2.GetMetaValue())`)
		}
		if len(res11) == 0 {
			t.Errorf(`failed to find any CommentMeta`)
		}
	} // end of if fileExists
}

func TestCommentMetaUpdaters(t *testing.T) {
	if fileExists(`../gopress.db.yml`) == false {
		return
	}
	a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
	defer a.Close()
	if err != nil {
		t.Errorf(`could not load ../gopress.db.yml %s`, err)
		return
	}
	file, err := os.OpenFile("adapter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Errorf("Failed to open log file %s", err)
		return
	}
	a.SetLogs(file)
	model := NewCommentMeta(a)

	model.SetCommentId(int64(randomInteger()))
	if model.GetCommentId() != model.CommentId {
		t.Errorf(`CommentMeta.GetCommentId() != CommentMeta.CommentId`)
	}
	if model.IsCommentIdDirty != true {
		t.Errorf(`CommentMeta.IsCommentIdDirty != true`)
		return
	}

	u0 := int64(randomInteger())
	_, err = model.UpdateCommentId(u0)
	if err != nil {
		t.Errorf(`failed UpdateCommentId(u0) %s`, err)
		return
	}

	if model.GetCommentId() != u0 {
		t.Errorf(`CommentMeta.GetCommentId() != u0 after UpdateCommentId`)
		return
	}
	model.Reload()
	if model.GetCommentId() != u0 {
		t.Errorf(`CommentMeta.GetCommentId() != u0 after Reload`)
		return
	}

	model.SetMetaKey(randomString(19))
	if model.GetMetaKey() != model.MetaKey {
		t.Errorf(`CommentMeta.GetMetaKey() != CommentMeta.MetaKey`)
	}
	if model.IsMetaKeyDirty != true {
		t.Errorf(`CommentMeta.IsMetaKeyDirty != true`)
		return
	}

	u1 := randomString(19)
	_, err = model.UpdateMetaKey(u1)
	if err != nil {
		t.Errorf(`failed UpdateMetaKey(u1) %s`, err)
		return
	}

	if model.GetMetaKey() != u1 {
		t.Errorf(`CommentMeta.GetMetaKey() != u1 after UpdateMetaKey`)
		return
	}
	model.Reload()
	if model.GetMetaKey() != u1 {
		t.Errorf(`CommentMeta.GetMetaKey() != u1 after Reload`)
		return
	}

	model.SetMetaValue(randomString(25))
	if model.GetMetaValue() != model.MetaValue {
		t.Errorf(`CommentMeta.GetMetaValue() != CommentMeta.MetaValue`)
	}
	if model.IsMetaValueDirty != true {
		t.Errorf(`CommentMeta.IsMetaValueDirty != true`)
		return
	}

	u2 := randomString(25)
	_, err = model.UpdateMetaValue(u2)
	if err != nil {
		t.Errorf(`failed UpdateMetaValue(u2) %s`, err)
		return
	}

	if model.GetMetaValue() != u2 {
		t.Errorf(`CommentMeta.GetMetaValue() != u2 after UpdateMetaValue`)
		return
	}
	model.Reload()
	if model.GetMetaValue() != u2 {
		t.Errorf(`CommentMeta.GetMetaValue() != u2 after Reload`)
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

func TestCommentCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		defer a.Close()
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
		model.CommentPostID = int64(randomInteger())
		model.CommentAuthor = randomString(25)
		model.CommentAuthorEmail = randomString(19)
		model.CommentAuthorUrl = randomString(19)
		model.CommentAuthorIP = randomString(19)
		model.CommentDate = randomDateTime(a)
		model.CommentDateGmt = randomDateTime(a)
		model.CommentContent = randomString(25)
		model.CommentKarma = int(randomInteger())
		model.CommentApproved = randomString(19)
		model.CommentAgent = randomString(19)
		model.CommentType = randomString(19)
		model.CommentParent = int64(randomInteger())
		model.UserId = int64(randomInteger())

		err = model.Create()
		if err != nil {
			t.Errorf(`failed to create model %s`, err)
			return
		}

		model2 := NewComment(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}

		if model.CommentPostID != model2.CommentPostID {
			t.Errorf(`model.CommentPostID[%d] != model2.CommentPostID[%d]`, model.CommentPostID, model2.CommentPostID)
			return
		}

		if model.CommentAuthor != model2.CommentAuthor {
			t.Errorf(`model.CommentAuthor[%s] != model2.CommentAuthor[%s]`, model.CommentAuthor, model2.CommentAuthor)
			return
		}

		if model.CommentAuthorEmail != model2.CommentAuthorEmail {
			t.Errorf(`model.CommentAuthorEmail[%s] != model2.CommentAuthorEmail[%s]`, model.CommentAuthorEmail, model2.CommentAuthorEmail)
			return
		}

		if model.CommentAuthorUrl != model2.CommentAuthorUrl {
			t.Errorf(`model.CommentAuthorUrl[%s] != model2.CommentAuthorUrl[%s]`, model.CommentAuthorUrl, model2.CommentAuthorUrl)
			return
		}

		if model.CommentAuthorIP != model2.CommentAuthorIP {
			t.Errorf(`model.CommentAuthorIP[%s] != model2.CommentAuthorIP[%s]`, model.CommentAuthorIP, model2.CommentAuthorIP)
			return
		}

		if model.CommentDate.Year != model2.CommentDate.Year ||
			model.CommentDate.Month != model2.CommentDate.Month ||
			model.CommentDate.Day != model2.CommentDate.Day ||
			model.CommentDate.Hours != model2.CommentDate.Hours ||
			model.CommentDate.Minutes != model2.CommentDate.Minutes ||
			model.CommentDate.Seconds != model2.CommentDate.Seconds {
			t.Errorf(`model.CommentDate != model2.CommentDate %+v --- %+v`, model.CommentDate, model2.CommentDate)
			return
		}

		if model.CommentDateGmt.Year != model2.CommentDateGmt.Year ||
			model.CommentDateGmt.Month != model2.CommentDateGmt.Month ||
			model.CommentDateGmt.Day != model2.CommentDateGmt.Day ||
			model.CommentDateGmt.Hours != model2.CommentDateGmt.Hours ||
			model.CommentDateGmt.Minutes != model2.CommentDateGmt.Minutes ||
			model.CommentDateGmt.Seconds != model2.CommentDateGmt.Seconds {
			t.Errorf(`model.CommentDateGmt != model2.CommentDateGmt %+v --- %+v`, model.CommentDateGmt, model2.CommentDateGmt)
			return
		}

		if model.CommentContent != model2.CommentContent {
			t.Errorf(`model.CommentContent[%s] != model2.CommentContent[%s]`, model.CommentContent, model2.CommentContent)
			return
		}

		if model.CommentKarma != model2.CommentKarma {
			t.Errorf(`model.CommentKarma[%d] != model2.CommentKarma[%d]`, model.CommentKarma, model2.CommentKarma)
			return
		}

		if model.CommentApproved != model2.CommentApproved {
			t.Errorf(`model.CommentApproved[%s] != model2.CommentApproved[%s]`, model.CommentApproved, model2.CommentApproved)
			return
		}

		if model.CommentAgent != model2.CommentAgent {
			t.Errorf(`model.CommentAgent[%s] != model2.CommentAgent[%s]`, model.CommentAgent, model2.CommentAgent)
			return
		}

		if model.CommentType != model2.CommentType {
			t.Errorf(`model.CommentType[%s] != model2.CommentType[%s]`, model.CommentType, model2.CommentType)
			return
		}

		if model.CommentParent != model2.CommentParent {
			t.Errorf(`model.CommentParent[%d] != model2.CommentParent[%d]`, model.CommentParent, model2.CommentParent)
			return
		}

		if model.UserId != model2.UserId {
			t.Errorf(`model.UserId[%d] != model2.UserId[%d]`, model.UserId, model2.UserId)
			return
		}
		model2.SetCommentPostID(int64(randomInteger()))
		model2.SetCommentAuthor(randomString(25))
		model2.SetCommentAuthorEmail(randomString(19))
		model2.SetCommentAuthorUrl(randomString(19))
		model2.SetCommentAuthorIP(randomString(19))
		model2.SetCommentDate(randomDateTime(a))
		model2.SetCommentDateGmt(randomDateTime(a))
		model2.SetCommentContent(randomString(25))
		model2.SetCommentKarma(int(randomInteger()))
		model2.SetCommentApproved(randomString(19))
		model2.SetCommentAgent(randomString(19))
		model2.SetCommentType(randomString(19))
		model2.SetCommentParent(int64(randomInteger()))
		model2.SetUserId(int64(randomInteger()))

		err = model2.Save()
		if err != nil {
			t.Errorf(`failed to save model2 %s`, err)
		}

		if model.CommentPostID == model2.CommentPostID {
			t.Errorf(`model.CommentPostID[%d] != model2.CommentPostID[%d]`, model.CommentPostID, model2.CommentPostID)
			return
		}

		if model.CommentAuthor == model2.CommentAuthor {
			t.Errorf(`model.CommentAuthor[%s] != model2.CommentAuthor[%s]`, model.CommentAuthor, model2.CommentAuthor)
			return
		}

		if model.CommentAuthorEmail == model2.CommentAuthorEmail {
			t.Errorf(`model.CommentAuthorEmail[%s] != model2.CommentAuthorEmail[%s]`, model.CommentAuthorEmail, model2.CommentAuthorEmail)
			return
		}

		if model.CommentAuthorUrl == model2.CommentAuthorUrl {
			t.Errorf(`model.CommentAuthorUrl[%s] != model2.CommentAuthorUrl[%s]`, model.CommentAuthorUrl, model2.CommentAuthorUrl)
			return
		}

		if model.CommentAuthorIP == model2.CommentAuthorIP {
			t.Errorf(`model.CommentAuthorIP[%s] != model2.CommentAuthorIP[%s]`, model.CommentAuthorIP, model2.CommentAuthorIP)
			return
		}

		if model.CommentDate.Year == model2.CommentDate.Year {
			t.Errorf(`model.CommentDate.Year == model2.CommentDate but should not!`, model.CommentDate, model2.CommentDate)
			return
		}

		if model.CommentDateGmt.Year == model2.CommentDateGmt.Year {
			t.Errorf(`model.CommentDateGmt.Year == model2.CommentDateGmt but should not!`, model.CommentDateGmt, model2.CommentDateGmt)
			return
		}

		if model.CommentContent == model2.CommentContent {
			t.Errorf(`model.CommentContent[%s] != model2.CommentContent[%s]`, model.CommentContent, model2.CommentContent)
			return
		}

		if model.CommentKarma == model2.CommentKarma {
			t.Errorf(`model.CommentKarma[%d] != model2.CommentKarma[%d]`, model.CommentKarma, model2.CommentKarma)
			return
		}

		if model.CommentApproved == model2.CommentApproved {
			t.Errorf(`model.CommentApproved[%s] != model2.CommentApproved[%s]`, model.CommentApproved, model2.CommentApproved)
			return
		}

		if model.CommentAgent == model2.CommentAgent {
			t.Errorf(`model.CommentAgent[%s] != model2.CommentAgent[%s]`, model.CommentAgent, model2.CommentAgent)
			return
		}

		if model.CommentType == model2.CommentType {
			t.Errorf(`model.CommentType[%s] != model2.CommentType[%s]`, model.CommentType, model2.CommentType)
			return
		}

		if model.CommentParent == model2.CommentParent {
			t.Errorf(`model.CommentParent[%d] != model2.CommentParent[%d]`, model.CommentParent, model2.CommentParent)
			return
		}

		if model.UserId == model2.UserId {
			t.Errorf(`model.UserId[%d] != model2.UserId[%d]`, model.UserId, model2.UserId)
			return
		}

		res38, err := model.FindByCommentPostID(model2.GetCommentPostID())
		if err != nil {
			t.Errorf(`failed model.FindByCommentPostID(model2.GetCommentPostID())`)
		}
		if len(res38) == 0 {
			t.Errorf(`failed to find any Comment`)
		}

		res39, err := model.FindByCommentAuthor(model2.GetCommentAuthor())
		if err != nil {
			t.Errorf(`failed model.FindByCommentAuthor(model2.GetCommentAuthor())`)
		}
		if len(res39) == 0 {
			t.Errorf(`failed to find any Comment`)
		}

		res40, err := model.FindByCommentAuthorEmail(model2.GetCommentAuthorEmail())
		if err != nil {
			t.Errorf(`failed model.FindByCommentAuthorEmail(model2.GetCommentAuthorEmail())`)
		}
		if len(res40) == 0 {
			t.Errorf(`failed to find any Comment`)
		}

		res41, err := model.FindByCommentAuthorUrl(model2.GetCommentAuthorUrl())
		if err != nil {
			t.Errorf(`failed model.FindByCommentAuthorUrl(model2.GetCommentAuthorUrl())`)
		}
		if len(res41) == 0 {
			t.Errorf(`failed to find any Comment`)
		}

		res42, err := model.FindByCommentAuthorIP(model2.GetCommentAuthorIP())
		if err != nil {
			t.Errorf(`failed model.FindByCommentAuthorIP(model2.GetCommentAuthorIP())`)
		}
		if len(res42) == 0 {
			t.Errorf(`failed to find any Comment`)
		}

		res43, err := model.FindByCommentDate(model2.GetCommentDate())
		if err != nil {
			t.Errorf(`failed model.FindByCommentDate(model2.GetCommentDate())`)
		}
		if len(res43) == 0 {
			t.Errorf(`failed to find any Comment`)
		}

		res44, err := model.FindByCommentDateGmt(model2.GetCommentDateGmt())
		if err != nil {
			t.Errorf(`failed model.FindByCommentDateGmt(model2.GetCommentDateGmt())`)
		}
		if len(res44) == 0 {
			t.Errorf(`failed to find any Comment`)
		}

		res45, err := model.FindByCommentContent(model2.GetCommentContent())
		if err != nil {
			t.Errorf(`failed model.FindByCommentContent(model2.GetCommentContent())`)
		}
		if len(res45) == 0 {
			t.Errorf(`failed to find any Comment`)
		}

		res46, err := model.FindByCommentKarma(model2.GetCommentKarma())
		if err != nil {
			t.Errorf(`failed model.FindByCommentKarma(model2.GetCommentKarma())`)
		}
		if len(res46) == 0 {
			t.Errorf(`failed to find any Comment`)
		}

		res47, err := model.FindByCommentApproved(model2.GetCommentApproved())
		if err != nil {
			t.Errorf(`failed model.FindByCommentApproved(model2.GetCommentApproved())`)
		}
		if len(res47) == 0 {
			t.Errorf(`failed to find any Comment`)
		}

		res48, err := model.FindByCommentAgent(model2.GetCommentAgent())
		if err != nil {
			t.Errorf(`failed model.FindByCommentAgent(model2.GetCommentAgent())`)
		}
		if len(res48) == 0 {
			t.Errorf(`failed to find any Comment`)
		}

		res49, err := model.FindByCommentType(model2.GetCommentType())
		if err != nil {
			t.Errorf(`failed model.FindByCommentType(model2.GetCommentType())`)
		}
		if len(res49) == 0 {
			t.Errorf(`failed to find any Comment`)
		}

		res50, err := model.FindByCommentParent(model2.GetCommentParent())
		if err != nil {
			t.Errorf(`failed model.FindByCommentParent(model2.GetCommentParent())`)
		}
		if len(res50) == 0 {
			t.Errorf(`failed to find any Comment`)
		}

		res51, err := model.FindByUserId(model2.GetUserId())
		if err != nil {
			t.Errorf(`failed model.FindByUserId(model2.GetUserId())`)
		}
		if len(res51) == 0 {
			t.Errorf(`failed to find any Comment`)
		}
	} // end of if fileExists
}

func TestCommentUpdaters(t *testing.T) {
	if fileExists(`../gopress.db.yml`) == false {
		return
	}
	a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
	defer a.Close()
	if err != nil {
		t.Errorf(`could not load ../gopress.db.yml %s`, err)
		return
	}
	file, err := os.OpenFile("adapter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Errorf("Failed to open log file %s", err)
		return
	}
	a.SetLogs(file)
	model := NewComment(a)

	model.SetCommentPostID(int64(randomInteger()))
	if model.GetCommentPostID() != model.CommentPostID {
		t.Errorf(`Comment.GetCommentPostID() != Comment.CommentPostID`)
	}
	if model.IsCommentPostIDDirty != true {
		t.Errorf(`Comment.IsCommentPostIDDirty != true`)
		return
	}

	u0 := int64(randomInteger())
	_, err = model.UpdateCommentPostID(u0)
	if err != nil {
		t.Errorf(`failed UpdateCommentPostID(u0) %s`, err)
		return
	}

	if model.GetCommentPostID() != u0 {
		t.Errorf(`Comment.GetCommentPostID() != u0 after UpdateCommentPostID`)
		return
	}
	model.Reload()
	if model.GetCommentPostID() != u0 {
		t.Errorf(`Comment.GetCommentPostID() != u0 after Reload`)
		return
	}

	model.SetCommentAuthor(randomString(25))
	if model.GetCommentAuthor() != model.CommentAuthor {
		t.Errorf(`Comment.GetCommentAuthor() != Comment.CommentAuthor`)
	}
	if model.IsCommentAuthorDirty != true {
		t.Errorf(`Comment.IsCommentAuthorDirty != true`)
		return
	}

	u1 := randomString(25)
	_, err = model.UpdateCommentAuthor(u1)
	if err != nil {
		t.Errorf(`failed UpdateCommentAuthor(u1) %s`, err)
		return
	}

	if model.GetCommentAuthor() != u1 {
		t.Errorf(`Comment.GetCommentAuthor() != u1 after UpdateCommentAuthor`)
		return
	}
	model.Reload()
	if model.GetCommentAuthor() != u1 {
		t.Errorf(`Comment.GetCommentAuthor() != u1 after Reload`)
		return
	}

	model.SetCommentAuthorEmail(randomString(19))
	if model.GetCommentAuthorEmail() != model.CommentAuthorEmail {
		t.Errorf(`Comment.GetCommentAuthorEmail() != Comment.CommentAuthorEmail`)
	}
	if model.IsCommentAuthorEmailDirty != true {
		t.Errorf(`Comment.IsCommentAuthorEmailDirty != true`)
		return
	}

	u2 := randomString(19)
	_, err = model.UpdateCommentAuthorEmail(u2)
	if err != nil {
		t.Errorf(`failed UpdateCommentAuthorEmail(u2) %s`, err)
		return
	}

	if model.GetCommentAuthorEmail() != u2 {
		t.Errorf(`Comment.GetCommentAuthorEmail() != u2 after UpdateCommentAuthorEmail`)
		return
	}
	model.Reload()
	if model.GetCommentAuthorEmail() != u2 {
		t.Errorf(`Comment.GetCommentAuthorEmail() != u2 after Reload`)
		return
	}

	model.SetCommentAuthorUrl(randomString(19))
	if model.GetCommentAuthorUrl() != model.CommentAuthorUrl {
		t.Errorf(`Comment.GetCommentAuthorUrl() != Comment.CommentAuthorUrl`)
	}
	if model.IsCommentAuthorUrlDirty != true {
		t.Errorf(`Comment.IsCommentAuthorUrlDirty != true`)
		return
	}

	u3 := randomString(19)
	_, err = model.UpdateCommentAuthorUrl(u3)
	if err != nil {
		t.Errorf(`failed UpdateCommentAuthorUrl(u3) %s`, err)
		return
	}

	if model.GetCommentAuthorUrl() != u3 {
		t.Errorf(`Comment.GetCommentAuthorUrl() != u3 after UpdateCommentAuthorUrl`)
		return
	}
	model.Reload()
	if model.GetCommentAuthorUrl() != u3 {
		t.Errorf(`Comment.GetCommentAuthorUrl() != u3 after Reload`)
		return
	}

	model.SetCommentAuthorIP(randomString(19))
	if model.GetCommentAuthorIP() != model.CommentAuthorIP {
		t.Errorf(`Comment.GetCommentAuthorIP() != Comment.CommentAuthorIP`)
	}
	if model.IsCommentAuthorIPDirty != true {
		t.Errorf(`Comment.IsCommentAuthorIPDirty != true`)
		return
	}

	u4 := randomString(19)
	_, err = model.UpdateCommentAuthorIP(u4)
	if err != nil {
		t.Errorf(`failed UpdateCommentAuthorIP(u4) %s`, err)
		return
	}

	if model.GetCommentAuthorIP() != u4 {
		t.Errorf(`Comment.GetCommentAuthorIP() != u4 after UpdateCommentAuthorIP`)
		return
	}
	model.Reload()
	if model.GetCommentAuthorIP() != u4 {
		t.Errorf(`Comment.GetCommentAuthorIP() != u4 after Reload`)
		return
	}

	model.SetCommentDate(randomDateTime(a))
	if model.GetCommentDate() != model.CommentDate {
		t.Errorf(`Comment.GetCommentDate() != Comment.CommentDate`)
	}
	if model.IsCommentDateDirty != true {
		t.Errorf(`Comment.IsCommentDateDirty != true`)
		return
	}

	u5 := randomDateTime(a)
	_, err = model.UpdateCommentDate(u5)
	if err != nil {
		t.Errorf(`failed UpdateCommentDate(u5) %s`, err)
		return
	}

	if model.GetCommentDate() != u5 {
		t.Errorf(`Comment.GetCommentDate() != u5 after UpdateCommentDate`)
		return
	}
	model.Reload()
	if model.GetCommentDate() != u5 {
		t.Errorf(`Comment.GetCommentDate() != u5 after Reload`)
		return
	}

	model.SetCommentDateGmt(randomDateTime(a))
	if model.GetCommentDateGmt() != model.CommentDateGmt {
		t.Errorf(`Comment.GetCommentDateGmt() != Comment.CommentDateGmt`)
	}
	if model.IsCommentDateGmtDirty != true {
		t.Errorf(`Comment.IsCommentDateGmtDirty != true`)
		return
	}

	u6 := randomDateTime(a)
	_, err = model.UpdateCommentDateGmt(u6)
	if err != nil {
		t.Errorf(`failed UpdateCommentDateGmt(u6) %s`, err)
		return
	}

	if model.GetCommentDateGmt() != u6 {
		t.Errorf(`Comment.GetCommentDateGmt() != u6 after UpdateCommentDateGmt`)
		return
	}
	model.Reload()
	if model.GetCommentDateGmt() != u6 {
		t.Errorf(`Comment.GetCommentDateGmt() != u6 after Reload`)
		return
	}

	model.SetCommentContent(randomString(25))
	if model.GetCommentContent() != model.CommentContent {
		t.Errorf(`Comment.GetCommentContent() != Comment.CommentContent`)
	}
	if model.IsCommentContentDirty != true {
		t.Errorf(`Comment.IsCommentContentDirty != true`)
		return
	}

	u7 := randomString(25)
	_, err = model.UpdateCommentContent(u7)
	if err != nil {
		t.Errorf(`failed UpdateCommentContent(u7) %s`, err)
		return
	}

	if model.GetCommentContent() != u7 {
		t.Errorf(`Comment.GetCommentContent() != u7 after UpdateCommentContent`)
		return
	}
	model.Reload()
	if model.GetCommentContent() != u7 {
		t.Errorf(`Comment.GetCommentContent() != u7 after Reload`)
		return
	}

	model.SetCommentKarma(int(randomInteger()))
	if model.GetCommentKarma() != model.CommentKarma {
		t.Errorf(`Comment.GetCommentKarma() != Comment.CommentKarma`)
	}
	if model.IsCommentKarmaDirty != true {
		t.Errorf(`Comment.IsCommentKarmaDirty != true`)
		return
	}

	u8 := int(randomInteger())
	_, err = model.UpdateCommentKarma(u8)
	if err != nil {
		t.Errorf(`failed UpdateCommentKarma(u8) %s`, err)
		return
	}

	if model.GetCommentKarma() != u8 {
		t.Errorf(`Comment.GetCommentKarma() != u8 after UpdateCommentKarma`)
		return
	}
	model.Reload()
	if model.GetCommentKarma() != u8 {
		t.Errorf(`Comment.GetCommentKarma() != u8 after Reload`)
		return
	}

	model.SetCommentApproved(randomString(19))
	if model.GetCommentApproved() != model.CommentApproved {
		t.Errorf(`Comment.GetCommentApproved() != Comment.CommentApproved`)
	}
	if model.IsCommentApprovedDirty != true {
		t.Errorf(`Comment.IsCommentApprovedDirty != true`)
		return
	}

	u9 := randomString(19)
	_, err = model.UpdateCommentApproved(u9)
	if err != nil {
		t.Errorf(`failed UpdateCommentApproved(u9) %s`, err)
		return
	}

	if model.GetCommentApproved() != u9 {
		t.Errorf(`Comment.GetCommentApproved() != u9 after UpdateCommentApproved`)
		return
	}
	model.Reload()
	if model.GetCommentApproved() != u9 {
		t.Errorf(`Comment.GetCommentApproved() != u9 after Reload`)
		return
	}

	model.SetCommentAgent(randomString(19))
	if model.GetCommentAgent() != model.CommentAgent {
		t.Errorf(`Comment.GetCommentAgent() != Comment.CommentAgent`)
	}
	if model.IsCommentAgentDirty != true {
		t.Errorf(`Comment.IsCommentAgentDirty != true`)
		return
	}

	u10 := randomString(19)
	_, err = model.UpdateCommentAgent(u10)
	if err != nil {
		t.Errorf(`failed UpdateCommentAgent(u10) %s`, err)
		return
	}

	if model.GetCommentAgent() != u10 {
		t.Errorf(`Comment.GetCommentAgent() != u10 after UpdateCommentAgent`)
		return
	}
	model.Reload()
	if model.GetCommentAgent() != u10 {
		t.Errorf(`Comment.GetCommentAgent() != u10 after Reload`)
		return
	}

	model.SetCommentType(randomString(19))
	if model.GetCommentType() != model.CommentType {
		t.Errorf(`Comment.GetCommentType() != Comment.CommentType`)
	}
	if model.IsCommentTypeDirty != true {
		t.Errorf(`Comment.IsCommentTypeDirty != true`)
		return
	}

	u11 := randomString(19)
	_, err = model.UpdateCommentType(u11)
	if err != nil {
		t.Errorf(`failed UpdateCommentType(u11) %s`, err)
		return
	}

	if model.GetCommentType() != u11 {
		t.Errorf(`Comment.GetCommentType() != u11 after UpdateCommentType`)
		return
	}
	model.Reload()
	if model.GetCommentType() != u11 {
		t.Errorf(`Comment.GetCommentType() != u11 after Reload`)
		return
	}

	model.SetCommentParent(int64(randomInteger()))
	if model.GetCommentParent() != model.CommentParent {
		t.Errorf(`Comment.GetCommentParent() != Comment.CommentParent`)
	}
	if model.IsCommentParentDirty != true {
		t.Errorf(`Comment.IsCommentParentDirty != true`)
		return
	}

	u12 := int64(randomInteger())
	_, err = model.UpdateCommentParent(u12)
	if err != nil {
		t.Errorf(`failed UpdateCommentParent(u12) %s`, err)
		return
	}

	if model.GetCommentParent() != u12 {
		t.Errorf(`Comment.GetCommentParent() != u12 after UpdateCommentParent`)
		return
	}
	model.Reload()
	if model.GetCommentParent() != u12 {
		t.Errorf(`Comment.GetCommentParent() != u12 after Reload`)
		return
	}

	model.SetUserId(int64(randomInteger()))
	if model.GetUserId() != model.UserId {
		t.Errorf(`Comment.GetUserId() != Comment.UserId`)
	}
	if model.IsUserIdDirty != true {
		t.Errorf(`Comment.IsUserIdDirty != true`)
		return
	}

	u13 := int64(randomInteger())
	_, err = model.UpdateUserId(u13)
	if err != nil {
		t.Errorf(`failed UpdateUserId(u13) %s`, err)
		return
	}

	if model.GetUserId() != u13 {
		t.Errorf(`Comment.GetUserId() != u13 after UpdateUserId`)
		return
	}
	model.Reload()
	if model.GetUserId() != u13 {
		t.Errorf(`Comment.GetUserId() != u13 after Reload`)
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

func TestLinkCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		defer a.Close()
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
		model.LinkUrl = randomString(19)
		model.LinkName = randomString(19)
		model.LinkImage = randomString(19)
		model.LinkTarget = randomString(19)
		model.LinkDescription = randomString(19)
		model.LinkVisible = randomString(19)
		model.LinkOwner = int64(randomInteger())
		model.LinkRating = int(randomInteger())
		model.LinkUpdated = randomDateTime(a)
		model.LinkRel = randomString(19)
		model.LinkNotes = randomString(25)
		model.LinkRss = randomString(19)

		err = model.Create()
		if err != nil {
			t.Errorf(`failed to create model %s`, err)
			return
		}

		model2 := NewLink(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}

		if model.LinkUrl != model2.LinkUrl {
			t.Errorf(`model.LinkUrl[%s] != model2.LinkUrl[%s]`, model.LinkUrl, model2.LinkUrl)
			return
		}

		if model.LinkName != model2.LinkName {
			t.Errorf(`model.LinkName[%s] != model2.LinkName[%s]`, model.LinkName, model2.LinkName)
			return
		}

		if model.LinkImage != model2.LinkImage {
			t.Errorf(`model.LinkImage[%s] != model2.LinkImage[%s]`, model.LinkImage, model2.LinkImage)
			return
		}

		if model.LinkTarget != model2.LinkTarget {
			t.Errorf(`model.LinkTarget[%s] != model2.LinkTarget[%s]`, model.LinkTarget, model2.LinkTarget)
			return
		}

		if model.LinkDescription != model2.LinkDescription {
			t.Errorf(`model.LinkDescription[%s] != model2.LinkDescription[%s]`, model.LinkDescription, model2.LinkDescription)
			return
		}

		if model.LinkVisible != model2.LinkVisible {
			t.Errorf(`model.LinkVisible[%s] != model2.LinkVisible[%s]`, model.LinkVisible, model2.LinkVisible)
			return
		}

		if model.LinkOwner != model2.LinkOwner {
			t.Errorf(`model.LinkOwner[%d] != model2.LinkOwner[%d]`, model.LinkOwner, model2.LinkOwner)
			return
		}

		if model.LinkRating != model2.LinkRating {
			t.Errorf(`model.LinkRating[%d] != model2.LinkRating[%d]`, model.LinkRating, model2.LinkRating)
			return
		}

		if model.LinkUpdated.Year != model2.LinkUpdated.Year ||
			model.LinkUpdated.Month != model2.LinkUpdated.Month ||
			model.LinkUpdated.Day != model2.LinkUpdated.Day ||
			model.LinkUpdated.Hours != model2.LinkUpdated.Hours ||
			model.LinkUpdated.Minutes != model2.LinkUpdated.Minutes ||
			model.LinkUpdated.Seconds != model2.LinkUpdated.Seconds {
			t.Errorf(`model.LinkUpdated != model2.LinkUpdated %+v --- %+v`, model.LinkUpdated, model2.LinkUpdated)
			return
		}

		if model.LinkRel != model2.LinkRel {
			t.Errorf(`model.LinkRel[%s] != model2.LinkRel[%s]`, model.LinkRel, model2.LinkRel)
			return
		}

		if model.LinkNotes != model2.LinkNotes {
			t.Errorf(`model.LinkNotes[%s] != model2.LinkNotes[%s]`, model.LinkNotes, model2.LinkNotes)
			return
		}

		if model.LinkRss != model2.LinkRss {
			t.Errorf(`model.LinkRss[%s] != model2.LinkRss[%s]`, model.LinkRss, model2.LinkRss)
			return
		}
		model2.SetLinkUrl(randomString(19))
		model2.SetLinkName(randomString(19))
		model2.SetLinkImage(randomString(19))
		model2.SetLinkTarget(randomString(19))
		model2.SetLinkDescription(randomString(19))
		model2.SetLinkVisible(randomString(19))
		model2.SetLinkOwner(int64(randomInteger()))
		model2.SetLinkRating(int(randomInteger()))
		model2.SetLinkUpdated(randomDateTime(a))
		model2.SetLinkRel(randomString(19))
		model2.SetLinkNotes(randomString(25))
		model2.SetLinkRss(randomString(19))

		err = model2.Save()
		if err != nil {
			t.Errorf(`failed to save model2 %s`, err)
		}

		if model.LinkUrl == model2.LinkUrl {
			t.Errorf(`model.LinkUrl[%s] != model2.LinkUrl[%s]`, model.LinkUrl, model2.LinkUrl)
			return
		}

		if model.LinkName == model2.LinkName {
			t.Errorf(`model.LinkName[%s] != model2.LinkName[%s]`, model.LinkName, model2.LinkName)
			return
		}

		if model.LinkImage == model2.LinkImage {
			t.Errorf(`model.LinkImage[%s] != model2.LinkImage[%s]`, model.LinkImage, model2.LinkImage)
			return
		}

		if model.LinkTarget == model2.LinkTarget {
			t.Errorf(`model.LinkTarget[%s] != model2.LinkTarget[%s]`, model.LinkTarget, model2.LinkTarget)
			return
		}

		if model.LinkDescription == model2.LinkDescription {
			t.Errorf(`model.LinkDescription[%s] != model2.LinkDescription[%s]`, model.LinkDescription, model2.LinkDescription)
			return
		}

		if model.LinkVisible == model2.LinkVisible {
			t.Errorf(`model.LinkVisible[%s] != model2.LinkVisible[%s]`, model.LinkVisible, model2.LinkVisible)
			return
		}

		if model.LinkOwner == model2.LinkOwner {
			t.Errorf(`model.LinkOwner[%d] != model2.LinkOwner[%d]`, model.LinkOwner, model2.LinkOwner)
			return
		}

		if model.LinkRating == model2.LinkRating {
			t.Errorf(`model.LinkRating[%d] != model2.LinkRating[%d]`, model.LinkRating, model2.LinkRating)
			return
		}

		if model.LinkUpdated.Year == model2.LinkUpdated.Year {
			t.Errorf(`model.LinkUpdated.Year == model2.LinkUpdated but should not!`, model.LinkUpdated, model2.LinkUpdated)
			return
		}

		if model.LinkRel == model2.LinkRel {
			t.Errorf(`model.LinkRel[%s] != model2.LinkRel[%s]`, model.LinkRel, model2.LinkRel)
			return
		}

		if model.LinkNotes == model2.LinkNotes {
			t.Errorf(`model.LinkNotes[%s] != model2.LinkNotes[%s]`, model.LinkNotes, model2.LinkNotes)
			return
		}

		if model.LinkRss == model2.LinkRss {
			t.Errorf(`model.LinkRss[%s] != model2.LinkRss[%s]`, model.LinkRss, model2.LinkRss)
			return
		}

		res34, err := model.FindByLinkUrl(model2.GetLinkUrl())
		if err != nil {
			t.Errorf(`failed model.FindByLinkUrl(model2.GetLinkUrl())`)
		}
		if len(res34) == 0 {
			t.Errorf(`failed to find any Link`)
		}

		res35, err := model.FindByLinkName(model2.GetLinkName())
		if err != nil {
			t.Errorf(`failed model.FindByLinkName(model2.GetLinkName())`)
		}
		if len(res35) == 0 {
			t.Errorf(`failed to find any Link`)
		}

		res36, err := model.FindByLinkImage(model2.GetLinkImage())
		if err != nil {
			t.Errorf(`failed model.FindByLinkImage(model2.GetLinkImage())`)
		}
		if len(res36) == 0 {
			t.Errorf(`failed to find any Link`)
		}

		res37, err := model.FindByLinkTarget(model2.GetLinkTarget())
		if err != nil {
			t.Errorf(`failed model.FindByLinkTarget(model2.GetLinkTarget())`)
		}
		if len(res37) == 0 {
			t.Errorf(`failed to find any Link`)
		}

		res38, err := model.FindByLinkDescription(model2.GetLinkDescription())
		if err != nil {
			t.Errorf(`failed model.FindByLinkDescription(model2.GetLinkDescription())`)
		}
		if len(res38) == 0 {
			t.Errorf(`failed to find any Link`)
		}

		res39, err := model.FindByLinkVisible(model2.GetLinkVisible())
		if err != nil {
			t.Errorf(`failed model.FindByLinkVisible(model2.GetLinkVisible())`)
		}
		if len(res39) == 0 {
			t.Errorf(`failed to find any Link`)
		}

		res40, err := model.FindByLinkOwner(model2.GetLinkOwner())
		if err != nil {
			t.Errorf(`failed model.FindByLinkOwner(model2.GetLinkOwner())`)
		}
		if len(res40) == 0 {
			t.Errorf(`failed to find any Link`)
		}

		res41, err := model.FindByLinkRating(model2.GetLinkRating())
		if err != nil {
			t.Errorf(`failed model.FindByLinkRating(model2.GetLinkRating())`)
		}
		if len(res41) == 0 {
			t.Errorf(`failed to find any Link`)
		}

		res42, err := model.FindByLinkUpdated(model2.GetLinkUpdated())
		if err != nil {
			t.Errorf(`failed model.FindByLinkUpdated(model2.GetLinkUpdated())`)
		}
		if len(res42) == 0 {
			t.Errorf(`failed to find any Link`)
		}

		res43, err := model.FindByLinkRel(model2.GetLinkRel())
		if err != nil {
			t.Errorf(`failed model.FindByLinkRel(model2.GetLinkRel())`)
		}
		if len(res43) == 0 {
			t.Errorf(`failed to find any Link`)
		}

		res44, err := model.FindByLinkNotes(model2.GetLinkNotes())
		if err != nil {
			t.Errorf(`failed model.FindByLinkNotes(model2.GetLinkNotes())`)
		}
		if len(res44) == 0 {
			t.Errorf(`failed to find any Link`)
		}

		res45, err := model.FindByLinkRss(model2.GetLinkRss())
		if err != nil {
			t.Errorf(`failed model.FindByLinkRss(model2.GetLinkRss())`)
		}
		if len(res45) == 0 {
			t.Errorf(`failed to find any Link`)
		}
	} // end of if fileExists
}

func TestLinkUpdaters(t *testing.T) {
	if fileExists(`../gopress.db.yml`) == false {
		return
	}
	a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
	defer a.Close()
	if err != nil {
		t.Errorf(`could not load ../gopress.db.yml %s`, err)
		return
	}
	file, err := os.OpenFile("adapter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Errorf("Failed to open log file %s", err)
		return
	}
	a.SetLogs(file)
	model := NewLink(a)

	model.SetLinkUrl(randomString(19))
	if model.GetLinkUrl() != model.LinkUrl {
		t.Errorf(`Link.GetLinkUrl() != Link.LinkUrl`)
	}
	if model.IsLinkUrlDirty != true {
		t.Errorf(`Link.IsLinkUrlDirty != true`)
		return
	}

	u0 := randomString(19)
	_, err = model.UpdateLinkUrl(u0)
	if err != nil {
		t.Errorf(`failed UpdateLinkUrl(u0) %s`, err)
		return
	}

	if model.GetLinkUrl() != u0 {
		t.Errorf(`Link.GetLinkUrl() != u0 after UpdateLinkUrl`)
		return
	}
	model.Reload()
	if model.GetLinkUrl() != u0 {
		t.Errorf(`Link.GetLinkUrl() != u0 after Reload`)
		return
	}

	model.SetLinkName(randomString(19))
	if model.GetLinkName() != model.LinkName {
		t.Errorf(`Link.GetLinkName() != Link.LinkName`)
	}
	if model.IsLinkNameDirty != true {
		t.Errorf(`Link.IsLinkNameDirty != true`)
		return
	}

	u1 := randomString(19)
	_, err = model.UpdateLinkName(u1)
	if err != nil {
		t.Errorf(`failed UpdateLinkName(u1) %s`, err)
		return
	}

	if model.GetLinkName() != u1 {
		t.Errorf(`Link.GetLinkName() != u1 after UpdateLinkName`)
		return
	}
	model.Reload()
	if model.GetLinkName() != u1 {
		t.Errorf(`Link.GetLinkName() != u1 after Reload`)
		return
	}

	model.SetLinkImage(randomString(19))
	if model.GetLinkImage() != model.LinkImage {
		t.Errorf(`Link.GetLinkImage() != Link.LinkImage`)
	}
	if model.IsLinkImageDirty != true {
		t.Errorf(`Link.IsLinkImageDirty != true`)
		return
	}

	u2 := randomString(19)
	_, err = model.UpdateLinkImage(u2)
	if err != nil {
		t.Errorf(`failed UpdateLinkImage(u2) %s`, err)
		return
	}

	if model.GetLinkImage() != u2 {
		t.Errorf(`Link.GetLinkImage() != u2 after UpdateLinkImage`)
		return
	}
	model.Reload()
	if model.GetLinkImage() != u2 {
		t.Errorf(`Link.GetLinkImage() != u2 after Reload`)
		return
	}

	model.SetLinkTarget(randomString(19))
	if model.GetLinkTarget() != model.LinkTarget {
		t.Errorf(`Link.GetLinkTarget() != Link.LinkTarget`)
	}
	if model.IsLinkTargetDirty != true {
		t.Errorf(`Link.IsLinkTargetDirty != true`)
		return
	}

	u3 := randomString(19)
	_, err = model.UpdateLinkTarget(u3)
	if err != nil {
		t.Errorf(`failed UpdateLinkTarget(u3) %s`, err)
		return
	}

	if model.GetLinkTarget() != u3 {
		t.Errorf(`Link.GetLinkTarget() != u3 after UpdateLinkTarget`)
		return
	}
	model.Reload()
	if model.GetLinkTarget() != u3 {
		t.Errorf(`Link.GetLinkTarget() != u3 after Reload`)
		return
	}

	model.SetLinkDescription(randomString(19))
	if model.GetLinkDescription() != model.LinkDescription {
		t.Errorf(`Link.GetLinkDescription() != Link.LinkDescription`)
	}
	if model.IsLinkDescriptionDirty != true {
		t.Errorf(`Link.IsLinkDescriptionDirty != true`)
		return
	}

	u4 := randomString(19)
	_, err = model.UpdateLinkDescription(u4)
	if err != nil {
		t.Errorf(`failed UpdateLinkDescription(u4) %s`, err)
		return
	}

	if model.GetLinkDescription() != u4 {
		t.Errorf(`Link.GetLinkDescription() != u4 after UpdateLinkDescription`)
		return
	}
	model.Reload()
	if model.GetLinkDescription() != u4 {
		t.Errorf(`Link.GetLinkDescription() != u4 after Reload`)
		return
	}

	model.SetLinkVisible(randomString(19))
	if model.GetLinkVisible() != model.LinkVisible {
		t.Errorf(`Link.GetLinkVisible() != Link.LinkVisible`)
	}
	if model.IsLinkVisibleDirty != true {
		t.Errorf(`Link.IsLinkVisibleDirty != true`)
		return
	}

	u5 := randomString(19)
	_, err = model.UpdateLinkVisible(u5)
	if err != nil {
		t.Errorf(`failed UpdateLinkVisible(u5) %s`, err)
		return
	}

	if model.GetLinkVisible() != u5 {
		t.Errorf(`Link.GetLinkVisible() != u5 after UpdateLinkVisible`)
		return
	}
	model.Reload()
	if model.GetLinkVisible() != u5 {
		t.Errorf(`Link.GetLinkVisible() != u5 after Reload`)
		return
	}

	model.SetLinkOwner(int64(randomInteger()))
	if model.GetLinkOwner() != model.LinkOwner {
		t.Errorf(`Link.GetLinkOwner() != Link.LinkOwner`)
	}
	if model.IsLinkOwnerDirty != true {
		t.Errorf(`Link.IsLinkOwnerDirty != true`)
		return
	}

	u6 := int64(randomInteger())
	_, err = model.UpdateLinkOwner(u6)
	if err != nil {
		t.Errorf(`failed UpdateLinkOwner(u6) %s`, err)
		return
	}

	if model.GetLinkOwner() != u6 {
		t.Errorf(`Link.GetLinkOwner() != u6 after UpdateLinkOwner`)
		return
	}
	model.Reload()
	if model.GetLinkOwner() != u6 {
		t.Errorf(`Link.GetLinkOwner() != u6 after Reload`)
		return
	}

	model.SetLinkRating(int(randomInteger()))
	if model.GetLinkRating() != model.LinkRating {
		t.Errorf(`Link.GetLinkRating() != Link.LinkRating`)
	}
	if model.IsLinkRatingDirty != true {
		t.Errorf(`Link.IsLinkRatingDirty != true`)
		return
	}

	u7 := int(randomInteger())
	_, err = model.UpdateLinkRating(u7)
	if err != nil {
		t.Errorf(`failed UpdateLinkRating(u7) %s`, err)
		return
	}

	if model.GetLinkRating() != u7 {
		t.Errorf(`Link.GetLinkRating() != u7 after UpdateLinkRating`)
		return
	}
	model.Reload()
	if model.GetLinkRating() != u7 {
		t.Errorf(`Link.GetLinkRating() != u7 after Reload`)
		return
	}

	model.SetLinkUpdated(randomDateTime(a))
	if model.GetLinkUpdated() != model.LinkUpdated {
		t.Errorf(`Link.GetLinkUpdated() != Link.LinkUpdated`)
	}
	if model.IsLinkUpdatedDirty != true {
		t.Errorf(`Link.IsLinkUpdatedDirty != true`)
		return
	}

	u8 := randomDateTime(a)
	_, err = model.UpdateLinkUpdated(u8)
	if err != nil {
		t.Errorf(`failed UpdateLinkUpdated(u8) %s`, err)
		return
	}

	if model.GetLinkUpdated() != u8 {
		t.Errorf(`Link.GetLinkUpdated() != u8 after UpdateLinkUpdated`)
		return
	}
	model.Reload()
	if model.GetLinkUpdated() != u8 {
		t.Errorf(`Link.GetLinkUpdated() != u8 after Reload`)
		return
	}

	model.SetLinkRel(randomString(19))
	if model.GetLinkRel() != model.LinkRel {
		t.Errorf(`Link.GetLinkRel() != Link.LinkRel`)
	}
	if model.IsLinkRelDirty != true {
		t.Errorf(`Link.IsLinkRelDirty != true`)
		return
	}

	u9 := randomString(19)
	_, err = model.UpdateLinkRel(u9)
	if err != nil {
		t.Errorf(`failed UpdateLinkRel(u9) %s`, err)
		return
	}

	if model.GetLinkRel() != u9 {
		t.Errorf(`Link.GetLinkRel() != u9 after UpdateLinkRel`)
		return
	}
	model.Reload()
	if model.GetLinkRel() != u9 {
		t.Errorf(`Link.GetLinkRel() != u9 after Reload`)
		return
	}

	model.SetLinkNotes(randomString(25))
	if model.GetLinkNotes() != model.LinkNotes {
		t.Errorf(`Link.GetLinkNotes() != Link.LinkNotes`)
	}
	if model.IsLinkNotesDirty != true {
		t.Errorf(`Link.IsLinkNotesDirty != true`)
		return
	}

	u10 := randomString(25)
	_, err = model.UpdateLinkNotes(u10)
	if err != nil {
		t.Errorf(`failed UpdateLinkNotes(u10) %s`, err)
		return
	}

	if model.GetLinkNotes() != u10 {
		t.Errorf(`Link.GetLinkNotes() != u10 after UpdateLinkNotes`)
		return
	}
	model.Reload()
	if model.GetLinkNotes() != u10 {
		t.Errorf(`Link.GetLinkNotes() != u10 after Reload`)
		return
	}

	model.SetLinkRss(randomString(19))
	if model.GetLinkRss() != model.LinkRss {
		t.Errorf(`Link.GetLinkRss() != Link.LinkRss`)
	}
	if model.IsLinkRssDirty != true {
		t.Errorf(`Link.IsLinkRssDirty != true`)
		return
	}

	u11 := randomString(19)
	_, err = model.UpdateLinkRss(u11)
	if err != nil {
		t.Errorf(`failed UpdateLinkRss(u11) %s`, err)
		return
	}

	if model.GetLinkRss() != u11 {
		t.Errorf(`Link.GetLinkRss() != u11 after UpdateLinkRss`)
		return
	}
	model.Reload()
	if model.GetLinkRss() != u11 {
		t.Errorf(`Link.GetLinkRss() != u11 after Reload`)
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

func TestOptionCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		defer a.Close()
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
		model.OptionName = randomString(19)
		model.OptionValue = randomString(25)
		model.Autoload = randomString(19)

		err = model.Create()
		if err != nil {
			t.Errorf(`failed to create model %s`, err)
			return
		}

		model2 := NewOption(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}

		if model.OptionName != model2.OptionName {
			t.Errorf(`model.OptionName[%s] != model2.OptionName[%s]`, model.OptionName, model2.OptionName)
			return
		}

		if model.OptionValue != model2.OptionValue {
			t.Errorf(`model.OptionValue[%s] != model2.OptionValue[%s]`, model.OptionValue, model2.OptionValue)
			return
		}

		if model.Autoload != model2.Autoload {
			t.Errorf(`model.Autoload[%s] != model2.Autoload[%s]`, model.Autoload, model2.Autoload)
			return
		}
		model2.SetOptionName(randomString(19))
		model2.SetOptionValue(randomString(25))
		model2.SetAutoload(randomString(19))

		err = model2.Save()
		if err != nil {
			t.Errorf(`failed to save model2 %s`, err)
		}

		if model.OptionName == model2.OptionName {
			t.Errorf(`model.OptionName[%s] != model2.OptionName[%s]`, model.OptionName, model2.OptionName)
			return
		}

		if model.OptionValue == model2.OptionValue {
			t.Errorf(`model.OptionValue[%s] != model2.OptionValue[%s]`, model.OptionValue, model2.OptionValue)
			return
		}

		if model.Autoload == model2.Autoload {
			t.Errorf(`model.Autoload[%s] != model2.Autoload[%s]`, model.Autoload, model2.Autoload)
			return
		}

		res9, err := model.FindByOptionName(model2.GetOptionName())
		if err != nil {
			t.Errorf(`failed model.FindByOptionName(model2.GetOptionName())`)
		}
		if len(res9) == 0 {
			t.Errorf(`failed to find any Option`)
		}

		res10, err := model.FindByOptionValue(model2.GetOptionValue())
		if err != nil {
			t.Errorf(`failed model.FindByOptionValue(model2.GetOptionValue())`)
		}
		if len(res10) == 0 {
			t.Errorf(`failed to find any Option`)
		}

		res11, err := model.FindByAutoload(model2.GetAutoload())
		if err != nil {
			t.Errorf(`failed model.FindByAutoload(model2.GetAutoload())`)
		}
		if len(res11) == 0 {
			t.Errorf(`failed to find any Option`)
		}
	} // end of if fileExists
}

func TestOptionUpdaters(t *testing.T) {
	if fileExists(`../gopress.db.yml`) == false {
		return
	}
	a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
	defer a.Close()
	if err != nil {
		t.Errorf(`could not load ../gopress.db.yml %s`, err)
		return
	}
	file, err := os.OpenFile("adapter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Errorf("Failed to open log file %s", err)
		return
	}
	a.SetLogs(file)
	model := NewOption(a)

	model.SetOptionName(randomString(19))
	if model.GetOptionName() != model.OptionName {
		t.Errorf(`Option.GetOptionName() != Option.OptionName`)
	}
	if model.IsOptionNameDirty != true {
		t.Errorf(`Option.IsOptionNameDirty != true`)
		return
	}

	u0 := randomString(19)
	_, err = model.UpdateOptionName(u0)
	if err != nil {
		t.Errorf(`failed UpdateOptionName(u0) %s`, err)
		return
	}

	if model.GetOptionName() != u0 {
		t.Errorf(`Option.GetOptionName() != u0 after UpdateOptionName`)
		return
	}
	model.Reload()
	if model.GetOptionName() != u0 {
		t.Errorf(`Option.GetOptionName() != u0 after Reload`)
		return
	}

	model.SetOptionValue(randomString(25))
	if model.GetOptionValue() != model.OptionValue {
		t.Errorf(`Option.GetOptionValue() != Option.OptionValue`)
	}
	if model.IsOptionValueDirty != true {
		t.Errorf(`Option.IsOptionValueDirty != true`)
		return
	}

	u1 := randomString(25)
	_, err = model.UpdateOptionValue(u1)
	if err != nil {
		t.Errorf(`failed UpdateOptionValue(u1) %s`, err)
		return
	}

	if model.GetOptionValue() != u1 {
		t.Errorf(`Option.GetOptionValue() != u1 after UpdateOptionValue`)
		return
	}
	model.Reload()
	if model.GetOptionValue() != u1 {
		t.Errorf(`Option.GetOptionValue() != u1 after Reload`)
		return
	}

	model.SetAutoload(randomString(19))
	if model.GetAutoload() != model.Autoload {
		t.Errorf(`Option.GetAutoload() != Option.Autoload`)
	}
	if model.IsAutoloadDirty != true {
		t.Errorf(`Option.IsAutoloadDirty != true`)
		return
	}

	u2 := randomString(19)
	_, err = model.UpdateAutoload(u2)
	if err != nil {
		t.Errorf(`failed UpdateAutoload(u2) %s`, err)
		return
	}

	if model.GetAutoload() != u2 {
		t.Errorf(`Option.GetAutoload() != u2 after UpdateAutoload`)
		return
	}
	model.Reload()
	if model.GetAutoload() != u2 {
		t.Errorf(`Option.GetAutoload() != u2 after Reload`)
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
		defer a.Close()
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
		model.PostId = int64(randomInteger())
		model.MetaKey = randomString(19)
		model.MetaValue = randomString(25)

		err = model.Create()
		if err != nil {
			t.Errorf(`failed to create model %s`, err)
			return
		}

		model2 := NewPostMeta(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}

		if model.PostId != model2.PostId {
			t.Errorf(`model.PostId[%d] != model2.PostId[%d]`, model.PostId, model2.PostId)
			return
		}

		if model.MetaKey != model2.MetaKey {
			t.Errorf(`model.MetaKey[%s] != model2.MetaKey[%s]`, model.MetaKey, model2.MetaKey)
			return
		}

		if model.MetaValue != model2.MetaValue {
			t.Errorf(`model.MetaValue[%s] != model2.MetaValue[%s]`, model.MetaValue, model2.MetaValue)
			return
		}
		model2.SetPostId(int64(randomInteger()))
		model2.SetMetaKey(randomString(19))
		model2.SetMetaValue(randomString(25))

		err = model2.Save()
		if err != nil {
			t.Errorf(`failed to save model2 %s`, err)
		}

		if model.PostId == model2.PostId {
			t.Errorf(`model.PostId[%d] != model2.PostId[%d]`, model.PostId, model2.PostId)
			return
		}

		if model.MetaKey == model2.MetaKey {
			t.Errorf(`model.MetaKey[%s] != model2.MetaKey[%s]`, model.MetaKey, model2.MetaKey)
			return
		}

		if model.MetaValue == model2.MetaValue {
			t.Errorf(`model.MetaValue[%s] != model2.MetaValue[%s]`, model.MetaValue, model2.MetaValue)
			return
		}

		res9, err := model.FindByPostId(model2.GetPostId())
		if err != nil {
			t.Errorf(`failed model.FindByPostId(model2.GetPostId())`)
		}
		if len(res9) == 0 {
			t.Errorf(`failed to find any PostMeta`)
		}

		res10, err := model.FindByMetaKey(model2.GetMetaKey())
		if err != nil {
			t.Errorf(`failed model.FindByMetaKey(model2.GetMetaKey())`)
		}
		if len(res10) == 0 {
			t.Errorf(`failed to find any PostMeta`)
		}

		res11, err := model.FindByMetaValue(model2.GetMetaValue())
		if err != nil {
			t.Errorf(`failed model.FindByMetaValue(model2.GetMetaValue())`)
		}
		if len(res11) == 0 {
			t.Errorf(`failed to find any PostMeta`)
		}
	} // end of if fileExists
}

func TestPostMetaUpdaters(t *testing.T) {
	if fileExists(`../gopress.db.yml`) == false {
		return
	}
	a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
	defer a.Close()
	if err != nil {
		t.Errorf(`could not load ../gopress.db.yml %s`, err)
		return
	}
	file, err := os.OpenFile("adapter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Errorf("Failed to open log file %s", err)
		return
	}
	a.SetLogs(file)
	model := NewPostMeta(a)

	model.SetPostId(int64(randomInteger()))
	if model.GetPostId() != model.PostId {
		t.Errorf(`PostMeta.GetPostId() != PostMeta.PostId`)
	}
	if model.IsPostIdDirty != true {
		t.Errorf(`PostMeta.IsPostIdDirty != true`)
		return
	}

	u0 := int64(randomInteger())
	_, err = model.UpdatePostId(u0)
	if err != nil {
		t.Errorf(`failed UpdatePostId(u0) %s`, err)
		return
	}

	if model.GetPostId() != u0 {
		t.Errorf(`PostMeta.GetPostId() != u0 after UpdatePostId`)
		return
	}
	model.Reload()
	if model.GetPostId() != u0 {
		t.Errorf(`PostMeta.GetPostId() != u0 after Reload`)
		return
	}

	model.SetMetaKey(randomString(19))
	if model.GetMetaKey() != model.MetaKey {
		t.Errorf(`PostMeta.GetMetaKey() != PostMeta.MetaKey`)
	}
	if model.IsMetaKeyDirty != true {
		t.Errorf(`PostMeta.IsMetaKeyDirty != true`)
		return
	}

	u1 := randomString(19)
	_, err = model.UpdateMetaKey(u1)
	if err != nil {
		t.Errorf(`failed UpdateMetaKey(u1) %s`, err)
		return
	}

	if model.GetMetaKey() != u1 {
		t.Errorf(`PostMeta.GetMetaKey() != u1 after UpdateMetaKey`)
		return
	}
	model.Reload()
	if model.GetMetaKey() != u1 {
		t.Errorf(`PostMeta.GetMetaKey() != u1 after Reload`)
		return
	}

	model.SetMetaValue(randomString(25))
	if model.GetMetaValue() != model.MetaValue {
		t.Errorf(`PostMeta.GetMetaValue() != PostMeta.MetaValue`)
	}
	if model.IsMetaValueDirty != true {
		t.Errorf(`PostMeta.IsMetaValueDirty != true`)
		return
	}

	u2 := randomString(25)
	_, err = model.UpdateMetaValue(u2)
	if err != nil {
		t.Errorf(`failed UpdateMetaValue(u2) %s`, err)
		return
	}

	if model.GetMetaValue() != u2 {
		t.Errorf(`PostMeta.GetMetaValue() != u2 after UpdateMetaValue`)
		return
	}
	model.Reload()
	if model.GetMetaValue() != u2 {
		t.Errorf(`PostMeta.GetMetaValue() != u2 after Reload`)
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

func TestPostCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		defer a.Close()
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
		model.PostAuthor = int64(randomInteger())
		model.PostDate = randomDateTime(a)
		model.PostDateGmt = randomDateTime(a)
		model.PostContent = randomString(25)
		model.PostTitle = randomString(25)
		model.PostExcerpt = randomString(25)
		model.PostStatus = randomString(19)
		model.CommentStatus = randomString(19)
		model.PingStatus = randomString(19)
		model.PostPassword = randomString(19)
		model.PostName = randomString(19)
		model.ToPing = randomString(25)
		model.Pinged = randomString(25)
		model.PostModified = randomDateTime(a)
		model.PostModifiedGmt = randomDateTime(a)
		model.PostContentFiltered = randomString(25)
		model.PostParent = int64(randomInteger())
		model.Guid = randomString(19)
		model.MenuOrder = int(randomInteger())
		model.PostType = randomString(19)
		model.PostMimeType = randomString(19)
		model.CommentCount = int64(randomInteger())

		err = model.Create()
		if err != nil {
			t.Errorf(`failed to create model %s`, err)
			return
		}

		model2 := NewPost(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}

		if model.PostAuthor != model2.PostAuthor {
			t.Errorf(`model.PostAuthor[%d] != model2.PostAuthor[%d]`, model.PostAuthor, model2.PostAuthor)
			return
		}

		if model.PostDate.Year != model2.PostDate.Year ||
			model.PostDate.Month != model2.PostDate.Month ||
			model.PostDate.Day != model2.PostDate.Day ||
			model.PostDate.Hours != model2.PostDate.Hours ||
			model.PostDate.Minutes != model2.PostDate.Minutes ||
			model.PostDate.Seconds != model2.PostDate.Seconds {
			t.Errorf(`model.PostDate != model2.PostDate %+v --- %+v`, model.PostDate, model2.PostDate)
			return
		}

		if model.PostDateGmt.Year != model2.PostDateGmt.Year ||
			model.PostDateGmt.Month != model2.PostDateGmt.Month ||
			model.PostDateGmt.Day != model2.PostDateGmt.Day ||
			model.PostDateGmt.Hours != model2.PostDateGmt.Hours ||
			model.PostDateGmt.Minutes != model2.PostDateGmt.Minutes ||
			model.PostDateGmt.Seconds != model2.PostDateGmt.Seconds {
			t.Errorf(`model.PostDateGmt != model2.PostDateGmt %+v --- %+v`, model.PostDateGmt, model2.PostDateGmt)
			return
		}

		if model.PostContent != model2.PostContent {
			t.Errorf(`model.PostContent[%s] != model2.PostContent[%s]`, model.PostContent, model2.PostContent)
			return
		}

		if model.PostTitle != model2.PostTitle {
			t.Errorf(`model.PostTitle[%s] != model2.PostTitle[%s]`, model.PostTitle, model2.PostTitle)
			return
		}

		if model.PostExcerpt != model2.PostExcerpt {
			t.Errorf(`model.PostExcerpt[%s] != model2.PostExcerpt[%s]`, model.PostExcerpt, model2.PostExcerpt)
			return
		}

		if model.PostStatus != model2.PostStatus {
			t.Errorf(`model.PostStatus[%s] != model2.PostStatus[%s]`, model.PostStatus, model2.PostStatus)
			return
		}

		if model.CommentStatus != model2.CommentStatus {
			t.Errorf(`model.CommentStatus[%s] != model2.CommentStatus[%s]`, model.CommentStatus, model2.CommentStatus)
			return
		}

		if model.PingStatus != model2.PingStatus {
			t.Errorf(`model.PingStatus[%s] != model2.PingStatus[%s]`, model.PingStatus, model2.PingStatus)
			return
		}

		if model.PostPassword != model2.PostPassword {
			t.Errorf(`model.PostPassword[%s] != model2.PostPassword[%s]`, model.PostPassword, model2.PostPassword)
			return
		}

		if model.PostName != model2.PostName {
			t.Errorf(`model.PostName[%s] != model2.PostName[%s]`, model.PostName, model2.PostName)
			return
		}

		if model.ToPing != model2.ToPing {
			t.Errorf(`model.ToPing[%s] != model2.ToPing[%s]`, model.ToPing, model2.ToPing)
			return
		}

		if model.Pinged != model2.Pinged {
			t.Errorf(`model.Pinged[%s] != model2.Pinged[%s]`, model.Pinged, model2.Pinged)
			return
		}

		if model.PostModified.Year != model2.PostModified.Year ||
			model.PostModified.Month != model2.PostModified.Month ||
			model.PostModified.Day != model2.PostModified.Day ||
			model.PostModified.Hours != model2.PostModified.Hours ||
			model.PostModified.Minutes != model2.PostModified.Minutes ||
			model.PostModified.Seconds != model2.PostModified.Seconds {
			t.Errorf(`model.PostModified != model2.PostModified %+v --- %+v`, model.PostModified, model2.PostModified)
			return
		}

		if model.PostModifiedGmt.Year != model2.PostModifiedGmt.Year ||
			model.PostModifiedGmt.Month != model2.PostModifiedGmt.Month ||
			model.PostModifiedGmt.Day != model2.PostModifiedGmt.Day ||
			model.PostModifiedGmt.Hours != model2.PostModifiedGmt.Hours ||
			model.PostModifiedGmt.Minutes != model2.PostModifiedGmt.Minutes ||
			model.PostModifiedGmt.Seconds != model2.PostModifiedGmt.Seconds {
			t.Errorf(`model.PostModifiedGmt != model2.PostModifiedGmt %+v --- %+v`, model.PostModifiedGmt, model2.PostModifiedGmt)
			return
		}

		if model.PostContentFiltered != model2.PostContentFiltered {
			t.Errorf(`model.PostContentFiltered[%s] != model2.PostContentFiltered[%s]`, model.PostContentFiltered, model2.PostContentFiltered)
			return
		}

		if model.PostParent != model2.PostParent {
			t.Errorf(`model.PostParent[%d] != model2.PostParent[%d]`, model.PostParent, model2.PostParent)
			return
		}

		if model.Guid != model2.Guid {
			t.Errorf(`model.Guid[%s] != model2.Guid[%s]`, model.Guid, model2.Guid)
			return
		}

		if model.MenuOrder != model2.MenuOrder {
			t.Errorf(`model.MenuOrder[%d] != model2.MenuOrder[%d]`, model.MenuOrder, model2.MenuOrder)
			return
		}

		if model.PostType != model2.PostType {
			t.Errorf(`model.PostType[%s] != model2.PostType[%s]`, model.PostType, model2.PostType)
			return
		}

		if model.PostMimeType != model2.PostMimeType {
			t.Errorf(`model.PostMimeType[%s] != model2.PostMimeType[%s]`, model.PostMimeType, model2.PostMimeType)
			return
		}

		if model.CommentCount != model2.CommentCount {
			t.Errorf(`model.CommentCount[%d] != model2.CommentCount[%d]`, model.CommentCount, model2.CommentCount)
			return
		}
		model2.SetPostAuthor(int64(randomInteger()))
		model2.SetPostDate(randomDateTime(a))
		model2.SetPostDateGmt(randomDateTime(a))
		model2.SetPostContent(randomString(25))
		model2.SetPostTitle(randomString(25))
		model2.SetPostExcerpt(randomString(25))
		model2.SetPostStatus(randomString(19))
		model2.SetCommentStatus(randomString(19))
		model2.SetPingStatus(randomString(19))
		model2.SetPostPassword(randomString(19))
		model2.SetPostName(randomString(19))
		model2.SetToPing(randomString(25))
		model2.SetPinged(randomString(25))
		model2.SetPostModified(randomDateTime(a))
		model2.SetPostModifiedGmt(randomDateTime(a))
		model2.SetPostContentFiltered(randomString(25))
		model2.SetPostParent(int64(randomInteger()))
		model2.SetGuid(randomString(19))
		model2.SetMenuOrder(int(randomInteger()))
		model2.SetPostType(randomString(19))
		model2.SetPostMimeType(randomString(19))
		model2.SetCommentCount(int64(randomInteger()))

		err = model2.Save()
		if err != nil {
			t.Errorf(`failed to save model2 %s`, err)
		}

		if model.PostAuthor == model2.PostAuthor {
			t.Errorf(`model.PostAuthor[%d] != model2.PostAuthor[%d]`, model.PostAuthor, model2.PostAuthor)
			return
		}

		if model.PostDate.Year == model2.PostDate.Year {
			t.Errorf(`model.PostDate.Year == model2.PostDate but should not!`, model.PostDate, model2.PostDate)
			return
		}

		if model.PostDateGmt.Year == model2.PostDateGmt.Year {
			t.Errorf(`model.PostDateGmt.Year == model2.PostDateGmt but should not!`, model.PostDateGmt, model2.PostDateGmt)
			return
		}

		if model.PostContent == model2.PostContent {
			t.Errorf(`model.PostContent[%s] != model2.PostContent[%s]`, model.PostContent, model2.PostContent)
			return
		}

		if model.PostTitle == model2.PostTitle {
			t.Errorf(`model.PostTitle[%s] != model2.PostTitle[%s]`, model.PostTitle, model2.PostTitle)
			return
		}

		if model.PostExcerpt == model2.PostExcerpt {
			t.Errorf(`model.PostExcerpt[%s] != model2.PostExcerpt[%s]`, model.PostExcerpt, model2.PostExcerpt)
			return
		}

		if model.PostStatus == model2.PostStatus {
			t.Errorf(`model.PostStatus[%s] != model2.PostStatus[%s]`, model.PostStatus, model2.PostStatus)
			return
		}

		if model.CommentStatus == model2.CommentStatus {
			t.Errorf(`model.CommentStatus[%s] != model2.CommentStatus[%s]`, model.CommentStatus, model2.CommentStatus)
			return
		}

		if model.PingStatus == model2.PingStatus {
			t.Errorf(`model.PingStatus[%s] != model2.PingStatus[%s]`, model.PingStatus, model2.PingStatus)
			return
		}

		if model.PostPassword == model2.PostPassword {
			t.Errorf(`model.PostPassword[%s] != model2.PostPassword[%s]`, model.PostPassword, model2.PostPassword)
			return
		}

		if model.PostName == model2.PostName {
			t.Errorf(`model.PostName[%s] != model2.PostName[%s]`, model.PostName, model2.PostName)
			return
		}

		if model.ToPing == model2.ToPing {
			t.Errorf(`model.ToPing[%s] != model2.ToPing[%s]`, model.ToPing, model2.ToPing)
			return
		}

		if model.Pinged == model2.Pinged {
			t.Errorf(`model.Pinged[%s] != model2.Pinged[%s]`, model.Pinged, model2.Pinged)
			return
		}

		if model.PostModified.Year == model2.PostModified.Year {
			t.Errorf(`model.PostModified.Year == model2.PostModified but should not!`, model.PostModified, model2.PostModified)
			return
		}

		if model.PostModifiedGmt.Year == model2.PostModifiedGmt.Year {
			t.Errorf(`model.PostModifiedGmt.Year == model2.PostModifiedGmt but should not!`, model.PostModifiedGmt, model2.PostModifiedGmt)
			return
		}

		if model.PostContentFiltered == model2.PostContentFiltered {
			t.Errorf(`model.PostContentFiltered[%s] != model2.PostContentFiltered[%s]`, model.PostContentFiltered, model2.PostContentFiltered)
			return
		}

		if model.PostParent == model2.PostParent {
			t.Errorf(`model.PostParent[%d] != model2.PostParent[%d]`, model.PostParent, model2.PostParent)
			return
		}

		if model.Guid == model2.Guid {
			t.Errorf(`model.Guid[%s] != model2.Guid[%s]`, model.Guid, model2.Guid)
			return
		}

		if model.MenuOrder == model2.MenuOrder {
			t.Errorf(`model.MenuOrder[%d] != model2.MenuOrder[%d]`, model.MenuOrder, model2.MenuOrder)
			return
		}

		if model.PostType == model2.PostType {
			t.Errorf(`model.PostType[%s] != model2.PostType[%s]`, model.PostType, model2.PostType)
			return
		}

		if model.PostMimeType == model2.PostMimeType {
			t.Errorf(`model.PostMimeType[%s] != model2.PostMimeType[%s]`, model.PostMimeType, model2.PostMimeType)
			return
		}

		if model.CommentCount == model2.CommentCount {
			t.Errorf(`model.CommentCount[%d] != model2.CommentCount[%d]`, model.CommentCount, model2.CommentCount)
			return
		}

		res58, err := model.FindByPostAuthor(model2.GetPostAuthor())
		if err != nil {
			t.Errorf(`failed model.FindByPostAuthor(model2.GetPostAuthor())`)
		}
		if len(res58) == 0 {
			t.Errorf(`failed to find any Post`)
		}

		res59, err := model.FindByPostDate(model2.GetPostDate())
		if err != nil {
			t.Errorf(`failed model.FindByPostDate(model2.GetPostDate())`)
		}
		if len(res59) == 0 {
			t.Errorf(`failed to find any Post`)
		}

		res60, err := model.FindByPostDateGmt(model2.GetPostDateGmt())
		if err != nil {
			t.Errorf(`failed model.FindByPostDateGmt(model2.GetPostDateGmt())`)
		}
		if len(res60) == 0 {
			t.Errorf(`failed to find any Post`)
		}

		res61, err := model.FindByPostContent(model2.GetPostContent())
		if err != nil {
			t.Errorf(`failed model.FindByPostContent(model2.GetPostContent())`)
		}
		if len(res61) == 0 {
			t.Errorf(`failed to find any Post`)
		}

		res62, err := model.FindByPostTitle(model2.GetPostTitle())
		if err != nil {
			t.Errorf(`failed model.FindByPostTitle(model2.GetPostTitle())`)
		}
		if len(res62) == 0 {
			t.Errorf(`failed to find any Post`)
		}

		res63, err := model.FindByPostExcerpt(model2.GetPostExcerpt())
		if err != nil {
			t.Errorf(`failed model.FindByPostExcerpt(model2.GetPostExcerpt())`)
		}
		if len(res63) == 0 {
			t.Errorf(`failed to find any Post`)
		}

		res64, err := model.FindByPostStatus(model2.GetPostStatus())
		if err != nil {
			t.Errorf(`failed model.FindByPostStatus(model2.GetPostStatus())`)
		}
		if len(res64) == 0 {
			t.Errorf(`failed to find any Post`)
		}

		res65, err := model.FindByCommentStatus(model2.GetCommentStatus())
		if err != nil {
			t.Errorf(`failed model.FindByCommentStatus(model2.GetCommentStatus())`)
		}
		if len(res65) == 0 {
			t.Errorf(`failed to find any Post`)
		}

		res66, err := model.FindByPingStatus(model2.GetPingStatus())
		if err != nil {
			t.Errorf(`failed model.FindByPingStatus(model2.GetPingStatus())`)
		}
		if len(res66) == 0 {
			t.Errorf(`failed to find any Post`)
		}

		res67, err := model.FindByPostPassword(model2.GetPostPassword())
		if err != nil {
			t.Errorf(`failed model.FindByPostPassword(model2.GetPostPassword())`)
		}
		if len(res67) == 0 {
			t.Errorf(`failed to find any Post`)
		}

		res68, err := model.FindByPostName(model2.GetPostName())
		if err != nil {
			t.Errorf(`failed model.FindByPostName(model2.GetPostName())`)
		}
		if len(res68) == 0 {
			t.Errorf(`failed to find any Post`)
		}

		res69, err := model.FindByToPing(model2.GetToPing())
		if err != nil {
			t.Errorf(`failed model.FindByToPing(model2.GetToPing())`)
		}
		if len(res69) == 0 {
			t.Errorf(`failed to find any Post`)
		}

		res70, err := model.FindByPinged(model2.GetPinged())
		if err != nil {
			t.Errorf(`failed model.FindByPinged(model2.GetPinged())`)
		}
		if len(res70) == 0 {
			t.Errorf(`failed to find any Post`)
		}

		res71, err := model.FindByPostModified(model2.GetPostModified())
		if err != nil {
			t.Errorf(`failed model.FindByPostModified(model2.GetPostModified())`)
		}
		if len(res71) == 0 {
			t.Errorf(`failed to find any Post`)
		}

		res72, err := model.FindByPostModifiedGmt(model2.GetPostModifiedGmt())
		if err != nil {
			t.Errorf(`failed model.FindByPostModifiedGmt(model2.GetPostModifiedGmt())`)
		}
		if len(res72) == 0 {
			t.Errorf(`failed to find any Post`)
		}

		res73, err := model.FindByPostContentFiltered(model2.GetPostContentFiltered())
		if err != nil {
			t.Errorf(`failed model.FindByPostContentFiltered(model2.GetPostContentFiltered())`)
		}
		if len(res73) == 0 {
			t.Errorf(`failed to find any Post`)
		}

		res74, err := model.FindByPostParent(model2.GetPostParent())
		if err != nil {
			t.Errorf(`failed model.FindByPostParent(model2.GetPostParent())`)
		}
		if len(res74) == 0 {
			t.Errorf(`failed to find any Post`)
		}

		res75, err := model.FindByGuid(model2.GetGuid())
		if err != nil {
			t.Errorf(`failed model.FindByGuid(model2.GetGuid())`)
		}
		if len(res75) == 0 {
			t.Errorf(`failed to find any Post`)
		}

		res76, err := model.FindByPostType(model2.GetPostType())
		if err != nil {
			t.Errorf(`failed model.FindByPostType(model2.GetPostType())`)
		}
		if len(res76) == 0 {
			t.Errorf(`failed to find any Post`)
		}

		res77, err := model.FindByPostMimeType(model2.GetPostMimeType())
		if err != nil {
			t.Errorf(`failed model.FindByPostMimeType(model2.GetPostMimeType())`)
		}
		if len(res77) == 0 {
			t.Errorf(`failed to find any Post`)
		}

		res78, err := model.FindByCommentCount(model2.GetCommentCount())
		if err != nil {
			t.Errorf(`failed model.FindByCommentCount(model2.GetCommentCount())`)
		}
		if len(res78) == 0 {
			t.Errorf(`failed to find any Post`)
		}
	} // end of if fileExists
}

func TestPostUpdaters(t *testing.T) {
	if fileExists(`../gopress.db.yml`) == false {
		return
	}
	a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
	defer a.Close()
	if err != nil {
		t.Errorf(`could not load ../gopress.db.yml %s`, err)
		return
	}
	file, err := os.OpenFile("adapter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Errorf("Failed to open log file %s", err)
		return
	}
	a.SetLogs(file)
	model := NewPost(a)

	model.SetPostAuthor(int64(randomInteger()))
	if model.GetPostAuthor() != model.PostAuthor {
		t.Errorf(`Post.GetPostAuthor() != Post.PostAuthor`)
	}
	if model.IsPostAuthorDirty != true {
		t.Errorf(`Post.IsPostAuthorDirty != true`)
		return
	}

	u0 := int64(randomInteger())
	_, err = model.UpdatePostAuthor(u0)
	if err != nil {
		t.Errorf(`failed UpdatePostAuthor(u0) %s`, err)
		return
	}

	if model.GetPostAuthor() != u0 {
		t.Errorf(`Post.GetPostAuthor() != u0 after UpdatePostAuthor`)
		return
	}
	model.Reload()
	if model.GetPostAuthor() != u0 {
		t.Errorf(`Post.GetPostAuthor() != u0 after Reload`)
		return
	}

	model.SetPostDate(randomDateTime(a))
	if model.GetPostDate() != model.PostDate {
		t.Errorf(`Post.GetPostDate() != Post.PostDate`)
	}
	if model.IsPostDateDirty != true {
		t.Errorf(`Post.IsPostDateDirty != true`)
		return
	}

	u1 := randomDateTime(a)
	_, err = model.UpdatePostDate(u1)
	if err != nil {
		t.Errorf(`failed UpdatePostDate(u1) %s`, err)
		return
	}

	if model.GetPostDate() != u1 {
		t.Errorf(`Post.GetPostDate() != u1 after UpdatePostDate`)
		return
	}
	model.Reload()
	if model.GetPostDate() != u1 {
		t.Errorf(`Post.GetPostDate() != u1 after Reload`)
		return
	}

	model.SetPostDateGmt(randomDateTime(a))
	if model.GetPostDateGmt() != model.PostDateGmt {
		t.Errorf(`Post.GetPostDateGmt() != Post.PostDateGmt`)
	}
	if model.IsPostDateGmtDirty != true {
		t.Errorf(`Post.IsPostDateGmtDirty != true`)
		return
	}

	u2 := randomDateTime(a)
	_, err = model.UpdatePostDateGmt(u2)
	if err != nil {
		t.Errorf(`failed UpdatePostDateGmt(u2) %s`, err)
		return
	}

	if model.GetPostDateGmt() != u2 {
		t.Errorf(`Post.GetPostDateGmt() != u2 after UpdatePostDateGmt`)
		return
	}
	model.Reload()
	if model.GetPostDateGmt() != u2 {
		t.Errorf(`Post.GetPostDateGmt() != u2 after Reload`)
		return
	}

	model.SetPostContent(randomString(25))
	if model.GetPostContent() != model.PostContent {
		t.Errorf(`Post.GetPostContent() != Post.PostContent`)
	}
	if model.IsPostContentDirty != true {
		t.Errorf(`Post.IsPostContentDirty != true`)
		return
	}

	u3 := randomString(25)
	_, err = model.UpdatePostContent(u3)
	if err != nil {
		t.Errorf(`failed UpdatePostContent(u3) %s`, err)
		return
	}

	if model.GetPostContent() != u3 {
		t.Errorf(`Post.GetPostContent() != u3 after UpdatePostContent`)
		return
	}
	model.Reload()
	if model.GetPostContent() != u3 {
		t.Errorf(`Post.GetPostContent() != u3 after Reload`)
		return
	}

	model.SetPostTitle(randomString(25))
	if model.GetPostTitle() != model.PostTitle {
		t.Errorf(`Post.GetPostTitle() != Post.PostTitle`)
	}
	if model.IsPostTitleDirty != true {
		t.Errorf(`Post.IsPostTitleDirty != true`)
		return
	}

	u4 := randomString(25)
	_, err = model.UpdatePostTitle(u4)
	if err != nil {
		t.Errorf(`failed UpdatePostTitle(u4) %s`, err)
		return
	}

	if model.GetPostTitle() != u4 {
		t.Errorf(`Post.GetPostTitle() != u4 after UpdatePostTitle`)
		return
	}
	model.Reload()
	if model.GetPostTitle() != u4 {
		t.Errorf(`Post.GetPostTitle() != u4 after Reload`)
		return
	}

	model.SetPostExcerpt(randomString(25))
	if model.GetPostExcerpt() != model.PostExcerpt {
		t.Errorf(`Post.GetPostExcerpt() != Post.PostExcerpt`)
	}
	if model.IsPostExcerptDirty != true {
		t.Errorf(`Post.IsPostExcerptDirty != true`)
		return
	}

	u5 := randomString(25)
	_, err = model.UpdatePostExcerpt(u5)
	if err != nil {
		t.Errorf(`failed UpdatePostExcerpt(u5) %s`, err)
		return
	}

	if model.GetPostExcerpt() != u5 {
		t.Errorf(`Post.GetPostExcerpt() != u5 after UpdatePostExcerpt`)
		return
	}
	model.Reload()
	if model.GetPostExcerpt() != u5 {
		t.Errorf(`Post.GetPostExcerpt() != u5 after Reload`)
		return
	}

	model.SetPostStatus(randomString(19))
	if model.GetPostStatus() != model.PostStatus {
		t.Errorf(`Post.GetPostStatus() != Post.PostStatus`)
	}
	if model.IsPostStatusDirty != true {
		t.Errorf(`Post.IsPostStatusDirty != true`)
		return
	}

	u6 := randomString(19)
	_, err = model.UpdatePostStatus(u6)
	if err != nil {
		t.Errorf(`failed UpdatePostStatus(u6) %s`, err)
		return
	}

	if model.GetPostStatus() != u6 {
		t.Errorf(`Post.GetPostStatus() != u6 after UpdatePostStatus`)
		return
	}
	model.Reload()
	if model.GetPostStatus() != u6 {
		t.Errorf(`Post.GetPostStatus() != u6 after Reload`)
		return
	}

	model.SetCommentStatus(randomString(19))
	if model.GetCommentStatus() != model.CommentStatus {
		t.Errorf(`Post.GetCommentStatus() != Post.CommentStatus`)
	}
	if model.IsCommentStatusDirty != true {
		t.Errorf(`Post.IsCommentStatusDirty != true`)
		return
	}

	u7 := randomString(19)
	_, err = model.UpdateCommentStatus(u7)
	if err != nil {
		t.Errorf(`failed UpdateCommentStatus(u7) %s`, err)
		return
	}

	if model.GetCommentStatus() != u7 {
		t.Errorf(`Post.GetCommentStatus() != u7 after UpdateCommentStatus`)
		return
	}
	model.Reload()
	if model.GetCommentStatus() != u7 {
		t.Errorf(`Post.GetCommentStatus() != u7 after Reload`)
		return
	}

	model.SetPingStatus(randomString(19))
	if model.GetPingStatus() != model.PingStatus {
		t.Errorf(`Post.GetPingStatus() != Post.PingStatus`)
	}
	if model.IsPingStatusDirty != true {
		t.Errorf(`Post.IsPingStatusDirty != true`)
		return
	}

	u8 := randomString(19)
	_, err = model.UpdatePingStatus(u8)
	if err != nil {
		t.Errorf(`failed UpdatePingStatus(u8) %s`, err)
		return
	}

	if model.GetPingStatus() != u8 {
		t.Errorf(`Post.GetPingStatus() != u8 after UpdatePingStatus`)
		return
	}
	model.Reload()
	if model.GetPingStatus() != u8 {
		t.Errorf(`Post.GetPingStatus() != u8 after Reload`)
		return
	}

	model.SetPostPassword(randomString(19))
	if model.GetPostPassword() != model.PostPassword {
		t.Errorf(`Post.GetPostPassword() != Post.PostPassword`)
	}
	if model.IsPostPasswordDirty != true {
		t.Errorf(`Post.IsPostPasswordDirty != true`)
		return
	}

	u9 := randomString(19)
	_, err = model.UpdatePostPassword(u9)
	if err != nil {
		t.Errorf(`failed UpdatePostPassword(u9) %s`, err)
		return
	}

	if model.GetPostPassword() != u9 {
		t.Errorf(`Post.GetPostPassword() != u9 after UpdatePostPassword`)
		return
	}
	model.Reload()
	if model.GetPostPassword() != u9 {
		t.Errorf(`Post.GetPostPassword() != u9 after Reload`)
		return
	}

	model.SetPostName(randomString(19))
	if model.GetPostName() != model.PostName {
		t.Errorf(`Post.GetPostName() != Post.PostName`)
	}
	if model.IsPostNameDirty != true {
		t.Errorf(`Post.IsPostNameDirty != true`)
		return
	}

	u10 := randomString(19)
	_, err = model.UpdatePostName(u10)
	if err != nil {
		t.Errorf(`failed UpdatePostName(u10) %s`, err)
		return
	}

	if model.GetPostName() != u10 {
		t.Errorf(`Post.GetPostName() != u10 after UpdatePostName`)
		return
	}
	model.Reload()
	if model.GetPostName() != u10 {
		t.Errorf(`Post.GetPostName() != u10 after Reload`)
		return
	}

	model.SetToPing(randomString(25))
	if model.GetToPing() != model.ToPing {
		t.Errorf(`Post.GetToPing() != Post.ToPing`)
	}
	if model.IsToPingDirty != true {
		t.Errorf(`Post.IsToPingDirty != true`)
		return
	}

	u11 := randomString(25)
	_, err = model.UpdateToPing(u11)
	if err != nil {
		t.Errorf(`failed UpdateToPing(u11) %s`, err)
		return
	}

	if model.GetToPing() != u11 {
		t.Errorf(`Post.GetToPing() != u11 after UpdateToPing`)
		return
	}
	model.Reload()
	if model.GetToPing() != u11 {
		t.Errorf(`Post.GetToPing() != u11 after Reload`)
		return
	}

	model.SetPinged(randomString(25))
	if model.GetPinged() != model.Pinged {
		t.Errorf(`Post.GetPinged() != Post.Pinged`)
	}
	if model.IsPingedDirty != true {
		t.Errorf(`Post.IsPingedDirty != true`)
		return
	}

	u12 := randomString(25)
	_, err = model.UpdatePinged(u12)
	if err != nil {
		t.Errorf(`failed UpdatePinged(u12) %s`, err)
		return
	}

	if model.GetPinged() != u12 {
		t.Errorf(`Post.GetPinged() != u12 after UpdatePinged`)
		return
	}
	model.Reload()
	if model.GetPinged() != u12 {
		t.Errorf(`Post.GetPinged() != u12 after Reload`)
		return
	}

	model.SetPostModified(randomDateTime(a))
	if model.GetPostModified() != model.PostModified {
		t.Errorf(`Post.GetPostModified() != Post.PostModified`)
	}
	if model.IsPostModifiedDirty != true {
		t.Errorf(`Post.IsPostModifiedDirty != true`)
		return
	}

	u13 := randomDateTime(a)
	_, err = model.UpdatePostModified(u13)
	if err != nil {
		t.Errorf(`failed UpdatePostModified(u13) %s`, err)
		return
	}

	if model.GetPostModified() != u13 {
		t.Errorf(`Post.GetPostModified() != u13 after UpdatePostModified`)
		return
	}
	model.Reload()
	if model.GetPostModified() != u13 {
		t.Errorf(`Post.GetPostModified() != u13 after Reload`)
		return
	}

	model.SetPostModifiedGmt(randomDateTime(a))
	if model.GetPostModifiedGmt() != model.PostModifiedGmt {
		t.Errorf(`Post.GetPostModifiedGmt() != Post.PostModifiedGmt`)
	}
	if model.IsPostModifiedGmtDirty != true {
		t.Errorf(`Post.IsPostModifiedGmtDirty != true`)
		return
	}

	u14 := randomDateTime(a)
	_, err = model.UpdatePostModifiedGmt(u14)
	if err != nil {
		t.Errorf(`failed UpdatePostModifiedGmt(u14) %s`, err)
		return
	}

	if model.GetPostModifiedGmt() != u14 {
		t.Errorf(`Post.GetPostModifiedGmt() != u14 after UpdatePostModifiedGmt`)
		return
	}
	model.Reload()
	if model.GetPostModifiedGmt() != u14 {
		t.Errorf(`Post.GetPostModifiedGmt() != u14 after Reload`)
		return
	}

	model.SetPostContentFiltered(randomString(25))
	if model.GetPostContentFiltered() != model.PostContentFiltered {
		t.Errorf(`Post.GetPostContentFiltered() != Post.PostContentFiltered`)
	}
	if model.IsPostContentFilteredDirty != true {
		t.Errorf(`Post.IsPostContentFilteredDirty != true`)
		return
	}

	u15 := randomString(25)
	_, err = model.UpdatePostContentFiltered(u15)
	if err != nil {
		t.Errorf(`failed UpdatePostContentFiltered(u15) %s`, err)
		return
	}

	if model.GetPostContentFiltered() != u15 {
		t.Errorf(`Post.GetPostContentFiltered() != u15 after UpdatePostContentFiltered`)
		return
	}
	model.Reload()
	if model.GetPostContentFiltered() != u15 {
		t.Errorf(`Post.GetPostContentFiltered() != u15 after Reload`)
		return
	}

	model.SetPostParent(int64(randomInteger()))
	if model.GetPostParent() != model.PostParent {
		t.Errorf(`Post.GetPostParent() != Post.PostParent`)
	}
	if model.IsPostParentDirty != true {
		t.Errorf(`Post.IsPostParentDirty != true`)
		return
	}

	u16 := int64(randomInteger())
	_, err = model.UpdatePostParent(u16)
	if err != nil {
		t.Errorf(`failed UpdatePostParent(u16) %s`, err)
		return
	}

	if model.GetPostParent() != u16 {
		t.Errorf(`Post.GetPostParent() != u16 after UpdatePostParent`)
		return
	}
	model.Reload()
	if model.GetPostParent() != u16 {
		t.Errorf(`Post.GetPostParent() != u16 after Reload`)
		return
	}

	model.SetGuid(randomString(19))
	if model.GetGuid() != model.Guid {
		t.Errorf(`Post.GetGuid() != Post.Guid`)
	}
	if model.IsGuidDirty != true {
		t.Errorf(`Post.IsGuidDirty != true`)
		return
	}

	u17 := randomString(19)
	_, err = model.UpdateGuid(u17)
	if err != nil {
		t.Errorf(`failed UpdateGuid(u17) %s`, err)
		return
	}

	if model.GetGuid() != u17 {
		t.Errorf(`Post.GetGuid() != u17 after UpdateGuid`)
		return
	}
	model.Reload()
	if model.GetGuid() != u17 {
		t.Errorf(`Post.GetGuid() != u17 after Reload`)
		return
	}

	model.SetMenuOrder(int(randomInteger()))
	if model.GetMenuOrder() != model.MenuOrder {
		t.Errorf(`Post.GetMenuOrder() != Post.MenuOrder`)
	}
	if model.IsMenuOrderDirty != true {
		t.Errorf(`Post.IsMenuOrderDirty != true`)
		return
	}

	u18 := int(randomInteger())
	_, err = model.UpdateMenuOrder(u18)
	if err != nil {
		t.Errorf(`failed UpdateMenuOrder(u18) %s`, err)
		return
	}

	if model.GetMenuOrder() != u18 {
		t.Errorf(`Post.GetMenuOrder() != u18 after UpdateMenuOrder`)
		return
	}
	model.Reload()
	if model.GetMenuOrder() != u18 {
		t.Errorf(`Post.GetMenuOrder() != u18 after Reload`)
		return
	}

	model.SetPostType(randomString(19))
	if model.GetPostType() != model.PostType {
		t.Errorf(`Post.GetPostType() != Post.PostType`)
	}
	if model.IsPostTypeDirty != true {
		t.Errorf(`Post.IsPostTypeDirty != true`)
		return
	}

	u19 := randomString(19)
	_, err = model.UpdatePostType(u19)
	if err != nil {
		t.Errorf(`failed UpdatePostType(u19) %s`, err)
		return
	}

	if model.GetPostType() != u19 {
		t.Errorf(`Post.GetPostType() != u19 after UpdatePostType`)
		return
	}
	model.Reload()
	if model.GetPostType() != u19 {
		t.Errorf(`Post.GetPostType() != u19 after Reload`)
		return
	}

	model.SetPostMimeType(randomString(19))
	if model.GetPostMimeType() != model.PostMimeType {
		t.Errorf(`Post.GetPostMimeType() != Post.PostMimeType`)
	}
	if model.IsPostMimeTypeDirty != true {
		t.Errorf(`Post.IsPostMimeTypeDirty != true`)
		return
	}

	u20 := randomString(19)
	_, err = model.UpdatePostMimeType(u20)
	if err != nil {
		t.Errorf(`failed UpdatePostMimeType(u20) %s`, err)
		return
	}

	if model.GetPostMimeType() != u20 {
		t.Errorf(`Post.GetPostMimeType() != u20 after UpdatePostMimeType`)
		return
	}
	model.Reload()
	if model.GetPostMimeType() != u20 {
		t.Errorf(`Post.GetPostMimeType() != u20 after Reload`)
		return
	}

	model.SetCommentCount(int64(randomInteger()))
	if model.GetCommentCount() != model.CommentCount {
		t.Errorf(`Post.GetCommentCount() != Post.CommentCount`)
	}
	if model.IsCommentCountDirty != true {
		t.Errorf(`Post.IsCommentCountDirty != true`)
		return
	}

	u21 := int64(randomInteger())
	_, err = model.UpdateCommentCount(u21)
	if err != nil {
		t.Errorf(`failed UpdateCommentCount(u21) %s`, err)
		return
	}

	if model.GetCommentCount() != u21 {
		t.Errorf(`Post.GetCommentCount() != u21 after UpdateCommentCount`)
		return
	}
	model.Reload()
	if model.GetCommentCount() != u21 {
		t.Errorf(`Post.GetCommentCount() != u21 after Reload`)
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

func TestTermRelationshipCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		defer a.Close()
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
		model.ObjectId = int64(randomInteger())
		model.TermTaxonomyId = int64(randomInteger())
		model.TermOrder = int(randomInteger())

		err = model.Create()
		if err != nil {
			t.Errorf(`failed to create model %s`, err)
			return
		}

		model2 := NewTermRelationship(a)
		found, err := model2.Find(model.TermTaxonomyId, model.ObjectId)
		if model.TermTaxonomyId == 0 {
			t.Errorf(`it's 0`)
		}

		if err != nil {
			t.Errorf(`did not find record for term_taxonomy_id = %d AND object_id = %d because of %s`, model.TermTaxonomyId, model.ObjectId, err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for term_taxonomy_id = %d AND object_id = %d because of %s`, model.TermTaxonomyId, model.ObjectId, err)
			return
		}

		if model.ObjectId != model2.ObjectId {
			t.Errorf(`model.ObjectId[%d] != model2.ObjectId[%d]`, model.ObjectId, model2.ObjectId)
			return
		}

		if model.TermTaxonomyId != model2.TermTaxonomyId {
			t.Errorf(`model.TermTaxonomyId[%d] != model2.TermTaxonomyId[%d]`, model.TermTaxonomyId, model2.TermTaxonomyId)
			return
		}

		if model.TermOrder != model2.TermOrder {
			t.Errorf(`model.TermOrder[%d] != model2.TermOrder[%d]`, model.TermOrder, model2.TermOrder)
			return
		}
		model2.SetObjectId(int64(randomInteger()))
		model2.SetTermTaxonomyId(int64(randomInteger()))
		model2.SetTermOrder(int(randomInteger()))

		err = model2.Save()
		if err != nil {
			t.Errorf(`failed to save model2 %s`, err)
		}

		if model.ObjectId == model2.ObjectId {
			t.Errorf(`model.ObjectId[%d] != model2.ObjectId[%d]`, model.ObjectId, model2.ObjectId)
			return
		}

		if model.TermTaxonomyId == model2.TermTaxonomyId {
			t.Errorf(`model.TermTaxonomyId[%d] != model2.TermTaxonomyId[%d]`, model.TermTaxonomyId, model2.TermTaxonomyId)
			return
		}

		if model.TermOrder == model2.TermOrder {
			t.Errorf(`model.TermOrder[%d] != model2.TermOrder[%d]`, model.TermOrder, model2.TermOrder)
			return
		}
	} // end of if fileExists
}

func TestTermRelationshipUpdaters(t *testing.T) {
	if fileExists(`../gopress.db.yml`) == false {
		return
	}
	a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
	defer a.Close()
	if err != nil {
		t.Errorf(`could not load ../gopress.db.yml %s`, err)
		return
	}
	file, err := os.OpenFile("adapter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Errorf("Failed to open log file %s", err)
		return
	}
	a.SetLogs(file)
	model := NewTermRelationship(a)

	model.SetObjectId(int64(randomInteger()))
	if model.GetObjectId() != model.ObjectId {
		t.Errorf(`TermRelationship.GetObjectId() != TermRelationship.ObjectId`)
	}
	if model.IsObjectIdDirty != true {
		t.Errorf(`TermRelationship.IsObjectIdDirty != true`)
		return
	}

	u0 := int64(randomInteger())
	_, err = model.UpdateObjectId(u0)
	if err != nil {
		t.Errorf(`failed UpdateObjectId(u0) %s`, err)
		return
	}

	if model.GetObjectId() != u0 {
		t.Errorf(`TermRelationship.GetObjectId() != u0 after UpdateObjectId`)
		return
	}
	model.Reload()
	if model.GetObjectId() != u0 {
		t.Errorf(`TermRelationship.GetObjectId() != u0 after Reload`)
		return
	}

	model.SetTermOrder(int(randomInteger()))
	if model.GetTermOrder() != model.TermOrder {
		t.Errorf(`TermRelationship.GetTermOrder() != TermRelationship.TermOrder`)
	}
	if model.IsTermOrderDirty != true {
		t.Errorf(`TermRelationship.IsTermOrderDirty != true`)
		return
	}

	u1 := int(randomInteger())
	_, err = model.UpdateTermOrder(u1)
	if err != nil {
		t.Errorf(`failed UpdateTermOrder(u1) %s`, err)
		return
	}

	if model.GetTermOrder() != u1 {
		t.Errorf(`TermRelationship.GetTermOrder() != u1 after UpdateTermOrder`)
		return
	}
	model.Reload()
	if model.GetTermOrder() != u1 {
		t.Errorf(`TermRelationship.GetTermOrder() != u1 after Reload`)
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

func TestTermTaxonomyCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		defer a.Close()
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
		model.TermId = int64(randomInteger())
		model.Taxonomy = randomString(19)
		model.Description = randomString(25)
		model.Parent = int64(randomInteger())
		model.Count = int64(randomInteger())

		err = model.Create()
		if err != nil {
			t.Errorf(`failed to create model %s`, err)
			return
		}

		model2 := NewTermTaxonomy(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}

		if model.TermId != model2.TermId {
			t.Errorf(`model.TermId[%d] != model2.TermId[%d]`, model.TermId, model2.TermId)
			return
		}

		if model.Taxonomy != model2.Taxonomy {
			t.Errorf(`model.Taxonomy[%s] != model2.Taxonomy[%s]`, model.Taxonomy, model2.Taxonomy)
			return
		}

		if model.Description != model2.Description {
			t.Errorf(`model.Description[%s] != model2.Description[%s]`, model.Description, model2.Description)
			return
		}

		if model.Parent != model2.Parent {
			t.Errorf(`model.Parent[%d] != model2.Parent[%d]`, model.Parent, model2.Parent)
			return
		}

		if model.Count != model2.Count {
			t.Errorf(`model.Count[%d] != model2.Count[%d]`, model.Count, model2.Count)
			return
		}
		model2.SetTermId(int64(randomInteger()))
		model2.SetTaxonomy(randomString(19))
		model2.SetDescription(randomString(25))
		model2.SetParent(int64(randomInteger()))
		model2.SetCount(int64(randomInteger()))

		err = model2.Save()
		if err != nil {
			t.Errorf(`failed to save model2 %s`, err)
		}

		if model.TermId == model2.TermId {
			t.Errorf(`model.TermId[%d] != model2.TermId[%d]`, model.TermId, model2.TermId)
			return
		}

		if model.Taxonomy == model2.Taxonomy {
			t.Errorf(`model.Taxonomy[%s] != model2.Taxonomy[%s]`, model.Taxonomy, model2.Taxonomy)
			return
		}

		if model.Description == model2.Description {
			t.Errorf(`model.Description[%s] != model2.Description[%s]`, model.Description, model2.Description)
			return
		}

		if model.Parent == model2.Parent {
			t.Errorf(`model.Parent[%d] != model2.Parent[%d]`, model.Parent, model2.Parent)
			return
		}

		if model.Count == model2.Count {
			t.Errorf(`model.Count[%d] != model2.Count[%d]`, model.Count, model2.Count)
			return
		}

		res15, err := model.FindByTermId(model2.GetTermId())
		if err != nil {
			t.Errorf(`failed model.FindByTermId(model2.GetTermId())`)
		}
		if len(res15) == 0 {
			t.Errorf(`failed to find any TermTaxonomy`)
		}

		res16, err := model.FindByTaxonomy(model2.GetTaxonomy())
		if err != nil {
			t.Errorf(`failed model.FindByTaxonomy(model2.GetTaxonomy())`)
		}
		if len(res16) == 0 {
			t.Errorf(`failed to find any TermTaxonomy`)
		}

		res17, err := model.FindByDescription(model2.GetDescription())
		if err != nil {
			t.Errorf(`failed model.FindByDescription(model2.GetDescription())`)
		}
		if len(res17) == 0 {
			t.Errorf(`failed to find any TermTaxonomy`)
		}

		res18, err := model.FindByParent(model2.GetParent())
		if err != nil {
			t.Errorf(`failed model.FindByParent(model2.GetParent())`)
		}
		if len(res18) == 0 {
			t.Errorf(`failed to find any TermTaxonomy`)
		}

		res19, err := model.FindByCount(model2.GetCount())
		if err != nil {
			t.Errorf(`failed model.FindByCount(model2.GetCount())`)
		}
		if len(res19) == 0 {
			t.Errorf(`failed to find any TermTaxonomy`)
		}
	} // end of if fileExists
}

func TestTermTaxonomyUpdaters(t *testing.T) {
	if fileExists(`../gopress.db.yml`) == false {
		return
	}
	a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
	defer a.Close()
	if err != nil {
		t.Errorf(`could not load ../gopress.db.yml %s`, err)
		return
	}
	file, err := os.OpenFile("adapter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Errorf("Failed to open log file %s", err)
		return
	}
	a.SetLogs(file)
	model := NewTermTaxonomy(a)

	model.SetTermId(int64(randomInteger()))
	if model.GetTermId() != model.TermId {
		t.Errorf(`TermTaxonomy.GetTermId() != TermTaxonomy.TermId`)
	}
	if model.IsTermIdDirty != true {
		t.Errorf(`TermTaxonomy.IsTermIdDirty != true`)
		return
	}

	u0 := int64(randomInteger())
	_, err = model.UpdateTermId(u0)
	if err != nil {
		t.Errorf(`failed UpdateTermId(u0) %s`, err)
		return
	}

	if model.GetTermId() != u0 {
		t.Errorf(`TermTaxonomy.GetTermId() != u0 after UpdateTermId`)
		return
	}
	model.Reload()
	if model.GetTermId() != u0 {
		t.Errorf(`TermTaxonomy.GetTermId() != u0 after Reload`)
		return
	}

	model.SetTaxonomy(randomString(19))
	if model.GetTaxonomy() != model.Taxonomy {
		t.Errorf(`TermTaxonomy.GetTaxonomy() != TermTaxonomy.Taxonomy`)
	}
	if model.IsTaxonomyDirty != true {
		t.Errorf(`TermTaxonomy.IsTaxonomyDirty != true`)
		return
	}

	u1 := randomString(19)
	_, err = model.UpdateTaxonomy(u1)
	if err != nil {
		t.Errorf(`failed UpdateTaxonomy(u1) %s`, err)
		return
	}

	if model.GetTaxonomy() != u1 {
		t.Errorf(`TermTaxonomy.GetTaxonomy() != u1 after UpdateTaxonomy`)
		return
	}
	model.Reload()
	if model.GetTaxonomy() != u1 {
		t.Errorf(`TermTaxonomy.GetTaxonomy() != u1 after Reload`)
		return
	}

	model.SetDescription(randomString(25))
	if model.GetDescription() != model.Description {
		t.Errorf(`TermTaxonomy.GetDescription() != TermTaxonomy.Description`)
	}
	if model.IsDescriptionDirty != true {
		t.Errorf(`TermTaxonomy.IsDescriptionDirty != true`)
		return
	}

	u2 := randomString(25)
	_, err = model.UpdateDescription(u2)
	if err != nil {
		t.Errorf(`failed UpdateDescription(u2) %s`, err)
		return
	}

	if model.GetDescription() != u2 {
		t.Errorf(`TermTaxonomy.GetDescription() != u2 after UpdateDescription`)
		return
	}
	model.Reload()
	if model.GetDescription() != u2 {
		t.Errorf(`TermTaxonomy.GetDescription() != u2 after Reload`)
		return
	}

	model.SetParent(int64(randomInteger()))
	if model.GetParent() != model.Parent {
		t.Errorf(`TermTaxonomy.GetParent() != TermTaxonomy.Parent`)
	}
	if model.IsParentDirty != true {
		t.Errorf(`TermTaxonomy.IsParentDirty != true`)
		return
	}

	u3 := int64(randomInteger())
	_, err = model.UpdateParent(u3)
	if err != nil {
		t.Errorf(`failed UpdateParent(u3) %s`, err)
		return
	}

	if model.GetParent() != u3 {
		t.Errorf(`TermTaxonomy.GetParent() != u3 after UpdateParent`)
		return
	}
	model.Reload()
	if model.GetParent() != u3 {
		t.Errorf(`TermTaxonomy.GetParent() != u3 after Reload`)
		return
	}

	model.SetCount(int64(randomInteger()))
	if model.GetCount() != model.Count {
		t.Errorf(`TermTaxonomy.GetCount() != TermTaxonomy.Count`)
	}
	if model.IsCountDirty != true {
		t.Errorf(`TermTaxonomy.IsCountDirty != true`)
		return
	}

	u4 := int64(randomInteger())
	_, err = model.UpdateCount(u4)
	if err != nil {
		t.Errorf(`failed UpdateCount(u4) %s`, err)
		return
	}

	if model.GetCount() != u4 {
		t.Errorf(`TermTaxonomy.GetCount() != u4 after UpdateCount`)
		return
	}
	model.Reload()
	if model.GetCount() != u4 {
		t.Errorf(`TermTaxonomy.GetCount() != u4 after Reload`)
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

func TestTermCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		defer a.Close()
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
		model.Name = randomString(19)
		model.Slug = randomString(19)
		model.TermGroup = int64(randomInteger())

		err = model.Create()
		if err != nil {
			t.Errorf(`failed to create model %s`, err)
			return
		}

		model2 := NewTerm(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}

		if model.Name != model2.Name {
			t.Errorf(`model.Name[%s] != model2.Name[%s]`, model.Name, model2.Name)
			return
		}

		if model.Slug != model2.Slug {
			t.Errorf(`model.Slug[%s] != model2.Slug[%s]`, model.Slug, model2.Slug)
			return
		}

		if model.TermGroup != model2.TermGroup {
			t.Errorf(`model.TermGroup[%d] != model2.TermGroup[%d]`, model.TermGroup, model2.TermGroup)
			return
		}
		model2.SetName(randomString(19))
		model2.SetSlug(randomString(19))
		model2.SetTermGroup(int64(randomInteger()))

		err = model2.Save()
		if err != nil {
			t.Errorf(`failed to save model2 %s`, err)
		}

		if model.Name == model2.Name {
			t.Errorf(`model.Name[%s] != model2.Name[%s]`, model.Name, model2.Name)
			return
		}

		if model.Slug == model2.Slug {
			t.Errorf(`model.Slug[%s] != model2.Slug[%s]`, model.Slug, model2.Slug)
			return
		}

		if model.TermGroup == model2.TermGroup {
			t.Errorf(`model.TermGroup[%d] != model2.TermGroup[%d]`, model.TermGroup, model2.TermGroup)
			return
		}

		res9, err := model.FindByName(model2.GetName())
		if err != nil {
			t.Errorf(`failed model.FindByName(model2.GetName())`)
		}
		if len(res9) == 0 {
			t.Errorf(`failed to find any Term`)
		}

		res10, err := model.FindBySlug(model2.GetSlug())
		if err != nil {
			t.Errorf(`failed model.FindBySlug(model2.GetSlug())`)
		}
		if len(res10) == 0 {
			t.Errorf(`failed to find any Term`)
		}

		res11, err := model.FindByTermGroup(model2.GetTermGroup())
		if err != nil {
			t.Errorf(`failed model.FindByTermGroup(model2.GetTermGroup())`)
		}
		if len(res11) == 0 {
			t.Errorf(`failed to find any Term`)
		}
	} // end of if fileExists
}

func TestTermUpdaters(t *testing.T) {
	if fileExists(`../gopress.db.yml`) == false {
		return
	}
	a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
	defer a.Close()
	if err != nil {
		t.Errorf(`could not load ../gopress.db.yml %s`, err)
		return
	}
	file, err := os.OpenFile("adapter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Errorf("Failed to open log file %s", err)
		return
	}
	a.SetLogs(file)
	model := NewTerm(a)

	model.SetName(randomString(19))
	if model.GetName() != model.Name {
		t.Errorf(`Term.GetName() != Term.Name`)
	}
	if model.IsNameDirty != true {
		t.Errorf(`Term.IsNameDirty != true`)
		return
	}

	u0 := randomString(19)
	_, err = model.UpdateName(u0)
	if err != nil {
		t.Errorf(`failed UpdateName(u0) %s`, err)
		return
	}

	if model.GetName() != u0 {
		t.Errorf(`Term.GetName() != u0 after UpdateName`)
		return
	}
	model.Reload()
	if model.GetName() != u0 {
		t.Errorf(`Term.GetName() != u0 after Reload`)
		return
	}

	model.SetSlug(randomString(19))
	if model.GetSlug() != model.Slug {
		t.Errorf(`Term.GetSlug() != Term.Slug`)
	}
	if model.IsSlugDirty != true {
		t.Errorf(`Term.IsSlugDirty != true`)
		return
	}

	u1 := randomString(19)
	_, err = model.UpdateSlug(u1)
	if err != nil {
		t.Errorf(`failed UpdateSlug(u1) %s`, err)
		return
	}

	if model.GetSlug() != u1 {
		t.Errorf(`Term.GetSlug() != u1 after UpdateSlug`)
		return
	}
	model.Reload()
	if model.GetSlug() != u1 {
		t.Errorf(`Term.GetSlug() != u1 after Reload`)
		return
	}

	model.SetTermGroup(int64(randomInteger()))
	if model.GetTermGroup() != model.TermGroup {
		t.Errorf(`Term.GetTermGroup() != Term.TermGroup`)
	}
	if model.IsTermGroupDirty != true {
		t.Errorf(`Term.IsTermGroupDirty != true`)
		return
	}

	u2 := int64(randomInteger())
	_, err = model.UpdateTermGroup(u2)
	if err != nil {
		t.Errorf(`failed UpdateTermGroup(u2) %s`, err)
		return
	}

	if model.GetTermGroup() != u2 {
		t.Errorf(`Term.GetTermGroup() != u2 after UpdateTermGroup`)
		return
	}
	model.Reload()
	if model.GetTermGroup() != u2 {
		t.Errorf(`Term.GetTermGroup() != u2 after Reload`)
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

func TestUserMetaCreate(t *testing.T) {
	if fileExists(`../gopress.db.yml`) {
		a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
		defer a.Close()
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
		model.UserId = int64(randomInteger())
		model.MetaKey = randomString(19)
		model.MetaValue = randomString(25)

		err = model.Create()
		if err != nil {
			t.Errorf(`failed to create model %s`, err)
			return
		}

		model2 := NewUserMeta(a)
		found, err := model2.Find(model.GetPrimaryKeyValue())
		if err != nil {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}
		if found == false {
			t.Errorf(`did not find record for %s = %d because of %s`, model.GetPrimaryKeyName(), model.GetPrimaryKeyValue(), err)
			return
		}

		if model.UserId != model2.UserId {
			t.Errorf(`model.UserId[%d] != model2.UserId[%d]`, model.UserId, model2.UserId)
			return
		}

		if model.MetaKey != model2.MetaKey {
			t.Errorf(`model.MetaKey[%s] != model2.MetaKey[%s]`, model.MetaKey, model2.MetaKey)
			return
		}

		if model.MetaValue != model2.MetaValue {
			t.Errorf(`model.MetaValue[%s] != model2.MetaValue[%s]`, model.MetaValue, model2.MetaValue)
			return
		}
		model2.SetUserId(int64(randomInteger()))
		model2.SetMetaKey(randomString(19))
		model2.SetMetaValue(randomString(25))

		err = model2.Save()
		if err != nil {
			t.Errorf(`failed to save model2 %s`, err)
		}

		if model.UserId == model2.UserId {
			t.Errorf(`model.UserId[%d] != model2.UserId[%d]`, model.UserId, model2.UserId)
			return
		}

		if model.MetaKey == model2.MetaKey {
			t.Errorf(`model.MetaKey[%s] != model2.MetaKey[%s]`, model.MetaKey, model2.MetaKey)
			return
		}

		if model.MetaValue == model2.MetaValue {
			t.Errorf(`model.MetaValue[%s] != model2.MetaValue[%s]`, model.MetaValue, model2.MetaValue)
			return
		}

		res9, err := model.FindByUserId(model2.GetUserId())
		if err != nil {
			t.Errorf(`failed model.FindByUserId(model2.GetUserId())`)
		}
		if len(res9) == 0 {
			t.Errorf(`failed to find any UserMeta`)
		}

		res10, err := model.FindByMetaKey(model2.GetMetaKey())
		if err != nil {
			t.Errorf(`failed model.FindByMetaKey(model2.GetMetaKey())`)
		}
		if len(res10) == 0 {
			t.Errorf(`failed to find any UserMeta`)
		}

		res11, err := model.FindByMetaValue(model2.GetMetaValue())
		if err != nil {
			t.Errorf(`failed model.FindByMetaValue(model2.GetMetaValue())`)
		}
		if len(res11) == 0 {
			t.Errorf(`failed to find any UserMeta`)
		}
	} // end of if fileExists
}

func TestUserMetaUpdaters(t *testing.T) {
	if fileExists(`../gopress.db.yml`) == false {
		return
	}
	a, err := NewMysqlAdapterEx(`../gopress.db.yml`)
	defer a.Close()
	if err != nil {
		t.Errorf(`could not load ../gopress.db.yml %s`, err)
		return
	}
	file, err := os.OpenFile("adapter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Errorf("Failed to open log file %s", err)
		return
	}
	a.SetLogs(file)
	model := NewUserMeta(a)

	model.SetUserId(int64(randomInteger()))
	if model.GetUserId() != model.UserId {
		t.Errorf(`UserMeta.GetUserId() != UserMeta.UserId`)
	}
	if model.IsUserIdDirty != true {
		t.Errorf(`UserMeta.IsUserIdDirty != true`)
		return
	}

	u0 := int64(randomInteger())
	_, err = model.UpdateUserId(u0)
	if err != nil {
		t.Errorf(`failed UpdateUserId(u0) %s`, err)
		return
	}

	if model.GetUserId() != u0 {
		t.Errorf(`UserMeta.GetUserId() != u0 after UpdateUserId`)
		return
	}
	model.Reload()
	if model.GetUserId() != u0 {
		t.Errorf(`UserMeta.GetUserId() != u0 after Reload`)
		return
	}

	model.SetMetaKey(randomString(19))
	if model.GetMetaKey() != model.MetaKey {
		t.Errorf(`UserMeta.GetMetaKey() != UserMeta.MetaKey`)
	}
	if model.IsMetaKeyDirty != true {
		t.Errorf(`UserMeta.IsMetaKeyDirty != true`)
		return
	}

	u1 := randomString(19)
	_, err = model.UpdateMetaKey(u1)
	if err != nil {
		t.Errorf(`failed UpdateMetaKey(u1) %s`, err)
		return
	}

	if model.GetMetaKey() != u1 {
		t.Errorf(`UserMeta.GetMetaKey() != u1 after UpdateMetaKey`)
		return
	}
	model.Reload()
	if model.GetMetaKey() != u1 {
		t.Errorf(`UserMeta.GetMetaKey() != u1 after Reload`)
		return
	}

	model.SetMetaValue(randomString(25))
	if model.GetMetaValue() != model.MetaValue {
		t.Errorf(`UserMeta.GetMetaValue() != UserMeta.MetaValue`)
	}
	if model.IsMetaValueDirty != true {
		t.Errorf(`UserMeta.IsMetaValueDirty != true`)
		return
	}

	u2 := randomString(25)
	_, err = model.UpdateMetaValue(u2)
	if err != nil {
		t.Errorf(`failed UpdateMetaValue(u2) %s`, err)
		return
	}

	if model.GetMetaValue() != u2 {
		t.Errorf(`UserMeta.GetMetaValue() != u2 after UpdateMetaValue`)
		return
	}
	model.Reload()
	if model.GetMetaValue() != u2 {
		t.Errorf(`UserMeta.GetMetaValue() != u2 after Reload`)
		return
	}

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
func TestAdapterFailures(t *testing.T) {
	_, err := NewMysqlAdapterEx(`file_that_does_not_exist123323`)
	if err == nil {
		t.Errorf(`Did not receive an error when file should not exist!`)
		return
	}
	// Load a nonsense yaml file
	_, err = NewMysqlAdapterEx(`test_data/nonsenseyaml.yml`)
	if err == nil {
		t.Errorf(`this should fail to load a nonsense yaml file`)
		return
	}
	// Load a test yaml with wrong Open
	_, err = NewMysqlAdapterEx(`test_data/adapter.yml`)
	if err == nil {
		t.Errorf(`this should fail with wrong login info`)
		return
	}
	// Load a silly yml file with wrong data
	_, err = NewMysqlAdapterEx(`test_data/silly.yml`)
	if err == nil {
		t.Errorf(`this should fail with wrong login info`)
		return
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

func TestAdapterInfoLogging(t *testing.T) {
	a := NewMysqlAdapter(`wp_`)
	var b bytes.Buffer
	r, err := regexp.Compile(`\[INFO\]:.+Hello World`)
	if err != nil {
		t.Errorf(`could not compile regex`)
		return
	}
	wr := bufio.NewWriter(&b)
	a.SetLogs(wr)
	a.LogInfo(`Hello World`)
	wr.Flush()
	if r.MatchString(b.String()) == false {
		t.Errorf(`failed to match info line`)
		return
	}
}
func TestAdapterEmptyInfoLogging(t *testing.T) {
	a := NewMysqlAdapter(`wp_`)
	var b bytes.Buffer
	wr := bufio.NewWriter(&b)
	a.SetLogs(wr)
	a.LogInfo(``)
	wr.Flush()
	if b.String() != `` {
		t.Errorf(`Info should not occur in this case`)
		return
	}
	a.SetLogFilter(func(tag string, val string) string {
		return ``
	})
	a.LogInfo(`Hello World`)
	wr.Flush()
	if b.String() != `` {
		t.Errorf(`Info should not occur due to filter in this case`)
		return
	}
}

func TestAdapterDebugLogging(t *testing.T) {
	a := NewMysqlAdapter(`wp_`)
	var b bytes.Buffer
	r, err := regexp.Compile(`\[DEBUG\]:.+Hello World`)
	if err != nil {
		t.Errorf(`could not compile regex`)
		return
	}
	wr := bufio.NewWriter(&b)
	a.SetLogs(wr)
	a.LogDebug(`Hello World`)
	wr.Flush()
	if r.MatchString(b.String()) == false {
		t.Errorf(`failed to match info line`)
		return
	}
}
func TestAdapterEmptyDebugLogging(t *testing.T) {
	a := NewMysqlAdapter(`wp_`)
	var b bytes.Buffer
	wr := bufio.NewWriter(&b)
	a.SetLogs(wr)
	a.LogDebug(``)
	wr.Flush()
	if b.String() != `` {
		t.Errorf(`Info should not occur in this case`)
		return
	}
	a.SetLogFilter(func(tag string, val string) string {
		return ``
	})
	a.LogDebug(`Hello World`)
	wr.Flush()
	if b.String() != `` {
		t.Errorf(`Info should not occur due to filter in this case`)
		return
	}
}

func TestAdapterErrorLogging(t *testing.T) {
	a := NewMysqlAdapter(`wp_`)

	r, err := regexp.Compile(`\[ERROR\]:.+Hello World`)
	if err != nil {
		t.Errorf(`could not compile regex`)
		return
	}
	var b bytes.Buffer
	wr := bufio.NewWriter(&b)
	a.SetLogs(wr)
	a.LogError(errors.New(`Hello World`))
	wr.Flush()
	if r.MatchString(b.String()) == false {
		t.Errorf(`failed to match info line`)
		return
	}

	var b2 bytes.Buffer
	wr2 := bufio.NewWriter(&b2)
	a.SetLogs(wr2)
	a.SetLogFilter(func(tag string, val string) string {
		return ``
	})
	a.LogError(errors.New(`Hello World`))
	wr2.Flush()
	if b2.String() != `` {
		t.Errorf(`Info should not occur due to filter in this case but equals %s`, b2.String())
		return
	}
}
