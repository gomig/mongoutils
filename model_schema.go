package mongoutils

// SchemaModel schema versioning field
type SchemaModel struct {
	SchemaVersion int `bson:"schema_version" json:"-"`
}
