package main

import (
	"context"
	"fmt"
	"plex_monitor/internal/services/datacollector"
)

func getDataCollectorRepositoryFromServiceName(svcName string) datacollector.DataCollectorRepository {
	switch svcName {
	case datacollector.SONARR_QUEUE:
		return &datacollector.SonarrQueue{}
	case datacollector.SONARR_CALENDAR:
		return &datacollector.SonarrCalendar{}
	default:
		return nil
	}
}

func main() {
	var availableServices = []string{datacollector.SONARR_QUEUE, datacollector.SONARR_CALENDAR}

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
