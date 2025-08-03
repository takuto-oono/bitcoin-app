package main

import (
	"flag"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"bitcoin-app-golang/api"
	"bitcoin-app-golang/config"
)

const (
	resultChanSize = 60 * 60

	GolangServerName  = "Golang Server"
	DRFServerName     = "DRF Server"
	FastAPIServerName = "FastAPI Server"
)

func main() {
	tomlFilePath := flag.String("toml", "toml/local.toml", "toml file path")
	envFilePath := flag.String("env", "env/.env.local", "env file path")
	flag.Parse()

	cfg, err := config.NewConfig(*tomlFilePath, *envFilePath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// golangServerに依存しないように直接APIフォルダのメソッドを叩く
	lineAPI, err := api.NewLineAPI(cfg)
	if err != nil {
		panic(err)
	}

	golangServer := api.NewGolangServerAPI(cfg)
	drf := api.NewDRFAPI(cfg)
	fastAPIServer := api.NewFastAPI(cfg)

	checkInterval := 5 * time.Second
	checkTicker := time.NewTicker(checkInterval)
	defer checkTicker.Stop()

	notificationInterval := 1 * time.Hour
	notificationTicker := time.NewTicker(notificationInterval)
	defer notificationTicker.Stop()

	golangServerChan := make(chan HealthCheckResult, resultChanSize)
	drfServerChan := make(chan HealthCheckResult, resultChanSize)
	fastAPIChan := make(chan HealthCheckResult, resultChanSize)

	editCheckResult := func(results []HealthCheckResult) string {
		editedResults := make([]string, 0, len(results))

		for i, r := range results {
			if i == 0 {
				result := ""
				if r.Error != nil {
					result = fmt.Sprintf("Time: %s, Status: NG", convertJST(r.TimeStamp))
				} else {
					result = fmt.Sprintf("Time: %s, Status: OK", convertJST(r.TimeStamp))
				}
				editedResults = append(editedResults, result)
				continue
			}

			if i == len(results)-1 {
				result := ""
				if r.Error != nil {
					result = fmt.Sprintf("Time: %s, Status: NG", convertJST(r.TimeStamp))
				} else {
					result = fmt.Sprintf("Time: %s, Status: OK", convertJST(r.TimeStamp))
				}
				editedResults = append(editedResults, result)
				continue
			}

			if (r.Error != nil) == (results[i-1].Error != nil) {
				continue
			}

			result := ""
			if r.Error != nil {
				result = fmt.Sprintf("Time: %s, Status: Down", convertJST(r.TimeStamp))
			} else {
				result = fmt.Sprintf("Time: %s, Status: Up", convertJST(r.TimeStamp))
			}
			editedResults = append(editedResults, result)
		}

		if len(editedResults) == 0 {
			return "No health check results available."
		}

		return strings.Join(editedResults, "\n")
	}

	for {
		select {
		case <-checkTicker.C:
			go func() {
				err := golangServer.GetHealthcheck()
				golangServerChan <- HealthCheckResult{
					ServerName: GolangServerName,
					TimeStamp:  time.Now().Unix(),
					Error:      err,
				}
			}()

			go func() {
				err := drf.GetHealthcheck()
				drfServerChan <- HealthCheckResult{
					ServerName: DRFServerName,
					TimeStamp:  time.Now().Unix(),
					Error:      err,
				}
			}()

			go func() {
				err := fastAPIServer.GetHealthcheck()
				fastAPIChan <- HealthCheckResult{
					ServerName: FastAPIServerName,
					TimeStamp:  time.Now().Unix(),
					Error:      err,
				}
			}()

		case <-notificationTicker.C:
			golangResult := make([]HealthCheckResult, 0, len(golangServerChan))
			drfResult := make([]HealthCheckResult, 0, len(drfServerChan))
			fastAPIResult := make([]HealthCheckResult, 0, len(fastAPIChan))

			for len(golangServerChan) > 0 {
				golangResult = append(golangResult, <-golangServerChan)
			}
			for len(drfServerChan) > 0 {
				drfResult = append(drfResult, <-drfServerChan)
			}
			for len(fastAPIChan) > 0 {
				fastAPIResult = append(fastAPIResult, <-fastAPIChan)
			}

			slices.SortFunc(golangResult, func(a, b HealthCheckResult) int {
				return sortResultFunc(a, b)
			})
			slices.SortFunc(drfResult, func(a, b HealthCheckResult) int {
				return sortResultFunc(a, b)
			})
			slices.SortFunc(fastAPIResult, func(a, b HealthCheckResult) int {
				return sortResultFunc(a, b)
			})

			golangHealth := editCheckResult(golangResult)
			drfHealth := editCheckResult(drfResult)
			fastAPIHealth := editCheckResult(fastAPIResult)

			healthMessage := fmt.Sprintf("%s Health:\n%s\n\n%s Health:\n%s\n\n%s Health:\n%s",
				GolangServerName, golangHealth, DRFServerName, drfHealth, FastAPIServerName, fastAPIHealth)

			if err := lineAPI.PostMessage(healthMessage); err != nil {
				log.Printf("Failed to send health check notification: %v", err)
			} else {
				log.Println("Health check notification sent successfully.")
			}
		}
	}
}

type HealthCheckResult struct {
	ServerName string
	TimeStamp  int64
	Error      error
}

func convertJST(unixTime int64) string {
	utcTime := time.Unix(unixTime, 0)

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}

	jstTime := utcTime.In(jst)
	return jstTime.Format("2006-01-02 15:04:05")
}

func sortResultFunc(a, b HealthCheckResult) int {
	if a.TimeStamp < b.TimeStamp {
		return -1
	} else if a.TimeStamp > b.TimeStamp {
		return 1
	}
	return 0
}
