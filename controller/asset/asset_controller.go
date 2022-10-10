package asset_controller

import (
	"context"
	"fmt"
	"sig_graph/controller"
	node_controller "sig_graph/controller/node"
	"sig_graph/model"
	"sig_graph/utility"
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

func (c *assetController) CreateMaterial(
	ctx context.Context,
	smartContract controller.SmartContractServiceI,
	time uint64,
	id string,
	signature string,
	ownerPublicKey string,
) error {
	if !c.nameService.IsIdValid(id) {
		return fmt.Errorf("invalid id: %w", utility.ErrInvalidArgument)
	}

	temp := map[any]any{}
	err := c.nodeController.GetNode(ctx, smartContract, id, &temp)
	if err == nil {
		return utility.ErrAlreadyExists
	}

	fullId, err := c.nameService.GenerateFullId(id)
	if err != nil {
		return err
	}

	node := model.NewDefaultNode(fullId, "asset", time, time, signature, ownerPublicKey)
	err = c.nodeController.SetNode(ctx, smartContract, time, &node)
	if err != nil {
		return err
	}

	return nil
}

func (c *assetController) GetMaterial(
	ctx context.Context,
	smartContract controller.SmartContractServiceI,
	id string,
) (*model.Asset, error) {
	asset := model.Asset{}
	err := c.nodeController.GetNode(ctx, smartContract, id, &asset)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}
