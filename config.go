package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
)

// Config 配置结构
type Config struct {
	DefaultDelay     int               `json:"default_delay"`
	DefaultClickType string            `json:"default_click_type"`
	Hotkeys          map[string]Hotkey `json:"hotkeys"`
	UI               UIConfig          `json:"ui"`
	Logging          LoggingConfig     `json:"logging"`
}

// Hotkey 快捷键配置
type Hotkey struct {
	Start string `json:"start"`
	Stop  string `json:"stop"`
}

// UIConfig UI配置
type UIConfig struct {
	WindowWidth  int    `json:"window_width"`
	WindowHeight int    `json:"window_height"`
	Theme        string `json:"theme"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Enabled bool   `json:"enabled"`
	Level   string `json:"level"`
}

// LoadConfig 加载配置文件
func LoadConfig() (*Config, error) {
	// 尝试读取配置文件
	data, err := os.ReadFile("config.json")
	if err != nil {
		// 如果文件不存在，创建默认配置
		return createDefaultConfig(), nil
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config, nil
}

// SaveConfig 保存配置文件
func SaveConfig(config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	err = os.WriteFile("config.json", data, 0644)
	if err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
}

// createDefaultConfig 创建默认配置
func createDefaultConfig() *Config {
	return &Config{
		DefaultDelay:     1000,
		DefaultClickType: "左键",
		Hotkeys: map[string]Hotkey{
			"macos": {
				Start: "Command+F",
				Stop:  "Command+G",
			},
			"windows": {
				Start: "Alt+F",
				Stop:  "Alt+G",
			},
			"linux": {
				Start: "Alt+F",
				Stop:  "Alt+G",
			},
		},
		UI: UIConfig{
			WindowWidth:  400,
			WindowHeight: 500,
			Theme:        "light",
		},
		Logging: LoggingConfig{
			Enabled: true,
			Level:   "info",
		},
	}
}

// GetCurrentOSHotkey 获取当前操作系统的快捷键
func (c *Config) GetCurrentOSHotkey() Hotkey {
	var osKey string
	switch runtime.GOOS {
	case "darwin":
		osKey = "macos"
	case "windows":
		osKey = "windows"
	default:
		osKey = "linux"
	}

	if hotkey, exists := c.Hotkeys[osKey]; exists {
		return hotkey
	}

	// 返回默认值
	return Hotkey{
		Start: "Alt+F",
		Stop:  "Alt+G",
	}
}
