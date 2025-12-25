package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Role struct {
	ID         int    `json:"id" gorm:"primaryKey"`
	Name       string `json:"name"`
	CreateUid  int    `json:"create_uid"`
	CreateTime int    `json:"create_time"`
	UpdateTime int    `json:"update_time"`
}

// GetRoles get role list data
func GetRoles(page int, limit int, keyword string, maps interface{}) ([]*Role, error) {
	var (
		roles []*Role
		err   error
	)

	query := db.Model(&Role{}).Where(maps)

	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	err = query.Find(&roles).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return roles, nil
}

// GetRoleTotal counts the total number of roles based on the constraint
func GetRoleTotal(keyword string, maps interface{}) (int, error) {
	var count int64

	query := db.Model(&Role{}).Where(maps)
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

// AddRole add a single role
func AddRole(name string, CreateUid int) error {
	role := Role{
		Name:       name,
		CreateUid:  CreateUid,
		CreateTime: int(time.Now().Unix()),
	}
	if err := db.Create(&role).Error; err != nil {
		return err
	}

	return nil
}

// EditRole modify a single role
func EditRole(id int, data interface{}) error {
	if err := db.Model(&Role{}).Where("id = ? ", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// DeleteRole delete a single role
func DeleteRole(id int) error {
	if err := db.Where("id = ?", id).Delete(Role{}).Error; err != nil {
		return err
	}

	return nil
}

// CleanAllRole clear all role
func CleanAllRole() error {
	if err := db.Unscoped().Delete(&Role{}).Error; err != nil {
		return err
	}

	return nil
}

// ExistRoleByID determines whether a role exists based on the ID
func ExistRoleByID(id int) (bool, error) {
	var count int64
	err := db.Model(&Role{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
