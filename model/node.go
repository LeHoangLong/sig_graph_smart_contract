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
	"sig_graph/utility"
)

type ENodeType = string

const (
	ENodeTypeAsset ENodeType = "asset"
)

type Header struct {
	NodeType        ENodeType                   `json:"type"`
	Id              string                      `json:"id"`
	IsFinalized     bool                        `json:"is_finalized"`
	OwnerPublicKeys []encrypt.Encrypted[string] `json:"owner_public_key"`
	Parents         Edges                       `json:"parents"`
	Children        Edges                       `json:"children"`
	CreatedTime     encrypt.Encrypted[int64]    `json:"created_time"`
	UpdatedTime     encrypt.Encrypted[int64]    `json:"updated_time"`
}

type Node[T any] struct {
	Header Header `json:"header"`
	Body   T      `json:"body"`
}

type Edges = []encrypt.Encrypted[string]

func NewHeader(
	NodeType ENodeType,
	Id string,
	IsFinalized bool,
	OwnerPublicKeys []encrypt.Encrypted[string],
	Parents Edges,
	Children Edges,
	CreatedTime *encrypt.Encrypted[int64],
	UpdatedTime *encrypt.Encrypted[int64],
) *Header {
	return &Header{
		NodeType:        NodeType,
		Id:              Id,
		IsFinalized:     IsFinalized,
		OwnerPublicKeys: OwnerPublicKeys,
		Parents:         Parents,
		Children:        Children,
		CreatedTime:     *CreatedTime,
		UpdatedTime:     *UpdatedTime,
	}
}

func NewNode[T any](iHeader *Header, iBody *T) (*Node[T], utility.Error) {
	return &Node[T]{
		Header: *iHeader,
		Body:   *iBody,
	}, nil
}
