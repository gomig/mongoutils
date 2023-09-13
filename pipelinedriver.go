package mongoutils

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type mPipe struct {
	data mongo.Pipeline
}

func (me *mPipe) Add(cb func(d MongoDoc) MongoDoc) MongoPipeline {
	me.data = append(me.data, cb(NewDoc()).Build())
	return me
}

func (me *mPipe) Match(filters any) MongoPipeline {
	if filters == nil {
		return me
	}
	return me.Add(func(d MongoDoc) MongoDoc {
		return d.Add("$match", filters)
	})
}

func (me *mPipe) In(key string, v any) MongoPipeline {
	return me.Add(func(d MongoDoc) MongoDoc {
		return d.Nested(key, "$in", v)
	})
}

func (me *mPipe) Limit(limit int64) MongoPipeline {
	if limit > 0 {
		me.Add(func(d MongoDoc) MongoDoc {
			return d.Add("$limit", limit)
		})
	}
	return me
}

func (me *mPipe) Skip(skip int64) MongoPipeline {
	if skip > 0 {
		me.Add(func(d MongoDoc) MongoDoc {
			return d.Add("$skip", skip)
		})
	}
	return me
}

func (me *mPipe) Sort(sorts any) MongoPipeline {
	if sorts != nil {
		me.Add(func(d MongoDoc) MongoDoc {
			return d.Add("$sort", sorts)
		})
	}
	return me
}

func (me *mPipe) Unwind(path string, prevNullAndEmpty bool) MongoPipeline {
	return me.Add(func(d MongoDoc) MongoDoc {
		return d.Doc("$unwind", func(d MongoDoc) MongoDoc {
			return d.
				Add("path", path).
				Add("preserveNullAndEmptyArrays", prevNullAndEmpty)
		})
	})
}

func (me *mPipe) Lookup(from string, local string, foreign string, as string) MongoPipeline {
	return me.Add(func(d MongoDoc) MongoDoc {
		return d.Doc("$lookup", func(d MongoDoc) MongoDoc {
			return d.
				Add("from", from).
				Add("localField", local).
				Add("foreignField", foreign).
				Add("as", as)
		})
	})
}

func (me *mPipe) Unwrap(field string, as string) MongoPipeline {
	return me.Add(func(d MongoDoc) MongoDoc {
		return d.Doc("$addFields", func(d MongoDoc) MongoDoc {
			return d.Nested(as, "$first", field)
		})
	})
}

func (me *mPipe) LoadRelation(from string, local string, foreign string, as string) MongoPipeline {
	me.Add(func(d MongoDoc) MongoDoc {
		return d.Doc("$lookup", func(d MongoDoc) MongoDoc {
			return d.
				Add("from", from).
				Add("localField", local).
				Add("foreignField", foreign).
				Add("as", as)
		})
	})
	me.Add(func(d MongoDoc) MongoDoc {
		return d.Doc("$addFields", func(d MongoDoc) MongoDoc {
			return d.Nested(as, "$first", "$"+as)
		})
	})
	return me
}

func (me *mPipe) Group(cb func(d MongoDoc) MongoDoc) MongoPipeline {
	return me.Add(func(d MongoDoc) MongoDoc {
		return d.Doc("$group", cb)
	})
}

func (me *mPipe) ReplaceRoot(v any) MongoPipeline {
	return me.Add(func(d MongoDoc) MongoDoc {
		return d.Doc("$replaceRoot", func(d MongoDoc) MongoDoc {
			return d.Add("newRoot", v)
		})
	})
}

func (me *mPipe) MergeRoot(fields ...any) MongoPipeline {
	return me.Add(func(d MongoDoc) MongoDoc {
		return d.Doc("$replaceRoot", func(d MongoDoc) MongoDoc {
			return d.Nested("newRoot", "$mergeObjects", fields)
		})
	})
}

func (me *mPipe) UnProject(fields ...string) MongoPipeline {
	return me.Add(func(d MongoDoc) MongoDoc {
		return d.Doc("$project", func(d MongoDoc) MongoDoc {
			for _, v := range fields {
				d.Add(v, 0)
			}
			return d
		})
	})
}

func (me *mPipe) Project(projects any) MongoPipeline {
	if projects == nil {
		return me
	}
	return me.Add(func(d MongoDoc) MongoDoc {
		return d.Add("$project", projects)
	})
}

func (me mPipe) Build() mongo.Pipeline {
	return me.data
}
