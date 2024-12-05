package models

type NuevoEspacioFisico struct {
	CamposDinamicos		[]*CamposDinamicos
	DependenciaPadre	int
	EspacioFisico		*EspacioFisico
	TipoEdificacion		int
	TipoEspacioFisico	int
	TipoTerreno			int
	TipoUso				int
}