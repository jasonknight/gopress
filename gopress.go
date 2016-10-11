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
	_table_   string
	_adapter_ Adapter
	_pkey_    string // 0 The name of the primary key in this table
	_conds_   []string
	metaId    int64
	commentId int64
	metaKey   string
	metaValue string
}

func NewCommentMeta(a Adapter) *CommentMeta {
	var o CommentMeta
	o._table_ = "wp_commentmeta"
	o._adapter_ = a
	o._pkey_ = "meta_id"
	return &o
}

func (o *CommentMeta) Find(_find_by_metaId int64) (CommentMeta, error) {

	var model_slice []CommentMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "meta_id", _find_by_metaId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "comment_id", _find_by_commentId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "meta_key", _find_by_metaKey)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "meta_value", _find_by_metaValue)
	results, err := o._adapter_.Query(q)
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

type Comment struct {
	_table_            string
	_adapter_          Adapter
	_pkey_             string // 0 The name of the primary key in this table
	_conds_            []string
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
	o._table_ = "wp_comments"
	o._adapter_ = a
	o._pkey_ = "comment_ID"
	return &o
}

func (o *Comment) Find(_find_by_commentID int64) (Comment, error) {

	var model_slice []Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "comment_ID", _find_by_commentID)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "comment_post_ID", _find_by_commentPostID)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "comment_author", _find_by_commentAuthor)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "comment_author_email", _find_by_commentAuthorEmail)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "comment_author_url", _find_by_commentAuthorUrl)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "comment_author_IP", _find_by_commentAuthorIP)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "comment_date", _find_by_commentDate)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "comment_date_gmt", _find_by_commentDateGmt)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "comment_content", _find_by_commentContent)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "comment_karma", _find_by_commentKarma)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "comment_approved", _find_by_commentApproved)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "comment_agent", _find_by_commentAgent)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "comment_type", _find_by_commentType)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "comment_parent", _find_by_commentParent)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "user_id", _find_by_userId)
	results, err := o._adapter_.Query(q)
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

type Link struct {
	_table_         string
	_adapter_       Adapter
	_pkey_          string // 0 The name of the primary key in this table
	_conds_         []string
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
	o._table_ = "wp_links"
	o._adapter_ = a
	o._pkey_ = "link_id"
	return &o
}

func (o *Link) Find(_find_by_linkId int64) (Link, error) {

	var model_slice []Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "link_id", _find_by_linkId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "link_url", _find_by_linkUrl)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "link_name", _find_by_linkName)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "link_image", _find_by_linkImage)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "link_target", _find_by_linkTarget)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "link_description", _find_by_linkDescription)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "link_visible", _find_by_linkVisible)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "link_owner", _find_by_linkOwner)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "link_rating", _find_by_linkRating)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "link_updated", _find_by_linkUpdated)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "link_rel", _find_by_linkRel)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "link_notes", _find_by_linkNotes)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "link_rss", _find_by_linkRss)
	results, err := o._adapter_.Query(q)
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

type Option struct {
	_table_     string
	_adapter_   Adapter
	_pkey_      string // 0 The name of the primary key in this table
	_conds_     []string
	optionId    int64
	optionName  string
	optionValue string
	autoload    string
}

func NewOption(a Adapter) *Option {
	var o Option
	o._table_ = "wp_options"
	o._adapter_ = a
	o._pkey_ = "option_id"
	return &o
}

func (o *Option) Find(_find_by_optionId int64) (Option, error) {

	var model_slice []Option
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "option_id", _find_by_optionId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "option_name", _find_by_optionName)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "option_value", _find_by_optionValue)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "autoload", _find_by_autoload)
	results, err := o._adapter_.Query(q)
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

type PostMeta struct {
	_table_   string
	_adapter_ Adapter
	_pkey_    string // 0 The name of the primary key in this table
	_conds_   []string
	metaId    int64
	stId      int64
	metaKey   string
	metaValue string
}

func NewPostMeta(a Adapter) *PostMeta {
	var o PostMeta
	o._table_ = "wp_postmeta"
	o._adapter_ = a
	o._pkey_ = "meta_id"
	return &o
}

func (o *PostMeta) Find(_find_by_metaId int64) (PostMeta, error) {

	var model_slice []PostMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "meta_id", _find_by_metaId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "post_id", _find_by_stId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "meta_key", _find_by_metaKey)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "meta_value", _find_by_metaValue)
	results, err := o._adapter_.Query(q)
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

type Post struct {
	_table_           string
	_adapter_         Adapter
	_pkey_            string // 0 The name of the primary key in this table
	_conds_           []string
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
	o._table_ = "wp_posts"
	o._adapter_ = a
	o._pkey_ = "ID"
	return &o
}

func (o *Post) Find(_find_by_iD int64) (Post, error) {

	var model_slice []Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "ID", _find_by_iD)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "post_author", _find_by_stAuthor)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "post_date", _find_by_stDate)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "post_date_gmt", _find_by_stDateGmt)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "post_content", _find_by_stContent)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "post_title", _find_by_stTitle)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "post_excerpt", _find_by_stExcerpt)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "post_status", _find_by_stStatus)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "comment_status", _find_by_commentStatus)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "ping_status", _find_by_pingStatus)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "post_password", _find_by_stPassword)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "post_name", _find_by_stName)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "to_ping", _find_by_toPing)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "pinged", _find_by_pinged)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "post_modified", _find_by_stModified)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "post_modified_gmt", _find_by_stModifiedGmt)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "post_content_filtered", _find_by_stContentFiltered)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "post_parent", _find_by_stParent)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "guid", _find_by_guid)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "menu_order", _find_by_menuOrder)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "post_type", _find_by_stType)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "post_mime_type", _find_by_stMimeType)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "comment_count", _find_by_commentCount)
	results, err := o._adapter_.Query(q)
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

type TermRelationship struct {
	_table_        string
	_adapter_      Adapter
	_pkey_         string // 0 The name of the primary key in this table
	_conds_        []string
	objectId       int64
	termTaxonomyId int64
	termOrder      int
}

func NewTermRelationship(a Adapter) *TermRelationship {
	var o TermRelationship
	o._table_ = "wp_term_relationships"
	o._adapter_ = a
	o._pkey_ = "term_taxonomy_id"
	return &o
}

func (o *TermRelationship) FindByObjectId(_find_by_objectId int64) ([]TermRelationship, error) {

	var model_slice []TermRelationship
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "object_id", _find_by_objectId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "term_taxonomy_id", _find_by_termTaxonomyId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "term_order", _find_by_termOrder)
	results, err := o._adapter_.Query(q)
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

type TermTaxonomy struct {
	_table_        string
	_adapter_      Adapter
	_pkey_         string // 0 The name of the primary key in this table
	_conds_        []string
	termTaxonomyId int64
	termId         int64
	taxonomy       string
	description    string
	parent         int64
	count          int64
}

func NewTermTaxonomy(a Adapter) *TermTaxonomy {
	var o TermTaxonomy
	o._table_ = "wp_term_taxonomy"
	o._adapter_ = a
	o._pkey_ = "term_taxonomy_id"
	return &o
}

func (o *TermTaxonomy) Find(_find_by_termTaxonomyId int64) (TermTaxonomy, error) {

	var model_slice []TermTaxonomy
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "term_taxonomy_id", _find_by_termTaxonomyId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "term_id", _find_by_termId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "taxonomy", _find_by_taxonomy)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "description", _find_by_description)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "parent", _find_by_parent)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "count", _find_by_count)
	results, err := o._adapter_.Query(q)
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

type Term struct {
	_table_   string
	_adapter_ Adapter
	_pkey_    string // 0 The name of the primary key in this table
	_conds_   []string
	termId    int64
	name      string
	slug      string
	termGroup int64
}

func NewTerm(a Adapter) *Term {
	var o Term
	o._table_ = "wp_terms"
	o._adapter_ = a
	o._pkey_ = "term_id"
	return &o
}

func (o *Term) Find(_find_by_termId int64) (Term, error) {

	var model_slice []Term
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "term_id", _find_by_termId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "name", _find_by_name)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "slug", _find_by_slug)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "term_group", _find_by_termGroup)
	results, err := o._adapter_.Query(q)
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

type UserMeta struct {
	_table_   string
	_adapter_ Adapter
	_pkey_    string // 0 The name of the primary key in this table
	_conds_   []string
	uMetaId   int64
	userId    int64
	metaKey   string
	metaValue string
}

func NewUserMeta(a Adapter) *UserMeta {
	var o UserMeta
	o._table_ = "wp_usermeta"
	o._adapter_ = a
	o._pkey_ = "umeta_id"
	return &o
}

func (o *UserMeta) Find(_find_by_uMetaId int64) (UserMeta, error) {

	var model_slice []UserMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "umeta_id", _find_by_uMetaId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "user_id", _find_by_userId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "meta_key", _find_by_metaKey)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "meta_value", _find_by_metaValue)
	results, err := o._adapter_.Query(q)
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

type User struct {
	_table_           string
	_adapter_         Adapter
	_pkey_            string // 0 The name of the primary key in this table
	_conds_           []string
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
	o._table_ = "wp_users"
	o._adapter_ = a
	o._pkey_ = "ID"
	return &o
}

func (o *User) Find(_find_by_iD int64) (User, error) {

	var model_slice []User
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "ID", _find_by_iD)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "user_login", _find_by_userLogin)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "user_pass", _find_by_userPass)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "user_nicename", _find_by_userNicename)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "user_email", _find_by_userEmail)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "user_url", _find_by_userUrl)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "user_registered", _find_by_userRegistered)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "user_activation_key", _find_by_userActivationKey)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "user_status", _find_by_userStatus)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "display_name", _find_by_displayName)
	results, err := o._adapter_.Query(q)
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

type WooAttrTaxonomie struct {
	_table_     string
	_adapter_   Adapter
	_pkey_      string // 0 The name of the primary key in this table
	_conds_     []string
	attrId      int64
	attrName    string
	attrLabel   string
	attrType    string
	attrOrderby string
}

func NewWooAttrTaxonomie(a Adapter) *WooAttrTaxonomie {
	var o WooAttrTaxonomie
	o._table_ = "wp_woocommerce_attribute_taxonomies"
	o._adapter_ = a
	o._pkey_ = "attribute_id"
	return &o
}

func (o *WooAttrTaxonomie) Find(_find_by_attrId int64) (WooAttrTaxonomie, error) {

	var model_slice []WooAttrTaxonomie
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "attribute_id", _find_by_attrId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "attribute_name", _find_by_attrName)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "attribute_label", _find_by_attrLabel)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "attribute_type", _find_by_attrType)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "attribute_orderby", _find_by_attrOrderby)
	results, err := o._adapter_.Query(q)
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

type WooDownloadableProductPerm struct {
	_table_            string
	_adapter_          Adapter
	_pkey_             string // 0 The name of the primary key in this table
	_conds_            []string
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
	o._table_ = "wp_woocommerce_downloadable_product_permissions"
	o._adapter_ = a
	o._pkey_ = "permission_id"
	return &o
}

func (o *WooDownloadableProductPerm) Find(_find_by_permissionId int64) (WooDownloadableProductPerm, error) {

	var model_slice []WooDownloadableProductPerm
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "permission_id", _find_by_permissionId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "download_id", _find_by_downloadId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "product_id", _find_by_productId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "order_id", _find_by_orderId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "order_key", _find_by_orderKey)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "user_email", _find_by_userEmail)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "user_id", _find_by_userId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "downloads_remaining", _find_by_downloadsRemaining)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "access_granted", _find_by_accessGranted)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "access_expires", _find_by_accessExpires)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "download_count", _find_by_downloadCount)
	results, err := o._adapter_.Query(q)
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

type WooOrderItemMeta struct {
	_table_     string
	_adapter_   Adapter
	_pkey_      string // 0 The name of the primary key in this table
	_conds_     []string
	metaId      int64
	orderItemId int64
	metaKey     string
	metaValue   string
}

func NewWooOrderItemMeta(a Adapter) *WooOrderItemMeta {
	var o WooOrderItemMeta
	o._table_ = "wp_woocommerce_order_itemmeta"
	o._adapter_ = a
	o._pkey_ = "meta_id"
	return &o
}

func (o *WooOrderItemMeta) Find(_find_by_metaId int64) (WooOrderItemMeta, error) {

	var model_slice []WooOrderItemMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "meta_id", _find_by_metaId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "order_item_id", _find_by_orderItemId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "meta_key", _find_by_metaKey)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "meta_value", _find_by_metaValue)
	results, err := o._adapter_.Query(q)
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

type WooOrderItem struct {
	_table_       string
	_adapter_     Adapter
	_pkey_        string // 0 The name of the primary key in this table
	_conds_       []string
	orderItemId   int64
	orderItemName string
	orderItemType string
	orderId       int64
}

func NewWooOrderItem(a Adapter) *WooOrderItem {
	var o WooOrderItem
	o._table_ = "wp_woocommerce_order_items"
	o._adapter_ = a
	o._pkey_ = "order_item_id"
	return &o
}

func (o *WooOrderItem) Find(_find_by_orderItemId int64) (WooOrderItem, error) {

	var model_slice []WooOrderItem
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "order_item_id", _find_by_orderItemId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "order_item_name", _find_by_orderItemName)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "order_item_type", _find_by_orderItemType)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "order_id", _find_by_orderId)
	results, err := o._adapter_.Query(q)
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

type WooTaxRateLocation struct {
	_table_      string
	_adapter_    Adapter
	_pkey_       string // 0 The name of the primary key in this table
	_conds_      []string
	locationId   int64
	locationCode string
	taxRateId    int64
	locationType string
}

func NewWooTaxRateLocation(a Adapter) *WooTaxRateLocation {
	var o WooTaxRateLocation
	o._table_ = "wp_woocommerce_tax_rate_locations"
	o._adapter_ = a
	o._pkey_ = "location_id"
	return &o
}

func (o *WooTaxRateLocation) Find(_find_by_locationId int64) (WooTaxRateLocation, error) {

	var model_slice []WooTaxRateLocation
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "location_id", _find_by_locationId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "location_code", _find_by_locationCode)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "tax_rate_id", _find_by_taxRateId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "location_type", _find_by_locationType)
	results, err := o._adapter_.Query(q)
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

type WooTaxRate struct {
	_table_         string
	_adapter_       Adapter
	_pkey_          string // 0 The name of the primary key in this table
	_conds_         []string
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
	o._table_ = "wp_woocommerce_tax_rates"
	o._adapter_ = a
	o._pkey_ = "tax_rate_id"
	return &o
}

func (o *WooTaxRate) Find(_find_by_taxRateId int64) (WooTaxRate, error) {

	var model_slice []WooTaxRate
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "tax_rate_id", _find_by_taxRateId)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "tax_rate_country", _find_by_taxRateCountry)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "tax_rate_state", _find_by_taxRateState)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "tax_rate", _find_by_taxRate)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "tax_rate_name", _find_by_taxRateName)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "tax_rate_priority", _find_by_taxRatePriority)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "tax_rate_compound", _find_by_taxRateCompound)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "tax_rate_shipping", _find_by_taxRateShipping)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table_, "tax_rate_order", _find_by_taxRateOrder)
	results, err := o._adapter_.Query(q)
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
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table_, "tax_rate_class", _find_by_taxRateClass)
	results, err := o._adapter_.Query(q)
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

type Adapter interface {
	Open(string, string, string, string) error
	Close()
	Query(string) ([]map[string]DBValue, error)
	Execute(string) error
	LastInsertedId() int64
	AffectedRows() int64
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

type MysqlAdapter struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Pass     string `yaml: "pass"`
	Database string `yaml:"database"`
	_conn_   *sql.DB
	_lid     int64
	_cnt     int64
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
