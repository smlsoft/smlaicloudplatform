package productadmin

import (
	"context"
	"smlaicloudplatform/internal/product/productbarcode/models"
	productBarcodeModel "smlaicloudplatform/internal/product/productbarcode/models"
	"smlaicloudplatform/internal/repositories"
	"smlaicloudplatform/pkg/microservice"

	micromodels "smlaicloudplatform/pkg/microservice/models"

	"github.com/smlsoft/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IProductAdminMongoRepository interface {
	FindProductBarcodeByShopId(ctx context.Context, shopID string) ([]productBarcodeModel.ProductBarcodeDoc, error)
	FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.ProductBarcodeDoc, mongopagination.PaginationData, error)
	DeleteProductBarcodeByShopId(ctx context.Context, shopID string, userName string, ids []string) error
}

type ProductAdminMongoRepository struct {
	pst microservice.IPersisterMongo
	repositories.SearchRepository[models.ProductBarcodeDoc]
}

func NewProductAdminMongoRepository(pst microservice.IPersisterMongo) IProductAdminMongoRepository {
	return &ProductAdminMongoRepository{
		pst:              pst,
		SearchRepository: repositories.NewSearchRepository[models.ProductBarcodeDoc](pst),
	}
}

func (r ProductAdminMongoRepository) FindProductBarcodeByShopId(ctx context.Context, shopID string) ([]productBarcodeModel.ProductBarcodeDoc, error) {

	docList := []productBarcodeModel.ProductBarcodeDoc{}
	err := r.pst.Find(ctx, &productBarcodeModel.ProductBarcodeDoc{}, bson.M{"shopid": shopID}, &docList)
	if err != nil {
		return nil, err
	}

	return docList, nil
}

func (r ProductAdminMongoRepository) FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.ProductBarcodeDoc, mongopagination.PaginationData, error) {

	results, pagination, err := r.SearchRepository.FindPage(ctx, shopID, searchInFields, pageable)

	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	return results, pagination, nil
}

func (r ProductAdminMongoRepository) DeleteProductBarcodeByShopId(ctx context.Context, shopID string, userName string, ids []string) error {
	// err := r.pst.Delete(&productBarcodeModel.ProductBarcodeDoc{}, bson.M{"shopid": shopID})
	// if err != nil {
	// 	return err
	// }

	// err := r.pst.DeleteByID(&productBarcodeModel.ProductBarcodeDoc{}, bson.M{"shopid": shopID})
	// if err != nil {
	// 	return err
	// }

	err := r.pst.SoftBatchDeleteByID(ctx, &productBarcodeModel.ProductBarcodeDoc{}, userName, ids)
	if err != nil {
		return err
	}

	return nil
}
