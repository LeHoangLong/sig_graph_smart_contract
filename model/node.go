package model

type ENodeType = string

const (
	ENodeTypeAsset ENodeType = "asset"
)

type Node struct {
	Id                       string          `json:"id" mapstructure:"id"`
	PublicParentsIds         map[string]bool `json:"public_parents_ids" mapstructure:"public_parents_ids"`
	PublicChildrenIds        map[string]bool `json:"public_children_ids" mapstructure:"public_children_ids"`
	PrivateParentsHashedIds  map[string]bool `json:"private_parents_hashed_ids" mapstructure:"private_parents_hashed_ids"`
	PrivateChildrenHashedIds map[string]bool `json:"private_children_hashed_ids" mapstructure:"private_children_hashed_ids"`
	IsFinalized              bool            `json:"is_finalized" mapstructure:"is_finalized"`
	NodeType                 ENodeType       `json:"type" mapstructure:"type"`
	CreatedTime              uint64          `json:"created_time" mapstructure:"created_time"`
	UpdatedTime              uint64          `json:"updated_time" mapstructure:"updated_time"`
	Signature                string          `json:"signature" mapstructure:"signature"`
	OwnerPublicKey           string          `json:"owner_public_key" mapstructure:"owner_public_key"`
}

func (n *Node) ClearAllEdges() {
	n.PublicParentsIds = map[string]bool{}
	n.PublicChildrenIds = map[string]bool{}
	n.PrivateParentsHashedIds = map[string]bool{}
	n.PrivateChildrenHashedIds = map[string]bool{}
}

func NewDefaultNode(
	id string,
	nodeType ENodeType,
	createdTime uint64,
	updatedTime uint64,
	signature string,
	ownerPublicKey string,
) Node {
	return Node{
		Id:                       id,
		PublicParentsIds:         map[string]bool{},
		PublicChildrenIds:        map[string]bool{},
		PrivateParentsHashedIds:  map[string]bool{},
		PrivateChildrenHashedIds: map[string]bool{},
		IsFinalized:              false,
		NodeType:                 nodeType,
		CreatedTime:              createdTime,
		UpdatedTime:              updatedTime,
		Signature:                signature,
		OwnerPublicKey:           ownerPublicKey,
	}
}
