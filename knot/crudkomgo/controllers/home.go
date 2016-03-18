package controllers

import (
	"github.com/eaciit/knot/knot.v1"
	. "github.com/ekobudy/learngo/knot/crudkomgo/modules"
)

type HomeController struct {
}

func (a *HomeController) Default(k *knot.WebContext) interface{} {
	FilterRequest(k)
	k.Config.OutputType = knot.OutputTemplate
	//	k.SetSession("username", "test")
	//	k.SetSession("lastlog", time.Now())
	//	log.Printf("Seesion %v \n", k.Session("username"))
	return ""
}
