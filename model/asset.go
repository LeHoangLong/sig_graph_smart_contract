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
package model

import (
	"sig_graph/encrypt"

	"github.com/shopspring/decimal"
)

type ECreationProcess = string

const (
	ECreationProcessCreate   = "create"
	ECreationProcessTransfer = "transfer"
)

type Asset struct {
	CreationProcess encrypt.Encrypted[ECreationProcess] `json:"creation_process" mapstructure:"creation_process"`
	Unit            encrypt.Encrypted[string]           `json:"unit"  mapstructure:"unit"`
	Quantity        encrypt.Encrypted[decimal.Decimal]  `json:"quantity"  mapstructure:"quantity"`
	MaterialName    encrypt.Encrypted[string]           `json:"material_name"  mapstructure:"material_name"`
}

type NodeAsset = Node[Asset]

func NewAsset(
	CreationProcess *encrypt.Encrypted[ECreationProcess],
	Unit *encrypt.Encrypted[string],
	Quantity *encrypt.Encrypted[decimal.Decimal],
	MaterialName *encrypt.Encrypted[string],
) *Asset {
	return &Asset{
		CreationProcess: *CreationProcess,
		Unit:            *Unit,
		Quantity:        *Quantity,
		MaterialName:    *MaterialName,
	}
}
