package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/shopcoupon/models"

	"github.com/userplant/mongopagination"
)

type IShopCouponRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.ShopCouponDoc) (string, error)
	CreateInBatch(docList []models.ShopCouponDoc) error
	Update(shopID string, guid string, doc models.ShopCouponDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.ShopCouponInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.ShopCouponDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.ShopCouponInfo, mongopagination.PaginationData, error)
}

type ShopCouponRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.ShopCouponDoc]
	repositories.SearchRepository[models.ShopCouponInfo]
	repositories.GuidRepository[models.ShopCouponItemGuid]
}

func NewShopCouponRepository(pst microservice.IPersisterMongo) *ShopCouponRepository {

	insRepo := &ShopCouponRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.ShopCouponDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.ShopCouponInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.ShopCouponItemGuid](pst)

	return insRepo
}
