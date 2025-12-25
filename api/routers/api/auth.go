package api

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hr-api/pkg/keyvault"
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"hr-api/pkg/app"
	"hr-api/pkg/e"
	"hr-api/pkg/util"
	"hr-api/service/auth_service"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

type Auth struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary Get Auth
// @Produce  json
// @Param username query string true "userName"
// @Param password query string true "password"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth [get]
func GetAuth(c *gin.Context) {
	appG := app.Gin{C: c}

	var auth Auth
	if err := c.ShouldBindJSON(&auth); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	authService := auth_service.Auth{Username: auth.Username, Password: util.EncodeMD5(auth.Password)}
	uid, err := authService.Check()
	if err != nil {
		appG.IntervalErrorResponse("系统内部错误, 请稍后再试")
		return
	}

	if uid == 0 {
		appG.FailResponse("用户名密码错误")
		return
	}

	token, err := util.GenerateToken(uid, auth.Username, auth.Password)
	if err != nil {
		appG.IntervalErrorResponse("Token生成失败, 请稍后再试")
		return
	}

	appG.SuccessResponse(map[string]string{
		"token": token,
	})
}

// internal/auth/entra.go

type EntraIDClient struct {
	config     *oauth2.Config
	oidcConfig *OIDCConfig
	httpClient *http.Client
}

type OIDCConfig struct {
	Issuer        string   `json:"issuer"`
	AuthURL       string   `json:"authorization_endpoint"`
	TokenURL      string   `json:"token_endpoint"`
	JWKSURL       string   `json:"jwks_uri"`
	UserInfoURL   string   `json:"userinfo_endpoint"`
	EndSessionURL string   `json:"end_session_endpoint"`
	Scopes        []string `json:"scopes_supported"`
}

type TokenClaims struct {
	jwt.RegisteredClaims
	Email            string                 `json:"email"`
	Name             string                 `json:"name"`
	Roles            []string               `json:"roles"`
	Groups           []string               `json:"groups"`
	Scopes           string                 `json:"scp"`
	UPN              string                 `json:"upn"`
	OID              string                 `json:"oid"`
	TID              string                 `json:"tid"`
	AppID            string                 `json:"appid"`
	IdentityProvider string                 `json:"idp"`
	AdditionalClaims map[string]interface{} `json:"-"`
}

func NewEntraIDClient() (*EntraIDClient, error) {
	//issuerURL := strings.Replace(setting.AppSetting.IssuerURL, "{tenantid}", setting.AppSetting.TenantID, 1)
	keyVault, err := keyvault.GetKeyVaultConf()
	s, _ := json.Marshal(keyVault)
	fmt.Println("keyvault ", string(s))
	if err != nil {
		return nil, fmt.Errorf("failed to create key vault client: %v", err)
	}

	oauthConfig := &oauth2.Config{
		ClientID:     keyVault.ClientID,
		ClientSecret: keyVault.ClientSecret,
		RedirectURL:  keyVault.RedirectURL,
		//Scopes:       keyVault.Scopes,
		Endpoint: microsoft.AzureADEndpoint(keyVault.TenantID),
	}

	// 获取 OIDC 配置
	oidcConfig, err := fetchOIDCConfig(keyVault.TenantID)
	if err != nil {
		return nil, err
	}

	return &EntraIDClient{
		config:     oauthConfig,
		oidcConfig: oidcConfig,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func fetchOIDCConfig(tenantID string) (*OIDCConfig, error) {
	url := fmt.Sprintf(
		"https://login.microsoftonline.com/%s/v2.0/.well-known/openid-configuration",
		tenantID,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var config OIDCConfig
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *EntraIDClient) GetAuthURL(state string) string {
	return c.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("scope", "User.Read openid profile offline_access email"))
}

func (c *EntraIDClient) ExchangeCode(ctx *gin.Context, code string) (*oauth2.Token, error) {
	return c.config.Exchange(ctx, code)
}

func (c *EntraIDClient) ValidateToken(ctx *gin.Context, tokenString string) (*TokenClaims, error) {
	// 解析 token 但不验证签名
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &TokenClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	keyVault, err := keyvault.GetKeyVaultConf()
	if err != nil {
		return nil, fmt.Errorf("failed to create key vault client: %v", err)
	}

	s, _ := json.Marshal(keyVault)
	fmt.Println("GetKeyVaultConf ", string(s))
	fmt.Println("tenant_id ", keyVault.TenantID)

	// 验证 issuer
	//expectedIssuer := fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", keyVault.TenantID)

	//fmt.Println("claims issuer ", claims.Issuer)
	//fmt.Println("expectedIssuer issuer ", expectedIssuer)
	//
	//if claims.Issuer != expectedIssuer {
	//	return nil, fmt.Errorf("invalid issuer: got %s, expected %s", claims.Issuer, expectedIssuer)
	//}

	// 验证 audience
	if !contains(claims.Audience, keyVault.ClientID) {
		return nil, errors.New("invalid audience")
	}

	// 验证 token 未过期
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	// 获取 JWKS 并验证签名
	key, err := c.getSigningKey(ctx, token)
	if err != nil {
		return nil, err
	}

	// 使用正确的 key 验证签名
	validatedToken, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := validatedToken.Claims.(*TokenClaims); ok && validatedToken.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (c *EntraIDClient) getSigningKey(ctx *gin.Context, token *jwt.Token) (interface{}, error) {
	// 从 token header 获取 kid
	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("kid not found in token header")
	}

	// 获取 JWKS
	resp, err := c.httpClient.Get(c.oidcConfig.JWKSURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jwks struct {
		Keys []struct {
			Kid string   `json:"kid"`
			Kty string   `json:"kty"`
			Alg string   `json:"alg"`
			Use string   `json:"use"`
			N   string   `json:"n"`
			E   string   `json:"e"`
			X5c []string `json:"x5c"`
		} `json:"keys"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	// 查找匹配的 key
	for _, key := range jwks.Keys {
		if key.Kid == kid {
			switch key.Kty {
			case "RSA":
				// 使用 x5c 证书链中的第一个证书
				if len(key.X5c) > 0 {
					certPEM := fmt.Sprintf("-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----",
						key.X5c[0])
					return jwt.ParseRSAPublicKeyFromPEM([]byte(certPEM))
				}

				// 或者从 n/e 构造 RSA 公钥
				nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
				if err != nil {
					return nil, err
				}

				eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
				if err != nil {
					return nil, err
				}

				n := new(big.Int).SetBytes(nBytes)
				e := new(big.Int).SetBytes(eBytes).Int64()

				return &rsa.PublicKey{
					N: n,
					E: int(e),
				}, nil
			}
		}
	}

	return nil, errors.New("signing key not found")
}

func (c *EntraIDClient) GetUserInfo(ctx *gin.Context, accessToken string) (map[string]interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.oidcConfig.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (c *EntraIDClient) RefreshToken(ctx *gin.Context, refreshToken string) (*oauth2.Token, error) {
	token := &oauth2.Token{
		RefreshToken: refreshToken,
		Expiry:       time.Now().Add(-time.Hour), // 设置为过期状态
	}

	return c.config.TokenSource(ctx, token).Token()
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
