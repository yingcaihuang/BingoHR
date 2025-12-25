package models

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
)

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (Auth) TableName() string {
	return "users"
}

// CheckAuth checks if authentication information exists
func CheckAuth(username, password string) (int, error) {
	var auth Auth
	err := db.Select("id").Where(Auth{Username: username, Password: password}).First(&auth).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || (strings.HasSuffix(err.Error(), "record not found")) {
			return 0, nil
		}
		return 0, err
	}

	return auth.ID, nil
}
