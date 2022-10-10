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
