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

	db "github.com/ifandonlyif-io/ifandonlyif-backend/db/sqlc"
	"github.com/labstack/echo/v4"
)

type User struct {
	walletAddress string
	nonce         string
}

var (
	max  *big.Int
	once sync.Once
)

var (
	hexRegex   *regexp.Regexp = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)
	nonceRegex *regexp.Regexp = regexp.MustCompile(`^[0-9]+$`)
)

type RegisterPayload struct {
	Address string `json:"address"`
}

func (p RegisterPayload) Validate() error {
	if !hexRegex.MatchString(p.Address) {
		return ErrInvalidAddress
	}
	return nil
}

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

func (server *Server) RegisterHandler(c echo.Context) (err error) {

	u := new(User)
	if err = c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var p RegisterPayload

	if err := (&echo.DefaultBinder{}).BindBody(c, &p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := p.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	nonce, err := GetNonce()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	user := db.CreateUserParams{
		WalletAddress: sql.NullString{String: strings.ToLower(p.Address), Valid: true},
		Nonce:         sql.NullString{String: nonce, Valid: true},
	}

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

// func Authenticate( address string, nonce string, sigHex string) (User, error) {

// 	if err != nil {
// 		return user, err
// 	}
// 	if user.Nonce != nonce {
// 		return user, ErrAuthError
// 	}

// 	sig := hexutil.MustDecode(sigHex)
// 	// https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L516
// 	// check here why I am subtracting 27 from the last byte
// 	sig[crypto.RecoveryIDOffset] -= 27
// 	msg := accounts.TextHash([]byte(nonce))
// 	recovered, err := crypto.SigToPub(msg, sig)
// 	if err != nil {
// 		return user, err
// 	}
// 	recoveredAddr := crypto.PubkeyToAddress(*recovered)

// 	if user.Address != strings.ToLower(recoveredAddr.Hex()) {
// 		return user, ErrAuthError
// 	}

// 	// update the nonce here so that the signature cannot be resused
// 	nonce, err = GetNonce()
// 	if err != nil {
// 		return user, err
// 	}
// 	user.Nonce = nonce
// 	storage.Update(user)

// 	return user, nil
// }

type SigninPayload struct {
	Address string `json:"address"`
	Nonce   string `json:"nonce"`
	Sig     string `json:"sig"`
}

// func getUserFromReqContext(r *http.Request) User {
// 	ctx := r.Context()
// 	key := ctx.Value("user").(User)
// 	return key
// }

// func (s SigninPayload) validateWalletAddress() error {
// 	if !hexRegex.MatchString(s.Address) {
// 		return ErrInvalidAddress
// 	}
// 	if !nonceRegex.MatchString(s.Nonce) {
// 		return ErrInvalidNonce
// 	}
// 	if len(s.Sig) == 0 {
// 		return ErrMissingSig
// 	}
// 	return nil
// }
