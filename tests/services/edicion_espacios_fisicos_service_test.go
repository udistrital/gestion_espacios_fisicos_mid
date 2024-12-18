package services_test

import (
	"testing"
	"errors"
	"strconv"
	"bou.ke/monkey"
	"github.com/udistrital/espacios_fisicos_mid/models"
	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/espacios_fisicos_mid/services"
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
			TipoEdificacion: "3",
			TipoTerreno: "3",
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
				Id: 5,
			},
			Activo: true,
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
		}}

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

		monkey.Patch(services.ActualizarEspacioFisico, func(espacioFisicoOriginal models.EspacioFisico, transaccion *models.EditarEspaciosFisicos) models.EspacioFisico {
			return models.EspacioFisico{
				Id: 1810,
				Nombre: "Espacio Modificado",
				Descripcion: "Espacio Modificado test",
			}
		})
		defer monkey.Unpatch(services.ActualizarEspacioFisico)

		monkey.Patch(services.ActualizarNuevoTipoUso, func(espacioModificado *models.EspacioFisico, transaccion *models.EditarEspaciosFisicos) models.TipoUsoEspacioFisico {
			return models.TipoUsoEspacioFisico{
				Id: 141,
				EspacioFisicoId: &models.EspacioFisico{
					Id: 1810,
					Nombre: "Espacio fisico Modificado",
				},
				TipoUsoId: &models.TipoUso{
					Id: 2,
					Nombre: "Tipo Uso Modificado",
				},
			}
		})
		defer monkey.Unpatch(services.ActualizarNuevoTipoUso)

		monkey.Patch(services.ActualizarTipoUsoExistente, func(espacioModificado *models.EspacioFisico, transaccion *models.EditarEspaciosFisicos) models.TipoUsoEspacioFisico {
			return models.TipoUsoEspacioFisico{
				Id: 142,
				EspacioFisicoId: &models.EspacioFisico{
					Id: 1810,
					Nombre: "Espacio fisico Modificado",
				},
				TipoUsoId: &models.TipoUso{
					Id: 3,
					Nombre: "Tipo Uso Existente anterior",
				},
			}
		})
		defer monkey.Unpatch(services.ActualizarTipoUsoExistente)

		monkey.Patch(services.ActualizarNuevaDependencia, func(espacioModificado *models.EspacioFisico, transaccion *models.EditarEspaciosFisicos) models.AsignacionEspacioFisicoDependencia {
			return models.AsignacionEspacioFisicoDependencia{
				Id: 140,
				EspacioFisicoId: &models.EspacioFisico{
					Id: 1810,
					Nombre: "Espacio fisico original",
				},
				DependenciaId: &models.Dependencia{
					Id: 22,
					Nombre: "Dependencia Modificada",
				},
			}
		})
		defer monkey.Unpatch(services.ActualizarNuevaDependencia)

		monkey.Patch(services.ActualizarNuevaDependenciaExistente, func(espacioModificado *models.EspacioFisico, transaccion *models.EditarEspaciosFisicos) models.AsignacionEspacioFisicoDependencia {
			return models.AsignacionEspacioFisicoDependencia{
				Id: 140,
				EspacioFisicoId: &models.EspacioFisico{
					Id: 1810,
					Nombre: "Espacio fisico original",
				},
				DependenciaId: &models.Dependencia{
					Id: 22,
					Nombre: "Dependencia Existente anterior",
				},
			}
		})
		defer monkey.Unpatch(services.ActualizarNuevaDependenciaExistente)

		monkey.Patch(services.ActualizarCampos, func( transaccion *models.EditarEspaciosFisicos) []models.EspacioFisicoCampo {
			return []models.EspacioFisicoCampo{
				
			}
		})
		defer monkey.Unpatch(services.ActualizarCampos)

		monkey.Patch(services.AgregarCampos, func(espacioModificado *models.EspacioFisico, transaccion *models.EditarEspaciosFisicos) []models.EspacioFisicoCampo {
			return []models.EspacioFisicoCampo{
				
			}
		})
		defer monkey.Unpatch(services.AgregarCampos)


        alerta, outputError := services.EditarEspacioFisico(transaccion)

        if len(alerta) == 0 || alerta[0] != "Success" {
            t.Errorf("Se esperaba una alerta con 'Success', pero se obtuvo: %v", alerta)
        }

        if outputError != nil {
            t.Errorf("Se esperaba outputError nulo, pero se obtuvo: %v", outputError)
        }

        t.Log("Editar de espacio fisico ejecutado exitosamente sin errores")
    })
  
}

