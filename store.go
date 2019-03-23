package main

import "database/sql"

type (

	// Store interface
	Store interface {
		CreateBird(bird *Bird) error
		GetBirds() ([]*Bird, error)
	}
	dbStore struct {
		db *sql.DB
	}
)

func (store *dbStore) CreateBird(bird *Bird) error {
	_, err := store.db.Query("INSERT INTO birds (species, description) VALUES ($1, $2)", bird.Species, bird.Description)
	return err
}

func (store *dbStore) GetBirds() ([]*Bird, error) {
	rows, err := store.db.Query("SELECT species, description from birds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	birds := []*Bird{}

	for rows.Next() {
		bird := &Bird{}

		if err := rows.Scan(&bird.Species, &bird.Description); err != nil {
			return nil, err
		}

		birds = append(birds, bird)
	}

	return birds, nil
}

// The store variable is a package level variable that will be available for
// use throughout our application code
var store Store

/* We will need to call the InitStore method to initialize the store. This will
typically be done at the beginning of our application (in this case, when the server starts up)
This can also be used to set up the store as a mock, which we will be observing
later on
*/

// InitStore ...
func InitStore(s Store) {
	store = s
}
