package keyvault

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"hr-api/pkg/setting"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

// 这些参数需要在keyvault中配置
//        "CLIENT-ID",
//        "CLIENT-SECRET",
//        "TENANT-ID": common,
//        "REDIRECT-URL",
//        "FRONTEND-URL",
//        "REDIS-URL",
//        "SESSION-SECRET",
//        "ALLOWED-GROUPS",

type KeyVaultClient struct {
	client *azsecrets.Client
	ctx    context.Context
}

func GetKeyVaultConf() (*setting.MicrosoftEntraIDConfig, error) {
	// 创建配置加载器
	loader, err := NewConfigLoader()
	if err != nil {
		return nil, fmt.Errorf("failed to create config loader: %v", err)
	}

	// 加载配置
	appConfig, err := loader.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	return appConfig, nil
}

func NewKeyVaultClient() (*KeyVaultClient, error) {
	keyVaultURL := setting.AppSetting.KeyVaultURL
	if len(keyVaultURL) == 0 {
		keyVaultURL = "https://hr-api-keyvault.vault.azure.net/"
	}

	ctx := context.Background()

	// 使用 DefaultAzureCredential，支持多种身份验证方式
	// 1. 环境变量 (AZURE_CLIENT_ID, AZURE_CLIENT_SECRET, AZURE_TENANT_ID)
	// 2. Managed Identity
	// 3. Azure CLI
	// 4. 等等...
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %v", err)
	}

	client, err := azsecrets.NewClient(keyVaultURL, credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create key vault client: %v", err)
	}

	return &KeyVaultClient{
		client: client,
		ctx:    ctx,
	}, nil
}

// GetSecret 获取秘密
func (kvc *KeyVaultClient) GetSecret(secretName string) (string, error) {
	return kvc.GetSecretWithVersion(secretName, "")
}

// GetSecretWithVersion 获取指定版本的秘密
func (kvc *KeyVaultClient) GetSecretWithVersion(secretName, version string) (string, error) {
	resp, err := kvc.client.GetSecret(kvc.ctx, secretName, version, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get secret '%s' (version: %s): %v",
			secretName, version, err)
	}

	if resp.Value == nil {
		return "", fmt.Errorf("secret '%s' has no value", secretName)
	}

	return *resp.Value, nil
}

// SetSecret 设置秘密（创建或更新）
func (kvc *KeyVaultClient) SetSecret(secretName, secretValue string) (string, error) {
	return kvc.SetSecretWithOptions(secretName, secretValue, nil)
}

// SetSecretWithOptions 设置秘密并指定选项
func (kvc *KeyVaultClient) SetSecretWithOptions(secretName, secretValue string,
	options *SecretOptions) (string, error) {

	params := azsecrets.SetSecretParameters{
		Value: to.Ptr(secretValue),
	}

	if options != nil {
		if !options.Expires.IsZero() {
			params.SecretAttributes = &azsecrets.SecretAttributes{
				Expires: &options.Expires,
			}
		}
		if !options.NotBefore.IsZero() {
			if params.SecretAttributes == nil {
				params.SecretAttributes = &azsecrets.SecretAttributes{}
			}
			params.SecretAttributes.NotBefore = &options.NotBefore
		}
		if options.Enabled != nil {
			if params.SecretAttributes == nil {
				params.SecretAttributes = &azsecrets.SecretAttributes{}
			}
			params.SecretAttributes.Enabled = options.Enabled
		}
		if options.ContentType != "" {
			params.ContentType = &options.ContentType
		}
		if len(options.Tags) > 0 {
			params.Tags = options.Tags
		}
	}

	resp, err := kvc.client.SetSecret(kvc.ctx, secretName, params, nil)
	if err != nil {
		return "", fmt.Errorf("failed to set secret '%s': %v", secretName, err)
	}

	if resp.ID == nil {
		return "", fmt.Errorf("secret ID is nil for secret '%s'", secretName)
	}

	return string(*resp.ID), nil
}

// UpdateSecret 更新秘密属性（不改变值）
func (kvc *KeyVaultClient) UpdateSecret(secretName string, options *SecretOptions) error {
	if options == nil {
		return fmt.Errorf("options cannot be nil for update")
	}

	secret, err := kvc.client.GetSecret(kvc.ctx, secretName, "", nil)
	if err != nil {
		return fmt.Errorf("failed to get secret for update '%s': %v", secretName, err)
	}

	params := azsecrets.UpdateSecretParameters{}

	if options.ContentType != "" {
		params.ContentType = &options.ContentType
	}

	if !options.Expires.IsZero() || !options.NotBefore.IsZero() || options.Enabled != nil {
		params.SecretAttributes = &azsecrets.SecretAttributes{}

		if !options.Expires.IsZero() {
			params.SecretAttributes.Expires = &options.Expires
		} else if secret.Attributes != nil && secret.Attributes.Expires != nil {
			params.SecretAttributes.Expires = secret.Attributes.Expires
		}

		if !options.NotBefore.IsZero() {
			params.SecretAttributes.NotBefore = &options.NotBefore
		} else if secret.Attributes != nil && secret.Attributes.NotBefore != nil {
			params.SecretAttributes.NotBefore = secret.Attributes.NotBefore
		}

		if options.Enabled != nil {
			params.SecretAttributes.Enabled = options.Enabled
		} else if secret.Attributes != nil && secret.Attributes.Enabled != nil {
			params.SecretAttributes.Enabled = secret.Attributes.Enabled
		}
	}

	if len(options.Tags) > 0 {
		params.Tags = options.Tags
	} else if secret.Tags != nil {
		params.Tags = secret.Tags
	}

	_, err = kvc.client.UpdateSecret(kvc.ctx, secretName, "", params, nil)
	if err != nil {
		return fmt.Errorf("failed to update secret '%s': %v", secretName, err)
	}

	return nil
}

// DeleteSecret 删除秘密
func (kvc *KeyVaultClient) DeleteSecret(secretName string) error {
	// 开始删除操作
	_, err := kvc.client.DeleteSecret(kvc.ctx, secretName, nil)
	if err != nil {
		return fmt.Errorf("failed to begin deletion of secret '%s': %v", secretName, err)
	}

	// 注意：BeginDeleteSecret 是异步操作
	// 返回后删除操作会在后台进行
	return nil
}

// DeleteAndPurgeSecret 删除并永久清除秘密（需要启用软删除和清除保护）
func (kvc *KeyVaultClient) DeleteAndPurgeSecret(secretName string) error {
	// 开始删除
	_, err := kvc.client.PurgeDeletedSecret(kvc.ctx, secretName, nil)
	if err != nil {
		return fmt.Errorf("failed to begin deletion of secret '%s': %v", secretName, err)
	}

	// 等待删除完成
	//_, err = poller.PollUntilDone(kvc.ctx, nil)
	//if err != nil {
	//	return fmt.Errorf("failed to delete secret '%s': %v", secretName, err)
	//}

	// 永久清除已删除的秘密
	_, err = kvc.client.PurgeDeletedSecret(kvc.ctx, secretName, nil)
	if err != nil {
		return fmt.Errorf("failed to purge deleted secret '%s': %v", secretName, err)
	}

	return nil
}

// GetDeletedSecret 获取已删除的秘密信息
func (kvc *KeyVaultClient) GetDeletedSecret(secretName string) (*DeletedSecretInfo, error) {
	resp, err := kvc.client.GetDeletedSecret(kvc.ctx, secretName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get deleted secret '%s': %v", secretName, err)
	}

	info := &DeletedSecretInfo{
		Name: secretName,
	}

	if resp.ID != nil {
		info.ID = string(*resp.ID)
	}

	if resp.Attributes != nil && resp.Attributes.Enabled != nil {
		info.Enabled = *resp.Attributes.Enabled
	}

	if resp.RecoveryID != nil {
		info.RecoveryID = *resp.RecoveryID
	}

	if resp.DeletedDate != nil {
		info.DeletedDate = resp.DeletedDate
	}

	if resp.ScheduledPurgeDate != nil {
		info.ScheduledPurgeDate = resp.ScheduledPurgeDate
	}

	return info, nil
}

// RecoverDeletedSecret 恢复已删除的秘密
func (kvc *KeyVaultClient) RecoverDeletedSecret(secretName string) error {
	_, err := kvc.client.RecoverDeletedSecret(kvc.ctx, secretName, nil)
	if err != nil {
		return fmt.Errorf("failed to begin recovery of deleted secret '%s': %v", secretName, err)
	}

	//_, err = poller.PollUntilDone(kvc.ctx, nil)
	//if err != nil {
	//	return fmt.Errorf("failed to recover deleted secret '%s': %v", secretName, err)
	//}

	return nil
}

// ListSecrets 列出所有当前版本的秘密
func (kvc *KeyVaultClient) ListSecrets() ([]SecretInfo, error) {
	var secrets []SecretInfo

	pager := kvc.client.NewListSecretsPager(nil)
	for pager.More() {
		page, err := pager.NextPage(kvc.ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list secrets: %v", err)
		}

		for _, secret := range page.Value {
			info := SecretInfo{}

			if secret.ID != nil {
				info.ID = string(*secret.ID)
				info.Name = extractSecretName(info.ID)
			}

			if secret.Attributes != nil {
				if secret.Attributes.Enabled != nil {
					info.Enabled = *secret.Attributes.Enabled
				}
				if secret.Attributes.Created != nil {
					info.Created = *secret.Attributes.Created
				}
				if secret.Attributes.Updated != nil {
					info.Updated = *secret.Attributes.Updated
				}
			}

			if secret.ContentType != nil {
				info.ContentType = *secret.ContentType
			}

			secrets = append(secrets, info)
		}
	}

	return secrets, nil
}

// ListDeletedSecrets 列出所有已删除的秘密
func (kvc *KeyVaultClient) ListDeletedSecrets() ([]DeletedSecretInfo, error) {
	var deletedSecrets []DeletedSecretInfo

	pager := kvc.client.NewListDeletedSecretsPager(nil)
	for pager.More() {
		page, err := pager.NextPage(kvc.ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list deleted secrets: %v", err)
		}

		for _, secret := range page.Value {
			info := DeletedSecretInfo{}

			if secret.ID != nil {
				info.ID = string(*secret.ID)
				info.Name = extractSecretName(info.ID)
			}

			if secret.Attributes != nil && secret.Attributes.Enabled != nil {
				info.Enabled = *secret.Attributes.Enabled
			}

			if secret.RecoveryID != nil {
				info.RecoveryID = *secret.RecoveryID
			}

			if secret.DeletedDate != nil {
				info.DeletedDate = secret.DeletedDate
			}

			if secret.ScheduledPurgeDate != nil {
				info.ScheduledPurgeDate = secret.ScheduledPurgeDate
			}

			deletedSecrets = append(deletedSecrets, info)
		}
	}

	return deletedSecrets, nil
}

// ListSecretVersions 列出秘密的所有版本
func (kvc *KeyVaultClient) ListSecretVersions(secretName string) ([]SecretVersionInfo, error) {
	var versions []SecretVersionInfo

	pager := kvc.client.NewListSecretVersionsPager(secretName, nil)
	for pager.More() {
		page, err := pager.NextPage(kvc.ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list versions for secret '%s': %v", secretName, err)
		}

		for _, version := range page.Value {
			info := SecretVersionInfo{}

			if version.ID != nil {
				info.ID = string(*version.ID)
				info.Version = extractSecretVersion(info.ID)
			}

			if version.Attributes != nil {
				if version.Attributes.Enabled != nil {
					info.Enabled = *version.Attributes.Enabled
				}
				if version.Attributes.Created != nil {
					info.Created = *version.Attributes.Created
				}
				if version.Attributes.Updated != nil {
					info.Updated = *version.Attributes.Updated
				}
				if version.Attributes.Expires != nil {
					info.Expires = version.Attributes.Expires
				}
			}

			versions = append(versions, info)
		}
	}

	return versions, nil
}

// BackupSecret 备份秘密
func (kvc *KeyVaultClient) BackupSecret(secretName string) ([]byte, error) {
	resp, err := kvc.client.BackupSecret(kvc.ctx, secretName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to backup secret '%s': %v", secretName, err)
	}

	return resp.Value, nil
}

// RestoreSecret 恢复秘密
func (kvc *KeyVaultClient) RestoreSecret(backupData []byte) (string, error) {
	params := azsecrets.RestoreSecretParameters{
		SecretBundleBackup: backupData,
	}

	resp, err := kvc.client.RestoreSecret(kvc.ctx, params, nil)
	if err != nil {
		return "", fmt.Errorf("failed to restore secret: %v", err)
	}

	if resp.ID == nil {
		return "", fmt.Errorf("restored secret ID is nil")
	}

	return extractSecretName(string(*resp.ID)), nil
}

// 辅助函数
func extractSecretName(secretID string) string {
	parts := strings.Split(secretID, "/")
	return parts[len(parts)-1]
}

func extractSecretVersion(secretID string) string {
	parts := strings.Split(secretID, "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-2]
}
