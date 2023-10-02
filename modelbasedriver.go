package mongoutils

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// BaseModel implementation with id and timestamp
type BaseModel struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  *time.Time         `bson:"updated_at" json:"updated_at"`
	LastBackup *time.Time         `bson:"last_backup" json:"last_backup"`
}

func (*BaseModel) TypeName() string {
	return "unspecified"
}

func (*BaseModel) Collection(db *mongo.Database) *mongo.Collection {
	panic("please override collection method")
}

func (*BaseModel) Index(db *mongo.Database) error {
	return nil
}

func (*BaseModel) Seed(db *mongo.Database) error {
	return nil
}

func (*BaseModel) Pipeline() MongoPipeline {
	return NewPipe()
}

func (model *BaseModel) NewId() {
	model.ID = primitive.NewObjectID()
}

func (model *BaseModel) SetID(id primitive.ObjectID) {
	model.ID = id
}

func (model *BaseModel) GetID() primitive.ObjectID {
	return model.ID
}

func (*BaseModel) IsEditable() bool {
	return true
}

func (*BaseModel) IsDeletable() bool {
	return false
}

func (model *BaseModel) NeedBackup() bool {
	return model.LastBackup == nil
}

func (model *BaseModel) MarkBackup() {
	t := time.Now().UTC()
	model.LastBackup = &t
}

func (model *BaseModel) UnMarkBackup() {
	model.LastBackup = nil
}

func (*BaseModel) Cleanup() {}

func (model *BaseModel) PrepareInsert() {
	model.CreatedAt = time.Now().UTC()
	model.UnMarkBackup()
}

func (model *BaseModel) PrepareUpdate(ghost bool) {
	if !ghost {
		now := time.Now().UTC()
		model.UpdatedAt = &now
	}
	model.UnMarkBackup()
}

func (*BaseModel) OnInsert(ctx context.Context, opt ...MongoOption) {}

func (*BaseModel) OnUpdate(ctx context.Context, opt ...MongoOption) {}

func (*BaseModel) OnDelete(ctx context.Context, opt ...MongoOption) {}

func (*BaseModel) OnInserted(ctx context.Context, opt ...MongoOption) {}

func (*BaseModel) OnUpdated(old any, ctx context.Context, opt ...MongoOption) {}

func (*BaseModel) OnDeleted(ctx context.Context, opt ...MongoOption) {}
