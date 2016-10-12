package gopress

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"regexp"
	"strconv"
)

type CommentMeta struct {
	_table    string
	_adapter  Adapter
	_pkey     string // 0 The name of the primary key in this table
	_conds    []string
	_new      bool
	metaId    int64
	commentId int64
	metaKey   string
	metaValue string
}

func NewCommentMeta(a Adapter) *CommentMeta {
	var o CommentMeta
	o._table = fmt.Sprintf("%scommentmeta", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "meta_id"
	o._new = false
	return &o
}

func (o *CommentMeta) Find(_find_by_metaId int64) (CommentMeta, error) {

	var model_slice []CommentMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "meta_id", _find_by_metaId)
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
func (o *CommentMeta) FindByCommentId(_find_by_commentId int64) ([]CommentMeta, error) {

	var model_slice []CommentMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "comment_id", _find_by_commentId)
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
func (o *CommentMeta) FindByMetaKey(_find_by_metaKey string) ([]CommentMeta, error) {

	var model_slice []CommentMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "meta_key", _find_by_metaKey)
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
func (o *CommentMeta) FindByMetaValue(_find_by_metaValue string) ([]CommentMeta, error) {

	var model_slice []CommentMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "meta_value", _find_by_metaValue)
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
	_metaId, err := m["meta_id"].AsInt64()
	if err != nil {
		return err
	}
	o.metaId = _metaId
	_commentId, err := m["comment_id"].AsInt64()
	if err != nil {
		return err
	}
	o.commentId = _commentId
	_metaKey, err := m["meta_key"].AsString()
	if err != nil {
		return err
	}
	o.metaKey = _metaKey
	_metaValue, err := m["meta_value"].AsString()
	if err != nil {
		return err
	}
	o.metaValue = _metaValue

	return nil
}

func (o *CommentMeta) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `comment_id` = '%d', `meta_key` = '%s', `meta_value` = '%s' WHERE %s = '%d' LIMIT 1", o._table, o.commentId, o.metaKey, o.metaValue, o._pkey, o.metaId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *CommentMeta) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`comment_id`, `meta_key`, `meta_value`) VALUES ('%d', '%s', '%s')", o._table, o.commentId, o.metaKey, o.metaValue)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *CommentMeta) UpdateCommentId(_upd_commentId int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_id` = '%d' WHERE `meta_id` = '%d'", o._table, _upd_commentId, o.commentId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.commentId = _upd_commentId
	return o._adapter.AffectedRows(), nil
}

func (o *CommentMeta) UpdateMetaKey(_upd_metaKey string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `meta_key` = '%s' WHERE `meta_id` = '%d'", o._table, _upd_metaKey, o.metaKey)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.metaKey = _upd_metaKey
	return o._adapter.AffectedRows(), nil
}

func (o *CommentMeta) UpdateMetaValue(_upd_metaValue string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `meta_value` = '%s' WHERE `meta_id` = '%d'", o._table, _upd_metaValue, o.metaValue)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.metaValue = _upd_metaValue
	return o._adapter.AffectedRows(), nil
}

type Comment struct {
	_table             string
	_adapter           Adapter
	_pkey              string // 0 The name of the primary key in this table
	_conds             []string
	_new               bool
	commentID          int64
	commentPostID      int64
	commentAuthor      string
	commentAuthorEmail string
	commentAuthorUrl   string
	commentAuthorIP    string
	commentDate        DateTime
	commentDateGmt     DateTime
	commentContent     string
	commentKarma       int
	commentApproved    string
	commentAgent       string
	commentType        string
	commentParent      int64
	userId             int64
}

func NewComment(a Adapter) *Comment {
	var o Comment
	o._table = fmt.Sprintf("%scomments", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "comment_ID"
	o._new = false
	return &o
}

func (o *Comment) Find(_find_by_commentID int64) (Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "comment_ID", _find_by_commentID)
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
func (o *Comment) FindByCommentPostID(_find_by_commentPostID int64) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "comment_post_ID", _find_by_commentPostID)
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
func (o *Comment) FindByCommentAuthor(_find_by_commentAuthor string) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "comment_author", _find_by_commentAuthor)
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
func (o *Comment) FindByCommentAuthorEmail(_find_by_commentAuthorEmail string) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "comment_author_email", _find_by_commentAuthorEmail)
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
func (o *Comment) FindByCommentAuthorUrl(_find_by_commentAuthorUrl string) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "comment_author_url", _find_by_commentAuthorUrl)
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
func (o *Comment) FindByCommentAuthorIP(_find_by_commentAuthorIP string) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "comment_author_IP", _find_by_commentAuthorIP)
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
func (o *Comment) FindByCommentDate(_find_by_commentDate DateTime) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "comment_date", _find_by_commentDate)
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
func (o *Comment) FindByCommentDateGmt(_find_by_commentDateGmt DateTime) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "comment_date_gmt", _find_by_commentDateGmt)
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
func (o *Comment) FindByCommentContent(_find_by_commentContent string) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "comment_content", _find_by_commentContent)
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
func (o *Comment) FindByCommentKarma(_find_by_commentKarma int) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "comment_karma", _find_by_commentKarma)
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
func (o *Comment) FindByCommentApproved(_find_by_commentApproved string) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "comment_approved", _find_by_commentApproved)
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
func (o *Comment) FindByCommentAgent(_find_by_commentAgent string) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "comment_agent", _find_by_commentAgent)
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
func (o *Comment) FindByCommentType(_find_by_commentType string) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "comment_type", _find_by_commentType)
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
func (o *Comment) FindByCommentParent(_find_by_commentParent int64) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "comment_parent", _find_by_commentParent)
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
func (o *Comment) FindByUserId(_find_by_userId int64) ([]Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "user_id", _find_by_userId)
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
	_commentID, err := m["comment_ID"].AsInt64()
	if err != nil {
		return err
	}
	o.commentID = _commentID
	_commentPostID, err := m["comment_post_ID"].AsInt64()
	if err != nil {
		return err
	}
	o.commentPostID = _commentPostID
	_commentAuthor, err := m["comment_author"].AsString()
	if err != nil {
		return err
	}
	o.commentAuthor = _commentAuthor
	_commentAuthorEmail, err := m["comment_author_email"].AsString()
	if err != nil {
		return err
	}
	o.commentAuthorEmail = _commentAuthorEmail
	_commentAuthorUrl, err := m["comment_author_url"].AsString()
	if err != nil {
		return err
	}
	o.commentAuthorUrl = _commentAuthorUrl
	_commentAuthorIP, err := m["comment_author_IP"].AsString()
	if err != nil {
		return err
	}
	o.commentAuthorIP = _commentAuthorIP
	_commentDate, err := m["comment_date"].AsDateTime()
	if err != nil {
		return err
	}
	o.commentDate = _commentDate
	_commentDateGmt, err := m["comment_date_gmt"].AsDateTime()
	if err != nil {
		return err
	}
	o.commentDateGmt = _commentDateGmt
	_commentContent, err := m["comment_content"].AsString()
	if err != nil {
		return err
	}
	o.commentContent = _commentContent
	_commentKarma, err := m["comment_karma"].AsInt()
	if err != nil {
		return err
	}
	o.commentKarma = _commentKarma
	_commentApproved, err := m["comment_approved"].AsString()
	if err != nil {
		return err
	}
	o.commentApproved = _commentApproved
	_commentAgent, err := m["comment_agent"].AsString()
	if err != nil {
		return err
	}
	o.commentAgent = _commentAgent
	_commentType, err := m["comment_type"].AsString()
	if err != nil {
		return err
	}
	o.commentType = _commentType
	_commentParent, err := m["comment_parent"].AsInt64()
	if err != nil {
		return err
	}
	o.commentParent = _commentParent
	_userId, err := m["user_id"].AsInt64()
	if err != nil {
		return err
	}
	o.userId = _userId

	return nil
}

func (o *Comment) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `comment_post_ID` = '%d', `comment_author` = '%s', `comment_author_email` = '%s', `comment_author_url` = '%s', `comment_author_IP` = '%s', `comment_date` = '%s', `comment_date_gmt` = '%s', `comment_content` = '%s', `comment_karma` = '%d', `comment_approved` = '%s', `comment_agent` = '%s', `comment_type` = '%s', `comment_parent` = '%d', `user_id` = '%d' WHERE %s = '%d' LIMIT 1", o._table, o.commentPostID, o.commentAuthor, o.commentAuthorEmail, o.commentAuthorUrl, o.commentAuthorIP, o.commentDate, o.commentDateGmt, o.commentContent, o.commentKarma, o.commentApproved, o.commentAgent, o.commentType, o.commentParent, o.userId, o._pkey, o.commentID)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *Comment) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`comment_post_ID`, `comment_author`, `comment_author_email`, `comment_author_url`, `comment_author_IP`, `comment_date`, `comment_date_gmt`, `comment_content`, `comment_karma`, `comment_approved`, `comment_agent`, `comment_type`, `comment_parent`, `user_id`) VALUES ('%d', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%d', '%s', '%s', '%s', '%d', '%d')", o._table, o.commentPostID, o.commentAuthor, o.commentAuthorEmail, o.commentAuthorUrl, o.commentAuthorIP, o.commentDate, o.commentDateGmt, o.commentContent, o.commentKarma, o.commentApproved, o.commentAgent, o.commentType, o.commentParent, o.userId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentPostID(_upd_commentPostID int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_post_ID` = '%d' WHERE `comment_ID` = '%d'", o._table, _upd_commentPostID, o.commentPostID)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.commentPostID = _upd_commentPostID
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentAuthor(_upd_commentAuthor string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_author` = '%s' WHERE `comment_ID` = '%d'", o._table, _upd_commentAuthor, o.commentAuthor)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.commentAuthor = _upd_commentAuthor
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentAuthorEmail(_upd_commentAuthorEmail string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_author_email` = '%s' WHERE `comment_ID` = '%d'", o._table, _upd_commentAuthorEmail, o.commentAuthorEmail)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.commentAuthorEmail = _upd_commentAuthorEmail
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentAuthorUrl(_upd_commentAuthorUrl string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_author_url` = '%s' WHERE `comment_ID` = '%d'", o._table, _upd_commentAuthorUrl, o.commentAuthorUrl)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.commentAuthorUrl = _upd_commentAuthorUrl
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentAuthorIP(_upd_commentAuthorIP string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_author_IP` = '%s' WHERE `comment_ID` = '%d'", o._table, _upd_commentAuthorIP, o.commentAuthorIP)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.commentAuthorIP = _upd_commentAuthorIP
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentDate(_upd_commentDate DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_date` = '%s' WHERE `comment_ID` = '%d'", o._table, _upd_commentDate, o.commentDate)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.commentDate = _upd_commentDate
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentDateGmt(_upd_commentDateGmt DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_date_gmt` = '%s' WHERE `comment_ID` = '%d'", o._table, _upd_commentDateGmt, o.commentDateGmt)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.commentDateGmt = _upd_commentDateGmt
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentContent(_upd_commentContent string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_content` = '%s' WHERE `comment_ID` = '%d'", o._table, _upd_commentContent, o.commentContent)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.commentContent = _upd_commentContent
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentKarma(_upd_commentKarma int) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_karma` = '%d' WHERE `comment_ID` = '%d'", o._table, _upd_commentKarma, o.commentKarma)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.commentKarma = _upd_commentKarma
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentApproved(_upd_commentApproved string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_approved` = '%s' WHERE `comment_ID` = '%d'", o._table, _upd_commentApproved, o.commentApproved)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.commentApproved = _upd_commentApproved
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentAgent(_upd_commentAgent string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_agent` = '%s' WHERE `comment_ID` = '%d'", o._table, _upd_commentAgent, o.commentAgent)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.commentAgent = _upd_commentAgent
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentType(_upd_commentType string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_type` = '%s' WHERE `comment_ID` = '%d'", o._table, _upd_commentType, o.commentType)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.commentType = _upd_commentType
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentParent(_upd_commentParent int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_parent` = '%d' WHERE `comment_ID` = '%d'", o._table, _upd_commentParent, o.commentParent)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.commentParent = _upd_commentParent
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateUserId(_upd_userId int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_id` = '%d' WHERE `comment_ID` = '%d'", o._table, _upd_userId, o.userId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.userId = _upd_userId
	return o._adapter.AffectedRows(), nil
}

type Link struct {
	_table          string
	_adapter        Adapter
	_pkey           string // 0 The name of the primary key in this table
	_conds          []string
	_new            bool
	linkId          int64
	linkUrl         string
	linkName        string
	linkImage       string
	linkTarget      string
	linkDescription string
	linkVisible     string
	linkOwner       int64
	linkRating      int
	linkUpdated     DateTime
	linkRel         string
	linkNotes       string
	linkRss         string
}

func NewLink(a Adapter) *Link {
	var o Link
	o._table = fmt.Sprintf("%slinks", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "link_id"
	o._new = false
	return &o
}

func (o *Link) Find(_find_by_linkId int64) (Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "link_id", _find_by_linkId)
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
func (o *Link) FindByLinkUrl(_find_by_linkUrl string) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "link_url", _find_by_linkUrl)
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
func (o *Link) FindByLinkName(_find_by_linkName string) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "link_name", _find_by_linkName)
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
func (o *Link) FindByLinkImage(_find_by_linkImage string) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "link_image", _find_by_linkImage)
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
func (o *Link) FindByLinkTarget(_find_by_linkTarget string) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "link_target", _find_by_linkTarget)
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
func (o *Link) FindByLinkDescription(_find_by_linkDescription string) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "link_description", _find_by_linkDescription)
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
func (o *Link) FindByLinkVisible(_find_by_linkVisible string) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "link_visible", _find_by_linkVisible)
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
func (o *Link) FindByLinkOwner(_find_by_linkOwner int64) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "link_owner", _find_by_linkOwner)
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
func (o *Link) FindByLinkRating(_find_by_linkRating int) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "link_rating", _find_by_linkRating)
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
func (o *Link) FindByLinkUpdated(_find_by_linkUpdated DateTime) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "link_updated", _find_by_linkUpdated)
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
func (o *Link) FindByLinkRel(_find_by_linkRel string) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "link_rel", _find_by_linkRel)
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
func (o *Link) FindByLinkNotes(_find_by_linkNotes string) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "link_notes", _find_by_linkNotes)
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
func (o *Link) FindByLinkRss(_find_by_linkRss string) ([]Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "link_rss", _find_by_linkRss)
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
	_linkId, err := m["link_id"].AsInt64()
	if err != nil {
		return err
	}
	o.linkId = _linkId
	_linkUrl, err := m["link_url"].AsString()
	if err != nil {
		return err
	}
	o.linkUrl = _linkUrl
	_linkName, err := m["link_name"].AsString()
	if err != nil {
		return err
	}
	o.linkName = _linkName
	_linkImage, err := m["link_image"].AsString()
	if err != nil {
		return err
	}
	o.linkImage = _linkImage
	_linkTarget, err := m["link_target"].AsString()
	if err != nil {
		return err
	}
	o.linkTarget = _linkTarget
	_linkDescription, err := m["link_description"].AsString()
	if err != nil {
		return err
	}
	o.linkDescription = _linkDescription
	_linkVisible, err := m["link_visible"].AsString()
	if err != nil {
		return err
	}
	o.linkVisible = _linkVisible
	_linkOwner, err := m["link_owner"].AsInt64()
	if err != nil {
		return err
	}
	o.linkOwner = _linkOwner
	_linkRating, err := m["link_rating"].AsInt()
	if err != nil {
		return err
	}
	o.linkRating = _linkRating
	_linkUpdated, err := m["link_updated"].AsDateTime()
	if err != nil {
		return err
	}
	o.linkUpdated = _linkUpdated
	_linkRel, err := m["link_rel"].AsString()
	if err != nil {
		return err
	}
	o.linkRel = _linkRel
	_linkNotes, err := m["link_notes"].AsString()
	if err != nil {
		return err
	}
	o.linkNotes = _linkNotes
	_linkRss, err := m["link_rss"].AsString()
	if err != nil {
		return err
	}
	o.linkRss = _linkRss

	return nil
}

func (o *Link) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `link_url` = '%s', `link_name` = '%s', `link_image` = '%s', `link_target` = '%s', `link_description` = '%s', `link_visible` = '%s', `link_owner` = '%d', `link_rating` = '%d', `link_updated` = '%s', `link_rel` = '%s', `link_notes` = '%s', `link_rss` = '%s' WHERE %s = '%d' LIMIT 1", o._table, o.linkUrl, o.linkName, o.linkImage, o.linkTarget, o.linkDescription, o.linkVisible, o.linkOwner, o.linkRating, o.linkUpdated, o.linkRel, o.linkNotes, o.linkRss, o._pkey, o.linkId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *Link) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`link_url`, `link_name`, `link_image`, `link_target`, `link_description`, `link_visible`, `link_owner`, `link_rating`, `link_updated`, `link_rel`, `link_notes`, `link_rss`) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%d', '%d', '%s', '%s', '%s', '%s')", o._table, o.linkUrl, o.linkName, o.linkImage, o.linkTarget, o.linkDescription, o.linkVisible, o.linkOwner, o.linkRating, o.linkUpdated, o.linkRel, o.linkNotes, o.linkRss)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkUrl(_upd_linkUrl string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_url` = '%s' WHERE `link_id` = '%d'", o._table, _upd_linkUrl, o.linkUrl)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.linkUrl = _upd_linkUrl
	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkName(_upd_linkName string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_name` = '%s' WHERE `link_id` = '%d'", o._table, _upd_linkName, o.linkName)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.linkName = _upd_linkName
	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkImage(_upd_linkImage string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_image` = '%s' WHERE `link_id` = '%d'", o._table, _upd_linkImage, o.linkImage)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.linkImage = _upd_linkImage
	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkTarget(_upd_linkTarget string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_target` = '%s' WHERE `link_id` = '%d'", o._table, _upd_linkTarget, o.linkTarget)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.linkTarget = _upd_linkTarget
	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkDescription(_upd_linkDescription string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_description` = '%s' WHERE `link_id` = '%d'", o._table, _upd_linkDescription, o.linkDescription)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.linkDescription = _upd_linkDescription
	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkVisible(_upd_linkVisible string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_visible` = '%s' WHERE `link_id` = '%d'", o._table, _upd_linkVisible, o.linkVisible)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.linkVisible = _upd_linkVisible
	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkOwner(_upd_linkOwner int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_owner` = '%d' WHERE `link_id` = '%d'", o._table, _upd_linkOwner, o.linkOwner)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.linkOwner = _upd_linkOwner
	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkRating(_upd_linkRating int) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_rating` = '%d' WHERE `link_id` = '%d'", o._table, _upd_linkRating, o.linkRating)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.linkRating = _upd_linkRating
	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkUpdated(_upd_linkUpdated DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_updated` = '%s' WHERE `link_id` = '%d'", o._table, _upd_linkUpdated, o.linkUpdated)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.linkUpdated = _upd_linkUpdated
	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkRel(_upd_linkRel string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_rel` = '%s' WHERE `link_id` = '%d'", o._table, _upd_linkRel, o.linkRel)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.linkRel = _upd_linkRel
	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkNotes(_upd_linkNotes string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_notes` = '%s' WHERE `link_id` = '%d'", o._table, _upd_linkNotes, o.linkNotes)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.linkNotes = _upd_linkNotes
	return o._adapter.AffectedRows(), nil
}

func (o *Link) UpdateLinkRss(_upd_linkRss string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `link_rss` = '%s' WHERE `link_id` = '%d'", o._table, _upd_linkRss, o.linkRss)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.linkRss = _upd_linkRss
	return o._adapter.AffectedRows(), nil
}

type Option struct {
	_table      string
	_adapter    Adapter
	_pkey       string // 0 The name of the primary key in this table
	_conds      []string
	_new        bool
	optionId    int64
	optionName  string
	optionValue string
	autoload    string
}

func NewOption(a Adapter) *Option {
	var o Option
	o._table = fmt.Sprintf("%soptions", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "option_id"
	o._new = false
	return &o
}

func (o *Option) Find(_find_by_optionId int64) (Option, error) {

	var model_slice []Option
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "option_id", _find_by_optionId)
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
func (o *Option) FindByOptionName(_find_by_optionName string) ([]Option, error) {

	var model_slice []Option
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "option_name", _find_by_optionName)
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
func (o *Option) FindByOptionValue(_find_by_optionValue string) ([]Option, error) {

	var model_slice []Option
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "option_value", _find_by_optionValue)
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
func (o *Option) FindByAutoload(_find_by_autoload string) ([]Option, error) {

	var model_slice []Option
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "autoload", _find_by_autoload)
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
	_optionId, err := m["option_id"].AsInt64()
	if err != nil {
		return err
	}
	o.optionId = _optionId
	_optionName, err := m["option_name"].AsString()
	if err != nil {
		return err
	}
	o.optionName = _optionName
	_optionValue, err := m["option_value"].AsString()
	if err != nil {
		return err
	}
	o.optionValue = _optionValue
	_autoload, err := m["autoload"].AsString()
	if err != nil {
		return err
	}
	o.autoload = _autoload

	return nil
}

func (o *Option) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `option_name` = '%s', `option_value` = '%s', `autoload` = '%s' WHERE %s = '%d' LIMIT 1", o._table, o.optionName, o.optionValue, o.autoload, o._pkey, o.optionId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *Option) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`option_name`, `option_value`, `autoload`) VALUES ('%s', '%s', '%s')", o._table, o.optionName, o.optionValue, o.autoload)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *Option) UpdateOptionName(_upd_optionName string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `option_name` = '%s' WHERE `option_id` = '%d'", o._table, _upd_optionName, o.optionName)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.optionName = _upd_optionName
	return o._adapter.AffectedRows(), nil
}

func (o *Option) UpdateOptionValue(_upd_optionValue string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `option_value` = '%s' WHERE `option_id` = '%d'", o._table, _upd_optionValue, o.optionValue)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.optionValue = _upd_optionValue
	return o._adapter.AffectedRows(), nil
}

func (o *Option) UpdateAutoload(_upd_autoload string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `autoload` = '%s' WHERE `option_id` = '%d'", o._table, _upd_autoload, o.autoload)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.autoload = _upd_autoload
	return o._adapter.AffectedRows(), nil
}

type PostMeta struct {
	_table    string
	_adapter  Adapter
	_pkey     string // 0 The name of the primary key in this table
	_conds    []string
	_new      bool
	metaId    int64
	stId      int64
	metaKey   string
	metaValue string
}

func NewPostMeta(a Adapter) *PostMeta {
	var o PostMeta
	o._table = fmt.Sprintf("%spostmeta", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "meta_id"
	o._new = false
	return &o
}

func (o *PostMeta) Find(_find_by_metaId int64) (PostMeta, error) {

	var model_slice []PostMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "meta_id", _find_by_metaId)
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
func (o *PostMeta) FindByStId(_find_by_stId int64) ([]PostMeta, error) {

	var model_slice []PostMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "post_id", _find_by_stId)
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
func (o *PostMeta) FindByMetaKey(_find_by_metaKey string) ([]PostMeta, error) {

	var model_slice []PostMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "meta_key", _find_by_metaKey)
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
func (o *PostMeta) FindByMetaValue(_find_by_metaValue string) ([]PostMeta, error) {

	var model_slice []PostMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "meta_value", _find_by_metaValue)
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
	_metaId, err := m["meta_id"].AsInt64()
	if err != nil {
		return err
	}
	o.metaId = _metaId
	_stId, err := m["post_id"].AsInt64()
	if err != nil {
		return err
	}
	o.stId = _stId
	_metaKey, err := m["meta_key"].AsString()
	if err != nil {
		return err
	}
	o.metaKey = _metaKey
	_metaValue, err := m["meta_value"].AsString()
	if err != nil {
		return err
	}
	o.metaValue = _metaValue

	return nil
}

func (o *PostMeta) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `post_id` = '%d', `meta_key` = '%s', `meta_value` = '%s' WHERE %s = '%d' LIMIT 1", o._table, o.stId, o.metaKey, o.metaValue, o._pkey, o.metaId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *PostMeta) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`post_id`, `meta_key`, `meta_value`) VALUES ('%d', '%s', '%s')", o._table, o.stId, o.metaKey, o.metaValue)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *PostMeta) UpdateStId(_upd_stId int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_id` = '%d' WHERE `meta_id` = '%d'", o._table, _upd_stId, o.stId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.stId = _upd_stId
	return o._adapter.AffectedRows(), nil
}

func (o *PostMeta) UpdateMetaKey(_upd_metaKey string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `meta_key` = '%s' WHERE `meta_id` = '%d'", o._table, _upd_metaKey, o.metaKey)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.metaKey = _upd_metaKey
	return o._adapter.AffectedRows(), nil
}

func (o *PostMeta) UpdateMetaValue(_upd_metaValue string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `meta_value` = '%s' WHERE `meta_id` = '%d'", o._table, _upd_metaValue, o.metaValue)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.metaValue = _upd_metaValue
	return o._adapter.AffectedRows(), nil
}

type Post struct {
	_table            string
	_adapter          Adapter
	_pkey             string // 0 The name of the primary key in this table
	_conds            []string
	_new              bool
	iD                int64
	stAuthor          int64
	stDate            DateTime
	stDateGmt         DateTime
	stContent         string
	stTitle           string
	stExcerpt         string
	stStatus          string
	commentStatus     string
	pingStatus        string
	stPassword        string
	stName            string
	toPing            string
	pinged            string
	stModified        DateTime
	stModifiedGmt     DateTime
	stContentFiltered string
	stParent          int64
	guid              string
	menuOrder         int
	stType            string
	stMimeType        string
	commentCount      int64
}

func NewPost(a Adapter) *Post {
	var o Post
	o._table = fmt.Sprintf("%sposts", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "ID"
	o._new = false
	return &o
}

func (o *Post) Find(_find_by_iD int64) (Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "ID", _find_by_iD)
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
func (o *Post) FindByStAuthor(_find_by_stAuthor int64) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "post_author", _find_by_stAuthor)
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
func (o *Post) FindByStDate(_find_by_stDate DateTime) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_date", _find_by_stDate)
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
func (o *Post) FindByStDateGmt(_find_by_stDateGmt DateTime) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_date_gmt", _find_by_stDateGmt)
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
func (o *Post) FindByStContent(_find_by_stContent string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_content", _find_by_stContent)
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
func (o *Post) FindByStTitle(_find_by_stTitle string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_title", _find_by_stTitle)
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
func (o *Post) FindByStExcerpt(_find_by_stExcerpt string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_excerpt", _find_by_stExcerpt)
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
func (o *Post) FindByStStatus(_find_by_stStatus string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_status", _find_by_stStatus)
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
func (o *Post) FindByCommentStatus(_find_by_commentStatus string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "comment_status", _find_by_commentStatus)
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
func (o *Post) FindByPingStatus(_find_by_pingStatus string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "ping_status", _find_by_pingStatus)
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
func (o *Post) FindByStPassword(_find_by_stPassword string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_password", _find_by_stPassword)
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
func (o *Post) FindByStName(_find_by_stName string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_name", _find_by_stName)
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
func (o *Post) FindByToPing(_find_by_toPing string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "to_ping", _find_by_toPing)
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
func (o *Post) FindByPinged(_find_by_pinged string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "pinged", _find_by_pinged)
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
func (o *Post) FindByStModified(_find_by_stModified DateTime) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_modified", _find_by_stModified)
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
func (o *Post) FindByStModifiedGmt(_find_by_stModifiedGmt DateTime) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_modified_gmt", _find_by_stModifiedGmt)
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
func (o *Post) FindByStContentFiltered(_find_by_stContentFiltered string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_content_filtered", _find_by_stContentFiltered)
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
func (o *Post) FindByStParent(_find_by_stParent int64) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "post_parent", _find_by_stParent)
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
func (o *Post) FindByGuid(_find_by_guid string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "guid", _find_by_guid)
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
func (o *Post) FindByMenuOrder(_find_by_menuOrder int) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "menu_order", _find_by_menuOrder)
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
func (o *Post) FindByStType(_find_by_stType string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_type", _find_by_stType)
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
func (o *Post) FindByStMimeType(_find_by_stMimeType string) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_mime_type", _find_by_stMimeType)
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
func (o *Post) FindByCommentCount(_find_by_commentCount int64) ([]Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "comment_count", _find_by_commentCount)
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
	_iD, err := m["ID"].AsInt64()
	if err != nil {
		return err
	}
	o.iD = _iD
	_stAuthor, err := m["post_author"].AsInt64()
	if err != nil {
		return err
	}
	o.stAuthor = _stAuthor
	_stDate, err := m["post_date"].AsDateTime()
	if err != nil {
		return err
	}
	o.stDate = _stDate
	_stDateGmt, err := m["post_date_gmt"].AsDateTime()
	if err != nil {
		return err
	}
	o.stDateGmt = _stDateGmt
	_stContent, err := m["post_content"].AsString()
	if err != nil {
		return err
	}
	o.stContent = _stContent
	_stTitle, err := m["post_title"].AsString()
	if err != nil {
		return err
	}
	o.stTitle = _stTitle
	_stExcerpt, err := m["post_excerpt"].AsString()
	if err != nil {
		return err
	}
	o.stExcerpt = _stExcerpt
	_stStatus, err := m["post_status"].AsString()
	if err != nil {
		return err
	}
	o.stStatus = _stStatus
	_commentStatus, err := m["comment_status"].AsString()
	if err != nil {
		return err
	}
	o.commentStatus = _commentStatus
	_pingStatus, err := m["ping_status"].AsString()
	if err != nil {
		return err
	}
	o.pingStatus = _pingStatus
	_stPassword, err := m["post_password"].AsString()
	if err != nil {
		return err
	}
	o.stPassword = _stPassword
	_stName, err := m["post_name"].AsString()
	if err != nil {
		return err
	}
	o.stName = _stName
	_toPing, err := m["to_ping"].AsString()
	if err != nil {
		return err
	}
	o.toPing = _toPing
	_pinged, err := m["pinged"].AsString()
	if err != nil {
		return err
	}
	o.pinged = _pinged
	_stModified, err := m["post_modified"].AsDateTime()
	if err != nil {
		return err
	}
	o.stModified = _stModified
	_stModifiedGmt, err := m["post_modified_gmt"].AsDateTime()
	if err != nil {
		return err
	}
	o.stModifiedGmt = _stModifiedGmt
	_stContentFiltered, err := m["post_content_filtered"].AsString()
	if err != nil {
		return err
	}
	o.stContentFiltered = _stContentFiltered
	_stParent, err := m["post_parent"].AsInt64()
	if err != nil {
		return err
	}
	o.stParent = _stParent
	_guid, err := m["guid"].AsString()
	if err != nil {
		return err
	}
	o.guid = _guid
	_menuOrder, err := m["menu_order"].AsInt()
	if err != nil {
		return err
	}
	o.menuOrder = _menuOrder
	_stType, err := m["post_type"].AsString()
	if err != nil {
		return err
	}
	o.stType = _stType
	_stMimeType, err := m["post_mime_type"].AsString()
	if err != nil {
		return err
	}
	o.stMimeType = _stMimeType
	_commentCount, err := m["comment_count"].AsInt64()
	if err != nil {
		return err
	}
	o.commentCount = _commentCount

	return nil
}

func (o *Post) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `post_author` = '%d', `post_date` = '%s', `post_date_gmt` = '%s', `post_content` = '%s', `post_title` = '%s', `post_excerpt` = '%s', `post_status` = '%s', `comment_status` = '%s', `ping_status` = '%s', `post_password` = '%s', `post_name` = '%s', `to_ping` = '%s', `pinged` = '%s', `post_modified` = '%s', `post_modified_gmt` = '%s', `post_content_filtered` = '%s', `post_parent` = '%d', `guid` = '%s', `menu_order` = '%d', `post_type` = '%s', `post_mime_type` = '%s', `comment_count` = '%d' WHERE %s = '%d' LIMIT 1", o._table, o.stAuthor, o.stDate, o.stDateGmt, o.stContent, o.stTitle, o.stExcerpt, o.stStatus, o.commentStatus, o.pingStatus, o.stPassword, o.stName, o.toPing, o.pinged, o.stModified, o.stModifiedGmt, o.stContentFiltered, o.stParent, o.guid, o.menuOrder, o.stType, o.stMimeType, o.commentCount, o._pkey, o.iD)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *Post) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`post_author`, `post_date`, `post_date_gmt`, `post_content`, `post_title`, `post_excerpt`, `post_status`, `comment_status`, `ping_status`, `post_password`, `post_name`, `to_ping`, `pinged`, `post_modified`, `post_modified_gmt`, `post_content_filtered`, `post_parent`, `guid`, `menu_order`, `post_type`, `post_mime_type`, `comment_count`) VALUES ('%d', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%d', '%s', '%d', '%s', '%s', '%d')", o._table, o.stAuthor, o.stDate, o.stDateGmt, o.stContent, o.stTitle, o.stExcerpt, o.stStatus, o.commentStatus, o.pingStatus, o.stPassword, o.stName, o.toPing, o.pinged, o.stModified, o.stModifiedGmt, o.stContentFiltered, o.stParent, o.guid, o.menuOrder, o.stType, o.stMimeType, o.commentCount)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateStAuthor(_upd_stAuthor int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_author` = '%d' WHERE `ID` = '%d'", o._table, _upd_stAuthor, o.stAuthor)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.stAuthor = _upd_stAuthor
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateStDate(_upd_stDate DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_date` = '%s' WHERE `ID` = '%d'", o._table, _upd_stDate, o.stDate)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.stDate = _upd_stDate
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateStDateGmt(_upd_stDateGmt DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_date_gmt` = '%s' WHERE `ID` = '%d'", o._table, _upd_stDateGmt, o.stDateGmt)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.stDateGmt = _upd_stDateGmt
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateStContent(_upd_stContent string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_content` = '%s' WHERE `ID` = '%d'", o._table, _upd_stContent, o.stContent)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.stContent = _upd_stContent
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateStTitle(_upd_stTitle string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_title` = '%s' WHERE `ID` = '%d'", o._table, _upd_stTitle, o.stTitle)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.stTitle = _upd_stTitle
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateStExcerpt(_upd_stExcerpt string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_excerpt` = '%s' WHERE `ID` = '%d'", o._table, _upd_stExcerpt, o.stExcerpt)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.stExcerpt = _upd_stExcerpt
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateStStatus(_upd_stStatus string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_status` = '%s' WHERE `ID` = '%d'", o._table, _upd_stStatus, o.stStatus)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.stStatus = _upd_stStatus
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateCommentStatus(_upd_commentStatus string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_status` = '%s' WHERE `ID` = '%d'", o._table, _upd_commentStatus, o.commentStatus)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.commentStatus = _upd_commentStatus
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdatePingStatus(_upd_pingStatus string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `ping_status` = '%s' WHERE `ID` = '%d'", o._table, _upd_pingStatus, o.pingStatus)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.pingStatus = _upd_pingStatus
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateStPassword(_upd_stPassword string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_password` = '%s' WHERE `ID` = '%d'", o._table, _upd_stPassword, o.stPassword)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.stPassword = _upd_stPassword
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateStName(_upd_stName string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_name` = '%s' WHERE `ID` = '%d'", o._table, _upd_stName, o.stName)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.stName = _upd_stName
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateToPing(_upd_toPing string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `to_ping` = '%s' WHERE `ID` = '%d'", o._table, _upd_toPing, o.toPing)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.toPing = _upd_toPing
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdatePinged(_upd_pinged string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `pinged` = '%s' WHERE `ID` = '%d'", o._table, _upd_pinged, o.pinged)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.pinged = _upd_pinged
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateStModified(_upd_stModified DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_modified` = '%s' WHERE `ID` = '%d'", o._table, _upd_stModified, o.stModified)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.stModified = _upd_stModified
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateStModifiedGmt(_upd_stModifiedGmt DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_modified_gmt` = '%s' WHERE `ID` = '%d'", o._table, _upd_stModifiedGmt, o.stModifiedGmt)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.stModifiedGmt = _upd_stModifiedGmt
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateStContentFiltered(_upd_stContentFiltered string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_content_filtered` = '%s' WHERE `ID` = '%d'", o._table, _upd_stContentFiltered, o.stContentFiltered)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.stContentFiltered = _upd_stContentFiltered
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateStParent(_upd_stParent int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_parent` = '%d' WHERE `ID` = '%d'", o._table, _upd_stParent, o.stParent)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.stParent = _upd_stParent
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateGuid(_upd_guid string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `guid` = '%s' WHERE `ID` = '%d'", o._table, _upd_guid, o.guid)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.guid = _upd_guid
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateMenuOrder(_upd_menuOrder int) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `menu_order` = '%d' WHERE `ID` = '%d'", o._table, _upd_menuOrder, o.menuOrder)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.menuOrder = _upd_menuOrder
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateStType(_upd_stType string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_type` = '%s' WHERE `ID` = '%d'", o._table, _upd_stType, o.stType)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.stType = _upd_stType
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateStMimeType(_upd_stMimeType string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_mime_type` = '%s' WHERE `ID` = '%d'", o._table, _upd_stMimeType, o.stMimeType)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.stMimeType = _upd_stMimeType
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdateCommentCount(_upd_commentCount int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_count` = '%d' WHERE `ID` = '%d'", o._table, _upd_commentCount, o.commentCount)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.commentCount = _upd_commentCount
	return o._adapter.AffectedRows(), nil
}

type TermRelationship struct {
	_table         string
	_adapter       Adapter
	_pkey          string // 0 The name of the primary key in this table
	_conds         []string
	_new           bool
	objectId       int64
	termTaxonomyId int64
	termOrder      int
}

func NewTermRelationship(a Adapter) *TermRelationship {
	var o TermRelationship
	o._table = fmt.Sprintf("%sterm_relationships", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "term_taxonomy_id"
	o._new = false
	return &o
}

func (o *TermRelationship) FindByObjectId(_find_by_objectId int64) ([]TermRelationship, error) {

	var model_slice []TermRelationship
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "object_id", _find_by_objectId)
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
func (o *TermRelationship) Find(_find_by_termTaxonomyId int64) (TermRelationship, error) {

	var model_slice []TermRelationship
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "term_taxonomy_id", _find_by_termTaxonomyId)
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
func (o *TermRelationship) FindByTermOrder(_find_by_termOrder int) ([]TermRelationship, error) {

	var model_slice []TermRelationship
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "term_order", _find_by_termOrder)
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
	_objectId, err := m["object_id"].AsInt64()
	if err != nil {
		return err
	}
	o.objectId = _objectId
	_termTaxonomyId, err := m["term_taxonomy_id"].AsInt64()
	if err != nil {
		return err
	}
	o.termTaxonomyId = _termTaxonomyId
	_termOrder, err := m["term_order"].AsInt()
	if err != nil {
		return err
	}
	o.termOrder = _termOrder

	return nil
}

func (o *TermRelationship) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `term_order` = '%d' WHERE %s = '%d' LIMIT 1", o._table, o.termOrder, o._pkey, o.termTaxonomyId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *TermRelationship) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`term_order`) VALUES ('%d')", o._table, o.termOrder)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *TermRelationship) UpdateTermOrder(_upd_termOrder int) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `term_order` = '%d' WHERE `term_taxonomy_id` = '%d'", o._table, _upd_termOrder, o.termOrder)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.termOrder = _upd_termOrder
	return o._adapter.AffectedRows(), nil
}

type TermTaxonomy struct {
	_table         string
	_adapter       Adapter
	_pkey          string // 0 The name of the primary key in this table
	_conds         []string
	_new           bool
	termTaxonomyId int64
	termId         int64
	taxonomy       string
	description    string
	parent         int64
	count          int64
}

func NewTermTaxonomy(a Adapter) *TermTaxonomy {
	var o TermTaxonomy
	o._table = fmt.Sprintf("%sterm_taxonomy", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "term_taxonomy_id"
	o._new = false
	return &o
}

func (o *TermTaxonomy) Find(_find_by_termTaxonomyId int64) (TermTaxonomy, error) {

	var model_slice []TermTaxonomy
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "term_taxonomy_id", _find_by_termTaxonomyId)
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
func (o *TermTaxonomy) FindByTermId(_find_by_termId int64) ([]TermTaxonomy, error) {

	var model_slice []TermTaxonomy
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "term_id", _find_by_termId)
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
func (o *TermTaxonomy) FindByTaxonomy(_find_by_taxonomy string) ([]TermTaxonomy, error) {

	var model_slice []TermTaxonomy
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "taxonomy", _find_by_taxonomy)
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
func (o *TermTaxonomy) FindByDescription(_find_by_description string) ([]TermTaxonomy, error) {

	var model_slice []TermTaxonomy
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "description", _find_by_description)
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
func (o *TermTaxonomy) FindByParent(_find_by_parent int64) ([]TermTaxonomy, error) {

	var model_slice []TermTaxonomy
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "parent", _find_by_parent)
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
func (o *TermTaxonomy) FindByCount(_find_by_count int64) ([]TermTaxonomy, error) {

	var model_slice []TermTaxonomy
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "count", _find_by_count)
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
	_termTaxonomyId, err := m["term_taxonomy_id"].AsInt64()
	if err != nil {
		return err
	}
	o.termTaxonomyId = _termTaxonomyId
	_termId, err := m["term_id"].AsInt64()
	if err != nil {
		return err
	}
	o.termId = _termId
	_taxonomy, err := m["taxonomy"].AsString()
	if err != nil {
		return err
	}
	o.taxonomy = _taxonomy
	_description, err := m["description"].AsString()
	if err != nil {
		return err
	}
	o.description = _description
	_parent, err := m["parent"].AsInt64()
	if err != nil {
		return err
	}
	o.parent = _parent
	_count, err := m["count"].AsInt64()
	if err != nil {
		return err
	}
	o.count = _count

	return nil
}

func (o *TermTaxonomy) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `term_id` = '%d', `taxonomy` = '%s', `description` = '%s', `parent` = '%d', `count` = '%d' WHERE %s = '%d' LIMIT 1", o._table, o.termId, o.taxonomy, o.description, o.parent, o.count, o._pkey, o.termTaxonomyId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *TermTaxonomy) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`term_id`, `taxonomy`, `description`, `parent`, `count`) VALUES ('%d', '%s', '%s', '%d', '%d')", o._table, o.termId, o.taxonomy, o.description, o.parent, o.count)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *TermTaxonomy) UpdateTermId(_upd_termId int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `term_id` = '%d' WHERE `term_taxonomy_id` = '%d'", o._table, _upd_termId, o.termId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.termId = _upd_termId
	return o._adapter.AffectedRows(), nil
}

func (o *TermTaxonomy) UpdateTaxonomy(_upd_taxonomy string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `taxonomy` = '%s' WHERE `term_taxonomy_id` = '%d'", o._table, _upd_taxonomy, o.taxonomy)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.taxonomy = _upd_taxonomy
	return o._adapter.AffectedRows(), nil
}

func (o *TermTaxonomy) UpdateDescription(_upd_description string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `description` = '%s' WHERE `term_taxonomy_id` = '%d'", o._table, _upd_description, o.description)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.description = _upd_description
	return o._adapter.AffectedRows(), nil
}

func (o *TermTaxonomy) UpdateParent(_upd_parent int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `parent` = '%d' WHERE `term_taxonomy_id` = '%d'", o._table, _upd_parent, o.parent)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.parent = _upd_parent
	return o._adapter.AffectedRows(), nil
}

func (o *TermTaxonomy) UpdateCount(_upd_count int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `count` = '%d' WHERE `term_taxonomy_id` = '%d'", o._table, _upd_count, o.count)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.count = _upd_count
	return o._adapter.AffectedRows(), nil
}

type Term struct {
	_table    string
	_adapter  Adapter
	_pkey     string // 0 The name of the primary key in this table
	_conds    []string
	_new      bool
	termId    int64
	name      string
	slug      string
	termGroup int64
}

func NewTerm(a Adapter) *Term {
	var o Term
	o._table = fmt.Sprintf("%sterms", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "term_id"
	o._new = false
	return &o
}

func (o *Term) Find(_find_by_termId int64) (Term, error) {

	var model_slice []Term
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "term_id", _find_by_termId)
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
func (o *Term) FindByName(_find_by_name string) ([]Term, error) {

	var model_slice []Term
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "name", _find_by_name)
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
func (o *Term) FindBySlug(_find_by_slug string) ([]Term, error) {

	var model_slice []Term
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "slug", _find_by_slug)
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
func (o *Term) FindByTermGroup(_find_by_termGroup int64) ([]Term, error) {

	var model_slice []Term
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "term_group", _find_by_termGroup)
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
	_termId, err := m["term_id"].AsInt64()
	if err != nil {
		return err
	}
	o.termId = _termId
	_name, err := m["name"].AsString()
	if err != nil {
		return err
	}
	o.name = _name
	_slug, err := m["slug"].AsString()
	if err != nil {
		return err
	}
	o.slug = _slug
	_termGroup, err := m["term_group"].AsInt64()
	if err != nil {
		return err
	}
	o.termGroup = _termGroup

	return nil
}

func (o *Term) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `name` = '%s', `slug` = '%s', `term_group` = '%d' WHERE %s = '%d' LIMIT 1", o._table, o.name, o.slug, o.termGroup, o._pkey, o.termId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *Term) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`name`, `slug`, `term_group`) VALUES ('%s', '%s', '%d')", o._table, o.name, o.slug, o.termGroup)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *Term) UpdateName(_upd_name string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `name` = '%s' WHERE `term_id` = '%d'", o._table, _upd_name, o.name)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.name = _upd_name
	return o._adapter.AffectedRows(), nil
}

func (o *Term) UpdateSlug(_upd_slug string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `slug` = '%s' WHERE `term_id` = '%d'", o._table, _upd_slug, o.slug)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.slug = _upd_slug
	return o._adapter.AffectedRows(), nil
}

func (o *Term) UpdateTermGroup(_upd_termGroup int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `term_group` = '%d' WHERE `term_id` = '%d'", o._table, _upd_termGroup, o.termGroup)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.termGroup = _upd_termGroup
	return o._adapter.AffectedRows(), nil
}

type UserMeta struct {
	_table    string
	_adapter  Adapter
	_pkey     string // 0 The name of the primary key in this table
	_conds    []string
	_new      bool
	uMetaId   int64
	userId    int64
	metaKey   string
	metaValue string
}

func NewUserMeta(a Adapter) *UserMeta {
	var o UserMeta
	o._table = fmt.Sprintf("%susermeta", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "umeta_id"
	o._new = false
	return &o
}

func (o *UserMeta) Find(_find_by_uMetaId int64) (UserMeta, error) {

	var model_slice []UserMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "umeta_id", _find_by_uMetaId)
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
func (o *UserMeta) FindByUserId(_find_by_userId int64) ([]UserMeta, error) {

	var model_slice []UserMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "user_id", _find_by_userId)
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
func (o *UserMeta) FindByMetaKey(_find_by_metaKey string) ([]UserMeta, error) {

	var model_slice []UserMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "meta_key", _find_by_metaKey)
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
func (o *UserMeta) FindByMetaValue(_find_by_metaValue string) ([]UserMeta, error) {

	var model_slice []UserMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "meta_value", _find_by_metaValue)
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
	_uMetaId, err := m["umeta_id"].AsInt64()
	if err != nil {
		return err
	}
	o.uMetaId = _uMetaId
	_userId, err := m["user_id"].AsInt64()
	if err != nil {
		return err
	}
	o.userId = _userId
	_metaKey, err := m["meta_key"].AsString()
	if err != nil {
		return err
	}
	o.metaKey = _metaKey
	_metaValue, err := m["meta_value"].AsString()
	if err != nil {
		return err
	}
	o.metaValue = _metaValue

	return nil
}

func (o *UserMeta) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `user_id` = '%d', `meta_key` = '%s', `meta_value` = '%s' WHERE %s = '%d' LIMIT 1", o._table, o.userId, o.metaKey, o.metaValue, o._pkey, o.uMetaId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *UserMeta) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`user_id`, `meta_key`, `meta_value`) VALUES ('%d', '%s', '%s')", o._table, o.userId, o.metaKey, o.metaValue)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *UserMeta) UpdateUserId(_upd_userId int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_id` = '%d' WHERE `umeta_id` = '%d'", o._table, _upd_userId, o.userId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.userId = _upd_userId
	return o._adapter.AffectedRows(), nil
}

func (o *UserMeta) UpdateMetaKey(_upd_metaKey string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `meta_key` = '%s' WHERE `umeta_id` = '%d'", o._table, _upd_metaKey, o.metaKey)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.metaKey = _upd_metaKey
	return o._adapter.AffectedRows(), nil
}

func (o *UserMeta) UpdateMetaValue(_upd_metaValue string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `meta_value` = '%s' WHERE `umeta_id` = '%d'", o._table, _upd_metaValue, o.metaValue)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.metaValue = _upd_metaValue
	return o._adapter.AffectedRows(), nil
}

type User struct {
	_table            string
	_adapter          Adapter
	_pkey             string // 0 The name of the primary key in this table
	_conds            []string
	_new              bool
	iD                int64
	userLogin         string
	userPass          string
	userNicename      string
	userEmail         string
	userUrl           string
	userRegistered    DateTime
	userActivationKey string
	userStatus        int
	displayName       string
}

func NewUser(a Adapter) *User {
	var o User
	o._table = fmt.Sprintf("%susers", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "ID"
	o._new = false
	return &o
}

func (o *User) Find(_find_by_iD int64) (User, error) {

	var model_slice []User
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "ID", _find_by_iD)
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
func (o *User) FindByUserLogin(_find_by_userLogin string) ([]User, error) {

	var model_slice []User
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "user_login", _find_by_userLogin)
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
func (o *User) FindByUserPass(_find_by_userPass string) ([]User, error) {

	var model_slice []User
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "user_pass", _find_by_userPass)
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
func (o *User) FindByUserNicename(_find_by_userNicename string) ([]User, error) {

	var model_slice []User
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "user_nicename", _find_by_userNicename)
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
func (o *User) FindByUserEmail(_find_by_userEmail string) ([]User, error) {

	var model_slice []User
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "user_email", _find_by_userEmail)
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
func (o *User) FindByUserUrl(_find_by_userUrl string) ([]User, error) {

	var model_slice []User
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "user_url", _find_by_userUrl)
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
func (o *User) FindByUserRegistered(_find_by_userRegistered DateTime) ([]User, error) {

	var model_slice []User
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "user_registered", _find_by_userRegistered)
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
func (o *User) FindByUserActivationKey(_find_by_userActivationKey string) ([]User, error) {

	var model_slice []User
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "user_activation_key", _find_by_userActivationKey)
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
func (o *User) FindByUserStatus(_find_by_userStatus int) ([]User, error) {

	var model_slice []User
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "user_status", _find_by_userStatus)
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
func (o *User) FindByDisplayName(_find_by_displayName string) ([]User, error) {

	var model_slice []User
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "display_name", _find_by_displayName)
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
	_iD, err := m["ID"].AsInt64()
	if err != nil {
		return err
	}
	o.iD = _iD
	_userLogin, err := m["user_login"].AsString()
	if err != nil {
		return err
	}
	o.userLogin = _userLogin
	_userPass, err := m["user_pass"].AsString()
	if err != nil {
		return err
	}
	o.userPass = _userPass
	_userNicename, err := m["user_nicename"].AsString()
	if err != nil {
		return err
	}
	o.userNicename = _userNicename
	_userEmail, err := m["user_email"].AsString()
	if err != nil {
		return err
	}
	o.userEmail = _userEmail
	_userUrl, err := m["user_url"].AsString()
	if err != nil {
		return err
	}
	o.userUrl = _userUrl
	_userRegistered, err := m["user_registered"].AsDateTime()
	if err != nil {
		return err
	}
	o.userRegistered = _userRegistered
	_userActivationKey, err := m["user_activation_key"].AsString()
	if err != nil {
		return err
	}
	o.userActivationKey = _userActivationKey
	_userStatus, err := m["user_status"].AsInt()
	if err != nil {
		return err
	}
	o.userStatus = _userStatus
	_displayName, err := m["display_name"].AsString()
	if err != nil {
		return err
	}
	o.displayName = _displayName

	return nil
}

func (o *User) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `user_login` = '%s', `user_pass` = '%s', `user_nicename` = '%s', `user_email` = '%s', `user_url` = '%s', `user_registered` = '%s', `user_activation_key` = '%s', `user_status` = '%d', `display_name` = '%s' WHERE %s = '%d' LIMIT 1", o._table, o.userLogin, o.userPass, o.userNicename, o.userEmail, o.userUrl, o.userRegistered, o.userActivationKey, o.userStatus, o.displayName, o._pkey, o.iD)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *User) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`user_login`, `user_pass`, `user_nicename`, `user_email`, `user_url`, `user_registered`, `user_activation_key`, `user_status`, `display_name`) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%d', '%s')", o._table, o.userLogin, o.userPass, o.userNicename, o.userEmail, o.userUrl, o.userRegistered, o.userActivationKey, o.userStatus, o.displayName)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *User) UpdateUserLogin(_upd_userLogin string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_login` = '%s' WHERE `ID` = '%d'", o._table, _upd_userLogin, o.userLogin)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.userLogin = _upd_userLogin
	return o._adapter.AffectedRows(), nil
}

func (o *User) UpdateUserPass(_upd_userPass string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_pass` = '%s' WHERE `ID` = '%d'", o._table, _upd_userPass, o.userPass)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.userPass = _upd_userPass
	return o._adapter.AffectedRows(), nil
}

func (o *User) UpdateUserNicename(_upd_userNicename string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_nicename` = '%s' WHERE `ID` = '%d'", o._table, _upd_userNicename, o.userNicename)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.userNicename = _upd_userNicename
	return o._adapter.AffectedRows(), nil
}

func (o *User) UpdateUserEmail(_upd_userEmail string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_email` = '%s' WHERE `ID` = '%d'", o._table, _upd_userEmail, o.userEmail)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.userEmail = _upd_userEmail
	return o._adapter.AffectedRows(), nil
}

func (o *User) UpdateUserUrl(_upd_userUrl string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_url` = '%s' WHERE `ID` = '%d'", o._table, _upd_userUrl, o.userUrl)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.userUrl = _upd_userUrl
	return o._adapter.AffectedRows(), nil
}

func (o *User) UpdateUserRegistered(_upd_userRegistered DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_registered` = '%s' WHERE `ID` = '%d'", o._table, _upd_userRegistered, o.userRegistered)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.userRegistered = _upd_userRegistered
	return o._adapter.AffectedRows(), nil
}

func (o *User) UpdateUserActivationKey(_upd_userActivationKey string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_activation_key` = '%s' WHERE `ID` = '%d'", o._table, _upd_userActivationKey, o.userActivationKey)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.userActivationKey = _upd_userActivationKey
	return o._adapter.AffectedRows(), nil
}

func (o *User) UpdateUserStatus(_upd_userStatus int) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_status` = '%d' WHERE `ID` = '%d'", o._table, _upd_userStatus, o.userStatus)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.userStatus = _upd_userStatus
	return o._adapter.AffectedRows(), nil
}

func (o *User) UpdateDisplayName(_upd_displayName string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `display_name` = '%s' WHERE `ID` = '%d'", o._table, _upd_displayName, o.displayName)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.displayName = _upd_displayName
	return o._adapter.AffectedRows(), nil
}

type WooAttrTaxonomie struct {
	_table      string
	_adapter    Adapter
	_pkey       string // 0 The name of the primary key in this table
	_conds      []string
	_new        bool
	attrId      int64
	attrName    string
	attrLabel   string
	attrType    string
	attrOrderby string
}

func NewWooAttrTaxonomie(a Adapter) *WooAttrTaxonomie {
	var o WooAttrTaxonomie
	o._table = fmt.Sprintf("%swoocommerce_attribute_taxonomies", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "attribute_id"
	o._new = false
	return &o
}

func (o *WooAttrTaxonomie) Find(_find_by_attrId int64) (WooAttrTaxonomie, error) {

	var model_slice []WooAttrTaxonomie
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "attribute_id", _find_by_attrId)
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
func (o *WooAttrTaxonomie) FindByAttrName(_find_by_attrName string) ([]WooAttrTaxonomie, error) {

	var model_slice []WooAttrTaxonomie
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "attribute_name", _find_by_attrName)
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
func (o *WooAttrTaxonomie) FindByAttrLabel(_find_by_attrLabel string) ([]WooAttrTaxonomie, error) {

	var model_slice []WooAttrTaxonomie
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "attribute_label", _find_by_attrLabel)
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
func (o *WooAttrTaxonomie) FindByAttrType(_find_by_attrType string) ([]WooAttrTaxonomie, error) {

	var model_slice []WooAttrTaxonomie
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "attribute_type", _find_by_attrType)
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
func (o *WooAttrTaxonomie) FindByAttrOrderby(_find_by_attrOrderby string) ([]WooAttrTaxonomie, error) {

	var model_slice []WooAttrTaxonomie
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "attribute_orderby", _find_by_attrOrderby)
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
	_attrId, err := m["attribute_id"].AsInt64()
	if err != nil {
		return err
	}
	o.attrId = _attrId
	_attrName, err := m["attribute_name"].AsString()
	if err != nil {
		return err
	}
	o.attrName = _attrName
	_attrLabel, err := m["attribute_label"].AsString()
	if err != nil {
		return err
	}
	o.attrLabel = _attrLabel
	_attrType, err := m["attribute_type"].AsString()
	if err != nil {
		return err
	}
	o.attrType = _attrType
	_attrOrderby, err := m["attribute_orderby"].AsString()
	if err != nil {
		return err
	}
	o.attrOrderby = _attrOrderby

	return nil
}

func (o *WooAttrTaxonomie) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `attribute_name` = '%s', `attribute_label` = '%s', `attribute_type` = '%s', `attribute_orderby` = '%s' WHERE %s = '%d' LIMIT 1", o._table, o.attrName, o.attrLabel, o.attrType, o.attrOrderby, o._pkey, o.attrId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *WooAttrTaxonomie) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`attribute_name`, `attribute_label`, `attribute_type`, `attribute_orderby`) VALUES ('%s', '%s', '%s', '%s')", o._table, o.attrName, o.attrLabel, o.attrType, o.attrOrderby)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *WooAttrTaxonomie) UpdateAttrName(_upd_attrName string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `attribute_name` = '%s' WHERE `attribute_id` = '%d'", o._table, _upd_attrName, o.attrName)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.attrName = _upd_attrName
	return o._adapter.AffectedRows(), nil
}

func (o *WooAttrTaxonomie) UpdateAttrLabel(_upd_attrLabel string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `attribute_label` = '%s' WHERE `attribute_id` = '%d'", o._table, _upd_attrLabel, o.attrLabel)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.attrLabel = _upd_attrLabel
	return o._adapter.AffectedRows(), nil
}

func (o *WooAttrTaxonomie) UpdateAttrType(_upd_attrType string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `attribute_type` = '%s' WHERE `attribute_id` = '%d'", o._table, _upd_attrType, o.attrType)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.attrType = _upd_attrType
	return o._adapter.AffectedRows(), nil
}

func (o *WooAttrTaxonomie) UpdateAttrOrderby(_upd_attrOrderby string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `attribute_orderby` = '%s' WHERE `attribute_id` = '%d'", o._table, _upd_attrOrderby, o.attrOrderby)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.attrOrderby = _upd_attrOrderby
	return o._adapter.AffectedRows(), nil
}

type WooDownloadableProductPerm struct {
	_table             string
	_adapter           Adapter
	_pkey              string // 0 The name of the primary key in this table
	_conds             []string
	_new               bool
	permissionId       int64
	downloadId         string
	productId          int64
	orderId            int64
	orderKey           string
	userEmail          string
	userId             int64
	downloadsRemaining string
	accessGranted      DateTime
	accessExpires      DateTime
	downloadCount      int64
}

func NewWooDownloadableProductPerm(a Adapter) *WooDownloadableProductPerm {
	var o WooDownloadableProductPerm
	o._table = fmt.Sprintf("%swoocommerce_downloadable_product_permissions", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "permission_id"
	o._new = false
	return &o
}

func (o *WooDownloadableProductPerm) Find(_find_by_permissionId int64) (WooDownloadableProductPerm, error) {

	var model_slice []WooDownloadableProductPerm
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "permission_id", _find_by_permissionId)
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
func (o *WooDownloadableProductPerm) FindByDownloadId(_find_by_downloadId string) ([]WooDownloadableProductPerm, error) {

	var model_slice []WooDownloadableProductPerm
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "download_id", _find_by_downloadId)
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
func (o *WooDownloadableProductPerm) FindByProductId(_find_by_productId int64) ([]WooDownloadableProductPerm, error) {

	var model_slice []WooDownloadableProductPerm
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "product_id", _find_by_productId)
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
func (o *WooDownloadableProductPerm) FindByOrderId(_find_by_orderId int64) ([]WooDownloadableProductPerm, error) {

	var model_slice []WooDownloadableProductPerm
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "order_id", _find_by_orderId)
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
func (o *WooDownloadableProductPerm) FindByOrderKey(_find_by_orderKey string) ([]WooDownloadableProductPerm, error) {

	var model_slice []WooDownloadableProductPerm
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "order_key", _find_by_orderKey)
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
func (o *WooDownloadableProductPerm) FindByUserEmail(_find_by_userEmail string) ([]WooDownloadableProductPerm, error) {

	var model_slice []WooDownloadableProductPerm
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "user_email", _find_by_userEmail)
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
func (o *WooDownloadableProductPerm) FindByUserId(_find_by_userId int64) ([]WooDownloadableProductPerm, error) {

	var model_slice []WooDownloadableProductPerm
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "user_id", _find_by_userId)
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
func (o *WooDownloadableProductPerm) FindByDownloadsRemaining(_find_by_downloadsRemaining string) ([]WooDownloadableProductPerm, error) {

	var model_slice []WooDownloadableProductPerm
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "downloads_remaining", _find_by_downloadsRemaining)
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
func (o *WooDownloadableProductPerm) FindByAccessGranted(_find_by_accessGranted DateTime) ([]WooDownloadableProductPerm, error) {

	var model_slice []WooDownloadableProductPerm
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "access_granted", _find_by_accessGranted)
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
func (o *WooDownloadableProductPerm) FindByAccessExpires(_find_by_accessExpires DateTime) ([]WooDownloadableProductPerm, error) {

	var model_slice []WooDownloadableProductPerm
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "access_expires", _find_by_accessExpires)
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
func (o *WooDownloadableProductPerm) FindByDownloadCount(_find_by_downloadCount int64) ([]WooDownloadableProductPerm, error) {

	var model_slice []WooDownloadableProductPerm
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "download_count", _find_by_downloadCount)
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
	_permissionId, err := m["permission_id"].AsInt64()
	if err != nil {
		return err
	}
	o.permissionId = _permissionId
	_downloadId, err := m["download_id"].AsString()
	if err != nil {
		return err
	}
	o.downloadId = _downloadId
	_productId, err := m["product_id"].AsInt64()
	if err != nil {
		return err
	}
	o.productId = _productId
	_orderId, err := m["order_id"].AsInt64()
	if err != nil {
		return err
	}
	o.orderId = _orderId
	_orderKey, err := m["order_key"].AsString()
	if err != nil {
		return err
	}
	o.orderKey = _orderKey
	_userEmail, err := m["user_email"].AsString()
	if err != nil {
		return err
	}
	o.userEmail = _userEmail
	_userId, err := m["user_id"].AsInt64()
	if err != nil {
		return err
	}
	o.userId = _userId
	_downloadsRemaining, err := m["downloads_remaining"].AsString()
	if err != nil {
		return err
	}
	o.downloadsRemaining = _downloadsRemaining
	_accessGranted, err := m["access_granted"].AsDateTime()
	if err != nil {
		return err
	}
	o.accessGranted = _accessGranted
	_accessExpires, err := m["access_expires"].AsDateTime()
	if err != nil {
		return err
	}
	o.accessExpires = _accessExpires
	_downloadCount, err := m["download_count"].AsInt64()
	if err != nil {
		return err
	}
	o.downloadCount = _downloadCount

	return nil
}

func (o *WooDownloadableProductPerm) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `download_id` = '%s', `product_id` = '%d', `order_id` = '%d', `order_key` = '%s', `user_email` = '%s', `user_id` = '%d', `downloads_remaining` = '%s', `access_granted` = '%s', `access_expires` = '%s', `download_count` = '%d' WHERE %s = '%d' LIMIT 1", o._table, o.downloadId, o.productId, o.orderId, o.orderKey, o.userEmail, o.userId, o.downloadsRemaining, o.accessGranted, o.accessExpires, o.downloadCount, o._pkey, o.permissionId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *WooDownloadableProductPerm) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`download_id`, `product_id`, `order_id`, `order_key`, `user_email`, `user_id`, `downloads_remaining`, `access_granted`, `access_expires`, `download_count`) VALUES ('%s', '%d', '%d', '%s', '%s', '%d', '%s', '%s', '%s', '%d')", o._table, o.downloadId, o.productId, o.orderId, o.orderKey, o.userEmail, o.userId, o.downloadsRemaining, o.accessGranted, o.accessExpires, o.downloadCount)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *WooDownloadableProductPerm) UpdateDownloadId(_upd_downloadId string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `download_id` = '%s' WHERE `permission_id` = '%d'", o._table, _upd_downloadId, o.downloadId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.downloadId = _upd_downloadId
	return o._adapter.AffectedRows(), nil
}

func (o *WooDownloadableProductPerm) UpdateProductId(_upd_productId int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `product_id` = '%d' WHERE `permission_id` = '%d'", o._table, _upd_productId, o.productId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.productId = _upd_productId
	return o._adapter.AffectedRows(), nil
}

func (o *WooDownloadableProductPerm) UpdateOrderId(_upd_orderId int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `order_id` = '%d' WHERE `permission_id` = '%d'", o._table, _upd_orderId, o.orderId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.orderId = _upd_orderId
	return o._adapter.AffectedRows(), nil
}

func (o *WooDownloadableProductPerm) UpdateOrderKey(_upd_orderKey string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `order_key` = '%s' WHERE `permission_id` = '%d'", o._table, _upd_orderKey, o.orderKey)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.orderKey = _upd_orderKey
	return o._adapter.AffectedRows(), nil
}

func (o *WooDownloadableProductPerm) UpdateUserEmail(_upd_userEmail string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_email` = '%s' WHERE `permission_id` = '%d'", o._table, _upd_userEmail, o.userEmail)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.userEmail = _upd_userEmail
	return o._adapter.AffectedRows(), nil
}

func (o *WooDownloadableProductPerm) UpdateUserId(_upd_userId int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `user_id` = '%d' WHERE `permission_id` = '%d'", o._table, _upd_userId, o.userId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.userId = _upd_userId
	return o._adapter.AffectedRows(), nil
}

func (o *WooDownloadableProductPerm) UpdateDownloadsRemaining(_upd_downloadsRemaining string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `downloads_remaining` = '%s' WHERE `permission_id` = '%d'", o._table, _upd_downloadsRemaining, o.downloadsRemaining)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.downloadsRemaining = _upd_downloadsRemaining
	return o._adapter.AffectedRows(), nil
}

func (o *WooDownloadableProductPerm) UpdateAccessGranted(_upd_accessGranted DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `access_granted` = '%s' WHERE `permission_id` = '%d'", o._table, _upd_accessGranted, o.accessGranted)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.accessGranted = _upd_accessGranted
	return o._adapter.AffectedRows(), nil
}

func (o *WooDownloadableProductPerm) UpdateAccessExpires(_upd_accessExpires DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `access_expires` = '%s' WHERE `permission_id` = '%d'", o._table, _upd_accessExpires, o.accessExpires)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.accessExpires = _upd_accessExpires
	return o._adapter.AffectedRows(), nil
}

func (o *WooDownloadableProductPerm) UpdateDownloadCount(_upd_downloadCount int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `download_count` = '%d' WHERE `permission_id` = '%d'", o._table, _upd_downloadCount, o.downloadCount)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.downloadCount = _upd_downloadCount
	return o._adapter.AffectedRows(), nil
}

type WooOrderItemMeta struct {
	_table      string
	_adapter    Adapter
	_pkey       string // 0 The name of the primary key in this table
	_conds      []string
	_new        bool
	metaId      int64
	orderItemId int64
	metaKey     string
	metaValue   string
}

func NewWooOrderItemMeta(a Adapter) *WooOrderItemMeta {
	var o WooOrderItemMeta
	o._table = fmt.Sprintf("%swoocommerce_order_itemmeta", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "meta_id"
	o._new = false
	return &o
}

func (o *WooOrderItemMeta) Find(_find_by_metaId int64) (WooOrderItemMeta, error) {

	var model_slice []WooOrderItemMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "meta_id", _find_by_metaId)
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
func (o *WooOrderItemMeta) FindByOrderItemId(_find_by_orderItemId int64) ([]WooOrderItemMeta, error) {

	var model_slice []WooOrderItemMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "order_item_id", _find_by_orderItemId)
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
func (o *WooOrderItemMeta) FindByMetaKey(_find_by_metaKey string) ([]WooOrderItemMeta, error) {

	var model_slice []WooOrderItemMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "meta_key", _find_by_metaKey)
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
func (o *WooOrderItemMeta) FindByMetaValue(_find_by_metaValue string) ([]WooOrderItemMeta, error) {

	var model_slice []WooOrderItemMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "meta_value", _find_by_metaValue)
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
	_metaId, err := m["meta_id"].AsInt64()
	if err != nil {
		return err
	}
	o.metaId = _metaId
	_orderItemId, err := m["order_item_id"].AsInt64()
	if err != nil {
		return err
	}
	o.orderItemId = _orderItemId
	_metaKey, err := m["meta_key"].AsString()
	if err != nil {
		return err
	}
	o.metaKey = _metaKey
	_metaValue, err := m["meta_value"].AsString()
	if err != nil {
		return err
	}
	o.metaValue = _metaValue

	return nil
}

func (o *WooOrderItemMeta) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `order_item_id` = '%d', `meta_key` = '%s', `meta_value` = '%s' WHERE %s = '%d' LIMIT 1", o._table, o.orderItemId, o.metaKey, o.metaValue, o._pkey, o.metaId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *WooOrderItemMeta) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`order_item_id`, `meta_key`, `meta_value`) VALUES ('%d', '%s', '%s')", o._table, o.orderItemId, o.metaKey, o.metaValue)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *WooOrderItemMeta) UpdateOrderItemId(_upd_orderItemId int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `order_item_id` = '%d' WHERE `meta_id` = '%d'", o._table, _upd_orderItemId, o.orderItemId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.orderItemId = _upd_orderItemId
	return o._adapter.AffectedRows(), nil
}

func (o *WooOrderItemMeta) UpdateMetaKey(_upd_metaKey string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `meta_key` = '%s' WHERE `meta_id` = '%d'", o._table, _upd_metaKey, o.metaKey)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.metaKey = _upd_metaKey
	return o._adapter.AffectedRows(), nil
}

func (o *WooOrderItemMeta) UpdateMetaValue(_upd_metaValue string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `meta_value` = '%s' WHERE `meta_id` = '%d'", o._table, _upd_metaValue, o.metaValue)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.metaValue = _upd_metaValue
	return o._adapter.AffectedRows(), nil
}

type WooOrderItem struct {
	_table        string
	_adapter      Adapter
	_pkey         string // 0 The name of the primary key in this table
	_conds        []string
	_new          bool
	orderItemId   int64
	orderItemName string
	orderItemType string
	orderId       int64
}

func NewWooOrderItem(a Adapter) *WooOrderItem {
	var o WooOrderItem
	o._table = fmt.Sprintf("%swoocommerce_order_items", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "order_item_id"
	o._new = false
	return &o
}

func (o *WooOrderItem) Find(_find_by_orderItemId int64) (WooOrderItem, error) {

	var model_slice []WooOrderItem
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "order_item_id", _find_by_orderItemId)
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
func (o *WooOrderItem) FindByOrderItemName(_find_by_orderItemName string) ([]WooOrderItem, error) {

	var model_slice []WooOrderItem
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "order_item_name", _find_by_orderItemName)
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
func (o *WooOrderItem) FindByOrderItemType(_find_by_orderItemType string) ([]WooOrderItem, error) {

	var model_slice []WooOrderItem
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "order_item_type", _find_by_orderItemType)
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
func (o *WooOrderItem) FindByOrderId(_find_by_orderId int64) ([]WooOrderItem, error) {

	var model_slice []WooOrderItem
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "order_id", _find_by_orderId)
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
	_orderItemId, err := m["order_item_id"].AsInt64()
	if err != nil {
		return err
	}
	o.orderItemId = _orderItemId
	_orderItemName, err := m["order_item_name"].AsString()
	if err != nil {
		return err
	}
	o.orderItemName = _orderItemName
	_orderItemType, err := m["order_item_type"].AsString()
	if err != nil {
		return err
	}
	o.orderItemType = _orderItemType
	_orderId, err := m["order_id"].AsInt64()
	if err != nil {
		return err
	}
	o.orderId = _orderId

	return nil
}

func (o *WooOrderItem) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `order_item_name` = '%s', `order_item_type` = '%s', `order_id` = '%d' WHERE %s = '%d' LIMIT 1", o._table, o.orderItemName, o.orderItemType, o.orderId, o._pkey, o.orderItemId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *WooOrderItem) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`order_item_name`, `order_item_type`, `order_id`) VALUES ('%s', '%s', '%d')", o._table, o.orderItemName, o.orderItemType, o.orderId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *WooOrderItem) UpdateOrderItemName(_upd_orderItemName string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `order_item_name` = '%s' WHERE `order_item_id` = '%d'", o._table, _upd_orderItemName, o.orderItemName)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.orderItemName = _upd_orderItemName
	return o._adapter.AffectedRows(), nil
}

func (o *WooOrderItem) UpdateOrderItemType(_upd_orderItemType string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `order_item_type` = '%s' WHERE `order_item_id` = '%d'", o._table, _upd_orderItemType, o.orderItemType)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.orderItemType = _upd_orderItemType
	return o._adapter.AffectedRows(), nil
}

func (o *WooOrderItem) UpdateOrderId(_upd_orderId int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `order_id` = '%d' WHERE `order_item_id` = '%d'", o._table, _upd_orderId, o.orderId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.orderId = _upd_orderId
	return o._adapter.AffectedRows(), nil
}

type WooTaxRateLocation struct {
	_table       string
	_adapter     Adapter
	_pkey        string // 0 The name of the primary key in this table
	_conds       []string
	_new         bool
	locationId   int64
	locationCode string
	taxRateId    int64
	locationType string
}

func NewWooTaxRateLocation(a Adapter) *WooTaxRateLocation {
	var o WooTaxRateLocation
	o._table = fmt.Sprintf("%swoocommerce_tax_rate_locations", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "location_id"
	o._new = false
	return &o
}

func (o *WooTaxRateLocation) Find(_find_by_locationId int64) (WooTaxRateLocation, error) {

	var model_slice []WooTaxRateLocation
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "location_id", _find_by_locationId)
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
func (o *WooTaxRateLocation) FindByLocationCode(_find_by_locationCode string) ([]WooTaxRateLocation, error) {

	var model_slice []WooTaxRateLocation
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "location_code", _find_by_locationCode)
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
func (o *WooTaxRateLocation) FindByTaxRateId(_find_by_taxRateId int64) ([]WooTaxRateLocation, error) {

	var model_slice []WooTaxRateLocation
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "tax_rate_id", _find_by_taxRateId)
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
func (o *WooTaxRateLocation) FindByLocationType(_find_by_locationType string) ([]WooTaxRateLocation, error) {

	var model_slice []WooTaxRateLocation
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "location_type", _find_by_locationType)
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
	_locationId, err := m["location_id"].AsInt64()
	if err != nil {
		return err
	}
	o.locationId = _locationId
	_locationCode, err := m["location_code"].AsString()
	if err != nil {
		return err
	}
	o.locationCode = _locationCode
	_taxRateId, err := m["tax_rate_id"].AsInt64()
	if err != nil {
		return err
	}
	o.taxRateId = _taxRateId
	_locationType, err := m["location_type"].AsString()
	if err != nil {
		return err
	}
	o.locationType = _locationType

	return nil
}

func (o *WooTaxRateLocation) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `location_code` = '%s', `tax_rate_id` = '%d', `location_type` = '%s' WHERE %s = '%d' LIMIT 1", o._table, o.locationCode, o.taxRateId, o.locationType, o._pkey, o.locationId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *WooTaxRateLocation) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`location_code`, `tax_rate_id`, `location_type`) VALUES ('%s', '%d', '%s')", o._table, o.locationCode, o.taxRateId, o.locationType)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRateLocation) UpdateLocationCode(_upd_locationCode string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `location_code` = '%s' WHERE `location_id` = '%d'", o._table, _upd_locationCode, o.locationCode)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.locationCode = _upd_locationCode
	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRateLocation) UpdateTaxRateId(_upd_taxRateId int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `tax_rate_id` = '%d' WHERE `location_id` = '%d'", o._table, _upd_taxRateId, o.taxRateId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.taxRateId = _upd_taxRateId
	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRateLocation) UpdateLocationType(_upd_locationType string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `location_type` = '%s' WHERE `location_id` = '%d'", o._table, _upd_locationType, o.locationType)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.locationType = _upd_locationType
	return o._adapter.AffectedRows(), nil
}

type WooTaxRate struct {
	_table          string
	_adapter        Adapter
	_pkey           string // 0 The name of the primary key in this table
	_conds          []string
	_new            bool
	taxRateId       int64
	taxRateCountry  string
	taxRateState    string
	taxRate         string
	taxRateName     string
	taxRatePriority int64
	taxRateCompound int
	taxRateShipping int
	taxRateOrder    int64
	taxRateClass    string
}

func NewWooTaxRate(a Adapter) *WooTaxRate {
	var o WooTaxRate
	o._table = fmt.Sprintf("%swoocommerce_tax_rates", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "tax_rate_id"
	o._new = false
	return &o
}

func (o *WooTaxRate) Find(_find_by_taxRateId int64) (WooTaxRate, error) {

	var model_slice []WooTaxRate
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "tax_rate_id", _find_by_taxRateId)
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
func (o *WooTaxRate) FindByTaxRateCountry(_find_by_taxRateCountry string) ([]WooTaxRate, error) {

	var model_slice []WooTaxRate
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "tax_rate_country", _find_by_taxRateCountry)
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
func (o *WooTaxRate) FindByTaxRateState(_find_by_taxRateState string) ([]WooTaxRate, error) {

	var model_slice []WooTaxRate
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "tax_rate_state", _find_by_taxRateState)
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
func (o *WooTaxRate) FindByTaxRate(_find_by_taxRate string) ([]WooTaxRate, error) {

	var model_slice []WooTaxRate
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "tax_rate", _find_by_taxRate)
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
func (o *WooTaxRate) FindByTaxRateName(_find_by_taxRateName string) ([]WooTaxRate, error) {

	var model_slice []WooTaxRate
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "tax_rate_name", _find_by_taxRateName)
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
func (o *WooTaxRate) FindByTaxRatePriority(_find_by_taxRatePriority int64) ([]WooTaxRate, error) {

	var model_slice []WooTaxRate
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "tax_rate_priority", _find_by_taxRatePriority)
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
func (o *WooTaxRate) FindByTaxRateCompound(_find_by_taxRateCompound int) ([]WooTaxRate, error) {

	var model_slice []WooTaxRate
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "tax_rate_compound", _find_by_taxRateCompound)
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
func (o *WooTaxRate) FindByTaxRateShipping(_find_by_taxRateShipping int) ([]WooTaxRate, error) {

	var model_slice []WooTaxRate
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "tax_rate_shipping", _find_by_taxRateShipping)
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
func (o *WooTaxRate) FindByTaxRateOrder(_find_by_taxRateOrder int64) ([]WooTaxRate, error) {

	var model_slice []WooTaxRate
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "tax_rate_order", _find_by_taxRateOrder)
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
func (o *WooTaxRate) FindByTaxRateClass(_find_by_taxRateClass string) ([]WooTaxRate, error) {

	var model_slice []WooTaxRate
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "tax_rate_class", _find_by_taxRateClass)
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
	_taxRateId, err := m["tax_rate_id"].AsInt64()
	if err != nil {
		return err
	}
	o.taxRateId = _taxRateId
	_taxRateCountry, err := m["tax_rate_country"].AsString()
	if err != nil {
		return err
	}
	o.taxRateCountry = _taxRateCountry
	_taxRateState, err := m["tax_rate_state"].AsString()
	if err != nil {
		return err
	}
	o.taxRateState = _taxRateState
	_taxRate, err := m["tax_rate"].AsString()
	if err != nil {
		return err
	}
	o.taxRate = _taxRate
	_taxRateName, err := m["tax_rate_name"].AsString()
	if err != nil {
		return err
	}
	o.taxRateName = _taxRateName
	_taxRatePriority, err := m["tax_rate_priority"].AsInt64()
	if err != nil {
		return err
	}
	o.taxRatePriority = _taxRatePriority
	_taxRateCompound, err := m["tax_rate_compound"].AsInt()
	if err != nil {
		return err
	}
	o.taxRateCompound = _taxRateCompound
	_taxRateShipping, err := m["tax_rate_shipping"].AsInt()
	if err != nil {
		return err
	}
	o.taxRateShipping = _taxRateShipping
	_taxRateOrder, err := m["tax_rate_order"].AsInt64()
	if err != nil {
		return err
	}
	o.taxRateOrder = _taxRateOrder
	_taxRateClass, err := m["tax_rate_class"].AsString()
	if err != nil {
		return err
	}
	o.taxRateClass = _taxRateClass

	return nil
}

func (o *WooTaxRate) Save() (int64, error) {
	if o._new == true {
		return o.Create()
	}
	frmt := fmt.Sprintf("UPDATE %s SET `tax_rate_country` = '%s', `tax_rate_state` = '%s', `tax_rate` = '%s', `tax_rate_name` = '%s', `tax_rate_priority` = '%d', `tax_rate_compound` = '%d', `tax_rate_shipping` = '%d', `tax_rate_order` = '%d', `tax_rate_class` = '%s' WHERE %s = '%d' LIMIT 1", o._table, o.taxRateCountry, o.taxRateState, o.taxRate, o.taxRateName, o.taxRatePriority, o.taxRateCompound, o.taxRateShipping, o.taxRateOrder, o.taxRateClass, o._pkey, o.taxRateId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}
func (o *WooTaxRate) Create() (int64, error) {
	frmt := fmt.Sprintf("INSERT INTO %s (`tax_rate_country`, `tax_rate_state`, `tax_rate`, `tax_rate_name`, `tax_rate_priority`, `tax_rate_compound`, `tax_rate_shipping`, `tax_rate_order`, `tax_rate_class`) VALUES ('%s', '%s', '%s', '%s', '%d', '%d', '%d', '%d', '%s')", o._table, o.taxRateCountry, o.taxRateState, o.taxRate, o.taxRateName, o.taxRatePriority, o.taxRateCompound, o.taxRateShipping, o.taxRateOrder, o.taxRateClass)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}

	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRate) UpdateTaxRateCountry(_upd_taxRateCountry string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `tax_rate_country` = '%s' WHERE `tax_rate_id` = '%d'", o._table, _upd_taxRateCountry, o.taxRateCountry)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.taxRateCountry = _upd_taxRateCountry
	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRate) UpdateTaxRateState(_upd_taxRateState string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `tax_rate_state` = '%s' WHERE `tax_rate_id` = '%d'", o._table, _upd_taxRateState, o.taxRateState)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.taxRateState = _upd_taxRateState
	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRate) UpdateTaxRate(_upd_taxRate string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `tax_rate` = '%s' WHERE `tax_rate_id` = '%d'", o._table, _upd_taxRate, o.taxRate)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.taxRate = _upd_taxRate
	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRate) UpdateTaxRateName(_upd_taxRateName string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `tax_rate_name` = '%s' WHERE `tax_rate_id` = '%d'", o._table, _upd_taxRateName, o.taxRateName)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.taxRateName = _upd_taxRateName
	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRate) UpdateTaxRatePriority(_upd_taxRatePriority int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `tax_rate_priority` = '%d' WHERE `tax_rate_id` = '%d'", o._table, _upd_taxRatePriority, o.taxRatePriority)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.taxRatePriority = _upd_taxRatePriority
	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRate) UpdateTaxRateCompound(_upd_taxRateCompound int) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `tax_rate_compound` = '%d' WHERE `tax_rate_id` = '%d'", o._table, _upd_taxRateCompound, o.taxRateCompound)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.taxRateCompound = _upd_taxRateCompound
	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRate) UpdateTaxRateShipping(_upd_taxRateShipping int) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `tax_rate_shipping` = '%d' WHERE `tax_rate_id` = '%d'", o._table, _upd_taxRateShipping, o.taxRateShipping)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.taxRateShipping = _upd_taxRateShipping
	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRate) UpdateTaxRateOrder(_upd_taxRateOrder int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `tax_rate_order` = '%d' WHERE `tax_rate_id` = '%d'", o._table, _upd_taxRateOrder, o.taxRateOrder)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.taxRateOrder = _upd_taxRateOrder
	return o._adapter.AffectedRows(), nil
}

func (o *WooTaxRate) UpdateTaxRateClass(_upd_taxRateClass string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `tax_rate_class` = '%s' WHERE `tax_rate_id` = '%d'", o._table, _upd_taxRateClass, o.taxRateClass)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.taxRateClass = _upd_taxRateClass
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
