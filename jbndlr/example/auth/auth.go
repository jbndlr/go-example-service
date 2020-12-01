package auth

import (
	"crypto/sha512"
	"errors"
	"jbndlr/example/conf"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	jwtKey           = []byte(conf.P.Auth.JWTKey)
	hash             = sha512.New384()
	expireAfter      = time.Duration(conf.P.Auth.ExpireAfter) * time.Second
	renewIfExpiresIn = time.Duration(conf.P.Auth.RenewWindow) * time.Second
)

var users = map[string]string{
	"user1": "pw1",
	"user2": "secret2",
}

// Subject : Data assigned an identity.
type Subject struct {
	Identity string    `json:"identity"`
	LastSeen time.Time `json:"lastSeen"`
}

// Credentials : Required information for auth request.
type Credentials struct {
	Secret   string `json:"secret"`
	Identity string `json:"identity"`
}

// Claims : Data to be encoded in JWT.
type Claims struct {
	Identity string `json:"identity"`
	jwt.StandardClaims
}

// Token : Contains shippable JWT.
type Token struct {
	JWT     string
	Expires time.Time
}

// PrematureRenewalError : If renewal attempted before `renewIfExpiresIn` is hit.
type PrematureRenewalError struct{}

func (m *PrematureRenewalError) Error() string {
	return "Premature JWT renewal attempt"
}

// HashString : Used for hashing prior to comparing secrets.
func HashString(pass string) string {
	return string(hash.Sum([]byte(pass)))
}

// Authenticate : Check credentials against persistence.
func Authenticate(creds *Credentials) (*Subject, *Token, error) {
	subject := &Subject{ // TODO: Load this from persistence
		Identity: creds.Identity,
		LastSeen: time.Now(),
	}
	expected, ok := users[creds.Identity]
	if !ok || expected != creds.Secret {
		return &Subject{}, &Token{}, errors.New("Authentication failed")
	}
	token, err := issueJWT(creds)
	if err != nil {
		return &Subject{}, &Token{}, err
	}
	return subject, token, nil
}

// Validate : Check/resolve pretended existing JWT authentication.
func Validate(jwts string) (*Subject, *Token, error) {
	claims, err := validateJWT(jwts)
	if err != nil {
		return &Subject{}, &Token{}, errors.New("Authentication failed")
	}
	subject := &Subject{ // TODO: Load this from persistence
		Identity: claims.Identity,
		LastSeen: time.Now(),
	}
	token := &Token{
		JWT:     jwts,
		Expires: time.Unix(claims.ExpiresAt, 0),
	}
	return subject, token, nil
}

// Renew : Check/resolve/renew pretended existing JWT authentication.
func Renew(jwts string) (*Subject, *Token, error) {
	claims, err := validateJWT(jwts)
	if err != nil {
		return &Subject{}, &Token{}, errors.New("Authentication failed")
	}
	subject := &Subject{ // TODO: Load this from persistence
		Identity: claims.Identity,
		LastSeen: time.Now(),
	}
	token, err := renewJWT(claims)
	if err != nil {
		return &Subject{}, &Token{}, err
	}
	return subject, token, nil
}

// issueJWT : Create and sign a new JWT.
func issueJWT(creds *Credentials) (*Token, error) {
	expiresAt := time.Now().Add(expireAfter)
	claims := &Claims{
		Identity: creds.Identity,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}
	return signClaims(claims)
}

// signClaims : Calculate signature hash and wrap into Token struct.
func signClaims(claims *Claims) (*Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	jwts, err := token.SignedString(jwtKey)
	if err != nil {
		return &Token{}, err
	}
	return &Token{
		JWT:     jwts,
		Expires: time.Unix(claims.ExpiresAt, 0),
	}, nil
}

// validateJWT : Check a signed JWT token string for validity.
func validateJWT(jwts string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(jwts, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return claims, err
}

// renewJWT : Adjourn expiration (without validating).
func renewJWT(claims *Claims) (*Token, error) {
	wouldExpireAt := time.Unix(claims.ExpiresAt, 0)
	if wouldExpireAt.Sub(time.Now()) > renewIfExpiresIn {
		return &Token{}, &PrematureRenewalError{}
	}
	expiresAt := time.Now().Add(expireAfter)
	claims.ExpiresAt = expiresAt.Unix()
	return signClaims(claims)
}
