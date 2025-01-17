package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/espacios_fisicos_mid/helpers"
	"github.com/udistrital/espacios_fisicos_mid/models"
	"github.com/udistrital/espacios_fisicos_mid/services"
)

// RegistroController operations for Registro
type RegistroEspaciosFisicosController struct {
	beego.Controller
}

// URLMapping ...
func (c *RegistroEspaciosFisicosController) URLMapping() {
	c.Mapping("RegistroEspacioFisico", c.RegistroEspacioFisicos)
}

// RegistroEspacioFisico ...
// @Title RegistroEspacioFisico
// @Description Registro Espacio Fisico
// @Param	body		body 	{}	true		"body for Registro Espacio Fisico content"
// @Success 201 {init}
// @Failure 400 the request contains incorrect syntax
// @router /RegistroEspacioFisico [post]
func (c *RegistroEspaciosFisicosController) RegistroEspacioFisicos() {
	fmt.Println("REGISTRO ESPACIO FISICO")
	defer helpers.ErrorController(c.Controller, "RegistroEspacioFisico")

	if v, e := helpers.ValidarBody(c.Ctx.Input.RequestBody); !v || e != nil {
		panic(map[string]interface{}{"funcion": "RegistroEspacioFisico", "err": helpers.ErrorBody, "status": "400"})
	}

	var v models.NuevoEspacioFisico

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if resultado, err := services.RegistrarEspacioFisico(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 201, "Message": "Resultado del registro de espacio fisico", "Data": resultado}
		} else {
			panic(err)
		}
	} else {
		panic(map[string]interface{}{"funcion": "RegistroEspacioFisicos", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}
