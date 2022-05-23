package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models/vfgl"
)

type IChartOfAccountPgRepository interface {
	CreateInBatch(docList []vfgl.ChartOfAccountPG) error
	Create(doc vfgl.ChartOfAccountPG) error
	Update(shopID string, accountCode string, doc vfgl.ChartOfAccountPG) error
	Delete(shopID string, accountCode string) error
}

type ChartOfAccountPgRepository struct {
	pst microservice.IPersister
}

func NewChartOfAccountPgRepository(pst microservice.IPersister) ChartOfAccountPgRepository {
	return ChartOfAccountPgRepository{
		pst: pst,
	}
}

func (repo ChartOfAccountPgRepository) CreateInBatch(docList []vfgl.ChartOfAccountPG) error {
	err := repo.pst.CreateInBatch(docList, len(docList))
	if err != nil {
		return err
	}
	return nil
}

func (repo ChartOfAccountPgRepository) Create(doc vfgl.ChartOfAccountPG) error {
	err := repo.pst.Create(doc)
	if err != nil {
		return err
	}
	return nil
}

func (repo ChartOfAccountPgRepository) Update(shopID string, accountCode string, doc vfgl.ChartOfAccountPG) error {
	err := repo.pst.Update(&doc, map[string]interface{}{
		"shopid":      shopID,
		"accountcode": accountCode,
	})

	if err != nil {
		return err
	}
	return nil
}

func (repo ChartOfAccountPgRepository) Delete(shopID string, accountCode string) error {
	err := repo.pst.Delete(vfgl.ChartOfAccountPG{}, map[string]interface{}{
		"shopid":      shopID,
		"accountcode": accountCode,
	})

	if err != nil {
		return err
	}
	return nil
}
