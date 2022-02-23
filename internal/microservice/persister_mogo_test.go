package microservice_test

import (
	"errors"
	"fmt"
	"smlcloudplatform/internal/microservice"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConfigDBTest struct{}

func (c *ConfigDBTest) MongodbURI() string {
	return "mongodb://root:rootx@localhost:27017/"
}

func (c *ConfigDBTest) DB() string {
	return "micro_test"
}

type Product struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProductCode string             `json:"product_code" bson:"product_code"`
	ProductName string             `json:"product_name" bson:"product_name"`
}

func (pdt *Product) CollectionName() string {
	return "product"
}

var (
	productCode = "pdt-01"
	productName = "product name 01"
)

func TestMongodbCount(t *testing.T) {

	cfg := &ConfigDBTest{}

	pst := microservice.NewPersisterMongo(cfg)

	count, err := pst.Count(&Product{}, bson.M{"product_code": "pdt-02"})
	if err != nil {
		t.Error(err.Error())
		return
	}

	fmt.Printf("Count :: %d \n", count)

	if count < 0 {
		t.Error("can't get Count ")
		return
	}

}

func TestMongodbCreate(t *testing.T) {

	cfg := &ConfigDBTest{}

	pst := microservice.NewPersisterMongo(cfg)

	objID, err := pst.Create(&Product{}, &Product{
		ProductCode: productCode,
		ProductName: productName,
	})

	t.Log(objID)

	if err != nil {
		t.Error(err.Error())
	}
}

func TestMongodbFind(t *testing.T) {
	cfg := &ConfigDBTest{}

	pst := microservice.NewPersisterMongo(cfg)

	// products := []Product{}
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

}

func TestMongodbFindPage(t *testing.T) {
	cfg := &ConfigDBTest{}

	pst := microservice.NewPersisterMongo(cfg)

	products := []Product{}
	pagination, err := pst.FindPage(&Product{}, 5, 1, bson.M{}, &products)

	if err != nil {
		t.Error(err.Error())
		return
	}

	fmt.Println(pagination)
	fmt.Println(products)

	if len(products) < 1 {
		t.Error("Find not found item")
	}

}

func TestMongodbFindOne(t *testing.T) {
	cfg := &ConfigDBTest{}

	pst := microservice.NewPersisterMongo(cfg)

	product := &Product{}
	err := pst.FindOne(&Product{}, bson.D{{Key: "product_code", Value: productCode}}, product)

	if err != nil {
		t.Error(err.Error())
		return
	}

	if product.ProductCode != productCode {
		t.Error("Product code not match")
		return
	}
}

func TestMongodbFindByID(t *testing.T) {
	cfg := &ConfigDBTest{}

	pst := microservice.NewPersisterMongo(cfg)

	productFind := &Product{}

	err := pst.FindOne(&Product{}, bson.M{"product_code": productCode}, productFind)

	t.Log(productFind.ID.Hex())
	if err != nil {
		t.Error(err.Error())
		return
	}

	product := &Product{}
	err = pst.FindByID(&Product{}, "_id", productFind.ID, product)

	if err != nil {
		t.Error(err.Error())
		return
	}

	if product.ProductCode != productFind.ProductCode {
		t.Error("Product code not match")
		return
	}
}

func TestMongodbUpdateOne(t *testing.T) {

	productNameModified := "product name modify"

	cfg := &ConfigDBTest{}

	pst := microservice.NewPersisterMongo(cfg)

	productFind := &Product{}
	err := pst.FindOne(&Product{}, bson.M{"product_code": productCode}, productFind)

	if err != nil {
		t.Error(err.Error())
	}

	err = pst.UpdateOne(&Product{}, "_id", productFind.ID, &Product{
		ProductCode: productCode,
		ProductName: productNameModified,
	})

	if err != nil {
		t.Error(err.Error())
	}

	productCheck := &Product{}
	err = pst.FindOne(&Product{}, bson.M{"product_code": productCode}, productCheck)

	if err != nil {
		t.Error(err.Error())
	}

	if productCheck.ProductName != productNameModified {
		t.Error("Product not modified")
	}
}

func TestMongodbUpdate(t *testing.T) {

	productNameModified := "test modify"

	cfg := &ConfigDBTest{}

	pst := microservice.NewPersisterMongo(cfg)

	productFind := &Product{}
	err := pst.FindOne(&Product{}, bson.M{"product_code": productCode}, productFind)

	if err != nil {
		t.Error(err.Error())
	}

	err = pst.Update(&Product{}, bson.M{"_id": productFind.ID}, bson.M{"$set": &Product{
		ProductCode: productCode,
		ProductName: productNameModified,
	}})

	if err != nil {
		t.Error(err.Error())
	}

	productCheck := &Product{}
	err = pst.FindOne(&Product{}, bson.M{"product_code": productCode}, productCheck)

	if err != nil {
		t.Error(err.Error())
	}

	if productCheck.ProductName != productNameModified {
		t.Error("Product not modified")
	}

}

func TestMongodbCreateInBatch(t *testing.T) {
	cfg := &ConfigDBTest{}

	pst := microservice.NewPersisterMongo(cfg)

	products := make([]interface{}, 0)

	for i := 1; i <= 25; i++ {
		products = append(products, Product{
			ProductCode: fmt.Sprintf("pdt-%02d", i),
			ProductName: fmt.Sprintf("pdt name %d", i),
		})
	}

	// x := make([]interface{}, len(products))

	// for i, v := range products {
	// 	x[i] = v
	// }
	err := pst.CreateInBatch(&Product{}, products)

	if err != nil {
		t.Error(err.Error())
	}
}

func TestMongodbDelete(t *testing.T) {
	cfg := &ConfigDBTest{}

	pst := microservice.NewPersisterMongo(cfg)

	productFind := &Product{}

	err := pst.FindOne(&Product{}, bson.M{"product_code": productCode}, productFind)

	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Log(productFind)
	t.Log(productFind.ID.Hex())

	product := &Product{}
	err = pst.DeleteByID(product, productFind.ID.Hex())

	if err != nil {
		t.Error(err.Error())
		return
	}

}

func TestMongodbSoftDeleteByID(t *testing.T) {
	cfg := &ConfigDBTest{}

	pst := microservice.NewPersisterMongo(cfg)

	product := &Product{}
	err := pst.DeleteByID(product, "6195af880e33cec3af136720")

	if err != nil {
		t.Error(err.Error())
		return
	}

}

func TestMongodbSoftDelete(t *testing.T) {
	cfg := &ConfigDBTest{}

	pst := microservice.NewPersisterMongo(cfg)

	err := pst.SoftBatchDeleteByID(&Product{}, []string{"6195af880e33cec3af136724", "6195af880e33cec3af136725"})

	if err != nil {
		t.Error(err.Error())
		return
	}

}

func TestMongodbAggregate(t *testing.T) {
	cfg := &ConfigDBTest{}

	pst := microservice.NewPersisterMongo(cfg)

	// pipeline := mongo.Pipeline{}
	products := []Product{}

	// query1 := bson.A{bson.D{{"$match", bson.D{{"product_code", "pdt-02"}}}}, bson.D{{"$count", "count"}}}
	// query2 := bson.A{bson.D{{"$match", bson.D{{"product_code", "pdt-03"}}}}, bson.D{{"$count", "count"}}}
	// query3 := bson.A{bson.D{{"$count", "total"}}}

	// facetStage := bson.D{{"$facet", bson.D{{"query1", query1}, {"query2", query2}, {"query3", query3}}}}

	err := pst.Aggregate(&Product{}, []bson.D{
		bson.D{{"$match", bson.M{"product_code": "pdt-01"}}},
	}, &products)

	if err != nil {
		fmt.Println("=====[Error]======")
		fmt.Println(err.Error())
	}

	t.Log("count :: ", products)
}

func TestMongodbAggregatePage(t *testing.T) {
	cfg := &ConfigDBTest{}

	pst := microservice.NewPersisterMongo(cfg)

	products := []Product{}

	aggPaginatedData, err := pst.AggregatePage(&Product{}, 2, 0, bson.M{"$match": bson.M{"product_code": "pdt-01"}}, &products)

	// var aggProductList []Product
	for _, raw := range aggPaginatedData.Data {
		var product *Product

		if marshallErr := bson.Unmarshal(raw, &product); marshallErr == nil {
			products = append(products, *product)
		}

	}

	if err != nil {
		fmt.Println("=====[Error]======")
		fmt.Println(err.Error())
	}
	t.Log("count :: ", products)
}

func TestConnectMongodbUrit(t *testing.T) {

	mongoPersisterConfig := microservice.NewMongoPersisterConfig()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	err := mongoPersister.TestConnect()
	if err != nil {
		t.Error(errors.New("Cannot connect Database"))
	}
}
