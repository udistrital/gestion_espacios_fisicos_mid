package models


type EditarEspaciosFisicos struct {
	EspacioId			int
	Nombre				string
	Descripcion 		string
	CodAbreviacion 		string
	DependenciaId		int
	TipoEspacioId		int
	TipoUsoId			int
	TipoEdificacion		int
	TipoTerreno			int
	CamposExistentes	*[]CamposEspacioFisico
	CamposNoExistentes	*[]CamposEspacioFisico
}