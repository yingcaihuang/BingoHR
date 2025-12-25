package models

import (
	"strings"
	"time"
)

type RolePerm struct {
	ID         int    `json:"id" gorm:"primaryKey"`
	RoleId     int    `json:"role_id"`
	Value      string `json:"value"`
	CreateUid  int    `json:"create_uid"`
	CreateTime int    `json:"create_time"`
}

// GetRolePerms get role perms list data base by role_id
func GetRolePerms(role_id int) ([]*RolePerm, error) {
	var rolePerms []*RolePerm

	err := db.Model(&RolePerm{}).Where("role_id = ?", role_id).First(&rolePerms).Error
	if err != nil {
		return nil, err
	}

	return rolePerms, nil
}

// GetRoleTotal counts the total number of roles based on the constraint
func GetRolePermsTotal(role_id int) (int, error) {
	var count int64

	err := db.Model(&RolePerm{}).Where("role_id = ?", role_id).Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// AddRole add a single role
func AddRolePerms(role_id int, perms []string, CreateUid int) error {
	// 先删除现有的权限
	var err error

	err = DeleteRolePerms(role_id)
	if err != nil {
		return err
	}

	if len(perms) == 0 {
		return nil
	}

	now := int(time.Now().Unix())
	var rolePerms []RolePerm
	for _, p := range perms {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		rolePerms = append(rolePerms, RolePerm{
			RoleId:     role_id,
			Value:      p,
			CreateUid:  CreateUid,
			CreateTime: now,
		})
	}

	if err = db.Create(&rolePerms).Error; err != nil {
		return err
	}

	return nil
}

// DeleteRolePerms delete a role's all perms based by role_id
func DeleteRolePerms(role_id int) error {
	if err := db.Where("role_id = ?", role_id).Delete(RolePerm{}).Error; err != nil {
		return err
	}

	return nil
}
