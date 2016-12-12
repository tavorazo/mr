package base

import(
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/base64"
	"strconv"
	"time"
	"math/rand"
	"strings"
	"mr/app/models"
)

func CreateToken(iss string, name string) string{

	/* 	Función que recopila la información del usuario y lo codifica para formar el token del tipo JWT.
		Recibe el ACCOUNT_ID y el nombre completo del usuario para codificarlos dentro del token */

	header := []byte(`{"alg":"HS256","typ":"JWT"}`)

	exp := int(time.Now().Unix()) + 86400    // Tiempo de expiración del token más 1 día
	payload := []byte(`{"iss":"`+iss+`","exp":`+strconv.Itoa(exp)+`,"name":"`+name+`"}`)

	signature := base64.StdEncoding.EncodeToString(header)+"."+base64.StdEncoding.EncodeToString(payload)

	// Forma un token del tipo [header].[payload].[signature]   								
	token := base64.StdEncoding.EncodeToString(header)+"."+base64.StdEncoding.EncodeToString(payload)+"."+EncryptToStringSha256(signature)
	return token
}

func CheckToken(token string) bool {

	/* Función que verifica la expliración del token y retorna true si es válido o false en caso contrario */

	tokSplit := strings.Split(token, ".")
	if (len(tokSplit) != 3) {
		return false
	}

	payloadJson, err := base64.StdEncoding.DecodeString(tokSplit[1]) // Divide el token y extrae la parte de payload

	if err != nil{
		return false
	}

	type Payload struct {
		Iss  	string   	`json:"iss"`
		Name 	string   	`json:"name"`// Se almacenarán los datos del payload
		Exp 	int 		`json:"exp"`
	}

	pl := Payload{}
	json.Unmarshal(payloadJson, &pl)

    col = session.DB(NameDB).C("usuarios")

    result := models.Usuario{}
    err = col.Find(bson.M{"_id":bson.ObjectIdHex(pl.Iss)}).One(&result)

    if err != nil {  //Si no se encuentra el usuario en la base de datos retorna falso
    	return false
    } else if pl.Exp < int(time.Now().Unix()) {  
		return false
	} else {
		return true
	}
	
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
func CreateUniqueId() string {
	return EncryptToString( strconv.Itoa( int(time.Now().Unix()) * rand.Int() ) )
}