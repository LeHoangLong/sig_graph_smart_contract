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
