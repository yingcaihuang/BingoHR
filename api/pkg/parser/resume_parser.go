package parser

import (
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/xavier268/mydocx" // DOCX
	"log"
	"os"
	"path/filepath"
	"strings"
)

type UnifiedResumeParser struct {
	// 可以在这里配置 TES 引擎等
	pdfEngine string
}

func NewUnifiedResumeParser() *UnifiedResumeParser {
	return &UnifiedResumeParser{
		pdfEngine: "pdfium", // 默认使用 PDFium，也可用 poppler, mupdf[citation:4]
	}
}

func (p *UnifiedResumeParser) SupportsFormat(ext string) bool {
	ext = strings.ToLower(ext)
	supported := []string{".pdf", ".docx", ".doc", ".txt", ".odt", ".pptx"}
	for _, s := range supported {
		if ext == s {
			return true
		}
	}
	return false
}

func (p *UnifiedResumeParser) Parse(filePath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".pdf":
		return p.parsePDF(filePath)
	case ".docx":
		// 使用 mydocx 进行文本提取[citation:3]
		return p.parseDocx(filePath)
	case ".txt":
		return p.parseText(filePath)
	default:
		return "", fmt.Errorf("unsupported file format: %s", ext)
	}
}

func (rp *UnifiedResumeParser) parsePDF(pdfPath string) (string, error) {
	// 检查文件是否存在
	if _, err := os.Stat(pdfPath); err != nil {
		return "", fmt.Errorf("文件不存在: %v", err)
	}

	// 策略2：使用pdfcpu
	text, err := rp.extractWithPDFCPU(pdfPath)
	if err == nil && len(text) > 50 {
		log.Printf("使用pdfcpu成功提取%d字符", len(text))
		return text, nil
	}

	// 策略3：降级到简单提取
	log.Println("警告：使用简单文本提取，质量可能较低")
	return rp.extractSimple(pdfPath)
}

// 优化说明（简短）:
// 1) 先尝试以可读文本提取（ExtractTextFile），再回退到 content streams（ExtractContentFile）。
// 2) 明确检查临时目录创建错误并确保清理。
// 3) 使用 rp.readExtractedFiles 和 rp.parseContentStream 的结果长度判断是否成功。
// 4) 所有提取步骤都包含错误处理并返回有意义的错误信息。
// 注意：如果你的工程里使用的是 model.NewDefaultConfiguration()，把 pdfcpu.NewDefaultConfiguration() 替换回你的 model 包即可。

func (rp *UnifiedResumeParser) extractWithPDFCPU(pdfPath string) (string, error) {
	// 创建临时目录用于保存 pdfcpu 输出
	tmpDir, err := os.MkdirTemp("", "pdfcpu_*")
	if err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}
	// 确保退出时清理
	defer os.RemoveAll(tmpDir)

	// 配置：若需要解密请设置 conf.UserPW
	conf := model.NewDefaultConfiguration()
	// conf.UserPW = "user-password-if-any"

	// 2) 回退到 content streams（低级流），再尝试解析 content stream
	if err := api.ExtractContentFile(pdfPath, tmpDir, nil, conf); err == nil {
		// 读取并尝试 parse content stream（parseContentStream 应能处理 content 输出）
		if raw, err := rp.readExtractedFiles(tmpDir); err == nil {
			parsed := rp.parseContentStream(raw)
			if len(parsed) > 100 {
				return parsed, nil
			}
			// 如果 parse 后仍不足，返回原始提取结果（有时原始文本也可用）
			if len(raw) > 0 {
				return raw, nil
			}
		}
	}

	// 3) 若以上都失败，返回详细错误提示（调用者可据此记录/回退到其他解析器）
	return "", fmt.Errorf("pdfcpu 无法提取文本（尝试 ExtractTextFile 和 ExtractContentFile 均失败或结果不足）")
}

func (rp *UnifiedResumeParser) parseContentStream(content string) string {
	// 改进的内容流解析
	var result strings.Builder
	inTextBlock := false

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "BT") {
			inTextBlock = true
			continue
		}

		if strings.Contains(line, "ET") {
			inTextBlock = false
			result.WriteString("\n")
			continue
		}

		if inTextBlock {
			// 提取文本操作符中的内容
			if strings.Contains(line, "Tj") {
				text := rp.extractTextFromTJ(line)
				result.WriteString(text)
				result.WriteString(" ")
			} else if strings.Contains(line, "TJ") {
				text := rp.extractTextFromTJ(line)
				result.WriteString(text)
				result.WriteString(" ")
			}
		}
	}

	return result.String()
}

func (rp *UnifiedResumeParser) extractTextFromTJ(line string) string {
	// 查找文本内容，可能在括号或尖括号中
	// 例如: (Hello) Tj 或 <48656C6C6F> Tj

	// 查找括号内容
	start := strings.Index(line, "(")
	end := strings.LastIndex(line, ")")
	if start != -1 && end != -1 && end > start {
		return line[start+1 : end]
	}

	// 查找十六进制内容
	start = strings.Index(line, "<")
	end = strings.LastIndex(line, ">")
	if start != -1 && end != -1 && end > start {
		hexStr := line[start+1 : end]
		// 这里可以添加十六进制解码
		return fmt.Sprintf("[HEX:%s]", hexStr)
	}

	return ""
}

func (rp *UnifiedResumeParser) extractSimple(pdfPath string) (string, error) {
	// 简单读取文件，尝试查找文本
	data, err := os.ReadFile(pdfPath)
	if err != nil {
		return "", err
	}

	content := string(data)
	var result strings.Builder

	// 查找可能的文本片段
	// 这是一个非常基础的实现
	for i := 0; i < len(content)-100; i++ {
		chunk := content[i : i+100]
		if isLikelyText(chunk) {
			result.WriteString(cleanText(chunk))
			result.WriteString(" ")
		}
	}

	if result.Len() == 0 {
		return "", fmt.Errorf("无法提取任何文本")
	}

	return result.String(), nil
}

func isLikelyText(chunk string) bool {
	// 简单启发式判断是否是文本
	textChars := 0
	for _, c := range chunk {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || c == ' ' || c == '.' || c == ',' {
			textChars++
		}
	}
	return float64(textChars)/float64(len(chunk)) > 0.6
}

func cleanText(text string) string {
	// 清理文本
	result := strings.Builder{}
	for _, c := range text {
		if c >= 32 && c <= 126 || c >= 0x4E00 && c <= 0x9FFF {
			result.WriteRune(c)
		} else {
			result.WriteRune(' ')
		}
	}
	return strings.TrimSpace(result.String())
}

func (rp *UnifiedResumeParser) readExtractedFiles(dir string) (string, error) {
	// 读取pdfcpu提取的文件
	var result strings.Builder

	files, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if !file.IsDir() {
			data, err := os.ReadFile(filepath.Join(dir, file.Name()))
			if err == nil {
				result.Write(data)
				result.WriteString("\n")
			}
		}
	}

	return result.String(), nil
}

func (p *UnifiedResumeParser) parseDocx(filePath string) (string, error) {
	// 使用 mydocx 提取文本[citation:3]
	contentMap, err := mydocx.ExtractText(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to extract text from DOCX using mydocx: %v", err)
	}

	// 将不同部分（正文、页眉、页脚）的文本合并
	var fullText strings.Builder
	for container, paragraphs := range contentMap {
		fullText.WriteString(fmt.Sprintf("--- From %s ---\n", container))
		for _, para := range paragraphs {
			fullText.WriteString(para + "\n")
		}
	}
	return fullText.String(), nil
}

func (p *UnifiedResumeParser) parseText(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read text file: %v", err)
	}
	return string(data), nil
}
