package api

import (
	"net/http"
	"strconv"
	"time"

	db "github.com/ifandonlyif-io/ifandonlyif-backend/db/sqlc"
	"github.com/labstack/echo/v4"
)

func (server *Server) getNFTs(c echo.Context) error {
	nfts, err := server.store.ListIffNfts(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nfts)
}

func (server Server) createNFT(e echo.Context) error {
	projectId := e.FormValue("projectId")
	parseInt, _ := strconv.ParseInt(projectId, 10, 64)

	created, err := server.store.CreateIffNft(e.Request().Context(), db.CreateIffNftParams{
		ProjectID:                  parseInt,
		UserWalletAddress:          e.FormValue("wallet"),
		NftProjectsContractAddress: e.FormValue("projectContractAddress"),
		NftProjectsCollectionName:  e.FormValue("collectionName"),
		MintDate:                   time.Now(),
		MintTransaction:            e.FormValue("mintTransaction"),
	})
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err)
	}
	return e.JSON(http.StatusCreated, created)
}

// nft godoc
// @Summary      nftProjects
// @Description  fetch limited nft projects
// @Tags         nftProjects
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /nftProjects [GET]
func (server *Server) ListNftProjects(c echo.Context) error {
	nftprojs, err := server.store.ListNftProjects(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nftprojs)
}
