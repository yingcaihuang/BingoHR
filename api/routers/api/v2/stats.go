package v2

import (
	"github.com/gin-gonic/gin"
	"hr-api/service/resume_analyze_record_service"

	"hr-api/pkg/app"
	"hr-api/pkg/util"
)

func StatusStatsHandler(c *gin.Context) {
	appG := app.Gin{C: c}

	page := util.GetPage(c)
	limit := util.GetLimit(c)
	cache_clear := util.GetCacheClear(c)

	service := resume_analyze_record_service.ResumeAnalyzeRecord{
		Page:       page,
		Limit:      limit,
		CacheClear: cache_clear,
		Ctx:        c.Request.Context(),
	}
	datas, err := service.CountByStatus()
	if err != nil {
		appG.FailResponse(err.Error())
		return
	}

	appG.SuccessResponse(map[string]interface{}{
		"lists": datas,
		"page":  page,
		"limit": limit,
	})
}

func DegreeStatsHandler(c *gin.Context) {
	appG := app.Gin{C: c}

	page := util.GetPage(c)
	limit := util.GetLimit(c)
	cache_clear := util.GetCacheClear(c)

	service := resume_analyze_record_service.ResumeAnalyzeRecord{
		Page:       page,
		Limit:      limit,
		CacheClear: cache_clear,
		Ctx:        c.Request.Context(),
	}
	datas, err := service.CountByDegree()
	if err != nil {
		appG.FailResponse(err.Error())
		return
	}

	appG.SuccessResponse(map[string]interface{}{
		"lists": datas,
		"page":  page,
		"limit": limit,
	})
}

func LocationStatsHandler(c *gin.Context) {
	appG := app.Gin{C: c}

	page := util.GetPage(c)
	limit := util.GetLimit(c)
	cache_clear := util.GetCacheClear(c)

	service := resume_analyze_record_service.ResumeAnalyzeRecord{
		Page:       page,
		Limit:      limit,
		CacheClear: cache_clear,
		Ctx:        c.Request.Context(),
	}
	datas, err := service.CountByLocation()
	if err != nil {
		appG.FailResponse(err.Error())
		return
	}

	appG.SuccessResponse(map[string]interface{}{
		"lists": datas,
		"page":  page,
		"limit": limit,
	})
}

func ScoreRangeStatistics(c *gin.Context) {
	appG := app.Gin{C: c}

	page := util.GetPage(c)
	limit := util.GetLimit(c)
	cache_clear := util.GetCacheClear(c)

	service := resume_analyze_record_service.ResumeAnalyzeRecord{
		Page:       page,
		Limit:      limit,
		CacheClear: cache_clear,
		Ctx:        c.Request.Context(),
	}
	datas, err := service.ScoreRangeStatistics()
	if err != nil {
		appG.FailResponse(err.Error())
		return
	}

	appG.SuccessResponse(map[string]interface{}{
		"lists": datas,
		"page":  page,
		"limit": limit,
	})
}

func TopInstitution(c *gin.Context) {
	appG := app.Gin{C: c}

	page := util.GetPage(c)
	limit := util.GetLimit(c)
	cache_clear := util.GetCacheClear(c)

	service := resume_analyze_record_service.ResumeAnalyzeRecord{
		Page:       page,
		Limit:      limit,
		CacheClear: cache_clear,
		Ctx:        c.Request.Context(),
	}
	datas, err := service.TopInstitution(limit)
	if err != nil {
		appG.FailResponse(err.Error())
		return
	}

	appG.SuccessResponse(map[string]interface{}{
		"lists": datas,
		"page":  page,
		"limit": limit,
	})
}

func TopLocation(c *gin.Context) {
	appG := app.Gin{C: c}

	page := util.GetPage(c)
	limit := util.GetLimit(c)
	cache_clear := util.GetCacheClear(c)

	service := resume_analyze_record_service.ResumeAnalyzeRecord{
		Page:       page,
		Limit:      limit,
		CacheClear: cache_clear,
		Ctx:        c.Request.Context(),
	}
	datas, err := service.TopLocation(limit)
	if err != nil {
		appG.FailResponse(err.Error())
		return
	}

	appG.SuccessResponse(map[string]interface{}{
		"lists": datas,
		"page":  page,
		"limit": limit,
	})
}
