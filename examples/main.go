package main

import (
	"collector/gcommon/cache"
	"collector/gcommon/config"
	"fmt"
)

func main() {
	x := config.New("./collector.json")
	v := x.Get("port")
	fmt.Println(v)
	c := cache.New("./cache", "db")
	c.Set("host", v)

	m := c.GetMap()
	fmt.Println(m)
}
