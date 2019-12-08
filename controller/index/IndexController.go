package index

import (
	"../../app"
	"fmt"
)

type IndexController struct {
	app.App //匿名组合了App这个结构体，就是继承
}

func (i *IndexController) IndexAction() {
	i.Data["message"] = "server start "
	//i.StatusCode=-1
	//i.Message="error"
	fmt.Println(i)
	i.Response()
}
func (i *IndexController) Index1Action() {
	i.Data["name"] = "333"
	i.Data["email"] = "333@gmail123.com"
	i.Data["num"] = 333
	i.Data["string"] = "333"
	i.Response()
}
