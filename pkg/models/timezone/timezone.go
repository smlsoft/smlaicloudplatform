package timezone

type Timezone struct {
	TimezoneLabel  string `json:"timezonelabel" bson:"timezonelabel"`
	TimezoneOffset int    `json:"timezoneoffset" bson:"timezoneoffset"`
}
