package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Database struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// Initiate database connection to MongoDB and returns Database instance
func Init(uri string, name string) (*Database, error) {
	client, err := connect(uri)
	if err != nil {
		return nil, fmt.Errorf("error while connecting to the database: %w", err)
	}
	instance := new(Database)
	db := client.Database(name)
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	if err := ensureIndex(db); err != nil {
		return nil, fmt.Errorf("error while checking for database indexes: %w", err)
	}

	instance.Client = client
	instance.DB = db

	return instance, nil
}

func connect(uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, nil
}

// EnsureIndex makes sure if the index is already created before doing db operations
func ensureIndex(db *mongo.Database) error {
	urlsColl := db.Collection("urls")
	models := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "code", Value: 1},
			},
			Options: options.Index().SetUnique(true).SetName("code_unique"),
		},
	}
	_, err := urlsColl.Indexes().CreateMany(context.TODO(), models)
	if err != nil {
		return err
	}
	return nil
}
