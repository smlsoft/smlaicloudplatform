package microservice

import "github.com/segmentio/ksuid"

func NewUUID() string {
	newid := ksuid.New()
	return newid.String()
}
