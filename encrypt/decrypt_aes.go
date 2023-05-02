package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"sig_graph/utility"
)

type DecryptAes[T any] struct {
	decrypted T
}

func newDecryptAes[T any](
	key []byte,
	iv []byte,
	encrypted []byte,
) (*DecryptAes[T], utility.Error) {
	block, err := aes.NewCipher(key)
	ret := DecryptAes[T]{}

	if err != nil {
		return &ret, utility.NewError(utility.ErrInvalidArgument).AddError(err).AddMessage("fail to create aes cipher from key")
	}
	{
		decrypted := make([]byte, len(encrypted))
		stream := cipher.NewCTR(block, iv)
		stream.XORKeyStream(decrypted, encrypted)
		if err != nil {
			return &ret, utility.NewError(utility.ErrInvalidArgument).AddError(err).AddMessage("fail to decrypt")
		}
		val, trackedErr := bytesToValue[T](decrypted)
		if trackedErr != nil {
			return &ret, trackedErr
		}
		ret.decrypted = val
		return &ret, nil
	}

}

func (d *DecryptAes[T]) Get() T {
	return d.decrypted
}
