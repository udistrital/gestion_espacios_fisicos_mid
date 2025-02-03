package models

type ObjetosSinEditar struct {
	EspacioFisico									*EspacioFisico
	TipoUsoEspacioFisicoActivo 						*TipoUsoEspacioFisico
	TipoUsoEspacioFisicoNoActivo					*TipoUsoEspacioFisico
	NuevoTipoUsoEspacioFisico						int
	AsignacionEspacioFisicoDependenciaActivo 		*AsignacionEspacioFisicoDependencia
	AsignacionEspacioFisicoDependenciaNoActivo		*AsignacionEspacioFisicoDependencia
	NuevaAsignacionEspacioFisicoDependencia			int
	CamposNoActivos 								[]EspacioFisicoCampo
	CamposActivos									[]EspacioFisicoCampo
	NuevoCampo										[]int
}