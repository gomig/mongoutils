package mongoutils

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// model interface
type Model interface {
	// TypeName get type string
	TypeName() string
	// Collection get model collection
	Collection(db *mongo.Database) *mongo.Collection
	// Indexes create model indexes
	Index(db *mongo.Database) error
	// Seed run model seed
	Seed(db *mongo.Database) error
	// Pipeline get model pipeline
	Pipeline() MongoPipeline
	// FillCreatedAt fill created_at parameter
	FillCreatedAt()
	// FillUpdatedAt fill updated_at parameter
	FillUpdatedAt()
	// NewId generate new id for model
	NewId()
	// SetID set model id
	SetID(id primitive.ObjectID)
	// ID get model id
	GetID() primitive.ObjectID
	// IsEditable check if document is editable
	// by default returns true on BaseModel
	IsEditable() bool
	// IsDeletable check if document is deletable
	// by default returns false on BaseModel
	IsDeletable() bool
	// Cleanup document before save
	// e.g set document field nil for ignore saving
	Cleanup()
	// OnInsert function to call before insert
	OnInsert(ctx context.Context, opt ...MongoOption)
	// OnUpdate function to call before update
	OnUpdate(ctx context.Context, opt ...MongoOption)
	// OnDelete function to call before delete
	OnDelete(ctx context.Context, opt ...MongoOption)
	// OnInserted function to call after insert
	OnInserted(ctx context.Context, opt ...MongoOption)
	// OnUpdated function to call after update
	OnUpdated(old any, ctx context.Context, opt ...MongoOption)
	// OnDeleted function to call after delete
	OnDeleted(ctx context.Context, opt ...MongoOption)
}

type SchemaVersioning interface {
	// GetVersion get schema version
	GetVersion() int
	// SetVersion set schema version
	SetVersion(int)
}

type SoftDelete interface {
	// SoftDelete set deleted_at field to current date
	SoftDelete()
	// Restore set deleted_at field to nil
	Restore()
	// IsDeleted check if item soft deleted
	IsDeleted() bool
}

type Backup interface {
	// ToMap get model as map for backup
	// return nil or empty map to skip backup
	ToMap() map[string]any
	// SetChecksum set model md5 checksum
	SetChecksum(string)
	// GetChecksum get model md5 checksum
	GetChecksum() string
	// NeedBackup check if record need backup
	NeedBackup() bool
	// MarkBackup set backup state to current date
	MarkBackup()
	// UnMarkBackup set backup state to nil
	UnMarkBackup()
}
