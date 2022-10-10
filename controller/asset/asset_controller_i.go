package asset_controller

import (
	"context"
	"sig_graph/controller"
	"sig_graph/model"
)

//go:generate mockgen -source=$GOFILE -destination ../../testutils/asset_controller.go -package mock
type AssetControllerI interface {
	// return ErrAlreadyExists f id already used
	CreateMaterial(ctx context.Context, smartContract controller.SmartContractServiceI, time uint64, id string, signature string, ownerPublicKey string) error
	// return ErrNotFound if no material with id
	GetMaterial(ctx context.Context, smartContract controller.SmartContractServiceI, id string) (*model.Asset, error)
}
