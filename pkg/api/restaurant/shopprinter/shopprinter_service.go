package shopprinter

import (
	"errors"
	"fmt"
	"smlcloudplatform/pkg/models/restaurant"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/utils"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IShopPrinterService interface {
	CreateShopPrinter(shopID string, authUsername string, doc restaurant.PrinterTerminal) (string, error)
	UpdateShopPrinter(guid string, shopID string, authUsername string, doc restaurant.PrinterTerminal) error
	DeleteShopPrinter(guid string, shopID string, authUsername string) error
	InfoShopPrinter(guid string, shopID string) (restaurant.PrinterTerminalInfo, error)
	SearchShopPrinter(shopID string, q string, page int, limit int) ([]restaurant.PrinterTerminalInfo, mongopagination.PaginationData, error)
}

type ShopPrinterService struct {
	crudRepo   repositories.CrudRepository[restaurant.PrinterTerminalDoc]
	searchRepo repositories.SearchRepository[restaurant.PrinterTerminalInfo]
}

func NewShopPrinterService(crudRepo repositories.CrudRepository[restaurant.PrinterTerminalDoc], searchRepo repositories.SearchRepository[restaurant.PrinterTerminalInfo]) ShopPrinterService {
	return ShopPrinterService{
		crudRepo:   crudRepo,
		searchRepo: searchRepo,
	}
}

func (svc ShopPrinterService) CreateShopPrinter(shopID string, authUsername string, doc restaurant.PrinterTerminal) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := restaurant.PrinterTerminalDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.PrinterTerminal = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.crudRepo.Create(docData)

	if err != nil {
		return "", err
	}
	return newGuidFixed, nil
}

func (svc ShopPrinterService) UpdateShopPrinter(guid string, shopID string, authUsername string, doc restaurant.PrinterTerminal) error {

	findDoc, err := svc.crudRepo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.PrinterTerminal = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.crudRepo.Update(guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc ShopPrinterService) DeleteShopPrinter(guid string, shopID string, authUsername string) error {
	err := svc.crudRepo.Delete(shopID, guid, authUsername)

	if err != nil {
		return err
	}
	return nil
}

func (svc ShopPrinterService) InfoShopPrinter(guid string, shopID string) (restaurant.PrinterTerminalInfo, error) {

	findDoc, err := svc.crudRepo.FindByGuid(shopID, guid)

	if err != nil {
		return restaurant.PrinterTerminalInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return restaurant.PrinterTerminalInfo{}, errors.New("document not found")
	}

	return findDoc.PrinterTerminalInfo, nil

}

func (svc ShopPrinterService) SearchShopPrinter(shopID string, q string, page int, limit int) ([]restaurant.PrinterTerminalInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"code",
	}

	for i := range [5]bool{} {
		searchCols = append(searchCols, fmt.Sprintf("name%d", (i+1)))
	}

	docList, pagination, err := svc.searchRepo.FindPage(shopID, searchCols, q, page, limit)

	if err != nil {
		return []restaurant.PrinterTerminalInfo{}, pagination, err
	}

	return docList, pagination, nil
}
