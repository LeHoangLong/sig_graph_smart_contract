package encrypt

import "sig_graph/utility"

type ToBeEncryptedGenerator[T any] interface {
	Generate(Val *T) (*ToBeEncrypted[T], utility.Error)
}
