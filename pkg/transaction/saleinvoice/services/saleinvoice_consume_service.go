package services

import (
	"encoding/json"
	"smlcloudplatform/pkg/transaction/saleinvoice/models"
	"smlcloudplatform/pkg/transaction/saleinvoice/repositories"

	"gorm.io/gorm"
)

type ISaleinvoiceConsumeService interface {
	Create(doc models.SaleinvoiceDoc) error
	Update(shopID string, docNo string, doc models.SaleinvoiceDoc) error
	Delete(shopID string, guid string) error
	SaveInBatch(docList []models.SaleinvoiceDoc) error
	UpSert(shopID string, docNo string, doc models.SaleinvoiceDoc) (*models.SaleinvoicePg, error)
}

type SaleinvoiceConsumeService struct {
	repo repositories.ISaleinvoicePgRepository
}

func NewSaleinvoiceConsumeService(repo repositories.ISaleinvoicePgRepository) SaleinvoiceConsumeService {
	return SaleinvoiceConsumeService{
		repo: repo,
	}
}

func (svc *SaleinvoiceConsumeService) Create(doc models.SaleinvoiceDoc) (*models.SaleinvoicePg, error) {
	pgDoc := models.SaleinvoicePg{}

	tmpJsonDoc, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(tmpJsonDoc), &pgDoc)
	if err != nil {
		return nil, err
	}

	err = svc.repo.Create(pgDoc)

	if err != nil {
		return nil, err
	}
	return &pgDoc, nil
}

func (svc *SaleinvoiceConsumeService) Update(shopID string, docNo string, doc models.SaleinvoiceDoc) error {
	pgDoc := models.SaleinvoicePg{}

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

func (svc *SaleinvoiceConsumeService) Delete(shopID string, docNo string) error {
	err := svc.repo.Delete(shopID, docNo)

	if err != nil {
		return err
	}
	return nil
}

func (svc *SaleinvoiceConsumeService) SaveInBatch(docList []models.SaleinvoiceDoc) error {
	pgDocList := []models.SaleinvoicePg{}

	for _, doc := range docList {
		tmpJsonDoc, err := json.Marshal(doc)
		if err != nil {
			return err
		}
		tmpDoc := models.SaleinvoicePg{}
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

func (svc *SaleinvoiceConsumeService) UpSert(shopID string, docNo string, doc models.SaleinvoiceDoc) (*models.SaleinvoicePg, error) {

	docPg := models.SaleinvoicePg{}

	tmpJsonDoc, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(tmpJsonDoc), &docPg)
	if err != nil {
		return nil, err
	}

	data, err := svc.repo.Get(shopID, docNo)
	if err == gorm.ErrRecordNotFound {
		data, err = svc.Create(doc)
		if err != nil {
			return nil, err
		}

	} else if data != nil {
		/*
			data.SaleinvoiceBody = doc.SaleinvoiceBody

			// check detail
			for tmpIdx, tmp := range *data.AccountBook {
				for detailIdx, detail := range *docPg.AccountBook {
					if tmpIdx == detailIdx && tmp.AccountCode == detail.AccountCode {
						tmpAccBook := *docPg.AccountBook
						tmpAccBook[detailIdx].ID = tmp.ID
					}
				}
			}
		*/
		if err = svc.repo.Update(shopID, doc.DocNo, docPg); err != nil {
			return nil, err
		}
	}

	return &docPg, nil
}
