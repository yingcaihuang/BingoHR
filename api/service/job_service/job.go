package job_service

import (
	"context"
	"time"

	"hr-api/models"
)

type Job struct {
	Id         int
	Name       string
	Demand     string
	Desc       string
	CreateUid  int
	CreateTime int
	UpdateTime int
	Ctx        context.Context

	Page       int
	Limit      int
	CacheClear int
}

func (j *Job) Add() error {
	job := map[string]interface{}{
		"name":       j.Name,
		"demand":     j.Demand,
		"desc":       j.Desc,
		"create_uid": j.CreateUid,
	}
	return models.AddJob(job)
}

func (j *Job) Edit() error {
	data := make(map[string]interface{})
	data["name"] = j.Name
	data["demand"] = j.Demand
	data["desc"] = j.Desc
	data["update_time"] = int(time.Now().Unix())

	return models.EditJob(j.Id, data)
}

func (r *Job) Delete() error {
	return models.DeleteJob(r.Id)
}

func (r *Job) Count() (int, error) {
	return models.GetJobTotal(r.Name, r.getMaps())
}

func (r *Job) ExistByID() (bool, error) {
	return models.ExistJobByID(r.Id)
}

func (r *Job) GetAll() ([]*models.Job, error) {
	var (
		datas []*models.Job
		err   error
	)

	datas, err = models.GetJobs(r.Page, r.Limit, r.Name, r.getMaps())
	if err != nil {
		return nil, err
	}

	return datas, nil
}

func (j *Job) GetJob() (*models.Job, error) {
	return models.GetJob(j.Id)
}

func (j *Job) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})

	if j.Id > 0 {
		maps["id"] = j.Id
	}

	return maps
}
