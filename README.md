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

#### type Post

```go
type Post struct {
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
    ...
    IsCommentCountDirty        bool
}
```

Post is a Object Relational Mapping to the database table that represents it. In
this case it is posts. The table name will be Sprintf'd to include the prefix
you define in your YAML configuration for the Adapter.

#### func  NewPost

```go
func NewPost(a Adapter) *Post
```
NewPost binds an Adapter to a new instance of Post and sets up the _table and
primary keys

#### func (*Post) Create

```go
func (o *Post) Create() error
```
Create inserts the model. Calling Save will call this function automatically for
new models

#### func (*Post) Find

```go
func (o *Post) Find(_findByID int64) (bool, error)
```
Find dynamic finder for ID -> bool,error Generic and programatically generator
finder for Post

Note that Fine returns a bool if found, not err, in the case of a return of
true, the instance data will be filled out. a call to find ALWAYS overwrites the
model you call Find on i.e. receiver is a pointer.

```go

    m := NewPost(a)
    found,err := m.Find(23)
    .. handle err
    if found == false {
        // handle found
    }
    ... do what you want with m here

```

#### func (*Post) FindByGuid

```go
func (o *Post) FindByGuid(_findByGuid string) ([]*Post, error)
```
FindByGuid dynamic finder for guid -> []*Post,error Generic and programatically
generator finder for Post

```go

    m := NewPost(a)
    results,err := m.FindByGuid(...)
    // handle err
    for i,r := results {
      // now r is an instance of Post
    }

```

#### func (*Post) FindByPostAuthor

```go
func (o *Post) FindByPostAuthor(_findByPostAuthor int64) ([]*Post, error)
```
FindByPostAuthor dynamic finder for post_author -> []*Post,error Generic and
programatically generator finder for Post

```go

    m := NewPost(a)
    results,err := m.FindByPostAuthor(...)
    // handle err
    for i,r := results {
      // now r is an instance of Post
    }

```

For more documentation of each supported type, [check out the docs](https://github.com/jasonknight/gopress/blob/master/docs.md).

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

