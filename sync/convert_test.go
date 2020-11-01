package sync

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/yqt/go-garmin2suunto/garmin"
	"os"
	"testing"
)

var (
	email    = ""
	password = ""
)

func TestConvert(t *testing.T) {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	err := garmin.Auth(email, password)
	assert.Nil(t, err)

	activityId := int64(12345678)

	activity, err := garmin.GetActivity(activityId)
	assert.Nil(t, err)

	activityDetails, err := garmin.GetActivityDetails(activityId)
	assert.Nil(t, err)

	activitySplits, err := garmin.GetActivitySplits(activityId)
	assert.Nil(t, err)

	move, err := Convert(activity, activityDetails, activitySplits, "test")
	assert.Nil(t, err)

	moveBytes, err := json.Marshal(move)
	assert.Nil(t, err)

	logrus.WithFields(logrus.Fields{
		"move": string(moveBytes),
	}).Info()
}
