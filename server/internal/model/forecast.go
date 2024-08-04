package model

import "time"

type Forecast struct {
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
	DetailedForecast string    `json:"detailedForecast"`
}
