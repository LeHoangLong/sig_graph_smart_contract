package controller

import "context"

type SmartContractServiceI interface {
	PutState(ctx context.Context, key string, value any) error
	// return ErrNotFound if not found
	GetState(ctx context.Context, key string, value any) error
}
