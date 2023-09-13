package mongoutils

import "go.mongodb.org/mongo-driver/bson/primitive"

// MongoDoc mongo document (primitive.D) builder
type MongoDoc interface {
	// Add add new element
	Add(k string, v any) MongoDoc
	// Doc add new element with nested doc value
	Doc(k string, cb func(d MongoDoc) MongoDoc) MongoDoc
	// Array add new element with array value
	Array(k string, v ...any) MongoDoc
	// DocArray add new array element with doc
	DocArray(k string, cb func(d MongoDoc) MongoDoc) MongoDoc
	// Nested add new nested element
	Nested(root string, k string, v any) MongoDoc
	// NestedDoc add new nested element with doc value
	NestedDoc(root string, k string, cb func(d MongoDoc) MongoDoc) MongoDoc
	// NestedArray add new nested element with array value
	NestedArray(root string, k string, v ...any) MongoDoc
	// NestedDocArray add new nested array element with doc
	NestedDocArray(root string, k string, cb func(d MongoDoc) MongoDoc) MongoDoc
	// Regex add new element with regex value
	Regex(k string, pattern string, opt string) MongoDoc
	// Map creates a map from the elements of the Doc
	Map() primitive.M
	// Build generate mongo doc
	Build() primitive.D
}
