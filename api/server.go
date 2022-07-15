package api

import (
	"fmt"
	"net/http"

	db "github.com/ifandonlyif-io/ifandonlyif-backend/db/sqlc"
	"github.com/ifandonlyif-io/ifandonlyif-backend/util"
	"github.com/labstack/echo/v4"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config util.Config
	store  db.Store
	Echo   *echo.Echo
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	server := &Server{
		config: config,
		store:  store,
	}

	server.setupRouter()
	fmt.Println("NewServer")
	return server, nil
}

func (server *Server) setupRouter() {
	e := echo.New()
	e.POST("/createUser", func(c echo.Context) error {

		createUser, err := server.store.CreateUser(c.Request().Context(), db.CreateUserParams{
			FullName:      c.FormValue("FullName"),
			WalletAddress: c.FormValue("WalletAddress"),
			CountryCode:   c.FormValue("CountryCode"),
			EmailAddress:  c.FormValue("EmailAddress"),
			TwitterName:   c.FormValue("TwitterName"),
			ImageUri:      c.FormValue("ImageUri"),
		})
		if err != nil {
			return err
		}
		fmt.Println("setupRouter")
		return c.JSON(http.StatusCreated, createUser)

	})

	server.Echo = e
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.Echo.Start(address)
}
