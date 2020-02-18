package flowdock

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Organization resource, as seen by GET /organization
type Organization struct {
	ID      int64  `json:"id"`
	APIName string `json:"parameterized_name"`
	Name    string `json:"name"`
	APIURL  string `json:"url"`
	Users   []User `json:"users"` // Maps user ID's to user objects.
	MESSAGE string `json:"message"`
}

func ResourceOrganization() *schema.Resource {
	return &schema.Resource{
		Create: resourceOrganiztionCreate,
		Read:   resourceOrganiztionRead,
		Update: resourceOrganiztionUpdate,
		Delete: resourceOrganiztionDelete,

		Schema: map[string]*schema.Schema{
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"message": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceOrganiztionCreate(d *schema.ResourceData, m interface{}) error {
	email := d.Get("email").(string)
	d.SetId(email)
	return resourceOrganiztionRead(d, m)
}

func resourceOrganiztionRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceOrganiztionUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceOrganiztionRead(d, m)
}

func resourceOrganiztionDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
