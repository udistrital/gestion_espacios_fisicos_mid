package services

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/espacios_fisicos_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

func PutActivarEspacioFisico(idEspacioFisico int) (alerta []string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "PutActivarEspacioFisico", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	alerta = append(alerta, "Success")
	var cambios models.CambiosActivarDesactivar
	CambiarEspacioFisico(idEspacioFisico,true)
	ActivarAsignacionEspacioFisicoDependencia(idEspacioFisico, &cambios)
	ActivarTipoUsoEspacioFisico(idEspacioFisico, &cambios)
	ActivarCampos(idEspacioFisico, &cambios)
	return alerta, outputError
}

func PutDesactivarEspacioFisico(idEspacioFisico int) (alerta []string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "PutDesactivarEspacioFisico", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	alerta = append(alerta, "Success")
	var cambios models.CambiosActivarDesactivar
	CambiarEspacioFisico(idEspacioFisico,false)
	DesactivarAsignacionEspacioFisicoDependencia(idEspacioFisico, &cambios)
	DesactivarTipoUsoEspacioFisico(idEspacioFisico, &cambios)
	DesactivarCampos(idEspacioFisico, &cambios)

	return alerta, outputError
}

func CambiarEspacioFisico(idEspacioFisico int, cambio bool){
	var espacioFisico []models.EspacioFisico
	url := beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico?query=Id:" + strconv.Itoa(idEspacioFisico)
	if err := request.GetJson(url, &espacioFisico); err != nil || espacioFisico[0].Id == 0 {
		logs.Error(err)
		panic(err.Error())
	}

	espacioFisico[0].Activo = cambio
	espacioFisico[0].FechaModificacion = time_bogota.TiempoBogotaFormato()
	url = beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico/" + strconv.Itoa(espacioFisico[0].Id)
	var respuestaEspacioFisicoModificado map[string]interface{}
	if err := request.SendJson(url, "PUT", &respuestaEspacioFisicoModificado, espacioFisico[0]); respuestaEspacioFisicoModificado["Status"] == "404" {
		logs.Error(err)
		panic(err.Error())
	}
}


func ActivarAsignacionEspacioFisicoDependencia(idEspacioFisico int, cambios *models.CambiosActivarDesactivar){
	var asignacionEspacioFisicoDependencia []models.AsignacionEspacioFisicoDependencia
	url := beego.AppConfig.String("OikosCrudUrl") + "asignacion_espacio_fisico_dependencia?query=EspacioFisicoId.Id:" + strconv.Itoa(idEspacioFisico) + "&sortby=FechaCreacion&order=desc"
	if err := request.GetJson(url, &asignacionEspacioFisicoDependencia); err != nil{
		CambiarEspacioFisico(idEspacioFisico, false)
		logs.Error(err)
		panic(err.Error())
	}
	if len(asignacionEspacioFisicoDependencia) > 0 {
		asignacionEspacioFisicoDependencia[0].Activo = true
		asignacionEspacioFisicoDependencia[0].FechaModificacion = time_bogota.TiempoBogotaFormato()
		url = beego.AppConfig.String("OikosCrudUrl") + "asignacion_espacio_fisico_dependencia/" + strconv.Itoa(asignacionEspacioFisicoDependencia[0].Id)
	
		var respuestaAsignacionEspacioFisicoDependenciaModificado map[string]interface{}
		if err := request.SendJson(url, "PUT", &respuestaAsignacionEspacioFisicoDependenciaModificado, asignacionEspacioFisicoDependencia[0]); respuestaAsignacionEspacioFisicoDependenciaModificado["Status"] == "404" {
			CambiarEspacioFisico(idEspacioFisico, false)
			logs.Error(err)
			panic(err.Error())
		}
		cambios.IdAsignacion = asignacionEspacioFisicoDependencia[0]
	}
}

func DesactivarAsignacionEspacioFisicoDependencia(idEspacioFisico int, cambios *models.CambiosActivarDesactivar){
	var asignacionEspacioFisicoDependencia []models.AsignacionEspacioFisicoDependencia
	url := beego.AppConfig.String("OikosCrudUrl") + "asignacion_espacio_fisico_dependencia?query=Activo:true,EspacioFisicoId.Id:" + strconv.Itoa(idEspacioFisico)
	if err := request.GetJson(url, &asignacionEspacioFisicoDependencia); err != nil{
		CambiarEspacioFisico(idEspacioFisico, true)
		logs.Error(err)
		panic(err.Error())
	}
	if len(asignacionEspacioFisicoDependencia) > 0 {
		asignacionEspacioFisicoDependencia[0].Activo = false
		asignacionEspacioFisicoDependencia[0].FechaModificacion = time_bogota.TiempoBogotaFormato()
		url = beego.AppConfig.String("OikosCrudUrl") + "asignacion_espacio_fisico_dependencia/" + strconv.Itoa(asignacionEspacioFisicoDependencia[0].Id)
		var respuestaAsignacionEspacioFisicoDependenciaModificado map[string]interface{}
		if err := request.SendJson(url, "PUT", &respuestaAsignacionEspacioFisicoDependenciaModificado, asignacionEspacioFisicoDependencia[0]); respuestaAsignacionEspacioFisicoDependenciaModificado["Status"] == "404" {
			CambiarEspacioFisico(idEspacioFisico, true)
			logs.Error(err)
			panic(err.Error())
		}
		cambios.IdAsignacion = asignacionEspacioFisicoDependencia[0]
	}
}

func ActivarTipoUsoEspacioFisico(idEspacioFisico int, cambios *models.CambiosActivarDesactivar){
	var tipoUsoEspacioFisico []models.TipoUsoEspacioFisico
	url := beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico?query=EspacioFisicoId.Id:" + strconv.Itoa(idEspacioFisico)  + "&sortby=FechaCreacion&order=desc"

	if err := request.GetJson(url, &tipoUsoEspacioFisico); err != nil{
		RollbackAsignacionEspacioFisicoDependenciaActivarDesactivar(idEspacioFisico, false, *cambios)
		logs.Error(err)
		panic(err.Error())
	}

	if len(tipoUsoEspacioFisico) > 0 {
		tipoUsoEspacioFisico[0].Activo = true
		tipoUsoEspacioFisico[0].FechaModificacion = time_bogota.TiempoBogotaFormato()
		url = beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico/" + strconv.Itoa(tipoUsoEspacioFisico[0].Id)
		var respuestaTipoUsoEspacioFisico map[string]interface{}
		if err := request.SendJson(url, "PUT", &respuestaTipoUsoEspacioFisico, tipoUsoEspacioFisico[0]); respuestaTipoUsoEspacioFisico["Status"] == "404" {
			RollbackAsignacionEspacioFisicoDependenciaActivarDesactivar(idEspacioFisico, false, *cambios)
			logs.Error(err)
			panic(err.Error())
		}
		cambios.IdTipoUso = tipoUsoEspacioFisico[0]
	}
}

func DesactivarTipoUsoEspacioFisico(idEspacioFisico int, cambios *models.CambiosActivarDesactivar){
	var tipoUsoEspacioFisico []models.TipoUsoEspacioFisico
	url := beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico?query=Activo:true,EspacioFisicoId.Id:" + strconv.Itoa(idEspacioFisico)

	if err := request.GetJson(url, &tipoUsoEspacioFisico); err != nil{
		RollbackAsignacionEspacioFisicoDependenciaActivarDesactivar(idEspacioFisico, true, *cambios)
		logs.Error(err)
		panic(err.Error())
	}

	if len(tipoUsoEspacioFisico) > 0 {
		tipoUsoEspacioFisico[0].Activo = false
		tipoUsoEspacioFisico[0].FechaModificacion = time_bogota.TiempoBogotaFormato()
		url = beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico/" + strconv.Itoa(tipoUsoEspacioFisico[0].Id)
		var respuestaTipoUsoEspacioFisico map[string]interface{}
		if err := request.SendJson(url, "PUT", &respuestaTipoUsoEspacioFisico, tipoUsoEspacioFisico[0]); respuestaTipoUsoEspacioFisico["Status"] == "404" {
			RollbackAsignacionEspacioFisicoDependenciaActivarDesactivar(idEspacioFisico, true, *cambios)
			logs.Error(err)
			panic(err.Error())
		}
		cambios.IdTipoUso = tipoUsoEspacioFisico[0]
	}
}

func ActivarCampos(idEspacioFisico int, cambios *models.CambiosActivarDesactivar){
	var camposEspacioFisico []models.EspacioFisicoCampo
	url := beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico_campo?query=EspacioFisicoId.Id:" + strconv.Itoa(idEspacioFisico) + "&sortby=FechaCreacion&order=desc"
	
	if err := request.GetJson(url, &camposEspacioFisico); err != nil{
		RollbackTipoUsoEspacioFisicoActivarDesactivar(idEspacioFisico,false, *cambios)
		logs.Error(err)
		panic(err.Error())
	}

	var fechaLimite string
	if len(camposEspacioFisico) > 0 {
		if len(camposEspacioFisico[0].FechaCreacion) >= 16 {
			fechaLimite = camposEspacioFisico[0].FechaCreacion[:16]
		}
	}
	for _, campo := range camposEspacioFisico{
		if campo.FechaCreacion[:16] == fechaLimite{
			campo.Activo = true
			campo.FechaModificacion = time_bogota.TiempoBogotaFormato()
			url = beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico_campo/" + strconv.Itoa(campo.Id)
			
			var respuestaEspacioFisicoCampo map[string]interface{}
			if err := request.SendJson(url, "PUT", &respuestaEspacioFisicoCampo, campo); respuestaEspacioFisicoCampo["Status"] == "404" {
				if len(cambios.IdsCampos) > 0{
					RollbackPutEspacioFisicoCampo(idEspacioFisico, *cambios, false)
				}else{
					RollbackTipoUsoEspacioFisicoActivarDesactivar(idEspacioFisico,false, *cambios)
				}
				logs.Error(err)
				panic(err.Error())
			}
			cambios.IdsCampos = append(cambios.IdsCampos, campo)
		}
	}

}

func DesactivarCampos(idEspacioFisico int, cambios *models.CambiosActivarDesactivar){
	var camposEspacioFisico []models.EspacioFisicoCampo
	url := beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico_campo?query=Activo:true,EspacioFisicoId.Id:" + strconv.Itoa(idEspacioFisico)
	
	if err := request.GetJson(url, &camposEspacioFisico); err != nil{
		RollbackTipoUsoEspacioFisicoActivarDesactivar(idEspacioFisico, true, *cambios)
		logs.Error(err)
		panic(err.Error())
	}

	for _, campo := range camposEspacioFisico{

		campo.Activo = false
		campo.FechaModificacion = time_bogota.TiempoBogotaFormato()
		url = beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico_campo/" + strconv.Itoa(campo.Id)
		var respuestaEspacioFisicoCampo map[string]interface{}
		if err := request.SendJson(url, "PUT", &respuestaEspacioFisicoCampo, campo); respuestaEspacioFisicoCampo["Status"] == "404" {
			if len(cambios.IdsCampos) > 0{
				RollbackPutEspacioFisicoCampo(idEspacioFisico, *cambios, true)
			}else{
				RollbackTipoUsoEspacioFisicoActivarDesactivar(idEspacioFisico,true, *cambios)
			}
			logs.Error(err)
			panic(err.Error())
		}
		cambios.IdsCampos = append(cambios.IdsCampos, campo)
	}

}

func RollbackPutEspacioFisicoCampo(idEspacioFisico int, cambios models.CambiosActivarDesactivar, cambio bool){
	for _, campo := range cambios.IdsCampos{
		campo.Activo = cambio
		campo.FechaModificacion = time_bogota.TiempoBogotaFormato()
		url := beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico_campo/" + strconv.Itoa(campo.Id)
		var respuestaEspacioFisicoCampo map[string]interface{}
		if err := request.SendJson(url, "PUT", &respuestaEspacioFisicoCampo, campo); respuestaEspacioFisicoCampo["Status"] == "404" {
			logs.Error(err)
			panic(err.Error())
		}
	}
	RollbackTipoUsoEspacioFisicoActivarDesactivar(idEspacioFisico,cambio, cambios)
}

func RollbackAsignacionEspacioFisicoDependenciaActivarDesactivar(idEspacioFisico int ,cambio bool ,cambios models.CambiosActivarDesactivar){
	cambios.IdAsignacion.Activo = cambio
	cambios.IdAsignacion.FechaModificacion = time_bogota.TiempoBogotaFormato()
	url := beego.AppConfig.String("OikosCrudUrl") + "asignacion_espacio_fisico_dependencia/" + strconv.Itoa(cambios.IdAsignacion.Id)
	var respuestaAsignacionEspacioFisicoDependencia map[string]interface{}
	if err := request.SendJson(url, "PUT", &respuestaAsignacionEspacioFisicoDependencia, cambios.IdAsignacion); respuestaAsignacionEspacioFisicoDependencia["Status"] == "404" {
		logs.Error(err)
		panic(err.Error())
	}
	CambiarEspacioFisico(idEspacioFisico, cambio)
}

func RollbackTipoUsoEspacioFisicoActivarDesactivar(idEspacioFisico int ,cambio bool ,cambios models.CambiosActivarDesactivar){
	cambios.IdTipoUso.Activo = cambio
	cambios.IdTipoUso.FechaModificacion = time_bogota.TiempoBogotaFormato()
	url := beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico/" + strconv.Itoa(cambios.IdTipoUso.Id)
	var respuestaTipoUsoEspacioFisico map[string]interface{}
	if err := request.SendJson(url, "PUT", &respuestaTipoUsoEspacioFisico, cambios.IdTipoUso); respuestaTipoUsoEspacioFisico["Status"] == "404" {
		logs.Error(err)
		panic(err.Error())
	}
	RollbackAsignacionEspacioFisicoDependenciaActivarDesactivar(idEspacioFisico, cambio, cambios)
}