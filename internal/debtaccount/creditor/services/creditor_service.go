package services

import (
	"smlcloudplatform/internal/debtaccount/creditor/models"
	"smlcloudplatform/internal/debtaccount/creditor/repositories"
)

type ICreditorConsumerService interface {
	Upsert(shopID string, code string, doc models.CreditorPG) error
	Delete(shopID string, code string) error
}

type CreditorConsumerService struct {
	repo repositories.ICreditorPostgresRepository
}

func NewCreditorConsumerService(repo repositories.ICreditorPostgresRepository) ICreditorConsumerService {
	return &CreditorConsumerService{
		repo: repo,
	}
}

func (s *CreditorConsumerService) Upsert(shopID string, code string, doc models.CreditorPG) error {
	findDoc, err := s.repo.Get(shopID, code)
	if err != nil || findDoc == nil {
		err = s.repo.Create(doc)
		if err != nil {
			return err
		}
	} else {

		isEqual := findDoc.CompareTo(&doc)

		if !isEqual {
			err = s.repo.Update(shopID, code, doc)
			if err != nil {
				return err
			}
		} else {
			// logger.GetLogger().Debug("Doc is equal, skip update")
		}
	}

	return nil
}

func (s *CreditorConsumerService) Delete(shopID string, code string) error {
	err := s.repo.Delete(shopID, code)
	if err != nil {
		return err
	}
	return nil
}
