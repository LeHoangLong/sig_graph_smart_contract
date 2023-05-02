package encrypt

type Decrypt[T any] interface {
	Get() T
}
