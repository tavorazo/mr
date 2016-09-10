package base

import (
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	"encoding/json"
	"crypto/md5"
	"encoding/hex"

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

// Función que recibe un texto, lo codifica y lo devuelve como un string
func EncryptToString(text string) string{
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}