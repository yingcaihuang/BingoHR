package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Resume struct {
	ID         int    `json:"id" gorm:"primaryKey"`
	JobId      int    `json:"job_id"`
	FileName   string `json:"filename" gorm:"column:filename"`
	Size       int    `json:"size"`
	CreateUid  int    `json:"create_uid"`
	CreateUser string `json:"create_user"`
	CreateTime int    `json:"create_time"`
	UpdateTime int    `json:"update_time"`
	JobName    string `json:"job_name" gorm:"->"`
}

// GetResumes get resume list data
func GetResumes(page int, limit int, keyword string, maps interface{}) ([]*Resume, error) {
	var (
		datas []*Resume
		err   error
	)

	query := db.Select("jobs.name AS job_name, resumes.*").Where(maps)
	query = query.Joins("LEFT JOIN jobs ON jobs.id = resumes.job_id")

	if keyword != "" {
		query = query.Where("resumes.filename LIKE ?", "%"+keyword+"%")
	}

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	err = query.Find(&datas).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if len(datas) > 0 {
		// 查询创建人用户名
		for index, job := range datas {
			user, err := GetUser(job.CreateUid)
			if err != nil {
				return nil, err
			}
			datas[index].CreateUser = user.Username
		}
	}

	return datas, nil
}

// GetResume Get a resume by id
func GetResume(id int) (*Resume, error) {
	var d Resume
	err := db.Where("id = ? ", id).First(&d).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &d, nil
}

// GetResumeTotal counts the total number of resumes based on the constraint
func GetResumeTotal(keyword string, maps interface{}) (int, error) {
	var count int64

	query := db.Model(&Resume{}).Where(maps)
	if keyword != "" {
		query = query.Where("filename LIKE ?", "%"+keyword+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

// AddResume add a single resume
func AddResume(data map[string]interface{}) error {
	now := int(time.Now().Unix())
	job := Resume{
		JobId:      data["job_id"].(int),
		FileName:   data["filename"].(string),
		Size:       data["size"].(int),
		CreateTime: now,
		CreateUid:  data["create_uid"].(int),
	}
	if err := db.Create(&job).Error; err != nil {
		return err
	}

	return nil
}

// EditResume modify a single resume
func EditResume(id int, data interface{}) error {
	if err := db.Model(&Resume{}).Where("id = ? ", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// DeleteResume delete a single resume
func DeleteResume(id int) error {
	if err := db.Where("id = ?", id).Delete(Resume{}).Error; err != nil {
		return err
	}

	return nil
}

// CleanAllResume clear all job
func CleanAllResume() error {
	if err := db.Unscoped().Delete(&Resume{}).Error; err != nil {
		return err
	}

	return nil
}

// ExistResumeByID determines whether a resume exists based on the ID
func ExistResumeByID(id int) (bool, error) {
	var count int64
	err := db.Model(&Resume{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
