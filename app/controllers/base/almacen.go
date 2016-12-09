package base

import (
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"time"

	"mr/app/models"
)

func NewProduct(account_id, reference_id string, token string, jsonStr []byte) (string, int) {

	/* Función que recibe los valores de nickname como string, y como JSON del producto nuevo que se insertará en la BD */

	if Connect() == false { // Conecta a la base de datos
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

	if CheckToken(token) == false {
		return "token no válido", 401   // Verifica que sea un token válido
	} else if UserExists("_id", account_id) == false{
		return "Usuario no encontrado", 403		//Verifica que el account_id exista en la base de datos
	}

	productVals := &models.Product{}
	json.Unmarshal(jsonStr, productVals)
	productVals.Deleted = 0 	// False indica que el producto no ha sido borrado del almacén


    col = session.DB(NameDB).C("almacenes")

    if StoreExists(reference_id, account_id) == false { // Si no existe el almacén en la base de datos crea uno nuevo

    	newReference := &models.Almacen{}
    	newReference.Reference_id = reference_id
    	newReference.Account_id = bson.ObjectIdHex(account_id)
    	err = col.Insert(newReference)

    }

    if ProductExists(productVals.N_serial, account_id, reference_id) == true {
		return "El número de serial del producto ya existe",400
	}

    colQuerier := bson.M{"account_id": bson.ObjectIdHex(account_id), "reference_id": reference_id}  // Busca el documento por account_id y referencia
	change := bson.M{"$push": bson.M{"products": productVals} } // Inserta en el array de productos
	err = col.Update(colQuerier, change)

	if err != nil {		
		return "Usuario no encontrado", 400
	}

	return "Producto agregado en el almacén", 201

}

func ProductExists(serial string, account_id, reference_id string) bool {

	/* Función que verifica si existe el número de serial del producto en la base de datos */

    col = session.DB(NameDB).C("almacenes")

    type Result struct{
    	Products []models.Product `json:"products"`
    }

    result := Result{}
    err = col.Find(bson.M{"account_id": bson.ObjectIdHex(account_id), "reference_id": reference_id, "products.n_serial": serial}).Select(bson.M{"products.n_serial": 1, "_id": 0}).One(&result)

    if err != nil{
    	return false
    } else{
    	return true
    }
}

func UpdateProductAmount(account_id, reference_id string, n_serial, token string, jsonStr []byte) (string, int) {

	/* Función que recibe los valores de ACCOUNT_ID como string, y como JSON la nueva cantidad del producto con n_serial */

	if Connect() == false { // Conecta a la base de datos
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

	if CheckToken(token) == false {
		return "token no válido", 401   // Verifica que sea un token válido
	} else if UserExists("_id", account_id) == false{
		return "Usuario no encontrado", 403		//Verifica que el account_id exista en la base de datos
	}

	productVals := &models.Product{}
	json.Unmarshal(jsonStr, productVals)

    col = session.DB(NameDB).C("almacenes")

    colQuerier := bson.M{"account_id": bson.ObjectIdHex(account_id), "reference_id": reference_id, "products.n_serial": n_serial}  // Busca el documento por ACCOUNT_ID
	change := bson.M{"$set": bson.M{"products.$.quantity": productVals.Quantity} } // Inserta en el array de productos
	err = col.Update(colQuerier, change)

	if err != nil {		
		return "Producto no encontrado", 400
	}

	return "Cantidad de productos actualizada", 200

}

func UpdateProduct(account_id, reference_id string, n_serial, token string, jsonStr []byte) (string, int){

	/* Función que actualiza un producto para un usuario en la base de datos 
		Se reciben el id de usuario y el número de serie unico del producto */

	if Connect() == false { // Conecta a la base de datos
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

	if CheckToken(token) == false {
		return "token no válido", 401   // Verifica que sea un token válido
	} else if UserExists("_id", account_id) == false{
		return "Usuario no encontrado", 403		//Verifica que el account_id exista en la base de datos
	}

	productVals := &models.Product{}
	json.Unmarshal(jsonStr, productVals)

	productVals.N_serial = n_serial

    col = session.DB(NameDB).C("almacenes")

    productVals.Deleted = 0  // Previene que se pueda actualizar el punto de restauración
    colQuerier := bson.M{"account_id": bson.ObjectIdHex(account_id), "reference_id": reference_id, "products.n_serial": n_serial}  // Busca el documento por ACCOUNT_ID
	change := bson.M{"$set": bson.M{"products.$": productVals} } // Inserta en el array de productos
	err = col.Update(colQuerier, change)

	if err != nil {		
		return "Producto no encontrado", 400
	}

	return "Datos de producto actualizados", 200

}

func EraseProduct(account_id, reference_id string, n_serial, token string) (string, int){

	/* Función que elimina un producto del usuario recibido de la base de datos
		se recibe el id de usuario y el numero de serie unico del producto */

	if Connect() == false { // Conecta a la base de datos
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

	if CheckToken(token) == false {
		return "token no válido", 401   // Verifica que sea un token válido
	} else if UserExists("_id", account_id) == false{
		return "Usuario no encontrado", 403		//Verifica que el account_id exista en la base de datos
	}

    col = session.DB(NameDB).C("almacenes")

    colQuerier := bson.M{"account_id": bson.ObjectIdHex(account_id), "reference_id": reference_id}  // Busca el documento por ACCOUNT_ID
	change := bson.M{"$pull": bson.M{"products": bson.M{"n_serial":n_serial } } } // Elimina en el array de productos en base al número de serial
	err = col.Update(colQuerier, change)

	if err != nil {		
		return "Producto no encontrado", 400
	}

	return "Producto eliminado totalmente de la base de datos", 200

}

func SaveDeletedProduct(account_id, reference_id string, n_serial, token string) (string, int){

	/* Función que actualiza un producto del usuario recibido de la base de datos para que se permita restaurar antes de ser eliminado
		se recibe el id de usuario y el numero de serie unico del producto  */

	if Connect() == false { // Conecta a la base de datos
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

	if CheckToken(token) == false {
		return "token no válido", 401   // Verifica que sea un token válido
	} else if UserExists("_id", account_id) == false{
		return "Usuario no encontrado", 403		//Verifica que el account_id exista en la base de datos
	}

    col = session.DB(NameDB).C("almacenes")

    colQuerier := bson.M{"account_id": bson.ObjectIdHex(account_id),  "reference_id": reference_id, "products.n_serial": n_serial }  // Busca el documento por ACCOUNT_ID
	change := bson.M{"$set": bson.M{"products.$.deleted": int(time.Now().Unix())  } } // El campo deleted se actualiza con el tiempo unix actual
	err = col.Update(colQuerier, change)

	if err != nil {		
		return "Producto no encontrado", 400
	}

	return "Producto eliminado", 200

}

func GetProducts(all bool, account_id, reference_id string, token, n_serial string) (string, int, interface{}) {

	/* 	Función que busca en la base de datos uno o más productos
		"all" es un valor booleano que indica si se quieren todos los productos o uno en específico
		Recibe también los valores de id de la cuenta y el numero de serial del producto en caso de que se requiera
	*/

	data := make(map[string]interface{})

	if Connect() == false { // Conecta a la base de datos
		return "No se ha conectado a la base de datos", 500, data
    }
    defer session.Close()

	if CheckToken(token) == false {
		return "token no válido", 401, data   // Verifica que sea un token válido
	} else if UserExists("_id", account_id) == false {
		return "Usuario no encontrado", 403, data		//Verifica que el account_id exista en la base de datos
	}

    col = session.DB(NameDB).C("almacenes")

    // type Result struct{
    // 	Id 			bson.ObjectId 		`json:"id" bson:"_id,omitempty"`
    // 	Products 	[]models.Product 	`json:"products"`
    // }

    result := models.Almacen{}

    if all == false {  // Si está desactivada la opción de todos los productos buscará uno en específico de acuerdo al n_serial indicado
    	err = col.Find(bson.M{"account_id": bson.ObjectIdHex(account_id), "reference_id": reference_id}).Select(bson.M{"products": bson.M{"$elemMatch": bson.M{"n_serial":n_serial, "deleted": 0} }, "_id":0 }).One(&result)
    } else{
    	filter := bson.M{"$filter": bson.M{"input": "$products", "as": "product", "cond":bson.M{"$eq": []interface{}{"$$product.deleted", 0} } }}  // Aggregation query para MGO
    	err = col.Pipe([]bson.M{{"$match":bson.M{"account_id": bson.ObjectIdHex(account_id), "reference_id": reference_id}}, {"$project": bson.M{"products": filter }} }).One(&result)
    }
    
    if err != nil  {
    	return "No se encontró el producto", 400, data
    } 	

	data["producto"] = result.Products
	productsFind := len(result.Products) // Cantidad de productos encontrados

	if productsFind == 0 {
		if all == true {
			return "No hay productos en el inventario", 206, data["producto"]  // Si realiza una búsqueda de todos los productos retorna exito pero con el array vacío
		} else {
			return "Producto no encontrado", 400, data["producto"]				// Si es una búsqueda de un solo producto retorna error al no ser encontrado
		}
	} else if productsFind == 1 {
		return "Producto encontrado", 200, data["producto"]			// Si hay un solo producto
	} else {
		return "Productos encontrados", 200, data["producto"]		// Si hay dos o más productos
	}
    
    
}

func StoreExists(reference_id, account_id string) bool {

	/* Función que verifica si existe un documento de almacén para el id de referencia indicado */

	col = session.DB(NameDB).C("almacenes")

    result := models.PatientsExt{}
    err := col.Find(bson.M{"reference_id": reference_id, "account_id": bson.ObjectIdHex(account_id)}).One(&result)

    if err != nil{
    	return false
    } else{
    	return true
    }

}