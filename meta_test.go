package mongoutils_test

import (
	"testing"
	"time"

	"github.com/bopher/mongoutils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMeta(t *testing.T) {
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	id3 := primitive.NewObjectID()
	date := time.Now().UTC()
	date2 := time.Now().UTC().Add(10 * time.Hour)

	setter := mongoutils.NewMetaSetter()
	setter.Add("test", "activity", &id1, date)
	setter.Add("test", "activity", nil, date2)
	setter.Add("test", "activity", &id2, date)
	setter.Add("test", "activity", &id3, date)
	setter.Add("test", "activity", &id3, nil)
	setter.Add("test", "activity", &id3, date2)

	t.Log([]primitive.ObjectID{id1, id2, id3})
	t.Fatalf("%+v", setter.Build())
}
