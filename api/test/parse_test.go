package test

import (
	"encoding/json"
	"fmt"
	"hr-api/pkg/analyzer"
	"hr-api/pkg/client"
	"hr-api/pkg/parser"
	"log"
	"testing"
)

func TestParsePDF(t *testing.T) {
	m := parser.NewUnifiedResumeParser()
	parse, err := m.Parse("/Users/captain/develop/verycloud/microsoft/hr-jianli/长沙/戴先生_36岁_智联简历_02158-深信服.docx")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(parse)
}

func TestParseDocx(t *testing.T) {
	analyzerConfig := &analyzer.AnalyzerConfig{
		OutputFormat: "html",
		OutputDir:    "/Users/captain/develop/verycloud/microsoft/hr-jianli/parse",
		SaveToFile:   true,
	}

	// 创建分析器
	resumeAnalyzer, err := analyzer.NewResumeAnalyzer(analyzerConfig)
	if err != nil {
		log.Fatalf("创建简历分析器失败: %v", err)
	}

	// 分析单个文件
	jobTitle := "c++软件工程师"
	jobRequirements := "工作职责:\n1、 负责开发操作医疗设备的软件；\n2、 负责二维或者三维图像的渲染，以及相关的交互；\n3、 根据功能要求完成相关的算法；\n4、 配合设备输入的图像进行功能开发；\n5、 根据公司技术文档规范编写相应的技术文档。"
	jobDescription := "任职资格:\n1、 熟悉C++以及基本的数据结构；\n2、 熟悉基本的设计模式，并且能够运用；\n3、 数学基础较好的优先；\n4、 熟悉嵌入式Linux操作系统，有医疗产品研发经验者优先；\n5、 对医疗行业了解，有HIS，PACS系统开发的优先；\n6、 有较强的责任心，良好团队协作能力，沟通能力，谦虚踏实。"

	analysis, err := resumeAnalyzer.AnalyzeFile(nil, jobTitle, jobRequirements, jobDescription, "/Users/captain/develop/verycloud/microsoft/hr-jianli/呼和浩特/杨先生_34岁_智联简历_00052-金万维.docx")
	if err != nil {
		log.Fatalf("分析失败: %v", err)
	}

	// 打印结果摘要
	printAnalysisSummary(analysis)

}

func printAnalysisSummary(analysis *client.ResumeAnalysis) {
	fmt.Printf("\n=== 简历分析报告 ===\n\n")

	s, _ := json.Marshal(analysis)
	fmt.Println(string(s))

	fmt.Printf("👤 候选人: %s\n", analysis.PersonalInfo.Name)
	fmt.Printf("📧 联系方式: %s | %s\n", analysis.PersonalInfo.Email, analysis.PersonalInfo.Phone)
	fmt.Printf("📍 地点: %s\n\n", analysis.PersonalInfo.Location)

	fmt.Printf("📊 匹配度评分: %d/100\n\n", analysis.Analysis.MatchScore)

	fmt.Printf("📝 职业摘要:\n%s\n\n", analysis.Summary)

	fmt.Printf("💼 工作经历 (%d 个):\n", len(analysis.WorkExperience))
	for i, exp := range analysis.WorkExperience {
		fmt.Printf("  %d. %s - %s (%s)\n", i+1, exp.Company, exp.Position, exp.Duration)
	}

	fmt.Printf("\n🎓 教育背景 (%d 个):\n", len(analysis.Education))
	for i, edu := range analysis.Education {
		fmt.Printf("  %d. %s - %s (%s)\n", i+1, edu.Institution, edu.Degree, edu.GraduationYear)
	}

	fmt.Printf("\n✅ 优势:\n")
	for _, strength := range analysis.Analysis.Strengths {
		fmt.Printf("  • %s\n", strength)
	}

	fmt.Printf("\n💡 改进建议:\n")
	for _, rec := range analysis.Analysis.Recommendations {
		fmt.Printf("  • %s\n", rec)
	}
}
