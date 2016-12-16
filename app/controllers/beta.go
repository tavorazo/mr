package controllers

import (
	"github.com/revel/revel"
	"mr/app/controllers/base"
	"io/ioutil"
)

type Beta struct {
	*revel.Controller
}

var (
	result 		string
	status 		int
	dataArray 	interface{}
	body 		[]byte
)

func (c Beta) SignUp(json_vals string) revel.Result {

	body, _  			= ioutil.ReadAll(c.Request.Body)  //Recibe de POST la cadena correspondiente a un JSON
	result,status 		= base.NewUser(body)
	c.Response.Status 	= status
	return c.RenderJson(OkResponse())

}

func (c Beta) Login() revel.Result{

	body, _ 			= ioutil.ReadAll(c.Request.Body)  //Recibe de POST la cadena correspondiente a un JSON
	result,status 		= base.Auth(body)
	c.Response.Status 	= status
	return c.RenderJson(OkResponse())
}

func (c Beta) RecoverAccount(mail string) revel.Result{

	// body, _ := ioutil.ReadAll(c.Request.Body)
	result,status 		= base.MailRecover(mail)
	c.Response.Status 	= status

	return c.RenderJson(OkResponse())

}

func (c Beta) UpdatePass(account_id string) revel.Result{

	body, _ 			= ioutil.ReadAll(c.Request.Body)
	result,status 		= base.NewPass(account_id, body)
	c.Response.Status 	= status
	return c.RenderJson(OkResponse())

}

func (c Beta) UpdateAccount(account_id string) revel.Result{

	body, _ 			= ioutil.ReadAll(c.Request.Body)  //Recibe de POST la cadena correspondiente a un JSON
	result,status 		= base.UserEdit(account_id, c.Request.Header.Get("token"), body)
	c.Response.Status 	= status
	return c.RenderJson(OkResponse())
}

func (c Beta) AddIp(json_vals string) revel.Result {

	body, _ 			= ioutil.ReadAll(c.Request.Body)  //Recibe de POST la cadena correspondiente a un JSON
	result,status 		= base.AddIp(body)
	c.Response.Status 	= status
	return c.RenderJson(OkResponse())

}


// Almacén

func (c Beta) AddProduct(account_id, reference_id string) revel.Result{

	body, _ 			= ioutil.ReadAll(c.Request.Body)  //Recibe de POST la cadena correspondiente a un JSON
	result,status 		= base.NewProduct(account_id, reference_id, c.Request.Header.Get("token"), body)
	c.Response.Status 	= status
	return c.RenderJson(OkResponse())
}

func (c Beta) EditAmount(account_id, reference_id string, product_id string) revel.Result {

	body, _ 			= ioutil.ReadAll(c.Request.Body)  //Recibe de POST la cadena correspondiente a un JSON
	result,status 		= base.UpdateProductAmount(account_id, reference_id, product_id, c.Request.Header.Get("token"), body)
	c.Response.Status 	= status
	return c.RenderJson(OkResponse())
}

func (c Beta) EditProduct(account_id, reference_id string, product_id string) revel.Result {

	body, _ 			= ioutil.ReadAll(c.Request.Body)  //Recibe de POST la cadena correspondiente a un JSON
	result,status 		= base.UpdateProduct(account_id, reference_id, product_id, c.Request.Header.Get("token"), body)
	c.Response.Status 	= status
	return c.RenderJson(OkResponse())

}

func (c Beta) DeleteProduct(account_id, reference_id string, product_id string) revel.Result {

	result,status 		= base.SaveDeletedProduct(account_id, reference_id, product_id, c.Request.Header.Get("token")) // Llama a la función que pondrá en punto de restauración el producto indicado
	c.Response.Status 	= status
	return c.RenderJson(OkResponse())

}

func (c Beta) ProductsAll(account_id, reference_id string) revel.Result {

	result, status, dataArray 	= base.GetProducts(true, account_id, reference_id, c.Request.Header.Get("token"), "") // True para activar la busqueda de todos los productos
	c.Response.Status 			= status
	return c.RenderJson(ShowArrayResponse("products"))
}

func (c Beta) ProductsOne(account_id, reference_id string, product_id string) revel.Result {

	result, status, dataArray 	= base.GetProducts(false, account_id, reference_id, c.Request.Header.Get("token"), product_id) // False para desactivar la busqueda de todos los productos
	c.Response.Status 			= status
	return c.RenderJson(ShowArrayResponse("products"))
}

// Proveedores

func (c Beta) AddCaterer(account_id string) revel.Result {

	body, _ 			= ioutil.ReadAll(c.Request.Body)  //Recibe de POST la cadena correspondiente a un JSON
	result,status 		= base.NewCaterer(account_id, c.Request.Header.Get("token"), body)
	c.Response.Status 	= status
	return c.RenderJson(OkResponse())
}

func (c Beta) EditCaterer(account_id string) revel.Result {

	body, _ 			= ioutil.ReadAll(c.Request.Body)  //Recibe de POST la cadena correspondiente a un JSON
	result,status 		= base.UpdateCaterer(account_id, c.Request.Header.Get("token"), body)
	c.Response.Status 	= status
	return c.RenderJson(OkResponse())

}

// PACIENTES

func (c Beta) AddPatient(account_id, reference_id string) revel.Result {

	body, _ 			= ioutil.ReadAll(c.Request.Body)  //Recibe de POST la cadena correspondiente a un JSON
	result,status 		= base.NewPatient(account_id, reference_id, c.Request.Header.Get("token"), body)
	c.Response.Status 	= status
	return c.RenderJson(OkResponse())
}

func (c Beta) GetPatients(account_id, reference_id string, patient_id string) revel.Result {

	result, status, dataArray 	= base.GetPatients(account_id, reference_id, c.Request.Header.Get("token"), patient_id) // True para activar la busqueda de todos los productos
	c.Response.Status 			= status
	return c.RenderJson(ShowArrayResponse("data"))
}

func (c Beta) EditPatient(account_id, reference_id string, patient_id string) revel.Result {

	body, _ 			= ioutil.ReadAll(c.Request.Body)  //Recibe de POST la cadena correspondiente a un JSON
	result,status 		= base.UpdatePatient(account_id, reference_id, patient_id, c.Request.Header.Get("token"), body)
	c.Response.Status 	= status
	return c.RenderJson(OkResponse())

}

func (c Beta) AddIntraPicture(account_id, reference_id string, patient_id string) revel.Result {

	result,status,dataArray = base.AddPicture(account_id, reference_id, patient_id, c.Request.Header.Get("token"),"intra")
	c.Response.Status 		= status
	return c.RenderJson(ShowArrayResponse("picture"))

}

func (c Beta) AddRadiography(account_id, reference_id string, patient_id string) revel.Result {

	result,status,dataArray = base.AddPicture(account_id, reference_id, patient_id, c.Request.Header.Get("token"),"radiography")
	c.Response.Status 		= status
	return c.RenderJson(ShowArrayResponse("picture"))

}

func (c Beta) AddPrescription(account_id, reference_id string, patient_id string) revel.Result {
	body, _ 			= ioutil.ReadAll(c.Request.Body)  //Recibe de POST la cadena correspondiente a un JSON
	result,status 		= base.AddPrescription(account_id, reference_id, patient_id, c.Request.Header.Get("token"), body)
	c.Response.Status 	= status
	return c.RenderJson(OkResponse())
}

// Funciones comunes de respuesta
func OkResponse() interface{} {

	data := make(map[string]interface{})
	if(status >= 400){
		data["error"]  	= 	result
	} else{
		data["OK"]  	= 	result
	}
	data["status"] = status

	return data

}
func ShowArrayResponse(ArrayName string) interface{} {

	data := make(map[string]interface{})
	if(status >= 400){
		data["error"] 	= 	result
	} else{
		data["OK"] 		= 	result
		data[ArrayName] = 	dataArray
	}
	data["status"] = status

	return data
}