package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrNoAuthHeaderIncluded = errors.New("not auth header included in request")

type TokenType string

const (
	TokenTypeAccess  TokenType = "chirpy-access"
	TokenTypeRefresh TokenType = "chirpy-refresh"
)

// 6. Authentication / 1. Authentication with Passwords
// Hash the password using the bcrypt.GenerateFromPassword function
// HashPassword -
func HashPassword(password string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(dat), nil
}

// 6. Authentication / 1. Authentication with Passwords
// Use the bcrypt.CompareHashAndPassword function
// to compare the password that the user entered in the HTTP request
// with the password that is stored in the database.
// CheckPasswordHash -
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// 6. Authentication / 6. Authentication with JWTs
// Create a JWT using JWT library
func MakeJWT(userID int, tokenSecret string, expiresIn time.Duration, tokenType TokenType) (string, error) {
	signingKey := []byte(tokenSecret)

	// 6. Authentication / 6. Authentication with JWTs
	// Use jwt.NewWithClaims to create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		// 6. Authentication / 6. Authentication with JWTs
		// Set the Issuer to "chirpy"
		Issuer: "chirpy",
		//Issuer:    string(tokenType),
		// 6. Authentication / 6. Authentication with JWTs
		// Set IssuedAt to the current time in UTC
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		// 6. Authentication / 6. Authentication with JWTs
		// Set ExpiresAt to the current time plus the expiration time
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		// 6. Authentication / 6. Authentication with JWTs
		// Set the Subject to a stringified version of the user's id
		Subject: fmt.Sprintf("%d", userID),
	})

	// 6. Authentication / 6. Authentication with JWTs
	// Use token.SignedString to sign the token with the secret key.
	return token.SignedString(signingKey)
}

func RefreshToken(tokenString, tokenSecret string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return "", err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}
	if issuer != string(TokenTypeRefresh) {
		return "", errors.New("invalid issuer")
	}

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		return "", err
	}

	newToken, err := MakeJWT(userID, tokenSecret, time.Hour, TokenTypeAccess)
	if err != nil {
		return "", err
	}

	return newToken, nil
}

// 6. Authentication / 6. Authentication with JWTs
func ValidateJWT(tokenString, tokenSecret string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}
	// 6. Authentication / 6. Authentication with JWTs
	// use the jwt.ParseWithClaims function to validate the signature of the JWT
	// and extract the claims into a *jwt.Token struct
	token, err := jwt.ParseWithClaims(tokenString, &claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return "", err
	}

	// 6. Authentication / 6. Authentication with JWTs
	// use the token.Claims interface
	// to get access to the user's id from the claims
	// (which should be stored in the Subject field)
	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	// issuer, err := token.Claims.GetIssuer()
	// if err != nil {
	// 	return "", err
	// }
	// if issuer != string(TokenTypeAccess) {
	// 	return "", errors.New("invalid issuer")
	// }

	return userIDString, nil
}

// 6. Authentication / 6. Authentication with JWTs
// This is our first authenticated endpoint,
// which means it will require a JWT to be present in the request headers
func GetBearerToken(headers http.Header) (string, error) {
	// 6. Authentication / 6. Authentication with JWTs
	// First, extract the token from the request headers r.Header.Get
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	// 6. Authentication / 6. Authentication with JWTs
	// Remember, you'll need to strip off the Bearer prefix
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}
