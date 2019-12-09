package crontab

import "../../app"

type CronController struct {
	app.App //匿名组合了App这个结构体，就是继承
}

func (c *CronController) MysqlTestAction() {
	c.Data["sql"] = "select * from ..."
	c.Response()
}
