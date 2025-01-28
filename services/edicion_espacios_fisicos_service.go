package services

import (
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
	if err := request.GetJson(url, &espacioFisicoOriginal); err != nil || espacioFisicoOriginal.Id == 0 {
		logs.Error(err)
		panic(err.Error())
	}

	var tiposUsoOriginal []models.TipoUsoEspacioFisico
	url = beego.AppConfig.String("OikosCrudUrl") + TIPO_USO_ESPACIO_ID + strconv.Itoa(transaccion.EspacioId)
	if err := request.GetJson(url, &tiposUsoOriginal); err != nil{
		logs.Error(err)
		panic(err.Error())
	}


	var dependenciasEspacioOriginal []models.AsignacionEspacioFisicoDependencia
	url = beego.AppConfig.String("OikosCrudUrl") + ASIGNACION_ESPACIO_ID + strconv.Itoa(transaccion.EspacioId)
	if err := request.GetJson(url, &dependenciasEspacioOriginal); err != nil{
		logs.Error(err)
		panic(err.Error())
	}




	var espacioModificado = ActualizarEspacioFisico(espacioFisicoOriginal, transaccion)

	if existeEnListaTipoUso(tiposUsoOriginal, transaccion.TipoUsoId) {
		ActualizarTipoUsoExistente(&espacioModificado, transaccion)
	} else {
		ActualizarNuevoTipoUso(&espacioModificado, transaccion)
	}

	if existeEnListaDependencia(dependenciasEspacioOriginal, transaccion.DependenciaId) {
		ActualizarNuevaDependenciaExistente(&espacioModificado, transaccion)
	} else {
		ActualizarNuevaDependencia(&espacioModificado, transaccion)
	}


	if len(*transaccion.CamposExistentes) > 0 {
		ActualizarCampos(transaccion)
	}
	if len(*transaccion.CamposNoExistentes) > 0 {
		AgregarCampos(&espacioModificado, transaccion)
	}


	return alerta, outputError
}

func existeEnListaTipoUso(tipos []models.TipoUsoEspacioFisico, id int) bool {
	for _, tipo := range tipos {
		if tipo.TipoUsoId.Id == id {
			return true
		}
	}
	return false
}

func existeEnListaDependencia(dependencias []models.AsignacionEspacioFisicoDependencia, id int) bool {
	for _, dependencia := range dependencias {
		if dependencia.DependenciaId.Id == id {
			return true
		}
	}
	return false
}

func ActualizarEspacioFisico(espacioFisicoOriginal models.EspacioFisico, transaccion *models.EditarEspaciosFisicos) (espacioModificado models.EspacioFisico) {
	espacioModificado = espacioFisicoOriginal
	espacioModificado.Nombre = transaccion.Nombre
	espacioModificado.CodigoAbreviacion = transaccion.CodAbreviacion
	espacioModificado.Descripcion = transaccion.Descripcion
	espacioModificado.TipoEdificacionId = transaccion.TipoEdificacion
	espacioModificado.TipoTerrenoId = transaccion.TipoTerreno
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

func ActualizarNuevoTipoUso(espacioModificado *models.EspacioFisico, transaccion *models.EditarEspaciosFisicos) {
	var nuevoTipoUso models.TipoUso
	url := beego.AppConfig.String("OikosCrudUrl") + "tipo_uso/" + strconv.Itoa(transaccion.TipoUsoId)
	if err := request.GetJson(url, &nuevoTipoUso); err != nil || nuevoTipoUso.Id == 0 {
		logs.Error(err)
		panic(err.Error())
	}
	var tipoUsoOriginal []models.TipoUsoEspacioFisico
	url = beego.AppConfig.String("OikosCrudUrl") + TIPO_USO_ESPACIO_ID + strconv.Itoa(transaccion.EspacioId) + ACTIVO_URL
	if err := request.GetJson(url, &tipoUsoOriginal); err != nil {
		logs.Error(err)
		panic(err.Error())
	}
	if len(tipoUsoOriginal) > 0 {
		tipoUsoOriginal[0].Activo = false
		tipoUsoOriginal[0].FechaModificacion = time_bogota.TiempoBogotaFormato()
		var err error
		url = beego.AppConfig.String("OikosCrudUrl") + TIPO_USO_ESPACIO_URL + strconv.Itoa(tipoUsoOriginal[0].Id)
		var respuestaTipoUsoOriginal map[string]interface{}
		if err = request.SendJson(url, "PUT", &respuestaTipoUsoOriginal, tipoUsoOriginal[0]); err != nil {
			logs.Error(err)
			panic(err.Error())
		}
	}
	var tipoUsoModificado models.TipoUsoEspacioFisico
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
}

func ActualizarTipoUsoExistente(espacioModificado *models.EspacioFisico, transaccion *models.EditarEspaciosFisicos) {

	var tipoUsoOriginal []models.TipoUsoEspacioFisico
	url := beego.AppConfig.String("OikosCrudUrl") + TIPO_USO_ESPACIO_ID + strconv.Itoa(transaccion.EspacioId) + ACTIVO_URL
	if err := request.GetJson(url, &tipoUsoOriginal); err != nil{
		logs.Error(err)
		panic(err.Error())
	}
	if len(tipoUsoOriginal) > 0{
		if tipoUsoOriginal[0].TipoUsoId.Id != transaccion.TipoUsoId {
			tipoUsoOriginal[0].Activo = false
			tipoUsoOriginal[0].FechaModificacion = time_bogota.TiempoBogotaFormato()
			var err error
			url = beego.AppConfig.String("OikosCrudUrl") + TIPO_USO_ESPACIO_URL + strconv.Itoa(tipoUsoOriginal[0].Id)
			var respuestaTipoUsoOriginal map[string]interface{}
			if err = request.SendJson(url, "PUT", &respuestaTipoUsoOriginal, tipoUsoOriginal[0]); err != nil {
				logs.Error(err)
				panic(err.Error())
			}
			ActivarTipoUso(transaccion)
		}
	}else{
		ActivarTipoUso(transaccion)
	}
}

func ActivarTipoUso(transaccion *models.EditarEspaciosFisicos) { 
	var nuevoTipoUsoActivo []models.TipoUsoEspacioFisico
	url := beego.AppConfig.String("OikosCrudUrl") + TIPO_USO_ESPACIO_ID + strconv.Itoa(transaccion.EspacioId) + ",TipoUsoId.Id:" + strconv.Itoa(transaccion.TipoUsoId)
	if err := request.GetJson(url, &nuevoTipoUsoActivo); err != nil || nuevoTipoUsoActivo[0].Id == 0 {
		logs.Error(err)
		panic(err.Error())
	}
	var tipoUsoModificado models.TipoUsoEspacioFisico
	tipoUsoModificado = nuevoTipoUsoActivo[0]
	tipoUsoModificado.Activo = true
	tipoUsoModificado.FechaModificacion = time_bogota.TiempoBogotaFormato()
	url = beego.AppConfig.String("OikosCrudUrl") + TIPO_USO_ESPACIO_URL + strconv.Itoa(tipoUsoModificado.Id)
	var respuestaTipoUsoExistente map[string]interface{}
	if err := request.SendJson(url, "PUT", &respuestaTipoUsoExistente, tipoUsoModificado); err != nil {
		logs.Error(err)
		panic(err.Error())
	}
}

func ActualizarNuevaDependencia(espacioModificado *models.EspacioFisico, transaccion *models.EditarEspaciosFisicos) {
	var nuevaDependencia models.Dependencia
	url := beego.AppConfig.String("OikosCrudUrl") + "dependencia/" + strconv.Itoa(transaccion.DependenciaId)
	if err := request.GetJson(url, &nuevaDependencia); err != nil || nuevaDependencia.Id == 0 {
		logs.Error(err)
		panic(err.Error())
	}

	var dependenciaEspacioOriginal []models.AsignacionEspacioFisicoDependencia
	url = beego.AppConfig.String("OikosCrudUrl") + ASIGNACION_ESPACIO_ID + strconv.Itoa(transaccion.EspacioId) + ACTIVO_URL
	if err := request.GetJson(url, &dependenciaEspacioOriginal); err != nil{
		logs.Error(err)
		panic(err.Error())
	}
	if len(dependenciaEspacioOriginal)> 0{
		dependenciaEspacioOriginal[0].Activo = false
		dependenciaEspacioOriginal[0].FechaModificacion = time_bogota.TiempoBogotaFormato()
		var err error
		url = beego.AppConfig.String("OikosCrudUrl") + ASIGNACION_ESPACIO_URL + strconv.Itoa(dependenciaEspacioOriginal[0].Id) + ACTIVO_URL
		var respuestaDependenciaEspacioOriginal map[string]interface{}
		if err = request.SendJson(url, "PUT", &respuestaDependenciaEspacioOriginal, dependenciaEspacioOriginal[0]); err != nil {
			logs.Error(err)
			panic(err.Error())
		}
	}
	var dependenciaEspacioModificado models.AsignacionEspacioFisicoDependencia
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
}

func ActualizarNuevaDependenciaExistente(espacioModificado *models.EspacioFisico, transaccion *models.EditarEspaciosFisicos) {
	var dependenciaEspacioOriginal []models.AsignacionEspacioFisicoDependencia
	url := beego.AppConfig.String("OikosCrudUrl") + ASIGNACION_ESPACIO_ID + strconv.Itoa(transaccion.EspacioId) + ACTIVO_URL
	if err := request.GetJson(url, &dependenciaEspacioOriginal); err != nil{
		logs.Error(err)
		panic(err.Error())
	}
	if len(dependenciaEspacioOriginal) > 0{
		if dependenciaEspacioOriginal[0].DependenciaId.Id != transaccion.DependenciaId {
			dependenciaEspacioOriginal[0].Activo = false
			dependenciaEspacioOriginal[0].FechaModificacion = time_bogota.TiempoBogotaFormato()
			var err error
			url = beego.AppConfig.String("OikosCrudUrl") + ASIGNACION_ESPACIO_URL + strconv.Itoa(dependenciaEspacioOriginal[0].Id)
			var respuestaDependenciaEspacioOriginal map[string]interface{}
			if err = request.SendJson(url, "PUT", &respuestaDependenciaEspacioOriginal, dependenciaEspacioOriginal[0]); err != nil {
				logs.Error(err)
				panic(err.Error())
			}
			ActivarDependenciaEspacioFisico(transaccion)
		}
	}else{
		ActivarDependenciaEspacioFisico(transaccion)
	}
}

func ActivarDependenciaEspacioFisico(transaccion *models.EditarEspaciosFisicos) {
	var nuevaDependenciaEspacioActiva []models.AsignacionEspacioFisicoDependencia
	url := beego.AppConfig.String("OikosCrudUrl") + ASIGNACION_ESPACIO_ID + strconv.Itoa(transaccion.EspacioId) + ",DependenciaId.Id:" + strconv.Itoa(transaccion.DependenciaId)
	if err := request.GetJson(url, &nuevaDependenciaEspacioActiva); err != nil || nuevaDependenciaEspacioActiva[0].Id == 0 {
		logs.Error(err)
		panic(err.Error())
	}
	var dependenciaEspacioModificado models.AsignacionEspacioFisicoDependencia
	dependenciaEspacioModificado = nuevaDependenciaEspacioActiva[0]

	dependenciaEspacioModificado.Activo = true
	dependenciaEspacioModificado.FechaModificacion = time_bogota.TiempoBogotaFormato()

	url = beego.AppConfig.String("OikosCrudUrl") + ASIGNACION_ESPACIO_URL + strconv.Itoa(dependenciaEspacioModificado.Id)
	var respuestaDependenciaEspacioExistente map[string]interface{}
	if err := request.SendJson(url, "PUT", &respuestaDependenciaEspacioExistente, dependenciaEspacioModificado); err != nil {
		logs.Error(err)
		panic(err.Error())
	}
}

func AgregarCampos(espacioModificado *models.EspacioFisico, transaccion *models.EditarEspaciosFisicos) {
	for _, campo := range *transaccion.CamposNoExistentes {
		var campoExistente []models.EspacioFisicoCampo
		url := beego.AppConfig.String("OikosCrudUrl") + ESPACIO_FISICO_CAMPO_ID + strconv.Itoa(transaccion.EspacioId) + ",CampoId.Id:" + strconv.Itoa(campo.IdCampo)
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
		}
	}
}

func CambiarEstadoCampoEspacio(campoExistente models.EspacioFisicoCampo, campo models.CamposEspacioFisico) {
	campoExistente.Activo = true
	campoExistente.Valor = campo.Valor
	campoExistente.FechaModificacion = time_bogota.TiempoBogotaFormato()
	var err error
	url := beego.AppConfig.String("OikosCrudUrl") + ESPACIO_FISICO_CAMPO_URL + strconv.Itoa(campoExistente.Id)
	var respuestaEspacioFisicoCampo map[string]interface{}
	if err = request.SendJson(url, "PUT", &respuestaEspacioFisicoCampo, campoExistente); err != nil {
		logs.Error(err)
		panic(err.Error())
	}
}

func ActualizarCampos(transaccion *models.EditarEspaciosFisicos) {

	var camposExistentesEspacio []models.EspacioFisicoCampo
	url := beego.AppConfig.String("OikosCrudUrl") + ESPACIO_FISICO_CAMPO_ID + strconv.Itoa(transaccion.EspacioId)
	if err := request.GetJson(url, &camposExistentesEspacio); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	var camposExistentesActivos []models.EspacioFisicoCampo
	for _, campo := range *transaccion.CamposExistentes {
		var campoExistente []models.EspacioFisicoCampo
		url := beego.AppConfig.String("OikosCrudUrl") + ESPACIO_FISICO_CAMPO_ID + strconv.Itoa(transaccion.EspacioId) + ",Id:" + strconv.Itoa(campo.IdCampo)
		if err := request.GetJson(url, &campoExistente); err != nil {
			logs.Error(err)
			panic(err.Error())
		}
		campoExistente[0].Valor = campo.Valor
		campoExistente[0].Activo = true
		campoExistente[0].FechaModificacion = time_bogota.TiempoBogotaFormato()
		var err error
		url = beego.AppConfig.String("OikosCrudUrl") + ESPACIO_FISICO_CAMPO_URL + strconv.Itoa(campoExistente[0].Id)
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
		url = beego.AppConfig.String("OikosCrudUrl") + ESPACIO_FISICO_CAMPO_URL + strconv.Itoa(campo.Id)
		var respuestaEspacioFisicoCampo map[string]interface{}
		if err = request.SendJson(url, "PUT", &respuestaEspacioFisicoCampo, campo); err != nil {
			logs.Error(err)
			panic(err.Error())
		}
	}
}
