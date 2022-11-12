package utility

import "errors"

var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("already exists")
var ErrInvalidArgument = errors.New("invalid argument")
var ErrInvalidNodeType = errors.New("invalid node type")
var ErrInvalidState = errors.New("invalid state")
