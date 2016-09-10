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

	c.Response.Status = status
	data["mensaje"] = result
	data["status"] = status

	return c.RenderJson(data)

}
