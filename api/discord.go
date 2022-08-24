package api

import (
	db "github.com/ifandonlyif-io/ifandonlyif-backend/db/sqlc"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (server *Server) report(e echo.Context) error {
	blocklist, err := server.store.CreateReportBlocklist(e.Request().Context(), db.CreateReportBlocklistParams{
		HttpAddress: e.FormValue("url"),
	})

	if err != nil {
		return e.JSON(http.StatusInternalServerError, err)
	}
	return e.JSON(http.StatusCreated, blocklist)
}
