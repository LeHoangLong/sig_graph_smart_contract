package encrypt

import (
	"crypto/rand"
	"sig_graph/utility"
)

type toBeEncryptedGeneratorAes[T any] struct {
}

func NewToBeEncryptedGeneratorAes[T any]() ToBeEncryptedGenerator[T] {
	return &toBeEncryptedGeneratorAes[T]{}
}

func (g *toBeEncryptedGeneratorAes[T]) Generate(Val *T) (*ToBeEncrypted[T], utility.Error) {
	ret := &ToBeEncrypted[T]{}
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return ret, utility.NewError(utility.ErrInternalError).AddMessage("fail to generate key")
	}

	iv := make([]byte, 16)
	_, err = rand.Read(iv)
	if err != nil {
		return ret, utility.NewError(utility.ErrInternalError).AddMessage("fail to generate iv")
	}

	ret.EncryptionType = ENCRYPT_TYPE_AES
	ret.PublicArg = map[string]any{
		"iv": iv,
	}
	ret.PrivateArg = map[string]any{
		"key": key,
	}
	ret.Value = *Val
	return ret, nil
}
