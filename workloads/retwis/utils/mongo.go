package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

const kLocalhostMongoUri = "mongodb://localhost:27017"

func getMongoUri() string {
	if uri, exists := os.LookupEnv("MONGODB_URI"); exists {
		return uri
	} else {
		return kLocalhostMongoUri
	}
}

func CreateMongoClientOrDie(ctx context.Context) *mongo.Client {
	uri := getMongoUri()
	opts := options.Client()
	opts.ApplyURI(uri)
	opts.SetReadConcern(readconcern.Majority())
	opts.SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
	opts.SetReadPreference(readpref.PrimaryPreferred())
	newCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if client, err := mongo.Connect(newCtx, opts); err != nil {
		log.Fatalf("[FATAL] Failed to connect to mongo %s: %v", uri, err)
		return nil
	} else {
		return client
	}
}

func MongoTxnOptions() *options.TransactionOptions {
	opts := options.Transaction()
	opts.SetReadConcern(readconcern.Snapshot())
	opts.SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
	opts.SetReadPreference(readpref.Primary())
	return opts
}

func MongoCreateCounter(ctx context.Context, db *mongo.Database, name string) error {
	collection := db.Collection("counters")
	_, err := collection.InsertOne(ctx, bson.D{{"name", name}, {"value", int32(0)}})
	return err
}

func MongoFetchAddCounter(ctx context.Context, db *mongo.Database, name string, delta int) (int, error) {
	collection := db.Collection("counters")
	filter := bson.D{{"name", name}}
	update := bson.D{{"$inc", bson.D{{"value", int32(delta)}}}}
	var updatedDocument bson.M
	err := collection.FindOneAndUpdate(ctx, filter, update).Decode(&updatedDocument)
	if err != nil {
		return 0, err
	}
	if value, ok := updatedDocument["value"].(int32); ok {
		return int(value), nil
	} else {
		return 0, fmt.Errorf("%s value cannot be converted to int32", name)
	}
}

func MongoCreateIndex(ctx context.Context, collection *mongo.Collection, key string, unique bool) error {
	indexOpts := options.Index().SetUnique(unique)
	mod := mongo.IndexModel{
		Keys:    bson.M{key: 1},
		Options: indexOpts,
	}
	_, err := collection.Indexes().CreateOne(ctx, mod)
	return err
}
