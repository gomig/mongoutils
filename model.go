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
	// PrepareInsert fill created_at before save
	PrepareInsert()
	// PrepareUpdate fill updated_at before save
	// in ghost mode updated_at field not changed
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
