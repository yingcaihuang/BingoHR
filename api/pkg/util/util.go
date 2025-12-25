package util

import (
	"hr-api/pkg/setting"
	"net/url"
	"path"
)

// Setup Initialize the util
func Setup() {
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
}

// Contains 字符串切片包含某个字符串
func Contains(slice []string, search string) bool {
	for _, s := range slice {
		if s == search {
			return true
		}
	}
	return false
}

// GetFilenameFromURL 从URL地址中提取出文件名
func GetFilenameFromURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	if u.Path == "" || u.Path == "/" {
		return ""
	}

	filename := path.Base(u.Path)
	if filename == "." || filename == "/" {
		return ""
	}

	return filename
}
