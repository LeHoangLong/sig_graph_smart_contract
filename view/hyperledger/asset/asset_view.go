package asset_hyperledger_view

import (
	"context"
	"encoding/json"
	asset_controller "sig_graph/controller/asset"
	"sig_graph/model"
	"sig_graph/service"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/shopspring/decimal"
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

type createAssetRequest struct {
	Time                 uint64   `json:"time"`
	Id                   string   `json:"id"`
	MaterialName         string   `json:"material_name"`
	Quantity             string   `json:"quantity"`
	Unit                 string   `json:"unit"`
	Signature            string   `json:"signature"`
	OwnerPublicKey       string   `json:"owner_public_key"`
	IngredientIds        []string `json:"ingredient_ids"`
	IngredientSecretIds  []string `json:"ingredient_secret_ids"`
	SecretIds            []string `json:"secret_ids"`
	IngredientSignatures []string `json:"ingredient_signatures"`
}

type asset struct {
	model.Node
	CreationProcess string `json:"creation_process"`
	Unit            string `json:"unit"`
	Quantity        string `json:"quantity"`
	MaterialName    string `json:"material_name"`
}

func FromModelAsset(modelAsset *model.Asset) asset {
	return asset{
		Node:            modelAsset.Node,
		CreationProcess: modelAsset.CreationProcess,
		Unit:            modelAsset.Unit,
		Quantity:        modelAsset.Quantity.String(),
		MaterialName:    modelAsset.MaterialName,
	}
}

func (c *assetView) CreateAsset(
	transaction contractapi.TransactionContextInterface,
	requestJson string,
) (*asset, error) {
	request := createAssetRequest{}
	err := json.Unmarshal([]byte(requestJson), &request)
	if err != nil {
		return nil, err
	}

	quantityDec, err := decimal.NewFromString(request.Quantity)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	service := service.NewSmartContractServiceHyperledger(transaction)
	modelAsset, err := c.controller.CreateAsset(
		ctx,
		service,
		request.Time,
		request.Id,
		request.MaterialName,
		quantityDec,
		request.Unit,
		request.Signature,
		request.OwnerPublicKey,
		request.IngredientIds,
		request.IngredientSecretIds,
		request.SecretIds,
		request.IngredientSignatures,
	)

	if err != nil {
		return nil, err
	}

	asset := FromModelAsset(modelAsset)
	return &asset, nil
}

func (c *assetView) GetAsset(
	transaction contractapi.TransactionContextInterface,
	id string,
) (*asset, error) {
	ctx := context.Background()
	service := service.NewSmartContractServiceHyperledger(transaction)
	modelAsset, err := c.controller.GetAsset(ctx, service, id)
	if err != nil {
		return nil, err
	}

	viewAsset := FromModelAsset(modelAsset)
	return &viewAsset, nil
}
