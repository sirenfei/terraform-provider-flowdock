package flowdock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"org": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceUserRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client)
	org := d.Get("org").(string)
	email := d.Get("email").(string)
	var url string

	if org != "" {
		url = fmt.Sprintf("https://%s@api.flowdock.com/organizations/%s/users", apiClient.ApiKey, org)
	} else {
		url = fmt.Sprintf("https://%s@api.flowdock.com/users", apiClient.ApiKey)
	}
	log.Printf("[DEBUG] Reading url: %s", url)
	res, err := apiClient.Http.Get(url)
	if err != nil {
		log.Printf("dataSourcesUserRead error:%s", err.Error())
	}
	defer res.Body.Close()

	var users []User
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	json.Unmarshal(body, &users)

	for _, user := range users {
		if user.Email == email {
			d.SetId(strconv.FormatInt(user.ID, 10))
			_ = d.Set("name", user.Name)
			_ = d.Set("email", user.Email)
			_ = d.Set("org", org)
			return nil
		}
	}
	return nil
}
