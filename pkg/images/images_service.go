package images

import (
	"bytes"
	"mime/multipart"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/images/models"
	inventoryModel "smlcloudplatform/pkg/product/inventory/models"
	inventoryRepo "smlcloudplatform/pkg/product/inventory/repositories"
	"smlcloudplatform/pkg/utils"
	"strings"

	"errors"
)

type IImagesService interface {
	UploadImage(shopId string, fh *multipart.FileHeader) (*models.Image, error)
	UploadImageToProduct(shopID string, fh *multipart.FileHeader) error
	GetImageByProductCode(shopid string, itemguid string, index int) (string, *bytes.Buffer, error)
}

type ImagesService struct {
	persisterImage *microservice.PersisterImage
	invRepo        inventoryRepo.IInventoryRepository
	NewGUIDFn      func() string
}

func NewImageService(persisterImage *microservice.PersisterImage,
	inventoryRepo inventoryRepo.IInventoryRepository,
) *ImagesService {

	// check config storage location

	return &ImagesService{
		persisterImage: persisterImage,
		invRepo:        inventoryRepo,
		NewGUIDFn:      func() string { return utils.NewGUID() },
	}
}

func (svc ImagesService) UploadImage(shopId string, fh *multipart.FileHeader) (*models.Image, error) {

	fileUploadMetadataSlice := strings.Split(fh.Filename, ".")
	fileName := svc.NewGUIDFn() //fileUploadMetadataSlice[0]
	fileExtension := fileUploadMetadataSlice[1]

	fileName, err := svc.persisterImage.Upload(fh, shopId+"/"+fileName, fileExtension)

	if err != nil {
		return nil, err
	}
	// create thumbnail bla bla
	image := &models.Image{
		Uri: fileName,
	}

	return image, nil
}

func (svc ImagesService) UploadImageToProduct(shopID string, fh *multipart.FileHeader) error {

	fileUploadMetadataSlice := strings.Split(fh.Filename, ".")

	// find product by code
	fileName := fileUploadMetadataSlice[0]
	fileExtension := fileUploadMetadataSlice[1]

	findDoc, err := svc.invRepo.FindByItemBarcode(shopID, fileName)
	if err != nil {
		return err
	}

	uploadFileName, err := svc.persisterImage.Upload(fh, shopID+"/"+fileName, fileExtension)
	if err != nil {
		return err
	}

	var imageSlice []inventoryModel.InventoryImage
	if findDoc.Images != nil {
		//imageSlice = make([]models.InventoryImage, 1)
		//} else {
		imageSlice = *findDoc.Images
	}
	productImage := inventoryModel.InventoryImage{
		Uri: uploadFileName,
	}
	// push image
	imageSlice = append(imageSlice, productImage)

	findDoc.Images = &imageSlice

	// save and return
	err = svc.invRepo.Update(shopID, findDoc.GuidFixed, findDoc)
	return err
}

func (svc ImagesService) GetImageByProductCode(shopid string, itemguid string, index int) (string, *bytes.Buffer, error) {

	findDoc, err := svc.invRepo.FindByItemGuid(shopid, itemguid)

	if err != nil {
		return "", nil, err
	}

	if findDoc.Images == nil {
		return "", nil, errors.New("No Image")
	}

	productImage := *findDoc.Images
	inventoryImageLength := len(productImage)
	if inventoryImageLength < index {
		return "", nil, errors.New("Not Found Image")
	}

	imgFileUrl := productImage[index-1].Uri

	imageUri, buffer, err := svc.persisterImage.FilePersister.LoadFile(imgFileUrl)
	if err != nil {
		return "", buffer, err
	}

	return imageUri, buffer, nil
}
