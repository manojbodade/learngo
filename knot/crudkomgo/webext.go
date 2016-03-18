// webext
package webext

import (
	"log"
	"os"
	"time"
	//	"strings"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	. "github.com/ekobudy/learngo/knot/crudkomgo/controllers"
	. "github.com/ekobudy/learngo/knot/crudkomgo/models"
	. "github.com/ekobudy/learngo/knot/crudkomgo/modules"
	"gopkg.in/mgo.v2/bson"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + "/../"
	}()
)

func init() {
	app := knot.NewApp("crudwithko")
	initUser()
	app.ViewsPath = wd + "views/"
	app.Controllers()
	app.Register(&HomeController{})
	app.Register(&LoginController{})
	app.Register(&EmpController{})
	app.Register(&LogoutController{})
	app.Static("static", wd+"assets")
	app.LayoutTemplate = "_layout.html"
	knot.RegisterApp(app)
}

func initUser() {
	log.Println("Init Database")
	conn, err := PrepareConnection()
	if err != nil {
		log.Println(err.Error())
		//		return err.Error()
	}
	defer conn.Close()

	csr, e := conn.NewQuery().Select().From("sys_users").Where(dbox.Eq("username", "admin")).Cursor(nil)
	if e != nil {
		log.Println("Error fetching Sys_Users")
		//		return e.Error()
	}
	if csr == nil {
		log.Println("Cursor is null")
	} else {
		defer csr.Close()
		results := make([]map[string]interface{}, 0)
		e = csr.Fetch(&results, 0, false)
		if e != nil {
			log.Printf(e.Error())
			//			return e.Error()
		} else {
			//			log.Printf("Fetch N1 OK. Result: %v \n", results)
			log.Println("Length of Results ", len(results))
			if len(results) == 0 {
				usr := SysUser{}
				usr.Id = bson.NewObjectId()
				usr.Username = "admin"
				usr.Password = toolkit.MD5String("password")
				usr.LastLogin = time.Now()
				log.Printf("init user=> %v\n", usr)
				e := conn.NewQuery().Insert().From("sys_users").Save().Exec(toolkit.M{"data": usr})
				if e != nil {
					log.Println("ERROR on save new record")
					//			return e.Error()
				}
			}
		}

	}
}
