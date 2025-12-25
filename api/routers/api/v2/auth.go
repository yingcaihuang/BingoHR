package v2

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"hr-api/models"
	"hr-api/pkg/app"
	"hr-api/pkg/keyvault"
	"hr-api/pkg/util"
	"hr-api/routers/api"
	"hr-api/service/user_service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"hr-api/pkg/cache"
)

func Login(c *gin.Context) {
	appG := app.Gin{C: c}
	// 生成随机 state 防止 CSRF
	state, err := generateRandomString(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}

	//fmt.Println(state)

	// 存储 state 到 session
	//session := sessions.Default(c)
	//session.Set("oauth_state", state)
	//session.Save()

	client, err1 := api.NewEntraIDClient()
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
		return
	}
	// 生成认证 URL
	authURL := client.GetAuthURL(state)

	appG.SuccessResponse(authURL)
	//c.Redirect(http.StatusFound, authURL)
}

// https://bingodev.gslb.vip/api/v2/auth/callback?code=1.Ab4AOp0nEXatb0q6roWW5q3gJYfaJMqjAOdHp-Ps0LKyHyojAdO-AA.AgABBAIAAABlMNzVhAPUTrARzfQjWPtKAwDs_wUA9P9Fdm9TdHNBcnRpZmFjdHMCAAAAAAAXWO_BorrV0JdrD770cKyUzJsEULPDqr5O5pJLmeu4W7UQVa6Xf7wVdzZQ1JoR6lqOxYHwF02Hgl25SFwz2hbhQOtTjwR11t1FOoRNy9slFtHAvBnd6PXlLGf4ZTDu7rAwgc8Kb-YEHSBc2ZQZ7QA6XvMerVehnb-6kqg4gyPHlwjVwBHOkh-GwWRrloNaSfQzeEptpfHITRoNIwCA3hbRwOPJ8gaqC-QTTLrECdXwNZhPGZgcVu1eWvxEJD_nmLAyIutRzymGaUWWbznaA_QXkEOskfwzzbkd5pIW-nUpLiA_71f_kVIUYsFi5mUTybsivdBFPGIWc4klHQB5busMCo3f7MI4OUZdlgyv5XVKDIAtBl7O0nrexYHG4XzJx6o5UwxXO-ti9aMF3MYULqo3h7uldUIBwX-wVQ4IEUHJjm6m62zuO6WhoXadPxLYaXMCCyUgu3HsssRaz3Q9_o9cEBln68ewfR66q13AGlWhQ3Fi39f7T-XdDCIqMkQ6MIHjwWzUjGDMBw52AK1DiwRGRB1BAyIF4vz5rqPJd23VRQB2q33yJ1kMLudKcxTnMuUMuOQ_ZxGwhFBOT0zcB4u3sOGFA4-T3Y7gjC3axy8QI4A8yFJPr8Bt_0nOCi0gzJyZm-Ciy8EoGH3mE8VEDpLIa2jDhi3H0hE8LGtRSIWUghaBbXnJ6Up3aw8MpdL9_eJ4U2GyH8Xe-Q5X2E4n0qEEVmzAOT46IDROLCYDscaz_Tn8ACahKT-yy7opvA&state=AccessTypeOnline&session_state=00a9b569-fbca-b4b9-a07e-e0df39afc17b#
func LoginCallback(c *gin.Context) {
	// 验证 state
	state := c.Query("state")
	fmt.Println(state)
	//session := sessions.Default(c)
	//storedState := session.Get("oauth_state")

	//if state != storedState {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
	//	return
	//}

	// 清除 state
	//session.Delete("oauth_state")

	// 获取授权码
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not found"})
		return
	}

	client, err1 := api.NewEntraIDClient()
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
		return
	}

	// 交换 token
	token, err := client.ExchangeCode(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token", "details": err.Error()})
		return
	}

	// 解析 ID token 获取用户信息
	idToken := token.Extra("id_token").(string)
	//fmt.Println("idToken ", idToken)
	claims, err := client.ValidateToken(c, idToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate token", "details": err.Error()})
		return
	}

	//fmt.Println(claims)

	// 获取用户信息
	u, err := client.GetUserInfo(c, token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	//s, _ := json.Marshal(u)
	//fmt.Println("GetUserInfo ", string(s))

	//创建 session
	sessionID := uuid.New().String()
	sessionData := SessionData{
		UserID:       claims.OID,
		Email:        claims.Email,
		Name:         claims.Name,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
		Groups:       claims.Groups,
		Roles:        claims.Roles,
	}

	//初始化缓存管理器
	redis, err := cache.GetInstance()
	//存储 session 到 Redis
	sessionJSON, _ := json.Marshal(sessionData)
	//fmt.Println(string(sessionJSON))
	err2 := redis.Set(c.Request.Context(), sessionID, string(sessionJSON), 24*time.Hour)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
		return
	}

	username := ""
	if uname, ok := u["name"]; ok {
		username = uname.(string)
	}

	if len(username) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "username is empty"})
		return
	}

	// 保存进数据库
	passwd := "nicaiba_88"
	service := user_service.User{
		Password: passwd,
		Email:    "",
		Roles:    []int{3},
		Username: username,
	}

	if ok, err := service.ExistUserByUsername(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if !ok {
		err = service.Add()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	user, err := models.GetUserByName(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	xtoken, err := util.GenerateToken(user.ID, username, passwd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//设置 session cookie
	c.SetCookie("session_id", sessionID, 86400, "/", "", false, true)
	c.SetCookie("x-token", xtoken, 86400, "/", "", false, true)

	keyVault, err := keyvault.GetKeyVaultConf()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 重定向到前端或返回 token
	c.Redirect(http.StatusFound, keyVault.FrontendURL+"dashboard")
}

func Logout(c *gin.Context) {
	// 清除 session
	sessionID, err := c.Cookie("session_id")
	redis := cache.RedisCache{}
	if err == nil {
		redis.Delete(c, sessionID)
	}

	// 清除 cookie
	c.SetCookie("session_id", "", -1, "/", "", false, true)

	keyVault, err := keyvault.GetKeyVaultConf()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 构建 Entra ID 登出 URL
	logoutURL := "https://login.microsoftonline.com/" + keyVault.TenantID + "/oauth2/v2.0/logout"
	logoutURL += "?post_logout_redirect_uri=" + keyVault.FrontendURL

	c.Redirect(http.StatusFound, logoutURL)
}

func generateRandomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

type SessionData struct {
	UserID       string    `json:"user_id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	Groups       []string  `json:"groups"`
	Roles        []string  `json:"roles"`
}
