package encrypt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncryptUnencrypted(t *testing.T) {

	t.Run("string", func(t *testing.T) {
		encryptAes := newEncryptUnencrypted[string]()
		toBeEncrypted := &ToBeEncrypted[string]{
			EncryptionType: ENCRYPT_TYPE_UNENCRYPTED,
			Value:          "abc",
		}
		val, err := encryptAes.Encrypt(toBeEncrypted)
		require.Nil(t, err)

		decrypted := val.Decrypted.Get()
		require.Equal(t, "abc", decrypted)
	})

	t.Run("bool", func(t *testing.T) {
		t.Run("true", func(t *testing.T) {
			encryptAes := newEncryptUnencrypted[bool]()
			toBeEncrypted := &ToBeEncrypted[bool]{
				EncryptionType: ENCRYPT_TYPE_UNENCRYPTED,
				Value:          true,
			}
			val, err := encryptAes.Encrypt(toBeEncrypted)
			require.Nil(t, err)

			decrypted := val.Decrypted.Get()
			require.Equal(t, true, decrypted)
		})

		t.Run("false", func(t *testing.T) {
			encryptAes := newEncryptUnencrypted[bool]()
			toBeEncrypted := &ToBeEncrypted[bool]{
				EncryptionType: ENCRYPT_TYPE_UNENCRYPTED,
				Value:          false,
			}
			val, err := encryptAes.Encrypt(toBeEncrypted)
			require.Nil(t, err)

			decrypted := val.Decrypted.Get()
			require.Equal(t, false, decrypted)
		})
	})

	t.Run("uint8", func(t *testing.T) {
		encryptAes := newEncryptUnencrypted[uint8]()
		toBeEncrypted := &ToBeEncrypted[uint8]{
			EncryptionType: ENCRYPT_TYPE_UNENCRYPTED,
			Value:          1,
		}
		val, err := encryptAes.Encrypt(toBeEncrypted)
		require.Nil(t, err)

		decrypted := val.Decrypted.Get()
		require.Equal(t, uint8(1), decrypted)
	})
}
