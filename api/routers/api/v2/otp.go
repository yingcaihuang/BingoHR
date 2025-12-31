package v2

import (
	"hr-api/models"
	"hr-api/pkg/app"
	"hr-api/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// 2FA 二次认证
func GetOtpTotp(c *gin.Context) {
	appG := app.Gin{C: c}
	username := c.DefaultQuery("username", "")
	if len(username) == 0 {
		appG.FailResponse("缺少参数username")
		return
	}

	// 将生成的secret和qurl存到数据库中
	u, _ := models.GetUserByName(username)
	if len(u.OtpUrl) > 0 {
		appG.SuccessResponse(u.OtpUrl)
		return
	}

	secret, qrurl := generateTOTP(username)
	if len(secret) == 0 {
		appG.FailResponse("secret生成错误")
		return
	}

	data := make(map[string]interface{})
	data["otp_secret"] = secret
	data["otp_url"] = qrurl
	data["update_time"] = int(time.Now().Unix())
	err := models.EditUser(u.ID, data, []int{})
	if err != nil {
		appG.FailResponse(err.Error())
		return
	}

	appG.SuccessResponse(qrurl)
}

func generateTOTP(user string) (string, string) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "BingoHR",
		AccountName: user,
	})
	if err != nil {
		return "", ""
	}
	// Secret (需要保存到数据库)
	secret := key.Secret()
	// QR Code URL（前端渲染即可）
	qrURL := key.URL()
	return secret, qrURL
}

func validate(code, secret string) (bool, error) {
	b, err := totp.ValidateCustom(code, secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	return b, err
}

func ValidateOtpTotp(c *gin.Context) {
	appG := app.Gin{C: c}
	username := c.DefaultQuery("username", "")
	if len(username) == 0 {
		appG.FailResponse("缺少参数username")
		return
	}

	code := c.DefaultQuery("code", "")
	if len(code) == 0 {
		appG.FailResponse("缺少参数code")
		return
	}

	u, _ := models.GetUserByName(username)
	ok, err := validate(code, u.OtpSecret)
	if err != nil {
		appG.FailResponse(err.Error())
		return
	}

	if !ok {
		appG.SuccessResponse(map[string]interface{}{
			"token":     "",
			"validated": ok,
		})
		return
	}

	token, err := util.GenerateToken(u.ID, u.Username, u.Password)
	if err != nil {
		appG.IntervalErrorResponse("Token生成失败, 请稍后再试")
		return
	}

	// 更新2fa标记字段
	data := make(map[string]interface{})
	data["otp_enable"] = 0
	data["update_time"] = int(time.Now().Unix())
	err = models.EditUser(u.ID, data, []int{})
	if err != nil {
		appG.FailResponse(err.Error())
		return
	}

	appG.SuccessResponse(map[string]interface{}{
		"token":     token,
		"validated": ok,
	})
}
