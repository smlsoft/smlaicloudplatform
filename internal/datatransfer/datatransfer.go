package datatransfer

import "context"

type IDataTransfer interface {
	StartTransfer(ctx context.Context, shopID string) error
	// CheckingBeforeTransfer(ctx context.Context, shopID string) (bool, error)
}

type IDBTransfer interface {
	BeginTransfer(shopID string)
}
