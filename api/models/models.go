package models

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"time"

	"hr-api/pkg/setting"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn  int `json:"deleted_on"`
}

// Setup initializes the database instance
func Setup() {
	var err error

	baseDSN := "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	if len(setting.DatabaseSetting.Cert) > 0 {
		caCert, err := os.ReadFile(setting.DatabaseSetting.Cert)
		if err != nil {
			log.Fatalf("models.Setup err: %v", err)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			log.Fatalf("models.Setup err: %v", err)
		}

		tlsConfig := &tls.Config{
			RootCAs:    caCertPool,
			ServerName: "webtest-server.mysql.database.azure.com",
			MinVersion: tls.VersionTLS12,
		}

		err = mysqlDriver.RegisterTLSConfig("mysql-cert", tlsConfig)
		if err != nil {
			log.Fatalf("models.Setup err: %v", err)
		}
		baseDSN += "&tls=mysql-cert"
	}

	dsn := fmt.Sprintf(baseDSN,
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		}})
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get raw database: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

}

// CloseDB closes database connection (unnecessary)
func CloseDB() error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get raw database: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	db = nil
	return nil
}
