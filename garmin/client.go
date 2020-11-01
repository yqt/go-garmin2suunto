package garmin

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/yqt/go-garmin2suunto/util"
	"regexp"
	"strconv"
	"strings"
)

type UserInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

const (
	ApiServiceHost   = "connect.garmin.com"
	ApiServicePrefix = "https://" + ApiServiceHost
	SsoPrefix        = "https://sso.garmin.com"
	UserAgent        = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36"
)

var (
	requestClient = util.NewCookieRequest()
)

func Auth(email string, password string) error {
	params := map[string]interface{}{
		"service":  ApiServicePrefix + "/modern",
		"clientId": "GarminConnect",
		//"gauthHost": ApiServicePrefix + "/modern",
		"gauthHost": SsoPrefix + "/sso",
		//"generateExtraServiceTicket": "true",
		"consumeServiceTicket": "false",
	}

	uri := SsoPrefix + "/sso/signin"

	headers := map[string]string{
		"User-Agent": UserAgent,
	}
	requestClient.SetHeaders(headers)

	respText, err := requestClient.Get(uri, nil)
	if err != nil {
		return err
	}
	csrfToken, err := extractCSRFToken(respText)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"csrfToken": csrfToken,
	}).Debug()

	formData := map[string]interface{}{
		"username": email,
		"password": password,
		"embed":    "false",
		"_csrf":    csrfToken,
	}
	headers["Origin"] = SsoPrefix
	requestClient.SetHeaders(headers)
	respText, err = requestClient.Post(uri, params, formData, nil, false)

	ticketUrl, err := extractTicketUrl(respText)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"ticketUrl": ticketUrl,
	}).Debug()

	respText, err = requestClient.Get(ticketUrl, nil)
	socialProfileText, err := extractSocialProfile(respText)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"socialProfileText": socialProfileText,
	}).Debug()

	return nil
}

func GetActivity(id int64) (Activity, error) {
	uri := ApiServicePrefix + "/modern/proxy/activity-service/activity/" + strconv.FormatInt(id, 10)
	activity := Activity{}
	err := requestClient.GetJson(uri, nil, &activity)
	if err != nil {
		return activity, err
	}

	return activity, nil
}

func GetActivityDetails(id int64) (ActivityDetail, error) {
	uri := ApiServicePrefix + "/modern/proxy/activity-service/activity/" + strconv.FormatInt(id, 10) + "/details"
	activityDetail := ActivityDetail{}
	err := requestClient.GetJson(uri, nil, &activityDetail)
	if err != nil {
		return activityDetail, err
	}

	return activityDetail, nil
}

func GetActivitySplits(id int64) (ActivitySplit, error) {
	uri := ApiServicePrefix + "/modern/proxy/activity-service/activity/" + strconv.FormatInt(id, 10) + "/splits"
	activitySplit := ActivitySplit{}
	err := requestClient.GetJson(uri, nil, &activitySplit)
	if err != nil {
		return activitySplit, err
	}

	return activitySplit, nil
}

func GetActivityItems(start int, limit int, startDate string) ([]ActivityItem, error) {
	uri := ApiServicePrefix + "/modern/proxy/activitylist-service/activities/search/activities"
	activityItems := make([]ActivityItem, 0)
	params := map[string]interface{}{
		"start": start,
		"limit": limit,
	}
	if startDate != "" {
		params["startDate"] = startDate
	}
	err := requestClient.GetJson(uri, params, &activityItems)
	if err != nil {
		return activityItems, err
	}

	return activityItems, nil
}

func extractCSRFToken(respText string) (string, error) {
	fragment := `<input type="hidden" name="_csrf" value="`
	startPos := strings.Index(respText, fragment)
	if startPos == -1 {
		return "", errors.New("CSRF token not found")
	}
	restText := respText[startPos:]
	endPos := strings.Index(restText, `" />`)
	if endPos == -1 {
		return "", errors.New("invalid CSRF token end")
	}
	restText = restText[len(fragment):endPos]
	return restText, nil
}

func extractTicketUrl(respText string) (string, error) {
	t := regexp.MustCompile(`https:\\\/\\\/` + ApiServiceHost + `\\\/modern(\\\/)?\?ticket=(([a-zA-Z0-9]|-)*)`)
	ticketUrl := t.FindString(respText)

	// NOTE: undo escaping
	ticketUrl = strings.Replace(ticketUrl, "\\/", "/", -1)

	if ticketUrl == "" {
		return "", errors.New("wrong credentials")
	}

	return ticketUrl, nil
}

func extractSocialProfile(respText string) (string, error) {
	fragment := `window.VIEWER_SOCIAL_PROFILE = JSON.parse("`
	startPos := strings.Index(respText, fragment)
	if startPos == -1 {
		return "", errors.New("social profile not found")
	}
	restText := respText[startPos:]
	endPos := strings.Index(restText, `");`)
	if endPos == -1 {
		return "", errors.New("invalid social profile end")
	}
	restText = restText[len(fragment):endPos]
	restText = strings.Replace(restText, "\\", "", -1)
	return restText, nil
}
