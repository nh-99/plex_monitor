package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"plex_monitor/internal/secrets"
	"sync"

	"github.com/sirupsen/logrus"
)

// AES is the struct that is used to encrypt and decrypt strings using AES encryption.
type AES struct {
	AESKey []byte
}

var (
	aesEncryption    *AES
	setupSecretsOnce sync.Once
)

// NewAESFromSecrets returns a new AES struct with the AES key set from the secrets manager.
func NewAESFromSecrets() *AES {
	setupSecretsOnce.Do(func() {
		secretsManager := secrets.NewEnvSecretManager()
		if secretsManager == nil {
			logrus.Error("failed to create secrets manager")
			return
		}

		aesKey, err := secretsManager.GetSecret("AES_KEY")
		if err != nil {
			logrus.Errorf("failed to get AES key: %v", err)
			return
		}

		if aesKey == "" {
			logrus.Error("AES key is not set - please set it in the secrets manager with the key 'AES_KEY'")
			return
		}

		aesEncryption = &AES{
			AESKey: []byte(aesKey),
		}
	})

	return aesEncryption
}

// IsAESKeySet returns true if the AES key is set.
func (a *AES) IsAESKeySet() bool {
	return len(a.AESKey) > 0
}

// SetAESKey sets the AES key.
func (a *AES) SetAESKey(key string) {
	a.AESKey = []byte(key)
}

// EncryptStringAES encrypts a string using AES encryption.
func (a *AES) EncryptStringAES(s string) (string, error) {
	aesBlock, err := aes.NewCipher(a.AESKey)
	if err != nil {
		logrus.Errorf("Failed to create AES cipher: %v", err)
		return "", err
	}

	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		logrus.Errorf("Failed to create GCM: %v", err)
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		logrus.Errorf("Failed to read nonce: %v", err)
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(s), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptStringAES decrypts a string using AES encryption.
func (a *AES) DecryptStringAES(s string) (string, error) {
	aesBlock, err := aes.NewCipher(a.AESKey)
	if err != nil {
		logrus.Errorf("Failed to create AES cipher: %v", err)
		return "", err
	}

	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		logrus.Errorf("Failed to create GCM: %v", err)
		return "", err
	}

	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		logrus.Errorf("Failed to decode string: %v", err)
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := decoded[:nonceSize], decoded[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		logrus.Errorf("Failed to decrypt string: %v", err)
		return "", err
	}

	return string(plaintext), nil
}

func mdHashing(input string) string {
	byteInput := []byte(input)
	md5Hash := md5.Sum(byteInput)
	return hex.EncodeToString(md5Hash[:]) // by referring to it as a string
}
