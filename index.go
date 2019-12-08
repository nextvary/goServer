package main

import (
	"./app"
	"./controller/index"
	mytest "./controller/test"
	"fmt"
	"os"
	"runtime"
)

func Init() {
	registerRoute()
	_, file, _, _ := runtime.Caller(0) //找到入口文件
	os.Setenv("indexFile", file)
	fmt.Println(app.Mapping)
}
func registerRoute() {
	app.Static["/assets"] = "./static"       //静态资源
	app.AutoRouter(&index.IndexController{}) //路由注册 /index/index/index1
	app.AutoRouter(&index.TestController{})
	app.Router("test2/test", &mytest.Test1Controller{})
}
func main() {
	Init()
	go app.RunOn(":8080")
	for {
		server := <-app.ServerCh //使用chan 阻塞server 主进程
		fmt.Println("server start")
		//go app.Stop(server)
		go app.Reload(server) //监听reload 信号 ，shutdown 服务，然后exec 重启服务
		server.ListenAndServe()
	}
}
