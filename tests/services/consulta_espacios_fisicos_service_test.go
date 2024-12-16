package services_test

import (
	"fmt"
	"testing"

	"bou.ke/monkey"
	"github.com/udistrital/espacios_fisicos_mid/models"
	"github.com/udistrital/espacios_fisicos_mid/services"
)

func TestBuscarEspacioFisico(t *testing.T) {
    t.Log("//////////////////////////////////")
    t.Log("Inicio BuscarEspacioFisico")
    t.Log("//////////////////////////////////")

    t.Run("Caso 1: Busqueda exitosa de espacios fisicos", func(t *testing.T) {
        transaccion := &models.BusquedaEspacioFisico{
            NombreEspacioFisico: "Espacio fisico",
			TipoEspacioFisicoId: 2,
			TipoUsoId: 2,
			DependenciaId: 22,
			BusquedaEstado: &models.Estado{
				Estado: true,
			},
        }

		monkey.Patch(services.BusquedaTipo, func(url string) []models.EspacioFisico{
			return []models.EspacioFisico{}
		})
		defer monkey.Unpatch(services.BusquedaTipo)

		monkey.Patch(services.BusquedaDepependencia, func(url string) []models.EspacioFisico{
			return []models.EspacioFisico{}
		})
		defer monkey.Unpatch(services.BusquedaDepependencia)

		monkey.Patch(services.BusquedaNombre, func(url string) []models.EspacioFisico{
			return []models.EspacioFisico{}
		})
		defer monkey.Unpatch(services.BusquedaNombre)

		monkey.Patch(services.EspacioFisicoIgual, func(a, b models.EspacioFisico) bool{
			return a.Id == b.Id
		})
		defer monkey.Unpatch(services.EspacioFisicoIgual)

		monkey.Patch(services.CrearRespuestaBusqueda, func(id models.EspacioFisico) models.RespuestaBusquedaEspacioFisico{
			return models.RespuestaBusquedaEspacioFisico{}
		})
		defer monkey.Unpatch(services.CrearRespuestaBusqueda)


        resultadoBusqueda, outputError := services.BuscarEspacioFisico(transaccion)

        // if len(resultadoBusqueda) == 0 {
        //     t.Errorf("Se esperaba una lista con los resultados, pero no se obtuvo")
        // }

        if outputError != nil {
            t.Errorf("Se esperaba outputError nulo, pero se obtuvo: %v", outputError)
        }

        t.Log("Busqueda de espacios fisicos ejecutada exitosamente sin errores")
		fmt.Println(resultadoBusqueda)
    })
  
}