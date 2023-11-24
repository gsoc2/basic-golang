package aliaspackagename

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type handler struct {}

func newHandler() *handler {
	return &handler{}
}

func (h *handler) health(c echo.Context) error {
	return errors.WithStack(c.JSONBlob(http.StatusOK, resp))
}
