package analyzer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"hr-api/pkg/client"
	"hr-api/pkg/parser"
)

// ResumeAnalyzer 简历分析器主类
type ResumeAnalyzer struct {
	aiClient *client.AzureOpenAIClient
	config   *AnalyzerConfig
}

type AnalyzerConfig struct {
	OutputFormat string
	OutputDir    string
	SaveToFile   bool
}

func NewResumeAnalyzer(analyzerConfig *AnalyzerConfig) (*ResumeAnalyzer, error) {
	// 创建AI客户端
	aiClient, err := client.NewAzureOpenAIClient()
	if err != nil {
		return nil, err
	}

	return &ResumeAnalyzer{
		aiClient: aiClient,
		config:   analyzerConfig,
	}, nil
}

// AnalyzeFile 分析单个简历文件
func (ra *ResumeAnalyzer) AnalyzeFile(ctx context.Context, jobTitle string,
	jobRequirements string,
	jobDescription string, filePath string) (*client.ResumeAnalysis, error) {
	// 1. 解析文件内容
	log.Printf("正在解析文件: %s", filePath)

	parserFactory := parser.NewUnifiedResumeParser()
	resumeText, err := parserFactory.Parse(filePath)
	if err != nil {
		return nil, fmt.Errorf("解析简历文件失败: %v", err)
	}

	log.Printf("成功提取文本，长度: %d 字符", len(resumeText))

	// 2. 使用AI分析内容
	log.Printf("正在使用Azure OpenAI分析简历...")

	analysis, err := ra.aiClient.AnalyzeResume(ctx, jobTitle, jobRequirements, jobDescription, resumeText)
	if err != nil {
		return nil, fmt.Errorf("AI分析失败: %v", err)
	}

	// 3. 生成输出
	if ra.config.SaveToFile {
		if err := ra.saveAnalysis(filePath, analysis); err != nil {
			log.Printf("保存分析结果失败: %v", err)
		}
	}

	return analysis, nil
}

// saveAnalysis 保存分析结果
func (ra *ResumeAnalyzer) saveAnalysis(originalPath string, analysis *client.ResumeAnalysis) error {
	// 创建输出目录
	if err := os.MkdirAll(ra.config.OutputDir, 0755); err != nil {
		return err
	}

	// 生成输出文件名
	baseName := filepath.Base(originalPath)
	nameWithoutExt := baseName[:len(baseName)-len(filepath.Ext(baseName))]
	outputFile := filepath.Join(ra.config.OutputDir, fmt.Sprintf("%s_analysis.json", nameWithoutExt))

	// 序列化为JSON
	jsonData, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		return err
	}

	// 写入文件
	if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
		return err
	}

	log.Printf("分析结果已保存至: %s", outputFile)

	// 同时生成HTML报告
	htmlFile := filepath.Join(ra.config.OutputDir, fmt.Sprintf("%s_report.html", nameWithoutExt))
	if err := ra.generateHTMLReport(htmlFile, analysis); err != nil {
		log.Printf("生成HTML报告失败: %v", err)
	}

	return nil
}

// generateHTMLReport 生成HTML格式的报告
func (ra *ResumeAnalyzer) generateHTMLReport(outputPath string, analysis *client.ResumeAnalysis) error {
	htmlTemplate := `<!DOCTYPE html>
<html>
<head>
    <title>简历分析报告 - %s</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .section { margin-bottom: 30px; }
        .section-title { color: #2c3e50; border-bottom: 2px solid #3498db; padding-bottom: 10px; }
        .info-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 20px; }
        .card { background: #f8f9fa; padding: 20px; border-radius: 8px; }
        .match-score { font-size: 24px; font-weight: bold; color: %s; }
        .strength { color: #27ae60; }
        .weakness { color: #e74c3c; }
    </style>
</head>
<body>
    <h1>简历智能分析报告</h1>
    
    <div class="section">
        <h2 class="section-title">候选人信息</h2>
        <div class="info-grid">
            <div class="card">
                <h3>基本信息</h3>
                <p><strong>姓名:</strong> %s</p>
                <p><strong>邮箱:</strong> %s</p>
                <p><strong>电话:</strong> %s</p>
                <p><strong>地点:</strong> %s</p>
            </div>
            
            <div class="card">
                <h3>匹配度评分</h3>
                <p class="match-score">%d/100</p>
                <p><strong>分析时间:</strong> %s</p>
                <p><strong>预估工作年限:</strong> %d 年</p>
            </div>
        </div>
    </div>
    
    <div class="section">
        <h2 class="section-title">职业摘要</h2>
        <p>%s</p>
    </div>
    
    <div class="section">
        <h2 class="section-title">工作经历</h2>
        %s
    </div>
    
    <div class="section">
        <h2 class="section-title">技能评估</h2>
        <div class="info-grid">
            <div class="card">
                <h3>优势</h3>
                <ul>%s</ul>
            </div>
            <div class="card">
                <h3">改进建议</h3>
                <ul>%s</ul>
            </div>
        </div>
    </div>
</body>
</html>`

	// 根据匹配分数决定颜色
	scoreColor := "#e74c3c" // 红色
	if analysis.Analysis.MatchScore >= 70 {
		scoreColor = "#27ae60" // 绿色
	} else if analysis.Analysis.MatchScore >= 50 {
		scoreColor = "#f39c12" // 橙色
	}

	// 生成工作经历HTML
	workExpHTML := ""
	for _, exp := range analysis.WorkExperience {
		workExpHTML += fmt.Sprintf(`
        <div class="card">
            <h3>%s - %s</h3>
            <p><strong>职位:</strong> %s</p>
            <p><strong>时长:</strong> %s</p>
            <p><strong>主要成就:</strong> %s</p>
        </div>`, exp.Company, exp.Duration, exp.Position, exp.Duration, exp.Achievements)
	}

	// 生成优势列表HTML
	strengthsHTML := ""
	for _, strength := range analysis.Analysis.Strengths {
		strengthsHTML += fmt.Sprintf("<li class=\"strength\">%s</li>", strength)
	}

	// 生成建议列表HTML
	recommendationsHTML := ""
	for _, rec := range analysis.Analysis.Recommendations {
		recommendationsHTML += fmt.Sprintf("<li>%s</li>", rec)
	}

	// 填充模板
	htmlContent := fmt.Sprintf(htmlTemplate,
		analysis.PersonalInfo.Name,
		scoreColor,
		analysis.PersonalInfo.Name,
		analysis.PersonalInfo.Email,
		analysis.PersonalInfo.Phone,
		analysis.PersonalInfo.Location,
		analysis.Analysis.MatchScore,
		time.Now().Format("2006-01-02 15:04:05"),
		analysis.Metadata.EstimatedYOE,
		analysis.Summary,
		workExpHTML,
		strengthsHTML,
		recommendationsHTML,
	)

	return os.WriteFile(outputPath, []byte(htmlContent), 0644)
}
