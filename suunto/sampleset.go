package suunto

type SampleSet struct {
	Altitude          float32 `json:"Altitude,omitempty"`
	BikeCadence       float32 `json:"BikeCadence"`
	Distance          int     `json:"Distance"`
	EnergyConsumption float32 `json:"EnergyConsumption"`
	HeartRate         int     `json:"HeartRate,omitempty"`
	LocalTime         string  `json:"LocalTime,omitempty"`
	Power             int     `json:"Power"`
	RunningCadence    float32 `json:"RunningCadence,omitempty"`
	Speed             float32 `json:"Speed"`
	Temperature       float32 `json:"Temperature,omitempty"`
	VerticalSpeed     float32 `json:"VerticalSpeed,omitempty"`

	Cadence float32 `json:"Cadence"`
	Epoc    float32 `json:"Epoc,omitempty"`
}
