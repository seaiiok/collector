package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	str := "解放了大数据咖啡碱撒是佛教四大金佛阿斯顿就是开发商贷款龙卷风爱神的箭欧式的积分按时到房间啊手动拍摄发生的卡龙卷风阿斯顿夫卡是的就放假啊速度快"

	for i := 0; i < 2; i++ {
		str = str + str
	}

	fmt.Println("str:", len(str))
	fmt.Fprint(w, str)
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	log.Fatal(http.ListenAndServe(":8080", router))
}
