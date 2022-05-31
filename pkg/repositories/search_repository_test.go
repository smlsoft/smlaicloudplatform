package repositories_test

import (
	"smlcloudplatform/pkg/models/restaurant"
	"smlcloudplatform/pkg/repositories"
	"testing"
)

type TestModel struct {
}

func TestSearchRepositories(t *testing.T) {
	repo := repositories.NewSearchRepository[restaurant.KitchenInfo](nil)
	xx := repo.SearchTextFilter([]string{"name_1", "sx"}, "qx ee")
	t.Log(xx)
}
