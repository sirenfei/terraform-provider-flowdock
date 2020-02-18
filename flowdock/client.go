package flowdock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// A Client is a Flowdock API client. It should be created
// using NewClient() and provided with a valid API key.
type Client struct {
	ApiKey string
	// HTTP client used to communicate with the API.
	Http *http.Client
	URL  string
}

// NewClient creates a new Client and automatically fetches
func NewClient(apiKey string) (*Client, error) {
	if len(strings.TrimSpace(apiKey)) == 0 {
		return nil, fmt.Errorf("can't run with an empty token")
	}
	client := &Client{
		ApiKey: apiKey,
		Http:   &http.Client{Timeout: 10 * time.Second},
		URL:    fmt.Sprintf("https://%s@api.flowdock.com", apiKey),
	}
	return client, nil
}

func (client *Client) getUserById(userId string) (*User, error) {
	url := fmt.Sprintf("%s/users/%s", client.URL, userId)
	res, err := client.Http.Get(url)
	if err != nil {
		log.Printf("getUserById error:%s", err.Error())
		return nil, fmt.Errorf("getUserById http request error: %s", userId)
	}
	defer res.Body.Close()
	user := &User{}
	json.NewDecoder(res.Body).Decode(user)
	if user.ID == 0 {
		return nil, fmt.Errorf("no matching user with userId %s", userId)
	}
	return user, nil
}

func (client *Client) getInvitationByInviteId(org string, flow string, inviteId string) (*Invitation, error) {
	url := fmt.Sprintf("%s/flows/%s/%s/invitations/%s", client.URL, org, flow, inviteId)

	res, err := client.Http.Get(url)
	if err != nil {
		log.Printf("getInvitationByInviteId error:%s", err.Error())
		return nil, fmt.Errorf("getInvitationByInviteId error")
	}
	defer res.Body.Close()
	invitation := &Invitation{}
	json.NewDecoder(res.Body).Decode(invitation)
	if invitation.ID == 0 {
		return nil, fmt.Errorf("no matching invitation with inviteId %s", inviteId)
	}
	return invitation, nil
}

func (client *Client) inviteNewUser(email string, message string,
	org string, flow string) (*Invitation, error) {

	params := url.Values{
		"email":   {email},
		"message": {message},
	}
	url := fmt.Sprintf("%s/flows/%s/%s/invitations", client.URL, org, flow)

	res, error := client.Http.PostForm(url, params)
	if error != nil {
		log.Printf("inviteNewUser http request error: %s", error)
	}
	defer res.Body.Close()

	invitation := &Invitation{}
	json.NewDecoder(res.Body).Decode(invitation)

	if invitation.ID == 0 {
		return nil, fmt.Errorf("inviteNewUser error, invitation id=0, response: %s", invitation.MESSAGE)
	}

	return invitation, nil
}

func (client *Client) deleteUserFromOrg(org string, id string) error {
	url := fmt.Sprintf("%s/organizations/%s/users/%s", client.URL, org, id)
	result := client.deleteByUrl(url)
	return result
}
func (client *Client) deleteInvitationById(org string, flow string, id string) error {
	url := fmt.Sprintf("%s/flows/%s/%s/invitations/%s", client.URL, org, flow, id)
	result := client.deleteByUrl(url)
	return result
}

func (client *Client) deleteByUrl(url string) error {
	// Create request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Printf("user Delete failed")
		return fmt.Errorf("user Delete failed")
	}
	// Fetch Request
	res, err := client.Http.Do(req)
	if err != nil || res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete request failed or user not found!")
	}
	defer res.Body.Close()

	return nil
}

func (client *Client) getUserIdByEmail(org string, email string) (string, error) {
	var url = fmt.Sprintf("%s/organizations/%s/users", client.URL, org)

	res, errHttp := client.Http.Get(url)
	if errHttp != nil {
		log.Printf("getUserIdByEmail http request error:%s", errHttp.Error())
		return "", errHttp
	}
	defer res.Body.Close()

	var users []User
	body, errIO := ioutil.ReadAll(res.Body)
	if errIO != nil {
		log.Printf("getUserIdByEmail ioutil.ReadAll encoding error response body: %s", res.Body)
		return "", errIO
	}

	errorEncode := json.Unmarshal(body, &users)
	if errorEncode != nil {
		log.Printf("unexpected encoding error: %s", body)
		return "", fmt.Errorf("unexpected encoding error:%s", body)
	}

	for _, user := range users {
		log.Printf("unmarshalled users, userId:%d, email:%s", user.ID, user.Email)
		if user.Email == email {
			return strconv.FormatInt(user.ID, 10), nil
		}
	}
	log.Printf("getUserIdByEmail didn't find matching email:%s in org:%s", email, org)
	return "", fmt.Errorf("no user found by the email")
}
