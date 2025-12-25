package user_service

import (
	"context"
	"encoding/json"
	"hr-api/pkg/cache"
	"time"

	"hr-api/models"
	"hr-api/pkg/util"
	"hr-api/service/cache_service"
)

type User struct {
	Id         int
	Username   string
	Password   string
	Email      string
	Roles      []int
	CreateUid  int
	CreateTime int
	UpdateTime int
	Ctx        context.Context

	Page       int
	Limit      int
	CacheClear int
}

func (u *User) Add() error {
	user := map[string]interface{}{
		"username":   u.Username,
		"password":   util.EncodeMD5(u.Password),
		"email":      u.Email,
		"create_uid": u.CreateUid,
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
		datas, cacheDatas []*models.User
		err               error
	)

	cacheService := cache_service.Cache{
		Name:    cache.CACHE_USER,
		Keyword: u.Username,
		Page:    u.Page,
		Limit:   u.Limit,
	}

	rd, err := cache.GetInstance()
	if err != nil {
		return []*models.User{}, nil
	}

	key := cacheService.GetTagsKey()
	if u.CacheClear > 0 {
		rd.Delete(u.Ctx, key)
	}

	exist, _ := rd.Exists(u.Ctx, key)
	if u.CacheClear == 0 && exist {
		var cacheData string
		if err := rd.Get(u.Ctx, key, &cacheData); err != nil {
			return nil, err
		} else {
			json.Unmarshal([]byte(cacheData), &cacheDatas)
			return cacheDatas, nil
		}
	}

	datas, err = models.GetUsers(u.Page, u.Limit, u.Username, u.getMaps())
	if err != nil {
		return nil, err
	}

	if len(datas) > 0 {
		rd.Set(u.Ctx, key, datas, 3600*time.Second)
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
