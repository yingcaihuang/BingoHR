package resume_analyze_record_service

import (
	"context"
	"hr-api/models"
)

type ResumeAnalyzeRecord struct {
	Id         int
	ResumeId   int
	Status     string
	Result     string
	CreateUid  int
	CreateTime int
	UpdateTime int
	Ctx        context.Context

	Page       int
	Limit      int
	CacheClear int
}

func (r *ResumeAnalyzeRecord) Add() error {
	data := map[string]interface{}{
		"resume_id":  r.ResumeId,
		"status":     "pending",
		"create_uid": r.CreateUid,
	}
	return models.AddRecord(data)
}

func (r *ResumeAnalyzeRecord) Edit() error {
	return models.EditRecord(r.Id, r.Result, r.Status)
}

func (r *ResumeAnalyzeRecord) Count() (int, error) {
	return models.GetRecordTotal(r.ResumeId)
}

func (r *ResumeAnalyzeRecord) ExistByID() (bool, error) {
	return models.ExistRecordByID(r.Id)
}

func (r *ResumeAnalyzeRecord) GetAll() ([]*models.ResumeAnalyzeRecord, error) {
	var (
		datas []*models.ResumeAnalyzeRecord
		err   error
	)

	datas, err = models.GetRecords(r.ResumeId, r.Page, r.Limit)
	if err != nil {
		return nil, err
	}

	return datas, nil
}

func (r *ResumeAnalyzeRecord) GetResume() (*models.Resume, error) {
	return models.GetResume(r.Id)
}

func (r *ResumeAnalyzeRecord) CountByStatus() ([]*models.StatusCount, error) {
	return models.CountByStatus()
}

func (r *ResumeAnalyzeRecord) CountByDegree() ([]*models.DegreeCount, error) {
	return models.CountByDegree()
}

func (r *ResumeAnalyzeRecord) CountByLocation() ([]*models.LocationCount, error) {
	return models.CountByLocation()
}

func (r *ResumeAnalyzeRecord) ScoreRangeStatistics() ([]*models.ScoreRangeStat, error) {
	return models.ScoreRangeStatistics()
}

func (r *ResumeAnalyzeRecord) TopInstitution(limit int) ([]*models.LocationCount, error) {
	return models.TopInstitution(limit)
}

func (r *ResumeAnalyzeRecord) TopLocation(limit int) ([]*models.LocationCount, error) {
	return models.TopLocation(limit)
}
