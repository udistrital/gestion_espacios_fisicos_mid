package services_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"bou.ke/monkey"
	"github.com/astaxie/beego"
	"github.com/udistrital/espacios_fisicos_mid/models"
	"github.com/udistrital/espacios_fisicos_mid/services"
	"github.com/udistrital/utils_oas/request"
)

func TestEditarEspacioFisico(t *testing.T) {
    t.Log("//////////////////////////////////")
    t.Log("Inicio PutActivarEspacioFisico")
    t.Log("//////////////////////////////////")

    t.Run("Caso 1: Activacion exitosa de espacio fisico", func(t *testing.T) {
        transaccion := &models.EditarEspaciosFisicos{
            EspacioId: 1810,
			Nombre: "ESPACIO DE PRUEBA",
			Descripcion: "Espacio para test",
			CodAbreviacion: "Test1",
			DependenciaId: 22,
			TipoEspacioId: 2,
			TipoUsoId: 2,
			TipoEdificacion: 3,
			TipoTerreno: 3,
			CamposExistentes: &[]models.CamposEspacioFisico{{
				IdCampo: 180,
				Valor: "50",
				Existente: true,
			}},
			CamposNoExistentes: &[]models.CamposEspacioFisico{{
				IdCampo: 2,
				Valor: "100",
				Existente: true,
			}},
        }

		espacioFisicoOriginal := models.EspacioFisico{
			Id: 1810,
			Nombre: "Espacio Modificado",
			Descripcion: "Espacio Modificado test",
		}

		tipoUsoOriginal := []models.TipoUsoEspacioFisico{{
			Id: 150,
			TipoUsoId: &models.TipoUso{
				Id: 2,
			},
			EspacioFisicoId: &models.EspacioFisico{
				Id: 1810,
			},
			Activo: true,
			FechaCreacion: "2024-10-30T15:40:12.909843Z",
			FechaModificacion: "2024-10-30T15:40:12.909843Z",
			},
			{Id: 50,
			TipoUsoId: &models.TipoUso{
				Id: 2,
			},
			EspacioFisicoId: &models.EspacioFisico{
				Id: 1810,
			},
			Activo: false,
			FechaCreacion: "2024-10-30T15:40:12.909843Z",
			FechaModificacion: "2024-10-30T15:40:12.909843Z",
		}}

		dependenciasEspacioOriginal := []models.AsignacionEspacioFisicoDependencia{{
			Id: 140,
			EspacioFisicoId: &models.EspacioFisico{
				Id: 1810,
				Nombre: "Espacio fisico original",
			},
			DependenciaId: &models.Dependencia{
				Id: 22,
				Nombre: "Dependencia Original",
			},
			Activo: true,
		},{
			Id: 141,
			EspacioFisicoId: &models.EspacioFisico{
				Id: 1810,
				Nombre: "Espacio fisico original",
			},
			DependenciaId: &models.Dependencia{
				Id: 25,
				Nombre: "Dependencia Original",
			},
			Activo: false,
		}}

		objetosOriginales := models.ObjetosSinEditar{
			EspacioFisico: &models.EspacioFisico{
				Id: 1810,
				Nombre: "Espacio de prueba",
			},
			TipoUsoEspacioFisicoActivo: &tipoUsoOriginal[0],
			TipoUsoEspacioFisicoNoActivo: &tipoUsoOriginal[1],
			NuevoTipoUsoEspacioFisico: 4,
			AsignacionEspacioFisicoDependenciaActivo: &dependenciasEspacioOriginal[0],
			AsignacionEspacioFisicoDependenciaNoActivo: &dependenciasEspacioOriginal[1],
			NuevaAsignacionEspacioFisicoDependencia: 26,
		}

		monkey.Patch(request.GetJson, func(url string, target interface{}) error {
			switch url{
			case  beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico/" + strconv.Itoa(transaccion.EspacioId):
				*target.(*models.EspacioFisico) = espacioFisicoOriginal
				return nil
			case beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico?query=EspacioFisicoId.Id:" + strconv.Itoa(transaccion.EspacioId):
				*target.(*[]models.TipoUsoEspacioFisico) = tipoUsoOriginal
				return nil
			
			case beego.AppConfig.String("OikosCrudUrl") + "asignacion_espacio_fisico_dependencia?query=EspacioFisicoId.Id:" + strconv.Itoa(transaccion.EspacioId):
				*target.(*[]models.AsignacionEspacioFisicoDependencia) = dependenciasEspacioOriginal
				return nil
			}
            return errors.New("URL no esperada")
        })
        defer monkey.Unpatch(request.GetJson)

		monkey.Patch(services.ExisteEnListaTipoUso, func(tipos []models.TipoUsoEspacioFisico, id int) bool {
			return true
		})
		defer monkey.Unpatch(services.ExisteEnListaTipoUso)

		monkey.Patch(services.ExisteEnListaDependencia, func(tipos []models.AsignacionEspacioFisicoDependencia, id int) bool {
			return true
		})
		defer monkey.Unpatch(services.ExisteEnListaDependencia)

		monkey.Patch(services.ActualizarEspacioFisico, func(espacioFisicoOriginal models.EspacioFisico, transaccion *models.EditarEspaciosFisicos) models.EspacioFisico {
			return models.EspacioFisico{
				Id: 1810,
				Nombre: "Espacio Modificado",
				Descripcion: "Espacio Modificado test",
			}
		})
		defer monkey.Unpatch(services.ActualizarEspacioFisico)

		monkey.Patch(services.ActualizarNuevoTipoUso, func(espacioModificado *models.EspacioFisico, transaccion *models.EditarEspaciosFisicos, objetosOriginales *models.ObjetosSinEditar) {

		})
		defer monkey.Unpatch(services.ActualizarNuevoTipoUso)

		monkey.Patch(services.ActualizarTipoUsoExistente, func(espacioModificado *models.EspacioFisico, transaccion *models.EditarEspaciosFisicos, objetosOriginales *models.ObjetosSinEditar) {
		})
		defer monkey.Unpatch(services.ActualizarTipoUsoExistente)

		monkey.Patch(services.ActualizarNuevaDependencia, func(espacioModificado *models.EspacioFisico, transaccion *models.EditarEspaciosFisicos, objetosOriginales *models.ObjetosSinEditar) {

		})
		defer monkey.Unpatch(services.ActualizarNuevaDependencia)

		monkey.Patch(services.ActualizarNuevaDependenciaExistente, func(espacioModificado *models.EspacioFisico, transaccion *models.EditarEspaciosFisicos, objetosOriginales *models.ObjetosSinEditar)  {
		})
		defer monkey.Unpatch(services.ActualizarNuevaDependenciaExistente)

		monkey.Patch(services.ActualizarCampos, func( transaccion *models.EditarEspaciosFisicos, objetosOriginales *models.ObjetosSinEditar) {
		})
		defer monkey.Unpatch(services.ActualizarCampos)

		monkey.Patch(services.AgregarCampos, func(espacioModificado *models.EspacioFisico, transaccion *models.EditarEspaciosFisicos, objetosOriginales *models.ObjetosSinEditar) {
		})
		defer monkey.Unpatch(services.AgregarCampos)

		monkey.Patch(services.EliminarCampos, func( transaccion *models.EditarEspaciosFisicos, objetosOriginales *models.ObjetosSinEditar) {
		})
		defer monkey.Unpatch(services.EliminarCampos)

		fmt.Println("ESPACIO FISICO")
		fmt.Println(objetosOriginales.EspacioFisico)
		fmt.Println("TIPO DE USO ACTIVO")
		fmt.Println(objetosOriginales.TipoUsoEspacioFisicoActivo)
		fmt.Println("TIPO DE USO NO ACTIVO")
		fmt.Println(objetosOriginales.TipoUsoEspacioFisicoNoActivo)
		fmt.Println("NUEVO TIPO USO")
		fmt.Println(objetosOriginales.NuevoTipoUsoEspacioFisico)
		fmt.Println("ASIGNACION ACTIVA")
		fmt.Println(objetosOriginales.AsignacionEspacioFisicoDependenciaActivo)
		fmt.Println("ASISGNACION NO ACTIVA")
		fmt.Println(objetosOriginales.AsignacionEspacioFisicoDependenciaNoActivo)
		fmt.Println("NUEVA ASIGNACION")
		fmt.Println(objetosOriginales.NuevaAsignacionEspacioFisicoDependencia)
		fmt.Println("CAMPOS ACTIVOS")

		alerta, outputError := services.EditarEspacioFisico(transaccion, transaccion.EspacioId)

        if len(alerta) == 0 || alerta[0] != "Success" {
            t.Errorf("Se esperaba una alerta con 'Success', pero se obtuvo: %v", alerta)
        }

        if outputError != nil {
            t.Errorf("Se esperaba outputError nulo, pero se obtuvo: %v", outputError)
        }

        t.Log("Editar de espacio fisico ejecutado exitosamente sin errores")
    })
  
}

