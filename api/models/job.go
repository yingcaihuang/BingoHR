package models

import (
	"time"

	"gorm.io/gorm"
)

type Job struct {
	ID         int    `json:"id" gorm:"primaryKey"`
	Name       string `json:"name"`
	Demand     string `json:"demand"`
	Desc       string `json:"desc"`
	CreateUid  int    `json:"create_uid"`
	CreateUser string `json:"create_user" gorm:"-"`
	CreateTime int    `json:"create_time"`
	UpdateTime int    `json:"update_time"`
}

// GetJobs get job list data
func GetJobs(page int, limit int, keyword string, maps interface{}) ([]*Job, error) {
	var (
		jobs []*Job
		err  error
	)

	query := db.Model(&Job{}).Where(maps)

	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	err = query.Find(&jobs).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if len(jobs) > 0 {
		// 查询创建人用户名
		for index, job := range jobs {
			user, err := GetUser(job.CreateUid)
			if err != nil {
				return nil, err
			}
			jobs[index].CreateUser = user.Username
		}
	}

	return jobs, nil
}

// GetJob Get a job by id
func GetJob(id int) (*Job, error) {
	var d Job
	err := db.Model(&Job{}).Where("id = ? ", id).First(&d).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &d, nil
}

// GetJobTotal counts the total number of jobs based on the constraint
func GetJobTotal(keyword string, maps interface{}) (int, error) {
	var count int64

	query := db.Model(&Job{}).Where(maps)
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

// AddJob add a single job
func AddJob(data map[string]interface{}) error {
	now := int(time.Now().Unix())
	job := Job{
		Name:       data["name"].(string),
		Demand:     data["demand"].(string),
		Desc:       data["desc"].(string),
		CreateTime: now,
		CreateUid:  data["create_uid"].(int),
	}
	if err := db.Debug().Create(&job).Error; err != nil {
		return err
	}

	return nil
}

// EditJob modify a single job
func EditJob(id int, data interface{}) error {
	if err := db.Model(&Job{}).Where("id = ? ", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// DeleteJob delete a single job
func DeleteJob(id int) error {
	if err := db.Where("id = ?", id).Delete(Job{}).Error; err != nil {
		return err
	}

	return nil
}

// CleanAllJob clear all job
func CleanAllJob() error {
	if err := db.Unscoped().Delete(&Job{}).Error; err != nil {
		return err
	}

	return nil
}

// ExistJobByID determines whether a job exists based on the ID
func ExistJobByID(id int) (bool, error) {
	var count int64
	err := db.Model(&Job{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
