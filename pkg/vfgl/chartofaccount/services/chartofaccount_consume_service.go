package services

import (
	"encoding/json"
	"smlcloudplatform/pkg/models/vfgl"
	"smlcloudplatform/pkg/vfgl/chartofaccount/repositories"
)

type IChartOfAccountConsumeService interface {
	Create(doc vfgl.ChartOfAccountDoc) error
	Update(shopID string, accountCode string, doc vfgl.ChartOfAccountDoc) error
	Delete(shopID string, guid string) error
	SaveInBatch(docList []vfgl.ChartOfAccountDoc) error
}

type ChartOfAccountConsumeService struct {
	repo repositories.IChartOfAccountPgRepository
}

func NewChartOfAccountConsumeService(repo repositories.IChartOfAccountPgRepository) ChartOfAccountConsumeService {
	return ChartOfAccountConsumeService{
		repo: repo,
	}
}

func (svc *ChartOfAccountConsumeService) Create(doc vfgl.ChartOfAccountDoc) error {
	pgDoc := vfgl.ChartOfAccountPG{}

	tmpJsonDoc, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	tmpDoc := vfgl.ChartOfAccountPG{}
	err = json.Unmarshal([]byte(tmpJsonDoc), &tmpDoc)
	if err != nil {
		return err
	}

	err = svc.repo.Create(pgDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc *ChartOfAccountConsumeService) Update(shopID string, accountCode string, doc vfgl.ChartOfAccountDoc) error {
	pgDoc := vfgl.ChartOfAccountPG{}

	tmpJsonDoc, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	tmpDoc := vfgl.ChartOfAccountPG{}
	err = json.Unmarshal([]byte(tmpJsonDoc), &tmpDoc)
	if err != nil {
		return err
	}

	err = svc.repo.Update(shopID, accountCode, pgDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc *ChartOfAccountConsumeService) Delete(shopID string, accountCode string) error {
	err := svc.repo.Delete(shopID, accountCode)

	if err != nil {
		return err
	}
	return nil
}

func (svc *ChartOfAccountConsumeService) SaveInBatch(docList []vfgl.ChartOfAccountDoc) error {
	pgDocList := []vfgl.ChartOfAccountPG{}

	for _, doc := range docList {
		tmpJsonDoc, err := json.Marshal(doc)
		if err != nil {
			return err
		}
		tmpDoc := vfgl.ChartOfAccountPG{}
		err = json.Unmarshal([]byte(tmpJsonDoc), &tmpDoc)
		if err != nil {
			return err
		}
		pgDocList = append(pgDocList, tmpDoc)
	}

	err := svc.repo.CreateInBatch(pgDocList)
	if err != nil {
		return err
	}

	return nil
}
