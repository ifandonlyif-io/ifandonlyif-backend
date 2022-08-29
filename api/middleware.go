package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// AuthMiddleware creates an echo middleware for authorization
func (server *Server) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		authorizationHeader := ctx.Request().Header
		if len(authorizationHeader) == 0 {
			fmt.Println("authorization header is not provided")
			err := errors.New("authorization header is not provided")
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		fields := strings.Fields(authorizationHeader.Get(authorizationHeaderKey))

		if len(fields) < 2 {
			fmt.Println("invalid authorization header format")
			err := errors.New("invalid authorization header format")
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			fmt.Println("unsupported authorization type")
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		accessToken := fields[1]
		payload, err := server.tokenMaker.VerifyToken(accessToken)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		ctx.Response().WriteHeader(http.StatusOK)
		return next(ctx)
	}
}
