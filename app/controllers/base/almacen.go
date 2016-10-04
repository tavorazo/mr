package base

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"

	"mr/app/models"
)

func NewProduct(account_id string, jsonStr []byte) (string, int) {

	/* Función que recibe los valores de nickname como string, y como JSON del producto nuevo que se insertará en la BD */

	productVals := &models.Product{}
	json.Unmarshal(jsonStr, productVals)

	if ProductExists(productVals.N_serial , account_id) == true {
		return "El número de serial del producto ya existe",400
	}

	session, err := mgo.Dial(HostDB)

	if err != nil {
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)
    con := session.DB(NameDB).C(CollectionDB)

    colQuerier := bson.M{"_id": bson.ObjectIdHex(account_id)}  // Busca el documento por nickname
	change := bson.M{"$push": bson.M{"products": productVals} } // Inserta en el array de productos
	err = con.Update(colQuerier, change)

	if err != nil {		
		return "Usuario no encontrado", 401
	}

	return "Producto agregado en el almacén", 200

}

func ProductExists(serial, account_id string) bool {

	/* Función que verifica si existe el número de serial del producto en la base de datos */

	session, err := mgo.Dial(HostDB)
	if err != nil {
		return false
    }
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)
    con := session.DB(NameDB).C(CollectionDB)

    type Result struct{
    	Products []models.Product `json:"products"`
    }

    result := Result{}
    err = con.Find(bson.M{"_id": bson.ObjectIdHex(account_id), "products.n_seria": serial}).Select(bson.M{"products.n_serial": 1, "_id": 0}).One(&result)

    if err != nil{
    	return false
    } else{
    	return true
    }
}

func UpdateProductAmount(account_id, n_serial string, jsonStr []byte) (string, int) {

	/* Función que recibe los valores de nickname como string, y como JSON del producto nuevo que se insertará en la BD */

	productVals := &models.Product{}
	json.Unmarshal(jsonStr, productVals)

	session, err := mgo.Dial(HostDB)

	if err != nil {
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)
    con := session.DB(NameDB).C(CollectionDB)

    colQuerier := bson.M{"_id": bson.ObjectIdHex(account_id), "products.n_serial": n_serial}  // Busca el documento por nickname
	change := bson.M{"$set": bson.M{"products.$.quantity": productVals.Quantity} } // Inserta en el array de productos
	err = con.Update(colQuerier, change)

	if err != nil {		
		return "Producto no encontrado", 401
	}

	return "Cantidad de productos actualizada", 200

}