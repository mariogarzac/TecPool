package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mariogarzac/tecpool/pkg/handlers"
)

func serveAndRoute(){

    // echo instance
    e := echo.New()

    // middleware
    e.Use(middleware.Recover())
    e.Use(middleware.Logger())

    // route
    e.GET("/",handlers.Repo.Login)

    // start server
    e.Start(portNumber)
}
