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
	"fmt"
	asset_controller "sig_graph/controller/asset"
	"sig_graph/encrypt"
	"sig_graph/model"
	"sig_graph/service"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type assetView struct {
	contractapi.Contract
	controller asset_controller.AssetController
}

func NewAssetView(
	controller asset_controller.AssetController,
) *assetView {
	return &assetView{
		controller: controller,
	}
}

type createAssetRequest struct {
	CreationProcess encrypt.ToBeEncrypted[model.ECreationProcess] `json:"creation_process"`
	TransactionTime encrypt.ToBeEncrypted[int64]                  `json:"transaction_time"`
	Asset           model.NodeAsset                               `json:"asset"`
}

type asset struct {
	CreationProcess string `json:"creation_process"`
	Unit            string `json:"unit"`
	Quantity        string `json:"quantity"`
	MaterialName    string `json:"material_name"`
}

func (c *assetView) CreateAsset(
	Transaction contractapi.TransactionContextInterface,
	RequestJson string,
) (string, error) {
	request := createAssetRequest{}
	err := json.Unmarshal([]byte(RequestJson), &request)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	service := service.NewSmartContractServiceHyperledger(Transaction)
	modelAsset, stackErr := c.controller.CreateAsset(
		ctx,
		service,
		&request.TransactionTime,
		&request.CreationProcess,
		&request.Asset,
	)
	if stackErr != nil {
		return "", fmt.Errorf(stackErr.String())
	}

	if err != nil {
		return "", err
	}

	assetJson, err := json.Marshal(modelAsset)
	if err != nil {
		return "", err
	}

	return string(assetJson), nil
}

type transefrAssetRequest struct {
	TimeMs uint64 `json:"time_ms"`

	CurrentId                 encrypt.ToBeEncrypted[string] `json:"current_id"`
	CurrentSignature          string                        `json:"current_signature"`
	CurrentTransactionTime_ms encrypt.ToBeEncrypted[int64]  `json:"current_transaction_time_ms"`
	NewId                     encrypt.ToBeEncrypted[string]
	NewSignature              string
	NewOwnerPublicKey         encrypt.ToBeEncrypted[string]
	NewTransactionTime_ms     encrypt.ToBeEncrypted[int64]
	NewCreationProcessType    encrypt.ToBeEncrypted[model.ECreationProcess]
}

func (v *assetView) TransferAsset(
	transaction contractapi.TransactionContextInterface,
	requestJson string,
) (string, error) {
	ctx := context.Background()
	service := service.NewSmartContractServiceHyperledger(transaction)

	request := transefrAssetRequest{}
	err := json.Unmarshal([]byte(requestJson), &request)
	if err != nil {
		return "", err
	}

	newAsset, stackErr := v.controller.TransferAsset(
		ctx,
		service,
		&request.CurrentId,
		request.CurrentSignature,
		&request.CurrentTransactionTime_ms,
		&request.NewId,
		request.NewSignature,
		&request.NewOwnerPublicKey,
		&request.NewTransactionTime_ms,
		&request.NewCreationProcessType,
	)

	if stackErr != nil {
		return "", fmt.Errorf(stackErr.String())
	}

	retJson, err := json.Marshal(newAsset)
	if err != nil {
		return "", err
	}

	return string(retJson), nil
}
