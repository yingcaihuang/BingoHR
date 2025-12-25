package keyvault

import (
	"fmt"
	"hr-api/pkg/setting"
	"strings"
)

type ConfigLoader struct {
	keyVaultClient *KeyVaultClient
}

func NewConfigLoader() (*ConfigLoader, error) {
	client, err := NewKeyVaultClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create key vault client: %v", err)
	}

	return &ConfigLoader{
		keyVaultClient: client,
	}, nil
}

func (cl *ConfigLoader) LoadConfig() (*setting.MicrosoftEntraIDConfig, error) {
	config := &setting.MicrosoftEntraIDConfig{}

	// 定义需要从 Key Vault 加载的配置项
	secrets := []string{
		"CLIENT-ID",
		"CLIENT-SECRET",
		"TENANT-ID",
		"REDIRECT-URL",
		"FRONTEND-URL",
		"REDIS-URL",
		"SESSION-SECRET",
		"ALLOWED-GROUPS",
		"OPENAPI-API-KEY",
		"OPENAPI-API-ENDPOINT",
		"OPENAPI-API-DEPLOYMENT-NAME",
		"OPENAPI-API-VERSION",
		"BLOB-ACCOUNT-NAME",
		"BLOB-CONTAINER-NAME",
		"BLOB-ACCESS-KEY",
	}

	// 从 Key Vault 加载每个秘密
	for _, secretName := range secrets {
		value, err := cl.keyVaultClient.GetSecret(secretName)
		if err != nil {
			//log.Printf("Warning: Failed to load secret '%s': %v", secretName, err)
			continue
		}

		// 根据秘密名称设置配置
		switch secretName {
		case "CLIENT-ID":
			config.ClientID = value
		case "CLIENT-SECRET":
			config.ClientSecret = value
		case "TENANT-ID":
			config.TenantID = value
		case "REDIRECT-URL":
			config.RedirectURL = value
		case "FRONTEND-URL":
			config.FrontendURL = value
		case "REDIS-URL":
			config.RedisURL = value
		case "SESSION-SECRET":
			config.SessionSecret = value
		case "ALLOWED-GROUPS":
			config.AllowedGroups = strings.Split(value, ",")
		case "OPENAPI-API-KEY":
			config.OpenapiApiKey = value
		case "OPENAPI-API-ENDPOINT":
			config.OpenapiApiEndpoint = value
		case "OPENAPI-API-DEPLOYMENT-NAME":
			config.OpenapiApiDeploymentName = value
		case "OPENAPI-API-VERSION":
			config.OpenapiApiVersion = value
		case "BLOB-ACCOUNT-NAME":
			config.BlobAccountName = value
		case "BLOB-CONTAINER-NAME":
			config.BlobContainerName = value
		case "BLOB-ACCESS-KEY":
			config.BlobAccessKey = value
		}
	}

	// 验证必要配置
	if err := cl.validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func (cl *ConfigLoader) validateConfig(config *setting.MicrosoftEntraIDConfig) error {
	var missing []string

	if config.ClientID == "" {
		missing = append(missing, "CLIENT-ID")
	}
	if config.ClientSecret == "" {
		missing = append(missing, "CLIENT-SECRET")
	}
	if config.TenantID == "" {
		missing = append(missing, "TENANT-ID")
	}
	if config.SessionSecret == "" {
		missing = append(missing, "SESSION-SECRET")
	}
	if config.OpenapiApiKey == "" {
		missing = append(missing, "OPENAPI-API-KEY")
	}
	if config.OpenapiApiEndpoint == "" {
		missing = append(missing, "OPENAPI-API-ENDPOINT")
	}
	if config.OpenapiApiDeploymentName == "" {
		missing = append(missing, "OPENAPI-API-DEPLOYMENT-NAME")
	}
	if config.OpenapiApiVersion == "" {
		missing = append(missing, "OPENAPI-API-VERSION")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required configuration: %s", strings.Join(missing, ", "))
	}

	return nil
}
