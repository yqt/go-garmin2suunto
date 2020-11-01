package suunto

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	email   = ""
	userKey = ""
)

func TestGetMoves(t *testing.T) {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	moveItems, err := GetMoveItems(email, userKey, "", -1)
	assert.Nil(t, err)

	moveItemsBytes, err := json.Marshal(moveItems)

	logrus.WithFields(logrus.Fields{
		"moveItems": string(moveItemsBytes),
	}).Info()
}
