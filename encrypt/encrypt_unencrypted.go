package encrypt

import (
	"encoding/base64"
	"fmt"
	"sig_graph/utility"
)

type encryptUnencrypted[T any] struct {
}

func newEncryptUnencrypted[T any]() Encrypt[T] {
	return &encryptUnencrypted[T]{}
}

func (e *encryptUnencrypted[T]) Encrypt(iVal *ToBeEncrypted[T]) (*Encrypted[T], utility.Error) {
	valStr := []byte(fmt.Sprintf("%+v", iVal.Value))
	valBase64 := base64.StdEncoding.EncodeToString(valStr)

	decrypt, errorStacktrace := BuildDecrypt[T](nil, valBase64, iVal.PrivateArg)
	if errorStacktrace != nil {
		return nil, errorStacktrace
	}

	return NewEncrypted(
		iVal.PublicArg,
		valBase64,
		decrypt,
	), nil
}
