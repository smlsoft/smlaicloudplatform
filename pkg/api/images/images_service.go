package images

import (
	"mime/multipart"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
)

type IImagesService interface {
	UploadImage(fh *multipart.FileHeader) (*models.Image, error)
}

type ImagesService struct {
	persisterImage *microservice.PersisterImage
}

func NewImageService(persisterImage *microservice.PersisterImage) *ImagesService {

	// check config storage location

	return &ImagesService{
		persisterImage: persisterImage,
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
