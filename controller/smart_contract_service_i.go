package controller

import "context"

//go:generate mockgen -source=$GOFILE -destination ../testutils/smart_contract_service.go -package mock
type SmartContractServiceI interface {
	PutState(ctx context.Context, key string, value any) error
	// return ErrNotFound if not found
	GetState(ctx context.Context, key string, value any) error
}
