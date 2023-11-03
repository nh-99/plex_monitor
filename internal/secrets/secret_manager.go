package secrets

// This is a thin abstraction for a secrets manager that allows implementing
// a secret manager in AWS Secrets Manager or the local environment.

// SecretManager is an interface that all secret managers must implement.
type SecretManager interface {
	// GetSecret returns the secret for the given key.
	GetSecret(key string) (string, error)
	// GetSecretOrDefault returns the secret for the given key or the default
	// value if the secret is not found.
	GetSecretOrDefault(key string, defaultValue string) string
}

// SecretManagerImpl is the implementation of the SecretManager interface.
type SecretManagerImpl struct {
	// SecretManager is the secret manager that is being used.
	SecretManager SecretManager
}

// GetSecret returns the secret for the given key.
func (smi *SecretManagerImpl) GetSecret(key string) (string, error) {
	return smi.SecretManager.GetSecret(key)
}

// NewSecretManager creates a new SecretManager.
func NewSecretManager(secretManager SecretManager) *SecretManagerImpl {
	return &SecretManagerImpl{SecretManager: secretManager}
}

// GetSecret returns the secret for the given key.
func GetSecret(key string) (string, error) {
	return secretManager.GetSecret(key)
}

// GetSecretOrDefault returns the secret for the given key or the default
// value if the secret is not found.
func GetSecretOrDefault(key string, defaultValue string) string {
	return secretManager.GetSecretOrDefault(key, defaultValue)
}

// secretManager is the singleton secret manager.
var secretManager SecretManager

// SetSecretManager sets the secret manager.
func SetSecretManager(sm SecretManager) {
	secretManager = sm
}

// GetSecretManager returns the secret manager.
func GetSecretManager() SecretManager {
	return secretManager
}
