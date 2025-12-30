package v2

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"hr-api/models"
	"hr-api/pkg/app"
	"hr-api/pkg/util"
	"hr-api/service/resume_analyze_record_service"
	"hr-api/service/resume_service"
)

// @Summary Get a resume's analyze result list
// @Produce  json
// @Param resume_id query string false
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/resume/list [get]
func GetRecords(c *gin.Context) {
	appG := app.Gin{C: c}
	resume_id := util.GetIntQuery("resume_id", c)

	page := util.GetPage(c)
	limit := util.GetLimit(c)
	cache_clear := util.GetCacheClear(c)

	service := resume_analyze_record_service.ResumeAnalyzeRecord{
		ResumeId:   resume_id,
		Page:       page,
		Limit:      limit,
		CacheClear: cache_clear,
		Ctx:        c.Request.Context(),
	}
	datas, err := service.GetAll()
	if err != nil {
		datas = []*models.ResumeAnalyzeRecord{}
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

type AnalyzeRecordAddBody struct {
	ResumeId int `json:"resume_id" binding:"required,min=1"`
}

// @Summary Add a resume analyze
// @Produce json
// @Param resume_id body int true "ResumeId"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/resume/create [post]
func AddResumeAnalyze(c *gin.Context) {
	var appG = app.Gin{C: c}

	var data AnalyzeRecordAddBody
	if err := c.ShouldBindJSON(&data); err != nil {
		appG.FailResponse(err.Error())
		return
	}

	resumeService := resume_service.Resume{Id: data.ResumeId}
	exists, err := resumeService.ExistByID()
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	if !exists {
		appG.FailResponse(fmt.Sprintf("简历不存在: %d", data.ResumeId))
		return
	}

	currentUid := util.GetCurrentUid(c)
	service := resume_analyze_record_service.ResumeAnalyzeRecord{
		ResumeId:  data.ResumeId,
		CreateUid: currentUid,
	}

	err = service.Add()
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	appG.SuccessResponse(data)
}
