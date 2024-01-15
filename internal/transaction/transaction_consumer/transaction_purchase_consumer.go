package transactionconsumer

import (
	"encoding/json"
	purchaseModels "smlcloudplatform/internal/transaction/purchase/models"
	"smlcloudplatform/pkg/microservice"
)

func (t *TransactionConsumer) ConsumeOnPurchaseDocCreated(ctx microservice.IContext) error {
	msg := ctx.ReadInput()

	purchaseDoc := purchaseModels.PurchaseDoc{}
	err := json.Unmarshal([]byte(msg), &purchaseDoc)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Unmarshal PurchaseDoc Message : %v", err.Error())
		return err
	}

	trx, err := t.phaser.PhasePurchaseDoc(&purchaseDoc)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Phase PurchaseDoc to StockTransaction : %v", err.Error())
		return err
	}

	err = t.svc.UpSert(trx.ShopID, trx.DocNo, *trx)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Insert StockTransaction : %v", err.Error())
		return err
	}
	return nil
}

func (t *TransactionConsumer) ConsumeOnPurchaseDocUpdated(ctx microservice.IContext) error {

	msg := ctx.ReadInput()

	purchaseDoc := purchaseModels.PurchaseDoc{}
	err := json.Unmarshal([]byte(msg), &purchaseDoc)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Unmarshal PurchaseDoc Message : %v", err.Error())
		return err
	}

	trx, err := t.phaser.PhasePurchaseDoc(&purchaseDoc)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Phase PurchaseDoc to StockTransaction : %v", err.Error())
		return err
	}

	err = t.svc.UpSert(trx.ShopID, trx.DocNo, *trx)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Insert StockTransaction : %v", err.Error())
		return err
	}
	return nil
}

func (t *TransactionConsumer) ConsumeOnPurchaseDocDeleted(ctx microservice.IContext) error {

	msg := ctx.ReadInput()

	purchaseDoc := purchaseModels.PurchaseDoc{}
	err := json.Unmarshal([]byte(msg), &purchaseDoc)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Unmarshal PurchaseDoc Message : %v", err.Error())
		return err
	}

	trx, err := t.phaser.PhasePurchaseDoc(&purchaseDoc)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Phase PurchaseDoc to StockTransaction : %v", err.Error())
		return err
	}

	err = t.svc.Delete(trx.ShopID, trx.DocNo)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Insert StockTransaction : %v", err.Error())
		return err
	}
	return nil
}
