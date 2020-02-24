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
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"manager": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ticket_number": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"message": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func invitationCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*Client)
	org := d.Get("org").(string)
	email := d.Get("email").(string)
	flow := d.Get("flow").(string)
	username := d.Get("username").(string)
	manager := d.Get("manager").(string)
	ticketNumber := d.Get("ticket_number").(string)

	userId, errorE := apiClient.getUserIdByEmail(org, email)

	if errorE != nil && errorE.Error() != missMatchEmail {
		log.Printf("invitationCreate communications between client and server error")
		return nil
	}
	if len(userId) > 0 {
		d.SetId("u" + userId)
		return nil
	}

	//not existing user invoke invitation
	var message = fmt.Sprintf(`Hi %s
	Please complete this process for access to Flowdock.
		If you require assistance then contact %s in the first instance and ref %s
	Regards,
		Kiwiops.`, username, manager, ticketNumber)
	d.Set("message", message)

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
		d.Set("username", user.Name)
		d.Set("manager", user.Name)
		d.Set("ticket_number", "xxxx")
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
	email := d.Get("email").(string)
	// Id is invitation id and the user has accepted the invitatoin, delete by email
	userId, errorE := apiClient.getUserIdByEmail(org, email)
	// If the user isn't exist in the org,the id must be invitation Id, delete by id
	if errorE != nil && errorE.Error() == missMatchEmail {
		apiClient.deleteInvitationById(org, flow, d.Id())
		return nil
	} else if len(userId) > 0 {
		// User exists in the org but the id is invitation Id, need to delete by userId
		apiClient.deleteUserFromOrg(org, userId)
	}
	return nil
}
