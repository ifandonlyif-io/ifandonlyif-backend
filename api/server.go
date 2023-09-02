package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	db "github.com/ifandonlyif-io/ifandonlyif-backend/db/sqlc"
	_ "github.com/ifandonlyif-io/ifandonlyif-backend/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/ifandonlyif-io/ifandonlyif-backend/token"
	"github.com/ifandonlyif-io/ifandonlyif-backend/util"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Server serves HTTP requests for our nft-platform service.
type Server struct {
	config            util.Config
	store             db.Store
	Echo              *echo.Echo
	accessTokenMaker  token.Maker
	refreshTokenMaker token.Maker
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	accessTokenMaker, accessErr := token.NewJWTMaker(config.AccessTokenSymmetricKey)
	if accessErr != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", accessErr)
	}

	refreshTokenMaker, refreshErr := token.NewJWTMaker(config.RefreshTokenSymmetricKey)
	if refreshErr != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", refreshErr)
	}

	server := &Server{
		config:            config,
		store:             store,
		accessTokenMaker:  accessTokenMaker,
		refreshTokenMaker: refreshTokenMaker,
	}
	server.setupRouter()
	server.RunCronFetchGas()
	return server, nil
}

func (server *Server) setupRouter() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS()) // dev setting for allow any origin

	// production setting
	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"http://localhost:3001", "http://219.84.184.170", "http://219.85.184.145"},
	// 	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	// }))

	api := e.Group("/api", middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:api-key",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == "valid-key", nil
		},
	}))

	adminApi := e.Group("/api/admin", middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:token",
		Validator: func(key string, c echo.Context) (bool, error) {
			_, err := server.store.GetUserIdByToken(c.Request().Context(), key)
			if err != nil {
				return false, err
			}
			return true, nil
		},
	}))

	// Routes
	auth := e.Group("/auth", server.AuthMiddleware)
	e.GET("/gasInfo", server.GasHandler)
	e.POST("/code", server.NonceHandler)
	e.POST("/login", server.LoginHandler)
	e.POST("/renewAccess", server.renewAccessToken)
	e.GET("/healthz", HealthCheck)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/nftProjects", server.ListNftProjects)
	api.POST("discord/report", server.report)
	api.POST("discord/apply", server.apply)
	api.GET("/discord/appliances", server.appliance)
	api.GET("/discord/channels", server.channels)
	api.PATCH("/discord/approve/:id", server.approve)
	api.PATCH("/discord/channel/lock/:id", server.lockChannel)
	api.PATCH("/discord/channel/unlock/:id", server.UnlockChannel)
	e.POST("/checkUri", server.GetBlocklistByUri)
	e.GET("/ethToUsd", server.EthToUsd)
	e.POST("/fetchIffNftById", server.FetchIffNftById)
	e.GET("/getIffNftMeta", server.FetchIffNftMeta)
	e.POST("/checkExistBlocklists", server.CheckExistBlocklists)
	e.POST("/checkBlocklists", server.CheckBlocklists)
	e.POST("/checkSpamContract", server.CheckSpamContract)
	// JWT - Authentication Middleware
	auth.POST("/fetchUserNft", server.FetchUserNfts)
	auth.POST("/fetchUserIffNft", server.FetchUserIffNfts)
	auth.POST("/fetchMinterIffNfts", server.fetchNftsByMinterAddress)
	// Key - Authentication Middleware
	api.GET("/getAllBlockLists", server.GetAllBlockLists)
	api.GET("/listDisprovedBlocklists", server.ListDisprovedBlocklists)
	api.GET("/listVerifiedBlocklists", server.ListVerifiedBlocklists)
	api.GET("/listUnreviewedBlocklists", server.ListUnreviewedBlocklists)
	api.POST("/fetchBlockListById", server.GetBlockListById)
	api.POST("/disproveBlocklist", server.DisproveBlocklist)
	api.POST("/verifyBlocklist", server.VerifyBlocklist)

	e.POST("api/admin/login", server.AdminLogin)
	adminApi.GET("/me", server.checkLoginStatus)
	adminApi.GET("/user-management/users", server.GetAllUsers)
	adminApi.POST("/user-management/user", server.NewUser)
	adminApi.DELETE("/user-management/user/:id", server.DeleteUser)

	server.Echo = e
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.Echo.Start(address)
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags health
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "Server is up and running",
	})
}
