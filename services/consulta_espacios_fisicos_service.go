package services

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/espacios_fisicos_mid/models"
	"github.com/udistrital/utils_oas/request"
)

func BuscarEspacioFisico(transaccion *models.BusquedaEspacioFisico) (resultadoBusqueda []models.RespuestaBusquedaEspacioFisico, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "BuscarEspacioFisico", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	var urlBusquedaTipo string
	var urlBusquedaDependencia string
	var urlBusquedaNombre string
	var urlBusquedaTipoEspacio string

	if transaccion.TipoUsoId != 0 {
		urlBusquedaTipo = beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico?limit=-1&query=TipoUsoId:" + strconv.Itoa(transaccion.TipoUsoId)
	}

	if transaccion.DependenciaId != 0 {
		urlBusquedaDependencia = beego.AppConfig.String("OikosCrudUrl") + "asignacion_espacio_fisico_dependencia?limit=-1&query=DependenciaId.Id:" + strconv.Itoa(transaccion.DependenciaId)
	}

	if transaccion.TipoEspacioFisicoId != 0 {
		if urlBusquedaTipo != "" {
			urlBusquedaTipo += ",EspacioFisicoId.TipoEspacioFisicoId.Id:" + strconv.Itoa(transaccion.TipoEspacioFisicoId)
		} else if urlBusquedaDependencia != "" {
			urlBusquedaDependencia += ",EspacioFisicoId.TipoEspacioFisicoId.Id:" + strconv.Itoa(transaccion.TipoEspacioFisicoId)
		} else {
			urlBusquedaTipoEspacio = beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico?limit=-1&query=TipoEspacioFisicoId.Id:" + strconv.Itoa(transaccion.TipoEspacioFisicoId)
		}
	}

	if transaccion.NombreEspacioFisico != "" {
		if urlBusquedaTipo != "" {
			urlBusquedaTipo += ",EspacioFisicoId.Nombre:" + url.QueryEscape(transaccion.NombreEspacioFisico)
		} else if urlBusquedaDependencia != "" {
			urlBusquedaDependencia += ",EspacioFisicoId.Nombre:" + url.QueryEscape(transaccion.NombreEspacioFisico)
		} else if urlBusquedaTipoEspacio != "" {
			urlBusquedaTipoEspacio += ",Nombre:" + url.QueryEscape(transaccion.NombreEspacioFisico)
		} else {
			urlBusquedaNombre = beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico?limit=-1&query=Nombre:" + url.QueryEscape(transaccion.NombreEspacioFisico)
		}
	}

	var idsTipo, idsDependencia, idNombre, idTipoEspacio []models.EspacioFisico
	listasNoVacias := 0
	if urlBusquedaTipo != ""{
		// urlBusquedaTipo += ",Activo:true"
		idsTipo = BusquedaTipo(urlBusquedaTipo)
		listasNoVacias++
	}
	if urlBusquedaDependencia != ""{
		// urlBusquedaDependencia += ",Activo:true"
		idsDependencia = BusquedaDepependencia(urlBusquedaDependencia)
		listasNoVacias++
	}
	if urlBusquedaNombre != ""{
		// urlBusquedaNombre += ",Activo:true"
		idNombre = BusquedaNombre(urlBusquedaNombre)
		listasNoVacias++
	}
	if urlBusquedaTipoEspacio != ""{
		// urlBusquedaTipoEspacio += ",Activo:true"
		idTipoEspacio = BusquedaNombre(urlBusquedaTipoEspacio)
		listasNoVacias++
	}

	contador := make(map[int]int)
	espaciosMap := make(map[int]models.EspacioFisico)

	listas := [][]models.EspacioFisico{idsTipo, idsDependencia, idNombre, idTipoEspacio}

	for _, lista := range listas {
		if len(lista) > 0 {
			for _, espacio := range lista {
				encontrado := false
				if existente, ok := espaciosMap[espacio.Id]; ok {
					if EspacioFisicoIgual(espacio, existente) {
						contador[espacio.Id]++
						encontrado = true
					}
				}
				if !encontrado {
					contador[espacio.Id] = 1
					espaciosMap[espacio.Id] = espacio
				}
			}
		}
	}

	fmt.Println(idsTipo)
	fmt.Println(idsDependencia)

	var repetidos []models.EspacioFisico
	for id, count := range contador {
		if count == listasNoVacias {
			repetidos = append(repetidos, espaciosMap[id])
		}
	}

	for _, id := range repetidos {
		var resultado = CrearRespuestaBusqueda(id)
		resultadoBusqueda = append(resultadoBusqueda, resultado)
	}

	return resultadoBusqueda, outputError
}

func EspacioFisicoIgual(a, b models.EspacioFisico) bool {
	return a.Id == b.Id
}

func BusquedaTipo(url string) (ids []models.EspacioFisico) {
	var respuesta []models.TipoUsoEspacioFisico
	fmt.Println(url)
	if err := request.GetJson(url, &respuesta); err != nil {
		logs.Error(err)
		panic(err.Error())
	}
	for _, espacio := range respuesta {
		ids = append(ids, *espacio.EspacioFisicoId)
	}
	return ids
}

func BusquedaDepependencia(url string) (ids []models.EspacioFisico) {
	var respuesta []models.AsignacionEspacioFisicoDependencia
	if err := request.GetJson(url, &respuesta); err != nil {
		logs.Error(err)
		panic(err.Error())
	}
	for _, espacio := range respuesta {
		ids = append(ids, *espacio.EspacioFisicoId)
	}
	return ids
}

func BusquedaNombre(url string) (ids []models.EspacioFisico) {
	var respuesta []models.EspacioFisico
	if err := request.GetJson(url, &respuesta); err != nil {
		logs.Error(err)
		panic(err.Error())
	}
	for _, espacio := range respuesta {
		ids = append(ids, espacio)
	}
	return ids
}

func CrearRespuestaBusqueda(id models.EspacioFisico) models.RespuestaBusquedaEspacioFisico {
	var resultado models.RespuestaBusquedaEspacioFisico

	resultado.EspacioFisico = &id
	resultado.TipoEspacioFisico = id.TipoEspacioFisicoId
	var tipoUsoEspacioFisico []models.TipoUsoEspacioFisico
	url := beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico?limit=-1&query=EspacioFisicoId.id:" + strconv.Itoa(id.Id) + ",Activo:true"
	if err := request.GetJson(url, &tipoUsoEspacioFisico); err != nil {
		logs.Error(err)
		panic(err.Error())
	}
	if len(tipoUsoEspacioFisico) > 0 {
		resultado.TipoUso = tipoUsoEspacioFisico[0].TipoUsoId
	} else {
		resultado.TipoUso = &models.TipoUso{}
	}

	// var camposEspacioFisico []models.EspacioFisicoCampo
	// url = beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico_campo?limit=-1&query=EspacioFisicoId.Id:" + strconv.Itoa(id.Id)
	// if err := request.GetJson(url, &camposEspacioFisico); err != nil {
	// 	logs.Error(err)
	// 	panic(err.Error())
	// }
	// if len(camposEspacioFisico) > 0{
	// 	var campos []models.Campo
	// 	for _, campo := range camposEspacioFisico{
	// 		campos = append(campos, *campo.CampoId)
	// 	}
	// 	resultado.Campos = &campos
	// }else{
	// 	resultado.Campos = &[]models.Campo{}
	// }

	// var dependenciaEspacioFisico []models.AsignacionEspacioFisicoDependencia
	// url = beego.AppConfig.String("OikosCrudUrl") + "asignacion_espacio_fisico_dependencia?limit=-1&query=EspacioFisicoId.id:" + strconv.Itoa(id.Id)
	// if err := request.GetJson(url, &dependenciaEspacioFisico); err != nil {
	// 	logs.Error(err)
	// 	panic(err.Error())
	// }
	// if len(dependenciaEspacioFisico)>0{
	// 	resultado.Dependencia = dependenciaEspacioFisico[0].DependenciaId
	// }else{
	// 	resultado.Dependencia = &models.Dependencia{}
	// }

	return resultado
}
