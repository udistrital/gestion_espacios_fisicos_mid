package models


type Creaciones struct {
	EspacioFisicoId							int
	AsignacionEspacioFisicoDependenciaId	int
	TipoUsoEspacioFisico					int
	CamposId								[]int
	EspacioFisicoCampoId					[]int
}