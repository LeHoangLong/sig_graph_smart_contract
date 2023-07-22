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
package asset_controller

import (
	"context"
	"sig_graph/controller"
	node_controller "sig_graph/controller/node"
	"sig_graph/encrypt"
	"sig_graph/model"
	"sig_graph/utility"
)

type assetControllerImpl struct {
	nodeController node_controller.NodeController[model.Asset]
	nameService    controller.NodeNameServiceI
	hashGenerator  controller.HashGeneratorI
}

func NewAssetController(
	nodeController node_controller.NodeController[model.Asset],
	nameService controller.NodeNameServiceI,
	hashGenerator controller.HashGeneratorI,
) AssetController {
	return &assetControllerImpl{
		nodeController: nodeController,
		nameService:    nameService,
		hashGenerator:  hashGenerator,
	}
}

func (c *assetControllerImpl) CreateAsset(
	Ctx context.Context,
	Args *CreateAssetArgs,
) (*model.NodeAsset, utility.Error) {
	if Args.CreationProcessType.Value != model.ECreationProcessCreate {
		return nil, utility.NewError(utility.ErrInvalidArgument).AddMessage("invalid creation process type")
	}
	Args.Asset.Header.NodeType = model.ENodeTypeAsset
	asset, err := c.nodeController.CreateNode(Ctx, Args.SmartContract, Args.TransactionTime, Args.Asset, Args.Signatures)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func (c *assetControllerImpl) getAsset(
	ctx context.Context,
	smartContract controller.SmartContractServiceI,
	id string,
) (*model.NodeAsset, utility.Error) {
	assets, err := c.nodeController.GetNodes(ctx, smartContract, map[string]bool{
		id: true,
	})
	if err != nil {
		return nil, err
	}
	if len(assets) == 0 {
		return nil, utility.NewError(utility.ErrNotFound)
	}

	asset := assets[id]
	return &asset, nil
}

func (c *assetControllerImpl) TransferAsset(
	Ctx context.Context,
	Args *TransferAssetArgs,
) (*model.NodeAsset, utility.Error) {
	var err utility.Error
	if Args.CurrentTransactionTime_ms.Value != Args.NewTransactionTime_ms.Value {
		return nil, utility.NewError(utility.ErrInvalidArgument).AddMessage("mismatched transaction time")
	}

	if Args.NewCreationProcessType.Value != model.ECreationProcessTransfer {
		return nil, utility.NewError(utility.ErrInvalidArgument).AddMessage("invalid transfer process type")
	}

	// build encrypted values first
	currentId := Args.CurrentId.Value
	encryptedCurrentId, err := encrypt.BuildEncrypted(Args.CurrentId)
	if err != nil {
		return nil, err
	}

	encryptedNewOwnerPublicKeys := make([]encrypt.Encrypted[string], len(Args.NewOwnerPublicKeys))
	for i := range Args.NewOwnerPublicKeys {
		encryptedNewOwnerPublicKey, err := encrypt.BuildEncrypted(&Args.NewOwnerPublicKeys[i])
		if err != nil {
			return nil, err
		}

		encryptedNewOwnerPublicKeys[i] = *encryptedNewOwnerPublicKey
	}

	encryptedTransactionTime, err := encrypt.BuildEncrypted(Args.NewTransactionTime_ms)
	if err != nil {
		return nil, err
	}

	encryptedCurrentTransactionTime, err := encrypt.BuildEncrypted(Args.CurrentTransactionTime_ms)
	if err != nil {
		return nil, err
	}

	newId := Args.NewId.Value
	encryptedNewId, err := encrypt.BuildEncrypted(Args.NewId)
	if err != nil {
		return nil, err
	}

	encryptedNewCreationProcess, err := encrypt.BuildEncrypted(Args.NewCreationProcessType)
	if err != nil {
		return nil, err
	}

	currentId, err = c.nameService.GenerateFullId(currentId)
	if err != nil {
		return nil, err
	}

	newId, err = c.nameService.GenerateFullId(newId)
	if err != nil {
		return nil, err
	}

	newIdExist, err := c.nodeController.DoesIdExist(Ctx, Args.SmartContract, newId)
	if newIdExist {
		return nil, utility.NewError(utility.ErrAlreadyExists)
	}

	currentAsset, err := c.getAsset(
		Ctx,
		Args.SmartContract,
		currentId,
	)
	if err != nil {
		return nil, err
	}

	// create new node first in case this fails, the new node simply just
	// become an orphan node
	newAsset := model.NewAsset(
		encryptedNewCreationProcess,
		&currentAsset.Body.Unit,
		&currentAsset.Body.Quantity,
		&currentAsset.Body.MaterialName,
	)

	newHeader := model.NewHeader(
		model.ENodeTypeAsset,
		newId,
		false,
		encryptedNewOwnerPublicKeys,
		model.Edges{},
		model.Edges{},
		encryptedTransactionTime,
		encryptedTransactionTime,
	)
	newHeader.Parents = append(newHeader.Parents, *encryptedCurrentId)

	newNodeAsset, err := model.NewNode(
		newHeader,
		newAsset,
	)

	newNodeAsset, err = c.nodeController.CreateNode(Ctx, Args.SmartContract, Args.NewTransactionTime_ms, newNodeAsset, Args.NewSignatures)
	if err != nil {
		return nil, err
	}

	// update current node
	updatedAsset, err := utility.DeepCopy(currentAsset)
	if err != nil {
		return nil, err
	}
	updatedAsset.Header.UpdatedTime = *encryptedCurrentTransactionTime
	updatedAsset.Header.CreatedTime = *encryptedCurrentTransactionTime
	updatedAsset.Header.IsFinalized = true
	updatedAsset.Header.Children = append(updatedAsset.Header.Children, *encryptedNewId)

	updatedAsset, err = c.nodeController.SetNode(Ctx, Args.SmartContract, Args.NewTransactionTime_ms, updatedAsset, Args.CurrentSignatures)
	if err != nil {
		return nil, err
	}

	return newNodeAsset, nil
}
