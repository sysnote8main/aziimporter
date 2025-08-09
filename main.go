package main

import (
	"fmt"
	"log/slog"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/sysnote8main/aziimporter/internal/config"
	"github.com/sysnote8main/aziimporter/internal/db"
	"github.com/sysnote8main/azisabaapi/pkg/aziapi"
)

func main() {
	slog.Info("Loading config...")
	const configPath = "config.yml"

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	aziClient := aziapi.NewAziApiClient(
		cfg.Token,
		cfg.Url,
	)

	dbCfg := cfg.DB
	dbClient := influxdb2.NewClientWithOptions(dbCfg.Url, dbCfg.Token, influxdb2.DefaultOptions().SetBatchSize(20))
	defer dbClient.Close()
	writeApi := dbClient.WriteAPI(dbCfg.Organization, dbCfg.Bucket)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		getApiAndPutData(aziClient, writeApi)
		slog.Info("Crawled!")
		<-ticker.C
	}
}

func getApiAndPutData(aziClient aziapi.AziApiClient, writeApi api.WriteAPI) {
	res, err := aziClient.GetCounts()
	if err != nil {
		slog.Error("Failed to get counts from azisaba api", slog.Any("error", err))
		return
	}
	writeApi.WritePoint(db.NewPointData("all", res.TotalPlayers))
	for name, data := range res.Games {
		writeApi.WritePoint(db.NewPointData(name, data.Players))
		for k, v := range data.Modes {
			writeApi.WritePoint(db.NewPointDataWithMode(fmt.Sprintf("%s-%s", name, k), v, &name))
		}
	}
	defer writeApi.Flush()
}
