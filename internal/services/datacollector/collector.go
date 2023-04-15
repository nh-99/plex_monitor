package datacollector

import (
	"context"
)

// Creates the connection to the database for storage.
type Database interface {
	connect() error
}

// Handles collecting and storing data in the database from various external service providers.
type DataCollectorRepository interface {
	collect() error
	store(db Database) error
}

// Executes the functions for data collection & storage.
type DataCollectionService struct {
	database   Database
	repository DataCollectorRepository
}

// Create a new service with the injected dependencies.
func NewDataCollectionService(db Database, repo DataCollectorRepository) *DataCollectionService {
	return &DataCollectionService{
		database:   db,
		repository: repo,
	}
}

// Run the data collection & storage.
func (s *DataCollectionService) Execute(ctx context.Context) error {
	// Run the collect method to populate the repository
	err := s.repository.collect()
	if err != nil {
		return err
	}

	// Run the store method to store the repository to the database
	err = s.repository.store(s.database)
	if err != nil {
		return err
	}

	return nil
}
