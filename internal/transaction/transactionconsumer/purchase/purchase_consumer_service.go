package purchase

import (
	pkgModels "smlcloudplatform/internal/models"
	"smlcloudplatform/internal/transaction/models"
	"smlcloudplatform/pkg/logger"
)

type IPurchaseTransactionConsumerService interface {
	Upsert(shopID string, docNo string, doc models.PurchaseTransactionPG) error
	Delete(shopID string, docNo string) error
}

type PurchaseTransactionConsumerService struct {
	repo IPurchaseTransactionPGRepository
}

func NewPurchaseTransactionService(repo IPurchaseTransactionPGRepository) IPurchaseTransactionConsumerService {
	return &PurchaseTransactionConsumerService{
		repo: repo,
	}
}

func (s *PurchaseTransactionConsumerService) Upsert(shopID string, docNo string, doc models.PurchaseTransactionPG) error {
	findDoc, err := s.repo.Get(shopID, docNo)
	if err != nil {
		err = s.repo.Create(doc)
		if err != nil {
			return err
		}
	} else {

		isEqual := findDoc.CompareTo(&doc)

		if isEqual == false {
			err = s.repo.Update(shopID, docNo, doc)
			if err != nil {
				return err
			}
		} else {
			logger.GetLogger().Debug("Doc is equal, skip update")
		}
	}

	return nil
}

func (s *PurchaseTransactionConsumerService) Delete(shopID string, docNo string) error {

	err := s.repo.DeleteData(shopID, docNo, models.PurchaseTransactionPG{
		TransactionPG: models.TransactionPG{
			ShopIdentity: pkgModels.ShopIdentity{
				ShopID: shopID,
			},
			DocNo: docNo,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
