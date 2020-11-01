package suunto

type MoveMark struct {
	DistanceFromPrevious  int    `json:"DistanceFromPrevious"`
	SplitTimeFromPrevious int    `json:"SplitTimeFromPrevious"`
	Type                  int    `json:"Type"`
	Notes                 string `json:"Notes,omitempty"`

	AvgCadence float32 `json:"AvgCadence,omitempty"`
	AvgHR      int     `json:"AvgHR,omitempty"`
	AvgSpeed   float32 `json:"AvgSpeed,omitempty"`
	MaxCadence float32 `json:"MaxCadence,omitempty"`
	MaxHR      int     `json:"MaxHR,omitempty"`
	MaxPower   float32 `json:"MaxPower"`
	MaxSpeed   float32 `json:"MaxSpeed,omitempty"`
	MinHR      int     `json:"MinHR,omitempty"`
	MinPower   float32 `json:"MinPower"`
	MinSpeed   float32 `json:"MinSpeed"`
}
