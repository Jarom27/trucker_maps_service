package models

import "time"

type GPSData struct {
	Device_id string
	Latitude  float64
	Longitude float64
	Altitude  float64
	Timestamp time.Time
}
