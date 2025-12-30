package bus

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"hr-api/models"
	"hr-api/pkg/analyzer"
	"hr-api/pkg/keyvault"
	"log"
	"os"
	"time"
)

func Run(ctx context.Context) error {
	client, err := NewBusClient(keyvault.ServiceBusNamespace)
	if err != nil {
		return fmt.Errorf("create bus client err: %s", err.Error())
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		receiver, err := client.NewQueueReceiver(keyvault.ServiceBusQueueName)
		if err != nil {
			return fmt.Errorf("NewQueueReceiver err: %s", err.Error())
		}
		err = receiver.ReceiveAndComplete(ctx, func(b []byte) error {
			//log.Println("收到了消息", string(b))

			bus := &models.ResumeBus{}
			if err := json.Unmarshal(b, bus); err != nil {
				return err
			}
			log.Println("收到了消息", bus.Url)

			// 下载文件
			localFile := fmt.Sprintf("/tmp/%s", bus.FileName)
			if err := DownloadFile(bus.FileName, localFile); err != nil {
				return err
			}

			// 分析文件
			analyzerConfig := &analyzer.AnalyzerConfig{
				OutputFormat: "html",
				OutputDir:    "/tmp",
				SaveToFile:   true,
			}

			// 创建分析器
			resumeAnalyzer, err := analyzer.NewResumeAnalyzer(analyzerConfig)
			if err != nil {
				return fmt.Errorf("创建简历分析器失败: %v", err)
			}
			analysis, err := resumeAnalyzer.AnalyzeFile(nil, bus.JobName, bus.JobDemand, bus.JobDesc, localFile)
			if err != nil {
				log.Fatalf("分析失败: %v", err)
			}

			// 分析完入库
			s, _ := json.Marshal(analysis)
			_, err = models.InsertRecord(models.ResumeAnalyzeRecord{
				ResumeId:    bus.ResumeId,
				Name:        analysis.PersonalInfo.Name,
				Location:    analysis.PersonalInfo.Location,
				Institution: analysis.Education[0].Institution,
				Degree:      analysis.Education[0].Degree,
				Result:      string(s),
				Score:       analysis.Analysis.MatchScore,
			})
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("ReceiveAndComplete err: %s", err.Error())
		}

		// 每30秒拉取一次
		time.Sleep(time.Second * 30)
	}
}

func DownloadFile(blobName, localPath string) error {
	conf, err := keyvault.GetBlobConf()
	accountName := conf.BlobAccountName     // BLOB-ACCOUNT-NAME
	containerName := conf.BlobContainerName // BLOB-CONTAINER-NAME
	// BLOB-ACCESS-KEY
	accoutKey := conf.BlobAccessKey
	// 创建 Shared Key 凭据
	cred, err := azblob.NewSharedKeyCredential(accountName, accoutKey)
	if err != nil {
		return fmt.Errorf("failed to create blob client: %v", err)
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create blob client: %v", err)
	}

	// 创建本地文件
	file, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer file.Close()
	// 使用 DownloadFile
	_, err = client.DownloadFile(
		context.Background(),
		containerName,
		blobName,
		file,
		nil,
	)
	return err
}
