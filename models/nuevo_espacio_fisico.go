package models

type NuevoEspacioFisico struct {
	CamposExistentes	[]*CamposEspacioFisico
	DependenciaPadre	int
	EspacioFisico		*EspacioFisico
	TipoEdificacion		int
	TipoEspacioFisico	int
	TipoTerreno			int
	TipoUso				int
}