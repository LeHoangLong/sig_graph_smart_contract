package controller

type NodeNameServiceI interface {
	/// append graph id to node id
	GenerateFullId(id string) (string, error)
	IsIdValid(id string) bool
	IsFullId(id string) bool
}
