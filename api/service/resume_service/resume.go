package resume_service

import (
	"context"
	"encoding/json"
	"hr-api/pkg/cache"
	"time"

	"hr-api/models"
	"hr-api/service/cache_service"
)

type Resume struct {
	Id         int
	JobId      int
	Url        string
	FileName   string
	Size       int
	CreateUid  int
	CreateTime int
	UpdateTime int
	Ctx        context.Context

	Page       int
	Limit      int
	CacheClear int
}

func (r *Resume) Add() error {
	resume := map[string]interface{}{
		"job_id":     r.JobId,
		"url":        r.Url,
		"size":       r.Size,
		"create_uid": r.CreateUid,
	}
	return models.AddResume(resume)
}

func (r *Resume) Edit() error {
	data := make(map[string]interface{})
	data["filename"] = r.FileName
	data["job_id"] = r.JobId
	data["update_time"] = int(time.Now().Unix())

	return models.EditResume(r.Id, data)
}

func (r *Resume) Delete() error {
	return models.DeleteResume(r.Id)
}

func (r *Resume) Count() (int, error) {
	return models.GetResumeTotal(r.FileName, r.getMaps())
}

func (r *Resume) ExistByID() (bool, error) {
	return models.ExistResumeByID(r.Id)
}

func (r *Resume) GetAll() ([]*models.Resume, error) {
	var (
		datas, cacheDatas []*models.Resume
		err               error
	)

	cacheService := cache_service.Cache{
		Name:    cache.CACHE_RESUME,
		Keyword: r.FileName,
		Page:    r.Page,
		Limit:   r.Limit,
	}

	rd, err := cache.GetInstance()
	if err != nil {
		return []*models.Resume{}, nil
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

	datas, err = models.GetResumes(r.Page, r.Limit, r.FileName, r.getMaps())
	if err != nil {
		return nil, err
	}

	if len(datas) > 0 {
		rd.Set(r.Ctx, key, datas, 3600*time.Second)
	}
	return datas, nil
}

func (r *Resume) GetResume() (*models.Resume, error) {
	return models.GetResume(r.Id)
}

func (r *Resume) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})

	if r.Id > 0 {
		maps["id"] = r.Id
	}

	return maps
}
