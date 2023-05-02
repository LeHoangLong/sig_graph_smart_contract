package encrypt

type EncryptedMeta = map[string]any
type Encrypted[T any] struct {
	Meta      EncryptedMeta `json:"meta"`
	Value     string        `json:"value"`
	Decrypted Decrypt[T]    `json:"-"`
}

func NewEncrypted[T any](
	Meta EncryptedMeta,
	Value string,
	Decrypted Decrypt[T],
) *Encrypted[T] {
	return &Encrypted[T]{
		Meta:      Meta,
		Value:     Value,
		Decrypted: Decrypted,
	}
}
