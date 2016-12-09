package base

import (
	"gopkg.in/mgo.v2"
)

var (
	HostDB 			string = "mongodb://mr-tooth:12qwaszx@ds044699.mlab.com:44699/mr"
	NameDB 			string = "mr"
	err 			error
	session 		*mgo.Session
	col 			*mgo.Collection
)

func Connect() bool {

	session, err = mgo.Dial(HostDB)
	if err != nil {
		return false
    }
    // defer session.Close()

    session.SetMode(mgo.Monotonic, true)

    return true
}