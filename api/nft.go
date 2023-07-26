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

// ToDo: Fetch IFF NFTs with user wallet
// mintit godoc
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

	params.Set("contractAddresses", "0x507AA149A42012AD74C1E40E076a3f2391E13b61")

	// set woner wallet address
	params.Set("owner", payload.Wallet)

	//main net
	//reqUrl := "https://eth-mainnet.g.alchemy.com/v2/uLe6RNK4s3INiolh-9N2t9hE2xpO2YGl/getNFTs?" + params.Encode()

	// goerli net
	// reqUrl := "https://eth-goerli.g.alchemy.com/v2/JJqZwPLyThiBz_TowjruMBZWKiL9UIae/getNFTs?" + params.Encode()

	// sepolia net

	reqUrl := "https://eth-sepolia.g.alchemy.com/v2/i8RTBcKFG3U1qEUbUprJXDatOggaZxcE/getNFTs?" + params.Encode()

	// request alchemy
	resp, err := client.R().
		EnableTrace().
		Get(reqUrl)
	if err != nil {
		fmt.Println("No response from request")
	}

	return c.JSON(http.StatusOK, resp.String())
}

// ToDo: Fetch ALL IFF NFTs Count: return Number
func (server *Server) FetchIffNftById(c echo.Context) (err error) {

	var p iffid

	if err := (&echo.DefaultBinder{}).BindBody(c, &p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Create a Resty Client
	client := resty.New()
	client.Header.Add("accept", "application/json")
	params := URL.Values{}

	params.Set("contractAddresses", "0x507AA149A42012AD74C1E40E076a3f2391E13b61")
	params.Set("withMetadata", "true")
	params.Set("startToken", p.Iffid)
	params.Set("limit", "1")
	//main net
	//reqUrl := "https://eth-mainnet.g.alchemy.com/v2/uLe6RNK4s3INiolh-9N2t9hE2xpO2YGl/getNFTs?" + params.Encode()

	// goerli net
	// reqUrl := "https://eth-goerli.g.alchemy.com/v2/JJqZwPLyThiBz_TowjruMBZWKiL9UIae/getNFTs?" + params.Encode()

	// sepolia net

	reqUrl := "https://eth-sepolia.g.alchemy.com/v2/i8RTBcKFG3U1qEUbUprJXDatOggaZxcE/getNFTsForCollection?" + params.Encode()

	// request alchemy
	resp, err := client.R().
		EnableTrace().
		Get(reqUrl)
	if err != nil {
		fmt.Println("No response from request")
	}

	return c.JSON(http.StatusOK, resp.String())
}
