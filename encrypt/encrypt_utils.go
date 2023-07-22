package encrypt

import "sig_graph/utility"

func BuildEncryptedFromRaw[T any](EncryptionType int, iVal *T) (*Encrypted[T], *ToBeEncrypted[T], utility.Error) {
	encrypter, stackErr := BuildToBeEncryptedGenerator[T](EncryptionType)
	if stackErr != nil {
		return nil, nil, stackErr
	}
	toBeEncryptedVal, stackErr := encrypter.Generate(iVal)
	if stackErr != nil {
		return nil, nil, stackErr
	}
	encryptedVal, stackErr := BuildEncrypted(toBeEncryptedVal)
	if stackErr != nil {
		return nil, nil, stackErr
	}

	return encryptedVal, toBeEncryptedVal, nil
}
