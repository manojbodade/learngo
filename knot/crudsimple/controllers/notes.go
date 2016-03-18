// notes
package controllers

import (
	"log"
	"strconv"
	"time"

	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
)

//Sample of Note Struct
type Note struct {
	Id          string `bson:"_id"`
	Title       string
	Description string
	CreateOn    time.Time
}

//Note storage using map
var noteStorage = make(map[string]Note)
var id int = 0

type NotesController struct {
}

func (a *NotesController) Index(k *knot.WebContext) interface{} {
	log.Println("Notes Index Controller")
	k.Config.OutputType = knot.OutputTemplate
	//init empty notes with sample
	if len(noteStorage) == 0 {
		//		s1 := Note{"Title Note1", "Description Note1", time.Now()}
		//		id++

		//		noteStorage[strconv.Itoa(id)] = s1
		//		s1 = Note{"Title Note2", "Description Note2", time.Now()}
		//		id++
		//		noteStorage[strconv.Itoa(id)] = s1
	}

	return (toolkit.M{}).Set("Title", "Index Layout").Set("data", noteStorage)
}

func (a *NotesController) Save(k *knot.WebContext) interface{} {
	log.Println("Notes Save Function")
	note := Note{}
	e := k.GetForms(&note)
	if e != nil {
		log.Println("ERROR ", e)
	}
	note.CreateOn = time.Now()
	mode := k.Request.Form.Get("mode")
	if mode == "add" {
		id++
		log.Println("New Id ==> ", id)
		noteStorage[strconv.Itoa(id)] = note
		log.Println("Notestorage now ->> ", noteStorage)
	} else {
		key := k.Request.Form.Get("id")
		if noteExist, ok := noteStorage[key]; ok {
			note.CreateOn = noteExist.CreateOn
			delete(noteStorage, key)
			noteStorage[key] = note
		}
	}
	log.Println("SAVING MODE=>", mode, "  NOTE SAVED==> ", note)
	k.Config.OutputType = knot.OutputTemplate
	k.Config.ViewName = "notes/index.html"

	return (toolkit.M{}).Set("Title", "Index Layout").Set("data", noteStorage)
}
func (a *NotesController) Add(k *knot.WebContext) interface{} {

	log.Println("Notes ADD Controller")
	k.Config.OutputType = knot.OutputTemplate
	return (toolkit.M{}).Set("mode", "add")
}

func (a *NotesController) Edit(k *knot.WebContext) interface{} {
	log.Println("Notes EDIT Controller")
	k.Config.OutputType = knot.OutputTemplate
	k.Config.ViewName = "notes/add.html"
	id := k.Query("id")
	if noteExist, ok := noteStorage[id]; ok {
		return (toolkit.M{}).Set("mode", "edit").Set("data", noteExist).Set("id", id)
	} else {
		return (toolkit.M{}).Set("error", "No Data FOUND")
	}

	//	k.Config.OutputType = knot.OutputJson
	//	var data struct {
	//		Id string
	//	}
	//	e := k.GetPayload(&data)
	//	log.Println("KEY==>", data.Id)
	//	if e != nil {
	//		log.Println(e)
	//	}
	//	if noteExist, ok := noteStorage[data.Id]; ok {
	//		return (toolkit.M{}).Set("key", data.Id).Set("data", noteExist)
	//	} else {
	//		return (toolkit.M{}).Set("error", "No Data FOUND")
	//	}
}

func (a *NotesController) Delete(k *knot.WebContext) interface{} {
	log.Println("Notes Delete Controller")
	k.Config.OutputType = knot.OutputTemplate
	k.Config.ViewName = "notes/index.html"
	id := k.Query("id")
	delete(noteStorage, id)
	return (toolkit.M{}).Set("Title", "Index Layout").Set("data", noteStorage)
}
