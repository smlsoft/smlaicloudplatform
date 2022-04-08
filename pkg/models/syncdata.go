package models

type SyncData struct {
	MyGuid    string      `json:"myguid"`
	TableCode string      `json:"tablecode"`
	Mode      int         `json:"mode"`
	Data      interface{} `json:"data"`
}

type BodyX struct {
	Guid    string   `json:"guid"`
	Details []Detail `json:"details"`
}

type Detail struct {
	Opt string `json:"opt"`
}

type SyncInventoryData struct {
	MyGuid    string    `json:"myguid"`
	TableCode string    `json:"tablecode"`
	Mode      int       `json:"mode"`
	Data      Inventory `json:"data"`
}
