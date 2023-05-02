package encrypt

import "sig_graph/utility"

type EncryptMeta = map[string]any

func BuildEncrypted[T any](iVal *ToBeEncrypted[T]) (*Encrypted[T], utility.Error) {
	var encrypter Encrypt[T]
	switch iVal.EncryptionType {
	case ENCRYPT_TYPE_AES:
		encrypter = newEncryptAes[T]()
	case ENCRYPT_TYPE_UNENCRYPTED:
		encrypter = newEncryptUnencrypted[T]()
	default:
		return nil, utility.NewError(utility.ErrInvalidArgument).AddMessage("unrecognized encryption type")
	}

	encrypted, err := encrypter.Encrypt(iVal)
	if err != nil {
		return nil, err.AddMessage("fail to build encrypt")
	}

	return encrypted, nil
}
