package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/vfgl/journal/models"
)

type IJournalPgRepository interface {
	CreateInBatch(docList []models.JournalPg) error
	Create(doc models.JournalPg) error
	Update(shopID string, accountCode string, doc models.JournalPg) error
	Delete(shopID string, accountCode string) error
}

type JournalPgRepository struct {
	pst microservice.IPersister
}

func NewJournalPgRepository(pst microservice.IPersister) JournalPgRepository {
	return JournalPgRepository{
		pst: pst,
	}
}

func (repo JournalPgRepository) CreateInBatch(docList []models.JournalPg) error {
	err := repo.pst.CreateInBatch(docList, len(docList))
	if err != nil {
		return err
	}
	return nil
}

func (repo JournalPgRepository) Create(doc models.JournalPg) error {
	err := repo.pst.Create(doc)
	if err != nil {
		return err
	}
	return nil
}

func (repo JournalPgRepository) Update(shopID string, docNo string, doc models.JournalPg) error {
	err := repo.pst.Update(&doc, map[string]interface{}{
		"shopid": shopID,
		"docno":  docNo,
	})

	if err != nil {
		return err
	}
	return nil
}

func (repo JournalPgRepository) Delete(shopID string, docNo string) error {
	err := repo.pst.Delete(models.JournalPg{}, map[string]interface{}{
		"shopid": shopID,
		"docno":  docNo,
	})

	if err != nil {
		return err
	}
	return nil
}
