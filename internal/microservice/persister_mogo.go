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
	Aggregate(model interface{}, pipeline interface{}, decode interface{}, opts ...*options.AggregateOptions) error
	Find(model interface{}, filter interface{}, decode interface{}, opts ...*options.FindOptions) error
	FindPage(model interface{}, limit int, page int, filter interface{}, decode interface{}) (paginate.PaginationData, error)
	FindOne(model interface{}, filter interface{}, decode interface{}) error
	FindByID(model interface{}, keyName string, id interface{}, decode interface{}) error
	Create(model interface{}, data interface{}) (primitive.ObjectID, error)
	Update(model interface{}, data interface{}, keyName string, id interface{}) error
	CreateInBatch(model interface{}, data []interface{}) error
	Count(model interface{}, filter interface{}) (int, error)
	Exec(model interface{}) (*mongo.Collection, error)
	Delete(model interface{}, args ...interface{}) error
	DeleteByID(model interface{}, id string) error
	SoftDelete(model interface{}, ids []string) error
	SoftDeleteByID(model interface{}, id string) error
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
	if err != nil {
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

func (pst *PersisterMongo) Update(model interface{}, data interface{}, keyName string, id interface{}) error {
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

	err := pst.Update(model, map[string]bool{"deleted": true}, "guidFixed", id)

	if err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) SoftDelete(model interface{}, ids []string) error {
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

func (pst *PersisterMongo) Aggregate(model interface{}, pipeline interface{}, decode interface{}, opts ...*options.AggregateOptions) error {
	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	filterCursor, err := db.Collection(collectionName).Aggregate(pst.ctx, pipeline, opts...)

	if err != nil {
		return err
	}

	if err = filterCursor.All(pst.ctx, decode); err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) AggregatePage(model interface{}, limit int, page int, pipeline interface{}, decode interface{}, opts ...*options.AggregateOptions) error {
	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	mongoPipeline := pipeline.(mongo.Pipeline)

	mongoPipeline = append(mongoPipeline, bson.D{{Key: "$limit", Value: 3}})

	filterCursor, err := db.Collection(collectionName).Aggregate(pst.ctx, mongoPipeline, opts...)

	if err != nil {
		return err
	}

	if err = filterCursor.All(pst.ctx, decode); err != nil {
		return err
	}

	return nil
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
