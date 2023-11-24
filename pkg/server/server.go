package server

import (
	"fmt"
	"net/http"

	apn "github.com/gsoc2/basic-golang/pkg/aliaspackagename"
	"github.com/gsoc2/basic-golang/pkg/logger"
	"github.com/gsoc2/basic-golang/pkg/regularpackagename"
	otherpackagename "github.com/gsoc2/basic-golang/pkg/wrongpackagename"
	"github.com/labstack/echo"
	elog "github.com/labstack/gommon/log"
)

type SomeType string

type SomeStruct struct{}

type SomeInterface interface{}

const SomeConst = "ConstValue"

var SomeVar SomeType = "VarValue"

const (
	ConstGroup1 = 1
	ConstGroup2 = 2
)

var somePrivateVar = 1

func New(port int) (*http.Server, error) {
	e := echo.New()

	e.Logger.SetLevel(elog.OFF)

	e.Use(logger.Middleware())

	regularpackagename.RegisterRoutes(e)
	otherpackagename.RegisterRoutes(e)
	apn.RegisterRoutes(e)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: e,
	}

	return srv, nil
}
