# MongoUtils

Mongodb helper functions, document and pipeline builder.

## Helpers

### MongoOperationCtx

Create context for mongo db operations for 10 sec.

```go
MongoOperationCtx() (context.Context, context.CancelFunc)
```

### ParseObjectID

Parse object id from string.

```go
ParseObjectID(id string) *primitive.ObjectID
```

### IsValidObjectId

Check if object id is valid and not zero.

```go
IsValidObjectId(id *primitive.ObjectID) bool
```

### FindOption

Generate find option with sorts params.

```go
FindOption(sort any, skip int64, limit int64) *options.FindOptions
```

### AggregateOption

Generate aggregation options.

```go
AggregateOption() *options.AggregateOptions
```

### TxOption

Generate transaction option with majority write and snapshot read.

```go
AggregateOption() *options.AggregateOptions
```

### Array

Generate `primitive.A` from parameters.

```go
Array(args ...any) primitive.A
```

### Map

Generate `primitive.M` from parameters. Parameters count must be even.

```go
// Signature:
Map(args ...any) primitive.M

// Example:
mongoutils.Map("name", "John", "age", 23) // { "name": "John", "age": 23 }
```

### Maps

Generate `[]primitive.M` from parameters. Parameters count must be even.

```go
// Signature:
Maps(args ...any) []primitive.M

// Example:
mongoutils.Maps("name", "John", "age", 23) // [{ "name": "John" }, { "age": 23 }]
```

### Doc

Generate primitive.D from parameters. Parameters count must be even.

```go
// Signature:
Doc(args ...any) primitive.D

// Example:
mongoutils.Doc("name", "John", "age", 23) // { "name": "John", "age": 23 }
```

### Regex

Generate mongo `Regex` doc.

```go
// Signature:
Regex(pattern string, opt string) primitive.Regex

// Example:
mongoutils.Regex("John.*", "i") // { pattern: "John.*", options: "i" }
```

### RegexFor

Generate map with regex parameter.

```go
// Signature
RegexFor(k string, pattern string, opt string) primitive.M

// Example:
mongoutils.RegexFor("name", "John.*", "i") // { "name": { pattern: "John.*", options: "i" } }
```

### In

Generate $in map `{k: {$in: v}}`.

```go
In(k string, v ...any) primitive.M
```

### Set

Generate simple set map `{$set: v}`.

```go
Set(v any) primitive.M
```

### SetNested

Generate nested set map `{$set: {k: v}}`.

```go
SetNested(k string, v any) primitive.M
```

### Match

Generate nested set map `{$match: v}`.

```go
Match(v any) primitive.M
```

## Model Interface

Base interface for mongodb model.

**Note:** You can inherit `EmptyModel` in your struct. `EmptyModel` implements all model methods.

**Note:** You can inherit `BaseModel` in your struct. `BaseModel` contains `_id`, `created_at`, `updated_at` fields and `SetID`, `NewId`, `PrepareInsert` and `PrepareUpdate`. you can override or implement other model method.

**Note**: Set `BaseModel` bson tag to `inline` for insert timestamps in document root.

```go
// Usage:
import "github.com/bopher/mongoutils"
type Person struct{
 mongoutils.BaseModel  `bson:",inline"`
}

// override methods
func (me Person) IsDeletable() bool{
    return true
}
```

### Available methods

**Note:** All method must implement with pointer receiver!

**Note**: Context can passed to model method on call for mongodb transaction mode!

```go
// TypeName get type string
TypeName() string
// Collection get model collection
Collection(db *mongo.Database) *mongo.Collection
// Indexes create model indexes
Index(db *mongo.Database) error
// Seed run model seed
Seed(db *mongo.Database) error
// Pipeline get model pipeline
Pipeline() MongoPipeline
// NewId generate new id for model (Pointer)
NewId()
// SetID set model id (Pointer)
SetID(id primitive.ObjectID)
// GetID get model id
GetID() primitive.ObjectID
// IsEditable check if document is editable
// by default returns true on BaseModel
IsEditable() bool
// IsDeletable check if document is deletable
// by default returns false on BaseModel
IsDeletable() bool
// Cleanup document before save
// e.g set document field nil for ignore saving
Cleanup()
// PrepareInsert fill created_at before save (Pointer)
PrepareInsert()
// PrepareUpdate fill updated_at before save (Pointer)
// in ghost mode updated_at field not changed
PrepareUpdate(ghost bool)
// OnInsert function to call before insert (Pointer)
OnInsert(ctx context.Context, opt ...MongoOption)
// OnUpdate function to call before update (Pointer)
OnUpdate(ctx context.Context, opt ...MongoOption)
// OnDelete function to call before delete (Pointer)
OnDelete(ctx context.Context, opt ...MongoOption)
// OnInserted function to call after insert
OnInserted(ctx context.Context, opt ...MongoOption)
// OnUpdated function to call after update
OnUpdated(old any, ctx context.Context, opt ...MongoOption)
// OnDeleted function to call after delete
OnDeleted(ctx context.Context, opt ...MongoOption)
```

### Required Methods

Two `PrepareInsert` and `PrepareUpdate` must called before save model to database.

**Note**: if `true` passed to `PrepareUpdate` method, `updated_at` method not updated.

## Doc Builder

Document builder is a helper type for creating mongo document (`primitive.D`) with _chained_ methods.

```go
import "github.com/bopher/mongoutils"
doc := mongoutils.NewDoc()
doc.
    Add("name", "John").
    Add("nick", "John2").
    Array("skills", "javascript", "go", "rust", "mongo")
fmt.Println(doc.Build())
// -> {
//   "name": "John",
//   "nick": "John2",
//   "skills": ["javascript","go","rust","mongo"]
// }
```

### Doc Methods

#### Add

Add new element.

```go
// Signature:
Add(k string, v any) MongoDoc

// Example:
doc.Add("name", "Kim")
```

#### Doc

Add new element with nested doc value.

```go
// Signature:
Doc(k string, cb func(d MongoDoc) MongoDoc) MongoDoc

// Example:
doc.Doc("age", func(d mongoutils.MongoDoc) mongoutils.MongoDoc {
    d.Add("$gt", 20)
    d.Add("$lte", 30)
    return d
}) // -> { "age": { "$gt": 20, "$lte": 30 } }
```

#### Array

Add new element with array value.

```go
// Signature:
Array(k string, v ...any) MongoDoc

// Example:
doc.Array("skills", "javascript", "golang") // -> { "skills": ["javascript", "golang"] }
```

#### DocArray

Add new array element with doc child.

```go
// Signature:
DocArray(k string, cb func(d MongoDoc) MongoDoc) MongoDoc

// Example:
doc.DocArray("$match", func(d mongoutils.MongoDoc) mongoutils.MongoDoc {
    return d.Add("name", "John")
            Add("Family", "Doe")
}) // -> { "$match": [{"name": "John"}, {"Family": "Doe"}] }
```

#### Nested

Add new nested element.

```go
// Signature:
Nested(root string, k string, v any) MongoDoc

// Example:
doc.Nested("$set", "name", "Jack") // { "$set": { "name": "Jack" } }
```

#### NestedDoc

Add new nested element with doc value.

```go
// Signature:
NestedDoc(root string, k string, cb func(d MongoDoc) MongoDoc) MongoDoc

// Example:
doc.NestedDoc("$set", "address", func(d mongoutils.MongoDoc) mongoutils.MongoDoc {
    d.
        Add("city", "London").
        Add("street", "12th")
    return d
}) // -> { "$set": { "address": { "city": "London", "street": "12th" } } }
```

#### NestedArray

Add new nested element with array value.

```go
// Signature:
NestedArray(root string, k string, v ...any) MongoDoc

// Example:
doc.NestedArray("skill", "$in", "mongo", "golang") // -> { "skill": { "$in": ["mongo", "golang"] } }
```

#### NestedDocArray

Add new nested array element with doc

```go
// Signature:
NestedDocArray(root string, k string, cb func(d MongoDoc) MongoDoc) MongoDoc

// Example:
doc.NestedDocArray("name", "$match", func(d mongoutils.MongoDoc) mongoutils.MongoDoc {
    return d.Add("first", "John")
            Add("last", "Doe")
}) // -> { "name" : {"$match": [{"name": "John"}, {"last": "Doe"}] } }
```

#### Regex

Add new element with regex value.

```go
// Signature:
Regex(k string, pattern string, opt string) MongoDoc

// Example:
doc.Regex("full_name", "John.*", "i") // -> { "full_name": { pattern: "John.*", options: "i" } }
```

#### Map

Creates a map from the elements of the Doc.

```go
Map() primitive.M
```

#### Build

Generate mongo doc.

```go
Build() primitive.D
```

## Pipeline Builder

Pipeline builder is a helper type for creating mongo pipeline (`[]primitive.D`) with _chained_ methods.

```go
import "github.com/bopher/mongoutils"
pipe := mongoutils.NewPipe()
pipe.
    Add(func(d mongoutils.MongoDoc) mongoutils.MongoDoc{
        d.Nested("$match", "name", "John")
        return d
    }).
    Group(func(d mongoutils.MongoDoc) mongoutils.MongoDoc{
        d.
            Add("_id", "$_id").
            Nested("name", "$first", "$name")
            Nested("total", "$sum", "$invoice")
        return d
    })
fmt.Println(pipe.Build())
// -> [
//   { "$match": { "name": "John"} },
//   { "$group": {
//       "_id": "$_id"
//       "name": { "$first": "$name" },
//       "total": { "$sum": "$invoice" }
//   }}
// ]
```

### Pipeline Methods

#### Add

Add new Doc.

```go
// Signature:
Add(cb func(d MongoDoc) MongoDoc) MongoPipeline

// Example:
pipe.Add(func(d mongoutils.MongoDoc) mongoutils.MongoDoc{
    d.Nested("$match", "name", "John")
    return d
}) // -> [ {"$match": { "name": "John"}} ]
```

### Match

Add $match stage. skip nil input

```go
// Signature:
Match(filters any) MongoPipeline

// Example:
pipe.Match(v)
```

### In

Add $in stage.

```go
// Signature:
In(key string, v any) MongoPipeline

// Example:
pipe.In("status", statuses)
```

### Limit

Add $limit stage (ignore negative and zero value).

```go
// Signature:
Limit(limit int64) MongoPipeline

// Example:
pipe.Limit(100)
```

### Skip

Add $skip stage (ignore negative and zero value).

```go
// Signature:
Skip(skip int64) MongoPipeline

// Example:
pipe.Skip(25)
```

### Sort

Add $sort stage (ignore nil value).

```go
// Signature:
Sort(sorts any) MongoPipeline

// Example:
pipe.Sort(primitive.M{"username": 1})
```

#### Unwind

Add $unwind stage.

```go
// Signature:
Unwind(path string, prevNullAndEmpty bool) MongoPipeline

// Example:
pipe.Unwind("services", true)
// -> [
//     {"$unwind": {
//         "path": "services",
//         "preserveNullAndEmptyArrays": true,
//     }}
// ]
```

#### Lookup

Add $lookup stage.

```go
// Signature:
Lookup(from string, local string, foreign string, as string) MongoPipeline

// Example:
pipe.Lookup("users", "user_id", "_id", "user")
// -> [
//     {"$lookup": {
//         "from": "users",
//         "localField": "user_id",
//         "foreignField": "_id",
//         "as": "user"
//     }}
// ]
```

#### Unwrap

Get first item of array and insert to doc using $addFields stage. When using lookup result returns as array, use me helper to unwrap lookup result as field.

```go
// Signature:
Unwrap(field string, as string) MongoPipeline

// Example:
pipe.
    Lookup("users", "user_id", "_id", "__user").
    Unwrap("$__user", "user")
// -> [
//     { "$lookup": {
//         "from": "users",
//         "localField": "user_id",
//         "foreignField": "_id",
//         "as": "user"
//     }},
//     { "$addFields": { "user" : { "$first": "$__user" } } }
// ]
```

### LoadRelation

Load related document using `$lookup` and `$addField` (Lookup and Unwrap method mix).

```go
// Signature:
LoadRelation(from string, local string, foreign string, as string) MongoPipeline

// Example:
pipe.LoadRelation("users", "user_id", "_id", "user")
```

#### Group

Add $group stage.

```go
// Signature:
Group(cb func(d MongoDoc) MongoDoc) MongoPipeline

// Example:
pipe.
    Group(func(d mongoutils.MongoDoc) mongoutils.MongoDoc{
        d.
            Add("_id", "$_id").
            Nested("name", "$first", "$name").
            Nested("total", "$sum", "$invoice")
        return d
    })
// -> [
//   { "$group": {
//       "_id": "$_id"
//       "name": { "$first": "$name" },
//       "total": { "$sum": "$invoice" }
//   }}
// ]
```

#### ReplaceRoot

Add $replaceRoot stage.

```go
// Signature:
ReplaceRoot(v any) MongoPipeline

// Example:
pipe.ReplaceRoot("$my_root")
// ->  [{ "$replaceRoot": {"newRoot": "$my_root" } }]
```

#### MergeRoot

Add $replaceRoot stage with $mergeObjects operator.

```go
// Signature:
MergeRoot(fields ...any) MongoPipeline

// Example:
pipe.MergeRoot("$my_root", "$$ROOT")
// -> [
//     {
//         "$replaceRoot": {
//             "newRoot": { "mergeObjects": ["$my_root", "$$ROOT"] }
//         }
//     }
// ]
```

#### UnProject

Generate $project stage to remove fields from result.

```go
// Signature:
UnProject(fields ...string) MongoPipeline

// Example:
pipe.UnProject("my_root", "__user")
// -> [
//     { "$project": { "my_root": 0, "__user": 0 } }
// ]
```

#### Project

Generate $project stage. skip nil input.

```go
// Signature:
Project(projects any) MongoPipeline

// Example:
pipe.Project(nil) // skiped
pipe.Project(primitve.M{"password": 0}) // remove password from result
```

#### Build

Generate mongo pipeline.

```go
Build() mongo.Pipeline
```

## MetaCounter

meta counter builder for mongo docs.

```go
import "github.com/bopher/mongoutils"
mCounter := mongoutils.NewMetaCounter()
mCounter.Add("services", "relations", id1, 2)
mCounter.Add("services", "relations", id1, 1)
mCounter.Add("services", "total", id2, 1)
mCounter.Add("services", "total", id2, 1)
mCounter.Add("services", "relations", id3, 3)
mCounter.Add("services", "relations", nil, 3) // ignored
mCounter.Add("services", "relations", id3, -3) // ignored because of 0
mCounter.Add("customers", "rel", id1, 4)
mCounter.Add("customers", "rel", id2, 3)
mCounter.Add("customers", "rel", id2, 1)
mCounter.Add("customers", "rel", id3, 4)
fmt.Println(mCounter.Build())
// ->
// [
//   {
//     "Col": "services",
//     "Ids": ["62763152a01b7d275ef58e00"],
//     "Values": {
//       "relations": 3
//     }
//   },
//   {
//     "Col": "services",
//     "Ids": ["62763152a01b7d275ef58e01"],
//     "Values": {
//       "total": 2
//     }
//   },
//   {
//     "Col": "customers",
//     "Ids": [
//       "62763152a01b7d275ef58e00",
//       "62763152a01b7d275ef58e01",
//       "62763152a01b7d275ef58e02"
//     ],
//     "Values": {
//       "rel": 4
//     }
//   }
// ]
```

## MetaSetter

meta setter builder for mongo docs.

```go
import "github.com/bopher/mongoutils"
setter := mongoutils.NewMetaSetter()
setter.Add("test", "activity", id1, date)
setter.Add("test", "activity", nil, date) // ignored
setter.Add("test", "activity", id2, date)
setter.Add("test", "activity", id3, date) // override next line
setter.Add("test", "activity", id3, nil) // nil used
fmt.Println(setter.Build())
// ->
// [
//   {
//     "Col": "test",
//     "Ids": ["6276509942d11385d52b7ae2", "6276509942d11385d52b7ae3"],
//     "Values": {
//       "activity": "2022-05-07 10:57:29.6228877 +0000 UTC"
//     }
//   },
//   {
//     "Col": "test",
//     "Ids": ["6276509942d11385d52b7ae4"],
//     "Values": {
//       "activity": nil
//     }
//   },
// ]
```

## Repository

Methods for work with data based on `mongoutils.Model` implementation!

You can define multiple Pipeline methods for your model and use them to fetch data by Pipeline option and params. If no pipeline option passed functions used `Pipeline()` method by default!

### Find

Find find records.

```go
// Signature
func Find[T any](
    filter any,
    sorts any,
    skip int64,
    limit int64,
    opts ...MongoOption,
) ([]T, error)

// Example
import "github.com/bopher/mongoutils"
type User struct{
    mongoutils.BaseModel `bson:",inline"`
    Name string `bson:"name" json:"name"`
}

func (*User) UserWithAccountPipe() mongoutils.MongoPipeline{
    // ...
}

users, err := mongoutils.Find[User](
    primitive.M{"name": "John"},
    primitive.M{"created_at": -1}, 0, 10,
    MongoOption{
        Debug: true,
        Pipeline: "UserWithAccountPipe",
        Params: []any{"1 June 1991", 12, 3},
    })
```

### FindOne

Find one record.

```go
// Signature
func FindOne[T any](
    filter any,
    sorts any,
    opts ...MongoOption,
) (*T, error)
```

### Insert

Insert new record.

```go
// Signature
func Insert[T any](
    v *T,
    opts ...MongoOption,
) (*mongo.InsertOneResult, error)
```

### Update

Update one record.

```go
// Signature
func Update[T any](
    v *T,
    silent bool,
    opts ...MongoOption,
) (*mongo.UpdateResult, error)
```

### Delete

Delete one record.

```go
// Signature
func Delete[T any](
    v *T,
    opts ...MongoOption,
) (*mongo.DeleteResult, error)
```

### Count

Get records count.

```go
// Signature
func Count[T any](
    filter any,
    opts ...MongoOption,
) (int64, error)
```

### BatchUpdate

Update multiple records.

```go
// Signature
func BatchUpdate[T any](
    condition any,
    updates any,
    opts ...MongoOption,
) (*mongo.UpdateResult, error)
```

### Patch

Partial update multiple records using $set

```go
// Signature
func Patch[T any](
    condition any,
    data primitive.M,
    silent bool,
    opts ...MongoOption,
) (*mongo.UpdateResult, error)
```

### Increment

Increment numeric data. Pass negative value for decrement.

```go
// Signature
func Increment[T any](
    condition any,
    data any,
    opts ...MongoOption,
) (*mongo.UpdateResult, error)
```
