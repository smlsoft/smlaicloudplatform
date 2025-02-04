package services_test

// import (
// 	"smlaicloudplatform/internal/config"
// 	"smlaicloudplatform/internal/transaction/saleinvoice/models"
// 	"smlaicloudplatform/internal/transaction/saleinvoice/repositories"
// 	"smlaicloudplatform/internal/transaction/saleinvoice/services"
// 	"smlaicloudplatform/pkg/microservice"
// 	"testing"

// 	mastersync "smlaicloudplatform/internal/mastersync/repositories"
// 	productbarcode_repositories "smlaicloudplatform/internal/product/productbarcode/repositories"
// 	trans_cache "smlaicloudplatform/internal/transaction/repositories"

// 	"github.com/stretchr/testify/require"
// )

// var svc *services.SaleInvoiceHttpService

// func init() {
// 	cfg := config.NewConfig()
// 	ms, err := microservice.NewMicroservice(cfg)

// 	if err != nil {
// 		panic(err)
// 	}

// 	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
// 	cache := ms.Cacher(cfg.CacherConfig())
// 	producer := ms.Producer(cfg.MQConfig())

// 	repo := repositories.NewSaleInvoiceRepository(pst)
// 	repoMq := repositories.NewSaleInvoiceMessageQueueRepository(producer)

// 	productBarcodeRepo := productbarcode_repositories.NewProductBarcodeRepository(pst, cache)

// 	transRepo := trans_cache.NewCacheRepository(cache)
// 	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
// 	svc = services.NewSaleInvoiceHttpService(repo, transRepo, productBarcodeRepo, repoMq, masterSyncCacheRepo, services.SaleInvocieParser{})
// }

// func TestCreate(t *testing.T) {

// 	docSaleInvoice := models.SaleInvoice{}
// 	docSaleInvoice.DocNo = "d111"
// 	docSaleInvoice.IsPOS = true

// 	_, _, err := svc.CreateSaleInvoice("2Gf5cN6DP1kX7TYq3EJ1m4DKsJC", "error404", docSaleInvoice)

// 	require.Error(t, err)
// }

// func TestUpdate(t *testing.T) {

// 	docSaleInvoice := models.SaleInvoice{}
// 	docSaleInvoice.DocNo = "d111x"
// 	docSaleInvoice.IsPOS = true

// 	err := svc.UpdateSaleInvoice("2Gf5cN6DP1kX7TYq3EJ1m4DKsJC", "2Xn9nmhWahThqckRoMGWurt3ypP", "error404", docSaleInvoice)

// 	require.Error(t, err)
// }
