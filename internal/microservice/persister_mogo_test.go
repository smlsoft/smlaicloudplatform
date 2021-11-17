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

	objID, err := pst.Create(&Product{
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

	products := []Product{}
	err := pst.Find(&Product{}, &products, bson.M{})

	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(products) < 1 {
		t.Error("Find not found item")
	}

}

func TestMongodbFindOne(t *testing.T) {
	cfg := &ConfigDBTest{}

	pst := NewPersisterMongo(cfg)

	product := &Product{}

	err := pst.FindOne(product, bson.M{"product_code": productCode})

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

	err := pst.FindOne(productFind, bson.M{"product_code": productCode})

	if err != nil {
		t.Error(err.Error())
		return
	}

	product := &Product{}
	err = pst.FindByID(product, productFind.ID.Hex())

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
	err := pst.FindOne(productFind, bson.M{"product_code": productCode})

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
	err = pst.FindOne(productCheck, bson.M{"product_code": productCode})

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

	products := []interface{}{
		Product{
			ProductCode: "pdt-01",
			ProductName: "pdt name 01",
		},
		Product{
			ProductCode: "pdt-02",
			ProductName: "pdt name 02",
		},
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

	err := pst.FindOne(productFind, bson.M{"product_code": productCode})

	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Log(productFind)
	t.Log(productFind.ID.Hex())

	product := &Product{}
	err = pst.Delete(product, productFind.ID.Hex())

	if err != nil {
		t.Error(err.Error())
		return
	}

}
