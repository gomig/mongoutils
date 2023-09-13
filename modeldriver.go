package mongoutils

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// EmptyModel only implement model methods
type EmptyModel struct{}

func (*EmptyModel) TypeName() string {
	return "unspecified"
}

func (*EmptyModel) Collection(db *mongo.Database) *mongo.Collection {
	panic("please override collection method")
}

func (*EmptyModel) Index(db *mongo.Database) error {
	return nil
}

func (*EmptyModel) Seed(db *mongo.Database) error {
	return nil
}

func (*EmptyModel) Pipeline() MongoPipeline {
	return NewPipe()
}

func (*EmptyModel) NewId() {}

func (*EmptyModel) SetID(id primitive.ObjectID) {}

func (*EmptyModel) GetID() primitive.ObjectID {
	return primitive.NilObjectID
}

func (*EmptyModel) IsEditable() bool {
	return true
}

func (*EmptyModel) IsDeletable() bool {
	return false
}

func (*EmptyModel) Cleanup() {}

func (*EmptyModel) PrepareInsert() {}

func (*EmptyModel) PrepareUpdate(ghost bool) {}

func (*EmptyModel) OnInsert(ctx context.Context, opt ...MongoOption) {}

func (*EmptyModel) OnUpdate(ctx context.Context, opt ...MongoOption) {}

func (*EmptyModel) OnDelete(ctx context.Context, opt ...MongoOption) {}

func (*EmptyModel) OnInserted(ctx context.Context, opt ...MongoOption) {}

func (*EmptyModel) OnUpdated(old any, ctx context.Context, opt ...MongoOption) {}

func (*EmptyModel) OnDeleted(ctx context.Context, opt ...MongoOption) {}
