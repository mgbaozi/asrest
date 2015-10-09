package rest

import (
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var session *mgo.Session

func Connect(url string) {
	var err error
	session, err = mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
}

func DB(name string) *mgo.Database {
	return session.DB(name)
}

func toObjectId(id interface{}) (object_id bson.ObjectId, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("objectid is not available")
		}
	}()
	var ok bool
	if object_id, ok = id.(bson.ObjectId); !ok {
		object_id = bson.ObjectIdHex(id.(string))
	}
	return
}
