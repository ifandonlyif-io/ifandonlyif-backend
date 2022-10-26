package api

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	_ "github.com/ifandonlyif-io/ifandonlyif-backend/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/ifandonlyif-io/ifandonlyif-backend/token"
	"github.com/labstack/echo/v4"
)

type FetchNftPayload struct {
	WalletAddress string `json:"walletAddress"`
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

	// Create a Resty Client
	client := resty.New()
	client.Header.Add("accept", "application/json")
	url1 := "https://eth-mainnet.g.alchemy.com/v2/uLe6RNK4s3INiolh-9N2t9hE2xpO2YGl/getNFTs?owner="

	url2 := "&pageKey=1&pageSize=15&contractAddresses[]=0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D&contractAddresses[]=0xb47e3cd837dDF8e4c57F05d70Ab865de6e193BBB&contractAddresses[]=0x60E4d786628Fea6478F785A6d7e704777c86a7c6&contractAddresses[]=0xa7d8d9ef8D8Ce8992Df33D8b8CF4Aebabd5bD270&contractAddresses[]=0xED5AF388653567Af2F388E6224dC7C4b3241C544&contractAddresses[]=0x34d85c9CDeB23FA97cb08333b511ac86E1C4E258&contractAddresses[]=0x49cf6f5d44e70224e2e23fdcdd2c053f30ada28b&contractAddresses[]=0x23581767a106ae21c074b2276D25e5C3e136a68b&contractAddresses[]=0x8a90CAb2b38dba80c64b7734e58Ee1dB38B8992e&contractAddresses[]=0x7Bd29408f11D2bFC23c34f18275bBf23bB716Bc7&contractAddresses[]=0xba30E5F9Bb24caa003E9f2f0497Ad287FDF95623&contractAddresses[]=0x0Cfb5d82BE2b949e8fa73A656dF91821E2aD99FD&withMetadata=true"
	url := url1 + payload.WalletAddress + url2

	resp, err := client.R().
		EnableTrace().
		Get(url)
	fmt.Println(resp)
	if err != nil {
		fmt.Println("No response from request")
	}

	return c.JSON(http.StatusOK, resp.String())
}
