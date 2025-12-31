package v2

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"

	"hr-api/models"
	"hr-api/pkg/app"
	"hr-api/pkg/bus"
	"hr-api/pkg/keyvault"
	"hr-api/pkg/util"
	"hr-api/service/job_service"
	"hr-api/service/resume_service"
)

// @Summary Get resume data list
// @Produce  json
// @Param keyword query string false
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/resume/list [get]
func GetResumes(c *gin.Context) {
	appG := app.Gin{C: c}
	keyword := c.DefaultQuery("keyword", "")

	page := util.GetPage(c)
	limit := util.GetLimit(c)
	cache_clear := util.GetCacheClear(c)

	service := resume_service.Resume{
		FileName:   keyword,
		Page:       page,
		Limit:      limit,
		CacheClear: cache_clear,
		Ctx:        c.Request.Context(),
	}
	datas, err := service.GetAll()
	if err != nil {
		datas = []*models.Resume{}
	}

	count, err := service.Count()
	if err != nil {
		count = 0
	}

	appG.SuccessResponse(map[string]interface{}{
		"lists": datas,
		"total": count,
		"page":  page,
		"limit": limit,
	})
}

type ResumeAddBody struct {
	JobId int    `json:"job_id" binding:"required,min=1"`
	Url   string `json:"url" binding:"required,max=1024"`
	Size  int    `json:"size" binding:"required,min=1"`
}

// @Summary Add a resume
// @Produce json
// @Param job_id body int true "JobId"
// @Param url body string true "Url"
// @Param size body int true "Size"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/resume/create [post]
func AddResume(c *gin.Context) {
	var appG = app.Gin{C: c}

	var data ResumeAddBody
	if err := c.ShouldBindJSON(&data); err != nil {
		appG.FailResponse(err.Error())
		return
	}

	jobService := job_service.Job{Id: data.JobId}
	job, err := jobService.GetJob()
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	if job.ID == 0 {
		appG.FailResponse(fmt.Sprintf("招聘需求不存在: %d", data.JobId))
		return
	}

	currentUid := util.GetCurrentUid(c)
	service := resume_service.Resume{
		JobId:     data.JobId,
		Url:       data.Url,
		Size:      data.Size,
		CreateUid: currentUid,
	}

	newId, err := service.Add()
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	// 将简历和职位信息推送给队列Bus
	resumeBus := &models.ResumeBus{
		ResumeId:  newId,
		FileName:  util.GetFilenameFromURL(data.Url),
		JobName:   job.Name,
		JobDemand: job.Demand,
		JobDesc:   job.Desc,
		Url:       data.Url,
	}
	err = uploadToBusQueue(c, resumeBus)
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	appG.SuccessResponse(data)
}

type ResumeEditBody struct {
	Id       int    `json:"id" binding:"required,min=1"`
	JobId    int    `json:"job_id"`
	FileName string `json:"filename"`
}

// @Summary Edit a resume
// @Produce json
// @Param id body int true "Id"
// @Param job_id body string false "JobId"
// @Param filename body string false "FileName"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/resume/update [put]
func EditResume(c *gin.Context) {
	var appG = app.Gin{C: c}

	var data ResumeEditBody
	if err := c.ShouldBindJSON(&data); err != nil {
		appG.FailResponse(err.Error())
		return
	}

	service := resume_service.Resume{Id: data.Id}

	existData, err := service.GetResume()
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	if existData.ID == 0 {
		appG.FailResponse(fmt.Sprintf("简历记录不存在: %d", data.Id))
		return
	}

	var jobChanged bool = false
	var job *models.Job

	if data.JobId > 0 && existData.JobId != data.JobId {
		jobService := job_service.Job{Id: data.JobId}
		job, err = jobService.GetJob()
		if err != nil {
			appG.IntervalErrorResponse(err.Error())
			return
		}

		if job.ID == 0 {
			appG.FailResponse(fmt.Sprintf("招聘需求不存在: %d", data.JobId))
			return
		}
		jobChanged = true
		service.JobId = data.JobId
	} else {
		data.JobId = existData.JobId
		service.JobId = existData.JobId
	}

	if len(data.FileName) > 0 && existData.FileName != data.FileName {
		service.FileName = data.FileName
	} else {
		data.FileName = existData.FileName
		service.FileName = existData.FileName
	}

	err = service.Edit()
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	if jobChanged {
		// 将简历和职位信息推送给队列Bus
		resumeBus := &models.ResumeBus{
			ResumeId: existData.ID,
			// 始终使用URL中解析出来的文件名
			FileName:  util.GetFilenameFromURL(existData.Url),
			JobName:   job.Name,
			JobDemand: job.Demand,
			JobDesc:   job.Desc,
			Url:       existData.Url,
		}
		err = uploadToBusQueue(c, resumeBus)
		if err != nil {
			appG.IntervalErrorResponse(err.Error())
			return
		}
	}

	appG.SuccessResponse(data)
}

type ResumeURI struct {
	Id int `uri:"id" binding:"required,min=1"`
}

// @Summary Delete a resume
// @Produce json
// @Param id path int true "Id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/tags/{id} [delete]
func DeleteResume(c *gin.Context) {
	appG := app.Gin{C: c}

	var uri ResumeURI
	if err := c.ShouldBindUri(&uri); err != nil {
		appG.FailResponse(err.Error())
		return
	}

	service := resume_service.Resume{Id: uri.Id}
	exists, err := service.ExistByID()
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	if !exists {
		appG.FailResponse(fmt.Sprintf("简历记录不存在: %d", uri.Id))
		return
	}

	if err := service.Delete(); err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	appG.SuccessResponse(uri)
}

func uploadToBusQueue(c *gin.Context, r *models.ResumeBus) error {
	client, err := bus.NewBusClient(keyvault.ServiceBusNamespace)
	if err != nil {
		return err
	}
	sender, err := client.NewQueueSender(keyvault.ServiceBusQueueName)
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(r)
	if err != nil {
		return err
	}
	err = sender.Send(c.Request.Context(), jsonData)
	if err != nil {
		return err
	}

	return nil
}
