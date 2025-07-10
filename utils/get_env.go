package utils

import (
	"fmt"
	"os"
	"sync"

	"github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"
)

var (
	logLevel    = os.Getenv("LOG_LEVEL")
	vaultOnce   sync.Once
	vaultClient *api.Client
)

func MustGetEnv(key string) string {
	if valueDotEnv, err := getDotEnv(key); err == nil && len(valueDotEnv) > 0 {
		if logLevel == "DEBUG" {
			fmt.Printf("found environment variable %s: %s", key, valueDotEnv)
		}
		return valueDotEnv
	}

	panic(fmt.Sprintf("missing environment variable %s", key))
}

func getDotEnv(key string) (string, error) {
	if err := godotenv.Load(); err != nil {
		err := fmt.Errorf("missing dotenv variable %s from .env", key)
		return "", err
	}

	value, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("no value for dotenv variable %s from .env", key)
	}

	return value, nil
}

func MustGetVaultEnv(key string) string {
	vaultValue, err := getFromVault(key)
	if err != nil || vaultValue == "" {
		panic(fmt.Sprintf("missing required environment variable '%s' from Vault: %v", key, err))
	}
	debugLog(key, vaultValue, "vault")
	return vaultValue
}

func debugLog(key, value, source string) {
	if logLevel == "DEBUG" {
		fmt.Printf("[DEBUG] found env %s from %s: %s\n", key, source, value)
	}
}

func getFromVault(key string) (string, error) {
	vaultOnce.Do(func() {
		config := api.DefaultConfig()
		config.Address = os.Getenv("VAULT_ADDR")

		client, err := api.NewClient(config)
		if err == nil {
			client.SetToken(os.Getenv("VAULT_TOKEN"))
			vaultClient = client
		}
	})

	if vaultClient == nil {
		return "", fmt.Errorf("vault client is not initialized")
	}

	secretPath := os.Getenv("VAULT_SECRET_PATH")
	secret, err := vaultClient.Logical().Read(secretPath)
	if err != nil || secret == nil || secret.Data == nil {
		return "", fmt.Errorf("vault: key '%s' not found in path '%s'", key, secretPath)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("vault: secret at path '%s' missing 'data'", secretPath)
	}

	value, ok := data[key].(string)
	if !ok || value == "" {
		return "", fmt.Errorf("vault: key '%s' is missing or not a string", key)
	}

	return value, nil
}
