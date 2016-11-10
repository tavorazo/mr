package models

import "gopkg.in/mgo.v2/bson"

type Patient struct{
	Id 		 		string 	    	`json:"id" bson:"_id,omitempty"`
	Name 			string			`json:"name"`
	Address			string			`json:"address"`
	Telephone		string			`json:"telephone"`
	Cellphone		string 			`json:"cellphone"`
	Mail			string			`json:"mail"`
	Birthdate		string			`json:"birthdate"`
	Gender	 		string			`json:"gender"`
	Comments 		string			`json:"comments"`
	Address_aux		string 			`json:"address_aux"`
	Telephone_aux	string			`json:"telephone_aux"`
	Name_aux		string			`json:"name_aux"`
}

type PatientsExt struct{
	Id 				bson.ObjectId 	`json:"id" bson:"_id,omitempty"`
	Reference_id	string			`json:"reference_id"`
	Patients 		[]Patient		`json:"patients"`
} 