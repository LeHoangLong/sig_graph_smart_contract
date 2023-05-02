package encrypt

import (
	"encoding/base64"
	"sig_graph/utility"

	"github.com/mitchellh/mapstructure"
)

func BuildDecrypt[T any](Meta EncryptedMeta, valueStr string, arg map[string]any) (Decrypt[T], utility.Error) {
	if Meta == nil {
		return BuildDecryptFromType[T](ENCRYPT_TYPE_UNENCRYPTED, Meta, valueStr, arg)
	}

	var typeValueInt int = ENCRYPT_TYPE_UNENCRYPTED
	var typeValue any
	var ok bool

	if typeValue, ok = Meta["type"]; ok {
		if typeValueInt, ok = typeValue.(int); !ok {
			return nil, utility.NewError(utility.ErrInvalidArgument).AddMessage("invalid type value, must be int")
		} else {
			return BuildDecryptFromType[T](typeValueInt, Meta, valueStr, arg)
		}
	}

	return nil, utility.NewError(utility.ErrInvalidArgument).AddMessage("missing type value in metadata")

}

func BuildDecryptFromType[T any](iType int, Meta EncryptedMeta, valueStr string, arg map[string]any) (Decrypt[T], utility.Error) {
	value, err := base64.StdEncoding.DecodeString(valueStr)
	if err != nil {
		return nil, utility.NewError(err).AddMessage("cannot base64 decode value")
	}

	builderMap := map[int]buildDecrypt[T]{
		ENCRYPT_TYPE_UNENCRYPTED: buildUnencrypted[T],
		ENCRYPT_TYPE_AES:         buildAes[T],
	}

	if builder, ok := builderMap[iType]; !ok {
		return nil, utility.NewError(utility.ErrInvalidArgument).AddMessage("unsupported decryption type")
	} else {
		return builder(Meta, value, arg)
	}
}

type buildDecrypt[T any] func(Meta EncryptedMeta, value []byte, arg map[string]any) (Decrypt[T], utility.Error)

func buildUnencrypted[T any](Meta EncryptedMeta, value []byte, arg map[string]any) (Decrypt[T], utility.Error) {
	val, err := newDecryptUnencrypted[T](value)
	if err != nil {
		return nil, err.AddMessage("failed to decrypt unencrypted value")
	}
	return val, nil
}

func buildAes[T any](Meta EncryptedMeta, value []byte, iArg map[string]any) (Decrypt[T], utility.Error) {
	aesArg := AesArgs{}
	err := mapstructure.Decode(iArg, &aesArg)
	if err != nil {
		return nil, utility.NewError(err).AddError(utility.ErrInvalidArgument).AddMessage("fail to decode aes decrypt arg")
	}

	iv := []byte{}
	if aesArg.Iv != nil {
		iv = aesArg.Iv
	} else if value, ok := Meta["iv"]; ok {
		if valueStr, ok := value.(string); !ok {
			return nil, utility.NewError(utility.ErrInvalidArgument).AddMessage("invalid iv type in metadata, expecting string")
		} else {
			iv, err = base64.StdEncoding.DecodeString(valueStr)
			if err != nil {
				return nil, utility.NewError(err).AddMessage("cannot base64 decode value")
			}
		}
	}

	return newDecryptAes[T]([]byte(aesArg.Key), iv, value)
}

type AesArgs struct {
	Key []byte `mapstructure:"key"`
	Iv  []byte `mapstructure:"iv"`
}
