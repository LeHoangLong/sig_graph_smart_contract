package model

type Asset struct {
	Node
	CreationProcess string `json:"creation_process"`
	Unit            string `json:"unit"`
	Quantity        string `json:"quantity"`
	MaterialName    string `json:"material_name"`
}
