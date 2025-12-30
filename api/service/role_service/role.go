package role_service

import (
	"context"
	"time"

	"hr-api/models"
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
		datas []*models.Role
		err   error
	)

	datas, err = models.GetRoles(r.Page, r.Limit, r.Name, r.getMaps())
	if err != nil {
		return nil, err
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
