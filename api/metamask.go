package api

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	db "github.com/ifandonlyif-io/ifandonlyif-backend/db/sqlc"
	"github.com/ifandonlyif-io/ifandonlyif-backend/util"
	"github.com/labstack/echo/v4"
)

type User struct {
	Wallet string `json:"wallet"`
	Nonce  string `json:"nonce"`
}

type code struct {
	Code string `json:"code"`
}

type returnToken struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RegisterPayload struct {
	Wallet string `json:"wallet"`
}

type SigninPayload struct {
	Wallet    string `json:"wallet"`
	Signature string `json:"signature"`
}

func (p RegisterPayload) Validate() error {
	if !hexRegex.MatchString(p.Wallet) {
		return ErrInvalidAddress
	}
	return nil
}

func (s SigninPayload) Validate() error {
	if !hexRegex.MatchString(s.Wallet) {
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

// code godoc
// @Summary      code
// @Description  register a new user
// @Tags         code
// @Accept       json
// @produce application/json
// @param wallet body string true "wallet"
// @Success      200  {string}  StatusOK
// @Success      201  {string}  StatusOK
// @Failure      400  {string}  StatusBadRequest
// @Failure      404  {string}  StatusNotFound
// @Failure      500  {string}  StatusInternalServerError
// @Router       /code [POST]
func (server *Server) NonceHandler(c echo.Context) (err error) {
	var p RegisterPayload

	if err := (&echo.DefaultBinder{}).BindBody(c, &p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// validate wallet address
	err = p.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidAddress)
	}

	user, err := server.store.GetUserByWalletAddress(c.Request().Context(), sql.NullString{String: p.Wallet, Valid: true})
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
		return c.JSON(http.StatusOK, resCode)
	}

	nonce, err := GetNonce()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidNonce)
	}

	createUser, err := server.store.CreateUser(c.Request().Context(), db.CreateUserParams{
		Wallet:   sql.NullString{String: p.Wallet, Valid: true},
		Nonce:    sql.NullString{String: nonce, Valid: true},
		FullName: sql.NullString{String: "", Valid: true},
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
// @param wallet body string true "wallet"
// @param signature body string true "signature"
// @Success      200  {string}  StatusOK
// @Success      201  {string}  StatusOK
// @Failure      400  {string}  StatusBadRequest
// @Failure      404  {string}  StatusNotFound
// @Failure      500  {string}  StatusInternalServerError
// @Router       /login [POST]
func (server *Server) LoginHandler(c echo.Context) (err error) {

	var p SigninPayload

	// parse payload
	if err = (&echo.DefaultBinder{}).BindBody(c, &p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// validate
	if err = p.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidAddress)
	}

	user, err := Authenticate(server, c, p.Wallet, p.Signature)
	switch err {
	case nil:
	case ErrAuthError:
		return echo.NewHTTPError(http.StatusUnauthorized, ErrAuthError)
	default:
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	accesstoken, _, accessErr := server.accessTokenMaker.CreateToken(user.FullName.String, user.Wallet.String, server.config.AccessTokenDuration)
	if accessErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	refreshtoken, _, refreshErr := server.refreshTokenMaker.CreateToken(user.FullName.String, user.Wallet.String, server.config.RefreshTokenDuration)
	if refreshErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	resToken := &returnToken{
		AccessToken:  accesstoken,
		RefreshToken: refreshtoken,
	}

	return c.JSON(http.StatusCreated, resToken)
}

func Authenticate(server *Server, c echo.Context, wallet string, sigHex string) (db.GetUserByWalletAddressRow, error) {
	user, err := server.store.GetUserByWalletAddress(c.Request().Context(), sql.NullString{String: wallet, Valid: true})
	if err != nil {
		return db.GetUserByWalletAddressRow{}, echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	fmt.Print(user)
	sig := hexutil.MustDecode(sigHex)
	// https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L516
	// check here why I am subtracting 27 from the last byte
	sig[crypto.RecoveryIDOffset] -= 27
	msg := accounts.TextHash([]byte("Welcome to IfAndOnlyIf.io!! We will enhance your WEB3/NFT experience. Security Nonce: " + user.Nonce.String))
	recovered, err := crypto.SigToPub(msg, sig)

	if err != nil {
		user.Nonce.String = ""
		return db.GetUserByWalletAddressRow{}, echo.NewHTTPError(http.StatusUnauthorized, ErrMissingSig)
	}
	recoveredAddr := crypto.PubkeyToAddress(*recovered)

	// check database string lower case
	if util.IsLower(user.Wallet.String) {
		// compare with lowercase from recovery
		if user.Wallet.String != strings.ToLower(recoveredAddr.Hex()) {
			return db.GetUserByWalletAddressRow{}, echo.NewHTTPError(http.StatusUnauthorized, ErrInvalidAddress)
		}
	}
	// check database string upper case
	if util.IsUpper(user.Wallet.String) {
		// compare with lowercase from recovery
		if user.Wallet.String != strings.ToUpper(recoveredAddr.Hex()) {
			return db.GetUserByWalletAddressRow{}, echo.NewHTTPError(http.StatusUnauthorized, ErrInvalidAddress)
		}
	}
	// check if Mixed with upper and lower
	if !util.IsLower(user.Wallet.String) {
		if !util.IsUpper(user.Wallet.String) {
			// check the original address
			if user.Wallet.String != recoveredAddr.Hex() {
				return db.GetUserByWalletAddressRow{}, echo.NewHTTPError(http.StatusUnauthorized, ErrInvalidAddress)
			}
		}
	}

	// update the nonce here so that the signature cannot be resused
	nonce, err := GetNonce()
	if err != nil {
		return user, echo.NewHTTPError(http.StatusUnauthorized, ErrInvalidNonce)
	}

	user.Nonce.String = nonce

	server.store.UpdateUserNonce(c.Request().Context(), db.UpdateUserNonceParams{
		Wallet: sql.NullString{String: user.Wallet.String, Valid: true},
		Nonce:  sql.NullString{String: user.Nonce.String, Valid: true},
	})

	return user, nil
}
