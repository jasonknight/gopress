# Gopress
[![Build Status](https://travis-ci.org/jasonknight/gopress.svg?branch=master)](https://travis-ci.org/jasonknight/gopress)
[![codecov](https://codecov.io/gh/jasonknight/gopress/branch/master/graph/badge.svg)](https://codecov.io/gh/jasonknight/gopress)
[![Go Report Card](https://goreportcard.com/badge/github.com/jasonknight/gopress?123)](https://goreportcard.com/report/github.com/jasonknight/gopress)

This is a Golang [ActiveRecord](https://en.wikipedia.org/wiki/Active_record_pattern) implementation of the Wordpress database Schema. This library is automatically generated. If you'd like to contribute, you can change the go files and I will merge them into the generator, or if you're brave you can try
to futz with the generator. 

I am still figuring out codecov.io which shows less coverage than
when local testing, current testing against a database in 75%. Needs
improvement - and I'll continue to work on that.

[Please see the wiki for details and examples](https://github.com/jasonknight/gopress/wiki), or look in the testing file. [Check out the docs](https://github.com/jasonknight/gopress/blob/master/docs.md).

The library implements basic CRUD functions (Create/Read(Find)/Update/Delete) and provides structs as "models". 

Each model must be provided with an adapter:

```go
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
    NewDBValue() DBValue
}
```

An adapter for MySQL is supplied:

```go
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

```

Which provides the required functions, however you can
supply your own if you have special needs. The interface
is generic enough that you can use it any weird way
you want.

[![Become A Patron](https://github.com/jasonknight/gobay/raw/master/assets/patreon.png)](https://www.patreon.com/user?u=4141497)

