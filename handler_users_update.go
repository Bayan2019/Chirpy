package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Bayan2019/chirpy/internal/auth"
)

// 6. Authentication / 6. Authentication with JWTs
func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
	}

	// 6. Authentication / 6. Authentication with JWTs
	// This is our first authenticated endpoint,
	// which means it will require a JWT to be present in the request headers
	// First, extract the token from the request headers r.Header.Get
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	// 6. Authentication / 6. Authentication with JWTs
	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	// 6. Authentication / 6. Authentication with JWTs
	// If the JWT was valid, you should now have the ID of the authenticated user.
	userIDInt, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	// 6. Authentication / 6. Authentication with JWTs
	// You'll probably need to add a new UpdateUser method to your database package
	user, err := cfg.DB.UpdateUser(userIDInt, params.Email, hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	// 6. Authentication / 6. Authentication with JWTs
	// After updating the user,
	// return a copy of the updated user resource
	// (without the password) and a 200 status code
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}
