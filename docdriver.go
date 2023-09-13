package mongoutils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mDoc struct {
	data primitive.D
}

func (me *mDoc) Add(k string, v any) MongoDoc {
	me.data = append(me.data, primitive.E{Key: k, Value: v})
	return me
}

func (me *mDoc) Doc(k string, cb func(d MongoDoc) MongoDoc) MongoDoc {
	return me.Add(k, cb(NewDoc()).Build())
}

func (me *mDoc) Array(k string, v ...any) MongoDoc {
	return me.Add(k, v)
}

func (me *mDoc) DocArray(k string, cb func(d MongoDoc) MongoDoc) MongoDoc {
	return me.Array(k, arrayOf(cb(NewDoc()))...)
}

func (me *mDoc) Nested(root string, k string, v any) MongoDoc {
	return me.Add(root, primitive.M{k: v})
}

func (me *mDoc) NestedDoc(root string, k string, cb func(d MongoDoc) MongoDoc) MongoDoc {
	return me.Add(root, primitive.M{k: cb(NewDoc()).Build()})
}

func (me *mDoc) NestedArray(root string, k string, v ...any) MongoDoc {
	return me.Add(root, primitive.M{k: v})
}

func (me *mDoc) NestedDocArray(root string, k string, cb func(d MongoDoc) MongoDoc) MongoDoc {
	return me.NestedArray(root, k, arrayOf(cb(NewDoc()))...)
}

func (me *mDoc) Regex(k string, pattern string, opt string) MongoDoc {
	return me.Add(k, primitive.Regex{Pattern: pattern, Options: opt})
}

func (me mDoc) Map() primitive.M {
	return me.data.Map()
}

func (me mDoc) Build() primitive.D {
	return me.data
}

func arrayOf(d MongoDoc) []any {
	data := d.Build()
	res := make([]any, len(data))
	for i, e := range data {
		res[i] = primitive.M{e.Key: e.Value}
	}
	return res
}
