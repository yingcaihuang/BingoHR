package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"hr-api/pkg/keyvault"
	"hr-api/pkg/setting"
	"log"
	"strings"
	"time"
)

const (
	defaultTimeout = 60 * time.Second
	maxRetry       = 3
)

// ================= Client =================

type AzureOpenAIClient struct {
	config *setting.MicrosoftEntraIDConfig
	client *openai.Client
}

// ================= Models =================

type ResumeAnalysis struct {
	PersonalInfo   PersonalInfo     `json:"personal_info"`
	Summary        string           `json:"summary"`
	WorkExperience []WorkExperience `json:"work_experience"`
	Education      []Education      `json:"education"`
	Skills         Skills           `json:"skills"`
	Analysis       JobAnalysis      `json:"analysis"`
	Metadata       AnalysisMetadata `json:"metadata"`
}

type PersonalInfo struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Location string   `json:"location"`
	Links    []string `json:"links"`
}

type WorkExperience struct {
	Company          string   `json:"company"`
	Position         string   `json:"position"`
	Duration         string   `json:"duration"`
	Responsibilities []string `json:"responsibilities"`
	Achievements     []string `json:"achievements"`
}

type Education struct {
	Institution    string `json:"institution"`
	Degree         string `json:"degree"`
	Field          string `json:"field"`
	GraduationYear string `json:"graduation_year"`
}

type Skills struct {
	Technical      []string    `json:"technical"`
	Soft           []string    `json:"soft"`
	Languages      interface{} `json:"languages"`
	Certifications []string    `json:"certifications"`
}

type JobAnalysis struct {
	Strengths       []string `json:"strengths"`
	Weaknesses      []string `json:"weaknesses"`
	Recommendations []string `json:"recommendations"`
	MatchScore      int      `json:"match_score"`
}

type AnalysisMetadata struct {
	AnalysisDate string `json:"analysis_date"`
	WordCount    int    `json:"word_count"`
	EstimatedYOE int    `json:"estimated_yoe"`
}

// ================= Config =================

func GetAzureOpenAIConf() (*setting.MicrosoftEntraIDConfig, error) {
	loader, err := keyvault.NewConfigLoader()
	if err != nil {
		return nil, err
	}
	return loader.LoadConfig()
}

// ================= Constructor =================

func NewAzureOpenAIClient() (*AzureOpenAIClient, error) {
	conf, err := GetAzureOpenAIConf()
	if err != nil {
		return nil, err
	}

	// endpoint keyvault这么配置 https://myjycloud.openai.azure.com/openai/deployments/gpt-5.2-chat
	// https://myjycloud.openai.azure.com/openai/deployments/gpt-5.2-chat/chat/completions#
	endpoint := strings.TrimSpace(conf.OpenapiApiEndpoint)
	if endpoint == "" {
		return nil, errors.New("Azure OpenAI endpoint 不能为空")
	}

	client := openai.NewClient(
		option.WithBaseURL(endpoint),
		option.WithAPIKey(conf.OpenapiApiKey),
		option.WithQuery("api-version", conf.OpenapiApiVersion),
	)

	return &AzureOpenAIClient{
		config: conf,
		client: &client,
	}, nil
}

// ================= Public API =================

func (c *AzureOpenAIClient) AnalyzeResume(
	ctx context.Context,
	jobTitle string,
	jobRequirements string,
	jobDescription string,
	resumeText string,
) (*ResumeAnalysis, error) {

	if strings.TrimSpace(resumeText) == "" {
		return nil, errors.New("简历内容不能为空")
	}

	ctx = normalizeContext(ctx)

	var lastErr error
	for attempt := 1; attempt <= maxRetry; attempt++ {
		result, err := c.callOnce(ctx, jobTitle, jobRequirements, jobDescription, resumeText)
		if err == nil {
			return result, nil
		}

		lastErr = err
		log.Printf("[AnalyzeResume] attempt=%d failed: %v", attempt, err)
		time.Sleep(time.Duration(attempt*attempt) * time.Second)
	}

	return nil, fmt.Errorf("AI分析失败（重试%d次）：%w", maxRetry, lastErr)
}

func normalizeContext(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, _ = context.WithTimeout(ctx, defaultTimeout)
	return ctx
}

// ================= Internal =================

func (c *AzureOpenAIClient) callOnce(
	ctx context.Context,
	jobTitle string,
	jobRequirements string,
	jobDescription string,
	resumeText string,
) (*ResumeAnalysis, error) {

	systemPrompt := `你是一个专业的简历分析师。
只允许输出 JSON，不允许任何多余内容。`

	userPrompt := fmt.Sprintf("【招聘需求】\n岗位名称：%s\n岗位要求：%s\n岗位描述：%s\n", jobTitle, jobRequirements, jobDescription)

	userPrompt += fmt.Sprintf(`请分析以下简历内容：

%s

严格返回如下 JSON 结构：
{
  "personal_info": { "name": "", "email": "", "phone": "", "location": "", "links": [] },
  "summary": "",
  "work_experience": [],
  "education": [],
  "skills": { "technical": [], "soft": [], "languages": [], "certifications": [] },
  "analysis": { "strengths": [], "weaknesses": [], "recommendations": [], "match_score": 0 },
  "metadata": { "analysis_date": "", "word_count": 0, "estimated_yoe": 0 }
}`, resumeText)

	var lastErr error

	for attempt := 0; attempt < 2; attempt++ {

		resp, err := c.client.Chat.Completions.New(
			ctx,
			openai.ChatCompletionNewParams{
				Model: c.config.OpenapiApiDeploymentName,
				Messages: []openai.ChatCompletionMessageParamUnion{
					{
						OfSystem: &openai.ChatCompletionSystemMessageParam{
							Content: openai.ChatCompletionSystemMessageParamContentUnion{
								OfString: openai.String(systemPrompt),
							},
						},
					},
					{
						OfUser: &openai.ChatCompletionUserMessageParam{
							Content: openai.ChatCompletionUserMessageParamContentUnion{
								OfString: openai.String(userPrompt),
							},
						},
					},
				},
				//Temperature:         openai.Float(0.7),
				MaxCompletionTokens: openai.Int(50000),
			},
		)
		if err != nil {
			lastErr = err
			continue
		}

		if len(resp.Choices) == 0 {
			lastErr = errors.New("空 choices")
			continue
		}

		raw := resp.Choices[0].Message.Content
		var result ResumeAnalysis

		if err := tryParseJSON(raw, &result); err != nil {
			lastErr = err
			log.Println("⚠️ JSON 解析失败，重试一次")
			log.Println(raw)
			continue
		}

		return &result, nil
	}

	return nil, fmt.Errorf("模型多次返回非法 JSON: %w", lastErr)
}

func tryParseJSON[T any](raw string, out *T) error {
	raw = strings.TrimSpace(raw)
	// 去掉常见包裹
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)
	return json.Unmarshal([]byte(raw), out)
}
