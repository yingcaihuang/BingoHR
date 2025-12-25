package role_perm_service

import (
	"context"
	"encoding/json"
	"hr-api/pkg/cache"
	"time"

	"hr-api/models"
	"hr-api/service/cache_service"
)

type RolePerm struct {
	RoleId     int
	Perms      []string
	CreateUid  int
	CacheClear int
	Ctx        context.Context
}

func (r *RolePerm) Add() error {
	return models.AddRolePerms(r.RoleId, r.Perms, r.CreateUid)
}

func (r *RolePerm) GetAll() ([]*models.RolePerm, error) {
	var (
		datas, cacheDatas []*models.RolePerm
		err               error
	)

	service := cache_service.Cache{
		Name: cache.CACHE_ROLE_PERM,
		Id:   r.RoleId,
	}

	rd, err := cache.GetInstance()
	if err != nil {
		return []*models.RolePerm{}, nil
	}

	key := service.GetTagsKey()
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

	datas, err = models.GetRolePerms(r.RoleId)
	if err != nil {
		return nil, err
	}

	if len(datas) > 0 {
		rd.Set(r.Ctx, key, datas, 3600*time.Second)
	}
	return datas, nil
}

func (r *RolePerm) Count() (int, error) {
	return models.GetRolePermsTotal(r.RoleId)
}
