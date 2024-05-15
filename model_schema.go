package mongoutils

// SchemaModel schema versioning field
type SchemaModel struct {
	SchemaVersion int `bson:"schema_version" json:"-"`
}

func (model SchemaModel) GetVersion() int {
	return model.SchemaVersion
}

func (model *SchemaModel) SetVersion(v int) {
	model.SchemaVersion = v
}
