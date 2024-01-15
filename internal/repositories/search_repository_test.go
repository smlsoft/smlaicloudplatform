package repositories_test

import (
	"smlcloudplatform/internal/repositories"
	"smlcloudplatform/internal/restaurant/kitchen/models"
	"testing"
)

type TestModel struct {
}

func TestSearchRepositories(t *testing.T) {
	repo := repositories.NewSearchRepository[models.KitchenInfo](nil)
	xx := repo.CreateTextFilter([]string{"name_1", "sx"}, "okลาก่อน")
	t.Log(xx)
}
