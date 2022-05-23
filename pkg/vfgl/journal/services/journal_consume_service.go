package services

import (
	"encoding/json"
	"smlcloudplatform/pkg/models/vfgl"
	"smlcloudplatform/pkg/vfgl/journal/repositories"
)

type IJournalConsumeService interface {
	Create(doc vfgl.JournalDoc) error
	Update(shopID string, docNo string, doc vfgl.JournalDoc) error
	Delete(shopID string, guid string) error
	SaveInBatch(docList []vfgl.JournalDoc) error
}

type JournalConsumeService struct {
	repo repositories.IJournalPgRepository
}

func NewJournalConsumeService(repo repositories.IJournalPgRepository) JournalConsumeService {
	return JournalConsumeService{
		repo: repo,
	}
}

func (svc *JournalConsumeService) Create(doc vfgl.JournalDoc) error {
	pgDoc := vfgl.JournalPg{}

	tmpJsonDoc, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	tmpDoc := vfgl.JournalPg{}
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

func (svc *JournalConsumeService) Update(shopID string, docNo string, doc vfgl.JournalDoc) error {
	pgDoc := vfgl.JournalPg{}

	tmpJsonDoc, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	tmpDoc := vfgl.JournalPg{}
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

func (svc *JournalConsumeService) SaveInBatch(docList []vfgl.JournalDoc) error {
	pgDocList := []vfgl.JournalPg{}

	for _, doc := range docList {
		tmpJsonDoc, err := json.Marshal(doc)
		if err != nil {
			return err
		}
		tmpDoc := vfgl.JournalPg{}
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
