package api

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"math/big"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	db "github.com/ifandonlyif-io/ifandonlyif-backend/db/sqlc"
	"github.com/ifandonlyif-io/ifandonlyif-backend/token"
	"github.com/labstack/echo/v4"
)

type User struct {
	WalletAddress string `json:"walletAddress"`
	Nonce         string `json:"nonce"`
}

type code struct {
	Code string `json:"code"`
}

type accessToken struct {
	AccessToken string `json:"accessToken"`
}

type RegisterPayload struct {
	WalletAddress string `json:"walletAddress"`
}

type SigninPayload struct {
	WalletAddress string `json:"address"`
	Signature     string `json:"Signature"`
}

func (p RegisterPayload) Validate() error {
	if !hexRegex.MatchString(p.WalletAddress) {
		return ErrInvalidAddress
	}
	return nil
}

func (s SigninPayload) Validate() error {
	if !hexRegex.MatchString(s.WalletAddress) {
		return ErrInvalidAddress
	}
	if len(s.Signature) == 0 {
		return ErrMissingSig
	}
	return nil
}

var (
	max  *big.Int
	once sync.Once
)

var (
	hexRegex   *regexp.Regexp = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)
	nonceRegex *regexp.Regexp = regexp.MustCompile(`^[0-9]+$`)
)

var (
	ErrUserNotExists  = errors.New("user does not exist")
	ErrUserExists     = errors.New("user already exists")
	ErrInvalidAddress = errors.New("invalid address")
	ErrInvalidNonce   = errors.New("invalid nonce")
	ErrMissingSig     = errors.New("signature is missing")
	ErrAuthError      = errors.New("authentication error")
)

func GetNonce() (string, error) {
	once.Do(func() {
		max = new(big.Int)
		max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))
	})
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return n.Text(10), nil
}

// register godoc
// @Summary      register
// @Description  register a new user
// @Tags         register
// @Accept       json
// @produce application/json
// @param walletAddress body string true "walletAddress"
// @Success      200  {string}  StatusOK
// @Success      201  {string}  StatusOK
// @Failure      400  {string}  StatusBadRequest
// @Failure      404  {string}  StatusNotFound
// @Failure      500  {string}  StatusInternalServerError
// @Router       /register [POST]
func (server *Server) NonceHandler(c echo.Context) (err error) {
	var p RegisterPayload

	if err := (&echo.DefaultBinder{}).BindBody(c, &p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = p.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidAddress)
	}

	user, err := server.store.GetUserByWalletAddress(c.Request().Context(), sql.NullString{String: p.WalletAddress, Valid: true})
	// return (echo.NewHTTPError(http.StatusInternalServerError, user))
	if err != nil && err != sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	// if err == sql.ErrNoRows {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, ErrUserNotExists)
	// }

	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, err)
	// }

	resCode := &code{
		Code: user.Nonce.String,
	}

	if len(user.Nonce.String) > 0 {
		return c.JSON(http.StatusFound, resCode)
	}

	nonce, err := GetNonce()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidNonce)
	}

	createUser, err := server.store.CreateUser(c.Request().Context(), db.CreateUserParams{
		WalletAddress: sql.NullString{String: p.WalletAddress, Valid: true},
		Nonce:         sql.NullString{String: nonce, Valid: true},
	})

	// return echo.NewHTTPError(http.StatusAccepted, createUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, ErrUserExists)
	}

	resCode = &code{
		Code: createUser.Nonce.String,
	}

	return c.JSON(http.StatusCreated, resCode)
}

// login godoc
// @Summary      login
// @Description  login
// @Tags         login
// @Accept       json
// @produce application/json
// @param walletAddress body string true "walletAddress"
// @param nonce body string true "nonce"
// @param signature body string true "signature"
// @Success      200  {string}  StatusOK
// @Success      201  {string}  StatusOK
// @Failure      400  {string}  StatusBadRequest
// @Failure      404  {string}  StatusNotFound
// @Failure      500  {string}  StatusInternalServerError
// @Router       /login [POST]
func (server *Server) LoginHandler(c echo.Context) (err error) {

	var p SigninPayload
	var tokenMaker token.Maker
	var duration time.Duration

	if err = (&echo.DefaultBinder{}).BindBody(c, &p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err = p.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidAddress)
	}

	address := strings.ToLower(p.WalletAddress)
	user, err := Authenticate(server, c, address, p.Signature)
	switch err {
	case nil:
	case ErrAuthError:
		return echo.NewHTTPError(http.StatusUnauthorized, ErrAuthError)
	default:
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	token, _, err := tokenMaker.CreateToken(user.FullName.String, duration)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	resToken := &accessToken{
		AccessToken: token,
	}
	return c.JSON(http.StatusCreated, resToken)
}

func Authenticate(server *Server, c echo.Context, walletAddress string, sigHex string) (db.User, error) {
	user, err := server.store.GetUserByWalletAddress(c.Request().Context(), sql.NullString{String: walletAddress, Valid: true})
	if err != nil {
		return user, echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	sig := hexutil.MustDecode(sigHex)
	// https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L516
	// check here why I am subtracting 27 from the last byte
	sig[crypto.RecoveryIDOffset] -= 27
	msg := accounts.TextHash([]byte(user.Nonce.String))
	recovered, err := crypto.SigToPub(msg, sig)

	if err != nil {
		user.Nonce = sql.NullString{}
		return user, echo.NewHTTPError(http.StatusUnauthorized, ErrMissingSig)
	}
	recoveredAddr := crypto.PubkeyToAddress(*recovered)

	if user.WalletAddress.String != strings.ToLower(recoveredAddr.Hex()) {
		return user, echo.NewHTTPError(http.StatusUnauthorized, ErrInvalidAddress)
	}

	// update the nonce here so that the signature cannot be resused
	nonce, err := GetNonce()
	if err != nil {
		return user, echo.NewHTTPError(http.StatusUnauthorized, ErrInvalidNonce)
	}
	user.Nonce.String = nonce

	server.store.UpdateUserNonce(c.Request().Context(), db.UpdateUserNonceParams{
		WalletAddress: user.WalletAddress,
		Nonce:         user.Nonce,
	})

	return user, nil
}
