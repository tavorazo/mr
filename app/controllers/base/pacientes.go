package base

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"encoding/json"
	"strconv"

	"mr/app/models"
)

func NewPatient(account_id, reference_id string, token string, jsonStr []byte) (string, int) {

	/* Función que recibe los valores de nickname como string, y como JSON del paciente nuevo que se insertará en la BD */

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

	patientVals := &models.Patient{}
	json.Unmarshal(jsonStr, patientVals)

    con := session.DB(NameDB).C("pacientes")

    patientId := 1

    if ReferenceExists(reference_id, session){

    	patientId = GetPatientsNumber(reference_id, session) + 1  // Obtiene el número de pacientes para asignarles el folio de paciente como id

    } else { // Si no existe el id de referencia (doctor/clinica) se crea un nuevo documento para el array de pacientes

    	newReference := &models.PatientsExt{}
    	newReference.Reference_id = reference_id
    	newReference.Account_id = bson.ObjectIdHex(account_id)
    	err = con.Insert(newReference)

    }

    if err != nil {		
		return "No se ha podido crear documento de pacientes", 400
	}

    patientVals.Id = strconv.Itoa(patientId) // Inserta el folio de paciente como id

    colQuerier := bson.M{"reference_id": reference_id}  // Busca el documento por id de referencia (doctor/clinica)
	change := bson.M{"$push": bson.M{"patients": patientVals} } // Inserta en el array de productos
	err = con.Update(colQuerier, change)

	if err != nil {		
		return "No se ha podido insertar el paciente", 400
	}

	return "Paciente agregado con éxito", 201
}

func ReferenceExists(reference_id string, session *mgo.Session) bool {

	con := session.DB(NameDB).C("pacientes")

    result := models.PatientsExt{}
    err := con.Find(bson.M{"reference_id": reference_id}).One(&result)

    if err != nil{
    	return false
    } else{
    	return true
    }

}

func GetPatientsNumber(reference_id string, session *mgo.Session) int {

	con := session.DB(NameDB).C("pacientes")

	result := models.PatientsExt{}

	err := con.Find(bson.M{"reference_id": reference_id}).Select(bson.M{"patients": 1, "_id": 0 }).One(&result)

	if err != nil  {
    	return 0
    }

    patientsFound := len(result.Patients) // Cantidad de pacientes encontrados

    return patientsFound

}

func GetPatients(all bool, account_id, reference_id string, token, patient_id string) (string, int, interface{}) {

	/* 	Función que busca en la base de datos uno o más pacientes
		"all" es un valor booleano que indica si se quieren todos los pacientes o uno en específico
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

    con := session.DB(NameDB).C("pacientes")

    result := models.PatientsExt{}

    if all == false {  // Si está desactivada la opción de todos los pacientes buscará uno en específico de acuerdo al id de referencia indicado
    	err = con.Find(bson.M{"reference_id": reference_id}).Select(bson.M{"products": bson.M{"$elemMatch": bson.M{"id":patient_id} }}).One(&result)
    } else{
    	err = con.Find(bson.M{"reference_id": reference_id}).One(&result)
    }
    
    if err != nil  {
    	return "No hay pacientes para el id de referencia", 206, data
    } 	

	data["paciente"] = result.Patients
	patientsFound := len(result.Patients) // Cantidad de pacientes encontrados

	if patientsFound == 0 {
		if all == true {
			return "No hay pacientes para el id de referencia", 206, data["paciente"]  // Si realiza una búsqueda de todos los pacientes retorna exito pero con el array vacío
		} else {
			return "Paciente no encontrado", 400, data["paciente"]				// Si es una búsqueda de un solo paciente retorna error al no ser encontrado
		}
	} else if patientsFound == 1 {
		return "Paciente encontrado", 200, data["paciente"]			// Si hay un solo paciente
	} else {
		return "Pacientes encontrados", 200, data["paciente"]		// Si hay dos o más pacientes
	}
    
}