package v2

import (
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"hr-api/pkg/app"
	"hr-api/pkg/cache"
	"time"
)

// 2FA 二次认证

func GetOtpTotp(c *gin.Context) {
	appG := app.Gin{C: c}
	username := c.DefaultQuery("username", "")
	if len(username) == 0 {
		appG.FailResponse("缺少参数username")
		return
	}

	secret, qrurl := generateTOTP(username)
	if len(secret) == 0 {
		appG.FailResponse("secret生成错误")
		return
	}

	redis, err := cache.GetInstance()
	if err != nil {
		appG.FailResponse(err.Error())
		return
	}
	//存储 session 到 Redis
	//fmt.Println(string(sessionJSON))
	err2 := redis.Set(c.Request.Context(), username, secret, 24*time.Hour)
	if err2 != nil {
		appG.FailResponse(err2.Error())
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

func validate(code, secret string) bool {
	b, _ := totp.ValidateCustom(code, secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	return b
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

	redis, err := cache.GetInstance()
	if err != nil {
		appG.FailResponse(err.Error())
		return
	}
	//存储 session 到 Redis
	//fmt.Println(string(sessionJSON))
	var secret string
	err2 := redis.Get(c.Request.Context(), username, &secret)
	if err2 != nil {
		appG.FailResponse(err2.Error())
		return
	}

	ok := validate(code, secret)

	appG.SuccessResponse(ok)
}
