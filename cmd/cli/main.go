package main

import (
	"fmt"
	"os"

	pmcli "plex_monitor/internal/cli"
	"plex_monitor/internal/database"
)

func main() {
	database.InitDB(os.Getenv("DATABASE_URL"), os.Getenv("DATABASE_NAME"))

	fmt.Println(`
	______ _            ___  ___            _ _             
	| ___ \ |           |  \/  |           (_) |            
	| |_/ / | _____  __ | .  . | ___  _ __  _| |_ ___  _ __ 
	|  __/| |/ _ \ \/ / | |\/| |/ _ \| '_ \| | __/ _ \| '__|
	| |   | |  __/>  <  | |  | | (_) | | | | | || (_) | |   
	\_|   |_|\___/_/\_\ \_|  |_/\___/|_| |_|_|\__\___/|_|   				 
	`)

	pmcli.GetCliApp().Run(os.Args)
}
