package main

import (
	"context"
	"fmt"
	"plex_monitor/internal/services/datacollector"
)

const (
	SONARR_QUEUE    = "sonarrQueue"
	SONARR_CALENDAR = "sonarrCalendar"
)

func getDataCollectorRepositoryFromServiceName(svcName string) datacollector.DataCollectorRepository {
	switch svcName {
	case SONARR_QUEUE:
		return &datacollector.SonarrQueue{}
	case SONARR_CALENDAR:
		return &datacollector.SonarrCalendar{}
	default:
		return nil
	}
}

func main() {
	var availableServices = []string{SONARR_QUEUE, SONARR_CALENDAR}

	// TODO: Go through for loop of all services and create collector service for each
	for _, serviceName := range availableServices {
		// Declare dependencies
		database := datacollector.MySQLDatabase{}
		dataCollector := getDataCollectorRepositoryFromServiceName(serviceName)

		if dataCollector == nil {
			fmt.Printf("\n[plex_monitor] Invalid data collector %s\n", serviceName)
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
