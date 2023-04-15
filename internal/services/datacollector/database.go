package datacollector

import (
	"database/sql"
	"os"
	"plex_monitor/internal/database"
)

type MySQLDatabase struct {
	db *sql.DB
}

func (mysqlDb MySQLDatabase) connect() error {
	database.InitDB(os.Getenv("DATABASE_URL"))
	mysqlDb.db = database.DB

	return nil
}
