package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

// Logger 日志记录器
type Logger struct {
	enabled bool
	level   string
	file    *os.File
	logger  *log.Logger
}

// NewLogger 创建新的日志记录器
func NewLogger(enabled bool, level string) *Logger {
	logger := &Logger{
		enabled: enabled,
		level:   level,
	}

	if enabled {
		// 创建日志文件
		logFile, err := os.OpenFile("clicker.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			// Windows 系统下不显示控制台日志
			if runtime.GOOS != "windows" {
				fmt.Printf("创建日志文件失败: %v\n", err)
			}
			return logger
		}

		logger.file = logFile
		logger.logger = log.New(logFile, "", log.LstdFlags)
	}

	return logger
}

// LogInfo 记录信息日志
func (l *Logger) LogInfo(format string, args ...interface{}) {
	if l.enabled && l.logger != nil {
		l.logger.Printf("[INFO] "+format, args...)
	}
}

// LogWarning 记录警告日志
func (l *Logger) LogWarning(format string, args ...interface{}) {
	if l.enabled && l.logger != nil {
		l.logger.Printf("[WARN] "+format, args...)
	}
}

// LogError 记录错误日志
func (l *Logger) LogError(format string, args ...interface{}) {
	if l.enabled && l.logger != nil {
		l.logger.Printf("[ERROR] "+format, args...)
	}
}

// LogClick 记录点击活动
func (l *Logger) LogClick(clickType string, delay int, timestamp time.Time) {
	if l.enabled && l.logger != nil {
		l.logger.Printf("[CLICK] 类型: %s, 延迟: %dms, 时间: %s",
			clickType, delay, timestamp.Format("2006-01-02 15:04:05.000"))
	}
}

// LogStart 记录开始点击
func (l *Logger) LogStart(delay int, clickType string) {
	if l.enabled && l.logger != nil {
		l.logger.Printf("[START] 开始连续点击 - 延迟: %dms, 类型: %s", delay, clickType)
	}
}

// LogStop 记录停止点击
func (l *Logger) LogStop(duration time.Duration, totalClicks int) {
	if l.enabled && l.logger != nil {
		l.logger.Printf("[STOP] 停止连续点击 - 运行时长: %v, 总点击次数: %d", duration, totalClicks)
	}
}

// Close 关闭日志记录器
func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}

// GetLogStats 获取日志统计信息
func (l *Logger) GetLogStats() map[string]interface{} {
	stats := make(map[string]interface{})

	if !l.enabled || l.file == nil {
		stats["enabled"] = false
		return stats
	}

	// 获取文件信息
	fileInfo, err := l.file.Stat()
	if err == nil {
		stats["file_size"] = fileInfo.Size()
		stats["last_modified"] = fileInfo.ModTime()
	}

	stats["enabled"] = true
	stats["level"] = l.level

	return stats
}
