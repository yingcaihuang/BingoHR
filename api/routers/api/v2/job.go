package v2

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"hr-api/models"
	"hr-api/pkg/app"
	"hr-api/pkg/util"
	"hr-api/service/job_service"
)

// @Summary Get job list
// @Produce json
// @Param keyword query string false
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/job/list [get]
func GetJobs(c *gin.Context) {
	appG := app.Gin{C: c}
	keyword := c.DefaultQuery("keyword", "")

	page := util.GetPage(c)
	limit := util.GetLimit(c)
	cache_clear := util.GetCacheClear(c)

	service := job_service.Job{
		Name:       keyword,
		Page:       page,
		Limit:      limit,
		CacheClear: cache_clear,
		Ctx:        c.Request.Context(),
	}
	datas, err := service.GetAll()
	if err != nil {
		datas = []*models.Job{}
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

type JobAddBody struct {
	Name   string `json:"name" binding:"required,max=64"`
	Demand string `json:"demand" binding:"required"`
	Desc   string `json:"desc"`
}

// @Summary Add a job
// @Produce  json
// @Param name body string true "Name"
// @Param demand body string true "Demand"
// @Param desc body string true "Desc"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/job/create [post]
func AddJob(c *gin.Context) {
	var appG = app.Gin{C: c}

	var bodyData JobAddBody
	if err := c.ShouldBindJSON(&bodyData); err != nil {
		appG.FailResponse(err.Error())
		return
	}

	currentUid := util.GetCurrentUid(c)
	service := job_service.Job{
		Name:      bodyData.Name,
		Demand:    bodyData.Demand,
		Desc:      bodyData.Desc,
		CreateUid: currentUid,
	}

	err := service.Add()
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	appG.SuccessResponse(bodyData)
}

type JobEditBody struct {
	Id     int    `json:"id" binding:"required,min=1"`
	Name   string `json:"name" binding:"required,max=64"`
	Demand string `json:"demand"`
	Desc   string `json:"desc"`
}

// @Summary Edit a job
// @Produce  json
// @Param name body string true "Name"
// @Param demand body string true "Demand"
// @Param desc body string true "Desc"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/job/update [put]
func EditJob(c *gin.Context) {
	var appG = app.Gin{C: c}

	var bodyData JobEditBody
	if err := c.ShouldBindJSON(&bodyData); err != nil {
		appG.FailResponse(err.Error())
		return
	}

	jobService := job_service.Job{
		Id:     bodyData.Id,
		Name:   bodyData.Name,
		Demand: bodyData.Demand,
		Desc:   bodyData.Desc,
	}

	exists, err := jobService.ExistByID()
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	if !exists {
		appG.FailResponse(fmt.Sprintf("招聘需求不存在: %d", bodyData.Id))
		return
	}

	err = jobService.Edit()
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	appG.SuccessResponse(bodyData)
}

type JobURI struct {
	Id int `uri:"id" binding:"required,min=1"`
}

// @Summary Delete a job
// @Produce  json
// @Param id path int true "Id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/job/delete/{id} [delete]
func DeleteJob(c *gin.Context) {
	appG := app.Gin{C: c}

	var uri JobURI
	if err := c.ShouldBindUri(&uri); err != nil {
		appG.FailResponse(err.Error())
		return
	}

	jobService := job_service.Job{Id: uri.Id}
	exists, err := jobService.ExistByID()
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	if !exists {
		appG.FailResponse(fmt.Sprintf("招聘需求不存在: %d", uri.Id))
		return
	}

	if err := jobService.Delete(); err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	appG.SuccessResponse(uri)
}
