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
package service

import (
	"context"
	"encoding/json"
	"sig_graph/controller"
	"sig_graph/utility"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type smartContractServiceHyperledger struct {
	transaction contractapi.TransactionContextInterface
}

func NewSmartContractServiceHyperledger(
	transaction contractapi.TransactionContextInterface,
) controller.SmartContractServiceI {
	return &smartContractServiceHyperledger{
		transaction: transaction,
	}
}

func (s *smartContractServiceHyperledger) PutState(ctx context.Context, key string, value any) error {
	valueStr, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.transaction.GetStub().PutState(key, valueStr)
}

func (s *smartContractServiceHyperledger) GetState(ctx context.Context, key string, value any) error {
	valueStr, err := s.transaction.GetStub().GetState(key)
	if err != nil {
		return err
	}

	if valueStr == nil {
		return utility.ErrNotFound
	}

	err = json.Unmarshal(valueStr, value)
	if err != nil {
		return err
	}

	return nil
}
