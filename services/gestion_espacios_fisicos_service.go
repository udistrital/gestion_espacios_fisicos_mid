package services

import (
	"fmt"

	"github.com/udistrital/gestion_espacios_fisicos_mid/models"
)

func BuscarEspacioFisico(transaccion *models.BusquedaEspacioFisico) (resultadoBusqueda []models.RespuestaBusquedaEspaciosFisicos, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "BuscarEspacioFisico", "err": err, "status": "500"}
			panic(outputError)
		}
	}()


	return resultadoBusqueda, outputError
}

func EditarDependencia(transaccion *models.EditarEspaciosFisicos) (alerta []string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "BuscarEspacioFisico", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	fmt.Println("BUSCA ESPACIO FISICO")
	return alerta, outputError
}
