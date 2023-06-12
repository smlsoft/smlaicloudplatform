package productbarcodeconsumer

import "smlcloudplatform/internal/microservice"

type ProductBarcodeConsumer interface {
	RegisterConsumer(ms *microservice.Microservice)
}
