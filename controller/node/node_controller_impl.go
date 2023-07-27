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
	"sig_graph/controller"
	"sig_graph/encrypt"
	"sig_graph/model"
	"sig_graph/utility"
)

type nodeControllerImpl[T any] struct {
	clock       utility.ClockI
	settings    utility.SettingsI
	nameService controller.NodeNameServiceI
}

func NewNodeController[T any](
	clock utility.ClockI,
	settings utility.SettingsI,
	nameService controller.NodeNameServiceI,
) NodeController[T] {
	return &nodeControllerImpl[T]{
		clock:       clock,
		settings:    settings,
		nameService: nameService,
	}
}

func (c *nodeControllerImpl[T]) CreateNode(
	Ctx context.Context,
	SmartContract controller.SmartContractServiceI,
	TransactionTime *encrypt.ToBeEncrypted[int64],
	Node *model.Node[T],
	Signatures []string,
) (*model.Node[T], utility.Error) {
	nodeId := Node.Header.Id
	if !c.nameService.IsFullId(nodeId) {
		var err utility.Error
		nodeId, err = c.nameService.GenerateFullId(nodeId)
		if err != nil {
			return nil, err
		}
	}

	exists, err := c.DoesIdExist(
		Ctx,
		SmartContract,
		Node.Header.Id,
	)

	if err != nil {
		return nil, err.AddMessage("fail to create node")
	}

	if exists {
		return nil, utility.NewError(utility.ErrAlreadyExists).AddMessage("id already exists")
	}

	createTime, err := encrypt.BuildEncrypted(TransactionTime)
	if err != nil {
		return nil, err.AddMessage("fail to create node")
	}

	newNode, err := utility.DeepCopy(Node)
	if err != nil {
		return nil, err.AddMessage("fail to create node")
	}
	newNode.Header.CreatedTime = *createTime

	return c.SetNode(Ctx, SmartContract, TransactionTime, Node, Signatures)
}

func (c *nodeControllerImpl[T]) SetNode(
	Ctx context.Context,
	SmartContract controller.SmartContractServiceI,
	TransactionTime *encrypt.ToBeEncrypted[int64],
	Node *model.Node[T],
	Signatures []string,
) (*model.Node[T], utility.Error) {
	nodeId := Node.Header.Id
	newNode, stackErr := utility.DeepCopy(Node)
	time_ms := TransactionTime.Value
	if stackErr != nil {
		return nil, stackErr.AddMessage("fail to set node")
	}

	if !c.nameService.IsFullId(nodeId) {
		nodeId, err := c.nameService.GenerateFullId(nodeId)
		if err != nil {
			return nil, err
		}
		newNode.Header.Id = nodeId
	}

	encryptedTransactionTime, stackErr := encrypt.BuildEncrypted(TransactionTime)
	if stackErr != nil {
		return nil, stackErr.AddMessage("fail to set node")
	}
	newNode.Header.UpdatedTime = *encryptedTransactionTime

	localTime := c.clock.Now_ms()
	if localTime < time_ms || uint64(localTime-time_ms) > c.settings.MaxTimeDifference_ms() {
		return nil, utility.NewError(ErrInvalidTimestamp).AddError(utility.ErrInvalidArgument).AddMessage(fmt.Sprintf("local time: %d, time_ms: %d, max diff: %d", localTime, time_ms, c.settings.MaxTimeDifference_ms()))
	}

	isNodePresentAndFinalized, errorWithStackTrace := c.isNodePresentAndFinalized(Ctx, SmartContract, nodeId)
	if errorWithStackTrace != nil {
		return nil, errorWithStackTrace
	} else if isNodePresentAndFinalized {
		return nil, utility.NewError(utility.ErrNotAllowed).AddMessage("node is finalized")
	}

	if len(Node.Header.OwnerPublicKeys) != len(Signatures) {
		return nil, utility.NewError(utility.ErrInvalidArgument).AddMessage("all signatures must be provided")
	}

	nodeJson, err := json.Marshal(Node)
	if err != nil {
		return nil, utility.NewError(err)
	}

	for i := range Node.Header.OwnerPublicKeys {
		decrypted := Node.Header.OwnerPublicKeys[i].Decrypted
		if decrypted == nil {
			return nil, utility.NewError(utility.ErrInvalidArgument).AddMessage("all public keys must be provided")
		}
		decryptedKey := decrypted.Get()
		signature := Signatures[i]
		stackErr := verify(string(nodeJson), decryptedKey, signature)
		if stackErr != nil {
			return nil, stackErr
		}
	}

	SmartContract.PutState(Ctx, nodeId, newNode)
	return newNode, nil
}

func (c *nodeControllerImpl[T]) AreIdsAvailable(
	ctx context.Context,
	smartContract controller.SmartContractServiceI,
	nodeIds map[string]bool,
) (bool, utility.Error) {
	for id := range nodeIds {
		node := map[string]any{}
		err := smartContract.GetState(ctx, id, &node)
		if err != nil {
			if err != utility.ErrNotFound {
				return false, utility.NewError(err).AddMessage("fail to get node")
			}
		} else {
			return false, nil
		}

	}
	return true, nil
}

func (c *nodeControllerImpl[T]) isNodePresentAndFinalized(
	ctx context.Context,
	smartContract controller.SmartContractServiceI,
	nodeId string,
) (bool, utility.Error) {
	node := model.Node[T]{}
	err := smartContract.GetState(ctx, nodeId, &node)
	if err != nil {
		if err == utility.ErrNotFound {
			return false, nil
		}

		return false, utility.NewError(err).AddMessage("fail to get state from smart contract")
	}

	return node.Header.IsFinalized, nil
}

func verify(data string, publicKey string, signature string) utility.Error {
	block, _ := pem.Decode([]byte(publicKey))
	publicKeyParsed, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return utility.NewError(err).AddMessage("fail to parse pkix public key")
	}

	hash := sha512.Sum512([]byte(data))

	signatureParsed, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return utility.NewError(utility.ErrInvalidArgument).AddMessage("cannot decode base64 signature")
	}

	if rsaKey, ok := publicKeyParsed.(*rsa.PublicKey); ok {
		err = rsa.VerifyPKCS1v15(rsaKey, crypto.SHA512, hash[:], []byte(signatureParsed))
		if err != nil {
			return utility.NewError(ErrInvalidSignature).AddMessage("fail to verify rsa signature")
		}

		return nil
	} else if ecdsaKey, ok := publicKeyParsed.(*ecdsa.PublicKey); ok {
		verified := ecdsa.VerifyASN1(ecdsaKey, hash[:], []byte(signatureParsed))
		if !verified {
			return utility.NewError(ErrInvalidSignature).AddMessage("fail to verify ecdsa signature")
		}

		return nil
	} else {
		return utility.NewError(utility.ErrInvalidArgument).AddMessage("unsupported signature algorithm")
	}
}

// return ErrNotFound if any one id does not exist
func (c *nodeControllerImpl[T]) GetNodes(
	ctx context.Context,
	smartContract controller.SmartContractServiceI,
	ids map[string]bool,
) (map[string]model.Node[T], utility.Error) {
	ret := map[string]model.Node[T]{}

	for id := range ids {
		node := model.Node[T]{}
		err := smartContract.GetState(ctx, id, &node)
		if err != nil {
			return nil, utility.NewError(err).AddMessage("fail to get node")
		}
		ret[id] = node
	}

	return ret, nil
}

func (c *nodeControllerImpl[T]) DoesIdExist(iCtx context.Context, iSmartContract controller.SmartContractServiceI, iNodeId string) (bool, utility.Error) {
	available, err := c.AreIdsAvailable(
		iCtx,
		iSmartContract,
		map[string]bool{
			iNodeId: true,
		},
	)

	if err != nil {
		return false, err
	}

	return !available, nil
}
