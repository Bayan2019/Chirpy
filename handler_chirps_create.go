package main

import (
	"encoding/json"
	"errors"
	"net/http"

	// "strconv"
	"strings"
	// encapsulating all of your database logic in an internal database package
	// "github.com/Bayan2019/chirpy/internal/auth"
)

// 5. Storage / 1. Storage
// If the chirp is valid, you should give it a unique id
type Chirp struct {
	ID int `json:"id"`
	// AuthorID int    `json:"author_id"`
	Body string `json:"body"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	// token, err := auth.GetBearerToken(r.Header)
	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
	// 	return
	// }

	// subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
	// 	return
	// }

	// userID, err := strconv.Atoi(subject)
	// if err != nil {
	// 	respondWithError(w, http.StatusBadRequest, "Couldn't parse user ID")
	// 	return
	// }

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 5. Storage / 1. Storage
	// CreateChirp creates a new chirp and saves it to disk
	chirp, err := cfg.DB.CreateChirp(cleaned)
	// chirp, err := cfg.DB.CreateChirp(cleaned, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID: chirp.ID,
		// AuthorID: chirp.AuthorID,
		Body: chirp.Body,
	})

	// 4. JSON / 2. JSON
	// If the Chirp is valid, respond with a 200 code and this body:

	// type returnVals struct {
	// 	// Valid bool `json:"valid"`
	// 	CleanBody string `json:"cleaned_body"`
	// }

	// 4. JSON 6. The Profane
	// your handler should return the cleaned version of the text in a JSON
	// respondWithJSON(w, http.StatusOK, returnVals{
	// 	CleanBody: cleaned,
	// })
}

func validateChirp(body string) (string, error) {
	// 4. JSON / 2. JSON
	// all Chirps must be 140 characters long or less.
	// if the Chirp is too long, respond with a 400 code
	const maxChirpLength = 140

	if len(body) > maxChirpLength {
		// return errors.New("Chirp is too long")
		return "", errors.New("Chirp is too long")
	}

	// 4. JSON 6. The Profane
	// replace all "profane" words with 4 asterisks: ****
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleaned := getCleanedBody(body, badWords)

	// return nil
	return cleaned, nil
}

// 4. JSON 6. The Profane
// replace all "profane" words with 4 asterisks: ****
func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")

	for i, word := range words {
		// 4. JSON 6. The Profane
		// Be sure to match against uppercase versions of the words as well
		loweredWord := strings.ToLower(word)
		_, ok := badWords[loweredWord]
		if ok {
			words[i] = "****"
		}
	}

	cleaned := strings.Join(words, " ")
	return cleaned
}
