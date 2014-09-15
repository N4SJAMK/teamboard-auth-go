package main

import (
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"gopkg.in/martini-contrib/render.v0"

	"code.google.com/p/go.crypto/bcrypt"

	"github.com/dchest/uniuri"
)

func Login(db *mgo.Database, c Credentials, r render.Render) {

	// find the user with given 'email'

	var (
		user      = User{}
		userQuery = db.C("users").Find(bson.M{
			"email": c.Email,
		})
	)

	if err := userQuery.One(&user); err != nil {
		if err == mgo.ErrNotFound {
			r.Error(http.StatusNotFound)
			return
		}
		r.Error(http.StatusInternalServerError)
		return
	}

	// compare the given password to the found user's password

	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password), []byte(c.Password))
	if err != nil {
		r.Error(http.StatusUnauthorized)
		return
	}

	// find the token belonging to that user, if the token does not exist
	// create a new one for that user

	var (
		token      = Token{}
		tokenQuery = db.C("tokens").Find(bson.M{
			"user_id": user.ID,
		})
	)

	if err := tokenQuery.One(&token); err != nil {
		if err == mgo.ErrNotFound {
			token = Token{
				Secret:    uniuri.NewLen(16),
				UserID:    user.ID,
				CreatedAt: time.Now(),
			}
			if err := db.C("tokens").Insert(&token); err != nil {
				r.Error(http.StatusInternalServerError)
				return
			}
		} else {
			r.Error(http.StatusInternalServerError)
			return
		}
	}

	// set the 'x-access-token' header to the generated token
	// and return the user

	r.Header().Set("x-access-token", token.Secret)
	r.JSON(http.StatusOK, &user)
	return
}
