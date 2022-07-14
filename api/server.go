package api

import (
	db "github.com/ifandonlyif-io/ifandonlyif-backend/db/sqlc"
	"github.com/ifandonlyif-io/ifandonlyif-backend/util"
	echo "github.com/labstack/echo/v4"
)

// Server serves HTTP requests for our iff service.
type Server struct {
	config util.Config
	store  db.Store
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {

	// server := &Server{
	// 	config: config,
	// 	store:  store,
	// }

	server := echo.New()

	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	// do something here
	// }

	return server, nil
}
