package flowdock

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccFlowdock_Invite_One_User(t *testing.T) {
	resourceName := "flowdock_invitation.gyles_polloso_1_stuff-kiwiops-projects"
	orgName := "test-terraform"
	flowName := "flow1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: checkItemBasic("1.1.6", os.Getenv("FLOWDOCK_TOKEN")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowdockItemExists(
						resourceName, orgName, flowName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "org", orgName),
					resource.TestCheckResourceAttr(resourceName, "flow", flowName),
					resource.TestCheckResourceAttr(resourceName, "email", "gyles.polloso@fairfaxmedia.co.nz"),
					resource.TestCheckResourceAttrSet(resourceName, "message"),
				),
			},
		},
	})
}
func checkItemBasic(token string, version string) string {
	return fmt.Sprintf(`
	provider "flowdock" {
		version = "%s"
		api_token = "%s"
	}
	resource "flowdock_invitation" "sirenfei_robot_1_test-terraform" {
		org = "test-terraform"
		flow = "flow1"
		email = "sirenfei.robot@gmail.com"
		username = "sirenfei"
		manager = "sirenfei"
		ticket_number = "xxxx"
	}
`, token, version)
}

func TestAccFlowdock_Invitation_Import_User_Resource(t *testing.T) {
	resourceName := "flowdock_invitation.reinxue_1_test-terraform"
	orgName := "test-terraform"
	flowName := "flow1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: checkItemImportBasic("1.1.6", os.Getenv("FLOWDOCK_TOKEN")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowdockItemExists(resourceName, orgName, flowName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "419225_flow1_test-terraform",
				ImportStateVerify: true,
			},
		},
	})
}

func checkItemImportBasic(token string, version string) string {
	return fmt.Sprintf(`
	provider "flowdock" {
		version = "%s"
		api_token = "%s"
	}
	resource "flowdock_invitation" "richard_xue_1_test-terraform" {
		org = "test-terraform"
		flow = "flow1"
		email = "richard.xue@fairfaxmedia.co.nz"
		username = "richard xue"
		manager = "richard xue"
		ticket_number = "xxxx"
	}
`, token, version)
}

func TestAccFlowdock_Invitation_Update_Resource(t *testing.T) {
	resourceName := "flowdock_invitation.gyles_polloso_1_stuff-kiwiops-projects"
	orgName := "stuff-kiwiops-projects"
	flowName := "kiwiops-projects"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: checkItemPreUpdate("1.1.6", os.Getenv("FLOWDOCK_TOKEN")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowdockItemExists(resourceName, orgName, flowName),
					resource.TestCheckResourceAttr(resourceName, "org", orgName),
					resource.TestCheckResourceAttr(resourceName, "flow", flowName),
					resource.TestCheckResourceAttr(resourceName, "email", "gyles.polloso@fairfaxmedia.co.nz"),
					resource.TestCheckResourceAttr(resourceName, "username", "Gyles Polloso"),
					resource.TestCheckResourceAttr(resourceName, "manager", "Darren"),
				),
			},
			{
				Config: checkItemPostUpdate("1.1.6", os.Getenv("FLOWDOCK_TOKEN")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowdockItemExists(resourceName, orgName, "testing2"),
					resource.TestCheckResourceAttr(resourceName, "flow", "testing2"),
					resource.TestCheckResourceAttr(resourceName, "email", "gyles.polloso@fairfaxmedia.co.nz"),
					resource.TestCheckResourceAttr(resourceName, "username", "Post-Gyles"),
					resource.TestCheckResourceAttr(resourceName, "manager", "Post-Darren"),
				),
			},
		},
	})
}

func checkItemPreUpdate(token string, version string) string {
	return fmt.Sprintf(`
	provider "flowdock" {
		version = "%s"
		api_token = "%s"
	}
	resource "flowdock_invitation" "gyles_polloso_1_stuff-kiwiops-projects" {
		org = "stuff-kiwiops-projects"
		flow = "kiwiops-projects"
		email = "gyles.polloso@fairfaxmedia.co.nz"
		username = "Gyles Polloso"
		manager = "Darren"
		ticket_number = "xxxx"
	}
`, token, version)
}
func checkItemPostUpdate(token string, version string) string {
	return fmt.Sprintf(`
	provider "flowdock" {
		version = "%s"
		api_token = "%s"
	}
	resource "flowdock_invitation" "gyles_polloso_1_stuff-kiwiops-projects" {
		org = "stuff-kiwiops-projects"
		flow = "testing2"
		email = "gyles.polloso@fairfaxmedia.co.nz"
		username = "Post-Gyles"
		manager = "Post-Darren"
		ticket_number = "xxxx"
	}
`, token, version)
}

func TestAccFlowdock_Invite_Multiple_Resources(t *testing.T) {
	resourceName1 := "flowdock_invitation.gyles_polloso_1_stuff-kiwiops-projects"
	resourceName2 := "flowdock_invitation.damian_mackle_1_stuff-kiwiops-projects"
	orgName := "stuff-kiwiops-projects"
	flowName := "kiwiops-projects"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFlowdockMultiple("1.1.6", os.Getenv("FLOWDOCK_TOKEN")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowdockItemExists(resourceName1, orgName, flowName),
					testAccCheckFlowdockItemExists(resourceName2, orgName, flowName),
				),
			},
		},
	})
}

func testAccCheckFlowdockMultiple(token string, version string) string {
	return fmt.Sprintf(`
	provider "flowdock" {
		version = "%s"
		api_token = "%s"
	}
	resource "flowdock_invitation" "gyles_polloso_1_stuff-kiwiops-projects" {
		org = "stuff-kiwiops-projects"
		flow = "kiwiops-projects"
		email = "gyles.polloso@fairfaxmedia.co.nz"
		username = "gyles"
		manager = "gyles"
		ticket_number = "xxxx"
	}

	resource "flowdock_invitation" "damian_mackle_1_stuff-kiwiops-projects" {
		org = "stuff-kiwiops-projects"
		flow = "kiwiops-projects"
		email = "damian.mackle@fairfaxmedia.co.nz"
		username = "Damian Mackle"
		manager = "Damian Mackle"
		ticket_number = "xxxx"
	}
	`, token, version)
}

func testAccCheckFlowdockItemExists(resource string, org string, flow string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Resource ID is set")
		}
		id := rs.Primary.ID
		isExist := isInvitationOrUserExistsInFlowdockServer(org, flow, id)
		if isExist == false {
			return fmt.Errorf("error testAccCheckFlowdockItemExists with resource %s", resource)
		}
		return nil
	}
}

// Check all domains specified in the configuration have been destroyed.
func checkItemDestroy(state *terraform.State) error {
	for _, item := range state.RootModule().Resources {

		if item.Type != "flowdock_invitation" {
			continue
		}

		id := item.Primary.ID
		// var org, flow string
		org := item.Primary.Attributes["org"]
		flow := item.Primary.Attributes["flow"]

		isExist := isInvitationOrUserExistsInFlowdockServer(org, flow, id)

		if isExist == true {
			return fmt.Errorf("Flowdock user '%s' still exists.", id)
		}
	}

	return nil
}

func isInvitationOrUserExistsInFlowdockServer(org string, flow string, id string) bool {
	client := testAccProvider.Meta().(*Client)
	// if id is userid, it means its prefix is u
	if strings.Contains(id, "u") {
		_, err := client.getUserById(id[1:])
		if err == nil {
			return true
		}
	} else {
		_, error := client.getInvitationByInviteId(org, flow, id)
		if error == nil {
			return true
		}
	}
	return false
}
