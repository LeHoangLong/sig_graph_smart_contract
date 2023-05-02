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
	"sig_graph/encrypt"
	"sig_graph/model"
	"sig_graph/utility"
)

//go:generate mockgen -source=$GOFILE -destination ../../testutils/asset_controller.go -package mock
type AssetController interface {
	// return ErrAlreadyExists f id already used
	CreateAsset(
		Ctx context.Context,
		SmartContract controller.SmartContractServiceI,
		TransactionTime *encrypt.ToBeEncrypted[int64],
		CreationProcessType *encrypt.ToBeEncrypted[model.ECreationProcess],
		Asset *model.NodeAsset,
	) (*model.NodeAsset, utility.Error)

	// return ErrNotFound if no material with currentId
	// return ErrAlreadyExists if newId already used
	// return new transferred asset
	TransferAsset(
		Ctx context.Context,
		SmartContract controller.SmartContractServiceI,
		CurrentId *encrypt.ToBeEncrypted[string],
		CurrentSignature string,
		CurrentTransactionTime_ms *encrypt.ToBeEncrypted[int64],
		NewId *encrypt.ToBeEncrypted[string],
		NewSignature string,
		NewOwnerPublicKey *encrypt.ToBeEncrypted[string],
		NewTransactionTime_ms *encrypt.ToBeEncrypted[int64],
		NewCreationProcessType *encrypt.ToBeEncrypted[model.ECreationProcess],
	) (*model.NodeAsset, utility.Error)
}
