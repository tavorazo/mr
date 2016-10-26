package base

import (
	"gopkg.in/mgo.v2"
)

var (
	HostDB 			string = "mongodb://mr-tooth:12qwaszx@ds044699.mlab.com:44699/mr"
	NameDB 			string = "mr"
	CollectionDB 	string = "usuarios"
)

func Connect() (*mgo.Session, error) {

	session, err := mgo.Dial(HostDB)
	if err != nil {
		return session, err
    }
    // defer session.Close()

    session.SetMode(mgo.Monotonic, true)

    return session, err
}