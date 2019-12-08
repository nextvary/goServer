package app

import (
	"encoding/json"
	"fmt"
	"github.com/ymzuiku/hit"
	"net/http"
	"time"
)

type IApp interface {
	Init(ctx *Context)
	W() http.ResponseWriter
	R() *http.Request
	Response()
	//Display(tpls ...string)
	//DisplayWithFuncs(funcs template.FuncMap, tpls ...string)
}

type App struct {
	ctx        *Context
	Data       map[string]interface{}
	StatusCode int
	Message    string
}

func (a *App) Init(ctx *Context) {
	a.ctx = ctx
	a.Data = make(map[string]interface{})
}

func (a *App) W() http.ResponseWriter {
	return a.ctx.w
}

func (a *App) R() *http.Request {
	return a.ctx.r
}

func (a *App) Response() {
	reponse := make(map[string]interface{})
	reponse["status_code"] = hit.If(a.StatusCode, a.StatusCode, 1)
	reponse["data"] = a.Data
	reponse["message"] = a.Message
	reponse["timestamps"] = time.Now().Unix()
	jsonData, _ := json.Marshal(reponse)
	//a.ctx.w.Header().Set("x-trace-id", strconv.Itoa(int(time.Now().Unix())))
	fmt.Fprint(a.ctx.w, string(jsonData))
}
