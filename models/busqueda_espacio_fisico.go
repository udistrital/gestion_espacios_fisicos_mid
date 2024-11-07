package models


type BusquedaEspacioFisico struct {
	NombreEspacioFisico		string
	TipoEspacioFisicoId		int
	TipoUsoId				int
	DependenciaId			int
	BusquedaEstado			*Estado
}