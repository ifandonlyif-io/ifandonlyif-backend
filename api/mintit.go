package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	URL "net/url"

	"github.com/go-resty/resty/v2"
	_ "github.com/ifandonlyif-io/ifandonlyif-backend/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/ifandonlyif-io/ifandonlyif-backend/token"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

type EthUsdPrice struct {
	Usd float64 `json:"USD"`
}

// mintit godoc
// @Summary      fetchUserNfts
// @Description  get mintable USER NFTS
// @Tags         fetchUserNfts
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /auth/fetchUserNfts [POST]
func (server *Server) FetchUserNfts(c echo.Context) (err error) {

	payload, ok := c.Get(AuthorizationPayloadKey).(*token.Payload)

	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid type for KEY")
	}

	nftprojs, err := server.store.ListNftProjects(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// Create a Resty Client
	client := resty.New()
	client.Header.Add("accept", "application/json")
	params := URL.Values{}

	// iterate array
	for i := range nftprojs {
		params.Add("contractAddresses[]", nftprojs[i].ContractAddress)
	}

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

// Eth to Usdt Price Api for UI usage
// https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD
// https://docs.alchemy.com/reference/alchemy-gettokenbalances

// ethToUsd godoc
// @Summary      ethToUsd
// @Description  get current 1 ETH to USD price
// @Tags         ethToUsd
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /ethToUsd [GET]
func (server *Server) EthToUsd(c echo.Context) (err error) {

	// perpare the request
	params := URL.Values{}
	params.Set("fsym", "ETH")
	params.Set("tsyms", "USD")

	// set the pricing api
	reqUrl := "https://min-api.cryptocompare.com/data/price?" + params.Encode()

	// Create a Resty Client
	client := resty.New()

	// request for current ETH:USD price
	resp, err := client.R().
		EnableTrace().
		Get(reqUrl)
	if err != nil {
		fmt.Println("No response from request")
	}
	var eg EthUsdPrice

	// parse usd json
	if err := json.Unmarshal(resp.Body(), &eg); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	resultEthUsd := decimal.NewFromFloat(eg.Usd)

	return c.JSON(http.StatusOK, resultEthUsd)
}
