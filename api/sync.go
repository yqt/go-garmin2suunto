package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/yqt/go-garmin2suunto/config"
	"github.com/yqt/go-garmin2suunto/garmin"
	"github.com/yqt/go-garmin2suunto/suunto"
	"github.com/yqt/go-garmin2suunto/sync"
	"net/http"
)

func genSyncHandler(c *gin.Context) {
	userInfo := sync.UserInfo{
		Suunto: suunto.UserInfo{
			Email:   config.SuuntoEmail,
			UserKey: config.SuuntoUserKey,
		},
		Garmin: garmin.UserInfo{
			Email:    config.GarminEmail,
			Password: config.GarminPassword,
		},
	}
	suc, err := sync.SynchronizeLatestActivities(userInfo)

	logrus.WithFields(logrus.Fields{
		"suc": suc,
		"err": err,
	}).Info("sync result")

	c.JSON(http.StatusOK, gin.H{
		"success": suc,
		"message": fmt.Sprintf("%v", err),
	})
}
