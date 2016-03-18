package modules

import (
	"log"
	"net/http"

	"github.com/eaciit/knot/knot.v1"
)

func FilterRequest(k *knot.WebContext) {
	if k.Session("username") == nil {
		log.Println("invalid session , redirecting to login/default")
		http.Redirect(k.Writer, k.Request, "/login/default", http.StatusTemporaryRedirect)
	}
	return
}
