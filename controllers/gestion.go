package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/gestion_espacios_fisicos_mid/helpers"
	"github.com/udistrital/gestion_espacios_fisicos_mid/models"
	"github.com/udistrital/gestion_espacios_fisicos_mid/services"
)

// GestionEspaciosFisicosController operations for GestionEspaciosFisicos
type GestionEspaciosFisicosController struct {
	beego.Controller
}

//URLMapping...
func (c *GestionEspaciosFisicosController) URLMapping(){
	c.Mapping("BuscarEspacioFisico", c.BuscarEspacioFisico)
}

// BuscarEspacioFisico ...
// @Title BuscarEspacioFisico
// @Description Buscar Espacio Fisico
// @Param	body		body 	{}	true		"body for Buscar Espacio Fisico content"
// @Success 201 {init} 
// @Failure 400 the request contains incorrect syntax
// @router /BuscarEspacioFisico [post]
func (c *GestionEspaciosFisicosController) BuscarEspacioFisico() {
	fmt.Println("BUSCA ESPACIO FISICO")
	defer helpers.ErrorController(c.Controller,"BuscarEspacioFisico")

	if v, e := helpers.ValidarBody(c.Ctx.Input.RequestBody); !v || e != nil {
		panic(map[string]interface{}{"funcion": "BuscarEspacioFisico", "err": helpers.ErrorBody, "status": "400"})
	}

	var v models.BusquedaEspacioFisico

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if resultado, err := services.BuscarEspacioFisico(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Resultado de busqueda", "Data": resultado}
		} else {
			panic(err)
		}
	} else {
		panic(map[string]interface{}{"funcion": "BuscarEspacioFisico", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}