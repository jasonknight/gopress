# Gopress
[![Build Status](https://travis-ci.org/jasonknight/gopress.svg?branch=master)](https://travis-ci.org/jasonknight/gopress)
[![codecov.io](https://codecov.io/gh/jasonknight/gopress/coverage.svg?branch=master)](https://codecov.io/gh/jasonknight/gopress)

This is a Golang [ActiveRecord](https://en.wikipedia.org/wiki/Active_record_pattern) implementation of the Wordpress database Schema (also includes WooCommerce).

```go
type Post struct {
    _table string
    _adapter Adapter
    _pkey string // 0 The name of the primary key in this table
    _conds []string
    _new bool
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
    ContentFiltered string
    Parent int64
    Guid string
    MenuOrder int
    Type string
    MimeType string
    CommentCount int64
}
```


