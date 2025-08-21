package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"time"

	hook "github.com/robotn/gohook"
)

// 按键映射常量
const (
	// 修饰键 rawcode
	CMD_LEFT  = 54
	CMD_RIGHT = 55
	SHIFT     = 56
	CAPS_LOCK = 57
	OPTION    = 58
	CTRL      = 59 // macOS/Linux 的 Control 键
	CTRL_ALT1 = 60
	CTRL_ALT2 = 61

	// Windows 系统的 Control 键
	CTRL_WIN = 162

	// Windows 系统的 Alt 键
	ALT_WIN = 164

	// 主键 rawcode
	KEY_A      = 0
	KEY_S      = 1
	KEY_D      = 2
	KEY_F      = 3
	KEY_H      = 4
	KEY_G      = 5
	KEY_Z      = 6
	KEY_X      = 7
	KEY_C      = 8
	KEY_V      = 9
	KEY_B      = 11
	KEY_Q      = 12
	KEY_W      = 13
	KEY_E      = 14
	KEY_R      = 15
	KEY_Y      = 16
	KEY_T      = 17
	KEY_RETURN = 36
	KEY_TAB    = 48
	KEY_SPACE  = 49
	KEY_DELETE = 51
	KEY_ESCAPE = 53
	KEY_UP     = 126
	KEY_DOWN   = 125
	KEY_LEFT   = 123
	KEY_RIGHT  = 124
)

// 根据操作系统返回正确的按键名称
func getModifierKeyName(keyType string) string {
	switch keyType {
	case "alt":
		if runtime.GOOS == "darwin" {
			return "Option"
		}
		return "Alt"
	default:
		return keyType
	}
}

// 获取按键名称
func getKeyName(rawcode uint16) string {
	// 根据操作系统使用不同的按键映射
	if runtime.GOOS == "windows" {
		return getWindowsKeyName(rawcode)
	}

	// macOS 和 Linux 使用原有的映射
	switch rawcode {
	case KEY_A:
		return "A"
	case KEY_S:
		return "S"
	case KEY_D:
		return "D"
	case KEY_F:
		return "F"
	case KEY_H:
		return "H"
	case KEY_G:
		return "G"
	case KEY_Z:
		return "Z"
	case KEY_X:
		return "X"
	case KEY_C:
		return "C"
	case KEY_V:
		return "V"
	case KEY_B:
		return "B"
	case KEY_Q:
		return "Q"
	case KEY_W:
		return "W"
	case KEY_E:
		return "E"
	case KEY_R:
		return "R"
	case KEY_Y:
		return "Y"
	case KEY_T:
		return "T"
	case KEY_RETURN:
		return "Return"
	case KEY_TAB:
		return "Tab"
	case KEY_SPACE:
		return "Space"
	case KEY_DELETE:
		return "Delete"
	case KEY_ESCAPE:
		return "Escape"
	case KEY_UP:
		return "Up"
	case KEY_DOWN:
		return "Down"
	case KEY_LEFT:
		return "Left"
	case KEY_RIGHT:
		return "Right"
	default:
		return fmt.Sprintf("Key%d", rawcode)
	}
}

// 获取 Windows 系统的按键名称
func getWindowsKeyName(rawcode uint16) string {
	switch rawcode {
	case 65: // A
		return "A"
	case 83: // S
		return "S"
	case 68: // D
		return "D"
	case 70: // F
		return "F"
	case 72: // H
		return "H"
	case 71: // G
		return "G"
	case 90: // Z
		return "Z"
	case 88: // X
		return "X"
	case 67: // C
		return "C"
	case 86: // V
		return "V"
	case 66: // B
		return "B"
	case 81: // Q
		return "Q"
	case 87: // W
		return "W"
	case 69: // E
		return "E"
	case 82: // R
		return "R"
	case 89: // Y
		return "Y"
	case 84: // T
		return "T"
	case 13: // Enter
		return "Enter"
	case 9: // Tab
		return "Tab"
	case 32: // Space
		return "Space"
	case 8: // Backspace
		return "Backspace"
	case 27: // Escape
		return "Escape"
	case 38: // Up
		return "Up"
	case 40: // Down
		return "Down"
	case 37: // Left
		return "Left"
	case 39: // Right
		return "Right"
	default:
		return fmt.Sprintf("Key%d", rawcode)
	}
}

// 获取修饰键名称
func getModifierKeyNameByRawcode(rawcode uint16) string {
	switch rawcode {
	case CMD_LEFT, CMD_RIGHT:
		return "Command"
	case SHIFT:
		return "Shift"
	case CAPS_LOCK:
		return "CapsLock"
	case OPTION, ALT_WIN:
		return getModifierKeyName("alt")
	case CTRL, CTRL_ALT1, CTRL_ALT2, CTRL_WIN:
		return "Control"
	default:
		return fmt.Sprintf("Mod%d", rawcode)
	}
}

// 构建组合键字符串
func buildCombinationKeyString(modifierKeys map[uint16]bool, keyName string) string {
	if len(modifierKeys) == 0 {
		return keyName
	}

	// 收集修饰键名称
	var modifiers []string
	for rawcode, pressed := range modifierKeys {
		if pressed {
			modifiers = append(modifiers, getModifierKeyNameByRawcode(rawcode))
		}
	}

	// 去重
	uniqueModifiers := make(map[string]bool)
	for _, mod := range modifiers {
		uniqueModifiers[mod] = true
	}

	// 转换为切片并排序
	var finalModifiers []string
	for mod := range uniqueModifiers {
		finalModifiers = append(finalModifiers, mod)
	}

	// 按固定顺序排序修饰键
	order := map[string]int{
		"Control":  1,
		"Alt":      2,
		"Option":   2, // macOS 中的 Alt 键
		"Shift":    3,
		"CapsLock": 4,
		"Command":  5,
	}
	sort.Slice(finalModifiers, func(i, j int) bool {
		return order[finalModifiers[i]] < order[finalModifiers[j]]
	})

	return strings.Join(finalModifiers, "+") + "+" + keyName
}

// SystemClicker 使用系统命令执行自定义按键连点
type SystemClicker struct {
	isRunning           bool
	delay               int
	boundKeys           string // 绑定的按键组合
	stopChan            chan bool
	logger              *Logger
	startTime           time.Time
	totalClicks         int
	isBinding           bool        // 是否正在绑定按键
	bindingCompleteChan chan string // 按键绑定完成通知通道
	stateChangeChan     chan bool   // 状态变化通知通道

	// 组合键绑定相关
	modifierKeys map[uint16]bool // 当前按下的修饰键
	lastMainKey  uint16          // 最后按下的主键
}

func NewSystemClicker() *SystemClicker {
	return &SystemClicker{
		isRunning:           false,
		delay:               1000,
		boundKeys:           "未绑定",
		stopChan:            make(chan bool),
		logger:              nil,
		startTime:           time.Time{},
		totalClicks:         0,
		isBinding:           false,
		bindingCompleteChan: make(chan string, 1),
		stateChangeChan:     make(chan bool, 1),

		// 初始化组合键绑定相关字段
		modifierKeys: make(map[uint16]bool),
		lastMainKey:  0,
	}
}

func (sc *SystemClicker) StartClicking() {
	if sc.isRunning {
		return
	}

	sc.isRunning = true
	sc.stopChan = make(chan bool)
	sc.startTime = time.Now()
	sc.totalClicks = 0

	// 发送状态变化通知
	select {
	case sc.stateChangeChan <- true:
	default:
	}

	// 记录开始日志
	if sc.logger != nil {
		sc.logger.LogStart(sc.delay, sc.boundKeys)
	}

	go func() {
		for {
			select {
			case <-sc.stopChan:
				return
			default:
				sc.performAction()
				sc.totalClicks++
				time.Sleep(time.Duration(sc.delay) * time.Millisecond)
			}
		}
	}()
}

func (sc *SystemClicker) StopClicking() {
	if !sc.isRunning {
		return
	}

	sc.isRunning = false
	close(sc.stopChan)

	// 发送状态变化通知
	select {
	case sc.stateChangeChan <- false:
	default:
	}

	// 记录停止日志
	if sc.logger != nil {
		duration := time.Since(sc.startTime)
		sc.logger.LogStop(duration, sc.totalClicks)
	}
}

func (sc *SystemClicker) performAction() {
	if sc.boundKeys == "未绑定" {
		// Windows 系统下不显示控制台日志
		if runtime.GOOS != "windows" {
			fmt.Println("未绑定任何按键，请先绑定按键")
		}
		return
	}

	var cmd *exec.Cmd

	// 根据绑定的按键类型自动判断执行方式
	switch runtime.GOOS {
	case "darwin":
		cmd = sc.getMacOSCommand()
	case "windows":
		// Windows 系统下，如果没有绑定按键，不执行任何操作
		if sc.boundKeys == "未绑定" {
			return
		}
		cmd = sc.getWindowsCommand()
	case "linux":
		cmd = sc.getLinuxCommand()
	default:
		// Windows 系统下不显示控制台日志
		if runtime.GOOS != "windows" {
			fmt.Printf("不支持的操作系统: %s\n", runtime.GOOS)
		}
		return
	}

	if cmd != nil {
		// 执行命令
		err := cmd.Run()
		if err != nil {
			// Windows 系统下不显示控制台日志
			if runtime.GOOS != "windows" {
				fmt.Printf("执行按键 %s 失败: %v\n", sc.boundKeys, err)
			}
			if sc.logger != nil {
				sc.logger.LogError("执行按键 %s 失败: %v", sc.boundKeys, err)
			}
		} else {
			// 命令成功执行
			// Windows 系统下不显示控制台日志
			if runtime.GOOS != "windows" {
				fmt.Printf("执行按键 %s 成功，延迟: %dms\n", sc.boundKeys, sc.delay)
			}
			if sc.logger != nil {
				sc.logger.LogClick(sc.boundKeys, sc.delay, time.Now())
			}
		}
	}
}

func (sc *SystemClicker) getMacOSCommand() *exec.Cmd {
	// 智能判断绑定的按键类型并执行相应命令
	if strings.Contains(sc.boundKeys, "左键") || strings.Contains(sc.boundKeys, "右键") || strings.Contains(sc.boundKeys, "中键") {
		// 鼠标点击
		var clickType string
		switch sc.boundKeys {
		case "左键":
			clickType = "c:." // 左键点击当前位置
		case "右键":
			clickType = "rc:." // 右键点击当前位置
		case "中键":
			clickType = "dc:." // 双击当前位置（模拟中键）
		default:
			clickType = "c:." // 默认左键
		}
		return exec.Command("cliclick", clickType)
	} else {
		// 键盘按键
		return sc.getMacOSKeyCommand()
	}
}

func (sc *SystemClicker) getWindowsCommand() *exec.Cmd {
	// 智能判断绑定的按键类型并执行相应命令
	if strings.Contains(sc.boundKeys, "左键") || strings.Contains(sc.boundKeys, "右键") || strings.Contains(sc.boundKeys, "中键") {
		// 鼠标点击 - 使用 Windows API 调用
		var psScript string
		switch sc.boundKeys {
		case "左键":
			psScript = `
Add-Type -TypeDefinition @"
using System;
using System.Runtime.InteropServices;
public class Win32 {
    [DllImport("user32.dll")]
    public static extern void mouse_event(uint dwFlags, uint dx, uint dy, uint dwData, UIntPtr dwExtraInfo);
    
    public const uint MOUSEEVENTF_LEFTDOWN = 0x0002;
    public const uint MOUSEEVENTF_LEFTUP = 0x0004;
}
"@
[Win32]::mouse_event([Win32]::MOUSEEVENTF_LEFTDOWN, 0, 0, 0, [UIntPtr]::Zero)
[Win32]::mouse_event([Win32]::MOUSEEVENTF_LEFTUP, 0, 0, 0, [UIntPtr]::Zero)
`
		case "右键":
			psScript = `
Add-Type -TypeDefinition @"
using System;
using System.Runtime.InteropServices;
public class Win32 {
    [DllImport("user32.dll")]
    public static extern void mouse_event(uint dwFlags, uint dx, uint dy, uint dwData, UIntPtr dwExtraInfo);
    
    public const uint MOUSEEVENTF_RIGHTDOWN = 0x0008;
    public const uint MOUSEEVENTF_RIGHTUP = 0x0010;
}
"@
[Win32]::mouse_event([Win32]::MOUSEEVENTF_RIGHTDOWN, 0, 0, 0, [UIntPtr]::Zero)
[Win32]::mouse_event([Win32]::MOUSEEVENTF_RIGHTUP, 0, 0, 0, [UIntPtr]::Zero)
`
		case "中键":
			psScript = `
Add-Type -TypeDefinition @"
using System;
using System.Runtime.InteropServices;
public class Win32 {
    [DllImport("user32.dll")]
    public static extern void mouse_event(uint dwFlags, uint dx, uint dy, uint dwData, UIntPtr dwExtraInfo);
    
    public const uint MOUSEEVENTF_MIDDLEDOWN = 0x0020;
    public const uint MOUSEEVENTF_MIDDLEUP = 0x0040;
}
"@
[Win32]::mouse_event([Win32]::MOUSEEVENTF_MIDDLEDOWN, 0, 0, 0, [UIntPtr]::Zero)
[Win32]::mouse_event([Win32]::MOUSEEVENTF_MIDDLEUP, 0, 0, 0, [UIntPtr]::Zero)
`
		default:
			psScript = `
Add-Type -TypeDefinition @"
using System;
using System.Runtime.InteropServices;
public class Win32 {
    [DllImport("user32.dll")]
    public static extern void mouse_event(uint dwFlags, uint dx, uint dy, uint dwData, UIntPtr dwExtraInfo);
    
    public const uint MOUSEEVENTF_LEFTDOWN = 0x0002;
    public const uint MOUSEEVENTF_LEFTUP = 0x0004;
}
"@
[Win32]::mouse_event([Win32]::MOUSEEVENTF_LEFTDOWN, 0, 0, 0, [UIntPtr]::Zero)
[Win32]::mouse_event([Win32]::MOUSEEVENTF_LEFTUP, 0, 0, 0, [UIntPtr]::Zero)
`
		}
		return exec.Command("powershell", "-Command", psScript)
	} else {
		// 键盘按键
		return sc.getWindowsKeyCommand()
	}
}

func (sc *SystemClicker) getLinuxCommand() *exec.Cmd {
	// 智能判断绑定的按键类型并执行相应命令
	if strings.Contains(sc.boundKeys, "左键") || strings.Contains(sc.boundKeys, "右键") || strings.Contains(sc.boundKeys, "中键") {
		// 鼠标点击
		var button string
		switch sc.boundKeys {
		case "左键":
			button = "1"
		case "右键":
			button = "3"
		case "中键":
			button = "2"
		default:
			button = "1"
		}

		return exec.Command("xdotool", "click", button)
	} else {
		// 键盘按键
		return sc.getLinuxKeyCommand()
	}
}

// macOS 键盘按键命令
func (sc *SystemClicker) getMacOSKeyCommand() *exec.Cmd {
	// 解析组合键，如 "Command+C", "Shift+A" 等
	keys := strings.Split(sc.boundKeys, "+")

	var script string
	if len(keys) > 1 {
		// 组合键
		modifiers := keys[:len(keys)-1]
		key := keys[len(keys)-1]

		modifierScript := ""
		for _, mod := range modifiers {
			switch strings.ToLower(mod) {
			case "command", "cmd":
				modifierScript += "command down, "
			case "shift":
				modifierScript += "shift down, "
			case "option", "alt":
				modifierScript += "option down, "
			case "control", "ctrl":
				modifierScript += "control down, "
			}
		}

		script = fmt.Sprintf(`tell application "System Events"
    try
        key code %s using {%s}
        return "success"
    on error errMsg
        return "error:" & errMsg
    end try
end tell`, sc.getKeyCode(key), modifierScript[:len(modifierScript)-2])
	} else {
		// 单键
		script = fmt.Sprintf(`tell application "System Events"
    try
        keystroke "%s"
        return "success"
    on error errMsg
        return "error:" & errMsg
    end try
end tell`, sc.boundKeys)
	}

	return exec.Command("osascript", "-e", script)
}

// Windows 键盘按键命令
func (sc *SystemClicker) getWindowsKeyCommand() *exec.Cmd {
	// 解析组合键
	keys := strings.Split(sc.boundKeys, "+")

	var key string
	if len(keys) > 1 {
		// 组合键
		modifiers := keys[:len(keys)-1]
		mainKey := keys[len(keys)-1]

		modifierStr := ""
		for _, mod := range modifiers {
			switch strings.ToLower(mod) {
			case "command", "cmd":
				modifierStr += "^" // Ctrl
			case "shift":
				modifierStr += "+"
			case "option", "alt":
				modifierStr += "%"
			case "control", "ctrl":
				modifierStr += "^"
			}
		}

		key = modifierStr + mainKey
	} else {
		// 单键
		key = sc.boundKeys
	}

	// 使用 PowerShell 执行键盘按键
	psScript := fmt.Sprintf(`Add-Type -AssemblyName System.Windows.Forms; [System.Windows.Forms.SendKeys]::SendWait("%s")`, key)
	return exec.Command("powershell", "-Command", psScript)
}

// Linux 键盘按键命令
func (sc *SystemClicker) getLinuxKeyCommand() *exec.Cmd {
	// 解析组合键
	keys := strings.Split(sc.boundKeys, "+")

	var key string
	if len(keys) > 1 {
		// 组合键
		modifiers := keys[:len(keys)-1]
		mainKey := keys[len(keys)-1]

		modifierStr := ""
		for _, mod := range modifiers {
			switch strings.ToLower(mod) {
			case "command", "cmd":
				modifierStr += "ctrl+"
			case "shift":
				modifierStr += "shift+"
			case "option", "alt":
				modifierStr += "alt+"
			case "control", "ctrl":
				modifierStr += "ctrl+"
			}
		}

		key = modifierStr + mainKey
	} else {
		// 单键
		key = sc.boundKeys
	}

	return exec.Command("xdotool", "key", key)
}

// 获取按键的keycode（macOS）
func (sc *SystemClicker) getKeyCode(key string) string {
	switch strings.ToLower(key) {
	case "a":
		return "0"
	case "s":
		return "1"
	case "d":
		return "2"
	case "f":
		return "3"
	case "h":
		return "4"
	case "g":
		return "5"
	case "z":
		return "6"
	case "x":
		return "7"
	case "c":
		return "8"
	case "v":
		return "9"
	case "b":
		return "11"
	case "q":
		return "12"
	case "w":
		return "13"
	case "e":
		return "14"
	case "r":
		return "15"
	case "y":
		return "16"
	case "t":
		return "17"
	case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0":
		return key
	case "return":
		return "36"
	case "tab":
		return "48"
	case "space":
		return "49"
	case "delete":
		return "51"
	case "escape":
		return "53"
	case "up":
		return "126"
	case "down":
		return "125"
	case "left":
		return "123"
	case "right":
		return "124"
	default:
		return "49" // 默认空格键
	}
}

// 设置延迟
func (sc *SystemClicker) SetDelay(delay int) {
	sc.delay = delay
}

// 设置绑定的按键
func (sc *SystemClicker) SetBoundKeys(keys string) {
	sc.boundKeys = keys
}

// 获取绑定的按键
func (sc *SystemClicker) GetBoundKeys() string {
	return sc.boundKeys
}

// 开始绑定按键
func (sc *SystemClicker) StartBinding() {
	sc.isBinding = true
	sc.modifierKeys = make(map[uint16]bool)
	sc.lastMainKey = 0
	// Windows 系统下不显示控制台日志
	if runtime.GOOS != "windows" {
		fmt.Println("请按下要绑定的按键组合...")
	}
}

// 获取绑定完成通道
func (sc *SystemClicker) GetBindingCompleteChan() <-chan string {
	return sc.bindingCompleteChan
}

// 获取状态变化通道
func (sc *SystemClicker) GetStateChangeChan() <-chan bool {
	return sc.stateChangeChan
}

// 停止绑定按键
func (sc *SystemClicker) StopBinding() {
	sc.isBinding = false
}

// 获取运行状态
func (sc *SystemClicker) IsRunning() bool {
	return sc.isRunning
}

// 设置日志器
func (sc *SystemClicker) SetLogger(logger *Logger) {
	sc.logger = logger
}

// 获取统计信息
func (sc *SystemClicker) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})
	stats["is_running"] = sc.isRunning
	stats["delay"] = sc.delay
	stats["bound_keys"] = sc.boundKeys
	stats["total_clicks"] = sc.totalClicks

	if !sc.startTime.IsZero() {
		stats["start_time"] = sc.startTime
		if sc.isRunning {
			stats["running_duration"] = time.Since(sc.startTime)
		}
	}

	return stats
}

// 启动快捷键监听
func (sc *SystemClicker) StartHotkeyListener() {
	go func() {
		// 监听全局快捷键
		evChan := hook.Start()
		defer hook.End()

		// 用于跟踪修饰键状态
		var modifierKeys map[uint16]bool = make(map[uint16]bool)

		for ev := range evChan {
			// 跟踪修饰键状态
			if ev.Kind == hook.KeyDown {
				modifierKeys[ev.Rawcode] = true
			} else if ev.Kind == hook.KeyUp {
				delete(modifierKeys, ev.Rawcode)
			}

			// 检查是否是快捷键组合
			if ev.Kind == hook.KeyDown {
				// 根据操作系统使用不同的快捷键
				if runtime.GOOS == "windows" {
					// Windows: Ctrl+F 开始连续点击
					if ev.Rawcode == 70 && (modifierKeys[CTRL_WIN] || modifierKeys[CTRL] || modifierKeys[CTRL_ALT1] || modifierKeys[CTRL_ALT2]) {
						if !sc.isRunning {
							// Windows 系统下不显示控制台日志
							if runtime.GOOS != "windows" {
								fmt.Println("检测到快捷键 Ctrl+F，开始连续点击")
							}
							sc.StartClicking()
						}
					}
					// Windows: Ctrl+G 停止连续点击
					if ev.Rawcode == 71 && (modifierKeys[CTRL_WIN] || modifierKeys[CTRL] || modifierKeys[CTRL_ALT1] || modifierKeys[CTRL_ALT2]) {
						if sc.isRunning {
							// Windows 系统下不显示控制台日志
							if runtime.GOOS != "windows" {
								fmt.Println("检测到快捷键 Ctrl+G，停止连续点击")
							}
							sc.StopClicking()
						}
					}
				} else {
					// macOS: Command+F 开始连续点击
					if ev.Rawcode == KEY_F && (modifierKeys[CMD_LEFT] || modifierKeys[CMD_RIGHT]) {
						if !sc.isRunning {
							fmt.Println("检测到快捷键 Command+F，开始连续点击")
							sc.StartClicking()
						}
					}
					// macOS: Command+G 停止连续点击
					if ev.Rawcode == KEY_G && (modifierKeys[CMD_LEFT] || modifierKeys[CMD_RIGHT]) {
						if sc.isRunning {
							fmt.Println("检测到快捷键 Command+G，停止连续点击")
							sc.StopClicking()
						}
					}
				}
			}

			// 组合键绑定：在按下时记录修饰键，在抬起主键时完成绑定
			if sc.isBinding {
				if ev.Kind == hook.KeyDown {
					sc.recordKeyPress(ev)
				} else if ev.Kind == hook.KeyUp && sc.isBinding {
					// 在按键抬起时检查是否完成绑定
					sc.checkKeyBindingComplete(ev)
				} else if ev.Kind == hook.MouseDown {
					sc.recordMousePress(ev)
				}
			}
		}
	}()
}

// 记录按键按下（用于组合键绑定）
func (sc *SystemClicker) recordKeyPress(ev hook.Event) {
	// 检测是否是修饰键
	isModifier := false
	var modifierName string

	switch ev.Rawcode {
	case CMD_LEFT, CMD_RIGHT:
		isModifier = true
		modifierName = "Command"
	case SHIFT:
		isModifier = true
		modifierName = "Shift"
	case CAPS_LOCK:
		isModifier = true
		modifierName = "CapsLock"
	case OPTION, ALT_WIN:
		isModifier = true
		modifierName = getModifierKeyName("alt")
	case CTRL, CTRL_ALT1, CTRL_ALT2, CTRL_WIN:
		isModifier = true
		modifierName = "Control"
	}

	if isModifier {
		// 记录修饰键
		sc.modifierKeys[ev.Rawcode] = true
		// Windows 系统下不显示控制台日志
		if runtime.GOOS != "windows" {
			fmt.Printf("记录修饰键: %s\n", modifierName)
		}
	} else {
		// 记录主键，但不立即完成绑定
		var keyName string

		// 使用统一的按键名称获取函数
		keyName = getKeyName(ev.Rawcode)

		// 构建组合键字符串
		sc.boundKeys = buildCombinationKeyString(sc.modifierKeys, keyName)

		// 完成绑定
		sc.isBinding = false
		// Windows 系统下不显示控制台日志
		if runtime.GOOS != "windows" {
			fmt.Printf("组合键绑定完成: %s\n", sc.boundKeys)
		}

		// 通知主线程更新界面（通过通道）
		if sc.bindingCompleteChan != nil {
			select {
			case sc.bindingCompleteChan <- sc.boundKeys:
			default:
			}
		}
	}
}

// 记录鼠标按下（用于组合键绑定）
func (sc *SystemClicker) recordMousePress(ev hook.Event) {
	// 鼠标按下时直接完成绑定
	switch ev.Button {
	case 1:
		sc.boundKeys = "左键"
	case 2:
		sc.boundKeys = "中键"
	case 3:
		sc.boundKeys = "右键"
	default:
		sc.boundKeys = fmt.Sprintf("鼠标按钮%d", ev.Button)
	}

	// 完成绑定
	sc.isBinding = false
	// Windows 系统下不显示控制台日志
	if runtime.GOOS != "windows" {
		fmt.Printf("鼠标绑定完成: %s\n", sc.boundKeys)
	}

	// 通知主线程更新界面（通过通道）
	if sc.bindingCompleteChan != nil {
		select {
		case sc.bindingCompleteChan <- sc.boundKeys:
		default:
		}
	}
}

// 检查按键绑定是否完成（在按键抬起时调用）
func (sc *SystemClicker) checkKeyBindingComplete(ev hook.Event) {
	// 如果没有记录主键，说明还在等待主键
	if sc.lastMainKey == 0 {
		return
	}

	// 检查抬起的是否是主键
	if ev.Rawcode == sc.lastMainKey {
		// 完成绑定
		var keyName string

		// 使用统一的按键名称获取函数
		keyName = getKeyName(sc.lastMainKey)

		// 构建组合键字符串
		sc.boundKeys = buildCombinationKeyString(sc.modifierKeys, keyName)

		// 完成绑定
		sc.isBinding = false
		// Windows 系统下不显示控制台日志
		if runtime.GOOS != "windows" {
			fmt.Printf("组合键绑定完成: %s\n", sc.boundKeys)
		}

		// 通知主线程更新界面（通过通道）
		if sc.bindingCompleteChan != nil {
			select {
			case sc.bindingCompleteChan <- sc.boundKeys:
			default:
			}
		}
	}
}
