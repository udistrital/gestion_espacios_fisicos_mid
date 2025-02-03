package services

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/espacios_fisicos_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

func RegistrarEspacioFisico(transaccion *models.NuevoEspacioFisico) (alerta []string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "RegistrarEspacioFisico", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	alerta = append(alerta, "Success")
	var creaciones models.Creaciones
	var tipoEspacioFisico models.TipoEspacioFisico
	url := beego.AppConfig.String("OikosCrudUrl") + "tipo_espacio_fisico/" + strconv.Itoa(transaccion.TipoEspacioFisico)
	if err := request.GetJson(url, &tipoEspacioFisico); err != nil || tipoEspacioFisico.Id == 0 {
		logs.Error(err)
		panic(err.Error())
	}

	var dependencia models.Dependencia
	url = beego.AppConfig.String("OikosCrudUrl") + "dependencia/" + strconv.Itoa(transaccion.DependenciaPadre)
	if err := request.GetJson(url, &dependencia); err != nil || dependencia.Id == 0 {
		logs.Error(err)
		panic(err.Error())
	}

	var tipoUso models.TipoUso
	url = beego.AppConfig.String("OikosCrudUrl") + "tipo_uso/" + strconv.Itoa(transaccion.TipoUso)
	if err := request.GetJson(url, &tipoUso); err != nil || tipoUso.Id == 0 {
		logs.Error(err)
		panic(err.Error())
	}

	
	var espacioFisico = CrearEspacioFisico(transaccion, tipoEspacioFisico, &creaciones)

	CrearAsignacionEspacioFisicoDependencia(transaccion, dependencia, espacioFisico, &creaciones)
	CrearTipoUsoEspacioFisico(transaccion, tipoUso, espacioFisico, &creaciones)


	if len(transaccion.CamposExistentes) != 0 {
		var espacioFisicoCampoCreados = CrearEspacioFisicoCampo(transaccion, espacioFisico, &creaciones)
		fmt.Println(espacioFisicoCampoCreados)
	}

	return alerta, outputError
}

func CrearEspacioFisico(transaccion *models.NuevoEspacioFisico, tipoEspacioFisico models.TipoEspacioFisico, creaciones *models.Creaciones) (nuevoEspacioFisico models.EspacioFisico) {
	nuevoEspacioFisico.Nombre = transaccion.EspacioFisico.Nombre
	nuevoEspacioFisico.Descripcion = transaccion.EspacioFisico.Descripcion
	nuevoEspacioFisico.CodigoAbreviacion = transaccion.EspacioFisico.CodigoAbreviacion
	nuevoEspacioFisico.TipoTerrenoId = transaccion.TipoTerreno
	nuevoEspacioFisico.TipoEdificacionId = transaccion.TipoEdificacion
	nuevoEspacioFisico.TipoEspacioFisicoId = &tipoEspacioFisico
	nuevoEspacioFisico.Activo = true
	nuevoEspacioFisico.FechaCreacion = time_bogota.TiempoBogotaFormato()
	nuevoEspacioFisico.FechaModificacion = time_bogota.TiempoBogotaFormato()
	url := beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico"
	var resEspacioFisicoRegistrado map[string]interface{}
	if err := request.SendJson(url, "POST", &resEspacioFisicoRegistrado, nuevoEspacioFisico); err != nil || resEspacioFisicoRegistrado["Id"] == nil {
		logs.Error(err)
		panic(err.Error())
	}
	creaciones.EspacioFisicoId = int(resEspacioFisicoRegistrado["Id"].(float64))
	return nuevoEspacioFisico
}

func CrearAsignacionEspacioFisicoDependencia(transaccion *models.NuevoEspacioFisico, dependencia models.Dependencia, espacioFisico models.EspacioFisico, creaciones *models.Creaciones) {
	var asignacionEspacioFisicoDependencia models.AsignacionEspacioFisicoDependencia
	asignacionEspacioFisicoDependencia.EspacioFisicoId = &espacioFisico
	asignacionEspacioFisicoDependencia.DependenciaId = &dependencia
	asignacionEspacioFisicoDependencia.DocumentoSoporte = 0
	asignacionEspacioFisicoDependencia.FechaInicio = time_bogota.TiempoBogotaFormato()
	asignacionEspacioFisicoDependencia.FechaFin = time_bogota.TiempoBogotaFormato()
	asignacionEspacioFisicoDependencia.Activo = true
	asignacionEspacioFisicoDependencia.FechaCreacion = time_bogota.TiempoBogotaFormato()
	asignacionEspacioFisicoDependencia.FechaModificacion = time_bogota.TiempoBogotaFormato()
	url := beego.AppConfig.String("OikosCrudUrl") + "asignacion_espacio_fisico_dependencia"
	var resAsignacionEspacioFisicoDependenciaRegistrado map[string]interface{}
	if err := request.SendJson(url, "POST", &resAsignacionEspacioFisicoDependenciaRegistrado, asignacionEspacioFisicoDependencia); err != nil || resAsignacionEspacioFisicoDependenciaRegistrado["Id"] == nil {
		rollbackEspacioFisicoCreado(creaciones)
		logs.Error(err)
		panic(err.Error())
	}
	creaciones.AsignacionEspacioFisicoDependenciaId = int(resAsignacionEspacioFisicoDependenciaRegistrado["Id"].(float64))
}

func CrearTipoUsoEspacioFisico(transaccion *models.NuevoEspacioFisico, tipoUso models.TipoUso, espacioFisico models.EspacioFisico, creaciones *models.Creaciones) {
	var tipoUsoEspacioFisico models.TipoUsoEspacioFisico
	tipoUsoEspacioFisico.EspacioFisicoId = &espacioFisico
	tipoUsoEspacioFisico.TipoUsoId = &tipoUso
	tipoUsoEspacioFisico.Activo = true
	tipoUsoEspacioFisico.FechaCreacion = time_bogota.TiempoBogotaFormato()
	tipoUsoEspacioFisico.FechaModificacion = time_bogota.TiempoBogotaFormato()
	url := beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico"
	var resTipoUsoEspacioFisicoRegistrado map[string]interface{}
	if err := request.SendJson(url, "POST", &resTipoUsoEspacioFisicoRegistrado, tipoUsoEspacioFisico); err != nil || resTipoUsoEspacioFisicoRegistrado["Id"] == nil {
		rollbackAsignacionEspacioFisicoDependencia(creaciones)
		logs.Error(err)
		panic(err.Error())
	}
	creaciones.TipoUsoEspacioFisico = int(resTipoUsoEspacioFisicoRegistrado["Id"].(float64))
}

func CrearEspacioFisicoCampo(transaccion *models.NuevoEspacioFisico, espacioFisico models.EspacioFisico, creaciones *models.Creaciones) (campos []models.EspacioFisicoCampo) {
	for _, campo := range transaccion.CamposExistentes {
		var nuevoCampo models.Campo
		url := beego.AppConfig.String("OikosCrudUrl") + "campo/" + strconv.Itoa(campo.IdCampo)
		if err := request.GetJson(url, &nuevoCampo); err != nil {
			if len(creaciones.EspacioFisicoCampoId) > 0 {
				rollbackEspacioFisicoCampo(creaciones)
			} else {
				rollbackTipoUsoEspacioFisico(creaciones)
			}
			logs.Error(err)
			panic(err.Error())
		}
		var espacioFisicoCampo models.EspacioFisicoCampo
		espacioFisicoCampo.Valor = campo.Valor
		espacioFisicoCampo.EspacioFisicoId = &espacioFisico
		espacioFisicoCampo.CampoId = &nuevoCampo
		espacioFisicoCampo.Activo = true
		espacioFisicoCampo.FechaInicio = time_bogota.TiempoBogotaFormato()
		espacioFisicoCampo.FechaFin = nil
		espacioFisicoCampo.FechaCreacion = time_bogota.TiempoBogotaFormato()
		espacioFisicoCampo.FechaModificacion = time_bogota.TiempoBogotaFormato()
		url = beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico_campo"
		var resEspacioFisicoCampoRegistrado map[string]interface{}
		if err := request.SendJson(url, "POST", &resEspacioFisicoCampoRegistrado, espacioFisicoCampo); err != nil || resEspacioFisicoCampoRegistrado["Id"] == nil {
			if len(creaciones.EspacioFisicoCampoId) > 0 {
				rollbackEspacioFisicoCampo(creaciones)
			} else {
				rollbackTipoUsoEspacioFisico(creaciones)
			}
			logs.Error(err)
			panic(err.Error())
		}
		creaciones.EspacioFisicoCampoId = append(creaciones.EspacioFisicoCampoId, int(resEspacioFisicoCampoRegistrado["Id"].(float64)))
		espacioFisicoCampo.Id = int(resEspacioFisicoCampoRegistrado["Id"].(float64))
		campos = append(campos, espacioFisicoCampo)
	}
	return campos
}

func rollbackEspacioFisicoCreado(creaciones *models.Creaciones) (outputError map[string]interface{}) {
	var respuesta map[string]interface{}
	url := beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico/" + strconv.Itoa(creaciones.EspacioFisicoId)
	if err := request.SendJson(url, "DELETE", &respuesta, nil); err != nil {
		panic("Rollback del espacio fisico" + err.Error())
	}
	return nil
}

func rollbackAsignacionEspacioFisicoDependencia(creaciones *models.Creaciones) (outputError map[string]interface{}) {
	var respuesta map[string]interface{}
	url := beego.AppConfig.String("OikosCrudUrl") + "asignacion_espacio_fisico_dependencia/" + strconv.Itoa(creaciones.AsignacionEspacioFisicoDependenciaId)
	if err := request.SendJson(url, "DELETE", &respuesta, nil); err != nil {
		panic("Rollback de la asignacion del espacio fisico dependencia" + err.Error())
	}
	rollbackEspacioFisicoCreado(creaciones)
	return nil
}

func rollbackTipoUsoEspacioFisico(creaciones *models.Creaciones) (outputError map[string]interface{}) {
	var respuesta map[string]interface{}
	url := beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico/" + strconv.Itoa(creaciones.TipoUsoEspacioFisico)
	if err := request.SendJson(url, "DELETE", &respuesta, nil); err != nil {
		panic("Rollback del tipo de uso del espacio fisico" + err.Error())
	}
	rollbackAsignacionEspacioFisicoDependencia(creaciones)
	return nil
}

func rollbackEspacioFisicoCampo(creaciones *models.Creaciones) (outputError map[string]interface{}) {
	for _, campo := range creaciones.EspacioFisicoCampoId {
		var respuesta map[string]interface{}
		url := beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico_campo/" + strconv.Itoa(campo)
		if err := request.SendJson(url, "DELETE", &respuesta, nil); err != nil {
			panic("Rollback del espacio fisico campo" + err.Error())
		}
	}
	rollbackTipoUsoEspacioFisico(creaciones)
	return nil
}
