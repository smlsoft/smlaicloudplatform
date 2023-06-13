package api

import (
	"smlcloudplatform/pkg/utils"
	journalModels "smlcloudplatform/pkg/vfgl/journal/models"
	journalRepo "smlcloudplatform/pkg/vfgl/journal/repositories"
)

func (m *MigrationService) ImportJournal(journals []journalModels.JournalDoc) error {
	journalpgRepo := journalRepo.NewJournalRepository(m.mongoPersister)
	journalMQRepo := journalRepo.NewJournalMqRepository(m.mqPersister)

	for _, journal := range journals {
		journal.GuidFixed = utils.NewGUID()
		_, err := journalpgRepo.Create(journal)
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
