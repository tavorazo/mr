package base

import(
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/base64"
	"strconv"
	"time"
)

func CreateToken(iss string, name string) string{

	/* 	Función que recopila la información del usuario y lo codifica para formar el token del tipo JWT.
		Recibe el nombre de usuario (nickname) y el nombre completo del usuario para codificarlos dentro del token */

	header := []byte(`{"typ": "JWT", "alg": "HS256"}`)

	exp := int(time.Now().Unix()) + 86400    // Tiempo de expiración del token más 1 día
	payload := []byte(`{"iss": "`+iss+`", "exp": "`+strconv.Itoa(exp)+`", "name": "`+name+`"}`)

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