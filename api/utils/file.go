package utils

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"

	"github.com/google/go-tika/tika"
)

func ReadFileContent(filePath string, tikaHost string) (string, error) {
	// for remote file, download it first
	if strings.HasPrefix(filePath, "http") {
		file, err := downloadFile(filePath)
		if err != nil {
			return "", err
		}
		filePath = file
	}
	// 创建 Tika 客户端
	client := tika.NewClient(nil, tikaHost)
	// 打开 PDF 文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error with open file: %v", err)
	}
	defer file.Close()

	// 使用 Tika 提取 PDF 文件的文本内容
	content, err := client.Parse(context.TODO(), file)
	if err != nil {
		return "", fmt.Errorf("error with parse file: %v", err)
	}

	ext := filepath.Ext(filePath)
	switch ext {
	case ".doc", ".docx", ".pdf", ".pptx", "ppt":
		return cleanBlankLine(cleanHtml(content, false)), nil
	case ".xls", ".xlsx":
		return cleanBlankLine(cleanHtml(content, true)), nil
	default:
		return cleanBlankLine(content), nil
	}

}

// 清理文本内容
func cleanHtml(html string, keepTable bool) string {
	// 清理 HTML 标签
	var policy *bluemonday.Policy
	if keepTable {
		policy = bluemonday.NewPolicy()
		policy.AllowElements("table", "thead", "tbody", "tfoot", "tr", "td", "th")
	} else {
		policy = bluemonday.StrictPolicy()
	}
	return policy.Sanitize(html)
}

func cleanBlankLine(content string) string {
	lines := strings.Split(content, "\n")
	texts := make([]string, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) < 2 {
			continue
		}
		// discard image
		if strings.HasSuffix(line, ".png") ||
			strings.HasSuffix(line, ".jpg") ||
			strings.HasSuffix(line, ".jpeg") {
			continue
		}
		texts = append(texts, line)
	}

	return strings.Join(texts, "\n")
}

// 下载文件
func downloadFile(url string) (string, error) {
	base := filepath.Base(url)
	dir := os.TempDir()
	filename := filepath.Join(dir, base)
	out, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// 获取数据
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 写入数据到文件
	_, err = io.Copy(out, resp.Body)
	return filename, err
}

// 判断是否是图片
// 图片类型：JPG、JPG2、PNG、GIF、WEBP、HEIC、HEIF、BMP、PCD、TIFF
func IsImage(url string) bool {
	return strings.HasSuffix(url, ".png") ||
		strings.HasSuffix(url, ".jpg") ||
		strings.HasSuffix(url, ".jpeg") ||
		strings.HasSuffix(url, ".gif") ||
		strings.HasSuffix(url, ".webp") ||
		strings.HasSuffix(url, ".heic") ||
		strings.HasSuffix(url, ".heif") ||
		strings.HasSuffix(url, ".bmp") ||
		strings.HasSuffix(url, ".pcd") ||
		strings.HasSuffix(url, ".tiff")
}

// 判断是否是音频
// 音频后缀 WAV、OGG_OPUS
func IsAudio(url string) bool {
	return strings.HasSuffix(url, ".mp3") ||
		strings.HasSuffix(url, ".wav") ||
		strings.HasSuffix(url, ".ogg")
}

// getExtensionFromContentType 根据Content-Type获取文件后缀
// 返回带点的后缀，例如 ".png"
func getExtensionFromContentType(contentType string) string {
	// 移除可能的参数，例如 "image/png; charset=utf-8" -> "image/png"
	if idx := strings.Index(contentType, ";"); idx != -1 {
		contentType = contentType[:idx]
	}
	contentType = strings.TrimSpace(strings.ToLower(contentType))

	// 常见的 MIME 类型到文件后缀的映射
	mimeToExt := map[string]string{
		// 图片类型
		"image/png":     ".png",
		"image/jpeg":    ".jpg",
		"image/jpg":     ".jpg",
		"image/gif":     ".gif",
		"image/webp":    ".webp",
		"image/bmp":     ".bmp",
		"image/tiff":    ".tiff",
		"image/x-icon":  ".ico",
		"image/svg+xml": ".svg",
		"image/heic":    ".heic",
		"image/heif":    ".heif",
		// 音频类型
		"audio/mpeg": ".mp3",
		"audio/wav":  ".wav",
		"audio/ogg":  ".ogg",
		"audio/mp4":  ".m4a",
		// 视频类型
		"video/mp4":       ".mp4",
		"video/mpeg":      ".mpeg",
		"video/quicktime": ".mov",
		"video/x-msvideo": ".avi",
		// 文档类型
		"application/pdf":    ".pdf",
		"application/msword": ".doc",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
		"application/vnd.ms-excel": ".xls",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         ".xlsx",
		"application/vnd.ms-powerpoint":                                             ".ppt",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": ".pptx",
		// 文本类型
		"text/plain": ".txt",
		"text/html":  ".html",
		"text/css":   ".css",
		"text/csv":   ".csv",
	}

	// 优先使用映射表
	if ext, ok := mimeToExt[contentType]; ok {
		return ext
	}

	// 尝试使用 Go 标准库的 mime 包
	exts, err := mime.ExtensionsByType(contentType)
	if err == nil && len(exts) > 0 {
		return exts[0]
	}

	return ""
}

// fetchContentType 通过HTTP HEAD请求获取URL的Content-Type
func fetchContentType(urlStr string) string {
	// 创建带超时的HTTP客户端
	client := &http.Client{
		Timeout: 3 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 允许最多10次重定向
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			return nil
		},
	}

	// 发送 HEAD 请求
	resp, err := client.Head(urlStr)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	// 获取 Content-Type
	return resp.Header.Get("Content-Type")
}

// ExtractFileSuffixFromURL 从URL中提取文件后缀
// 处理带查询参数的URL，例如：https://example.com/file.pptx?auth_key=xxx 返回 .pptx
// 对于短链接（URL路径中没有后缀），会尝试发送HTTP HEAD请求获取Content-Type来推断后缀
func ExtractFileSuffixFromURL(urlStr string) string {
	// 使用标准库解析URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		// 解析失败时回退到简单字符串处理
		if idx := strings.Index(urlStr, "?"); idx != -1 {
			urlStr = urlStr[:idx]
		}
		if idx := strings.Index(urlStr, "#"); idx != -1 {
			urlStr = urlStr[:idx]
		}
		ext := filepath.Ext(urlStr)
		return strings.ToLower(ext)
	}

	// 从URL路径中提取文件后缀
	ext := filepath.Ext(parsedURL.Path)
	ext = strings.ToLower(ext)

	// 如果有后缀，直接返回
	if ext != "" {
		return ext
	}

	// 如果没有后缀，尝试通过HTTP HEAD请求获取Content-Type
	contentType := fetchContentType(urlStr)
	if contentType != "" {
		ext = getExtensionFromContentType(contentType)
	}

	return ext
}

// ExtractURLsFromText 从文本中提取所有URL
// 使用正则表达式匹配 http:// 或 https:// 开头的URL
func ExtractURLsFromText(text string) []string {
	urls := make([]string, 0)
	// 简单的URL匹配逻辑：查找 http:// 或 https:// 开头的字符串
	words := strings.Fields(text)
	for _, word := range words {
		if strings.HasPrefix(word, "http://") || strings.HasPrefix(word, "https://") {
			// 移除可能的尾部标点符号
			word = strings.TrimRight(word, ".,;:!?\"')")
			urls = append(urls, word)
		}
	}
	return urls
}
