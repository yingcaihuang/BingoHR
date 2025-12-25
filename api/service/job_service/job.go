package job_service

import (
	"context"
	"encoding/json"
	"hr-api/pkg/cache"
	"time"

	"hr-api/models"
	"hr-api/service/cache_service"
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
		datas, cacheDatas []*models.Job
		err               error
	)

	cacheService := cache_service.Cache{
		Name:    cache.CACHE_JOB,
		Keyword: r.Name,

		Page:  r.Page,
		Limit: r.Limit,
	}

	rd, err := cache.GetInstance()
	if err != nil {
		return []*models.Job{}, nil
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

	datas, err = models.GetJobs(r.Page, r.Limit, r.Name, r.getMaps())
	if err != nil {
		return nil, err
	}

	if len(datas) > 0 {
		rd.Set(r.Ctx, key, datas, 3600*time.Second)
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
