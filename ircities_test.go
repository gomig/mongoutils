package mongoutils_test

import (
	"context"
	"testing"

	"github.com/gomig/mongoutils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestIrCity(t *testing.T) {
	host := "mongodb://127.0.0.1:27017/?directConnection=true"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(host))
	if err != nil {
		t.Fatal(err)
	}

	irc := new(mongoutils.IrCity)
	db := client.Database("test")
	err = irc.Index(db)
	if err != nil {
		t.Fatal(err)
	}
	err = irc.Seed(db)
	if err != nil {
		t.Fatal(err)
	}
}
