package main

import (
	"context"
	"fmt"
	"plex_monitor/internal/database/models"
	"plex_monitor/internal/services/datacollector"
)

func getDataCollectorRepositoryFromServiceName(svcName string) datacollector.DataCollectorRepository {
	switch svcName {
	case datacollector.REPOSITORY_SONARR_QUEUE:
		return &datacollector.SonarrQueue{}
	case datacollector.REPOSITORY_SONARR_CALENDAR:
		return &datacollector.SonarrCalendar{}
	case datacollector.REPOSITORY_TRANSMISSION:
		return &datacollector.Transmission{}
	default:
		return nil
	}
}

func main() {
	database := datacollector.MySQLDatabase{}
	// Connect the database
	err := database.Connect()
	if err != nil {
		panic(err)
	}

	availableServices, err := models.GetScannableMonitoredServices()
	if err != nil {
		panic(err)
	}

	for _, service := range availableServices {
		// Declare dependencies
		dataCollector := getDataCollectorRepositoryFromServiceName(service.Identifier)

		if dataCollector == nil {
			panic(fmt.Sprintf("\n[plex_monitor] Invalid data collector %s\n", service.Identifier))
		}

		// Create service
		svc := datacollector.NewDataCollectionService(database, dataCollector)

		// Execute service once
		err := svc.Execute(context.TODO())

		// Handle errors
		if err != nil {
			panic(err)
		}
	}
}
