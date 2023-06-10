package microservice

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/config"
)

// IPersister is interface for persister
type IPersisterMongo interface {
	Aggregate(model interface{}, pipeline interface{}, decode interface{}) error
	AggregatePage(model interface{}, pageable models.Pageable, criteria ...interface{}) (*mongopagination.PaginatedData, error)
	Find(model interface{}, filter interface{}, decode interface{}, opts ...*options.FindOptions) error
	FindPage(model interface{}, filter interface{}, pageable models.Pageable, decode interface{}) (mongopagination.PaginationData, error)
	FindSelectPage(model interface{}, selectFields interface{}, filter interface{}, pageable models.Pageable, decode interface{}) (mongopagination.PaginationData, error)
	// FindPage(model interface{}, limit int, page int, filter interface{}, decode interface{}) (mongopagination.PaginationData, error)
	// FindPageSort(model interface{}, limit int, page int, filter interface{}, sorts map[string]int, decode interface{}) (mongopagination.PaginationData, error)
	FindOne(model interface{}, filter interface{}, decode interface{}, opts ...*options.FindOneOptions) error
	FindByID(model interface{}, keyName string, id interface{}, decode interface{}) error
	Create(model interface{}, data interface{}) (primitive.ObjectID, error)
	UpdateOne(model interface{}, filterConditions map[string]interface{}, data interface{}) error
	Update(model interface{}, filter interface{}, data interface{}, opts ...*options.UpdateOptions) error
	CreateInBatch(model interface{}, data []interface{}) error
	Count(model interface{}, filter interface{}) (int, error)
	Exec(model interface{}) (*mongo.Collection, error)
	Delete(model interface{}, filter interface{}) error
	DeleteByID(model interface{}, id string) error
	SoftDelete(model interface{}, username string, filter interface{}) error
	SoftDeleteLastUpdate(model interface{}, username string, filter interface{}) error
	SoftBatchDeleteByID(model interface{}, username string, ids []string) error
	SoftDeleteByID(model interface{}, id string, username string) error
	Transaction(queryFunc func() error) error
	Cleanup() error
	TestConnect() error
	Healthcheck() error
	CreateIndex(model interface{}, indexName string, keys interface{}) (string, error)
}

type MongoModel interface {
	CollectionName() string
}

type PersisterMongo struct {
	config    config.IPersisterMongoConfig
	db        *mongo.Database
	dbMutex   sync.Mutex
	client    *mongo.Client
	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewPersisterMongo(config config.IPersisterMongoConfig) *PersisterMongo {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	ctx := context.Background()
	// defer cancel()

	return &PersisterMongo{
		config:    config,
		ctx:       ctx,
		ctxCancel: nil,
	}
}

func NewPersisterMongoWithDBContext(db *mongo.Database) *PersisterMongo {
	ctx := context.Background()
	return &PersisterMongo{
		db:  db,
		ctx: ctx,
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

	if pst.config.Debug() {
		cmdMonitor := &event.CommandMonitor{
			Started: func(_ context.Context, evt *event.CommandStartedEvent) {
				log.Print(evt.Command)
			},
		}

		pst.client, err = mongo.NewClient(options.Client().ApplyURI(connectionStr).SetMonitor(cmdMonitor))
		if err != nil {
			return nil, err
		}

	} else {
		pst.client, err = mongo.NewClient(options.Client().ApplyURI(connectionStr))
		if err != nil {
			return nil, err
		}
	}

	err = pst.client.Connect(pst.ctx)
	if err != nil {
		return nil, err
	}

	// check connection
	err = pst.client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	// defer client.Disconnect(ctx)

	// databases, err := client.ListDatabaseNames(ctx, bson.M{})
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(databases)

	db := pst.client.Database(pst.config.DB())

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

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return 0, err
	}
	return pst.PersisterCount(collectionName, filter)
}

func (pst *PersisterMongo) PersisterCount(collectionName string, filter interface{}) (int, error) {

	db, err := pst.getClient()
	if err != nil {
		return 0, err
	}

	count, err := db.Collection(collectionName).CountDocuments(pst.ctx, filter)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (pst *PersisterMongo) FindPage(model interface{}, filter interface{}, pageable models.Pageable, decode interface{}) (mongopagination.PaginationData, error) {
	db, err := pst.getClient()

	emptyPage := mongopagination.PaginationData{}

	if err != nil {
		return emptyPage, err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return emptyPage, err
	}

	var limit64 int64 = int64(pageable.Limit)
	var page64 int64 = int64(pageable.Page)

	pagingQuery := mongopagination.New(db.Collection(collectionName)).Context(pst.ctx).Limit(limit64).Page(page64).Filter(filter)

	for _, tempSort := range pageable.Sorts {
		tempSortVal := 1
		if tempSort.Value < 1 {
			tempSortVal = -1
		}
		pagingQuery = pagingQuery.Sort(tempSort.Key, tempSortVal)
	}

	paginatedData, err := pagingQuery.Decode(decode).Find()
	if err != nil {
		return emptyPage, err
	}

	return paginatedData.Pagination, nil
}

func (pst *PersisterMongo) FindSelectPage(model interface{}, selectFields interface{}, filter interface{}, pageable models.Pageable, decode interface{}) (mongopagination.PaginationData, error) {
	db, err := pst.getClient()

	emptyPage := mongopagination.PaginationData{}

	if err != nil {
		return emptyPage, err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return emptyPage, err
	}

	var limit64 int64 = int64(pageable.Limit)
	var page64 int64 = int64(pageable.Page)

	pagingQuery := mongopagination.New(db.Collection(collectionName)).Context(pst.ctx).Select(selectFields).Limit(limit64).Page(page64).Filter(filter)

	for _, tempSort := range pageable.Sorts {
		tempSortVal := 1
		if tempSort.Value < 1 {
			tempSortVal = -1
		}
		pagingQuery = pagingQuery.Sort(tempSort.Key, tempSortVal)
	}

	paginatedData, err := pagingQuery.Decode(decode).Find()
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

func (pst *PersisterMongo) FindOne(model interface{}, filter interface{}, decode interface{}, opts ...*options.FindOneOptions) error {
	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	result := db.Collection(collectionName).FindOne(context.TODO(), filter, opts...)

	err = result.Decode(decode)
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

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return pst.PersisterCreate(collectionName, data)
}

func (pst *PersisterMongo) PersisterCreate(collectionName string, data interface{}) (primitive.ObjectID, error) {

	db, err := pst.getClient()
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

func (pst *PersisterMongo) UpdateOne(model interface{}, filterConditions map[string]interface{}, data interface{}) error {
	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	updateDoc, err := pst.toDoc(data)
	if err != nil {
		return err
	}

	filterDoc := bson.M{}

	for key, val := range filterConditions {
		filterDoc[key] = val
	}

	_, err = db.Collection(collectionName).UpdateOne(
		pst.ctx,
		filterDoc,
		bson.D{
			{Key: "$set", Value: updateDoc},
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) SoftDeleteByID(model interface{}, id string, username string) error {
	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	deletedAt := time.Now()
	// _, err := pst.UpdateOne(model, "guidfixed", id, map[string]interface{}{"deletedat": deletedAt})

	_, err = db.Collection(collectionName).UpdateOne(
		pst.ctx,
		bson.D{{
			Key:   "guidfixed",
			Value: id,
		}},
		bson.D{
			{Key: "$set",
				Value: bson.D{
					{Key: "deletedby", Value: username},
					{Key: "deletedat", Value: deletedAt},
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) SoftDelete(model interface{}, username string, filter interface{}) error {

	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	deletedAt := time.Now()

	_, err = db.Collection(collectionName).UpdateMany(pst.ctx, filter, bson.D{
		{Key: "$set", Value: bson.M{"deletedat": deletedAt, "deletedby": username}},
	})

	if err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) SoftDeleteLastUpdate(model interface{}, username string, filter interface{}) error {

	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	deletedAt := time.Now()

	_, err = db.Collection(collectionName).UpdateMany(pst.ctx, filter, bson.D{
		{Key: "$set", Value: bson.M{"deletedat": deletedAt, "deletedby": username, "lastupdatedat": deletedAt}},
	})

	if err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) SoftBatchDeleteByID(model interface{}, username string, ids []string) error {
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

	deletedAt := time.Now()

	_, err = db.Collection(collectionName).UpdateMany(pst.ctx,
		bson.M{"_id": bson.M{"$in": objIDs}},
		bson.D{
			{Key: "$set", Value: bson.M{"deletedat": deletedAt, "deletedby": username}},
		})

	if err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) Delete(model interface{}, filter interface{}) error {
	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	_, err = db.Collection(collectionName).DeleteMany(pst.ctx, filter)
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

func (pst *PersisterMongo) Aggregate(model interface{}, pipeline interface{}, decode interface{}) error {
	db, err := pst.getClient()
	if err != nil {
		return err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return err
	}

	// var aggregationFilter []bson.M
	// for _, filter := range filters {
	// 	aggregationFilter = append(aggregationFilter, filter.(bson.M))
	// }

	filterCursor, err := db.Collection(collectionName).Aggregate(pst.ctx, pipeline)

	if err != nil {
		return err
	}

	if err = filterCursor.All(pst.ctx, decode); err != nil {
		return err
	}

	return nil
}

func (pst *PersisterMongo) AggregatePage(model interface{}, pageable models.Pageable, criteria ...interface{}) (*mongopagination.PaginatedData, error) {
	db, err := pst.getClient()

	emptyPage := &mongopagination.PaginatedData{}

	if err != nil {
		return emptyPage, err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return emptyPage, err
	}

	var page64 int64 = int64(pageable.Page)
	var limit64 int64 = int64(pageable.Limit)

	paginatedData, err := mongopagination.New(db.Collection(collectionName)).Context(pst.ctx).Limit(limit64).Page(page64).Aggregate(criteria...)
	if err != nil {
		return emptyPage, err
	}

	return paginatedData, nil
}

func (pst *PersisterMongo) Cleanup() error {
	err := pst.client.Disconnect(pst.ctx)
	if err != nil {
		return err
	}

	if pst != nil {
		pst.ctxCancel()
	}

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

func (pst *PersisterMongo) Transaction(queryFunc func() error) error {
	pst.getClient()
	client := pst.client

	session, err := client.StartSession()
	if err != nil {
		return err
	}

	if err := session.StartTransaction(); err != nil {
		return err
	}

	if err := mongo.WithSession(pst.ctx, session, func(sc mongo.SessionContext) error {
		err := queryFunc()

		if err != nil {
			return err
		}

		if err = session.CommitTransaction(sc); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	session.EndSession(pst.ctx)

	return nil
}

func (pst *PersisterMongo) CreateIndex(model interface{}, indexName string, keys interface{}) (string, error) {
	db, err := pst.getClient()
	if err != nil {
		return "", err
	}

	collectionName, err := pst.getCollectionName(model)
	if err != nil {
		return "", err
	}

	indexModel := mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetUnique(true).SetName(indexName),
	}

	mongoCollection := db.Collection(collectionName)

	resultIndexName, err := mongoCollection.Indexes().CreateOne(pst.ctx, indexModel)

	if err != nil {
		return "", err
	}

	return resultIndexName, nil
}
