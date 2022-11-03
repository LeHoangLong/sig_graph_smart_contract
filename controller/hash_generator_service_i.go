package controller

type HashGeneratorI interface {
	// generate hash in base64 format
	New(id string, secret string) string
}
