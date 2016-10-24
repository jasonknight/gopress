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
");
include "mysql_adapter.php";
puts("
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
include "model_functions.php";