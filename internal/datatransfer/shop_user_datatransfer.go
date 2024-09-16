package datatransfer

import (
	"context"
	shopModule "smlcloudplatform/internal/shop"
)

type ShopUserDataTransfer struct {
	transferConnection IDataTransferConnection
}

func NewShopUserDataTransfer(transferConnection IDataTransferConnection) IDataTransfer {
	return &ShopUserDataTransfer{
		transferConnection: transferConnection,
	}
}

func (sdt *ShopUserDataTransfer) StartTransfer(ctx context.Context, shopID string) error {

	shopSourceRepository := shopModule.NewShopUserRepository(sdt.transferConnection.GetSourceConnection())

	shopUserList, err := shopSourceRepository.FindByShopID(ctx, shopID)
	if err != nil {
		return err
	}

	shopTargetRepository := shopModule.NewShopUserRepository(sdt.transferConnection.GetTargetConnection())

	for _, shopUser := range *shopUserList {
		err := shopTargetRepository.Create(ctx, &shopUser)
		if err != nil {
			return err
		}
	}

	return nil

}
