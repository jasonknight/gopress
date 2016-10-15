<?php
puts("
type LogFilter func (string,string)string
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
    NewDBValue() DBValue
}

type MysqlAdapter struct {
    Host string `yaml:\"host\"`
    User string `yaml:\"user\"`
    Pass string `yaml: \"pass\"`
    Database string `yaml:\"database\"`
    DBPrefix string `yaml:\"prefix\"`
    _info_log *log.Logger
    _error_log *log.Logger
    _debug_log *log.Logger
    _conn_ *sql.DB
    _lid int64
    _cnt int64
    _opened bool
    _log_filter LogFilter
}

func NewMysqlAdapter(pre string) *MysqlAdapter {
    return &MysqlAdapter{DBPrefix: pre}
} 
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
func (a *MysqlAdapter) SetLogFilter(f LogFilter) {
    a._log_filter = f
}
func (a *MysqlAdapter) SetInfoLog(t io.Writer) {
    a._info_log = log.New(t,`[INFO]:`,log.Ldate|log.Ltime|log.Lshortfile)
}
func (a *MysqlAdapter) SetErrorLog(t io.Writer) {
    a._error_log = log.New(t,`[ERROR]:`,log.Ldate|log.Ltime|log.Lshortfile)
}
func (a *MysqlAdapter) SetDebugLog(t io.Writer) {
    a._debug_log = log.New(t,`[DEBUG]:`,log.Ldate|log.Ltime|log.Lshortfile)
}
func (a *MysqlAdapter) SetLogs(t io.Writer) {
    a.SetInfoLog(t)
    a.SetErrorLog(t)
    a.SetDebugLog(t)
}

func (a *MysqlAdapter) LogInfo(s string) {
    if a._log_filter != nil {
        s = a._log_filter(`INFO`,s)
    }
    if s == \"\" {
        return
    }
    a._info_log.Println(s)
}

func (a *MysqlAdapter) LogError(s error) {
    if a._log_filter != nil {
        ns := a._log_filter(`ERROR`,fmt.Sprintf(`%s`,s))
        a._error_log.Println(ns)
        return
    }
    a._error_log.Println(s)
}

func (a *MysqlAdapter) LogDebug(s string) {
    if a._log_filter != nil {
        s = a._log_filter(`DEBUG`,s)
    }
    if s == \"\" {
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
    return yaml.Unmarshal(b,a)
}

func (a *MysqlAdapter) Open(h,u,p,d string) error {
    if ( h != \"localhost\") {
        l := fmt.Sprintf(\"%s:%s@tcp(%s)/%s\",u,p,h,d)
        tc, err := sql.Open(\"mysql\",l)
        if err != nil {
            return a.Oops(fmt.Sprintf(`%s with %s`,err,l))
        }
        a._conn_ = tc
    } else {
        l := fmt.Sprintf(\"%s:%s@/%s\",u,p,d)
        tc, err := sql.Open(\"mysql\",l)
        if err != nil {
            return a.Oops(fmt.Sprintf(`%s with %s`,err,l))
        }
        a._conn_ = tc
    }
    a._opened = true
    return nil

}
func (a *MysqlAdapter) Close() {
    a._conn_.Close()
}

func (a *MysqlAdapter) Query(q string) ([]map[string]DBValue,error) {
    if a._opened != true {
        return nil,a.Oops(`you must first open the connection`)
    }
    results := new([]map[string]DBValue)
    a.LogInfo(q)
    rows, err := a._conn_.Query(q)
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
func (a *MysqlAdapter) LastInsertedId() int64 {
    return a._lid
}
func (a *MysqlAdapter) AffectedRows() int64 {
    return a._cnt
}

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

type MysqlValue struct {
    _v string
    _k string
    _adapter Adapter
}
func (v *MysqlValue) SetInternalValue(key,value string) {
    v._v = value
    v._k = key

}
func (v *MysqlValue) AsString() (string,error) {
    return v._v,nil
}
func (v *MysqlValue) AsInt() (int,error) {
    i,err := strconv.ParseInt(v._v,10,32)
    return int(i),err
}
func (v *MysqlValue) AsInt32() (int32,error) {
    i,err := strconv.ParseInt(v._v,10,32)
    return int32(i),err
}
func (v *MysqlValue) AsInt64() (int64,error) {
    i,err := strconv.ParseInt(v._v,10,64)
    return i,err
}
func (v *MysqlValue) AsFloat32() (float32,error) {
    i,err := strconv.ParseFloat(v._v,32)
    if err != nil {
        return 0.0,err
    }
    return float32(i),err
}
func (v *MysqlValue) AsFloat64() (float64,error) {
    i,err := strconv.ParseFloat(v._v,64)
    if err != nil {
        return 0.0,err
    }
    return i,err
}

func (v *MysqlValue) AsDateTime() (*DateTime,error) {
    dt := NewDateTime(v._adapter)
    err := dt.FromString(v._v)
    if err != nil {
        return &DateTime{}, err
    }
    return dt,nil
}

func NewMysqlValue(a Adapter) *MysqlValue {
    return &MysqlValue{_adapter: a}
}
type DateTime struct {
    Day int
    Month int
    Year int
    Hours int
    Minutes int
    Seconds int
    _adapter Adapter
}
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
                return d._adapter.Oops(fmt.Sprintf(\"failed to convert %s in %s received %s\",n[i],es,err))
            }
        }
        if n1[i] == \"month\" {
            _Month,err := strconv.ParseInt(n,10,32)
            d.Month = int(_Month)
            if err != nil {
                return d._adapter.Oops(fmt.Sprintf(\"failed to convert %s in %s received %s\",n[i],es,err))
            }
        }
        if n1[i] == \"day\" {
            _Day,err := strconv.ParseInt(n,10,32)
            d.Day = int(_Day)
            if err != nil {
                return d._adapter.Oops(fmt.Sprintf(\"failed to convert %s in %s received %s\",n[i],es,err))
            }
        }
        if n1[i] == \"hours\" {
            _Hours,err := strconv.ParseInt(n,10,32)
            d.Hours = int(_Hours)
            if err != nil {
                return d._adapter.Oops(fmt.Sprintf(\"failed to convert %s in %s received %s\",n[i],es,err))
            }
        }
        if n1[i] == \"minutes\" {
            _Minutes,err := strconv.ParseInt(n,10,32)
            d.Minutes = int(_Minutes)
            if err != nil {
                return d._adapter.Oops(fmt.Sprintf(\"failed to convert %s in %s received %s\",n[i],es,err))
            }
        }
        if n1[i] == \"seconds\" {
            _Seconds,err := strconv.ParseInt(n,10,32)
            d.Seconds = int(_Seconds)
            if err != nil {
                return d._adapter.Oops(fmt.Sprintf(\"failed to convert %s in %s received %s\",n[i],es,err))
            }
        }
    }
    return nil
}
func (d *DateTime) ToString() string {
    return fmt.Sprintf(\"%d-%02d-%02d %02d:%02d:%02d\",d.Year,d.Month,d.Day,d.Hours,d.Minutes,d.Seconds)
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
");