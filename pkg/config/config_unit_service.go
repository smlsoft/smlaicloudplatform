package config

import "os"

type IUnitServiceConfig interface {
	ProductHost() string
}

type UnitServiceConfig struct{}

func NewUnitServiceConfig() *UnitServiceConfig {
	return &UnitServiceConfig{}
}

func (cfg *UnitServiceConfig) ProductHost() string {
	return os.Getenv("UNIT_SERVICE_PRODUCT_HOST")
}
