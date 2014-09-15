package main

import (
	"os"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"gopkg.in/martini.v1"

	"gopkg.in/martini-contrib/binding.v0"
	"gopkg.in/martini-contrib/render.v0"
)

type (
	Credentials struct {
		Email    string `json:"email"    binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	NewUser struct {
		Credentials
		Username string `json:"username"`
	}

	User struct {
		ID           bson.ObjectId `bson:"_id"           json:"id"`
		Email        string        `bson:"email"         json:"email"`
		Username     string        `bson:"username"      json:"username"`
		Password     string        `bson:"password"      json:"password"`
		RegisteredAt time.Time     `bson:"registered_at" json:"registeredAt"`
	}

	Token struct {
		Secret    string        `bson:"secret"`
		UserID    bson.ObjectId `bson:"user_id"`
		CreatedAt time.Time     `bson:"created_at"`
	}
)

const (
	DEFAULT_DB_URL  = "mongodb://localhost"
	DEFAULT_DB_NAME = "teamboard-dev-go"
)

func main() {
	var (
		dburl  = os.Getenv("MONGODB_URL")
		dbname = os.Getenv("MONGODB_NAME")
	)

	if len(dburl) == 0 {
		dburl = DEFAULT_DB_URL
	}

	if len(dbname) == 0 {
		dbname = DEFAULT_DB_NAME
	}

	var (
		m  = martini.Classic()
		db = NewMongoDB(dburl, dbname)
	)

	db.session.DB(db.name).C("tokens").EnsureIndex(mgo.Index{
		Key:         []string{"created_at"},
		ExpireAfter: time.Minute * 2,
	})

	db.session.DB(db.name).C("users").EnsureIndex(mgo.Index{
		Key:    []string{"email"},
		Unique: true,
	})

	m.Use(db.Session())
	m.Use(render.Renderer())

	m.Post("/login", binding.Bind(Credentials{}), Login)
	m.Post("/register", binding.Bind(NewUser{}), Register)

	m.Run()
}
