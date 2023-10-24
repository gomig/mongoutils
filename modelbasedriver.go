package mongoutils

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// BaseModel implementation with id and timestamp
type BaseModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt *time.Time         `bson:"updated_at" json:"updated_at"`
}

func (BaseModel) TypeName() string {
	return "unspecified"
}

func (BaseModel) Collection(db *mongo.Database) *mongo.Collection {
	panic("please override collection method")
}

func (BaseModel) Index(db *mongo.Database) error {
	return nil
}

func (BaseModel) Seed(db *mongo.Database) error {
	return nil
}

func (BaseModel) Pipeline() MongoPipeline {
	return NewPipe()
}

func (model *BaseModel) NewId() {
	model.ID = primitive.NewObjectID()
}

func (model *BaseModel) SetID(id primitive.ObjectID) {
	model.ID = id
}

func (model BaseModel) GetID() primitive.ObjectID {
	return model.ID
}

func (BaseModel) IsEditable() bool {
	return true
}

func (BaseModel) IsDeletable() bool {
	return false
}

func (*BaseModel) Cleanup() {}

func (model *BaseModel) PrepareInsert() {
	if model.CreatedAt.IsZero() {
		model.CreatedAt = time.Now().UTC()
	}
	if b, ok := parseAsBackup(model); ok && b.CanBackup() {
		b.SetChecksum(b.MD5())
		b.UnMarkBackup()
	}
}

func (model *BaseModel) PrepareUpdate(ghost bool) {
	isChanged := true
	if b, ok := parseAsBackup(model); ok && b.CanBackup() {
		newCS := b.MD5()
		if b.GetChecksum() != newCS {
			isChanged = true
			b.SetChecksum(newCS)
			b.UnMarkBackup()
		} else {
			isChanged = false
		}
	}
	if !ghost && isChanged {
		now := time.Now().UTC()
		model.UpdatedAt = &now
	}
}

func (*BaseModel) OnInsert(ctx context.Context, opt ...MongoOption) {}

func (*BaseModel) OnUpdate(ctx context.Context, opt ...MongoOption) {}

func (*BaseModel) OnDelete(ctx context.Context, opt ...MongoOption) {}

func (*BaseModel) OnInserted(ctx context.Context, opt ...MongoOption) {}

func (*BaseModel) OnUpdated(old any, ctx context.Context, opt ...MongoOption) {}

func (*BaseModel) OnDeleted(ctx context.Context, opt ...MongoOption) {}
