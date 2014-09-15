package main

import (
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"gopkg.in/martini-contrib/render.v0"

	"code.google.com/p/go.crypto/bcrypt"
)

func Register(db *mgo.Database, nu NewUser, r render.Render) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(nu.Password), bcrypt.MinCost)
	if err != nil {
		r.Error(http.StatusInternalServerError)
		return
	}

	newUser := User{
		ID:           bson.NewObjectId(),
		Email:        nu.Email,
		Username:     nu.Username,
		Password:     string(hash),
		RegisteredAt: time.Now(),
	}

	if err := db.C("users").Insert(&newUser); err != nil {
		if mgo.IsDup(err) {
			r.Error(http.StatusConflict)
			return
		}
		r.Error(http.StatusInternalServerError)
		return
	}

	r.JSON(201, &newUser)
	return
}
