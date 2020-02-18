package flowdock

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// invitations resource, as seen by GET /invitations
type Invitation struct {
	ID      int64  `json:"id"`
	Email   string `json:"email"`
	State   string `json:"state"`
	URL     string `json:"url"`
	MESSAGE string `json:"message"`
}

func ResourceInvitation() *schema.Resource {
	return &schema.Resource{
		Create: invitationCreate,
		Read:   invitationRead,
		Update: invitationUpdate,
		Delete: invitationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"org": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"flow": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"message": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func invitationCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client)
	org := d.Get("org").(string)
	email := d.Get("email").(string)
	flow := d.Get("flow").(string)
	message := d.Get("message").(string)

	userId, errorEmail := apiClient.getUserIdByEmail(org, email)
	if errorEmail != nil {
		return nil
	}
	if len(userId) > 0 {
		d.SetId("u" + userId)
		return nil
	}
	if len(message) > 0 {
		d.Set("message", message)
	}

	invitation, error := apiClient.inviteNewUser(email, message, org, flow)
	if error != nil {
		return fmt.Errorf("invitationCreate failed, response: %s", error)
	}

	d.SetId(strconv.FormatInt(invitation.ID, 10))
	return nil
}

func invitationRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client)

	// imported data format:userId_flowName_orgName
	if strings.Contains(d.Id(), "_") {
		userInfo := strings.Split(d.Id(), "_")
		userId, flow, org := userInfo[0], userInfo[1], userInfo[2]

		url := fmt.Sprintf("%s/users/%s", apiClient.URL, userId)
		res, err := apiClient.Http.Get(url)
		if err != nil {
			log.Printf("invitationRead error,get user error:%v", err)
			return err
		}
		defer res.Body.Close()
		user := &User{}
		json.NewDecoder(res.Body).Decode(user)

		d.SetId(userId)
		d.Set("org", org)
		d.Set("flow", flow)
		d.Set("email", user.Email)
		d.Set("message", user.Name)
	}
	d.SetId(d.Id())
	return nil
}

func invitationUpdate(d *schema.ResourceData, m interface{}) error {
	return invitationRead(d, m)
}

func invitationDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client)
	org := d.Get("org").(string)
	flow := d.Get("flow").(string)

	if strings.Contains(d.Id(), "u") {
		apiClient.deleteUserFromOrg(org, d.Id()[1:])
		return nil
	}
	apiClient.deleteInvitationById(org, flow, d.Id())
	return nil
}
