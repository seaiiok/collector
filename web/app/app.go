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

func New() IAbout {
	ctx, cancel := context.WithCancel(context.Background())
	return &tpl{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (this *tpl) about(w http.ResponseWriter, r *http.Request) {
	this.updateInfo()
	// t, _ := template.ParseFiles("./index.tpl")
	t := template.Must(template.New("").Parse(tplString))
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

var tplString = `
<!DOCTYPE html>

<head>
    <meta charset="UTF-8">
    <title>Collect</title>
</head>

<body>
    <style type="text/css">
        table.gridtable {
            font-family: verdana, arial, sans-serif;
            font-size: 11px;
            color: #303133;
            width: 80%
        }

        table.gridtable th {
            padding: 10px;
            color: #ffffff;
            background-color: #303133;
        }

        table.gridtable td {
            padding: 8px;
            background-color: #F2F6FC;
        }
    </style>

    <!-- Table goes JF room offline File Info by MES-->
    <table class="gridtable">
        <tr>
            <th>日 期</th>
            <th>文件地址</th>
            <th>文件路径</th>
            <th>状 态</th>
        </tr>

        {{ range $key,$value:=. }}
        <tr class="ctl">
            <td>{{$value.Date}}</td>
            <td>{{$value.IP}}</td>
            <td>{{$value.Path}}</td>
          
            <td id="{{$key}}" style="color: grey"></td>

            <script>
                var x = document.getElementById("{{$key}}")
                if ({{ $value.Status }} == "1") {
                    x.style.color = "green"
                    x.innerHTML = "完成"
                } else if ({{ $value.Status }} == "0") {
                    x.style.color = "red"
                    x.innerHTML = "未完成"
                }else {
                    x.innerHTML = "未知"
                }
            </script>

        </tr>
        {{ end }}
    </table>

</body>

</html>
`
