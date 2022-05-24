package db

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

type KeyValue map[string]interface{}

type MongoStore struct {
	db      *mongo.Database
	Client  *mongo.Client
	Logger  *logrus.Logger
	Context context.Context
}

func connectMongo(uri string) (*mongo.Client, error) {
	serverApiOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(uri).SetServerAPIOptions(serverApiOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	return client, err
}

func (c *MongoStore) Insert(collectionName string, document interface{}) error {
	collection := c.db.Collection(collectionName)
	_, err := collection.InsertOne(c.Context, document)
	return err
}

func (c *MongoStore) UpdateOne(
	collectionName string,
	filter KeyValue,
	document interface{},
	opt options.UpdateOptions) error {
	collection := c.db.Collection(collectionName)
	_, err := collection.UpdateOne(c.Context, filter, document, &opt)
	return err
}

func (c *MongoStore) Get(collectionName string, filter KeyValue, document interface{}) error {
	collection := c.db.Collection(collectionName)
	return collection.FindOne(c.Context, filter).Decode(document)
}

func (c *MongoStore) FindOneAndUpdate(
	collectionName string,
	filter KeyValue,
	document interface{},
	opt options.FindOneAndUpdateOptions) (bson.M, error) {
	collection := c.db.Collection(collectionName)
	result := collection.FindOneAndUpdate(c.Context, filter, document, &opt)
	if result.Err() != nil {
		return nil, result.Err()
	}

	doc := bson.M{}
	decodeErr := result.Decode(&doc)
	return doc, decodeErr
}

func (c *MongoStore) Delete(collectionName string, filter KeyValue) error {
	collection := c.db.Collection(collectionName)
	_, err := collection.DeleteOne(c.Context, filter)
	return err
}

func (c *MongoStore) GetAll(collectionName string, filter KeyValue, documents interface{}) error {
	collection := c.db.Collection(collectionName)
	cursor, err := collection.Find(c.Context, filter)
	if err != nil {
		return err
	}
	return cursor.All(c.Context, documents)
}

func NewMongoStore() *MongoStore {
	var connectOnce sync.Once
	var client *mongo.Client
	var err error
	ctx := context.Background()
	// Ensure we connect only once
	connectOnce.Do(func() {
		client, err = connectMongo("mongodb+srv://gobizdb:admin@gobiz.tiecp.mongodb.net/?retryWrites=true&w=majority")
	})

	if err != nil {
		logrus.Fatalf("Failed to connect to database %s:", err)
	}
	return &MongoStore{
		db:      client.Database("prodo"),
		Client:  client,
		Logger:  &logrus.Logger{},
		Context: ctx,
	}
}
