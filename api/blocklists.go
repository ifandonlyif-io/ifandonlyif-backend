package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type GetListPayload struct {
	UUID string `json:"uuid"`
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

	queryUuid, errUuid := uuid.Parse(p.UUID)
	if errUuid != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errUuid)
	}

	blocklist, errGetList := server.store.GetReportBlocklist(c.Request().Context(), queryUuid)

	if errGetList != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errGetList)
	}

	return c.JSON(http.StatusOK, blocklist)
}

// blocklists godoc
// @Summary      checkUri
// @Description  fetch blocklist by uri
// @Tags         checkUri
// @param uri body string true "uri"
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /checkUri [POST]
func (server *Server) GetBlocklistByUri(c echo.Context) (err error) {

	var u CheckUriPayload

	if errPayload := (&echo.DefaultBinder{}).BindBody(c, &u); errPayload != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errPayload)
	}

	test, errGetList := server.store.GetBlocklistByUri(c.Request().Context(), u.Uri)

	fmt.Print(test)
	if errGetList != nil {

		if errGetList == sql.ErrNoRows {
			return c.JSON(http.StatusOK, false)
		} else {
			return echo.NewHTTPError(http.StatusBadRequest, errGetList)
		}
	}

	return c.JSON(http.StatusOK, true)
}
