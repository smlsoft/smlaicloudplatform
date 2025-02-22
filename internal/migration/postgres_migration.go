package migration

import (
	"smlaicloudplatform/internal/config"
	vfgl "smlaicloudplatform/internal/vfgl/journal/models"
	"smlaicloudplatform/pkg/microservice"
)

func StartMigrateModel(ms *microservice.Microservice, cfg config.IConfig) error {
	pst := ms.Persister(cfg.PersisterConfig())

	// pst.DropTable(&models.InventoryData{}, &models.InventoryOption{}, &models.Option{}, &models.InventoryImage{}, &models.InventoryTag{}, &models.Choice{})

	// if err := pst.SetupJoinTable(&models.InventoryData{}, "Options", &models.InventoryOption{}); err != nil {
	// 	fmt.Printf("Failed to setup join table , got error %v \n", err)
	// 	return err
	// }

	// pst.DropTable(vfgl.JournalPg{}, vfgl.JournalDetailPg{})

	pst.AutoMigrate(
		// &saleinvoice.SaleinvoiceTable{},
		// &saleinvoice.SaleinvoiceDetailTable{},
		// &models.InventoryImage{},
		// &models.InventoryTag{},

		// &models.CategoryData{},

		// &models.InventoryData{},
		// &models.InventoryOption{},
		// &models.Option{},
		// &models.Choice{},
		// &models.InventoryIndex{},
		// models.Trans{},
		// models.TransItemDetail{},
		vfgl.JournalPg{},
		vfgl.JournalDetailPg{},
	)

	// pst.AutoMigrate()

	return nil
}
