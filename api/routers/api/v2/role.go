package v2

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"hr-api/models"
	"hr-api/pkg/app"
	"hr-api/pkg/util"
	"hr-api/service/role_service"
)

// @Summary Get multiple roles
// @Produce  json
// @Param keyword query string false
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/role/list [get]
func GetRoles(c *gin.Context) {
	appG := app.Gin{C: c}
	keyword := c.DefaultQuery("keyword", "")

	page := util.GetPage(c)
	limit := util.GetLimit(c)
	cache_clear := util.GetCacheClear(c)

	service := role_service.Role{
		Name:       keyword,
		Page:       page,
		Limit:      limit,
		CacheClear: cache_clear,
		Ctx:        c.Request.Context(),
	}
	datas, err := service.GetAll()
	if err != nil {
		datas = []*models.Role{}
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

type RoleAddBody struct {
	Name string `json:"name" binding:"required,max=32"`
}

// @Summary Add a role
// @Produce  json
// @Param name body string true "Name"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/role/create [post]
func AddRole(c *gin.Context) {
	var appG = app.Gin{C: c}

	var roleData RoleAddBody
	if err := c.ShouldBindJSON(&roleData); err != nil {
		appG.FailResponse(err.Error())
		return
	}

	currentUid := util.GetCurrentUid(c)
	service := role_service.Role{
		Name:      roleData.Name,
		CreateUid: currentUid,
	}

	err := service.Add()
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	appG.SuccessResponse(roleData)
}

type RoleEditBody struct {
	Id   int    `json:"id" binding:"required,min=1"`
	Name string `json:"name" binding:"required,max=32"`
}

// @Summary Edit a role
// @Produce  json
// @Param name body string true "Name"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/role/update [put]
func EditRole(c *gin.Context) {
	var appG = app.Gin{C: c}

	var roleData RoleEditBody
	if err := c.ShouldBindJSON(&roleData); err != nil {
		appG.FailResponse(err.Error())
		return
	}

	roleService := role_service.Role{
		Id:   roleData.Id,
		Name: roleData.Name,
	}

	exists, err := roleService.ExistByID()
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	if !exists {
		appG.FailResponse(fmt.Sprintf("角色不存在: %d", roleData.Id))
		return
	}

	err = roleService.Edit()
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	appG.SuccessResponse(roleData)
}

type RoleURI struct {
	Id int `uri:"id" binding:"required,min=1"`
}

// @Summary Delete a role
// @Produce  json
// @Param id path int true "Id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/tags/{id} [delete]
func DeleteRole(c *gin.Context) {
	appG := app.Gin{C: c}

	var uri RoleURI
	if err := c.ShouldBindUri(&uri); err != nil {
		appG.FailResponse(err.Error())
		return
	}

	roleService := role_service.Role{Id: uri.Id}
	exists, err := roleService.ExistByID()
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	if !exists {
		appG.FailResponse(fmt.Sprintf("角色不存在: %d", uri.Id))
		return
	}

	if err := roleService.Delete(); err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	appG.SuccessResponse(uri)
}
