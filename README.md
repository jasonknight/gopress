# Gopress
[![Build Status](https://travis-ci.org/jasonknight/gopress.svg?branch=master)](https://travis-ci.org/jasonknight/gopress)
[![codecov.io](https://codecov.io/gh/jasonknight/gopress/coverage.svg?branch=master)](https://codecov.io/gh/jasonknight/gopress)

This is a Golang [ActiveRecord](https://en.wikipedia.org/wiki/Active_record_pattern) implementation of the Wordpress database Schema (also includes WooCommerce).

```go
type Post struct {
    ...
    ID int64
    Author int64
    Date DateTime
    DateGmt DateTime
    Content string
    Title string
    Excerpt string
    Status string
    CommentStatus string
    PingStatus string
    Password string
    Name string
    ToPing string
    Pinged string
    Modified DateTime
    ModifiedGmt DateTime
    ...
}
```

The library implements basic CRUD functions (Create/Find/Update/Delete) and provides sane structs as models. 

Each model must be provided with an adapter:

```go
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
```

An adapter for MySQL is supplied:

```go
type MysqlAdapter struct {
    Host     string `yaml:"host"`
    User     string `yaml:"user"`
    Pass     string `yaml: "pass"`
    Database string `yaml:"database"`
    DBPrefix string `yaml:"prefix"`
    _conn   *sql.DB
    _lid     int64
    _cnt     int64
}
```

Which provides the required functions, however you can
supply your own if you have special needs. The interface
is generic enough that you can use it any weird way
you want.

I would point out that you may want to migrate your wordpress to another database, like PostgreSQL, but this is not supported by Wordpress. Ideally, this library would let you swap, and it is intended to be used in my planned drop in replacement for wordpress named Caddy (gophers, Caddy Shack...haha). The biggest issue with wordpress is that you are trapped in their bad decisions (same with WooCommerce). 


