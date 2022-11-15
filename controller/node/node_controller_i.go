package node_controller

import (
	"context"
	"fmt"
	"sig_graph/controller"
	"sig_graph/utility"
)

var ErrInvalidSignature = fmt.Errorf("%w: invalid signature", utility.ErrInvalidArgument)
var ErrInvalidTimestamp = fmt.Errorf("%w: invalid timestamp", utility.ErrInvalidArgument)

//go:generate mockgen -source=$GOFILE -destination ../../testutils/node_controller.go -package mock
type NodeControllerI interface {
	// verify that signature matches, timestamp is within allowed limit
	// return ErrInvalidTimestamp if timestamp fails
	// return ErrInvalidSignature if signature fails
	// return ErrInvalidArgument for other argument failure
	// set updatedTime to timestamp
	// set id to full id if it is not
	SetNode(ctx context.Context, smartContract controller.SmartContractServiceI, time_ms uint64, node any) error
	// check if node ids exist
	// for every key in nodeIds, a corresponding entry in the output will indicate whether the
	// id already exists
	DoNodeIdsExist(ctx context.Context, smartContract controller.SmartContractServiceI, nodeIds map[string]bool) (map[string]bool, error)
	// return ErrNotFound if any one id does not exist
	GetNodes(ctx context.Context, smartContract controller.SmartContractServiceI, ids map[string]bool) (nodes map[string]any, err error)
}
