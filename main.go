package main

import (
	"log"
	asset_controller "sig_graph/controller/asset"
	node_controller "sig_graph/controller/node"
	"sig_graph/service"
	"sig_graph/utility"
	asset_hyperledger_view "sig_graph/view/hyperledger/asset"
	node_hyperledger_view "sig_graph/view/hyperledger/node"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	// intialize utility
	settings := utility.NewSettings("sgp://hyper:[http://localhost:7051,http://localhost:9051]:public")
	clockReal := utility.NewClockRealtime()
	cloner := utility.NewCloner()

	// initialize servcies
	nameService := service.NewNodeNameService(settings)
	hashGeneratorService := service.NewHashGenerator()

	// initialize controllers
	nodeController := node_controller.NewNodeController(clockReal, settings, nameService)
	assetController := asset_controller.NewAssetController(nodeController, nameService, hashGeneratorService, cloner)

	// initialize views
	assetView := asset_hyperledger_view.NewAssetView(
		assetController,
	)
	nodeView := node_hyperledger_view.NewNodeView(
		nodeController,
	)

	//
	// initialize and start contract
	chaincode, err := contractapi.NewChaincode(
		assetView,
		nodeView,
	)

	if err != nil {
		log.Panicf("Error creating chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting chaincode: %v", err)
	}
}
