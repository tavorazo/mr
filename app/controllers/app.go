package controllers

import "github.com/revel/revel"

type DatosJson struct{
	NombreApp 	string 	` json:"nombreApp" `
	Mensaje 	string	` json:"mensaje" `
	Version 	int 	` json:"version" `
}

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Mensaje() revel.Result {
	
	data := make(map[string]interface{})
	data["error"] = nil

	json := DatosJson{NombreApp: "Mr-Tooth", Mensaje: "Hola mundo!", Version:0}
	data["datos"] = json

	return c.RenderJson(data)
}

