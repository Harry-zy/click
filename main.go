package main

import (
	"fmt"
	"image/color"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
		fmt.Printf("加载配置失败: %v，使用默认配置\n", err)
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
	log.Println("已加载中文字体")

	runMainApplication(a)
}

// 以下保持原 main.go 界面逻辑
func runMainApplication(a fyne.App) {
	w := a.NewWindow("连续点击器")
	w.Resize(fyne.NewSize(float32(config.UI.WindowWidth), float32(config.UI.WindowHeight)))
	w.CenterOnScreen()

	checkSystemPermissions(w)

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

				bindDialog := dialog.NewInformation("按键绑定", "请按下要绑定的按键组合\n支持：鼠标、键盘、组合键\n按下并释放按键即可完成绑定", w)
				bindDialog.Show()

				go func() {
					select {
					case boundKeys := <-clicker.GetBindingCompleteChan():
						bindDialog.Hide()
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

	startButton := widget.NewButton("开始连续点击", func() {
		clicker.StartClicking()
		dialog.ShowInformation("提示", "连续点击已开始！\n使用快捷键停止。", w)
	})
	stopButton := widget.NewButton("停止连续点击", func() {
		clicker.StopClicking()
		dialog.ShowInformation("提示", "连续点击已停止！", w)
	})

	testButton := widget.NewButton("测试点击", func() {
		testSingleClick(w)
	})

	statusLabel := widget.NewLabel("状态: 未运行")
	statusLabel.TextStyle = fyne.TextStyle{Bold: true}

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

		// 检查系统权限后再启动快捷键监听器
		if checkSystemPermissionsSilent() {
			fmt.Println("启动快捷键监听器...")
			clicker.StartHotkeyListener()
		} else {
			fmt.Println("系统权限不足，快捷键功能不可用")
		}
	}()

	w.ShowAndRun()
}

func checkSystemPermissions(w fyne.Window) {
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("osascript", "-e", `tell application "System Events" to return "test"`)
		if _, err := cmd.Output(); err != nil {
			showPermissionGuide(w)
		}
	}
}

func showPermissionGuide(w fyne.Window) {
	guideText := `需要系统权限设置
在macOS上，连续点击器需要"辅助功能"权限才能正常工作。
设置步骤：
1. 打开"系统偏好设置"
2. 点击"安全性与隐私"
3. 选择"隐私"标签
4. 在左侧列表中选择"辅助功能"
5. 点击锁图标解锁设置
6. 将"连续点击器"添加到允许列表
7. 重启程序
设置完成后，程序就能正常执行点击操作了。`
	dialog.ShowInformation("权限设置指导", guideText, w)
}

func testSingleClick(w fyne.Window) {
	clicker.performAction()
	dialog.ShowInformation("测试结果", "测试动作已执行！\n如果看到相应效果，说明权限设置正确。", w)
}

// checkSystemPermissionsSilent 静默检查系统权限，不显示对话框
func checkSystemPermissionsSilent() bool {
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("osascript", "-e", `tell application "System Events" to return "test"`)
		if _, err := cmd.Output(); err != nil {
			return false
		}
		return true
	}
	return true
}
