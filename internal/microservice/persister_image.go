package microservice

import (
	"fmt"
	"mime/multipart"
	"net/url"
	"path"
	"strings"

	"github.com/google/uuid"
)

type IPersisterImage interface {
	CreateThumbnail() error
	Upload() error
}

type PersisterImage struct {
	FilePersister *PersisterFile
}

func NewPersisterImage(filePersister *PersisterFile) *PersisterImage {
	return &PersisterImage{
		FilePersister: filePersister,
	}
}

func (pst *PersisterImage) Upload(fh *multipart.FileHeader) (string, error) {

	//
	u, err := url.Parse(pst.FilePersister.StoreDataUri)
	if err != nil {
		return "", err
	}

	uniqueId := uuid.New()
	filename := strings.Replace(uniqueId.String(), "-", "", -1)
	fileExt := strings.Split(fh.Filename, ".")[1]
	imageFileName := fmt.Sprintf("%s.%s", filename, fileExt)

	err = pst.FilePersister.Save(fh, filename, fileExt)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, imageFileName)
	imageUri := u.String()

	return imageUri, nil
}
