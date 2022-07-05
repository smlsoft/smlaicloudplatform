package member_test

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/member"
	"smlcloudplatform/pkg/member/models"
	"testing"
)

func newPgRepo() member.MemberPGRepository {
	persisterConfig := mock.NewPersisterPostgresqlConfig()
	pst := microservice.NewPersister(persisterConfig)
	repo := member.NewMemberPGRepository(pst)
	return repo
}

func TestCreate(t *testing.T) {
	repo := newPgRepo()

	idx := models.MemberIndex{}
	idx.ID = "134567"
	idx.ShopID = "shopidx001"
	idx.GuidFixed = "fixguid"
	err := repo.Create(idx)

	if err != nil {
		t.Error(err)
	}
}

func TestCount(t *testing.T) {
	repo := newPgRepo()

	count, err := repo.Count("shopidx001", "fixguid")

	if err != nil {
		t.Error(err)
	}

	t.Log(count)
}

func TestFindByGuid(t *testing.T) {
	repo := newPgRepo()
	inv, err := repo.FindByGuid("shopidx001", "fixguid")

	if err != nil {
		t.Error(err)
	}

	t.Log(inv)
}

func TestDelete(t *testing.T) {
	repo := newPgRepo()

	err := repo.Delete("shopidx001", "fixguid")

	if err != nil {
		t.Error(err)
	}
}
