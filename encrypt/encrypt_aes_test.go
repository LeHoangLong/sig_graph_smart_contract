package encrypt

import (
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncryptAes(t *testing.T) {
	key, _ := hex.DecodeString("6E6BAFF908E8F63831DF1ED9D37E16DAF9F0C8ED081F9E4E49780F5C13111858")
	iv, _ := hex.DecodeString("29BE740A7FA2E6C99BFF2FC46B685E90")

	t.Run("string", func(t *testing.T) {
		encryptAes := newEncryptAes[string]()
		toBeEncrypted := &ToBeEncrypted[string]{
			EncryptionType: ENCRYPT_TYPE_AES,
			Value:          "abc",
			PublicArg: map[string]any{
				"iv": base64.StdEncoding.EncodeToString(iv),
			},
			PrivateArg: map[string]any{
				"key": base64.StdEncoding.EncodeToString(key),
			},
		}
		val, err := encryptAes.Encrypt(toBeEncrypted)
		require.Nil(t, err)

		require.Equal(t, "WfIP", val.Value)

		decrypted := val.Decrypted.Get()
		require.Equal(t, "abc", decrypted)
	})

	t.Run("bool_true", func(t *testing.T) {
		encryptAes := newEncryptAes[bool]()
		toBeEncrypted := &ToBeEncrypted[bool]{
			EncryptionType: ENCRYPT_TYPE_AES,
			Value:          true,
			PublicArg: map[string]any{
				"iv": base64.StdEncoding.EncodeToString(iv),
			},
			PrivateArg: map[string]any{
				"key": base64.StdEncoding.EncodeToString(key),
			},
		}
		val, err := encryptAes.Encrypt(toBeEncrypted)
		require.Nil(t, err)

		require.Equal(t, "TOIZFw==", val.Value)

		decrypted := val.Decrypted.Get()
		require.Equal(t, true, decrypted)
	})

	t.Run("bool_false", func(t *testing.T) {
		encryptAes := newEncryptAes[bool]()
		toBeEncrypted := &ToBeEncrypted[bool]{
			EncryptionType: ENCRYPT_TYPE_AES,
			Value:          false,
			PublicArg: map[string]any{
				"iv": base64.StdEncoding.EncodeToString(iv),
			},
			PrivateArg: map[string]any{
				"key": base64.StdEncoding.EncodeToString(key),
			},
		}
		val, err := encryptAes.Encrypt(toBeEncrypted)
		require.Nil(t, err)

		require.Equal(t, "XvEAAfw=", val.Value)

		decrypted := val.Decrypted.Get()
		require.Equal(t, false, decrypted)
	})

	t.Run("bool_uint8", func(t *testing.T) {
		encryptAes := newEncryptAes[uint8]()
		toBeEncrypted := &ToBeEncrypted[uint8]{
			EncryptionType: ENCRYPT_TYPE_AES,
			Value:          1,
			PublicArg: map[string]any{
				"iv": base64.StdEncoding.EncodeToString(iv),
			},
			PrivateArg: map[string]any{
				"key": base64.StdEncoding.EncodeToString(key),
			},
		}
		val, err := encryptAes.Encrypt(toBeEncrypted)
		require.Nil(t, err)

		require.Equal(t, "CQ==", val.Value)

		decrypted := val.Decrypted.Get()
		require.Equal(t, uint8(1), decrypted)
	})
}
