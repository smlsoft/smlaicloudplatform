package services

import (
	"encoding/json"
	"smlcloudplatform/pkg/vfgl/journal/models"
	"smlcloudplatform/pkg/vfgl/journal/repositories"
)

type IJournalConsumeService interface {
	Create(doc models.JournalDoc) error
	Update(shopID string, docNo string, doc models.JournalDoc) error
	Delete(shopID string, guid string) error
	SaveInBatch(docList []models.JournalDoc) error
}

type JournalConsumeService struct {
	repo repositories.IJournalPgRepository
}

func NewJournalConsumeService(repo repositories.IJournalPgRepository) JournalConsumeService {
	return JournalConsumeService{
		repo: repo,
	}
}

func (svc *JournalConsumeService) Create(doc models.JournalDoc) error {
	pgDoc := models.JournalPg{}

	tmpJsonDoc, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	tmpDoc := models.JournalPg{}
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

func (svc *JournalConsumeService) Update(shopID string, docNo string, doc models.JournalDoc) error {
	pgDoc := models.JournalPg{}

	tmpJsonDoc, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	tmpDoc := models.JournalPg{}
	err = json.Unmarshal([]byte(tmpJsonDoc), &tmpDoc)
	if err != nil {
		return err
	}

	err = svc.repo.Update(shopID, docNo, pgDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc *JournalConsumeService) Delete(shopID string, docNo string) error {
	err := svc.repo.Delete(shopID, docNo)

	if err != nil {
		return err
	}
	return nil
}

func (svc *JournalConsumeService) SaveInBatch(docList []models.JournalDoc) error {
	pgDocList := []models.JournalPg{}

	for _, doc := range docList {
		tmpJsonDoc, err := json.Marshal(doc)
		if err != nil {
			return err
		}
		tmpDoc := models.JournalPg{}
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
