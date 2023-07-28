package api

import (
	"fmt"
	"net/http"
	URL "net/url"

	"github.com/go-resty/resty/v2"
	"github.com/ifandonlyif-io/ifandonlyif-backend/token"
	"github.com/labstack/echo/v4"
)

type iffid struct {
	Iffid string `json:"iffid"`
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

// nft godoc
// @Summary      fetchUserIffNfts
// @Description  get USER IffNFTS
// @Tags         fetchUserIffNfts
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /auth/fetchUserIffNfts [POST]
func (server *Server) FetchUserIffNfts(c echo.Context) (err error) {

	payload, ok := c.Get(AuthorizationPayloadKey).(*token.Payload)

	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid type for KEY")
	}

	// Create a Resty Client
	client := resty.New()
	client.Header.Add("accept", "application/json")
	params := URL.Values{}

	params.Set("contractAddresses", server.config.IFFNftContractAddress)

	// set woner wallet address
	params.Set("owner", payload.Wallet)

	//main net
	//reqUrl := "https://eth-mainnet.g.alchemy.com/v2/uLe6RNK4s3INiolh-9N2t9hE2xpO2YGl/getNFTs?" + params.Encode()

	// goerli net
	// reqUrl := "https://eth-goerli.g.alchemy.com/v2/JJqZwPLyThiBz_TowjruMBZWKiL9UIae/getNFTs?" + params.Encode()

	// sepolia net
	reqUrl := server.config.AlchemyApiUrl + "getNFTs?" + params.Encode()

	// request alchemy
	resp, err := client.R().
		EnableTrace().
		Get(reqUrl)
	if err != nil {
		fmt.Println("No response from request")
	}

	return c.JSON(http.StatusOK, resp.String())
}

// nft godoc
// @Summary      getIffNftById
// @Description  fetch limited IffNft
// @Tags         getIffNftById
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /getIffNftById [GET]
func (server *Server) FetchIffNftById(c echo.Context) (err error) {

	var p iffid

	if err := (&echo.DefaultBinder{}).BindBody(c, &p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Create a Resty Client
	client := resty.New()
	client.Header.Add("accept", "application/json")
	params := URL.Values{}

	params.Set("contractAddresses", server.config.IFFNftContractAddress)
	params.Set("withMetadata", "true")
	params.Set("startToken", p.Iffid)
	params.Set("limit", "1")
	//main net
	//reqUrl := "https://eth-mainnet.g.alchemy.com/v2/uLe6RNK4s3INiolh-9N2t9hE2xpO2YGl/getNFTs?" + params.Encode()

	// goerli net
	// reqUrl := "https://eth-goerli.g.alchemy.com/v2/JJqZwPLyThiBz_TowjruMBZWKiL9UIae/getNFTs?" + params.Encode()

	// sepolia net

	reqUrl := server.config.AlchemyApiUrl + "getNFTsForCollection?" + params.Encode()

	// request alchemy
	resp, err := client.R().
		EnableTrace().
		Get(reqUrl)
	if err != nil {
		fmt.Println("No response from request")
	}

	return c.JSON(http.StatusOK, resp.String())
}

// Todo list nft counts
