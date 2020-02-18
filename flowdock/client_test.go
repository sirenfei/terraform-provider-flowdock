package flowdock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Flowdock_Client_NewClient(t *testing.T) {
	cases := []struct {
		name           string
		args           string
		ExpectedApiKey string
		ExpectedURL    string
		expectError    bool
	}{
		{
			name:           "Init client with apiKey",
			args:           "test",
			ExpectedApiKey: "test",
			ExpectedURL:    "https://test@api.flowdock.com",
			expectError:    false,
		},
		{
			name:        "Init client without apiKey",
			args:        "",
			expectError: true,
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			got, gotErr := NewClient(cc.args)

			if cc.expectError {
				assert.Error(t, gotErr, "Init NewClient error")
				return
			}
			if got != nil {
				assert.Equal(t, cc.ExpectedApiKey, "test", "apiKey should be equal")
				assert.Equal(t, cc.ExpectedURL, "https://test@api.flowdock.com", "URL should be equal")
			}
		})
	}
}

func Test_Flowdock_Client_inviteNewUser_AccessDenied(t *testing.T) {
	client, _ := NewClient("apiKey")

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusForbidden)
		res.Write([]byte(inviteNewUserMockAccessDenied()))
	}))
	defer ts.Close()
	client.URL = ts.URL
	_, err := client.inviteNewUser("xxxxxxx@gmail.com", "message", "org", "flow")
	fmt.Printf(err.Error())
	assert.Error(t, err)
}

func inviteNewUserMockAccessDenied() string {
	return fmt.Sprintf(`{"message":"Access denied"}`)
}

func Test_Flowdock_Client_inviteNewUser_basic(t *testing.T) {
	client, _ := NewClient("apiKey")

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		fmt.Fprintln(res, inviteNewUserMockBasic())
	}))
	defer ts.Close()

	client.URL = ts.URL
	result, err := client.inviteNewUser("xxxxxxx@gmail.com", "message", "org", "flow")
	assert.NoError(t, err)
	assert.Equal(t, int64(1413413), result.ID)
	assert.Equal(t, "xxxxxxx@gmail.com", result.Email)
}

func inviteNewUserMockBasic() string {
	return fmt.Sprintf(`
	{
		"id": 1413413,
		"email": "xxxxxxx@gmail.com",
		"state": "pending",
		"url": "https://api.flowdock.com/flows/test-terraform/flow1/invitations/1413413"
	}
	`)
}

func Test_Flowdock_Client_deleteByUrl_basic(t *testing.T) {
	client, _ := NewClient("apiKey")
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	result := client.deleteByUrl(ts.URL)
	assert.NoError(t, result)
}

func Test_Flowdock_Client_deleteByUrl_NotFound(t *testing.T) {
	client, _ := NewClient("apiKey")
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(deleteUserFromOrgMockNotFound()))
	}))
	defer ts.Close()

	result := client.deleteByUrl(ts.URL)
	assert.Error(t, result)
}

func deleteUserFromOrgMockNotFound() string {
	return fmt.Sprintf(`{ "message": "not found"}`)
}

func Test_Flowdock_Client_getUserIdByEmail_UserExists(t *testing.T) {
	userList := []User{
		{
			ID:      123456,
			Email:   "xxxxx@gmail.com",
			Name:    "xxxxx",
			Nick:    "xxxxx",
			MESSAGE: "",
		},
		{
			ID:      654321,
			Email:   "yyyyy@gmail.com",
			Name:    "yyyyy",
			Nick:    "yyyyy",
			MESSAGE: "",
		},
	}

	output, err := json.Marshal(userList)
	if err != nil {
		t.Errorf("unexpected encoding error: %v", err)
		return
	}

	client, _ := NewClient("apiKey")
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Header().Set("Content-Type", "application/json")
		res.Write(output)
	}))
	defer ts.Close()
	client.URL = ts.URL

	result, _ := client.getUserIdByEmail("org", "xxxxx@gmail.com")
	assert.Equal(t, "123456", result)

	result1, _ := client.getUserIdByEmail("org", "yyyyy@gmail.com")
	assert.Equal(t, "654321", result1)

	noResult, _ := client.getUserIdByEmail("org", "zzzzz@gmail.com")
	assert.Equal(t, "", noResult)
}

func Test_Flowdock_Client_getUserIdByEmail_BadResponse(t *testing.T) {
	client, _ := NewClient("apiKey")
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusInternalServerError)
		res.Header().Set("Content-Type", "application/json")
	}))
	defer ts.Close()
	client.URL = ts.URL
	result, err := client.getUserIdByEmail("org", "zzzzz@gmail.com")

	assert.Error(t, err)
	assert.Equal(t, "", result)

}
