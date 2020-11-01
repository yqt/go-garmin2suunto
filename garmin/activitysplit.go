package garmin

type ActivitySplit struct {
	ActivityId int64 `json:"activityId"`
	Laps       []Lap `json:"lapDTOs"`
}

type Lap struct {
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`

	AverageRunCadence float64 `json:"averageRunCadence"`
	MaxRunCadence     float64 `json:"maxRunCadence"`
	AverageHR         float64 `json:"averageHR"`
	AverageSpeed      float64 `json:"averageSpeed"`
	MaxHR             float64 `json:"maxHR"`
	MaxSpeed          float64 `json:"maxSpeed"`
	Calories          float64 `json:"calories"`
}
