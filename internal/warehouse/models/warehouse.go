package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"smlcloudplatform/internal/models"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const warehouseCollectionName = "warehouse"

type Warehouse struct {
	models.PartitionIdentity `bson:"inline"`
	Code                     string          `json:"code" bson:"code"`
	Names                    *[]models.NameX `json:"names" bson:"names" validate:"required,min=1,unique=Code,dive"`
	Location                 *[]Location     `json:"location" bson:"location" validate:"omitempty,unique=Code,dive"`
}

type Location struct {
	Code  string          `json:"code" bson:"code"`
	Names *[]models.NameX `json:"names" bson:"names" validate:"required,min=1,unique=Code,dive"`
	Shelf *[]Shelf        `json:"shelf" bson:"shelf" validate:"omitempty,unique=Code,dive"`
}

type Shelf struct {
	Code string `json:"code" bson:"code"`
	Name string `json:"name" bson:"name" validate:"required,min=1"`
}

type LocationInfo struct {
	GuidFixed      string          `json:"guidfixed" bson:"guidfixed"`
	WarehouseCode  string          `json:"warehousecode" bson:"warehousecode"`
	WarehouseNames *[]models.NameX `json:"warehousenames" bson:"warehousenames"`
	LocationCode   string          `json:"locationcode" bson:"locationcode"`
	LocationNames  *[]models.NameX `json:"locationnames" bson:"locationnames"`
	Shelf          []Shelf         `json:"shelf" bson:"shelf"`
}

func (LocationInfo) CollectionName() string {
	return warehouseCollectionName
}

type ShelfInfo struct {
	GuidFixed      string          `json:"guidfixed" bson:"guidfixed"`
	WarehouseCode  string          `json:"warehousecode" bson:"warehousecode"`
	WarehouseNames *[]models.NameX `json:"warehousenames" bson:"warehousenames"`
	LocationCode   string          `json:"locationcode" bson:"locationcode"`
	LocationNames  *[]models.NameX `json:"locationnames" bson:"locationnames"`
	ShelfCode      string          `json:"shelfcode" bson:"shelfcode"`
	ShelfName      string          `json:"shelfname" bson:"shelfname"`
}

func (ShelfInfo) CollectionName() string {
	return warehouseCollectionName
}

type LocationRequest struct {
	WarehouseCode string          `json:"warehousecode" bson:"warehousecode" validate:"required"`
	Code          string          `json:"locationcode" bson:"locationcode" validate:"required"`
	Names         *[]models.NameX `json:"locationnames" bson:"locationnames"`
	Shelf         []Shelf         `json:"shelf" bson:"shelf"`
}

type ShelfRequest struct {
	WarehouseCode string `json:"warehousecode" bson:"warehousecode" validate:"required"`
	LocationCode  string `json:"locationcode" bson:"locationcode" validate:"required"`
	Code          string `json:"shelfcode" bson:"shelfcode" validate:"required"`
	Name          string `json:"shelfname" bson:"shelfname"`
}

type WarehouseInfo struct {
	models.DocIdentity `bson:"inline"`
	Warehouse          `bson:"inline"`
}

func (WarehouseInfo) CollectionName() string {
	return warehouseCollectionName
}

type WarehouseData struct {
	models.ShopIdentity `bson:"inline"`
	WarehouseInfo       `bson:"inline"`
}

type WarehouseDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	WarehouseData      `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (WarehouseDoc) CollectionName() string {
	return warehouseCollectionName
}

type WarehouseItemGuid struct {
	Code string `json:"code" bson:"code"`
}

func (WarehouseItemGuid) CollectionName() string {
	return warehouseCollectionName
}

type WarehouseActivity struct {
	WarehouseData       `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (WarehouseActivity) CollectionName() string {
	return warehouseCollectionName
}

type WarehouseDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (WarehouseDeleteActivity) CollectionName() string {
	return warehouseCollectionName
}

type WarehouseMessageQueue struct {
	models.ShopIdentity `bson:"inline"`
	models.DocIdentity  `bson:"inline"`
	Warehouse           `bson:"inline"`
}

type WarehousePG struct {
	models.ShopIdentity `gorm:"embedded;"`
	GuidFixed           string       `json:"guidfixed" gorm:"column:guidfixed;primaryKey"`
	Code                string       `json:"code" gorm:"column:code"`
	Names               models.JSONB `json:"names" gorm:"column:names;type:jsonb"`
	Location            LocationsPG  `json:"location" gorm:"column:location;type:jsonb"`
}

func (WarehousePG) TableName() string {
	return "warehouse"
}

func (jd *WarehousePG) BeforeCreate(tx *gorm.DB) error {

	tx.Statement.AddClause(clause.OnConflict{
		UpdateAll: true,
	})
	return nil
}

func (s *WarehousePG) CompareTo(other *WarehousePG) bool {

	diff := cmp.Diff(s, other,
		cmpopts.IgnoreFields(WarehousePG{}, "guidfixed"),
	)

	return diff == ""
}

type LocationsPG []Location

func (a LocationsPG) Value() (driver.Value, error) {

	j, err := json.Marshal(a)
	return j, err
}

func (a *LocationsPG) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
