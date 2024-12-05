package models


type EspacioFisico struct {
	Id						int
	Nombre			 		string
	Descripcion				string
	CodigoAbreviacion		string
	TipoTerrenoId			int
	TipoEdificacionId		int
	TipoEspacioFisicoId 	*TipoEspacioFisico
	Activo					bool
	FechaCreacion			string
	FechaModificacion		string
}