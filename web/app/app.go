package app

import (
	"context"
	"fmt"
	"gcom/gtools/gdb"
	"html/template"
	"net/http"
	"time"
)

type IAbout interface {
	Serve()
	Stop()
}
type tpl struct {
	t      string
	list   []FilesInfo
	ctx    context.Context
	cancel context.CancelFunc
}

type FilesInfo struct {
	Date   string
	IP     string
	Path   string
	Status string
}

func New(t string) IAbout {
	ctx, cancel := context.WithCancel(context.Background())
	return &tpl{
		t:      t,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (this *tpl) about(w http.ResponseWriter, r *http.Request) {
	this.updateInfo()
	t, _ := template.ParseFiles(this.t)
	t.Execute(w, this.list)
}

func (this *tpl) Serve() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", this.about)

	svr := &http.Server{
		Addr:    ":497",
		Handler: mux,
	}

	go func() {
		for {
			select {
			case <-this.ctx.Done():
				svr.Shutdown(this.ctx)
				return
			}
		}

	}()

	svr.ListenAndServe()
}

func (this *tpl) Stop() {
	this.cancel()
}

func (this *tpl) updateInfo() {
	ntime := time.Now().Format("2006-01-02")
	queryssql := fmt.Sprintf("SELECT FModTime,FID,FilePath,FFinish FROM [iFixsvr_JF_Files] WHERE FModTime like '%s%s%s'", "%", ntime, "%")
	res, err := gdb.Querys(queryssql)

	if err != nil {
		return
	}

	fis := make([]FilesInfo, 0)

	for _, v := range res {
		if len(v) == 4 {
			fi := FilesInfo{}
			fi.Date = v[0]
			fi.IP = v[1]
			fi.Path = v[2]
			fi.Status = v[3]

			fis = append(fis, fi)
		}
	}
	this.list = fis
}
