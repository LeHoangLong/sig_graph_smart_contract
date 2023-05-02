package utility

import "github.com/barkimedes/go-deepcopy"

func DeepCopy[T any](x T) (T, Error) {
	cloned, err := deepcopy.AnythingType(x)
	if err != nil {
		return cloned, NewError(err).AddError(ErrInvalidArgument).AddMessage("fail to clone object")
	}

	return cloned, nil
}
