package encryption

// EncryptionType is the type of the encryption.
type EncryptionType string

const (
	// EncryptionTypeAES is the type of the AES encryption.
	EncryptionTypeAES EncryptionType = "AES"
)

// Encryption is the interface that all encryption types must implement.
type Encryption interface {
	EncryptString(s string) (string, error)
	DecryptString(s string) (string, error)
}
