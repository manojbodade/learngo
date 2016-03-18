package controllers

import (
	"net/http"

	"github.com/eaciit/knot/knot.v1"
)

type LogoutController struct {
}

func (a *LogoutController) Default(k *knot.WebContext) interface{} {
	k.SetSession("username", nil)
	k.SetSession("lastlog", nil)
	http.Redirect(k.Writer, k.Request, "/home/default", http.StatusTemporaryRedirect)
	return ""
}
