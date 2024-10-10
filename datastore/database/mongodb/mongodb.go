package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gnanasuryateja/golib/constants"
	"github.com/gnanasuryateja/golib/datastore/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStoreConfig struct {
	Uri      string
	DbName   string
	Username string
	Password string
}

// validates the input params
func (msc MongoStoreConfig) validate() error {
	if msc.Uri == "" || msc.DbName == "" || msc.Username == "" || msc.Password == "" {
		return fmt.Errorf("MongoStoreConfig cannot have empty values")
	}
	return nil
}

type mongoStore struct {
	client   *mongo.Client
	database *mongo.Database
}

// creates a new mongoStore client
func NewMongoStoreClient(ctx context.Context, mongoStoreConfig MongoStoreConfig) (database.Database, context.CancelFunc, error) {

	// validate the mongoStoreConfig
	err := mongoStoreConfig.validate()
	if err != nil {
		return nil, nil, err
	}

	// update the uri with username and password
	mongoStoreConfig.Uri = strings.Replace(mongoStoreConfig.Uri, constants.USERNAME, mongoStoreConfig.Username, 1)
	mongoStoreConfig.Uri = strings.Replace(mongoStoreConfig.Uri, constants.PASSWORD, mongoStoreConfig.Password, 1)

	// get the ctx and cancel
	ctx, cancel := context.WithCancel(ctx)

	// build the clientOptions
	clientOptions := options.Client().ApplyURI(mongoStoreConfig.Uri)

	// get the mongo client
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, cancel, err
	}

	// check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, cancel, err
	}

	// get the database
	database := client.Database(mongoStoreConfig.DbName)

	// return the mongoStore
	return mongoStore{
		client:   client,
		database: database,
	}, cancel, nil
}

func (db mongoStore) CloseDB(ctx context.Context, cancel context.CancelFunc) {

	// CancelFunc to cancel to context
	defer cancel()

	// client provides a method to close
	// a mongoDB connection.
	defer func() {

		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := db.client.Disconnect(ctx); err != nil {
			fmt.Println("error disconnecting from db: ", err)
			// panic(err)
		}
	}()
}

// checks the connection to db and returns an error if any
func (db mongoStore) HealthCheck(ctx context.Context) error {
	// check the connection
	err := db.client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

// inserts data into the db
func (db mongoStore) AddData(ctx context.Context, args ...any) (string, error) {

	// validate the passed args
	if len(args) < 2 {
		return "", fmt.Errorf("collection name or data is(are) missing")
	}
	if len(args) > 2 {
		return "", fmt.Errorf("more params are passed than expected")
	}

	// get the collection string from args
	coll, ok := args[0].(string)
	if !ok {
		return "", fmt.Errorf("invalid collection name is passed")
	}

	// get the mongo collection
	collection := db.database.Collection(coll)

	// get the data from args
	data := args[1]

	// add the data to the collection
	insertOneResult, err := collection.InsertOne(ctx, data)
	if err != nil {
		return "", fmt.Errorf("error inserting the data %v", err)
	}

	// return the inserted id
	return fmt.Sprintf("Inserted the data with %v", insertOneResult.InsertedID), nil
}

// inserts mutiple data into the db
func (db mongoStore) AddMultipleData(ctx context.Context, args ...any) ([]string, error) {

	// validate the passed args
	if len(args) < 2 {
		return nil, fmt.Errorf("collection name or data is(are) missing")
	}
	if len(args) > 2 {
		return nil, fmt.Errorf("more params are passed than expected")
	}

	// get the collection string from args
	coll, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid collection name is passed")
	}

	// get the mongo collection
	collection := db.database.Collection(coll)

	// get the data from args
	data := args[1]

	// marshal the data to bytes
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// unmarshal the bytes to a slice
	var sliceData []any
	err = json.Unmarshal(dataBytes, &sliceData)
	if err != nil {
		return nil, err
	}

	// add the data to the collection
	insertManyResult, err := collection.InsertMany(ctx, sliceData)
	if err != nil {
		return nil, fmt.Errorf("error inserting the data %v", err)
	}

	// get the inserted ids
	var insertedIds []string
	for id := range insertManyResult.InsertedIDs {
		insertedIds = append(insertedIds, fmt.Sprintf("Inserted the data with %v", insertManyResult.InsertedIDs[id]))
	}

	// return the inserted ids
	return insertedIds, nil
}

// gets the data from db
func (db mongoStore) GetData(ctx context.Context, args ...any) (any, error) {

	// validate the passed args
	if len(args) < 2 {
		return nil, fmt.Errorf("collection name or filter is(are) missing")
	}
	if len(args) > 2 {
		return nil, fmt.Errorf("more params are passed than expected")
	}

	// get the collection string from args
	coll, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid collection name is passed")
	}

	// get the mongo collection
	collection := db.database.Collection(coll)

	// get the filter from args
	filter := args[1]

	// get the data from collection
	var data any
	err := collection.FindOne(ctx, filter).Decode(&data)
	if err != nil {
		return nil, err
	}

	// return the data
	return data, nil
}

// gets multiple data from db
func (db mongoStore) GetMultipleData(ctx context.Context, args ...any) ([]any, error) {

	// validate the passed args
	if len(args) < 2 {
		return nil, fmt.Errorf("collection name or filter is(are) missing")
	}
	if len(args) > 2 {
		return nil, fmt.Errorf("more params are passed than expected")
	}

	// get the collection string from args
	coll, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid collection name is passed")
	}

	// get the mongo collection
	collection := db.database.Collection(coll)

	// get the filter from args
	filter := args[1]

	// get the data from collection
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// check for cursor error
	if cursor.Err() != nil {
		return nil, err
	}

	// get the data from cursor
	var result []any
	for cursor.Next(ctx) {
		var data any
		err = cursor.Decode(&data)
		if err != nil {
			return nil, err
		}
		result = append(result, data)
	}

	// return the data slice
	return result, nil
}

// updates the data in db
func (db mongoStore) UpdateData(ctx context.Context, args ...any) (any, error) {

	// validate the passed args
	if len(args) < 3 {
		return nil, fmt.Errorf("collection name, filter or updateData is(are) missing")
	}
	if len(args) > 3 {
		return nil, fmt.Errorf("more params are passed than expected")
	}

	// get the collection string from args
	coll, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid collection name is passed")
	}

	// get the mongo collection
	collection := db.database.Collection(coll)

	// get the filter from args
	filter := args[1]

	// get the updateData from args
	updateData := args[2]

	// update the data in collection
	updateResult, err := collection.UpdateOne(ctx, filter, updateData)
	if err != nil {
		return nil, err
	}

	// return the modified count
	return fmt.Sprintf("%v document(s) have been updated", updateResult.ModifiedCount), nil
}

// updates multiple data in db
func (db mongoStore) UpdateMultipleData(ctx context.Context, args ...any) (any, error) {

	// validate the passed args
	if len(args) < 3 {
		return nil, fmt.Errorf("collection name, filter or updateData is(are) missing")
	}
	if len(args) > 3 {
		return nil, fmt.Errorf("more params are passed than expected")
	}

	// get the collection string from args
	coll, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid collection name is passed")
	}

	// get the mongo collection
	collection := db.database.Collection(coll)

	// get the filter from args
	filter := args[1]

	// get the updateData from args
	updateData := args[2]

	// update the data in collection
	updateResult, err := collection.UpdateMany(ctx, filter, updateData)
	if err != nil {
		return nil, err
	}

	// return the modified count
	return fmt.Sprintf("%v document(s) have been updated", updateResult.ModifiedCount), nil
}

// deletes the data from db
func (db mongoStore) DeleteData(ctx context.Context, args ...any) (any, error) {

	// validate the passed args
	if len(args) < 2 {
		return nil, fmt.Errorf("collection name or filter is(are) missing")
	}
	if len(args) > 2 {
		return nil, fmt.Errorf("more params are passed than expected")
	}

	// get the collection string from args
	coll, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid collection name is passed")
	}

	// get the mongo collection
	collection := db.database.Collection(coll)

	// get the filter from args
	filter := args[1]

	// delete the data from collection
	deleteResult, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	// return the deleted count
	return fmt.Sprintf("%v document(s) have been updated", deleteResult.DeletedCount), nil
}

func (db mongoStore) DeleteMultipleData(ctx context.Context, args ...any) (any, error) {

	// validate the passed args
	if len(args) < 2 {
		return nil, fmt.Errorf("collection name or filter is(are) missing")
	}
	if len(args) > 2 {
		return nil, fmt.Errorf("more params are passed than expected")
	}

	// get the collection string from args
	coll, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid collection name is passed")
	}

	// get the mongo collection
	collection := db.database.Collection(coll)

	// get the filter from args
	filter := args[1]

	// delete the data from collection
	deleteResult, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return nil, err
	}

	// return the deleted count
	return fmt.Sprintf("%v document(s) have been updated", deleteResult.DeletedCount), nil
}
