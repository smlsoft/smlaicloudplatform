package microservice

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IPersister is interface for persister
type IPersisterMongo interface {
	Find(model interface{}, filter interface{}) ([]interface{}, error)
	FindOne(model interface{}, id string) (interface{}, error)
	Create(model interface{}) error
	Update(model interface{}, id string) error

	CreateInBatch(model interface{}, data []interface{}) error
	Count(model interface{}, args ...interface{}) (int64, error)
	Cleanup() error
}

type MongoModel interface {
	CollectionName() string
}

type PersisterMongo struct {
	config  IPersisterConfig
	db      *mongo.Database
	dbMutex sync.Mutex
	client  *mongo.Client
	ctx     context.Context
}

func NewPersisterMongo(config IPersisterConfig) *PersisterMongo {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return &PersisterMongo{
		config: config,
		ctx:    ctx,
	}
}

func (pst *PersisterMongo) getConnectionString() (string, error) {
	cfg := pst.config

	return fmt.Sprintf("mongodb://%s:%s@%s:%s/",

		cfg.Username(),
		cfg.Password(),
		cfg.Host(),
		cfg.Port(),
	), nil
}

func (pst *PersisterMongo) getClient() (*mongo.Database, error) {
	if pst.db != nil {
		return pst.db, nil
	}

	pst.dbMutex.Lock()
	defer pst.dbMutex.Unlock()

	connectionStr, err := pst.getConnectionString()
	if err != nil {
		return nil, err
	}

	fmt.Println(connectionStr)

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionStr))
	if err != nil {
		return nil, err
	}

	pst.client = client

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	// defer client.Disconnect(ctx)

	// databases, err := client.ListDatabaseNames(ctx, bson.M{})
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(databases)

	db := client.Database(pst.config.DB())

	pst.db = db

	return db, nil
}

func (pst *PersisterMongo) getCollectionName(model interface{}) (string, error) {
	mongoModel, ok := model.(MongoModel)
	if ok {
		return mongoModel.CollectionName(), nil
	}
	return "", fmt.Errorf("model is not implement MongoModel")
}

func (pst *PersisterMongo) toDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}

func (pst *PersisterMongo) Count(model interface{}, args ...interface{}) (int64, error) {
	db, err := pst.getClient()
	if err != nil {
		return 0, err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return 0, err
	}

	count, err := db.Collection(collectionName).CountDocuments(pst.ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (pst *PersisterMongo) Find(model interface{}, filter interface{}) ([]interface{}, error) {
	db, err := pst.getClient()
	if err != nil {
		return nil, err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return nil, err
	}

	filterCursor, err := db.Collection(collectionName).Find(pst.ctx, filter)
	if err != nil {
		return nil, err
	}

	var results []interface{}
	if err = filterCursor.All(pst.ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (pst *PersisterMongo) FindOne(model interface{}, id string) (interface{}, error) {
	db, err := pst.getClient()
	if err != nil {
		return nil, err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return nil, err
	}

	idx, _ := primitive.ObjectIDFromHex(id)
	err = db.Collection(collectionName).FindOne(pst.ctx, bson.M{"_id": idx}).Decode(&model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (pst *PersisterMongo) Create(model interface{}) error {
	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	_, err = db.Collection(collectionName).InsertOne(pst.ctx, &model)

	if err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) CreateInBatch(model interface{}, data []interface{}) error {
	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	_, err = db.Collection(collectionName).InsertMany(pst.ctx, data)

	if err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) Update(model interface{}, id string) error {
	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	idx, _ := primitive.ObjectIDFromHex(id)

	updateDoc, err := pst.toDoc(model)
	if err != nil {
		return err
	}

	_, err = db.Collection(collectionName).UpdateOne(
		pst.ctx,
		bson.M{"_id": idx},
		bson.D{
			{"$set", updateDoc},
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) Cleanup() error {
	err := pst.client.Disconnect(pst.ctx)
	if err != nil {
		return err
	}
	return nil
}
