package base

import (
	"gopkg.in/mgo.v2/bson"
	"encoding/json"

	"mr/app/models"
	"mr/app/controllers/mailing"
)

func NewUser(jsonStr []byte) (string, int) {

	/*	Función que recibe el valor en JSON y lo inserta en la Base de datos
		Devuelve Un mensaje, el estátus del servidor y el error si existe */

	if Connect() == false { // Conecta a la base de datos
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

	usr := &models.Usuario{}
	json.Unmarshal(jsonStr, usr)			// Recibe el valor json y lo almacena en la estructura

	if(usr.Pass != usr.Confirm_pass) {
		return "Las contraseñas no coinciden", 409
	} else if(usr.Mail != usr.Confirm_mail) {
		return "Los correos no coinciden", 409
	} else if UserExists("mail",usr.Mail) == true {
		return "Correo ya existente en la base de datos", 409
	} else if UserExists("nickname",usr.Nickname,) == true {
		return "Nombre de usuario ya existente en la base de datos", 409
	}

	usr.Pass = EncryptToString(usr.Pass)	// Encripta la contraseña
	usr.Confirm_pass = EncryptToString(usr.Confirm_pass)

	col = session.DB(NameDB).C("usuarios")
	err = col.Insert(usr)

	if err != nil {
		return "No se ha insertado", 500
	}

	return "Nuevo usuario creado", 201
}

func NewPass(account_id string, jsonStr []byte) (string, int) {

	/* Función que recibe los valores de ID_ACCOUNT como string, y como JSON de nuevo pass y confirm_pass 
		para actualizarlos en la BD */

	if Connect() == false { // Conecta a la base de datos
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

	passValues := &models.Usuario{}
	json.Unmarshal(jsonStr, passValues)

	if UserExists("_id", account_id) == false{
		return "Usuario no encontrado", 400		//Verifica que el account_id exista en la base de datos
	}

	if(passValues.Pass != passValues.Confirm_pass) {
		return "Las contraseñas no coinciden", 409
	}

	passValues.Pass = EncryptToString(passValues.Pass)	// Encripta la contraseña
	passValues.Confirm_pass = EncryptToString(passValues.Confirm_pass)

    col = session.DB(NameDB).C("usuarios")

    colQuerier := bson.M{"_id": bson.ObjectIdHex(account_id)}  // Busca el documento por ID_ACCOUNT
	change := bson.M{"$set": bson.M{"pass": passValues.Pass, "confirm_pass": passValues.Confirm_pass}}
	err = col.Update(colQuerier, change)

	if err != nil {		
		return "Usuario no encontrado", 400
	}

	return "Password actualizada", 200

}

func Auth(jsonStr []byte) (string, int) {

	/*	Función que recibe el los valores en JSON de autenticación y el número del intento
		Devuelve el token en caso de exito o un mensaje en caso contrario, y el estátus del servidor */

	type LogValues struct { // Datos del host solicitante
		Nickname 	string   	`json:"nickname"`
		Pass  		string   	`json:"pass"`
		Ip  		string   	`json:"ip"`
	}

	logValues := LogValues{}  // Llama al modelo de usuario para almacenar los valores recibidos en JSON
    json.Unmarshal(jsonStr, &logValues)

	if Connect() == false { // Conecta a la base de datos
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

    if IsMaliciousIp(logValues.Ip) {
    	return "La IP desde donde se intenta acceder se encuentra bloqueada", 401
    }

    col = session.DB(NameDB).C("usuarios")

    result := models.Usuario{}
    pass := EncryptToString(logValues.Pass)
    err = col.Find(bson.M{"nickname": logValues.Nickname, "pass": pass}).One(&result) // Busca un nombre en la colección y lo almacena en result

	if err != nil {
		return "Usuario no encontrado",400
	}

    token := CreateToken(bson.ObjectId(result.Id).Hex(), result.Firstname+" "+result.Lastname)

	return token,200

}

func MailRecover(mail string) (string, int){

	/* Función que recibe el correo al que se enviará el link de recuperación de pass por medio de parámetro URL 
		Verifica que el correo exista y devuelve un mensaje en caso de que el correo haya sido enviado o de error en caso contrario */

	if Connect() == false { // Conecta a la base de datos
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

    col = session.DB(NameDB).C("usuarios")

    result := models.Usuario{}
    err = col.Find(bson.M{"mail": mail}).One(&result) // Busca un nombre en la colección y lo almacena en result

	if err != nil {
		return "Usuario no encontrado",400
	}

	mailing.PassRecoverMail(result.Mail)  // Llama a la función que envía el correo con el link de recuperación

	return "Se ha enviado un correo al usuario: "+result.Nickname+", correo: " +result.Mail, 200 

}

func UserExists(userBy, textUser string) bool {

	/* Función que verifica si existe algún valor en la base de datos recibe la opción a buscar (userBy) y el valor a buscar(textUser)
		Retorna true en caso de encontrarlo o false cuando no se encuentra*/

    col = session.DB(NameDB).C("usuarios")

    findBson := bson.M{userBy:textUser}
    if userBy == "_id" {
    	if len(textUser)  != 24 {
    		return false
    	} else {
    		findBson = bson.M{userBy: bson.ObjectIdHex(textUser)}
    	}
    } else{
    	findBson = bson.M{userBy:textUser}
    }
    
    result := models.Usuario{}
    err := col.Find(findBson).One(&result) // Busca un nombre en la colección y lo almacena en result

    if err != nil{
    	return false
    } else{
    	return true
    }
}

func UserEdit(account_id, token string,jsonStr []byte) (string, int){

	if Connect() == false { // Conecta a la base de datos
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

	if CheckToken(token) == false {
		return "token no válido", 401
	} else if UserExists("_id", account_id) == false{
		return "Usuario no encontrado", 403		//Verifica que el account_id exista en la base de datos
	}

	editValues := &models.Usuario{}
	json.Unmarshal(jsonStr, editValues)

    col = session.DB(NameDB).C("usuarios")

    colQuerier := bson.M{"_id": bson.ObjectIdHex(account_id)}  // Busca el documento por ACCOUNT_ID
	change := bson.M{"$set": bson.M{"firstname": editValues.Firstname, "lastname": editValues.Lastname, "age": editValues.Age, "country":editValues.Country,"state": editValues.State, "address": editValues.Address, "tel": editValues.Tel}}
	err = col.Update(colQuerier, change)

	if err != nil {		
		return "Usuario no encontrado", 400
	}

	return "Datos de usuario almacenados", 200

}

func AddIp(jsonStr []byte) (string, int) {

	type Applicant struct { // Datos del host solicitante
		Ip  		string   	`json:"ip"`
		nickname 	string   	`json:"nickname"`
	}

	applicant := Applicant{}

    json.Unmarshal(jsonStr, &applicant)

	if Connect() == false { // Conecta a la base de datos
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

    col = session.DB(NameDB).C("blacklist")
	err = col.Insert(bson.M{"ip":applicant.Ip})

	if err != nil {
		return "No se ha insertado", 500
	}

	/* 	En esta parte se enviará un correo al usuario para informarle
		SendAlertmail(applicant.Nickname)
	*/

	return "Ip maliciosa agregada a lista negra", 201

}

func IsMaliciousIp(ip string) bool {

	type Ip_blacklist struct {
		nickname 	string   	`json:"id" bson:"_id,omitempty"`
		Ip  		string   	`json:"ip"`
	}

	result := Ip_blacklist{}

	col = session.DB(NameDB).C("blacklist")

	err = col.Find(bson.M{"ip": ip}).One(&result)

	if err != nil{
    	return false
    } else{
    	return true
    }

}
