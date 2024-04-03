package database

import "errors"

// 5. Storage / 7. Users
// For now, a user will just have an id (integer) and an email (string).
type User struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	// IsChirpyRed    bool   `json:"is_chirpy_red"`
}

var ErrAlreadyExists = errors.New("already exists")

// 5. Storage / 7. Users
// func (db *DB) CreateUser(email string) (User, error) {
// 6. Authentication / 1. Authentication with Passwords
func (db *DB) CreateUser(email, hashedPassword string) (User, error) {
	_, err := db.GetUserByEmail(email)

	if !errors.Is(err, ErrNotExist) {
		return User{}, ErrAlreadyExists
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID:             id,
		Email:          email,
		HashedPassword: hashedPassword,
		// IsChirpyRed:    false,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUser(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}

	return user, nil
}

// 6. Authentication / 1. Authentication with Passwords
// you don't have access to an ID here
func (db *DB) GetUserByEmail(email string) (User, error) {

	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, ErrNotExist
}

// 6. Authentication / 6. Authentication with JWTs
// You'll probably need to add a new UpdateUser method to your database package
func (db *DB) UpdateUser(id int, email, hashedPassword string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}

	user.Email = email
	user.HashedPassword = hashedPassword
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// func (db *DB) UpgradeChirpyRed(id int) (User, error) {
// 	dbStructure, err := db.loadDB()
// 	if err != nil {
// 		return User{}, err
// 	}

// 	user, ok := dbStructure.Users[id]
// 	if !ok {
// 		return User{}, ErrNotExist
// 	}

// 	user.IsChirpyRed = true
// 	dbStructure.Users[id] = user

// 	err = db.writeDB(dbStructure)
// 	if err != nil {
// 		return User{}, err
// 	}

// 	return user, nil
// }
