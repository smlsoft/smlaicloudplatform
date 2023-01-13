package repositories_test

import (
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/restaurant/kitchen/models"
	"testing"
)

type TestModel struct {
}

func TestSearchRepositories(t *testing.T) {
	repo := repositories.NewSearchRepository[models.KitchenInfo](nil)
	xx := repo.SearchTextFilter([]string{"name_1", "sx"}, "qx ee")
	t.Log(xx)
}
