package api

import (
	"fmt"
	"hr-api/pkg/keyvault"
	"hr-api/pkg/setting"
	"io"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gin-gonic/gin"

	"hr-api/pkg/app"
)

func GetBlobConf() (*setting.MicrosoftEntraIDConfig, error) {
	loader, err := keyvault.NewConfigLoader()
	if err != nil {
		return nil, err
	}
	return loader.LoadConfig()
}

// @Summary Upload file to microsoft azure blob storage service
// @Produce  json
// @Param file formData file true "File"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags/import [post]
func UploadFile(c *gin.Context) {
	appG := app.Gin{C: c}

	conf, err := GetBlobConf()
	accountName := conf.BlobAccountName     // BLOB-ACCOUNT-NAME
	containerName := conf.BlobContainerName // BLOB-CONTAINER-NAME
	// BLOB-ACCESS-KEY
	accoutKey := conf.BlobAccessKey
	// 创建 Shared Key 凭据
	cred, err := azblob.NewSharedKeyCredential(accountName, accoutKey)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to create blob client: %v", err)
		appG.IntervalErrorResponse(errMsg)
		return
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to create blob client: %v", err)
		appG.IntervalErrorResponse(errMsg)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		appG.IntervalErrorResponse(err.Error())
		return
	}

	fp, err := file.Open()
	if err != nil {
		appG.IntervalErrorResponse(fmt.Sprintf("Failed to open file: %v", err))
		return
	}
	defer fp.Close()

	fileBytes, err := io.ReadAll(fp)
	if err != nil {
		appG.IntervalErrorResponse(fmt.Sprintf("Failed to read file: %v", err))
		return
	}

	ctx := c.Request.Context()

	_, err = client.CreateContainer(ctx, containerName, nil)
	if err != nil {
		log.Printf("CreateContainer returned (likely exists): %v", err)
	}

	blobName := file.Filename
	_, err = client.UploadBuffer(ctx, containerName, blobName, fileBytes, nil)
	if err != nil {
		errMsg := fmt.Sprintf("Upload failed: %v", err)
		appG.IntervalErrorResponse(errMsg)
		return
	}

	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", accountName, containerName, blobName)

	appG.SuccessResponse(map[string]interface{}{
		"container": containerName,
		"filename":  file.Filename,
		"size":      file.Size,
		"url":       blobURL,
	})
}
