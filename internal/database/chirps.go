package database

type Chirp struct {
	// AuthorID int    `json:"author_id"`
	Body string `json:"body"`
	ID   int    `json:"id"`
}

// 5. Storage / 1. Storage
// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	// func (db *DB) CreateChirp(body string, authorID int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	// 5. Storage / 1. Storage
	// For now, just use integers for the id field,
	// and increment the id by 1 for each new chirp
	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID: id,
		// AuthorID: authorID,
		Body: body,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

// 5. Storage / 1. Storage
// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

// 5. Storage / 4. Get
// Add a new endpoint to your server that allows users to get a single chirp by ID.
func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	// If the chirp does not exist, the server should return a 404 status code.
	if !ok {
		return Chirp{}, ErrNotExist
	}

	return chirp, nil
}

func (db *DB) DeleteChirp(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(dbStructure.Chirps, id)
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}
