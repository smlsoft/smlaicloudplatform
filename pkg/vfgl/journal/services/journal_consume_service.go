package services

import (
	"encoding/json"
	"smlcloudplatform/pkg/vfgl/journal/models"
	"smlcloudplatform/pkg/vfgl/journal/repositories"

	"gorm.io/gorm"
)

type IJournalConsumeService interface {
	Create(doc models.JournalDoc) error
	Update(shopID string, docNo string, doc models.JournalDoc) error
	Delete(shopID string, guid string) error
	SaveInBatch(docList []models.JournalDoc) error
	UpSert(shopID string, docNo string, doc models.JournalDoc) error
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

	err = json.Unmarshal([]byte(tmpJsonDoc), &pgDoc)
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

	err = json.Unmarshal([]byte(tmpJsonDoc), &pgDoc)
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

func (svc *JournalConsumeService) UpSert(shopID string, docNo string, doc models.JournalDoc) error {

	docPg := models.JournalPg{}

	tmpJsonDoc, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(tmpJsonDoc), &docPg)
	if err != nil {
		return err
	}

	data, err := svc.repo.Get(shopID, docNo)
	if err == gorm.ErrRecordNotFound {
		if err = svc.Create(doc); err != nil {
			return err
		}
	} else if data != nil {

		data.JournalBody = doc.JournalBody

		// check detail
		for tmpIdx, tmp := range *data.AccountBook {
			for detailIdx, detail := range *docPg.AccountBook {
				if tmpIdx == detailIdx {
					detail.AccountCode = tmp.AccountCode
					detail.AccountName = tmp.AccountName
					detail.CreditAmount = tmp.CreditAmount
					detail.DebitAmount = tmp.CreditAmount
					detail.ID = tmp.ID
				}
			}
		}

		if err = svc.repo.Update(shopID, doc.DocNo, *data); err != nil {
			return err
		}
	}

	return nil
}
