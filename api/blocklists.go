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

type VerifyPayload struct {
	UUID string `json:"uuid"`
}

type DisprovePayload struct {
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
	fmt.Println(queryUuid)
	blocklist, errGetList := server.store.GetReportBlocklist(c.Request().Context(), queryUuid)

	if errGetList != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errGetList)
	}

	return c.JSON(http.StatusOK, blocklist)
}

// blocklists godoc
// @Summary      ListVerifiedBlocklists
// @Description  fetch verified blocklists
// @Tags         ListVerifiedBlocklists
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /api/listVerifiedBlocklists [GET]
func (server *Server) ListVerifiedBlocklists(c echo.Context) (err error) {

	blocklists, errGetList := server.store.ListVerifiedBlocklists(c.Request().Context())

	if errGetList != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errGetList)
	}

	return c.JSON(http.StatusOK, blocklists)
}

// blocklists godoc
// @Summary      ListDisprovedBlocklists
// @Description  get all disproved blocklists
// @Tags         ListDisprovedBlocklists
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /api/listDisprovedBlocklists [GET]
func (server *Server) ListDisprovedBlocklists(c echo.Context) (err error) {

	blocklists, errGetList := server.store.ListDisprovedBlocklists(c.Request().Context())

	if errGetList != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errGetList)
	}

	return c.JSON(http.StatusOK, blocklists)
}

// blocklists godoc
// @Summary      VerifyBlocklist
// @Description  verify blocklist by uuid
// @Tags         VerifyBlocklist
// @param uuid body string true "uuid"
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /api/verifyBlocklist [POST]
func (server *Server) VerifyBlocklist(c echo.Context) (err error) {

	var v VerifyPayload

	if errPayload := (&echo.DefaultBinder{}).BindBody(c, &v); errPayload != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errPayload)
	}

	queryUuid, errUuid := uuid.Parse(v.UUID)
	if errUuid != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errUuid)
	}

	blocklist, errGetList := server.store.VerifyBlocklist(c.Request().Context(), queryUuid)

	if errGetList != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errGetList)
	}

	return c.JSON(http.StatusOK, blocklist)
}

// blocklists godoc
// @Summary      DisproveBlocklist
// @Description  disprove blocklist by uuid
// @Tags         DisproveBlocklist
// @param uuid body string true "uuid"
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /api/disproveBlocklist [POST]
func (server *Server) DisproveBlocklist(c echo.Context) (err error) {

	var d DisprovePayload

	if errPayload := (&echo.DefaultBinder{}).BindBody(c, &d); errPayload != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errPayload)
	}

	queryUuid, errUuid := uuid.Parse(d.UUID)
	if errUuid != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errUuid)
	}

	blocklist, errGetList := server.store.DisproveBlocklist(c.Request().Context(), queryUuid)

	if errGetList != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errGetList)
	}

	return c.JSON(http.StatusOK, blocklist)
}

// blocklists godoc
// @Summary      ListUnreviewedBlocklists
// @Description  get all unreviewed blocklists
// @Tags         ListUnreviewedBlocklists
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /api/listUnreviewedBlocklists [GET]
func (server *Server) ListUnreviewedBlocklists(c echo.Context) (err error) {

	blocklists, errGetList := server.store.ListUnreviewedBlocklists(c.Request().Context())

	if errGetList != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errGetList)
	}

	return c.JSON(http.StatusOK, blocklists)
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

	_, errGetList := server.store.GetBlocklistByUri(c.Request().Context(), u.Uri)

	if errGetList != nil {

		if errGetList == sql.ErrNoRows {
			return c.JSON(http.StatusOK, false)
		} else {
			return echo.NewHTTPError(http.StatusBadRequest, errGetList)
		}
	}

	return c.JSON(http.StatusOK, true)
}
