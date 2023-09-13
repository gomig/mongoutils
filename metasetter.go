package mongoutils

import "go.mongodb.org/mongo-driver/bson/primitive"

type MetaSetter interface {
	// Add new meta
	Add(_col, _meta string, id *primitive.ObjectID, value any) MetaSetter
	// Build get combined meta with query
	Build() []MetaSetterResult
}

type MetaSetterResult struct {
	Col    string
	Ids    []primitive.ObjectID
	Values map[string]any
}
