package sync

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/yqt/go-garmin2suunto/garmin"
	"github.com/yqt/go-garmin2suunto/suunto"
	"os"
	"testing"
)

var (
	userInfo = UserInfo{
		Suunto: suunto.UserInfo{
			Email:   "",
			UserKey: "",
		},
		Garmin: garmin.UserInfo{
			Email:    "",
			Password: "",
		},
	}
)

func TestSynchronizeLatestActivities(t *testing.T) {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	suc, err := SynchronizeLatestActivities(userInfo)
	assert.True(t, suc)
	logrus.WithFields(logrus.Fields{
		"err": err,
	}).Info()
}

func TestSaveActivities(t *testing.T) {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	err := garmin.Auth(userInfo.Garmin.Email, userInfo.Garmin.Password)
	if !assert.Nil(t, err) {
		return
	}

	activityId := int64(12345678)

	activity, err := garmin.GetActivity(activityId)
	assert.Nil(t, err)

	activityDetails, err := garmin.GetActivityDetails(activityId)
	assert.Nil(t, err)

	activitySplits, err := garmin.GetActivitySplits(activityId)
	assert.Nil(t, err)

	move, err := Convert(activity, activityDetails, activitySplits, "")
	assert.Nil(t, err)

	moveResult, err := suunto.SaveMove(userInfo.Suunto.Email, userInfo.Suunto.UserKey, move)
	assert.Nil(t, err)

	moveResultBytes, err := json.Marshal(moveResult)
	assert.Nil(t, err)

	logrus.WithFields(logrus.Fields{
		"moveResult": string(moveResultBytes),
	}).Info()
}
