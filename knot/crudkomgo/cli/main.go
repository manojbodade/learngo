// main
package main

import (
	"log"

	"github.com/eaciit/knot/knot.v1"
	_ "github.com/ekobudy/learngo/knot/crudkomgo"
)

func main() {
	app := knot.GetApp("crudwithko")
	if app == nil {
		log.Println("APP is NULL")
		return
	}
	knot.StartApp(app, "localhost:9999")
}
