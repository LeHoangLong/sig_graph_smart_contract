package encrypt

import "sig_graph/utility"

type Encrypt[T any] interface {
	Encrypt(val *ToBeEncrypted[T]) (*Encrypted[T], utility.Error)
}
