package base

import (
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	"encoding/json"

	"mr/app/models"
)

func NewCaterer(account_id, token string, jsonStr []byte) (string, int) {

	/* Función que recibe los valores de nickname como string, y como JSON del producto nuevo que se insertará en la BD */

	if CheckToken(token) == false {
		return "token no válido", 401   // Verifica que sea un token válido
	} else if UserExists("_id", account_id) == false{
		return "Usuario no encontrado", 403		//Verifica que el account_id exista en la base de datos
	}

	caterer := &models.Proveedor{}
	json.Unmarshal(jsonStr, caterer)

	session, err := mgo.Dial(HostDB)

	if err != nil {
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)
    con := session.DB(NameDB).C("proveedores")

    err = con.Insert(caterer)

	if err != nil {
		return "No se ha insertado", 500
	}

	return "Nuevo proveedor agregado", 201

}