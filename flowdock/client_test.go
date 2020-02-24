package flowdock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Flowdock_Init_NewClient_When_System_First_Launched(t *testing.T) {
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

func Test_inviteNewUser_Should_Return_Error_When_Get_AccessDenied_From_Server(t *testing.T) {
	client, _ := NewClient("apiKey")

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusForbidden)
		res.Write([]byte(inviteNewUserMockAccessDenied()))
	}))
	defer ts.Close()
	client.URL = ts.URL
	_, err := client.inviteNewUser("xxxxxxx@fairfaxmedia.co.nz",
		"message", "org", "flow")
	assert.Error(t, err)
}

func inviteNewUserMockAccessDenied() string {
	return fmt.Sprintf(`{"message":"Access denied"}`)
}

func Test_Should_Ship_Invitation_With_Valid_Email_Org_Flow(t *testing.T) {
	client, _ := NewClient("apiKey")

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		fmt.Fprintln(res, inviteNewUserMockBasic())
	}))
	defer ts.Close()

	client.URL = ts.URL
	result, err := client.inviteNewUser("xxxxxxx@fairfaxmedia.co.nz",
		"message", "org", "flow")
	assert.NoError(t, err)
	assert.Equal(t, int64(1413413), result.ID)
	assert.Equal(t, "xxxxxxx@fairfaxmedia.co.nz", result.Email)
}

func inviteNewUserMockBasic() string {
	return fmt.Sprintf(`
	{
		"id": 1413413,
		"email": "xxxxxxx@fairfaxmedia.co.nz",
		"state": "pending",
		"url": "https://api.flowdock.com/flows/test-terraform/flow1/invitations/1413413"
	}
	`)
}

func Test_Should_Trigger_Http_Delete_When_Given_Org_And_UserId(t *testing.T) {
	client, _ := NewClient("apiKey")
	org := "org1"
	id := "123456"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusNoContent)
		assert.Equal(t, req.Method, "DELETE")
		assert.Equal(t, strings.Contains(req.URL.Path, org), true)
		assert.Equal(t, strings.Contains(req.URL.Path, id), true)
	}))
	defer ts.Close()
	client.URL = ts.URL + client.URL
	client.deleteUserFromOrg(org, id)
}

func Test_Should_Delete_User_Success_When_Given_Valid_URL(t *testing.T) {
	client, _ := NewClient("apiKey")

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusNoContent)
		fmt.Printf(req.URL.RawPath)
	}))
	defer ts.Close()
	result := client.deleteByUrl(ts.URL)

	assert.NoError(t, result)
}

func Test_Should_Delete_User_Failed_When_Get_Error_From_Server(t *testing.T) {
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

func Test_Should_getUserId_By_Email_And_Org_When_User_Exists_Or_Empty_If_Not(t *testing.T) {
	userList := []User{
		{
			ID:      123456,
			Email:   "xxxxx@fairfaxmedia.co.nz",
			Name:    "xxxxx",
			Nick:    "xxxxx",
			MESSAGE: "",
		},
		{
			ID:      654321,
			Email:   "yyyyy@fairfaxmedia.co.nz",
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

	result, _ := client.getUserIdByEmail("org", "xxxxx@fairfaxmedia.co.nz")
	assert.Equal(t, "123456", result)

	result1, _ := client.getUserIdByEmail("org", "yyyyy@fairfaxmedia.co.nz")
	assert.Equal(t, "654321", result1)

	noResult, _ := client.getUserIdByEmail("org", "zzzzz@fairfaxmedia.co.nz")
	assert.Equal(t, "", noResult)
}

func Test_getUserIdByEmail_Should_Get_Error_When_Internal_Server_Error_Happens(t *testing.T) {
	client, _ := NewClient("apiKey")
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusInternalServerError)
		res.Header().Set("Content-Type", "application/json")
	}))
	defer ts.Close()
	client.URL = ts.URL
	result, err := client.getUserIdByEmail("org", "zzzzz@fairfaxmedia.co.nz")

	assert.Error(t, err)
	assert.Equal(t, "", result)

}
