package suunto

import (
	"encoding/json"
	"github.com/yqt/go-garmin2suunto/util"
)

type UserInfo struct {
	Email   string `json:"email"`
	UserKey string `json:"user_key"`
}

const (
	ApiServicePrefix = "https://uiservices.movescount.com"
	AppKey           = "uFiPE28bwLykgnTlYyvlS7GzgaAcIRP3I85FCMFJUDFwTa7hcAihvk7x9SJEC5CP"
)

var (
	requestClient = util.NewCookieRequest()
)

func GetMoveItems(email string, userKey string, startDate string, maxCount int) ([]MoveItem, error) {
	uri := ApiServicePrefix + "/moves/private"
	params := map[string]interface{}{
		"appkey":  AppKey,
		"email":   email,
		"userkey": userKey,
	}
	if startDate != "" {
		params["startDate"] = startDate
	}
	if maxCount > 0 {
		params["maxcount"] = maxCount
	}

	moveItems := make([]MoveItem, 0)
	err := requestClient.GetJson(uri, params, &moveItems)
	if err != nil {
		return moveItems, err
	}

	return moveItems, nil
}

func SaveMove(email string, userKey string, move Move) (MoveResult, error) {
	uri := ApiServicePrefix + "/moves"
	moveResult := MoveResult{}
	params := map[string]interface{}{
		"appkey":  AppKey,
		"email":   email,
		"userkey": userKey,
	}
	moveData, err := json.Marshal(move)
	if err != nil {
		return moveResult, nil
	}
	err = requestClient.PostJson(uri, params, nil, moveData, true, &moveResult)
	if err != nil {
		return moveResult, err
	}

	return moveResult, nil
}
