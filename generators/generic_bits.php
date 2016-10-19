<?php
puts("
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
    Host string `yaml:\"host\"`
    // The database username
    User string `yaml:\"user\"`
    // The database password
    Pass string `yaml:\"pass\"`
    // The database name
    Database string `yaml:\"database\"`
    // A prefix, if any - can be blank
    DBPrefix string `yaml:\"prefix\"`
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
// NewMysqlAdapterEx
// Args: fname is a string path to a YAML config file
// This function will attempt to Open the database
// defined in that file. Example file:
//     host: \"localhost\"
//     user: \"dbuser\"
//     pass: \"dbuserpass\"
//     database: \"my_db\"
//     prefix: \"wp_\"
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
    if s == \"\" {
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
    if s == \"\" {
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
    if ( h != \"localhost\") {
        l := fmt.Sprintf(\"%s:%s@tcp(%s)/%s\",u,p,h,d)
        tc, err := sql.Open(\"mysql\",l)
        if err != nil {
            return a.Oops(fmt.Sprintf(`%s with %s`,err,l))
        }
        a._conn = tc
    } else {
        l := fmt.Sprintf(\"%s:%s@/%s\",u,p,d)
        tc, err := sql.Open(\"mysql\",l)
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
    re := regexp.MustCompile(\"(?P<year>[\\\d]{4})-(?P<month>[\\\d]{2})-(?P<day>[\\\d]{2}) (?P<hours>[\\\d]{2}):(?P<minutes>[\\\d]{2}):(?P<seconds>[\\\d]{2})\")
    n1 := re.SubexpNames()
    ir2 := re.FindAllStringSubmatch(es, -1)
    if len(ir2) == 0 {
        return d._adapter.Oops(fmt.Sprintf(\"found no data to capture in %s\",es))
    }
    r2 := ir2[0]
    for i, n := range r2 {
        if n1[i] == \"year\" {
            _Year,err := strconv.ParseInt(n,10,32)
            d.Year = int(_Year)
            if err != nil {
                return d._adapter.Oops(fmt.Sprintf(\"failed to convert %d in %v received %s\",n[i],es,err))
            }
        }
        if n1[i] == \"month\" {
            _Month,err := strconv.ParseInt(n,10,32)
            d.Month = int(_Month)
            if err != nil {
                return d._adapter.Oops(fmt.Sprintf(\"failed to convert %d in %v received %s\",n[i],es,err))
            }
        }
        if n1[i] == \"day\" {
            _Day,err := strconv.ParseInt(n,10,32)
            d.Day = int(_Day)
            if err != nil {
                return d._adapter.Oops(fmt.Sprintf(\"failed to convert %d in %v received %s\",n[i],es,err))
            }
        }
        if n1[i] == \"hours\" {
            _Hours,err := strconv.ParseInt(n,10,32)
            d.Hours = int(_Hours)
            if err != nil {
                return d._adapter.Oops(fmt.Sprintf(\"failed to convert %d in %v received %s\",n[i],es,err))
            }
        }
        if n1[i] == \"minutes\" {
            _Minutes,err := strconv.ParseInt(n,10,32)
            d.Minutes = int(_Minutes)
            if err != nil {
                return d._adapter.Oops(fmt.Sprintf(\"failed to convert %d in %v received %s\",n[i],es,err))
            }
        }
        if n1[i] == \"seconds\" {
            _Seconds,err := strconv.ParseInt(n,10,32)
            d.Seconds = int(_Seconds)
            if err != nil {
                return d._adapter.Oops(fmt.Sprintf(\"failed to convert %d in %v received %s\",n[i],es,err))
            }
        }
    }
    return nil
}
// ToString For backwards compat... Never use this, use String() instead.
func (d *DateTime) ToString() string {
    return fmt.Sprintf(\"%d-%02d-%02d %02d:%02d:%02d\",d.Year,d.Month,d.Day,d.Hours,d.Minutes,d.Seconds)
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
");