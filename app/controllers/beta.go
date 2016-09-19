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

func (c Beta) UpdatePass(nickname string) revel.Result{

	/* Recibe como parámetro URL el usuario al que se actualizará la contraseña 
	   y de Body (ocultos) el JSON con los valores de pass y confirm_pass */

	data := make(map[string]interface{})

	body, _ := ioutil.ReadAll(c.Request.Body)
	result,status := base.NewPass(nickname, body)

	if(status != 200){
		data["error"] = result
	} else{
		data["OK"] = result
	}

	c.Response.Status = status
	data["status"] = status

	return c.RenderJson(data)
}
