package api

import (
	"fmt"
	"net/http"
	URL "net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/ifandonlyif-io/ifandonlyif-backend/token"
	"github.com/labstack/echo/v4"
)

type IffIdPayload struct {
	IffId string `json:"iffid"`
}

type CheckPayload struct {
	Address string `json:"address"`
}

type Response struct {
	Total       interface{} `json:"total"`
	Page        int         `json:"page"`
	PageSize    int         `json:"page_size"`
	Cursor      interface{} `json:"cursor"`
	Result      []Result    `json:"result"`
	BlockExists bool        `json:"block_exists"`
}

type Result struct {
	TokenID     string `json:"token_id"`
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
}

type Token struct {
	TokenId         string `json:"tokenId"`
	ContractAddress string `json:"contractAddress"`
}

type NFTMetadataRequest struct {
	RefreshCache bool    `json:"refreshCache"`
	Tokens       []Token `json:"tokens"`
}

func (ch CheckPayload) Validate() error {
	if !hexRegex.MatchString(ch.Address) {
		return ErrInvalidAddress
	}
	return nil
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

	params.Set("contractAddresses[]", server.config.IFFNftContractAddress)

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
// @Summary      fetchIffNftById
// @Description  fetch limited IffNft
// @Tags         getIffNftById
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /fetchIffNftById [POST]
func (server *Server) FetchIffNftById(c echo.Context) (err error) {
	var p IffIdPayload

	if err := (&echo.DefaultBinder{}).BindBody(c, &p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidAddress)
	}

	// Create a Resty Client
	client := resty.New()
	client.Header.Add("accept", "application/json")
	params := URL.Values{}

	params.Set("contractAddress", server.config.IFFNftContractAddress)
	params.Set("tokenId", p.IffId)
	params.Set("refreshCache", "false")
	//main net
	//reqUrl := "https://eth-mainnet.g.alchemy.com/v2/uLe6RNK4s3INiolh-9N2t9hE2xpO2YGl/getNFTs?" + params.Encode()

	// goerli net
	// reqUrl := "https://eth-goerli.g.alchemy.com/v2/JJqZwPLyThiBz_TowjruMBZWKiL9UIae/getNFTs?" + params.Encode()

	// sepolia net

	reqUrl := server.config.AlchemyNftApiUrl + "getNFTMetadata?" + params.Encode()

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
// @Summary      getIffNftMeta
// @Description  fetch limited IffNft
// @Tags         getIffNftMeta
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /getIffNftMeta [GET]
func (server *Server) FetchIffNftMeta(c echo.Context) (err error) {

	// Create a Resty Client
	client := resty.New()
	client.Header.Add("accept", "application/json")
	params := URL.Values{}

	params.Set("contractAddress", server.config.IFFNftContractAddress)
	//main net
	//reqUrl := "https://eth-mainnet.g.alchemy.com/v2/uLe6RNK4s3INiolh-9N2t9hE2xpO2YGl/getNFTs?" + params.Encode()

	// goerli net
	// reqUrl := "https://eth-goerli.g.alchemy.com/v2/JJqZwPLyThiBz_TowjruMBZWKiL9UIae/getNFTs?" + params.Encode()

	// sepolia net

	reqUrl := server.config.AlchemyApiUrl + "getContractMetadata?" + params.Encode()

	// request alchemy
	resp, err := client.R().
		EnableTrace().
		Get(reqUrl)
	if err != nil {
		fmt.Println("No response from request")
	}

	return c.JSON(http.StatusOK, resp.String())
}

// checkspamcontract

// nft godoc
// @Summary      checkSpamContract
// @Description  fetch limited IffNft
// @Tags         checkSpamContract
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /checkSpamContract [POST]
func (server *Server) CheckSpamContract(c echo.Context) (err error) {

	var check CheckPayload

	if err := (&echo.DefaultBinder{}).BindBody(c, &check); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// validate wallet address
	err = check.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidAddress)
	}

	// Create a Resty Client
	client := resty.New()
	client.Header.Add("accept", "application/json")
	params := URL.Values{}

	params.Set("contractAddress", check.Address)

	reqUrl := server.config.AlchemyApiUrl + "isSpamContract?" + params.Encode()

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
// @Summary      fetchNftsByMinterAddress
// @Description  fetch IffNftby Minter Address
// @Tags         fetchNftsByMinterAddress
// @Accept */*
// @produce application/json
// @Success      200  {string}  StatusOK
// @Router       /fetchNftsByMinterAddress [POST]
func (server *Server) fetchNftsByMinterAddress(c echo.Context) (err error) {

	payload, ok := c.Get(AuthorizationPayloadKey).(*token.Payload)

	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid type for KEY")
	}

	params := URL.Values{}

	params.Set("chain", server.config.MoralisEthNetwork)
	params.Set("format", "decimal")
	reqUrl := server.config.MoralisApiUrl + server.config.IFFNftContractAddress + "/transfers?" + params.Encode()

	// Create a Resty Client
	client := resty.New()
	client.Header.Add("accept", "application/json")
	client.Header.Add("X-API-Key", server.config.MoralisApiKey)

	// To unmarshal resutls from moralis
	var results Response

	// request moralis
	resp, err := client.R().
		SetResult(&results).
		Get(reqUrl)
	if err != nil {
		fmt.Printf("response failed: %s", err)
	}

	fmt.Scan(resp)

	var minterIffTokens []Token

	// filter token id of user from minting
	for _, s := range results.Result {
		if s.FromAddress == server.config.BlackholeAddress &&
			s.ToAddress == strings.ToLower(payload.Wallet) {
			minterIffTokens = append(minterIffTokens, Token{
				TokenId:         s.TokenID,
				ContractAddress: server.config.IFFNftContractAddress,
			})
		}
	}

	alchemyClient := resty.New()

	requestBody := NFTMetadataRequest{
		RefreshCache: false,
		Tokens:       minterIffTokens,
	}

	alchemyReqUrl := server.config.AlchemyNftApiUrl + "getNFTMetadataBatch"

	alchemyResp, alchemyErr := alchemyClient.R().
		SetBody(requestBody).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		EnableTrace().
		Post(alchemyReqUrl)

	if alchemyErr != nil {
		fmt.Printf("response failed: %s", alchemyErr)
	}

	return c.JSON(http.StatusOK, alchemyResp.String())

}
