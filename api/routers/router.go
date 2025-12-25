package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "hr-api/docs"

	swaggerFiles "github.com/swaggo/files" // 注意这里的包路径
	ginSwagger "github.com/swaggo/gin-swagger"

	"hr-api/middleware/jwt"
	"hr-api/pkg/export"
	"hr-api/pkg/namedroute"
	"hr-api/pkg/qrcode"
	"hr-api/pkg/upload"
	"hr-api/routers/api"
	v2 "hr-api/routers/api/v2"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(namedroute.Middleware())

	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	// r.POST("/auth", api.GetAuth)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiV2 := namedroute.New(r).Group("/api/v2")

	// 登录
	apiV2.POST("/oauth/login", api.GetAuth)

	// Microsoft Entra ID 相关
	apiV2.GET("/login", v2.Login)
	apiV2.GET("/auth/callback", v2.LoginCallback)

	// 公共上传接口
	apiV2.POST("/upload", api.UploadFile)

	authGroup := apiV2.Group("")
	authGroup.Use(jwt.JWT())
	{
		// 角色相关接口
		authGroup.GET("/role/list", v2.GetRoles).Name("rest.role.list")
		authGroup.POST("/role/create", v2.AddRole).Name("rest.role.create")
		authGroup.PUT("/role/update", v2.EditRole).Name("rest.role.update")
		authGroup.DELETE("/role/delete/:id", v2.DeleteRole).Name("rest.role.delete")
		// 角色权限相关接口
		authGroup.GET("/role/perms", v2.GetRolePerms).Name("rest.role.perms.list")
		authGroup.POST("/role/perms", v2.AddRolePerms).Name("rest.role.perms.create")

		// Microsoft Entra ID 相关
		authGroup.GET("/logout", v2.Logout).Name("rest.entraid.logout")

		// 用户相关接口
		authGroup.GET("/user/list", v2.GetUsers).Name("rest.user.list")
		authGroup.POST("/user/create", v2.AddUser).Name("rest.user.create")
		authGroup.PUT("/user/update", v2.EditUser).Name("rest.user.update")
		authGroup.DELETE("/user/delete/:id", v2.DeleteUser).Name("rest.user.delete")
		authGroup.GET("/user/perms", v2.GetUserPerms).Name("rest.user.perms")

		// 招聘需求相关接口
		authGroup.GET("/job/list", v2.GetJobs).Name("rest.job.list")
		authGroup.POST("/job/create", v2.AddJob).Name("rest.job.create")
		authGroup.PUT("/job/update", v2.EditJob).Name("rest.job.update")
		authGroup.DELETE("/job/delete/:id", v2.DeleteJob).Name("rest.job.delete")

		// 简历相关接口
		authGroup.GET("/resume/list", v2.GetResumes).Name("rest.resume.list")
		authGroup.POST("/resume/create", v2.AddResume).Name("rest.resume.create")
		authGroup.PUT("/resume/update", v2.EditResume).Name("rest.resume.update")
		authGroup.DELETE("/resume/delete/:id", v2.DeleteResume).Name("rest.resume.delete")

		// 返回所有接口地址的名称
		authGroup.GET("/all/perms", func(c *gin.Context) {
			mapedPerms := namedroute.GetRouteNameMap()
			var perms []string
			for _, v := range mapedPerms {
				perms = append(perms, v)
			}
			c.JSON(http.StatusOK, struct {
				Code int      `json:"code"`
				Data []string `json:"data"`
			}{
				Code: 0,
				Data: perms,
			})
		})
	}

	return r
}
