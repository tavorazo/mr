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

	session, err := mgo.Dial("mongodb://mr-tooth:12qwaszx@ds044699.mlab.com:44699/mr")

	if err != nil {
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)
	con := session.DB("mr").C("usuarios")
	err = con.Insert(usr)

	if err != nil {
		return "No se ha insertado", 500
	}

	return "Nuevo usuario creado", 201
}

func Auth(jsonStr []byte) (string, int) {

	/*	Función que recibe el los valores en JSON de autenticación y el número del intento
		Devuelve el token en caso de exito o un mensaje en caso contrario, y el estátus del servidor */

	logValues := &models.Usuario{}  // Llama al modelo de usuario para almacenar los valores recibidos en JSON
    json.Unmarshal(jsonStr, logValues)

	session, err := mgo.Dial("mongodb://mr-tooth:12qwaszx@ds044699.mlab.com:44699/mr")

	if err != nil {
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)
    con := session.DB("mr").C("usuarios")

    result := models.Usuario{}
    pass := EncryptToString(logValues.Pass)
    err = con.Find(bson.M{"nickname": logValues.Nickname, "pass": pass}).One(&result) // Busca un nombre en la colección y lo almacena en result

	if err != nil {
		return "Usuario no encontrado",401
	}

    token := CreateToken(result.Nickname,result.Firstname+" "+result.Lastname)

	return token,200

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