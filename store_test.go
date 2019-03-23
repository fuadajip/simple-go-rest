package main

import (
	"database/sql"

	"github.com/stretchr/testify/suite"
)

type StoreSuite struct {
	suite.Suite
	/*
		The suite is defined as a struct, with the store and db as its
		attributes. Any variables that are to be shared between tests in a
		suite should be stored as attributes of the suite instance
	*/
	store *dbStore
	db    *sql.DB
}

func (s *StoreSuite) SetupSuite() {
	/*
		The database connection is opened in the setup, and
		stored as an instance variable,
		as is the higher level `store`, that wraps the `db`
	*/

	connString := "dbname=sample_rest_go sslmode=disable"
	db, err := sql.Open("postgres", connString)
	if err != nil {
		s.T().Fatal(err)
	}

	s.db = db
	s.store = &dbStore{db: db}
}

func (s *StoreSuite) SetupTest() {
	/*
		We delete all entries from the table before each test runs, to ensure a
		consistent state before our tests run. In more complex applications, this
		is sometimes achieved in the form of migrations
	*/
	_, err := s.db.Query("DELETE from birds")
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *StoreSuite) TearDownSuite() {
	// close the connection after all test in the suite finish
	s.db.Close()
}

func (s *StoreSuite) TestCreateBird() {
	s.store.CreateBird(&Bird{
		Description: "test description",
		Species:     "test species",
	})

	// Query database for the entry we just created
	res, err := s.db.Query("select count(*) from birds where description='test description'")
	if err != nil {
		s.T().Fatal(err)
	}

	var count int
	for res.Next() {
		err := res.Scan(&count)
		if err != nil {
			s.T().Fatal(err)
		}
	}
	// Assert that there must be one entry with the properties of the bird that
	// we just inserted (since the database was empty before this)
	if count != 1 {
		s.T().Errorf("incorrect count: wanted 1, got %d", count)
	}
}

func (s *StoreSuite) TestGetBirdHandler() {
	// Insert sample bird into `birds` table
	_, err := s.db.Query("INSERT INTP birds (species, description) VALUES ('bird','description')")
	if err != nil {
		s.T().Fatal(err)
	}

	birds, err := s.store.GetBirds()
	if err != nil {
		s.T().Fatal(err)
	}

	// assert that the count of birds received must be 1
	nBirds := len(birds)
	if nBirds != 1 {
		s.T().Errorf("incorrect count, wanted 1 got %d", nBirds)
	}

	// assert that the details of bird is same as the one we inserted
	expectedBird := Bird{"bird", "description"}
	if *birds[0] != expectedBird {
		s.T().Errorf("incorect details, expected %v got %v", expectedBird, *birds[0])
	}
}
