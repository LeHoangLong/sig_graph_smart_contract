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
	"sig_graph/model"

	"github.com/shopspring/decimal"
)

//go:generate mockgen -source=$GOFILE -destination ../../testutils/asset_controller.go -package mock
type AssetControllerI interface {
	// return ErrAlreadyExists f id already used
	CreateAsset(
		ctx context.Context,
		smartContract controller.SmartContractServiceI,
		time_ms uint64,
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
	) (*model.Asset, error)
	// return ErrNotFound if no material with id
	GetAsset(ctx context.Context, smartContract controller.SmartContractServiceI, id string) (*model.Asset, error)
	// return ErrNotFound if no material with currentId
	// return ErrAlreadyExists if newId already used
	// return new transferred asset
	TransferAsset(
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
	) (*model.Asset, error)
}
