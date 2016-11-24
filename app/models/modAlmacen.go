package models

import "gopkg.in/mgo.v2/bson"

type Product struct {
	Name 			string			`json:"name"`
	Type 			string 			`json:"type"`
	Description 	string			`json:"description"`
	Quantity 		int 			`json:"quantity"`
	Min_quantity 	int				`json:"min_quantity"`
	Sale			bool			`json:"sale"`
	N_serial 		string 			`json:"n_serial"`
	Caterer			bson.ObjectId 	`json:"caterer" bson:"caterer,omitempty"`
	Deleted			int 			`json:"deleted"`
}

type Almacen struct{
	Id 				bson.ObjectId 	`json:"id" bson:"_id,omitempty"`
	Reference_id	string			`json:"reference_id"`  // A qué doctor/clinica está asignado el paciente
	Account_id		bson.ObjectId	`json:"account_id"`
	Products 		[]Product		`json:"patients"`
} 