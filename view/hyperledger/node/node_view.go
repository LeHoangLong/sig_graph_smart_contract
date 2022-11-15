package node_hyperledger_view

import (
	"context"
	"encoding/json"
	"fmt"
	node_controller "sig_graph/controller/node"
	"sig_graph/service"
	"sig_graph/utility"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type nodeView struct {
	contractapi.Contract
	controller node_controller.NodeControllerI
}

func NewNodeView(
	controller node_controller.NodeControllerI,
) *nodeView {
	return &nodeView{
		controller: controller,
	}
}

func (c *nodeView) DoNodeIdsExist(
	transaction contractapi.TransactionContextInterface,
	request string,
) (map[string]bool, error) {
	ids := map[string]bool{}
	err := json.Unmarshal([]byte(request), &ids)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", utility.ErrInvalidArgument, err.Error())
	}

	ctx := context.Background()
	service := service.NewSmartContractServiceHyperledger(transaction)
	if len(ids) > 100 {
		// avoid scanning / dos attack, so we only allow 100 ids to be checked at a time
		return nil, utility.ErrInvalidArgument
	}

	return c.controller.DoNodeIdsExist(ctx, service, ids)

}

type getNodesByIdRequest struct {
	Ids map[string]bool `json:"ids"`
}

func (v *nodeView) GetNodesById(
	transaction contractapi.TransactionContextInterface,
	requestJson string,
) (string, error) {
	ctx := context.Background()
	service := service.NewSmartContractServiceHyperledger(transaction)

	request := getNodesByIdRequest{}
	err := json.Unmarshal([]byte(requestJson), &request)
	if err != nil {
		return "", err
	}

	nodes, err := v.controller.GetNodes(
		ctx,
		service,
		request.Ids,
	)
	if err != nil {
		return "", err
	}

	nodesJson, err := json.Marshal(nodes)
	if err != nil {
		return "", err
	}

	return string(nodesJson), nil
}
