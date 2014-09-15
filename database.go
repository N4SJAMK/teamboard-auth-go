package main

import "gopkg.in/mgo.v2"
import "gopkg.in/martini.v1"

type MongoDB struct {
	name    string
	session *mgo.Session
}

func NewMongoDB(url, name string) *MongoDB {
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	return &MongoDB{
		name:    name,
		session: session,
	}
}

func (mdb *MongoDB) Session() martini.Handler {
	return func(ctx martini.Context) {
		session := mdb.session.Clone()
		defer session.Close()

		ctx.Map(session.DB(mdb.name))

		ctx.Next()
		return
	}
}
