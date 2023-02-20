// Copyright (C) 2022 Le Hoang Long
// This file is part of SigGraph smart contract <https://github.com/LeHoangLong/sig_graph_smart_contract>.
// 
// SigGraph is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// 
// SigGraph is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with SigGraph.  If not, see <http://www.gnu.org/licenses/>.
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

type transefrAssetRequest struct {
	TimeMs uint64 `json:"time_ms"`

	CurrentId        string `json:"current_id"`
	CurrentSignature string `json:"current_signature"`
	CurrentSecret    string `json:"current_secret"`

	NewId        string `json:"new_id"`
	NewSignature string `json:"new_signature"`
	NewSecret    string `json:"new_secret"`

	NewOwnerPublicKey string `json:"new_owner_public_key"`
}

func (v *assetView) TransferAsset(
	transaction contractapi.TransactionContextInterface,
	requestJson string,
) (*asset, error) {
	ctx := context.Background()
	service := service.NewSmartContractServiceHyperledger(transaction)

	request := transefrAssetRequest{}
	err := json.Unmarshal([]byte(requestJson), &request)
	if err != nil {
		return nil, err
	}

	newAsset, err := v.controller.TransferAsset(
		ctx,
		service,
		request.TimeMs,
		request.CurrentId,
		request.CurrentSecret,
		request.CurrentSignature,
		request.NewId,
		request.NewSecret,
		request.NewSignature,
		request.NewOwnerPublicKey,
	)

	if err != nil {
		return nil, err
	}

	viewAsset := FromModelAsset(newAsset)
	return &viewAsset, nil
}
