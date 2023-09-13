package mongoutils_test

import (
	"encoding/json"
	"testing"

	"github.com/bopher/mongoutils"
)

func pretty(v any) (string, error) {
	bytes, err := json.Marshal(v)
	return string(bytes), err
}

func TestDoc(t *testing.T) {
	var v string
	var err error

	// Add
	v, err = pretty(mongoutils.NewDoc().Add("name", "John").Map())
	if err != nil {
		t.Fatal(err)
	}
	if v != `{"name":"John"}` {
		t.Log(v)
		t.Fatal("fail Add")
	}

	// Doc
	v, err = pretty(mongoutils.NewDoc().Doc("$set", func(d mongoutils.MongoDoc) mongoutils.MongoDoc {
		return d.Add("name", "Jack")
	}).Map())
	if err != nil {
		t.Fatal(err)
	}
	if v != `{"$set":[{"Key":"name","Value":"Jack"}]}` {
		t.Log(v)
		t.Fatal("fail Doc")
	}

	// Array
	v, err = pretty(mongoutils.NewDoc().Array("skills", "js", "go", "mongo").Map())
	if err != nil {
		t.Fatal(err)
	}
	if v != `{"skills":["js","go","mongo"]}` {
		t.Log(v)
		t.Fatal("fail Array")
	}

	// DocArray
	v, err = pretty(mongoutils.NewDoc().DocArray("pipeline", func(d mongoutils.MongoDoc) mongoutils.MongoDoc {
		return d.Add("name", "Jack").
			Add("family", "Ma")
	}).Map())
	if err != nil {
		t.Fatal(err)
	}
	if v != `{"pipeline":[{"name":"Jack"},{"family":"Ma"}]}` {
		t.Log(v)
		t.Fatal("fail DocArray")
	}

	// Nested
	v, err = pretty(mongoutils.NewDoc().Nested("$set", "name", "Kim").Map())
	if err != nil {
		t.Fatal(err)
	}
	if v != `{"$set":{"name":"Kim"}}` {
		t.Log(v)
		t.Fatal("fail Nested")
	}

	// NestedDoc
	v, err = pretty(mongoutils.NewDoc().NestedDoc("$set", "name", func(d mongoutils.MongoDoc) mongoutils.MongoDoc {
		return d.Add("first", "jack").Add("last", "ma")
	}).Map())
	if err != nil {
		t.Fatal(err)
	}
	if v != `{"$set":{"name":[{"Key":"first","Value":"jack"},{"Key":"last","Value":"ma"}]}}` {
		t.Log(v)
		t.Fatal("fail NestedDoc")
	}

	// NestedArray
	v, err = pretty(mongoutils.NewDoc().NestedArray("skill", "$in", "js", "mongo").Map())
	if err != nil {
		t.Fatal(err)
	}
	if v != `{"skill":{"$in":["js","mongo"]}}` {
		t.Log(v)
		t.Fatal("fail NestedArray")
	}

	// NestedDocArray
	v, err = pretty(mongoutils.NewDoc().NestedDocArray("let", "name", func(d mongoutils.MongoDoc) mongoutils.MongoDoc {
		return d.Add("name", "Jack").
			Add("family", "Ma")
	}).Map())
	if err != nil {
		t.Fatal(err)
	}
	if v != `{"let":{"name":[{"name":"Jack"},{"family":"Ma"}]}}` {
		t.Log(v)
		t.Fatal("fail NestedDocArray")
	}

	// Regex
	v, err = pretty(mongoutils.NewDoc().Regex("name", "Jo.*", "i").Map())
	if err != nil {
		t.Fatal(err)
	}
	if v != `{"name":{"Pattern":"Jo.*","Options":"i"}}` {
		t.Log(v)
		t.Fatal("fail Regex")
	}
}
