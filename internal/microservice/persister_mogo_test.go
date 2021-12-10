package microservice

import (
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConfigDBTest struct{}

func (c *ConfigDBTest) Host() string {
	return "localhost"
}

func (c *ConfigDBTest) DB() string {
	return "micro_test"
}

func (c *ConfigDBTest) Port() string {
	return "27017"
}

func (c *ConfigDBTest) Username() string {
	return "root"
}
func (c *ConfigDBTest) Password() string {
	return "rootx"
}

func (c *ConfigDBTest) SSLMode() string {
	return ""
}
func (c *ConfigDBTest) TimeZone() string {
	return ""
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
	productCode = "pdtx01"
	productName = "product name 01"
)

func TestMongodbCount(t *testing.T) {

	cfg := &ConfigDBTest{}

	pst := NewPersisterMongo(cfg)

	count, err := pst.Count(&Product{})
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

	pst := NewPersisterMongo(cfg)

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

	pst := NewPersisterMongo(cfg)

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

}

func TestMongodbFindPage(t *testing.T) {
	cfg := &ConfigDBTest{}

	pst := NewPersisterMongo(cfg)

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

	pst := NewPersisterMongo(cfg)

	product := &Product{}

	err := pst.FindOne(&Product{}, bson.M{"product_code": productCode}, product)

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

	pst := NewPersisterMongo(cfg)

	productFind := &Product{}

	err := pst.FindOne(&Product{}, bson.M{"product_code": productCode}, productFind)

	if err != nil {
		t.Error(err.Error())
		return
	}

	product := &Product{}
	err = pst.FindByID(&Product{}, productFind.ID.Hex(), product)

	if err != nil {
		t.Error(err.Error())
		return
	}

	if product.ProductCode != productFind.ProductCode {
		t.Error("Product code not match")
		return
	}
}

func TestMongodbUpdate(t *testing.T) {

	productNameModified := "product name modify"

	cfg := &ConfigDBTest{}

	pst := NewPersisterMongo(cfg)

	productFind := &Product{}
	err := pst.FindOne(&Product{}, bson.M{"product_code": productCode}, productFind)

	if err != nil {
		t.Error(err.Error())
	}

	err = pst.Update(&Product{
		ProductCode: productCode,
		ProductName: productNameModified,
	}, productFind.ID.Hex())

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

	pst := NewPersisterMongo(cfg)

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

	pst := NewPersisterMongo(cfg)

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
