package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models/vfgl"
)

type IJournalPgRepository interface {
	CreateInBatch(docList []vfgl.JournalPg) error
	Create(doc vfgl.JournalPg) error
	Update(shopID string, accountCode string, doc vfgl.JournalPg) error
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

func (repo JournalPgRepository) CreateInBatch(docList []vfgl.JournalPg) error {
	err := repo.pst.CreateInBatch(docList, len(docList))
	if err != nil {
		return err
	}
	return nil
}

func (repo JournalPgRepository) Create(doc vfgl.JournalPg) error {
	err := repo.pst.Create(doc)
	if err != nil {
		return err
	}
	return nil
}

func (repo JournalPgRepository) Update(shopID string, docNo string, doc vfgl.JournalPg) error {
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
	err := repo.pst.Delete(vfgl.JournalPg{}, map[string]interface{}{
		"shopid": shopID,
		"docno":  docNo,
	})

	if err != nil {
		return err
	}
	return nil
}
