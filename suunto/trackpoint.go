package suunto

type TrackPoint struct {
	Altitude  float32 `json:"Altitude,omitempty"`
	Bearing   float32 `json:"Bearing,omitempty"`
	EHPE      float32 `json:"EHPE,omitempty"`
	Latitude  float64 `json:"Latitude,omitempty"`
	LocalTime string  `json:"LocalTime,omitempty"`
	Longitude float64 `json:"Longitude,omitempty"`
	Speed     float64 `json:"Speed,omitempty"`
}
