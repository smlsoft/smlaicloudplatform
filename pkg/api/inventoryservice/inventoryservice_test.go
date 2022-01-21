package inventoryservice

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestInventory(t *testing.T) {
	t.Log(primitive.NewObjectID())
	t.Log(time.Now())

}
