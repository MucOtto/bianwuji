package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func getMeasurements(url string) ([]float64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(body)))
	measurements := make([]float64, 0)
	for scanner.Scan() {
		line := scanner.Text()
		values := strings.Fields(line)
		for _, v := range values {
			value, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, err
			}
			measurements = append(measurements, value)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return measurements, nil
}

func main() {
	serviceURL := os.Getenv("SERVICE_URL")
	if serviceURL == "" {
		log.Fatalf("SERVICE_URL environment variable is required")
	}

	intervalStr := os.Getenv("INTERVAL_SECONDS")
	if intervalStr == "" {
		log.Fatalf("INTERVAL_SECONDS environment variable is required")
	}

	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		log.Fatalf("Invalid interval: %v", err)
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			measurements, err := getMeasurements(serviceURL)
			if err != nil {
				log.Printf("Error getting measurements: %v\n", err)
				continue
			}

			if len(measurements) == 0 {
				log.Println("No valid measurements retrieved")
				continue
			}

			sum := 0.0
			for _, value := range measurements {
				sum += value
			}
			average := sum / float64(len(measurements))
			log.Printf("Average measurement: %.2f\n", average)
		}
	}
}
