package sync

import (
	"bytes"
	"github.com/yqt/go-garmin2suunto/garmin"
	"github.com/yqt/go-garmin2suunto/suunto"
	"math"
	"sort"
	"strconv"
	"strings"
)

const (
	InvalidHeartRateHigh = 300
	InvalidSpeedHigh     = 1000
	SourceDefaultPrefix  = "go-garmin2suunto "
)

func Convert(
	activity garmin.Activity,
	activityDetail garmin.ActivityDetail,
	activitySplit garmin.ActivitySplit,
	sourcePrefix string) (suunto.Move, error) {

	move := suunto.Move{}

	move.Notes = genNotes(activity)

	move.LocalStartTime = activity.Summary.StartTimeLocal
	move.StartLatitude = activity.Summary.StartLatitude
	move.StartLongitude = activity.Summary.StartLongitude
	move.AscentAltitude = activity.Summary.ElevationGain
	move.DescentAltitude = activity.Summary.ElevationLoss

	move.Distance = int(activity.Summary.Distance)
	move.Duration = activity.Summary.MovingDuration
	if move.Duration == 0 {
		move.Duration = activity.Summary.Duration
	}

	move.AvgHR = int(activity.Summary.AverageHR)
	move.PeakHR = int(activity.Summary.MaxHR)

	move.Energy = int(activity.Summary.Calories)

	move.HighAltitude = float32(activity.Summary.MaxElevation)
	move.LowAltitude = float32(activity.Summary.MinElevation)

	move.MaxSpeed = float32(activity.Summary.MaxSpeed)
	move.AvgSpeed = float32(activity.Summary.AverageSpeed)

	move.MaxPower = float32(activity.Summary.MaxPower)
	move.AvgPower = float32(activity.Summary.AveragePower)

	move.AvgTemp = float32(activity.Summary.AverageTemperature)
	move.MinTemp = float32(activity.Summary.MinTemperature)
	move.MaxTemp = float32(activity.Summary.MaxTemperature)

	// NOTE: cadence in Garmin activity summary is doubled
	move.AvgBikeCadence = float32(activity.Summary.AverageBikeCadence / 2)
	move.MaxBikeCadence = float32(activity.Summary.MaxBikeCadence / 2)

	if move.AvgBikeCadence == 0 {
		move.AvgBikeCadence = float32(activity.Summary.AverageRunCadence / 2)
	}
	if move.MaxBikeCadence == 0 {
		move.MaxBikeCadence = float32(activity.Summary.MaxRunCadence / 2)
	}

	move.AvgRunningCadence = float32(activity.Summary.AverageRunCadence / 2)
	move.MaxRunningCadence = float32(activity.Summary.MaxRunCadence / 2)

	move.TrainingEffect = float32(activity.Summary.TrainingEffect)
	move.PeakTrainingEffect = float32(activity.Summary.TrainingEffect)

	move.AvgCadence = float32(activity.Summary.AverageRunCadence / 2)
	move.MaxCadence = float32(activity.Summary.MaxRunCadence / 2)

	if move.AvgCadence == 0 {
		move.AvgCadence = float32(activity.Summary.AverageBikeCadence / 2)
	}
	if move.MaxCadence == 0 {
		move.MaxCadence = float32(activity.Summary.MaxBikeCadence / 2)
	}

	// NOTE: sequence bellow cares
	move.Samples = convertToSamples(activityDetail, &move)
	move.Marks = convertToMoveMarks(activitySplit, &move)
	move.Track = convertToTrack(activityDetail)

	move.ActivityID = int64(getActivityType(activity, len(move.Track.TrackPoints) != 0))

	if sourcePrefix == "" {
		sourcePrefix = SourceDefaultPrefix
	}
	move.Source = sourcePrefix + getDeviceName(activity)

	return move, nil
}

func genNotes(activity garmin.Activity) string {
	var notes bytes.Buffer
	notes.WriteString(suunto.NoteActivityIdBeginTag)
	notes.WriteString(strconv.FormatInt(activity.ActivityId, 10))
	notes.WriteString(suunto.NoteActivityIdEndTag)
	notes.WriteString("\n")

	if strings.TrimSpace(activity.ActivityName) != "" {
		notes.WriteString(activity.ActivityName)
		notes.WriteString("\n")
	}

	if strings.TrimSpace(activity.Description) != "" {
		notes.WriteString(activity.Description)
		notes.WriteString("\n")
	}

	return notes.String()
}

func convertToMoveMarks(split garmin.ActivitySplit, move *suunto.Move) []suunto.MoveMark {
	moveMarks := make([]suunto.MoveMark, 0)
	moveDistanceInterval := make([]int, 0)
	sumDistance := 0
	for _, lap := range split.Laps {
		moveMark := suunto.MoveMark{
			Type:                  0,
			SplitTimeFromPrevious: int(lap.Duration),
			DistanceFromPrevious:  int(lap.Distance),

			AvgCadence: float32(lap.AverageRunCadence),
			AvgHR:      int(lap.AverageHR),
			AvgSpeed:   float32(lap.AverageSpeed),
			MaxCadence: float32(lap.MaxRunCadence),
			MaxHR:      int(lap.MaxHR),
			MaxSpeed:   float32(lap.MaxSpeed),
		}
		sumDistance += moveMark.DistanceFromPrevious
		moveDistanceInterval = append(moveDistanceInterval, sumDistance)
		moveMarks = append(moveMarks, moveMark)
	}

	// NOTE: get minHR etc of each lap
	targetIdx := 0
	lapStats := make(map[int]map[string]float32)
	for _, ss := range move.Samples.SampleSets {
		for idx, v := range moveDistanceInterval[targetIdx:] {
			if ss.Distance > v {
				targetIdx += idx + 1
				break
			}
			break
		}
		if _, found := lapStats[targetIdx]; !found {
			lapStats[targetIdx] = map[string]float32{
				"MinHR":    InvalidHeartRateHigh,
				"MinSpeed": InvalidSpeedHigh,
			}
		}

		hr := float32(ss.HeartRate)
		if hr != 0 && lapStats[targetIdx]["MinHR"] > hr {
			lapStats[targetIdx]["MinHR"] = hr
		}
		if ss.Speed != 0 && lapStats[targetIdx]["MinSpeed"] > ss.Speed {
			lapStats[targetIdx]["MinSpeed"] = ss.Speed
		}
	}

	for idx, _ := range moveMarks {
		if _, found := lapStats[idx]; found {
			moveMarks[idx].MinHR = int(lapStats[idx]["MinHR"])
			moveMarks[idx].MinSpeed = lapStats[idx]["MinSpeed"]
		}
	}

	return moveMarks
}

func convertToSamples(detail garmin.ActivityDetail, move *suunto.Move) suunto.Samples {
	sampleSets := make([]suunto.SampleSet, 0)

	startTs := int64(math.MaxInt64)
	endTs := int64(math.MinInt64)

	metricIndex := garmin.NewMetricIndex(detail.MetricDescriptors)
	minHR := InvalidHeartRateHigh
	for _, metric := range detail.ActivityDetailMetrics {
		cadence := float32(metric.GetCadence(metricIndex))
		ts := metric.GetTimeStamp(metricIndex)
		if ts < startTs {
			startTs = ts
		}
		if ts > endTs {
			endTs = ts
		}
		ss := suunto.SampleSet{
			BikeCadence: cadence,
			//Cadence:       cadence,
			// NOTE:
			//RunningCadence: cadence,
			Distance:  int(math.Round(metric.GetDistance(metricIndex))),
			HeartRate: int(math.Round(metric.GetHeartRate(metricIndex))),
			LocalTime: metric.GetLocalTime(metricIndex),
			Speed:     float32(metric.GetSpeed(metricIndex)),
			Timestamp: ts,
		}
		sampleSets = append(sampleSets, ss)

		// NOTE: get min heart rate from samples/metrics
		if ss.HeartRate != 0 && ss.HeartRate < minHR {
			minHR = ss.HeartRate
		}
	}

	// NOTE: set minHR into move
	if minHR != InvalidHeartRateHigh && minHR < move.PeakHR {
		move.MinHR = minHR
	}
	// NOTE: set duration if not set in details
	if move.Duration == 0 {
		move.Duration = float64(endTs-startTs) / 1000
	}

	sort.SliceStable(sampleSets, func(i, j int) bool {
		return sampleSets[i].Timestamp < sampleSets[j].Timestamp
	})

	samples := suunto.Samples{
		SampleSets: sampleSets,
	}

	return samples
}

func convertToTrack(detail garmin.ActivityDetail) suunto.Track {
	trackPoints := make([]suunto.TrackPoint, 0)

	for _, track := range detail.GeoPolyline.Polyline {
		tp := suunto.TrackPoint{
			Altitude:  float32(track.Altitude),
			Latitude:  track.Latitude,
			Longitude: track.Longitude,
			LocalTime: track.GetTimeLocale(),
			Speed:     track.Speed,
			Timestamp: track.Time,
		}

		trackPoints = append(trackPoints, tp)
	}

	sort.SliceStable(trackPoints, func(i, j int) bool {
		return trackPoints[i].Timestamp < trackPoints[j].Timestamp
	})

	track := suunto.Track{
		TrackPoints: trackPoints,
	}

	return track
}

func getActivityType(activity garmin.Activity, hasTrack bool) int {
	garminActivityType := activity.ActivityType.TypeKey
	switch garminActivityType {
	case "running":
		if hasTrack {
			return suunto.Run
		} else {
			return suunto.Treadmill
		}
	case "cycling":
		return suunto.Cycling
	case "indoor_cycling":
		return suunto.IndoorCycling
	case "treadmill_running":
		return suunto.Treadmill
	case "lap_swimming":
		return suunto.Swimming
	case "open_water_swimming":
		return suunto.OpenWaterSwimming
	default:
		return suunto.NotSpecifiedSport
	}
}

func getDeviceName(activity garmin.Activity) string {
	deviceName := activity.MetadataDTO.Manufacturer
	for _, sensor := range activity.MetadataDTO.Sensors {
		if sensor.SKU != "" {
			deviceID, found := garmin.DevicePN2ID[sensor.SKU]
			if !found {
				continue
			}
			deviceName += " - " + strings.ToUpper(deviceID)
			break
		}
	}
	return deviceName
}
