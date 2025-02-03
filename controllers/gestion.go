package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/astaxie/beego"
	"github.com/udistrital/espacios_fisicos_mid/helpers"
	"github.com/udistrital/espacios_fisicos_mid/models"
	"github.com/udistrital/espacios_fisicos_mid/services"
)

// GestionEspaciosFisicosController operations for GestionEspaciosFisicos
type GestionEspaciosFisicosController struct {
	beego.Controller
}

//URLMapping...
func (c *GestionEspaciosFisicosController) URLMapping() {
	c.Mapping("EditarEspacioFisico", c.EditarEspacioFisico)
	c.Mapping("ActivarEspacioFisico", c.PutActivarEspacioFisico)
	c.Mapping("DesactivarEspacioFisico", c.PutDesactivarEspacioFisico)
}

// EditarEspacioFisico ...
// @Title EditarEspacioFisico
// @Description Editar espacio fisico
// @Param	body		body 	{}	true		"body for Editar Espcaio Fisico content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /EditarEspacioFisico/:id [put]
func (c *GestionEspaciosFisicosController) EditarEspacioFisico() {
	defer helpers.ErrorController(c.Controller, "EditarEspacioFisico")
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)

	if v, e := helpers.ValidarBody(c.Ctx.Input.RequestBody); !v || e != nil {
		panic(map[string]interface{}{"funcion": "EditarEspacioFisico", "err": helpers.ErrorBody, "status": "400"})
	}
	var v models.EditarEspaciosFisicos
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if resultado, err := services.EditarEspacioFisico(&v,id); err == nil {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Espacio Fisico editado con exito", "Data": resultado}
		} else {
			panic(err)
		}
	} else {
		panic(map[string]interface{}{"funcion": "EditarEspacioFisico", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}

// PutActivarEspacioFisico ...
// @Title PutActivarEspacioFisico
// @Description Activar el espacio fisico
// @Param   body        body    {}  true        "body Activar Espacio Fisico content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /ActivarEspacioFisico/:id [put]
func (c *GestionEspaciosFisicosController) PutActivarEspacioFisico() {
	defer helpers.ErrorController(c.Controller,"ActivarEspacioFisico")
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)

	resultado, err := services.PutActivarEspacioFisico(id)

	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Espacio Fisico activado con exito", "Data": resultado}
	} else {
		panic(map[string]interface{}{"funcion": "PutActivarEspacioFisico", "err": err, "status": "400"})
	}

	c.ServeJSON()
}

// PutDesactivarEspacioFisico ...
// @Title PutActivarEspacioFisico
// @Description Activar el espacio fisico
// @Param   body        body    {}  true        "body Desactivar Espacio Fisico content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /DesactivarEspacioFisico/:id [put]
func (c *GestionEspaciosFisicosController) PutDesactivarEspacioFisico() {
	defer helpers.ErrorController(c.Controller,"DesactivarEspacioFisico")
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)

	resultado, err := services.PutDesactivarEspacioFisico(id)

	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Espacio Fisico desactivado con exito", "Data": resultado}
	} else {
		panic(map[string]interface{}{"funcion": "PutDesactivarEspacioFisico", "err": err, "status": "400"})
	}

	c.ServeJSON()
}
