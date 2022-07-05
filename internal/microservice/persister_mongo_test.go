package microservice_test

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"smlcloudplatform/internal/microservice"
	"testing"

	"github.com/tj/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// type ConfigDBTest struct{}

// func (c *ConfigDBTest) MongodbURI() string {
// 	return "mongodb://192.168.2.209:27017/"
// }

// func (c *ConfigDBTest) DB() string {
// 	return "micro_test"
// }

type Product struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProductCode string             `json:"product_code" bson:"product_code"`
	ProductName string             `json:"product_name" bson:"product_name"`
}

func (pdt *Product) CollectionName() string {
	return "product"
}

var (
	productCode  = "pdt-01"
	productName  = "product name 01"
	noClientOpts = mtest.NewOptions().CreateClient(false)
)

type negateCodec struct {
	ID int64 `bson:"_id"`
}

func (e *negateCodec) EncodeValue(ectx bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	return vw.WriteInt64(val.Int())
}

// DecodeValue negates the value of ID when reading
func (e *negateCodec) DecodeValue(ectx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	i, err := vr.ReadInt64()
	if err != nil {
		return err
	}

	val.SetInt(i * -1)
	return nil
}

// init

func initCollection(mt *mtest.T, coll *mongo.Collection) {
	mt.Helper()

	var docs []interface{}
	for i := 1; i <= 5; i++ {
		docs = append(docs, bson.D{{"product_code", int32(i)}})
	}

	_, err := coll.InsertMany(context.Background(), docs)
	assert.Nil(mt, err, "InsertMany error for initial data: %v", err)
}

func TestMongoPersisterCountWithmtest(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("Success", func(mt *mtest.T) {

		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"n", 1},
		})

		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)

		mt.AddMockResponses(first, killCursors)

		pst := microservice.NewPersisterMongoWithDBContext(mt.DB)
		count, err := pst.Count(&Product{}, bson.M{"product_code": "0001"})
		if err != nil {
			t.Error(err.Error())
			return
		}

		fmt.Printf("Count :: %d \n", count)

		if count < 0 {
			t.Error("can't get Count ")
			return
		}
	})
}

func TestMongodbCreate(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("Success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		pst := microservice.NewPersisterMongoWithDBContext(mt.DB)
		objID, err := pst.Create(&Product{}, &Product{
			ProductCode: productCode,
			ProductName: productName,
		})

		t.Log(objID)

		if err != nil {
			t.Error(err.Error())
		}
	})
}

func TestMongodbCreateInBatchWithMtest(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("Success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		pst := microservice.NewPersisterMongoWithDBContext(mt.DB)
		products := make([]interface{}, 0)

		for i := 1; i <= 25; i++ {
			products = append(products, Product{
				ProductCode: fmt.Sprintf("pdt-%02d", i),
				ProductName: fmt.Sprintf("pdt name %d", i),
			})
		}

		err := pst.CreateInBatch(&Product{}, products)

		if err != nil {
			t.Error(err.Error())
		}
	})
}

func TestMongoDBFindMockData(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("Success", func(mt *mtest.T) {

		id1 := primitive.NewObjectID()
		id2 := primitive.NewObjectID()

		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", id1},
			{"product_code", "0001"},
			{"product_name", "name 0001"},
		})

		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{"_id", id2},
			{"product_code", "0002"},
			{"product_name", "name 0002"},
		})

		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)

		mt.AddMockResponses(first, second, killCursors)

		pst := microservice.NewPersisterMongoWithDBContext(mt.DB)
		products := []Product{}
		err := pst.Find(&Product{}, bson.M{}, &products)

		if err != nil {
			t.Error(err.Error())
			return
		}

		if len(products) < 1 {
			t.Error("Find not found item")
		}

		t.Log("product all :: ", len(products))
	})
}

func TestMongodbFindPage(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	// // registryOpts := options.Client().
	// // 	SetRegistry(bson.NewRegistryBuilder().RegisterCodec(reflect.TypeOf(int64(0)), &negateCodec{}).Build())

	// mt.RunOpts("registry passed to cursors", mtest.NewOptions().CreateClient(false), func(mt *mtest.T) {
	// 	_, err := mt.Coll.InsertOne(context.Background(), negateCodec{ID: 10})
	// 	assert.Nil(mt, err, "InsertOne error: %v", err)
	// 	var got negateCodec
	// 	err = mt.Coll.FindOne(context.Background(), bson.D{}).Decode(&got)
	// 	assert.Nil(mt, err, "Find error: %v", err)

	// 	assert.Equal(mt, int64(-10), got.ID, "expected ID -10, got %v", got.ID)
	// })

	// mt.RunOpts("try find", noClientOpts, func(mt *mtest.T) {
	if os.Getenv("SERVERLESS") == "serverless" {
		mt.Skip("skipping as serverless forbids capped collections")
	}

	mt.Run("Success", func(mt *mtest.T) {

		//initCollection(mt, mt.Coll)
		id1 := primitive.NewObjectID()
		id2 := primitive.NewObjectID()

		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", id1},
			{"product_code", "0001"},
			{"product_name", "name 0001"},
		})

		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{"_id", id2},
			{"product_code", "0002"},
			{"product_name", "name 0002"},
		})

		//killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)

		mt.AddMockResponses(first, second)
		pst := microservice.NewPersisterMongoWithDBContext(mt.DB)

		products := []Product{}
		pagination, err := pst.FindPage(&Product{}, 1, 1, bson.M{}, &products)

		if err != nil {
			t.Error(err.Error())
			return
		}

		fmt.Println(pagination)
		fmt.Println(products)

		if len(products) < 1 {
			t.Error("Find not found item")
		}
	})
	// })
}

// func TestMongodbFindOne(t *testing.T) {
// 	cfg := &ConfigDBTest{}

// 	pst := microservice.NewPersisterMongo(cfg)

// 	product := &Product{}
// 	err := pst.FindOne(&Product{}, bson.D{{Key: "product_code", Value: productCode}}, product)

// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}

// 	if product.ProductCode != productCode {
// 		t.Error("Product code not match")
// 		return
// 	}
// }

// func TestMongodbFindByID(t *testing.T) {
// 	cfg := &ConfigDBTest{}

// 	pst := microservice.NewPersisterMongo(cfg)

// 	productFind := &Product{}

// 	err := pst.FindOne(&Product{}, bson.M{"product_code": productCode}, productFind)

// 	t.Log(productFind.ID.Hex())
// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}

// 	product := &Product{}
// 	err = pst.FindByID(&Product{}, "_id", productFind.ID, product)

// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}

// 	if product.ProductCode != productFind.ProductCode {
// 		t.Error("Product code not match")
// 		return
// 	}
// }

// func TestMongodbUpdateOne(t *testing.T) {

// 	productNameModified := "product name modifyx"

// 	cfg := &ConfigDBTest{}

// 	pst := microservice.NewPersisterMongo(cfg)

// 	productFind := &Product{}
// 	err := pst.FindOne(&Product{}, bson.M{"product_code": productCode}, productFind)

// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	err = pst.UpdateOne(
// 		&Product{},
// 		map[string]interface{}{
// 			"_id":          productFind.ID,
// 			"product_code": productFind.ProductCode,
// 		}, &Product{
// 			ProductCode: productCode,
// 			ProductName: productNameModified,
// 		})

// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	productCheck := &Product{}
// 	err = pst.FindOne(&Product{}, bson.M{"product_code": productCode}, productCheck)

// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	if productCheck.ProductName != productNameModified {
// 		t.Error("Product not modified")
// 	}
// }

// func TestMongodbUpdate(t *testing.T) {

// 	productNameModified := "test modify"

// 	cfg := &ConfigDBTest{}

// 	pst := microservice.NewPersisterMongo(cfg)

// 	productFind := &Product{}
// 	err := pst.FindOne(&Product{}, bson.M{"product_code": productCode}, productFind)

// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	err = pst.Update(&Product{}, bson.M{"_id": productFind.ID}, bson.M{"$set": &Product{
// 		ProductCode: productCode,
// 		ProductName: productNameModified,
// 	}})

// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	productCheck := &Product{}
// 	err = pst.FindOne(&Product{}, bson.M{"product_code": productCode}, productCheck)

// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	if productCheck.ProductName != productNameModified {
// 		t.Error("Product not modified")
// 	}

// }

// func TestMongodbDelete(t *testing.T) {
// 	cfg := &ConfigDBTest{}

// 	pst := microservice.NewPersisterMongo(cfg)

// 	productFind := &Product{}

// 	err := pst.FindOne(&Product{}, bson.M{"product_code": productCode}, productFind)

// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}

// 	t.Log(productFind)
// 	t.Log(productFind.ID.Hex())

// 	product := &Product{}
// 	err = pst.DeleteByID(product, productFind.ID.Hex())

// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}

// }

// func TestMongodbSoftDeleteByID(t *testing.T) {
// 	cfg := &ConfigDBTest{}

// 	pst := microservice.NewPersisterMongo(cfg)

// 	product := &Product{}
// 	err := pst.DeleteByID(product, "6195af880e33cec3af136720")

// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}

// }

// func TestMongodbSoftBatchDelete(t *testing.T) {
// 	cfg := &ConfigDBTest{}

// 	pst := microservice.NewPersisterMongo(cfg)

// 	err := pst.SoftBatchDeleteByID(&Product{}, "", []string{"6195af880e33cec3af136724", "6195af880e33cec3af136725"})

// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}

// }

// func TestMongodbSoftDelete(t *testing.T) {
// 	cfg := &ConfigDBTest{}

// 	pst := microservice.NewPersisterMongo(cfg)

// 	err := pst.SoftDeleteByID(&Product{}, "x123x", "dev")

// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}

// }

// func TestMongodbAggregate(t *testing.T) {
// 	cfg := &ConfigDBTest{}

// 	pst := microservice.NewPersisterMongo(cfg)

// 	// pipeline := mongo.Pipeline{}
// 	products := []Product{}

// 	// query1 := bson.A{bson.D{{"$match", bson.D{{"product_code", "pdt-02"}}}}, bson.D{{"$count", "count"}}}
// 	// query2 := bson.A{bson.D{{"$match", bson.D{{"product_code", "pdt-03"}}}}, bson.D{{"$count", "count"}}}
// 	// query3 := bson.A{bson.D{{"$count", "total"}}}

// 	// facetStage := bson.D{{"$facet", bson.D{{"query1", query1}, {"query2", query2}, {"query3", query3}}}}

// 	err := pst.Aggregate(&Product{}, []bson.D{
// 		bson.D{{"$match", bson.M{"product_code": "pdt-01"}}},
// 	}, &products)

// 	if err != nil {
// 		fmt.Println("=====[Error]======")
// 		fmt.Println(err.Error())
// 	}

// 	t.Log("count :: ", products)
// }

// func TestMongodbAggregatePage(t *testing.T) {
// 	cfg := &ConfigDBTest{}

// 	pst := microservice.NewPersisterMongo(cfg)

// 	aggPaginatedData, err := pst.AggregatePage(&Product{}, 2, 0, bson.M{"$match": bson.M{"product_code": "pdt-01"}})

// 	products := []Product{}
// 	// var aggProductList []Product
// 	for _, raw := range aggPaginatedData.Data {
// 		var product *Product

// 		if marshallErr := bson.Unmarshal(raw, &product); marshallErr == nil {
// 			products = append(products, *product)
// 		}

// 	}

// 	if err != nil {
// 		fmt.Println("=====[Error]======")
// 		fmt.Println(err.Error())
// 	}
// 	t.Log("count :: ", products)
// }

// func TestMongodbFindPageX(t *testing.T) {
// 	cfg := &ConfigDBTest{}

// 	pst := microservice.NewPersisterMongo(cfg)

// 	products := []Product{}
// 	pagination, err := pst.FindPageSort(&Product{}, 5, 1, bson.M{}, map[string]int{
// 		"product_code": 5,
// 	}, &products)

// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}

// 	fmt.Println(pagination)
// 	fmt.Println(products)

// 	if len(products) < 1 {
// 		t.Error("Find not found item")
// 	}
// }
