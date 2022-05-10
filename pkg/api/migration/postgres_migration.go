package migration

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
)

func StartMigrateModel(ms *microservice.Microservice, cfg microservice.IConfig) error {
	pst := ms.Persister(cfg.PersisterConfig())

	// pst.DropTable(&models.InventoryData{}, &models.InventoryOption{}, &models.Option{}, &models.InventoryImage{}, &models.InventoryTag{}, &models.Choice{})

	// if err := pst.SetupJoinTable(&models.InventoryData{}, "Options", &models.InventoryOption{}); err != nil {
	// 	fmt.Printf("Failed to setup join table , got error %v \n", err)
	// 	return err
	// }

	// pst.DropTable(models.Trans{}, models.TransItemDetail{})

	pst.AutoMigrate(
		// &saleinvoice.SaleInvoiceTable{},
		// &saleinvoice.SaleInvoiceDetailTable{},
		// &models.InventoryImage{},
		// &models.InventoryTag{},

		// &models.CategoryData{},

		// &models.InventoryData{},
		// &models.InventoryOption{},
		// &models.Option{},
		// &models.Choice{},
		// &models.InventoryIndex{},
		models.Trans{},
		models.TransItemDetail{},
	)

	pst.AutoMigrate()

	return nil
}
