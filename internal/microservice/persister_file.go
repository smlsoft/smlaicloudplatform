package microservice

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
)

type IPersisterFile interface {
	Save(file File, fileName string) error
}

type PersisterFile struct {
	StoreFilePath string
	StoreDataUri  string
}

type File struct {
	FileName    string
	ContentType string
	Data        []byte
	Size        int
}

func NewPersisterFile(cfg *StorageFileConfig) *PersisterFile {
	return &PersisterFile{
		StoreFilePath: cfg.StorageDataPath(),
		StoreDataUri:  cfg.StorageUriAtlas(),
	}
}

func (pst *PersisterFile) Save(fh *multipart.FileHeader, fileName string, fileExtension string) error {

	imageFileName := fmt.Sprintf("%s.%s", fileName, fileExtension)
	file, err := fh.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	// upload Now
	tempFile, err := ioutil.TempFile(pst.StoreFilePath, "upload-*."+fileExtension)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)

	// rename
	tmpFileName := tempFile.Name()
	tempFile.Close()

	if fileName != "" {
		//fmt.Println(fileName)
		// dir, _ := filepath.Split(pst.StoreFilePath)
		uploadFileName := filepath.Join(pst.StoreFilePath, imageFileName)
		//fmt.Println(uploadFileName)
		_ = os.Rename(tmpFileName, uploadFileName)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }
		fileName = uploadFileName
	}

	return nil

	//storeFilePath := filepath.Join(pst.StoreFilePath, fileName)
	//return ioutil.WriteFile(storeFilePath, file.Data, 0600)
}
