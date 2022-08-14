package crypto

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"regexp"
	"sync"
)

type User struct {
	Address string
	Nonce   string
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

func bindReqBody(r *http.Request, obj any) error {
	return json.NewDecoder(r.Body).Decode(obj)
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

func getUserFromReqContext(r *http.Request) User {
	ctx := r.Context()
	key := ctx.Value("user").(User)
	return key
}

func renderJson(r *http.Request, w http.ResponseWriter, statusCode int, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8 ")
	var body []byte
	if res != nil {
		var err error
		body, err = json.Marshal(res)
		if err != nil { // TODO handle me better
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	w.WriteHeader(statusCode)
	if len(body) > 0 {
		w.Write(body)
	}
}
