package encrypt

import "sig_graph/utility"

func BuildToBeEncryptedGenerator[T any](Type int) (ToBeEncryptedGenerator[T], utility.Error) {
	switch Type {
	case ENCRYPT_TYPE_AES:
		return NewToBeEncryptedGeneratorAes[T](), nil
	case ENCRYPT_TYPE_UNENCRYPTED:
		return NewToBeEncryptedGeneratorUnencrypted[T](), nil
	default:
		return NewToBeEncryptedGeneratorUnencrypted[T](), utility.NewError(utility.ErrInvalidArgument).AddMessage("unsupported encryption type")
	}
}
