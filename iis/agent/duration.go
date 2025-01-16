package agent

import (
	"encoding/json"
	"time"
)

type DurationMinutes int
type DurationSeconds int

type timeSpan struct {
	Ticks             float64 `json:"Ticks"`
	Days              float64 `json:"Days"`
	Hours             float64 `json:"Hours"`
	Milliseconds      float64 `json:"Milliseconds"`
	Minutes           float64 `json:"Minutes"`
	Seconds           float64 `json:"Seconds"`
	TotalDays         float64 `json:"TotalDays"`
	TotalHours        float64 `json:"TotalHours"`
	TotalMilliseconds float64 `json:"TotalMilliseconds"`
	TotalMinutes      float64 `json:"TotalMinutes"`
	TotalSeconds      float64 `json:"TotalSeconds"`
}

func (duration DurationMinutes) toTimeString() string {
	var t time.Time
	t = t.Add(time.Minute * time.Duration(duration))
	return t.Format("15:04:05")
}

func (duration DurationSeconds) toTimeString() string {
	var t time.Time
	t = t.Add(time.Second * time.Duration(duration))
	return t.Format("15:04:05")
}

func (state *DurationMinutes) UnmarshalJSON(data []byte) error {
	var rawValue timeSpan
	err := json.Unmarshal(data, &rawValue)
	if err != nil {
		return err
	}

	*state = DurationMinutes(int(rawValue.TotalMinutes))
	return nil
}

func (state *DurationSeconds) UnmarshalJSON(data []byte) error {
	var rawValue timeSpan
	err := json.Unmarshal(data, &rawValue)
	if err != nil {
		return err
	}

	*state = DurationSeconds(int(rawValue.TotalSeconds))
	return nil
}
