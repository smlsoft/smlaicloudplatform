package datatransfer

import (
	"context"
	"fmt"
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
	if err != nil {
		panic(err)
	}

	fmt.Println("Start transfer data")

	// start transfer shop
	todo := context.TODO()

	// shop
	fmt.Println("Start transfer shop")
	shopDataTransfer := NewShopDataTransfer(connection)
	err = shopDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// shopuser
	fmt.Println("Start transfer shop user")
	shopUserDataTransfer := NewShopUserDataTransfer(connection)
	err = shopUserDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// shop employee
	fmt.Println("Start transfer shop employee")
	shopEmployeeDataTransfer := NewShopEmployeeDataTransfer(connection)
	err = shopEmployeeDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// productbarcode
	fmt.Println("Start transfer Product Barcode")
	productBarcodeDataTransfer := NewProductBarcodeDataTransfer(connection)
	err = productBarcodeDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// category
	fmt.Println("Start transfer Product Category")
	productCategoryTransfer := NewProductCategoryDataTransfer(connection)
	err = productCategoryTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// kitchen
	fmt.Println("Start transfer Restaurant Kitchen")
	kitchenDataTransfer := NewRestaurantKitchenDataTransfer(connection)
	err = kitchenDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// restaurant setting
	fmt.Println("Start transfer Restaurant Setting")
	restaurantSettingDataTransfer := NewRestaurantSettingDataTransfer(connection)
	err = restaurantSettingDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// bank master
	fmt.Println("Start transfer Bank Master")
	bankMasterDataTransfer := NewBankMasterDataTransfer(connection)
	err = bankMasterDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// book bank
	fmt.Println("Start transfer Book Bank")
	bookbankDataTransfer := NewBookBankDataTransfer(connection)
	err = bookbankDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// qr payment
	fmt.Println("Start transfer QR Payment")
	qrPaymentDataTransfer := NewQRPaymentDataTransfer(connection)
	err = qrPaymentDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// order device
	fmt.Println("Start transfer Order Device")
	orderDeviceDataTransfer := NewOrderDeviceDataTransfer(connection)
	err = orderDeviceDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// order device setting
	fmt.Println("Start transfer Order Device Setting")
	orderDeviceSettingDataTransfer := NewOrderDeviceSettingDataTransfer(connection)
	err = orderDeviceSettingDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// organization branch
	fmt.Println("Start transfer Organization Branch")
	organizationBranchDataTransfer := NewOrganizationBranchDataTransfer(connection)
	err = organizationBranchDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// organization business type
	fmt.Println("Start transfer Organization Business Type")
	organizationBusinessTypeDataTransfer := NewOrganizationBusinessTypeDataTransfer(connection)
	err = organizationBusinessTypeDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// organization department
	fmt.Println("Start transfer Organization Department")
	organizationDepartmentDataTransfer := NewOrganizationDepartmentDataTransfer(connection)
	err = organizationDepartmentDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// pos media
	fmt.Println("Start transfer Pos Media")
	posMediaDataTransfer := NewPosMediaDataTransfer(connection)
	err = posMediaDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// pos setting
	fmt.Println("Start transfer Pos Setting")
	posSettingDataTransfer := NewPosSettingDataTransfer(connection)
	err = posSettingDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// product barcode bom
	fmt.Println("Start transfer Product Barcode BOM")
	productbarcodeBOMDataTransfer := NewProductbarcodeBOMDataTransfer(connection)
	err = productbarcodeBOMDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// product group
	fmt.Println("Start transfer Product Group")
	productGroupDataTransfer := NewProductGroupDataTransfer(connection)
	err = productGroupDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// product unit
	fmt.Println("Start transfer Product Unit")
	productUnitDataTransfer := NewProductUnitDataTransfer(connection)
	err = productUnitDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// order type
	fmt.Println("Start transfer Order Type")
	orderTypeDataTransfer := NewOrderTypeDataTransfer(connection)
	err = orderTypeDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// restaurant table
	fmt.Println("Start transfer Restaurant Table")
	restaurantTableDataTransfer := NewRestaurantTableDataTransfer(connection)
	err = restaurantTableDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// restaurant zone
	fmt.Println("Start transfer Restaurant Zone")
	restaurantZoneDataTransfer := NewRestaurantZoneDataTransfer(connection)
	err = restaurantZoneDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// restaurant staff
	fmt.Println("Start transfer Restaurant Staff")
	restaurantStaffDataTransfer := NewRestaurantStaffDataTransfer(connection)
	err = restaurantStaffDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// sale channel
	fmt.Println("Start transfer Sale Channel")
	saleChannelDataTransfer := NewSaleChannelDataTransfer(connection)
	err = saleChannelDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// transport channel
	fmt.Println("Start transfer Transport Channel")
	transportChannelDataTransfer := NewSaleTransportDataTransfer(connection)
	err = transportChannelDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// slip image
	slipImageDataTransfer := NewSlipImageDataTransfer(connection)
	err = slipImageDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// debtor
	fmt.Println("Start transfer Debtor")
	debtorDataTransfer := NewDebtorDataTransfer(connection)
	err = debtorDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// debor group
	fmt.Println("Start transfer Debtor Group")
	debtorGroupDataTransfer := NewDebtorGroupDataTransfer(connection)
	err = debtorGroupDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// creditor
	fmt.Println("Start transfer Creditor")
	creditorDataTransfer := NewCreditorDataTransfer(connection)
	err = creditorDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// creditor group
	fmt.Println("Start transfer Creditor Group")
	creditorGroupDataTransfer := NewCreditorGroupDataTransfer(connection)
	err = creditorGroupDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// warehouse
	fmt.Println("Start transfer Warehouse")
	warehouseDataTransfer := NewProductWarehouseDataTransfer(connection)
	err = warehouseDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// transaction paid
	transactionPaidDataTransfer := NewTransactionPaidDataTransfer(connection)
	err = transactionPaidDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// transaction pay
	transactionPayDataTransfer := NewTransactionPayDataTransfer(connection)
	err = transactionPayDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// transaction purchase
	transactionPurchaseDataTransfer := NewTransactionPurchaseDataTransfer(connection)
	err = transactionPurchaseDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// transaction purchase return
	purchaseReturnDataTransfer := NewPurchaseReturnDataTransfer(connection)
	err = purchaseReturnDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// transaction sale
	transactionSaleDataTransfer := NewSaleInvoiceDataTransfer(connection)
	err = transactionSaleDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// transaction sale invoice bom price
	transactionSaleInvoiceBomPriceDataTransfer := NewSaleInvoiceBomPricesDataTransfer(connection)
	err = transactionSaleInvoiceBomPriceDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// sale invoice return
	saleInvoiceReturnDataTransfer := NewSaleInvoiceReturnDataTransfer(connection)
	err = saleInvoiceReturnDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// transaction stock balance
	transactionStockBalanceDataTransfer := NewStockBalanceDataTransfer(connection)
	err = transactionStockBalanceDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// transaction stock balance detail
	transactionStockBalanceDetailDataTransfer := NewStockBalanceDetailDataTransfer(connection)
	err = transactionStockBalanceDetailDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// transaction stock adjustment
	transactionStockAdjustmentDataTransfer := NewStockAdjustmentDataTransfer(connection)
	err = transactionStockAdjustmentDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// transaction stock pickup
	transactionStockPickupDataTransfer := NewStockPickupProductDataTransfer(connection)
	err = transactionStockPickupDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// transaction stock receive
	transactionStockReceiveDataTransfer := NewStockReceiveProductDataTransfer(connection)
	err = transactionStockReceiveDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// transaction stock return
	transactionStockReturnDataTransfer := NewStockReturnProductDataTransfer(connection)
	err = transactionStockReturnDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

	// transaction stock transfer
	transactionStockTransferDataTransfer := NewStockTransferDataTransfer(connection)
	err = transactionStockTransferDataTransfer.StartTransfer(todo, shopID)
	if err != nil {
		panic(err)
	}

}
