package models

type RespuestaBusquedaEspacioFisico struct {
	EspacioFisico		*EspacioFisico
	TipoEspacioFisico   *TipoEspacioFisico
    TipoUso				*TipoUso
}