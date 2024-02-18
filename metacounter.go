package mongoutils

import "go.mongodb.org/mongo-driver/bson/primitive"

type MetaCounter interface {
	// Add increase meta amount
	Add(_col, _meta string, id *primitive.ObjectID, amount int64) MetaCounter
	// Sub decrease meta amount
	Sub(_col, _meta string, id *primitive.ObjectID, amount int64) MetaCounter
	// Build get combined meta with query
	Build() []MetaCounterResult
}

type MetaCounterResult struct {
	Col string
	Ids []primitive.ObjectID
	// data to update
	Values map[string]int64
}
