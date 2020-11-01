package garmin

type ActivityItem struct {
	ActivityId                            int64        `json:"activityId"`
	ActivityName                          string       `json:"activityName"`
	StartTimeLocal                        string       `json:"startTimeLocal"`
	ActivityType                          ActivityType `json:"activityType"`
	Calories                              float64      `json:"calories"`
	AverageHR                             float64      `json:"averageHR"`
	MaxHR                                 float64      `json:"maxHR"`
	AverageRunningCadenceInStepsPerMinute float64      `json:"averageRunningCadenceInStepsPerMinute"`
	MaxRunningCadenceInStepsPerMinute     float64      `json:"maxRunningCadenceInStepsPerMinute"`
	Distance                              float64      `json:"distance"`
	Duration                              float64      `json:"duration"`
	AverageSpeed                          float64      `json:"averageSpeed"`
	MaxSpeed                              float64      `json:"maxSpeed"`
	StartLatitude                         float64      `json:"startLatitude"`
	StartLongitude                        float64      `json:"startLongitude"`
	HasPolyline                           bool         `json:"hasPolyline"`
}

func (a *ActivityItem) Equals(obj interface{}) bool {
	if a == obj {
		return true
	}

	obj1, ok := obj.(ActivityItem)
	if !ok {
		return false
	}

	return a.ActivityId == obj1.ActivityId
}
