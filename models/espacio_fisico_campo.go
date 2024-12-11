package models


type EspacioFisicoCampo struct {
	Id					int
	Valor			 	string
	EspacioFisicoId		*EspacioFisico
	CampoId				*Campo
	Activo				bool
	FechaInicio			string
	FechaFin			*string
	FechaCreacion		string
	FechaModificacion	string
}