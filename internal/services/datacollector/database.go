package datacollector

import (
	"os"
	"plex_monitor/internal/database"
)

type MySQLDatabase struct{}

func (mysqlDb MySQLDatabase) Connect() error {
	database.InitDB(os.Getenv("DATABASE_URL"))

	return nil
}
