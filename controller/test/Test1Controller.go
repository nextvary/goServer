package test

import "fmt"
import "../../app"

type Test1Controller struct {
	app.App //匿名组合了App这个结构体，就是继承
}

func (i *Test1Controller) BanxiaAction() {
	fmt.Println(i.R().Method)
	i.Data["string"] = "test"
	fmt.Println(i)
	i.Response()
}
