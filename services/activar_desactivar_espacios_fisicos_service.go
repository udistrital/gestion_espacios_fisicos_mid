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
	CambiarEspacioFisico(idEspacioFisico,true, false)
	CambiarAsignacionEspacioFisicoDependencia(idEspacioFisico,true, false)
	CambiarTipoUsoEspacioFisico(idEspacioFisico,true, false)
	CambiarCampos(idEspacioFisico,true)
	return alerta, outputError
}

func PutDesactivarEspacioFisico(idEspacioFisico int) (alerta []string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "PutActivarEspacioFisico", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	alerta = append(alerta, "Success")
	CambiarEspacioFisico(idEspacioFisico,false,false)
	CambiarAsignacionEspacioFisicoDependencia(idEspacioFisico,false, false)
	CambiarTipoUsoEspacioFisico(idEspacioFisico,false, false)
	CambiarCampos(idEspacioFisico,false)
	return alerta, outputError
}

func CambiarEspacioFisico(idEspacioFisico int, cambio bool, rollback bool){
	var espacioFisico []models.EspacioFisico
	url := beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico?query=Id:" + strconv.Itoa(idEspacioFisico)
	if err := request.GetJson(url, &espacioFisico); err != nil || espacioFisico[0].Id == 0 {
		if !rollback {
			CambiarEspacioFisico(idEspacioFisico, !cambio, true)
		}else{
			CambiarEspacioFisico(idEspacioFisico, cambio, true)
		}
		logs.Error(err)
		panic(err.Error())
	}

	espacioFisico[0].Activo = cambio
	espacioFisico[0].FechaModificacion = time_bogota.TiempoBogotaFormato()
	url = beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico/" + strconv.Itoa(espacioFisico[0].Id)
	var respuestaEspacioFisicoModificado map[string]interface{}
	if err := request.SendJson(url, "PUT", &respuestaEspacioFisicoModificado, espacioFisico[0]); respuestaEspacioFisicoModificado["Status"] == "404" {
		if !rollback {
			CambiarEspacioFisico(idEspacioFisico, !cambio, true)
		}else{
			CambiarEspacioFisico(idEspacioFisico, cambio, true)
		}
		logs.Error(err)
		panic(err.Error())
	}
}

func CambiarAsignacionEspacioFisicoDependencia(idEspacioFisico int, cambio bool, rollback bool){
	var asignacionEspacioFisicoDependencia []models.AsignacionEspacioFisicoDependencia
	url := beego.AppConfig.String("OikosCrudUrl") + "asignacion_espacio_fisico_dependencia?query=EspacioFisicoId.Id:" + strconv.Itoa(idEspacioFisico)
	if err := request.GetJson(url, &asignacionEspacioFisicoDependencia); err != nil || asignacionEspacioFisicoDependencia[0].Id == 0 {
		if !rollback {
			CambiarAsignacionEspacioFisicoDependencia(idEspacioFisico, !cambio, true)
			CambiarEspacioFisico(idEspacioFisico, !cambio, true)
		}else{
			CambiarAsignacionEspacioFisicoDependencia(idEspacioFisico, cambio, true)
			CambiarEspacioFisico(idEspacioFisico, cambio, true)
		}
		logs.Error(err)
		panic(err.Error())
	}


	asignacionEspacioFisicoDependencia[0].Activo = cambio
	asignacionEspacioFisicoDependencia[0].FechaModificacion = time_bogota.TiempoBogotaFormato()
	url = beego.AppConfig.String("OikosCrudUrl") + "asignacion_espacio_fisico_dependencia/" + strconv.Itoa(asignacionEspacioFisicoDependencia[0].Id)

	var respuestaAsignacionEspacioFisicoDependenciaModificado map[string]interface{}
	if err := request.SendJson(url, "PUT", &respuestaAsignacionEspacioFisicoDependenciaModificado, asignacionEspacioFisicoDependencia[0]); respuestaAsignacionEspacioFisicoDependenciaModificado["Status"] == "404" {
		if !rollback {
			CambiarAsignacionEspacioFisicoDependencia(idEspacioFisico, !cambio, true)
			CambiarEspacioFisico(idEspacioFisico, !cambio, true)
		}else{
			CambiarAsignacionEspacioFisicoDependencia(idEspacioFisico, cambio, true)
			CambiarEspacioFisico(idEspacioFisico, cambio, true)
		}
		logs.Error(err)
		panic(err.Error())
	}
}

func CambiarTipoUsoEspacioFisico(idEspacioFisico int, cambio bool, rollback bool){
	var tipoUsoEspacioFisico []models.TipoUsoEspacioFisico
	url := beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico?query=EspacioFisicoId.Id:" + strconv.Itoa(idEspacioFisico)

	if err := request.GetJson(url, &tipoUsoEspacioFisico); err != nil || tipoUsoEspacioFisico[0].Id == 0 {
		if !rollback {
			CambiarTipoUsoEspacioFisico(idEspacioFisico, !cambio, true)
			CambiarAsignacionEspacioFisicoDependencia(idEspacioFisico, !cambio, true)
			CambiarEspacioFisico(idEspacioFisico, !cambio, true)
		}else{
			CambiarTipoUsoEspacioFisico(idEspacioFisico, cambio, true)
			CambiarAsignacionEspacioFisicoDependencia(idEspacioFisico, cambio, true)
			CambiarEspacioFisico(idEspacioFisico, cambio, true)
		}
		logs.Error(err)
		panic(err.Error())
	}


	tipoUsoEspacioFisico[0].Activo = cambio
	tipoUsoEspacioFisico[0].FechaModificacion = time_bogota.TiempoBogotaFormato()
	url = beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico/" + strconv.Itoa(tipoUsoEspacioFisico[0].Id)
	if rollback {
		url = beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico/" + strconv.Itoa(tipoUsoEspacioFisico[0].Id)
	}
	var respuestaTipoUsoEspacioFisico map[string]interface{}
	if err := request.SendJson(url, "PUT", &respuestaTipoUsoEspacioFisico, tipoUsoEspacioFisico[0]); respuestaTipoUsoEspacioFisico["Status"] == "404" {
		if !rollback {
			CambiarTipoUsoEspacioFisico(idEspacioFisico, !cambio, true)
			CambiarAsignacionEspacioFisicoDependencia(idEspacioFisico, !cambio, true)
			CambiarEspacioFisico(idEspacioFisico, !cambio, true)
		}else{
			CambiarTipoUsoEspacioFisico(idEspacioFisico, cambio, true)
			CambiarAsignacionEspacioFisicoDependencia(idEspacioFisico, cambio, true)
			CambiarEspacioFisico(idEspacioFisico, cambio, true)
		}
		logs.Error(err)
		panic(err.Error())
	}
}

func CambiarCampos(idEspacioFisico int, cambio bool){
	var camposEspacioFisico []models.EspacioFisicoCampo
	url := beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico_campo?query=EspacioFisicoId.Id:" + strconv.Itoa(idEspacioFisico)
	
	if err := request.GetJson(url, &camposEspacioFisico); err != nil{
		CambiarTipoUsoEspacioFisico(idEspacioFisico, !cambio, true)
		CambiarAsignacionEspacioFisicoDependencia(idEspacioFisico, !cambio, true)
		CambiarEspacioFisico(idEspacioFisico, !cambio, true)
		logs.Error(err)
		panic(err.Error())
	}
	var modificaciones []models.EspacioFisicoCampo
	for _, campo := range camposEspacioFisico{
		campo.Activo = cambio
		campo.FechaModificacion = time_bogota.TiempoBogotaFormato()
		url = beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico_campo/" + strconv.Itoa(campo.Id)
		
		
		var respuestaEspacioFisicoCampo map[string]interface{}
		if err := request.SendJson(url, "PUT", &respuestaEspacioFisicoCampo, campo); respuestaEspacioFisicoCampo["Status"] == "404" {
			if (len(modificaciones)>0) {
				RollbackPutEspacioFisicoCampo(modificaciones, !cambio)
			}
			CambiarTipoUsoEspacioFisico(idEspacioFisico, !cambio, true)
			CambiarAsignacionEspacioFisicoDependencia(idEspacioFisico, !cambio, true)
			CambiarEspacioFisico(idEspacioFisico, !cambio, true)
			logs.Error(err)
			panic(err.Error())
		}
		modificaciones = append(modificaciones, campo)
	}

}

func RollbackPutEspacioFisicoCampo(modificaciones []models.EspacioFisicoCampo, cambio bool){
	for _, campo := range modificaciones{
		campo.Activo = cambio
		campo.FechaModificacion = time_bogota.TiempoBogotaFormato()
		url := beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico_campo/" + strconv.Itoa(campo.Id)
		var respuestaEspacioFisicoCampo map[string]interface{}
		if err := request.SendJson(url, "PUT", &respuestaEspacioFisicoCampo, campo); respuestaEspacioFisicoCampo["Status"] == "404" {
			logs.Error(err)
			panic(err.Error())
		}
	}
}