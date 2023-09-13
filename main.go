package mongoutils

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// NewPipe new mongo pipe builder
func NewPipe() MongoPipeline {
	return new(mPipe)
}

// NewDoc new mongo doc builder
func NewDoc() MongoDoc {
	return new(mDoc)
}

// NewMetaCounter new mongo meta counter
func NewMetaCounter() MetaCounter {
	res := new(metaCounter)
	res.Data = make(map[string][]meta)
	return res
}

// NewMetaSetter new mongo meta setter
func NewMetaSetter() MetaSetter {
	res := new(metaSetter)
	res.Data = make(map[string][]metaV)
	return res
}

// MongoOperationCtx create context for mongo db operations for 10 sec
func MongoOperationCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.TODO(), 10*time.Second)
}

// ParseObjectID parse object id from string
func ParseObjectID(id string) *primitive.ObjectID {
	if oId, err := primitive.ObjectIDFromHex(id); err == nil && !oId.IsZero() {
		return &oId
	}
	return nil
}

// IsValidObjectId check if object id is valid and not zero
func IsValidObjectId(id *primitive.ObjectID) bool {
	return id != nil && !id.IsZero()
}

// FindOption generate find option with sorts params
func FindOption(sort any, skip int64, limit int64) *options.FindOptions {
	opt := new(options.FindOptions)
	opt.SetAllowDiskUse(true)
	opt.SetSkip(skip)
	if limit > 0 {
		opt.SetLimit(limit)
	}
	if sort != nil {
		opt.SetSort(sort)
	}
	return opt
}

// AggregateOption generate aggregation options
func AggregateOption() *options.AggregateOptions {
	return new(options.AggregateOptions).
		SetAllowDiskUse(true)
}

// TxOption generate transaction option with majority write and snapshot read
func TxOption() *options.TransactionOptions {
	return options.
		Transaction().
		SetWriteConcern(writeconcern.New(writeconcern.WMajority())).
		SetReadConcern(readconcern.Snapshot())
}
