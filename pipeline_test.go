package mongoutils_test

import (
	"testing"

	"github.com/bopher/mongoutils"
)

func TestPipeline(t *testing.T) {
	var v string
	var err error

	// Add
	v, err = pretty(mongoutils.NewPipe().Add(func(d mongoutils.MongoDoc) mongoutils.MongoDoc {
		return d.Nested("$match", "name", "Jack")
	}).Build())
	if err != nil {
		t.Fatal(err)
	}
	if v != `[[{"Key":"$match","Value":{"name":"Jack"}}]]` {
		t.Log(v)
		t.Fatal("fail Add")
	}

	// Match
	v, err = pretty(mongoutils.NewPipe().Match("TEST").Build())
	if err != nil {
		t.Fatal(err)
	}
	if v != `[[{"Key":"$match","Value":"TEST"}]]` {
		t.Log(v)
		t.Fatal("fail Match")
	}

	// In
	v, err = pretty(mongoutils.NewPipe().In("name", []string{"a", "b", "c"}).Build())
	if err != nil {
		t.Fatal(err)
	}
	if v != `[[{"Key":"name","Value":{"$in":["a","b","c"]}}]]` {
		t.Log(v)
		t.Fatal("fail In")
	}

	// Limit
	v, err = pretty(mongoutils.NewPipe().Limit(4).Build())
	if err != nil {
		t.Fatal(err)
	}
	if v != `[[{"Key":"$limit","Value":4}]]` {
		t.Log(v)
		t.Fatal("fail Limit")
	}

	// Skip
	v, err = pretty(mongoutils.NewPipe().Skip(4).Build())
	if err != nil {
		t.Fatal(err)
	}
	if v != `[[{"Key":"$skip","Value":4}]]` {
		t.Log(v)
		t.Fatal("fail Skip")
	}

	// Sort
	v, err = pretty(mongoutils.NewPipe().Sort("sorts...").Build())
	if err != nil {
		t.Fatal(err)
	}
	if v != `[[{"Key":"$sort","Value":"sorts..."}]]` {
		t.Log(v)
		t.Fatal("fail Sort")
	}

	// Unwind
	v, err = pretty(mongoutils.NewPipe().Unwind("services", true).Build())
	if err != nil {
		t.Fatal(err)
	}
	if v != `[[{"Key":"$unwind","Value":[{"Key":"path","Value":"services"},{"Key":"preserveNullAndEmptyArrays","Value":true}]}]]` {
		t.Log(v)
		t.Fatal("fail Unwind")
	}

	// Lookup
	v, err = pretty(mongoutils.NewPipe().Lookup("users", "user_id", "_id", "user").Build())
	if err != nil {
		t.Fatal(err)
	}
	if v != `[[{"Key":"$lookup","Value":[{"Key":"from","Value":"users"},{"Key":"localField","Value":"user_id"},{"Key":"foreignField","Value":"_id"},{"Key":"as","Value":"user"}]}]]` {
		t.Log(v)
		t.Fatal("fail Lookup")
	}

	// Unwrap
	v, err = pretty(mongoutils.NewPipe().Unwrap("_user", "user").Build())
	if err != nil {
		t.Fatal(err)
	}
	if v != `[[{"Key":"$addFields","Value":[{"Key":"user","Value":{"$first":"_user"}}]}]]` {
		t.Log(v)
		t.Fatal("fail Unwrap")
	}

	// LoadRelation
	v, err = pretty(mongoutils.NewPipe().LoadRelation("users", "user_id", "_id", "user").Build())
	if err != nil {
		t.Fatal(err)
	}
	if v != `[[{"Key":"$lookup","Value":[{"Key":"from","Value":"users"},{"Key":"localField","Value":"user_id"},{"Key":"foreignField","Value":"_id"},{"Key":"as","Value":"user"}]}],[{"Key":"$addFields","Value":[{"Key":"user","Value":{"$first":"$user"}}]}]]` {
		t.Log(v)
		t.Fatal("fail LoadRelation")
	}

	// Group
	v, err = pretty(mongoutils.NewPipe().Group(func(d mongoutils.MongoDoc) mongoutils.MongoDoc {
		d.
			Add("_id", "$_id").
			Nested("name", "$first", "$name").
			Nested("total", "$sum", "$invoice")
		return d
	}).Build())
	if err != nil {
		t.Fatal(err)
	}
	if v != `[[{"Key":"$group","Value":[{"Key":"_id","Value":"$_id"},{"Key":"name","Value":{"$first":"$name"}},{"Key":"total","Value":{"$sum":"$invoice"}}]}]]` {
		t.Log(v)
		t.Fatal("fail Group")
	}

	// ReplaceRoot
	v, err = pretty(mongoutils.NewPipe().ReplaceRoot("$my_root").Build())
	if err != nil {
		t.Fatal(err)
	}
	if v != `[[{"Key":"$replaceRoot","Value":[{"Key":"newRoot","Value":"$my_root"}]}]]` {
		t.Log(v)
		t.Fatal("fail ReplaceRoot")
	}

	// MergeRoot
	v, err = pretty(mongoutils.NewPipe().MergeRoot("$my_root", "$$ROOT").Build())
	if err != nil {
		t.Fatal(err)
	}
	if v != `[[{"Key":"$replaceRoot","Value":[{"Key":"newRoot","Value":{"$mergeObjects":["$my_root","$$ROOT"]}}]}]]` {
		t.Log(v)
		t.Fatal("fail MergeRoot")
	}

	// UnProject
	v, err = pretty(mongoutils.NewPipe().UnProject("my_root", "__user").Build())
	if err != nil {
		t.Fatal(err)
	}
	if v != `[[{"Key":"$project","Value":[{"Key":"my_root","Value":0},{"Key":"__user","Value":0}]}]]` {
		t.Log(v)
		t.Fatal("fail UnProject")
	}
}
