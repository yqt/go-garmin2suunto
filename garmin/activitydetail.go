package garmin

import (
	"strings"
	"time"
)

type ActivityDetail struct {
	ActivityId            int64                  `json:"activityId"`
	DetailsAvailable      bool                   `json:"detailsAvailable"`
	MeasurementCount      int                    `json:"measurementCount"`
	MetricsCount          int                    `json:"metricsCount"`
	MetricDescriptors     []MetricDescriptor     `json:"metricDescriptors"`
	ActivityDetailMetrics []ActivityDetailMetric `json:"activityDetailMetrics"`
	GeoPolyline           GeoPolyline            `json:"geoPolylineDTO"`
}

type MetricDescriptor struct {
	MetricsIndex int    `json:"metricsIndex"`
	Key          string `json:"key"`
	Unit         Unit   `json:"unit"`
}

type Unit struct {
	ID     int     `json:"id"`
	Key    string  `json:"key"`
	Factor float64 `json:"factor"`
}

type ActivityDetailMetric struct {
	Metrics []float64 `json:"metrics"`
}

func (a *ActivityDetailMetric) GetCadence(index *MetricIndex) float64 {
	return a.Metrics[index.RunCadence]
}

func (a *ActivityDetailMetric) GetRunningCadence(index *MetricIndex) float64 {
	return a.Metrics[index.DoubleCadence]
}

func (a *ActivityDetailMetric) GetSpeed(index *MetricIndex) float64 {
	return a.Metrics[index.Speed]
}

func (a *ActivityDetailMetric) GetDistance(index *MetricIndex) float64 {
	return a.Metrics[index.Distance]
}

func (a *ActivityDetailMetric) GetHeartRate(index *MetricIndex) float64 {
	return a.Metrics[index.HearRate]
}

func (a *ActivityDetailMetric) GetLocalTime(index *MetricIndex) string {
	idx := index.Time
	if idx < 0 {
		return ""
	}
	if a.Metrics[idx] == 0 {
		return ""
	}
	ts := int64(a.Metrics[idx])
	tmp := time.Unix(0, ts*int64(time.Millisecond)).Format(time.RFC3339)
	tailStartPos := strings.Index(tmp, "+")
	if tailStartPos != -1 {
		tmp = tmp[:tailStartPos]
	}
	return tmp
}

type MetricIndex struct {
	RunCadence    int `json:"runCadence"`
	DoubleCadence int `json:"doubleCadence"`
	Speed         int `json:"speed"`
	Time          int `json:"time"`
	HearRate      int `json:"hearRate"`
	Distance      int `json:"distance"`
}

func NewMetricIndex(descriptors []MetricDescriptor) *MetricIndex {
	m := &MetricIndex{
		RunCadence:    -1,
		DoubleCadence: -1,
		Speed:         -1,
		Time:          -1,
		HearRate:      -1,
		Distance:      -1,
	}
	if len(descriptors) == 0 {
		return m
	}
	for _, descriptor := range descriptors {
		switch descriptor.Key {
		case "directRunCadence":
			m.RunCadence = descriptor.MetricsIndex
		case "directHeartRate":
			m.HearRate = descriptor.MetricsIndex
		case "directSpeed":
			m.Speed = descriptor.MetricsIndex
		case "directDoubleCadence":
			m.DoubleCadence = descriptor.MetricsIndex
		case "directTimestamp":
			m.Time = descriptor.MetricsIndex
		case "sumDistance":
			m.Distance = descriptor.MetricsIndex
		}
	}

	return m
}

type GeoPolyline struct {
	Polyline []Track `json:"polyline"`
}

type Track struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Altitude  float64 `json:"altitude"`
	Speed     float64 `json:"speed"`
	Time      int64   `json:"time"`
}

func (t *Track) GetTimeLocale() string {
	if t.Time == 0 {
		return ""
	}
	tmp := time.Unix(0, t.Time*int64(time.Millisecond)).Format(time.RFC3339)
	tailStartPos := strings.Index(tmp, "+")
	if tailStartPos != -1 {
		tmp = tmp[:tailStartPos]
	}
	return tmp
}
