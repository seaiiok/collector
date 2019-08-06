package app

import "context"

type App struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func New() *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (this *App) Serve() {
	for {
		select {
		case <-this.ctx.Done():
			return
		// default:
		}
	}
}

func (this *App) Exit() {
	this.cancel()
}
