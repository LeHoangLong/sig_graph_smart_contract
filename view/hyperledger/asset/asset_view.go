package asset_hyperledger_view

import (
	"context"
	asset_controller "sig_graph/controller/asset"
	"sig_graph/model"
	"sig_graph/service"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type assetView struct {
	contractapi.Contract
	controller asset_controller.AssetControllerI
}

func NewAssetView(
	controller asset_controller.AssetControllerI,
) *assetView {
	return &assetView{
		controller: controller,
	}
}

func (c *assetView) CreateMaterial(
	transaction contractapi.TransactionContextInterface,
	time uint64,
	id string,
	signature string,
	ownerPublicKey string,
) error {
	ctx := context.Background()
	service := service.NewSmartContractServiceHyperledger(transaction)
	return c.controller.CreateMaterial(ctx, service, time, id, signature, ownerPublicKey)
}

func (c *assetView) GetMaterial(
	ctx contractapi.TransactionContextInterface,
	id string,
) (*model.Asset, error) {
	return nil, nil
}
