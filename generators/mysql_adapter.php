<?php
puts("
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
// NewMysqlAdapterEx sets everything up based on your YAML config
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
");