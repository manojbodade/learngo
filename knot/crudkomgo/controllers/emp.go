package controllers

import (
	"fmt"
	"log"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	. "github.com/ekobudy/learngo/knot/crudkomgo/models"
	. "github.com/ekobudy/learngo/knot/crudkomgo/modules"
	"gopkg.in/mgo.v2/bson"
)

type EmpController struct{}

func (a *EmpController) Default(k *knot.WebContext) interface{} {
	FilterRequest(k)
	k.Config.OutputType = knot.OutputTemplate
	return ""
}

func (a *EmpController) Read(k *knot.WebContext) interface{} {
	FilterRequest(k)
	k.Config.OutputType = knot.OutputJson
	conn, err := PrepareConnection()
	if err != nil {
		return (toolkit.M{}).Set("status", "nok").Set("message", "System Error, cannot create connection")
	}
	defer conn.Close()
	cursor, err := conn.NewQuery().Select().From("emp").Cursor(nil)
	if err != nil {
		log.Println("Error fetching emp")
		return (toolkit.M{}).Set("status", "nok").Set("message", "System Error, cannot read EMP")
	}
	defer cursor.Close()
	if cursor == nil {
		log.Println("Cursor EMP is empty")
	}
	results := make([]Emp, 0)
	e := cursor.Fetch(&results, 0, false)
	if e != nil {
		log.Printf(e.Error())
		errMsg := fmt.Sprintf("{status:\"nok\",message: %v}", e.Error())
		return errMsg
	}
	return results
}

func (a *EmpController) Get(k *knot.WebContext) interface{} {
	FilterRequest(k)
	k.Config.OutputType = knot.OutputJson
	qId := struct {
		Id string
	}{}
	err := k.GetPayload(&qId)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Find rec=> ", qId)

	conn, err := PrepareConnection()
	if err != nil {
		return (toolkit.M{}).Set("status", "nok").Set("message", "System Error, cannot create connection")
	}
	defer conn.Close()
	cursor, err := conn.NewQuery().Select().From("emp").Where(dbox.Eq("_id", bson.ObjectIdHex(qId.Id))).Cursor(nil)
	if err != nil {
		log.Println("Error fetching emp")
		return (toolkit.M{}).Set("status", "nok").Set("message", "System Error, cannot read EMP")
	}
	defer cursor.Close()
	if cursor == nil {
		log.Println("Cursor EMP is empty")
	}
	results := make([]Emp, 0)
	e := cursor.Fetch(&results, 0, false)
	if e != nil {
		log.Printf(e.Error())
		errMsg := fmt.Sprintf("{status:\"nok\",message: %v}", e.Error())
		return errMsg
	}
	return results[0]
}
