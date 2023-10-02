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
	// PrepareInsert called by mongoutils Repository before insert
	// this method fill created_at if not set on BaseModel
	// this method fill Checksum and LastBackup if LastBackup if model implement Backup
	PrepareInsert()
	// PrepareInsert called by mongoutils Repository before update
	// this method fill updated_at on BaseModel
	// updated_at not changed if model implement Backup and backup data not change
	// updated_at not changed if ghost mode is true
	// this method fill Checksum and LastBackup if LastBackup if model implement Backup
	PrepareUpdate(ghost bool)
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
