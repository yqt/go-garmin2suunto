package garmin

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/yqt/go-garmin2suunto/config"
	"os"
	"testing"
)

var (
	email    = config.GarminEmail
	password = config.GarminPassword
)

func TestAuth(t *testing.T) {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	err := Auth(email, password)
	assert.Nil(t, err)
}

func TestGetActivity(t *testing.T) {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	err := Auth(email, password)
	assert.Nil(t, err)

	activity, err := GetActivity(12345678)
	assert.Nil(t, err)
	activityBytes, err := json.Marshal(activity)
	assert.Nil(t, err)

	logrus.WithFields(logrus.Fields{
		"activity": string(activityBytes),
	}).Info()
}

func TestGetActivityDetails(t *testing.T) {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	err := Auth(email, password)
	assert.Nil(t, err)

	activityDetails, err := GetActivityDetails(12345678)
	assert.Nil(t, err)
	activityDetailsBytes, err := json.Marshal(activityDetails)
	assert.Nil(t, err)

	logrus.WithFields(logrus.Fields{
		"activityDetails": string(activityDetailsBytes),
	}).Info()
}

func TestGetActivitySplits(t *testing.T) {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	err := Auth(email, password)
	assert.Nil(t, err)

	activitySplits, err := GetActivitySplits(12345678)
	assert.Nil(t, err)
	activitySplitsBytes, err := json.Marshal(activitySplits)
	assert.Nil(t, err)

	logrus.WithFields(logrus.Fields{
		"activitySplits": string(activitySplitsBytes),
	}).Info()
}

func TestGetActivityItems(t *testing.T) {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	err := Auth(email, password)
	assert.Nil(t, err)

	activityItems, err := GetActivityItems(0, 10, "2020-10-29T23:02:12+08:00")
	assert.Nil(t, err)
	activityItemsBytes, err := json.Marshal(activityItems)
	assert.Nil(t, err)

	logrus.WithFields(logrus.Fields{
		"activityItems": string(activityItemsBytes),
	}).Info()
}
