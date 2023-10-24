package mongoutils

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// EmptyModel only implement model methods
type EmptyModel struct {
	ID primitive.ObjectID `bson:"_id" json:"_id"`
}

func (EmptyModel) TypeName() string {
	return "unspecified"
}

func (EmptyModel) Collection(db *mongo.Database) *mongo.Collection {
	panic("please override collection method")
}

func (EmptyModel) Index(db *mongo.Database) error {
	return nil
}

func (EmptyModel) Seed(db *mongo.Database) error {
	return nil
}

func (EmptyModel) Pipeline() MongoPipeline {
	return NewPipe()
}

func (*EmptyModel) FillCreatedAt() {}

func (*EmptyModel) FillUpdatedAt() {}

func (model *EmptyModel) NewId() {
	model.ID = primitive.NewObjectID()
}

func (model *EmptyModel) SetID(id primitive.ObjectID) {
	model.ID = id
}

func (model EmptyModel) GetID() primitive.ObjectID {
	return model.ID
}

func (EmptyModel) IsEditable() bool {
	return true
}

func (EmptyModel) IsDeletable() bool {
	return false
}

func (*EmptyModel) Cleanup() {}

func (*EmptyModel) OnInsert(ctx context.Context, opt ...MongoOption) {}

func (*EmptyModel) OnUpdate(ctx context.Context, opt ...MongoOption) {}

func (*EmptyModel) OnDelete(ctx context.Context, opt ...MongoOption) {}

func (EmptyModel) OnInserted(ctx context.Context, opt ...MongoOption) {}

func (EmptyModel) OnUpdated(old any, ctx context.Context, opt ...MongoOption) {}

func (EmptyModel) OnDeleted(ctx context.Context, opt ...MongoOption) {}
