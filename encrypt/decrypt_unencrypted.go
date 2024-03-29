package encrypt

import "sig_graph/utility"

type DecryptUnencrypted[T any] struct {
	Decrypted T
}

func newDecryptUnencrypted[T any](
	encrypted []byte,
) (*DecryptUnencrypted[T], utility.Error) {
	val, err := bytesToValue[T](encrypted)
	ret := &DecryptUnencrypted[T]{}
	if err != nil {
		return ret, err
	}

	ret.Decrypted = val
	return ret, nil
}

func (d *DecryptUnencrypted[T]) Get() T {
	return d.Decrypted
}
