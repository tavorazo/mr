package models

import "gopkg.in/mgo.v2/bson"

type Proveedor struct {
	Id 			bson.ObjectId 	`json:"id" bson:"_id,omitempty"`
	Name 		string 			`json:"name"`
	Agent_name 	string			`json:"agent_name"`
	Mail 		string			`json:"mail"`
	Address		string			`json:"address"`
	Phone	 	string 			`json:"phone"`
	Cellphone	string 			`json:"cellphone"`
}