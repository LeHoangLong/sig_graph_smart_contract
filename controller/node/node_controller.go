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
	"sig_graph/utility"
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

	localTime := c.clock.Now_ms()
	if localTime < time_ms || localTime-time_ms > c.settings.MaxTimeDifference_ms() {
		return ErrInvalidTimestamp
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

func (c *nodeController) GetNode(ctx context.Context, smartContract controller.SmartContractServiceI, nodeId string, node any) error {
	return smartContract.GetState(ctx, nodeId, node)
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
