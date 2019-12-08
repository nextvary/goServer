package index

import (
	"../../app"
	"fmt"
)

type TestController struct {
	app.App //匿名组合了App这个结构体，就是继承
}

func (i *TestController) TestAction() {
	fmt.Println(i.R().Method)
	i.Data["name"] = "111"
	i.Data["email"] = "111@gmail123.com"
	i.Data["num"] = 111
	i.Data["string"] = "111"
	//i.StatusCode=-1
	//i.Message="error"
	fmt.Println(i)
	i.Response()
}
