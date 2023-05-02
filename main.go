// Copyright (C) 2022 Le Hoang Long
// This file is part of SigGraph smart contract <https://github.com/LeHoangLong/sig_graph_smart_contract>.
//
// SigGraph is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// SigGraph is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with SigGraph.  If not, see <http://www.gnu.org/licenses/>.
package main

import (
	"log"
	asset_controller "sig_graph/controller/asset"
	node_controller "sig_graph/controller/node"
	"sig_graph/model"
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

	// initialize servcies
	nameService := service.NewNodeNameService(settings)
	hashGeneratorService := service.NewHashGenerator()

	// initialize controllers
	nodeController := node_controller.NewNodeController[any](clockReal, settings, nameService)
	nodeAssetController := node_controller.NewNodeController[model.Asset](clockReal, settings, nameService)
	assetController := asset_controller.NewAssetController(nodeAssetController, nameService, hashGeneratorService)

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
