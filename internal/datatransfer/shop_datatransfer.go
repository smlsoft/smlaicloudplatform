package datatransfer

import (
	"context"
	"fmt"
	shopModule "smlcloudplatform/internal/shop"
)

type ShopDataTransfer struct {
	transferConnection IDataTransferConnection
}

func NewShopDataTransfer(transferConnection IDataTransferConnection) IDataTransfer {
	return &ShopDataTransfer{
		transferConnection: transferConnection,
	}
}

func (sdt *ShopDataTransfer) StartTransfer(ctx context.Context, shopID string) error {

	shopSourceRepository := shopModule.NewShopRepository(sdt.transferConnection.GetSourceConnection())

	showDoc, err := shopSourceRepository.FindByGuid(ctx, shopID)
	if err != nil {
		return err
	}

	shopTargetRepository := shopModule.NewShopRepository(sdt.transferConnection.GetTargetConnection())
	id, err := shopTargetRepository.Create(ctx, showDoc)
	if err != nil {
		return err
	}

	fmt.Println("Shop ID: ", id)
	return nil

}
