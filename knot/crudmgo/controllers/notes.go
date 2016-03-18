// notes
package controllers

import (
	"log"
	//	"strconv"
	"time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	. "github.com/ekobudy/learngo/knot/crudmgo/modules"
	"gopkg.in/mgo.v2/bson"
)

type NotesController struct {
}
type Note struct {
	Id          bson.ObjectId `bson:"_id"`
	Title       string        `bson:"Title"`
	Description string        `bson:"Description"`
	CreateOn    time.Time     `bson:"CreateOn"`
}

func (a *NotesController) Index(k *knot.WebContext) interface{} {
	log.Println("Notes Index Controller")
	k.Config.OutputType = knot.OutputTemplate
	c, e := PrepareConnection()
	if e != nil {
		return e.Error()
	}
	defer c.Close()

	csr, e := c.NewQuery().Select().From("Notes").Cursor(nil)
	if e != nil {
		return e.Error()
	}
	if csr == nil {
		return "Cursor is nil"
	}
	defer csr.Close()
	results := make([]map[string]interface{}, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return e.Error()
	} else {
		log.Printf("Fetch N1 OK. Result: %v \n", results)
	}
	return (toolkit.M{}).Set("Title", "Index Layout").Set("data", results)
}

func (a *NotesController) Edit(k *knot.WebContext) interface{} {
	log.Println("Notes Edit Controller")
	//	k.Config.OutputType = knot.OutputJson
	k.Config.OutputType = knot.OutputTemplate
	k.Config.ViewName = "notes/add.html"
	id := k.Query("id")
	c, e := PrepareConnection()
	if e != nil {
		return e.Error()
	}
	defer c.Close()

	csr, e := c.NewQuery().Select().Where(dbox.Eq("_id", bson.ObjectIdHex(id))).From("Notes").Cursor(nil)
	//	csr, e := c.NewQuery().Select().Where(dbox.Eq("Title", "Learning MongoDB")).From("Notes").Cursor(nil)
	if e != nil {
		return e.Error()
	}
	if csr == nil {
		return "Cursor is nil"
	}
	defer csr.Close()
	results := make([]map[string]interface{}, 0)
	//	results := make([]Note, 0)
	e = csr.Fetch(&results, 10, false)
	if e != nil {
		log.Println("No Cursor Found")
		return e.Error()
	} else {
		log.Printf("Fetch N1 OK. Result: %v \n", results)
	}
	return (toolkit.M{}).Set("mode", "edit").Set("data", results[0]).Set("id", id)
}

func (a *NotesController) Save(k *knot.WebContext) interface{} {
	log.Println("Notes Save Function")
	note := Note{}
	e := k.GetForms(&note)
	c, e := PrepareConnection()
	if e != nil {
		return e.Error()
	}
	defer c.Close()
	if e != nil {
		log.Println("ERROR ", e)
	}
	mode := k.Request.Form.Get("mode")
	log.Printf("NOTES %v", note)
	if mode == "add" {
		note.Id = bson.NewObjectId()
		e := c.NewQuery().Insert().From("Notes").Save().Exec(toolkit.M{"data": note})
		if e != nil {
			log.Println("ERROR on save new record")
			return e.Error()
		}
	} else {
		note.Id = bson.ObjectIdHex(k.Request.Form.Get("id"))
		e := c.NewQuery().Update().From("Notes").Exec(toolkit.M{"data": note})
		if e != nil {
			log.Println("ERROR on update record")
			return e.Error()
		}
	}
	log.Println("SAVING MODE=>", mode, "  NOTE SAVED==> ", note)
	k.Config.OutputType = knot.OutputTemplate
	k.Config.ViewName = "notes/index.html"

	csr, e := c.NewQuery().Select().From("Notes").Cursor(nil)
	if e != nil {
		return e.Error()
	}
	if csr == nil {
		return "Cursor is nil"
	}
	defer csr.Close()
	results := make([]map[string]interface{}, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return e.Error()
	} else {
		log.Printf("Fetch N1 OK. Result: %v \n", results)
	}

	return (toolkit.M{}).Set("Title", "Index Layout").Set("data", results)
}

func (a *NotesController) Add(k *knot.WebContext) interface{} {
	log.Println("Notes ADD Controller")
	k.Config.OutputType = knot.OutputTemplate
	return (toolkit.M{}).Set("mode", "add")
}

func (a *NotesController) Delete(k *knot.WebContext) interface{} {
	log.Println("Notes Delete Controller")
	k.Config.OutputType = knot.OutputTemplate
	k.Config.ViewName = "notes/index.html"
	id := k.Query("id")

	c, e := PrepareConnection()
	if e != nil {
		return e.Error()
	}
	defer c.Close()
	data := Note{}
	data.Id = bson.ObjectIdHex(id)
	e = c.NewQuery().Delete().From("Notes").Exec(toolkit.M{"data": data})
	if e != nil {
		log.Println("ERROR on update record")
		return e.Error()
	}
	csr, e := c.NewQuery().Select().From("Notes").Cursor(nil)
	if e != nil {
		return e.Error()
	}
	if csr == nil {
		return "Cursor is nil"
	}
	defer csr.Close()
	results := make([]map[string]interface{}, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return e.Error()
	} else {
		log.Printf("Fetch N1 OK. Result: %v \n", results)
	}
	return (toolkit.M{}).Set("Title", "Index Layout").Set("data", results)
}
