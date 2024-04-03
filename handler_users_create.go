package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Bayan2019/chirpy/internal/auth"
	"github.com/Bayan2019/chirpy/internal/database"
)

// 5. Storage / 7. Users
// For now, a user will just have an id (integer) and an email (string).
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	// IsChirpyRed bool   `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	// 6. Authentication / 1. Authentication with Passwords
	// Update the body parameters for this endpoint to include a new password field:
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	// 6. Authentication / 1. Authentication with Passwords
	// Hash the password using the bcrypt.GenerateFromPassword function
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	// 5. Storage / 7. Users
	// For now, a user will just have an id (integer) and an email (string).
	// user, err := cfg.DB.CreateUser(params.Email)

	// 6. Authentication / 1. Authentication with Passwords
	// Be sure to store the hashed password in the database as you create the user.
	user, err := cfg.DB.CreateUser(params.Email, hashedPassword)
	if err != nil {
		if errors.Is(err, database.ErrAlreadyExists) {
			respondWithError(w, http.StatusConflict, "User already exists")
			return
		}

		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
			// IsChirpyRed: user.IsChirpyRed,
		},
	})
}
