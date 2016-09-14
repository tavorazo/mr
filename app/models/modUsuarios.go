package models

import "gopkg.in/mgo.v2/bson"

type Usuario struct{
	Id 				bson.ObjectId 	`json:"id" bson:"_id,omitempty"`
	Firstname 		string			`json:"firstname"`
	Lastname		string			`json:"lastname"`
	Nickname 		string			`json:"nickname"`
	Age 			int				`json:"age"`
	Country 		string			`json:"country"`
	State 			string 			`json:"state"`
	Address			string			`json:"address"`
	Tel				string			`json:"tel"`
	Mail			string			`json:"mail"`
	Confirm_mail	string			`json:"confirm_mail"`
	Pass			string 			`json:"pass"`
	Confirm_pass	string			`json:"confirm_pass"`
	Clinica			bool			`json:"clinica"`
}