package db

import (
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func NewPointData(serverName string, playerCount int) *write.Point {
	return NewPointDataWithMode(serverName, playerCount, nil)
}

func NewPointDataWithMode(serverName string, playerCount int, parentServer *string) *write.Point {
	tagData := map[string]string{
		"server_name": serverName,
	}

	if parentServer != nil {
		tagData["parent_server"] = *parentServer
	}

	return influxdb2.NewPoint(
		"status",
		tagData,
		map[string]any{
			"players": playerCount,
		},
		time.Now(),
	)
}
