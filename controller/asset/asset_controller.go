package asset_controller

import (
	"context"
	"fmt"
	"sig_graph/controller"
	node_controller "sig_graph/controller/node"
	"sig_graph/model"
	"sig_graph/utility"

	"github.com/shopspring/decimal"
)

type assetController struct {
	nodeController node_controller.NodeControllerI
	nameService    controller.NodeNameServiceI
	hashGenerator  controller.HashGeneratorI
	cloner         utility.ClonerI
}

func NewAssetController(
	nodeController node_controller.NodeControllerI,
	nameService controller.NodeNameServiceI,
	hashGenerator controller.HashGeneratorI,
	cloner utility.ClonerI,
) AssetControllerI {
	return &assetController{
		nodeController: nodeController,
		nameService:    nameService,
		hashGenerator:  hashGenerator,
		cloner:         cloner,
	}
}

func (c *assetController) CreateAsset(
	ctx context.Context,
	smartContract controller.SmartContractServiceI,
	time uint64,
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
) (*model.Asset, error) {
	fullId, err := c.nameService.GenerateFullId(id)
	if err != nil {
		return nil, err
	}

	exist, err := c.nodeController.DoNodeIdsExist(ctx, smartContract, map[string]bool{
		fullId: true,
	})
	if err != nil {
		return nil, err
	}
	if exist[fullId] {
		return nil, utility.ErrAlreadyExists
	}

	node := model.NewDefaultNode(fullId, model.ENodeTypeAsset, time, time, signature, ownerPublicKey)

	asset := model.Asset{
		Node:            node,
		CreationProcess: model.ECreationProcessCreate,
		Unit:            unit,
		Quantity:        quantity,
		MaterialName:    materialName,
	}

	err = c.nodeController.SetNode(ctx, smartContract, time, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (c *assetController) GetAsset(
	ctx context.Context,
	smartContract controller.SmartContractServiceI,
	id string,
) (*model.Asset, error) {
	if !c.nameService.IsFullId(id) {
		if c.nameService.IsIdValid(id) {
			var err error
			id, err = c.nameService.GenerateFullId(id)
			if err != nil {
				return nil, err
			}
		}
	}

	nodes, err := c.nodeController.GetNodes(ctx, smartContract, map[string]bool{
		id: true,
	})
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, utility.ErrNotFound
	}

	if asset, ok := nodes[id].(model.Asset); !ok {
		fmt.Println("invalid type 2")
		return nil, utility.ErrInvalidNodeType
	} else {
		return &asset, nil
	}
}

func (c *assetController) TransferAsset(
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
) (*model.Asset, error) {
	if !c.nameService.IsFullId(currentId) {
		if c.nameService.IsIdValid(currentId) {
			var err error
			currentId, err = c.nameService.GenerateFullId(currentId)
			if err != nil {
				return nil, err
			}
		}
	}

	if !c.nameService.IsFullId(newId) {
		if c.nameService.IsIdValid(newId) {
			var err error
			newId, err = c.nameService.GenerateFullId(newId)
			if err != nil {
				return nil, err
			}
		}
	}

	exist, err := c.nodeController.DoNodeIdsExist(
		ctx,
		smartContract,
		map[string]bool{
			currentId: true,
			newId:     true,
		},
	)

	if err != nil {
		return nil, err
	}

	if !exist[currentId] {
		return nil, utility.ErrNotFound
	}

	if exist[newId] {
		return nil, utility.ErrAlreadyExists
	}

	asset, err := c.GetAsset(
		ctx,
		smartContract,
		currentId,
	)

	if err != nil {
		return nil, err
	}

	// create new node first in case this fails, the new node simply just
	// become an orphan node
	newAsset := &model.Asset{}
	err = c.cloner.Clone(ctx, asset, newAsset)
	if err != nil {
		return nil, err
	}
	newAsset.CreatedTime = time_ms
	newAsset.UpdatedTime = time_ms
	newAsset.OwnerPublicKey = newOwnerPublicKey
	newAsset.CreationProcess = model.ECreationProcessTransfer
	newAsset.ClearAllEdges()
	newAsset.Signature = newSignature
	newAsset.Id = newId
	if currentSecret != "" {
		currentHash := c.hashGenerator.New(currentId, currentSecret)
		newAsset.PrivateParentsHashedIds[currentHash] = true
	} else {
		newAsset.PublicParentsIds[currentId] = true
	}

	err = c.nodeController.SetNode(ctx, smartContract, time_ms, &newAsset)
	if err != nil {
		return nil, err
	}

	// update current node
	updatedAsset := &model.Asset{}
	err = c.cloner.Clone(ctx, asset, updatedAsset)
	if err != nil {
		return nil, err
	}
	updatedAsset.UpdatedTime = time_ms
	updatedAsset.Signature = currentSignature
	updatedAsset.IsFinalized = true
	if newSecret != "" {
		newHash := c.hashGenerator.New(newId, newSecret)
		updatedAsset.PrivateChildrenHashedIds[newHash] = true
	} else {
		updatedAsset.PublicChildrenIds[newId] = true
	}

	err = c.nodeController.SetNode(ctx, smartContract, time_ms, updatedAsset)
	if err != nil {
		return nil, err
	}

	return newAsset, nil
}
