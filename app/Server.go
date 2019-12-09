package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"sync"
)

var reload = make(chan bool, 1)
var StopCh = make(chan bool, 1)
var ServerCh = make(chan *http.Server)

type handler struct {
	p sync.Pool
}

func newHandler() *handler {
	h := &handler{}
	h.p.New = func() interface{} {
		return &Context{}
	}

	return h
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if serveStatic(w, r) {
		return
	}

	ctx := h.p.Get().(*Context)
	defer h.p.Put(ctx)
	ctx.Config(w, r)

	controllerName, methodName := h.findControllerInfo(r)
	fmt.Println(controllerName, methodName)
	if strings.HasPrefix(methodName, "_") {
		switch methodName {
		case "_reload":
			fmt.Fprintln(ctx.w, "stop success")
			reload <- true
			return
		}
	}

	controllerT, ok := Mapping[controllerName[0]]
	if !ok {
		controllerT, ok = Mapping[controllerName[1]]
		if !ok {
			app := App{ctx, nil, -1, "controller not found1"}
			//http.NotFound(w, r)
			app.Response()
			return
		}
	}

	refV := reflect.New(controllerT)

	method := refV.MethodByName(methodName)
	if !method.IsValid() {
		//http.NotFound(w, r)
		tempMethodName := ""
		for i := 0; i < refV.Type().NumMethod(); i++ {
			tempMethodName = refV.Type().Method(i).Name
			if strings.ToLower(tempMethodName) == strings.ToLower(methodName) {
				method = refV.Method(i)
			}
		}
		if !method.IsValid() {
			app := App{ctx, nil, -1, "action not found2"}
			app.Response()
			return
		}
	}

	controller := refV.Interface().(IApp)
	controller.Init(ctx)
	method.Call(nil)
}

func RunOn(port string) {
	server := &http.Server{
		Handler: newHandler(),
		Addr:    port,
	}
	ServerCh <- server
}

func Stop(server *http.Server) {
	ret := <-StopCh
	fmt.Println("开始reload", ret)
	if err := server.Shutdown(context.Background()); err != nil {
		fmt.Println(err)
	}
	reload <- true
}

func Reload(server *http.Server) {
	for {
		<-reload
		fmt.Println("结束server")
		if err := server.Shutdown(context.Background()); err != nil {
			fmt.Println(err)
		}

		indexFile := os.Getenv("indexFile")
		cmd := exec.Command("go", "run", indexFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		fmt.Println("重启server: " + indexFile)

		err := cmd.Start()
		if err != nil {
			log.Fatalf("Restart: Failed to launch, error: %v", err)
		}
	}
	return
}
