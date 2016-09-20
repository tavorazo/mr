package base

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/base64"
	"strconv"
	"time"

	"mr/app/models"
	"mr/app/controllers/mailing"
)

var (
	HostDB 			string = "mongodb://mr-tooth:12qwaszx@ds044699.mlab.com:44699/mr"
	NameDB 			string = "mr"
	CollectionDB 	string = "usuarios"
)

func NewUser(jsonStr []byte) (string, int) {

	/*	Función que recibe el valor en JSON y lo inserta en la Base de datos
		Devuelve Un mensaje, el estátus del servidor y el error si existe */

	usr := &models.Usuario{}
	json.Unmarshal(jsonStr, usr)			// Recibe el valor json y lo almacena en la estructura

	if(usr.Pass != usr.Confirm_pass) {
		return "Las contraseñas no coinciden", 400
	} else if(usr.Mail != usr.Confirm_mail) {
		return "Los correos no coinciden", 400
	}

	usr.Pass = EncryptToString(usr.Pass)	// Encripta la contraseña
	usr.Confirm_pass = EncryptToString(usr.Confirm_pass)

	session, err := mgo.Dial(HostDB)

	if err != nil {
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)
	con := session.DB(NameDB).C(CollectionDB)
	err = con.Insert(usr)

	if err != nil {
		return "No se ha insertado", 500
	}

	return "Nuevo usuario creado", 201
}

func NewPass(nickname string,jsonStr []byte) (string, int) {

	/* Función que recibe los valores de nickname como string, y como JSON de nuevo pass y confirm_pass 
		para actualizarlos en la BD */

	passValues := &models.Usuario{}
	json.Unmarshal(jsonStr, passValues)

	if(passValues.Pass != passValues.Confirm_pass) {
		return "Las contraseñas no coinciden", 400
	}

	passValues.Pass = EncryptToString(passValues.Pass)	// Encripta la contraseña
	passValues.Confirm_pass = EncryptToString(passValues.Confirm_pass)

	session, err := mgo.Dial(HostDB)

	if err != nil {
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)
    con := session.DB(NameDB).C(CollectionDB)

    colQuerier := bson.M{"nickname": nickname}  // Busca el documento por nickname
	change := bson.M{"$set": bson.M{"pass": passValues.Pass, "confirm_pass": passValues.Confirm_pass}}
	err = con.Update(colQuerier, change)

	if err != nil {		
		return "Usuario no encontrado", 401
	}

	return "Password actualizada", 200

}

func Auth(jsonStr []byte) (string, int) {

	/*	Función que recibe el los valores en JSON de autenticación y el número del intento
		Devuelve el token en caso de exito o un mensaje en caso contrario, y el estátus del servidor */

	logValues := &models.Usuario{}  // Llama al modelo de usuario para almacenar los valores recibidos en JSON
    json.Unmarshal(jsonStr, logValues)

	session, err := mgo.Dial(HostDB)

	if err != nil {
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)
    con := session.DB(NameDB).C(CollectionDB)

    result := models.Usuario{}
    pass := EncryptToString(logValues.Pass)
    err = con.Find(bson.M{"nickname": logValues.Nickname, "pass": pass}).One(&result) // Busca un nombre en la colección y lo almacena en result

	if err != nil {
		return "Usuario no encontrado",401
	}

    token := CreateToken(result.Nickname,result.Firstname+" "+result.Lastname)

	return token,200

}

func MailRecover(mail string) (string, int){

	/* Función que recibe el correo al que se enviará el link de recuperación de pass por medio de parámetro URL 
		Verifica que el correo exista y devuelve un mensaje en caso de que el correo haya sido enviado o de error en caso contrario */

	session, err := mgo.Dial(HostDB)

	if err != nil {
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)
    con := session.DB(NameDB).C(CollectionDB)

    result := models.Usuario{}
    err = con.Find(bson.M{"mail": mail}).One(&result) // Busca un nombre en la colección y lo almacena en result

	if err != nil {
		return "Usuario no encontrado",401
	}

	mailing.PassRecoverMail(result.Mail)  // Llama a la función que envía el correo con el link de recuperación

	return "Se ha enviado un correo al usuario: "+result.Nickname+", correo: " +result.Mail, 200 

}

func CreateToken(iss string, name string) string{

	/* 	Función que recopila la información del usuario y lo codifica para formar el token del tipo JWT.
		Recibe el nombre de usuario (nickname) y el nombre completo del usuario para codificarlos dentro del token */

	header := []byte(`{"typ": "JWT", "alg": "HS256"}`)

	exp := int(time.Now().Unix()) + 86400    // Tiempo de expiración del token más 1 día
	payload := []byte(`{"iss": "`+iss+`", "exp": `+strconv.Itoa(exp)+`, "name": "`+name+`"}`)

	signature := base64.StdEncoding.EncodeToString(header)+"."+base64.StdEncoding.EncodeToString(payload)

	// Forma un token del tipo [header].[payload].[signature]   								
	token := base64.StdEncoding.EncodeToString(header)+"."+base64.StdEncoding.EncodeToString(payload)+"."+EncryptToStringSha256(signature)
	return token
}

// Funciónes que reciben un texto, lo codifica y lo devuelve como un string
func EncryptToString(text string) string{
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
func EncryptToStringSha256(text string) string{
	hasher := sha256.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}