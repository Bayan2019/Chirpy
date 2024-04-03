package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Bayan2019/chirpy/internal/auth"
)

// 6. Authentication / 1. Authentication with Passwords
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	// 6. Authentication / 1. Authentication with Passwords
	// 6. Authentication / 6. Authentication with JWTs
	// This endpoint should accept a new optional expires_in_seconds field
	// in the request body
	type parameters struct {
		Password  string        `json:"password"`
		Email     string        `json:"email"`
		ExpiresIn time.Duration `json:"expires_in_seconds"`
	}

	type response struct {
		User
		Token string `json:"token"`
		// RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	// 6. Authentication / 1. Authentication with Passwords
	// you don't have access to an ID here
	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	// 6. Authentication / 1. Authentication with Passwords
	// Use the bcrypt.CompareHashAndPassword function
	// to compare the password that the user entered in the HTTP request
	// with the password that is stored in the database.
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	// 6. Authentication / 1. Authentication with Passwords
	// If the passwords do not match, return a 401 Unauthorized response.
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	// 6. Authentication / 6. Authentication with JWTs
	// expires_in_seconds is an optional parameter.
	// If it's specified by the client, use it as the expiration time.
	// If it's not specified, use a default expiration time of 24 hours.
	// If the client specified a number over 24 hours,
	// use 24 hours as the expiration time.
	if params.ExpiresIn == 0 || params.ExpiresIn > 24*time.Hour {
		params.ExpiresIn = 24 * time.Hour
	}
	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, params.ExpiresIn, auth.TokenTypeAccess)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT")
		return
	}

	// refreshToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour*24*30*6, auth.TokenTypeRefresh)
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh JWT")
	// 	return
	// }

	// 6. Authentication / 1. Authentication with Passwords
	// If the passwords match, return a 200 OK response and a copy of the user resource
	// 6. Authentication / 6. Authentication with JWTs
	// Once you have the token, respond to the request with a 200 code and token
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
		Token: accessToken,
	})
	// respondWithJSON(w, http.StatusOK, Response{
	// 	User: User{
	// 		ID:          user.ID,
	// 		Email:       user.Email,
	// 		IsChirpyRed: user.IsChirpyRed,
	// 	},
	// 	Token:        accessToken,
	// 	RefreshToken: refreshToken,
	// })
}
