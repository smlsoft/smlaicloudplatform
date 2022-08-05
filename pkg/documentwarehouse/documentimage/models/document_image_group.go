package models

const documentImageGroupCollectionName = "documentImages"

type DocumentImageGroup struct {
	DocumentRef    string   `json:"documentref" bson:"documentref"`
	DocumentImages []string `json:"documentimages" bson:"documentimages"`
}

func (DocumentImageGroup) CollectionName() string {
	return documentImageGroupCollectionName
}
