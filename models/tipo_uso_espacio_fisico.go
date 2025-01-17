package models


type TipoUsoEspacioFisico struct {
	Id					int
	TipoUsoId			*TipoUso
	EspacioFisicoId		*EspacioFisico
	Activo				bool
	FechaCreacion		string
	FechaModificacion	string
}