package user_service

import (
	"context"
	"time"

	"hr-api/models"
	"hr-api/pkg/util"
)

type User struct {
	Id       int
	Username string
	Password string
	Email    string
	// 是否来自Microsoft Entra ID登录, 0否1是
	IsMicrosoft int
	Roles       []int
	CreateUid   int
	CreateTime  int
	UpdateTime  int
	Ctx         context.Context

	Page       int
	Limit      int
	CacheClear int
}

func (u *User) Add() error {
	user := map[string]interface{}{
		"username":     u.Username,
		"password":     util.EncodeMD5(u.Password),
		"email":        u.Email,
		"create_uid":   u.CreateUid,
		"is_microsoft": u.IsMicrosoft,
	}
	return models.AddUser(user, u.Roles)
}

func (u *User) Edit() error {
	data := make(map[string]interface{})
	data["email"] = u.Email
	if u.Password != "" {
		data["password"] = util.EncodeMD5(u.Password)
	}
	data["update_time"] = int(time.Now().Unix())

	return models.EditUser(u.Id, data, u.Roles)
}

func (u *User) Delete() error {
	return models.DeleteUser(u.Id)
}

func (u *User) Count() (int, error) {
	return models.GetUserTotal(u.Username, u.getMaps())
}

func (u *User) ExistUserByUsername() (bool, error) {
	return models.ExistUserByUsername(u.Username)
}

func (u *User) ExistByID() (bool, error) {
	return models.ExistUserByID(u.Id)
}

func (u *User) GetAll() ([]*models.User, error) {
	var (
		datas []*models.User
		err   error
	)

	datas, err = models.GetUsers(u.Page, u.Limit, u.Username, u.getMaps())
	if err != nil {
		return nil, err
	}

	return datas, nil
}

func (u *User) GetUser() (*models.User, error) {
	return models.GetUser(u.Id)
}

func (u *User) GetUserPerms() []string {
	return models.GetUserPerms(u.Id)
}

func (u *User) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})

	if u.Id > 0 {
		maps["id"] = u.Id
	}

	return maps
}
