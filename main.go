package main

import (
	"fmt"
	"image/color"
	"log"
	"runtime"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// 全局变量
var (
	clicker *SystemClicker
	config  *Config
)

// 自定义主题
type myTheme struct {
	font fyne.Resource
}

func (m *myTheme) Font(style fyne.TextStyle) fyne.Resource {
	return m.font
}

func (m *myTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}

func (m *myTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (m *myTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}

// 初始化全局点击器和配置
func init() {
	var err error
	config, err = LoadConfig()
	if err != nil {
		// Windows 系统下不显示控制台日志
		if runtime.GOOS != "windows" {
			fmt.Printf("加载配置失败: %v，使用默认配置\n", err)
		}
		config = createDefaultConfig()
	}

	logger := NewLogger(config.Logging.Enabled, config.Logging.Level)
	clicker = NewSystemClicker()
	clicker.SetDelay(config.DefaultDelay)
	clicker.SetBoundKeys("未绑定")
	clicker.SetLogger(logger)
	// 不在 init 中启动快捷键监听器，而是在界面显示后启动
}

func main() {
	a := app.New()

	// 使用嵌入字体
	fontResource := fyne.NewStaticResource("SourceHanSansSC-Bold", SourceHanSansSCBold)
	a.Settings().SetTheme(&myTheme{font: fontResource})
	// Windows 系统下不显示控制台日志
	if runtime.GOOS != "windows" {
		log.Println("已加载中文字体")
	}

	runMainApplication(a)
}

// 以下保持原 main.go 界面逻辑
func runMainApplication(a fyne.App) {
	w := a.NewWindow("连续点击器")
	w.Resize(fyne.NewSize(float32(config.UI.WindowWidth), float32(config.UI.WindowHeight)))
	w.CenterOnScreen()

	delayLabel := widget.NewLabel("点击延迟 (毫秒):")
	delayEntry := widget.NewEntry()
	delayEntry.SetText(strconv.Itoa(config.DefaultDelay))
	delayEntry.OnChanged = func(text string) {
		if delay, err := strconv.Atoi(text); err == nil {
			clicker.SetDelay(delay)
		}
	}

	bindKeyButton := widget.NewButton("点击绑定按键", nil)
	currentKeysLabel := widget.NewLabel("当前绑定: 未绑定")
	currentKeysLabel.TextStyle = fyne.TextStyle{Bold: true}

	bindKeyButton.OnTapped = func() {
		if clicker.GetBoundKeys() == "未绑定" {
			go func() {
				time.Sleep(100 * time.Millisecond)
				clicker.StartBinding()
				w.Canvas().Refresh(bindKeyButton)
				bindKeyButton.SetText("正在监听... 请按下要绑定的按键")
				bindKeyButton.Disable()

				go func() {
					select {
					case boundKeys := <-clicker.GetBindingCompleteChan():
						w.Canvas().Refresh(currentKeysLabel)
						w.Canvas().Refresh(bindKeyButton)
						currentKeysLabel.SetText(fmt.Sprintf("当前绑定: %s", boundKeys))
						bindKeyButton.SetText(fmt.Sprintf("重新绑定 (%s)", boundKeys))
						bindKeyButton.Enable()
					}
				}()
			}()
		} else {
			clicker.SetBoundKeys("未绑定")
			bindKeyButton.SetText("点击绑定按键")
			currentKeysLabel.SetText("当前绑定: 未绑定")
		}
	}

	hotkey := config.GetCurrentOSHotkey()
	startKeyLabel := widget.NewLabel(fmt.Sprintf("开始快捷键: %s", hotkey.Start))
	stopKeyLabel := widget.NewLabel(fmt.Sprintf("停止快捷键: %s", hotkey.Stop))

	statusLabel := widget.NewLabel("状态: 未运行")
	statusLabel.TextStyle = fyne.TextStyle{Bold: true}

	startButton := widget.NewButton("开始连续点击", func() {
		clicker.StartClicking()
		// 只更新状态标签，不直接操作按钮
		statusLabel.SetText("状态: 运行中")
	})
	stopButton := widget.NewButton("停止连续点击", func() {
		clicker.StopClicking()
		// 只更新状态标签，不直接操作按钮
		statusLabel.SetText("状态: 未运行")
	})

	// 创建统一的状态更新函数，通过面板状态控制按钮
	updateButtonStates := func(isRunning bool) {
		if isRunning {
			statusLabel.SetText("状态: 运行中")
			startButton.Disable()
			stopButton.Enable()
			bindKeyButton.Disable() // 运行时禁用重新绑定按钮
		} else {
			statusLabel.SetText("状态: 未运行")
			startButton.Enable()
			stopButton.Disable()
			bindKeyButton.Enable() // 停止时启用重新绑定按钮
		}
	}

	// 初始化时设置停止按钮为不可用状态
	stopButton.Disable()

	testButton := widget.NewButton("测试点击", func() {
		// Windows 系统下，如果没有绑定按键，不执行任何操作
		if runtime.GOOS == "windows" && clicker.GetBoundKeys() == "未绑定" {
			return
		}
		clicker.performAction()
	})

	instructions := widget.NewLabel(`使用说明:
1. 设置连点延迟时间（毫秒）
2. 点击"绑定按键"按钮，然后按下并释放要绑定的按键组合
3. 支持鼠标、键盘、组合键
4. 点击"开始连续点击"按钮
5. 使用快捷键 Command+F 开始，Command+G 停止
注意: 程序运行期间请确保目标窗口处于活动状态`)

	content := container.NewVBox(
		widget.NewLabel("连续点击器"),
		widget.NewSeparator(),
		delayLabel,
		delayEntry,
		widget.NewSeparator(),
		bindKeyButton,
		currentKeysLabel,
		widget.NewSeparator(),
		startKeyLabel,
		stopKeyLabel,
		widget.NewSeparator(),
		startButton,
		stopButton,
		testButton,
		widget.NewSeparator(),
		statusLabel,
		widget.NewSeparator(),
		instructions,
	)

	w.SetContent(content)

	// 在界面显示后安全地启动快捷键监听器
	go func() {
		// 等待界面完全显示
		time.Sleep(500 * time.Millisecond)

		// 启动快捷键监听器
		clicker.StartHotkeyListener()
	}()

	// 启动状态变化监听器
	go func() {
		for {
			select {
			case isRunning := <-clicker.GetStateChangeChan():
				// 在主线程中更新界面状态
				updateButtonStates(isRunning)
			}
		}
	}()

	w.ShowAndRun()
}
