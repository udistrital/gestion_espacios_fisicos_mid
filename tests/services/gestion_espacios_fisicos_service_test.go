package services_test

import (
	"errors"
	"testing"
	"bou.ke/monkey"
	"strconv"
	"github.com/astaxie/beego"
	"github.com/udistrital/gestion_espacios_fisicos_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/gestion_espacios_fisicos_mid/services"
)

func TestRegistrarEspacioFisico(t *testing.T) {
    t.Log("//////////////////////////////////")
    t.Log("Inicio RegistrarEspacioFisico")
    t.Log("//////////////////////////////////")

    t.Run("Caso 1: Registro exitoso de espacio fisico", func(t *testing.T) {
        transaccion := &models.NuevoEspacioFisico{
            CamposExistentes: []*models.CamposEspacioFisico{
				{
					IdCampo:   1,
					Valor:     "50",
					Existente: true,
				},
				{
					IdCampo:   2,
					Valor:     "30",
					Existente: true,
				},
			},
            DependenciaPadre: 5,
			EspacioFisico: &models.EspacioFisico{
				Nombre: "NUEVO ESPACIO PRUEBA CAMPOS 8", 
				Descripcion: "descripcion", 
				CodigoAbreviacion: "codigo",
			},
			TipoEdificacion: 2,
			TipoTerreno: 2,
			TipoUso: 2,
        }

        monkey.Patch(request.GetJson, func(url string, target interface{}) error {
			switch url{
			case  beego.AppConfig.String("OikosCrudUrl") + "tipo_espacio_fisico/" + strconv.Itoa(transaccion.TipoEspacioFisico):
				*target.(*models.TipoEspacioFisico) = models.TipoEspacioFisico{Id: transaccion.TipoEspacioFisico}
                return nil
			
			case beego.AppConfig.String("OikosCrudUrl") + "dependencia/" + strconv.Itoa(transaccion.DependenciaPadre):
				*target.(*models.Dependencia) = models.Dependencia{Id: transaccion.DependenciaPadre}
				return nil
			
			case beego.AppConfig.String("OikosCrudUrl") + "tipo_uso/" + strconv.Itoa(transaccion.TipoUso):
				*target.(*models.TipoUso) = models.TipoUso{Id: transaccion.TipoUso}
				return nil
			}
            return errors.New("URL no esperada")
        })
        defer monkey.Unpatch(request.GetJson)

        monkey.Patch(services.CrearEspacioFisico, func(transaccion *models.NuevoEspacioFisico, tipoEspacio models.TipoEspacioFisico, creaciones *models.Creaciones) models.EspacioFisico {
			return models.EspacioFisico{} 
		})
        defer monkey.Unpatch(services.CrearEspacioFisico)

        monkey.Patch(services.CrearAsignacionEspacioFisicoDependencia, func(transaccion *models.NuevoEspacioFisico, dependencia models.Dependencia, espacioFisico models.EspacioFisico, creaciones *models.Creaciones) {
            
        })
        defer monkey.Unpatch(services.CrearAsignacionEspacioFisicoDependencia)

		monkey.Patch(services.CrearTipoUsoEspacioFisico, func(transaccion *models.NuevoEspacioFisico, tipoUso models.TipoUso , espacioFisico models.EspacioFisico, creaciones *models.Creaciones) {
            
        })
        defer monkey.Unpatch(services.CrearTipoUsoEspacioFisico)

		monkey.Patch(services.CrearEspacioFisicoCampo, func(transaccion *models.NuevoEspacioFisico, espacioFisico models.EspacioFisico, creaciones *models.Creaciones) []models.EspacioFisicoCampo {
			return []models.EspacioFisicoCampo{}
		})		
        defer monkey.Unpatch(services.CrearEspacioFisicoCampo)

        alerta, outputError := services.RegistrarEspacioFisico(transaccion)

        if len(alerta) == 0 || alerta[0] != "Success" {
            t.Errorf("Se esperaba una alerta con 'Success', pero se obtuvo: %v", alerta)
        }

        if outputError != nil {
            t.Errorf("Se esperaba outputError nulo, pero se obtuvo: %v", outputError)
        }

        t.Log("Registro de espacio fisico ejecutado exitosamente sin errores")
    })
  
}