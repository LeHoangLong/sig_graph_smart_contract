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
	"sig_graph/utility"

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
	CreationProcess *encrypt.ToBeEncrypted[model.ECreationProcess] `json:"creation_process"`
	TransactionTime encrypt.ToBeEncrypted[int64]                   `json:"transaction_time"`
	OwnerPublicKeys []encrypt.ToBeEncrypted[string]                `json:"owner_public_keys"`
	Signatures      []string                                       `json:"signatures"`
	Asset           model.NodeAsset                                `json:"asset"`
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
	var stackErr utility.Error
	request := createAssetRequest{}
	err := json.Unmarshal([]byte(RequestJson), &request)
	if err != nil {
		return "", err
	}

	creationProcessType := encrypt.ToBeEncrypted[model.ECreationProcess]{
		EncryptionType: encrypt.ENCRYPT_TYPE_UNENCRYPTED,
		Value:          model.ECreationProcessCreate,
	}

	if request.CreationProcess != nil {
		creationProcessType = *request.CreationProcess
	}

	request.Asset.Header.OwnerPublicKeys = make([]encrypt.Encrypted[string], len(request.OwnerPublicKeys))
	for i := range request.OwnerPublicKeys {
		encryptedKey, stackErr := encrypt.BuildEncrypted(&request.OwnerPublicKeys[i])
		if stackErr != nil {
			return "", fmt.Errorf(stackErr.String())
		}
		request.Asset.Header.OwnerPublicKeys[i] = *encryptedKey
	}

	ctx := context.Background()
	service := service.NewSmartContractServiceHyperledger(Transaction)

	createAssetArgs := asset_controller.CreateAssetArgs{
		SmartContract:       service,
		TransactionTime:     &request.TransactionTime,
		CreationProcessType: &creationProcessType,
		Asset:               &request.Asset,
		Signatures:          request.Signatures,
	}
	modelAsset, stackErr := c.controller.CreateAsset(
		ctx,
		&createAssetArgs,
	)
	if stackErr != nil {
		return "", fmt.Errorf(stackErr.String())
	}

	assetJson, err := json.Marshal(modelAsset)
	if err != nil {
		return "", err
	}

	return string(assetJson), nil
}

type transefrAssetRequest struct {
	CurrentId                 encrypt.ToBeEncrypted[string]                  `json:"current_id"`
	CurrentOwnerPublicKeys    []encrypt.ToBeEncrypted[string]                `json:"current_owner_public_keys"`
	CurrentSignatures         []string                                       `json:"current_signatures"`
	CurrentTransactionTime_ms encrypt.ToBeEncrypted[int64]                   `json:"current_transaction_time_ms"`
	NewId                     encrypt.ToBeEncrypted[string]                  `json:"new_id"`
	NewOwnerPublicKeys        []encrypt.ToBeEncrypted[string]                `json:"new_owner_public_keys"`
	NewSignatures             []string                                       `json:"new_signatures"`
	NewTransactionTime_ms     encrypt.ToBeEncrypted[int64]                   `json:"new_transaction_time_ms"`
	NewCreationProcessType    *encrypt.ToBeEncrypted[model.ECreationProcess] `json:"new_creation_process_type"`
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

	newCreationProcessType := encrypt.ToBeEncrypted[model.ECreationProcess]{
		EncryptionType: encrypt.ENCRYPT_TYPE_UNENCRYPTED,
		Value:          model.ECreationProcessTransfer,
	}

	if request.NewCreationProcessType != nil {
		newCreationProcessType = *request.NewCreationProcessType
	}

	transferAssetArg := asset_controller.TransferAssetArgs{
		CurrentId:                 &request.CurrentId,
		CurrentSignatures:         request.CurrentSignatures,
		CurrentTransactionTime_ms: &request.CurrentTransactionTime_ms,
		NewId:                     &request.NewId,
		NewSignatures:             request.NewSignatures,
		NewOwnerPublicKeys:        request.NewOwnerPublicKeys,
		NewTransactionTime_ms:     &request.NewTransactionTime_ms,
		NewCreationProcessType:    &newCreationProcessType,
		SmartContract:             service,
	}

	newAsset, stackErr := v.controller.TransferAsset(
		ctx,
		&transferAssetArg,
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
