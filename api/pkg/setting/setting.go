package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

type App struct {
	JwtSecret string
	PageSize  int
	PrefixUrl string

	RuntimeRootPath string

	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExts []string

	ExportSavePath string
	QrCodeSavePath string
	FontSavePath   string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string

	KeyVaultURL string
}

type MicrosoftEntraIDConfig struct {
	// Microsoft Entra ID (Azure AD)
	ClientID     string `json:"CLIENT_ID"`
	ClientSecret string `json:"CLIENT_SECRET"`
	TenantID     string `json:"TENANT_ID"`
	RedirectURL  string `json:"REDIRECT_URL"`

	// Application URLs
	FrontendURL string `json:"FRONTEND_URL"`

	// openapi-api-key
	OpenapiApiKey string `json:"openapi_api_key"`
	// OPENAPI-API-ENDPOINT
	OpenapiApiEndpoint string `json:"openapi_api_endpoint"`
	// OPENAPI-API-DEPLOYMENT-NAME deployment_name
	OpenapiApiDeploymentName string `json:"openapi_api_deployment_name"`
	// OPENAPI-API-VERSION api_version
	OpenapiApiVersion string `json:"openapi_api_version"`
	OpenapiUseAzureAD bool   `json:"openapi_use_azure_ad"`

	// blob
	BlobAccessKey     string `json:"blob_access_key"`
	BlobContainerName string `json:"blob_container_name"`
	BlobAccountName   string `json:"blob_account_name"`

	// Redis
	RedisURL string `json:"REDIS_URL"`

	// Session
	SessionSecret string `json:"SESSION_SECRET"`

	// Authorization
	AllowedGroups []string `json:"ALLOWED_GROUPS"`
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Cert        string
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &Database{}

type Redis struct {
	Host         string
	Password     string
	MaxIdle      int
	MaxActive    int
	IdleTimeout  time.Duration
	DB           int
	MaxRetries   int
	PoolSize     int
	PoolTimeout  time.Duration
	Prefix       string
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var RedisSetting = &Redis{}

var cfg *ini.File

// Setup initialize the configuration instance
func Setup() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
		return
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)

	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
	RedisSetting.PoolTimeout = RedisSetting.PoolTimeout * time.Second
	RedisSetting.DialTimeout = RedisSetting.DialTimeout * time.Second
	RedisSetting.ReadTimeout = RedisSetting.ReadTimeout * time.Second
	RedisSetting.WriteTimeout = RedisSetting.WriteTimeout * time.Second
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
