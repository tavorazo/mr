package base

import (
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"strconv"
	"time"

	"mr/app/models"
)

func NewPatient(account_id, reference_id string, token string, jsonStr []byte) (string, int) {

	/* Función que recibe los valores de nickname como string, y como JSON del paciente nuevo que se insertará en la BD */

	if Connect() == false { // Conecta a la base de datos
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

	if CheckToken(token) == false {
		return "token no válido", 401   // Verifica que sea un token válido
	} else if UserExists("_id", account_id) == false{
		return "Usuario no encontrado", 403		//Verifica que el account_id exista en la base de datos
	}

	patientVals := &models.Patient{}
	json.Unmarshal(jsonStr, patientVals)

    col = session.DB(NameDB).C("pacientes")

    patientId := 1

    if ReferenceExists(reference_id, account_id){
    	patientId = GetLastPatientId(reference_id, account_id) + 1  // Obtiene el número de pacientes para asignarles el folio de paciente como id
    } else { // Si no existe el id de referencia (doctor/clinica) se crea un nuevo documento para el array de pacientes

    	newReference := &models.PatientsExt{}
    	newReference.Reference_id = reference_id
    	newReference.Account_id = bson.ObjectIdHex(account_id)
    	err = col.Insert(newReference)

    }

    if err != nil {		
		return "No se ha podido crear documento de pacientes", 400
	}

    patientVals.Id = strconv.Itoa(patientId) // Inserta el folio de paciente como id

    colQuerier := bson.M{"reference_id": reference_id}  // Busca el documento por id de referencia (doctor/clinica)
	change := bson.M{"$push": bson.M{"patients": patientVals} } // Inserta en el array de productos
	err = col.Update(colQuerier, change)

	if err != nil {		
		return "No se ha podido insertar el paciente", 400
	}

	return "Paciente agregado con éxito", 201
}

func ReferenceExists(reference_id, account_id string) bool {

	/* Función que verifica si existe un array de pacientes para el id de referencia indicado */

	col = session.DB(NameDB).C("pacientes")

    result := models.PatientsExt{}
    err = col.Find(bson.M{"reference_id": reference_id, "account_id": bson.ObjectIdHex(account_id)}).One(&result)

    if err != nil{
    	return false
    } else{
    	return true
    }

}

func GetLastPatientId(reference_id, account_id string) int {

	/* Función que devuelve en numero entero el último ID de paciente de determinado id de referencia */

	col = session.DB(NameDB).C("pacientes")

	result := models.PatientsExt{}

	err = col.Find(bson.M{"reference_id": reference_id, "account_id": bson.ObjectIdHex(account_id)}).Select(bson.M{"patients": bson.M{"$slice": -1 }, "_id": 0 }).One(&result)

	if err != nil  {
    	return 0
    }

    if len(result.Patients) == 0 { // Cantidad de pacientes encontrados
    	return 0
    } 

    var lastId int
    lastId, err = strconv.Atoi(result.Patients[0].Id) // Cantidad de pacientes encontrados

    if err != nil  {
    	return 0
    }

    return lastId

}

func GetPatients(all bool, account_id, reference_id string, token, patient_id string) (string, int, interface{}) {

	/* 	Función que busca en la base de datos uno o más pacientes
		"all" es un valor booleano que indica si se quieren todos los pacientes o uno en específico
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

    col = session.DB(NameDB).C("pacientes")

    result := models.PatientsExt{}

    if all == false {  // Si está desactivada la opción de todos los pacientes buscará uno en específico de acuerdo al id de referencia indicado
    	err = col.Find(bson.M{"reference_id": reference_id, "account_id": bson.ObjectIdHex(account_id)}).Select(bson.M{"products": bson.M{"$elemMatch": bson.M{"id":patient_id} }}).One(&result)
    } else{
    	err = col.Find(bson.M{"reference_id": reference_id, "account_id": bson.ObjectIdHex(account_id)}).One(&result)
    }
    
    if err != nil  {
    	return "No se encontraron pacientes", 206, data
    } 	

	data["paciente"] = result.Patients
	patientsFound := len(result.Patients) // Cantidad de pacientes encontrados

	if patientsFound == 0 {
		if all == true {
			return "No se encontraron pacientes", 206, data["paciente"]  // Si realiza una búsqueda de todos los pacientes retorna exito pero con el array vacío
		} else {
			return "Paciente no encontrado", 400, data["paciente"]				// Si es una búsqueda de un solo paciente retorna error al no ser encontrado
		}
	} else if patientsFound == 1 {
		return "Paciente encontrado", 200, data["paciente"]			// Si hay un solo paciente
	} else {
		return "Pacientes encontrados", 200, data["paciente"]		// Si hay dos o más pacientes
	}
    
}

func UpdatePatient(account_id, reference_id string, patient_id, token string, jsonStr []byte) (string, int){

	/* Función que actualiza un paciente para un usuario, con la referencia (doctor/clínica) indicada en la base de datos 
		Se reciben el id de usuario, el id de referencia y el id de paciente (folio) */

	if Connect() == false { // Conecta a la base de datos
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

	if CheckToken(token) == false {
		return "token no válido", 401   // Verifica que sea un token válido
	} else if UserExists("_id", account_id) == false{
		return "Usuario no encontrado", 403		//Verifica que el account_id exista en la base de datos
	}

	patientVals := &models.Patient{}
	json.Unmarshal(jsonStr, patientVals)

	patientVals.Id = patient_id 	// Evita que se borre el id de paciente

    col = session.DB(NameDB).C("pacientes")

    colQuerier := bson.M{"reference_id": reference_id, "account_id": bson.ObjectIdHex(account_id), "patients._id": patient_id }  // Busca el documento
	change := bson.M{"$set": bson.M{"patients.$": patientVals} } // Inserta en el array de productos
	err = col.Update(colQuerier, change)

	if err != nil {		
		return "Paciente no encontrado", 400
	}

	return "Datos del paciente actualizados", 200

}

func AddPicture(account_id, reference_id string, patient_id, token string, pictureType string) (string, int, interface{}){

	/* Función que actualiza un paciente para agregar una imagen al array de fotos intra o radiografías */

	data := make(map[string]interface{})

	if Connect() == false { // Conecta a la base de datos
		return "No se ha conectado a la base de datos", 500, data
    }
    defer session.Close()

	if CheckToken(token) == false {
		return "token no válido", 401, data   // Verifica que sea un token válido
	} else if UserExists("_id", account_id) == false{
		return "Usuario no encontrado", 403, data		//Verifica que el account_id exista en la base de datos
	}

	col = session.DB(NameDB).C("pacientes")

	fileName := "/data/paciente/"+ pictureType +"/" + patient_id+ "_"+ CreateUniqueId() +".jpg"

	colQuerier := bson.M{"reference_id": reference_id, "account_id": bson.ObjectIdHex(account_id), "patients._id": patient_id }  // Busca el documento por id de referencia (doctor/clinica)
	change := bson.M{"$push": bson.M{"patients.$."+pictureType: fileName} } // Inserta en el array de productos
	err = col.Update(colQuerier, change)

	if err != nil {		
		return "Paciente no encontrado", 400, data
	}

	data["type"] = pictureType
	data["name"] = fileName

	return "Imagen agregada", 201, data
}

func AddPrescription(account_id, reference_id string, patient_id, token string, jsonStr []byte) (string, int) {

	if Connect() == false { // Conecta a la base de datos
		return "No se ha conectado a la base de datos", 500
    }
    defer session.Close()

	if CheckToken(token) == false {
		return "token no válido", 401  // Verifica que sea un token válido
	} else if UserExists("_id", account_id) == false{
		return "Usuario no encontrado", 404	//Verifica que el account_id exista en la base de datos
	} else if ReferenceExists(reference_id, account_id) == false {
		return "ID de clínica/doctor no existe",404
	}

	col = session.DB(NameDB).C("pacientes")

	prescriptionVals := &models.Prescription{}
	json.Unmarshal(jsonStr, prescriptionVals)

	prescriptionVals.Patient_id = patient_id
	prescriptionVals.Date = time.Now().Format("02/01/2006")

	colQuerier := bson.M{"reference_id": reference_id}  // Busca el documento por id de referencia (doctor/clinica)
	change := bson.M{"$push": bson.M{"prescriptions": prescriptionVals} } // Inserta en el array de productos
	err = col.Update(colQuerier, change)

	if err != nil {		
		return "No se ha agregado la receta del paciente", 500
	}

	return "Receta agregada al paciente", 201
}