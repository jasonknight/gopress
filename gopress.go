package gopress

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

type CommentMeta struct {
	_table    string
	_adapter  Adapter
	_pkey     string // 0 The name of the primary key in this table
	_conds    []string
	_new      bool
	MetaId    int64
	CommentId int64
	MetaKey   string
	MetaValue string
}

func NewCommentMeta(a Adapter) *CommentMeta {
	var o CommentMeta
	o._table = fmt.Sprintf("%scommentmeta", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "meta_id"
	o._new = false
	return &o
}

func (o *CommentMeta) Find(_find_by_MetaId int64) (CommentMeta, error) {

	var model_slice []CommentMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "meta_id", _find_by_MetaId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return CommentMeta{}, err
	}

	for _, result := range results {
		ro := CommentMeta{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return CommentMeta{}, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return CommentMeta{}, errors.New("not found")
	}
	return model_slice[0], nil

}
func (o *CommentMeta) FindByCommentId(_find_by_CommentId int64) ([]CommentMeta, error) {

	var model_slice []CommentMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "comment_id", _find_by_CommentId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := CommentMeta{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *CommentMeta) FindByMetaKey(_find_by_MetaKey string) ([]CommentMeta, error) {

	var model_slice []CommentMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "meta_key", _find_by_MetaKey)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := CommentMeta{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *CommentMeta) FindByMetaValue(_find_by_MetaValue string) ([]CommentMeta, error) {

	var model_slice []CommentMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "meta_value", _find_by_MetaValue)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := CommentMeta{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}

func (o *CommentMeta) FromDBValueMap(m map[string]DBValue) error {
	_MetaId, err := m["meta_id"].AsInt64()
	if err != nil {
		return err
	}
	o.MetaId = _MetaId
	_CommentId, err := m["comment_id"].AsInt64()
	if err != nil {
		return err
	}
	o.CommentId = _CommentId
	_MetaKey, err := m["meta_key"].AsString()
	if err != nil {
		return err
	}
	o.MetaKey = _MetaKey
	_MetaValue, err := m["meta_value"].AsString()
	if err != nil {
		return err
	}
	o.MetaValue = _MetaValue

	return nil
}

func (o *CommentMeta) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `comment_id` = '%d', `meta_key` = '%s', `meta_value` = '%s' WHERE %s = '%d' LIMIT 1", o._table, o.CommentId, o.MetaKey, o.MetaValue, o._pkey, o.MetaId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *CommentMeta) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`comment_id`, `meta_key`, `meta_value`) VALUES ('%d', '%s', '%s')", o._table, o.CommentId, o.MetaKey, o.MetaValue)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *CommentMeta) UpdateCommentId(_upd_CommentId int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_id` = '%d' WHERE `meta_id` = '%d'", o._table, _upd_CommentId, o.CommentId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.CommentId = _upd_CommentId
	return o._adapter.AffectedRows(), nil
}

func (o *CommentMeta) UpdateMetaKey(_upd_MetaKey string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `meta_key` = '%s' WHERE `meta_id` = '%d'", o._table, _upd_MetaKey, o.MetaKey)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.MetaKey = _upd_MetaKey
	return o._adapter.AffectedRows(), nil
}

func (o *CommentMeta) UpdateMetaValue(_upd_MetaValue string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `meta_value` = '%s' WHERE `meta_id` = '%d'", o._table, _upd_MetaValue, o.MetaValue)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.MetaValue = _upd_MetaValue
	return o._adapter.AffectedRows(), nil
}

type Comment struct {
	_table             string
	_adapter           Adapter
	_pkey              string // 0 The name of the primary key in this table
	_conds             []string
	_new               bool
	CommentID          int64
	CommentPostID      int64
	CommentAuthor      string
	CommentAuthorEmail string
	CommentAuthorUrl   string
	CommentAuthorIP    string
	CommentDate        DateTime
	CommentDateGmt     DateTime
	CommentContent     string
	CommentKarma       int
	CommentApproved    string
	CommentAgent       string
	CommentType        string
	CommentParent      int64
	UserId             int64
}

func NewComment(a Adapter) *Comment {
	var o Comment
	o._table = fmt.Sprintf("%scomments", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "comment_ID"
	o._new = false
	return &o
}

func (o *Comment) Find(_find_by_CommentID int64) (Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "comment_ID", _find_by_CommentID)
	results, err := o._adapter.Query(q)
	if err != nil {
		return Comment{}, err
	}

	for _, result := range results {
		ro := Comment{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return Comment{}, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return Comment{}, errors.New("not found")
	}
	return model_slice[0], nil

}
func (o *Comment) FindByCommentPostID(_find_by_CommentPostID int64) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "comment_post_ID", _find_by_CommentPostID)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Comment{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentAuthor(_find_by_CommentAuthor string) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "comment_author", _find_by_CommentAuthor)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Comment{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentAuthorEmail(_find_by_CommentAuthorEmail string) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "comment_author_email", _find_by_CommentAuthorEmail)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Comment{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentAuthorUrl(_find_by_CommentAuthorUrl string) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "comment_author_url", _find_by_CommentAuthorUrl)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Comment{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentAuthorIP(_find_by_CommentAuthorIP string) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "comment_author_IP", _find_by_CommentAuthorIP)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Comment{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentDate(_find_by_CommentDate DateTime) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "comment_date", _find_by_CommentDate)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Comment{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentDateGmt(_find_by_CommentDateGmt DateTime) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "comment_date_gmt", _find_by_CommentDateGmt)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Comment{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentContent(_find_by_CommentContent string) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "comment_content", _find_by_CommentContent)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Comment{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentKarma(_find_by_CommentKarma int) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "comment_karma", _find_by_CommentKarma)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Comment{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentApproved(_find_by_CommentApproved string) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "comment_approved", _find_by_CommentApproved)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Comment{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentAgent(_find_by_CommentAgent string) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "comment_agent", _find_by_CommentAgent)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Comment{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentType(_find_by_CommentType string) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "comment_type", _find_by_CommentType)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Comment{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentParent(_find_by_CommentParent int64) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "comment_parent", _find_by_CommentParent)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Comment{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Comment) FindByUserId(_find_by_UserId int64) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "user_id", _find_by_UserId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Comment{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}

func (o *Comment) FromDBValueMap(m map[string]DBValue) error {
	_CommentID, err := m["comment_ID"].AsInt64()
	if err != nil {
		return err
	}
	o.CommentID = _CommentID
	_CommentPostID, err := m["comment_post_ID"].AsInt64()
	if err != nil {
		return err
	}
	o.CommentPostID = _CommentPostID
	_CommentAuthor, err := m["comment_author"].AsString()
	if err != nil {
		return err
	}
	o.CommentAuthor = _CommentAuthor
	_CommentAuthorEmail, err := m["comment_author_email"].AsString()
	if err != nil {
		return err
	}
	o.CommentAuthorEmail = _CommentAuthorEmail
	_CommentAuthorUrl, err := m["comment_author_url"].AsString()
	if err != nil {
		return err
	}
	o.CommentAuthorUrl = _CommentAuthorUrl
	_CommentAuthorIP, err := m["comment_author_IP"].AsString()
	if err != nil {
		return err
	}
	o.CommentAuthorIP = _CommentAuthorIP
	_CommentDate, err := m["comment_date"].AsDateTime()
	if err != nil {
		return err
	}
	o.CommentDate = _CommentDate
	_CommentDateGmt, err := m["comment_date_gmt"].AsDateTime()
	if err != nil {
		return err
	}
	o.CommentDateGmt = _CommentDateGmt
	_CommentContent, err := m["comment_content"].AsString()
	if err != nil {
		return err
	}
	o.CommentContent = _CommentContent
	_CommentKarma, err := m["comment_karma"].AsInt()
	if err != nil {
		return err
	}
	o.CommentKarma = _CommentKarma
	_CommentApproved, err := m["comment_approved"].AsString()
	if err != nil {
		return err
	}
	o.CommentApproved = _CommentApproved
	_CommentAgent, err := m["comment_agent"].AsString()
	if err != nil {
		return err
	}
	o.CommentAgent = _CommentAgent
	_CommentType, err := m["comment_type"].AsString()
	if err != nil {
		return err
	}
	o.CommentType = _CommentType
	_CommentParent, err := m["comment_parent"].AsInt64()
	if err != nil {
		return err
	}
	o.CommentParent = _CommentParent
	_UserId, err := m["user_id"].AsInt64()
	if err != nil {
		return err
	}
	o.UserId = _UserId

	return nil
}

func (o *Comment) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `comment_post_ID` = '%d', `comment_author` = '%s', `comment_author_email` = '%s', `comment_author_url` = '%s', `comment_author_IP` = '%s', `comment_date` = '%s', `comment_date_gmt` = '%s', `comment_content` = '%s', `comment_karma` = '%d', `comment_approved` = '%s', `comment_agent` = '%s', `comment_type` = '%s', `comment_parent` = '%d', `user_id` = '%d' WHERE %s = '%d' LIMIT 1", o._table, o.CommentPostID, o.CommentAuthor, o.CommentAuthorEmail, o.CommentAuthorUrl, o.CommentAuthorIP, o.CommentDate, o.CommentDateGmt, o.CommentContent, o.CommentKarma, o.CommentApproved, o.CommentAgent, o.CommentType, o.CommentParent, o.UserId, o._pkey, o.CommentID)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *Comment) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`comment_post_ID`, `comment_author`, `comment_author_email`, `comment_author_url`, `comment_author_IP`, `comment_date`, `comment_date_gmt`, `comment_content`, `comment_karma`, `comment_approved`, `comment_agent`, `comment_type`, `comment_parent`, `user_id`) VALUES ('%d', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%d', '%s', '%s', '%s', '%d', '%d')", o._table, o.CommentPostID, o.CommentAuthor, o.CommentAuthorEmail, o.CommentAuthorUrl, o.CommentAuthorIP, o.CommentDate, o.CommentDateGmt, o.CommentContent, o.CommentKarma, o.CommentApproved, o.CommentAgent, o.CommentType, o.CommentParent, o.UserId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentPostID(_upd_CommentPostID int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_post_ID` = '%d' WHERE `comment_ID` = '%d'", o._table, _upd_CommentPostID, o.CommentPostID)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.CommentPostID = _upd_CommentPostID
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentAuthor(_upd_CommentAuthor string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_author` = '%s' WHERE `comment_ID` = '%d'", o._table, _upd_CommentAuthor, o.CommentAuthor)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.CommentAuthor = _upd_CommentAuthor
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentAuthorEmail(_upd_CommentAuthorEmail string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_author_email` = '%s' WHERE `comment_ID` = '%d'", o._table, _upd_CommentAuthorEmail, o.CommentAuthorEmail)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.CommentAuthorEmail = _upd_CommentAuthorEmail
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentAuthorUrl(_upd_CommentAuthorUrl string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_author_url` = '%s' WHERE `comment_ID` = '%d'", o._table, _upd_CommentAuthorUrl, o.CommentAuthorUrl)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.CommentAuthorUrl = _upd_CommentAuthorUrl
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentAuthorIP(_upd_CommentAuthorIP string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_author_IP` = '%s' WHERE `comment_ID` = '%d'", o._table, _upd_CommentAuthorIP, o.CommentAuthorIP)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.CommentAuthorIP = _upd_CommentAuthorIP
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentDate(_upd_CommentDate DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_date` = '%s' WHERE `comment_ID` = '%d'", o._table, _upd_CommentDate, o.CommentDate)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.CommentDate = _upd_CommentDate
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentDateGmt(_upd_CommentDateGmt DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_date_gmt` = '%s' WHERE `comment_ID` = '%d'", o._table, _upd_CommentDateGmt, o.CommentDateGmt)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.CommentDateGmt = _upd_CommentDateGmt
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentContent(_upd_CommentContent string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_content` = '%s' WHERE `comment_ID` = '%d'", o._table, _upd_CommentContent, o.CommentContent)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.CommentContent = _upd_CommentContent
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentKarma(_upd_CommentKarma int) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_karma` = '%d' WHERE `comment_ID` = '%d'", o._table, _upd_CommentKarma, o.CommentKarma)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.CommentKarma = _upd_CommentKarma
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentApproved(_upd_CommentApproved string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_approved` = '%s' WHERE `comment_ID` = '%d'", o._table, _upd_CommentApproved, o.CommentApproved)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.CommentApproved = _upd_CommentApproved
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentAgent(_upd_CommentAgent string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_agent` = '%s' WHERE `comment_ID` = '%d'", o._table, _upd_CommentAgent, o.CommentAgent)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.CommentAgent = _upd_CommentAgent
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentType(_upd_CommentType string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_type` = '%s' WHERE `comment_ID` = '%d'", o._table, _upd_CommentType, o.CommentType)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.CommentType = _upd_CommentType
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentParent(_upd_CommentParent int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_parent` = '%d' WHERE `comment_ID` = '%d'", o._table, _upd_CommentParent, o.CommentParent)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.CommentParent = _upd_CommentParent
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateUserId(_upd_UserId int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_id` = '%d' WHERE `comment_ID` = '%d'", o._table, _upd_UserId, o.UserId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.UserId = _upd_UserId
	return o._adapter.AffectedRows(), nil
}

type Link struct {
	_table          string
	_adapter        Adapter
	_pkey           string // 0 The name of the primary key in this table
	_conds          []string
	_new            bool
	LinkId          int64
	LinkUrl         string
	LinkName        string
	LinkImage       string
	LinkTarget      string
	LinkDescription string
	LinkVisible     string
	LinkOwner       int64
	LinkRating      int
	LinkUpdated     DateTime
	LinkRel         string
	LinkNotes       string
	LinkRss         string
}

func NewLink(a Adapter) *Link {
	var o Link
	o._table = fmt.Sprintf("%slinks", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "link_id"
	o._new = false
	return &o
}

func (o *Link) Find(_find_by_LinkId int64) (Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "link_id", _find_by_LinkId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return Link{}, err
	}

	for _, result := range results {
		ro := Link{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return Link{}, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return Link{}, errors.New("not found")
	}
	return model_slice[0], nil

}
func (o *Link) FindByLinkUrl(_find_by_LinkUrl string) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "link_url", _find_by_LinkUrl)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Link{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Link) FindByLinkName(_find_by_LinkName string) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "link_name", _find_by_LinkName)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Link{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Link) FindByLinkImage(_find_by_LinkImage string) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "link_image", _find_by_LinkImage)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Link{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Link) FindByLinkTarget(_find_by_LinkTarget string) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "link_target", _find_by_LinkTarget)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Link{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Link) FindByLinkDescription(_find_by_LinkDescription string) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "link_description", _find_by_LinkDescription)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Link{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Link) FindByLinkVisible(_find_by_LinkVisible string) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "link_visible", _find_by_LinkVisible)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Link{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Link) FindByLinkOwner(_find_by_LinkOwner int64) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "link_owner", _find_by_LinkOwner)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Link{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Link) FindByLinkRating(_find_by_LinkRating int) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "link_rating", _find_by_LinkRating)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Link{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Link) FindByLinkUpdated(_find_by_LinkUpdated DateTime) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "link_updated", _find_by_LinkUpdated)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Link{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Link) FindByLinkRel(_find_by_LinkRel string) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "link_rel", _find_by_LinkRel)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Link{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Link) FindByLinkNotes(_find_by_LinkNotes string) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "link_notes", _find_by_LinkNotes)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Link{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Link) FindByLinkRss(_find_by_LinkRss string) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "link_rss", _find_by_LinkRss)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Link{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}

func (o *Link) FromDBValueMap(m map[string]DBValue) error {
	_LinkId, err := m["link_id"].AsInt64()
	if err != nil {
		return err
	}
	o.LinkId = _LinkId
	_LinkUrl, err := m["link_url"].AsString()
	if err != nil {
		return err
	}
	o.LinkUrl = _LinkUrl
	_LinkName, err := m["link_name"].AsString()
	if err != nil {
		return err
	}
	o.LinkName = _LinkName
	_LinkImage, err := m["link_image"].AsString()
	if err != nil {
		return err
	}
	o.LinkImage = _LinkImage
	_LinkTarget, err := m["link_target"].AsString()
	if err != nil {
		return err
	}
	o.LinkTarget = _LinkTarget
	_LinkDescription, err := m["link_description"].AsString()
	if err != nil {
		return err
	}
	o.LinkDescription = _LinkDescription
	_LinkVisible, err := m["link_visible"].AsString()
	if err != nil {
		return err
	}
	o.LinkVisible = _LinkVisible
	_LinkOwner, err := m["link_owner"].AsInt64()
	if err != nil {
		return err
	}
	o.LinkOwner = _LinkOwner
	_LinkRating, err := m["link_rating"].AsInt()
	if err != nil {
		return err
	}
	o.LinkRating = _LinkRating
	_LinkUpdated, err := m["link_updated"].AsDateTime()
	if err != nil {
		return err
	}
	o.LinkUpdated = _LinkUpdated
	_LinkRel, err := m["link_rel"].AsString()
	if err != nil {
		return err
	}
	o.LinkRel = _LinkRel
	_LinkNotes, err := m["link_notes"].AsString()
	if err != nil {
		return err
	}
	o.LinkNotes = _LinkNotes
	_LinkRss, err := m["link_rss"].AsString()
	if err != nil {
		return err
	}
	o.LinkRss = _LinkRss

	return nil
}

func (o *Link) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `link_url` = '%s', `link_name` = '%s', `link_image` = '%s', `link_target` = '%s', `link_description` = '%s', `link_visible` = '%s', `link_owner` = '%d', `link_rating` = '%d', `link_updated` = '%s', `link_rel` = '%s', `link_notes` = '%s', `link_rss` = '%s' WHERE %s = '%d' LIMIT 1", o._table, o.LinkUrl, o.LinkName, o.LinkImage, o.LinkTarget, o.LinkDescription, o.LinkVisible, o.LinkOwner, o.LinkRating, o.LinkUpdated, o.LinkRel, o.LinkNotes, o.LinkRss, o._pkey, o.LinkId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *Link) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`link_url`, `link_name`, `link_image`, `link_target`, `link_description`, `link_visible`, `link_owner`, `link_rating`, `link_updated`, `link_rel`, `link_notes`, `link_rss`) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%d', '%d', '%s', '%s', '%s', '%s')", o._table, o.LinkUrl, o.LinkName, o.LinkImage, o.LinkTarget, o.LinkDescription, o.LinkVisible, o.LinkOwner, o.LinkRating, o.LinkUpdated, o.LinkRel, o.LinkNotes, o.LinkRss)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkUrl(_upd_LinkUrl string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_url` = '%s' WHERE `link_id` = '%d'", o._table, _upd_LinkUrl, o.LinkUrl)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.LinkUrl = _upd_LinkUrl
	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkName(_upd_LinkName string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_name` = '%s' WHERE `link_id` = '%d'", o._table, _upd_LinkName, o.LinkName)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.LinkName = _upd_LinkName
	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkImage(_upd_LinkImage string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_image` = '%s' WHERE `link_id` = '%d'", o._table, _upd_LinkImage, o.LinkImage)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.LinkImage = _upd_LinkImage
	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkTarget(_upd_LinkTarget string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_target` = '%s' WHERE `link_id` = '%d'", o._table, _upd_LinkTarget, o.LinkTarget)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.LinkTarget = _upd_LinkTarget
	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkDescription(_upd_LinkDescription string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_description` = '%s' WHERE `link_id` = '%d'", o._table, _upd_LinkDescription, o.LinkDescription)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.LinkDescription = _upd_LinkDescription
	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkVisible(_upd_LinkVisible string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_visible` = '%s' WHERE `link_id` = '%d'", o._table, _upd_LinkVisible, o.LinkVisible)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.LinkVisible = _upd_LinkVisible
	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkOwner(_upd_LinkOwner int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_owner` = '%d' WHERE `link_id` = '%d'", o._table, _upd_LinkOwner, o.LinkOwner)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.LinkOwner = _upd_LinkOwner
	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkRating(_upd_LinkRating int) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_rating` = '%d' WHERE `link_id` = '%d'", o._table, _upd_LinkRating, o.LinkRating)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.LinkRating = _upd_LinkRating
	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkUpdated(_upd_LinkUpdated DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_updated` = '%s' WHERE `link_id` = '%d'", o._table, _upd_LinkUpdated, o.LinkUpdated)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.LinkUpdated = _upd_LinkUpdated
	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkRel(_upd_LinkRel string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_rel` = '%s' WHERE `link_id` = '%d'", o._table, _upd_LinkRel, o.LinkRel)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.LinkRel = _upd_LinkRel
	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkNotes(_upd_LinkNotes string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_notes` = '%s' WHERE `link_id` = '%d'", o._table, _upd_LinkNotes, o.LinkNotes)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.LinkNotes = _upd_LinkNotes
	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkRss(_upd_LinkRss string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_rss` = '%s' WHERE `link_id` = '%d'", o._table, _upd_LinkRss, o.LinkRss)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.LinkRss = _upd_LinkRss
	return o._adapter.AffectedRows(), nil
}

type Option struct {
	_table      string
	_adapter    Adapter
	_pkey       string // 0 The name of the primary key in this table
	_conds      []string
	_new        bool
	OptionId    int64
	OptionName  string
	OptionValue string
	Autoload    string
}

func NewOption(a Adapter) *Option {
	var o Option
	o._table = fmt.Sprintf("%soptions", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "option_id"
	o._new = false
	return &o
}

func (o *Option) Find(_find_by_OptionId int64) (Option, error) {

	var model_slice []Option
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "option_id", _find_by_OptionId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return Option{}, err
	}

	for _, result := range results {
		ro := Option{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return Option{}, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return Option{}, errors.New("not found")
	}
	return model_slice[0], nil

}
func (o *Option) FindByOptionName(_find_by_OptionName string) ([]Option, error) {

	var model_slice []Option
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "option_name", _find_by_OptionName)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Option{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Option) FindByOptionValue(_find_by_OptionValue string) ([]Option, error) {

	var model_slice []Option
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "option_value", _find_by_OptionValue)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Option{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Option) FindByAutoload(_find_by_Autoload string) ([]Option, error) {

	var model_slice []Option
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "autoload", _find_by_Autoload)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Option{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}

func (o *Option) FromDBValueMap(m map[string]DBValue) error {
	_OptionId, err := m["option_id"].AsInt64()
	if err != nil {
		return err
	}
	o.OptionId = _OptionId
	_OptionName, err := m["option_name"].AsString()
	if err != nil {
		return err
	}
	o.OptionName = _OptionName
	_OptionValue, err := m["option_value"].AsString()
	if err != nil {
		return err
	}
	o.OptionValue = _OptionValue
	_Autoload, err := m["autoload"].AsString()
	if err != nil {
		return err
	}
	o.Autoload = _Autoload

	return nil
}

func (o *Option) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `option_name` = '%s', `option_value` = '%s', `autoload` = '%s' WHERE %s = '%d' LIMIT 1", o._table, o.OptionName, o.OptionValue, o.Autoload, o._pkey, o.OptionId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *Option) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`option_name`, `option_value`, `autoload`) VALUES ('%s', '%s', '%s')", o._table, o.OptionName, o.OptionValue, o.Autoload)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *Option) UpdateOptionName(_upd_OptionName string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `option_name` = '%s' WHERE `option_id` = '%d'", o._table, _upd_OptionName, o.OptionName)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.OptionName = _upd_OptionName
	return o._adapter.AffectedRows(), nil
}

func (o *Option) UpdateOptionValue(_upd_OptionValue string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `option_value` = '%s' WHERE `option_id` = '%d'", o._table, _upd_OptionValue, o.OptionValue)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.OptionValue = _upd_OptionValue
	return o._adapter.AffectedRows(), nil
}

func (o *Option) UpdateAutoload(_upd_Autoload string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `autoload` = '%s' WHERE `option_id` = '%d'", o._table, _upd_Autoload, o.Autoload)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.Autoload = _upd_Autoload
	return o._adapter.AffectedRows(), nil
}

type PostMeta struct {
	_table    string
	_adapter  Adapter
	_pkey     string // 0 The name of the primary key in this table
	_conds    []string
	_new      bool
	MetaId    int64
	Id        int64
	MetaKey   string
	MetaValue string
}

func NewPostMeta(a Adapter) *PostMeta {
	var o PostMeta
	o._table = fmt.Sprintf("%spostmeta", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "meta_id"
	o._new = false
	return &o
}

func (o *PostMeta) Find(_find_by_MetaId int64) (PostMeta, error) {

	var model_slice []PostMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "meta_id", _find_by_MetaId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return PostMeta{}, err
	}

	for _, result := range results {
		ro := PostMeta{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return PostMeta{}, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return PostMeta{}, errors.New("not found")
	}
	return model_slice[0], nil

}
func (o *PostMeta) FindById(_find_by_Id int64) ([]PostMeta, error) {

	var model_slice []PostMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "post_id", _find_by_Id)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := PostMeta{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *PostMeta) FindByMetaKey(_find_by_MetaKey string) ([]PostMeta, error) {

	var model_slice []PostMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "meta_key", _find_by_MetaKey)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := PostMeta{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *PostMeta) FindByMetaValue(_find_by_MetaValue string) ([]PostMeta, error) {

	var model_slice []PostMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "meta_value", _find_by_MetaValue)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := PostMeta{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}

func (o *PostMeta) FromDBValueMap(m map[string]DBValue) error {
	_MetaId, err := m["meta_id"].AsInt64()
	if err != nil {
		return err
	}
	o.MetaId = _MetaId
	_Id, err := m["post_id"].AsInt64()
	if err != nil {
		return err
	}
	o.Id = _Id
	_MetaKey, err := m["meta_key"].AsString()
	if err != nil {
		return err
	}
	o.MetaKey = _MetaKey
	_MetaValue, err := m["meta_value"].AsString()
	if err != nil {
		return err
	}
	o.MetaValue = _MetaValue

	return nil
}

func (o *PostMeta) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `post_id` = '%d', `meta_key` = '%s', `meta_value` = '%s' WHERE %s = '%d' LIMIT 1", o._table, o.Id, o.MetaKey, o.MetaValue, o._pkey, o.MetaId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *PostMeta) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`post_id`, `meta_key`, `meta_value`) VALUES ('%d', '%s', '%s')", o._table, o.Id, o.MetaKey, o.MetaValue)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *PostMeta) UpdateId(_upd_Id int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_id` = '%d' WHERE `meta_id` = '%d'", o._table, _upd_Id, o.Id)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.Id = _upd_Id
	return o._adapter.AffectedRows(), nil
}

func (o *PostMeta) UpdateMetaKey(_upd_MetaKey string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `meta_key` = '%s' WHERE `meta_id` = '%d'", o._table, _upd_MetaKey, o.MetaKey)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.MetaKey = _upd_MetaKey
	return o._adapter.AffectedRows(), nil
}

func (o *PostMeta) UpdateMetaValue(_upd_MetaValue string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `meta_value` = '%s' WHERE `meta_id` = '%d'", o._table, _upd_MetaValue, o.MetaValue)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.MetaValue = _upd_MetaValue
	return o._adapter.AffectedRows(), nil
}

type Post struct {
	_table          string
	_adapter        Adapter
	_pkey           string // 0 The name of the primary key in this table
	_conds          []string
	_new            bool
	ID              int64
	Author          int64
	Date            DateTime
	DateGmt         DateTime
	Content         string
	Title           string
	Excerpt         string
	Status          string
	CommentStatus   string
	PingStatus      string
	Password        string
	Name            string
	ToPing          string
	Pinged          string
	Modified        DateTime
	ModifiedGmt     DateTime
	ContentFiltered string
	Parent          int64
	Guid            string
	MenuOrder       int
	Type            string
	MimeType        string
	CommentCount    int64
}

func NewPost(a Adapter) *Post {
	var o Post
	o._table = fmt.Sprintf("%sposts", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "ID"
	o._new = false
	return &o
}

func (o *Post) Find(_find_by_ID int64) (Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "ID", _find_by_ID)
	results, err := o._adapter.Query(q)
	if err != nil {
		return Post{}, err
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return Post{}, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return Post{}, errors.New("not found")
	}
	return model_slice[0], nil

}
func (o *Post) FindByAuthor(_find_by_Author int64) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "post_author", _find_by_Author)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Post) FindByDate(_find_by_Date DateTime) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_date", _find_by_Date)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Post) FindByDateGmt(_find_by_DateGmt DateTime) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_date_gmt", _find_by_DateGmt)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Post) FindByContent(_find_by_Content string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_content", _find_by_Content)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Post) FindByTitle(_find_by_Title string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_title", _find_by_Title)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Post) FindByExcerpt(_find_by_Excerpt string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_excerpt", _find_by_Excerpt)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Post) FindByStatus(_find_by_Status string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_status", _find_by_Status)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Post) FindByCommentStatus(_find_by_CommentStatus string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "comment_status", _find_by_CommentStatus)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Post) FindByPingStatus(_find_by_PingStatus string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "ping_status", _find_by_PingStatus)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Post) FindByPassword(_find_by_Password string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_password", _find_by_Password)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Post) FindByName(_find_by_Name string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_name", _find_by_Name)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Post) FindByToPing(_find_by_ToPing string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "to_ping", _find_by_ToPing)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Post) FindByPinged(_find_by_Pinged string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "pinged", _find_by_Pinged)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Post) FindByModified(_find_by_Modified DateTime) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_modified", _find_by_Modified)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Post) FindByModifiedGmt(_find_by_ModifiedGmt DateTime) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_modified_gmt", _find_by_ModifiedGmt)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Post) FindByContentFiltered(_find_by_ContentFiltered string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_content_filtered", _find_by_ContentFiltered)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Post) FindByParent(_find_by_Parent int64) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "post_parent", _find_by_Parent)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Post) FindByGuid(_find_by_Guid string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "guid", _find_by_Guid)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Post) FindByMenuOrder(_find_by_MenuOrder int) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "menu_order", _find_by_MenuOrder)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Post) FindByType(_find_by_Type string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_type", _find_by_Type)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Post) FindByMimeType(_find_by_MimeType string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_mime_type", _find_by_MimeType)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Post) FindByCommentCount(_find_by_CommentCount int64) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "comment_count", _find_by_CommentCount)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}

func (o *Post) FromDBValueMap(m map[string]DBValue) error {
	_ID, err := m["ID"].AsInt64()
	if err != nil {
		return err
	}
	o.ID = _ID
	_Author, err := m["post_author"].AsInt64()
	if err != nil {
		return err
	}
	o.Author = _Author
	_Date, err := m["post_date"].AsDateTime()
	if err != nil {
		return err
	}
	o.Date = _Date
	_DateGmt, err := m["post_date_gmt"].AsDateTime()
	if err != nil {
		return err
	}
	o.DateGmt = _DateGmt
	_Content, err := m["post_content"].AsString()
	if err != nil {
		return err
	}
	o.Content = _Content
	_Title, err := m["post_title"].AsString()
	if err != nil {
		return err
	}
	o.Title = _Title
	_Excerpt, err := m["post_excerpt"].AsString()
	if err != nil {
		return err
	}
	o.Excerpt = _Excerpt
	_Status, err := m["post_status"].AsString()
	if err != nil {
		return err
	}
	o.Status = _Status
	_CommentStatus, err := m["comment_status"].AsString()
	if err != nil {
		return err
	}
	o.CommentStatus = _CommentStatus
	_PingStatus, err := m["ping_status"].AsString()
	if err != nil {
		return err
	}
	o.PingStatus = _PingStatus
	_Password, err := m["post_password"].AsString()
	if err != nil {
		return err
	}
	o.Password = _Password
	_Name, err := m["post_name"].AsString()
	if err != nil {
		return err
	}
	o.Name = _Name
	_ToPing, err := m["to_ping"].AsString()
	if err != nil {
		return err
	}
	o.ToPing = _ToPing
	_Pinged, err := m["pinged"].AsString()
	if err != nil {
		return err
	}
	o.Pinged = _Pinged
	_Modified, err := m["post_modified"].AsDateTime()
	if err != nil {
		return err
	}
	o.Modified = _Modified
	_ModifiedGmt, err := m["post_modified_gmt"].AsDateTime()
	if err != nil {
		return err
	}
	o.ModifiedGmt = _ModifiedGmt
	_ContentFiltered, err := m["post_content_filtered"].AsString()
	if err != nil {
		return err
	}
	o.ContentFiltered = _ContentFiltered
	_Parent, err := m["post_parent"].AsInt64()
	if err != nil {
		return err
	}
	o.Parent = _Parent
	_Guid, err := m["guid"].AsString()
	if err != nil {
		return err
	}
	o.Guid = _Guid
	_MenuOrder, err := m["menu_order"].AsInt()
	if err != nil {
		return err
	}
	o.MenuOrder = _MenuOrder
	_Type, err := m["post_type"].AsString()
	if err != nil {
		return err
	}
	o.Type = _Type
	_MimeType, err := m["post_mime_type"].AsString()
	if err != nil {
		return err
	}
	o.MimeType = _MimeType
	_CommentCount, err := m["comment_count"].AsInt64()
	if err != nil {
		return err
	}
	o.CommentCount = _CommentCount

	return nil
}

func (o *Post) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `post_author` = '%d', `post_date` = '%s', `post_date_gmt` = '%s', `post_content` = '%s', `post_title` = '%s', `post_excerpt` = '%s', `post_status` = '%s', `comment_status` = '%s', `ping_status` = '%s', `post_password` = '%s', `post_name` = '%s', `to_ping` = '%s', `pinged` = '%s', `post_modified` = '%s', `post_modified_gmt` = '%s', `post_content_filtered` = '%s', `post_parent` = '%d', `guid` = '%s', `menu_order` = '%d', `post_type` = '%s', `post_mime_type` = '%s', `comment_count` = '%d' WHERE %s = '%d' LIMIT 1", o._table, o.Author, o.Date, o.DateGmt, o.Content, o.Title, o.Excerpt, o.Status, o.CommentStatus, o.PingStatus, o.Password, o.Name, o.ToPing, o.Pinged, o.Modified, o.ModifiedGmt, o.ContentFiltered, o.Parent, o.Guid, o.MenuOrder, o.Type, o.MimeType, o.CommentCount, o._pkey, o.ID)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *Post) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`post_author`, `post_date`, `post_date_gmt`, `post_content`, `post_title`, `post_excerpt`, `post_status`, `comment_status`, `ping_status`, `post_password`, `post_name`, `to_ping`, `pinged`, `post_modified`, `post_modified_gmt`, `post_content_filtered`, `post_parent`, `guid`, `menu_order`, `post_type`, `post_mime_type`, `comment_count`) VALUES ('%d', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%d', '%s', '%d', '%s', '%s', '%d')", o._table, o.Author, o.Date, o.DateGmt, o.Content, o.Title, o.Excerpt, o.Status, o.CommentStatus, o.PingStatus, o.Password, o.Name, o.ToPing, o.Pinged, o.Modified, o.ModifiedGmt, o.ContentFiltered, o.Parent, o.Guid, o.MenuOrder, o.Type, o.MimeType, o.CommentCount)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateAuthor(_upd_Author int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_author` = '%d' WHERE `ID` = '%d'", o._table, _upd_Author, o.Author)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.Author = _upd_Author
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateDate(_upd_Date DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_date` = '%s' WHERE `ID` = '%d'", o._table, _upd_Date, o.Date)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.Date = _upd_Date
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateDateGmt(_upd_DateGmt DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_date_gmt` = '%s' WHERE `ID` = '%d'", o._table, _upd_DateGmt, o.DateGmt)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.DateGmt = _upd_DateGmt
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateContent(_upd_Content string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_content` = '%s' WHERE `ID` = '%d'", o._table, _upd_Content, o.Content)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.Content = _upd_Content
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateTitle(_upd_Title string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_title` = '%s' WHERE `ID` = '%d'", o._table, _upd_Title, o.Title)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.Title = _upd_Title
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateExcerpt(_upd_Excerpt string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_excerpt` = '%s' WHERE `ID` = '%d'", o._table, _upd_Excerpt, o.Excerpt)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.Excerpt = _upd_Excerpt
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateStatus(_upd_Status string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_status` = '%s' WHERE `ID` = '%d'", o._table, _upd_Status, o.Status)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.Status = _upd_Status
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateCommentStatus(_upd_CommentStatus string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_status` = '%s' WHERE `ID` = '%d'", o._table, _upd_CommentStatus, o.CommentStatus)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.CommentStatus = _upd_CommentStatus
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdatePingStatus(_upd_PingStatus string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `ping_status` = '%s' WHERE `ID` = '%d'", o._table, _upd_PingStatus, o.PingStatus)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.PingStatus = _upd_PingStatus
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdatePassword(_upd_Password string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_password` = '%s' WHERE `ID` = '%d'", o._table, _upd_Password, o.Password)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.Password = _upd_Password
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateName(_upd_Name string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_name` = '%s' WHERE `ID` = '%d'", o._table, _upd_Name, o.Name)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.Name = _upd_Name
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateToPing(_upd_ToPing string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `to_ping` = '%s' WHERE `ID` = '%d'", o._table, _upd_ToPing, o.ToPing)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.ToPing = _upd_ToPing
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdatePinged(_upd_Pinged string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `pinged` = '%s' WHERE `ID` = '%d'", o._table, _upd_Pinged, o.Pinged)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.Pinged = _upd_Pinged
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateModified(_upd_Modified DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_modified` = '%s' WHERE `ID` = '%d'", o._table, _upd_Modified, o.Modified)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.Modified = _upd_Modified
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateModifiedGmt(_upd_ModifiedGmt DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_modified_gmt` = '%s' WHERE `ID` = '%d'", o._table, _upd_ModifiedGmt, o.ModifiedGmt)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.ModifiedGmt = _upd_ModifiedGmt
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateContentFiltered(_upd_ContentFiltered string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_content_filtered` = '%s' WHERE `ID` = '%d'", o._table, _upd_ContentFiltered, o.ContentFiltered)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.ContentFiltered = _upd_ContentFiltered
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateParent(_upd_Parent int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_parent` = '%d' WHERE `ID` = '%d'", o._table, _upd_Parent, o.Parent)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.Parent = _upd_Parent
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateGuid(_upd_Guid string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `guid` = '%s' WHERE `ID` = '%d'", o._table, _upd_Guid, o.Guid)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.Guid = _upd_Guid
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateMenuOrder(_upd_MenuOrder int) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `menu_order` = '%d' WHERE `ID` = '%d'", o._table, _upd_MenuOrder, o.MenuOrder)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.MenuOrder = _upd_MenuOrder
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateType(_upd_Type string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_type` = '%s' WHERE `ID` = '%d'", o._table, _upd_Type, o.Type)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.Type = _upd_Type
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateMimeType(_upd_MimeType string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_mime_type` = '%s' WHERE `ID` = '%d'", o._table, _upd_MimeType, o.MimeType)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.MimeType = _upd_MimeType
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateCommentCount(_upd_CommentCount int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_count` = '%d' WHERE `ID` = '%d'", o._table, _upd_CommentCount, o.CommentCount)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.CommentCount = _upd_CommentCount
	return o._adapter.AffectedRows(), nil
}

type TermRelationship struct {
	_table         string
	_adapter       Adapter
	_pkey          string // 0 The name of the primary key in this table
	_conds         []string
	_new           bool
	ObjectId       int64
	TermTaxonomyId int64
	TermOrder      int
}

func NewTermRelationship(a Adapter) *TermRelationship {
	var o TermRelationship
	o._table = fmt.Sprintf("%sterm_relationships", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "term_taxonomy_id"
	o._new = false
	return &o
}

func (o *TermRelationship) FindByObjectId(_find_by_ObjectId int64) ([]TermRelationship, error) {

	var model_slice []TermRelationship
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "object_id", _find_by_ObjectId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := TermRelationship{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *TermRelationship) Find(_find_by_TermTaxonomyId int64) (TermRelationship, error) {

	var model_slice []TermRelationship
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "term_taxonomy_id", _find_by_TermTaxonomyId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return TermRelationship{}, err
	}

	for _, result := range results {
		ro := TermRelationship{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return TermRelationship{}, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return TermRelationship{}, errors.New("not found")
	}
	return model_slice[0], nil

}
func (o *TermRelationship) FindByTermOrder(_find_by_TermOrder int) ([]TermRelationship, error) {

	var model_slice []TermRelationship
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "term_order", _find_by_TermOrder)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := TermRelationship{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}

func (o *TermRelationship) FromDBValueMap(m map[string]DBValue) error {
	_ObjectId, err := m["object_id"].AsInt64()
	if err != nil {
		return err
	}
	o.ObjectId = _ObjectId
	_TermTaxonomyId, err := m["term_taxonomy_id"].AsInt64()
	if err != nil {
		return err
	}
	o.TermTaxonomyId = _TermTaxonomyId
	_TermOrder, err := m["term_order"].AsInt()
	if err != nil {
		return err
	}
	o.TermOrder = _TermOrder

	return nil
}

func (o *TermRelationship) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `term_order` = '%d' WHERE %s = '%d' LIMIT 1", o._table, o.TermOrder, o._pkey, o.TermTaxonomyId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *TermRelationship) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`term_order`) VALUES ('%d')", o._table, o.TermOrder)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *TermRelationship) UpdateTermOrder(_upd_TermOrder int) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `term_order` = '%d' WHERE `term_taxonomy_id` = '%d'", o._table, _upd_TermOrder, o.TermOrder)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.TermOrder = _upd_TermOrder
	return o._adapter.AffectedRows(), nil
}

type TermTaxonomy struct {
	_table         string
	_adapter       Adapter
	_pkey          string // 0 The name of the primary key in this table
	_conds         []string
	_new           bool
	TermTaxonomyId int64
	TermId         int64
	Taxonomy       string
	Description    string
	Parent         int64
	Count          int64
}

func NewTermTaxonomy(a Adapter) *TermTaxonomy {
	var o TermTaxonomy
	o._table = fmt.Sprintf("%sterm_taxonomy", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "term_taxonomy_id"
	o._new = false
	return &o
}

func (o *TermTaxonomy) Find(_find_by_TermTaxonomyId int64) (TermTaxonomy, error) {

	var model_slice []TermTaxonomy
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "term_taxonomy_id", _find_by_TermTaxonomyId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return TermTaxonomy{}, err
	}

	for _, result := range results {
		ro := TermTaxonomy{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return TermTaxonomy{}, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return TermTaxonomy{}, errors.New("not found")
	}
	return model_slice[0], nil

}
func (o *TermTaxonomy) FindByTermId(_find_by_TermId int64) ([]TermTaxonomy, error) {

	var model_slice []TermTaxonomy
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "term_id", _find_by_TermId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := TermTaxonomy{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *TermTaxonomy) FindByTaxonomy(_find_by_Taxonomy string) ([]TermTaxonomy, error) {

	var model_slice []TermTaxonomy
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "taxonomy", _find_by_Taxonomy)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := TermTaxonomy{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *TermTaxonomy) FindByDescription(_find_by_Description string) ([]TermTaxonomy, error) {

	var model_slice []TermTaxonomy
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "description", _find_by_Description)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := TermTaxonomy{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *TermTaxonomy) FindByParent(_find_by_Parent int64) ([]TermTaxonomy, error) {

	var model_slice []TermTaxonomy
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "parent", _find_by_Parent)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := TermTaxonomy{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *TermTaxonomy) FindByCount(_find_by_Count int64) ([]TermTaxonomy, error) {

	var model_slice []TermTaxonomy
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "count", _find_by_Count)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := TermTaxonomy{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}

func (o *TermTaxonomy) FromDBValueMap(m map[string]DBValue) error {
	_TermTaxonomyId, err := m["term_taxonomy_id"].AsInt64()
	if err != nil {
		return err
	}
	o.TermTaxonomyId = _TermTaxonomyId
	_TermId, err := m["term_id"].AsInt64()
	if err != nil {
		return err
	}
	o.TermId = _TermId
	_Taxonomy, err := m["taxonomy"].AsString()
	if err != nil {
		return err
	}
	o.Taxonomy = _Taxonomy
	_Description, err := m["description"].AsString()
	if err != nil {
		return err
	}
	o.Description = _Description
	_Parent, err := m["parent"].AsInt64()
	if err != nil {
		return err
	}
	o.Parent = _Parent
	_Count, err := m["count"].AsInt64()
	if err != nil {
		return err
	}
	o.Count = _Count

	return nil
}

func (o *TermTaxonomy) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `term_id` = '%d', `taxonomy` = '%s', `description` = '%s', `parent` = '%d', `count` = '%d' WHERE %s = '%d' LIMIT 1", o._table, o.TermId, o.Taxonomy, o.Description, o.Parent, o.Count, o._pkey, o.TermTaxonomyId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *TermTaxonomy) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`term_id`, `taxonomy`, `description`, `parent`, `count`) VALUES ('%d', '%s', '%s', '%d', '%d')", o._table, o.TermId, o.Taxonomy, o.Description, o.Parent, o.Count)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *TermTaxonomy) UpdateTermId(_upd_TermId int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `term_id` = '%d' WHERE `term_taxonomy_id` = '%d'", o._table, _upd_TermId, o.TermId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.TermId = _upd_TermId
	return o._adapter.AffectedRows(), nil
}

func (o *TermTaxonomy) UpdateTaxonomy(_upd_Taxonomy string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `taxonomy` = '%s' WHERE `term_taxonomy_id` = '%d'", o._table, _upd_Taxonomy, o.Taxonomy)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.Taxonomy = _upd_Taxonomy
	return o._adapter.AffectedRows(), nil
}

func (o *TermTaxonomy) UpdateDescription(_upd_Description string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `description` = '%s' WHERE `term_taxonomy_id` = '%d'", o._table, _upd_Description, o.Description)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.Description = _upd_Description
	return o._adapter.AffectedRows(), nil
}

func (o *TermTaxonomy) UpdateParent(_upd_Parent int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `parent` = '%d' WHERE `term_taxonomy_id` = '%d'", o._table, _upd_Parent, o.Parent)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.Parent = _upd_Parent
	return o._adapter.AffectedRows(), nil
}

func (o *TermTaxonomy) UpdateCount(_upd_Count int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `count` = '%d' WHERE `term_taxonomy_id` = '%d'", o._table, _upd_Count, o.Count)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.Count = _upd_Count
	return o._adapter.AffectedRows(), nil
}

type Term struct {
	_table    string
	_adapter  Adapter
	_pkey     string // 0 The name of the primary key in this table
	_conds    []string
	_new      bool
	TermId    int64
	Name      string
	Slug      string
	TermGroup int64
}

func NewTerm(a Adapter) *Term {
	var o Term
	o._table = fmt.Sprintf("%sterms", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "term_id"
	o._new = false
	return &o
}

func (o *Term) Find(_find_by_TermId int64) (Term, error) {

	var model_slice []Term
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "term_id", _find_by_TermId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return Term{}, err
	}

	for _, result := range results {
		ro := Term{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return Term{}, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return Term{}, errors.New("not found")
	}
	return model_slice[0], nil

}
func (o *Term) FindByName(_find_by_Name string) ([]Term, error) {

	var model_slice []Term
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "name", _find_by_Name)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Term{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Term) FindBySlug(_find_by_Slug string) ([]Term, error) {

	var model_slice []Term
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "slug", _find_by_Slug)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Term{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *Term) FindByTermGroup(_find_by_TermGroup int64) ([]Term, error) {

	var model_slice []Term
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "term_group", _find_by_TermGroup)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := Term{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}

func (o *Term) FromDBValueMap(m map[string]DBValue) error {
	_TermId, err := m["term_id"].AsInt64()
	if err != nil {
		return err
	}
	o.TermId = _TermId
	_Name, err := m["name"].AsString()
	if err != nil {
		return err
	}
	o.Name = _Name
	_Slug, err := m["slug"].AsString()
	if err != nil {
		return err
	}
	o.Slug = _Slug
	_TermGroup, err := m["term_group"].AsInt64()
	if err != nil {
		return err
	}
	o.TermGroup = _TermGroup

	return nil
}

func (o *Term) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `name` = '%s', `slug` = '%s', `term_group` = '%d' WHERE %s = '%d' LIMIT 1", o._table, o.Name, o.Slug, o.TermGroup, o._pkey, o.TermId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *Term) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`name`, `slug`, `term_group`) VALUES ('%s', '%s', '%d')", o._table, o.Name, o.Slug, o.TermGroup)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *Term) UpdateName(_upd_Name string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `name` = '%s' WHERE `term_id` = '%d'", o._table, _upd_Name, o.Name)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.Name = _upd_Name
	return o._adapter.AffectedRows(), nil
}

func (o *Term) UpdateSlug(_upd_Slug string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `slug` = '%s' WHERE `term_id` = '%d'", o._table, _upd_Slug, o.Slug)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.Slug = _upd_Slug
	return o._adapter.AffectedRows(), nil
}

func (o *Term) UpdateTermGroup(_upd_TermGroup int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `term_group` = '%d' WHERE `term_id` = '%d'", o._table, _upd_TermGroup, o.TermGroup)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.TermGroup = _upd_TermGroup
	return o._adapter.AffectedRows(), nil
}

type UserMeta struct {
	_table    string
	_adapter  Adapter
	_pkey     string // 0 The name of the primary key in this table
	_conds    []string
	_new      bool
	UMetaId   int64
	UserId    int64
	MetaKey   string
	MetaValue string
}

func NewUserMeta(a Adapter) *UserMeta {
	var o UserMeta
	o._table = fmt.Sprintf("%susermeta", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "umeta_id"
	o._new = false
	return &o
}

func (o *UserMeta) Find(_find_by_UMetaId int64) (UserMeta, error) {

	var model_slice []UserMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "umeta_id", _find_by_UMetaId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return UserMeta{}, err
	}

	for _, result := range results {
		ro := UserMeta{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return UserMeta{}, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return UserMeta{}, errors.New("not found")
	}
	return model_slice[0], nil

}
func (o *UserMeta) FindByUserId(_find_by_UserId int64) ([]UserMeta, error) {

	var model_slice []UserMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "user_id", _find_by_UserId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := UserMeta{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *UserMeta) FindByMetaKey(_find_by_MetaKey string) ([]UserMeta, error) {

	var model_slice []UserMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "meta_key", _find_by_MetaKey)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := UserMeta{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *UserMeta) FindByMetaValue(_find_by_MetaValue string) ([]UserMeta, error) {

	var model_slice []UserMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "meta_value", _find_by_MetaValue)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := UserMeta{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}

func (o *UserMeta) FromDBValueMap(m map[string]DBValue) error {
	_UMetaId, err := m["umeta_id"].AsInt64()
	if err != nil {
		return err
	}
	o.UMetaId = _UMetaId
	_UserId, err := m["user_id"].AsInt64()
	if err != nil {
		return err
	}
	o.UserId = _UserId
	_MetaKey, err := m["meta_key"].AsString()
	if err != nil {
		return err
	}
	o.MetaKey = _MetaKey
	_MetaValue, err := m["meta_value"].AsString()
	if err != nil {
		return err
	}
	o.MetaValue = _MetaValue

	return nil
}

func (o *UserMeta) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `user_id` = '%d', `meta_key` = '%s', `meta_value` = '%s' WHERE %s = '%d' LIMIT 1", o._table, o.UserId, o.MetaKey, o.MetaValue, o._pkey, o.UMetaId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *UserMeta) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`user_id`, `meta_key`, `meta_value`) VALUES ('%d', '%s', '%s')", o._table, o.UserId, o.MetaKey, o.MetaValue)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *UserMeta) UpdateUserId(_upd_UserId int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_id` = '%d' WHERE `umeta_id` = '%d'", o._table, _upd_UserId, o.UserId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.UserId = _upd_UserId
	return o._adapter.AffectedRows(), nil
}

func (o *UserMeta) UpdateMetaKey(_upd_MetaKey string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `meta_key` = '%s' WHERE `umeta_id` = '%d'", o._table, _upd_MetaKey, o.MetaKey)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.MetaKey = _upd_MetaKey
	return o._adapter.AffectedRows(), nil
}

func (o *UserMeta) UpdateMetaValue(_upd_MetaValue string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `meta_value` = '%s' WHERE `umeta_id` = '%d'", o._table, _upd_MetaValue, o.MetaValue)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.MetaValue = _upd_MetaValue
	return o._adapter.AffectedRows(), nil
}

type User struct {
	_table            string
	_adapter          Adapter
	_pkey             string // 0 The name of the primary key in this table
	_conds            []string
	_new              bool
	ID                int64
	UserLogin         string
	UserPass          string
	UserNicename      string
	UserEmail         string
	UserUrl           string
	UserRegistered    DateTime
	UserActivationKey string
	UserStatus        int
	DisplayName       string
}

func NewUser(a Adapter) *User {
	var o User
	o._table = fmt.Sprintf("%susers", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "ID"
	o._new = false
	return &o
}

func (o *User) Find(_find_by_ID int64) (User, error) {

	var model_slice []User
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "ID", _find_by_ID)
	results, err := o._adapter.Query(q)
	if err != nil {
		return User{}, err
	}

	for _, result := range results {
		ro := User{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return User{}, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return User{}, errors.New("not found")
	}
	return model_slice[0], nil

}
func (o *User) FindByUserLogin(_find_by_UserLogin string) ([]User, error) {

	var model_slice []User
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "user_login", _find_by_UserLogin)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := User{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *User) FindByUserPass(_find_by_UserPass string) ([]User, error) {

	var model_slice []User
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "user_pass", _find_by_UserPass)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := User{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *User) FindByUserNicename(_find_by_UserNicename string) ([]User, error) {

	var model_slice []User
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "user_nicename", _find_by_UserNicename)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := User{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *User) FindByUserEmail(_find_by_UserEmail string) ([]User, error) {

	var model_slice []User
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "user_email", _find_by_UserEmail)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := User{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *User) FindByUserUrl(_find_by_UserUrl string) ([]User, error) {

	var model_slice []User
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "user_url", _find_by_UserUrl)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := User{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *User) FindByUserRegistered(_find_by_UserRegistered DateTime) ([]User, error) {

	var model_slice []User
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "user_registered", _find_by_UserRegistered)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := User{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *User) FindByUserActivationKey(_find_by_UserActivationKey string) ([]User, error) {

	var model_slice []User
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "user_activation_key", _find_by_UserActivationKey)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := User{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *User) FindByUserStatus(_find_by_UserStatus int) ([]User, error) {

	var model_slice []User
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "user_status", _find_by_UserStatus)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := User{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *User) FindByDisplayName(_find_by_DisplayName string) ([]User, error) {

	var model_slice []User
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "display_name", _find_by_DisplayName)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := User{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}

func (o *User) FromDBValueMap(m map[string]DBValue) error {
	_ID, err := m["ID"].AsInt64()
	if err != nil {
		return err
	}
	o.ID = _ID
	_UserLogin, err := m["user_login"].AsString()
	if err != nil {
		return err
	}
	o.UserLogin = _UserLogin
	_UserPass, err := m["user_pass"].AsString()
	if err != nil {
		return err
	}
	o.UserPass = _UserPass
	_UserNicename, err := m["user_nicename"].AsString()
	if err != nil {
		return err
	}
	o.UserNicename = _UserNicename
	_UserEmail, err := m["user_email"].AsString()
	if err != nil {
		return err
	}
	o.UserEmail = _UserEmail
	_UserUrl, err := m["user_url"].AsString()
	if err != nil {
		return err
	}
	o.UserUrl = _UserUrl
	_UserRegistered, err := m["user_registered"].AsDateTime()
	if err != nil {
		return err
	}
	o.UserRegistered = _UserRegistered
	_UserActivationKey, err := m["user_activation_key"].AsString()
	if err != nil {
		return err
	}
	o.UserActivationKey = _UserActivationKey
	_UserStatus, err := m["user_status"].AsInt()
	if err != nil {
		return err
	}
	o.UserStatus = _UserStatus
	_DisplayName, err := m["display_name"].AsString()
	if err != nil {
		return err
	}
	o.DisplayName = _DisplayName

	return nil
}

func (o *User) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `user_login` = '%s', `user_pass` = '%s', `user_nicename` = '%s', `user_email` = '%s', `user_url` = '%s', `user_registered` = '%s', `user_activation_key` = '%s', `user_status` = '%d', `display_name` = '%s' WHERE %s = '%d' LIMIT 1", o._table, o.UserLogin, o.UserPass, o.UserNicename, o.UserEmail, o.UserUrl, o.UserRegistered, o.UserActivationKey, o.UserStatus, o.DisplayName, o._pkey, o.ID)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *User) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`user_login`, `user_pass`, `user_nicename`, `user_email`, `user_url`, `user_registered`, `user_activation_key`, `user_status`, `display_name`) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%d', '%s')", o._table, o.UserLogin, o.UserPass, o.UserNicename, o.UserEmail, o.UserUrl, o.UserRegistered, o.UserActivationKey, o.UserStatus, o.DisplayName)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *User) UpdateUserLogin(_upd_UserLogin string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_login` = '%s' WHERE `ID` = '%d'", o._table, _upd_UserLogin, o.UserLogin)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.UserLogin = _upd_UserLogin
	return o._adapter.AffectedRows(), nil
}

func (o *User) UpdateUserPass(_upd_UserPass string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_pass` = '%s' WHERE `ID` = '%d'", o._table, _upd_UserPass, o.UserPass)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.UserPass = _upd_UserPass
	return o._adapter.AffectedRows(), nil
}

func (o *User) UpdateUserNicename(_upd_UserNicename string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_nicename` = '%s' WHERE `ID` = '%d'", o._table, _upd_UserNicename, o.UserNicename)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.UserNicename = _upd_UserNicename
	return o._adapter.AffectedRows(), nil
}

func (o *User) UpdateUserEmail(_upd_UserEmail string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_email` = '%s' WHERE `ID` = '%d'", o._table, _upd_UserEmail, o.UserEmail)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.UserEmail = _upd_UserEmail
	return o._adapter.AffectedRows(), nil
}

func (o *User) UpdateUserUrl(_upd_UserUrl string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_url` = '%s' WHERE `ID` = '%d'", o._table, _upd_UserUrl, o.UserUrl)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.UserUrl = _upd_UserUrl
	return o._adapter.AffectedRows(), nil
}

func (o *User) UpdateUserRegistered(_upd_UserRegistered DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_registered` = '%s' WHERE `ID` = '%d'", o._table, _upd_UserRegistered, o.UserRegistered)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.UserRegistered = _upd_UserRegistered
	return o._adapter.AffectedRows(), nil
}

func (o *User) UpdateUserActivationKey(_upd_UserActivationKey string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_activation_key` = '%s' WHERE `ID` = '%d'", o._table, _upd_UserActivationKey, o.UserActivationKey)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.UserActivationKey = _upd_UserActivationKey
	return o._adapter.AffectedRows(), nil
}

func (o *User) UpdateUserStatus(_upd_UserStatus int) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_status` = '%d' WHERE `ID` = '%d'", o._table, _upd_UserStatus, o.UserStatus)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.UserStatus = _upd_UserStatus
	return o._adapter.AffectedRows(), nil
}

func (o *User) UpdateDisplayName(_upd_DisplayName string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `display_name` = '%s' WHERE `ID` = '%d'", o._table, _upd_DisplayName, o.DisplayName)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.DisplayName = _upd_DisplayName
	return o._adapter.AffectedRows(), nil
}

type WooAttrTaxonomie struct {
	_table      string
	_adapter    Adapter
	_pkey       string // 0 The name of the primary key in this table
	_conds      []string
	_new        bool
	AttrId      int64
	AttrName    string
	AttrLabel   string
	AttrType    string
	AttrOrderby string
}

func NewWooAttrTaxonomie(a Adapter) *WooAttrTaxonomie {
	var o WooAttrTaxonomie
	o._table = fmt.Sprintf("%swoocommerce_attribute_taxonomies", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "attribute_id"
	o._new = false
	return &o
}

func (o *WooAttrTaxonomie) Find(_find_by_AttrId int64) (WooAttrTaxonomie, error) {

	var model_slice []WooAttrTaxonomie
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "attribute_id", _find_by_AttrId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return WooAttrTaxonomie{}, err
	}

	for _, result := range results {
		ro := WooAttrTaxonomie{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return WooAttrTaxonomie{}, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return WooAttrTaxonomie{}, errors.New("not found")
	}
	return model_slice[0], nil

}
func (o *WooAttrTaxonomie) FindByAttrName(_find_by_AttrName string) ([]WooAttrTaxonomie, error) {

	var model_slice []WooAttrTaxonomie
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "attribute_name", _find_by_AttrName)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooAttrTaxonomie{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooAttrTaxonomie) FindByAttrLabel(_find_by_AttrLabel string) ([]WooAttrTaxonomie, error) {

	var model_slice []WooAttrTaxonomie
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "attribute_label", _find_by_AttrLabel)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooAttrTaxonomie{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooAttrTaxonomie) FindByAttrType(_find_by_AttrType string) ([]WooAttrTaxonomie, error) {

	var model_slice []WooAttrTaxonomie
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "attribute_type", _find_by_AttrType)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooAttrTaxonomie{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooAttrTaxonomie) FindByAttrOrderby(_find_by_AttrOrderby string) ([]WooAttrTaxonomie, error) {

	var model_slice []WooAttrTaxonomie
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "attribute_orderby", _find_by_AttrOrderby)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooAttrTaxonomie{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}

func (o *WooAttrTaxonomie) FromDBValueMap(m map[string]DBValue) error {
	_AttrId, err := m["attribute_id"].AsInt64()
	if err != nil {
		return err
	}
	o.AttrId = _AttrId
	_AttrName, err := m["attribute_name"].AsString()
	if err != nil {
		return err
	}
	o.AttrName = _AttrName
	_AttrLabel, err := m["attribute_label"].AsString()
	if err != nil {
		return err
	}
	o.AttrLabel = _AttrLabel
	_AttrType, err := m["attribute_type"].AsString()
	if err != nil {
		return err
	}
	o.AttrType = _AttrType
	_AttrOrderby, err := m["attribute_orderby"].AsString()
	if err != nil {
		return err
	}
	o.AttrOrderby = _AttrOrderby

	return nil
}

func (o *WooAttrTaxonomie) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `attribute_name` = '%s', `attribute_label` = '%s', `attribute_type` = '%s', `attribute_orderby` = '%s' WHERE %s = '%d' LIMIT 1", o._table, o.AttrName, o.AttrLabel, o.AttrType, o.AttrOrderby, o._pkey, o.AttrId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *WooAttrTaxonomie) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`attribute_name`, `attribute_label`, `attribute_type`, `attribute_orderby`) VALUES ('%s', '%s', '%s', '%s')", o._table, o.AttrName, o.AttrLabel, o.AttrType, o.AttrOrderby)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *WooAttrTaxonomie) UpdateAttrName(_upd_AttrName string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `attribute_name` = '%s' WHERE `attribute_id` = '%d'", o._table, _upd_AttrName, o.AttrName)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.AttrName = _upd_AttrName
	return o._adapter.AffectedRows(), nil
}

func (o *WooAttrTaxonomie) UpdateAttrLabel(_upd_AttrLabel string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `attribute_label` = '%s' WHERE `attribute_id` = '%d'", o._table, _upd_AttrLabel, o.AttrLabel)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.AttrLabel = _upd_AttrLabel
	return o._adapter.AffectedRows(), nil
}

func (o *WooAttrTaxonomie) UpdateAttrType(_upd_AttrType string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `attribute_type` = '%s' WHERE `attribute_id` = '%d'", o._table, _upd_AttrType, o.AttrType)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.AttrType = _upd_AttrType
	return o._adapter.AffectedRows(), nil
}

func (o *WooAttrTaxonomie) UpdateAttrOrderby(_upd_AttrOrderby string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `attribute_orderby` = '%s' WHERE `attribute_id` = '%d'", o._table, _upd_AttrOrderby, o.AttrOrderby)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.AttrOrderby = _upd_AttrOrderby
	return o._adapter.AffectedRows(), nil
}

type WooDownloadableProductPerm struct {
	_table             string
	_adapter           Adapter
	_pkey              string // 0 The name of the primary key in this table
	_conds             []string
	_new               bool
	PermissionId       int64
	DownloadId         string
	ProductId          int64
	OrderId            int64
	OrderKey           string
	UserEmail          string
	UserId             int64
	DownloadsRemaining string
	AccessGranted      DateTime
	AccessExpires      DateTime
	DownloadCount      int64
}

func NewWooDownloadableProductPerm(a Adapter) *WooDownloadableProductPerm {
	var o WooDownloadableProductPerm
	o._table = fmt.Sprintf("%swoocommerce_downloadable_product_permissions", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "permission_id"
	o._new = false
	return &o
}

func (o *WooDownloadableProductPerm) Find(_find_by_PermissionId int64) (WooDownloadableProductPerm, error) {

	var model_slice []WooDownloadableProductPerm
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "permission_id", _find_by_PermissionId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return WooDownloadableProductPerm{}, err
	}

	for _, result := range results {
		ro := WooDownloadableProductPerm{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return WooDownloadableProductPerm{}, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return WooDownloadableProductPerm{}, errors.New("not found")
	}
	return model_slice[0], nil

}
func (o *WooDownloadableProductPerm) FindByDownloadId(_find_by_DownloadId string) ([]WooDownloadableProductPerm, error) {

	var model_slice []WooDownloadableProductPerm
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "download_id", _find_by_DownloadId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooDownloadableProductPerm{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooDownloadableProductPerm) FindByProductId(_find_by_ProductId int64) ([]WooDownloadableProductPerm, error) {

	var model_slice []WooDownloadableProductPerm
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "product_id", _find_by_ProductId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooDownloadableProductPerm{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooDownloadableProductPerm) FindByOrderId(_find_by_OrderId int64) ([]WooDownloadableProductPerm, error) {

	var model_slice []WooDownloadableProductPerm
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "order_id", _find_by_OrderId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooDownloadableProductPerm{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooDownloadableProductPerm) FindByOrderKey(_find_by_OrderKey string) ([]WooDownloadableProductPerm, error) {

	var model_slice []WooDownloadableProductPerm
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "order_key", _find_by_OrderKey)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooDownloadableProductPerm{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooDownloadableProductPerm) FindByUserEmail(_find_by_UserEmail string) ([]WooDownloadableProductPerm, error) {

	var model_slice []WooDownloadableProductPerm
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "user_email", _find_by_UserEmail)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooDownloadableProductPerm{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooDownloadableProductPerm) FindByUserId(_find_by_UserId int64) ([]WooDownloadableProductPerm, error) {

	var model_slice []WooDownloadableProductPerm
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "user_id", _find_by_UserId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooDownloadableProductPerm{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooDownloadableProductPerm) FindByDownloadsRemaining(_find_by_DownloadsRemaining string) ([]WooDownloadableProductPerm, error) {

	var model_slice []WooDownloadableProductPerm
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "downloads_remaining", _find_by_DownloadsRemaining)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooDownloadableProductPerm{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooDownloadableProductPerm) FindByAccessGranted(_find_by_AccessGranted DateTime) ([]WooDownloadableProductPerm, error) {

	var model_slice []WooDownloadableProductPerm
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "access_granted", _find_by_AccessGranted)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooDownloadableProductPerm{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooDownloadableProductPerm) FindByAccessExpires(_find_by_AccessExpires DateTime) ([]WooDownloadableProductPerm, error) {

	var model_slice []WooDownloadableProductPerm
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "access_expires", _find_by_AccessExpires)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooDownloadableProductPerm{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooDownloadableProductPerm) FindByDownloadCount(_find_by_DownloadCount int64) ([]WooDownloadableProductPerm, error) {

	var model_slice []WooDownloadableProductPerm
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "download_count", _find_by_DownloadCount)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooDownloadableProductPerm{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}

func (o *WooDownloadableProductPerm) FromDBValueMap(m map[string]DBValue) error {
	_PermissionId, err := m["permission_id"].AsInt64()
	if err != nil {
		return err
	}
	o.PermissionId = _PermissionId
	_DownloadId, err := m["download_id"].AsString()
	if err != nil {
		return err
	}
	o.DownloadId = _DownloadId
	_ProductId, err := m["product_id"].AsInt64()
	if err != nil {
		return err
	}
	o.ProductId = _ProductId
	_OrderId, err := m["order_id"].AsInt64()
	if err != nil {
		return err
	}
	o.OrderId = _OrderId
	_OrderKey, err := m["order_key"].AsString()
	if err != nil {
		return err
	}
	o.OrderKey = _OrderKey
	_UserEmail, err := m["user_email"].AsString()
	if err != nil {
		return err
	}
	o.UserEmail = _UserEmail
	_UserId, err := m["user_id"].AsInt64()
	if err != nil {
		return err
	}
	o.UserId = _UserId
	_DownloadsRemaining, err := m["downloads_remaining"].AsString()
	if err != nil {
		return err
	}
	o.DownloadsRemaining = _DownloadsRemaining
	_AccessGranted, err := m["access_granted"].AsDateTime()
	if err != nil {
		return err
	}
	o.AccessGranted = _AccessGranted
	_AccessExpires, err := m["access_expires"].AsDateTime()
	if err != nil {
		return err
	}
	o.AccessExpires = _AccessExpires
	_DownloadCount, err := m["download_count"].AsInt64()
	if err != nil {
		return err
	}
	o.DownloadCount = _DownloadCount

	return nil
}

func (o *WooDownloadableProductPerm) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `download_id` = '%s', `product_id` = '%d', `order_id` = '%d', `order_key` = '%s', `user_email` = '%s', `user_id` = '%d', `downloads_remaining` = '%s', `access_granted` = '%s', `access_expires` = '%s', `download_count` = '%d' WHERE %s = '%d' LIMIT 1", o._table, o.DownloadId, o.ProductId, o.OrderId, o.OrderKey, o.UserEmail, o.UserId, o.DownloadsRemaining, o.AccessGranted, o.AccessExpires, o.DownloadCount, o._pkey, o.PermissionId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *WooDownloadableProductPerm) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`download_id`, `product_id`, `order_id`, `order_key`, `user_email`, `user_id`, `downloads_remaining`, `access_granted`, `access_expires`, `download_count`) VALUES ('%s', '%d', '%d', '%s', '%s', '%d', '%s', '%s', '%s', '%d')", o._table, o.DownloadId, o.ProductId, o.OrderId, o.OrderKey, o.UserEmail, o.UserId, o.DownloadsRemaining, o.AccessGranted, o.AccessExpires, o.DownloadCount)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *WooDownloadableProductPerm) UpdateDownloadId(_upd_DownloadId string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `download_id` = '%s' WHERE `permission_id` = '%d'", o._table, _upd_DownloadId, o.DownloadId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.DownloadId = _upd_DownloadId
	return o._adapter.AffectedRows(), nil
}

func (o *WooDownloadableProductPerm) UpdateProductId(_upd_ProductId int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `product_id` = '%d' WHERE `permission_id` = '%d'", o._table, _upd_ProductId, o.ProductId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.ProductId = _upd_ProductId
	return o._adapter.AffectedRows(), nil
}

func (o *WooDownloadableProductPerm) UpdateOrderId(_upd_OrderId int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `order_id` = '%d' WHERE `permission_id` = '%d'", o._table, _upd_OrderId, o.OrderId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.OrderId = _upd_OrderId
	return o._adapter.AffectedRows(), nil
}

func (o *WooDownloadableProductPerm) UpdateOrderKey(_upd_OrderKey string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `order_key` = '%s' WHERE `permission_id` = '%d'", o._table, _upd_OrderKey, o.OrderKey)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.OrderKey = _upd_OrderKey
	return o._adapter.AffectedRows(), nil
}

func (o *WooDownloadableProductPerm) UpdateUserEmail(_upd_UserEmail string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_email` = '%s' WHERE `permission_id` = '%d'", o._table, _upd_UserEmail, o.UserEmail)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.UserEmail = _upd_UserEmail
	return o._adapter.AffectedRows(), nil
}

func (o *WooDownloadableProductPerm) UpdateUserId(_upd_UserId int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_id` = '%d' WHERE `permission_id` = '%d'", o._table, _upd_UserId, o.UserId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.UserId = _upd_UserId
	return o._adapter.AffectedRows(), nil
}

func (o *WooDownloadableProductPerm) UpdateDownloadsRemaining(_upd_DownloadsRemaining string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `downloads_remaining` = '%s' WHERE `permission_id` = '%d'", o._table, _upd_DownloadsRemaining, o.DownloadsRemaining)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.DownloadsRemaining = _upd_DownloadsRemaining
	return o._adapter.AffectedRows(), nil
}

func (o *WooDownloadableProductPerm) UpdateAccessGranted(_upd_AccessGranted DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `access_granted` = '%s' WHERE `permission_id` = '%d'", o._table, _upd_AccessGranted, o.AccessGranted)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.AccessGranted = _upd_AccessGranted
	return o._adapter.AffectedRows(), nil
}

func (o *WooDownloadableProductPerm) UpdateAccessExpires(_upd_AccessExpires DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `access_expires` = '%s' WHERE `permission_id` = '%d'", o._table, _upd_AccessExpires, o.AccessExpires)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.AccessExpires = _upd_AccessExpires
	return o._adapter.AffectedRows(), nil
}

func (o *WooDownloadableProductPerm) UpdateDownloadCount(_upd_DownloadCount int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `download_count` = '%d' WHERE `permission_id` = '%d'", o._table, _upd_DownloadCount, o.DownloadCount)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.DownloadCount = _upd_DownloadCount
	return o._adapter.AffectedRows(), nil
}

type WooOrderItemMeta struct {
	_table      string
	_adapter    Adapter
	_pkey       string // 0 The name of the primary key in this table
	_conds      []string
	_new        bool
	MetaId      int64
	OrderItemId int64
	MetaKey     string
	MetaValue   string
}

func NewWooOrderItemMeta(a Adapter) *WooOrderItemMeta {
	var o WooOrderItemMeta
	o._table = fmt.Sprintf("%swoocommerce_order_itemmeta", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "meta_id"
	o._new = false
	return &o
}

func (o *WooOrderItemMeta) Find(_find_by_MetaId int64) (WooOrderItemMeta, error) {

	var model_slice []WooOrderItemMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "meta_id", _find_by_MetaId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return WooOrderItemMeta{}, err
	}

	for _, result := range results {
		ro := WooOrderItemMeta{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return WooOrderItemMeta{}, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return WooOrderItemMeta{}, errors.New("not found")
	}
	return model_slice[0], nil

}
func (o *WooOrderItemMeta) FindByOrderItemId(_find_by_OrderItemId int64) ([]WooOrderItemMeta, error) {

	var model_slice []WooOrderItemMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "order_item_id", _find_by_OrderItemId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooOrderItemMeta{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooOrderItemMeta) FindByMetaKey(_find_by_MetaKey string) ([]WooOrderItemMeta, error) {

	var model_slice []WooOrderItemMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "meta_key", _find_by_MetaKey)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooOrderItemMeta{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooOrderItemMeta) FindByMetaValue(_find_by_MetaValue string) ([]WooOrderItemMeta, error) {

	var model_slice []WooOrderItemMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "meta_value", _find_by_MetaValue)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooOrderItemMeta{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}

func (o *WooOrderItemMeta) FromDBValueMap(m map[string]DBValue) error {
	_MetaId, err := m["meta_id"].AsInt64()
	if err != nil {
		return err
	}
	o.MetaId = _MetaId
	_OrderItemId, err := m["order_item_id"].AsInt64()
	if err != nil {
		return err
	}
	o.OrderItemId = _OrderItemId
	_MetaKey, err := m["meta_key"].AsString()
	if err != nil {
		return err
	}
	o.MetaKey = _MetaKey
	_MetaValue, err := m["meta_value"].AsString()
	if err != nil {
		return err
	}
	o.MetaValue = _MetaValue

	return nil
}

func (o *WooOrderItemMeta) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `order_item_id` = '%d', `meta_key` = '%s', `meta_value` = '%s' WHERE %s = '%d' LIMIT 1", o._table, o.OrderItemId, o.MetaKey, o.MetaValue, o._pkey, o.MetaId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *WooOrderItemMeta) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`order_item_id`, `meta_key`, `meta_value`) VALUES ('%d', '%s', '%s')", o._table, o.OrderItemId, o.MetaKey, o.MetaValue)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *WooOrderItemMeta) UpdateOrderItemId(_upd_OrderItemId int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `order_item_id` = '%d' WHERE `meta_id` = '%d'", o._table, _upd_OrderItemId, o.OrderItemId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.OrderItemId = _upd_OrderItemId
	return o._adapter.AffectedRows(), nil
}

func (o *WooOrderItemMeta) UpdateMetaKey(_upd_MetaKey string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `meta_key` = '%s' WHERE `meta_id` = '%d'", o._table, _upd_MetaKey, o.MetaKey)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.MetaKey = _upd_MetaKey
	return o._adapter.AffectedRows(), nil
}

func (o *WooOrderItemMeta) UpdateMetaValue(_upd_MetaValue string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `meta_value` = '%s' WHERE `meta_id` = '%d'", o._table, _upd_MetaValue, o.MetaValue)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.MetaValue = _upd_MetaValue
	return o._adapter.AffectedRows(), nil
}

type WooOrderItem struct {
	_table        string
	_adapter      Adapter
	_pkey         string // 0 The name of the primary key in this table
	_conds        []string
	_new          bool
	OrderItemId   int64
	OrderItemName string
	OrderItemType string
	OrderId       int64
}

func NewWooOrderItem(a Adapter) *WooOrderItem {
	var o WooOrderItem
	o._table = fmt.Sprintf("%swoocommerce_order_items", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "order_item_id"
	o._new = false
	return &o
}

func (o *WooOrderItem) Find(_find_by_OrderItemId int64) (WooOrderItem, error) {

	var model_slice []WooOrderItem
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "order_item_id", _find_by_OrderItemId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return WooOrderItem{}, err
	}

	for _, result := range results {
		ro := WooOrderItem{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return WooOrderItem{}, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return WooOrderItem{}, errors.New("not found")
	}
	return model_slice[0], nil

}
func (o *WooOrderItem) FindByOrderItemName(_find_by_OrderItemName string) ([]WooOrderItem, error) {

	var model_slice []WooOrderItem
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "order_item_name", _find_by_OrderItemName)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooOrderItem{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooOrderItem) FindByOrderItemType(_find_by_OrderItemType string) ([]WooOrderItem, error) {

	var model_slice []WooOrderItem
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "order_item_type", _find_by_OrderItemType)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooOrderItem{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooOrderItem) FindByOrderId(_find_by_OrderId int64) ([]WooOrderItem, error) {

	var model_slice []WooOrderItem
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "order_id", _find_by_OrderId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooOrderItem{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}

func (o *WooOrderItem) FromDBValueMap(m map[string]DBValue) error {
	_OrderItemId, err := m["order_item_id"].AsInt64()
	if err != nil {
		return err
	}
	o.OrderItemId = _OrderItemId
	_OrderItemName, err := m["order_item_name"].AsString()
	if err != nil {
		return err
	}
	o.OrderItemName = _OrderItemName
	_OrderItemType, err := m["order_item_type"].AsString()
	if err != nil {
		return err
	}
	o.OrderItemType = _OrderItemType
	_OrderId, err := m["order_id"].AsInt64()
	if err != nil {
		return err
	}
	o.OrderId = _OrderId

	return nil
}

func (o *WooOrderItem) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `order_item_name` = '%s', `order_item_type` = '%s', `order_id` = '%d' WHERE %s = '%d' LIMIT 1", o._table, o.OrderItemName, o.OrderItemType, o.OrderId, o._pkey, o.OrderItemId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *WooOrderItem) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`order_item_name`, `order_item_type`, `order_id`) VALUES ('%s', '%s', '%d')", o._table, o.OrderItemName, o.OrderItemType, o.OrderId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *WooOrderItem) UpdateOrderItemName(_upd_OrderItemName string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `order_item_name` = '%s' WHERE `order_item_id` = '%d'", o._table, _upd_OrderItemName, o.OrderItemName)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.OrderItemName = _upd_OrderItemName
	return o._adapter.AffectedRows(), nil
}

func (o *WooOrderItem) UpdateOrderItemType(_upd_OrderItemType string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `order_item_type` = '%s' WHERE `order_item_id` = '%d'", o._table, _upd_OrderItemType, o.OrderItemType)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.OrderItemType = _upd_OrderItemType
	return o._adapter.AffectedRows(), nil
}

func (o *WooOrderItem) UpdateOrderId(_upd_OrderId int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `order_id` = '%d' WHERE `order_item_id` = '%d'", o._table, _upd_OrderId, o.OrderId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.OrderId = _upd_OrderId
	return o._adapter.AffectedRows(), nil
}

type WooTaxRateLocation struct {
	_table       string
	_adapter     Adapter
	_pkey        string // 0 The name of the primary key in this table
	_conds       []string
	_new         bool
	LocationId   int64
	LocationCode string
	TaxRateId    int64
	LocationType string
}

func NewWooTaxRateLocation(a Adapter) *WooTaxRateLocation {
	var o WooTaxRateLocation
	o._table = fmt.Sprintf("%swoocommerce_tax_rate_locations", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "location_id"
	o._new = false
	return &o
}

func (o *WooTaxRateLocation) Find(_find_by_LocationId int64) (WooTaxRateLocation, error) {

	var model_slice []WooTaxRateLocation
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "location_id", _find_by_LocationId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return WooTaxRateLocation{}, err
	}

	for _, result := range results {
		ro := WooTaxRateLocation{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return WooTaxRateLocation{}, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return WooTaxRateLocation{}, errors.New("not found")
	}
	return model_slice[0], nil

}
func (o *WooTaxRateLocation) FindByLocationCode(_find_by_LocationCode string) ([]WooTaxRateLocation, error) {

	var model_slice []WooTaxRateLocation
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "location_code", _find_by_LocationCode)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooTaxRateLocation{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooTaxRateLocation) FindByTaxRateId(_find_by_TaxRateId int64) ([]WooTaxRateLocation, error) {

	var model_slice []WooTaxRateLocation
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "tax_rate_id", _find_by_TaxRateId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooTaxRateLocation{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooTaxRateLocation) FindByLocationType(_find_by_LocationType string) ([]WooTaxRateLocation, error) {

	var model_slice []WooTaxRateLocation
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "location_type", _find_by_LocationType)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooTaxRateLocation{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}

func (o *WooTaxRateLocation) FromDBValueMap(m map[string]DBValue) error {
	_LocationId, err := m["location_id"].AsInt64()
	if err != nil {
		return err
	}
	o.LocationId = _LocationId
	_LocationCode, err := m["location_code"].AsString()
	if err != nil {
		return err
	}
	o.LocationCode = _LocationCode
	_TaxRateId, err := m["tax_rate_id"].AsInt64()
	if err != nil {
		return err
	}
	o.TaxRateId = _TaxRateId
	_LocationType, err := m["location_type"].AsString()
	if err != nil {
		return err
	}
	o.LocationType = _LocationType

	return nil
}

func (o *WooTaxRateLocation) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `location_code` = '%s', `tax_rate_id` = '%d', `location_type` = '%s' WHERE %s = '%d' LIMIT 1", o._table, o.LocationCode, o.TaxRateId, o.LocationType, o._pkey, o.LocationId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *WooTaxRateLocation) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`location_code`, `tax_rate_id`, `location_type`) VALUES ('%s', '%d', '%s')", o._table, o.LocationCode, o.TaxRateId, o.LocationType)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRateLocation) UpdateLocationCode(_upd_LocationCode string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `location_code` = '%s' WHERE `location_id` = '%d'", o._table, _upd_LocationCode, o.LocationCode)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.LocationCode = _upd_LocationCode
	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRateLocation) UpdateTaxRateId(_upd_TaxRateId int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `tax_rate_id` = '%d' WHERE `location_id` = '%d'", o._table, _upd_TaxRateId, o.TaxRateId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.TaxRateId = _upd_TaxRateId
	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRateLocation) UpdateLocationType(_upd_LocationType string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `location_type` = '%s' WHERE `location_id` = '%d'", o._table, _upd_LocationType, o.LocationType)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.LocationType = _upd_LocationType
	return o._adapter.AffectedRows(), nil
}

type WooTaxRate struct {
	_table          string
	_adapter        Adapter
	_pkey           string // 0 The name of the primary key in this table
	_conds          []string
	_new            bool
	TaxRateId       int64
	TaxRateCountry  string
	TaxRateState    string
	TaxRate         string
	TaxRateName     string
	TaxRatePriority int64
	TaxRateCompound int
	TaxRateShipping int
	TaxRateOrder    int64
	TaxRateClass    string
}

func NewWooTaxRate(a Adapter) *WooTaxRate {
	var o WooTaxRate
	o._table = fmt.Sprintf("%swoocommerce_tax_rates", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "tax_rate_id"
	o._new = false
	return &o
}

func (o *WooTaxRate) Find(_find_by_TaxRateId int64) (WooTaxRate, error) {

	var model_slice []WooTaxRate
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "tax_rate_id", _find_by_TaxRateId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return WooTaxRate{}, err
	}

	for _, result := range results {
		ro := WooTaxRate{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return WooTaxRate{}, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return WooTaxRate{}, errors.New("not found")
	}
	return model_slice[0], nil

}
func (o *WooTaxRate) FindByTaxRateCountry(_find_by_TaxRateCountry string) ([]WooTaxRate, error) {

	var model_slice []WooTaxRate
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "tax_rate_country", _find_by_TaxRateCountry)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooTaxRate{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooTaxRate) FindByTaxRateState(_find_by_TaxRateState string) ([]WooTaxRate, error) {

	var model_slice []WooTaxRate
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "tax_rate_state", _find_by_TaxRateState)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooTaxRate{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooTaxRate) FindByTaxRate(_find_by_TaxRate string) ([]WooTaxRate, error) {

	var model_slice []WooTaxRate
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "tax_rate", _find_by_TaxRate)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooTaxRate{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooTaxRate) FindByTaxRateName(_find_by_TaxRateName string) ([]WooTaxRate, error) {

	var model_slice []WooTaxRate
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "tax_rate_name", _find_by_TaxRateName)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooTaxRate{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooTaxRate) FindByTaxRatePriority(_find_by_TaxRatePriority int64) ([]WooTaxRate, error) {

	var model_slice []WooTaxRate
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "tax_rate_priority", _find_by_TaxRatePriority)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooTaxRate{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooTaxRate) FindByTaxRateCompound(_find_by_TaxRateCompound int) ([]WooTaxRate, error) {

	var model_slice []WooTaxRate
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "tax_rate_compound", _find_by_TaxRateCompound)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooTaxRate{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooTaxRate) FindByTaxRateShipping(_find_by_TaxRateShipping int) ([]WooTaxRate, error) {

	var model_slice []WooTaxRate
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "tax_rate_shipping", _find_by_TaxRateShipping)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooTaxRate{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooTaxRate) FindByTaxRateOrder(_find_by_TaxRateOrder int64) ([]WooTaxRate, error) {

	var model_slice []WooTaxRate
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "tax_rate_order", _find_by_TaxRateOrder)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooTaxRate{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}
func (o *WooTaxRate) FindByTaxRateClass(_find_by_TaxRateClass string) ([]WooTaxRate, error) {

	var model_slice []WooTaxRate
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "tax_rate_class", _find_by_TaxRateClass)
	results, err := o._adapter.Query(q)
	if err != nil {
		return model_slice, err
	}

	for _, result := range results {
		ro := WooTaxRate{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return model_slice, err
		}
		model_slice = append(model_slice, ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, errors.New("no results")
	}
	return model_slice, nil

}

func (o *WooTaxRate) FromDBValueMap(m map[string]DBValue) error {
	_TaxRateId, err := m["tax_rate_id"].AsInt64()
	if err != nil {
		return err
	}
	o.TaxRateId = _TaxRateId
	_TaxRateCountry, err := m["tax_rate_country"].AsString()
	if err != nil {
		return err
	}
	o.TaxRateCountry = _TaxRateCountry
	_TaxRateState, err := m["tax_rate_state"].AsString()
	if err != nil {
		return err
	}
	o.TaxRateState = _TaxRateState
	_TaxRate, err := m["tax_rate"].AsString()
	if err != nil {
		return err
	}
	o.TaxRate = _TaxRate
	_TaxRateName, err := m["tax_rate_name"].AsString()
	if err != nil {
		return err
	}
	o.TaxRateName = _TaxRateName
	_TaxRatePriority, err := m["tax_rate_priority"].AsInt64()
	if err != nil {
		return err
	}
	o.TaxRatePriority = _TaxRatePriority
	_TaxRateCompound, err := m["tax_rate_compound"].AsInt()
	if err != nil {
		return err
	}
	o.TaxRateCompound = _TaxRateCompound
	_TaxRateShipping, err := m["tax_rate_shipping"].AsInt()
	if err != nil {
		return err
	}
	o.TaxRateShipping = _TaxRateShipping
	_TaxRateOrder, err := m["tax_rate_order"].AsInt64()
	if err != nil {
		return err
	}
	o.TaxRateOrder = _TaxRateOrder
	_TaxRateClass, err := m["tax_rate_class"].AsString()
	if err != nil {
		return err
	}
	o.TaxRateClass = _TaxRateClass

	return nil
}

func (o *WooTaxRate) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `tax_rate_country` = '%s', `tax_rate_state` = '%s', `tax_rate` = '%s', `tax_rate_name` = '%s', `tax_rate_priority` = '%d', `tax_rate_compound` = '%d', `tax_rate_shipping` = '%d', `tax_rate_order` = '%d', `tax_rate_class` = '%s' WHERE %s = '%d' LIMIT 1", o._table, o.TaxRateCountry, o.TaxRateState, o.TaxRate, o.TaxRateName, o.TaxRatePriority, o.TaxRateCompound, o.TaxRateShipping, o.TaxRateOrder, o.TaxRateClass, o._pkey, o.TaxRateId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *WooTaxRate) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`tax_rate_country`, `tax_rate_state`, `tax_rate`, `tax_rate_name`, `tax_rate_priority`, `tax_rate_compound`, `tax_rate_shipping`, `tax_rate_order`, `tax_rate_class`) VALUES ('%s', '%s', '%s', '%s', '%d', '%d', '%d', '%d', '%s')", o._table, o.TaxRateCountry, o.TaxRateState, o.TaxRate, o.TaxRateName, o.TaxRatePriority, o.TaxRateCompound, o.TaxRateShipping, o.TaxRateOrder, o.TaxRateClass)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRate) UpdateTaxRateCountry(_upd_TaxRateCountry string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `tax_rate_country` = '%s' WHERE `tax_rate_id` = '%d'", o._table, _upd_TaxRateCountry, o.TaxRateCountry)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.TaxRateCountry = _upd_TaxRateCountry
	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRate) UpdateTaxRateState(_upd_TaxRateState string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `tax_rate_state` = '%s' WHERE `tax_rate_id` = '%d'", o._table, _upd_TaxRateState, o.TaxRateState)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.TaxRateState = _upd_TaxRateState
	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRate) UpdateTaxRate(_upd_TaxRate string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `tax_rate` = '%s' WHERE `tax_rate_id` = '%d'", o._table, _upd_TaxRate, o.TaxRate)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.TaxRate = _upd_TaxRate
	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRate) UpdateTaxRateName(_upd_TaxRateName string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `tax_rate_name` = '%s' WHERE `tax_rate_id` = '%d'", o._table, _upd_TaxRateName, o.TaxRateName)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.TaxRateName = _upd_TaxRateName
	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRate) UpdateTaxRatePriority(_upd_TaxRatePriority int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `tax_rate_priority` = '%d' WHERE `tax_rate_id` = '%d'", o._table, _upd_TaxRatePriority, o.TaxRatePriority)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.TaxRatePriority = _upd_TaxRatePriority
	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRate) UpdateTaxRateCompound(_upd_TaxRateCompound int) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `tax_rate_compound` = '%d' WHERE `tax_rate_id` = '%d'", o._table, _upd_TaxRateCompound, o.TaxRateCompound)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.TaxRateCompound = _upd_TaxRateCompound
	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRate) UpdateTaxRateShipping(_upd_TaxRateShipping int) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `tax_rate_shipping` = '%d' WHERE `tax_rate_id` = '%d'", o._table, _upd_TaxRateShipping, o.TaxRateShipping)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.TaxRateShipping = _upd_TaxRateShipping
	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRate) UpdateTaxRateOrder(_upd_TaxRateOrder int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `tax_rate_order` = '%d' WHERE `tax_rate_id` = '%d'", o._table, _upd_TaxRateOrder, o.TaxRateOrder)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.TaxRateOrder = _upd_TaxRateOrder
	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRate) UpdateTaxRateClass(_upd_TaxRateClass string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `tax_rate_class` = '%s' WHERE `tax_rate_id` = '%d'", o._table, _upd_TaxRateClass, o.TaxRateClass)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.TaxRateClass = _upd_TaxRateClass
	return o._adapter.AffectedRows(), nil
}

type Adapter interface {
	Open(string, string, string, string) error
	Close()
	Query(string) ([]map[string]DBValue, error)
	Execute(string) error
	LastInsertedId() int64
	AffectedRows() int64
	DatabasePrefix() string
	NewDBValue() DBValue
}

type DBValue interface {
	AsInt() (int, error)
	AsInt32() (int32, error)
	AsInt64() (int64, error)
	AsFloat32() (float32, error)
	AsString() (string, error)
	AsDateTime() (DateTime, error)
	SetInternalValue(string, string)
}

type MysqlValue struct {
	_v string
	_k string
}

func (v *MysqlValue) SetInternalValue(key, value string) {
	v._v = value
	v._k = key

}
func (v *MysqlValue) AsString() (string, error) {
	return v._v, nil
}
func (v *MysqlValue) AsInt() (int, error) {
	i, err := strconv.ParseInt(v._v, 10, 32)
	return int(i), err
}
func (v *MysqlValue) AsInt32() (int32, error) {
	i, err := strconv.ParseInt(v._v, 10, 32)
	return int32(i), err
}
func (v *MysqlValue) AsInt64() (int64, error) {
	i, err := strconv.ParseInt(v._v, 10, 64)
	return i, err
}
func (v *MysqlValue) AsFloat32() (float32, error) {
	i, err := strconv.ParseFloat(v._v, 32)
	if err != nil {
		return 0.0, err
	}
	return float32(i), err
}
func (v *MysqlValue) AsFloat64() (float64, error) {
	i, err := strconv.ParseFloat(v._v, 64)
	if err != nil {
		return 0.0, err
	}
	return i, err
}

func (v *MysqlValue) AsDateTime() (DateTime, error) {
	var dt DateTime
	err := dt.FromString(v._v)
	if err != nil {
		return DateTime{}, err
	}
	return dt, nil
}

func NewMysqlValue() *MysqlValue {
	return &MysqlValue{}
}

type MysqlAdapter struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Pass     string `yaml: "pass"`
	Database string `yaml:"database"`
	DBPrefix string `yaml:"prefix"`
	_conn_   *sql.DB
	_lid     int64
	_cnt     int64
}

func NewMysqlAdapter(pre string) *MysqlAdapter {
	return &MysqlAdapter{DBPrefix: pre}
}
func (a *MysqlAdapter) NewDBValue() DBValue {
	return NewMysqlValue()
}
func (a *MysqlAdapter) DatabasePrefix() string {
	return a.DBPrefix
}
func (a *MysqlAdapter) FromYAML(b []byte) error {
	return yaml.Unmarshal(b, a)
}

func (a *MysqlAdapter) Open(h, u, p, d string) error {
	if h != "localhost" {
		tc, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", u, p, h, d))
		if err != nil {
			return err
		}
		a._conn_ = tc
	} else {
		tc, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", u, p, d))
		if err != nil {
			return err
		}
		a._conn_ = tc
	}
	return nil

}
func (a *MysqlAdapter) Close() {
	a._conn_.Close()
}

func (a *MysqlAdapter) Query(q string) ([]map[string]DBValue, error) {
	results := new([]map[string]DBValue)
	rows, err := a._conn_.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		res := make(map[string]DBValue)
		for i, col := range values {
			k := columns[i]
			res[k] = a.NewDBValue()
			res[k].SetInternalValue(k, string(col))
		}
		*results = append(*results, res)
	}
	return *results, nil
}
func (a *MysqlAdapter) Execute(q string) error {
	tx, err := a._conn_.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(q)
	if err != nil {
		return err
	}
	a._lid, err = res.LastInsertId()
	if err != nil {
		return err
	}
	a._cnt, err = res.RowsAffected()
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}
func (a *MysqlAdapter) LastInsertedId() int64 {
	return a._lid
}
func (a *MysqlAdapter) AffectedRows() int64 {
	return a._cnt
}

type DateTime struct {
	Day     int
	Month   int
	Year    int
	Hours   int
	Minutes int
	Seconds int
	Zone    string
	Offset  int
}

func (d *DateTime) FromString(s string) error {
	es := s
	re := regexp.MustCompile("(?P<year>[\\d]{4})-(?P<month>[\\d]{2})-(?P<day>[\\d]{2}) (?P<hours>[\\d]{2}):(?P<minutes>[\\d]{2}):(?P<seconds>[\\d]{2})\\.(?P<offset>[\\d]+)(?P<zone>[\\w]+)")
	n1 := re.SubexpNames()
	ir2 := re.FindAllStringSubmatch(es, -1)
	if len(ir2) == 0 {
		return errors.New(fmt.Sprintf("found now data to capture in %s", es))
	}
	r2 := ir2[0]
	for i, n := range r2 {
		if n1[i] == "year" {
			_Year, err := strconv.ParseInt(n, 10, 32)
			d.Year = int(_Year)
			if err != nil {
				return errors.New(fmt.Sprintf("failed to convert %s in %s received %s", n[i], es, err))
			}
		}
		if n1[i] == "month" {
			_Month, err := strconv.ParseInt(n, 10, 32)
			d.Month = int(_Month)
			if err != nil {
				return errors.New(fmt.Sprintf("failed to convert %s in %s received %s", n[i], es, err))
			}
		}
		if n1[i] == "day" {
			_Day, err := strconv.ParseInt(n, 10, 32)
			d.Day = int(_Day)
			if err != nil {
				return errors.New(fmt.Sprintf("failed to convert %s in %s received %s", n[i], es, err))
			}
		}
		if n1[i] == "hours" {
			_Hours, err := strconv.ParseInt(n, 10, 32)
			d.Hours = int(_Hours)
			if err != nil {
				return errors.New(fmt.Sprintf("failed to convert %s in %s received %s", n[i], es, err))
			}
		}
		if n1[i] == "minutes" {
			_Minutes, err := strconv.ParseInt(n, 10, 32)
			d.Minutes = int(_Minutes)
			if err != nil {
				return errors.New(fmt.Sprintf("failed to convert %s in %s received %s", n[i], es, err))
			}
		}
		if n1[i] == "seconds" {
			_Seconds, err := strconv.ParseInt(n, 10, 32)
			d.Seconds = int(_Seconds)
			if err != nil {
				return errors.New(fmt.Sprintf("failed to convert %s in %s received %s", n[i], es, err))
			}
		}
		if n1[i] == "offset" {
			_Offset, err := strconv.ParseInt(n, 10, 32)
			d.Offset = int(_Offset)
			if err != nil {
				return errors.New(fmt.Sprintf("failed to convert %s in %s received %s", n[i], es, err))
			}
		}
		if n1[i] == "zone" {
			d.Zone = n
		}
	}
	return nil
}
func (d *DateTime) ToString() string {
	return fmt.Sprintf("%d-%d-%d %d:%d:%d.%d.%s", d.Year, d.Month, d.Day, d.Hours, d.Minutes, d.Seconds, d.Offset, d.Zone)
}
func fileExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}
func filePutContents(p string, txt string) error {
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	_, err = w.WriteString(txt)
	w.Flush()
	return nil
}
func fileGetContents(p string) ([]byte, error) {
	return ioutil.ReadFile(p)
}
