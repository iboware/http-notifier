//go:build unit

package notification

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotificationTask_Execute(t *testing.T) {
	tests := []struct {
		name                   string
		wantErr                bool
		message                string
		expectedServerResponse int
	}{
		{
			name:                   "want success",
			wantErr:                false,
			message:                "foo",
			expectedServerResponse: http.StatusOK,
		},
		{
			name:                   "500 error from the server",
			wantErr:                true,
			message:                "foo",
			expectedServerResponse: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// generate a test server so we can capture and inspect the request
			testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(tt.expectedServerResponse)
				res.Write([]byte(""))
			}))
			defer testServer.Close()
			nt := &NotificationTask{
				url:        testServer.URL,
				message:    tt.message,
				httpClient: testServer.Client(),
				failure:    nil,
			}
			err := nt.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("NotificationTask.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
