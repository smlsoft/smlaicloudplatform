package microservice

import (
	"context"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/apex/log"
)

// imprement PersisterFile Interface
type PersisterAzureBlob struct{}

func NewPersisterAzureBlob() *PersisterAzureBlob {
	return &PersisterAzureBlob{}
}

func (p *PersisterAzureBlob) Save(fh *multipart.FileHeader, fileName string, fileExtension string) (string, error) {

	file, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	imageFileName := fmt.Sprintf("%s.%s", fileName, fileExtension)

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	AZURE_STORAGE_ACCOUNT_NAME := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	AZURE_STORAGE_ACCOUNT_KEY := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	url := fmt.Sprintf("https://%s.blob.core.windows.net/", AZURE_STORAGE_ACCOUNT_NAME)

	credential, err := azblob.NewSharedKeyCredential(AZURE_STORAGE_ACCOUNT_NAME, AZURE_STORAGE_ACCOUNT_KEY)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}

	service, err := azblob.NewServiceClientWithSharedKey(url, credential, nil)
	if err != nil {
		log.Fatal("Invalid Service Client with error: " + err.Error())
	}

	container, err := service.NewContainerClient("dedeposproductimage")
	if err != nil {
		log.Fatal("Invalid Create Containnner with error: " + err.Error())
	}

	blockBlob, err := container.NewBlockBlobClient(imageFileName)
	if err != nil {
		log.Fatal("Invalid New BlockBlob with error: " + err.Error())
	}

	context := context.TODO()
	//blockOptions := azblob.HighLevelUploadToBlockBlobOption{}
	//var blockOptions azblob.HighLevelUploadToBlockBlobOption{}
	uploadOption := azblob.UploadOption{}

	_, err = blockBlob.UploadBuffer(context, fileBytes, uploadOption)
	if err != nil {
		return "", err
	}

	blobUri := blockBlob.URL()

	return blobUri, nil
}
