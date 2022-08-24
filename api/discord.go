package api

import (
	db "github.com/ifandonlyif-io/ifandonlyif-backend/db/sqlc"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (server *Server) report(e echo.Context) error {
	blocklists, _ := server.store.ListReportBlocklists(e.Request().Context())
	for i := range blocklists {
		if blocklists[i].HttpAddress == e.FormValue("url") {
			return e.JSON(http.StatusOK, "this project had been reported")
		}
	}

	blocklist, err := server.store.CreateReportBlocklist(e.Request().Context(), db.CreateReportBlocklistParams{
		HttpAddress: e.FormValue("url"),
	})

	if err != nil {
		return e.JSON(http.StatusInternalServerError, err)
	}
	return e.JSON(http.StatusCreated, blocklist)
}

func (server *Server) getReportNFTs(e echo.Context) error {
	blocklists, _ := server.store.ListReportBlocklists(e.Request().Context())
	return e.JSON(http.StatusOK, blocklists)
}
