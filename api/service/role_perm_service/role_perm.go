package role_perm_service

import (
	"context"

	"hr-api/models"
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
		datas []*models.RolePerm
		err   error
	)

	datas, err = models.GetRolePerms(r.RoleId)
	if err != nil {
		return nil, err
	}

	return datas, nil
}

func (r *RolePerm) Count() (int, error) {
	return models.GetRolePermsTotal(r.RoleId)
}
