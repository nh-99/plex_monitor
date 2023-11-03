package secrets

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// AwsSecretManager is a struct that implements the SecretManager interface.
type AwsSecretManager struct {
	secretsManager *secretsmanager.SecretsManager
}

// GetSecret returns the secret for the given key.
func (asm *AwsSecretManager) GetSecret(key string) (string, error) {
	result, err := asm.secretsManager.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: &key})
	if err != nil {
		log.Fatal(err.Error())
	}
	return *result.SecretString, nil
}

// GetSecretOrDefault returns the secret for the given key or the default
// value if the secret is not found.
func (asm *AwsSecretManager) GetSecretOrDefault(key string, defaultValue string) string {
	result, err := asm.secretsManager.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: &key})
	if err != nil {
		log.Fatal(err.Error())
		return defaultValue
	}
	if result.SecretString == nil {
		return defaultValue
	}

	return *result.SecretString
}

// NewAwsSecretManager creates a new AwsSecretManager.
func NewAwsSecretManager(region string) *AwsSecretManager {
	sess := session.Must(session.NewSession())
	svc := secretsmanager.New(sess, aws.NewConfig().WithRegion(region))

	return &AwsSecretManager{
		secretsManager: svc,
	}
}
