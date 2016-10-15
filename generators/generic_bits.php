<?php
puts("
type Adapter interface {
    Open(string,string,string,string) error
    Close()
    Query(string) ([]map[string]DBValue,error)
    Execute(string) error
    LastInsertedId() int64
    AffectedRows() int64
    DatabasePrefix() string
    NewDBValue() DBValue
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
    dt := NewDateTime()
    err := dt.FromString(v._v)
    if err != nil {
        return &DateTime{}, err
    }
    return dt,nil
}

func NewMysqlValue() *MysqlValue {
    return &MysqlValue{}
}

type MysqlAdapter struct {
    Host string `yaml:\"host\"`
    User string `yaml:\"user\"`
    Pass string `yaml: \"pass\"`
    Database string `yaml:\"database\"`
    DBPrefix string `yaml:\"prefix\"`
    _conn_ *sql.DB
    _lid int64
    _cnt int64
    _opened bool
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
    return a,err

}
func (a *MysqlAdapter) NewDBValue() DBValue {
    return NewMysqlValue()
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
            return errors.New(fmt.Sprintf(`%s with %s`,err,l))
        }
        a._conn_ = tc
    } else {
        l := fmt.Sprintf(\"%s:%s@/%s\",u,p,d)
        tc, err := sql.Open(\"mysql\",l)
        if err != nil {
            return errors.New(fmt.Sprintf(`%s with %s`,err,l))
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
        return nil,errors.New(`you must first open the connection`)
    }
    results := new([]map[string]DBValue)
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
func (a *MysqlAdapter) Execute(q string) error {
    if a._opened != true {
        return errors.New(`you must first open the connection`)
    }
    tx, err := a._conn_.Begin()
    if err != nil {
        return errors.New(fmt.Sprintf(`could not Begin Transaction %s`,err))
    }
    defer tx.Rollback();
    stmt, err := tx.Prepare(q)
    if err != nil {
        return errors.New(fmt.Sprintf(`could not Prepare Statement %s`,err))
    }
    defer stmt.Close()
    res,err := stmt.Exec()
    if err != nil {
        return errors.New(fmt.Sprintf(`could not Exec stmt %s`,err))
    }
    a._lid,err = res.LastInsertId()
    if err != nil {
        return errors.New(fmt.Sprintf(`could not get LastInsertId %s`,err))
    }
    a._cnt,err = res.RowsAffected()
    if err != nil {
        return errors.New(fmt.Sprintf(`could not get RowsAffected %s`,err))
    }
    err = tx.Commit()
    if err != nil {
        return errors.New(fmt.Sprintf(`could not Commit Transaction %s`,err))
    }
    return nil
}
func (a *MysqlAdapter) LastInsertedId() int64 {
    return a._lid
}
func (a *MysqlAdapter) AffectedRows() int64 {
    return a._cnt
}
type DateTime struct {
    Day int
    Month int
    Year int
    Hours int
    Minutes int
    Seconds int
}
func (d *DateTime) FromString(s string) error {
    es := s
    re := regexp.MustCompile(\"(?P<year>[\\\d]{4})-(?P<month>[\\\d]{2})-(?P<day>[\\\d]{2}) (?P<hours>[\\\d]{2}):(?P<minutes>[\\\d]{2}):(?P<seconds>[\\\d]{2})\")
    n1 := re.SubexpNames()
    ir2 := re.FindAllStringSubmatch(es, -1)
    if len(ir2) == 0 {
        return errors.New(fmt.Sprintf(\"found no data to capture in %s\",es))
    }
    r2 := ir2[0]
    for i, n := range r2 {
        if n1[i] == \"year\" {
            _Year,err := strconv.ParseInt(n,10,32)
            d.Year = int(_Year)
            if err != nil {
                return errors.New(fmt.Sprintf(\"failed to convert %s in %s received %s\",n[i],es,err))
            }
        }
        if n1[i] == \"month\" {
            _Month,err := strconv.ParseInt(n,10,32)
            d.Month = int(_Month)
            if err != nil {
                return errors.New(fmt.Sprintf(\"failed to convert %s in %s received %s\",n[i],es,err))
            }
        }
        if n1[i] == \"day\" {
            _Day,err := strconv.ParseInt(n,10,32)
            d.Day = int(_Day)
            if err != nil {
                return errors.New(fmt.Sprintf(\"failed to convert %s in %s received %s\",n[i],es,err))
            }
        }
        if n1[i] == \"hours\" {
            _Hours,err := strconv.ParseInt(n,10,32)
            d.Hours = int(_Hours)
            if err != nil {
                return errors.New(fmt.Sprintf(\"failed to convert %s in %s received %s\",n[i],es,err))
            }
        }
        if n1[i] == \"minutes\" {
            _Minutes,err := strconv.ParseInt(n,10,32)
            d.Minutes = int(_Minutes)
            if err != nil {
                return errors.New(fmt.Sprintf(\"failed to convert %s in %s received %s\",n[i],es,err))
            }
        }
        if n1[i] == \"seconds\" {
            _Seconds,err := strconv.ParseInt(n,10,32)
            d.Seconds = int(_Seconds)
            if err != nil {
                return errors.New(fmt.Sprintf(\"failed to convert %s in %s received %s\",n[i],es,err))
            }
        }
    }
    return nil
}
func (d *DateTime) ToString() string {
    return fmt.Sprintf(\"%d-%02d-%02d %02d:%02d:%02d\",d.Year,d.Month,d.Day,d.Hours,d.Minutes,d.Seconds)
}
func NewDateTime() *DateTime {
    d := &DateTime{}
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