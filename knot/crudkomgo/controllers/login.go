package controllers

import (
	"fmt"
	"log"
	"time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	. "github.com/ekobudy/learngo/knot/crudkomgo/models"
	. "github.com/ekobudy/learngo/knot/crudkomgo/modules"
)

type LoginController struct {
}

func (a *LoginController) Default(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputTemplate
	//	k.SetSession("username", "test")
	//	k.SetSession("lastlog", time.Now())
	//	log.Printf("Seesion %v \n", k.Session("username"))
	return ""
}
func (a *LoginController) Check(k *knot.WebContext) interface{} {
	FilterRequest(k)
	k.Config.OutputType = knot.OutputTemplate
	log.Println("current session, uname =>", k.Session("username"))
	return (toolkit.M{}).Set("cses", k.Session("username"))
}
func (a *LoginController) DoLogin(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	lgForm := struct {
		Uname  string
		Passwd string
	}{}
	er := k.GetPayload(&lgForm)
	if er != nil {
		log.Println(er.Error())
	}
	log.Printf("Login Forms => %v \n", lgForm)
	conn, er := PrepareConnection()
	if er != nil {
		return (toolkit.M{}).Set("status", "nok").Set("message", "System Error")
	}
	cursor1, er := conn.NewQuery().Select().From("sys_users").Where(dbox.Eq("username", lgForm.Uname)).Cursor(nil)
	if er != nil {
		log.Println("Error fetching Sys_Users")
		return (toolkit.M{}).Set("status", "nok").Set("message", "System Error")
	}
	if cursor1 == nil {
		return (toolkit.M{}).Set("status", "nok").Set("message", "Invalid Username/Password")
	} else {
		defer cursor1.Close()
		//		results := make([]map[string]interface{}, 0)
		results := make([]SysUser, 0)
		e := cursor1.Fetch(&results, 0, false)
		if e != nil {
			log.Printf(e.Error())
			errMsg := fmt.Sprintf("{status:\"nok\",message: %v}", e.Error())
			return errMsg
		}
		if len(results) > 0 {
			log.Printf("ResutsRow %v\n", results[0])
			if results[0].Password == toolkit.MD5String(lgForm.Passwd) {
				k.SetSession("username", results[0].Username)
				k.SetSession("lastlog", time.Now())
				return (toolkit.M{}).Set("status", "ok")
			} else {
				return (toolkit.M{}).Set("status", "nok").Set("message", "Invalid Username/Password")
			}
		} else {
			return (toolkit.M{}).Set("status", "nok").Set("message", "Invalid Username/Password")
		}
	}
}
