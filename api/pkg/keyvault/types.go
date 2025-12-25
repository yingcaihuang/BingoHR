package keyvault

import (
	"time"
)

// SecretOptions 秘密选项
type SecretOptions struct {
	// 过期时间
	Expires time.Time

	// 激活时间（在此之前不可用）
	NotBefore time.Time

	// 是否启用
	Enabled *bool

	// 内容类型
	ContentType string

	// 标签
	Tags map[string]*string
}

// SecretInfo 秘密信息
type SecretInfo struct {
	ID          string
	Name        string
	Enabled     bool
	Created     time.Time
	Updated     time.Time
	ContentType string
}

// SecretVersionInfo 秘密版本信息
type SecretVersionInfo struct {
	ID      string
	Version string
	Enabled bool
	Created time.Time
	Updated time.Time
	Expires *time.Time
}

// DeletedSecretInfo 已删除的秘密信息
type DeletedSecretInfo struct {
	ID                 string
	Name               string
	Enabled            bool
	RecoveryID         string
	DeletedDate        *time.Time
	ScheduledPurgeDate *time.Time
}
