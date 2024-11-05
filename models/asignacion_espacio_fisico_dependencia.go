package models


type AsignacionEspacioFisicoDependencia struct {
	Id					int
	EspacioFisicoId		*EspacioFisico
	DependenciaId		*Dependencia
	DocumentoSoporte	int
	Activo				bool
	FechaInicio			string
	FechaFin			string
	FechaCreacion		string
	FechaModificacion	string
}