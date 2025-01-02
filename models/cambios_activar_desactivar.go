package models


type CambiosActivarDesactivar struct {
	IdEspacioFisico 	int
	IdAsignacion		AsignacionEspacioFisicoDependencia
	IdTipoUso			TipoUsoEspacioFisico
	IdsCampos			[]EspacioFisicoCampo
}