package base

import (
	"gopkg.in/mgo.v2/bson"
	"encoding/json"

	"mr/app/models"
)

func NewCaterer(account_id, token string, jsonStr []byte) (string, int) {

	/* Función que recibe los valores de nickname como string, y como JSON del proveedor nuevo que se insertará en la BD */

	session, err := Connect() // Conecta a la base de datos
	if err != nil {
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

	if CheckToken(token, session) == false {
		return "token no válido", 401   // Verifica que sea un token válido
	} else if UserExists("_id", account_id) == false{
		return "Usuario no encontrado", 403		//Verifica que el account_id exista en la base de datos
	}

	caterer := &models.Proveedor{}
	json.Unmarshal(jsonStr, caterer)

	if CatererExists(caterer.Name){
		return "El nombre del proveedor ya existe", 409
	}

    con := session.DB(NameDB).C("proveedores")

    err = con.Insert(caterer)

	if err != nil {
		return "No se ha insertado", 500
	}

	return "Nuevo proveedor agregado", 201

}

func CatererExists(name string) bool {

	/* Función que verifica si existe el número de serial del proveedor en la base de datos */

	session, err := Connect()
	if err != nil {
		return false
    }
    defer session.Close()
    con := session.DB(NameDB).C("proveedores")

    result := &models.Proveedor{}
    err = con.Find(bson.M{"name": name}).One(&result)

    if err != nil{
    	return false
    } else{
    	return true
    }
}

func UpdateCaterer(account_id, token string, jsonStr []byte) (string, int){

	/* Función que actualiza un proveedor en la base de datos 
		Se reciben el id de usuario y el id de la BD del proveedkr */

	session, err := Connect()
	if err != nil {
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

	if CheckToken(token, session) == false {
		return "token no válido", 401   // Verifica que sea un token válido
	} else if UserExists("_id", account_id) == false{
		return "Usuario no encontrado", 403		//Verifica que el account_id exista en la base de datos
	}

	caterer := &models.Proveedor{}
	json.Unmarshal(jsonStr, caterer)

    con := session.DB(NameDB).C("proveedores")

    colQuerier := bson.M{"_id": caterer.Id }  // Busca el documento por ACCOUNT_ID
	change := bson.M{"$set": caterer } // Inserta en el array de proveedores
	err = con.Update(colQuerier, change)

	if err != nil {

		return "Proveedor no encontrado", 400
	}

	return "Datos de proveedor cambiados", 200

}

