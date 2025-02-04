package api

import (
	"context"
	"smlaicloudplatform/internal/utils"
	journalModels "smlaicloudplatform/internal/vfgl/journal/models"
	journalRepo "smlaicloudplatform/internal/vfgl/journal/repositories"
)

func (m *MigrationService) ImportJournal(journals []journalModels.JournalDoc) error {
	journalpgRepo := journalRepo.NewJournalRepository(m.mongoPersister)
	journalMQRepo := journalRepo.NewJournalMqRepository(m.mqPersister)

	for _, journal := range journals {
		journal.GuidFixed = utils.NewGUID()
		_, err := journalpgRepo.Create(context.Background(), journal)
		if err != nil {
			return err
		}

		err = journalMQRepo.Create(journal)
		if err != nil {
			return err
		}
	}
	return nil
}
