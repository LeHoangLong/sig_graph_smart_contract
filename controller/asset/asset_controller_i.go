package asset_controller

import (
	"context"
	"sig_graph/controller"
	"sig_graph/model"

	"github.com/shopspring/decimal"
)

//go:generate mockgen -source=$GOFILE -destination ../../testutils/asset_controller.go -package mock
type AssetControllerI interface {
	// return ErrAlreadyExists f id already used
	CreateAsset(
		ctx context.Context,
		smartContract controller.SmartContractServiceI,
		time_ms uint64,
		id string,
		materialName string,
		quantity decimal.Decimal,
		unit string,
		signature string,
		ownerPublicKey string,
		ingredientIds []string,
		ingredientSecretIds []string,
		secretIds []string,
		ingredientSignatures []string,
	) (*model.Asset, error)
	// return ErrNotFound if no material with id
	GetAsset(ctx context.Context, smartContract controller.SmartContractServiceI, id string) (*model.Asset, error)
	// return ErrNotFound if no material with currentId
	// return ErrAlreadyExists if newId already used
	// return new transferred asset
	TransferAsset(
		ctx context.Context,
		smartContract controller.SmartContractServiceI,
		time_ms uint64,
		currentId string,
		currentSecret string,
		currentSignature string,
		newId string,
		newSecret string,
		newSignature string,
		newOwnerPublicKey string,
	) (*model.Asset, error)
}
