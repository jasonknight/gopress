# gopress
--
    import "github.com/jasonknight/gopress"


## Usage

#### type Adapter

```go
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
	SafeString(string) string
	NewDBValue() DBValue
}
```

Adapter is the main Database interface which helps to separate the DB from the
Models. This is not 100% just yet, and may never be. Eventually the Adapter will
probably receive some arguments and a value map and build the Query internally

#### type Comment

```go
type Comment struct {
	CommentID          int64
	CommentPostID      int64
	CommentAuthor      string
	CommentAuthorEmail string
	CommentAuthorUrl   string
	CommentAuthorIP    string
	CommentDate        *DateTime
	CommentDateGmt     *DateTime
	CommentContent     string
	CommentKarma       int
	CommentApproved    string
	CommentAgent       string
	CommentType        string
	CommentParent      int64
	UserId             int64
	// Dirty markers for smart updates
	IsCommentIDDirty          bool
	IsCommentPostIDDirty      bool
	IsCommentAuthorDirty      bool
	IsCommentAuthorEmailDirty bool
	IsCommentAuthorUrlDirty   bool
	IsCommentAuthorIPDirty    bool
	IsCommentDateDirty        bool
	IsCommentDateGmtDirty     bool
	IsCommentContentDirty     bool
	IsCommentKarmaDirty       bool
	IsCommentApprovedDirty    bool
	IsCommentAgentDirty       bool
	IsCommentTypeDirty        bool
	IsCommentParentDirty      bool
	IsUserIdDirty             bool
}
```

Comment is a Object Relational Mapping to the database table that represents it.
In this case it is comments. The table name will be Sprintf'd to include the
prefix you define in your YAML configuration for the Adapter.

#### func  NewComment

```go
func NewComment(a Adapter) *Comment
```
NewComment binds an Adapter to a new instance of Comment and sets up the _table
and primary keys

#### func (*Comment) Create

```go
func (o *Comment) Create() error
```
Create inserts the model. Calling Save will call this function automatically for
new models

#### func (*Comment) Find

```go
func (o *Comment) Find(_findByCommentID int64) (bool, error)
```
Find dynamic finder for comment_ID -> bool,error Generic and programatically
generator finder for Comment

Note that Fine returns a bool if found, not err, in the case of a return of
true, the instance data will be filled out. a call to find ALWAYS overwrites the
model you call Find on i.e. receiver is a pointer.

```go

    m := NewComment(a)
    found,err := m.Find(23)
    .. handle err
    if found == false {
        // handle found
    }
    ... do what you want with m here

```

#### func (*Comment) FindByCommentAgent

```go
func (o *Comment) FindByCommentAgent(_findByCommentAgent string) ([]*Comment, error)
```
FindByCommentAgent dynamic finder for comment_agent -> []*Comment,error Generic
and programatically generator finder for Comment

```go

    m := NewComment(a)
    results,err := m.FindByCommentAgent(...)
    // handle err
    for i,r := results {
      // now r is an instance of Comment
    }

```

#### func (*Comment) FindByCommentApproved

```go
func (o *Comment) FindByCommentApproved(_findByCommentApproved string) ([]*Comment, error)
```
FindByCommentApproved dynamic finder for comment_approved -> []*Comment,error
Generic and programatically generator finder for Comment

```go

    m := NewComment(a)
    results,err := m.FindByCommentApproved(...)
    // handle err
    for i,r := results {
      // now r is an instance of Comment
    }

```

#### func (*Comment) FindByCommentAuthor

```go
func (o *Comment) FindByCommentAuthor(_findByCommentAuthor string) ([]*Comment, error)
```
FindByCommentAuthor dynamic finder for comment_author -> []*Comment,error
Generic and programatically generator finder for Comment

```go

    m := NewComment(a)
    results,err := m.FindByCommentAuthor(...)
    // handle err
    for i,r := results {
      // now r is an instance of Comment
    }

```

#### func (*Comment) FindByCommentAuthorEmail

```go
func (o *Comment) FindByCommentAuthorEmail(_findByCommentAuthorEmail string) ([]*Comment, error)
```
FindByCommentAuthorEmail dynamic finder for comment_author_email ->
[]*Comment,error Generic and programatically generator finder for Comment

```go

    m := NewComment(a)
    results,err := m.FindByCommentAuthorEmail(...)
    // handle err
    for i,r := results {
      // now r is an instance of Comment
    }

```

#### func (*Comment) FindByCommentAuthorIP

```go
func (o *Comment) FindByCommentAuthorIP(_findByCommentAuthorIP string) ([]*Comment, error)
```
FindByCommentAuthorIP dynamic finder for comment_author_IP -> []*Comment,error
Generic and programatically generator finder for Comment

```go

    m := NewComment(a)
    results,err := m.FindByCommentAuthorIP(...)
    // handle err
    for i,r := results {
      // now r is an instance of Comment
    }

```

#### func (*Comment) FindByCommentAuthorUrl

```go
func (o *Comment) FindByCommentAuthorUrl(_findByCommentAuthorUrl string) ([]*Comment, error)
```
FindByCommentAuthorUrl dynamic finder for comment_author_url -> []*Comment,error
Generic and programatically generator finder for Comment

```go

    m := NewComment(a)
    results,err := m.FindByCommentAuthorUrl(...)
    // handle err
    for i,r := results {
      // now r is an instance of Comment
    }

```

#### func (*Comment) FindByCommentContent

```go
func (o *Comment) FindByCommentContent(_findByCommentContent string) ([]*Comment, error)
```
FindByCommentContent dynamic finder for comment_content -> []*Comment,error
Generic and programatically generator finder for Comment

```go

    m := NewComment(a)
    results,err := m.FindByCommentContent(...)
    // handle err
    for i,r := results {
      // now r is an instance of Comment
    }

```

#### func (*Comment) FindByCommentDate

```go
func (o *Comment) FindByCommentDate(_findByCommentDate *DateTime) ([]*Comment, error)
```
FindByCommentDate dynamic finder for comment_date -> []*Comment,error Generic
and programatically generator finder for Comment

```go

    m := NewComment(a)
    results,err := m.FindByCommentDate(...)
    // handle err
    for i,r := results {
      // now r is an instance of Comment
    }

```

#### func (*Comment) FindByCommentDateGmt

```go
func (o *Comment) FindByCommentDateGmt(_findByCommentDateGmt *DateTime) ([]*Comment, error)
```
FindByCommentDateGmt dynamic finder for comment_date_gmt -> []*Comment,error
Generic and programatically generator finder for Comment

```go

    m := NewComment(a)
    results,err := m.FindByCommentDateGmt(...)
    // handle err
    for i,r := results {
      // now r is an instance of Comment
    }

```

#### func (*Comment) FindByCommentKarma

```go
func (o *Comment) FindByCommentKarma(_findByCommentKarma int) ([]*Comment, error)
```
FindByCommentKarma dynamic finder for comment_karma -> []*Comment,error Generic
and programatically generator finder for Comment

```go

    m := NewComment(a)
    results,err := m.FindByCommentKarma(...)
    // handle err
    for i,r := results {
      // now r is an instance of Comment
    }

```

#### func (*Comment) FindByCommentParent

```go
func (o *Comment) FindByCommentParent(_findByCommentParent int64) ([]*Comment, error)
```
FindByCommentParent dynamic finder for comment_parent -> []*Comment,error
Generic and programatically generator finder for Comment

```go

    m := NewComment(a)
    results,err := m.FindByCommentParent(...)
    // handle err
    for i,r := results {
      // now r is an instance of Comment
    }

```

#### func (*Comment) FindByCommentPostID

```go
func (o *Comment) FindByCommentPostID(_findByCommentPostID int64) ([]*Comment, error)
```
FindByCommentPostID dynamic finder for comment_post_ID -> []*Comment,error
Generic and programatically generator finder for Comment

```go

    m := NewComment(a)
    results,err := m.FindByCommentPostID(...)
    // handle err
    for i,r := results {
      // now r is an instance of Comment
    }

```

#### func (*Comment) FindByCommentType

```go
func (o *Comment) FindByCommentType(_findByCommentType string) ([]*Comment, error)
```
FindByCommentType dynamic finder for comment_type -> []*Comment,error Generic
and programatically generator finder for Comment

```go

    m := NewComment(a)
    results,err := m.FindByCommentType(...)
    // handle err
    for i,r := results {
      // now r is an instance of Comment
    }

```

#### func (*Comment) FindByUserId

```go
func (o *Comment) FindByUserId(_findByUserId int64) ([]*Comment, error)
```
FindByUserId dynamic finder for user_id -> []*Comment,error Generic and
programatically generator finder for Comment

```go

    m := NewComment(a)
    results,err := m.FindByUserId(...)
    // handle err
    for i,r := results {
      // now r is an instance of Comment
    }

```

#### func (*Comment) FromComment

```go
func (o *Comment) FromComment(m *Comment)
```
FromComment A kind of Clone function for Comment

#### func (*Comment) FromDBValueMap

```go
func (o *Comment) FromDBValueMap(m map[string]DBValue) error
```
FromDBValueMap Converts a DBValueMap returned from Adapter.Query to a Comment

#### func (*Comment) GetCommentAgent

```go
func (o *Comment) GetCommentAgent() string
```
GetCommentAgent returns the value of Comment.CommentAgent

#### func (*Comment) GetCommentApproved

```go
func (o *Comment) GetCommentApproved() string
```
GetCommentApproved returns the value of Comment.CommentApproved

#### func (*Comment) GetCommentAuthor

```go
func (o *Comment) GetCommentAuthor() string
```
GetCommentAuthor returns the value of Comment.CommentAuthor

#### func (*Comment) GetCommentAuthorEmail

```go
func (o *Comment) GetCommentAuthorEmail() string
```
GetCommentAuthorEmail returns the value of Comment.CommentAuthorEmail

#### func (*Comment) GetCommentAuthorIP

```go
func (o *Comment) GetCommentAuthorIP() string
```
GetCommentAuthorIP returns the value of Comment.CommentAuthorIP

#### func (*Comment) GetCommentAuthorUrl

```go
func (o *Comment) GetCommentAuthorUrl() string
```
GetCommentAuthorUrl returns the value of Comment.CommentAuthorUrl

#### func (*Comment) GetCommentContent

```go
func (o *Comment) GetCommentContent() string
```
GetCommentContent returns the value of Comment.CommentContent

#### func (*Comment) GetCommentDate

```go
func (o *Comment) GetCommentDate() *DateTime
```
GetCommentDate returns the value of Comment.CommentDate

#### func (*Comment) GetCommentDateGmt

```go
func (o *Comment) GetCommentDateGmt() *DateTime
```
GetCommentDateGmt returns the value of Comment.CommentDateGmt

#### func (*Comment) GetCommentID

```go
func (o *Comment) GetCommentID() int64
```
GetCommentID returns the value of Comment.CommentID

#### func (*Comment) GetCommentKarma

```go
func (o *Comment) GetCommentKarma() int
```
GetCommentKarma returns the value of Comment.CommentKarma

#### func (*Comment) GetCommentParent

```go
func (o *Comment) GetCommentParent() int64
```
GetCommentParent returns the value of Comment.CommentParent

#### func (*Comment) GetCommentPostID

```go
func (o *Comment) GetCommentPostID() int64
```
GetCommentPostID returns the value of Comment.CommentPostID

#### func (*Comment) GetCommentType

```go
func (o *Comment) GetCommentType() string
```
GetCommentType returns the value of Comment.CommentType

#### func (*Comment) GetPrimaryKeyName

```go
func (o *Comment) GetPrimaryKeyName() string
```
GetPrimaryKeyName returns the DB field name

#### func (*Comment) GetPrimaryKeyValue

```go
func (o *Comment) GetPrimaryKeyValue() int64
```
GetPrimaryKeyValue returns the value, usually int64 of the PrimaryKey

#### func (*Comment) GetUserId

```go
func (o *Comment) GetUserId() int64
```
GetUserId returns the value of Comment.UserId

#### func (*Comment) Reload

```go
func (o *Comment) Reload() error
```
Reload A function to forcibly reload Comment

#### func (*Comment) Save

```go
func (o *Comment) Save() error
```
Save is a dynamic saver 'inherited' by all models

#### func (*Comment) SetCommentAgent

```go
func (o *Comment) SetCommentAgent(arg string)
```
SetCommentAgent sets and marks as dirty the value of Comment.CommentAgent

#### func (*Comment) SetCommentApproved

```go
func (o *Comment) SetCommentApproved(arg string)
```
SetCommentApproved sets and marks as dirty the value of Comment.CommentApproved

#### func (*Comment) SetCommentAuthor

```go
func (o *Comment) SetCommentAuthor(arg string)
```
SetCommentAuthor sets and marks as dirty the value of Comment.CommentAuthor

#### func (*Comment) SetCommentAuthorEmail

```go
func (o *Comment) SetCommentAuthorEmail(arg string)
```
SetCommentAuthorEmail sets and marks as dirty the value of
Comment.CommentAuthorEmail

#### func (*Comment) SetCommentAuthorIP

```go
func (o *Comment) SetCommentAuthorIP(arg string)
```
SetCommentAuthorIP sets and marks as dirty the value of Comment.CommentAuthorIP

#### func (*Comment) SetCommentAuthorUrl

```go
func (o *Comment) SetCommentAuthorUrl(arg string)
```
SetCommentAuthorUrl sets and marks as dirty the value of
Comment.CommentAuthorUrl

#### func (*Comment) SetCommentContent

```go
func (o *Comment) SetCommentContent(arg string)
```
SetCommentContent sets and marks as dirty the value of Comment.CommentContent

#### func (*Comment) SetCommentDate

```go
func (o *Comment) SetCommentDate(arg *DateTime)
```
SetCommentDate sets and marks as dirty the value of Comment.CommentDate

#### func (*Comment) SetCommentDateGmt

```go
func (o *Comment) SetCommentDateGmt(arg *DateTime)
```
SetCommentDateGmt sets and marks as dirty the value of Comment.CommentDateGmt

#### func (*Comment) SetCommentID

```go
func (o *Comment) SetCommentID(arg int64)
```
SetCommentID sets and marks as dirty the value of Comment.CommentID

#### func (*Comment) SetCommentKarma

```go
func (o *Comment) SetCommentKarma(arg int)
```
SetCommentKarma sets and marks as dirty the value of Comment.CommentKarma

#### func (*Comment) SetCommentParent

```go
func (o *Comment) SetCommentParent(arg int64)
```
SetCommentParent sets and marks as dirty the value of Comment.CommentParent

#### func (*Comment) SetCommentPostID

```go
func (o *Comment) SetCommentPostID(arg int64)
```
SetCommentPostID sets and marks as dirty the value of Comment.CommentPostID

#### func (*Comment) SetCommentType

```go
func (o *Comment) SetCommentType(arg string)
```
SetCommentType sets and marks as dirty the value of Comment.CommentType

#### func (*Comment) SetUserId

```go
func (o *Comment) SetUserId(arg int64)
```
SetUserId sets and marks as dirty the value of Comment.UserId

#### func (*Comment) Update

```go
func (o *Comment) Update() error
```
Update is a dynamic updater, it considers whether or not a field is 'dirty' and
needs to be updated. Will only work if you use the Getters and Setters

#### func (*Comment) UpdateCommentAgent

```go
func (o *Comment) UpdateCommentAgent(_updCommentAgent string) (int64, error)
```
UpdateCommentAgent an immediate DB Query to update a single column, in this case
comment_agent

#### func (*Comment) UpdateCommentApproved

```go
func (o *Comment) UpdateCommentApproved(_updCommentApproved string) (int64, error)
```
UpdateCommentApproved an immediate DB Query to update a single column, in this
case comment_approved

#### func (*Comment) UpdateCommentAuthor

```go
func (o *Comment) UpdateCommentAuthor(_updCommentAuthor string) (int64, error)
```
UpdateCommentAuthor an immediate DB Query to update a single column, in this
case comment_author

#### func (*Comment) UpdateCommentAuthorEmail

```go
func (o *Comment) UpdateCommentAuthorEmail(_updCommentAuthorEmail string) (int64, error)
```
UpdateCommentAuthorEmail an immediate DB Query to update a single column, in
this case comment_author_email

#### func (*Comment) UpdateCommentAuthorIP

```go
func (o *Comment) UpdateCommentAuthorIP(_updCommentAuthorIP string) (int64, error)
```
UpdateCommentAuthorIP an immediate DB Query to update a single column, in this
case comment_author_IP

#### func (*Comment) UpdateCommentAuthorUrl

```go
func (o *Comment) UpdateCommentAuthorUrl(_updCommentAuthorUrl string) (int64, error)
```
UpdateCommentAuthorUrl an immediate DB Query to update a single column, in this
case comment_author_url

#### func (*Comment) UpdateCommentContent

```go
func (o *Comment) UpdateCommentContent(_updCommentContent string) (int64, error)
```
UpdateCommentContent an immediate DB Query to update a single column, in this
case comment_content

#### func (*Comment) UpdateCommentDate

```go
func (o *Comment) UpdateCommentDate(_updCommentDate *DateTime) (int64, error)
```
UpdateCommentDate an immediate DB Query to update a single column, in this case
comment_date

#### func (*Comment) UpdateCommentDateGmt

```go
func (o *Comment) UpdateCommentDateGmt(_updCommentDateGmt *DateTime) (int64, error)
```
UpdateCommentDateGmt an immediate DB Query to update a single column, in this
case comment_date_gmt

#### func (*Comment) UpdateCommentKarma

```go
func (o *Comment) UpdateCommentKarma(_updCommentKarma int) (int64, error)
```
UpdateCommentKarma an immediate DB Query to update a single column, in this case
comment_karma

#### func (*Comment) UpdateCommentParent

```go
func (o *Comment) UpdateCommentParent(_updCommentParent int64) (int64, error)
```
UpdateCommentParent an immediate DB Query to update a single column, in this
case comment_parent

#### func (*Comment) UpdateCommentPostID

```go
func (o *Comment) UpdateCommentPostID(_updCommentPostID int64) (int64, error)
```
UpdateCommentPostID an immediate DB Query to update a single column, in this
case comment_post_ID

#### func (*Comment) UpdateCommentType

```go
func (o *Comment) UpdateCommentType(_updCommentType string) (int64, error)
```
UpdateCommentType an immediate DB Query to update a single column, in this case
comment_type

#### func (*Comment) UpdateUserId

```go
func (o *Comment) UpdateUserId(_updUserId int64) (int64, error)
```
UpdateUserId an immediate DB Query to update a single column, in this case
user_id

#### type CommentMeta

```go
type CommentMeta struct {
	MetaId    int64
	CommentId int64
	MetaKey   string
	MetaValue string
	// Dirty markers for smart updates
	IsMetaIdDirty    bool
	IsCommentIdDirty bool
	IsMetaKeyDirty   bool
	IsMetaValueDirty bool
}
```

CommentMeta is a Object Relational Mapping to the database table that represents
it. In this case it is commentmeta. The table name will be Sprintf'd to include
the prefix you define in your YAML configuration for the Adapter.

#### func  NewCommentMeta

```go
func NewCommentMeta(a Adapter) *CommentMeta
```
NewCommentMeta binds an Adapter to a new instance of CommentMeta and sets up the
_table and primary keys

#### func (*CommentMeta) Create

```go
func (o *CommentMeta) Create() error
```
Create inserts the model. Calling Save will call this function automatically for
new models

#### func (*CommentMeta) Find

```go
func (o *CommentMeta) Find(_findByMetaId int64) (bool, error)
```
Find dynamic finder for meta_id -> bool,error Generic and programatically
generator finder for CommentMeta

Note that Fine returns a bool if found, not err, in the case of a return of
true, the instance data will be filled out. a call to find ALWAYS overwrites the
model you call Find on i.e. receiver is a pointer.

```go

    m := NewCommentMeta(a)
    found,err := m.Find(23)
    .. handle err
    if found == false {
        // handle found
    }
    ... do what you want with m here

```

#### func (*CommentMeta) FindByCommentId

```go
func (o *CommentMeta) FindByCommentId(_findByCommentId int64) ([]*CommentMeta, error)
```
FindByCommentId dynamic finder for comment_id -> []*CommentMeta,error Generic
and programatically generator finder for CommentMeta

```go

    m := NewCommentMeta(a)
    results,err := m.FindByCommentId(...)
    // handle err
    for i,r := results {
      // now r is an instance of CommentMeta
    }

```

#### func (*CommentMeta) FindByMetaKey

```go
func (o *CommentMeta) FindByMetaKey(_findByMetaKey string) ([]*CommentMeta, error)
```
FindByMetaKey dynamic finder for meta_key -> []*CommentMeta,error Generic and
programatically generator finder for CommentMeta

```go

    m := NewCommentMeta(a)
    results,err := m.FindByMetaKey(...)
    // handle err
    for i,r := results {
      // now r is an instance of CommentMeta
    }

```

#### func (*CommentMeta) FindByMetaValue

```go
func (o *CommentMeta) FindByMetaValue(_findByMetaValue string) ([]*CommentMeta, error)
```
FindByMetaValue dynamic finder for meta_value -> []*CommentMeta,error Generic
and programatically generator finder for CommentMeta

```go

    m := NewCommentMeta(a)
    results,err := m.FindByMetaValue(...)
    // handle err
    for i,r := results {
      // now r is an instance of CommentMeta
    }

```

#### func (*CommentMeta) FromCommentMeta

```go
func (o *CommentMeta) FromCommentMeta(m *CommentMeta)
```
FromCommentMeta A kind of Clone function for CommentMeta

#### func (*CommentMeta) FromDBValueMap

```go
func (o *CommentMeta) FromDBValueMap(m map[string]DBValue) error
```
FromDBValueMap Converts a DBValueMap returned from Adapter.Query to a
CommentMeta

#### func (*CommentMeta) GetCommentId

```go
func (o *CommentMeta) GetCommentId() int64
```
GetCommentId returns the value of CommentMeta.CommentId

#### func (*CommentMeta) GetMetaId

```go
func (o *CommentMeta) GetMetaId() int64
```
GetMetaId returns the value of CommentMeta.MetaId

#### func (*CommentMeta) GetMetaKey

```go
func (o *CommentMeta) GetMetaKey() string
```
GetMetaKey returns the value of CommentMeta.MetaKey

#### func (*CommentMeta) GetMetaValue

```go
func (o *CommentMeta) GetMetaValue() string
```
GetMetaValue returns the value of CommentMeta.MetaValue

#### func (*CommentMeta) GetPrimaryKeyName

```go
func (o *CommentMeta) GetPrimaryKeyName() string
```
GetPrimaryKeyName returns the DB field name

#### func (*CommentMeta) GetPrimaryKeyValue

```go
func (o *CommentMeta) GetPrimaryKeyValue() int64
```
GetPrimaryKeyValue returns the value, usually int64 of the PrimaryKey

#### func (*CommentMeta) Reload

```go
func (o *CommentMeta) Reload() error
```
Reload A function to forcibly reload CommentMeta

#### func (*CommentMeta) Save

```go
func (o *CommentMeta) Save() error
```
Save is a dynamic saver 'inherited' by all models

#### func (*CommentMeta) SetCommentId

```go
func (o *CommentMeta) SetCommentId(arg int64)
```
SetCommentId sets and marks as dirty the value of CommentMeta.CommentId

#### func (*CommentMeta) SetMetaId

```go
func (o *CommentMeta) SetMetaId(arg int64)
```
SetMetaId sets and marks as dirty the value of CommentMeta.MetaId

#### func (*CommentMeta) SetMetaKey

```go
func (o *CommentMeta) SetMetaKey(arg string)
```
SetMetaKey sets and marks as dirty the value of CommentMeta.MetaKey

#### func (*CommentMeta) SetMetaValue

```go
func (o *CommentMeta) SetMetaValue(arg string)
```
SetMetaValue sets and marks as dirty the value of CommentMeta.MetaValue

#### func (*CommentMeta) Update

```go
func (o *CommentMeta) Update() error
```
Update is a dynamic updater, it considers whether or not a field is 'dirty' and
needs to be updated. Will only work if you use the Getters and Setters

#### func (*CommentMeta) UpdateCommentId

```go
func (o *CommentMeta) UpdateCommentId(_updCommentId int64) (int64, error)
```
UpdateCommentId an immediate DB Query to update a single column, in this case
comment_id

#### func (*CommentMeta) UpdateMetaKey

```go
func (o *CommentMeta) UpdateMetaKey(_updMetaKey string) (int64, error)
```
UpdateMetaKey an immediate DB Query to update a single column, in this case
meta_key

#### func (*CommentMeta) UpdateMetaValue

```go
func (o *CommentMeta) UpdateMetaValue(_updMetaValue string) (int64, error)
```
UpdateMetaValue an immediate DB Query to update a single column, in this case
meta_value

#### type DBValue

```go
type DBValue interface {
	AsInt() (int, error)
	AsInt32() (int32, error)
	AsInt64() (int64, error)
	AsFloat32() (float32, error)
	AsFloat64() (float64, error)
	AsString() (string, error)
	AsDateTime() (*DateTime, error)
	SetInternalValue(string, string)
}
```

DBValue Provides a tidy way to convert string values from the DB into go values

#### type DateTime

```go
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
}
```

DateTime A simple struct to represent DateTime fields

#### func  NewDateTime

```go
func NewDateTime(a Adapter) *DateTime
```
NewDateTime Returns a basic DateTime value

#### func (*DateTime) FromString

```go
func (d *DateTime) FromString(s string) error
```
FromString Converts a string like 0000-00-00 00:00:00 into a DateTime

#### func (*DateTime) String

```go
func (d *DateTime) String() string
```
String The Stringer for DateTime to avoid having to call ToString all the time.

#### func (*DateTime) ToString

```go
func (d *DateTime) ToString() string
```
ToString For backwards compat... Never use this, use String() instead.

#### type Link

```go
type Link struct {
	LinkId          int64
	LinkUrl         string
	LinkName        string
	LinkImage       string
	LinkTarget      string
	LinkDescription string
	LinkVisible     string
	LinkOwner       int64
	LinkRating      int
	LinkUpdated     *DateTime
	LinkRel         string
	LinkNotes       string
	LinkRss         string
	// Dirty markers for smart updates
	IsLinkIdDirty          bool
	IsLinkUrlDirty         bool
	IsLinkNameDirty        bool
	IsLinkImageDirty       bool
	IsLinkTargetDirty      bool
	IsLinkDescriptionDirty bool
	IsLinkVisibleDirty     bool
	IsLinkOwnerDirty       bool
	IsLinkRatingDirty      bool
	IsLinkUpdatedDirty     bool
	IsLinkRelDirty         bool
	IsLinkNotesDirty       bool
	IsLinkRssDirty         bool
}
```

Link is a Object Relational Mapping to the database table that represents it. In
this case it is links. The table name will be Sprintf'd to include the prefix
you define in your YAML configuration for the Adapter.

#### func  NewLink

```go
func NewLink(a Adapter) *Link
```
NewLink binds an Adapter to a new instance of Link and sets up the _table and
primary keys

#### func (*Link) Create

```go
func (o *Link) Create() error
```
Create inserts the model. Calling Save will call this function automatically for
new models

#### func (*Link) Find

```go
func (o *Link) Find(_findByLinkId int64) (bool, error)
```
Find dynamic finder for link_id -> bool,error Generic and programatically
generator finder for Link

Note that Fine returns a bool if found, not err, in the case of a return of
true, the instance data will be filled out. a call to find ALWAYS overwrites the
model you call Find on i.e. receiver is a pointer.

```go

    m := NewLink(a)
    found,err := m.Find(23)
    .. handle err
    if found == false {
        // handle found
    }
    ... do what you want with m here

```

#### func (*Link) FindByLinkDescription

```go
func (o *Link) FindByLinkDescription(_findByLinkDescription string) ([]*Link, error)
```
FindByLinkDescription dynamic finder for link_description -> []*Link,error
Generic and programatically generator finder for Link

```go

    m := NewLink(a)
    results,err := m.FindByLinkDescription(...)
    // handle err
    for i,r := results {
      // now r is an instance of Link
    }

```

#### func (*Link) FindByLinkImage

```go
func (o *Link) FindByLinkImage(_findByLinkImage string) ([]*Link, error)
```
FindByLinkImage dynamic finder for link_image -> []*Link,error Generic and
programatically generator finder for Link

```go

    m := NewLink(a)
    results,err := m.FindByLinkImage(...)
    // handle err
    for i,r := results {
      // now r is an instance of Link
    }

```

#### func (*Link) FindByLinkName

```go
func (o *Link) FindByLinkName(_findByLinkName string) ([]*Link, error)
```
FindByLinkName dynamic finder for link_name -> []*Link,error Generic and
programatically generator finder for Link

```go

    m := NewLink(a)
    results,err := m.FindByLinkName(...)
    // handle err
    for i,r := results {
      // now r is an instance of Link
    }

```

#### func (*Link) FindByLinkNotes

```go
func (o *Link) FindByLinkNotes(_findByLinkNotes string) ([]*Link, error)
```
FindByLinkNotes dynamic finder for link_notes -> []*Link,error Generic and
programatically generator finder for Link

```go

    m := NewLink(a)
    results,err := m.FindByLinkNotes(...)
    // handle err
    for i,r := results {
      // now r is an instance of Link
    }

```

#### func (*Link) FindByLinkOwner

```go
func (o *Link) FindByLinkOwner(_findByLinkOwner int64) ([]*Link, error)
```
FindByLinkOwner dynamic finder for link_owner -> []*Link,error Generic and
programatically generator finder for Link

```go

    m := NewLink(a)
    results,err := m.FindByLinkOwner(...)
    // handle err
    for i,r := results {
      // now r is an instance of Link
    }

```

#### func (*Link) FindByLinkRating

```go
func (o *Link) FindByLinkRating(_findByLinkRating int) ([]*Link, error)
```
FindByLinkRating dynamic finder for link_rating -> []*Link,error Generic and
programatically generator finder for Link

```go

    m := NewLink(a)
    results,err := m.FindByLinkRating(...)
    // handle err
    for i,r := results {
      // now r is an instance of Link
    }

```

#### func (*Link) FindByLinkRel

```go
func (o *Link) FindByLinkRel(_findByLinkRel string) ([]*Link, error)
```
FindByLinkRel dynamic finder for link_rel -> []*Link,error Generic and
programatically generator finder for Link

```go

    m := NewLink(a)
    results,err := m.FindByLinkRel(...)
    // handle err
    for i,r := results {
      // now r is an instance of Link
    }

```

#### func (*Link) FindByLinkRss

```go
func (o *Link) FindByLinkRss(_findByLinkRss string) ([]*Link, error)
```
FindByLinkRss dynamic finder for link_rss -> []*Link,error Generic and
programatically generator finder for Link

```go

    m := NewLink(a)
    results,err := m.FindByLinkRss(...)
    // handle err
    for i,r := results {
      // now r is an instance of Link
    }

```

#### func (*Link) FindByLinkTarget

```go
func (o *Link) FindByLinkTarget(_findByLinkTarget string) ([]*Link, error)
```
FindByLinkTarget dynamic finder for link_target -> []*Link,error Generic and
programatically generator finder for Link

```go

    m := NewLink(a)
    results,err := m.FindByLinkTarget(...)
    // handle err
    for i,r := results {
      // now r is an instance of Link
    }

```

#### func (*Link) FindByLinkUpdated

```go
func (o *Link) FindByLinkUpdated(_findByLinkUpdated *DateTime) ([]*Link, error)
```
FindByLinkUpdated dynamic finder for link_updated -> []*Link,error Generic and
programatically generator finder for Link

```go

    m := NewLink(a)
    results,err := m.FindByLinkUpdated(...)
    // handle err
    for i,r := results {
      // now r is an instance of Link
    }

```

#### func (*Link) FindByLinkUrl

```go
func (o *Link) FindByLinkUrl(_findByLinkUrl string) ([]*Link, error)
```
FindByLinkUrl dynamic finder for link_url -> []*Link,error Generic and
programatically generator finder for Link

```go

    m := NewLink(a)
    results,err := m.FindByLinkUrl(...)
    // handle err
    for i,r := results {
      // now r is an instance of Link
    }

```

#### func (*Link) FindByLinkVisible

```go
func (o *Link) FindByLinkVisible(_findByLinkVisible string) ([]*Link, error)
```
FindByLinkVisible dynamic finder for link_visible -> []*Link,error Generic and
programatically generator finder for Link

```go

    m := NewLink(a)
    results,err := m.FindByLinkVisible(...)
    // handle err
    for i,r := results {
      // now r is an instance of Link
    }

```

#### func (*Link) FromDBValueMap

```go
func (o *Link) FromDBValueMap(m map[string]DBValue) error
```
FromDBValueMap Converts a DBValueMap returned from Adapter.Query to a Link

#### func (*Link) FromLink

```go
func (o *Link) FromLink(m *Link)
```
FromLink A kind of Clone function for Link

#### func (*Link) GetLinkDescription

```go
func (o *Link) GetLinkDescription() string
```
GetLinkDescription returns the value of Link.LinkDescription

#### func (*Link) GetLinkId

```go
func (o *Link) GetLinkId() int64
```
GetLinkId returns the value of Link.LinkId

#### func (*Link) GetLinkImage

```go
func (o *Link) GetLinkImage() string
```
GetLinkImage returns the value of Link.LinkImage

#### func (*Link) GetLinkName

```go
func (o *Link) GetLinkName() string
```
GetLinkName returns the value of Link.LinkName

#### func (*Link) GetLinkNotes

```go
func (o *Link) GetLinkNotes() string
```
GetLinkNotes returns the value of Link.LinkNotes

#### func (*Link) GetLinkOwner

```go
func (o *Link) GetLinkOwner() int64
```
GetLinkOwner returns the value of Link.LinkOwner

#### func (*Link) GetLinkRating

```go
func (o *Link) GetLinkRating() int
```
GetLinkRating returns the value of Link.LinkRating

#### func (*Link) GetLinkRel

```go
func (o *Link) GetLinkRel() string
```
GetLinkRel returns the value of Link.LinkRel

#### func (*Link) GetLinkRss

```go
func (o *Link) GetLinkRss() string
```
GetLinkRss returns the value of Link.LinkRss

#### func (*Link) GetLinkTarget

```go
func (o *Link) GetLinkTarget() string
```
GetLinkTarget returns the value of Link.LinkTarget

#### func (*Link) GetLinkUpdated

```go
func (o *Link) GetLinkUpdated() *DateTime
```
GetLinkUpdated returns the value of Link.LinkUpdated

#### func (*Link) GetLinkUrl

```go
func (o *Link) GetLinkUrl() string
```
GetLinkUrl returns the value of Link.LinkUrl

#### func (*Link) GetLinkVisible

```go
func (o *Link) GetLinkVisible() string
```
GetLinkVisible returns the value of Link.LinkVisible

#### func (*Link) GetPrimaryKeyName

```go
func (o *Link) GetPrimaryKeyName() string
```
GetPrimaryKeyName returns the DB field name

#### func (*Link) GetPrimaryKeyValue

```go
func (o *Link) GetPrimaryKeyValue() int64
```
GetPrimaryKeyValue returns the value, usually int64 of the PrimaryKey

#### func (*Link) Reload

```go
func (o *Link) Reload() error
```
Reload A function to forcibly reload Link

#### func (*Link) Save

```go
func (o *Link) Save() error
```
Save is a dynamic saver 'inherited' by all models

#### func (*Link) SetLinkDescription

```go
func (o *Link) SetLinkDescription(arg string)
```
SetLinkDescription sets and marks as dirty the value of Link.LinkDescription

#### func (*Link) SetLinkId

```go
func (o *Link) SetLinkId(arg int64)
```
SetLinkId sets and marks as dirty the value of Link.LinkId

#### func (*Link) SetLinkImage

```go
func (o *Link) SetLinkImage(arg string)
```
SetLinkImage sets and marks as dirty the value of Link.LinkImage

#### func (*Link) SetLinkName

```go
func (o *Link) SetLinkName(arg string)
```
SetLinkName sets and marks as dirty the value of Link.LinkName

#### func (*Link) SetLinkNotes

```go
func (o *Link) SetLinkNotes(arg string)
```
SetLinkNotes sets and marks as dirty the value of Link.LinkNotes

#### func (*Link) SetLinkOwner

```go
func (o *Link) SetLinkOwner(arg int64)
```
SetLinkOwner sets and marks as dirty the value of Link.LinkOwner

#### func (*Link) SetLinkRating

```go
func (o *Link) SetLinkRating(arg int)
```
SetLinkRating sets and marks as dirty the value of Link.LinkRating

#### func (*Link) SetLinkRel

```go
func (o *Link) SetLinkRel(arg string)
```
SetLinkRel sets and marks as dirty the value of Link.LinkRel

#### func (*Link) SetLinkRss

```go
func (o *Link) SetLinkRss(arg string)
```
SetLinkRss sets and marks as dirty the value of Link.LinkRss

#### func (*Link) SetLinkTarget

```go
func (o *Link) SetLinkTarget(arg string)
```
SetLinkTarget sets and marks as dirty the value of Link.LinkTarget

#### func (*Link) SetLinkUpdated

```go
func (o *Link) SetLinkUpdated(arg *DateTime)
```
SetLinkUpdated sets and marks as dirty the value of Link.LinkUpdated

#### func (*Link) SetLinkUrl

```go
func (o *Link) SetLinkUrl(arg string)
```
SetLinkUrl sets and marks as dirty the value of Link.LinkUrl

#### func (*Link) SetLinkVisible

```go
func (o *Link) SetLinkVisible(arg string)
```
SetLinkVisible sets and marks as dirty the value of Link.LinkVisible

#### func (*Link) Update

```go
func (o *Link) Update() error
```
Update is a dynamic updater, it considers whether or not a field is 'dirty' and
needs to be updated. Will only work if you use the Getters and Setters

#### func (*Link) UpdateLinkDescription

```go
func (o *Link) UpdateLinkDescription(_updLinkDescription string) (int64, error)
```
UpdateLinkDescription an immediate DB Query to update a single column, in this
case link_description

#### func (*Link) UpdateLinkImage

```go
func (o *Link) UpdateLinkImage(_updLinkImage string) (int64, error)
```
UpdateLinkImage an immediate DB Query to update a single column, in this case
link_image

#### func (*Link) UpdateLinkName

```go
func (o *Link) UpdateLinkName(_updLinkName string) (int64, error)
```
UpdateLinkName an immediate DB Query to update a single column, in this case
link_name

#### func (*Link) UpdateLinkNotes

```go
func (o *Link) UpdateLinkNotes(_updLinkNotes string) (int64, error)
```
UpdateLinkNotes an immediate DB Query to update a single column, in this case
link_notes

#### func (*Link) UpdateLinkOwner

```go
func (o *Link) UpdateLinkOwner(_updLinkOwner int64) (int64, error)
```
UpdateLinkOwner an immediate DB Query to update a single column, in this case
link_owner

#### func (*Link) UpdateLinkRating

```go
func (o *Link) UpdateLinkRating(_updLinkRating int) (int64, error)
```
UpdateLinkRating an immediate DB Query to update a single column, in this case
link_rating

#### func (*Link) UpdateLinkRel

```go
func (o *Link) UpdateLinkRel(_updLinkRel string) (int64, error)
```
UpdateLinkRel an immediate DB Query to update a single column, in this case
link_rel

#### func (*Link) UpdateLinkRss

```go
func (o *Link) UpdateLinkRss(_updLinkRss string) (int64, error)
```
UpdateLinkRss an immediate DB Query to update a single column, in this case
link_rss

#### func (*Link) UpdateLinkTarget

```go
func (o *Link) UpdateLinkTarget(_updLinkTarget string) (int64, error)
```
UpdateLinkTarget an immediate DB Query to update a single column, in this case
link_target

#### func (*Link) UpdateLinkUpdated

```go
func (o *Link) UpdateLinkUpdated(_updLinkUpdated *DateTime) (int64, error)
```
UpdateLinkUpdated an immediate DB Query to update a single column, in this case
link_updated

#### func (*Link) UpdateLinkUrl

```go
func (o *Link) UpdateLinkUrl(_updLinkUrl string) (int64, error)
```
UpdateLinkUrl an immediate DB Query to update a single column, in this case
link_url

#### func (*Link) UpdateLinkVisible

```go
func (o *Link) UpdateLinkVisible(_updLinkVisible string) (int64, error)
```
UpdateLinkVisible an immediate DB Query to update a single column, in this case
link_visible

#### type LogFilter

```go
type LogFilter func(string, string) string
```

LogFilter is an anonymous function that that receives the log tag and string and
allows you to filter out extraneous lines when trying to find bugs.

#### type MysqlAdapter

```go
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
}
```

MysqlAdapter is the MySql implementation

#### func  NewMysqlAdapter

```go
func NewMysqlAdapter(pre string) *MysqlAdapter
```
NewMysqlAdapter returns a pointer to MysqlAdapter

#### func  NewMysqlAdapterEx

```go
func NewMysqlAdapterEx(fname string) (*MysqlAdapter, error)
```
NewMysqlAdapterEx sets everything up based on your YAML config Args: fname is a
string path to a YAML config file This function will attempt to Open the
database defined in that file. Example file:

    host: "localhost"
    user: "dbuser"
    pass: "dbuserpass"
    database: "my_db"
    prefix: "wp_"

#### func (*MysqlAdapter) AffectedRows

```go
func (a *MysqlAdapter) AffectedRows() int64
```
AffectedRows Grab the number of AffectedRows

#### func (*MysqlAdapter) Close

```go
func (a *MysqlAdapter) Close()
```
Close This should be called in your application with a defer a.Close() or
something similar. Closing is not automatic!

#### func (*MysqlAdapter) DatabasePrefix

```go
func (a *MysqlAdapter) DatabasePrefix() string
```
DatabasePrefix Get the DatabasePrefix from the Adapter

#### func (*MysqlAdapter) Execute

```go
func (a *MysqlAdapter) Execute(q string) error
```
Execute For UPDATE and INSERT calls, i.e. nothing that returns a result set.

#### func (*MysqlAdapter) FromYAML

```go
func (a *MysqlAdapter) FromYAML(b []byte) error
```
FromYAML Set the Adapter's members from a YAML file

#### func (*MysqlAdapter) LastInsertedId

```go
func (a *MysqlAdapter) LastInsertedId() int64
```
LastInsertedId Grab the last auto_incremented id

#### func (*MysqlAdapter) LogDebug

```go
func (a *MysqlAdapter) LogDebug(s string)
```
LogDebug Tags the string with DEBUG and puts it into _debugLog.

#### func (*MysqlAdapter) LogError

```go
func (a *MysqlAdapter) LogError(s error)
```
LogError Tags the string with ERROR and puts it into _errorLog.

#### func (*MysqlAdapter) LogInfo

```go
func (a *MysqlAdapter) LogInfo(s string)
```
LogInfo Tags the string with INFO and puts it into _infoLog.

#### func (*MysqlAdapter) NewDBValue

```go
func (a *MysqlAdapter) NewDBValue() DBValue
```
NewDBValue Creates a new DBValue, mostly used internally, but you may wish to
use it in special circumstances.

#### func (*MysqlAdapter) Oops

```go
func (a *MysqlAdapter) Oops(s string) error
```
Oops A function for catching errors generated by the library and funneling them
to the log files

#### func (*MysqlAdapter) Open

```go
func (a *MysqlAdapter) Open(h, u, p, d string) error
```
Open Opens the database connection. Be sure to use a.Close() as closing is NOT
handled for you.

#### func (*MysqlAdapter) Query

```go
func (a *MysqlAdapter) Query(q string) ([]map[string]DBValue, error)
```
Query The generay Query function, i.e. SQL that returns results, as opposed to
an INSERT or UPDATE which uses Execute.

#### func (*MysqlAdapter) SafeString

```go
func (a *MysqlAdapter) SafeString(s string) string
```
SafeString Not implemented yet, but soon.

#### func (*MysqlAdapter) SetDebugLog

```go
func (a *MysqlAdapter) SetDebugLog(t io.Writer)
```
SetDebugLog Sets the _debugLog to the io.Writer, use ioutil.Discard if you don't
want this one at all.

#### func (*MysqlAdapter) SetErrorLog

```go
func (a *MysqlAdapter) SetErrorLog(t io.Writer)
```
SetErrorLog Sets the _errorLog to the io.Writer, use ioutil.Discard if you don't
want this one at all.

#### func (*MysqlAdapter) SetInfoLog

```go
func (a *MysqlAdapter) SetInfoLog(t io.Writer)
```
SetInfoLog Sets the _infoLog to the io.Writer, use ioutil.Discard if you don't
want this one at all.

#### func (*MysqlAdapter) SetLogFilter

```go
func (a *MysqlAdapter) SetLogFilter(f LogFilter)
```
SetLogFilter sets the LogFilter to a function. This is only useful if you are
debugging, or you want to reformat the log data.

#### func (*MysqlAdapter) SetLogs

```go
func (a *MysqlAdapter) SetLogs(t io.Writer)
```
SetLogs Sets ALL logs to the io.Writer, use ioutil.Discard if you don't want
this one at all.

#### type MysqlValue

```go
type MysqlValue struct {
}
```

MysqlValue Implements DBValue for MySQL, you'll generally not interact directly
with this type, but it is there for special cases.

#### func  NewMysqlValue

```go
func NewMysqlValue(a Adapter) *MysqlValue
```
NewMysqlValue A function for largely internal use, but basically in order to use
a DBValue, it needs to have its Adapter setup, this is because some values have
Adapter specific issues. The implementing adapter may need to provide some
information, or logging etc

#### func (*MysqlValue) AsDateTime

```go
func (v *MysqlValue) AsDateTime() (*DateTime, error)
```
AsDateTime Tries to convert the string to a DateTime, parsing may fail.

#### func (*MysqlValue) AsFloat32

```go
func (v *MysqlValue) AsFloat32() (float32, error)
```
AsFloat32 Tries to convert the internal string to a float32

#### func (*MysqlValue) AsFloat64

```go
func (v *MysqlValue) AsFloat64() (float64, error)
```
AsFloat64 Tries to convert the internal string to a float64

#### func (*MysqlValue) AsInt

```go
func (v *MysqlValue) AsInt() (int, error)
```
AsInt Attempts to convert the internal string to an Int

#### func (*MysqlValue) AsInt32

```go
func (v *MysqlValue) AsInt32() (int32, error)
```
AsInt32 Tries to convert the internal string to an int32

#### func (*MysqlValue) AsInt64

```go
func (v *MysqlValue) AsInt64() (int64, error)
```
AsInt64 Tries to convert the internal string to an int64 (i.e. BIGINT)

#### func (*MysqlValue) AsString

```go
func (v *MysqlValue) AsString() (string, error)
```
AsString Simply returns the internal string representation.

#### func (*MysqlValue) SetInternalValue

```go
func (v *MysqlValue) SetInternalValue(key, value string)
```
SetInternalValue Sets the internal value of the DBValue to the string provided.
key isn't really used, but it may be.

#### type Option

```go
type Option struct {
	OptionId    int64
	OptionName  string
	OptionValue string
	Autoload    string
	// Dirty markers for smart updates
	IsOptionIdDirty    bool
	IsOptionNameDirty  bool
	IsOptionValueDirty bool
	IsAutoloadDirty    bool
}
```

Option is a Object Relational Mapping to the database table that represents it.
In this case it is options. The table name will be Sprintf'd to include the
prefix you define in your YAML configuration for the Adapter.

#### func  NewOption

```go
func NewOption(a Adapter) *Option
```
NewOption binds an Adapter to a new instance of Option and sets up the _table
and primary keys

#### func (*Option) Create

```go
func (o *Option) Create() error
```
Create inserts the model. Calling Save will call this function automatically for
new models

#### func (*Option) Find

```go
func (o *Option) Find(_findByOptionId int64) (bool, error)
```
Find dynamic finder for option_id -> bool,error Generic and programatically
generator finder for Option

Note that Fine returns a bool if found, not err, in the case of a return of
true, the instance data will be filled out. a call to find ALWAYS overwrites the
model you call Find on i.e. receiver is a pointer.

```go

    m := NewOption(a)
    found,err := m.Find(23)
    .. handle err
    if found == false {
        // handle found
    }
    ... do what you want with m here

```

#### func (*Option) FindByAutoload

```go
func (o *Option) FindByAutoload(_findByAutoload string) ([]*Option, error)
```
FindByAutoload dynamic finder for autoload -> []*Option,error Generic and
programatically generator finder for Option

```go

    m := NewOption(a)
    results,err := m.FindByAutoload(...)
    // handle err
    for i,r := results {
      // now r is an instance of Option
    }

```

#### func (*Option) FindByOptionName

```go
func (o *Option) FindByOptionName(_findByOptionName string) ([]*Option, error)
```
FindByOptionName dynamic finder for option_name -> []*Option,error Generic and
programatically generator finder for Option

```go

    m := NewOption(a)
    results,err := m.FindByOptionName(...)
    // handle err
    for i,r := results {
      // now r is an instance of Option
    }

```

#### func (*Option) FindByOptionValue

```go
func (o *Option) FindByOptionValue(_findByOptionValue string) ([]*Option, error)
```
FindByOptionValue dynamic finder for option_value -> []*Option,error Generic and
programatically generator finder for Option

```go

    m := NewOption(a)
    results,err := m.FindByOptionValue(...)
    // handle err
    for i,r := results {
      // now r is an instance of Option
    }

```

#### func (*Option) FromDBValueMap

```go
func (o *Option) FromDBValueMap(m map[string]DBValue) error
```
FromDBValueMap Converts a DBValueMap returned from Adapter.Query to a Option

#### func (*Option) FromOption

```go
func (o *Option) FromOption(m *Option)
```
FromOption A kind of Clone function for Option

#### func (*Option) GetAutoload

```go
func (o *Option) GetAutoload() string
```
GetAutoload returns the value of Option.Autoload

#### func (*Option) GetOptionId

```go
func (o *Option) GetOptionId() int64
```
GetOptionId returns the value of Option.OptionId

#### func (*Option) GetOptionName

```go
func (o *Option) GetOptionName() string
```
GetOptionName returns the value of Option.OptionName

#### func (*Option) GetOptionValue

```go
func (o *Option) GetOptionValue() string
```
GetOptionValue returns the value of Option.OptionValue

#### func (*Option) GetPrimaryKeyName

```go
func (o *Option) GetPrimaryKeyName() string
```
GetPrimaryKeyName returns the DB field name

#### func (*Option) GetPrimaryKeyValue

```go
func (o *Option) GetPrimaryKeyValue() int64
```
GetPrimaryKeyValue returns the value, usually int64 of the PrimaryKey

#### func (*Option) Reload

```go
func (o *Option) Reload() error
```
Reload A function to forcibly reload Option

#### func (*Option) Save

```go
func (o *Option) Save() error
```
Save is a dynamic saver 'inherited' by all models

#### func (*Option) SetAutoload

```go
func (o *Option) SetAutoload(arg string)
```
SetAutoload sets and marks as dirty the value of Option.Autoload

#### func (*Option) SetOptionId

```go
func (o *Option) SetOptionId(arg int64)
```
SetOptionId sets and marks as dirty the value of Option.OptionId

#### func (*Option) SetOptionName

```go
func (o *Option) SetOptionName(arg string)
```
SetOptionName sets and marks as dirty the value of Option.OptionName

#### func (*Option) SetOptionValue

```go
func (o *Option) SetOptionValue(arg string)
```
SetOptionValue sets and marks as dirty the value of Option.OptionValue

#### func (*Option) Update

```go
func (o *Option) Update() error
```
Update is a dynamic updater, it considers whether or not a field is 'dirty' and
needs to be updated. Will only work if you use the Getters and Setters

#### func (*Option) UpdateAutoload

```go
func (o *Option) UpdateAutoload(_updAutoload string) (int64, error)
```
UpdateAutoload an immediate DB Query to update a single column, in this case
autoload

#### func (*Option) UpdateOptionName

```go
func (o *Option) UpdateOptionName(_updOptionName string) (int64, error)
```
UpdateOptionName an immediate DB Query to update a single column, in this case
option_name

#### func (*Option) UpdateOptionValue

```go
func (o *Option) UpdateOptionValue(_updOptionValue string) (int64, error)
```
UpdateOptionValue an immediate DB Query to update a single column, in this case
option_value

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
	IsPostDateDirty            bool
	IsPostDateGmtDirty         bool
	IsPostContentDirty         bool
	IsPostTitleDirty           bool
	IsPostExcerptDirty         bool
	IsPostStatusDirty          bool
	IsCommentStatusDirty       bool
	IsPingStatusDirty          bool
	IsPostPasswordDirty        bool
	IsPostNameDirty            bool
	IsToPingDirty              bool
	IsPingedDirty              bool
	IsPostModifiedDirty        bool
	IsPostModifiedGmtDirty     bool
	IsPostContentFilteredDirty bool
	IsPostParentDirty          bool
	IsGuidDirty                bool
	IsMenuOrderDirty           bool
	IsPostTypeDirty            bool
	IsPostMimeTypeDirty        bool
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

#### func (*Post) FindByCommentCount

```go
func (o *Post) FindByCommentCount(_findByCommentCount int64) ([]*Post, error)
```
FindByCommentCount dynamic finder for comment_count -> []*Post,error Generic and
programatically generator finder for Post

```go

    m := NewPost(a)
    results,err := m.FindByCommentCount(...)
    // handle err
    for i,r := results {
      // now r is an instance of Post
    }

```

#### func (*Post) FindByCommentStatus

```go
func (o *Post) FindByCommentStatus(_findByCommentStatus string) ([]*Post, error)
```
FindByCommentStatus dynamic finder for comment_status -> []*Post,error Generic
and programatically generator finder for Post

```go

    m := NewPost(a)
    results,err := m.FindByCommentStatus(...)
    // handle err
    for i,r := results {
      // now r is an instance of Post
    }

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

#### func (*Post) FindByMenuOrder

```go
func (o *Post) FindByMenuOrder(_findByMenuOrder int) ([]*Post, error)
```
FindByMenuOrder dynamic finder for menu_order -> []*Post,error Generic and
programatically generator finder for Post

```go

    m := NewPost(a)
    results,err := m.FindByMenuOrder(...)
    // handle err
    for i,r := results {
      // now r is an instance of Post
    }

```

#### func (*Post) FindByPingStatus

```go
func (o *Post) FindByPingStatus(_findByPingStatus string) ([]*Post, error)
```
FindByPingStatus dynamic finder for ping_status -> []*Post,error Generic and
programatically generator finder for Post

```go

    m := NewPost(a)
    results,err := m.FindByPingStatus(...)
    // handle err
    for i,r := results {
      // now r is an instance of Post
    }

```

#### func (*Post) FindByPinged

```go
func (o *Post) FindByPinged(_findByPinged string) ([]*Post, error)
```
FindByPinged dynamic finder for pinged -> []*Post,error Generic and
programatically generator finder for Post

```go

    m := NewPost(a)
    results,err := m.FindByPinged(...)
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

#### func (*Post) FindByPostContent

```go
func (o *Post) FindByPostContent(_findByPostContent string) ([]*Post, error)
```
FindByPostContent dynamic finder for post_content -> []*Post,error Generic and
programatically generator finder for Post

```go

    m := NewPost(a)
    results,err := m.FindByPostContent(...)
    // handle err
    for i,r := results {
      // now r is an instance of Post
    }

```

#### func (*Post) FindByPostContentFiltered

```go
func (o *Post) FindByPostContentFiltered(_findByPostContentFiltered string) ([]*Post, error)
```
FindByPostContentFiltered dynamic finder for post_content_filtered ->
[]*Post,error Generic and programatically generator finder for Post

```go

    m := NewPost(a)
    results,err := m.FindByPostContentFiltered(...)
    // handle err
    for i,r := results {
      // now r is an instance of Post
    }

```

#### func (*Post) FindByPostDate

```go
func (o *Post) FindByPostDate(_findByPostDate *DateTime) ([]*Post, error)
```
FindByPostDate dynamic finder for post_date -> []*Post,error Generic and
programatically generator finder for Post

```go

    m := NewPost(a)
    results,err := m.FindByPostDate(...)
    // handle err
    for i,r := results {
      // now r is an instance of Post
    }

```

#### func (*Post) FindByPostDateGmt

```go
func (o *Post) FindByPostDateGmt(_findByPostDateGmt *DateTime) ([]*Post, error)
```
FindByPostDateGmt dynamic finder for post_date_gmt -> []*Post,error Generic and
programatically generator finder for Post

```go

    m := NewPost(a)
    results,err := m.FindByPostDateGmt(...)
    // handle err
    for i,r := results {
      // now r is an instance of Post
    }

```

#### func (*Post) FindByPostExcerpt

```go
func (o *Post) FindByPostExcerpt(_findByPostExcerpt string) ([]*Post, error)
```
FindByPostExcerpt dynamic finder for post_excerpt -> []*Post,error Generic and
programatically generator finder for Post

```go

    m := NewPost(a)
    results,err := m.FindByPostExcerpt(...)
    // handle err
    for i,r := results {
      // now r is an instance of Post
    }

```

#### func (*Post) FindByPostMimeType

```go
func (o *Post) FindByPostMimeType(_findByPostMimeType string) ([]*Post, error)
```
FindByPostMimeType dynamic finder for post_mime_type -> []*Post,error Generic
and programatically generator finder for Post

```go

    m := NewPost(a)
    results,err := m.FindByPostMimeType(...)
    // handle err
    for i,r := results {
      // now r is an instance of Post
    }

```

#### func (*Post) FindByPostModified

```go
func (o *Post) FindByPostModified(_findByPostModified *DateTime) ([]*Post, error)
```
FindByPostModified dynamic finder for post_modified -> []*Post,error Generic and
programatically generator finder for Post

```go

    m := NewPost(a)
    results,err := m.FindByPostModified(...)
    // handle err
    for i,r := results {
      // now r is an instance of Post
    }

```

#### func (*Post) FindByPostModifiedGmt

```go
func (o *Post) FindByPostModifiedGmt(_findByPostModifiedGmt *DateTime) ([]*Post, error)
```
FindByPostModifiedGmt dynamic finder for post_modified_gmt -> []*Post,error
Generic and programatically generator finder for Post

```go

    m := NewPost(a)
    results,err := m.FindByPostModifiedGmt(...)
    // handle err
    for i,r := results {
      // now r is an instance of Post
    }

```

#### func (*Post) FindByPostName

```go
func (o *Post) FindByPostName(_findByPostName string) ([]*Post, error)
```
FindByPostName dynamic finder for post_name -> []*Post,error Generic and
programatically generator finder for Post

```go

    m := NewPost(a)
    results,err := m.FindByPostName(...)
    // handle err
    for i,r := results {
      // now r is an instance of Post
    }

```

#### func (*Post) FindByPostParent

```go
func (o *Post) FindByPostParent(_findByPostParent int64) ([]*Post, error)
```
FindByPostParent dynamic finder for post_parent -> []*Post,error Generic and
programatically generator finder for Post

```go

    m := NewPost(a)
    results,err := m.FindByPostParent(...)
    // handle err
    for i,r := results {
      // now r is an instance of Post
    }

```

#### func (*Post) FindByPostPassword

```go
func (o *Post) FindByPostPassword(_findByPostPassword string) ([]*Post, error)
```
FindByPostPassword dynamic finder for post_password -> []*Post,error Generic and
programatically generator finder for Post

```go

    m := NewPost(a)
    results,err := m.FindByPostPassword(...)
    // handle err
    for i,r := results {
      // now r is an instance of Post
    }

```

#### func (*Post) FindByPostStatus

```go
func (o *Post) FindByPostStatus(_findByPostStatus string) ([]*Post, error)
```
FindByPostStatus dynamic finder for post_status -> []*Post,error Generic and
programatically generator finder for Post

```go

    m := NewPost(a)
    results,err := m.FindByPostStatus(...)
    // handle err
    for i,r := results {
      // now r is an instance of Post
    }

```

#### func (*Post) FindByPostTitle

```go
func (o *Post) FindByPostTitle(_findByPostTitle string) ([]*Post, error)
```
FindByPostTitle dynamic finder for post_title -> []*Post,error Generic and
programatically generator finder for Post

```go

    m := NewPost(a)
    results,err := m.FindByPostTitle(...)
    // handle err
    for i,r := results {
      // now r is an instance of Post
    }

```

#### func (*Post) FindByPostType

```go
func (o *Post) FindByPostType(_findByPostType string) ([]*Post, error)
```
FindByPostType dynamic finder for post_type -> []*Post,error Generic and
programatically generator finder for Post

```go

    m := NewPost(a)
    results,err := m.FindByPostType(...)
    // handle err
    for i,r := results {
      // now r is an instance of Post
    }

```

#### func (*Post) FindByToPing

```go
func (o *Post) FindByToPing(_findByToPing string) ([]*Post, error)
```
FindByToPing dynamic finder for to_ping -> []*Post,error Generic and
programatically generator finder for Post

```go

    m := NewPost(a)
    results,err := m.FindByToPing(...)
    // handle err
    for i,r := results {
      // now r is an instance of Post
    }

```

#### func (*Post) FromDBValueMap

```go
func (o *Post) FromDBValueMap(m map[string]DBValue) error
```
FromDBValueMap Converts a DBValueMap returned from Adapter.Query to a Post

#### func (*Post) FromPost

```go
func (o *Post) FromPost(m *Post)
```
FromPost A kind of Clone function for Post

#### func (*Post) GetCommentCount

```go
func (o *Post) GetCommentCount() int64
```
GetCommentCount returns the value of Post.CommentCount

#### func (*Post) GetCommentStatus

```go
func (o *Post) GetCommentStatus() string
```
GetCommentStatus returns the value of Post.CommentStatus

#### func (*Post) GetGuid

```go
func (o *Post) GetGuid() string
```
GetGuid returns the value of Post.Guid

#### func (*Post) GetID

```go
func (o *Post) GetID() int64
```
GetID returns the value of Post.ID

#### func (*Post) GetMenuOrder

```go
func (o *Post) GetMenuOrder() int
```
GetMenuOrder returns the value of Post.MenuOrder

#### func (*Post) GetPingStatus

```go
func (o *Post) GetPingStatus() string
```
GetPingStatus returns the value of Post.PingStatus

#### func (*Post) GetPinged

```go
func (o *Post) GetPinged() string
```
GetPinged returns the value of Post.Pinged

#### func (*Post) GetPostAuthor

```go
func (o *Post) GetPostAuthor() int64
```
GetPostAuthor returns the value of Post.PostAuthor

#### func (*Post) GetPostContent

```go
func (o *Post) GetPostContent() string
```
GetPostContent returns the value of Post.PostContent

#### func (*Post) GetPostContentFiltered

```go
func (o *Post) GetPostContentFiltered() string
```
GetPostContentFiltered returns the value of Post.PostContentFiltered

#### func (*Post) GetPostDate

```go
func (o *Post) GetPostDate() *DateTime
```
GetPostDate returns the value of Post.PostDate

#### func (*Post) GetPostDateGmt

```go
func (o *Post) GetPostDateGmt() *DateTime
```
GetPostDateGmt returns the value of Post.PostDateGmt

#### func (*Post) GetPostExcerpt

```go
func (o *Post) GetPostExcerpt() string
```
GetPostExcerpt returns the value of Post.PostExcerpt

#### func (*Post) GetPostMimeType

```go
func (o *Post) GetPostMimeType() string
```
GetPostMimeType returns the value of Post.PostMimeType

#### func (*Post) GetPostModified

```go
func (o *Post) GetPostModified() *DateTime
```
GetPostModified returns the value of Post.PostModified

#### func (*Post) GetPostModifiedGmt

```go
func (o *Post) GetPostModifiedGmt() *DateTime
```
GetPostModifiedGmt returns the value of Post.PostModifiedGmt

#### func (*Post) GetPostName

```go
func (o *Post) GetPostName() string
```
GetPostName returns the value of Post.PostName

#### func (*Post) GetPostParent

```go
func (o *Post) GetPostParent() int64
```
GetPostParent returns the value of Post.PostParent

#### func (*Post) GetPostPassword

```go
func (o *Post) GetPostPassword() string
```
GetPostPassword returns the value of Post.PostPassword

#### func (*Post) GetPostStatus

```go
func (o *Post) GetPostStatus() string
```
GetPostStatus returns the value of Post.PostStatus

#### func (*Post) GetPostTitle

```go
func (o *Post) GetPostTitle() string
```
GetPostTitle returns the value of Post.PostTitle

#### func (*Post) GetPostType

```go
func (o *Post) GetPostType() string
```
GetPostType returns the value of Post.PostType

#### func (*Post) GetPrimaryKeyName

```go
func (o *Post) GetPrimaryKeyName() string
```
GetPrimaryKeyName returns the DB field name

#### func (*Post) GetPrimaryKeyValue

```go
func (o *Post) GetPrimaryKeyValue() int64
```
GetPrimaryKeyValue returns the value, usually int64 of the PrimaryKey

#### func (*Post) GetToPing

```go
func (o *Post) GetToPing() string
```
GetToPing returns the value of Post.ToPing

#### func (*Post) Reload

```go
func (o *Post) Reload() error
```
Reload A function to forcibly reload Post

#### func (*Post) Save

```go
func (o *Post) Save() error
```
Save is a dynamic saver 'inherited' by all models

#### func (*Post) SetCommentCount

```go
func (o *Post) SetCommentCount(arg int64)
```
SetCommentCount sets and marks as dirty the value of Post.CommentCount

#### func (*Post) SetCommentStatus

```go
func (o *Post) SetCommentStatus(arg string)
```
SetCommentStatus sets and marks as dirty the value of Post.CommentStatus

#### func (*Post) SetGuid

```go
func (o *Post) SetGuid(arg string)
```
SetGuid sets and marks as dirty the value of Post.Guid

#### func (*Post) SetID

```go
func (o *Post) SetID(arg int64)
```
SetID sets and marks as dirty the value of Post.ID

#### func (*Post) SetMenuOrder

```go
func (o *Post) SetMenuOrder(arg int)
```
SetMenuOrder sets and marks as dirty the value of Post.MenuOrder

#### func (*Post) SetPingStatus

```go
func (o *Post) SetPingStatus(arg string)
```
SetPingStatus sets and marks as dirty the value of Post.PingStatus

#### func (*Post) SetPinged

```go
func (o *Post) SetPinged(arg string)
```
SetPinged sets and marks as dirty the value of Post.Pinged

#### func (*Post) SetPostAuthor

```go
func (o *Post) SetPostAuthor(arg int64)
```
SetPostAuthor sets and marks as dirty the value of Post.PostAuthor

#### func (*Post) SetPostContent

```go
func (o *Post) SetPostContent(arg string)
```
SetPostContent sets and marks as dirty the value of Post.PostContent

#### func (*Post) SetPostContentFiltered

```go
func (o *Post) SetPostContentFiltered(arg string)
```
SetPostContentFiltered sets and marks as dirty the value of
Post.PostContentFiltered

#### func (*Post) SetPostDate

```go
func (o *Post) SetPostDate(arg *DateTime)
```
SetPostDate sets and marks as dirty the value of Post.PostDate

#### func (*Post) SetPostDateGmt

```go
func (o *Post) SetPostDateGmt(arg *DateTime)
```
SetPostDateGmt sets and marks as dirty the value of Post.PostDateGmt

#### func (*Post) SetPostExcerpt

```go
func (o *Post) SetPostExcerpt(arg string)
```
SetPostExcerpt sets and marks as dirty the value of Post.PostExcerpt

#### func (*Post) SetPostMimeType

```go
func (o *Post) SetPostMimeType(arg string)
```
SetPostMimeType sets and marks as dirty the value of Post.PostMimeType

#### func (*Post) SetPostModified

```go
func (o *Post) SetPostModified(arg *DateTime)
```
SetPostModified sets and marks as dirty the value of Post.PostModified

#### func (*Post) SetPostModifiedGmt

```go
func (o *Post) SetPostModifiedGmt(arg *DateTime)
```
SetPostModifiedGmt sets and marks as dirty the value of Post.PostModifiedGmt

#### func (*Post) SetPostName

```go
func (o *Post) SetPostName(arg string)
```
SetPostName sets and marks as dirty the value of Post.PostName

#### func (*Post) SetPostParent

```go
func (o *Post) SetPostParent(arg int64)
```
SetPostParent sets and marks as dirty the value of Post.PostParent

#### func (*Post) SetPostPassword

```go
func (o *Post) SetPostPassword(arg string)
```
SetPostPassword sets and marks as dirty the value of Post.PostPassword

#### func (*Post) SetPostStatus

```go
func (o *Post) SetPostStatus(arg string)
```
SetPostStatus sets and marks as dirty the value of Post.PostStatus

#### func (*Post) SetPostTitle

```go
func (o *Post) SetPostTitle(arg string)
```
SetPostTitle sets and marks as dirty the value of Post.PostTitle

#### func (*Post) SetPostType

```go
func (o *Post) SetPostType(arg string)
```
SetPostType sets and marks as dirty the value of Post.PostType

#### func (*Post) SetToPing

```go
func (o *Post) SetToPing(arg string)
```
SetToPing sets and marks as dirty the value of Post.ToPing

#### func (*Post) Update

```go
func (o *Post) Update() error
```
Update is a dynamic updater, it considers whether or not a field is 'dirty' and
needs to be updated. Will only work if you use the Getters and Setters

#### func (*Post) UpdateCommentCount

```go
func (o *Post) UpdateCommentCount(_updCommentCount int64) (int64, error)
```
UpdateCommentCount an immediate DB Query to update a single column, in this case
comment_count

#### func (*Post) UpdateCommentStatus

```go
func (o *Post) UpdateCommentStatus(_updCommentStatus string) (int64, error)
```
UpdateCommentStatus an immediate DB Query to update a single column, in this
case comment_status

#### func (*Post) UpdateGuid

```go
func (o *Post) UpdateGuid(_updGuid string) (int64, error)
```
UpdateGuid an immediate DB Query to update a single column, in this case guid

#### func (*Post) UpdateMenuOrder

```go
func (o *Post) UpdateMenuOrder(_updMenuOrder int) (int64, error)
```
UpdateMenuOrder an immediate DB Query to update a single column, in this case
menu_order

#### func (*Post) UpdatePingStatus

```go
func (o *Post) UpdatePingStatus(_updPingStatus string) (int64, error)
```
UpdatePingStatus an immediate DB Query to update a single column, in this case
ping_status

#### func (*Post) UpdatePinged

```go
func (o *Post) UpdatePinged(_updPinged string) (int64, error)
```
UpdatePinged an immediate DB Query to update a single column, in this case
pinged

#### func (*Post) UpdatePostAuthor

```go
func (o *Post) UpdatePostAuthor(_updPostAuthor int64) (int64, error)
```
UpdatePostAuthor an immediate DB Query to update a single column, in this case
post_author

#### func (*Post) UpdatePostContent

```go
func (o *Post) UpdatePostContent(_updPostContent string) (int64, error)
```
UpdatePostContent an immediate DB Query to update a single column, in this case
post_content

#### func (*Post) UpdatePostContentFiltered

```go
func (o *Post) UpdatePostContentFiltered(_updPostContentFiltered string) (int64, error)
```
UpdatePostContentFiltered an immediate DB Query to update a single column, in
this case post_content_filtered

#### func (*Post) UpdatePostDate

```go
func (o *Post) UpdatePostDate(_updPostDate *DateTime) (int64, error)
```
UpdatePostDate an immediate DB Query to update a single column, in this case
post_date

#### func (*Post) UpdatePostDateGmt

```go
func (o *Post) UpdatePostDateGmt(_updPostDateGmt *DateTime) (int64, error)
```
UpdatePostDateGmt an immediate DB Query to update a single column, in this case
post_date_gmt

#### func (*Post) UpdatePostExcerpt

```go
func (o *Post) UpdatePostExcerpt(_updPostExcerpt string) (int64, error)
```
UpdatePostExcerpt an immediate DB Query to update a single column, in this case
post_excerpt

#### func (*Post) UpdatePostMimeType

```go
func (o *Post) UpdatePostMimeType(_updPostMimeType string) (int64, error)
```
UpdatePostMimeType an immediate DB Query to update a single column, in this case
post_mime_type

#### func (*Post) UpdatePostModified

```go
func (o *Post) UpdatePostModified(_updPostModified *DateTime) (int64, error)
```
UpdatePostModified an immediate DB Query to update a single column, in this case
post_modified

#### func (*Post) UpdatePostModifiedGmt

```go
func (o *Post) UpdatePostModifiedGmt(_updPostModifiedGmt *DateTime) (int64, error)
```
UpdatePostModifiedGmt an immediate DB Query to update a single column, in this
case post_modified_gmt

#### func (*Post) UpdatePostName

```go
func (o *Post) UpdatePostName(_updPostName string) (int64, error)
```
UpdatePostName an immediate DB Query to update a single column, in this case
post_name

#### func (*Post) UpdatePostParent

```go
func (o *Post) UpdatePostParent(_updPostParent int64) (int64, error)
```
UpdatePostParent an immediate DB Query to update a single column, in this case
post_parent

#### func (*Post) UpdatePostPassword

```go
func (o *Post) UpdatePostPassword(_updPostPassword string) (int64, error)
```
UpdatePostPassword an immediate DB Query to update a single column, in this case
post_password

#### func (*Post) UpdatePostStatus

```go
func (o *Post) UpdatePostStatus(_updPostStatus string) (int64, error)
```
UpdatePostStatus an immediate DB Query to update a single column, in this case
post_status

#### func (*Post) UpdatePostTitle

```go
func (o *Post) UpdatePostTitle(_updPostTitle string) (int64, error)
```
UpdatePostTitle an immediate DB Query to update a single column, in this case
post_title

#### func (*Post) UpdatePostType

```go
func (o *Post) UpdatePostType(_updPostType string) (int64, error)
```
UpdatePostType an immediate DB Query to update a single column, in this case
post_type

#### func (*Post) UpdateToPing

```go
func (o *Post) UpdateToPing(_updToPing string) (int64, error)
```
UpdateToPing an immediate DB Query to update a single column, in this case
to_ping

#### type PostMeta

```go
type PostMeta struct {
	MetaId    int64
	PostId    int64
	MetaKey   string
	MetaValue string
	// Dirty markers for smart updates
	IsMetaIdDirty    bool
	IsPostIdDirty    bool
	IsMetaKeyDirty   bool
	IsMetaValueDirty bool
}
```

PostMeta is a Object Relational Mapping to the database table that represents
it. In this case it is postmeta. The table name will be Sprintf'd to include the
prefix you define in your YAML configuration for the Adapter.

#### func  NewPostMeta

```go
func NewPostMeta(a Adapter) *PostMeta
```
NewPostMeta binds an Adapter to a new instance of PostMeta and sets up the
_table and primary keys

#### func (*PostMeta) Create

```go
func (o *PostMeta) Create() error
```
Create inserts the model. Calling Save will call this function automatically for
new models

#### func (*PostMeta) Find

```go
func (o *PostMeta) Find(_findByMetaId int64) (bool, error)
```
Find dynamic finder for meta_id -> bool,error Generic and programatically
generator finder for PostMeta

Note that Fine returns a bool if found, not err, in the case of a return of
true, the instance data will be filled out. a call to find ALWAYS overwrites the
model you call Find on i.e. receiver is a pointer.

```go

    m := NewPostMeta(a)
    found,err := m.Find(23)
    .. handle err
    if found == false {
        // handle found
    }
    ... do what you want with m here

```

#### func (*PostMeta) FindByMetaKey

```go
func (o *PostMeta) FindByMetaKey(_findByMetaKey string) ([]*PostMeta, error)
```
FindByMetaKey dynamic finder for meta_key -> []*PostMeta,error Generic and
programatically generator finder for PostMeta

```go

    m := NewPostMeta(a)
    results,err := m.FindByMetaKey(...)
    // handle err
    for i,r := results {
      // now r is an instance of PostMeta
    }

```

#### func (*PostMeta) FindByMetaValue

```go
func (o *PostMeta) FindByMetaValue(_findByMetaValue string) ([]*PostMeta, error)
```
FindByMetaValue dynamic finder for meta_value -> []*PostMeta,error Generic and
programatically generator finder for PostMeta

```go

    m := NewPostMeta(a)
    results,err := m.FindByMetaValue(...)
    // handle err
    for i,r := results {
      // now r is an instance of PostMeta
    }

```

#### func (*PostMeta) FindByPostId

```go
func (o *PostMeta) FindByPostId(_findByPostId int64) ([]*PostMeta, error)
```
FindByPostId dynamic finder for post_id -> []*PostMeta,error Generic and
programatically generator finder for PostMeta

```go

    m := NewPostMeta(a)
    results,err := m.FindByPostId(...)
    // handle err
    for i,r := results {
      // now r is an instance of PostMeta
    }

```

#### func (*PostMeta) FromDBValueMap

```go
func (o *PostMeta) FromDBValueMap(m map[string]DBValue) error
```
FromDBValueMap Converts a DBValueMap returned from Adapter.Query to a PostMeta

#### func (*PostMeta) FromPostMeta

```go
func (o *PostMeta) FromPostMeta(m *PostMeta)
```
FromPostMeta A kind of Clone function for PostMeta

#### func (*PostMeta) GetMetaId

```go
func (o *PostMeta) GetMetaId() int64
```
GetMetaId returns the value of PostMeta.MetaId

#### func (*PostMeta) GetMetaKey

```go
func (o *PostMeta) GetMetaKey() string
```
GetMetaKey returns the value of PostMeta.MetaKey

#### func (*PostMeta) GetMetaValue

```go
func (o *PostMeta) GetMetaValue() string
```
GetMetaValue returns the value of PostMeta.MetaValue

#### func (*PostMeta) GetPostId

```go
func (o *PostMeta) GetPostId() int64
```
GetPostId returns the value of PostMeta.PostId

#### func (*PostMeta) GetPrimaryKeyName

```go
func (o *PostMeta) GetPrimaryKeyName() string
```
GetPrimaryKeyName returns the DB field name

#### func (*PostMeta) GetPrimaryKeyValue

```go
func (o *PostMeta) GetPrimaryKeyValue() int64
```
GetPrimaryKeyValue returns the value, usually int64 of the PrimaryKey

#### func (*PostMeta) Reload

```go
func (o *PostMeta) Reload() error
```
Reload A function to forcibly reload PostMeta

#### func (*PostMeta) Save

```go
func (o *PostMeta) Save() error
```
Save is a dynamic saver 'inherited' by all models

#### func (*PostMeta) SetMetaId

```go
func (o *PostMeta) SetMetaId(arg int64)
```
SetMetaId sets and marks as dirty the value of PostMeta.MetaId

#### func (*PostMeta) SetMetaKey

```go
func (o *PostMeta) SetMetaKey(arg string)
```
SetMetaKey sets and marks as dirty the value of PostMeta.MetaKey

#### func (*PostMeta) SetMetaValue

```go
func (o *PostMeta) SetMetaValue(arg string)
```
SetMetaValue sets and marks as dirty the value of PostMeta.MetaValue

#### func (*PostMeta) SetPostId

```go
func (o *PostMeta) SetPostId(arg int64)
```
SetPostId sets and marks as dirty the value of PostMeta.PostId

#### func (*PostMeta) Update

```go
func (o *PostMeta) Update() error
```
Update is a dynamic updater, it considers whether or not a field is 'dirty' and
needs to be updated. Will only work if you use the Getters and Setters

#### func (*PostMeta) UpdateMetaKey

```go
func (o *PostMeta) UpdateMetaKey(_updMetaKey string) (int64, error)
```
UpdateMetaKey an immediate DB Query to update a single column, in this case
meta_key

#### func (*PostMeta) UpdateMetaValue

```go
func (o *PostMeta) UpdateMetaValue(_updMetaValue string) (int64, error)
```
UpdateMetaValue an immediate DB Query to update a single column, in this case
meta_value

#### func (*PostMeta) UpdatePostId

```go
func (o *PostMeta) UpdatePostId(_updPostId int64) (int64, error)
```
UpdatePostId an immediate DB Query to update a single column, in this case
post_id

#### type SafeStringFilter

```go
type SafeStringFilter func(string) string
```

SafeStringFilter is the function that escapes possible SQL Injection code.

#### type Term

```go
type Term struct {
	TermId    int64
	Name      string
	Slug      string
	TermGroup int64
	// Dirty markers for smart updates
	IsTermIdDirty    bool
	IsNameDirty      bool
	IsSlugDirty      bool
	IsTermGroupDirty bool
}
```

Term is a Object Relational Mapping to the database table that represents it. In
this case it is terms. The table name will be Sprintf'd to include the prefix
you define in your YAML configuration for the Adapter.

#### func  NewTerm

```go
func NewTerm(a Adapter) *Term
```
NewTerm binds an Adapter to a new instance of Term and sets up the _table and
primary keys

#### func (*Term) Create

```go
func (o *Term) Create() error
```
Create inserts the model. Calling Save will call this function automatically for
new models

#### func (*Term) Find

```go
func (o *Term) Find(_findByTermId int64) (bool, error)
```
Find dynamic finder for term_id -> bool,error Generic and programatically
generator finder for Term

Note that Fine returns a bool if found, not err, in the case of a return of
true, the instance data will be filled out. a call to find ALWAYS overwrites the
model you call Find on i.e. receiver is a pointer.

```go

    m := NewTerm(a)
    found,err := m.Find(23)
    .. handle err
    if found == false {
        // handle found
    }
    ... do what you want with m here

```

#### func (*Term) FindByName

```go
func (o *Term) FindByName(_findByName string) ([]*Term, error)
```
FindByName dynamic finder for name -> []*Term,error Generic and programatically
generator finder for Term

```go

    m := NewTerm(a)
    results,err := m.FindByName(...)
    // handle err
    for i,r := results {
      // now r is an instance of Term
    }

```

#### func (*Term) FindBySlug

```go
func (o *Term) FindBySlug(_findBySlug string) ([]*Term, error)
```
FindBySlug dynamic finder for slug -> []*Term,error Generic and programatically
generator finder for Term

```go

    m := NewTerm(a)
    results,err := m.FindBySlug(...)
    // handle err
    for i,r := results {
      // now r is an instance of Term
    }

```

#### func (*Term) FindByTermGroup

```go
func (o *Term) FindByTermGroup(_findByTermGroup int64) ([]*Term, error)
```
FindByTermGroup dynamic finder for term_group -> []*Term,error Generic and
programatically generator finder for Term

```go

    m := NewTerm(a)
    results,err := m.FindByTermGroup(...)
    // handle err
    for i,r := results {
      // now r is an instance of Term
    }

```

#### func (*Term) FromDBValueMap

```go
func (o *Term) FromDBValueMap(m map[string]DBValue) error
```
FromDBValueMap Converts a DBValueMap returned from Adapter.Query to a Term

#### func (*Term) FromTerm

```go
func (o *Term) FromTerm(m *Term)
```
FromTerm A kind of Clone function for Term

#### func (*Term) GetName

```go
func (o *Term) GetName() string
```
GetName returns the value of Term.Name

#### func (*Term) GetPrimaryKeyName

```go
func (o *Term) GetPrimaryKeyName() string
```
GetPrimaryKeyName returns the DB field name

#### func (*Term) GetPrimaryKeyValue

```go
func (o *Term) GetPrimaryKeyValue() int64
```
GetPrimaryKeyValue returns the value, usually int64 of the PrimaryKey

#### func (*Term) GetSlug

```go
func (o *Term) GetSlug() string
```
GetSlug returns the value of Term.Slug

#### func (*Term) GetTermGroup

```go
func (o *Term) GetTermGroup() int64
```
GetTermGroup returns the value of Term.TermGroup

#### func (*Term) GetTermId

```go
func (o *Term) GetTermId() int64
```
GetTermId returns the value of Term.TermId

#### func (*Term) Reload

```go
func (o *Term) Reload() error
```
Reload A function to forcibly reload Term

#### func (*Term) Save

```go
func (o *Term) Save() error
```
Save is a dynamic saver 'inherited' by all models

#### func (*Term) SetName

```go
func (o *Term) SetName(arg string)
```
SetName sets and marks as dirty the value of Term.Name

#### func (*Term) SetSlug

```go
func (o *Term) SetSlug(arg string)
```
SetSlug sets and marks as dirty the value of Term.Slug

#### func (*Term) SetTermGroup

```go
func (o *Term) SetTermGroup(arg int64)
```
SetTermGroup sets and marks as dirty the value of Term.TermGroup

#### func (*Term) SetTermId

```go
func (o *Term) SetTermId(arg int64)
```
SetTermId sets and marks as dirty the value of Term.TermId

#### func (*Term) Update

```go
func (o *Term) Update() error
```
Update is a dynamic updater, it considers whether or not a field is 'dirty' and
needs to be updated. Will only work if you use the Getters and Setters

#### func (*Term) UpdateName

```go
func (o *Term) UpdateName(_updName string) (int64, error)
```
UpdateName an immediate DB Query to update a single column, in this case name

#### func (*Term) UpdateSlug

```go
func (o *Term) UpdateSlug(_updSlug string) (int64, error)
```
UpdateSlug an immediate DB Query to update a single column, in this case slug

#### func (*Term) UpdateTermGroup

```go
func (o *Term) UpdateTermGroup(_updTermGroup int64) (int64, error)
```
UpdateTermGroup an immediate DB Query to update a single column, in this case
term_group

#### type TermRelationship

```go
type TermRelationship struct {
	ObjectId       int64
	TermTaxonomyId int64
	TermOrder      int
	// Dirty markers for smart updates
	IsObjectIdDirty       bool
	IsTermTaxonomyIdDirty bool
	IsTermOrderDirty      bool
}
```

TermRelationship is a Object Relational Mapping to the database table that
represents it. In this case it is term_relationships. The table name will be
Sprintf'd to include the prefix you define in your YAML configuration for the
Adapter.

#### func  NewTermRelationship

```go
func NewTermRelationship(a Adapter) *TermRelationship
```
NewTermRelationship binds an Adapter to a new instance of TermRelationship and
sets up the _table and primary keys

#### func (*TermRelationship) Create

```go
func (o *TermRelationship) Create() error
```
Create inserts the model. Calling Save will call this function automatically for
new models

#### func (*TermRelationship) Find

```go
func (o *TermRelationship) Find(termId int64, objectId int64) (bool, error)
```
Find for the TermRelationship is a bit tricky, as it has no primary key as such,
but a composite key.

#### func (*TermRelationship) FindByObjectId

```go
func (o *TermRelationship) FindByObjectId(_findByObjectId int64) ([]*TermRelationship, error)
```
FindByObjectId dynamic finder for object_id -> []*TermRelationship,error Generic
and programatically generator finder for TermRelationship

```go

    m := NewTermRelationship(a)
    results,err := m.FindByObjectId(...)
    // handle err
    for i,r := results {
      // now r is an instance of TermRelationship
    }

```

#### func (*TermRelationship) FindByTermOrder

```go
func (o *TermRelationship) FindByTermOrder(_findByTermOrder int) ([]*TermRelationship, error)
```
FindByTermOrder dynamic finder for term_order -> []*TermRelationship,error
Generic and programatically generator finder for TermRelationship

```go

    m := NewTermRelationship(a)
    results,err := m.FindByTermOrder(...)
    // handle err
    for i,r := results {
      // now r is an instance of TermRelationship
    }

```

#### func (*TermRelationship) FromDBValueMap

```go
func (o *TermRelationship) FromDBValueMap(m map[string]DBValue) error
```
FromDBValueMap Converts a DBValueMap returned from Adapter.Query to a
TermRelationship

#### func (*TermRelationship) FromTermRelationship

```go
func (o *TermRelationship) FromTermRelationship(m *TermRelationship)
```
FromTermRelationship A kind of Clone function for TermRelationship

#### func (*TermRelationship) GetObjectId

```go
func (o *TermRelationship) GetObjectId() int64
```
GetObjectId returns the value of TermRelationship.ObjectId

#### func (*TermRelationship) GetPrimaryKeyName

```go
func (o *TermRelationship) GetPrimaryKeyName() string
```
GetPrimaryKeyName returns the DB field name

#### func (*TermRelationship) GetPrimaryKeyValue

```go
func (o *TermRelationship) GetPrimaryKeyValue() int64
```
GetPrimaryKeyValue returns the value, usually int64 of the PrimaryKey

#### func (*TermRelationship) GetTermOrder

```go
func (o *TermRelationship) GetTermOrder() int
```
GetTermOrder returns the value of TermRelationship.TermOrder

#### func (*TermRelationship) GetTermTaxonomyId

```go
func (o *TermRelationship) GetTermTaxonomyId() int64
```
GetTermTaxonomyId returns the value of TermRelationship.TermTaxonomyId

#### func (*TermRelationship) Reload

```go
func (o *TermRelationship) Reload() error
```
Reload A function to forcibly reload TermRelationship

#### func (*TermRelationship) Save

```go
func (o *TermRelationship) Save() error
```
Save is a dynamic saver 'inherited' by all models

#### func (*TermRelationship) SetObjectId

```go
func (o *TermRelationship) SetObjectId(arg int64)
```
SetObjectId sets and marks as dirty the value of TermRelationship.ObjectId

#### func (*TermRelationship) SetTermOrder

```go
func (o *TermRelationship) SetTermOrder(arg int)
```
SetTermOrder sets and marks as dirty the value of TermRelationship.TermOrder

#### func (*TermRelationship) SetTermTaxonomyId

```go
func (o *TermRelationship) SetTermTaxonomyId(arg int64)
```
SetTermTaxonomyId sets and marks as dirty the value of
TermRelationship.TermTaxonomyId

#### func (*TermRelationship) Update

```go
func (o *TermRelationship) Update() error
```
Update is a dynamic updater, it considers whether or not a field is 'dirty' and
needs to be updated. Will only work if you use the Getters and Setters

#### func (*TermRelationship) UpdateObjectId

```go
func (o *TermRelationship) UpdateObjectId(_updObjectId int64) (int64, error)
```
UpdateObjectId an immediate DB Query to update a single column, in this case
object_id

#### func (*TermRelationship) UpdateTermOrder

```go
func (o *TermRelationship) UpdateTermOrder(_updTermOrder int) (int64, error)
```
UpdateTermOrder an immediate DB Query to update a single column, in this case
term_order

#### type TermTaxonomy

```go
type TermTaxonomy struct {
	TermTaxonomyId int64
	TermId         int64
	Taxonomy       string
	Description    string
	Parent         int64
	Count          int64
	// Dirty markers for smart updates
	IsTermTaxonomyIdDirty bool
	IsTermIdDirty         bool
	IsTaxonomyDirty       bool
	IsDescriptionDirty    bool
	IsParentDirty         bool
	IsCountDirty          bool
}
```

TermTaxonomy is a Object Relational Mapping to the database table that
represents it. In this case it is term_taxonomy. The table name will be
Sprintf'd to include the prefix you define in your YAML configuration for the
Adapter.

#### func  NewTermTaxonomy

```go
func NewTermTaxonomy(a Adapter) *TermTaxonomy
```
NewTermTaxonomy binds an Adapter to a new instance of TermTaxonomy and sets up
the _table and primary keys

#### func (*TermTaxonomy) Create

```go
func (o *TermTaxonomy) Create() error
```
Create inserts the model. Calling Save will call this function automatically for
new models

#### func (*TermTaxonomy) Find

```go
func (o *TermTaxonomy) Find(_findByTermTaxonomyId int64) (bool, error)
```
Find dynamic finder for term_taxonomy_id -> bool,error Generic and
programatically generator finder for TermTaxonomy

Note that Fine returns a bool if found, not err, in the case of a return of
true, the instance data will be filled out. a call to find ALWAYS overwrites the
model you call Find on i.e. receiver is a pointer.

```go

    m := NewTermTaxonomy(a)
    found,err := m.Find(23)
    .. handle err
    if found == false {
        // handle found
    }
    ... do what you want with m here

```

#### func (*TermTaxonomy) FindByCount

```go
func (o *TermTaxonomy) FindByCount(_findByCount int64) ([]*TermTaxonomy, error)
```
FindByCount dynamic finder for count -> []*TermTaxonomy,error Generic and
programatically generator finder for TermTaxonomy

```go

    m := NewTermTaxonomy(a)
    results,err := m.FindByCount(...)
    // handle err
    for i,r := results {
      // now r is an instance of TermTaxonomy
    }

```

#### func (*TermTaxonomy) FindByDescription

```go
func (o *TermTaxonomy) FindByDescription(_findByDescription string) ([]*TermTaxonomy, error)
```
FindByDescription dynamic finder for description -> []*TermTaxonomy,error
Generic and programatically generator finder for TermTaxonomy

```go

    m := NewTermTaxonomy(a)
    results,err := m.FindByDescription(...)
    // handle err
    for i,r := results {
      // now r is an instance of TermTaxonomy
    }

```

#### func (*TermTaxonomy) FindByParent

```go
func (o *TermTaxonomy) FindByParent(_findByParent int64) ([]*TermTaxonomy, error)
```
FindByParent dynamic finder for parent -> []*TermTaxonomy,error Generic and
programatically generator finder for TermTaxonomy

```go

    m := NewTermTaxonomy(a)
    results,err := m.FindByParent(...)
    // handle err
    for i,r := results {
      // now r is an instance of TermTaxonomy
    }

```

#### func (*TermTaxonomy) FindByTaxonomy

```go
func (o *TermTaxonomy) FindByTaxonomy(_findByTaxonomy string) ([]*TermTaxonomy, error)
```
FindByTaxonomy dynamic finder for taxonomy -> []*TermTaxonomy,error Generic and
programatically generator finder for TermTaxonomy

```go

    m := NewTermTaxonomy(a)
    results,err := m.FindByTaxonomy(...)
    // handle err
    for i,r := results {
      // now r is an instance of TermTaxonomy
    }

```

#### func (*TermTaxonomy) FindByTermId

```go
func (o *TermTaxonomy) FindByTermId(_findByTermId int64) ([]*TermTaxonomy, error)
```
FindByTermId dynamic finder for term_id -> []*TermTaxonomy,error Generic and
programatically generator finder for TermTaxonomy

```go

    m := NewTermTaxonomy(a)
    results,err := m.FindByTermId(...)
    // handle err
    for i,r := results {
      // now r is an instance of TermTaxonomy
    }

```

#### func (*TermTaxonomy) FromDBValueMap

```go
func (o *TermTaxonomy) FromDBValueMap(m map[string]DBValue) error
```
FromDBValueMap Converts a DBValueMap returned from Adapter.Query to a
TermTaxonomy

#### func (*TermTaxonomy) FromTermTaxonomy

```go
func (o *TermTaxonomy) FromTermTaxonomy(m *TermTaxonomy)
```
FromTermTaxonomy A kind of Clone function for TermTaxonomy

#### func (*TermTaxonomy) GetCount

```go
func (o *TermTaxonomy) GetCount() int64
```
GetCount returns the value of TermTaxonomy.Count

#### func (*TermTaxonomy) GetDescription

```go
func (o *TermTaxonomy) GetDescription() string
```
GetDescription returns the value of TermTaxonomy.Description

#### func (*TermTaxonomy) GetParent

```go
func (o *TermTaxonomy) GetParent() int64
```
GetParent returns the value of TermTaxonomy.Parent

#### func (*TermTaxonomy) GetPrimaryKeyName

```go
func (o *TermTaxonomy) GetPrimaryKeyName() string
```
GetPrimaryKeyName returns the DB field name

#### func (*TermTaxonomy) GetPrimaryKeyValue

```go
func (o *TermTaxonomy) GetPrimaryKeyValue() int64
```
GetPrimaryKeyValue returns the value, usually int64 of the PrimaryKey

#### func (*TermTaxonomy) GetTaxonomy

```go
func (o *TermTaxonomy) GetTaxonomy() string
```
GetTaxonomy returns the value of TermTaxonomy.Taxonomy

#### func (*TermTaxonomy) GetTermId

```go
func (o *TermTaxonomy) GetTermId() int64
```
GetTermId returns the value of TermTaxonomy.TermId

#### func (*TermTaxonomy) GetTermTaxonomyId

```go
func (o *TermTaxonomy) GetTermTaxonomyId() int64
```
GetTermTaxonomyId returns the value of TermTaxonomy.TermTaxonomyId

#### func (*TermTaxonomy) Reload

```go
func (o *TermTaxonomy) Reload() error
```
Reload A function to forcibly reload TermTaxonomy

#### func (*TermTaxonomy) Save

```go
func (o *TermTaxonomy) Save() error
```
Save is a dynamic saver 'inherited' by all models

#### func (*TermTaxonomy) SetCount

```go
func (o *TermTaxonomy) SetCount(arg int64)
```
SetCount sets and marks as dirty the value of TermTaxonomy.Count

#### func (*TermTaxonomy) SetDescription

```go
func (o *TermTaxonomy) SetDescription(arg string)
```
SetDescription sets and marks as dirty the value of TermTaxonomy.Description

#### func (*TermTaxonomy) SetParent

```go
func (o *TermTaxonomy) SetParent(arg int64)
```
SetParent sets and marks as dirty the value of TermTaxonomy.Parent

#### func (*TermTaxonomy) SetTaxonomy

```go
func (o *TermTaxonomy) SetTaxonomy(arg string)
```
SetTaxonomy sets and marks as dirty the value of TermTaxonomy.Taxonomy

#### func (*TermTaxonomy) SetTermId

```go
func (o *TermTaxonomy) SetTermId(arg int64)
```
SetTermId sets and marks as dirty the value of TermTaxonomy.TermId

#### func (*TermTaxonomy) SetTermTaxonomyId

```go
func (o *TermTaxonomy) SetTermTaxonomyId(arg int64)
```
SetTermTaxonomyId sets and marks as dirty the value of
TermTaxonomy.TermTaxonomyId

#### func (*TermTaxonomy) Update

```go
func (o *TermTaxonomy) Update() error
```
Update is a dynamic updater, it considers whether or not a field is 'dirty' and
needs to be updated. Will only work if you use the Getters and Setters

#### func (*TermTaxonomy) UpdateCount

```go
func (o *TermTaxonomy) UpdateCount(_updCount int64) (int64, error)
```
UpdateCount an immediate DB Query to update a single column, in this case count

#### func (*TermTaxonomy) UpdateDescription

```go
func (o *TermTaxonomy) UpdateDescription(_updDescription string) (int64, error)
```
UpdateDescription an immediate DB Query to update a single column, in this case
description

#### func (*TermTaxonomy) UpdateParent

```go
func (o *TermTaxonomy) UpdateParent(_updParent int64) (int64, error)
```
UpdateParent an immediate DB Query to update a single column, in this case
parent

#### func (*TermTaxonomy) UpdateTaxonomy

```go
func (o *TermTaxonomy) UpdateTaxonomy(_updTaxonomy string) (int64, error)
```
UpdateTaxonomy an immediate DB Query to update a single column, in this case
taxonomy

#### func (*TermTaxonomy) UpdateTermId

```go
func (o *TermTaxonomy) UpdateTermId(_updTermId int64) (int64, error)
```
UpdateTermId an immediate DB Query to update a single column, in this case
term_id

#### type UserMeta

```go
type UserMeta struct {
	UMetaId   int64
	UserId    int64
	MetaKey   string
	MetaValue string
	// Dirty markers for smart updates
	IsUMetaIdDirty   bool
	IsUserIdDirty    bool
	IsMetaKeyDirty   bool
	IsMetaValueDirty bool
}
```

UserMeta is a Object Relational Mapping to the database table that represents
it. In this case it is usermeta. The table name will be Sprintf'd to include the
prefix you define in your YAML configuration for the Adapter.

#### func  NewUserMeta

```go
func NewUserMeta(a Adapter) *UserMeta
```
NewUserMeta binds an Adapter to a new instance of UserMeta and sets up the
_table and primary keys

#### func (*UserMeta) Create

```go
func (o *UserMeta) Create() error
```
Create inserts the model. Calling Save will call this function automatically for
new models

#### func (*UserMeta) Find

```go
func (o *UserMeta) Find(_findByUMetaId int64) (bool, error)
```
Find dynamic finder for umeta_id -> bool,error Generic and programatically
generator finder for UserMeta

Note that Fine returns a bool if found, not err, in the case of a return of
true, the instance data will be filled out. a call to find ALWAYS overwrites the
model you call Find on i.e. receiver is a pointer.

```go

    m := NewUserMeta(a)
    found,err := m.Find(23)
    .. handle err
    if found == false {
        // handle found
    }
    ... do what you want with m here

```

#### func (*UserMeta) FindByMetaKey

```go
func (o *UserMeta) FindByMetaKey(_findByMetaKey string) ([]*UserMeta, error)
```
FindByMetaKey dynamic finder for meta_key -> []*UserMeta,error Generic and
programatically generator finder for UserMeta

```go

    m := NewUserMeta(a)
    results,err := m.FindByMetaKey(...)
    // handle err
    for i,r := results {
      // now r is an instance of UserMeta
    }

```

#### func (*UserMeta) FindByMetaValue

```go
func (o *UserMeta) FindByMetaValue(_findByMetaValue string) ([]*UserMeta, error)
```
FindByMetaValue dynamic finder for meta_value -> []*UserMeta,error Generic and
programatically generator finder for UserMeta

```go

    m := NewUserMeta(a)
    results,err := m.FindByMetaValue(...)
    // handle err
    for i,r := results {
      // now r is an instance of UserMeta
    }

```

#### func (*UserMeta) FindByUserId

```go
func (o *UserMeta) FindByUserId(_findByUserId int64) ([]*UserMeta, error)
```
FindByUserId dynamic finder for user_id -> []*UserMeta,error Generic and
programatically generator finder for UserMeta

```go

    m := NewUserMeta(a)
    results,err := m.FindByUserId(...)
    // handle err
    for i,r := results {
      // now r is an instance of UserMeta
    }

```

#### func (*UserMeta) FromDBValueMap

```go
func (o *UserMeta) FromDBValueMap(m map[string]DBValue) error
```
FromDBValueMap Converts a DBValueMap returned from Adapter.Query to a UserMeta

#### func (*UserMeta) FromUserMeta

```go
func (o *UserMeta) FromUserMeta(m *UserMeta)
```
FromUserMeta A kind of Clone function for UserMeta

#### func (*UserMeta) GetMetaKey

```go
func (o *UserMeta) GetMetaKey() string
```
GetMetaKey returns the value of UserMeta.MetaKey

#### func (*UserMeta) GetMetaValue

```go
func (o *UserMeta) GetMetaValue() string
```
GetMetaValue returns the value of UserMeta.MetaValue

#### func (*UserMeta) GetPrimaryKeyName

```go
func (o *UserMeta) GetPrimaryKeyName() string
```
GetPrimaryKeyName returns the DB field name

#### func (*UserMeta) GetPrimaryKeyValue

```go
func (o *UserMeta) GetPrimaryKeyValue() int64
```
GetPrimaryKeyValue returns the value, usually int64 of the PrimaryKey

#### func (*UserMeta) GetUMetaId

```go
func (o *UserMeta) GetUMetaId() int64
```
GetUMetaId returns the value of UserMeta.UMetaId

#### func (*UserMeta) GetUserId

```go
func (o *UserMeta) GetUserId() int64
```
GetUserId returns the value of UserMeta.UserId

#### func (*UserMeta) Reload

```go
func (o *UserMeta) Reload() error
```
Reload A function to forcibly reload UserMeta

#### func (*UserMeta) Save

```go
func (o *UserMeta) Save() error
```
Save is a dynamic saver 'inherited' by all models

#### func (*UserMeta) SetMetaKey

```go
func (o *UserMeta) SetMetaKey(arg string)
```
SetMetaKey sets and marks as dirty the value of UserMeta.MetaKey

#### func (*UserMeta) SetMetaValue

```go
func (o *UserMeta) SetMetaValue(arg string)
```
SetMetaValue sets and marks as dirty the value of UserMeta.MetaValue

#### func (*UserMeta) SetUMetaId

```go
func (o *UserMeta) SetUMetaId(arg int64)
```
SetUMetaId sets and marks as dirty the value of UserMeta.UMetaId

#### func (*UserMeta) SetUserId

```go
func (o *UserMeta) SetUserId(arg int64)
```
SetUserId sets and marks as dirty the value of UserMeta.UserId

#### func (*UserMeta) Update

```go
func (o *UserMeta) Update() error
```
Update is a dynamic updater, it considers whether or not a field is 'dirty' and
needs to be updated. Will only work if you use the Getters and Setters

#### func (*UserMeta) UpdateMetaKey

```go
func (o *UserMeta) UpdateMetaKey(_updMetaKey string) (int64, error)
```
UpdateMetaKey an immediate DB Query to update a single column, in this case
meta_key

#### func (*UserMeta) UpdateMetaValue

```go
func (o *UserMeta) UpdateMetaValue(_updMetaValue string) (int64, error)
```
UpdateMetaValue an immediate DB Query to update a single column, in this case
meta_value

#### func (*UserMeta) UpdateUserId

```go
func (o *UserMeta) UpdateUserId(_updUserId int64) (int64, error)
```
UpdateUserId an immediate DB Query to update a single column, in this case
user_id
