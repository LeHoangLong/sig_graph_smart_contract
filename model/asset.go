package model

import "github.com/shopspring/decimal"

type ECreationProcess = string

const (
	ECreationProcessCreate   = "create"
	ECreationProcessTransfer = "transfer"
)

type Asset struct {
	Node            `mapstructure:",squash"`
	CreationProcess ECreationProcess `json:"creation_process" mapstructure:"creation_process"`
	Unit            string           `json:"unit"  mapstructure:"unit"`
	Quantity        decimal.Decimal  `json:"quantity"  mapstructure:"quantity"`
	MaterialName    string           `json:"material_name"  mapstructure:"material_name"`
}
