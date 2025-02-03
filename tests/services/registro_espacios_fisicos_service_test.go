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
	"github.com/udistrital/utils_oas/time_bogota"
)

func TestRegistrarEspacioFisico(t *testing.T) {
    t.Log("//////////////////////////////////")
    t.Log("Inicio RegistrarEspacioFisico")
    t.Log("//////////////////////////////////")

    t.Run("Caso 1: Registro exitoso de espacio fisico", func(t *testing.T) {
		// Transaccion incial
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
            DependenciaPadre: 7,
			EspacioFisico: &models.EspacioFisico{
				Nombre: "NUEVO ESPACIO PRUEBA CAMPOS 8", 
				Descripcion: "descripcion", 
				CodigoAbreviacion: "codigo",
			},
			TipoEdificacion: 2,
			TipoTerreno: 2,
			TipoUso: 2,
			TipoEspacioFisico: 2,
        }

		// RESPUESTAS BUSQUEDAS
		// Tipo espacio
		tipoEspacio := models.TipoEspacioFisico{
			Id: transaccion.TipoEspacioFisico,
			Nombre: "AULAS DE CLASE",
			Descripcion: "AULAS DE CLASE",
			CodigoAbreviacion: "TIPO_2",
			Activo: true,
			FechaCreacion: "2023-06-29 19:19:19.544",
			FechaModificacion: "2023-06-29 19:19:19.544",
		}
		// Dependencia
		dependencia := models.Dependencia{
			Id: transaccion.DependenciaPadre,
			Nombre: "RECTORIA",
			TelefonoDependencia: "3239300 Secretaria: 1001 - Secretaria asesor 1009 - Comunicaciones Institucionales 1909 Asesores: 1024 - 6506 - abogado: 6507",
			CorreoElectronico: "rectoria@udistrital.edu.co",
			Activo: true,
			FechaCreacion: "2022-07-06 19:45:10.995513",
			FechaModificacion: "2022-07-06 19:45:10.995513",
		}
		// Tipo Uso
		tipoUso := models.TipoUso{
			Id: transaccion.TipoUso,
			Nombre: "Administrativo",
			Descripcion: "",
			CodigoAbreviacion: "",
			Activo: true,
			FechaCreacion: "2022-07-06T19:34:24.454448Z",
			FechaModificacion: "2022-07-06T19:34:24.454448Z",
		}

        monkey.Patch(request.GetJson, func(url string, target interface{}) error {
			switch url{
			case  beego.AppConfig.String("OikosCrudUrl") + "tipo_espacio_fisico/" + strconv.Itoa(transaccion.TipoEspacioFisico):
				*target.(*models.TipoEspacioFisico) = tipoEspacio
				return nil
			case beego.AppConfig.String("OikosCrudUrl") + "dependencia/" + strconv.Itoa(transaccion.DependenciaPadre):
				*target.(*models.Dependencia) = dependencia
				return nil
			
			case beego.AppConfig.String("OikosCrudUrl") + "tipo_uso/" + strconv.Itoa(transaccion.TipoUso):
				*target.(*models.TipoUso) = tipoUso
				return nil
			}
            return errors.New("URL no esperada")
        })
        defer monkey.Unpatch(request.GetJson)

		// CREACIONES
		// Espacio Fisico
		nuevoEspacioFisico := models.EspacioFisico{
			Id: 1852,
			Nombre: transaccion.EspacioFisico.Nombre,
			Descripcion: transaccion.EspacioFisico.Descripcion,
			CodigoAbreviacion: transaccion.EspacioFisico.CodigoAbreviacion,
			TipoTerrenoId: transaccion.TipoTerreno,
			TipoEdificacionId: transaccion.TipoEdificacion,
			TipoEspacioFisicoId: &tipoEspacio,
			Activo: true,
			FechaCreacion: time_bogota.TiempoBogotaFormato(),
			FechaModificacion: time_bogota.TiempoBogotaFormato(),
		}
		// Asignacion Espacio Fisico Dependencia
		nuevoEspacioFisicoDependencia := models.AsignacionEspacioFisicoDependencia{
			Id: 1557,
			EspacioFisicoId: &nuevoEspacioFisico,
			DependenciaId: &dependencia,
			DocumentoSoporte: 0,
			FechaInicio: time_bogota.TiempoBogotaFormato(),
			FechaFin: time_bogota.TiempoBogotaFormato(),
			Activo: true,
			FechaCreacion: time_bogota.TiempoBogotaFormato(),
			FechaModificacion: time_bogota.TiempoBogotaFormato(),
		}
		fmt.Println(nuevoEspacioFisicoDependencia)

        monkey.Patch(services.CrearEspacioFisico, func(transaccion *models.NuevoEspacioFisico, tipoEspacio models.TipoEspacioFisico, creaciones *models.Creaciones) models.EspacioFisico {
			return nuevoEspacioFisico
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