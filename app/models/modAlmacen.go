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
}