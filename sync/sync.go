package sync

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yqt/go-garmin2suunto/config"
	"github.com/yqt/go-garmin2suunto/garmin"
	"github.com/yqt/go-garmin2suunto/suunto"
	"strconv"
	"strings"
	"time"
)

type UserInfo struct {
	Suunto suunto.UserInfo `json:"suunto"`
	Garmin garmin.UserInfo `json:"garmin"`
}

var (
	SupportedGarminSportType = map[string]bool{
		"treadmill_running": true,
		"running":           true,
	}
)

func SynchronizeLatestActivities(userInfo UserInfo) (bool, error) {
	currentTime := time.Now().AddDate(0, 0, -1)
	currentDate := currentTime.Format("2006-01-02")

	err := garmin.Auth(userInfo.Garmin.Email, userInfo.Garmin.Password)
	if err != nil {
		return false, err
	}

	activityItems, err := garmin.GetActivityItems(0, 3, currentDate)
	if err != nil {
		return false, err
	}

	if len(activityItems) == 0 {
		return true, errors.New("empty activities")
	}

	existedMoveItems, err := suunto.GetMoveItems(userInfo.Suunto.Email, userInfo.Suunto.UserKey, currentDate, 10)
	if err != nil {
		return false, err
	}

	succeedActivityIds := make([]int64, 0)
	failedActivityIds := make([]int64, 0)
	skippedActivityIds := make([]int64, 0)

	for _, activityItem := range activityItems {
		// TODO: support cycling and swimming
		if _, found := SupportedGarminSportType[activityItem.ActivityType.TypeKey]; !found {
			skippedActivityIds = append(skippedActivityIds, activityItem.ActivityId)
			continue
		}

		skip := false
		for _, moveItem := range existedMoveItems {
			garminActivityId := moveItem.GetGarminActivityId()
			if strings.EqualFold(strconv.FormatInt(activityItem.ActivityId, 10), garminActivityId) {
				skip = true
				break
			}
		}
		if skip {
			skippedActivityIds = append(skippedActivityIds, activityItem.ActivityId)
			continue
		}

		moveResult, err := syncSingleActivity(userInfo, activityItem.ActivityId)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"activityId": activityItem.ActivityId,
				"moveResult": moveResult,
				"err":        err,
			}).Error("sync single activity failed")
			failedActivityIds = append(failedActivityIds, activityItem.ActivityId)
			continue
		}
		succeedActivityIds = append(succeedActivityIds, activityItem.ActivityId)
	}

	logrus.WithFields(logrus.Fields{
		"succeedActivityIds": succeedActivityIds,
		"failedActivityIds":  failedActivityIds,
		"skippedActivityIds": skippedActivityIds,
		"err":                err,
	}).Debug("sync detail")

	suc := true
	if len(succeedActivityIds) == 0 && len(failedActivityIds) != 0 {
		suc = false
	}
	return suc, errors.New(
		fmt.Sprintf(
			"id[%v] succeeded. id[%v] failed. id[%v] skipped.",
			succeedActivityIds, failedActivityIds, skippedActivityIds))
}

func syncSingleActivity(userInfo UserInfo, activityId int64) (suunto.MoveResult, error) {
	var ret suunto.MoveResult
	activity, err := garmin.GetActivity(activityId)
	if err != nil {
		return ret, err
	}

	activityDetail, err := garmin.GetActivityDetails(activityId)
	if err != nil {
		return ret, err
	}

	activitySplit, err := garmin.GetActivitySplits(activityId)
	if err != nil {
		return ret, err
	}

	move, err := Convert(activity, activityDetail, activitySplit, config.SourcePrefix)
	if err != nil {
		return ret, err
	}

	return suunto.SaveMove(userInfo.Suunto.Email, userInfo.Suunto.UserKey, move)
}
