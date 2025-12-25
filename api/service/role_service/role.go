package role_service

import (
	"context"
	"encoding/json"
	"hr-api/pkg/cache"
	"time"

	"hr-api/models"
	"hr-api/service/cache_service"
)

type Role struct {
	Id         int
	Name       string
	CreateUid  int
	CreateTime int
	UpdateTime int
	Ctx        context.Context

	Page       int
	Limit      int
	CacheClear int
}

func (r *Role) Add() error {
	return models.AddRole(r.Name, r.CreateUid)
}

func (r *Role) Edit() error {
	data := make(map[string]interface{})
	data["name"] = r.Name
	data["update_time"] = int(time.Now().Unix())

	return models.EditRole(r.Id, data)
}

func (r *Role) Delete() error {
	return models.DeleteRole(r.Id)
}

func (r *Role) Count() (int, error) {
	return models.GetRoleTotal(r.Name, r.getMaps())
}

func (r *Role) ExistByID() (bool, error) {
	return models.ExistRoleByID(r.Id)
}

func (r *Role) GetAll() ([]*models.Role, error) {
	var (
		datas, cacheDatas []*models.Role
		err               error
	)

	cacheService := cache_service.Cache{
		Name:    cache.CACHE_ROLE,
		Keyword: r.Name,

		Page:  r.Page,
		Limit: r.Limit,
	}

	rd, err := cache.GetInstance()
	if err != nil {
		return []*models.Role{}, nil
	}

	key := cacheService.GetTagsKey()
	if r.CacheClear > 0 {
		rd.Delete(r.Ctx, key)
	}

	exist, _ := rd.Exists(r.Ctx, key)
	if r.CacheClear == 0 && exist {
		var cacheData string
		if err := rd.Get(r.Ctx, key, &cacheData); err != nil {
			return nil, err
		} else {
			json.Unmarshal([]byte(cacheData), &cacheDatas)
			return cacheDatas, nil
		}
	}

	datas, err = models.GetRoles(r.Page, r.Limit, r.Name, r.getMaps())
	if err != nil {
		return nil, err
	}

	if len(datas) > 0 {
		rd.Set(r.Ctx, key, datas, 3600*time.Second)
	}
	return datas, nil
}

func (r *Role) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})

	if r.Id > 0 {
		maps["id"] = r.Id
	}

	return maps
}
