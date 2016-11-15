package controllers

import (
	"github.com/revel/revel"
	"mr/app/controllers/base"
	"io/ioutil"
)

type Beta struct {
	*revel.Controller
}

func (c Beta) SignUp(json_vals string) revel.Result {

	data := make(map[string]interface{})

	body, _ := ioutil.ReadAll(c.Request.Body)  //Recibe de POST la cadena correspondiente a un JSON
	result,status := base.NewUser(body)

	
	if(status != 201){
		data["error"] = result
	} else{
		data["OK"] = result
	}
	data["status"] = status

	c.Response.Status = status
	return c.RenderJson(data)

}

func (c Beta) Login() revel.Result{

	data := make(map[string]interface{})

	body, _ := ioutil.ReadAll(c.Request.Body)  //Recibe de POST la cadena correspondiente a un JSON
	result,status := base.Auth(body)

	if(status != 200){
		data["error"] = result
	} else{
		data["OK"] = "Ha ingresado correctamente"
		data["token"] = result
	}

	c.Response.Status = status
	data["status"] = status

	return c.RenderJson(data)
}

func (c Beta) RecoverAccount(mail string) revel.Result{

	data := make(map[string]interface{})

	// body, _ := ioutil.ReadAll(c.Request.Body)
	result,status := base.MailRecover(mail)

	if(status != 200){
		data["error"] = result
	} else{
		data["OK"] = result
	}

	c.Response.Status = status
	data["status"] = status

	return c.RenderJson(data)

}

func (c Beta) UpdatePass(account_id string) revel.Result{

	/* Recibe como parámetro URL el usuario al que se actualizará la contraseña 
	   y de Body (ocultos) el JSON con los valores de pass y confirm_pass */

	data := make(map[string]interface{})

	body, _ := ioutil.ReadAll(c.Request.Body)
	result,status := base.NewPass(account_id, body)

	if(status != 200){
		data["error"] = result
	} else{
		data["OK"] = result
	}

	c.Response.Status = status
	data["status"] = status

	return c.RenderJson(data)
}

func (c Beta) UpdateAccount(account_id string) revel.Result{

	data := make(map[string]interface{})

	body, _ := ioutil.ReadAll(c.Request.Body)  //Recibe de POST la cadena correspondiente a un JSON
	token := c.Request.Header.Get("token") 		// Lee token desde header

	result,status := base.UserEdit(account_id, token, body)
	
	if(status != 200){
		data["error"] = result
	} else{
		data["OK"] = result
	}

	c.Response.Status = status
	data["status"] = status

	return c.RenderJson(data)
}

func (c Beta) AddProduct(account_id string) revel.Result{

	data := make(map[string]interface{})

	body, _ := ioutil.ReadAll(c.Request.Body)  //Recibe de POST la cadena correspondiente a un JSON
	token := c.Request.Header.Get("token")		// Lee token desde header

	result,status := base.NewProduct(account_id, token, body)

	if(status != 201){
		data["error"] = result
	} else{
		data["OK"] = result
	}

	c.Response.Status = status
	data["status"] = status

	return c.RenderJson(data)
}

func (c Beta) EditAmount(account_id, product_id string) revel.Result {
	data := make(map[string]interface{})

	body, _ := ioutil.ReadAll(c.Request.Body)  //Recibe de POST la cadena correspondiente a un JSON
	token := c.Request.Header.Get("token")		// Lee token desde header

	result,status := base.UpdateProductAmount(account_id, product_id, token, body)

	if(status != 200){
		data["error"] = result
	} else{
		data["OK"] = result
	}

	c.Response.Status = status
	data["status"] = status

	return c.RenderJson(data)
}

func (c Beta) EditProduct(account_id, product_id string) revel.Result {

	data := make(map[string]interface{})

	body, _ := ioutil.ReadAll(c.Request.Body)  //Recibe de POST la cadena correspondiente a un JSON
	token := c.Request.Header.Get("token")		// Lee token desde header

	result,status := base.UpdateProduct(account_id, product_id, token, body)

	if(status != 200){
		data["error"] = result
	} else{
		data["OK"] = result
	}

	c.Response.Status = status
	data["status"] = status

	return c.RenderJson(data)

}

func (c Beta) DeleteProduct(account_id, product_id string) revel.Result {

	data := make(map[string]interface{})

	token := c.Request.Header.Get("token")		// Lee token desde header

	result,status := base.SaveDeletedProduct(account_id, product_id, token) // Llama a la función que pondrá en punto de restauración el producto indicado

	if(status != 200){
		data["error"] = result
	} else{
		data["OK"] = result
	}

	c.Response.Status = status
	data["status"] = status

	return c.RenderJson(data)

}

func (c Beta) ProductsAll(account_id string) revel.Result {

	data := make(map[string]interface{})

	token := c.Request.Header.Get("token")		// Lee token desde header

	result, status, datos := base.GetProducts(true, account_id, token, "") // True para activar la busqueda de todos los productos

	if(status != 200){
		data["error"] = result
	} else{
		data["OK"] = result
		data["products"] = datos
	}

	c.Response.Status = status
	data["status"] = status

	return c.RenderJson(data)
}

func (c Beta) ProductsOne(account_id, product_id string) revel.Result {

	data := make(map[string]interface{})

	token := c.Request.Header.Get("token")		// Lee token desde header

	result, status, datos := base.GetProducts(false, account_id, token, product_id) // False para desactivar la busqueda de todos los productos

	if(status != 200){
		data["error"] = result
	} else{
		data["OK"] = result
		data["products"] = datos
	}

	c.Response.Status = status
	data["status"] = status

	return c.RenderJson(data)
}

// Proveedores

func (c Beta) AddCaterer(account_id string) revel.Result {

	data := make(map[string]interface{})

	body, _ := ioutil.ReadAll(c.Request.Body)  //Recibe de POST la cadena correspondiente a un JSON
	token := c.Request.Header.Get("token")		// Lee token desde header

	result,status := base.NewCaterer(account_id, token, body)

	if(status != 201){
		data["error"] = result
	} else{
		data["OK"] = result
	}

	c.Response.Status = status
	data["status"] = status

	return c.RenderJson(data)
}

func (c Beta) EditCaterer(account_id string) revel.Result {

	data := make(map[string]interface{})

	body, _ := ioutil.ReadAll(c.Request.Body)  //Recibe de POST la cadena correspondiente a un JSON
	token := c.Request.Header.Get("token")		// Lee token desde header

	result,status := base.UpdateCaterer(account_id, token, body)

	if(status != 200){
		data["error"] = result
	} else{
		data["OK"] = result
	}

	c.Response.Status = status
	data["status"] = status

	return c.RenderJson(data)

}

// PACIENTES

func (c Beta) AddPatient(account_id, reference_id string) revel.Result {

	data := make(map[string]interface{})

	body, _ := ioutil.ReadAll(c.Request.Body)  //Recibe de POST la cadena correspondiente a un JSON
	token := c.Request.Header.Get("token")		// Lee token desde header

	result,status := base.NewPatient(account_id, reference_id, token, body)

	if(status != 201){
		data["error"] = result
	} else{
		data["OK"] = result
	}

	c.Response.Status = status
	data["status"] = status

	return c.RenderJson(data)
}

func (c Beta) PatientsAll(account_id, reference_id string) revel.Result {

	data := make(map[string]interface{})

	token := c.Request.Header.Get("token")		// Lee token desde header

	result, status, datos := base.GetPatients(true, account_id, reference_id, token, "") // True para activar la busqueda de todos los productos

	if(status != 200){
		data["error"] = result
	} else{
		data["OK"] = result
		data["patients"] = datos
	}

	c.Response.Status = status
	data["status"] = status

	return c.RenderJson(data)
}

func (c Beta) EditPatient(account_id, reference_id string, patient_id string) revel.Result {

	data := make(map[string]interface{})

	body, _ := ioutil.ReadAll(c.Request.Body)  //Recibe de POST la cadena correspondiente a un JSON
	token := c.Request.Header.Get("token")		// Lee token desde header

	result,status := base.UpdatePatient(account_id, reference_id, patient_id, token, body)

	if(status != 200){
		data["error"] = result
	} else{
		data["OK"] = result
	}

	c.Response.Status = status
	data["status"] = status

	return c.RenderJson(data)

}