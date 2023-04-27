package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string `json:"accessToken"`
	AccessTokenExpiresAt int64  `json:"accessTokenExpiresAt"`
}

// token godoc
// @Summary      renewAccess
// @Description  renewAccess
// @Tags         renewAccess
// @Accept       json
// @produce application/json
// @param refreshToken body string true "refreshToken"
// @Success      200  {object}  renewAccessTokenResponse
// @Success      201  {string}  StatusOK
// @Failure      400  {string}  StatusBadRequest
// @Failure      404  {string}  StatusNotFound
// @Failure      500  {string}  StatusInternalServerError
// @Router       /renewAccess [POST]
func (server *Server) renewAccessToken(ctx echo.Context) (errEcho error) {
	var req renewAccessTokenRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	refreshPayload, err := server.refreshTokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	// session management
	// session, err := server.store.GetSession(ctx.Request().Context(), refreshPayload.ID)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return echo.NewHTTPError(http.StatusNotFound, err)
	// 	}
	// 	return echo.NewHTTPError(http.StatusInternalServerError, err)
	// }

	// if session.IsBlocked {
	// 	err := fmt.Errorf("blocked session")
	// 	return echo.NewHTTPError(http.StatusUnauthorized, err)
	// }

	// if session.Wallet != refreshPayload.Wallet {
	// 	err := fmt.Errorf("incorrect session wallet")

	// 	return echo.NewHTTPError(http.StatusUnauthorized, err)
	// }

	// if session.RefreshToken != req.RefreshToken {
	// 	err := fmt.Errorf("mismatched session token")
	// 	return echo.NewHTTPError(http.StatusUnauthorized, err)
	// }

	// if time.Now().After(session.ExpiresAt) {
	// 	err := fmt.Errorf("expired session")
	// 	return echo.NewHTTPError(http.StatusUnauthorized, err)
	// }

	accessToken, accessPayload, err := server.accessTokenMaker.CreateToken(
		refreshPayload.UserName,
		refreshPayload.Wallet,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	return ctx.JSON(http.StatusOK, rsp)
}
