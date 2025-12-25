package models

import (
	"errors"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         int    `json:"id" gorm:"primaryKey"`
	Username   string `json:"username"`
	Password   string `json:"-"`
	Email      string `json:"email"`
	CreateUid  int    `json:"create_uid"`
	CreateTime int    `json:"create_time"`
	UpdateTime int    `json:"update_time"`
	Roles      []Role `json:"roles" gorm:"-"`
}

type UserRole struct {
	ID         int `json:"id" gorm:"primaryKey"`
	Uid        int `json:"uid"`
	RoleId     int `json:"role_id"`
	CreateTime int `json:"create_time"`
}

// GetUsers get user list data
func GetUsers(page int, limit int, keyword string, maps interface{}) ([]*User, error) {
	var (
		users []*User
		err   error
	)

	query := db.Model(&User{}).Where(maps)

	if keyword != "" {
		query = query.Where("users.name LIKE ?", "%"+keyword+"%")
	}

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	err = query.Find(&users).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 循环用户数据查询其角色数据
	for i := range users {
		var roles []Role
		sql := `SELECT r.* FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.uid = ?`
		db.Raw(sql, users[i].ID).Scan(&roles)
		if roles == nil {
			roles = []Role{}
		}
		users[i].Roles = roles
	}

	return users, nil
}

// GetUserTotal gets the total number of users based on the constraints
func GetUserTotal(keyword string, maps interface{}) (int, error) {
	var count int64

	query := db.Model(&User{}).Where(maps)
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

// GetUser Get a single user based on ID
func GetUser(id int) (*User, error) {
	var d User
	err := db.Model(&User{}).Where("id = ? ", id).First(&d).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &d, nil
}

func GetUserByName(name string) (*User, error) {
	var user User
	err := db.Table("users").Where("name = ? ", name).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &user, nil
}

// AddUser Add a user
func AddUser(data map[string]interface{}, roles []int) error {
	// 先创建用户
	now := int(time.Now().Unix())
	user := User{
		Username:   data["username"].(string),
		Password:   data["password"].(string),
		Email:      data["email"].(string),
		CreateTime: now,
		CreateUid:  data["create_uid"].(int),
	}

	if err := db.Create(&user).Error; err != nil {
		return err
	}

	if len(roles) > 0 {
		// 关联角色ID
		var userRoles []UserRole
		for _, rid := range roles {
			userRoles = append(userRoles, UserRole{
				Uid:        user.ID,
				RoleId:     rid,
				CreateTime: now,
			})
		}
		if err := db.Create(&userRoles).Error; err != nil {
			return err
		}
	}

	return nil
}

// EditUser modify a single user
func EditUser(id int, data map[string]interface{}, roles []int) error {
	// 更新基础数据
	if len(data) > 0 {
		if err := db.Model(&User{}).Where("id = ? ", id).Updates(data).Error; err != nil {
			return err
		}
	}

	// 更新角色数据
	if len(roles) > 0 {
		if err := DeleteUserRoles(id); err != nil {
			return nil
		}

		now := int(time.Now().Unix())
		// 关联角色ID
		var userRoles []UserRole
		for _, rid := range roles {
			userRoles = append(userRoles, UserRole{
				Uid:        id,
				RoleId:     rid,
				CreateTime: now,
			})
		}
		if err := db.Create(&userRoles).Error; err != nil {
			return err
		}
	}

	return nil
}

// DeleteUser delete a single user
func DeleteUser(id int) error {
	if err := db.Where("id = ?", id).Delete(User{}).Error; err != nil {
		return err
	}

	return nil
}

// DeleteUserRoles Delete user's roles
func DeleteUserRoles(uid int) error {
	if err := db.Where("uid = ?", uid).Delete(UserRole{}).Error; err != nil {
		return err
	}

	return nil
}

// CleanAllUser clear all user
func CleanAllUser() error {
	if err := db.Unscoped().Delete(&User{}).Error; err != nil {
		return err
	}

	return nil
}

// GetUserPerms Get user's perms by his binded roles
func GetUserPerms(uid int) []string {
	// 先查询这个用户关联的角色
	var roles []Role
	sql := `SELECT r.* FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.uid = ?`
	db.Raw(sql, uid).Scan(&roles)
	if roles == nil {
		roles = []Role{}
	}

	if len(roles) == 0 {
		return []string{}
	}

	var role_ids []string
	for _, role := range roles {
		role_ids = append(role_ids, strconv.Itoa(role.ID))
	}

	// 根据角色查询权限
	var role_perms []RolePerm
	sql = `SELECT DISTINCT value FROM role_perms WHERE role_id IN ?`
	db.Raw(sql, role_ids).Scan(&role_perms)

	var perms []string
	for _, perm := range role_perms {
		perms = append(perms, perm.Value)
	}

	return perms
}

// ExistUserByUsername determines whether a user exists based on username
func ExistUserByUsername(username string) (bool, error) {
	var count int64
	err := db.Model(&User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistUserByID determines whether a user exists based on the ID
func ExistUserByID(id int) (bool, error) {
	var count int64
	err := db.Model(&User{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
