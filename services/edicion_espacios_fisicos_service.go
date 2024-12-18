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

func EditarEspacioFisico(transaccion *models.EditarEspaciosFisicos) (alerta []string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "EditarEspacioFisico", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	alerta = append(alerta, "Success")

	var espacioFisicoOriginal models.EspacioFisico
	url := beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico/" + strconv.Itoa(transaccion.EspacioId)
	fmt.Println("URL ", url)
	if err := request.GetJson(url, &espacioFisicoOriginal); err != nil || espacioFisicoOriginal.Id == 0 {
		logs.Error(err)
		panic(err.Error())
	}

	var tiposUsoOriginal []models.TipoUsoEspacioFisico
	url = beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico?query=EspacioFisicoId.Id:" + strconv.Itoa(transaccion.EspacioId)
	if err := request.GetJson(url, &tiposUsoOriginal); err != nil || tiposUsoOriginal[0].Id == 0 {
		logs.Error(err)
		panic(err.Error())
	}

	var dependenciasEspacioOriginal []models.AsignacionEspacioFisicoDependencia
	url = beego.AppConfig.String("OikosCrudUrl") + "asignacion_espacio_fisico_dependencia?query=EspacioFisicoId.Id:" + strconv.Itoa(transaccion.EspacioId)
	if err := request.GetJson(url, &dependenciasEspacioOriginal); err != nil || dependenciasEspacioOriginal[0].Id == 0 {
		logs.Error(err)
		panic(err.Error())
	}

	var tipoUsoModificado models.TipoUsoEspacioFisico
	var dependenciaEspacioModificado models.AsignacionEspacioFisicoDependencia

	var espacioModificado = ActualizarEspacioFisico(espacioFisicoOriginal, transaccion)
	var existeTipoUsoAnterior = false
	var existeDependenciaAnterior = false
	for _, tipo := range tiposUsoOriginal {
		if tipo.TipoUsoId.Id == transaccion.TipoUsoId {
			existeTipoUsoAnterior = true
		}
	}

	for _, dependencia := range dependenciasEspacioOriginal {
		if dependencia.DependenciaId.Id == transaccion.DependenciaId {
			existeDependenciaAnterior = true
		}
	}

	if !existeTipoUsoAnterior {
		tipoUsoModificado = ActualizarNuevoTipoUso(&espacioModificado, transaccion)
	} else {
		tipoUsoModificado = ActualizarTipoUsoExistente(&espacioModificado, transaccion)
	}

	if !existeDependenciaAnterior {
		dependenciaEspacioModificado = ActualizarNuevaDependencia(&espacioModificado, transaccion)
	} else {
		dependenciaEspacioModificado = ActualizarNuevaDependenciaExistente(&espacioModificado, transaccion)
	}

	fmt.Println(espacioModificado)
	fmt.Println(tipoUsoModificado)
	fmt.Println(dependenciaEspacioModificado)

	fmt.Println(transaccion.CamposExistentes)
	fmt.Println(transaccion.CamposNoExistentes)
	var nuevosCampos []models.EspacioFisicoCampo
	var camposExistentes []models.EspacioFisicoCampo
	if len(*transaccion.CamposExistentes) > 0 {
		camposExistentes = ActualizarCampos(transaccion)
	}
	if len(*transaccion.CamposNoExistentes) > 0 {
		nuevosCampos = AgregarCampos(&espacioModificado, transaccion)
	}
	fmt.Println(nuevosCampos)
	fmt.Println(camposExistentes)

	return alerta, outputError
}

func ActualizarEspacioFisico(espacioFisicoOriginal models.EspacioFisico, transaccion *models.EditarEspaciosFisicos) (espacioModificado models.EspacioFisico) {
	espacioModificado = espacioFisicoOriginal
	espacioModificado.Nombre = transaccion.Nombre
	espacioModificado.CodigoAbreviacion = transaccion.CodAbreviacion
	espacioModificado.Descripcion = transaccion.Descripcion
	espacioModificado.FechaModificacion = time_bogota.TiempoBogotaFormato()
	if espacioFisicoOriginal.TipoEspacioFisicoId.Id != transaccion.TipoEspacioId {
		var nuevoTipoEspacio models.TipoEspacioFisico
		url := beego.AppConfig.String("OikosCrudUrl") + "tipo_espacio_fisico/" + strconv.Itoa(transaccion.TipoEspacioId)
		if err := request.GetJson(url, &nuevoTipoEspacio); err != nil || nuevoTipoEspacio.Id == 0 {
			logs.Error(err)
			panic(err.Error())
		}
		espacioModificado.TipoEspacioFisicoId = &nuevoTipoEspacio
	} else {
		espacioModificado.TipoEspacioFisicoId = espacioFisicoOriginal.TipoEspacioFisicoId
	}

	var err error
	url := beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico/" + strconv.Itoa(transaccion.EspacioId)
	var respuestaEspacioFisicoModificado map[string]interface{}
	if err = request.SendJson(url, "PUT", &respuestaEspacioFisicoModificado, espacioModificado); err != nil {
		logs.Error(err)
		panic(err.Error())
	}
	return espacioModificado
}

func ActualizarNuevoTipoUso(espacioModificado *models.EspacioFisico, transaccion *models.EditarEspaciosFisicos) (tipoUsoModificado models.TipoUsoEspacioFisico) {
	var nuevoTipoUso models.TipoUso
	url := beego.AppConfig.String("OikosCrudUrl") + "tipo_uso/" + strconv.Itoa(transaccion.TipoUsoId)
	if err := request.GetJson(url, &nuevoTipoUso); err != nil || nuevoTipoUso.Id == 0 {
		logs.Error(err)
		panic(err.Error())
	}
	var tipoUsoOriginal []models.TipoUsoEspacioFisico
	url = beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico?query=EspacioFisicoId.Id:" + strconv.Itoa(transaccion.EspacioId) + ",Activo:true"
	if err := request.GetJson(url, &tipoUsoOriginal); err != nil || tipoUsoOriginal[0].Id == 0 {
		logs.Error(err)
		panic(err.Error())
	}

	tipoUsoOriginal[0].Activo = false
	tipoUsoOriginal[0].FechaModificacion = time_bogota.TiempoBogotaFormato()
	var err error
	url = beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico/" + strconv.Itoa(tipoUsoOriginal[0].Id)
	var respuestaTipoUsoOriginal map[string]interface{}
	if err = request.SendJson(url, "PUT", &respuestaTipoUsoOriginal, tipoUsoOriginal[0]); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	tipoUsoModificado.Activo = true
	tipoUsoModificado.FechaCreacion = time_bogota.TiempoBogotaFormato()
	tipoUsoModificado.FechaModificacion = time_bogota.TiempoBogotaFormato()
	tipoUsoModificado.TipoUsoId = &nuevoTipoUso
	tipoUsoModificado.EspacioFisicoId = espacioModificado
	url = beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico"
	var res map[string]interface{}
	if err := request.SendJson(url, "POST", &res, tipoUsoModificado); err != nil || res["Id"] == nil {
		logs.Error(err)
		panic(err.Error())
	}
	return tipoUsoModificado
}

func ActualizarTipoUsoExistente(espacioModificado *models.EspacioFisico, transaccion *models.EditarEspaciosFisicos) (tipoUsoModificado models.TipoUsoEspacioFisico) {

	var tipoUsoOriginal []models.TipoUsoEspacioFisico
	url := beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico?query=EspacioFisicoId.Id:" + strconv.Itoa(transaccion.EspacioId) + ",Activo:true"
	if err := request.GetJson(url, &tipoUsoOriginal); err != nil || tipoUsoOriginal[0].Id == 0 {
		logs.Error(err)
		panic(err.Error())
	}

	if tipoUsoOriginal[0].TipoUsoId.Id != transaccion.TipoUsoId {
		tipoUsoOriginal[0].Activo = false
		tipoUsoOriginal[0].FechaModificacion = time_bogota.TiempoBogotaFormato()
		var err error
		url = beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico/" + strconv.Itoa(tipoUsoOriginal[0].Id)
		var respuestaTipoUsoOriginal map[string]interface{}
		if err = request.SendJson(url, "PUT", &respuestaTipoUsoOriginal, tipoUsoOriginal[0]); err != nil {
			logs.Error(err)
			panic(err.Error())
		}
		var nuevoTipoUso models.TipoUso
		url := beego.AppConfig.String("OikosCrudUrl") + "tipo_uso/" + strconv.Itoa(transaccion.TipoUsoId)
		if err := request.GetJson(url, &nuevoTipoUso); err != nil || nuevoTipoUso.Id == 0 {
			logs.Error(err)
			panic(err.Error())
		}
		var nuevoTipoUsoActivo []models.TipoUsoEspacioFisico
		url = beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico?query=EspacioFisicoId.Id:" + strconv.Itoa(transaccion.EspacioId) + ",TipoUsoId.Id:" + strconv.Itoa(transaccion.TipoUsoId)
		if err := request.GetJson(url, &nuevoTipoUsoActivo); err != nil || nuevoTipoUsoActivo[0].Id == 0 {
			logs.Error(err)
			panic(err.Error())
		}

		tipoUsoModificado = nuevoTipoUsoActivo[0]

		tipoUsoModificado.Activo = true
		tipoUsoModificado.FechaModificacion = time_bogota.TiempoBogotaFormato()
		url = beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico/" + strconv.Itoa(tipoUsoModificado.Id)
		var respuestaTipoUsoExistente map[string]interface{}
		if err = request.SendJson(url, "PUT", &respuestaTipoUsoExistente, tipoUsoModificado); err != nil {
			logs.Error(err)
			panic(err.Error())
		}
	}

	return tipoUsoModificado
}

func ActualizarNuevaDependencia(espacioModificado *models.EspacioFisico, transaccion *models.EditarEspaciosFisicos) (dependenciaEspacioModificado models.AsignacionEspacioFisicoDependencia) {
	var nuevaDependencia models.Dependencia
	url := beego.AppConfig.String("OikosCrudUrl") + "dependencia/" + strconv.Itoa(transaccion.DependenciaId)
	if err := request.GetJson(url, &nuevaDependencia); err != nil || nuevaDependencia.Id == 0 {
		logs.Error(err)
		panic(err.Error())
	}

	var dependenciaEspacioOriginal []models.AsignacionEspacioFisicoDependencia
	url = beego.AppConfig.String("OikosCrudUrl") + "asignacion_espacio_fisico_dependencia?query=EspacioFisicoId.Id:" + strconv.Itoa(transaccion.EspacioId) + ",Activo:true"
	if err := request.GetJson(url, &dependenciaEspacioOriginal); err != nil || dependenciaEspacioOriginal[0].Id == 0 {
		logs.Error(err)
		panic(err.Error())
	}

	dependenciaEspacioOriginal[0].Activo = false
	dependenciaEspacioOriginal[0].FechaModificacion = time_bogota.TiempoBogotaFormato()
	var err error
	url = beego.AppConfig.String("OikosCrudUrl") + "asignacion_espacio_fisico_dependencia/" + strconv.Itoa(dependenciaEspacioOriginal[0].Id) + ",Activo:true"
	var respuestaDependenciaEspacioOriginal map[string]interface{}
	if err = request.SendJson(url, "PUT", &respuestaDependenciaEspacioOriginal, dependenciaEspacioOriginal[0]); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	dependenciaEspacioModificado.Activo = true
	dependenciaEspacioModificado.FechaCreacion = time_bogota.TiempoBogotaFormato()
	dependenciaEspacioModificado.FechaModificacion = time_bogota.TiempoBogotaFormato()
	dependenciaEspacioModificado.DependenciaId = &nuevaDependencia
	dependenciaEspacioModificado.EspacioFisicoId = espacioModificado
	dependenciaEspacioModificado.DocumentoSoporte = 0
	dependenciaEspacioModificado.FechaInicio = time_bogota.TiempoBogotaFormato()
	dependenciaEspacioModificado.FechaFin = time_bogota.TiempoBogotaFormato()
	url = beego.AppConfig.String("OikosCrudUrl") + "asignacion_espacio_fisico_dependencia"
	var res map[string]interface{}
	if err := request.SendJson(url, "POST", &res, dependenciaEspacioModificado); err != nil || res["Id"] == nil {
		logs.Error(err)
		panic(err.Error())
	}
	return dependenciaEspacioModificado
}

func ActualizarNuevaDependenciaExistente(espacioModificado *models.EspacioFisico, transaccion *models.EditarEspaciosFisicos) (dependenciaEspacioModificado models.AsignacionEspacioFisicoDependencia) {
	var dependenciaEspacioOriginal []models.AsignacionEspacioFisicoDependencia
	url := beego.AppConfig.String("OikosCrudUrl") + "asignacion_espacio_fisico_dependencia?query=EspacioFisicoId.Id:" + strconv.Itoa(transaccion.EspacioId) + ",Activo:true"
	if err := request.GetJson(url, &dependenciaEspacioOriginal); err != nil || dependenciaEspacioOriginal[0].Id == 0 {
		logs.Error(err)
		panic(err.Error())
	}

	if dependenciaEspacioOriginal[0].DependenciaId.Id != transaccion.DependenciaId {
		dependenciaEspacioOriginal[0].Activo = false
		dependenciaEspacioOriginal[0].FechaModificacion = time_bogota.TiempoBogotaFormato()
		var err error
		url = beego.AppConfig.String("OikosCrudUrl") + "asignacion_espacio_fisico_dependencia/" + strconv.Itoa(dependenciaEspacioOriginal[0].Id)
		var respuestaDependenciaEspacioOriginal map[string]interface{}
		if err = request.SendJson(url, "PUT", &respuestaDependenciaEspacioOriginal, dependenciaEspacioOriginal[0]); err != nil {
			logs.Error(err)
			panic(err.Error())
		}
		var nuevaDependenciaEspacioActiva []models.AsignacionEspacioFisicoDependencia
		url := beego.AppConfig.String("OikosCrudUrl") + "asignacion_espacio_fisico_dependencia?query=EspacioFisicoId.Id:" + strconv.Itoa(transaccion.EspacioId) + ",DependenciaId.Id:" + strconv.Itoa(transaccion.DependenciaId)
		if err := request.GetJson(url, &nuevaDependenciaEspacioActiva); err != nil || nuevaDependenciaEspacioActiva[0].Id == 0 {
			logs.Error(err)
			panic(err.Error())
		}
		dependenciaEspacioModificado = nuevaDependenciaEspacioActiva[0]

		dependenciaEspacioModificado.Activo = true
		dependenciaEspacioModificado.FechaModificacion = time_bogota.TiempoBogotaFormato()

		url = beego.AppConfig.String("OikosCrudUrl") + "asignacion_espacio_fisico_dependencia/" + strconv.Itoa(dependenciaEspacioModificado.Id)
		var respuestaDependenciaEspacioExistente map[string]interface{}
		if err = request.SendJson(url, "PUT", &respuestaDependenciaEspacioExistente, dependenciaEspacioModificado); err != nil {
			logs.Error(err)
			panic(err.Error())
		}
	}
	return dependenciaEspacioModificado
}

func AgregarCampos(espacioModificado *models.EspacioFisico, transaccion *models.EditarEspaciosFisicos) (nuevosCampos []models.EspacioFisicoCampo) {
	for _, campo := range *transaccion.CamposNoExistentes {
		var campoExistente []models.EspacioFisicoCampo
		url := beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico_campo?query=EspacioFisicoId.Id:" + strconv.Itoa(transaccion.EspacioId) + ",CampoId.Id:" + strconv.Itoa(campo.IdCampo)
		if err := request.GetJson(url, &campoExistente); err != nil {
			logs.Error(err)
			panic(err.Error())
		}
		if len(campoExistente) > 0 {
			CambiarEstadoCampoEspacio(campoExistente[0], campo)
		} else {
			var nuevoCampo models.Campo
			url := beego.AppConfig.String("OikosCrudUrl") + "campo/" + strconv.Itoa(campo.IdCampo)
			if err := request.GetJson(url, &nuevoCampo); err != nil {
				logs.Error(err)
				panic(err.Error())
			}
			fmt.Println(espacioModificado.Id)
			var nuevoCampoEspacio models.EspacioFisicoCampo
			nuevoCampoEspacio.Valor = campo.Valor
			nuevoCampoEspacio.EspacioFisicoId = espacioModificado
			nuevoCampoEspacio.CampoId = &nuevoCampo
			nuevoCampoEspacio.Activo = true
			nuevoCampoEspacio.FechaInicio = time_bogota.TiempoBogotaFormato()
			fechaFin := time_bogota.TiempoBogotaFormato()
			nuevoCampoEspacio.FechaFin = &fechaFin
			nuevoCampoEspacio.FechaCreacion = time_bogota.TiempoBogotaFormato()
			nuevoCampoEspacio.FechaModificacion = time_bogota.TiempoBogotaFormato()
			url = beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico_campo"
			var res map[string]interface{}
			if err := request.SendJson(url, "POST", &res, nuevoCampoEspacio); err != nil || res["Id"] == nil {
				logs.Error(err)
				panic(err.Error())
			}
			nuevosCampos = append(nuevosCampos, nuevoCampoEspacio)
		}
	}
	return nuevosCampos
}

func CambiarEstadoCampoEspacio(campoExistente models.EspacioFisicoCampo, campo models.CamposEspacioFisico) {
	campoExistente.Activo = true
	campoExistente.Valor = campo.Valor
	campoExistente.FechaModificacion = time_bogota.TiempoBogotaFormato()
	var err error
	url := beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico_campo/" + strconv.Itoa(campoExistente.Id)
	var respuestaEspacioFisicoCampo map[string]interface{}
	if err = request.SendJson(url, "PUT", &respuestaEspacioFisicoCampo, campoExistente); err != nil {
		logs.Error(err)
		panic(err.Error())
	}
}

func ActualizarCampos(transaccion *models.EditarEspaciosFisicos) (camposExistentes []models.EspacioFisicoCampo) {

	var camposExistentesEspacio []models.EspacioFisicoCampo
	url := beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico_campo?query=EspacioFisicoId.Id:" + strconv.Itoa(transaccion.EspacioId)
	if err := request.GetJson(url, &camposExistentesEspacio); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	var camposExistentesActivos []models.EspacioFisicoCampo
	for _, campo := range *transaccion.CamposExistentes {
		fmt.Println(campo.IdCampo)
		var campoExistente []models.EspacioFisicoCampo
		url := beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico_campo?query=EspacioFisicoId.Id:" + strconv.Itoa(transaccion.EspacioId) + ",Id:" + strconv.Itoa(campo.IdCampo)
		if err := request.GetJson(url, &campoExistente); err != nil {
			logs.Error(err)
			panic(err.Error())
		}
		campoExistente[0].Valor = campo.Valor
		campoExistente[0].Activo = true
		campoExistente[0].FechaModificacion = time_bogota.TiempoBogotaFormato()
		var err error
		url = beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico_campo/" + strconv.Itoa(campoExistente[0].Id)
		var respuestaEspacioFisicoCampo map[string]interface{}
		if err = request.SendJson(url, "PUT", &respuestaEspacioFisicoCampo, campoExistente[0]); err != nil {
			logs.Error(err)
			panic(err.Error())
		}
		camposExistentesActivos = append(camposExistentesActivos, campoExistente[0])
	}

	var camposExistentesNoActivos []models.EspacioFisicoCampo
	activosMap := make(map[int]bool)
	for _, activo := range camposExistentesActivos {
		activosMap[activo.Id] = true
	}

	for _, espacio := range camposExistentesEspacio {
		if !activosMap[espacio.Id] {
			camposExistentesNoActivos = append(camposExistentesNoActivos, espacio)
		}
	}

	for _, campo := range camposExistentesNoActivos {
		campo.Activo = false
		campo.FechaModificacion = time_bogota.TiempoBogotaFormato()
		var err error
		url = beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico_campo/" + strconv.Itoa(campo.Id)
		var respuestaEspacioFisicoCampo map[string]interface{}
		if err = request.SendJson(url, "PUT", &respuestaEspacioFisicoCampo, campo); err != nil {
			logs.Error(err)
			panic(err.Error())
		}
	}
	return camposExistentes
}
