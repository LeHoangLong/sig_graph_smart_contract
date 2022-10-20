package asset_controller

import (
	"context"
	"sig_graph/controller"
	node_controller "sig_graph/controller/node"
	"sig_graph/model"
	"sig_graph/utility"

	"github.com/shopspring/decimal"
)

type assetController struct {
	nodeController node_controller.NodeControllerI
	nameService    controller.NodeNameServiceI
}

func NewAssetController(
	nodeController node_controller.NodeControllerI,
	nameService controller.NodeNameServiceI,
) AssetControllerI {
	return &assetController{
		nodeController: nodeController,
		nameService:    nameService,
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

	temp := map[any]any{}
	err = c.nodeController.GetNode(ctx, smartContract, fullId, &temp)
	if err == nil {
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

	asset := model.Asset{}
	err := c.nodeController.GetNode(ctx, smartContract, id, &asset)
	if err != nil {
		return nil, err
	}

	if asset.NodeType != model.ENodeTypeAsset {
		return nil, utility.ErrInvalidNodeType
	}
	return &asset, nil
}
