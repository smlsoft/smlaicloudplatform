package datatransfer

import (
	"context"
	"smlcloudplatform/pkg/microservice"
)

type DBTransfer struct {
	sourceDatabase microservice.IPersisterMongo
	targetDatabase microservice.IPersisterMongo
}

func NewDBTransfer(sourceDatabase microservice.IPersisterMongo, targetDatabase microservice.IPersisterMongo) IDBTransfer {

	return &DBTransfer{
		sourceDatabase: sourceDatabase,
		targetDatabase: targetDatabase,
	}
}

func (db *DBTransfer) BeginTransfer(shopID string) {

	connection := NewDataTransferConnection(db.sourceDatabase, db.targetDatabase)

	_, err := connection.TestConnect()

	// start transfer shop
	todo := context.TODO()

	// // shop
	// shopDataTransfer := NewShopDataTransfer(connection)
	// err = shopDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // shopuser
	// shopUserDataTransfer := NewShopUserDataTransfer(connection)
	// err = shopUserDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // productbarcode
	// productBarcodeDataTransfer := NewProductBarcodeDataTransfer(connection)
	// err = productBarcodeDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // category
	// productCategoryTransfer := NewProductCategoryDataTransfer(connection)
	// err = productCategoryTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // kitchen
	// kitchenDataTransfer := NewRestaurantKitchenDataTransfer(connection)
	// err = kitchenDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // restaurant setting
	// restaurantSettingDataTransfer := NewRestaurantSettingDataTransfer(connection)
	// err = restaurantSettingDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // bank master
	// bankMasterDataTransfer := NewBankMasterDataTransfer(connection)
	// err = bankMasterDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // book bank
	// bookbankDataTransfer := NewBookBankDataTransfer(connection)
	// err = bookbankDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // order device
	// orderDeviceDataTransfer := NewOrderDeviceDataTransfer(connection)
	// err = orderDeviceDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // order device setting
	// orderDeviceSettingDataTransfer := NewOrderDeviceSettingDataTransfer(connection)
	// err = orderDeviceSettingDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // organization branch
	// organizationBranchDataTransfer := NewOrganizationBranchDataTransfer(connection)
	// err = organizationBranchDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // organization business type
	// organizationBusinessTypeDataTransfer := NewOrganizationBusinessTypeDataTransfer(connection)
	// err = organizationBusinessTypeDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // organization department
	// organizationDepartmentDataTransfer := NewOrganizationDepartmentDataTransfer(connection)
	// err = organizationDepartmentDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // pos media
	// posMediaDataTransfer := NewPosMediaDataTransfer(connection)
	// err = posMediaDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // pos setting
	// posSettingDataTransfer := NewPosSettingDataTransfer(connection)
	// err = posSettingDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // product barcode bom
	// productbarcodeBOMDataTransfer := NewProductbarcodeBOMDataTransfer(connection)
	// err = productbarcodeBOMDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // product group
	// productGroupDataTransfer := NewProductGroupDataTransfer(connection)
	// err = productGroupDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// product unit
	// productUnitDataTransfer := NewProductUnitDataTransfer(connection)
	// err = productUnitDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // order type
	// orderTypeDataTransfer := NewOrderTypeDataTransfer(connection)
	// err = orderTypeDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// restaurant table
	// restaurantTableDataTransfer := NewRestaurantTableDataTransfer(connection)
	// err = restaurantTableDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // restaurant zone
	// restaurantZoneDataTransfer := NewRestaurantZoneDataTransfer(connection)
	// err = restaurantZoneDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // sale channel
	// saleChannelDataTransfer := NewSaleChannelDataTransfer(connection)
	// err = saleChannelDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // transport channel
	// transportChannelDataTransfer := NewSaleTransportDataTransfer(connection)
	// err = transportChannelDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // slip image
	// slipImageDataTransfer := NewSlipImageDataTransfer(connection)
	// err = slipImageDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // warehouse
	// warehouseDataTransfer := NewProductWarehouseDataTransfer(connection)
	// err = warehouseDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // transaction paid
	// transactionPaidDataTransfer := NewTransactionPaidDataTransfer(connection)
	// err = transactionPaidDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // transaction pay
	// transactionPayDataTransfer := NewTransactionPayDataTransfer(connection)
	// err = transactionPayDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // transaction purchase
	// transactionPurchaseDataTransfer := NewTransactionPurchaseDataTransfer(connection)
	// err = transactionPurchaseDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // transaction purchase return
	// purchaseReturnDataTransfer := NewPurchaseReturnDataTransfer(connection)
	// err = purchaseReturnDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // transaction sale
	// transactionSaleDataTransfer := NewSaleInvoiceDataTransfer(connection)
	// err = transactionSaleDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // transaction sale invoice bom price
	// transactionSaleInvoiceBomPriceDataTransfer := NewSaleInvoiceBomPricesDataTransfer(connection)
	// err = transactionSaleInvoiceBomPriceDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // sale invoice return
	// saleInvoiceReturnDataTransfer := NewSaleInvoiceReturnDataTransfer(connection)
	// err = saleInvoiceReturnDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // transaction stock balance
	// transactionStockBalanceDataTransfer := NewStockBalanceDataTransfer(connection)
	// err = transactionStockBalanceDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // transaction stock balance detail
	transactionStockBalanceDetailDataTransfer := NewStockBalanceDetailDataTransfer(connection)
	err = transactionStockBalanceDetailDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// // transaction stock adjustment
	// transactionStockAdjustmentDataTransfer := NewStockAdjustmentDataTransfer(connection)
	// err = transactionStockAdjustmentDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // transaction stock pickup
	// transactionStockPickupDataTransfer := NewStockPickupProductDataTransfer(connection)
	// err = transactionStockPickupDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // transaction stock receive
	// transactionStockReceiveDataTransfer := NewStockReceiveProductDataTransfer(connection)
	// err = transactionStockReceiveDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // transaction stock return
	// transactionStockReturnDataTransfer := NewStockReturnProductDataTransfer(connection)
	// err = transactionStockReturnDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

	// // transaction stock transfer
	// transactionStockTransferDataTransfer := NewStockTransferDataTransfer(connection)
	// err = transactionStockTransferDataTransfer.StartTransfer(todo, shopID)
	// if err != nil {
	// 	panic(err)
	// }

}
