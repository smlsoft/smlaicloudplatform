package usecases_test

import (
	"encoding/json"
	"smlaicloudplatform/internal/product/productbarcode/models"
	"smlaicloudplatform/internal/product/productbarcode/usecases"
	"testing"

	"github.com/stretchr/testify/assert"
)

var jsonStr = `{"id":"000000000000000000000000","shopid":"2QoilMQkX9i6vtAE88ilEubnrhz","guidfixed":"2QotSoH9vZZ1LW7loBAZbuCt1pW","itemcode":"","barcode":"KSC-001","groupcode":"","groupnames":[],"names":[{"code":"th","name":"SIZE S ปีกกลาง 4 ชิ้น","isauto":false,"isdelete":false},{"code":"en","name":"CHICKEN WINGS SIZE S (4 PCS.)","isauto":false,"isdelete":false},{"code":"ko","name":"미들 윙 4 개","isauto":false,"isdelete":false}],"xsorts":[],"itemunitcode":"PLATE","itemunitnames":[{"code":"th","name":"จาน","isauto":false,"isdelete":false}],"prices":[{"keynumber":1,"price":99},{"keynumber":2,"price":0}],"imageuri":"","options":[{"guid":"a6270c74-5a86-488a-951b-c63d651e52bb","choicetype":0,"maxselect":2,"minselect":1,"names":[{"code":"th","name":"โซลมายด์ ชิกเก้นท์ *เลือกได้ 2 ซอส","isauto":false,"isdelete":false},{"code":"en","name":"SEOULMIND CHICKEN ","isauto":false,"isdelete":false},{"code":"ko","name":"서울 마인드 치킨(2가지 소스 선택할 수 있음)","isauto":false,"isdelete":false}],"choices":[{"guid":"0e3e9cd9-fa5c-496e-b391-f95287b54393","refbarcode":"","refproductcode":"","refunitcode":"","names":[{"code":"th","name":"ซอสเกาหลี","isauto":false,"isdelete":false},{"code":"en","name":"KOREA SAUCE","isauto":false,"isdelete":false},{"code":"ko","name":"한국 소스","isauto":false,"isdelete":false}],"price":"","qty":0,"isstock":false,"isdefault":false},{"guid":"8580a534-591e-46f6-801e-1504502d3e42","refbarcode":"","refproductcode":"","refunitcode":"","names":[{"code":"th","name":"ซอสกระเทียม ","isauto":false,"isdelete":false},{"code":"en","name":"GARLIC SAUCE","isauto":false,"isdelete":false},{"code":"ko","name":"마늘 소스","isauto":false,"isdelete":false}],"price":"","qty":0,"isstock":false,"isdefault":false},{"guid":"8074cfa4-a82b-4dbd-bf32-cb392cbad345","refbarcode":"","refproductcode":"","refunitcode":"","names":[{"code":"th","name":"ซอสหมาล่า ","isauto":false,"isdelete":false},{"code":"en","name":"MALA SAUCE ","isauto":false,"isdelete":false},{"code":"ko","name":"마라 소스","isauto":false,"isdelete":false}],"price":"","qty":0,"isstock":false,"isdefault":false}]}],"images":null,"useimageorcolor":true,"colorselect":"","colorselecthex":"","condition":false,"dividevalue":1,"standvalue":1,"isusesubbarcodes":false,"itemtype":0,"taxtype":0,"vattype":0,"issumpoint":false,"maxdiscount":"","isdividend":false,"refunitnames":null,"stockbarcode":"","qty":0,"refdividevalue":0,"refstandvalue":0,"vatcal":0,"refbarcodes":[],"bom":[{"guidfixed":"some-guid-fixed-value","names":[{"code":"name-code-1","name":"Name 1","isauto":true,"isdelete":false},{"code":"name-code-2","name":"Name 2","isauto":false,"isdelete":true}],"itemunitcode":"item-unit-code","itemunitnames":[{"code":"item-unit-code-1","name":"Item Unit Name 1","isauto":false,"isdelete":false}],"barcode":"1234567890","condition":true,"dividevalue":2.5,"standvalue":10,"qty":100}]}`

var giveProductBarcodeDoc models.ProductBarcodeDoc

func init() {

	err := json.Unmarshal([]byte(jsonStr), &giveProductBarcodeDoc)
	if err != nil {
		panic(err)
	}
}
func TestDeserializeJsonProductBarcode(t *testing.T) {
	assert.Equal(t, "2QoilMQkX9i6vtAE88ilEubnrhz", giveProductBarcodeDoc.ShopID)

	phaser := usecases.ProductBarcodePhaser{}
	got, err := phaser.PhaseProductBarcodeDoc(&giveProductBarcodeDoc)

	assert.NoError(t, err)
	assert.Equal(t, "2QoilMQkX9i6vtAE88ilEubnrhz", got.ShopID)
	assert.Equal(t, "PLATE", got.UnitCode)
	assert.Equal(t, "th", *got.UnitNames[0].Code)
	assert.Equal(t, "จาน", *got.UnitNames[0].Name)
	assert.Equal(t, 1, len(got.BOM))

	assert.Equal(t, "some-guid-fixed-value", (*got).BOM[0].BarcodeGuidFixed)
	assert.Equal(t, "1234567890", (*got).BOM[0].Barcode)

}
