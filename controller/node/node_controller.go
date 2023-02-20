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
package node_controller

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"reflect"
	"sig_graph/controller"
	"sig_graph/model"
	"sig_graph/utility"

	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
)

type nodeController struct {
	clock       utility.ClockI
	settings    utility.SettingsI
	nameService controller.NodeNameServiceI
}

func NewNodeController(
	clock utility.ClockI,
	settings utility.SettingsI,
	nameService controller.NodeNameServiceI,
) NodeControllerI {
	return &nodeController{
		clock:       clock,
		settings:    settings,
		nameService: nameService,
	}
}

func (c *nodeController) SetNode(
	ctx context.Context,
	smartContract controller.SmartContractServiceI,
	time_ms uint64,
	node any,
) error {
	nodeMap := map[string]any{}
	nodeJson, err := json.Marshal(node)
	if err != nil {
		return err
	}

	err = json.Unmarshal(nodeJson, &nodeMap)
	if err != nil {
		return err
	}

	nodeId := ""
	if tempNodeId, ok := nodeMap["id"].(string); !ok {
		return utility.ErrInvalidArgument
	} else {
		nodeId = tempNodeId
	}

	if !c.nameService.IsFullId(nodeId) {
		nodeId, err := c.nameService.GenerateFullId(nodeId)
		if err != nil {
			return err
		}
		nodeMap["id"] = nodeId
	}

	isNodeFinalized, err := c.isNodeFinalized(ctx, smartContract, nodeId)
	if err != nil {
		if err != utility.ErrNotFound {
			return err
		}
	} else if isNodeFinalized {
		return fmt.Errorf("%w: node is finalized", utility.ErrInvalidState)
	}

	localTime := c.clock.Now_ms()
	if localTime < time_ms || localTime-time_ms > c.settings.MaxTimeDifference_ms() {
		return fmt.Errorf("%w: local time: %d, time_ms: %d, max diff: %d", ErrInvalidTimestamp, localTime, time_ms, c.settings.MaxTimeDifference_ms())
	}

	ownerPublicKey := ""
	if tempOwnerPublicKey, ok := nodeMap["owner_public_key"].(string); !ok || tempOwnerPublicKey == "" {
		return utility.ErrInvalidArgument
	} else {
		ownerPublicKey = tempOwnerPublicKey
	}

	signature := ""
	if tempSignature, ok := nodeMap["signature"]; !ok {
		return utility.ErrInvalidArgument
	} else {

		if tempSignature, ok := tempSignature.(string); !ok {
			return utility.ErrInvalidArgument
		} else {
			signature = tempSignature
		}
	}

	delete(nodeMap, "signature")

	nodeWithoutSignatureJson, err := json.Marshal(nodeMap)
	if err != nil {
		return err
	}

	err = c.verify(string(nodeWithoutSignatureJson), ownerPublicKey, signature)
	if err != nil {
		return err
	}

	nodeMap["signature"] = signature
	smartContract.PutState(ctx, nodeId, nodeMap)

	return nil
}

func (c *nodeController) DoNodeIdsExist(ctx context.Context, smartContract controller.SmartContractServiceI, nodeIds map[string]bool) (map[string]bool, error) {
	ret := map[string]bool{}
	for id := range nodeIds {
		node := map[string]any{}
		err := smartContract.GetState(ctx, id, &node)
		if err != nil {
			if err != utility.ErrNotFound {
				return nil, err
			}

			ret[id] = false
		} else {
			ret[id] = true
		}

	}
	return ret, nil
}

func (c *nodeController) isNodeFinalized(
	ctx context.Context,
	smartContract controller.SmartContractServiceI,
	nodeId string,
) (bool, error) {
	node := map[string]any{}
	err := smartContract.GetState(ctx, nodeId, &node)
	if err != nil {
		return false, err
	}

	if isFinalized, ok := node["is_finalized"].(bool); !ok {
		return false, utility.ErrInvalidState
	} else {
		return isFinalized, nil
	}
}

func (c *nodeController) verify(data string, publicKey string, signature string) error {
	block, _ := pem.Decode([]byte(publicKey))
	publicKeyParsed, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	hash := sha512.Sum512([]byte(data))

	signatureParsed, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("%w: cannot decode base64 signature", utility.ErrInvalidArgument)
	}

	if rsaKey, ok := publicKeyParsed.(*rsa.PublicKey); ok {
		err = rsa.VerifyPKCS1v15(rsaKey, crypto.SHA512, hash[:], []byte(signatureParsed))
		if err != nil {
			return ErrInvalidSignature
		}

		return nil
	} else if ecdsaKey, ok := publicKeyParsed.(*ecdsa.PublicKey); ok {
		verified := ecdsa.VerifyASN1(ecdsaKey, hash[:], []byte(signatureParsed))
		if !verified {
			return ErrInvalidSignature
		}

		return nil
	} else {
		return fmt.Errorf("%w: unsupported signature algorithm", utility.ErrInvalidArgument)
	}
}

func StringToDecimalHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(decimal.NewFromInt(0)) {
			return data, nil
		}

		// Convert it by parsing
		return decimal.NewFromString(data.(string))
	}
}

// return ErrNotFound if any one id does not exist
func (c *nodeController) GetNodes(
	ctx context.Context,
	smartContract controller.SmartContractServiceI,
	ids map[string]bool,
) (map[string]any, error) {
	ret := map[string]any{}

	for id := range ids {
		nodeMap := map[string]any{}
		err := smartContract.GetState(ctx, id, &nodeMap)
		if err != nil {
			return nil, err
		}
		mapJson, _ := json.Marshal(nodeMap)

		nodeType := nodeMap["type"]
		fmt.Printf("nodeMap %+v\n", nodeMap)
		fmt.Println("nodeType ", nodeType)
		switch nodeType {
		case model.ENodeTypeAsset:
			asset := model.Asset{}
			err = json.Unmarshal(mapJson, &asset)
			if err != nil {
				return nil, err
			}

			fmt.Printf("asset: %+v\n", asset)
			ret[id] = asset
		default:
			fmt.Println("invalid type 1")
			return nil, utility.ErrInvalidNodeType
		}
	}

	return ret, nil
}
