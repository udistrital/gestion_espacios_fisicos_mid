package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/espacios_fisicos_mid/controllers:GestionEspaciosFisicosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/espacios_fisicos_mid/controllers:GestionEspaciosFisicosController"],
        beego.ControllerComments{
            Method: "PutActivarEspacioFisico",
            Router: "/ActivarEspacioFisico/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/espacios_fisicos_mid/controllers:GestionEspaciosFisicosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/espacios_fisicos_mid/controllers:GestionEspaciosFisicosController"],
        beego.ControllerComments{
            Method: "PutDesactivarEspacioFisico",
            Router: "/DesactivarEspacioFisico/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/espacios_fisicos_mid/controllers:GestionEspaciosFisicosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/espacios_fisicos_mid/controllers:GestionEspaciosFisicosController"],
        beego.ControllerComments{
            Method: "EditarEspacioFisico",
            Router: "/EditarEspacioFisico/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/espacios_fisicos_mid/controllers:RegistroEspaciosFisicosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/espacios_fisicos_mid/controllers:RegistroEspaciosFisicosController"],
        beego.ControllerComments{
            Method: "RegistroEspacioFisicos",
            Router: "/RegistroEspacioFisico",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
