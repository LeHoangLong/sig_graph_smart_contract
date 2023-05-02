package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"sig_graph/utility"

	"github.com/mitchellh/mapstructure"
)

type encryptAes[T any] struct {
}

type encryptAesEncryptPrivateArg struct {
	Key string `mapstructure:"key"`
}

type encryptAesEncryptPublicArg struct {
	Iv string `mapstructure:"iv"`
}

func newEncryptAes[T any]() Encrypt[T] {
	return &encryptAes[T]{}
}

func (e *encryptAes[T]) Encrypt(iVal *ToBeEncrypted[T]) (*Encrypted[T], utility.Error) {
	privateArg := encryptAesEncryptPrivateArg{}

	err := mapstructure.Decode(iVal.PrivateArg, &privateArg)
	if err != nil {
		return nil, utility.NewError(err).AddError(utility.ErrInvalidArgument).AddMessage("fail to decode private arg")
	}

	publicArg := encryptAesEncryptPublicArg{}
	err = mapstructure.Decode(iVal.PublicArg, &publicArg)
	if err != nil {
		return nil, utility.NewError(err).AddError(utility.ErrInvalidArgument).AddMessage("fail to decode public arg")
	}

	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateArg.Key)
	if err != nil {
		return nil, utility.NewError(err).AddError(utility.ErrInvalidArgument).AddMessage("fail to base64 decode key")
	}

	block, err := aes.NewCipher(privateKeyBytes)
	if err != nil {
		return nil, utility.NewError(err).AddError(utility.ErrInvalidArgument).AddMessage("invalid aes key")
	}

	ivBytes, err := base64.StdEncoding.DecodeString(publicArg.Iv)
	if err != nil {
		return nil, utility.NewError(err).AddError(utility.ErrInvalidArgument).AddMessage("fail to base64 decode iv")
	}

	valStr := []byte(fmt.Sprintf("%+v", iVal.Value))
	encrypted := make([]byte, len(valStr))
	mode := cipher.NewCTR(block, []byte(ivBytes))
	mode.XORKeyStream(encrypted, valStr)

	encryptedBase64 := base64.StdEncoding.EncodeToString(encrypted)

	decrypt, errorStacktrace := BuildDecryptFromType[T](ENCRYPT_TYPE_AES, iVal.PublicArg, encryptedBase64, map[string]any{
		"iv":  ivBytes,
		"key": privateKeyBytes,
	})
	if errorStacktrace != nil {
		return nil, errorStacktrace
	}

	return NewEncrypted(
		iVal.PublicArg,
		encryptedBase64,
		decrypt,
	), nil
}
