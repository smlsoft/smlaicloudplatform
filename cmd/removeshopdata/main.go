package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	tf "smlaicloudplatform/internal/datatransfer"
	"smlaicloudplatform/internal/shop"
	"smlaicloudplatform/pkg/microservice"

	shopUserModels "smlaicloudplatform/internal/authentication/models"
	saleChannelModels "smlaicloudplatform/internal/channel/salechannel/models"
	transportChannelModels "smlaicloudplatform/internal/channel/transportchannel/models"
	creditorModels "smlaicloudplatform/internal/debtaccount/creditor/models"
	creditorGroupModels "smlaicloudplatform/internal/debtaccount/creditorgroup/models"
	debtorModels "smlaicloudplatform/internal/debtaccount/debtor/models"
	debtorGroupModels "smlaicloudplatform/internal/debtaccount/debtorgroup/models"
	orderDeviceModels "smlaicloudplatform/internal/order/device/models"
	orderSettingModels "smlaicloudplatform/internal/order/setting/models"
	organizationBranchModels "smlaicloudplatform/internal/organization/branch/models"
	organizationBusinessTypeModels "smlaicloudplatform/internal/organization/businesstype/models"
	organizationDepartmentModels "smlaicloudplatform/internal/organization/department/models"
	paymentBankMasterModels "smlaicloudplatform/internal/payment/bankmaster/models"
	paymentBookBankMasterModels "smlaicloudplatform/internal/payment/bookbank/models"
	paymentQRPaymentMasterModels "smlaicloudplatform/internal/payment/qrpayment/models"
	posMediaModels "smlaicloudplatform/internal/pos/media/models"
	posSettingModels "smlaicloudplatform/internal/pos/setting/models"
	productBarcodeBOMModels "smlaicloudplatform/internal/product/bom/models"
	orderTypeModels "smlaicloudplatform/internal/product/ordertype/models"
	productBarcodeModels "smlaicloudplatform/internal/product/productbarcode/models"
	productCategoryModels "smlaicloudplatform/internal/product/productcategory/models"
	productGroupModels "smlaicloudplatform/internal/product/productgroup/models"
	productUnitModels "smlaicloudplatform/internal/product/unit/models"
	restaurantKitchenModels "smlaicloudplatform/internal/restaurant/kitchen/models"
	restaurantSettingModels "smlaicloudplatform/internal/restaurant/settings/models"
	restaurantStaffModels "smlaicloudplatform/internal/restaurant/staff/models"
	restaurentTableModels "smlaicloudplatform/internal/restaurant/table/models"
	restaurantZoneModels "smlaicloudplatform/internal/restaurant/zone/models"
	shopEmployeeModels "smlaicloudplatform/internal/shop/employee/models"
	slipImageModels "smlaicloudplatform/internal/slipimage/models"
	transactionPaidModels "smlaicloudplatform/internal/transaction/paid/models"
	transactionPayModels "smlaicloudplatform/internal/transaction/pay/models"
	transactionPurchaseModels "smlaicloudplatform/internal/transaction/purchase/models"
	transactionPurcahseReturnModels "smlaicloudplatform/internal/transaction/purchasereturn/models"
	transactionSaleModels "smlaicloudplatform/internal/transaction/saleinvoice/models"
	transactionSaleBOMModels "smlaicloudplatform/internal/transaction/saleinvoicebomprice/models"
	transactionSaleInvoiceReturnModels "smlaicloudplatform/internal/transaction/saleinvoicereturn/models"
	transactionStockAdjustModels "smlaicloudplatform/internal/transaction/stockadjustment/models"
	transactionStockBalanceModels "smlaicloudplatform/internal/transaction/stockbalance/models"
	transactionStockBalanceDetailModels "smlaicloudplatform/internal/transaction/stockbalancedetail/models"
	transactionStockPickupModels "smlaicloudplatform/internal/transaction/stockpickupproduct/models"
	transactionStockReceiveModels "smlaicloudplatform/internal/transaction/stockreceiveproduct/models"
	transactionStockReturnModels "smlaicloudplatform/internal/transaction/stockreturnproduct/models"
	transactionStockTransferModels "smlaicloudplatform/internal/transaction/stocktransfer/models"
	warehouseModels "smlaicloudplatform/internal/warehouse/models"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	shopID = flag.String("shopid", "", "shopID to transfer")
)

func main() {

	godotenv.Load()
	flag.Parse()

	destinationDBConfig := tf.DestinationDatabaseConfig{}
	targetDatabase := microservice.NewPersisterMongo(destinationDBConfig)
	ctx := context.TODO()

	shopRepo := shop.NewShopRepository(targetDatabase)
	findShop, err := shopRepo.FindByGuid(ctx, *shopID)
	if err != nil {
		fmt.Println("Find Shop Error")
		return
	}

	if findShop.GuidFixed == "" {
		fmt.Println("Shop not found")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Are you sure to Remove Shop", *findShop.Names[0].Name, "-", findShop.GuidFixed, "by", findShop.CreatedBy, " ? (y/n)")
	text, _ := reader.ReadString('\n')

	if text != "y\n" {
		fmt.Println("Transfer is cancelled")
		return
	}

	// remove data in shop user {"shopid": "2PrIIqTWxoBXv16K310sNwfHmfY", "username": {"$ne": "maxkorn"}}
	fmt.Println("Remove Shop User")
	err = targetDatabase.Delete(ctx, &shopUserModels.ShopUser{}, bson.M{
		"shopid":   shopID,
		"username": bson.M{"$ne": findShop.CreatedBy},
	})
	if err != nil {
		fmt.Println("Remove Shop User Error", err)
		return
	}

	// remove shop employee
	fmt.Println("Remove Shop Employee")
	err = targetDatabase.Delete(ctx, &shopEmployeeModels.EmployeeDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Shop Employee Error", err)
		return
	}

	// remove product barcode
	fmt.Println("Remove Product Barcode")
	err = targetDatabase.Delete(ctx, &productBarcodeModels.ProductBarcodeDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Product Barcode Error", err)
		return
	}

	// remove product categories
	fmt.Println("Remove Product Categories")
	err = targetDatabase.Delete(ctx, &productCategoryModels.ProductCategoryDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Product Categories Error", err)
		return
	}

	// remove kitchen
	fmt.Println("Remove Product Kitchen")
	err = targetDatabase.Delete(ctx, &restaurantKitchenModels.KitchenDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Product Kitchen Error", err)
		return
	}

	// remove restaurant setting
	fmt.Println("Remove Restaurant Setting")
	err = targetDatabase.Delete(ctx, &restaurantSettingModels.RestaurantSettingsDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Restaurant Setting Error", err)
		return
	}

	// remove bank master
	fmt.Println("Remove Bank Master")
	err = targetDatabase.Delete(ctx, &paymentBankMasterModels.BankMasterDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Bank Master Error", err)
		return
	}

	// remove bookbank master
	fmt.Println("Remove Book Bank Master")
	err = targetDatabase.Delete(ctx, &paymentBookBankMasterModels.BookBankDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Book Bank Master Error", err)
		return
	}

	// remove qrpayment master
	fmt.Println("Remove QR Payment Master")
	err = targetDatabase.Delete(ctx, &paymentQRPaymentMasterModels.QrPaymentDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove QR Payment Master Error", err)
		return
	}

	// remove order device
	fmt.Println("Remove Order Device")
	err = targetDatabase.Delete(ctx, &orderDeviceModels.OrderDeviceDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Device Error", err)
		return
	}

	// remove order setting
	fmt.Println("Remove Order Setting")
	err = targetDatabase.Delete(ctx, &orderSettingModels.SettingDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// organization branch
	fmt.Println("Remove Organization Branch")
	err = targetDatabase.Delete(ctx, &organizationBranchModels.BranchDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// organization business type
	fmt.Println("Remove Organization Business Type")
	err = targetDatabase.Delete(ctx, &organizationBusinessTypeModels.BusinessTypeDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// organization department
	fmt.Println("Remove Organization Department")
	err = targetDatabase.Delete(ctx, &organizationDepartmentModels.DepartmentDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// pos media
	fmt.Println("Remove POS Media")
	err = targetDatabase.Delete(ctx, &posMediaModels.MediaDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// pos setting
	fmt.Println("Remove POS Setting")
	err = targetDatabase.Delete(ctx, &posSettingModels.SettingDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// product barcode bom
	fmt.Println("Remove Product Barcode BOM")
	err = targetDatabase.Delete(ctx, &productBarcodeBOMModels.ProductBarcodeBOMViewDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// product group
	fmt.Println("Remove Product Group")
	err = targetDatabase.Delete(ctx, &productGroupModels.ProductGroupDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// product unit
	fmt.Println("Remove Product Unit")
	err = targetDatabase.Delete(ctx, &productUnitModels.UnitDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// order type
	fmt.Println("Remove Order Type")
	err = targetDatabase.Delete(ctx, &orderTypeModels.OrderTypeDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove restaurant table
	fmt.Println("Remove Restaurant Table")
	err = targetDatabase.Delete(ctx, &restaurentTableModels.TableDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove restaurant zone
	fmt.Println("Remove Restaurant Zone")
	err = targetDatabase.Delete(ctx, &restaurantZoneModels.ZoneDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove restaurant staff
	fmt.Println("Remove Restaurant Staff")
	err = targetDatabase.Delete(ctx, &restaurantStaffModels.StaffDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove sale channel
	fmt.Println("Remove Sale Channel")
	err = targetDatabase.Delete(ctx, &saleChannelModels.SaleChannelDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove transport channel
	fmt.Println("Remove Transport Channel")
	err = targetDatabase.Delete(ctx, &transportChannelModels.TransportChannelDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove slip image
	fmt.Println("Remove Slip Image")
	err = targetDatabase.Delete(ctx, &slipImageModels.SlipImageDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove debtor
	fmt.Println("Remove Debtor")
	err = targetDatabase.Delete(ctx, &debtorModels.DebtorDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove debtor group
	fmt.Println("Remove Debtor Group")
	err = targetDatabase.Delete(ctx, &debtorGroupModels.DebtorGroupDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove creditor
	fmt.Println("Remove Creditor")
	err = targetDatabase.Delete(ctx, &creditorModels.CreditorDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove creditor group
	fmt.Println("Remove Creditor Group")
	err = targetDatabase.Delete(ctx, &creditorGroupModels.CreditorGroupDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove warehouse
	fmt.Println("Remove Warehouse")
	err = targetDatabase.Delete(ctx, &warehouseModels.WarehouseDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove transaction paid
	fmt.Println("Remove Transaction Paid")
	err = targetDatabase.Delete(ctx, &transactionPaidModels.PaidDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove transaction pay
	fmt.Println("Remove Transaction Pay")
	err = targetDatabase.Delete(ctx, &transactionPayModels.PayDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove transaction purchase
	fmt.Println("Remove Transaction Purchase")
	err = targetDatabase.Delete(ctx, &transactionPurchaseModels.PurchaseDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove transaction purchase return
	fmt.Println("Remove Transaction Purchase Return")
	err = targetDatabase.Delete(ctx, &transactionPurcahseReturnModels.PurchaseReturnDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove transaction sale
	fmt.Println("Remove Transaction Sale Invoice")
	err = targetDatabase.Delete(ctx, &transactionSaleModels.SaleInvoiceDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove transaction sale bom
	fmt.Println("Remove Transaction Sale BOM")
	err = targetDatabase.Delete(ctx, &transactionSaleBOMModels.SaleInvoiceBomPriceDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove transaction sale return
	fmt.Println("Remove Transaction Sale Return")
	err = targetDatabase.Delete(ctx, &transactionSaleInvoiceReturnModels.SaleInvoiceReturnDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove transaction stock balance
	fmt.Println("Remove Transaction Stock Balance")
	err = targetDatabase.Delete(ctx, &transactionStockBalanceModels.StockBalanceDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove transaction stock balance detail
	fmt.Println("Remove Transaction Stock Balance Detail")
	err = targetDatabase.Delete(ctx, &transactionStockBalanceDetailModels.StockBalanceDetailDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove transaction stock adjust
	fmt.Println("Remove Transaction Stock Adjust")
	err = targetDatabase.Delete(ctx, &transactionStockAdjustModels.StockAdjustmentDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove transaction stock pickup
	fmt.Println("Remove Transaction Stock Pickup")
	err = targetDatabase.Delete(ctx, &transactionStockPickupModels.StockPickupProductDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove transaction stock receive
	fmt.Println("Remove Transaction Stock Receive")
	err = targetDatabase.Delete(ctx, &transactionStockReceiveModels.StockReceiveProductDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove transaction stock return
	fmt.Println("Remove Transaction Stock Return")
	err = targetDatabase.Delete(ctx, &transactionStockReturnModels.StockReturnProductDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

	// remove transaction stock transfer
	fmt.Println("Remove Transaction Stock Transfer")
	err = targetDatabase.Delete(ctx, &transactionStockTransferModels.StockTransferDoc{}, bson.M{
		"shopid": shopID,
	})
	if err != nil {
		fmt.Println("Remove Order Setting Error", err)
		return
	}

}
