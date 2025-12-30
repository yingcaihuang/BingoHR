package resume_service

import (
	"context"
	"time"

	"hr-api/models"
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

func (r *Resume) Add() (int, error) {
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
		datas []*models.Resume
		err   error
	)

	datas, err = models.GetResumes(r.Page, r.Limit, r.FileName, r.getMaps())
	if err != nil {
		return nil, err
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
