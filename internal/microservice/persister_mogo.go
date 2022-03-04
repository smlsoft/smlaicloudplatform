package microservice

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IPersister is interface for persister
type IPersisterMongo interface {
	Aggregate(model interface{}, pipeline []bson.D, decode interface{}, opts ...*options.AggregateOptions) error
	Find(model interface{}, filter interface{}, decode interface{}, opts ...*options.FindOptions) error
	FindPage(model interface{}, limit int, page int, filter interface{}, decode interface{}) (paginate.PaginationData, error)
	FindOne(model interface{}, filter interface{}, decode interface{}) error
	FindByID(model interface{}, keyName string, id interface{}, decode interface{}) error
	Create(model interface{}, data interface{}) (primitive.ObjectID, error)
	UpdateOne(model interface{}, keyName string, id interface{}, data interface{}) error
	Update(model interface{}, filter interface{}, data interface{}, opts ...*options.UpdateOptions) error
	CreateInBatch(model interface{}, data []interface{}) error
	Count(model interface{}, filter interface{}) (int, error)
	Exec(model interface{}) (*mongo.Collection, error)
	Delete(model interface{}, args ...interface{}) error
	DeleteByID(model interface{}, id string) error
	SoftDelete(model interface{}, filter interface{}) error
	SoftBatchDeleteByID(model interface{}, ids []string) error
	SoftDeleteByID(model interface{}, id string) error
	Cleanup() error
	TestConnect() error
	Healthcheck() error
}

// IPersisterConfig is interface for persister
type IPersisterMongoConfig interface {
	MongodbURI() string
	DB() string
}

type MongoModel interface {
	CollectionName() string
}

type PersisterMongo struct {
	config    IPersisterMongoConfig
	db        *mongo.Database
	dbMutex   sync.Mutex
	client    *mongo.Client
	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewPersisterMongo(config IPersisterMongoConfig) *PersisterMongo {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	ctx := context.Background()
	// defer cancel()

	return &PersisterMongo{
		config:    config,
		ctx:       ctx,
		ctxCancel: nil,
	}
}

func (pst *PersisterMongo) getConnectionString() (string, error) {
	cfg := pst.config

	return cfg.MongodbURI(), nil
}

func (pst *PersisterMongo) TestConnect() error {
	_, err := pst.getClient()

	if err != nil {
		return err
	}

	return nil
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

	// check connection
	err = client.Ping(context.TODO(), nil)
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

func (pst *PersisterMongo) Count(model interface{}, filter interface{}) (int, error) {
	db, err := pst.getClient()
	if err != nil {
		return 0, err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return 0, err
	}

	count, err := db.Collection(collectionName).CountDocuments(pst.ctx, filter)
	if err != nil {
		return 0, err
	}
	return int(count), nil

}

func (pst *PersisterMongo) FindPage(model interface{}, limit int, page int, filter interface{}, decode interface{}) (paginate.PaginationData, error) {
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

	paginatedData, err := paginate.New(db.Collection(collectionName)).Context(pst.ctx).Limit(limit64).Page(page64).Filter(filter).Decode(decode).Find()
	if err != nil {
		return emptyPage, err
	}

	return paginatedData.Pagination, nil
}

func (pst *PersisterMongo) Find(model interface{}, filter interface{}, decode interface{}, opts ...*options.FindOptions) error {
	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	filterCursor, err := db.Collection(collectionName).Find(pst.ctx, filter, opts...)
	if err != nil {
		return err
	}

	if err = filterCursor.All(pst.ctx, decode); err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) FindOne(model interface{}, filter interface{}, decode interface{}) error {
	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	err = db.Collection(collectionName).FindOne(context.TODO(), filter).Decode(decode)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return err
	}

	return nil
}

func (pst *PersisterMongo) FindByID(model interface{}, keyName string, id interface{}, decode interface{}) error {
	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	err = db.Collection(collectionName).FindOne(pst.ctx, bson.D{{Key: keyName, Value: id}}).Decode(decode)
	if err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) Create(model interface{}, data interface{}) (primitive.ObjectID, error) {
	db, err := pst.getClient()
	if err != nil {
		return primitive.NilObjectID, err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return primitive.NilObjectID, err
	}

	result, err := db.Collection(collectionName).InsertOne(pst.ctx, &data)

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

func (pst *PersisterMongo) Update(model interface{}, filter interface{}, data interface{}, opts ...*options.UpdateOptions) error {
	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	_, err = db.Collection(collectionName).UpdateMany(
		pst.ctx,
		filter,
		data,
		opts...,
	)

	if err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) UpdateOne(model interface{}, keyName string, id interface{}, data interface{}) error {
	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	// idx, _ := primitive.ObjectIDFromHex(id)

	updateDoc, err := pst.toDoc(data)
	if err != nil {
		return err
	}

	_, err = db.Collection(collectionName).UpdateOne(
		pst.ctx,
		bson.D{{
			Key:   keyName,
			Value: id,
		}},
		bson.D{
			{Key: "$set", Value: updateDoc},
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) SoftDeleteByID(model interface{}, id string) error {

	err := pst.UpdateOne(model, "guidFixed", id, map[string]bool{"deleted": true})

	if err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) SoftDelete(model interface{}, filter interface{}) error {

	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	_, err = db.Collection(collectionName).UpdateMany(pst.ctx, filter, bson.D{
		{Key: "$set", Value: bson.M{"deleted": true}},
	})

	if err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) SoftBatchDeleteByID(model interface{}, ids []string) error {
	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	objIDs := []primitive.ObjectID{}

	for _, id := range ids {
		idx, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}
		objIDs = append(objIDs, idx)
	}

	_, err = db.Collection(collectionName).UpdateMany(pst.ctx, bson.M{"_id": bson.M{"$in": objIDs}}, bson.D{
		{Key: "$set", Value: bson.M{"deleted": true}},
	})

	if err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) Delete(model interface{}, args ...interface{}) error {
	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	_, err = db.Collection(collectionName).DeleteMany(pst.ctx, args)
	if err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) DeleteByID(model interface{}, id string) error {
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

func (pst *PersisterMongo) Aggregate(model interface{}, pipeline []bson.D, decode interface{}, opts ...*options.AggregateOptions) error {
	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	// Value: bson.D{
	// 	bson.E{Key: "meta", Value: bson.E{Key: "$count", Value: "total"}},
	// }

	// query1 := bson.D{
	// 	{"meta", bson.A{
	// 		bson.E{Key: "$count", Value: "total"},
	// 	}},
	// }

	// pageFilter := []bson.D{
	// 	bson.D{
	// 		{"$facet", bson.D{
	// 			{"meta", bson.A{bson.D{{"$count", "total"}}}},
	// 			{"data", bson.A{bson.D{{"$limit", 2}}}},
	// 		}},
	// 	},
	// }

	//** var pipelinePage primitive.D
	// if len(pipeline) > 0 {
	// 	pageFilter = append(pipeline, pageFilter...)
	// }

	// fmt.Printf("%s\n\n", pageFilter)
	fmt.Printf("%s\n\n", pipeline[0])

	//facetStage := bson.D{{"$facet", query1}}

	filterCursor, err := db.Collection(collectionName).Aggregate(pst.ctx, mongo.Pipeline(pipeline), opts...)

	if err != nil {
		return err
	}

	var resultx []interface{}
	if err = filterCursor.All(pst.ctx, &resultx); err != nil {
		return err
	}

	rx, err := json.Marshal(resultx[0])
	if err != nil {
		return err
	}

	fmt.Printf("\n\n%s\n", rx)

	fmt.Printf("\n===result===\n%v\n\n", resultx)

	return nil
}

func (pst *PersisterMongo) AggregatePage(model interface{}, limit int, page int, filter interface{}, decode interface{}) (*paginate.PaginatedData, error) {
	db, err := pst.getClient()

	emptyPage := &paginate.PaginatedData{}

	if err != nil {
		return emptyPage, err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return emptyPage, err
	}

	var limit64 int64 = int64(limit)
	var page64 int64 = int64(page)

	paginatedData, err := paginate.New(db.Collection(collectionName)).Context(pst.ctx).Limit(limit64).Page(page64).Aggregate(filter)
	if err != nil {
		return emptyPage, err
	}

	return paginatedData, nil
}

func (pst *PersisterMongo) Cleanup() error {
	// err := pst.client.Disconnect(pst.ctx)
	// if err != nil {
	// 	return err
	// }

	// if pst != nil {
	// 	pst.ctxCancel()
	// }

	return nil
}

func (pst *PersisterMongo) Healthcheck() error {
	retry := 5
	// We will try to getClient 5 times
	for {
		if retry <= 0 {
			return fmt.Errorf("mongodb healthcheck failed")
		}
		retry--

		_, err := pst.getClient()
		if err != nil {
			// Healthcheck failed, wait 250ms then try again
			time.Sleep(250 * time.Millisecond)
			continue
		}
		return nil
	}
}
