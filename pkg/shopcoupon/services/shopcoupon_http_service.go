package services

import (
	"errors"
	"smlcloudplatform/pkg/shopcoupon/models"
	"smlcloudplatform/pkg/shopcoupon/repositories"
	"smlcloudplatform/pkg/utils"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IShopCouponHttpService interface {
	CreateShopCoupon(shopID string, authUsername string, doc models.ShopCoupon) (string, error)
	UpdateShopCoupon(shopID string, guid string, authUsername string, doc models.ShopCoupon) error
	DeleteShopCoupon(shopID string, guid string, authUsername string) error
	InfoShopCoupon(shopID string, guid string) (models.ShopCouponInfo, error)
	SearchShopCoupon(shopID string, filters map[string]interface{}, q string, page int, limit int, sort map[string]int) ([]models.ShopCouponInfo, mongopagination.PaginationData, error)
}

type ShopCouponHttpService struct {
	repo repositories.IShopCouponRepository
}

func NewShopCouponHttpService(repo repositories.IShopCouponRepository) *ShopCouponHttpService {

	return &ShopCouponHttpService{
		repo: repo,
	}
}

func (svc ShopCouponHttpService) CreateShopCoupon(shopID string, authUsername string, doc models.ShopCoupon) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := models.ShopCouponDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.ShopCoupon = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc ShopCouponHttpService) UpdateShopCoupon(shopID string, guid string, authUsername string, doc models.ShopCoupon) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.ShopCoupon = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc ShopCouponHttpService) DeleteShopCoupon(shopID string, guid string, authUsername string) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	err = svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)
	if err != nil {
		return err
	}

	return nil
}

func (svc ShopCouponHttpService) InfoShopCoupon(shopID string, guid string) (models.ShopCouponInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.ShopCouponInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.ShopCouponInfo{}, errors.New("document not found")
	}

	return findDoc.ShopCouponInfo, nil

}

func (svc ShopCouponHttpService) SearchShopCoupon(shopID string, filters map[string]interface{}, q string, page int, limit int, sort map[string]int) ([]models.ShopCouponInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"name1",
	}

	docList, pagination, err := svc.repo.FindPageFilterSort(shopID, filters, searchCols, q, page, limit, sort)

	if err != nil {
		return []models.ShopCouponInfo{}, pagination, err
	}

	return docList, pagination, nil
}
