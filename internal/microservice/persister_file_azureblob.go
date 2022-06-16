package microservice

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

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
	AZURE_STORAGE_CONTAINER_NAME := os.Getenv("AZURE_STORAGE_CONTAINER_NAME")

	url := fmt.Sprintf("https://%s.blob.core.windows.net/", AZURE_STORAGE_ACCOUNT_NAME)

	credential, err := azblob.NewSharedKeyCredential(AZURE_STORAGE_ACCOUNT_NAME, AZURE_STORAGE_ACCOUNT_KEY)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}

	service, err := azblob.NewServiceClientWithSharedKey(url, credential, nil)
	if err != nil {
		log.Fatal("Invalid Service Client with error: " + err.Error())
	}

	container, err := service.NewContainerClient(AZURE_STORAGE_CONTAINER_NAME)
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

func (p *PersisterAzureBlob) LoadFile(fileName string) (string, *bytes.Buffer, error) {
	//imgFileName := strings.Replace(fileName, pst.StoreDataUri, "", -1)
	//storateFileName := filepath.Join(pst.StoreFilePath, imgFileName)

	if strings.HasPrefix(fileName, "http") {
		resp, err := http.Get(fileName)
		if err != nil {
			return "", nil, err
		}
		defer resp.Body.Close()

		downloadedData := &bytes.Buffer{}
		downloadedData.ReadFrom(resp.Body)

		return fileName, downloadedData, nil
	}

	// download blob from azure
	AZURE_STORAGE_ACCOUNT_NAME := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	AZURE_STORAGE_ACCOUNT_KEY := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	AZURE_STORAGE_CONTAINER_NAME := os.Getenv("AZURE_STORAGE_CONTAINER_NAME")

	url := fmt.Sprintf("https://%s.blob.core.windows.net/", AZURE_STORAGE_ACCOUNT_NAME)

	credential, err := azblob.NewSharedKeyCredential(AZURE_STORAGE_ACCOUNT_NAME, AZURE_STORAGE_ACCOUNT_KEY)
	if err != nil {
		return "", nil, errors.New("Invalid credentials with error: " + err.Error())
	}

	service, err := azblob.NewServiceClientWithSharedKey(url, credential, nil)
	if err != nil {
		return "", nil, errors.New("Invalid Service Client with error: " + err.Error())
	}

	container, err := service.NewContainerClient(AZURE_STORAGE_CONTAINER_NAME)
	if err != nil {
		return "", nil, errors.New("Invalid Create Containnner with error: " + err.Error())
	}

	blockBlob, err := container.NewBlockBlobClient(fileName)
	if err != nil {
		return "", nil, errors.New("Invalid New BlockBlob with error: " + err.Error())
	}

	context := context.TODO()
	get, err := blockBlob.Download(context, nil)
	if err != nil {
		return "", nil, errors.New(err.Error())
	}
	downloadedData := &bytes.Buffer{}
	reader := get.Body(&azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(reader)
	if err != nil {
		return "", nil, errors.New(err.Error())
	}
	err = reader.Close()
	if err != nil {
		return "", nil, errors.New(err.Error())
	}

	return fileName, downloadedData, nil
}
