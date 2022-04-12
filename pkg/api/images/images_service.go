package images

import (
	"mime/multipart"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/inventory"
	"smlcloudplatform/pkg/models"

	"errors"
)

type IImagesService interface {
	UploadImage(fh *multipart.FileHeader) (*models.Image, error)
	GetImageByProductCode(shopid string, itemguid string, index int) (string, error)
	GetStoragePath() string
}

type ImagesService struct {
	persisterImage *microservice.PersisterImage
	invRepo        inventory.IInventoryRepository
}

func NewImageService(persisterImage *microservice.PersisterImage,
	inventoryRepo inventory.IInventoryRepository,
) *ImagesService {

	// check config storage location

	return &ImagesService{
		persisterImage: persisterImage,
		invRepo:        inventoryRepo,
	}
}

func (svc ImagesService) UploadImage(fh *multipart.FileHeader) (*models.Image, error) {

	fileName, err := svc.persisterImage.Upload(fh)

	if err != nil {
		return nil, err
	}
	// create thumbnail bla bla
	image := &models.Image{
		Uri: fileName,
	}

	return image, nil
}

func (svc ImagesService) GetImageByProductCode(shopid string, itemguid string, index int) (string, error) {

	findDoc, err := svc.invRepo.FindByItemGuid(itemguid, shopid)

	if err != nil {
		return "", err
	}

	if findDoc.Images == nil {
		return "", errors.New("No Image")
	}

	productImage := *findDoc.Images
	inventoryImageLength := len(productImage)
	if inventoryImageLength < index {
		return "", errors.New("Not Found Image")
	}

	imgFileName := productImage[index-1].Url

	return imgFileName, nil
}

func (svc ImagesService) GetStoragePath() string {

	storageCfg := microservice.NewStorageFileConfig()
	return storageCfg.StorageDataPath()
}
