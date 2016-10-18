package gopress

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type LogFilter func(string, string) string
type Adapter interface {
	Open(string, string, string, string) error
	Close()
	Query(string) ([]map[string]DBValue, error)
	Execute(string) error
	LastInsertedId() int64
	AffectedRows() int64
	DatabasePrefix() string
	LogInfo(string)
	LogError(error)
	LogDebug(string)
	SetLogs(io.Writer)
	SetLogFilter(LogFilter)
	Oops(string) error
	SafeString(string) string
	NewDBValue() DBValue
}

type MysqlAdapter struct {
	Host        string `yaml:"host"`
	User        string `yaml:"user"`
	Pass        string `yaml: "pass"`
	Database    string `yaml:"database"`
	DBPrefix    string `yaml:"prefix"`
	_info_log   *log.Logger
	_error_log  *log.Logger
	_debug_log  *log.Logger
	_conn_      *sql.DB
	_lid        int64
	_cnt        int64
	_opened     bool
	_log_filter LogFilter
}

func NewMysqlAdapter(pre string) *MysqlAdapter {
	return &MysqlAdapter{DBPrefix: pre}
}
func NewMysqlAdapterEx(fname string) (*MysqlAdapter, error) {
	a := NewMysqlAdapter(``)
	y, err := fileGetContents(fname)
	if err != nil {
		return nil, err
	}
	err = a.FromYAML(y)
	if err != nil {
		return nil, err
	}
	err = a.Open(a.Host, a.User, a.Pass, a.Database)
	if err != nil {
		return nil, err
	}
	a.SetLogs(ioutil.Discard)
	return a, nil
}
func (a *MysqlAdapter) SetLogFilter(f LogFilter) {
	a._log_filter = f
}
func (a *MysqlAdapter) SafeString(s string) string {
	return s
}
func (a *MysqlAdapter) SetInfoLog(t io.Writer) {
	a._info_log = log.New(t, `[INFO]:`, log.Ldate|log.Ltime|log.Lshortfile)
}
func (a *MysqlAdapter) SetErrorLog(t io.Writer) {
	a._error_log = log.New(t, `[ERROR]:`, log.Ldate|log.Ltime|log.Lshortfile)
}
func (a *MysqlAdapter) SetDebugLog(t io.Writer) {
	a._debug_log = log.New(t, `[DEBUG]:`, log.Ldate|log.Ltime|log.Lshortfile)
}
func (a *MysqlAdapter) SetLogs(t io.Writer) {
	a.SetInfoLog(t)
	a.SetErrorLog(t)
	a.SetDebugLog(t)
}

func (a *MysqlAdapter) LogInfo(s string) {
	if a._log_filter != nil {
		s = a._log_filter(`INFO`, s)
	}
	if s == "" {
		return
	}
	a._info_log.Println(s)
}

func (a *MysqlAdapter) LogError(s error) {
	if a._log_filter != nil {
		ns := a._log_filter(`ERROR`, fmt.Sprintf(`%s`, s))
		a._error_log.Println(ns)
		return
	}
	a._error_log.Println(s)
}

func (a *MysqlAdapter) LogDebug(s string) {
	if a._log_filter != nil {
		s = a._log_filter(`DEBUG`, s)
	}
	if s == "" {
		return
	}
	a._debug_log.Println(s)
}

func (a *MysqlAdapter) NewDBValue() DBValue {
	return NewMysqlValue(a)
}
func (a *MysqlAdapter) DatabasePrefix() string {
	return a.DBPrefix
}
func (a *MysqlAdapter) FromYAML(b []byte) error {
	return yaml.Unmarshal(b, a)
}

func (a *MysqlAdapter) Open(h, u, p, d string) error {
	if h != "localhost" {
		l := fmt.Sprintf("%s:%s@tcp(%s)/%s", u, p, h, d)
		tc, err := sql.Open("mysql", l)
		if err != nil {
			return a.Oops(fmt.Sprintf(`%s with %s`, err, l))
		}
		a._conn_ = tc
	} else {
		l := fmt.Sprintf("%s:%s@/%s", u, p, d)
		tc, err := sql.Open("mysql", l)
		if err != nil {
			return a.Oops(fmt.Sprintf(`%s with %s`, err, l))
		}
		a._conn_ = tc
	}
	a._opened = true
	return nil

}
func (a *MysqlAdapter) Close() {
	a._conn_.Close()
}

func (a *MysqlAdapter) Query(q string) ([]map[string]DBValue, error) {
	if a._opened != true {
		return nil, a.Oops(`you must first open the connection`)
	}
	results := new([]map[string]DBValue)
	a.LogInfo(q)
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
func (a *MysqlAdapter) Oops(s string) error {
	e := errors.New(s)
	a.LogError(e)
	return e
}
func (a *MysqlAdapter) Execute(q string) error {
	if a._opened != true {
		return a.Oops(`you must first open the connection`)
	}
	tx, err := a._conn_.Begin()
	if err != nil {
		return a.Oops(fmt.Sprintf(`could not Begin Transaction %s`, err))
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(q)
	if err != nil {
		return a.Oops(fmt.Sprintf(`could not Prepare Statement %s`, err))
	}
	defer stmt.Close()
	a.LogInfo(q)
	res, err := stmt.Exec()
	if err != nil {
		return a.Oops(fmt.Sprintf(`could not Exec stmt %s`, err))
	}
	a._lid, err = res.LastInsertId()
	a.LogInfo(fmt.Sprintf(`LastInsertedId is %d`, a._lid))
	if err != nil {
		return a.Oops(fmt.Sprintf(`could not get LastInsertId %s`, err))
	}
	a._cnt, err = res.RowsAffected()
	if err != nil {
		return a.Oops(fmt.Sprintf(`could not get RowsAffected %s`, err))
	}
	err = tx.Commit()
	if err != nil {
		return a.Oops(fmt.Sprintf(`could not Commit Transaction %s`, err))
	}
	return nil
}
func (a *MysqlAdapter) LastInsertedId() int64 {
	return a._lid
}
func (a *MysqlAdapter) AffectedRows() int64 {
	return a._cnt
}

type DBValue interface {
	AsInt() (int, error)
	AsInt32() (int32, error)
	AsInt64() (int64, error)
	AsFloat32() (float32, error)
	AsFloat64() (float64, error)
	AsString() (string, error)
	AsDateTime() (*DateTime, error)
	SetInternalValue(string, string)
}

type MysqlValue struct {
	_v       string
	_k       string
	_adapter Adapter
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

func (v *MysqlValue) AsDateTime() (*DateTime, error) {
	dt := NewDateTime(v._adapter)
	err := dt.FromString(v._v)
	if err != nil {
		return &DateTime{}, err
	}
	return dt, nil
}

func NewMysqlValue(a Adapter) *MysqlValue {
	return &MysqlValue{_adapter: a}
}

type DateTime struct {
	Day      int
	Month    int
	Year     int
	Hours    int
	Minutes  int
	Seconds  int
	_adapter Adapter
}

func (d *DateTime) FromString(s string) error {
	es := s
	re := regexp.MustCompile("(?P<year>[\\d]{4})-(?P<month>[\\d]{2})-(?P<day>[\\d]{2}) (?P<hours>[\\d]{2}):(?P<minutes>[\\d]{2}):(?P<seconds>[\\d]{2})")
	n1 := re.SubexpNames()
	ir2 := re.FindAllStringSubmatch(es, -1)
	if len(ir2) == 0 {
		return d._adapter.Oops(fmt.Sprintf("found no data to capture in %s", es))
	}
	r2 := ir2[0]
	for i, n := range r2 {
		if n1[i] == "year" {
			_Year, err := strconv.ParseInt(n, 10, 32)
			d.Year = int(_Year)
			if err != nil {
				return d._adapter.Oops(fmt.Sprintf("failed to convert %s in %s received %s", n[i], es, err))
			}
		}
		if n1[i] == "month" {
			_Month, err := strconv.ParseInt(n, 10, 32)
			d.Month = int(_Month)
			if err != nil {
				return d._adapter.Oops(fmt.Sprintf("failed to convert %s in %s received %s", n[i], es, err))
			}
		}
		if n1[i] == "day" {
			_Day, err := strconv.ParseInt(n, 10, 32)
			d.Day = int(_Day)
			if err != nil {
				return d._adapter.Oops(fmt.Sprintf("failed to convert %s in %s received %s", n[i], es, err))
			}
		}
		if n1[i] == "hours" {
			_Hours, err := strconv.ParseInt(n, 10, 32)
			d.Hours = int(_Hours)
			if err != nil {
				return d._adapter.Oops(fmt.Sprintf("failed to convert %s in %s received %s", n[i], es, err))
			}
		}
		if n1[i] == "minutes" {
			_Minutes, err := strconv.ParseInt(n, 10, 32)
			d.Minutes = int(_Minutes)
			if err != nil {
				return d._adapter.Oops(fmt.Sprintf("failed to convert %s in %s received %s", n[i], es, err))
			}
		}
		if n1[i] == "seconds" {
			_Seconds, err := strconv.ParseInt(n, 10, 32)
			d.Seconds = int(_Seconds)
			if err != nil {
				return d._adapter.Oops(fmt.Sprintf("failed to convert %s in %s received %s", n[i], es, err))
			}
		}
	}
	return nil
}
func (d *DateTime) ToString() string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", d.Year, d.Month, d.Day, d.Hours, d.Minutes, d.Seconds)
}
func (d *DateTime) String() string {
	return d.ToString()
}
func NewDateTime(a Adapter) *DateTime {
	d := &DateTime{_adapter: a}
	return d
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
	// Dirty markers for smart updates
	IsMetaIdDirty    bool
	IsCommentIdDirty bool
	IsMetaKeyDirty   bool
	IsMetaValueDirty bool
	// Relationships
}

func NewCommentMeta(a Adapter) *CommentMeta {
	var o CommentMeta
	o._table = fmt.Sprintf("%scommentmeta", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "meta_id"
	o._new = false
	return &o
}

func (m *CommentMeta) GetPrimaryKeyValue() int64 {
	return m.MetaId
}
func (m *CommentMeta) GetPrimaryKeyName() string {
	return `meta_id`
}

func (m *CommentMeta) GetMetaId() int64 {
	return m.MetaId
}
func (m *CommentMeta) SetMetaId(arg int64) {
	m.MetaId = arg
	m.IsMetaIdDirty = true
}

func (m *CommentMeta) GetCommentId() int64 {
	return m.CommentId
}
func (m *CommentMeta) SetCommentId(arg int64) {
	m.CommentId = arg
	m.IsCommentIdDirty = true
}

func (m *CommentMeta) GetMetaKey() string {
	return m.MetaKey
}
func (m *CommentMeta) SetMetaKey(arg string) {
	m.MetaKey = arg
	m.IsMetaKeyDirty = true
}

func (m *CommentMeta) GetMetaValue() string {
	return m.MetaValue
}
func (m *CommentMeta) SetMetaValue(arg string) {
	m.MetaValue = arg
	m.IsMetaValueDirty = true
}

func (o *CommentMeta) Find(_find_by_MetaId int64) (bool, error) {

	var model_slice []*CommentMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "meta_id", _find_by_MetaId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return false, o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}

	for _, result := range results {
		ro := CommentMeta{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return false, o._adapter.Oops(fmt.Sprintf(`%s`, err))
		}
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return false, o._adapter.Oops(`not found`)
	}
	o.FromCommentMeta(model_slice[0])
	return true, nil

}
func (o *CommentMeta) FindByCommentId(_find_by_CommentId int64) ([]*CommentMeta, error) {

	var model_slice []*CommentMeta
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *CommentMeta) FindByMetaKey(_find_by_MetaKey string) ([]*CommentMeta, error) {

	var model_slice []*CommentMeta
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *CommentMeta) FindByMetaValue(_find_by_MetaValue string) ([]*CommentMeta, error) {

	var model_slice []*CommentMeta
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}

func (o *CommentMeta) FromDBValueMap(m map[string]DBValue) error {
	_MetaId, err := m["meta_id"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.MetaId = _MetaId
	_CommentId, err := m["comment_id"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.CommentId = _CommentId
	_MetaKey, err := m["meta_key"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.MetaKey = _MetaKey
	_MetaValue, err := m["meta_value"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.MetaValue = _MetaValue

	return nil
}
func (o *CommentMeta) FromCommentMeta(m *CommentMeta) {
	o.MetaId = m.MetaId
	o.CommentId = m.CommentId
	o.MetaKey = m.MetaKey
	o.MetaValue = m.MetaValue

}

func (o *CommentMeta) Save() error {
	if o._new == true {
		return o.Create()
	}
	var sets []string

	if o.IsCommentIdDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_id = '%d'`, o.CommentId))
	}

	if o.IsMetaKeyDirty == true {
		sets = append(sets, fmt.Sprintf(`meta_key = '%s'`, o._adapter.SafeString(o.MetaKey)))
	}

	if o.IsMetaValueDirty == true {
		sets = append(sets, fmt.Sprintf(`meta_value = '%s'`, o._adapter.SafeString(o.MetaValue)))
	}

	frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'", o._table, strings.Join(sets, `,`))
	err := o._adapter.Execute(frmt)
	if err != nil {
		return err
	}
	return nil
}
func (o *CommentMeta) Update() error {
	var sets []string

	if o.IsCommentIdDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_id = '%d'`, o.CommentId))
	}

	if o.IsMetaKeyDirty == true {
		sets = append(sets, fmt.Sprintf(`meta_key = '%s'`, o._adapter.SafeString(o.MetaKey)))
	}

	if o.IsMetaValueDirty == true {
		sets = append(sets, fmt.Sprintf(`meta_value = '%s'`, o._adapter.SafeString(o.MetaValue)))
	}

	frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'", o._table, strings.Join(sets, `,`))
	err := o._adapter.Execute(frmt)
	if err != nil {
		return err
	}
	return nil
}
func (o *CommentMeta) Create() error {
	frmt := fmt.Sprintf("INSERT INTO %s (`comment_id`, `meta_key`, `meta_value`) VALUES ('%d', '%s', '%s')", o._table, o.CommentId, o.MetaKey, o.MetaValue)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s led to %s`, frmt, err))
	}
	o.MetaId = o._adapter.LastInsertedId()
	o._new = false
	return nil
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
	CommentDate        *DateTime
	CommentDateGmt     *DateTime
	CommentContent     string
	CommentKarma       int
	CommentApproved    string
	CommentAgent       string
	CommentType        string
	CommentParent      int64
	UserId             int64
	// Dirty markers for smart updates
	IsCommentIDDirty          bool
	IsCommentPostIDDirty      bool
	IsCommentAuthorDirty      bool
	IsCommentAuthorEmailDirty bool
	IsCommentAuthorUrlDirty   bool
	IsCommentAuthorIPDirty    bool
	IsCommentDateDirty        bool
	IsCommentDateGmtDirty     bool
	IsCommentContentDirty     bool
	IsCommentKarmaDirty       bool
	IsCommentApprovedDirty    bool
	IsCommentAgentDirty       bool
	IsCommentTypeDirty        bool
	IsCommentParentDirty      bool
	IsUserIdDirty             bool
	// Relationships
}

func NewComment(a Adapter) *Comment {
	var o Comment
	o._table = fmt.Sprintf("%scomments", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "comment_ID"
	o._new = false
	return &o
}

func (m *Comment) GetPrimaryKeyValue() int64 {
	return m.CommentID
}
func (m *Comment) GetPrimaryKeyName() string {
	return `comment_ID`
}

func (m *Comment) GetCommentID() int64 {
	return m.CommentID
}
func (m *Comment) SetCommentID(arg int64) {
	m.CommentID = arg
	m.IsCommentIDDirty = true
}

func (m *Comment) GetCommentPostID() int64 {
	return m.CommentPostID
}
func (m *Comment) SetCommentPostID(arg int64) {
	m.CommentPostID = arg
	m.IsCommentPostIDDirty = true
}

func (m *Comment) GetCommentAuthor() string {
	return m.CommentAuthor
}
func (m *Comment) SetCommentAuthor(arg string) {
	m.CommentAuthor = arg
	m.IsCommentAuthorDirty = true
}

func (m *Comment) GetCommentAuthorEmail() string {
	return m.CommentAuthorEmail
}
func (m *Comment) SetCommentAuthorEmail(arg string) {
	m.CommentAuthorEmail = arg
	m.IsCommentAuthorEmailDirty = true
}

func (m *Comment) GetCommentAuthorUrl() string {
	return m.CommentAuthorUrl
}
func (m *Comment) SetCommentAuthorUrl(arg string) {
	m.CommentAuthorUrl = arg
	m.IsCommentAuthorUrlDirty = true
}

func (m *Comment) GetCommentAuthorIP() string {
	return m.CommentAuthorIP
}
func (m *Comment) SetCommentAuthorIP(arg string) {
	m.CommentAuthorIP = arg
	m.IsCommentAuthorIPDirty = true
}

func (m *Comment) GetCommentDate() *DateTime {
	return m.CommentDate
}
func (m *Comment) SetCommentDate(arg *DateTime) {
	m.CommentDate = arg
	m.IsCommentDateDirty = true
}

func (m *Comment) GetCommentDateGmt() *DateTime {
	return m.CommentDateGmt
}
func (m *Comment) SetCommentDateGmt(arg *DateTime) {
	m.CommentDateGmt = arg
	m.IsCommentDateGmtDirty = true
}

func (m *Comment) GetCommentContent() string {
	return m.CommentContent
}
func (m *Comment) SetCommentContent(arg string) {
	m.CommentContent = arg
	m.IsCommentContentDirty = true
}

func (m *Comment) GetCommentKarma() int {
	return m.CommentKarma
}
func (m *Comment) SetCommentKarma(arg int) {
	m.CommentKarma = arg
	m.IsCommentKarmaDirty = true
}

func (m *Comment) GetCommentApproved() string {
	return m.CommentApproved
}
func (m *Comment) SetCommentApproved(arg string) {
	m.CommentApproved = arg
	m.IsCommentApprovedDirty = true
}

func (m *Comment) GetCommentAgent() string {
	return m.CommentAgent
}
func (m *Comment) SetCommentAgent(arg string) {
	m.CommentAgent = arg
	m.IsCommentAgentDirty = true
}

func (m *Comment) GetCommentType() string {
	return m.CommentType
}
func (m *Comment) SetCommentType(arg string) {
	m.CommentType = arg
	m.IsCommentTypeDirty = true
}

func (m *Comment) GetCommentParent() int64 {
	return m.CommentParent
}
func (m *Comment) SetCommentParent(arg int64) {
	m.CommentParent = arg
	m.IsCommentParentDirty = true
}

func (m *Comment) GetUserId() int64 {
	return m.UserId
}
func (m *Comment) SetUserId(arg int64) {
	m.UserId = arg
	m.IsUserIdDirty = true
}

func (o *Comment) Find(_find_by_CommentID int64) (bool, error) {

	var model_slice []*Comment
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "comment_ID", _find_by_CommentID)
	results, err := o._adapter.Query(q)
	if err != nil {
		return false, o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}

	for _, result := range results {
		ro := Comment{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return false, o._adapter.Oops(fmt.Sprintf(`%s`, err))
		}
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return false, o._adapter.Oops(`not found`)
	}
	o.FromComment(model_slice[0])
	return true, nil

}
func (o *Comment) FindByCommentPostID(_find_by_CommentPostID int64) ([]*Comment, error) {

	var model_slice []*Comment
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentAuthor(_find_by_CommentAuthor string) ([]*Comment, error) {

	var model_slice []*Comment
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentAuthorEmail(_find_by_CommentAuthorEmail string) ([]*Comment, error) {

	var model_slice []*Comment
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentAuthorUrl(_find_by_CommentAuthorUrl string) ([]*Comment, error) {

	var model_slice []*Comment
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentAuthorIP(_find_by_CommentAuthorIP string) ([]*Comment, error) {

	var model_slice []*Comment
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentDate(_find_by_CommentDate *DateTime) ([]*Comment, error) {

	var model_slice []*Comment
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentDateGmt(_find_by_CommentDateGmt *DateTime) ([]*Comment, error) {

	var model_slice []*Comment
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentContent(_find_by_CommentContent string) ([]*Comment, error) {

	var model_slice []*Comment
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentKarma(_find_by_CommentKarma int) ([]*Comment, error) {

	var model_slice []*Comment
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentApproved(_find_by_CommentApproved string) ([]*Comment, error) {

	var model_slice []*Comment
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentAgent(_find_by_CommentAgent string) ([]*Comment, error) {

	var model_slice []*Comment
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentType(_find_by_CommentType string) ([]*Comment, error) {

	var model_slice []*Comment
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Comment) FindByCommentParent(_find_by_CommentParent int64) ([]*Comment, error) {

	var model_slice []*Comment
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Comment) FindByUserId(_find_by_UserId int64) ([]*Comment, error) {

	var model_slice []*Comment
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}

func (o *Comment) FromDBValueMap(m map[string]DBValue) error {
	_CommentID, err := m["comment_ID"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.CommentID = _CommentID
	_CommentPostID, err := m["comment_post_ID"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.CommentPostID = _CommentPostID
	_CommentAuthor, err := m["comment_author"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.CommentAuthor = _CommentAuthor
	_CommentAuthorEmail, err := m["comment_author_email"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.CommentAuthorEmail = _CommentAuthorEmail
	_CommentAuthorUrl, err := m["comment_author_url"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.CommentAuthorUrl = _CommentAuthorUrl
	_CommentAuthorIP, err := m["comment_author_IP"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.CommentAuthorIP = _CommentAuthorIP
	_CommentDate, err := m["comment_date"].AsDateTime()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.CommentDate = _CommentDate
	_CommentDateGmt, err := m["comment_date_gmt"].AsDateTime()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.CommentDateGmt = _CommentDateGmt
	_CommentContent, err := m["comment_content"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.CommentContent = _CommentContent
	_CommentKarma, err := m["comment_karma"].AsInt()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.CommentKarma = _CommentKarma
	_CommentApproved, err := m["comment_approved"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.CommentApproved = _CommentApproved
	_CommentAgent, err := m["comment_agent"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.CommentAgent = _CommentAgent
	_CommentType, err := m["comment_type"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.CommentType = _CommentType
	_CommentParent, err := m["comment_parent"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.CommentParent = _CommentParent
	_UserId, err := m["user_id"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.UserId = _UserId

	return nil
}
func (o *Comment) FromComment(m *Comment) {
	o.CommentID = m.CommentID
	o.CommentPostID = m.CommentPostID
	o.CommentAuthor = m.CommentAuthor
	o.CommentAuthorEmail = m.CommentAuthorEmail
	o.CommentAuthorUrl = m.CommentAuthorUrl
	o.CommentAuthorIP = m.CommentAuthorIP
	o.CommentDate = m.CommentDate
	o.CommentDateGmt = m.CommentDateGmt
	o.CommentContent = m.CommentContent
	o.CommentKarma = m.CommentKarma
	o.CommentApproved = m.CommentApproved
	o.CommentAgent = m.CommentAgent
	o.CommentType = m.CommentType
	o.CommentParent = m.CommentParent
	o.UserId = m.UserId

}

func (o *Comment) Save() error {
	if o._new == true {
		return o.Create()
	}
	var sets []string

	if o.IsCommentPostIDDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_post_ID = '%d'`, o.CommentPostID))
	}

	if o.IsCommentAuthorDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_author = '%s'`, o._adapter.SafeString(o.CommentAuthor)))
	}

	if o.IsCommentAuthorEmailDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_author_email = '%s'`, o._adapter.SafeString(o.CommentAuthorEmail)))
	}

	if o.IsCommentAuthorUrlDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_author_url = '%s'`, o._adapter.SafeString(o.CommentAuthorUrl)))
	}

	if o.IsCommentAuthorIPDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_author_IP = '%s'`, o._adapter.SafeString(o.CommentAuthorIP)))
	}

	if o.IsCommentDateDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_date = '%s'`, o.CommentDate))
	}

	if o.IsCommentDateGmtDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_date_gmt = '%s'`, o.CommentDateGmt))
	}

	if o.IsCommentContentDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_content = '%s'`, o._adapter.SafeString(o.CommentContent)))
	}

	if o.IsCommentKarmaDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_karma = '%d'`, o.CommentKarma))
	}

	if o.IsCommentApprovedDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_approved = '%s'`, o._adapter.SafeString(o.CommentApproved)))
	}

	if o.IsCommentAgentDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_agent = '%s'`, o._adapter.SafeString(o.CommentAgent)))
	}

	if o.IsCommentTypeDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_type = '%s'`, o._adapter.SafeString(o.CommentType)))
	}

	if o.IsCommentParentDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_parent = '%d'`, o.CommentParent))
	}

	if o.IsUserIdDirty == true {
		sets = append(sets, fmt.Sprintf(`user_id = '%d'`, o.UserId))
	}

	frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'", o._table, strings.Join(sets, `,`))
	err := o._adapter.Execute(frmt)
	if err != nil {
		return err
	}
	return nil
}
func (o *Comment) Update() error {
	var sets []string

	if o.IsCommentPostIDDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_post_ID = '%d'`, o.CommentPostID))
	}

	if o.IsCommentAuthorDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_author = '%s'`, o._adapter.SafeString(o.CommentAuthor)))
	}

	if o.IsCommentAuthorEmailDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_author_email = '%s'`, o._adapter.SafeString(o.CommentAuthorEmail)))
	}

	if o.IsCommentAuthorUrlDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_author_url = '%s'`, o._adapter.SafeString(o.CommentAuthorUrl)))
	}

	if o.IsCommentAuthorIPDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_author_IP = '%s'`, o._adapter.SafeString(o.CommentAuthorIP)))
	}

	if o.IsCommentDateDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_date = '%s'`, o.CommentDate))
	}

	if o.IsCommentDateGmtDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_date_gmt = '%s'`, o.CommentDateGmt))
	}

	if o.IsCommentContentDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_content = '%s'`, o._adapter.SafeString(o.CommentContent)))
	}

	if o.IsCommentKarmaDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_karma = '%d'`, o.CommentKarma))
	}

	if o.IsCommentApprovedDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_approved = '%s'`, o._adapter.SafeString(o.CommentApproved)))
	}

	if o.IsCommentAgentDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_agent = '%s'`, o._adapter.SafeString(o.CommentAgent)))
	}

	if o.IsCommentTypeDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_type = '%s'`, o._adapter.SafeString(o.CommentType)))
	}

	if o.IsCommentParentDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_parent = '%d'`, o.CommentParent))
	}

	if o.IsUserIdDirty == true {
		sets = append(sets, fmt.Sprintf(`user_id = '%d'`, o.UserId))
	}

	frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'", o._table, strings.Join(sets, `,`))
	err := o._adapter.Execute(frmt)
	if err != nil {
		return err
	}
	return nil
}
func (o *Comment) Create() error {
	frmt := fmt.Sprintf("INSERT INTO %s (`comment_post_ID`, `comment_author`, `comment_author_email`, `comment_author_url`, `comment_author_IP`, `comment_date`, `comment_date_gmt`, `comment_content`, `comment_karma`, `comment_approved`, `comment_agent`, `comment_type`, `comment_parent`, `user_id`) VALUES ('%d', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%d', '%s', '%s', '%s', '%d', '%d')", o._table, o.CommentPostID, o.CommentAuthor, o.CommentAuthorEmail, o.CommentAuthorUrl, o.CommentAuthorIP, o.CommentDate.ToString(), o.CommentDateGmt.ToString(), o.CommentContent, o.CommentKarma, o.CommentApproved, o.CommentAgent, o.CommentType, o.CommentParent, o.UserId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s led to %s`, frmt, err))
	}
	o.CommentID = o._adapter.LastInsertedId()
	o._new = false
	return nil
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

func (o *Comment) UpdateCommentDate(_upd_CommentDate *DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `comment_date` = '%s' WHERE `comment_ID` = '%d'", o._table, _upd_CommentDate, o.CommentDate)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.CommentDate = _upd_CommentDate
	return o._adapter.AffectedRows(), nil
}

func (o *Comment) UpdateCommentDateGmt(_upd_CommentDateGmt *DateTime) (int64, error) {
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
	LinkUpdated     *DateTime
	LinkRel         string
	LinkNotes       string
	LinkRss         string
	// Dirty markers for smart updates
	IsLinkIdDirty          bool
	IsLinkUrlDirty         bool
	IsLinkNameDirty        bool
	IsLinkImageDirty       bool
	IsLinkTargetDirty      bool
	IsLinkDescriptionDirty bool
	IsLinkVisibleDirty     bool
	IsLinkOwnerDirty       bool
	IsLinkRatingDirty      bool
	IsLinkUpdatedDirty     bool
	IsLinkRelDirty         bool
	IsLinkNotesDirty       bool
	IsLinkRssDirty         bool
	// Relationships
}

func NewLink(a Adapter) *Link {
	var o Link
	o._table = fmt.Sprintf("%slinks", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "link_id"
	o._new = false
	return &o
}

func (m *Link) GetPrimaryKeyValue() int64 {
	return m.LinkId
}
func (m *Link) GetPrimaryKeyName() string {
	return `link_id`
}

func (m *Link) GetLinkId() int64 {
	return m.LinkId
}
func (m *Link) SetLinkId(arg int64) {
	m.LinkId = arg
	m.IsLinkIdDirty = true
}

func (m *Link) GetLinkUrl() string {
	return m.LinkUrl
}
func (m *Link) SetLinkUrl(arg string) {
	m.LinkUrl = arg
	m.IsLinkUrlDirty = true
}

func (m *Link) GetLinkName() string {
	return m.LinkName
}
func (m *Link) SetLinkName(arg string) {
	m.LinkName = arg
	m.IsLinkNameDirty = true
}

func (m *Link) GetLinkImage() string {
	return m.LinkImage
}
func (m *Link) SetLinkImage(arg string) {
	m.LinkImage = arg
	m.IsLinkImageDirty = true
}

func (m *Link) GetLinkTarget() string {
	return m.LinkTarget
}
func (m *Link) SetLinkTarget(arg string) {
	m.LinkTarget = arg
	m.IsLinkTargetDirty = true
}

func (m *Link) GetLinkDescription() string {
	return m.LinkDescription
}
func (m *Link) SetLinkDescription(arg string) {
	m.LinkDescription = arg
	m.IsLinkDescriptionDirty = true
}

func (m *Link) GetLinkVisible() string {
	return m.LinkVisible
}
func (m *Link) SetLinkVisible(arg string) {
	m.LinkVisible = arg
	m.IsLinkVisibleDirty = true
}

func (m *Link) GetLinkOwner() int64 {
	return m.LinkOwner
}
func (m *Link) SetLinkOwner(arg int64) {
	m.LinkOwner = arg
	m.IsLinkOwnerDirty = true
}

func (m *Link) GetLinkRating() int {
	return m.LinkRating
}
func (m *Link) SetLinkRating(arg int) {
	m.LinkRating = arg
	m.IsLinkRatingDirty = true
}

func (m *Link) GetLinkUpdated() *DateTime {
	return m.LinkUpdated
}
func (m *Link) SetLinkUpdated(arg *DateTime) {
	m.LinkUpdated = arg
	m.IsLinkUpdatedDirty = true
}

func (m *Link) GetLinkRel() string {
	return m.LinkRel
}
func (m *Link) SetLinkRel(arg string) {
	m.LinkRel = arg
	m.IsLinkRelDirty = true
}

func (m *Link) GetLinkNotes() string {
	return m.LinkNotes
}
func (m *Link) SetLinkNotes(arg string) {
	m.LinkNotes = arg
	m.IsLinkNotesDirty = true
}

func (m *Link) GetLinkRss() string {
	return m.LinkRss
}
func (m *Link) SetLinkRss(arg string) {
	m.LinkRss = arg
	m.IsLinkRssDirty = true
}

func (o *Link) Find(_find_by_LinkId int64) (bool, error) {

	var model_slice []*Link
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "link_id", _find_by_LinkId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return false, o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}

	for _, result := range results {
		ro := Link{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return false, o._adapter.Oops(fmt.Sprintf(`%s`, err))
		}
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return false, o._adapter.Oops(`not found`)
	}
	o.FromLink(model_slice[0])
	return true, nil

}
func (o *Link) FindByLinkUrl(_find_by_LinkUrl string) ([]*Link, error) {

	var model_slice []*Link
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Link) FindByLinkName(_find_by_LinkName string) ([]*Link, error) {

	var model_slice []*Link
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Link) FindByLinkImage(_find_by_LinkImage string) ([]*Link, error) {

	var model_slice []*Link
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Link) FindByLinkTarget(_find_by_LinkTarget string) ([]*Link, error) {

	var model_slice []*Link
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Link) FindByLinkDescription(_find_by_LinkDescription string) ([]*Link, error) {

	var model_slice []*Link
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Link) FindByLinkVisible(_find_by_LinkVisible string) ([]*Link, error) {

	var model_slice []*Link
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Link) FindByLinkOwner(_find_by_LinkOwner int64) ([]*Link, error) {

	var model_slice []*Link
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Link) FindByLinkRating(_find_by_LinkRating int) ([]*Link, error) {

	var model_slice []*Link
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Link) FindByLinkUpdated(_find_by_LinkUpdated *DateTime) ([]*Link, error) {

	var model_slice []*Link
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Link) FindByLinkRel(_find_by_LinkRel string) ([]*Link, error) {

	var model_slice []*Link
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Link) FindByLinkNotes(_find_by_LinkNotes string) ([]*Link, error) {

	var model_slice []*Link
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Link) FindByLinkRss(_find_by_LinkRss string) ([]*Link, error) {

	var model_slice []*Link
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}

func (o *Link) FromDBValueMap(m map[string]DBValue) error {
	_LinkId, err := m["link_id"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.LinkId = _LinkId
	_LinkUrl, err := m["link_url"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.LinkUrl = _LinkUrl
	_LinkName, err := m["link_name"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.LinkName = _LinkName
	_LinkImage, err := m["link_image"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.LinkImage = _LinkImage
	_LinkTarget, err := m["link_target"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.LinkTarget = _LinkTarget
	_LinkDescription, err := m["link_description"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.LinkDescription = _LinkDescription
	_LinkVisible, err := m["link_visible"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.LinkVisible = _LinkVisible
	_LinkOwner, err := m["link_owner"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.LinkOwner = _LinkOwner
	_LinkRating, err := m["link_rating"].AsInt()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.LinkRating = _LinkRating
	_LinkUpdated, err := m["link_updated"].AsDateTime()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.LinkUpdated = _LinkUpdated
	_LinkRel, err := m["link_rel"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.LinkRel = _LinkRel
	_LinkNotes, err := m["link_notes"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.LinkNotes = _LinkNotes
	_LinkRss, err := m["link_rss"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.LinkRss = _LinkRss

	return nil
}
func (o *Link) FromLink(m *Link) {
	o.LinkId = m.LinkId
	o.LinkUrl = m.LinkUrl
	o.LinkName = m.LinkName
	o.LinkImage = m.LinkImage
	o.LinkTarget = m.LinkTarget
	o.LinkDescription = m.LinkDescription
	o.LinkVisible = m.LinkVisible
	o.LinkOwner = m.LinkOwner
	o.LinkRating = m.LinkRating
	o.LinkUpdated = m.LinkUpdated
	o.LinkRel = m.LinkRel
	o.LinkNotes = m.LinkNotes
	o.LinkRss = m.LinkRss

}

func (o *Link) Save() error {
	if o._new == true {
		return o.Create()
	}
	var sets []string

	if o.IsLinkUrlDirty == true {
		sets = append(sets, fmt.Sprintf(`link_url = '%s'`, o._adapter.SafeString(o.LinkUrl)))
	}

	if o.IsLinkNameDirty == true {
		sets = append(sets, fmt.Sprintf(`link_name = '%s'`, o._adapter.SafeString(o.LinkName)))
	}

	if o.IsLinkImageDirty == true {
		sets = append(sets, fmt.Sprintf(`link_image = '%s'`, o._adapter.SafeString(o.LinkImage)))
	}

	if o.IsLinkTargetDirty == true {
		sets = append(sets, fmt.Sprintf(`link_target = '%s'`, o._adapter.SafeString(o.LinkTarget)))
	}

	if o.IsLinkDescriptionDirty == true {
		sets = append(sets, fmt.Sprintf(`link_description = '%s'`, o._adapter.SafeString(o.LinkDescription)))
	}

	if o.IsLinkVisibleDirty == true {
		sets = append(sets, fmt.Sprintf(`link_visible = '%s'`, o._adapter.SafeString(o.LinkVisible)))
	}

	if o.IsLinkOwnerDirty == true {
		sets = append(sets, fmt.Sprintf(`link_owner = '%d'`, o.LinkOwner))
	}

	if o.IsLinkRatingDirty == true {
		sets = append(sets, fmt.Sprintf(`link_rating = '%d'`, o.LinkRating))
	}

	if o.IsLinkUpdatedDirty == true {
		sets = append(sets, fmt.Sprintf(`link_updated = '%s'`, o.LinkUpdated))
	}

	if o.IsLinkRelDirty == true {
		sets = append(sets, fmt.Sprintf(`link_rel = '%s'`, o._adapter.SafeString(o.LinkRel)))
	}

	if o.IsLinkNotesDirty == true {
		sets = append(sets, fmt.Sprintf(`link_notes = '%s'`, o._adapter.SafeString(o.LinkNotes)))
	}

	if o.IsLinkRssDirty == true {
		sets = append(sets, fmt.Sprintf(`link_rss = '%s'`, o._adapter.SafeString(o.LinkRss)))
	}

	frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'", o._table, strings.Join(sets, `,`))
	err := o._adapter.Execute(frmt)
	if err != nil {
		return err
	}
	return nil
}
func (o *Link) Update() error {
	var sets []string

	if o.IsLinkUrlDirty == true {
		sets = append(sets, fmt.Sprintf(`link_url = '%s'`, o._adapter.SafeString(o.LinkUrl)))
	}

	if o.IsLinkNameDirty == true {
		sets = append(sets, fmt.Sprintf(`link_name = '%s'`, o._adapter.SafeString(o.LinkName)))
	}

	if o.IsLinkImageDirty == true {
		sets = append(sets, fmt.Sprintf(`link_image = '%s'`, o._adapter.SafeString(o.LinkImage)))
	}

	if o.IsLinkTargetDirty == true {
		sets = append(sets, fmt.Sprintf(`link_target = '%s'`, o._adapter.SafeString(o.LinkTarget)))
	}

	if o.IsLinkDescriptionDirty == true {
		sets = append(sets, fmt.Sprintf(`link_description = '%s'`, o._adapter.SafeString(o.LinkDescription)))
	}

	if o.IsLinkVisibleDirty == true {
		sets = append(sets, fmt.Sprintf(`link_visible = '%s'`, o._adapter.SafeString(o.LinkVisible)))
	}

	if o.IsLinkOwnerDirty == true {
		sets = append(sets, fmt.Sprintf(`link_owner = '%d'`, o.LinkOwner))
	}

	if o.IsLinkRatingDirty == true {
		sets = append(sets, fmt.Sprintf(`link_rating = '%d'`, o.LinkRating))
	}

	if o.IsLinkUpdatedDirty == true {
		sets = append(sets, fmt.Sprintf(`link_updated = '%s'`, o.LinkUpdated))
	}

	if o.IsLinkRelDirty == true {
		sets = append(sets, fmt.Sprintf(`link_rel = '%s'`, o._adapter.SafeString(o.LinkRel)))
	}

	if o.IsLinkNotesDirty == true {
		sets = append(sets, fmt.Sprintf(`link_notes = '%s'`, o._adapter.SafeString(o.LinkNotes)))
	}

	if o.IsLinkRssDirty == true {
		sets = append(sets, fmt.Sprintf(`link_rss = '%s'`, o._adapter.SafeString(o.LinkRss)))
	}

	frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'", o._table, strings.Join(sets, `,`))
	err := o._adapter.Execute(frmt)
	if err != nil {
		return err
	}
	return nil
}
func (o *Link) Create() error {
	frmt := fmt.Sprintf("INSERT INTO %s (`link_url`, `link_name`, `link_image`, `link_target`, `link_description`, `link_visible`, `link_owner`, `link_rating`, `link_updated`, `link_rel`, `link_notes`, `link_rss`) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%d', '%d', '%s', '%s', '%s', '%s')", o._table, o.LinkUrl, o.LinkName, o.LinkImage, o.LinkTarget, o.LinkDescription, o.LinkVisible, o.LinkOwner, o.LinkRating, o.LinkUpdated.ToString(), o.LinkRel, o.LinkNotes, o.LinkRss)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s led to %s`, frmt, err))
	}
	o.LinkId = o._adapter.LastInsertedId()
	o._new = false
	return nil
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

func (o *Link) UpdateLinkUpdated(_upd_LinkUpdated *DateTime) (int64, error) {
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
	// Dirty markers for smart updates
	IsOptionIdDirty    bool
	IsOptionNameDirty  bool
	IsOptionValueDirty bool
	IsAutoloadDirty    bool
	// Relationships
}

func NewOption(a Adapter) *Option {
	var o Option
	o._table = fmt.Sprintf("%soptions", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "option_id"
	o._new = false
	return &o
}

func (m *Option) GetPrimaryKeyValue() int64 {
	return m.OptionId
}
func (m *Option) GetPrimaryKeyName() string {
	return `option_id`
}

func (m *Option) GetOptionId() int64 {
	return m.OptionId
}
func (m *Option) SetOptionId(arg int64) {
	m.OptionId = arg
	m.IsOptionIdDirty = true
}

func (m *Option) GetOptionName() string {
	return m.OptionName
}
func (m *Option) SetOptionName(arg string) {
	m.OptionName = arg
	m.IsOptionNameDirty = true
}

func (m *Option) GetOptionValue() string {
	return m.OptionValue
}
func (m *Option) SetOptionValue(arg string) {
	m.OptionValue = arg
	m.IsOptionValueDirty = true
}

func (m *Option) GetAutoload() string {
	return m.Autoload
}
func (m *Option) SetAutoload(arg string) {
	m.Autoload = arg
	m.IsAutoloadDirty = true
}

func (o *Option) Find(_find_by_OptionId int64) (bool, error) {

	var model_slice []*Option
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "option_id", _find_by_OptionId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return false, o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}

	for _, result := range results {
		ro := Option{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return false, o._adapter.Oops(fmt.Sprintf(`%s`, err))
		}
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return false, o._adapter.Oops(`not found`)
	}
	o.FromOption(model_slice[0])
	return true, nil

}
func (o *Option) FindByOptionName(_find_by_OptionName string) ([]*Option, error) {

	var model_slice []*Option
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Option) FindByOptionValue(_find_by_OptionValue string) ([]*Option, error) {

	var model_slice []*Option
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Option) FindByAutoload(_find_by_Autoload string) ([]*Option, error) {

	var model_slice []*Option
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}

func (o *Option) FromDBValueMap(m map[string]DBValue) error {
	_OptionId, err := m["option_id"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.OptionId = _OptionId
	_OptionName, err := m["option_name"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.OptionName = _OptionName
	_OptionValue, err := m["option_value"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.OptionValue = _OptionValue
	_Autoload, err := m["autoload"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.Autoload = _Autoload

	return nil
}
func (o *Option) FromOption(m *Option) {
	o.OptionId = m.OptionId
	o.OptionName = m.OptionName
	o.OptionValue = m.OptionValue
	o.Autoload = m.Autoload

}

func (o *Option) Save() error {
	if o._new == true {
		return o.Create()
	}
	var sets []string

	if o.IsOptionNameDirty == true {
		sets = append(sets, fmt.Sprintf(`option_name = '%s'`, o._adapter.SafeString(o.OptionName)))
	}

	if o.IsOptionValueDirty == true {
		sets = append(sets, fmt.Sprintf(`option_value = '%s'`, o._adapter.SafeString(o.OptionValue)))
	}

	if o.IsAutoloadDirty == true {
		sets = append(sets, fmt.Sprintf(`autoload = '%s'`, o._adapter.SafeString(o.Autoload)))
	}

	frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'", o._table, strings.Join(sets, `,`))
	err := o._adapter.Execute(frmt)
	if err != nil {
		return err
	}
	return nil
}
func (o *Option) Update() error {
	var sets []string

	if o.IsOptionNameDirty == true {
		sets = append(sets, fmt.Sprintf(`option_name = '%s'`, o._adapter.SafeString(o.OptionName)))
	}

	if o.IsOptionValueDirty == true {
		sets = append(sets, fmt.Sprintf(`option_value = '%s'`, o._adapter.SafeString(o.OptionValue)))
	}

	if o.IsAutoloadDirty == true {
		sets = append(sets, fmt.Sprintf(`autoload = '%s'`, o._adapter.SafeString(o.Autoload)))
	}

	frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'", o._table, strings.Join(sets, `,`))
	err := o._adapter.Execute(frmt)
	if err != nil {
		return err
	}
	return nil
}
func (o *Option) Create() error {
	frmt := fmt.Sprintf("INSERT INTO %s (`option_name`, `option_value`, `autoload`) VALUES ('%s', '%s', '%s')", o._table, o.OptionName, o.OptionValue, o.Autoload)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s led to %s`, frmt, err))
	}
	o.OptionId = o._adapter.LastInsertedId()
	o._new = false
	return nil
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
	PostId    int64
	MetaKey   string
	MetaValue string
	// Dirty markers for smart updates
	IsMetaIdDirty    bool
	IsPostIdDirty    bool
	IsMetaKeyDirty   bool
	IsMetaValueDirty bool
	// Relationships
}

func NewPostMeta(a Adapter) *PostMeta {
	var o PostMeta
	o._table = fmt.Sprintf("%spostmeta", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "meta_id"
	o._new = false
	return &o
}

func (m *PostMeta) GetPrimaryKeyValue() int64 {
	return m.MetaId
}
func (m *PostMeta) GetPrimaryKeyName() string {
	return `meta_id`
}

func (m *PostMeta) GetMetaId() int64 {
	return m.MetaId
}
func (m *PostMeta) SetMetaId(arg int64) {
	m.MetaId = arg
	m.IsMetaIdDirty = true
}

func (m *PostMeta) GetPostId() int64 {
	return m.PostId
}
func (m *PostMeta) SetPostId(arg int64) {
	m.PostId = arg
	m.IsPostIdDirty = true
}

func (m *PostMeta) GetMetaKey() string {
	return m.MetaKey
}
func (m *PostMeta) SetMetaKey(arg string) {
	m.MetaKey = arg
	m.IsMetaKeyDirty = true
}

func (m *PostMeta) GetMetaValue() string {
	return m.MetaValue
}
func (m *PostMeta) SetMetaValue(arg string) {
	m.MetaValue = arg
	m.IsMetaValueDirty = true
}

func (o *PostMeta) Find(_find_by_MetaId int64) (bool, error) {

	var model_slice []*PostMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "meta_id", _find_by_MetaId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return false, o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}

	for _, result := range results {
		ro := PostMeta{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return false, o._adapter.Oops(fmt.Sprintf(`%s`, err))
		}
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return false, o._adapter.Oops(`not found`)
	}
	o.FromPostMeta(model_slice[0])
	return true, nil

}
func (o *PostMeta) FindByPostId(_find_by_PostId int64) ([]*PostMeta, error) {

	var model_slice []*PostMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "post_id", _find_by_PostId)
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *PostMeta) FindByMetaKey(_find_by_MetaKey string) ([]*PostMeta, error) {

	var model_slice []*PostMeta
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *PostMeta) FindByMetaValue(_find_by_MetaValue string) ([]*PostMeta, error) {

	var model_slice []*PostMeta
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}

func (o *PostMeta) FromDBValueMap(m map[string]DBValue) error {
	_MetaId, err := m["meta_id"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.MetaId = _MetaId
	_PostId, err := m["post_id"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.PostId = _PostId
	_MetaKey, err := m["meta_key"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.MetaKey = _MetaKey
	_MetaValue, err := m["meta_value"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.MetaValue = _MetaValue

	return nil
}
func (o *PostMeta) FromPostMeta(m *PostMeta) {
	o.MetaId = m.MetaId
	o.PostId = m.PostId
	o.MetaKey = m.MetaKey
	o.MetaValue = m.MetaValue

}

func (o *PostMeta) Save() error {
	if o._new == true {
		return o.Create()
	}
	var sets []string

	if o.IsPostIdDirty == true {
		sets = append(sets, fmt.Sprintf(`post_id = '%d'`, o.PostId))
	}

	if o.IsMetaKeyDirty == true {
		sets = append(sets, fmt.Sprintf(`meta_key = '%s'`, o._adapter.SafeString(o.MetaKey)))
	}

	if o.IsMetaValueDirty == true {
		sets = append(sets, fmt.Sprintf(`meta_value = '%s'`, o._adapter.SafeString(o.MetaValue)))
	}

	frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'", o._table, strings.Join(sets, `,`))
	err := o._adapter.Execute(frmt)
	if err != nil {
		return err
	}
	return nil
}
func (o *PostMeta) Update() error {
	var sets []string

	if o.IsPostIdDirty == true {
		sets = append(sets, fmt.Sprintf(`post_id = '%d'`, o.PostId))
	}

	if o.IsMetaKeyDirty == true {
		sets = append(sets, fmt.Sprintf(`meta_key = '%s'`, o._adapter.SafeString(o.MetaKey)))
	}

	if o.IsMetaValueDirty == true {
		sets = append(sets, fmt.Sprintf(`meta_value = '%s'`, o._adapter.SafeString(o.MetaValue)))
	}

	frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'", o._table, strings.Join(sets, `,`))
	err := o._adapter.Execute(frmt)
	if err != nil {
		return err
	}
	return nil
}
func (o *PostMeta) Create() error {
	frmt := fmt.Sprintf("INSERT INTO %s (`post_id`, `meta_key`, `meta_value`) VALUES ('%d', '%s', '%s')", o._table, o.PostId, o.MetaKey, o.MetaValue)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s led to %s`, frmt, err))
	}
	o.MetaId = o._adapter.LastInsertedId()
	o._new = false
	return nil
}

func (o *PostMeta) UpdatePostId(_upd_PostId int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_id` = '%d' WHERE `meta_id` = '%d'", o._table, _upd_PostId, o.PostId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.PostId = _upd_PostId
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
	_table              string
	_adapter            Adapter
	_pkey               string // 0 The name of the primary key in this table
	_conds              []string
	_new                bool
	ID                  int64
	PostAuthor          int64
	PostDate            *DateTime
	PostDateGmt         *DateTime
	PostContent         string
	PostTitle           string
	PostExcerpt         string
	PostStatus          string
	CommentStatus       string
	PingStatus          string
	PostPassword        string
	PostName            string
	ToPing              string
	Pinged              string
	PostModified        *DateTime
	PostModifiedGmt     *DateTime
	PostContentFiltered string
	PostParent          int64
	Guid                string
	MenuOrder           int
	PostType            string
	PostMimeType        string
	CommentCount        int64
	// Dirty markers for smart updates
	IsIDDirty                  bool
	IsPostAuthorDirty          bool
	IsPostDateDirty            bool
	IsPostDateGmtDirty         bool
	IsPostContentDirty         bool
	IsPostTitleDirty           bool
	IsPostExcerptDirty         bool
	IsPostStatusDirty          bool
	IsCommentStatusDirty       bool
	IsPingStatusDirty          bool
	IsPostPasswordDirty        bool
	IsPostNameDirty            bool
	IsToPingDirty              bool
	IsPingedDirty              bool
	IsPostModifiedDirty        bool
	IsPostModifiedGmtDirty     bool
	IsPostContentFilteredDirty bool
	IsPostParentDirty          bool
	IsGuidDirty                bool
	IsMenuOrderDirty           bool
	IsPostTypeDirty            bool
	IsPostMimeTypeDirty        bool
	IsCommentCountDirty        bool
	// Relationships
}

func NewPost(a Adapter) *Post {
	var o Post
	o._table = fmt.Sprintf("%sposts", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "ID"
	o._new = false
	return &o
}

func (m *Post) GetPrimaryKeyValue() int64 {
	return m.ID
}
func (m *Post) GetPrimaryKeyName() string {
	return `ID`
}

func (m *Post) GetID() int64 {
	return m.ID
}
func (m *Post) SetID(arg int64) {
	m.ID = arg
	m.IsIDDirty = true
}

func (m *Post) GetPostAuthor() int64 {
	return m.PostAuthor
}
func (m *Post) SetPostAuthor(arg int64) {
	m.PostAuthor = arg
	m.IsPostAuthorDirty = true
}

func (m *Post) GetPostDate() *DateTime {
	return m.PostDate
}
func (m *Post) SetPostDate(arg *DateTime) {
	m.PostDate = arg
	m.IsPostDateDirty = true
}

func (m *Post) GetPostDateGmt() *DateTime {
	return m.PostDateGmt
}
func (m *Post) SetPostDateGmt(arg *DateTime) {
	m.PostDateGmt = arg
	m.IsPostDateGmtDirty = true
}

func (m *Post) GetPostContent() string {
	return m.PostContent
}
func (m *Post) SetPostContent(arg string) {
	m.PostContent = arg
	m.IsPostContentDirty = true
}

func (m *Post) GetPostTitle() string {
	return m.PostTitle
}
func (m *Post) SetPostTitle(arg string) {
	m.PostTitle = arg
	m.IsPostTitleDirty = true
}

func (m *Post) GetPostExcerpt() string {
	return m.PostExcerpt
}
func (m *Post) SetPostExcerpt(arg string) {
	m.PostExcerpt = arg
	m.IsPostExcerptDirty = true
}

func (m *Post) GetPostStatus() string {
	return m.PostStatus
}
func (m *Post) SetPostStatus(arg string) {
	m.PostStatus = arg
	m.IsPostStatusDirty = true
}

func (m *Post) GetCommentStatus() string {
	return m.CommentStatus
}
func (m *Post) SetCommentStatus(arg string) {
	m.CommentStatus = arg
	m.IsCommentStatusDirty = true
}

func (m *Post) GetPingStatus() string {
	return m.PingStatus
}
func (m *Post) SetPingStatus(arg string) {
	m.PingStatus = arg
	m.IsPingStatusDirty = true
}

func (m *Post) GetPostPassword() string {
	return m.PostPassword
}
func (m *Post) SetPostPassword(arg string) {
	m.PostPassword = arg
	m.IsPostPasswordDirty = true
}

func (m *Post) GetPostName() string {
	return m.PostName
}
func (m *Post) SetPostName(arg string) {
	m.PostName = arg
	m.IsPostNameDirty = true
}

func (m *Post) GetToPing() string {
	return m.ToPing
}
func (m *Post) SetToPing(arg string) {
	m.ToPing = arg
	m.IsToPingDirty = true
}

func (m *Post) GetPinged() string {
	return m.Pinged
}
func (m *Post) SetPinged(arg string) {
	m.Pinged = arg
	m.IsPingedDirty = true
}

func (m *Post) GetPostModified() *DateTime {
	return m.PostModified
}
func (m *Post) SetPostModified(arg *DateTime) {
	m.PostModified = arg
	m.IsPostModifiedDirty = true
}

func (m *Post) GetPostModifiedGmt() *DateTime {
	return m.PostModifiedGmt
}
func (m *Post) SetPostModifiedGmt(arg *DateTime) {
	m.PostModifiedGmt = arg
	m.IsPostModifiedGmtDirty = true
}

func (m *Post) GetPostContentFiltered() string {
	return m.PostContentFiltered
}
func (m *Post) SetPostContentFiltered(arg string) {
	m.PostContentFiltered = arg
	m.IsPostContentFilteredDirty = true
}

func (m *Post) GetPostParent() int64 {
	return m.PostParent
}
func (m *Post) SetPostParent(arg int64) {
	m.PostParent = arg
	m.IsPostParentDirty = true
}

func (m *Post) GetGuid() string {
	return m.Guid
}
func (m *Post) SetGuid(arg string) {
	m.Guid = arg
	m.IsGuidDirty = true
}

func (m *Post) GetMenuOrder() int {
	return m.MenuOrder
}
func (m *Post) SetMenuOrder(arg int) {
	m.MenuOrder = arg
	m.IsMenuOrderDirty = true
}

func (m *Post) GetPostType() string {
	return m.PostType
}
func (m *Post) SetPostType(arg string) {
	m.PostType = arg
	m.IsPostTypeDirty = true
}

func (m *Post) GetPostMimeType() string {
	return m.PostMimeType
}
func (m *Post) SetPostMimeType(arg string) {
	m.PostMimeType = arg
	m.IsPostMimeTypeDirty = true
}

func (m *Post) GetCommentCount() int64 {
	return m.CommentCount
}
func (m *Post) SetCommentCount(arg int64) {
	m.CommentCount = arg
	m.IsCommentCountDirty = true
}

func (o *Post) Find(_find_by_ID int64) (bool, error) {

	var model_slice []*Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "ID", _find_by_ID)
	results, err := o._adapter.Query(q)
	if err != nil {
		return false, o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}

	for _, result := range results {
		ro := Post{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return false, o._adapter.Oops(fmt.Sprintf(`%s`, err))
		}
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return false, o._adapter.Oops(`not found`)
	}
	o.FromPost(model_slice[0])
	return true, nil

}
func (o *Post) FindByPostAuthor(_find_by_PostAuthor int64) ([]*Post, error) {

	var model_slice []*Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "post_author", _find_by_PostAuthor)
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Post) FindByPostDate(_find_by_PostDate *DateTime) ([]*Post, error) {

	var model_slice []*Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_date", _find_by_PostDate)
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Post) FindByPostDateGmt(_find_by_PostDateGmt *DateTime) ([]*Post, error) {

	var model_slice []*Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_date_gmt", _find_by_PostDateGmt)
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Post) FindByPostContent(_find_by_PostContent string) ([]*Post, error) {

	var model_slice []*Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_content", _find_by_PostContent)
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Post) FindByPostTitle(_find_by_PostTitle string) ([]*Post, error) {

	var model_slice []*Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_title", _find_by_PostTitle)
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Post) FindByPostExcerpt(_find_by_PostExcerpt string) ([]*Post, error) {

	var model_slice []*Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_excerpt", _find_by_PostExcerpt)
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Post) FindByPostStatus(_find_by_PostStatus string) ([]*Post, error) {

	var model_slice []*Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_status", _find_by_PostStatus)
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Post) FindByCommentStatus(_find_by_CommentStatus string) ([]*Post, error) {

	var model_slice []*Post
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Post) FindByPingStatus(_find_by_PingStatus string) ([]*Post, error) {

	var model_slice []*Post
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Post) FindByPostPassword(_find_by_PostPassword string) ([]*Post, error) {

	var model_slice []*Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_password", _find_by_PostPassword)
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Post) FindByPostName(_find_by_PostName string) ([]*Post, error) {

	var model_slice []*Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_name", _find_by_PostName)
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Post) FindByToPing(_find_by_ToPing string) ([]*Post, error) {

	var model_slice []*Post
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Post) FindByPinged(_find_by_Pinged string) ([]*Post, error) {

	var model_slice []*Post
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Post) FindByPostModified(_find_by_PostModified *DateTime) ([]*Post, error) {

	var model_slice []*Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_modified", _find_by_PostModified)
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Post) FindByPostModifiedGmt(_find_by_PostModifiedGmt *DateTime) ([]*Post, error) {

	var model_slice []*Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_modified_gmt", _find_by_PostModifiedGmt)
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Post) FindByPostContentFiltered(_find_by_PostContentFiltered string) ([]*Post, error) {

	var model_slice []*Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_content_filtered", _find_by_PostContentFiltered)
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Post) FindByPostParent(_find_by_PostParent int64) ([]*Post, error) {

	var model_slice []*Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "post_parent", _find_by_PostParent)
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Post) FindByGuid(_find_by_Guid string) ([]*Post, error) {

	var model_slice []*Post
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Post) FindByMenuOrder(_find_by_MenuOrder int) ([]*Post, error) {

	var model_slice []*Post
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Post) FindByPostType(_find_by_PostType string) ([]*Post, error) {

	var model_slice []*Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_type", _find_by_PostType)
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Post) FindByPostMimeType(_find_by_PostMimeType string) ([]*Post, error) {

	var model_slice []*Post
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'", o._table, "post_mime_type", _find_by_PostMimeType)
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Post) FindByCommentCount(_find_by_CommentCount int64) ([]*Post, error) {

	var model_slice []*Post
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}

func (o *Post) FromDBValueMap(m map[string]DBValue) error {
	_ID, err := m["ID"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.ID = _ID
	_PostAuthor, err := m["post_author"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.PostAuthor = _PostAuthor
	_PostDate, err := m["post_date"].AsDateTime()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.PostDate = _PostDate
	_PostDateGmt, err := m["post_date_gmt"].AsDateTime()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.PostDateGmt = _PostDateGmt
	_PostContent, err := m["post_content"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.PostContent = _PostContent
	_PostTitle, err := m["post_title"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.PostTitle = _PostTitle
	_PostExcerpt, err := m["post_excerpt"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.PostExcerpt = _PostExcerpt
	_PostStatus, err := m["post_status"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.PostStatus = _PostStatus
	_CommentStatus, err := m["comment_status"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.CommentStatus = _CommentStatus
	_PingStatus, err := m["ping_status"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.PingStatus = _PingStatus
	_PostPassword, err := m["post_password"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.PostPassword = _PostPassword
	_PostName, err := m["post_name"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.PostName = _PostName
	_ToPing, err := m["to_ping"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.ToPing = _ToPing
	_Pinged, err := m["pinged"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.Pinged = _Pinged
	_PostModified, err := m["post_modified"].AsDateTime()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.PostModified = _PostModified
	_PostModifiedGmt, err := m["post_modified_gmt"].AsDateTime()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.PostModifiedGmt = _PostModifiedGmt
	_PostContentFiltered, err := m["post_content_filtered"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.PostContentFiltered = _PostContentFiltered
	_PostParent, err := m["post_parent"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.PostParent = _PostParent
	_Guid, err := m["guid"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.Guid = _Guid
	_MenuOrder, err := m["menu_order"].AsInt()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.MenuOrder = _MenuOrder
	_PostType, err := m["post_type"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.PostType = _PostType
	_PostMimeType, err := m["post_mime_type"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.PostMimeType = _PostMimeType
	_CommentCount, err := m["comment_count"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.CommentCount = _CommentCount

	return nil
}
func (o *Post) FromPost(m *Post) {
	o.ID = m.ID
	o.PostAuthor = m.PostAuthor
	o.PostDate = m.PostDate
	o.PostDateGmt = m.PostDateGmt
	o.PostContent = m.PostContent
	o.PostTitle = m.PostTitle
	o.PostExcerpt = m.PostExcerpt
	o.PostStatus = m.PostStatus
	o.CommentStatus = m.CommentStatus
	o.PingStatus = m.PingStatus
	o.PostPassword = m.PostPassword
	o.PostName = m.PostName
	o.ToPing = m.ToPing
	o.Pinged = m.Pinged
	o.PostModified = m.PostModified
	o.PostModifiedGmt = m.PostModifiedGmt
	o.PostContentFiltered = m.PostContentFiltered
	o.PostParent = m.PostParent
	o.Guid = m.Guid
	o.MenuOrder = m.MenuOrder
	o.PostType = m.PostType
	o.PostMimeType = m.PostMimeType
	o.CommentCount = m.CommentCount

}

func (o *Post) Save() error {
	if o._new == true {
		return o.Create()
	}
	var sets []string

	if o.IsPostAuthorDirty == true {
		sets = append(sets, fmt.Sprintf(`post_author = '%d'`, o.PostAuthor))
	}

	if o.IsPostDateDirty == true {
		sets = append(sets, fmt.Sprintf(`post_date = '%s'`, o.PostDate))
	}

	if o.IsPostDateGmtDirty == true {
		sets = append(sets, fmt.Sprintf(`post_date_gmt = '%s'`, o.PostDateGmt))
	}

	if o.IsPostContentDirty == true {
		sets = append(sets, fmt.Sprintf(`post_content = '%s'`, o._adapter.SafeString(o.PostContent)))
	}

	if o.IsPostTitleDirty == true {
		sets = append(sets, fmt.Sprintf(`post_title = '%s'`, o._adapter.SafeString(o.PostTitle)))
	}

	if o.IsPostExcerptDirty == true {
		sets = append(sets, fmt.Sprintf(`post_excerpt = '%s'`, o._adapter.SafeString(o.PostExcerpt)))
	}

	if o.IsPostStatusDirty == true {
		sets = append(sets, fmt.Sprintf(`post_status = '%s'`, o._adapter.SafeString(o.PostStatus)))
	}

	if o.IsCommentStatusDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_status = '%s'`, o._adapter.SafeString(o.CommentStatus)))
	}

	if o.IsPingStatusDirty == true {
		sets = append(sets, fmt.Sprintf(`ping_status = '%s'`, o._adapter.SafeString(o.PingStatus)))
	}

	if o.IsPostPasswordDirty == true {
		sets = append(sets, fmt.Sprintf(`post_password = '%s'`, o._adapter.SafeString(o.PostPassword)))
	}

	if o.IsPostNameDirty == true {
		sets = append(sets, fmt.Sprintf(`post_name = '%s'`, o._adapter.SafeString(o.PostName)))
	}

	if o.IsToPingDirty == true {
		sets = append(sets, fmt.Sprintf(`to_ping = '%s'`, o._adapter.SafeString(o.ToPing)))
	}

	if o.IsPingedDirty == true {
		sets = append(sets, fmt.Sprintf(`pinged = '%s'`, o._adapter.SafeString(o.Pinged)))
	}

	if o.IsPostModifiedDirty == true {
		sets = append(sets, fmt.Sprintf(`post_modified = '%s'`, o.PostModified))
	}

	if o.IsPostModifiedGmtDirty == true {
		sets = append(sets, fmt.Sprintf(`post_modified_gmt = '%s'`, o.PostModifiedGmt))
	}

	if o.IsPostContentFilteredDirty == true {
		sets = append(sets, fmt.Sprintf(`post_content_filtered = '%s'`, o._adapter.SafeString(o.PostContentFiltered)))
	}

	if o.IsPostParentDirty == true {
		sets = append(sets, fmt.Sprintf(`post_parent = '%d'`, o.PostParent))
	}

	if o.IsGuidDirty == true {
		sets = append(sets, fmt.Sprintf(`guid = '%s'`, o._adapter.SafeString(o.Guid)))
	}

	if o.IsMenuOrderDirty == true {
		sets = append(sets, fmt.Sprintf(`menu_order = '%d'`, o.MenuOrder))
	}

	if o.IsPostTypeDirty == true {
		sets = append(sets, fmt.Sprintf(`post_type = '%s'`, o._adapter.SafeString(o.PostType)))
	}

	if o.IsPostMimeTypeDirty == true {
		sets = append(sets, fmt.Sprintf(`post_mime_type = '%s'`, o._adapter.SafeString(o.PostMimeType)))
	}

	if o.IsCommentCountDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_count = '%d'`, o.CommentCount))
	}

	frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'", o._table, strings.Join(sets, `,`))
	err := o._adapter.Execute(frmt)
	if err != nil {
		return err
	}
	return nil
}
func (o *Post) Update() error {
	var sets []string

	if o.IsPostAuthorDirty == true {
		sets = append(sets, fmt.Sprintf(`post_author = '%d'`, o.PostAuthor))
	}

	if o.IsPostDateDirty == true {
		sets = append(sets, fmt.Sprintf(`post_date = '%s'`, o.PostDate))
	}

	if o.IsPostDateGmtDirty == true {
		sets = append(sets, fmt.Sprintf(`post_date_gmt = '%s'`, o.PostDateGmt))
	}

	if o.IsPostContentDirty == true {
		sets = append(sets, fmt.Sprintf(`post_content = '%s'`, o._adapter.SafeString(o.PostContent)))
	}

	if o.IsPostTitleDirty == true {
		sets = append(sets, fmt.Sprintf(`post_title = '%s'`, o._adapter.SafeString(o.PostTitle)))
	}

	if o.IsPostExcerptDirty == true {
		sets = append(sets, fmt.Sprintf(`post_excerpt = '%s'`, o._adapter.SafeString(o.PostExcerpt)))
	}

	if o.IsPostStatusDirty == true {
		sets = append(sets, fmt.Sprintf(`post_status = '%s'`, o._adapter.SafeString(o.PostStatus)))
	}

	if o.IsCommentStatusDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_status = '%s'`, o._adapter.SafeString(o.CommentStatus)))
	}

	if o.IsPingStatusDirty == true {
		sets = append(sets, fmt.Sprintf(`ping_status = '%s'`, o._adapter.SafeString(o.PingStatus)))
	}

	if o.IsPostPasswordDirty == true {
		sets = append(sets, fmt.Sprintf(`post_password = '%s'`, o._adapter.SafeString(o.PostPassword)))
	}

	if o.IsPostNameDirty == true {
		sets = append(sets, fmt.Sprintf(`post_name = '%s'`, o._adapter.SafeString(o.PostName)))
	}

	if o.IsToPingDirty == true {
		sets = append(sets, fmt.Sprintf(`to_ping = '%s'`, o._adapter.SafeString(o.ToPing)))
	}

	if o.IsPingedDirty == true {
		sets = append(sets, fmt.Sprintf(`pinged = '%s'`, o._adapter.SafeString(o.Pinged)))
	}

	if o.IsPostModifiedDirty == true {
		sets = append(sets, fmt.Sprintf(`post_modified = '%s'`, o.PostModified))
	}

	if o.IsPostModifiedGmtDirty == true {
		sets = append(sets, fmt.Sprintf(`post_modified_gmt = '%s'`, o.PostModifiedGmt))
	}

	if o.IsPostContentFilteredDirty == true {
		sets = append(sets, fmt.Sprintf(`post_content_filtered = '%s'`, o._adapter.SafeString(o.PostContentFiltered)))
	}

	if o.IsPostParentDirty == true {
		sets = append(sets, fmt.Sprintf(`post_parent = '%d'`, o.PostParent))
	}

	if o.IsGuidDirty == true {
		sets = append(sets, fmt.Sprintf(`guid = '%s'`, o._adapter.SafeString(o.Guid)))
	}

	if o.IsMenuOrderDirty == true {
		sets = append(sets, fmt.Sprintf(`menu_order = '%d'`, o.MenuOrder))
	}

	if o.IsPostTypeDirty == true {
		sets = append(sets, fmt.Sprintf(`post_type = '%s'`, o._adapter.SafeString(o.PostType)))
	}

	if o.IsPostMimeTypeDirty == true {
		sets = append(sets, fmt.Sprintf(`post_mime_type = '%s'`, o._adapter.SafeString(o.PostMimeType)))
	}

	if o.IsCommentCountDirty == true {
		sets = append(sets, fmt.Sprintf(`comment_count = '%d'`, o.CommentCount))
	}

	frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'", o._table, strings.Join(sets, `,`))
	err := o._adapter.Execute(frmt)
	if err != nil {
		return err
	}
	return nil
}
func (o *Post) Create() error {
	frmt := fmt.Sprintf("INSERT INTO %s (`post_author`, `post_date`, `post_date_gmt`, `post_content`, `post_title`, `post_excerpt`, `post_status`, `comment_status`, `ping_status`, `post_password`, `post_name`, `to_ping`, `pinged`, `post_modified`, `post_modified_gmt`, `post_content_filtered`, `post_parent`, `guid`, `menu_order`, `post_type`, `post_mime_type`, `comment_count`) VALUES ('%d', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%d', '%s', '%d', '%s', '%s', '%d')", o._table, o.PostAuthor, o.PostDate.ToString(), o.PostDateGmt.ToString(), o.PostContent, o.PostTitle, o.PostExcerpt, o.PostStatus, o.CommentStatus, o.PingStatus, o.PostPassword, o.PostName, o.ToPing, o.Pinged, o.PostModified.ToString(), o.PostModifiedGmt.ToString(), o.PostContentFiltered, o.PostParent, o.Guid, o.MenuOrder, o.PostType, o.PostMimeType, o.CommentCount)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s led to %s`, frmt, err))
	}
	o.ID = o._adapter.LastInsertedId()
	o._new = false
	return nil
}

func (o *Post) UpdatePostAuthor(_upd_PostAuthor int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_author` = '%d' WHERE `ID` = '%d'", o._table, _upd_PostAuthor, o.PostAuthor)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.PostAuthor = _upd_PostAuthor
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdatePostDate(_upd_PostDate *DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_date` = '%s' WHERE `ID` = '%d'", o._table, _upd_PostDate, o.PostDate)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.PostDate = _upd_PostDate
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdatePostDateGmt(_upd_PostDateGmt *DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_date_gmt` = '%s' WHERE `ID` = '%d'", o._table, _upd_PostDateGmt, o.PostDateGmt)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.PostDateGmt = _upd_PostDateGmt
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdatePostContent(_upd_PostContent string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_content` = '%s' WHERE `ID` = '%d'", o._table, _upd_PostContent, o.PostContent)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.PostContent = _upd_PostContent
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdatePostTitle(_upd_PostTitle string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_title` = '%s' WHERE `ID` = '%d'", o._table, _upd_PostTitle, o.PostTitle)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.PostTitle = _upd_PostTitle
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdatePostExcerpt(_upd_PostExcerpt string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_excerpt` = '%s' WHERE `ID` = '%d'", o._table, _upd_PostExcerpt, o.PostExcerpt)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.PostExcerpt = _upd_PostExcerpt
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdatePostStatus(_upd_PostStatus string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_status` = '%s' WHERE `ID` = '%d'", o._table, _upd_PostStatus, o.PostStatus)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.PostStatus = _upd_PostStatus
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

func (o *Post) UpdatePostPassword(_upd_PostPassword string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_password` = '%s' WHERE `ID` = '%d'", o._table, _upd_PostPassword, o.PostPassword)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.PostPassword = _upd_PostPassword
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdatePostName(_upd_PostName string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_name` = '%s' WHERE `ID` = '%d'", o._table, _upd_PostName, o.PostName)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.PostName = _upd_PostName
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

func (o *Post) UpdatePostModified(_upd_PostModified *DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_modified` = '%s' WHERE `ID` = '%d'", o._table, _upd_PostModified, o.PostModified)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.PostModified = _upd_PostModified
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdatePostModifiedGmt(_upd_PostModifiedGmt *DateTime) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_modified_gmt` = '%s' WHERE `ID` = '%d'", o._table, _upd_PostModifiedGmt, o.PostModifiedGmt)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.PostModifiedGmt = _upd_PostModifiedGmt
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdatePostContentFiltered(_upd_PostContentFiltered string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_content_filtered` = '%s' WHERE `ID` = '%d'", o._table, _upd_PostContentFiltered, o.PostContentFiltered)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.PostContentFiltered = _upd_PostContentFiltered
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdatePostParent(_upd_PostParent int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_parent` = '%d' WHERE `ID` = '%d'", o._table, _upd_PostParent, o.PostParent)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.PostParent = _upd_PostParent
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

func (o *Post) UpdatePostType(_upd_PostType string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_type` = '%s' WHERE `ID` = '%d'", o._table, _upd_PostType, o.PostType)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.PostType = _upd_PostType
	return o._adapter.AffectedRows(), nil
}

func (o *Post) UpdatePostMimeType(_upd_PostMimeType string) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `post_mime_type` = '%s' WHERE `ID` = '%d'", o._table, _upd_PostMimeType, o.PostMimeType)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.PostMimeType = _upd_PostMimeType
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
	// Dirty markers for smart updates
	IsObjectIdDirty       bool
	IsTermTaxonomyIdDirty bool
	IsTermOrderDirty      bool
	// Relationships
}

func NewTermRelationship(a Adapter) *TermRelationship {
	var o TermRelationship
	o._table = fmt.Sprintf("%sterm_relationships", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "term_taxonomy_id"
	o._new = false
	return &o
}

func (m *TermRelationship) GetPrimaryKeyValue() int64 {
	return m.TermTaxonomyId
}
func (m *TermRelationship) GetPrimaryKeyName() string {
	return `term_taxonomy_id`
}

func (m *TermRelationship) GetObjectId() int64 {
	return m.ObjectId
}
func (m *TermRelationship) SetObjectId(arg int64) {
	m.ObjectId = arg
	m.IsObjectIdDirty = true
}

func (m *TermRelationship) GetTermTaxonomyId() int64 {
	return m.TermTaxonomyId
}
func (m *TermRelationship) SetTermTaxonomyId(arg int64) {
	m.TermTaxonomyId = arg
	m.IsTermTaxonomyIdDirty = true
}

func (m *TermRelationship) GetTermOrder() int {
	return m.TermOrder
}
func (m *TermRelationship) SetTermOrder(arg int) {
	m.TermOrder = arg
	m.IsTermOrderDirty = true
}

func (o *TermRelationship) FindByObjectId(_find_by_ObjectId int64) ([]*TermRelationship, error) {

	var model_slice []*TermRelationship
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *TermRelationship) Find(termId int64, objectId int64) (bool, error) {

	var model_slice []*TermRelationship
	q := fmt.Sprintf("SELECT * FROM %s WHERE `term_taxonomy_id` = '%d' AND `object_id` = '%d'", o._table, termId, objectId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return false, err
	}

	for _, result := range results {
		ro := TermRelationship{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return false, err
		}
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return false, errors.New("not found")
	}
	o.FromTermRelationship(model_slice[0])
	return true, nil

}
func (o *TermRelationship) FindByTermOrder(_find_by_TermOrder int) ([]*TermRelationship, error) {

	var model_slice []*TermRelationship
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}

func (o *TermRelationship) FromDBValueMap(m map[string]DBValue) error {
	_ObjectId, err := m["object_id"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.ObjectId = _ObjectId
	_TermTaxonomyId, err := m["term_taxonomy_id"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.TermTaxonomyId = _TermTaxonomyId
	_TermOrder, err := m["term_order"].AsInt()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.TermOrder = _TermOrder

	return nil
}
func (o *TermRelationship) FromTermRelationship(m *TermRelationship) {
	o.ObjectId = m.ObjectId
	o.TermTaxonomyId = m.TermTaxonomyId
	o.TermOrder = m.TermOrder

}

func (o *TermRelationship) Save() error {
	if o._new == true {
		return o.Create()
	}
	var sets []string

	if o.IsObjectIdDirty == true {
		sets = append(sets, fmt.Sprintf(`object_id = '%d'`, o.ObjectId))
	}

	if o.IsTermTaxonomyIdDirty == true {
		sets = append(sets, fmt.Sprintf(`term_taxonomy_id = '%d'`, o.TermTaxonomyId))
	}

	if o.IsTermOrderDirty == true {
		sets = append(sets, fmt.Sprintf(`term_order = '%d'`, o.TermOrder))
	}

	frmt := fmt.Sprintf("UPDATE %s SET %s WHERE `term_taxonomy_id` = '%d' AND object_id = '%d'", o._table, strings.Join(sets, `,`))
	err := o._adapter.Execute(frmt)
	if err != nil {
		return err
	}
	return nil
}
func (o *TermRelationship) Update() error {
	var sets []string

	if o.IsObjectIdDirty == true {
		sets = append(sets, fmt.Sprintf(`object_id = '%d'`, o.ObjectId))
	}

	if o.IsTermTaxonomyIdDirty == true {
		sets = append(sets, fmt.Sprintf(`term_taxonomy_id = '%d'`, o.TermTaxonomyId))
	}

	if o.IsTermOrderDirty == true {
		sets = append(sets, fmt.Sprintf(`term_order = '%d'`, o.TermOrder))
	}

	frmt := fmt.Sprintf("UPDATE %s SET %s WHERE `term_taxonomy_id` = '%d' AND object_id = '%d'", o._table, strings.Join(sets, `,`))
	err := o._adapter.Execute(frmt)
	if err != nil {
		return err
	}
	return nil
}
func (o *TermRelationship) Create() error {
	frmt := fmt.Sprintf("INSERT INTO %s (`object_id`, `term_taxonomy_id`, `term_order`) VALUES ('%d', '%d', '%d')", o._table, o.ObjectId, o.TermTaxonomyId, o.TermOrder)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s led to %s`, frmt, err))
	}

	o._new = false
	return nil
}

func (o *TermRelationship) UpdateObjectId(_upd_ObjectId int64) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `object_id` = '%d' WHERE term_taxonomy_id = '%d' AND object_id = '%d'", o._table, _upd_ObjectId, o.TermTaxonomyId, o.ObjectId)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return 0, err
	}
	o.ObjectId = _upd_ObjectId
	return o._adapter.AffectedRows(), nil
}

func (o *TermRelationship) UpdateTermOrder(_upd_TermOrder int) (int64, error) {
	frmt := fmt.Sprintf("UPDATE %s SET `term_order` = '%d' WHERE term_taxonomy_id = '%d' AND object_id = '%d'", o._table, _upd_TermOrder, o.TermTaxonomyId, o.ObjectId)
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
	// Dirty markers for smart updates
	IsTermTaxonomyIdDirty bool
	IsTermIdDirty         bool
	IsTaxonomyDirty       bool
	IsDescriptionDirty    bool
	IsParentDirty         bool
	IsCountDirty          bool
	// Relationships
}

func NewTermTaxonomy(a Adapter) *TermTaxonomy {
	var o TermTaxonomy
	o._table = fmt.Sprintf("%sterm_taxonomy", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "term_taxonomy_id"
	o._new = false
	return &o
}

func (m *TermTaxonomy) GetPrimaryKeyValue() int64 {
	return m.TermTaxonomyId
}
func (m *TermTaxonomy) GetPrimaryKeyName() string {
	return `term_taxonomy_id`
}

func (m *TermTaxonomy) GetTermTaxonomyId() int64 {
	return m.TermTaxonomyId
}
func (m *TermTaxonomy) SetTermTaxonomyId(arg int64) {
	m.TermTaxonomyId = arg
	m.IsTermTaxonomyIdDirty = true
}

func (m *TermTaxonomy) GetTermId() int64 {
	return m.TermId
}
func (m *TermTaxonomy) SetTermId(arg int64) {
	m.TermId = arg
	m.IsTermIdDirty = true
}

func (m *TermTaxonomy) GetTaxonomy() string {
	return m.Taxonomy
}
func (m *TermTaxonomy) SetTaxonomy(arg string) {
	m.Taxonomy = arg
	m.IsTaxonomyDirty = true
}

func (m *TermTaxonomy) GetDescription() string {
	return m.Description
}
func (m *TermTaxonomy) SetDescription(arg string) {
	m.Description = arg
	m.IsDescriptionDirty = true
}

func (m *TermTaxonomy) GetParent() int64 {
	return m.Parent
}
func (m *TermTaxonomy) SetParent(arg int64) {
	m.Parent = arg
	m.IsParentDirty = true
}

func (m *TermTaxonomy) GetCount() int64 {
	return m.Count
}
func (m *TermTaxonomy) SetCount(arg int64) {
	m.Count = arg
	m.IsCountDirty = true
}

func (o *TermTaxonomy) Find(_find_by_TermTaxonomyId int64) (bool, error) {

	var model_slice []*TermTaxonomy
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "term_taxonomy_id", _find_by_TermTaxonomyId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return false, o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}

	for _, result := range results {
		ro := TermTaxonomy{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return false, o._adapter.Oops(fmt.Sprintf(`%s`, err))
		}
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return false, o._adapter.Oops(`not found`)
	}
	o.FromTermTaxonomy(model_slice[0])
	return true, nil

}
func (o *TermTaxonomy) FindByTermId(_find_by_TermId int64) ([]*TermTaxonomy, error) {

	var model_slice []*TermTaxonomy
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *TermTaxonomy) FindByTaxonomy(_find_by_Taxonomy string) ([]*TermTaxonomy, error) {

	var model_slice []*TermTaxonomy
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *TermTaxonomy) FindByDescription(_find_by_Description string) ([]*TermTaxonomy, error) {

	var model_slice []*TermTaxonomy
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *TermTaxonomy) FindByParent(_find_by_Parent int64) ([]*TermTaxonomy, error) {

	var model_slice []*TermTaxonomy
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *TermTaxonomy) FindByCount(_find_by_Count int64) ([]*TermTaxonomy, error) {

	var model_slice []*TermTaxonomy
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}

func (o *TermTaxonomy) FromDBValueMap(m map[string]DBValue) error {
	_TermTaxonomyId, err := m["term_taxonomy_id"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.TermTaxonomyId = _TermTaxonomyId
	_TermId, err := m["term_id"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.TermId = _TermId
	_Taxonomy, err := m["taxonomy"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.Taxonomy = _Taxonomy
	_Description, err := m["description"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.Description = _Description
	_Parent, err := m["parent"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.Parent = _Parent
	_Count, err := m["count"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.Count = _Count

	return nil
}
func (o *TermTaxonomy) FromTermTaxonomy(m *TermTaxonomy) {
	o.TermTaxonomyId = m.TermTaxonomyId
	o.TermId = m.TermId
	o.Taxonomy = m.Taxonomy
	o.Description = m.Description
	o.Parent = m.Parent
	o.Count = m.Count

}

func (o *TermTaxonomy) Save() error {
	if o._new == true {
		return o.Create()
	}
	var sets []string

	if o.IsTermIdDirty == true {
		sets = append(sets, fmt.Sprintf(`term_id = '%d'`, o.TermId))
	}

	if o.IsTaxonomyDirty == true {
		sets = append(sets, fmt.Sprintf(`taxonomy = '%s'`, o._adapter.SafeString(o.Taxonomy)))
	}

	if o.IsDescriptionDirty == true {
		sets = append(sets, fmt.Sprintf(`description = '%s'`, o._adapter.SafeString(o.Description)))
	}

	if o.IsParentDirty == true {
		sets = append(sets, fmt.Sprintf(`parent = '%d'`, o.Parent))
	}

	if o.IsCountDirty == true {
		sets = append(sets, fmt.Sprintf(`count = '%d'`, o.Count))
	}

	frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'", o._table, strings.Join(sets, `,`))
	err := o._adapter.Execute(frmt)
	if err != nil {
		return err
	}
	return nil
}
func (o *TermTaxonomy) Update() error {
	var sets []string

	if o.IsTermIdDirty == true {
		sets = append(sets, fmt.Sprintf(`term_id = '%d'`, o.TermId))
	}

	if o.IsTaxonomyDirty == true {
		sets = append(sets, fmt.Sprintf(`taxonomy = '%s'`, o._adapter.SafeString(o.Taxonomy)))
	}

	if o.IsDescriptionDirty == true {
		sets = append(sets, fmt.Sprintf(`description = '%s'`, o._adapter.SafeString(o.Description)))
	}

	if o.IsParentDirty == true {
		sets = append(sets, fmt.Sprintf(`parent = '%d'`, o.Parent))
	}

	if o.IsCountDirty == true {
		sets = append(sets, fmt.Sprintf(`count = '%d'`, o.Count))
	}

	frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'", o._table, strings.Join(sets, `,`))
	err := o._adapter.Execute(frmt)
	if err != nil {
		return err
	}
	return nil
}
func (o *TermTaxonomy) Create() error {
	frmt := fmt.Sprintf("INSERT INTO %s (`term_id`, `taxonomy`, `description`, `parent`, `count`) VALUES ('%d', '%s', '%s', '%d', '%d')", o._table, o.TermId, o.Taxonomy, o.Description, o.Parent, o.Count)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s led to %s`, frmt, err))
	}
	o.TermTaxonomyId = o._adapter.LastInsertedId()
	o._new = false
	return nil
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
	// Dirty markers for smart updates
	IsTermIdDirty    bool
	IsNameDirty      bool
	IsSlugDirty      bool
	IsTermGroupDirty bool
	// Relationships
}

func NewTerm(a Adapter) *Term {
	var o Term
	o._table = fmt.Sprintf("%sterms", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "term_id"
	o._new = false
	return &o
}

func (m *Term) GetPrimaryKeyValue() int64 {
	return m.TermId
}
func (m *Term) GetPrimaryKeyName() string {
	return `term_id`
}

func (m *Term) GetTermId() int64 {
	return m.TermId
}
func (m *Term) SetTermId(arg int64) {
	m.TermId = arg
	m.IsTermIdDirty = true
}

func (m *Term) GetName() string {
	return m.Name
}
func (m *Term) SetName(arg string) {
	m.Name = arg
	m.IsNameDirty = true
}

func (m *Term) GetSlug() string {
	return m.Slug
}
func (m *Term) SetSlug(arg string) {
	m.Slug = arg
	m.IsSlugDirty = true
}

func (m *Term) GetTermGroup() int64 {
	return m.TermGroup
}
func (m *Term) SetTermGroup(arg int64) {
	m.TermGroup = arg
	m.IsTermGroupDirty = true
}

func (o *Term) Find(_find_by_TermId int64) (bool, error) {

	var model_slice []*Term
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "term_id", _find_by_TermId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return false, o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}

	for _, result := range results {
		ro := Term{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return false, o._adapter.Oops(fmt.Sprintf(`%s`, err))
		}
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return false, o._adapter.Oops(`not found`)
	}
	o.FromTerm(model_slice[0])
	return true, nil

}
func (o *Term) FindByName(_find_by_Name string) ([]*Term, error) {

	var model_slice []*Term
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Term) FindBySlug(_find_by_Slug string) ([]*Term, error) {

	var model_slice []*Term
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *Term) FindByTermGroup(_find_by_TermGroup int64) ([]*Term, error) {

	var model_slice []*Term
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}

func (o *Term) FromDBValueMap(m map[string]DBValue) error {
	_TermId, err := m["term_id"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.TermId = _TermId
	_Name, err := m["name"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.Name = _Name
	_Slug, err := m["slug"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.Slug = _Slug
	_TermGroup, err := m["term_group"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.TermGroup = _TermGroup

	return nil
}
func (o *Term) FromTerm(m *Term) {
	o.TermId = m.TermId
	o.Name = m.Name
	o.Slug = m.Slug
	o.TermGroup = m.TermGroup

}

func (o *Term) Save() error {
	if o._new == true {
		return o.Create()
	}
	var sets []string

	if o.IsNameDirty == true {
		sets = append(sets, fmt.Sprintf(`name = '%s'`, o._adapter.SafeString(o.Name)))
	}

	if o.IsSlugDirty == true {
		sets = append(sets, fmt.Sprintf(`slug = '%s'`, o._adapter.SafeString(o.Slug)))
	}

	if o.IsTermGroupDirty == true {
		sets = append(sets, fmt.Sprintf(`term_group = '%d'`, o.TermGroup))
	}

	frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'", o._table, strings.Join(sets, `,`))
	err := o._adapter.Execute(frmt)
	if err != nil {
		return err
	}
	return nil
}
func (o *Term) Update() error {
	var sets []string

	if o.IsNameDirty == true {
		sets = append(sets, fmt.Sprintf(`name = '%s'`, o._adapter.SafeString(o.Name)))
	}

	if o.IsSlugDirty == true {
		sets = append(sets, fmt.Sprintf(`slug = '%s'`, o._adapter.SafeString(o.Slug)))
	}

	if o.IsTermGroupDirty == true {
		sets = append(sets, fmt.Sprintf(`term_group = '%d'`, o.TermGroup))
	}

	frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'", o._table, strings.Join(sets, `,`))
	err := o._adapter.Execute(frmt)
	if err != nil {
		return err
	}
	return nil
}
func (o *Term) Create() error {
	frmt := fmt.Sprintf("INSERT INTO %s (`name`, `slug`, `term_group`) VALUES ('%s', '%s', '%d')", o._table, o.Name, o.Slug, o.TermGroup)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s led to %s`, frmt, err))
	}
	o.TermId = o._adapter.LastInsertedId()
	o._new = false
	return nil
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
	// Dirty markers for smart updates
	IsUMetaIdDirty   bool
	IsUserIdDirty    bool
	IsMetaKeyDirty   bool
	IsMetaValueDirty bool
	// Relationships
}

func NewUserMeta(a Adapter) *UserMeta {
	var o UserMeta
	o._table = fmt.Sprintf("%susermeta", a.DatabasePrefix())
	o._adapter = a
	o._pkey = "umeta_id"
	o._new = false
	return &o
}

func (m *UserMeta) GetPrimaryKeyValue() int64 {
	return m.UMetaId
}
func (m *UserMeta) GetPrimaryKeyName() string {
	return `umeta_id`
}

func (m *UserMeta) GetUMetaId() int64 {
	return m.UMetaId
}
func (m *UserMeta) SetUMetaId(arg int64) {
	m.UMetaId = arg
	m.IsUMetaIdDirty = true
}

func (m *UserMeta) GetUserId() int64 {
	return m.UserId
}
func (m *UserMeta) SetUserId(arg int64) {
	m.UserId = arg
	m.IsUserIdDirty = true
}

func (m *UserMeta) GetMetaKey() string {
	return m.MetaKey
}
func (m *UserMeta) SetMetaKey(arg string) {
	m.MetaKey = arg
	m.IsMetaKeyDirty = true
}

func (m *UserMeta) GetMetaValue() string {
	return m.MetaValue
}
func (m *UserMeta) SetMetaValue(arg string) {
	m.MetaValue = arg
	m.IsMetaValueDirty = true
}

func (o *UserMeta) Find(_find_by_UMetaId int64) (bool, error) {

	var model_slice []*UserMeta
	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "umeta_id", _find_by_UMetaId)
	results, err := o._adapter.Query(q)
	if err != nil {
		return false, o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}

	for _, result := range results {
		ro := UserMeta{}
		err = ro.FromDBValueMap(result)
		if err != nil {
			return false, o._adapter.Oops(fmt.Sprintf(`%s`, err))
		}
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return false, o._adapter.Oops(`not found`)
	}
	o.FromUserMeta(model_slice[0])
	return true, nil

}
func (o *UserMeta) FindByUserId(_find_by_UserId int64) ([]*UserMeta, error) {

	var model_slice []*UserMeta
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *UserMeta) FindByMetaKey(_find_by_MetaKey string) ([]*UserMeta, error) {

	var model_slice []*UserMeta
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}
func (o *UserMeta) FindByMetaValue(_find_by_MetaValue string) ([]*UserMeta, error) {

	var model_slice []*UserMeta
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
		model_slice = append(model_slice, &ro)
	}

	if len(model_slice) == 0 {
		// there was an error!
		return nil, o._adapter.Oops(`no results`)
	}
	return model_slice, nil

}

func (o *UserMeta) FromDBValueMap(m map[string]DBValue) error {
	_UMetaId, err := m["umeta_id"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.UMetaId = _UMetaId
	_UserId, err := m["user_id"].AsInt64()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.UserId = _UserId
	_MetaKey, err := m["meta_key"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.MetaKey = _MetaKey
	_MetaValue, err := m["meta_value"].AsString()
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s`, err))
	}
	o.MetaValue = _MetaValue

	return nil
}
func (o *UserMeta) FromUserMeta(m *UserMeta) {
	o.UMetaId = m.UMetaId
	o.UserId = m.UserId
	o.MetaKey = m.MetaKey
	o.MetaValue = m.MetaValue

}

func (o *UserMeta) Save() error {
	if o._new == true {
		return o.Create()
	}
	var sets []string

	if o.IsUserIdDirty == true {
		sets = append(sets, fmt.Sprintf(`user_id = '%d'`, o.UserId))
	}

	if o.IsMetaKeyDirty == true {
		sets = append(sets, fmt.Sprintf(`meta_key = '%s'`, o._adapter.SafeString(o.MetaKey)))
	}

	if o.IsMetaValueDirty == true {
		sets = append(sets, fmt.Sprintf(`meta_value = '%s'`, o._adapter.SafeString(o.MetaValue)))
	}

	frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'", o._table, strings.Join(sets, `,`))
	err := o._adapter.Execute(frmt)
	if err != nil {
		return err
	}
	return nil
}
func (o *UserMeta) Update() error {
	var sets []string

	if o.IsUserIdDirty == true {
		sets = append(sets, fmt.Sprintf(`user_id = '%d'`, o.UserId))
	}

	if o.IsMetaKeyDirty == true {
		sets = append(sets, fmt.Sprintf(`meta_key = '%s'`, o._adapter.SafeString(o.MetaKey)))
	}

	if o.IsMetaValueDirty == true {
		sets = append(sets, fmt.Sprintf(`meta_value = '%s'`, o._adapter.SafeString(o.MetaValue)))
	}

	frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'", o._table, strings.Join(sets, `,`))
	err := o._adapter.Execute(frmt)
	if err != nil {
		return err
	}
	return nil
}
func (o *UserMeta) Create() error {
	frmt := fmt.Sprintf("INSERT INTO %s (`user_id`, `meta_key`, `meta_value`) VALUES ('%d', '%s', '%s')", o._table, o.UserId, o.MetaKey, o.MetaValue)
	err := o._adapter.Execute(frmt)
	if err != nil {
		return o._adapter.Oops(fmt.Sprintf(`%s led to %s`, frmt, err))
	}
	o.UMetaId = o._adapter.LastInsertedId()
	o._new = false
	return nil
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
