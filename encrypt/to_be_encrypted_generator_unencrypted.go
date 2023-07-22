package encrypt

import "sig_graph/utility"

type toBeEncryptedGeneratorUnencrypted[T any] struct {
}

func NewToBeEncryptedGeneratorUnencrypted[T any]() ToBeEncryptedGenerator[T] {
	return &toBeEncryptedGeneratorUnencrypted[T]{}
}

func (g *toBeEncryptedGeneratorUnencrypted[T]) Generate(Val *T) (*ToBeEncrypted[T], utility.Error) {
	ret := ToBeEncrypted[T]{}
	ret.EncryptionType = ENCRYPT_TYPE_UNENCRYPTED
	ret.PublicArg = nil
	ret.PrivateArg = nil
	ret.Value = *Val
	return &ret, nil
}
