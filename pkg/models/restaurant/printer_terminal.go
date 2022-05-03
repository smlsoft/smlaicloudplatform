package restaurant

type PrinterTerminal struct {
	Code    string `json:"code" bson:"code"`
	Name1   string `json:"name1" bson:"name1" gorm:"name1"`
	Name2   string `json:"name2,omitempty" bson:"name2,omitempty"`
	Name3   string `json:"name3,omitempty" bson:"name3,omitempty"`
	Name4   string `json:"name4,omitempty" bson:"name4,omitempty"`
	Name5   string `json:"name5,omitempty" bson:"name5,omitempty"`
	Address string `json:"address" bson:"address" `
	Type    int8   `json:"type" bson:"type"`
}
