package notification

import (
	"fmt"
	"net/http"
	"strings"
)

type NotificationTask struct {
	url        string
	message    string
	httpClient *http.Client
	failure    func(error)
}

func NewNotificationTask(url, message string, httpClient *http.Client, failure func(error)) *NotificationTask {
	return &NotificationTask{
		url:        url,
		message:    message,
		httpClient: httpClient,
		failure:    failure,
	}
}

// Execute performs the work
func (nt *NotificationTask) Execute() error {
	req, err := http.NewRequest("POST", nt.url, strings.NewReader(nt.message))
	if err != nil {
		return err
	}
	response, err := nt.httpClient.Do(req)
	if err != nil {
		return err
	}

	statusOK := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOK {
		return fmt.Errorf("Non-OK HTTP status:%d", response.StatusCode)
	}
	return nil
}

// OnFailure handles any error returned from Execute()
func (nt *NotificationTask) OnFailure(err error) {
	nt.failure(err)
}
