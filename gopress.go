package gopress
import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql" // This is standard for this library.
    "strconv"
    "gopkg.in/yaml.v2"
    "regexp"
    "errors"
    "os"
    "io"
    "io/ioutil"
    "bufio"
    "log"
    "strings"
)



// LogFilter is an anonymous function that
// that receives the log tag and string and
// allows you to filter out extraneous lines
// when trying to find bugs.
type LogFilter func (string,string)string
// SafeStringFilter is the function that escapes
// possible SQL Injection code. 
type SafeStringFilter func(string)string
// Adapter is the main Database interface which helps
// to separate the DB from the Models. This is not
// 100% just yet, and may never be. Eventually the
// Adapter will probably receive some arguments and
// a value map and build the Query internally
type Adapter interface {
    Open(string,string,string,string) error
    Close()
    Query(string) ([]map[string]DBValue,error)
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
    SafeString(string)string
    NewDBValue() DBValue
}
// MysqlAdapter is the MySql implementation
type MysqlAdapter struct {
    // The host, localhost is valid here, or 127.0.0.1
    // if you use localhost, the system won't use TCP
    Host string `yaml:"host"`
    // The database username
    User string `yaml:"user"`
    // The database password
    Pass string `yaml:"pass"`
    // The database name
    Database string `yaml:"database"`
    // A prefix, if any - can be blank
    DBPrefix string `yaml:"prefix"`
    _infoLog *log.Logger
    _errorLog *log.Logger
    _debugLog *log.Logger
    _conn *sql.DB
    _lid int64
    _cnt int64
    _opened bool
    _logFilter LogFilter
    _safeStringFilter SafeStringFilter
}
// NewMysqlAdapter returns a pointer to MysqlAdapter
func NewMysqlAdapter(pre string) *MysqlAdapter {
    return &MysqlAdapter{DBPrefix: pre}
} 
// NewMysqlAdapterEx sets everything up based on your YAML config
// Args: fname is a string path to a YAML config file
// This function will attempt to Open the database
// defined in that file. Example file:
//     host: "localhost"
//     user: "dbuser"
//     pass: "dbuserpass"
//     database: "my_db"
//     prefix: "wp_"
func NewMysqlAdapterEx(fname string) (*MysqlAdapter,error) {
    a := NewMysqlAdapter(``)
    y,err := fileGetContents(fname)
    if err != nil {
        return nil,err
    }
    err = a.FromYAML(y)
    if err != nil {
        return nil,err
    }
    err = a.Open(a.Host,a.User,a.Pass,a.Database)
    if err != nil {
        return nil,err
    }
    a.SetLogs(ioutil.Discard)
    return a,nil
}
// SetLogFilter sets the LogFilter to a function. This is only
// useful if you are debugging, or you want to
// reformat the log data.
func (a *MysqlAdapter) SetLogFilter(f LogFilter) {
    a._logFilter = f
}
// SafeString Not implemented yet, but soon.
func (a *MysqlAdapter) SafeString(s string) string {
    return s
}
// SetInfoLog Sets the _infoLog to the io.Writer, use ioutil.Discard if you
// don't want this one at all.
func (a *MysqlAdapter) SetInfoLog(t io.Writer) {
    a._infoLog = log.New(t,`[INFO]:`,log.Ldate|log.Ltime|log.Lshortfile)
}
// SetErrorLog Sets the _errorLog to the io.Writer, use ioutil.Discard if you
// don't want this one at all.
func (a *MysqlAdapter) SetErrorLog(t io.Writer) {
    a._errorLog = log.New(t,`[ERROR]:`,log.Ldate|log.Ltime|log.Lshortfile)
}
// SetDebugLog Sets the _debugLog to the io.Writer, use ioutil.Discard if you
// don't want this one at all.
func (a *MysqlAdapter) SetDebugLog(t io.Writer) {
    a._debugLog = log.New(t,`[DEBUG]:`,log.Ldate|log.Ltime|log.Lshortfile)
}
// SetLogs Sets ALL logs to the io.Writer, use ioutil.Discard if you
// don't want this one at all.
func (a *MysqlAdapter) SetLogs(t io.Writer) {
    a.SetInfoLog(t)
    a.SetErrorLog(t)
    a.SetDebugLog(t)
}
// LogInfo Tags the string with INFO and puts it into _infoLog.
func (a *MysqlAdapter) LogInfo(s string) {
    if a._logFilter != nil {
        s = a._logFilter(`INFO`,s)
    }
    if s == "" {
        return
    }
    a._infoLog.Println(s)
}
// LogError Tags the string with ERROR and puts it into _errorLog.
func (a *MysqlAdapter) LogError(s error) {
    if a._logFilter != nil {
        ns := a._logFilter(`ERROR`,fmt.Sprintf(`%s`,s))
        if ns == `` {
            return
        }
        a._errorLog.Println(ns)
        return
    }
    a._errorLog.Println(s)
}
// LogDebug Tags the string with DEBUG and puts it into _debugLog.
func (a *MysqlAdapter) LogDebug(s string) {
    if a._logFilter != nil {
        s = a._logFilter(`DEBUG`,s)
    }
    if s == "" {
        return
    }
    a._debugLog.Println(s)
}
// NewDBValue Creates a new DBValue, mostly used internally, but
// you may wish to use it in special circumstances.
func (a *MysqlAdapter) NewDBValue() DBValue {
    return NewMysqlValue(a)
}
// DatabasePrefix Get the DatabasePrefix from the Adapter
func (a *MysqlAdapter) DatabasePrefix() string {
    return a.DBPrefix
}
// FromYAML Set the Adapter's members from a YAML file
func (a *MysqlAdapter) FromYAML(b []byte) error {
    return yaml.Unmarshal(b,a)
}
// Open Opens the database connection. Be sure to use 
// a.Close() as closing is NOT handled for you.
func (a *MysqlAdapter) Open(h,u,p,d string) error {
    if ( h != "localhost") {
        l := fmt.Sprintf("%s:%s@tcp(%s)/%s",u,p,h,d)
        tc, err := sql.Open("mysql",l)
        if err != nil {
            return a.Oops(fmt.Sprintf(`%s with %s`,err,l))
        }
        a._conn = tc
    } else {
        l := fmt.Sprintf("%s:%s@/%s",u,p,d)
        tc, err := sql.Open("mysql",l)
        if err != nil {
            return a.Oops(fmt.Sprintf(`%s with %s`,err,l))
        }
        a._conn = tc
    }
    err := a._conn.Ping()
    if err != nil {
        return err
    }
    a._opened = true
    return nil

}
// Close This should be called in your application with a defer a.Close() 
// or something similar. Closing is not automatic!
func (a *MysqlAdapter) Close() {
    a._conn.Close()
}
// Query The generay Query function, i.e. SQL that returns results, as
// opposed to an INSERT or UPDATE which uses Execute.
func (a *MysqlAdapter) Query(q string) ([]map[string]DBValue,error) {
    if a._opened != true {
        return nil,a.Oops(`you must first open the connection`)
    }
    results := new([]map[string]DBValue)
    a.LogInfo(q)
    rows, err := a._conn.Query(q)
    if err != nil {
        return nil,err
    }
    defer rows.Close()
    columns, err := rows.Columns()
    if err != nil {
        return nil, err
    }
    values := make([]sql.RawBytes, len(columns))
    scanArgs := make([]interface{},len(values))
    for i := range values {
        scanArgs[i] = &values[i]
    }
    for rows.Next() {
        err = rows.Scan(scanArgs...)
        if err != nil {
            return nil,err
        }
        res := make(map[string]DBValue)
        for i,col := range values {
            k := columns[i]
            res[k] = a.NewDBValue()
            res[k].SetInternalValue(k,string(col))
        }
        *results = append(*results,res)
    }
    return *results,nil
}
// Oops A function for catching errors generated by
// the library and funneling them to the log files
func (a *MysqlAdapter) Oops(s string) error {
    e := errors.New(s)
    a.LogError(e)
    return e
}
// Execute For UPDATE and INSERT calls, i.e. nothing that
// returns a result set.
func (a *MysqlAdapter) Execute(q string) error {
    if a._opened != true {
        return a.Oops(`you must first open the connection`)
    }
    tx, err := a._conn.Begin()
    if err != nil {
        return a.Oops(fmt.Sprintf(`could not Begin Transaction %s`,err))
    }
    defer tx.Rollback();
    stmt, err := tx.Prepare(q)
    if err != nil {
        return a.Oops(fmt.Sprintf(`could not Prepare Statement %s`,err))
    }
    defer stmt.Close()
    a.LogInfo(q)
    res,err := stmt.Exec()
    if err != nil {
        return a.Oops(fmt.Sprintf(`could not Exec stmt %s`,err))
    }
    a._lid,err = res.LastInsertId()
    a.LogInfo(fmt.Sprintf(`LastInsertedId is %d`,a._lid))
    if err != nil {
        return a.Oops(fmt.Sprintf(`could not get LastInsertId %s`,err))
    }
    a._cnt,err = res.RowsAffected()
    if err != nil {
        return a.Oops(fmt.Sprintf(`could not get RowsAffected %s`,err))
    }
    err = tx.Commit()
    if err != nil {
        return a.Oops(fmt.Sprintf(`could not Commit Transaction %s`,err))
    }
    return nil
}
// LastInsertedId Grab the last auto_incremented id
func (a *MysqlAdapter) LastInsertedId() int64 {
    return a._lid
}
// AffectedRows Grab the number of AffectedRows
func (a *MysqlAdapter) AffectedRows() int64 {
    return a._cnt
}
// DBValue Provides a tidy way to convert string
// values from the DB into go values
type DBValue interface {
    AsInt() (int,error)
    AsInt32() (int32,error)
    AsInt64() (int64,error)
    AsFloat32() (float32,error)
    AsFloat64() (float64,error)
    AsString() (string,error)
    AsDateTime() (*DateTime,error)
    SetInternalValue(string,string)
}
// MysqlValue Implements DBValue for MySQL, you'll generally
// not interact directly with this type, but it
// is there for special cases.
type MysqlValue struct {
    _v string
    _k string
    _adapter Adapter
}
// SetInternalValue Sets the internal value of the DBValue to the string
// provided. key isn't really used, but it may be.
func (v *MysqlValue) SetInternalValue(key,value string) {
    v._v = value
    v._k = key

}
// AsString Simply returns the internal string representation.
func (v *MysqlValue) AsString() (string,error) {
    return v._v,nil
}
// AsInt Attempts to convert the internal string to an Int
func (v *MysqlValue) AsInt() (int,error) {
    i,err := strconv.ParseInt(v._v,10,32)
    return int(i),err
}
// AsInt32 Tries to convert the internal string to an int32
func (v *MysqlValue) AsInt32() (int32,error) {
    i,err := strconv.ParseInt(v._v,10,32)
    return int32(i),err
}
// AsInt64 Tries to convert the internal string to an int64 (i.e. BIGINT)
func (v *MysqlValue) AsInt64() (int64,error) {
    i,err := strconv.ParseInt(v._v,10,64)
    return i,err
}
// AsFloat32 Tries to convert the internal string to a float32
func (v *MysqlValue) AsFloat32() (float32,error) {
    i,err := strconv.ParseFloat(v._v,32)
    if err != nil {
        return 0.0,err
    }
    return float32(i),err
}
// AsFloat64 Tries to convert the internal string to a float64
func (v *MysqlValue) AsFloat64() (float64,error) {
    i,err := strconv.ParseFloat(v._v,64)
    if err != nil {
        return 0.0,err
    }
    return i,err
}
// AsDateTime Tries to convert the string to a DateTime,
// parsing may fail.
func (v *MysqlValue) AsDateTime() (*DateTime,error) {
    dt := NewDateTime(v._adapter)
    err := dt.FromString(v._v)
    if err != nil {
        return &DateTime{}, err
    }
    return dt,nil
}
// NewMysqlValue A function for largely internal use, but
// basically in order to use a DBValue, it 
// needs to have its Adapter setup, this is
// because some values have Adapter specific
// issues. The implementing adapter may need
// to provide some information, or logging etc
func NewMysqlValue(a Adapter) *MysqlValue {
    return &MysqlValue{_adapter: a}
}
// DateTime A simple struct to represent DateTime fields
type DateTime struct {
    // The day as an int
    Day int
    // the month, as an int
    Month int
    // The year, as an int
    Year int
    // the hours, in 24 hour format
    Hours int
    // the minutes
    Minutes int
    // the seconds
    Seconds int
    _adapter Adapter
}
// FromString Converts a string like 0000-00-00 00:00:00 into a DateTime
func (d *DateTime) FromString(s string) error {
    es := s
    re := regexp.MustCompile("(?P<year>[\\d]{4})-(?P<month>[\\d]{2})-(?P<day>[\\d]{2}) (?P<hours>[\\d]{2}):(?P<minutes>[\\d]{2}):(?P<seconds>[\\d]{2})")
    n1 := re.SubexpNames()
    ir2 := re.FindAllStringSubmatch(es, -1)
    if len(ir2) == 0 {
        return d._adapter.Oops(fmt.Sprintf("found no data to capture in %s",es))
    }
    r2 := ir2[0]
    for i, n := range r2 {
        if n1[i] == "year" {
            _Year,err := strconv.ParseInt(n,10,32)
            d.Year = int(_Year)
            if err != nil {
                return d._adapter.Oops(fmt.Sprintf("failed to convert %d in %v received %s",n[i],es,err))
            }
        }
        if n1[i] == "month" {
            _Month,err := strconv.ParseInt(n,10,32)
            d.Month = int(_Month)
            if err != nil {
                return d._adapter.Oops(fmt.Sprintf("failed to convert %d in %v received %s",n[i],es,err))
            }
        }
        if n1[i] == "day" {
            _Day,err := strconv.ParseInt(n,10,32)
            d.Day = int(_Day)
            if err != nil {
                return d._adapter.Oops(fmt.Sprintf("failed to convert %d in %v received %s",n[i],es,err))
            }
        }
        if n1[i] == "hours" {
            _Hours,err := strconv.ParseInt(n,10,32)
            d.Hours = int(_Hours)
            if err != nil {
                return d._adapter.Oops(fmt.Sprintf("failed to convert %d in %v received %s",n[i],es,err))
            }
        }
        if n1[i] == "minutes" {
            _Minutes,err := strconv.ParseInt(n,10,32)
            d.Minutes = int(_Minutes)
            if err != nil {
                return d._adapter.Oops(fmt.Sprintf("failed to convert %d in %v received %s",n[i],es,err))
            }
        }
        if n1[i] == "seconds" {
            _Seconds,err := strconv.ParseInt(n,10,32)
            d.Seconds = int(_Seconds)
            if err != nil {
                return d._adapter.Oops(fmt.Sprintf("failed to convert %d in %v received %s",n[i],es,err))
            }
        }
    }
    return nil
}
// ToString For backwards compat... Never use this, use String() instead.
func (d *DateTime) ToString() string {
    return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",d.Year,d.Month,d.Day,d.Hours,d.Minutes,d.Seconds)
}
// String The Stringer for DateTime to avoid having to call ToString all the time.
func (d *DateTime) String() string {
    return d.ToString()
}
// NewDateTime Returns a basic DateTime value
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

// CommentMeta is a Object Relational Mapping to
// the database table that represents it. In this case it is
// commentmeta. The table name will be Sprintf'd to include
// the prefix you define in your YAML configuration for the
// Adapter.
type CommentMeta struct {
    _table string
    _adapter Adapter
    _pkey string // 0 The name of the primary key in this table
    _conds []string
    _new bool
    MetaId int64
    CommentId int64
    MetaKey string
    MetaValue string
	// Dirty markers for smart updates
    IsMetaIdDirty bool
    IsCommentIdDirty bool
    IsMetaKeyDirty bool
    IsMetaValueDirty bool
	// Relationships
}

// NewCommentMeta binds an Adapter to a new instance
// of CommentMeta and sets up the _table and primary keys
func NewCommentMeta(a Adapter) *CommentMeta {
    var o CommentMeta
    o._table = fmt.Sprintf("%scommentmeta",a.DatabasePrefix())
    o._adapter = a
    o._pkey = "meta_id"
    o._new = false
    return &o
}


// GetPrimaryKeyValue returns the value, usually int64 of
// the PrimaryKey
func (o *CommentMeta) GetPrimaryKeyValue() int64 {
    return o.MetaId
}
// GetPrimaryKeyName returns the DB field name
func (o *CommentMeta) GetPrimaryKeyName() string {
    return `meta_id`
}

// GetMetaId returns the value of 
// CommentMeta.MetaId
func (o *CommentMeta) GetMetaId() int64 {
    return o.MetaId
}
// SetMetaId sets and marks as dirty the value of
// CommentMeta.MetaId
func (o *CommentMeta) SetMetaId(arg int64) {
    o.MetaId = arg
    o.IsMetaIdDirty = true
}

// GetCommentId returns the value of 
// CommentMeta.CommentId
func (o *CommentMeta) GetCommentId() int64 {
    return o.CommentId
}
// SetCommentId sets and marks as dirty the value of
// CommentMeta.CommentId
func (o *CommentMeta) SetCommentId(arg int64) {
    o.CommentId = arg
    o.IsCommentIdDirty = true
}

// GetMetaKey returns the value of 
// CommentMeta.MetaKey
func (o *CommentMeta) GetMetaKey() string {
    return o.MetaKey
}
// SetMetaKey sets and marks as dirty the value of
// CommentMeta.MetaKey
func (o *CommentMeta) SetMetaKey(arg string) {
    o.MetaKey = arg
    o.IsMetaKeyDirty = true
}

// GetMetaValue returns the value of 
// CommentMeta.MetaValue
func (o *CommentMeta) GetMetaValue() string {
    return o.MetaValue
}
// SetMetaValue sets and marks as dirty the value of
// CommentMeta.MetaValue
func (o *CommentMeta) SetMetaValue(arg string) {
    o.MetaValue = arg
    o.IsMetaValueDirty = true
}

// Find dynamic finder for meta_id -> bool,error
// Generic and programatically generator finder for CommentMeta
// Note that Fine returns a bool if found, not err, in the case of
// a return of true, the instance data will be filled out.
// a call to find ALWAYS overwrites the model you call Find on
// i.e. receiver is a pointer. 
//```go
//      m := NewCommentMeta(a)
//      found,err := m.Find(23)
//      .. handle err
//      if found == false {
//          // handle found
//      }
//      ... do what you want with m here
//```
        func (o *CommentMeta) Find(_findByMetaId int64) (bool,error) {

    var _modelSlice []*CommentMeta
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "meta_id", _findByMetaId)
    results, err := o._adapter.Query(q)
    if err != nil {
        return false,o._adapter.Oops(fmt.Sprintf(`%s`,err))
    }
    
    for _,result := range results {
        ro := CommentMeta{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return false,o._adapter.Oops(fmt.Sprintf(`%s`,err))
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return false, o._adapter.Oops(`not found`)
    }
    o.FromCommentMeta(_modelSlice[0])
    return true,nil

}
// FindByCommentId dynamic finder for comment_id -> []*CommentMeta,error
// Generic and programatically generator finder for CommentMeta
//```go  
//    m := NewCommentMeta(a)
//    results,err := m.FindByCommentId(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of CommentMeta
//    }
//```  
        func (o *CommentMeta) FindByCommentId(_findByCommentId int64) ([]*CommentMeta,error) {

    var _modelSlice []*CommentMeta
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "comment_id", _findByCommentId)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := CommentMeta{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByMetaKey dynamic finder for meta_key -> []*CommentMeta,error
// Generic and programatically generator finder for CommentMeta
//```go  
//    m := NewCommentMeta(a)
//    results,err := m.FindByMetaKey(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of CommentMeta
//    }
//```  
        func (o *CommentMeta) FindByMetaKey(_findByMetaKey string) ([]*CommentMeta,error) {

    var _modelSlice []*CommentMeta
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "meta_key", _findByMetaKey)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := CommentMeta{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByMetaValue dynamic finder for meta_value -> []*CommentMeta,error
// Generic and programatically generator finder for CommentMeta
//```go  
//    m := NewCommentMeta(a)
//    results,err := m.FindByMetaValue(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of CommentMeta
//    }
//```  
        func (o *CommentMeta) FindByMetaValue(_findByMetaValue string) ([]*CommentMeta,error) {

    var _modelSlice []*CommentMeta
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "meta_value", _findByMetaValue)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := CommentMeta{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}

// FromDBValueMap Converts a DBValueMap returned from Adapter.Query to a CommentMeta
func (o *CommentMeta) FromDBValueMap(m map[string]DBValue) error {
	_MetaId,err := m["meta_id"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.MetaId = _MetaId
	_CommentId,err := m["comment_id"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.CommentId = _CommentId
	_MetaKey,err := m["meta_key"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.MetaKey = _MetaKey
	_MetaValue,err := m["meta_value"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.MetaValue = _MetaValue

 	return nil
}
// FromCommentMeta A kind of Clone function for CommentMeta
func (o *CommentMeta) FromCommentMeta(m *CommentMeta) {
	o.MetaId = m.MetaId
	o.CommentId = m.CommentId
	o.MetaKey = m.MetaKey
	o.MetaValue = m.MetaValue

}
// Reload A function to forcibly reload CommentMeta
func (o *CommentMeta) Reload() error {
    _,err := o.Find(o.GetPrimaryKeyValue())
    return err
}

// Save is a dynamic saver 'inherited' by all models
func (o *CommentMeta) Save() error {
    if o._new == true {
        return o.Create()
    }
    var sets []string
    
    if o.IsCommentIdDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_id = '%d'`,o.CommentId))
    }

    if o.IsMetaKeyDirty == true {
        sets = append(sets,fmt.Sprintf(`meta_key = '%s'`,o._adapter.SafeString(o.MetaKey)))
    }

    if o.IsMetaValueDirty == true {
        sets = append(sets,fmt.Sprintf(`meta_value = '%s'`,o._adapter.SafeString(o.MetaValue)))
    }

    frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'",o._table,strings.Join(sets,`,`),o._pkey, o.MetaId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return err
    }
    return nil
}
// Update is a dynamic updater, it considers whether or not
// a field is 'dirty' and needs to be updated. Will only work
// if you use the Getters and Setters
func (o *CommentMeta) Update() error {
    var sets []string
    
    if o.IsCommentIdDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_id = '%d'`,o.CommentId))
    }

    if o.IsMetaKeyDirty == true {
        sets = append(sets,fmt.Sprintf(`meta_key = '%s'`,o._adapter.SafeString(o.MetaKey)))
    }

    if o.IsMetaValueDirty == true {
        sets = append(sets,fmt.Sprintf(`meta_value = '%s'`,o._adapter.SafeString(o.MetaValue)))
    }

    frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'",o._table,strings.Join(sets,`,`),o._pkey, o.MetaId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return err
    }
    return nil
}
// Create inserts the model. Calling Save will call this function
// automatically for new models
func (o *CommentMeta) Create() error {
    frmt := fmt.Sprintf("INSERT INTO %s (`comment_id`, `meta_key`, `meta_value`) VALUES ('%d', '%s', '%s')",o._table,o.CommentId, o.MetaKey, o.MetaValue)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return o._adapter.Oops(fmt.Sprintf(`%s led to %s`,frmt,err))
    }
    o.MetaId = o._adapter.LastInsertedId()
    o._new = false
    return nil
}


// UpdateCommentId an immediate DB Query to update a single column, in this
// case comment_id
func (o *CommentMeta) UpdateCommentId(_updCommentId int64) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `comment_id` = '%d' WHERE `meta_id` = '%d'",o._table,_updCommentId,o.MetaId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.CommentId = _updCommentId
    return o._adapter.AffectedRows(),nil
}

// UpdateMetaKey an immediate DB Query to update a single column, in this
// case meta_key
func (o *CommentMeta) UpdateMetaKey(_updMetaKey string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `meta_key` = '%s' WHERE `meta_id` = '%d'",o._table,_updMetaKey,o.MetaId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.MetaKey = _updMetaKey
    return o._adapter.AffectedRows(),nil
}

// UpdateMetaValue an immediate DB Query to update a single column, in this
// case meta_value
func (o *CommentMeta) UpdateMetaValue(_updMetaValue string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `meta_value` = '%s' WHERE `meta_id` = '%d'",o._table,_updMetaValue,o.MetaId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.MetaValue = _updMetaValue
    return o._adapter.AffectedRows(),nil
}

// Comment is a Object Relational Mapping to
// the database table that represents it. In this case it is
// comments. The table name will be Sprintf'd to include
// the prefix you define in your YAML configuration for the
// Adapter.
type Comment struct {
    _table string
    _adapter Adapter
    _pkey string // 0 The name of the primary key in this table
    _conds []string
    _new bool
    CommentID int64
    CommentPostID int64
    CommentAuthor string
    CommentAuthorEmail string
    CommentAuthorUrl string
    CommentAuthorIP string
    CommentDate *DateTime
    CommentDateGmt *DateTime
    CommentContent string
    CommentKarma int
    CommentApproved string
    CommentAgent string
    CommentType string
    CommentParent int64
    UserId int64
	// Dirty markers for smart updates
    IsCommentIDDirty bool
    IsCommentPostIDDirty bool
    IsCommentAuthorDirty bool
    IsCommentAuthorEmailDirty bool
    IsCommentAuthorUrlDirty bool
    IsCommentAuthorIPDirty bool
    IsCommentDateDirty bool
    IsCommentDateGmtDirty bool
    IsCommentContentDirty bool
    IsCommentKarmaDirty bool
    IsCommentApprovedDirty bool
    IsCommentAgentDirty bool
    IsCommentTypeDirty bool
    IsCommentParentDirty bool
    IsUserIdDirty bool
	// Relationships
}

// NewComment binds an Adapter to a new instance
// of Comment and sets up the _table and primary keys
func NewComment(a Adapter) *Comment {
    var o Comment
    o._table = fmt.Sprintf("%scomments",a.DatabasePrefix())
    o._adapter = a
    o._pkey = "comment_ID"
    o._new = false
    return &o
}


// GetPrimaryKeyValue returns the value, usually int64 of
// the PrimaryKey
func (o *Comment) GetPrimaryKeyValue() int64 {
    return o.CommentID
}
// GetPrimaryKeyName returns the DB field name
func (o *Comment) GetPrimaryKeyName() string {
    return `comment_ID`
}

// GetCommentID returns the value of 
// Comment.CommentID
func (o *Comment) GetCommentID() int64 {
    return o.CommentID
}
// SetCommentID sets and marks as dirty the value of
// Comment.CommentID
func (o *Comment) SetCommentID(arg int64) {
    o.CommentID = arg
    o.IsCommentIDDirty = true
}

// GetCommentPostID returns the value of 
// Comment.CommentPostID
func (o *Comment) GetCommentPostID() int64 {
    return o.CommentPostID
}
// SetCommentPostID sets and marks as dirty the value of
// Comment.CommentPostID
func (o *Comment) SetCommentPostID(arg int64) {
    o.CommentPostID = arg
    o.IsCommentPostIDDirty = true
}

// GetCommentAuthor returns the value of 
// Comment.CommentAuthor
func (o *Comment) GetCommentAuthor() string {
    return o.CommentAuthor
}
// SetCommentAuthor sets and marks as dirty the value of
// Comment.CommentAuthor
func (o *Comment) SetCommentAuthor(arg string) {
    o.CommentAuthor = arg
    o.IsCommentAuthorDirty = true
}

// GetCommentAuthorEmail returns the value of 
// Comment.CommentAuthorEmail
func (o *Comment) GetCommentAuthorEmail() string {
    return o.CommentAuthorEmail
}
// SetCommentAuthorEmail sets and marks as dirty the value of
// Comment.CommentAuthorEmail
func (o *Comment) SetCommentAuthorEmail(arg string) {
    o.CommentAuthorEmail = arg
    o.IsCommentAuthorEmailDirty = true
}

// GetCommentAuthorUrl returns the value of 
// Comment.CommentAuthorUrl
func (o *Comment) GetCommentAuthorUrl() string {
    return o.CommentAuthorUrl
}
// SetCommentAuthorUrl sets and marks as dirty the value of
// Comment.CommentAuthorUrl
func (o *Comment) SetCommentAuthorUrl(arg string) {
    o.CommentAuthorUrl = arg
    o.IsCommentAuthorUrlDirty = true
}

// GetCommentAuthorIP returns the value of 
// Comment.CommentAuthorIP
func (o *Comment) GetCommentAuthorIP() string {
    return o.CommentAuthorIP
}
// SetCommentAuthorIP sets and marks as dirty the value of
// Comment.CommentAuthorIP
func (o *Comment) SetCommentAuthorIP(arg string) {
    o.CommentAuthorIP = arg
    o.IsCommentAuthorIPDirty = true
}

// GetCommentDate returns the value of 
// Comment.CommentDate
func (o *Comment) GetCommentDate() *DateTime {
    return o.CommentDate
}
// SetCommentDate sets and marks as dirty the value of
// Comment.CommentDate
func (o *Comment) SetCommentDate(arg *DateTime) {
    o.CommentDate = arg
    o.IsCommentDateDirty = true
}

// GetCommentDateGmt returns the value of 
// Comment.CommentDateGmt
func (o *Comment) GetCommentDateGmt() *DateTime {
    return o.CommentDateGmt
}
// SetCommentDateGmt sets and marks as dirty the value of
// Comment.CommentDateGmt
func (o *Comment) SetCommentDateGmt(arg *DateTime) {
    o.CommentDateGmt = arg
    o.IsCommentDateGmtDirty = true
}

// GetCommentContent returns the value of 
// Comment.CommentContent
func (o *Comment) GetCommentContent() string {
    return o.CommentContent
}
// SetCommentContent sets and marks as dirty the value of
// Comment.CommentContent
func (o *Comment) SetCommentContent(arg string) {
    o.CommentContent = arg
    o.IsCommentContentDirty = true
}

// GetCommentKarma returns the value of 
// Comment.CommentKarma
func (o *Comment) GetCommentKarma() int {
    return o.CommentKarma
}
// SetCommentKarma sets and marks as dirty the value of
// Comment.CommentKarma
func (o *Comment) SetCommentKarma(arg int) {
    o.CommentKarma = arg
    o.IsCommentKarmaDirty = true
}

// GetCommentApproved returns the value of 
// Comment.CommentApproved
func (o *Comment) GetCommentApproved() string {
    return o.CommentApproved
}
// SetCommentApproved sets and marks as dirty the value of
// Comment.CommentApproved
func (o *Comment) SetCommentApproved(arg string) {
    o.CommentApproved = arg
    o.IsCommentApprovedDirty = true
}

// GetCommentAgent returns the value of 
// Comment.CommentAgent
func (o *Comment) GetCommentAgent() string {
    return o.CommentAgent
}
// SetCommentAgent sets and marks as dirty the value of
// Comment.CommentAgent
func (o *Comment) SetCommentAgent(arg string) {
    o.CommentAgent = arg
    o.IsCommentAgentDirty = true
}

// GetCommentType returns the value of 
// Comment.CommentType
func (o *Comment) GetCommentType() string {
    return o.CommentType
}
// SetCommentType sets and marks as dirty the value of
// Comment.CommentType
func (o *Comment) SetCommentType(arg string) {
    o.CommentType = arg
    o.IsCommentTypeDirty = true
}

// GetCommentParent returns the value of 
// Comment.CommentParent
func (o *Comment) GetCommentParent() int64 {
    return o.CommentParent
}
// SetCommentParent sets and marks as dirty the value of
// Comment.CommentParent
func (o *Comment) SetCommentParent(arg int64) {
    o.CommentParent = arg
    o.IsCommentParentDirty = true
}

// GetUserId returns the value of 
// Comment.UserId
func (o *Comment) GetUserId() int64 {
    return o.UserId
}
// SetUserId sets and marks as dirty the value of
// Comment.UserId
func (o *Comment) SetUserId(arg int64) {
    o.UserId = arg
    o.IsUserIdDirty = true
}

// Find dynamic finder for comment_ID -> bool,error
// Generic and programatically generator finder for Comment
// Note that Fine returns a bool if found, not err, in the case of
// a return of true, the instance data will be filled out.
// a call to find ALWAYS overwrites the model you call Find on
// i.e. receiver is a pointer. 
//```go
//      m := NewComment(a)
//      found,err := m.Find(23)
//      .. handle err
//      if found == false {
//          // handle found
//      }
//      ... do what you want with m here
//```
        func (o *Comment) Find(_findByCommentID int64) (bool,error) {

    var _modelSlice []*Comment
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "comment_ID", _findByCommentID)
    results, err := o._adapter.Query(q)
    if err != nil {
        return false,o._adapter.Oops(fmt.Sprintf(`%s`,err))
    }
    
    for _,result := range results {
        ro := Comment{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return false,o._adapter.Oops(fmt.Sprintf(`%s`,err))
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return false, o._adapter.Oops(`not found`)
    }
    o.FromComment(_modelSlice[0])
    return true,nil

}
// FindByCommentPostID dynamic finder for comment_post_ID -> []*Comment,error
// Generic and programatically generator finder for Comment
//```go  
//    m := NewComment(a)
//    results,err := m.FindByCommentPostID(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Comment
//    }
//```  
        func (o *Comment) FindByCommentPostID(_findByCommentPostID int64) ([]*Comment,error) {

    var _modelSlice []*Comment
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "comment_post_ID", _findByCommentPostID)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Comment{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByCommentAuthor dynamic finder for comment_author -> []*Comment,error
// Generic and programatically generator finder for Comment
//```go  
//    m := NewComment(a)
//    results,err := m.FindByCommentAuthor(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Comment
//    }
//```  
        func (o *Comment) FindByCommentAuthor(_findByCommentAuthor string) ([]*Comment,error) {

    var _modelSlice []*Comment
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "comment_author", _findByCommentAuthor)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Comment{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByCommentAuthorEmail dynamic finder for comment_author_email -> []*Comment,error
// Generic and programatically generator finder for Comment
//```go  
//    m := NewComment(a)
//    results,err := m.FindByCommentAuthorEmail(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Comment
//    }
//```  
        func (o *Comment) FindByCommentAuthorEmail(_findByCommentAuthorEmail string) ([]*Comment,error) {

    var _modelSlice []*Comment
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "comment_author_email", _findByCommentAuthorEmail)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Comment{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByCommentAuthorUrl dynamic finder for comment_author_url -> []*Comment,error
// Generic and programatically generator finder for Comment
//```go  
//    m := NewComment(a)
//    results,err := m.FindByCommentAuthorUrl(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Comment
//    }
//```  
        func (o *Comment) FindByCommentAuthorUrl(_findByCommentAuthorUrl string) ([]*Comment,error) {

    var _modelSlice []*Comment
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "comment_author_url", _findByCommentAuthorUrl)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Comment{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByCommentAuthorIP dynamic finder for comment_author_IP -> []*Comment,error
// Generic and programatically generator finder for Comment
//```go  
//    m := NewComment(a)
//    results,err := m.FindByCommentAuthorIP(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Comment
//    }
//```  
        func (o *Comment) FindByCommentAuthorIP(_findByCommentAuthorIP string) ([]*Comment,error) {

    var _modelSlice []*Comment
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "comment_author_IP", _findByCommentAuthorIP)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Comment{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByCommentDate dynamic finder for comment_date -> []*Comment,error
// Generic and programatically generator finder for Comment
//```go  
//    m := NewComment(a)
//    results,err := m.FindByCommentDate(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Comment
//    }
//```  
        func (o *Comment) FindByCommentDate(_findByCommentDate *DateTime) ([]*Comment,error) {

    var _modelSlice []*Comment
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "comment_date", _findByCommentDate)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Comment{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByCommentDateGmt dynamic finder for comment_date_gmt -> []*Comment,error
// Generic and programatically generator finder for Comment
//```go  
//    m := NewComment(a)
//    results,err := m.FindByCommentDateGmt(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Comment
//    }
//```  
        func (o *Comment) FindByCommentDateGmt(_findByCommentDateGmt *DateTime) ([]*Comment,error) {

    var _modelSlice []*Comment
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "comment_date_gmt", _findByCommentDateGmt)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Comment{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByCommentContent dynamic finder for comment_content -> []*Comment,error
// Generic and programatically generator finder for Comment
//```go  
//    m := NewComment(a)
//    results,err := m.FindByCommentContent(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Comment
//    }
//```  
        func (o *Comment) FindByCommentContent(_findByCommentContent string) ([]*Comment,error) {

    var _modelSlice []*Comment
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "comment_content", _findByCommentContent)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Comment{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByCommentKarma dynamic finder for comment_karma -> []*Comment,error
// Generic and programatically generator finder for Comment
//```go  
//    m := NewComment(a)
//    results,err := m.FindByCommentKarma(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Comment
//    }
//```  
        func (o *Comment) FindByCommentKarma(_findByCommentKarma int) ([]*Comment,error) {

    var _modelSlice []*Comment
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "comment_karma", _findByCommentKarma)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Comment{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByCommentApproved dynamic finder for comment_approved -> []*Comment,error
// Generic and programatically generator finder for Comment
//```go  
//    m := NewComment(a)
//    results,err := m.FindByCommentApproved(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Comment
//    }
//```  
        func (o *Comment) FindByCommentApproved(_findByCommentApproved string) ([]*Comment,error) {

    var _modelSlice []*Comment
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "comment_approved", _findByCommentApproved)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Comment{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByCommentAgent dynamic finder for comment_agent -> []*Comment,error
// Generic and programatically generator finder for Comment
//```go  
//    m := NewComment(a)
//    results,err := m.FindByCommentAgent(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Comment
//    }
//```  
        func (o *Comment) FindByCommentAgent(_findByCommentAgent string) ([]*Comment,error) {

    var _modelSlice []*Comment
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "comment_agent", _findByCommentAgent)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Comment{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByCommentType dynamic finder for comment_type -> []*Comment,error
// Generic and programatically generator finder for Comment
//```go  
//    m := NewComment(a)
//    results,err := m.FindByCommentType(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Comment
//    }
//```  
        func (o *Comment) FindByCommentType(_findByCommentType string) ([]*Comment,error) {

    var _modelSlice []*Comment
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "comment_type", _findByCommentType)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Comment{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByCommentParent dynamic finder for comment_parent -> []*Comment,error
// Generic and programatically generator finder for Comment
//```go  
//    m := NewComment(a)
//    results,err := m.FindByCommentParent(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Comment
//    }
//```  
        func (o *Comment) FindByCommentParent(_findByCommentParent int64) ([]*Comment,error) {

    var _modelSlice []*Comment
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "comment_parent", _findByCommentParent)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Comment{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByUserId dynamic finder for user_id -> []*Comment,error
// Generic and programatically generator finder for Comment
//```go  
//    m := NewComment(a)
//    results,err := m.FindByUserId(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Comment
//    }
//```  
        func (o *Comment) FindByUserId(_findByUserId int64) ([]*Comment,error) {

    var _modelSlice []*Comment
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "user_id", _findByUserId)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Comment{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}

// FromDBValueMap Converts a DBValueMap returned from Adapter.Query to a Comment
func (o *Comment) FromDBValueMap(m map[string]DBValue) error {
	_CommentID,err := m["comment_ID"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.CommentID = _CommentID
	_CommentPostID,err := m["comment_post_ID"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.CommentPostID = _CommentPostID
	_CommentAuthor,err := m["comment_author"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.CommentAuthor = _CommentAuthor
	_CommentAuthorEmail,err := m["comment_author_email"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.CommentAuthorEmail = _CommentAuthorEmail
	_CommentAuthorUrl,err := m["comment_author_url"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.CommentAuthorUrl = _CommentAuthorUrl
	_CommentAuthorIP,err := m["comment_author_IP"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.CommentAuthorIP = _CommentAuthorIP
	_CommentDate,err := m["comment_date"].AsDateTime()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.CommentDate = _CommentDate
	_CommentDateGmt,err := m["comment_date_gmt"].AsDateTime()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.CommentDateGmt = _CommentDateGmt
	_CommentContent,err := m["comment_content"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.CommentContent = _CommentContent
	_CommentKarma,err := m["comment_karma"].AsInt()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.CommentKarma = _CommentKarma
	_CommentApproved,err := m["comment_approved"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.CommentApproved = _CommentApproved
	_CommentAgent,err := m["comment_agent"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.CommentAgent = _CommentAgent
	_CommentType,err := m["comment_type"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.CommentType = _CommentType
	_CommentParent,err := m["comment_parent"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.CommentParent = _CommentParent
	_UserId,err := m["user_id"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.UserId = _UserId

 	return nil
}
// FromComment A kind of Clone function for Comment
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
// Reload A function to forcibly reload Comment
func (o *Comment) Reload() error {
    _,err := o.Find(o.GetPrimaryKeyValue())
    return err
}

// Save is a dynamic saver 'inherited' by all models
func (o *Comment) Save() error {
    if o._new == true {
        return o.Create()
    }
    var sets []string
    
    if o.IsCommentPostIDDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_post_ID = '%d'`,o.CommentPostID))
    }

    if o.IsCommentAuthorDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_author = '%s'`,o._adapter.SafeString(o.CommentAuthor)))
    }

    if o.IsCommentAuthorEmailDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_author_email = '%s'`,o._adapter.SafeString(o.CommentAuthorEmail)))
    }

    if o.IsCommentAuthorUrlDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_author_url = '%s'`,o._adapter.SafeString(o.CommentAuthorUrl)))
    }

    if o.IsCommentAuthorIPDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_author_IP = '%s'`,o._adapter.SafeString(o.CommentAuthorIP)))
    }

    if o.IsCommentDateDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_date = '%s'`,o.CommentDate))
    }

    if o.IsCommentDateGmtDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_date_gmt = '%s'`,o.CommentDateGmt))
    }

    if o.IsCommentContentDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_content = '%s'`,o._adapter.SafeString(o.CommentContent)))
    }

    if o.IsCommentKarmaDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_karma = '%d'`,o.CommentKarma))
    }

    if o.IsCommentApprovedDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_approved = '%s'`,o._adapter.SafeString(o.CommentApproved)))
    }

    if o.IsCommentAgentDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_agent = '%s'`,o._adapter.SafeString(o.CommentAgent)))
    }

    if o.IsCommentTypeDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_type = '%s'`,o._adapter.SafeString(o.CommentType)))
    }

    if o.IsCommentParentDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_parent = '%d'`,o.CommentParent))
    }

    if o.IsUserIdDirty == true {
        sets = append(sets,fmt.Sprintf(`user_id = '%d'`,o.UserId))
    }

    frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'",o._table,strings.Join(sets,`,`),o._pkey, o.CommentID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return err
    }
    return nil
}
// Update is a dynamic updater, it considers whether or not
// a field is 'dirty' and needs to be updated. Will only work
// if you use the Getters and Setters
func (o *Comment) Update() error {
    var sets []string
    
    if o.IsCommentPostIDDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_post_ID = '%d'`,o.CommentPostID))
    }

    if o.IsCommentAuthorDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_author = '%s'`,o._adapter.SafeString(o.CommentAuthor)))
    }

    if o.IsCommentAuthorEmailDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_author_email = '%s'`,o._adapter.SafeString(o.CommentAuthorEmail)))
    }

    if o.IsCommentAuthorUrlDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_author_url = '%s'`,o._adapter.SafeString(o.CommentAuthorUrl)))
    }

    if o.IsCommentAuthorIPDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_author_IP = '%s'`,o._adapter.SafeString(o.CommentAuthorIP)))
    }

    if o.IsCommentDateDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_date = '%s'`,o.CommentDate))
    }

    if o.IsCommentDateGmtDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_date_gmt = '%s'`,o.CommentDateGmt))
    }

    if o.IsCommentContentDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_content = '%s'`,o._adapter.SafeString(o.CommentContent)))
    }

    if o.IsCommentKarmaDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_karma = '%d'`,o.CommentKarma))
    }

    if o.IsCommentApprovedDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_approved = '%s'`,o._adapter.SafeString(o.CommentApproved)))
    }

    if o.IsCommentAgentDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_agent = '%s'`,o._adapter.SafeString(o.CommentAgent)))
    }

    if o.IsCommentTypeDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_type = '%s'`,o._adapter.SafeString(o.CommentType)))
    }

    if o.IsCommentParentDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_parent = '%d'`,o.CommentParent))
    }

    if o.IsUserIdDirty == true {
        sets = append(sets,fmt.Sprintf(`user_id = '%d'`,o.UserId))
    }

    frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'",o._table,strings.Join(sets,`,`),o._pkey, o.CommentID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return err
    }
    return nil
}
// Create inserts the model. Calling Save will call this function
// automatically for new models
func (o *Comment) Create() error {
    frmt := fmt.Sprintf("INSERT INTO %s (`comment_post_ID`, `comment_author`, `comment_author_email`, `comment_author_url`, `comment_author_IP`, `comment_date`, `comment_date_gmt`, `comment_content`, `comment_karma`, `comment_approved`, `comment_agent`, `comment_type`, `comment_parent`, `user_id`) VALUES ('%d', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%d', '%s', '%s', '%s', '%d', '%d')",o._table,o.CommentPostID, o.CommentAuthor, o.CommentAuthorEmail, o.CommentAuthorUrl, o.CommentAuthorIP, o.CommentDate.ToString(), o.CommentDateGmt.ToString(), o.CommentContent, o.CommentKarma, o.CommentApproved, o.CommentAgent, o.CommentType, o.CommentParent, o.UserId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return o._adapter.Oops(fmt.Sprintf(`%s led to %s`,frmt,err))
    }
    o.CommentID = o._adapter.LastInsertedId()
    o._new = false
    return nil
}


// UpdateCommentPostID an immediate DB Query to update a single column, in this
// case comment_post_ID
func (o *Comment) UpdateCommentPostID(_updCommentPostID int64) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `comment_post_ID` = '%d' WHERE `comment_ID` = '%d'",o._table,_updCommentPostID,o.CommentID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.CommentPostID = _updCommentPostID
    return o._adapter.AffectedRows(),nil
}

// UpdateCommentAuthor an immediate DB Query to update a single column, in this
// case comment_author
func (o *Comment) UpdateCommentAuthor(_updCommentAuthor string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `comment_author` = '%s' WHERE `comment_ID` = '%d'",o._table,_updCommentAuthor,o.CommentID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.CommentAuthor = _updCommentAuthor
    return o._adapter.AffectedRows(),nil
}

// UpdateCommentAuthorEmail an immediate DB Query to update a single column, in this
// case comment_author_email
func (o *Comment) UpdateCommentAuthorEmail(_updCommentAuthorEmail string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `comment_author_email` = '%s' WHERE `comment_ID` = '%d'",o._table,_updCommentAuthorEmail,o.CommentID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.CommentAuthorEmail = _updCommentAuthorEmail
    return o._adapter.AffectedRows(),nil
}

// UpdateCommentAuthorUrl an immediate DB Query to update a single column, in this
// case comment_author_url
func (o *Comment) UpdateCommentAuthorUrl(_updCommentAuthorUrl string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `comment_author_url` = '%s' WHERE `comment_ID` = '%d'",o._table,_updCommentAuthorUrl,o.CommentID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.CommentAuthorUrl = _updCommentAuthorUrl
    return o._adapter.AffectedRows(),nil
}

// UpdateCommentAuthorIP an immediate DB Query to update a single column, in this
// case comment_author_IP
func (o *Comment) UpdateCommentAuthorIP(_updCommentAuthorIP string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `comment_author_IP` = '%s' WHERE `comment_ID` = '%d'",o._table,_updCommentAuthorIP,o.CommentID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.CommentAuthorIP = _updCommentAuthorIP
    return o._adapter.AffectedRows(),nil
}

// UpdateCommentDate an immediate DB Query to update a single column, in this
// case comment_date
func (o *Comment) UpdateCommentDate(_updCommentDate *DateTime) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `comment_date` = '%s' WHERE `comment_ID` = '%d'",o._table,_updCommentDate,o.CommentID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.CommentDate = _updCommentDate
    return o._adapter.AffectedRows(),nil
}

// UpdateCommentDateGmt an immediate DB Query to update a single column, in this
// case comment_date_gmt
func (o *Comment) UpdateCommentDateGmt(_updCommentDateGmt *DateTime) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `comment_date_gmt` = '%s' WHERE `comment_ID` = '%d'",o._table,_updCommentDateGmt,o.CommentID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.CommentDateGmt = _updCommentDateGmt
    return o._adapter.AffectedRows(),nil
}

// UpdateCommentContent an immediate DB Query to update a single column, in this
// case comment_content
func (o *Comment) UpdateCommentContent(_updCommentContent string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `comment_content` = '%s' WHERE `comment_ID` = '%d'",o._table,_updCommentContent,o.CommentID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.CommentContent = _updCommentContent
    return o._adapter.AffectedRows(),nil
}

// UpdateCommentKarma an immediate DB Query to update a single column, in this
// case comment_karma
func (o *Comment) UpdateCommentKarma(_updCommentKarma int) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `comment_karma` = '%d' WHERE `comment_ID` = '%d'",o._table,_updCommentKarma,o.CommentID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.CommentKarma = _updCommentKarma
    return o._adapter.AffectedRows(),nil
}

// UpdateCommentApproved an immediate DB Query to update a single column, in this
// case comment_approved
func (o *Comment) UpdateCommentApproved(_updCommentApproved string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `comment_approved` = '%s' WHERE `comment_ID` = '%d'",o._table,_updCommentApproved,o.CommentID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.CommentApproved = _updCommentApproved
    return o._adapter.AffectedRows(),nil
}

// UpdateCommentAgent an immediate DB Query to update a single column, in this
// case comment_agent
func (o *Comment) UpdateCommentAgent(_updCommentAgent string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `comment_agent` = '%s' WHERE `comment_ID` = '%d'",o._table,_updCommentAgent,o.CommentID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.CommentAgent = _updCommentAgent
    return o._adapter.AffectedRows(),nil
}

// UpdateCommentType an immediate DB Query to update a single column, in this
// case comment_type
func (o *Comment) UpdateCommentType(_updCommentType string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `comment_type` = '%s' WHERE `comment_ID` = '%d'",o._table,_updCommentType,o.CommentID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.CommentType = _updCommentType
    return o._adapter.AffectedRows(),nil
}

// UpdateCommentParent an immediate DB Query to update a single column, in this
// case comment_parent
func (o *Comment) UpdateCommentParent(_updCommentParent int64) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `comment_parent` = '%d' WHERE `comment_ID` = '%d'",o._table,_updCommentParent,o.CommentID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.CommentParent = _updCommentParent
    return o._adapter.AffectedRows(),nil
}

// UpdateUserId an immediate DB Query to update a single column, in this
// case user_id
func (o *Comment) UpdateUserId(_updUserId int64) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `user_id` = '%d' WHERE `comment_ID` = '%d'",o._table,_updUserId,o.CommentID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.UserId = _updUserId
    return o._adapter.AffectedRows(),nil
}

// Link is a Object Relational Mapping to
// the database table that represents it. In this case it is
// links. The table name will be Sprintf'd to include
// the prefix you define in your YAML configuration for the
// Adapter.
type Link struct {
    _table string
    _adapter Adapter
    _pkey string // 0 The name of the primary key in this table
    _conds []string
    _new bool
    LinkId int64
    LinkUrl string
    LinkName string
    LinkImage string
    LinkTarget string
    LinkDescription string
    LinkVisible string
    LinkOwner int64
    LinkRating int
    LinkUpdated *DateTime
    LinkRel string
    LinkNotes string
    LinkRss string
	// Dirty markers for smart updates
    IsLinkIdDirty bool
    IsLinkUrlDirty bool
    IsLinkNameDirty bool
    IsLinkImageDirty bool
    IsLinkTargetDirty bool
    IsLinkDescriptionDirty bool
    IsLinkVisibleDirty bool
    IsLinkOwnerDirty bool
    IsLinkRatingDirty bool
    IsLinkUpdatedDirty bool
    IsLinkRelDirty bool
    IsLinkNotesDirty bool
    IsLinkRssDirty bool
	// Relationships
}

// NewLink binds an Adapter to a new instance
// of Link and sets up the _table and primary keys
func NewLink(a Adapter) *Link {
    var o Link
    o._table = fmt.Sprintf("%slinks",a.DatabasePrefix())
    o._adapter = a
    o._pkey = "link_id"
    o._new = false
    return &o
}


// GetPrimaryKeyValue returns the value, usually int64 of
// the PrimaryKey
func (o *Link) GetPrimaryKeyValue() int64 {
    return o.LinkId
}
// GetPrimaryKeyName returns the DB field name
func (o *Link) GetPrimaryKeyName() string {
    return `link_id`
}

// GetLinkId returns the value of 
// Link.LinkId
func (o *Link) GetLinkId() int64 {
    return o.LinkId
}
// SetLinkId sets and marks as dirty the value of
// Link.LinkId
func (o *Link) SetLinkId(arg int64) {
    o.LinkId = arg
    o.IsLinkIdDirty = true
}

// GetLinkUrl returns the value of 
// Link.LinkUrl
func (o *Link) GetLinkUrl() string {
    return o.LinkUrl
}
// SetLinkUrl sets and marks as dirty the value of
// Link.LinkUrl
func (o *Link) SetLinkUrl(arg string) {
    o.LinkUrl = arg
    o.IsLinkUrlDirty = true
}

// GetLinkName returns the value of 
// Link.LinkName
func (o *Link) GetLinkName() string {
    return o.LinkName
}
// SetLinkName sets and marks as dirty the value of
// Link.LinkName
func (o *Link) SetLinkName(arg string) {
    o.LinkName = arg
    o.IsLinkNameDirty = true
}

// GetLinkImage returns the value of 
// Link.LinkImage
func (o *Link) GetLinkImage() string {
    return o.LinkImage
}
// SetLinkImage sets and marks as dirty the value of
// Link.LinkImage
func (o *Link) SetLinkImage(arg string) {
    o.LinkImage = arg
    o.IsLinkImageDirty = true
}

// GetLinkTarget returns the value of 
// Link.LinkTarget
func (o *Link) GetLinkTarget() string {
    return o.LinkTarget
}
// SetLinkTarget sets and marks as dirty the value of
// Link.LinkTarget
func (o *Link) SetLinkTarget(arg string) {
    o.LinkTarget = arg
    o.IsLinkTargetDirty = true
}

// GetLinkDescription returns the value of 
// Link.LinkDescription
func (o *Link) GetLinkDescription() string {
    return o.LinkDescription
}
// SetLinkDescription sets and marks as dirty the value of
// Link.LinkDescription
func (o *Link) SetLinkDescription(arg string) {
    o.LinkDescription = arg
    o.IsLinkDescriptionDirty = true
}

// GetLinkVisible returns the value of 
// Link.LinkVisible
func (o *Link) GetLinkVisible() string {
    return o.LinkVisible
}
// SetLinkVisible sets and marks as dirty the value of
// Link.LinkVisible
func (o *Link) SetLinkVisible(arg string) {
    o.LinkVisible = arg
    o.IsLinkVisibleDirty = true
}

// GetLinkOwner returns the value of 
// Link.LinkOwner
func (o *Link) GetLinkOwner() int64 {
    return o.LinkOwner
}
// SetLinkOwner sets and marks as dirty the value of
// Link.LinkOwner
func (o *Link) SetLinkOwner(arg int64) {
    o.LinkOwner = arg
    o.IsLinkOwnerDirty = true
}

// GetLinkRating returns the value of 
// Link.LinkRating
func (o *Link) GetLinkRating() int {
    return o.LinkRating
}
// SetLinkRating sets and marks as dirty the value of
// Link.LinkRating
func (o *Link) SetLinkRating(arg int) {
    o.LinkRating = arg
    o.IsLinkRatingDirty = true
}

// GetLinkUpdated returns the value of 
// Link.LinkUpdated
func (o *Link) GetLinkUpdated() *DateTime {
    return o.LinkUpdated
}
// SetLinkUpdated sets and marks as dirty the value of
// Link.LinkUpdated
func (o *Link) SetLinkUpdated(arg *DateTime) {
    o.LinkUpdated = arg
    o.IsLinkUpdatedDirty = true
}

// GetLinkRel returns the value of 
// Link.LinkRel
func (o *Link) GetLinkRel() string {
    return o.LinkRel
}
// SetLinkRel sets and marks as dirty the value of
// Link.LinkRel
func (o *Link) SetLinkRel(arg string) {
    o.LinkRel = arg
    o.IsLinkRelDirty = true
}

// GetLinkNotes returns the value of 
// Link.LinkNotes
func (o *Link) GetLinkNotes() string {
    return o.LinkNotes
}
// SetLinkNotes sets and marks as dirty the value of
// Link.LinkNotes
func (o *Link) SetLinkNotes(arg string) {
    o.LinkNotes = arg
    o.IsLinkNotesDirty = true
}

// GetLinkRss returns the value of 
// Link.LinkRss
func (o *Link) GetLinkRss() string {
    return o.LinkRss
}
// SetLinkRss sets and marks as dirty the value of
// Link.LinkRss
func (o *Link) SetLinkRss(arg string) {
    o.LinkRss = arg
    o.IsLinkRssDirty = true
}

// Find dynamic finder for link_id -> bool,error
// Generic and programatically generator finder for Link
// Note that Fine returns a bool if found, not err, in the case of
// a return of true, the instance data will be filled out.
// a call to find ALWAYS overwrites the model you call Find on
// i.e. receiver is a pointer. 
//```go
//      m := NewLink(a)
//      found,err := m.Find(23)
//      .. handle err
//      if found == false {
//          // handle found
//      }
//      ... do what you want with m here
//```
        func (o *Link) Find(_findByLinkId int64) (bool,error) {

    var _modelSlice []*Link
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "link_id", _findByLinkId)
    results, err := o._adapter.Query(q)
    if err != nil {
        return false,o._adapter.Oops(fmt.Sprintf(`%s`,err))
    }
    
    for _,result := range results {
        ro := Link{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return false,o._adapter.Oops(fmt.Sprintf(`%s`,err))
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return false, o._adapter.Oops(`not found`)
    }
    o.FromLink(_modelSlice[0])
    return true,nil

}
// FindByLinkUrl dynamic finder for link_url -> []*Link,error
// Generic and programatically generator finder for Link
//```go  
//    m := NewLink(a)
//    results,err := m.FindByLinkUrl(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Link
//    }
//```  
        func (o *Link) FindByLinkUrl(_findByLinkUrl string) ([]*Link,error) {

    var _modelSlice []*Link
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "link_url", _findByLinkUrl)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Link{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByLinkName dynamic finder for link_name -> []*Link,error
// Generic and programatically generator finder for Link
//```go  
//    m := NewLink(a)
//    results,err := m.FindByLinkName(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Link
//    }
//```  
        func (o *Link) FindByLinkName(_findByLinkName string) ([]*Link,error) {

    var _modelSlice []*Link
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "link_name", _findByLinkName)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Link{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByLinkImage dynamic finder for link_image -> []*Link,error
// Generic and programatically generator finder for Link
//```go  
//    m := NewLink(a)
//    results,err := m.FindByLinkImage(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Link
//    }
//```  
        func (o *Link) FindByLinkImage(_findByLinkImage string) ([]*Link,error) {

    var _modelSlice []*Link
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "link_image", _findByLinkImage)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Link{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByLinkTarget dynamic finder for link_target -> []*Link,error
// Generic and programatically generator finder for Link
//```go  
//    m := NewLink(a)
//    results,err := m.FindByLinkTarget(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Link
//    }
//```  
        func (o *Link) FindByLinkTarget(_findByLinkTarget string) ([]*Link,error) {

    var _modelSlice []*Link
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "link_target", _findByLinkTarget)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Link{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByLinkDescription dynamic finder for link_description -> []*Link,error
// Generic and programatically generator finder for Link
//```go  
//    m := NewLink(a)
//    results,err := m.FindByLinkDescription(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Link
//    }
//```  
        func (o *Link) FindByLinkDescription(_findByLinkDescription string) ([]*Link,error) {

    var _modelSlice []*Link
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "link_description", _findByLinkDescription)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Link{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByLinkVisible dynamic finder for link_visible -> []*Link,error
// Generic and programatically generator finder for Link
//```go  
//    m := NewLink(a)
//    results,err := m.FindByLinkVisible(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Link
//    }
//```  
        func (o *Link) FindByLinkVisible(_findByLinkVisible string) ([]*Link,error) {

    var _modelSlice []*Link
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "link_visible", _findByLinkVisible)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Link{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByLinkOwner dynamic finder for link_owner -> []*Link,error
// Generic and programatically generator finder for Link
//```go  
//    m := NewLink(a)
//    results,err := m.FindByLinkOwner(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Link
//    }
//```  
        func (o *Link) FindByLinkOwner(_findByLinkOwner int64) ([]*Link,error) {

    var _modelSlice []*Link
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "link_owner", _findByLinkOwner)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Link{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByLinkRating dynamic finder for link_rating -> []*Link,error
// Generic and programatically generator finder for Link
//```go  
//    m := NewLink(a)
//    results,err := m.FindByLinkRating(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Link
//    }
//```  
        func (o *Link) FindByLinkRating(_findByLinkRating int) ([]*Link,error) {

    var _modelSlice []*Link
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "link_rating", _findByLinkRating)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Link{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByLinkUpdated dynamic finder for link_updated -> []*Link,error
// Generic and programatically generator finder for Link
//```go  
//    m := NewLink(a)
//    results,err := m.FindByLinkUpdated(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Link
//    }
//```  
        func (o *Link) FindByLinkUpdated(_findByLinkUpdated *DateTime) ([]*Link,error) {

    var _modelSlice []*Link
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "link_updated", _findByLinkUpdated)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Link{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByLinkRel dynamic finder for link_rel -> []*Link,error
// Generic and programatically generator finder for Link
//```go  
//    m := NewLink(a)
//    results,err := m.FindByLinkRel(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Link
//    }
//```  
        func (o *Link) FindByLinkRel(_findByLinkRel string) ([]*Link,error) {

    var _modelSlice []*Link
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "link_rel", _findByLinkRel)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Link{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByLinkNotes dynamic finder for link_notes -> []*Link,error
// Generic and programatically generator finder for Link
//```go  
//    m := NewLink(a)
//    results,err := m.FindByLinkNotes(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Link
//    }
//```  
        func (o *Link) FindByLinkNotes(_findByLinkNotes string) ([]*Link,error) {

    var _modelSlice []*Link
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "link_notes", _findByLinkNotes)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Link{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByLinkRss dynamic finder for link_rss -> []*Link,error
// Generic and programatically generator finder for Link
//```go  
//    m := NewLink(a)
//    results,err := m.FindByLinkRss(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Link
//    }
//```  
        func (o *Link) FindByLinkRss(_findByLinkRss string) ([]*Link,error) {

    var _modelSlice []*Link
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "link_rss", _findByLinkRss)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Link{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}

// FromDBValueMap Converts a DBValueMap returned from Adapter.Query to a Link
func (o *Link) FromDBValueMap(m map[string]DBValue) error {
	_LinkId,err := m["link_id"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.LinkId = _LinkId
	_LinkUrl,err := m["link_url"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.LinkUrl = _LinkUrl
	_LinkName,err := m["link_name"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.LinkName = _LinkName
	_LinkImage,err := m["link_image"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.LinkImage = _LinkImage
	_LinkTarget,err := m["link_target"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.LinkTarget = _LinkTarget
	_LinkDescription,err := m["link_description"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.LinkDescription = _LinkDescription
	_LinkVisible,err := m["link_visible"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.LinkVisible = _LinkVisible
	_LinkOwner,err := m["link_owner"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.LinkOwner = _LinkOwner
	_LinkRating,err := m["link_rating"].AsInt()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.LinkRating = _LinkRating
	_LinkUpdated,err := m["link_updated"].AsDateTime()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.LinkUpdated = _LinkUpdated
	_LinkRel,err := m["link_rel"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.LinkRel = _LinkRel
	_LinkNotes,err := m["link_notes"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.LinkNotes = _LinkNotes
	_LinkRss,err := m["link_rss"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.LinkRss = _LinkRss

 	return nil
}
// FromLink A kind of Clone function for Link
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
// Reload A function to forcibly reload Link
func (o *Link) Reload() error {
    _,err := o.Find(o.GetPrimaryKeyValue())
    return err
}

// Save is a dynamic saver 'inherited' by all models
func (o *Link) Save() error {
    if o._new == true {
        return o.Create()
    }
    var sets []string
    
    if o.IsLinkUrlDirty == true {
        sets = append(sets,fmt.Sprintf(`link_url = '%s'`,o._adapter.SafeString(o.LinkUrl)))
    }

    if o.IsLinkNameDirty == true {
        sets = append(sets,fmt.Sprintf(`link_name = '%s'`,o._adapter.SafeString(o.LinkName)))
    }

    if o.IsLinkImageDirty == true {
        sets = append(sets,fmt.Sprintf(`link_image = '%s'`,o._adapter.SafeString(o.LinkImage)))
    }

    if o.IsLinkTargetDirty == true {
        sets = append(sets,fmt.Sprintf(`link_target = '%s'`,o._adapter.SafeString(o.LinkTarget)))
    }

    if o.IsLinkDescriptionDirty == true {
        sets = append(sets,fmt.Sprintf(`link_description = '%s'`,o._adapter.SafeString(o.LinkDescription)))
    }

    if o.IsLinkVisibleDirty == true {
        sets = append(sets,fmt.Sprintf(`link_visible = '%s'`,o._adapter.SafeString(o.LinkVisible)))
    }

    if o.IsLinkOwnerDirty == true {
        sets = append(sets,fmt.Sprintf(`link_owner = '%d'`,o.LinkOwner))
    }

    if o.IsLinkRatingDirty == true {
        sets = append(sets,fmt.Sprintf(`link_rating = '%d'`,o.LinkRating))
    }

    if o.IsLinkUpdatedDirty == true {
        sets = append(sets,fmt.Sprintf(`link_updated = '%s'`,o.LinkUpdated))
    }

    if o.IsLinkRelDirty == true {
        sets = append(sets,fmt.Sprintf(`link_rel = '%s'`,o._adapter.SafeString(o.LinkRel)))
    }

    if o.IsLinkNotesDirty == true {
        sets = append(sets,fmt.Sprintf(`link_notes = '%s'`,o._adapter.SafeString(o.LinkNotes)))
    }

    if o.IsLinkRssDirty == true {
        sets = append(sets,fmt.Sprintf(`link_rss = '%s'`,o._adapter.SafeString(o.LinkRss)))
    }

    frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'",o._table,strings.Join(sets,`,`),o._pkey, o.LinkId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return err
    }
    return nil
}
// Update is a dynamic updater, it considers whether or not
// a field is 'dirty' and needs to be updated. Will only work
// if you use the Getters and Setters
func (o *Link) Update() error {
    var sets []string
    
    if o.IsLinkUrlDirty == true {
        sets = append(sets,fmt.Sprintf(`link_url = '%s'`,o._adapter.SafeString(o.LinkUrl)))
    }

    if o.IsLinkNameDirty == true {
        sets = append(sets,fmt.Sprintf(`link_name = '%s'`,o._adapter.SafeString(o.LinkName)))
    }

    if o.IsLinkImageDirty == true {
        sets = append(sets,fmt.Sprintf(`link_image = '%s'`,o._adapter.SafeString(o.LinkImage)))
    }

    if o.IsLinkTargetDirty == true {
        sets = append(sets,fmt.Sprintf(`link_target = '%s'`,o._adapter.SafeString(o.LinkTarget)))
    }

    if o.IsLinkDescriptionDirty == true {
        sets = append(sets,fmt.Sprintf(`link_description = '%s'`,o._adapter.SafeString(o.LinkDescription)))
    }

    if o.IsLinkVisibleDirty == true {
        sets = append(sets,fmt.Sprintf(`link_visible = '%s'`,o._adapter.SafeString(o.LinkVisible)))
    }

    if o.IsLinkOwnerDirty == true {
        sets = append(sets,fmt.Sprintf(`link_owner = '%d'`,o.LinkOwner))
    }

    if o.IsLinkRatingDirty == true {
        sets = append(sets,fmt.Sprintf(`link_rating = '%d'`,o.LinkRating))
    }

    if o.IsLinkUpdatedDirty == true {
        sets = append(sets,fmt.Sprintf(`link_updated = '%s'`,o.LinkUpdated))
    }

    if o.IsLinkRelDirty == true {
        sets = append(sets,fmt.Sprintf(`link_rel = '%s'`,o._adapter.SafeString(o.LinkRel)))
    }

    if o.IsLinkNotesDirty == true {
        sets = append(sets,fmt.Sprintf(`link_notes = '%s'`,o._adapter.SafeString(o.LinkNotes)))
    }

    if o.IsLinkRssDirty == true {
        sets = append(sets,fmt.Sprintf(`link_rss = '%s'`,o._adapter.SafeString(o.LinkRss)))
    }

    frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'",o._table,strings.Join(sets,`,`),o._pkey, o.LinkId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return err
    }
    return nil
}
// Create inserts the model. Calling Save will call this function
// automatically for new models
func (o *Link) Create() error {
    frmt := fmt.Sprintf("INSERT INTO %s (`link_url`, `link_name`, `link_image`, `link_target`, `link_description`, `link_visible`, `link_owner`, `link_rating`, `link_updated`, `link_rel`, `link_notes`, `link_rss`) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%d', '%d', '%s', '%s', '%s', '%s')",o._table,o.LinkUrl, o.LinkName, o.LinkImage, o.LinkTarget, o.LinkDescription, o.LinkVisible, o.LinkOwner, o.LinkRating, o.LinkUpdated.ToString(), o.LinkRel, o.LinkNotes, o.LinkRss)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return o._adapter.Oops(fmt.Sprintf(`%s led to %s`,frmt,err))
    }
    o.LinkId = o._adapter.LastInsertedId()
    o._new = false
    return nil
}


// UpdateLinkUrl an immediate DB Query to update a single column, in this
// case link_url
func (o *Link) UpdateLinkUrl(_updLinkUrl string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `link_url` = '%s' WHERE `link_id` = '%d'",o._table,_updLinkUrl,o.LinkId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.LinkUrl = _updLinkUrl
    return o._adapter.AffectedRows(),nil
}

// UpdateLinkName an immediate DB Query to update a single column, in this
// case link_name
func (o *Link) UpdateLinkName(_updLinkName string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `link_name` = '%s' WHERE `link_id` = '%d'",o._table,_updLinkName,o.LinkId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.LinkName = _updLinkName
    return o._adapter.AffectedRows(),nil
}

// UpdateLinkImage an immediate DB Query to update a single column, in this
// case link_image
func (o *Link) UpdateLinkImage(_updLinkImage string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `link_image` = '%s' WHERE `link_id` = '%d'",o._table,_updLinkImage,o.LinkId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.LinkImage = _updLinkImage
    return o._adapter.AffectedRows(),nil
}

// UpdateLinkTarget an immediate DB Query to update a single column, in this
// case link_target
func (o *Link) UpdateLinkTarget(_updLinkTarget string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `link_target` = '%s' WHERE `link_id` = '%d'",o._table,_updLinkTarget,o.LinkId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.LinkTarget = _updLinkTarget
    return o._adapter.AffectedRows(),nil
}

// UpdateLinkDescription an immediate DB Query to update a single column, in this
// case link_description
func (o *Link) UpdateLinkDescription(_updLinkDescription string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `link_description` = '%s' WHERE `link_id` = '%d'",o._table,_updLinkDescription,o.LinkId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.LinkDescription = _updLinkDescription
    return o._adapter.AffectedRows(),nil
}

// UpdateLinkVisible an immediate DB Query to update a single column, in this
// case link_visible
func (o *Link) UpdateLinkVisible(_updLinkVisible string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `link_visible` = '%s' WHERE `link_id` = '%d'",o._table,_updLinkVisible,o.LinkId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.LinkVisible = _updLinkVisible
    return o._adapter.AffectedRows(),nil
}

// UpdateLinkOwner an immediate DB Query to update a single column, in this
// case link_owner
func (o *Link) UpdateLinkOwner(_updLinkOwner int64) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `link_owner` = '%d' WHERE `link_id` = '%d'",o._table,_updLinkOwner,o.LinkId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.LinkOwner = _updLinkOwner
    return o._adapter.AffectedRows(),nil
}

// UpdateLinkRating an immediate DB Query to update a single column, in this
// case link_rating
func (o *Link) UpdateLinkRating(_updLinkRating int) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `link_rating` = '%d' WHERE `link_id` = '%d'",o._table,_updLinkRating,o.LinkId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.LinkRating = _updLinkRating
    return o._adapter.AffectedRows(),nil
}

// UpdateLinkUpdated an immediate DB Query to update a single column, in this
// case link_updated
func (o *Link) UpdateLinkUpdated(_updLinkUpdated *DateTime) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `link_updated` = '%s' WHERE `link_id` = '%d'",o._table,_updLinkUpdated,o.LinkId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.LinkUpdated = _updLinkUpdated
    return o._adapter.AffectedRows(),nil
}

// UpdateLinkRel an immediate DB Query to update a single column, in this
// case link_rel
func (o *Link) UpdateLinkRel(_updLinkRel string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `link_rel` = '%s' WHERE `link_id` = '%d'",o._table,_updLinkRel,o.LinkId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.LinkRel = _updLinkRel
    return o._adapter.AffectedRows(),nil
}

// UpdateLinkNotes an immediate DB Query to update a single column, in this
// case link_notes
func (o *Link) UpdateLinkNotes(_updLinkNotes string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `link_notes` = '%s' WHERE `link_id` = '%d'",o._table,_updLinkNotes,o.LinkId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.LinkNotes = _updLinkNotes
    return o._adapter.AffectedRows(),nil
}

// UpdateLinkRss an immediate DB Query to update a single column, in this
// case link_rss
func (o *Link) UpdateLinkRss(_updLinkRss string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `link_rss` = '%s' WHERE `link_id` = '%d'",o._table,_updLinkRss,o.LinkId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.LinkRss = _updLinkRss
    return o._adapter.AffectedRows(),nil
}

// Option is a Object Relational Mapping to
// the database table that represents it. In this case it is
// options. The table name will be Sprintf'd to include
// the prefix you define in your YAML configuration for the
// Adapter.
type Option struct {
    _table string
    _adapter Adapter
    _pkey string // 0 The name of the primary key in this table
    _conds []string
    _new bool
    OptionId int64
    OptionName string
    OptionValue string
    Autoload string
	// Dirty markers for smart updates
    IsOptionIdDirty bool
    IsOptionNameDirty bool
    IsOptionValueDirty bool
    IsAutoloadDirty bool
	// Relationships
}

// NewOption binds an Adapter to a new instance
// of Option and sets up the _table and primary keys
func NewOption(a Adapter) *Option {
    var o Option
    o._table = fmt.Sprintf("%soptions",a.DatabasePrefix())
    o._adapter = a
    o._pkey = "option_id"
    o._new = false
    return &o
}


// GetPrimaryKeyValue returns the value, usually int64 of
// the PrimaryKey
func (o *Option) GetPrimaryKeyValue() int64 {
    return o.OptionId
}
// GetPrimaryKeyName returns the DB field name
func (o *Option) GetPrimaryKeyName() string {
    return `option_id`
}

// GetOptionId returns the value of 
// Option.OptionId
func (o *Option) GetOptionId() int64 {
    return o.OptionId
}
// SetOptionId sets and marks as dirty the value of
// Option.OptionId
func (o *Option) SetOptionId(arg int64) {
    o.OptionId = arg
    o.IsOptionIdDirty = true
}

// GetOptionName returns the value of 
// Option.OptionName
func (o *Option) GetOptionName() string {
    return o.OptionName
}
// SetOptionName sets and marks as dirty the value of
// Option.OptionName
func (o *Option) SetOptionName(arg string) {
    o.OptionName = arg
    o.IsOptionNameDirty = true
}

// GetOptionValue returns the value of 
// Option.OptionValue
func (o *Option) GetOptionValue() string {
    return o.OptionValue
}
// SetOptionValue sets and marks as dirty the value of
// Option.OptionValue
func (o *Option) SetOptionValue(arg string) {
    o.OptionValue = arg
    o.IsOptionValueDirty = true
}

// GetAutoload returns the value of 
// Option.Autoload
func (o *Option) GetAutoload() string {
    return o.Autoload
}
// SetAutoload sets and marks as dirty the value of
// Option.Autoload
func (o *Option) SetAutoload(arg string) {
    o.Autoload = arg
    o.IsAutoloadDirty = true
}

// Find dynamic finder for option_id -> bool,error
// Generic and programatically generator finder for Option
// Note that Fine returns a bool if found, not err, in the case of
// a return of true, the instance data will be filled out.
// a call to find ALWAYS overwrites the model you call Find on
// i.e. receiver is a pointer. 
//```go
//      m := NewOption(a)
//      found,err := m.Find(23)
//      .. handle err
//      if found == false {
//          // handle found
//      }
//      ... do what you want with m here
//```
        func (o *Option) Find(_findByOptionId int64) (bool,error) {

    var _modelSlice []*Option
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "option_id", _findByOptionId)
    results, err := o._adapter.Query(q)
    if err != nil {
        return false,o._adapter.Oops(fmt.Sprintf(`%s`,err))
    }
    
    for _,result := range results {
        ro := Option{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return false,o._adapter.Oops(fmt.Sprintf(`%s`,err))
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return false, o._adapter.Oops(`not found`)
    }
    o.FromOption(_modelSlice[0])
    return true,nil

}
// FindByOptionName dynamic finder for option_name -> []*Option,error
// Generic and programatically generator finder for Option
//```go  
//    m := NewOption(a)
//    results,err := m.FindByOptionName(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Option
//    }
//```  
        func (o *Option) FindByOptionName(_findByOptionName string) ([]*Option,error) {

    var _modelSlice []*Option
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "option_name", _findByOptionName)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Option{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByOptionValue dynamic finder for option_value -> []*Option,error
// Generic and programatically generator finder for Option
//```go  
//    m := NewOption(a)
//    results,err := m.FindByOptionValue(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Option
//    }
//```  
        func (o *Option) FindByOptionValue(_findByOptionValue string) ([]*Option,error) {

    var _modelSlice []*Option
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "option_value", _findByOptionValue)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Option{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByAutoload dynamic finder for autoload -> []*Option,error
// Generic and programatically generator finder for Option
//```go  
//    m := NewOption(a)
//    results,err := m.FindByAutoload(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Option
//    }
//```  
        func (o *Option) FindByAutoload(_findByAutoload string) ([]*Option,error) {

    var _modelSlice []*Option
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "autoload", _findByAutoload)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Option{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}

// FromDBValueMap Converts a DBValueMap returned from Adapter.Query to a Option
func (o *Option) FromDBValueMap(m map[string]DBValue) error {
	_OptionId,err := m["option_id"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.OptionId = _OptionId
	_OptionName,err := m["option_name"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.OptionName = _OptionName
	_OptionValue,err := m["option_value"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.OptionValue = _OptionValue
	_Autoload,err := m["autoload"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.Autoload = _Autoload

 	return nil
}
// FromOption A kind of Clone function for Option
func (o *Option) FromOption(m *Option) {
	o.OptionId = m.OptionId
	o.OptionName = m.OptionName
	o.OptionValue = m.OptionValue
	o.Autoload = m.Autoload

}
// Reload A function to forcibly reload Option
func (o *Option) Reload() error {
    _,err := o.Find(o.GetPrimaryKeyValue())
    return err
}

// Save is a dynamic saver 'inherited' by all models
func (o *Option) Save() error {
    if o._new == true {
        return o.Create()
    }
    var sets []string
    
    if o.IsOptionNameDirty == true {
        sets = append(sets,fmt.Sprintf(`option_name = '%s'`,o._adapter.SafeString(o.OptionName)))
    }

    if o.IsOptionValueDirty == true {
        sets = append(sets,fmt.Sprintf(`option_value = '%s'`,o._adapter.SafeString(o.OptionValue)))
    }

    if o.IsAutoloadDirty == true {
        sets = append(sets,fmt.Sprintf(`autoload = '%s'`,o._adapter.SafeString(o.Autoload)))
    }

    frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'",o._table,strings.Join(sets,`,`),o._pkey, o.OptionId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return err
    }
    return nil
}
// Update is a dynamic updater, it considers whether or not
// a field is 'dirty' and needs to be updated. Will only work
// if you use the Getters and Setters
func (o *Option) Update() error {
    var sets []string
    
    if o.IsOptionNameDirty == true {
        sets = append(sets,fmt.Sprintf(`option_name = '%s'`,o._adapter.SafeString(o.OptionName)))
    }

    if o.IsOptionValueDirty == true {
        sets = append(sets,fmt.Sprintf(`option_value = '%s'`,o._adapter.SafeString(o.OptionValue)))
    }

    if o.IsAutoloadDirty == true {
        sets = append(sets,fmt.Sprintf(`autoload = '%s'`,o._adapter.SafeString(o.Autoload)))
    }

    frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'",o._table,strings.Join(sets,`,`),o._pkey, o.OptionId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return err
    }
    return nil
}
// Create inserts the model. Calling Save will call this function
// automatically for new models
func (o *Option) Create() error {
    frmt := fmt.Sprintf("INSERT INTO %s (`option_name`, `option_value`, `autoload`) VALUES ('%s', '%s', '%s')",o._table,o.OptionName, o.OptionValue, o.Autoload)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return o._adapter.Oops(fmt.Sprintf(`%s led to %s`,frmt,err))
    }
    o.OptionId = o._adapter.LastInsertedId()
    o._new = false
    return nil
}


// UpdateOptionName an immediate DB Query to update a single column, in this
// case option_name
func (o *Option) UpdateOptionName(_updOptionName string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `option_name` = '%s' WHERE `option_id` = '%d'",o._table,_updOptionName,o.OptionId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.OptionName = _updOptionName
    return o._adapter.AffectedRows(),nil
}

// UpdateOptionValue an immediate DB Query to update a single column, in this
// case option_value
func (o *Option) UpdateOptionValue(_updOptionValue string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `option_value` = '%s' WHERE `option_id` = '%d'",o._table,_updOptionValue,o.OptionId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.OptionValue = _updOptionValue
    return o._adapter.AffectedRows(),nil
}

// UpdateAutoload an immediate DB Query to update a single column, in this
// case autoload
func (o *Option) UpdateAutoload(_updAutoload string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `autoload` = '%s' WHERE `option_id` = '%d'",o._table,_updAutoload,o.OptionId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.Autoload = _updAutoload
    return o._adapter.AffectedRows(),nil
}

// PostMeta is a Object Relational Mapping to
// the database table that represents it. In this case it is
// postmeta. The table name will be Sprintf'd to include
// the prefix you define in your YAML configuration for the
// Adapter.
type PostMeta struct {
    _table string
    _adapter Adapter
    _pkey string // 0 The name of the primary key in this table
    _conds []string
    _new bool
    MetaId int64
    PostId int64
    MetaKey string
    MetaValue string
	// Dirty markers for smart updates
    IsMetaIdDirty bool
    IsPostIdDirty bool
    IsMetaKeyDirty bool
    IsMetaValueDirty bool
	// Relationships
}

// NewPostMeta binds an Adapter to a new instance
// of PostMeta and sets up the _table and primary keys
func NewPostMeta(a Adapter) *PostMeta {
    var o PostMeta
    o._table = fmt.Sprintf("%spostmeta",a.DatabasePrefix())
    o._adapter = a
    o._pkey = "meta_id"
    o._new = false
    return &o
}


// GetPrimaryKeyValue returns the value, usually int64 of
// the PrimaryKey
func (o *PostMeta) GetPrimaryKeyValue() int64 {
    return o.MetaId
}
// GetPrimaryKeyName returns the DB field name
func (o *PostMeta) GetPrimaryKeyName() string {
    return `meta_id`
}

// GetMetaId returns the value of 
// PostMeta.MetaId
func (o *PostMeta) GetMetaId() int64 {
    return o.MetaId
}
// SetMetaId sets and marks as dirty the value of
// PostMeta.MetaId
func (o *PostMeta) SetMetaId(arg int64) {
    o.MetaId = arg
    o.IsMetaIdDirty = true
}

// GetPostId returns the value of 
// PostMeta.PostId
func (o *PostMeta) GetPostId() int64 {
    return o.PostId
}
// SetPostId sets and marks as dirty the value of
// PostMeta.PostId
func (o *PostMeta) SetPostId(arg int64) {
    o.PostId = arg
    o.IsPostIdDirty = true
}

// GetMetaKey returns the value of 
// PostMeta.MetaKey
func (o *PostMeta) GetMetaKey() string {
    return o.MetaKey
}
// SetMetaKey sets and marks as dirty the value of
// PostMeta.MetaKey
func (o *PostMeta) SetMetaKey(arg string) {
    o.MetaKey = arg
    o.IsMetaKeyDirty = true
}

// GetMetaValue returns the value of 
// PostMeta.MetaValue
func (o *PostMeta) GetMetaValue() string {
    return o.MetaValue
}
// SetMetaValue sets and marks as dirty the value of
// PostMeta.MetaValue
func (o *PostMeta) SetMetaValue(arg string) {
    o.MetaValue = arg
    o.IsMetaValueDirty = true
}

// Find dynamic finder for meta_id -> bool,error
// Generic and programatically generator finder for PostMeta
// Note that Fine returns a bool if found, not err, in the case of
// a return of true, the instance data will be filled out.
// a call to find ALWAYS overwrites the model you call Find on
// i.e. receiver is a pointer. 
//```go
//      m := NewPostMeta(a)
//      found,err := m.Find(23)
//      .. handle err
//      if found == false {
//          // handle found
//      }
//      ... do what you want with m here
//```
        func (o *PostMeta) Find(_findByMetaId int64) (bool,error) {

    var _modelSlice []*PostMeta
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "meta_id", _findByMetaId)
    results, err := o._adapter.Query(q)
    if err != nil {
        return false,o._adapter.Oops(fmt.Sprintf(`%s`,err))
    }
    
    for _,result := range results {
        ro := PostMeta{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return false,o._adapter.Oops(fmt.Sprintf(`%s`,err))
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return false, o._adapter.Oops(`not found`)
    }
    o.FromPostMeta(_modelSlice[0])
    return true,nil

}
// FindByPostId dynamic finder for post_id -> []*PostMeta,error
// Generic and programatically generator finder for PostMeta
//```go  
//    m := NewPostMeta(a)
//    results,err := m.FindByPostId(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of PostMeta
//    }
//```  
        func (o *PostMeta) FindByPostId(_findByPostId int64) ([]*PostMeta,error) {

    var _modelSlice []*PostMeta
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "post_id", _findByPostId)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := PostMeta{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByMetaKey dynamic finder for meta_key -> []*PostMeta,error
// Generic and programatically generator finder for PostMeta
//```go  
//    m := NewPostMeta(a)
//    results,err := m.FindByMetaKey(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of PostMeta
//    }
//```  
        func (o *PostMeta) FindByMetaKey(_findByMetaKey string) ([]*PostMeta,error) {

    var _modelSlice []*PostMeta
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "meta_key", _findByMetaKey)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := PostMeta{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByMetaValue dynamic finder for meta_value -> []*PostMeta,error
// Generic and programatically generator finder for PostMeta
//```go  
//    m := NewPostMeta(a)
//    results,err := m.FindByMetaValue(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of PostMeta
//    }
//```  
        func (o *PostMeta) FindByMetaValue(_findByMetaValue string) ([]*PostMeta,error) {

    var _modelSlice []*PostMeta
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "meta_value", _findByMetaValue)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := PostMeta{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}

// FromDBValueMap Converts a DBValueMap returned from Adapter.Query to a PostMeta
func (o *PostMeta) FromDBValueMap(m map[string]DBValue) error {
	_MetaId,err := m["meta_id"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.MetaId = _MetaId
	_PostId,err := m["post_id"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.PostId = _PostId
	_MetaKey,err := m["meta_key"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.MetaKey = _MetaKey
	_MetaValue,err := m["meta_value"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.MetaValue = _MetaValue

 	return nil
}
// FromPostMeta A kind of Clone function for PostMeta
func (o *PostMeta) FromPostMeta(m *PostMeta) {
	o.MetaId = m.MetaId
	o.PostId = m.PostId
	o.MetaKey = m.MetaKey
	o.MetaValue = m.MetaValue

}
// Reload A function to forcibly reload PostMeta
func (o *PostMeta) Reload() error {
    _,err := o.Find(o.GetPrimaryKeyValue())
    return err
}

// Save is a dynamic saver 'inherited' by all models
func (o *PostMeta) Save() error {
    if o._new == true {
        return o.Create()
    }
    var sets []string
    
    if o.IsPostIdDirty == true {
        sets = append(sets,fmt.Sprintf(`post_id = '%d'`,o.PostId))
    }

    if o.IsMetaKeyDirty == true {
        sets = append(sets,fmt.Sprintf(`meta_key = '%s'`,o._adapter.SafeString(o.MetaKey)))
    }

    if o.IsMetaValueDirty == true {
        sets = append(sets,fmt.Sprintf(`meta_value = '%s'`,o._adapter.SafeString(o.MetaValue)))
    }

    frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'",o._table,strings.Join(sets,`,`),o._pkey, o.MetaId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return err
    }
    return nil
}
// Update is a dynamic updater, it considers whether or not
// a field is 'dirty' and needs to be updated. Will only work
// if you use the Getters and Setters
func (o *PostMeta) Update() error {
    var sets []string
    
    if o.IsPostIdDirty == true {
        sets = append(sets,fmt.Sprintf(`post_id = '%d'`,o.PostId))
    }

    if o.IsMetaKeyDirty == true {
        sets = append(sets,fmt.Sprintf(`meta_key = '%s'`,o._adapter.SafeString(o.MetaKey)))
    }

    if o.IsMetaValueDirty == true {
        sets = append(sets,fmt.Sprintf(`meta_value = '%s'`,o._adapter.SafeString(o.MetaValue)))
    }

    frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'",o._table,strings.Join(sets,`,`),o._pkey, o.MetaId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return err
    }
    return nil
}
// Create inserts the model. Calling Save will call this function
// automatically for new models
func (o *PostMeta) Create() error {
    frmt := fmt.Sprintf("INSERT INTO %s (`post_id`, `meta_key`, `meta_value`) VALUES ('%d', '%s', '%s')",o._table,o.PostId, o.MetaKey, o.MetaValue)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return o._adapter.Oops(fmt.Sprintf(`%s led to %s`,frmt,err))
    }
    o.MetaId = o._adapter.LastInsertedId()
    o._new = false
    return nil
}


// UpdatePostId an immediate DB Query to update a single column, in this
// case post_id
func (o *PostMeta) UpdatePostId(_updPostId int64) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `post_id` = '%d' WHERE `meta_id` = '%d'",o._table,_updPostId,o.MetaId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.PostId = _updPostId
    return o._adapter.AffectedRows(),nil
}

// UpdateMetaKey an immediate DB Query to update a single column, in this
// case meta_key
func (o *PostMeta) UpdateMetaKey(_updMetaKey string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `meta_key` = '%s' WHERE `meta_id` = '%d'",o._table,_updMetaKey,o.MetaId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.MetaKey = _updMetaKey
    return o._adapter.AffectedRows(),nil
}

// UpdateMetaValue an immediate DB Query to update a single column, in this
// case meta_value
func (o *PostMeta) UpdateMetaValue(_updMetaValue string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `meta_value` = '%s' WHERE `meta_id` = '%d'",o._table,_updMetaValue,o.MetaId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.MetaValue = _updMetaValue
    return o._adapter.AffectedRows(),nil
}

// Post is a Object Relational Mapping to
// the database table that represents it. In this case it is
// posts. The table name will be Sprintf'd to include
// the prefix you define in your YAML configuration for the
// Adapter.
type Post struct {
    _table string
    _adapter Adapter
    _pkey string // 0 The name of the primary key in this table
    _conds []string
    _new bool
    ID int64
    PostAuthor int64
    PostDate *DateTime
    PostDateGmt *DateTime
    PostContent string
    PostTitle string
    PostExcerpt string
    PostStatus string
    CommentStatus string
    PingStatus string
    PostPassword string
    PostName string
    ToPing string
    Pinged string
    PostModified *DateTime
    PostModifiedGmt *DateTime
    PostContentFiltered string
    PostParent int64
    Guid string
    MenuOrder int
    PostType string
    PostMimeType string
    CommentCount int64
	// Dirty markers for smart updates
    IsIDDirty bool
    IsPostAuthorDirty bool
    IsPostDateDirty bool
    IsPostDateGmtDirty bool
    IsPostContentDirty bool
    IsPostTitleDirty bool
    IsPostExcerptDirty bool
    IsPostStatusDirty bool
    IsCommentStatusDirty bool
    IsPingStatusDirty bool
    IsPostPasswordDirty bool
    IsPostNameDirty bool
    IsToPingDirty bool
    IsPingedDirty bool
    IsPostModifiedDirty bool
    IsPostModifiedGmtDirty bool
    IsPostContentFilteredDirty bool
    IsPostParentDirty bool
    IsGuidDirty bool
    IsMenuOrderDirty bool
    IsPostTypeDirty bool
    IsPostMimeTypeDirty bool
    IsCommentCountDirty bool
	// Relationships
}

// NewPost binds an Adapter to a new instance
// of Post and sets up the _table and primary keys
func NewPost(a Adapter) *Post {
    var o Post
    o._table = fmt.Sprintf("%sposts",a.DatabasePrefix())
    o._adapter = a
    o._pkey = "ID"
    o._new = false
    return &o
}


// GetPrimaryKeyValue returns the value, usually int64 of
// the PrimaryKey
func (o *Post) GetPrimaryKeyValue() int64 {
    return o.ID
}
// GetPrimaryKeyName returns the DB field name
func (o *Post) GetPrimaryKeyName() string {
    return `ID`
}

// GetID returns the value of 
// Post.ID
func (o *Post) GetID() int64 {
    return o.ID
}
// SetID sets and marks as dirty the value of
// Post.ID
func (o *Post) SetID(arg int64) {
    o.ID = arg
    o.IsIDDirty = true
}

// GetPostAuthor returns the value of 
// Post.PostAuthor
func (o *Post) GetPostAuthor() int64 {
    return o.PostAuthor
}
// SetPostAuthor sets and marks as dirty the value of
// Post.PostAuthor
func (o *Post) SetPostAuthor(arg int64) {
    o.PostAuthor = arg
    o.IsPostAuthorDirty = true
}

// GetPostDate returns the value of 
// Post.PostDate
func (o *Post) GetPostDate() *DateTime {
    return o.PostDate
}
// SetPostDate sets and marks as dirty the value of
// Post.PostDate
func (o *Post) SetPostDate(arg *DateTime) {
    o.PostDate = arg
    o.IsPostDateDirty = true
}

// GetPostDateGmt returns the value of 
// Post.PostDateGmt
func (o *Post) GetPostDateGmt() *DateTime {
    return o.PostDateGmt
}
// SetPostDateGmt sets and marks as dirty the value of
// Post.PostDateGmt
func (o *Post) SetPostDateGmt(arg *DateTime) {
    o.PostDateGmt = arg
    o.IsPostDateGmtDirty = true
}

// GetPostContent returns the value of 
// Post.PostContent
func (o *Post) GetPostContent() string {
    return o.PostContent
}
// SetPostContent sets and marks as dirty the value of
// Post.PostContent
func (o *Post) SetPostContent(arg string) {
    o.PostContent = arg
    o.IsPostContentDirty = true
}

// GetPostTitle returns the value of 
// Post.PostTitle
func (o *Post) GetPostTitle() string {
    return o.PostTitle
}
// SetPostTitle sets and marks as dirty the value of
// Post.PostTitle
func (o *Post) SetPostTitle(arg string) {
    o.PostTitle = arg
    o.IsPostTitleDirty = true
}

// GetPostExcerpt returns the value of 
// Post.PostExcerpt
func (o *Post) GetPostExcerpt() string {
    return o.PostExcerpt
}
// SetPostExcerpt sets and marks as dirty the value of
// Post.PostExcerpt
func (o *Post) SetPostExcerpt(arg string) {
    o.PostExcerpt = arg
    o.IsPostExcerptDirty = true
}

// GetPostStatus returns the value of 
// Post.PostStatus
func (o *Post) GetPostStatus() string {
    return o.PostStatus
}
// SetPostStatus sets and marks as dirty the value of
// Post.PostStatus
func (o *Post) SetPostStatus(arg string) {
    o.PostStatus = arg
    o.IsPostStatusDirty = true
}

// GetCommentStatus returns the value of 
// Post.CommentStatus
func (o *Post) GetCommentStatus() string {
    return o.CommentStatus
}
// SetCommentStatus sets and marks as dirty the value of
// Post.CommentStatus
func (o *Post) SetCommentStatus(arg string) {
    o.CommentStatus = arg
    o.IsCommentStatusDirty = true
}

// GetPingStatus returns the value of 
// Post.PingStatus
func (o *Post) GetPingStatus() string {
    return o.PingStatus
}
// SetPingStatus sets and marks as dirty the value of
// Post.PingStatus
func (o *Post) SetPingStatus(arg string) {
    o.PingStatus = arg
    o.IsPingStatusDirty = true
}

// GetPostPassword returns the value of 
// Post.PostPassword
func (o *Post) GetPostPassword() string {
    return o.PostPassword
}
// SetPostPassword sets and marks as dirty the value of
// Post.PostPassword
func (o *Post) SetPostPassword(arg string) {
    o.PostPassword = arg
    o.IsPostPasswordDirty = true
}

// GetPostName returns the value of 
// Post.PostName
func (o *Post) GetPostName() string {
    return o.PostName
}
// SetPostName sets and marks as dirty the value of
// Post.PostName
func (o *Post) SetPostName(arg string) {
    o.PostName = arg
    o.IsPostNameDirty = true
}

// GetToPing returns the value of 
// Post.ToPing
func (o *Post) GetToPing() string {
    return o.ToPing
}
// SetToPing sets and marks as dirty the value of
// Post.ToPing
func (o *Post) SetToPing(arg string) {
    o.ToPing = arg
    o.IsToPingDirty = true
}

// GetPinged returns the value of 
// Post.Pinged
func (o *Post) GetPinged() string {
    return o.Pinged
}
// SetPinged sets and marks as dirty the value of
// Post.Pinged
func (o *Post) SetPinged(arg string) {
    o.Pinged = arg
    o.IsPingedDirty = true
}

// GetPostModified returns the value of 
// Post.PostModified
func (o *Post) GetPostModified() *DateTime {
    return o.PostModified
}
// SetPostModified sets and marks as dirty the value of
// Post.PostModified
func (o *Post) SetPostModified(arg *DateTime) {
    o.PostModified = arg
    o.IsPostModifiedDirty = true
}

// GetPostModifiedGmt returns the value of 
// Post.PostModifiedGmt
func (o *Post) GetPostModifiedGmt() *DateTime {
    return o.PostModifiedGmt
}
// SetPostModifiedGmt sets and marks as dirty the value of
// Post.PostModifiedGmt
func (o *Post) SetPostModifiedGmt(arg *DateTime) {
    o.PostModifiedGmt = arg
    o.IsPostModifiedGmtDirty = true
}

// GetPostContentFiltered returns the value of 
// Post.PostContentFiltered
func (o *Post) GetPostContentFiltered() string {
    return o.PostContentFiltered
}
// SetPostContentFiltered sets and marks as dirty the value of
// Post.PostContentFiltered
func (o *Post) SetPostContentFiltered(arg string) {
    o.PostContentFiltered = arg
    o.IsPostContentFilteredDirty = true
}

// GetPostParent returns the value of 
// Post.PostParent
func (o *Post) GetPostParent() int64 {
    return o.PostParent
}
// SetPostParent sets and marks as dirty the value of
// Post.PostParent
func (o *Post) SetPostParent(arg int64) {
    o.PostParent = arg
    o.IsPostParentDirty = true
}

// GetGuid returns the value of 
// Post.Guid
func (o *Post) GetGuid() string {
    return o.Guid
}
// SetGuid sets and marks as dirty the value of
// Post.Guid
func (o *Post) SetGuid(arg string) {
    o.Guid = arg
    o.IsGuidDirty = true
}

// GetMenuOrder returns the value of 
// Post.MenuOrder
func (o *Post) GetMenuOrder() int {
    return o.MenuOrder
}
// SetMenuOrder sets and marks as dirty the value of
// Post.MenuOrder
func (o *Post) SetMenuOrder(arg int) {
    o.MenuOrder = arg
    o.IsMenuOrderDirty = true
}

// GetPostType returns the value of 
// Post.PostType
func (o *Post) GetPostType() string {
    return o.PostType
}
// SetPostType sets and marks as dirty the value of
// Post.PostType
func (o *Post) SetPostType(arg string) {
    o.PostType = arg
    o.IsPostTypeDirty = true
}

// GetPostMimeType returns the value of 
// Post.PostMimeType
func (o *Post) GetPostMimeType() string {
    return o.PostMimeType
}
// SetPostMimeType sets and marks as dirty the value of
// Post.PostMimeType
func (o *Post) SetPostMimeType(arg string) {
    o.PostMimeType = arg
    o.IsPostMimeTypeDirty = true
}

// GetCommentCount returns the value of 
// Post.CommentCount
func (o *Post) GetCommentCount() int64 {
    return o.CommentCount
}
// SetCommentCount sets and marks as dirty the value of
// Post.CommentCount
func (o *Post) SetCommentCount(arg int64) {
    o.CommentCount = arg
    o.IsCommentCountDirty = true
}

// Find dynamic finder for ID -> bool,error
// Generic and programatically generator finder for Post
// Note that Fine returns a bool if found, not err, in the case of
// a return of true, the instance data will be filled out.
// a call to find ALWAYS overwrites the model you call Find on
// i.e. receiver is a pointer. 
//```go
//      m := NewPost(a)
//      found,err := m.Find(23)
//      .. handle err
//      if found == false {
//          // handle found
//      }
//      ... do what you want with m here
//```
        func (o *Post) Find(_findByID int64) (bool,error) {

    var _modelSlice []*Post
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "ID", _findByID)
    results, err := o._adapter.Query(q)
    if err != nil {
        return false,o._adapter.Oops(fmt.Sprintf(`%s`,err))
    }
    
    for _,result := range results {
        ro := Post{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return false,o._adapter.Oops(fmt.Sprintf(`%s`,err))
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return false, o._adapter.Oops(`not found`)
    }
    o.FromPost(_modelSlice[0])
    return true,nil

}
// FindByPostAuthor dynamic finder for post_author -> []*Post,error
// Generic and programatically generator finder for Post
//```go  
//    m := NewPost(a)
//    results,err := m.FindByPostAuthor(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Post
//    }
//```  
        func (o *Post) FindByPostAuthor(_findByPostAuthor int64) ([]*Post,error) {

    var _modelSlice []*Post
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "post_author", _findByPostAuthor)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Post{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByPostDate dynamic finder for post_date -> []*Post,error
// Generic and programatically generator finder for Post
//```go  
//    m := NewPost(a)
//    results,err := m.FindByPostDate(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Post
//    }
//```  
        func (o *Post) FindByPostDate(_findByPostDate *DateTime) ([]*Post,error) {

    var _modelSlice []*Post
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "post_date", _findByPostDate)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Post{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByPostDateGmt dynamic finder for post_date_gmt -> []*Post,error
// Generic and programatically generator finder for Post
//```go  
//    m := NewPost(a)
//    results,err := m.FindByPostDateGmt(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Post
//    }
//```  
        func (o *Post) FindByPostDateGmt(_findByPostDateGmt *DateTime) ([]*Post,error) {

    var _modelSlice []*Post
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "post_date_gmt", _findByPostDateGmt)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Post{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByPostContent dynamic finder for post_content -> []*Post,error
// Generic and programatically generator finder for Post
//```go  
//    m := NewPost(a)
//    results,err := m.FindByPostContent(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Post
//    }
//```  
        func (o *Post) FindByPostContent(_findByPostContent string) ([]*Post,error) {

    var _modelSlice []*Post
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "post_content", _findByPostContent)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Post{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByPostTitle dynamic finder for post_title -> []*Post,error
// Generic and programatically generator finder for Post
//```go  
//    m := NewPost(a)
//    results,err := m.FindByPostTitle(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Post
//    }
//```  
        func (o *Post) FindByPostTitle(_findByPostTitle string) ([]*Post,error) {

    var _modelSlice []*Post
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "post_title", _findByPostTitle)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Post{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByPostExcerpt dynamic finder for post_excerpt -> []*Post,error
// Generic and programatically generator finder for Post
//```go  
//    m := NewPost(a)
//    results,err := m.FindByPostExcerpt(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Post
//    }
//```  
        func (o *Post) FindByPostExcerpt(_findByPostExcerpt string) ([]*Post,error) {

    var _modelSlice []*Post
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "post_excerpt", _findByPostExcerpt)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Post{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByPostStatus dynamic finder for post_status -> []*Post,error
// Generic and programatically generator finder for Post
//```go  
//    m := NewPost(a)
//    results,err := m.FindByPostStatus(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Post
//    }
//```  
        func (o *Post) FindByPostStatus(_findByPostStatus string) ([]*Post,error) {

    var _modelSlice []*Post
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "post_status", _findByPostStatus)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Post{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByCommentStatus dynamic finder for comment_status -> []*Post,error
// Generic and programatically generator finder for Post
//```go  
//    m := NewPost(a)
//    results,err := m.FindByCommentStatus(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Post
//    }
//```  
        func (o *Post) FindByCommentStatus(_findByCommentStatus string) ([]*Post,error) {

    var _modelSlice []*Post
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "comment_status", _findByCommentStatus)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Post{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByPingStatus dynamic finder for ping_status -> []*Post,error
// Generic and programatically generator finder for Post
//```go  
//    m := NewPost(a)
//    results,err := m.FindByPingStatus(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Post
//    }
//```  
        func (o *Post) FindByPingStatus(_findByPingStatus string) ([]*Post,error) {

    var _modelSlice []*Post
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "ping_status", _findByPingStatus)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Post{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByPostPassword dynamic finder for post_password -> []*Post,error
// Generic and programatically generator finder for Post
//```go  
//    m := NewPost(a)
//    results,err := m.FindByPostPassword(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Post
//    }
//```  
        func (o *Post) FindByPostPassword(_findByPostPassword string) ([]*Post,error) {

    var _modelSlice []*Post
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "post_password", _findByPostPassword)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Post{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByPostName dynamic finder for post_name -> []*Post,error
// Generic and programatically generator finder for Post
//```go  
//    m := NewPost(a)
//    results,err := m.FindByPostName(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Post
//    }
//```  
        func (o *Post) FindByPostName(_findByPostName string) ([]*Post,error) {

    var _modelSlice []*Post
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "post_name", _findByPostName)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Post{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByToPing dynamic finder for to_ping -> []*Post,error
// Generic and programatically generator finder for Post
//```go  
//    m := NewPost(a)
//    results,err := m.FindByToPing(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Post
//    }
//```  
        func (o *Post) FindByToPing(_findByToPing string) ([]*Post,error) {

    var _modelSlice []*Post
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "to_ping", _findByToPing)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Post{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByPinged dynamic finder for pinged -> []*Post,error
// Generic and programatically generator finder for Post
//```go  
//    m := NewPost(a)
//    results,err := m.FindByPinged(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Post
//    }
//```  
        func (o *Post) FindByPinged(_findByPinged string) ([]*Post,error) {

    var _modelSlice []*Post
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "pinged", _findByPinged)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Post{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByPostModified dynamic finder for post_modified -> []*Post,error
// Generic and programatically generator finder for Post
//```go  
//    m := NewPost(a)
//    results,err := m.FindByPostModified(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Post
//    }
//```  
        func (o *Post) FindByPostModified(_findByPostModified *DateTime) ([]*Post,error) {

    var _modelSlice []*Post
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "post_modified", _findByPostModified)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Post{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByPostModifiedGmt dynamic finder for post_modified_gmt -> []*Post,error
// Generic and programatically generator finder for Post
//```go  
//    m := NewPost(a)
//    results,err := m.FindByPostModifiedGmt(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Post
//    }
//```  
        func (o *Post) FindByPostModifiedGmt(_findByPostModifiedGmt *DateTime) ([]*Post,error) {

    var _modelSlice []*Post
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "post_modified_gmt", _findByPostModifiedGmt)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Post{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByPostContentFiltered dynamic finder for post_content_filtered -> []*Post,error
// Generic and programatically generator finder for Post
//```go  
//    m := NewPost(a)
//    results,err := m.FindByPostContentFiltered(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Post
//    }
//```  
        func (o *Post) FindByPostContentFiltered(_findByPostContentFiltered string) ([]*Post,error) {

    var _modelSlice []*Post
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "post_content_filtered", _findByPostContentFiltered)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Post{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByPostParent dynamic finder for post_parent -> []*Post,error
// Generic and programatically generator finder for Post
//```go  
//    m := NewPost(a)
//    results,err := m.FindByPostParent(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Post
//    }
//```  
        func (o *Post) FindByPostParent(_findByPostParent int64) ([]*Post,error) {

    var _modelSlice []*Post
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "post_parent", _findByPostParent)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Post{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByGuid dynamic finder for guid -> []*Post,error
// Generic and programatically generator finder for Post
//```go  
//    m := NewPost(a)
//    results,err := m.FindByGuid(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Post
//    }
//```  
        func (o *Post) FindByGuid(_findByGuid string) ([]*Post,error) {

    var _modelSlice []*Post
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "guid", _findByGuid)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Post{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByMenuOrder dynamic finder for menu_order -> []*Post,error
// Generic and programatically generator finder for Post
//```go  
//    m := NewPost(a)
//    results,err := m.FindByMenuOrder(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Post
//    }
//```  
        func (o *Post) FindByMenuOrder(_findByMenuOrder int) ([]*Post,error) {

    var _modelSlice []*Post
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "menu_order", _findByMenuOrder)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Post{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByPostType dynamic finder for post_type -> []*Post,error
// Generic and programatically generator finder for Post
//```go  
//    m := NewPost(a)
//    results,err := m.FindByPostType(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Post
//    }
//```  
        func (o *Post) FindByPostType(_findByPostType string) ([]*Post,error) {

    var _modelSlice []*Post
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "post_type", _findByPostType)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Post{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByPostMimeType dynamic finder for post_mime_type -> []*Post,error
// Generic and programatically generator finder for Post
//```go  
//    m := NewPost(a)
//    results,err := m.FindByPostMimeType(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Post
//    }
//```  
        func (o *Post) FindByPostMimeType(_findByPostMimeType string) ([]*Post,error) {

    var _modelSlice []*Post
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "post_mime_type", _findByPostMimeType)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Post{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByCommentCount dynamic finder for comment_count -> []*Post,error
// Generic and programatically generator finder for Post
//```go  
//    m := NewPost(a)
//    results,err := m.FindByCommentCount(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Post
//    }
//```  
        func (o *Post) FindByCommentCount(_findByCommentCount int64) ([]*Post,error) {

    var _modelSlice []*Post
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "comment_count", _findByCommentCount)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Post{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}

// FromDBValueMap Converts a DBValueMap returned from Adapter.Query to a Post
func (o *Post) FromDBValueMap(m map[string]DBValue) error {
	_ID,err := m["ID"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.ID = _ID
	_PostAuthor,err := m["post_author"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.PostAuthor = _PostAuthor
	_PostDate,err := m["post_date"].AsDateTime()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.PostDate = _PostDate
	_PostDateGmt,err := m["post_date_gmt"].AsDateTime()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.PostDateGmt = _PostDateGmt
	_PostContent,err := m["post_content"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.PostContent = _PostContent
	_PostTitle,err := m["post_title"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.PostTitle = _PostTitle
	_PostExcerpt,err := m["post_excerpt"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.PostExcerpt = _PostExcerpt
	_PostStatus,err := m["post_status"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.PostStatus = _PostStatus
	_CommentStatus,err := m["comment_status"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.CommentStatus = _CommentStatus
	_PingStatus,err := m["ping_status"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.PingStatus = _PingStatus
	_PostPassword,err := m["post_password"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.PostPassword = _PostPassword
	_PostName,err := m["post_name"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.PostName = _PostName
	_ToPing,err := m["to_ping"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.ToPing = _ToPing
	_Pinged,err := m["pinged"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.Pinged = _Pinged
	_PostModified,err := m["post_modified"].AsDateTime()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.PostModified = _PostModified
	_PostModifiedGmt,err := m["post_modified_gmt"].AsDateTime()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.PostModifiedGmt = _PostModifiedGmt
	_PostContentFiltered,err := m["post_content_filtered"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.PostContentFiltered = _PostContentFiltered
	_PostParent,err := m["post_parent"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.PostParent = _PostParent
	_Guid,err := m["guid"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.Guid = _Guid
	_MenuOrder,err := m["menu_order"].AsInt()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.MenuOrder = _MenuOrder
	_PostType,err := m["post_type"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.PostType = _PostType
	_PostMimeType,err := m["post_mime_type"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.PostMimeType = _PostMimeType
	_CommentCount,err := m["comment_count"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.CommentCount = _CommentCount

 	return nil
}
// FromPost A kind of Clone function for Post
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
// Reload A function to forcibly reload Post
func (o *Post) Reload() error {
    _,err := o.Find(o.GetPrimaryKeyValue())
    return err
}

// Save is a dynamic saver 'inherited' by all models
func (o *Post) Save() error {
    if o._new == true {
        return o.Create()
    }
    var sets []string
    
    if o.IsPostAuthorDirty == true {
        sets = append(sets,fmt.Sprintf(`post_author = '%d'`,o.PostAuthor))
    }

    if o.IsPostDateDirty == true {
        sets = append(sets,fmt.Sprintf(`post_date = '%s'`,o.PostDate))
    }

    if o.IsPostDateGmtDirty == true {
        sets = append(sets,fmt.Sprintf(`post_date_gmt = '%s'`,o.PostDateGmt))
    }

    if o.IsPostContentDirty == true {
        sets = append(sets,fmt.Sprintf(`post_content = '%s'`,o._adapter.SafeString(o.PostContent)))
    }

    if o.IsPostTitleDirty == true {
        sets = append(sets,fmt.Sprintf(`post_title = '%s'`,o._adapter.SafeString(o.PostTitle)))
    }

    if o.IsPostExcerptDirty == true {
        sets = append(sets,fmt.Sprintf(`post_excerpt = '%s'`,o._adapter.SafeString(o.PostExcerpt)))
    }

    if o.IsPostStatusDirty == true {
        sets = append(sets,fmt.Sprintf(`post_status = '%s'`,o._adapter.SafeString(o.PostStatus)))
    }

    if o.IsCommentStatusDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_status = '%s'`,o._adapter.SafeString(o.CommentStatus)))
    }

    if o.IsPingStatusDirty == true {
        sets = append(sets,fmt.Sprintf(`ping_status = '%s'`,o._adapter.SafeString(o.PingStatus)))
    }

    if o.IsPostPasswordDirty == true {
        sets = append(sets,fmt.Sprintf(`post_password = '%s'`,o._adapter.SafeString(o.PostPassword)))
    }

    if o.IsPostNameDirty == true {
        sets = append(sets,fmt.Sprintf(`post_name = '%s'`,o._adapter.SafeString(o.PostName)))
    }

    if o.IsToPingDirty == true {
        sets = append(sets,fmt.Sprintf(`to_ping = '%s'`,o._adapter.SafeString(o.ToPing)))
    }

    if o.IsPingedDirty == true {
        sets = append(sets,fmt.Sprintf(`pinged = '%s'`,o._adapter.SafeString(o.Pinged)))
    }

    if o.IsPostModifiedDirty == true {
        sets = append(sets,fmt.Sprintf(`post_modified = '%s'`,o.PostModified))
    }

    if o.IsPostModifiedGmtDirty == true {
        sets = append(sets,fmt.Sprintf(`post_modified_gmt = '%s'`,o.PostModifiedGmt))
    }

    if o.IsPostContentFilteredDirty == true {
        sets = append(sets,fmt.Sprintf(`post_content_filtered = '%s'`,o._adapter.SafeString(o.PostContentFiltered)))
    }

    if o.IsPostParentDirty == true {
        sets = append(sets,fmt.Sprintf(`post_parent = '%d'`,o.PostParent))
    }

    if o.IsGuidDirty == true {
        sets = append(sets,fmt.Sprintf(`guid = '%s'`,o._adapter.SafeString(o.Guid)))
    }

    if o.IsMenuOrderDirty == true {
        sets = append(sets,fmt.Sprintf(`menu_order = '%d'`,o.MenuOrder))
    }

    if o.IsPostTypeDirty == true {
        sets = append(sets,fmt.Sprintf(`post_type = '%s'`,o._adapter.SafeString(o.PostType)))
    }

    if o.IsPostMimeTypeDirty == true {
        sets = append(sets,fmt.Sprintf(`post_mime_type = '%s'`,o._adapter.SafeString(o.PostMimeType)))
    }

    if o.IsCommentCountDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_count = '%d'`,o.CommentCount))
    }

    frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'",o._table,strings.Join(sets,`,`),o._pkey, o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return err
    }
    return nil
}
// Update is a dynamic updater, it considers whether or not
// a field is 'dirty' and needs to be updated. Will only work
// if you use the Getters and Setters
func (o *Post) Update() error {
    var sets []string
    
    if o.IsPostAuthorDirty == true {
        sets = append(sets,fmt.Sprintf(`post_author = '%d'`,o.PostAuthor))
    }

    if o.IsPostDateDirty == true {
        sets = append(sets,fmt.Sprintf(`post_date = '%s'`,o.PostDate))
    }

    if o.IsPostDateGmtDirty == true {
        sets = append(sets,fmt.Sprintf(`post_date_gmt = '%s'`,o.PostDateGmt))
    }

    if o.IsPostContentDirty == true {
        sets = append(sets,fmt.Sprintf(`post_content = '%s'`,o._adapter.SafeString(o.PostContent)))
    }

    if o.IsPostTitleDirty == true {
        sets = append(sets,fmt.Sprintf(`post_title = '%s'`,o._adapter.SafeString(o.PostTitle)))
    }

    if o.IsPostExcerptDirty == true {
        sets = append(sets,fmt.Sprintf(`post_excerpt = '%s'`,o._adapter.SafeString(o.PostExcerpt)))
    }

    if o.IsPostStatusDirty == true {
        sets = append(sets,fmt.Sprintf(`post_status = '%s'`,o._adapter.SafeString(o.PostStatus)))
    }

    if o.IsCommentStatusDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_status = '%s'`,o._adapter.SafeString(o.CommentStatus)))
    }

    if o.IsPingStatusDirty == true {
        sets = append(sets,fmt.Sprintf(`ping_status = '%s'`,o._adapter.SafeString(o.PingStatus)))
    }

    if o.IsPostPasswordDirty == true {
        sets = append(sets,fmt.Sprintf(`post_password = '%s'`,o._adapter.SafeString(o.PostPassword)))
    }

    if o.IsPostNameDirty == true {
        sets = append(sets,fmt.Sprintf(`post_name = '%s'`,o._adapter.SafeString(o.PostName)))
    }

    if o.IsToPingDirty == true {
        sets = append(sets,fmt.Sprintf(`to_ping = '%s'`,o._adapter.SafeString(o.ToPing)))
    }

    if o.IsPingedDirty == true {
        sets = append(sets,fmt.Sprintf(`pinged = '%s'`,o._adapter.SafeString(o.Pinged)))
    }

    if o.IsPostModifiedDirty == true {
        sets = append(sets,fmt.Sprintf(`post_modified = '%s'`,o.PostModified))
    }

    if o.IsPostModifiedGmtDirty == true {
        sets = append(sets,fmt.Sprintf(`post_modified_gmt = '%s'`,o.PostModifiedGmt))
    }

    if o.IsPostContentFilteredDirty == true {
        sets = append(sets,fmt.Sprintf(`post_content_filtered = '%s'`,o._adapter.SafeString(o.PostContentFiltered)))
    }

    if o.IsPostParentDirty == true {
        sets = append(sets,fmt.Sprintf(`post_parent = '%d'`,o.PostParent))
    }

    if o.IsGuidDirty == true {
        sets = append(sets,fmt.Sprintf(`guid = '%s'`,o._adapter.SafeString(o.Guid)))
    }

    if o.IsMenuOrderDirty == true {
        sets = append(sets,fmt.Sprintf(`menu_order = '%d'`,o.MenuOrder))
    }

    if o.IsPostTypeDirty == true {
        sets = append(sets,fmt.Sprintf(`post_type = '%s'`,o._adapter.SafeString(o.PostType)))
    }

    if o.IsPostMimeTypeDirty == true {
        sets = append(sets,fmt.Sprintf(`post_mime_type = '%s'`,o._adapter.SafeString(o.PostMimeType)))
    }

    if o.IsCommentCountDirty == true {
        sets = append(sets,fmt.Sprintf(`comment_count = '%d'`,o.CommentCount))
    }

    frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'",o._table,strings.Join(sets,`,`),o._pkey, o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return err
    }
    return nil
}
// Create inserts the model. Calling Save will call this function
// automatically for new models
func (o *Post) Create() error {
    frmt := fmt.Sprintf("INSERT INTO %s (`post_author`, `post_date`, `post_date_gmt`, `post_content`, `post_title`, `post_excerpt`, `post_status`, `comment_status`, `ping_status`, `post_password`, `post_name`, `to_ping`, `pinged`, `post_modified`, `post_modified_gmt`, `post_content_filtered`, `post_parent`, `guid`, `menu_order`, `post_type`, `post_mime_type`, `comment_count`) VALUES ('%d', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%d', '%s', '%d', '%s', '%s', '%d')",o._table,o.PostAuthor, o.PostDate.ToString(), o.PostDateGmt.ToString(), o.PostContent, o.PostTitle, o.PostExcerpt, o.PostStatus, o.CommentStatus, o.PingStatus, o.PostPassword, o.PostName, o.ToPing, o.Pinged, o.PostModified.ToString(), o.PostModifiedGmt.ToString(), o.PostContentFiltered, o.PostParent, o.Guid, o.MenuOrder, o.PostType, o.PostMimeType, o.CommentCount)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return o._adapter.Oops(fmt.Sprintf(`%s led to %s`,frmt,err))
    }
    o.ID = o._adapter.LastInsertedId()
    o._new = false
    return nil
}


// UpdatePostAuthor an immediate DB Query to update a single column, in this
// case post_author
func (o *Post) UpdatePostAuthor(_updPostAuthor int64) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `post_author` = '%d' WHERE `ID` = '%d'",o._table,_updPostAuthor,o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.PostAuthor = _updPostAuthor
    return o._adapter.AffectedRows(),nil
}

// UpdatePostDate an immediate DB Query to update a single column, in this
// case post_date
func (o *Post) UpdatePostDate(_updPostDate *DateTime) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `post_date` = '%s' WHERE `ID` = '%d'",o._table,_updPostDate,o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.PostDate = _updPostDate
    return o._adapter.AffectedRows(),nil
}

// UpdatePostDateGmt an immediate DB Query to update a single column, in this
// case post_date_gmt
func (o *Post) UpdatePostDateGmt(_updPostDateGmt *DateTime) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `post_date_gmt` = '%s' WHERE `ID` = '%d'",o._table,_updPostDateGmt,o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.PostDateGmt = _updPostDateGmt
    return o._adapter.AffectedRows(),nil
}

// UpdatePostContent an immediate DB Query to update a single column, in this
// case post_content
func (o *Post) UpdatePostContent(_updPostContent string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `post_content` = '%s' WHERE `ID` = '%d'",o._table,_updPostContent,o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.PostContent = _updPostContent
    return o._adapter.AffectedRows(),nil
}

// UpdatePostTitle an immediate DB Query to update a single column, in this
// case post_title
func (o *Post) UpdatePostTitle(_updPostTitle string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `post_title` = '%s' WHERE `ID` = '%d'",o._table,_updPostTitle,o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.PostTitle = _updPostTitle
    return o._adapter.AffectedRows(),nil
}

// UpdatePostExcerpt an immediate DB Query to update a single column, in this
// case post_excerpt
func (o *Post) UpdatePostExcerpt(_updPostExcerpt string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `post_excerpt` = '%s' WHERE `ID` = '%d'",o._table,_updPostExcerpt,o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.PostExcerpt = _updPostExcerpt
    return o._adapter.AffectedRows(),nil
}

// UpdatePostStatus an immediate DB Query to update a single column, in this
// case post_status
func (o *Post) UpdatePostStatus(_updPostStatus string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `post_status` = '%s' WHERE `ID` = '%d'",o._table,_updPostStatus,o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.PostStatus = _updPostStatus
    return o._adapter.AffectedRows(),nil
}

// UpdateCommentStatus an immediate DB Query to update a single column, in this
// case comment_status
func (o *Post) UpdateCommentStatus(_updCommentStatus string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `comment_status` = '%s' WHERE `ID` = '%d'",o._table,_updCommentStatus,o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.CommentStatus = _updCommentStatus
    return o._adapter.AffectedRows(),nil
}

// UpdatePingStatus an immediate DB Query to update a single column, in this
// case ping_status
func (o *Post) UpdatePingStatus(_updPingStatus string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `ping_status` = '%s' WHERE `ID` = '%d'",o._table,_updPingStatus,o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.PingStatus = _updPingStatus
    return o._adapter.AffectedRows(),nil
}

// UpdatePostPassword an immediate DB Query to update a single column, in this
// case post_password
func (o *Post) UpdatePostPassword(_updPostPassword string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `post_password` = '%s' WHERE `ID` = '%d'",o._table,_updPostPassword,o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.PostPassword = _updPostPassword
    return o._adapter.AffectedRows(),nil
}

// UpdatePostName an immediate DB Query to update a single column, in this
// case post_name
func (o *Post) UpdatePostName(_updPostName string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `post_name` = '%s' WHERE `ID` = '%d'",o._table,_updPostName,o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.PostName = _updPostName
    return o._adapter.AffectedRows(),nil
}

// UpdateToPing an immediate DB Query to update a single column, in this
// case to_ping
func (o *Post) UpdateToPing(_updToPing string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `to_ping` = '%s' WHERE `ID` = '%d'",o._table,_updToPing,o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.ToPing = _updToPing
    return o._adapter.AffectedRows(),nil
}

// UpdatePinged an immediate DB Query to update a single column, in this
// case pinged
func (o *Post) UpdatePinged(_updPinged string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `pinged` = '%s' WHERE `ID` = '%d'",o._table,_updPinged,o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.Pinged = _updPinged
    return o._adapter.AffectedRows(),nil
}

// UpdatePostModified an immediate DB Query to update a single column, in this
// case post_modified
func (o *Post) UpdatePostModified(_updPostModified *DateTime) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `post_modified` = '%s' WHERE `ID` = '%d'",o._table,_updPostModified,o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.PostModified = _updPostModified
    return o._adapter.AffectedRows(),nil
}

// UpdatePostModifiedGmt an immediate DB Query to update a single column, in this
// case post_modified_gmt
func (o *Post) UpdatePostModifiedGmt(_updPostModifiedGmt *DateTime) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `post_modified_gmt` = '%s' WHERE `ID` = '%d'",o._table,_updPostModifiedGmt,o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.PostModifiedGmt = _updPostModifiedGmt
    return o._adapter.AffectedRows(),nil
}

// UpdatePostContentFiltered an immediate DB Query to update a single column, in this
// case post_content_filtered
func (o *Post) UpdatePostContentFiltered(_updPostContentFiltered string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `post_content_filtered` = '%s' WHERE `ID` = '%d'",o._table,_updPostContentFiltered,o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.PostContentFiltered = _updPostContentFiltered
    return o._adapter.AffectedRows(),nil
}

// UpdatePostParent an immediate DB Query to update a single column, in this
// case post_parent
func (o *Post) UpdatePostParent(_updPostParent int64) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `post_parent` = '%d' WHERE `ID` = '%d'",o._table,_updPostParent,o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.PostParent = _updPostParent
    return o._adapter.AffectedRows(),nil
}

// UpdateGuid an immediate DB Query to update a single column, in this
// case guid
func (o *Post) UpdateGuid(_updGuid string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `guid` = '%s' WHERE `ID` = '%d'",o._table,_updGuid,o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.Guid = _updGuid
    return o._adapter.AffectedRows(),nil
}

// UpdateMenuOrder an immediate DB Query to update a single column, in this
// case menu_order
func (o *Post) UpdateMenuOrder(_updMenuOrder int) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `menu_order` = '%d' WHERE `ID` = '%d'",o._table,_updMenuOrder,o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.MenuOrder = _updMenuOrder
    return o._adapter.AffectedRows(),nil
}

// UpdatePostType an immediate DB Query to update a single column, in this
// case post_type
func (o *Post) UpdatePostType(_updPostType string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `post_type` = '%s' WHERE `ID` = '%d'",o._table,_updPostType,o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.PostType = _updPostType
    return o._adapter.AffectedRows(),nil
}

// UpdatePostMimeType an immediate DB Query to update a single column, in this
// case post_mime_type
func (o *Post) UpdatePostMimeType(_updPostMimeType string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `post_mime_type` = '%s' WHERE `ID` = '%d'",o._table,_updPostMimeType,o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.PostMimeType = _updPostMimeType
    return o._adapter.AffectedRows(),nil
}

// UpdateCommentCount an immediate DB Query to update a single column, in this
// case comment_count
func (o *Post) UpdateCommentCount(_updCommentCount int64) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `comment_count` = '%d' WHERE `ID` = '%d'",o._table,_updCommentCount,o.ID)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.CommentCount = _updCommentCount
    return o._adapter.AffectedRows(),nil
}

// TermRelationship is a Object Relational Mapping to
// the database table that represents it. In this case it is
// term_relationships. The table name will be Sprintf'd to include
// the prefix you define in your YAML configuration for the
// Adapter.
type TermRelationship struct {
    _table string
    _adapter Adapter
    _pkey string // 0 The name of the primary key in this table
    _conds []string
    _new bool
    ObjectId int64
    TermTaxonomyId int64
    TermOrder int
	// Dirty markers for smart updates
    IsObjectIdDirty bool
    IsTermTaxonomyIdDirty bool
    IsTermOrderDirty bool
	// Relationships
}

// NewTermRelationship binds an Adapter to a new instance
// of TermRelationship and sets up the _table and primary keys
func NewTermRelationship(a Adapter) *TermRelationship {
    var o TermRelationship
    o._table = fmt.Sprintf("%sterm_relationships",a.DatabasePrefix())
    o._adapter = a
    o._pkey = "term_taxonomy_id"
    o._new = false
    return &o
}


// GetPrimaryKeyValue returns the value, usually int64 of
// the PrimaryKey
func (o *TermRelationship) GetPrimaryKeyValue() int64 {
    return o.TermTaxonomyId
}
// GetPrimaryKeyName returns the DB field name
func (o *TermRelationship) GetPrimaryKeyName() string {
    return `term_taxonomy_id`
}

// GetObjectId returns the value of 
// TermRelationship.ObjectId
func (o *TermRelationship) GetObjectId() int64 {
    return o.ObjectId
}
// SetObjectId sets and marks as dirty the value of
// TermRelationship.ObjectId
func (o *TermRelationship) SetObjectId(arg int64) {
    o.ObjectId = arg
    o.IsObjectIdDirty = true
}

// GetTermTaxonomyId returns the value of 
// TermRelationship.TermTaxonomyId
func (o *TermRelationship) GetTermTaxonomyId() int64 {
    return o.TermTaxonomyId
}
// SetTermTaxonomyId sets and marks as dirty the value of
// TermRelationship.TermTaxonomyId
func (o *TermRelationship) SetTermTaxonomyId(arg int64) {
    o.TermTaxonomyId = arg
    o.IsTermTaxonomyIdDirty = true
}

// GetTermOrder returns the value of 
// TermRelationship.TermOrder
func (o *TermRelationship) GetTermOrder() int {
    return o.TermOrder
}
// SetTermOrder sets and marks as dirty the value of
// TermRelationship.TermOrder
func (o *TermRelationship) SetTermOrder(arg int) {
    o.TermOrder = arg
    o.IsTermOrderDirty = true
}

// FindByObjectId dynamic finder for object_id -> []*TermRelationship,error
// Generic and programatically generator finder for TermRelationship
//```go  
//    m := NewTermRelationship(a)
//    results,err := m.FindByObjectId(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of TermRelationship
//    }
//```  
        func (o *TermRelationship) FindByObjectId(_findByObjectId int64) ([]*TermRelationship,error) {

    var _modelSlice []*TermRelationship
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "object_id", _findByObjectId)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := TermRelationship{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// Find for the TermRelationship is a bit tricky, as it has no
// primary key as such, but a composite key.
func (o *TermRelationship) Find(termId int64,objectId int64) (bool,error) {

    var _modelSlice []*TermRelationship
    q := fmt.Sprintf("SELECT * FROM %s WHERE `term_taxonomy_id` = '%d' AND `object_id` = '%d'",o._table, termId,objectId)
    results, err := o._adapter.Query(q)
    if err != nil {
        return false,err
    }
    
    for _,result := range results {
        ro := TermRelationship{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return false,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return false, errors.New("not found")
    }
    o.FromTermRelationship(_modelSlice[0])
    return true,nil

}
// FindByTermOrder dynamic finder for term_order -> []*TermRelationship,error
// Generic and programatically generator finder for TermRelationship
//```go  
//    m := NewTermRelationship(a)
//    results,err := m.FindByTermOrder(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of TermRelationship
//    }
//```  
        func (o *TermRelationship) FindByTermOrder(_findByTermOrder int) ([]*TermRelationship,error) {

    var _modelSlice []*TermRelationship
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "term_order", _findByTermOrder)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := TermRelationship{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}

// FromDBValueMap Converts a DBValueMap returned from Adapter.Query to a TermRelationship
func (o *TermRelationship) FromDBValueMap(m map[string]DBValue) error {
	_ObjectId,err := m["object_id"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.ObjectId = _ObjectId
	_TermTaxonomyId,err := m["term_taxonomy_id"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.TermTaxonomyId = _TermTaxonomyId
	_TermOrder,err := m["term_order"].AsInt()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.TermOrder = _TermOrder

 	return nil
}
// FromTermRelationship A kind of Clone function for TermRelationship
func (o *TermRelationship) FromTermRelationship(m *TermRelationship) {
	o.ObjectId = m.ObjectId
	o.TermTaxonomyId = m.TermTaxonomyId
	o.TermOrder = m.TermOrder

}
// Reload A function to forcibly reload TermRelationship
func (o *TermRelationship) Reload() error {
    _,err := o.Find(o.TermTaxonomyId ,o.ObjectId)
    return err
}

// Save is a dynamic saver 'inherited' by all models
func (o *TermRelationship) Save() error {
    if o._new == true {
        return o.Create()
    }
    var sets []string
    
    if o.IsObjectIdDirty == true {
        sets = append(sets,fmt.Sprintf(`object_id = '%d'`,o.ObjectId))
    }

    if o.IsTermTaxonomyIdDirty == true {
        sets = append(sets,fmt.Sprintf(`term_taxonomy_id = '%d'`,o.TermTaxonomyId))
    }

    if o.IsTermOrderDirty == true {
        sets = append(sets,fmt.Sprintf(`term_order = '%d'`,o.TermOrder))
    }

    frmt := fmt.Sprintf("UPDATE %s SET %s WHERE `term_taxonomy_id` = '%d' AND object_id = '%d'",o._table,strings.Join(sets,`,`),o.TermTaxonomyId, o.ObjectId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return err
    }
    return nil
}
// Update is a dynamic updater, it considers whether or not
// a field is 'dirty' and needs to be updated. Will only work
// if you use the Getters and Setters
func (o *TermRelationship) Update() error {
    var sets []string
    
    if o.IsObjectIdDirty == true {
        sets = append(sets,fmt.Sprintf(`object_id = '%d'`,o.ObjectId))
    }

    if o.IsTermTaxonomyIdDirty == true {
        sets = append(sets,fmt.Sprintf(`term_taxonomy_id = '%d'`,o.TermTaxonomyId))
    }

    if o.IsTermOrderDirty == true {
        sets = append(sets,fmt.Sprintf(`term_order = '%d'`,o.TermOrder))
    }

    frmt := fmt.Sprintf("UPDATE %s SET %s WHERE `term_taxonomy_id` = '%d' AND object_id = '%d'",o._table,strings.Join(sets,`,`),o.TermTaxonomyId, o.ObjectId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return err
    }
    return nil
}
// Create inserts the model. Calling Save will call this function
// automatically for new models
func (o *TermRelationship) Create() error {
    frmt := fmt.Sprintf("INSERT INTO %s (`object_id`, `term_taxonomy_id`, `term_order`) VALUES ('%d', '%d', '%d')",o._table,o.ObjectId, o.TermTaxonomyId, o.TermOrder)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return o._adapter.Oops(fmt.Sprintf(`%s led to %s`,frmt,err))
    }
    
    o._new = false
    return nil
}


// UpdateObjectId an immediate DB Query to update a single column, in this
// case object_id
func (o *TermRelationship) UpdateObjectId(_updObjectId int64) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `object_id` = '%d' WHERE term_taxonomy_id = '%d' AND object_id = '%d'",o._table,_updObjectId,o.TermTaxonomyId,o.ObjectId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.ObjectId = _updObjectId
    return o._adapter.AffectedRows(),nil
}

// UpdateTermOrder an immediate DB Query to update a single column, in this
// case term_order
func (o *TermRelationship) UpdateTermOrder(_updTermOrder int) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `term_order` = '%d' WHERE term_taxonomy_id = '%d' AND object_id = '%d'",o._table,_updTermOrder,o.TermTaxonomyId,o.ObjectId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.TermOrder = _updTermOrder
    return o._adapter.AffectedRows(),nil
}

// TermTaxonomy is a Object Relational Mapping to
// the database table that represents it. In this case it is
// term_taxonomy. The table name will be Sprintf'd to include
// the prefix you define in your YAML configuration for the
// Adapter.
type TermTaxonomy struct {
    _table string
    _adapter Adapter
    _pkey string // 0 The name of the primary key in this table
    _conds []string
    _new bool
    TermTaxonomyId int64
    TermId int64
    Taxonomy string
    Description string
    Parent int64
    Count int64
	// Dirty markers for smart updates
    IsTermTaxonomyIdDirty bool
    IsTermIdDirty bool
    IsTaxonomyDirty bool
    IsDescriptionDirty bool
    IsParentDirty bool
    IsCountDirty bool
	// Relationships
}

// NewTermTaxonomy binds an Adapter to a new instance
// of TermTaxonomy and sets up the _table and primary keys
func NewTermTaxonomy(a Adapter) *TermTaxonomy {
    var o TermTaxonomy
    o._table = fmt.Sprintf("%sterm_taxonomy",a.DatabasePrefix())
    o._adapter = a
    o._pkey = "term_taxonomy_id"
    o._new = false
    return &o
}


// GetPrimaryKeyValue returns the value, usually int64 of
// the PrimaryKey
func (o *TermTaxonomy) GetPrimaryKeyValue() int64 {
    return o.TermTaxonomyId
}
// GetPrimaryKeyName returns the DB field name
func (o *TermTaxonomy) GetPrimaryKeyName() string {
    return `term_taxonomy_id`
}

// GetTermTaxonomyId returns the value of 
// TermTaxonomy.TermTaxonomyId
func (o *TermTaxonomy) GetTermTaxonomyId() int64 {
    return o.TermTaxonomyId
}
// SetTermTaxonomyId sets and marks as dirty the value of
// TermTaxonomy.TermTaxonomyId
func (o *TermTaxonomy) SetTermTaxonomyId(arg int64) {
    o.TermTaxonomyId = arg
    o.IsTermTaxonomyIdDirty = true
}

// GetTermId returns the value of 
// TermTaxonomy.TermId
func (o *TermTaxonomy) GetTermId() int64 {
    return o.TermId
}
// SetTermId sets and marks as dirty the value of
// TermTaxonomy.TermId
func (o *TermTaxonomy) SetTermId(arg int64) {
    o.TermId = arg
    o.IsTermIdDirty = true
}

// GetTaxonomy returns the value of 
// TermTaxonomy.Taxonomy
func (o *TermTaxonomy) GetTaxonomy() string {
    return o.Taxonomy
}
// SetTaxonomy sets and marks as dirty the value of
// TermTaxonomy.Taxonomy
func (o *TermTaxonomy) SetTaxonomy(arg string) {
    o.Taxonomy = arg
    o.IsTaxonomyDirty = true
}

// GetDescription returns the value of 
// TermTaxonomy.Description
func (o *TermTaxonomy) GetDescription() string {
    return o.Description
}
// SetDescription sets and marks as dirty the value of
// TermTaxonomy.Description
func (o *TermTaxonomy) SetDescription(arg string) {
    o.Description = arg
    o.IsDescriptionDirty = true
}

// GetParent returns the value of 
// TermTaxonomy.Parent
func (o *TermTaxonomy) GetParent() int64 {
    return o.Parent
}
// SetParent sets and marks as dirty the value of
// TermTaxonomy.Parent
func (o *TermTaxonomy) SetParent(arg int64) {
    o.Parent = arg
    o.IsParentDirty = true
}

// GetCount returns the value of 
// TermTaxonomy.Count
func (o *TermTaxonomy) GetCount() int64 {
    return o.Count
}
// SetCount sets and marks as dirty the value of
// TermTaxonomy.Count
func (o *TermTaxonomy) SetCount(arg int64) {
    o.Count = arg
    o.IsCountDirty = true
}

// Find dynamic finder for term_taxonomy_id -> bool,error
// Generic and programatically generator finder for TermTaxonomy
// Note that Fine returns a bool if found, not err, in the case of
// a return of true, the instance data will be filled out.
// a call to find ALWAYS overwrites the model you call Find on
// i.e. receiver is a pointer. 
//```go
//      m := NewTermTaxonomy(a)
//      found,err := m.Find(23)
//      .. handle err
//      if found == false {
//          // handle found
//      }
//      ... do what you want with m here
//```
        func (o *TermTaxonomy) Find(_findByTermTaxonomyId int64) (bool,error) {

    var _modelSlice []*TermTaxonomy
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "term_taxonomy_id", _findByTermTaxonomyId)
    results, err := o._adapter.Query(q)
    if err != nil {
        return false,o._adapter.Oops(fmt.Sprintf(`%s`,err))
    }
    
    for _,result := range results {
        ro := TermTaxonomy{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return false,o._adapter.Oops(fmt.Sprintf(`%s`,err))
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return false, o._adapter.Oops(`not found`)
    }
    o.FromTermTaxonomy(_modelSlice[0])
    return true,nil

}
// FindByTermId dynamic finder for term_id -> []*TermTaxonomy,error
// Generic and programatically generator finder for TermTaxonomy
//```go  
//    m := NewTermTaxonomy(a)
//    results,err := m.FindByTermId(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of TermTaxonomy
//    }
//```  
        func (o *TermTaxonomy) FindByTermId(_findByTermId int64) ([]*TermTaxonomy,error) {

    var _modelSlice []*TermTaxonomy
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "term_id", _findByTermId)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := TermTaxonomy{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByTaxonomy dynamic finder for taxonomy -> []*TermTaxonomy,error
// Generic and programatically generator finder for TermTaxonomy
//```go  
//    m := NewTermTaxonomy(a)
//    results,err := m.FindByTaxonomy(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of TermTaxonomy
//    }
//```  
        func (o *TermTaxonomy) FindByTaxonomy(_findByTaxonomy string) ([]*TermTaxonomy,error) {

    var _modelSlice []*TermTaxonomy
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "taxonomy", _findByTaxonomy)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := TermTaxonomy{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByDescription dynamic finder for description -> []*TermTaxonomy,error
// Generic and programatically generator finder for TermTaxonomy
//```go  
//    m := NewTermTaxonomy(a)
//    results,err := m.FindByDescription(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of TermTaxonomy
//    }
//```  
        func (o *TermTaxonomy) FindByDescription(_findByDescription string) ([]*TermTaxonomy,error) {

    var _modelSlice []*TermTaxonomy
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "description", _findByDescription)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := TermTaxonomy{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByParent dynamic finder for parent -> []*TermTaxonomy,error
// Generic and programatically generator finder for TermTaxonomy
//```go  
//    m := NewTermTaxonomy(a)
//    results,err := m.FindByParent(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of TermTaxonomy
//    }
//```  
        func (o *TermTaxonomy) FindByParent(_findByParent int64) ([]*TermTaxonomy,error) {

    var _modelSlice []*TermTaxonomy
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "parent", _findByParent)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := TermTaxonomy{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByCount dynamic finder for count -> []*TermTaxonomy,error
// Generic and programatically generator finder for TermTaxonomy
//```go  
//    m := NewTermTaxonomy(a)
//    results,err := m.FindByCount(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of TermTaxonomy
//    }
//```  
        func (o *TermTaxonomy) FindByCount(_findByCount int64) ([]*TermTaxonomy,error) {

    var _modelSlice []*TermTaxonomy
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "count", _findByCount)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := TermTaxonomy{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}

// FromDBValueMap Converts a DBValueMap returned from Adapter.Query to a TermTaxonomy
func (o *TermTaxonomy) FromDBValueMap(m map[string]DBValue) error {
	_TermTaxonomyId,err := m["term_taxonomy_id"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.TermTaxonomyId = _TermTaxonomyId
	_TermId,err := m["term_id"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.TermId = _TermId
	_Taxonomy,err := m["taxonomy"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.Taxonomy = _Taxonomy
	_Description,err := m["description"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.Description = _Description
	_Parent,err := m["parent"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.Parent = _Parent
	_Count,err := m["count"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.Count = _Count

 	return nil
}
// FromTermTaxonomy A kind of Clone function for TermTaxonomy
func (o *TermTaxonomy) FromTermTaxonomy(m *TermTaxonomy) {
	o.TermTaxonomyId = m.TermTaxonomyId
	o.TermId = m.TermId
	o.Taxonomy = m.Taxonomy
	o.Description = m.Description
	o.Parent = m.Parent
	o.Count = m.Count

}
// Reload A function to forcibly reload TermTaxonomy
func (o *TermTaxonomy) Reload() error {
    _,err := o.Find(o.GetPrimaryKeyValue())
    return err
}

// Save is a dynamic saver 'inherited' by all models
func (o *TermTaxonomy) Save() error {
    if o._new == true {
        return o.Create()
    }
    var sets []string
    
    if o.IsTermIdDirty == true {
        sets = append(sets,fmt.Sprintf(`term_id = '%d'`,o.TermId))
    }

    if o.IsTaxonomyDirty == true {
        sets = append(sets,fmt.Sprintf(`taxonomy = '%s'`,o._adapter.SafeString(o.Taxonomy)))
    }

    if o.IsDescriptionDirty == true {
        sets = append(sets,fmt.Sprintf(`description = '%s'`,o._adapter.SafeString(o.Description)))
    }

    if o.IsParentDirty == true {
        sets = append(sets,fmt.Sprintf(`parent = '%d'`,o.Parent))
    }

    if o.IsCountDirty == true {
        sets = append(sets,fmt.Sprintf(`count = '%d'`,o.Count))
    }

    frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'",o._table,strings.Join(sets,`,`),o._pkey, o.TermTaxonomyId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return err
    }
    return nil
}
// Update is a dynamic updater, it considers whether or not
// a field is 'dirty' and needs to be updated. Will only work
// if you use the Getters and Setters
func (o *TermTaxonomy) Update() error {
    var sets []string
    
    if o.IsTermIdDirty == true {
        sets = append(sets,fmt.Sprintf(`term_id = '%d'`,o.TermId))
    }

    if o.IsTaxonomyDirty == true {
        sets = append(sets,fmt.Sprintf(`taxonomy = '%s'`,o._adapter.SafeString(o.Taxonomy)))
    }

    if o.IsDescriptionDirty == true {
        sets = append(sets,fmt.Sprintf(`description = '%s'`,o._adapter.SafeString(o.Description)))
    }

    if o.IsParentDirty == true {
        sets = append(sets,fmt.Sprintf(`parent = '%d'`,o.Parent))
    }

    if o.IsCountDirty == true {
        sets = append(sets,fmt.Sprintf(`count = '%d'`,o.Count))
    }

    frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'",o._table,strings.Join(sets,`,`),o._pkey, o.TermTaxonomyId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return err
    }
    return nil
}
// Create inserts the model. Calling Save will call this function
// automatically for new models
func (o *TermTaxonomy) Create() error {
    frmt := fmt.Sprintf("INSERT INTO %s (`term_id`, `taxonomy`, `description`, `parent`, `count`) VALUES ('%d', '%s', '%s', '%d', '%d')",o._table,o.TermId, o.Taxonomy, o.Description, o.Parent, o.Count)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return o._adapter.Oops(fmt.Sprintf(`%s led to %s`,frmt,err))
    }
    o.TermTaxonomyId = o._adapter.LastInsertedId()
    o._new = false
    return nil
}


// UpdateTermId an immediate DB Query to update a single column, in this
// case term_id
func (o *TermTaxonomy) UpdateTermId(_updTermId int64) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `term_id` = '%d' WHERE `term_taxonomy_id` = '%d'",o._table,_updTermId,o.TermTaxonomyId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.TermId = _updTermId
    return o._adapter.AffectedRows(),nil
}

// UpdateTaxonomy an immediate DB Query to update a single column, in this
// case taxonomy
func (o *TermTaxonomy) UpdateTaxonomy(_updTaxonomy string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `taxonomy` = '%s' WHERE `term_taxonomy_id` = '%d'",o._table,_updTaxonomy,o.TermTaxonomyId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.Taxonomy = _updTaxonomy
    return o._adapter.AffectedRows(),nil
}

// UpdateDescription an immediate DB Query to update a single column, in this
// case description
func (o *TermTaxonomy) UpdateDescription(_updDescription string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `description` = '%s' WHERE `term_taxonomy_id` = '%d'",o._table,_updDescription,o.TermTaxonomyId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.Description = _updDescription
    return o._adapter.AffectedRows(),nil
}

// UpdateParent an immediate DB Query to update a single column, in this
// case parent
func (o *TermTaxonomy) UpdateParent(_updParent int64) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `parent` = '%d' WHERE `term_taxonomy_id` = '%d'",o._table,_updParent,o.TermTaxonomyId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.Parent = _updParent
    return o._adapter.AffectedRows(),nil
}

// UpdateCount an immediate DB Query to update a single column, in this
// case count
func (o *TermTaxonomy) UpdateCount(_updCount int64) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `count` = '%d' WHERE `term_taxonomy_id` = '%d'",o._table,_updCount,o.TermTaxonomyId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.Count = _updCount
    return o._adapter.AffectedRows(),nil
}

// Term is a Object Relational Mapping to
// the database table that represents it. In this case it is
// terms. The table name will be Sprintf'd to include
// the prefix you define in your YAML configuration for the
// Adapter.
type Term struct {
    _table string
    _adapter Adapter
    _pkey string // 0 The name of the primary key in this table
    _conds []string
    _new bool
    TermId int64
    Name string
    Slug string
    TermGroup int64
	// Dirty markers for smart updates
    IsTermIdDirty bool
    IsNameDirty bool
    IsSlugDirty bool
    IsTermGroupDirty bool
	// Relationships
}

// NewTerm binds an Adapter to a new instance
// of Term and sets up the _table and primary keys
func NewTerm(a Adapter) *Term {
    var o Term
    o._table = fmt.Sprintf("%sterms",a.DatabasePrefix())
    o._adapter = a
    o._pkey = "term_id"
    o._new = false
    return &o
}


// GetPrimaryKeyValue returns the value, usually int64 of
// the PrimaryKey
func (o *Term) GetPrimaryKeyValue() int64 {
    return o.TermId
}
// GetPrimaryKeyName returns the DB field name
func (o *Term) GetPrimaryKeyName() string {
    return `term_id`
}

// GetTermId returns the value of 
// Term.TermId
func (o *Term) GetTermId() int64 {
    return o.TermId
}
// SetTermId sets and marks as dirty the value of
// Term.TermId
func (o *Term) SetTermId(arg int64) {
    o.TermId = arg
    o.IsTermIdDirty = true
}

// GetName returns the value of 
// Term.Name
func (o *Term) GetName() string {
    return o.Name
}
// SetName sets and marks as dirty the value of
// Term.Name
func (o *Term) SetName(arg string) {
    o.Name = arg
    o.IsNameDirty = true
}

// GetSlug returns the value of 
// Term.Slug
func (o *Term) GetSlug() string {
    return o.Slug
}
// SetSlug sets and marks as dirty the value of
// Term.Slug
func (o *Term) SetSlug(arg string) {
    o.Slug = arg
    o.IsSlugDirty = true
}

// GetTermGroup returns the value of 
// Term.TermGroup
func (o *Term) GetTermGroup() int64 {
    return o.TermGroup
}
// SetTermGroup sets and marks as dirty the value of
// Term.TermGroup
func (o *Term) SetTermGroup(arg int64) {
    o.TermGroup = arg
    o.IsTermGroupDirty = true
}

// Find dynamic finder for term_id -> bool,error
// Generic and programatically generator finder for Term
// Note that Fine returns a bool if found, not err, in the case of
// a return of true, the instance data will be filled out.
// a call to find ALWAYS overwrites the model you call Find on
// i.e. receiver is a pointer. 
//```go
//      m := NewTerm(a)
//      found,err := m.Find(23)
//      .. handle err
//      if found == false {
//          // handle found
//      }
//      ... do what you want with m here
//```
        func (o *Term) Find(_findByTermId int64) (bool,error) {

    var _modelSlice []*Term
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "term_id", _findByTermId)
    results, err := o._adapter.Query(q)
    if err != nil {
        return false,o._adapter.Oops(fmt.Sprintf(`%s`,err))
    }
    
    for _,result := range results {
        ro := Term{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return false,o._adapter.Oops(fmt.Sprintf(`%s`,err))
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return false, o._adapter.Oops(`not found`)
    }
    o.FromTerm(_modelSlice[0])
    return true,nil

}
// FindByName dynamic finder for name -> []*Term,error
// Generic and programatically generator finder for Term
//```go  
//    m := NewTerm(a)
//    results,err := m.FindByName(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Term
//    }
//```  
        func (o *Term) FindByName(_findByName string) ([]*Term,error) {

    var _modelSlice []*Term
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "name", _findByName)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Term{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindBySlug dynamic finder for slug -> []*Term,error
// Generic and programatically generator finder for Term
//```go  
//    m := NewTerm(a)
//    results,err := m.FindBySlug(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Term
//    }
//```  
        func (o *Term) FindBySlug(_findBySlug string) ([]*Term,error) {

    var _modelSlice []*Term
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "slug", _findBySlug)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Term{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByTermGroup dynamic finder for term_group -> []*Term,error
// Generic and programatically generator finder for Term
//```go  
//    m := NewTerm(a)
//    results,err := m.FindByTermGroup(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of Term
//    }
//```  
        func (o *Term) FindByTermGroup(_findByTermGroup int64) ([]*Term,error) {

    var _modelSlice []*Term
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "term_group", _findByTermGroup)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := Term{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}

// FromDBValueMap Converts a DBValueMap returned from Adapter.Query to a Term
func (o *Term) FromDBValueMap(m map[string]DBValue) error {
	_TermId,err := m["term_id"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.TermId = _TermId
	_Name,err := m["name"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.Name = _Name
	_Slug,err := m["slug"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.Slug = _Slug
	_TermGroup,err := m["term_group"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.TermGroup = _TermGroup

 	return nil
}
// FromTerm A kind of Clone function for Term
func (o *Term) FromTerm(m *Term) {
	o.TermId = m.TermId
	o.Name = m.Name
	o.Slug = m.Slug
	o.TermGroup = m.TermGroup

}
// Reload A function to forcibly reload Term
func (o *Term) Reload() error {
    _,err := o.Find(o.GetPrimaryKeyValue())
    return err
}

// Save is a dynamic saver 'inherited' by all models
func (o *Term) Save() error {
    if o._new == true {
        return o.Create()
    }
    var sets []string
    
    if o.IsNameDirty == true {
        sets = append(sets,fmt.Sprintf(`name = '%s'`,o._adapter.SafeString(o.Name)))
    }

    if o.IsSlugDirty == true {
        sets = append(sets,fmt.Sprintf(`slug = '%s'`,o._adapter.SafeString(o.Slug)))
    }

    if o.IsTermGroupDirty == true {
        sets = append(sets,fmt.Sprintf(`term_group = '%d'`,o.TermGroup))
    }

    frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'",o._table,strings.Join(sets,`,`),o._pkey, o.TermId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return err
    }
    return nil
}
// Update is a dynamic updater, it considers whether or not
// a field is 'dirty' and needs to be updated. Will only work
// if you use the Getters and Setters
func (o *Term) Update() error {
    var sets []string
    
    if o.IsNameDirty == true {
        sets = append(sets,fmt.Sprintf(`name = '%s'`,o._adapter.SafeString(o.Name)))
    }

    if o.IsSlugDirty == true {
        sets = append(sets,fmt.Sprintf(`slug = '%s'`,o._adapter.SafeString(o.Slug)))
    }

    if o.IsTermGroupDirty == true {
        sets = append(sets,fmt.Sprintf(`term_group = '%d'`,o.TermGroup))
    }

    frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'",o._table,strings.Join(sets,`,`),o._pkey, o.TermId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return err
    }
    return nil
}
// Create inserts the model. Calling Save will call this function
// automatically for new models
func (o *Term) Create() error {
    frmt := fmt.Sprintf("INSERT INTO %s (`name`, `slug`, `term_group`) VALUES ('%s', '%s', '%d')",o._table,o.Name, o.Slug, o.TermGroup)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return o._adapter.Oops(fmt.Sprintf(`%s led to %s`,frmt,err))
    }
    o.TermId = o._adapter.LastInsertedId()
    o._new = false
    return nil
}


// UpdateName an immediate DB Query to update a single column, in this
// case name
func (o *Term) UpdateName(_updName string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `name` = '%s' WHERE `term_id` = '%d'",o._table,_updName,o.TermId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.Name = _updName
    return o._adapter.AffectedRows(),nil
}

// UpdateSlug an immediate DB Query to update a single column, in this
// case slug
func (o *Term) UpdateSlug(_updSlug string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `slug` = '%s' WHERE `term_id` = '%d'",o._table,_updSlug,o.TermId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.Slug = _updSlug
    return o._adapter.AffectedRows(),nil
}

// UpdateTermGroup an immediate DB Query to update a single column, in this
// case term_group
func (o *Term) UpdateTermGroup(_updTermGroup int64) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `term_group` = '%d' WHERE `term_id` = '%d'",o._table,_updTermGroup,o.TermId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.TermGroup = _updTermGroup
    return o._adapter.AffectedRows(),nil
}

// UserMeta is a Object Relational Mapping to
// the database table that represents it. In this case it is
// usermeta. The table name will be Sprintf'd to include
// the prefix you define in your YAML configuration for the
// Adapter.
type UserMeta struct {
    _table string
    _adapter Adapter
    _pkey string // 0 The name of the primary key in this table
    _conds []string
    _new bool
    UMetaId int64
    UserId int64
    MetaKey string
    MetaValue string
	// Dirty markers for smart updates
    IsUMetaIdDirty bool
    IsUserIdDirty bool
    IsMetaKeyDirty bool
    IsMetaValueDirty bool
	// Relationships
}

// NewUserMeta binds an Adapter to a new instance
// of UserMeta and sets up the _table and primary keys
func NewUserMeta(a Adapter) *UserMeta {
    var o UserMeta
    o._table = fmt.Sprintf("%susermeta",a.DatabasePrefix())
    o._adapter = a
    o._pkey = "umeta_id"
    o._new = false
    return &o
}


// GetPrimaryKeyValue returns the value, usually int64 of
// the PrimaryKey
func (o *UserMeta) GetPrimaryKeyValue() int64 {
    return o.UMetaId
}
// GetPrimaryKeyName returns the DB field name
func (o *UserMeta) GetPrimaryKeyName() string {
    return `umeta_id`
}

// GetUMetaId returns the value of 
// UserMeta.UMetaId
func (o *UserMeta) GetUMetaId() int64 {
    return o.UMetaId
}
// SetUMetaId sets and marks as dirty the value of
// UserMeta.UMetaId
func (o *UserMeta) SetUMetaId(arg int64) {
    o.UMetaId = arg
    o.IsUMetaIdDirty = true
}

// GetUserId returns the value of 
// UserMeta.UserId
func (o *UserMeta) GetUserId() int64 {
    return o.UserId
}
// SetUserId sets and marks as dirty the value of
// UserMeta.UserId
func (o *UserMeta) SetUserId(arg int64) {
    o.UserId = arg
    o.IsUserIdDirty = true
}

// GetMetaKey returns the value of 
// UserMeta.MetaKey
func (o *UserMeta) GetMetaKey() string {
    return o.MetaKey
}
// SetMetaKey sets and marks as dirty the value of
// UserMeta.MetaKey
func (o *UserMeta) SetMetaKey(arg string) {
    o.MetaKey = arg
    o.IsMetaKeyDirty = true
}

// GetMetaValue returns the value of 
// UserMeta.MetaValue
func (o *UserMeta) GetMetaValue() string {
    return o.MetaValue
}
// SetMetaValue sets and marks as dirty the value of
// UserMeta.MetaValue
func (o *UserMeta) SetMetaValue(arg string) {
    o.MetaValue = arg
    o.IsMetaValueDirty = true
}

// Find dynamic finder for umeta_id -> bool,error
// Generic and programatically generator finder for UserMeta
// Note that Fine returns a bool if found, not err, in the case of
// a return of true, the instance data will be filled out.
// a call to find ALWAYS overwrites the model you call Find on
// i.e. receiver is a pointer. 
//```go
//      m := NewUserMeta(a)
//      found,err := m.Find(23)
//      .. handle err
//      if found == false {
//          // handle found
//      }
//      ... do what you want with m here
//```
        func (o *UserMeta) Find(_findByUMetaId int64) (bool,error) {

    var _modelSlice []*UserMeta
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "umeta_id", _findByUMetaId)
    results, err := o._adapter.Query(q)
    if err != nil {
        return false,o._adapter.Oops(fmt.Sprintf(`%s`,err))
    }
    
    for _,result := range results {
        ro := UserMeta{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return false,o._adapter.Oops(fmt.Sprintf(`%s`,err))
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return false, o._adapter.Oops(`not found`)
    }
    o.FromUserMeta(_modelSlice[0])
    return true,nil

}
// FindByUserId dynamic finder for user_id -> []*UserMeta,error
// Generic and programatically generator finder for UserMeta
//```go  
//    m := NewUserMeta(a)
//    results,err := m.FindByUserId(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of UserMeta
//    }
//```  
        func (o *UserMeta) FindByUserId(_findByUserId int64) ([]*UserMeta,error) {

    var _modelSlice []*UserMeta
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'",o._table, "user_id", _findByUserId)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := UserMeta{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByMetaKey dynamic finder for meta_key -> []*UserMeta,error
// Generic and programatically generator finder for UserMeta
//```go  
//    m := NewUserMeta(a)
//    results,err := m.FindByMetaKey(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of UserMeta
//    }
//```  
        func (o *UserMeta) FindByMetaKey(_findByMetaKey string) ([]*UserMeta,error) {

    var _modelSlice []*UserMeta
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "meta_key", _findByMetaKey)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := UserMeta{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}
// FindByMetaValue dynamic finder for meta_value -> []*UserMeta,error
// Generic and programatically generator finder for UserMeta
//```go  
//    m := NewUserMeta(a)
//    results,err := m.FindByMetaValue(...)
//    // handle err
//    for i,r := results {
//      // now r is an instance of UserMeta
//    }
//```  
        func (o *UserMeta) FindByMetaValue(_findByMetaValue string) ([]*UserMeta,error) {

    var _modelSlice []*UserMeta
    q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%s'",o._table, "meta_value", _findByMetaValue)
    results, err := o._adapter.Query(q)
    if err != nil {
        return _modelSlice,err
    }
    
    for _,result := range results {
        ro := UserMeta{}
        err = ro.FromDBValueMap(result)
        if err != nil {
            return _modelSlice,err
        }
        _modelSlice = append(_modelSlice,&ro)
    }

    if len(_modelSlice) == 0 {
        // there was an error!
        return nil, o._adapter.Oops(`no results`)
    }
    return _modelSlice,nil

}

// FromDBValueMap Converts a DBValueMap returned from Adapter.Query to a UserMeta
func (o *UserMeta) FromDBValueMap(m map[string]DBValue) error {
	_UMetaId,err := m["umeta_id"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.UMetaId = _UMetaId
	_UserId,err := m["user_id"].AsInt64()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.UserId = _UserId
	_MetaKey,err := m["meta_key"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.MetaKey = _MetaKey
	_MetaValue,err := m["meta_value"].AsString()
	if err != nil {
 		return o._adapter.Oops(fmt.Sprintf(`%s`,err))
	}
	o.MetaValue = _MetaValue

 	return nil
}
// FromUserMeta A kind of Clone function for UserMeta
func (o *UserMeta) FromUserMeta(m *UserMeta) {
	o.UMetaId = m.UMetaId
	o.UserId = m.UserId
	o.MetaKey = m.MetaKey
	o.MetaValue = m.MetaValue

}
// Reload A function to forcibly reload UserMeta
func (o *UserMeta) Reload() error {
    _,err := o.Find(o.GetPrimaryKeyValue())
    return err
}

// Save is a dynamic saver 'inherited' by all models
func (o *UserMeta) Save() error {
    if o._new == true {
        return o.Create()
    }
    var sets []string
    
    if o.IsUserIdDirty == true {
        sets = append(sets,fmt.Sprintf(`user_id = '%d'`,o.UserId))
    }

    if o.IsMetaKeyDirty == true {
        sets = append(sets,fmt.Sprintf(`meta_key = '%s'`,o._adapter.SafeString(o.MetaKey)))
    }

    if o.IsMetaValueDirty == true {
        sets = append(sets,fmt.Sprintf(`meta_value = '%s'`,o._adapter.SafeString(o.MetaValue)))
    }

    frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'",o._table,strings.Join(sets,`,`),o._pkey, o.UMetaId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return err
    }
    return nil
}
// Update is a dynamic updater, it considers whether or not
// a field is 'dirty' and needs to be updated. Will only work
// if you use the Getters and Setters
func (o *UserMeta) Update() error {
    var sets []string
    
    if o.IsUserIdDirty == true {
        sets = append(sets,fmt.Sprintf(`user_id = '%d'`,o.UserId))
    }

    if o.IsMetaKeyDirty == true {
        sets = append(sets,fmt.Sprintf(`meta_key = '%s'`,o._adapter.SafeString(o.MetaKey)))
    }

    if o.IsMetaValueDirty == true {
        sets = append(sets,fmt.Sprintf(`meta_value = '%s'`,o._adapter.SafeString(o.MetaValue)))
    }

    frmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%d'",o._table,strings.Join(sets,`,`),o._pkey, o.UMetaId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return err
    }
    return nil
}
// Create inserts the model. Calling Save will call this function
// automatically for new models
func (o *UserMeta) Create() error {
    frmt := fmt.Sprintf("INSERT INTO %s (`user_id`, `meta_key`, `meta_value`) VALUES ('%d', '%s', '%s')",o._table,o.UserId, o.MetaKey, o.MetaValue)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return o._adapter.Oops(fmt.Sprintf(`%s led to %s`,frmt,err))
    }
    o.UMetaId = o._adapter.LastInsertedId()
    o._new = false
    return nil
}


// UpdateUserId an immediate DB Query to update a single column, in this
// case user_id
func (o *UserMeta) UpdateUserId(_updUserId int64) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `user_id` = '%d' WHERE `umeta_id` = '%d'",o._table,_updUserId,o.UMetaId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.UserId = _updUserId
    return o._adapter.AffectedRows(),nil
}

// UpdateMetaKey an immediate DB Query to update a single column, in this
// case meta_key
func (o *UserMeta) UpdateMetaKey(_updMetaKey string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `meta_key` = '%s' WHERE `umeta_id` = '%d'",o._table,_updMetaKey,o.UMetaId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.MetaKey = _updMetaKey
    return o._adapter.AffectedRows(),nil
}

// UpdateMetaValue an immediate DB Query to update a single column, in this
// case meta_value
func (o *UserMeta) UpdateMetaValue(_updMetaValue string) (int64,error) {
    frmt := fmt.Sprintf("UPDATE %s SET `meta_value` = '%s' WHERE `umeta_id` = '%d'",o._table,_updMetaValue,o.UMetaId)
    err := o._adapter.Execute(frmt)
    if err != nil {
        return 0,err
    }
    o.MetaValue = _updMetaValue
    return o._adapter.AffectedRows(),nil
}

