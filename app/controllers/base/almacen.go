package base

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"encoding/json"
	"time"

	"mr/app/models"
)

func NewProduct(account_id, token string, jsonStr []byte) (string, int) {

	/* Función que recibe los valores de nickname como string, y como JSON del producto nuevo que se insertará en la BD */

	session, err := Connect() // Conecta a la base de datos
	if err != nil {
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

	if CheckToken(token, session) == false {
		return "token no válido", 401   // Verifica que sea un token válido
	} else if UserExists("_id", account_id, session) == false{
		return "Usuario no encontrado", 403		//Verifica que el account_id exista en la base de datos
	}

	productVals := &models.Product{}
	json.Unmarshal(jsonStr, productVals)
	productVals.Deleted = 0 	// False indica que el producto no ha sido borrado del almacén

	if ProductExists(productVals.N_serial , account_id, session) == true {
		return "El número de serial del producto ya existe",400
	}

    con := session.DB(NameDB).C(CollectionDB)

    colQuerier := bson.M{"_id": bson.ObjectIdHex(account_id)}  // Busca el documento por nickname
	change := bson.M{"$push": bson.M{"products": productVals} } // Inserta en el array de productos
	err = con.Update(colQuerier, change)

	if err != nil {		
		return "Usuario no encontrado", 400
	}

	return "Producto agregado en el almacén", 201

}

func ProductExists(serial, account_id string, session *mgo.Session) bool {

	/* Función que verifica si existe el número de serial del producto en la base de datos */

    con := session.DB(NameDB).C(CollectionDB)

    type Result struct{
    	Products []models.Product `json:"products"`
    }

    result := Result{}
    err := con.Find(bson.M{"_id": bson.ObjectIdHex(account_id), "products.n_serial": serial}).Select(bson.M{"products.n_serial": 1, "_id": 0}).One(&result)

    if err != nil{
    	return false
    } else{
    	return true
    }
}

func UpdateProductAmount(account_id, n_serial string, token string, jsonStr []byte) (string, int) {

	/* Función que recibe los valores de ACCOUNT_ID como string, y como JSON la nueva cantidad del producto con n_serial */

	session, err := Connect() // Conecta a la base de datos
	if err != nil {
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

	if CheckToken(token, session) == false {
		return "token no válido", 401   // Verifica que sea un token válido
	} else if UserExists("_id", account_id, session) == false{
		return "Usuario no encontrado", 403		//Verifica que el account_id exista en la base de datos
	}

	productVals := &models.Product{}
	json.Unmarshal(jsonStr, productVals)

    con := session.DB(NameDB).C(CollectionDB)

    colQuerier := bson.M{"_id": bson.ObjectIdHex(account_id), "products.n_serial": n_serial}  // Busca el documento por ACCOUNT_ID
	change := bson.M{"$set": bson.M{"products.$.quantity": productVals.Quantity} } // Inserta en el array de productos
	err = con.Update(colQuerier, change)

	if err != nil {		
		return "Producto no encontrado", 400
	}

	return "Cantidad de productos actualizada", 200

}

func UpdateProduct(account_id, n_serial string, token string, jsonStr []byte) (string, int){

	/* Función que actualiza un producto para un usuario en la base de datos 
		Se reciben el id de usuario y el número de serie unico del producto */

	session, err := Connect() // Conecta a la base de datos
	if err != nil {
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

	if CheckToken(token, session) == false {
		return "token no válido", 401   // Verifica que sea un token válido
	} else if UserExists("_id", account_id, session) == false{
		return "Usuario no encontrado", 403		//Verifica que el account_id exista en la base de datos
	}

	productVals := &models.Product{}
	json.Unmarshal(jsonStr, productVals)

	if ProductExists(productVals.N_serial , account_id, session) == true {
		return "El número de serial del producto ya existe",400
	}

    con := session.DB(NameDB).C(CollectionDB)

    productVals.Deleted = 0  // Previene que se pueda actualizar el punto de restauración
    colQuerier := bson.M{"_id": bson.ObjectIdHex(account_id), "products.n_serial": n_serial}  // Busca el documento por ACCOUNT_ID
	change := bson.M{"$set": bson.M{"products.$": productVals} } // Inserta en el array de productos
	err = con.Update(colQuerier, change)

	if err != nil {		
		return "Producto no encontrado", 400
	}

	return "Datos de producto actualizados", 200

}

func EraseProduct(account_id, n_serial string, token string) (string, int){

	/* Función que elmina un producto del usuario recibido de la base de datos
		se recibe el id de usuario y el numero de serie unico del producto */

	session, err := Connect() // Conecta a la base de datos
	if err != nil {
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

	if CheckToken(token, session) == false {
		return "token no válido", 401   // Verifica que sea un token válido
	} else if UserExists("_id", account_id, session) == false{
		return "Usuario no encontrado", 403		//Verifica que el account_id exista en la base de datos
	}

    con := session.DB(NameDB).C(CollectionDB)

    colQuerier := bson.M{"_id": bson.ObjectIdHex(account_id)}  // Busca el documento por ACCOUNT_ID
	change := bson.M{"$pull": bson.M{"products": bson.M{"n_serial":n_serial } } } // Elimina en el array de productos en base al número de serial
	err = con.Update(colQuerier, change)

	if err != nil {		
		return "Producto no encontrado", 400
	}

	return "Producto eliminado totalmente de la base de datos", 200

}

func SaveDeletedProduct(account_id, n_serial string, token string) (string, int){

	/* Función que actualiza un producto del usuario recibido de la base de datos para que se permita restaurar antes de ser eliminado
		se recibe el id de usuario y el numero de serie unico del producto  */

	session, err := Connect() // Conecta a la base de datos
	if err != nil {
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

	if CheckToken(token, session) == false {
		return "token no válido", 401   // Verifica que sea un token válido
	} else if UserExists("_id", account_id, session) == false{
		return "Usuario no encontrado", 403		//Verifica que el account_id exista en la base de datos
	}

    con := session.DB(NameDB).C(CollectionDB)

    colQuerier := bson.M{"_id": bson.ObjectIdHex(account_id) , "products.n_serial": n_serial }  // Busca el documento por ACCOUNT_ID
	change := bson.M{"$set": bson.M{"products.$.deleted": int(time.Now().Unix())  } } // El campo deleted se actualiza con el tiempo unix actual
	err = con.Update(colQuerier, change)

	if err != nil {		
		return "Producto no encontrado", 400
	}

	return "Producto eliminado", 200

}

func GetProducts(all bool, account_id string, token, n_serial string) (string, int, interface{}) {

	/* 	Función que busca en la base de datos uno o más productos
		"all" es un valor booleano que indica si se quieren todos los productos o uno en específico
		Recibe también los valores de id de la cuenta y el numero de serial del producto en caso de que se requiera
	*/

	data := make(map[string]interface{})

	session, err := Connect() // Conecta a la base de datos
	if err != nil {
		return "No se ha conectado a la base de datos", 500, data
    }
    defer session.Close()

	if CheckToken(token, session) == false {
		return "token no válido", 401, data   // Verifica que sea un token válido
	} else if UserExists("_id", account_id, session) == false {
		return "Usuario no encontrado", 403, data		//Verifica que el account_id exista en la base de datos
	}

    con := session.DB(NameDB).C(CollectionDB)

    type Result struct{
    	Id 			bson.ObjectId 		`json:"id" bson:"_id,omitempty"`
    	Products 	[]models.Product 	`json:"products"`
    }

    result := Result{}

    if all == false {  // Si está desactivada la opción de todos los productos buscará uno en específico de acuerdo al n_serial indicado
    	err = con.Find(bson.M{"_id": bson.ObjectIdHex(account_id)}).Select(bson.M{"products": bson.M{"$elemMatch": bson.M{"n_serial":n_serial, "deleted": 0} }, "_id":0 }).One(&result)
    } else{
    	filter := bson.M{"$filter": bson.M{"input": "$products", "as": "product", "cond":bson.M{"$eq": []interface{}{"$$product.deleted", 0} } }}  // Aggregation query para MGO
    	err = con.Pipe([]bson.M{{"$match":bson.M{"_id": bson.ObjectIdHex(account_id)}}, {"$project": bson.M{"products": filter }} }).One(&result)
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