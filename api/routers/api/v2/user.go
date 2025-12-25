package v2

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"hr-api/models"
	"hr-api/pkg/app"
	"hr-api/pkg/util"
	"hr-api/service/user_service"
)

// @Summary Get user data list
// @Produce  json
// @Param keyword query string false
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/user/list [get]
func GetUsers(c *gin.Context) {
	appG := app.Gin{C: c}
	keyword := c.DefaultQuery("keyword", "")

	page := util.GetPage(c)
	limit := util.GetLimit(c)
	cache_clear := util.GetCacheClear(c)

	service := user_service.User{
		Username:   keyword,
		Page:       page,
		Limit:      limit,
		CacheClear: cache_clear,
		Ctx:        c.Request.Context(),
	}
	datas, err := service.GetAll()
	if err != nil {
		datas = []*models.User{}
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

type UserAddBody struct {
	Username string `json:"username" binding:"required,max=32"`
	Password string `json:"password" binding:"required,max=32"`
	Email    string `json:"email" binding:"required,max=128"`
	Roles    []int  `json:"roles" binding:"required,dive,required,min=1"`
}

// @Summary Add a user
// @Produce  json
// @Param username body string true "Username"
// @Param password body string true "Password"
// @Param email body string true "Email"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/user/create [post]
func AddUser(c *gin.Context) {
	var appG = app.Gin{C: c}

	var data UserAddBody
	if err := c.ShouldBindJSON(&data); err != nil {
		appG.FailResponse(err.Error())
		return
	}

	currentUid := util.GetCurrentUid(c)
	service := user_service.User{
		Username:  data.Username,
		Password:  data.Password,
		Email:     data.Email,
		CreateUid: currentUid,
		Roles:     data.Roles,
	}

	exist, err := service.ExistUserByUsername()
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}
	if exist {
		appG.FailResponse(fmt.Sprintf("用户名 %s 已经存在", data.Username))
		return
	}

	err = service.Add()
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	appG.SuccessResponse(data)
}

type UserEditBody struct {
	Id       int    `json:"id" binding:"required,min=1"`
	Password string `json:"password" binding:"required,max=32"`
	Email    string `json:"email" binding:"required,max=128"`
	Roles    []int  `json:"roles" binding:"required,dive,required,min=1"`
}

// @Summary Edit a role
// @Produce  json
// @Param id body int true "Id"
// @Param password body string false "Password"
// @Param email body string false "Email"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/user/update [put]
func EditUser(c *gin.Context) {
	var appG = app.Gin{C: c}

	var data UserEditBody
	if err := c.ShouldBindJSON(&data); err != nil {
		appG.FailResponse(err.Error())
		return
	}

	service := user_service.User{
		Id:       data.Id,
		Password: data.Password,
		Roles:    data.Roles,
	}

	existsData, err := service.GetUser()
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	if existsData.ID == 0 {
		appG.FailResponse(fmt.Sprintf("用户不存在: %d", data.Id))
		return
	}

	if data.Email != existsData.Email {
		service.Email = data.Email
	} else {
		service.Email = existsData.Email
	}

	err = service.Edit()
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	appG.SuccessResponse(data)
}

type UserURI struct {
	Id int `uri:"id" binding:"required,min=1"`
}

// @Summary Delete a user
// @Produce  json
// @Param id path int true "Id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/tags/{id} [delete]
func DeleteUser(c *gin.Context) {
	appG := app.Gin{C: c}

	var uri UserURI
	if err := c.ShouldBindUri(&uri); err != nil {
		appG.FailResponse(err.Error())
		return
	}

	// 用户ID为1的禁止删除
	if uri.Id == 1 {
		appG.PermDeniedResponse("Permission deined")
		return
	}

	service := user_service.User{Id: uri.Id}
	exists, err := service.ExistByID()
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	if !exists {
		appG.FailResponse(fmt.Sprintf("用户不存在: %d", uri.Id))
		return
	}

	if err := service.Delete(); err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	appG.SuccessResponse(uri)
}

// @Summary Get current user's perms
// @Produce json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/user/perms [get]
func GetUserPerms(c *gin.Context) {
	var appG = app.Gin{C: c}

	currentUid := util.GetCurrentUid(c)
	service := user_service.User{Id: currentUid}
	perms := service.GetUserPerms()

	appG.SuccessResponse(map[string]interface{}{
		"perms": perms,
	})
}
