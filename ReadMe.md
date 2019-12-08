#go http server with reload
####一、 简介
* 写完代码每次都要重新编译，开发阶段很是繁琐，本demo使用 /_reload 进行重启服务
* 核心简述：使用channel阻塞主进程，然后监听 _reload chan 信号，调用exec 重启服务
```
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




func Reload(server *http.Server) {
	for {
		<-reload
		fmt.Println("结束server")
		if err := server.Shutdown(context.Background()); err != nil {
			fmt.Println(err)
		}
		env := append(
			os.Environ(),
			"ENDLESS_CONTINUE=1",
		)

		indexFile := os.Getenv("indexFile")
		cmd := exec.Command("go", "run", indexFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = env
		fmt.Println("重启server: " + indexFile)

		err := cmd.Start()
		if err != nil {
			log.Fatalf("Restart: Failed to launch, error: %v", err)
		}
	}
	return
}


http://localhost:8080/_reload
http://localhost:8080/index/test/test

```
