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
	// return ErrNotFound if no node with id
	GetNode(ctx context.Context, smartContract controller.SmartContractServiceI, nodeId string, node any) error
}
