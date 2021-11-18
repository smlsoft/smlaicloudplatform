package microservice

import (
	"context"
	"fmt"
	"sync"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IPersister is interface for persister
type IPersisterMongo interface {
	Find(model interface{}, data interface{}, filter interface{}) error
	FindPage(model interface{}, data interface{}, limit int, page int, filter interface{}) (paginate.PaginationData, error)
	FindOne(model interface{}, filter interface{}) error
	FindByID(model interface{}, id string) error
	Create(model interface{}) (primitive.ObjectID, error)
	Update(model interface{}, id string) error
	CreateInBatch(model interface{}, data []interface{}) error
	Count(model interface{}, args ...interface{}) (int64, error)
	Exec(model interface{}) (*mongo.Collection, error)
	Delete(model interface{}, id string) error
	Cleanup() error
}

type MongoModel interface {
	CollectionName() string
}

type PersisterMongo struct {
	config    IPersisterConfig
	db        *mongo.Database
	dbMutex   sync.Mutex
	client    *mongo.Client
	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewPersisterMongo(config IPersisterConfig) *PersisterMongo {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	ctx := context.TODO()
	// defer cancel()

	return &PersisterMongo{
		config:    config,
		ctx:       ctx,
		ctxCancel: nil,
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

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionStr))
	if err != nil {
		return nil, err
	}

	pst.client = client

	err = client.Connect(pst.ctx)
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
	return "", fmt.Errorf("struct is not implement MongoModel")
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

func (pst *PersisterMongo) FindPage(model interface{}, data interface{}, limit int, page int, filter interface{}) (paginate.PaginationData, error) {
	db, err := pst.getClient()

	emptyPage := paginate.PaginationData{}

	if err != nil {
		return emptyPage, err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return emptyPage, err
	}

	var limit64 int64 = int64(limit)
	var page64 int64 = int64(page)

	paginatedData, err := paginate.New(db.Collection(collectionName)).Context(pst.ctx).Limit(limit64).Page(page64).Filter(filter).Decode(data).Find()
	if err != nil {
		return emptyPage, err
	}

	return paginatedData.Pagination, nil
}

func (pst *PersisterMongo) Find(model interface{}, data interface{}, filter interface{}) error {
	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	filterCursor, err := db.Collection(collectionName).Find(pst.ctx, filter)
	if err != nil {
		return err
	}

	if err = filterCursor.All(pst.ctx, data); err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) FindOne(model interface{}, filter interface{}) error {
	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	err = db.Collection(collectionName).FindOne(pst.ctx, filter).Decode(model)
	if err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) FindByID(model interface{}, id string) error {
	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	idx, _ := primitive.ObjectIDFromHex(id)
	err = db.Collection(collectionName).FindOne(pst.ctx, bson.M{"_id": idx}).Decode(model)
	if err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) Create(model interface{}) (primitive.ObjectID, error) {
	db, err := pst.getClient()
	if err != nil {
		return primitive.NilObjectID, err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return primitive.NilObjectID, err
	}

	result, err := db.Collection(collectionName).InsertOne(pst.ctx, &model)

	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
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
			{Key: "$set", Value: updateDoc},
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) Delete(model interface{}, id string) error {
	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	idx, _ := primitive.ObjectIDFromHex(id)
	_, err = db.Collection(collectionName).DeleteOne(pst.ctx, bson.M{"_id": idx})
	if err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) Exec(model interface{}) (*mongo.Collection, error) {
	db, err := pst.getClient()
	if err != nil {
		return nil, err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return nil, err
	}

	mongoCollection := db.Collection(collectionName)

	return mongoCollection, nil
}

func (pst *PersisterMongo) Cleanup() error {
	err := pst.client.Disconnect(pst.ctx)
	if err != nil {
		return err
	}

	pst.ctxCancel()

	return nil
}
