package controller

//go:generate mockgen -source=$GOFILE -destination ../testutils/node_name_service.go -package mock
type NodeNameServiceI interface {
	/// append graph id to node id
	GenerateFullId(id string) (string, error)
	IsIdValid(id string) bool
	IsFullId(id string) bool
}
