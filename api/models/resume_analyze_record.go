package models

import (
	"github.com/pkg/errors"
	"time"

	"gorm.io/gorm"
)

type ResumeAnalyzeRecord struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	ResumeId    int    `json:"resume_id"`
	FileName    string `json:"filename" gorm:"-"`
	Name        string `json:"name"`
	Location    string `json:"location"`
	Institution string `json:"institution"`
	Degree      string `json:"degree"`
	Status      string `json:"status"`
	Result      string `json:"result"`
	Score       int    `json:"score"`
	CreateUser  string `json:"create_user" gorm:"-"`
	CreateTime  int    `json:"create_time"`
}

// GetRecords get resume analyze record list
func GetRecords(resume_id, page, limit int) ([]*ResumeAnalyzeRecord, error) {
	var (
		records []*ResumeAnalyzeRecord
		err     error
	)

	query := db.Select("resumes.filename as filename, resume_analyze_records.*")
	query = query.Joins("LEFT JOIN resumes ON resumes.id = resume_analyze_records.resume_id")

	if resume_id > 0 {
		query = query.Where("resume_analyze_records.resume_id = ?", resume_id)
	}

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	err = query.Find(&records).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return records, nil
}

// GetRecortd Get a resume analyze record by id
func GetRecortd(id int) (*ResumeAnalyzeRecord, error) {
	var d ResumeAnalyzeRecord
	err := db.Model(&ResumeAnalyzeRecord{}).Where("id = ? ", id).First(&d).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &d, nil
}

// GetRecordTotal counts the total number of jobs based on the constraint
func GetRecordTotal(resume_id int) (int, error) {
	var count int64

	query := db.Model(&ResumeAnalyzeRecord{})
	if resume_id > 0 {
		query = query.Where("resume_analyze_records.resume_id = ?", resume_id)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func InsertRecord(req ResumeAnalyzeRecord) (*ResumeAnalyzeRecord, error) {
	record := ResumeAnalyzeRecord{
		ResumeId:    req.ResumeId,
		Name:        req.Name,
		Location:    req.Location,
		Institution: req.Institution,
		Degree:      req.Degree,
		Result:      req.Result,
		Status:      "success",
		Score:       req.Score,
		//CreateUid:  req.CreateUID,
		CreateTime: int(time.Now().Unix()),
	}
	if err := db.Create(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// AddRecord add a resume analyze record
func AddRecord(data map[string]interface{}) error {
	now := int(time.Now().Unix())
	newData := ResumeAnalyzeRecord{
		ResumeId: data["resume_id"].(int),
		Status:   "pending",
		//CreateUid:  data["create_uid"].(int),
		CreateTime: now,
	}
	if err := db.Create(&newData).Error; err != nil {
		return err
	}

	return nil
}

// EditRecord modify a resume analyze record(for update result after ai analyze)
func EditRecord(id int, result, status string) error {
	var update = make(map[string]string)
	if len(result) > 0 {
		update["result"] = result
	}
	if len(status) > 0 {
		update["status"] = status
	}

	if len(update) == 0 {
		return nil
	}

	if err := db.Model(&ResumeAnalyzeRecord{}).Where("id = ? ", id).Updates(update).Error; err != nil {
		return err
	}

	return nil
}

// ExistRecordByID determines whether a analyze record exists based on the ID
func ExistRecordByID(id int) (bool, error) {
	var count int64
	err := db.Model(&ResumeAnalyzeRecord{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

func CountByStatus() ([]*StatusCount, error) {
	var result []*StatusCount
	err := db.Model(&ResumeAnalyzeRecord{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&result).Error
	return result, err
}

type DegreeCount struct {
	Degree string `json:"degree"`
	Count  int    `json:"count"`
}

func CountByDegree() ([]*DegreeCount, error) {
	var result []*DegreeCount
	err := db.Model(&ResumeAnalyzeRecord{}).
		Where("degree <> ''").
		Select("degree, count(*) as count").
		Group("degree").
		Scan(&result).Error
	return result, err
}

type LocationCount struct {
	Location string `json:"location"`
	Count    int    `json:"count"`
}

func CountByLocation() ([]*LocationCount, error) {
	var result []*LocationCount
	err := db.Model(&ResumeAnalyzeRecord{}).
		Where("location <> ''").
		Select("location, count(*) as count").
		Group("location").
		Scan(&result).Error
	return result, err
}

type ScoreRangeStat struct {
	ScoreRange string `json:"score_range"`
	Count      int64  `json:"count"`
}

func ScoreRangeStatistics() ([]*ScoreRangeStat, error) {
	sql := `
	SELECT '60-70' AS score_range, COUNT(*) as count 
FROM resume_analyze_records
WHERE CAST(score AS UNSIGNED) BETWEEN 60 AND 69

UNION ALL

SELECT '70-80' AS score_range, COUNT(*) as count 
FROM resume_analyze_records
WHERE CAST(score AS UNSIGNED) BETWEEN 70 AND 79

UNION ALL

SELECT '80-90' AS score_range, COUNT(*) as count 
FROM resume_analyze_records
WHERE CAST(score AS UNSIGNED) BETWEEN 80 AND 89

UNION ALL

SELECT '90+' AS score_range, COUNT(*) as count 
FROM resume_analyze_records
WHERE CAST(score AS UNSIGNED) >= 90;

	`
	var res []*ScoreRangeStat
	err := db.Raw(sql).Scan(&res).Error
	return res, err
}

func TopInstitution(limit int) ([]*LocationCount, error) {
	var res []*LocationCount
	err := db.Model(&ResumeAnalyzeRecord{}).
		Where("institution <> ''").
		Select("institution as location, count(*) as count").
		Group("institution").
		Order("count desc").
		Limit(limit).
		Scan(&res).Error
	return res, err
}

func TopLocation(limit int) ([]*LocationCount, error) {
	var res []*LocationCount
	err := db.Model(&ResumeAnalyzeRecord{}).
		Select("location, count(*) as count").
		Group("location").
		Order("count desc").
		Limit(limit).
		Scan(&res).Error
	return res, err
}
