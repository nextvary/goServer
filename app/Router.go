package app

import (
	"net/http"
	"reflect"
	"strings"
)

const (
	defController = "index.index"
	defMethod     = "IndexAction"
)

//静态资源
var Static map[string]string = make(map[string]string)

//路由注册
var Mapping map[string]reflect.Type = make(map[string]reflect.Type)

func router(pattern string, t reflect.Type) {
	Mapping[strings.ToLower(pattern)] = t
}

func Router(pattern string, app IApp) {
	refV := reflect.ValueOf(app)
	refT := reflect.Indirect(refV).Type()
	router(pattern, refT)
}

func AutoRouter(app IApp) {
	refV := reflect.ValueOf(app)
	refT := reflect.Indirect(refV).Type()
	refName := strings.TrimSuffix(strings.ToLower(refT.String()), "controller")
	router(refName, refT)
}

func (h *handler) findControllerInfo(r *http.Request) ([]string, string) {
	path := strings.Trim(r.URL.Path, "/")
	if strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}
	pathInfo := strings.Split(path, "/")
	length := len(pathInfo)
	controllerName := []string{"", ""}
	methodName := defMethod
	if length > 1 {
		methodName = strings.Title(strings.ToLower(pathInfo[length-1]))
		methodName += "Action"
	}
	if length == 3 {
		controllerName[0] = strings.ToLower(pathInfo[0]) + "." + strings.ToLower(pathInfo[1])
		controllerName[1] = strings.ToLower(pathInfo[0]) + "/" + strings.ToLower(pathInfo[1])
	}
	if strings.HasPrefix(pathInfo[0], "_") {
		methodName = strings.ToLower(pathInfo[0])
	}
	if pathInfo[0] == "" {
		controllerName[0] = defController
	}

	return controllerName, methodName
}

func serveStatic(w http.ResponseWriter, r *http.Request) bool {
	for prefix, static := range Static {
		if strings.HasPrefix(r.URL.Path, prefix) {
			file := static + r.URL.Path[len(prefix):]
			http.ServeFile(w, r, file)
			return true
		}
	}

	return false
}
