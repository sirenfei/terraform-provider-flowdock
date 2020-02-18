package flowdock

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// users resource, as seen by GET /users/:id
type User struct {
	ID      int64  `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Nick    string `json:"nick"`
	MESSAGE string `json:"message"`
}

func ResourceUser() *schema.Resource {
	return &schema.Resource{
		Create: userCreate,
		Read:   userRead,
		Update: userUpdate,
		Delete: userDelete,

		Schema: map[string]*schema.Schema{
			"org": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"flow": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"user_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func userCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client)
	org := d.Get("org").(string)
	flow := d.Get("flow").(string)
	userId := d.Get("user_id").(string)

	params := url.Values{
		"id": {userId},
	}

	url := fmt.Sprintf("https://%s@api.flowdock.com/flows/%s/%s/users", apiClient.ApiKey, org, flow)

	res, error := apiClient.Http.PostForm(url, params)
	if error != nil {
		log.Printf("error:%v", error)
	}
	defer res.Body.Close()
	d.SetId(userId)
	return userRead(d, meta)
}

func userRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client)

	url := fmt.Sprintf("https://%s@api.flowdock.com/users/%s", apiClient.ApiKey, d.Id())

	res, err := apiClient.Http.Get(url)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	user := &User{}
	json.NewDecoder(res.Body).Decode(user)

	d.SetId(strconv.FormatInt(user.ID, 10))
	d.Set("email", user.Email)
	d.Set("name", user.Name)
	d.Set("nick", user.Nick)
	return nil
}

func userUpdate(d *schema.ResourceData, meta interface{}) error {
	return userRead(d, meta)
}

func userDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client)
	org := d.Get("org").(string)

	url := fmt.Sprintf("https://%s@api.flowdock.com/organizations/%s/users/%s", apiClient.ApiKey, org, d.Id())

	// Create request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Printf("user Delete failed")
		return err
	}
	// Fetch Request
	res, err := apiClient.Http.Do(req)
	if err != nil {
		log.Printf("user Delete failed")
		return err
	}
	defer res.Body.Close()

	return nil
}
