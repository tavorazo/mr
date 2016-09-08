package base

import (
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Persona struct{
	Id 			bson.ObjectId 	`json:"id" bson:"_id,omitempty"`
	Usuario 	string			`json:"usuario"`
	Nombre 		string			`json:"nombre"`
	Apellido 	string			`json:"apellido"`
	Edad 		string			`json:"edad"`
	Contrasena 	string			`json:"contrasena"`
	Correo 		string 			`json:"correo"`
}

func Busca() Persona{

	session, err := mgo.Dial("clinic.westus.cloudapp.azure.com") //Host del servidor mongo
	if err != nil{
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	con := session.DB("mr-tooth").C("usuarios")

	result := Persona{}
	err = con.Find(bson.M{"usuario": "mrtooth"}).One(&result) // Busca un nombre en la colección y lo almacena en result

	if err != nil {
		log.Fatal(err)
	}
	
	return result     //Regresa la estructura con los valores
}