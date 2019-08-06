package appclient

import (
	"collector/assets/ico"
	"collector/internal/interfaces"
	"collector/internal/pkg/producter/collectzip"
	"context"
	"gcom/gtools/systray"
	"gcom/gwin"
	"os/exec"
)

const (
	apptitle = "MES Collect"
	apptip   = "MES Collect-Zip"
)

type app struct {
	ctx    context.Context
	cancel context.CancelFunc

	log       interfaces.ILog
	app       interfaces.IApp
	cache     interfaces.ICache
	config    interfaces.IConfig
	producter interfaces.IProducter
}

func New(log interfaces.ILog, cache interfaces.ICache, config interfaces.IConfig) *app {
	ctx, cancel := context.WithCancel(context.Background())
	return &app{
		log:    log,
		ctx:    ctx,
		cancel: cancel,

		producter: collectzip.New(log, cache, config, nil),
	}
}

//Run 程序运行
func (this *app) Serve() {

	go systray.Run(this.onReady, this.onExit)

	//阻塞主程
	this.listen()
}

//Run 程序运行
func (this *app) Run() {
	//TODO 配置更新 网络组件启动 文件扫描启动
	
	// 文件扫描启动
	this.producter.Run()
}

//Run 程序运行
func (this *app) Stop() {
	//TODO 释放资源
}

//Run 程序退出
func (this *app) Exit() {
	//TODO 释放资源
	this.cancel()
}

//Run 程序退出
func (this *app) listen() {
	for {
		select {
		case <-this.ctx.Done():
			return
		}
	}
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
				this.Run()
				mStart.Disable()
				mStop.Enable()

			case <-mStop.ClickedCh:
				this.log.Info("app click stop")
				this.Stop()
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
