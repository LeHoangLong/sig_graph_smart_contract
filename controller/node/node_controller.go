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
	"sig_graph/encrypt"
	"sig_graph/model"
	"sig_graph/utility"
)

var ErrInvalidSignature = fmt.Errorf("%w: invalid signature", utility.ErrInvalidArgument)
var ErrInvalidTimestamp = fmt.Errorf("%w: invalid timestamp", utility.ErrInvalidArgument)

//go:generate mockgen -source=$GOFILE -destination ../../testutils/node_controller.go -package mock
type NodeController[T any] interface {
	// verify node id doesnot exist. If it does, return AlreadyExists error
	// then call SetNode
	CreateNode(ctx context.Context, smartContract controller.SmartContractServiceI, time *encrypt.ToBeEncrypted[int64], node *model.Node[T], Signatures []string) (*model.Node[T], utility.Error)
	// verify that signature matches, timestamp is within allowed limit
	// return ErrInvalidTimestamp if timestamp fails
	// return ErrInvalidSignature if signature fails
	// return ErrInvalidArgument for other argument failure
	// set updatedTime to timestamp
	// set id to full id if it is not
	SetNode(ctx context.Context, smartContract controller.SmartContractServiceI, time *encrypt.ToBeEncrypted[int64], node *model.Node[T], Signatures []string) (*model.Node[T], utility.Error)
	// check if node ids are available
	// returns false if any one id does not exist
	AreIdsAvailable(ctx context.Context, smartContract controller.SmartContractServiceI, nodeIds map[string]bool) (bool, utility.Error)
	// check if id exist
	DoesIdExist(ctx context.Context, smartContract controller.SmartContractServiceI, nodeId string) (bool, utility.Error)
	// return ErrNotFound if any one id does not exist
	GetNodes(ctx context.Context, smartContract controller.SmartContractServiceI, ids map[string]bool) (nodes map[string]model.Node[T], err utility.Error)
}
