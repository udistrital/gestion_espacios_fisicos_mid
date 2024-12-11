package services

import (

	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/gestion_espacios_fisicos_mid/models"
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
	cambiarEspacioFisico(idEspacioFisico,true, false)
	cambiarAsignacionEspacioFisicoDependencia(idEspacioFisico,true, false)
	cambiarTipoUsoEspacioFisico(idEspacioFisico,true, false)
	cambiarCampos(idEspacioFisico,true)
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
	cambiarEspacioFisico(idEspacioFisico,false,false)
	cambiarAsignacionEspacioFisicoDependencia(idEspacioFisico,false, false)
	cambiarTipoUsoEspacioFisico(idEspacioFisico,false, false)
	cambiarCampos(idEspacioFisico,false)
	return alerta, outputError
}

func cambiarEspacioFisico(idEspacioFisico int, cambio bool, rollback bool){
	var espacioFisico []models.EspacioFisico
	url := beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico?query=Id:" + strconv.Itoa(idEspacioFisico)
	if err := request.GetJson(url, &espacioFisico); err != nil || espacioFisico[0].Id == 0 {
		if !rollback {
			cambiarEspacioFisico(idEspacioFisico, !cambio, true)
		}else{
			cambiarEspacioFisico(idEspacioFisico, cambio, true)
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
			cambiarEspacioFisico(idEspacioFisico, !cambio, true)
		}else{
			cambiarEspacioFisico(idEspacioFisico, cambio, true)
		}
		logs.Error(err)
		panic(err.Error())
	}
}

func cambiarAsignacionEspacioFisicoDependencia(idEspacioFisico int, cambio bool, rollback bool){
	var asignacionEspacioFisicoDependencia []models.AsignacionEspacioFisicoDependencia
	url := beego.AppConfig.String("OikosCrudUrl") + "asignacion_espacio_fisico_dependencia?query=EspacioFisicoId.Id:" + strconv.Itoa(idEspacioFisico)
	if err := request.GetJson(url, &asignacionEspacioFisicoDependencia); err != nil || asignacionEspacioFisicoDependencia[0].Id == 0 {
		if !rollback {
			cambiarAsignacionEspacioFisicoDependencia(idEspacioFisico, !cambio, true)
			cambiarEspacioFisico(idEspacioFisico, !cambio, true)
		}else{
			cambiarAsignacionEspacioFisicoDependencia(idEspacioFisico, cambio, true)
			cambiarEspacioFisico(idEspacioFisico, cambio, true)
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
			cambiarAsignacionEspacioFisicoDependencia(idEspacioFisico, !cambio, true)
			cambiarEspacioFisico(idEspacioFisico, !cambio, true)
		}else{
			cambiarAsignacionEspacioFisicoDependencia(idEspacioFisico, cambio, true)
			cambiarEspacioFisico(idEspacioFisico, cambio, true)
		}
		logs.Error(err)
		panic(err.Error())
	}
}

func cambiarTipoUsoEspacioFisico(idEspacioFisico int, cambio bool, rollback bool){
	var tipoUsoEspacioFisico []models.TipoUsoEspacioFisico
	url := beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico?query=EspacioFisicoId.Id:" + strconv.Itoa(idEspacioFisico)

	if err := request.GetJson(url, &tipoUsoEspacioFisico); err != nil || tipoUsoEspacioFisico[0].Id == 0 {
		if !rollback {
			cambiarTipoUsoEspacioFisico(idEspacioFisico, !cambio, true)
			cambiarAsignacionEspacioFisicoDependencia(idEspacioFisico, !cambio, true)
			cambiarEspacioFisico(idEspacioFisico, !cambio, true)
		}else{
			cambiarTipoUsoEspacioFisico(idEspacioFisico, cambio, true)
			cambiarAsignacionEspacioFisicoDependencia(idEspacioFisico, cambio, true)
			cambiarEspacioFisico(idEspacioFisico, cambio, true)
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
			cambiarTipoUsoEspacioFisico(idEspacioFisico, !cambio, true)
			cambiarAsignacionEspacioFisicoDependencia(idEspacioFisico, !cambio, true)
			cambiarEspacioFisico(idEspacioFisico, !cambio, true)
		}else{
			cambiarTipoUsoEspacioFisico(idEspacioFisico, cambio, true)
			cambiarAsignacionEspacioFisicoDependencia(idEspacioFisico, cambio, true)
			cambiarEspacioFisico(idEspacioFisico, cambio, true)
		}
		logs.Error(err)
		panic(err.Error())
	}
}

func cambiarCampos(idEspacioFisico int, cambio bool){
	var camposEspacioFisico []models.EspacioFisicoCampo
	url := beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico_campo?query=EspacioFisicoId.Id:" + strconv.Itoa(idEspacioFisico)
	
	if err := request.GetJson(url, &camposEspacioFisico); err != nil{
		cambiarTipoUsoEspacioFisico(idEspacioFisico, !cambio, true)
		cambiarAsignacionEspacioFisicoDependencia(idEspacioFisico, !cambio, true)
		cambiarEspacioFisico(idEspacioFisico, !cambio, true)
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
				rollbackPutEspacioFisicoCampo(modificaciones, !cambio)
			}
			cambiarTipoUsoEspacioFisico(idEspacioFisico, !cambio, true)
			cambiarAsignacionEspacioFisicoDependencia(idEspacioFisico, !cambio, true)
			cambiarEspacioFisico(idEspacioFisico, !cambio, true)
			logs.Error(err)
			panic(err.Error())
		}
		modificaciones = append(modificaciones, campo)
	}

}

func rollbackPutEspacioFisicoCampo(modificaciones []models.EspacioFisicoCampo, cambio bool){
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