package suunto

type Move struct {
	ActivityID        int64   `json:"ActivityID"`
	AscentAltitude    float64 `json:"AscentAltitude,omitempty"`
	AvgBikeCadence    float32 `json:"AvgBikeCadence,omitempty"`
	AvgRunningCadence float32 `json:"AvgRunningCadence,omitempty"`
	AvgHR             int     `json:"AvgHR,omitempty"`
	AvgSpeed          float32 `json:"AvgSpeed,omitempty"`
	AvgTemp           float32 `json:"AvgTemp,omitempty"`
	AvgPower          float32 `json:"AvgPower,omitempty"`
	DescentAltitude   float64 `json:"DescentAltitude,omitempty"`
	DescentTime       float64 `json:"DescentTime,omitempty"`
	Distance          int     `json:"Distance,omitempty"`
	Duration          float64 `json:"Duration,omitempty"`
	Energy            int     `json:"Energy,omitempty"`
	Feeling           int     `json:"Feeling,omitempty"`
	FlatTime          float64 `json:"FlatTime,omitempty"`
	HighAltitude      float32 `json:"HighAltitude,omitempty"`
	LastModifiedDate  string  `json:"LastModifiedDate,omitempty"`
	LocalStartTime    string  `json:"LocalStartTime,omitempty"`
	LowAltitude       float32 `json:"LowAltitude,omitempty"`
	MaxBikeCadence    float32 `json:"MaxBikeCadence,omitempty"`
	MaxRunningCadence float32 `json:"MaxRunningCadence,omitempty"`
	MaxSpeed          float32 `json:"MaxSpeed,omitempty"`
	MaxTemp           float32 `json:"MaxTemp,omitempty"`
	MinTemp           float32 `json:"MinTemp,omitempty"`
	MaxPower          float32 `json:"MaxPower"`
	MemberID          int     `json:"MemberID,omitempty"`
	MinHR             int     `json:"MinHR,omitempty"`
	MIntegeremp       float32 `json:"MIntegeremp,omitempty"`
	MoveID            int64   `json:"MoveID,omitempty"`
	Notes             string  `json:"Notes,omitempty"`
	PeakHR            int     `json:"PeakHR,omitempty"`
	SessionName       string  `json:"SessionName,omitempty"`
	StartLatitude     float64 `json:"StartLatitude,omitempty"`
	StartLongitude    float64 `json:"StartLongitute,omitempty"`
	Tags              string  `json:"Tags,omitempty"`
	TrainingEffect    float32 `json:"TrainingEffect,omitempty"`

	Weather            int        `json:"Weather,omitempty"`
	DeviceName         string     `json:"DeviceName,omitempty"`
	DeviceSerialNumber string     `json:"DeviceSerialNumber,omitempty"`
	Samples            Samples    `json:"Samples"`
	Track              Track      `json:"Track"`
	Marks              []MoveMark `json:"Marks"`

	AvgCadence         float32 `json:"AvgCadence"`
	MaxCadence         float32 `json:"MaxCadence"`
	PeakTrainingEffect float32 `json:"PeakTrainingEffect"`
	Source             string  `json:"Source"`
}

type Samples struct {
	SampleSets []SampleSet `json:"SampleSets"`
}

type Track struct {
	TrackPoints []TrackPoint `json:"TrackPoints"`
}
