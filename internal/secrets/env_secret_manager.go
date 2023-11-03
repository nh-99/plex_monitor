package secrets

import "os"

// EnvSecretManager is a struct that implements the SecretManager interface.
type EnvSecretManager struct{}

// GetSecret returns the secret for the given key.
func (esm *EnvSecretManager) GetSecret(key string) (string, error) {
	return os.Getenv(key), nil
}

// GetSecretOrDefault returns the secret for the given key or the default
// value if the secret is not found.
func (esm *EnvSecretManager) GetSecretOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// NewEnvSecretManager creates a new EnvSecretManager.
func NewEnvSecretManager() *EnvSecretManager {
	return &EnvSecretManager{}
}
