package apptray

import (
	"collector/pkg/interfaces"
	"collector/web/app"
	"context"
	myapp "gcom/gtools/myapp"
	"sync"

	"github.com/skratchdot/open-golang/open"
)

var once sync.Once

type apptray struct {
	ctx    context.Context
	cancel context.CancelFunc
	config interfaces.IConfig
	about  app.IAbout
}

func New(ctx context.Context, cancel context.CancelFunc, config interfaces.IConfig) *apptray {
	return &apptray{
		ctx:    ctx,
		cancel: cancel,
		config: config,
		about:  app.New(),
	}
}

func (this *apptray) AppTray() {
	go this.about.Serve()
	myapp.Run(this.onReady, this.onExit)
}

func (this *apptray) onReady() {
	appname := this.config.Get("appname")
	myapp.SetIcon(ico)
	myapp.SetTitle(appname)
	myapp.SetTooltip(appname)

	go func() {
		//AddMenuItem
		mAbout := myapp.AddMenuItem("About", "About myapp")
		myapp.AddSeparator()
		mQuit := myapp.AddMenuItem("Quit", "Quit myapp")

		for {
			select {
			case <-mQuit.ClickedCh:
				this.about.Stop()
				myapp.Quit()
				return
			case <-mAbout.ClickedCh:
				go func() {
					open.Run("http://localhost:497")
				}()
			}
		}
	}()

}

func (this *apptray) onExit() {
	this.cancel()
}
