package encrypt

type ToBeEncrypted[T any] struct {
	EncryptionType int            `json:"encryption_type"`
	Value          T              `json:"value"`
	PublicArg      map[string]any `json:"public_arg"`
	PrivateArg     map[string]any `json:"private_arg"`
}

func NewToBeEncrypted[T any](
	EncryptionType int,
	Value T,
	PublicArg map[string]any,
	PrivateArg map[string]any,
) ToBeEncrypted[T] {
	return ToBeEncrypted[T]{
		EncryptionType: EncryptionType,
		Value:          Value,
		PublicArg:      PublicArg,
		PrivateArg:     PrivateArg,
	}
}
