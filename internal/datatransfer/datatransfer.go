package datatransfer

import "context"

type IDataTransfer interface {
	StartTransfer(ctx context.Context, shopID string) error
}

type IDBTransfer interface {
	BeginTransfer(shopID string)
}
