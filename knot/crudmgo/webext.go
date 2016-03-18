// webext
package webext

import (
	"os"

	. "github.com/ekobudy/learngo/knot/crudmgo/controllers"

	"github.com/eaciit/knot/knot.v1"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + "/../"
	}()
)

func init() {
	app := knot.NewApp("crudmgo")
	app.ViewsPath = wd + "views/"
	app.Controllers()
	app.Register(&NotesController{})
	//	app.Static("static", wd+"assets")
	app.LayoutTemplate = "_layout.html"
	knot.RegisterApp(app)
}
