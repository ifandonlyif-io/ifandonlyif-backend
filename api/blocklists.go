package api

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type GetListPayload struct {
	UUID uuid.UUID `json:"uuid"`
}

type CheckUriPayload struct {
	Uri string `json:"uri"`
}

// blocklists godoc
// @Summary      getAllBlockLists
// @Description  get all blocklists
// @Tags         getAllBlockLists
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /api/getAllBlockLists [GET]
func (server *Server) GetAllBlockLists(c echo.Context) (err error) {
	blocklists, err := server.store.ListReportBlocklists(c.Request().Context())

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, blocklists)
}

// blocklists godoc
// @Summary      fetchBlockListById
// @Description  fetch blocklist by uuid
// @Tags         fetchBlockListById
// @param uuid body string true "uuid"
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /api/fetchBlockListById [POST]
func (server *Server) GetBlockListById(c echo.Context) (err error) {

	var p GetListPayload

	if errPayload := (&echo.DefaultBinder{}).BindBody(c, &p); errPayload != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errPayload)
	}

	blocklist, errGetList := server.store.GetReportBlocklist(c.Request().Context(), p.UUID)

	if errGetList != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errGetList)
	}

	return c.JSON(http.StatusOK, blocklist)
}

// blocklists godoc
// @Summary      fetchBlockListByUri
// @Description  fetch blocklist by uri
// @Tags         fetchBlockListByUri
// @param uri body string true "uri"
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /fetchBlockListByUri [POST]
func (server *Server) GetBlocklistByUri(c echo.Context) (err error) {

	var u CheckUriPayload

	if errPayload := (&echo.DefaultBinder{}).BindBody(c, &u); errPayload != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errPayload)
	}
	fmt.Println(u.Uri)
	blocklist, errGetList := server.store.GetBlocklistByUri(c.Request().Context(), u.Uri)

	if errGetList != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errGetList)
	}

	return c.JSON(http.StatusOK, blocklist)
}
