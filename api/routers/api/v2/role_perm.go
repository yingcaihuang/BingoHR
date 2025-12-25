package v2

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"hr-api/models"
	"hr-api/pkg/app"
	"hr-api/pkg/util"
	"hr-api/service/role_perm_service"
	"hr-api/service/role_service"
)

// @Summary Get multiple perms based on role id
// @Produce  json
// @Param role_id query string false
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/role/perms [get]
func GetRolePerms(c *gin.Context) {
	appG := app.Gin{C: c}
	rid := c.DefaultQuery("role_id", "0")
	role_id := com.StrTo(rid).MustInt()

	roleService := role_service.Role{
		Id: role_id,
	}

	exists, err := roleService.ExistByID()
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	if !exists {
		appG.FailResponse(fmt.Sprintf("角色不存在: %d", role_id))
		return
	}

	cache_clear := util.GetCacheClear(c)

	service := role_perm_service.RolePerm{
		RoleId:     role_id,
		CacheClear: cache_clear,
		Ctx:        c.Request.Context(),
	}

	rolePerms, err := service.GetAll()
	if err != nil {
		rolePerms = []*models.RolePerm{}
	}

	count, err := service.Count()
	if err != nil {
		count = 0
	}

	appG.SuccessResponse(map[string]interface{}{
		"lists": rolePerms,
		"total": count,
	})
}

type RoleAddPerms struct {
	RoleId int      `json:"role_id" binding:"required,min=1"`
	Perms  []string `json:"perms" binding:"required,dive,required,max=128"`
}

// @Summary Role add perms
// @Produce  json
// @Param role_id body string true "RoleId"
// @Param perms body []string true "Perms"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/role/perms [post]
func AddRolePerms(c *gin.Context) {
	var appG = app.Gin{C: c}

	var data RoleAddPerms
	if err := c.ShouldBindJSON(&data); err != nil {
		appG.FailResponse(err.Error())
		return
	}

	service := role_service.Role{
		Id: data.RoleId,
	}

	exists, err := service.ExistByID()
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	if !exists {
		appG.FailResponse(fmt.Sprintf("角色不存在: %d", data.RoleId))
		return
	}

	currentUid := util.GetCurrentUid(c)

	srv := role_perm_service.RolePerm{
		RoleId:    data.RoleId,
		Perms:     data.Perms,
		CreateUid: currentUid,
	}

	err = srv.Add()
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	appG.SuccessResponse(data)
}
