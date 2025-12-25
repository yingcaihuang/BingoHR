package test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hr-api/routers"
	"hr-api/routers/api"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEntraId(t *testing.T) {
	entraManager, err := api.NewEntraIDClient()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//Opts may include [AccessTypeOnline] or [AccessTypeOffline], as well
	//// as [ApprovalForce].
	s := entraManager.GetAuthURL("AccessTypeOnline")
	fmt.Println(s)

	// https://login.microsoftonline.com/common/oauth2/v2.0/authorize?access_type=offline&client_id=ca24da87-00a3-47e7-a7e3-ecd0b2b21f2a&redirect_uri=https%3A%2F%2Fbingodev.gslb.vip%2Fapi%2Fv2%2Fauth%2Fcallback&response_type=code&scope=User.Read+openid+profile+offline_access&state=AccessTypeOnline
	// 跳转url
	// https://bingodev.gslb.vip/api/v2/auth/callback?code=1.Ab4AOp0nEXatb0q6roWW5q3gJYfaJMqjAOdHp-Ps0LKyHyojAdO-AA.AgABBAIAAABlMNzVhAPUTrARzfQjWPtKAwDs_wUA9P9Fdm9TdHNBcnRpZmFjdHMCAAAAAAAXWO_BorrV0JdrD770cKyUzJsEULPDqr5O5pJLmeu4W7UQVa6Xf7wVdzZQ1JoR6lqOxYHwF02Hgl25SFwz2hbhQOtTjwR11t1FOoRNy9slFtHAvBnd6PXlLGf4ZTDu7rAwgc8Kb-YEHSBc2ZQZ7QA6XvMerVehnb-6kqg4gyPHlwjVwBHOkh-GwWRrloNaSfQzeEptpfHITRoNIwCA3hbRwOPJ8gaqC-QTTLrECdXwNZhPGZgcVu1eWvxEJD_nmLAyIutRzymGaUWWbznaA_QXkEOskfwzzbkd5pIW-nUpLiA_71f_kVIUYsFi5mUTybsivdBFPGIWc4klHQB5busMCo3f7MI4OUZdlgyv5XVKDIAtBl7O0nrexYHG4XzJx6o5UwxXO-ti9aMF3MYULqo3h7uldUIBwX-wVQ4IEUHJjm6m62zuO6WhoXadPxLYaXMCCyUgu3HsssRaz3Q9_o9cEBln68ewfR66q13AGlWhQ3Fi39f7T-XdDCIqMkQ6MIHjwWzUjGDMBw52AK1DiwRGRB1BAyIF4vz5rqPJd23VRQB2q33yJ1kMLudKcxTnMuUMuOQ_ZxGwhFBOT0zcB4u3sOGFA4-T3Y7gjC3axy8QI4A8yFJPr8Bt_0nOCi0gzJyZm-Ciy8EoGH3mE8VEDpLIa2jDhi3H0hE8LGtRSIWUghaBbXnJ6Up3aw8MpdL9_eJ4U2GyH8Xe-Q5X2E4n0qEEVmzAOT46IDROLCYDscaz_Tn8ACahKT-yy7opvA&state=AccessTypeOnline&session_state=00a9b569-fbca-b4b9-a07e-e0df39afc17b#

	// 第三方
	//https://login.microsoftonline.com/11279d3a-ad76-4a6f-baae-8596e6ade025/oauth2/v2.0/authorize?client_id=739fe76e-b150-4dd1-a4e0-70257e9466fd&scope=User.Read%20openid%20profile%20offline_access&redirect_uri=https%3A%2F%2Fscijhm.logto.app%2Fcallback%2Ftmvvqc3cgfqgliii3bbkt&client-request-id=31864e54-cc45-448c-a541-1afb6dec55ca&response_mode=query&response_type=code&x-client-SKU=msal.js.node&x-client-VER=2.6.4&x-client-OS=linux&x-client-CPU=x64&client_info=1&prompt=login&state=teg_SeTORAo
}

func TestLoginRoute(t *testing.T) {
	// 将Gin设置为测试模式
	gin.SetMode(gin.TestMode)

	// 设置路由
	router := routers.InitRouter()

	// 创建一个响应记录器
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v2/login", nil)
	router.ServeHTTP(w, req)

	fmt.Println(w.Body.String())

	// 断言响应状态码为200
	//assert.Equal(t, 200, w.Code)
	// 断言响应体为"pong"
	//assert.Equal(t, "pong", w.Body.String())
}

func TestLoginCallbackRoute(t *testing.T) {
	// 将Gin设置为测试模式
	gin.SetMode(gin.TestMode)

	// 设置路由
	router := routers.InitRouter()

	// 创建一个响应记录器
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v2/auth/callback?code=1.Ab4AOp0nEXatb0q6roWW5q3gJYfaJMqjAOdHp-Ps0LKyHyojAdO-AA.AgABBAIAAABlMNzVhAPUTrARzfQjWPtKAwDs_wUA9P9Fdm9TdHNBcnRpZmFjdHMCAAAAAABsI2T99oMEmhNfcUT44KxaAp0I0jfNx1k3hxjumcDhrlq2gEyyeyWrZNMRMCRoTTxgJKwiWBllEsL4KWvinF5O55NNWCmGxWS08DRQnZyo6D0baLoznYwcKOu71d0rkM3bXws8AvLLU6G-no654MZd06_yB8QzGaCl5q9BlyEoKbLLTJjl8grxX35nm-SoofAtvphAorUnuMI2qTq6r2-tHnxfCLOSoHTzT2cKRPM1j4mmRlyw2ng-x_HcEMl9gFhc-cGBCb3EmVYCLQG5o2QMUmR6oRzRvF92qofv1TTcvFQvXD-fZy8vmptlW1I3KPMG4E1EpAnugwEf4p99z-Ioqi1ipIt9rdPNyQj7IkmVZy9tHANfJ3UgrKv-_Rrv2l_hcEAwlG49Ya1ZcV-ur8ukRNXnyHLdcaeF_yX19qkM9E5BWkMF-4-m-O40kOy5B7iMnCbycD2K_zoQh9GtKef6L0WPSo8gPih7uNgbon0ehCVElorQr8b9Db8GpcBM3DffcaLdU6BicWW79fHH4yqEGwiaBpQlm4IjPU5O8EJ3c5Fuskn1_Ze5RPn_coJ5LiOwGNkfXLI3xR8Te00xg5lsv0IXZS_VxrGcxUoo4X1tRF0nOKPzglepWyIpXGHjaxqCo4tS9ZTpcRAhdsbrUGBmaiAKwCzWm5bDQznpvFRBRailpvYikC40fKTfr6NzORt3jmPxOYhkBdDNXUeSd_0LrEnU7SVkzFEdpCt5L6p7Xm-nrkyFaheZbwMpH-sCQR5USD915EYZRuDBEw&state=AccessTypeOnline&session_state=00a9b569-fbca-b4b9-a07e-e0df39afc17b#", nil)
	router.ServeHTTP(w, req)

	fmt.Println(w.Code)
	fmt.Println(w.Body.String())

	// 断言响应状态码为200
	//assert.Equal(t, 200, w.Code)
	// 断言响应体为"pong"
	//assert.Equal(t, "pong", w.Body.String())
}
