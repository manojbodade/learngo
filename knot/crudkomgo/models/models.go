package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type SysUser struct {
	Id        bson.ObjectId `bson:"_id"`
	Username  string
	Password  string
	LastLogin time.Time `bson:last_login`
}

type Emp struct {
	Id        bson.ObjectId `bson:"_id"`
	Nik       string
	FirstName string `bson:"first_name"`
	LastName  string `bson:"last_name"`
	Address   string
	City      string
	Province  string
	Skills    []string
}
