package sever

import (
	"encoding/json"
	"fmt"
	"github.com/DisposaBoy/JsonConfigReader"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// removeComments 使用正则表达式去除 JSON 字符串中的注释
func removeComments(jsonStr string) (string, error) {
	// 将字符串分割成行
	lines := strings.Split(jsonStr, "\n")
	// 用于存储清理后的行
	var cleanedLines []string

	for _, line := range lines {
		// 如果是单行注释，且不是 URL 的一部分，则移除
		if strings.HasPrefix(strings.TrimSpace(line), "//") && !strings.Contains(line, "://") {
			continue
		}
		// 检查多行注释的开头和结尾
		if strings.Contains(line, "/*") {
			// 移除多行注释的开始部分直至结束部分
			startIndex := strings.Index(line, "/*")
			endIndex := strings.Index(line, "*/")
			if endIndex != -1 && endIndex+2 < len(line) {
				// 注释在行内结束
				line = line[:startIndex] + line[endIndex+2:]
			} else {
				// 注释跨行，移除开始部分，并从下一行开始继续检查
				line = line[:startIndex]
				for _, line = range lines {
					endIndex := strings.Index(line, "*/")
					if endIndex != -1 {
						// 注释结束，继续处理剩余部分
						line = line[endIndex+2:]
						break
					}
				}
			}
		}
		// 将清理后的行添加到结果中
		cleanedLines = append(cleanedLines, line)
	}

	// 将清理后的行重新组合成字符串
	return strings.Join(cleanedLines, "\n"), nil
}

func SetSetting(pUrl, enterprise string) error {
	var settingsPath string

	// 获取当前用户的主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取主目录出错: %v\n", err)
	}

	// 根据操作系统确定 settings.json 文件的路径
	switch runtime.GOOS {
	case "windows":
		settingsPath = filepath.Join(homeDir, "AppData", "Roaming", "Code", "User", "settings.json")
	case "darwin":
		settingsPath = filepath.Join(homeDir, "Library", "Application Support", "Code", "User", "settings.json")
	default:
		return fmt.Errorf("不支持的操作系统: %s\n", runtime.GOOS)
	}

	// 打开 settings.json 文件
	file, err := os.Open(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，创建一个具有默认设置的新文件
			defaultSettings := map[string]interface{}{
				"": "",
			}
			content, err := json.MarshalIndent(defaultSettings, "", "    ")
			if err != nil {
				return fmt.Errorf("创建默认 settings.json 文件出错: %v\n", err)
			}

			err = os.WriteFile(settingsPath, content, 0644)
			if err != nil {
				return fmt.Errorf("写入默认 settings.json 文件出错: %v\n", err)
			}
			// 打开新创建的 settings.json 文件
			file, err = os.Open(settingsPath)
			if err != nil {
				return fmt.Errorf("打开新创建的 settings.json 文件出错: %v\n", err)
			}
		} else {
			return fmt.Errorf("读取 settings.json 文件出错: %v\n", err)
		}
	}
	defer file.Close()

	// 使用 JsonConfigReader 创建一个新的Reader
	reader := JsonConfigReader.New(file)

	// 解析JSON文件内容到一个map中
	var settings map[string]interface{}
	err = json.NewDecoder(reader).Decode(&settings)
	if err != nil {
		return fmt.Errorf("解析 JSON 出错: %v\n", err)
	}

	// 检查是否存在特定的键，并删除它们
	delete(settings, "github.copilot.advanced")
	delete(settings, "github-enterprise.uri")

	// 添加新的内容
	settings["github.copilot.advanced"] = map[string]interface{}{
		"authProvider":               "github-enterprise",
		"debug.overrideProxyUrl":     pUrl,
		"debug.chatOverrideProxyUrl": pUrl + "/chat",
		"debug.overrideChatEngine":   "gpt-4",
	}
	settings["github-enterprise.uri"] = enterprise

	// 将修改后的map转换回JSON格式
	updatedContent, err := json.MarshalIndent(settings, "", "    ")
	if err != nil {
		return fmt.Errorf("编组 JSON 出错: %v\n", err)
	}

	// 将更新后的内容写回文件
	err = os.WriteFile(settingsPath, updatedContent, 0644)
	if err != nil {
		return fmt.Errorf("写入更新后的 settings.json 文件出错: %v\n", err)
	}

	return nil
}

func SetJbHost(token, pUrl, name string) error {
	var settingsPath string

	// 获取当前用户的主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("Error getting home directory: %v\n", err)
	}

	// 根据操作系统确定 hosts.json 文件的路径
	switch runtime.GOOS {
	case "windows":
		settingsPath = filepath.Join(homeDir, "AppData", "Local", "github-copilot", "hosts.json")
	case "darwin":
		settingsPath = filepath.Join(homeDir, ".config", "github-copilot", "hosts.json")
	default:

		return fmt.Errorf("Unsupported operating system: %s\n", runtime.GOOS)
	}
	settingsDir := filepath.Dir(settingsPath)
	if err := os.MkdirAll(settingsDir, 0755); err != nil {

		return fmt.Errorf("Error creating directory: %v\n", err)
	}
	// 新的hosts.json内容
	newSettings := map[string]interface{}{
		"github.com": map[string]interface{}{
			"user":        name,
			"oauth_token": token,
			"dev_override": map[string]interface{}{
				"copilot_token_url": pUrl + "/copilot_internal/v2/token",
			},
		},
	}

	// 将新的设置转换为JSON格式
	newContent, err := json.MarshalIndent(newSettings, "", "    ")
	if err != nil {

		return fmt.Errorf("Error marshaling new settings to JSON: %v\n", err)
	}

	// 将新的JSON内容写入hosts.json文件
	err = os.WriteFile(settingsPath, newContent, 0644)
	if err != nil {

		return fmt.Errorf("Error writing new settings to hosts.json file: %v\n", err)
	}

	fmt.Println("hosts.json updated successfully.")

	return nil
}
