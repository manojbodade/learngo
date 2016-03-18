// main.go
package main

import (
	//	"fmt"
	"log"

	"github.com/eaciit/knot/knot.v1"
)

func main() {
	knot.DefaultOutputType = knot.OutputHtml
	ks := new(knot.Server)
	ks.Route("notes/", NotesHandler)
	//		ks.Route("notes/{id}", NotesHandler)
	ks.Listen()
}
func NotesHandler(r *knot.WebContext) interface{} {
	id := r.Query("id")
	if id != "" {
		log.Println("ID is not null =>> ", id)
	}
	switch r.Request.Method {
	case "GET":
		return "GET Method called"
	case "POST":
		return "POST Method called"
	case "PUT":
		return "PUT Method called"
	case "DELETE":
		return "DELETE Method called"
	default:
		return "Unknown Method"
	}
}
