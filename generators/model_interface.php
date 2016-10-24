<?php
$txt = "
// Model is the generic interface that can be 
// passed to an Adapter to be executed.
type Model interface {
    // Default is *, can use any Model field
    Select(string) Model
    // a string like: ID = 2 AND post_title LIKE '%The Rain in%'
    Where(string) Model
    // The full, comma separated SET line, like SET post_title = 'xxx',post_status = 'yyy'
    // e.g. m.Set(`post_title`,'xxx').Set('post_status','draft').Go()
    Set(string,string) Model
    // Set the LIMIT
    Limit(string) Model
    // Set the order by
    Order(string)
    // Set the cols for an insert
    Columns(string)
    // Set VALUES for the exec
    Values(string)
    // Return the Set value
    GetSet() string
    // Return the Select value
    GetSelect() string
    // Return the Where value
    GetWhere() string
    // Return the table name
    GetTable() string
    // Return the Columns
    GetColumns() string
    // Return the Values
    GetValues()
    // Return the LIMIT
    GetLimit() string
    // Return Order
    GetOrder() string
    // Query the results
    FindAll() []Model,err
    // Find
    FindOne() Model,err
    // Execute the statement
    Go() err
}
";

puts($txt);