package mongoutils

import "go.mongodb.org/mongo-driver/bson/primitive"

type MetaCounter interface {
	// Add new meta
	Add(_col, _meta string, id *primitive.ObjectID, amount int) MetaCounter
	// Build get combined meta with query
	Build() []MetaCounterResult
}

type MetaCounterResult struct {
	Col string
	Ids []primitive.ObjectID
	// data to update
	Values map[string]int
}
