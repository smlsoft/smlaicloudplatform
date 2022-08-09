package models

const documentImageGroupCollectionName = "documentImages"

type DocumentImageGroup struct {
	DocumentRef    string                      `json:"documentref" bson:"documentref"`
	DocumentImages *[]DocumentImageGroupDetail `json:"documentimages" bson:"documentimages"`
}

func (DocumentImageGroup) CollectionName() string {
	return documentImageGroupCollectionName
}

type DocumentImageGroupDetail struct {
	Name       string `json:"name" bson:"name"`
	ImageUri   string `json:"imageuri" bson:"imageuri"`
	Module     string `json:"module" bson:"module"`
	DocGUIDRef string `json:"docguidref" bson:"docguidref"`
	Status     int8   `json:"status" bson:"status"`
}

type DocumentImageGroupRequest struct {
	DocumentRef    string   `json:"documentref" bson:"documentref"`
	DocumentImages []string `json:"documentimages" bson:"documentimages"`
}

func (DocumentImageGroupRequest) CollectionName() string {
	return documentImageGroupCollectionName
}
