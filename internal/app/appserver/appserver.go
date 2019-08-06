package appserver

import (
	"collector/assets/ico"
	"collector/internal/interfaces"
	"gcom/gtools/systray"
	"gcom/gwin"
	"os/exec"
)

const (
	apptitle = "MES Collect"
	apptip   = "MES Collect-Zip"
)

type app struct {
	log interfaces.ILog
	app interfaces.IApp
}

func New(log interfaces.ILog, apps interfaces.IApp) *app {
	return &app{
		log: log,
		app: apps,
	}
}

//Run 程序运行
func (this *app) Serve() {
	//程序托盘图标
	go systray.Run(this.onReady, this.onExit)

	//阻塞主程
	this.app.Serve()
}

//Run 程序运行
func (this *app) Run() {
	//TODO 配置更新 网络组件启动 文件扫描启动
}

//Run 程序运行
func (this *app) Stop() {
	//TODO 释放资源
}

func (this *app) onReady() {

	systray.SetIcon(ico.Data)
	systray.SetTitle(apptitle)
	systray.SetTooltip(apptip)
	go func() {
		//关于 ...
		mAbout := systray.AddMenuItem("关于", "App About")
		systray.AddSeparator()

		//启动 ...
		mStart := systray.AddMenuItem("启动", "App Start")
		systray.AddSeparator()

		//停止 ...
		mStop := systray.AddMenuItem("停止", "App Stop")
		systray.AddSeparator()

		//退出 ...
		mQuit := systray.AddMenuItem("退出", "Quit")

		//设置初始状态
		go func() {
			this.log.Info("app init")
			mStart.ClickedCh <- struct{}{}
		}()

		for {
			select {
			case <-mQuit.ClickedCh:
				result := gwin.MessageBox(apptitle, "确认退出程序吗", gwin.MB_OKCANCEL|gwin.MB_ICONWARNING)
				if result == 1 {
					this.log.Info("app click quit")
					this.Stop()
					systray.Quit()
					return
				}

			case <-mAbout.ClickedCh:
				cmd := exec.Command("cmd", "/c start www.baidu.com")
				if err := cmd.Start(); err != nil {
					gwin.MessageBox(apptitle, err.Error(), gwin.MB_OK)
				}

			case <-mStart.ClickedCh:
				this.log.Info("app click start")
				mStart.Disable()
				mStop.Enable()

			case <-mStop.ClickedCh:
				this.log.Info("app click stop")
				mStart.Enable()
				mStop.Disable()

			}
		}
	}()

}

func (this *app) onExit() {
	//释放主程
	this.app.Exit()
}
