package flowdock

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	orgName  = "test-terraform"
	flowName = "flow1"
)

func TestAccFlowdock_Invite_One_User(t *testing.T) {
	resourceName := "flowdock_invitation.sirenfei_robot_1_test-terraform"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: checkItemBasic(clientVersion, os.Getenv("FLOWDOCK_TOKEN")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowdockItemExists(
						resourceName, orgName, flowName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "org", orgName),
					resource.TestCheckResourceAttr(resourceName, "flow", flowName),
					resource.TestCheckResourceAttr(resourceName, "email", "sirenfei.robot@gmail.com"),
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
	message = "sirenfei"
}
`, token, version)
}

func TestAccFlowdock_Invitation_Import_User_Resource(t *testing.T) {
	resourceName := "flowdock_invitation.richard-xue_1_test-terraform"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: checkItemImportBasic(clientVersion, os.Getenv("FLOWDOCK_TOKEN")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowdockItemExists(resourceName, orgName, flowName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "350495_flow1_test-terraform",
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
	resource "flowdock_invitation" "richard-xue_1_test-terraform" {
		org = "test-terraform"
		flow = "flow1"
		email = "richard.xue@fairfaxmedia.co.nz"
		message = "Richard Xue"
	}
`, token, version)
}

func TestAccFlowdock_Invitation_Update_Resource(t *testing.T) {
	resourceName := "flowdock_invitation.gyles-polloso_1_test-terraform"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: checkItemPreUpdate(clientVersion, os.Getenv("FLOWDOCK_TOKEN")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowdockItemExists(resourceName, orgName, flowName),
					resource.TestCheckResourceAttr(resourceName, "org", orgName),
					resource.TestCheckResourceAttr(resourceName, "flow", flowName),
					resource.TestCheckResourceAttr(resourceName, "email", "gyles.polloso@fairfaxmedia.co.nz"),
					resource.TestCheckResourceAttr(resourceName, "message", "Gyles Polloso"),
				),
			},
			{
				Config: checkItemPostUpdate(clientVersion, os.Getenv("FLOWDOCK_TOKEN")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowdockItemExists(resourceName, orgName, flowName),
					resource.TestCheckResourceAttr(resourceName, "org", orgName),
					resource.TestCheckResourceAttr(resourceName, "flow", flowName),
					resource.TestCheckResourceAttr(resourceName, "email", "gyles.polloso@fairfaxmedia.co.nz"),
					resource.TestCheckResourceAttr(resourceName, "message", "Post-Gyles"),
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
	resource "flowdock_invitation" "gyles-polloso_1_test-terraform" {
		org = "test-terraform"
		flow = "flow1"
		email = "gyles.polloso@fairfaxmedia.co.nz"
		message = "Gyles Polloso"
	}
`, token, version)
}
func checkItemPostUpdate(token string, version string) string {
	return fmt.Sprintf(`
	provider "flowdock" {
		version = "%s"
		api_token = "%s"
	}
	resource "flowdock_invitation" "gyles-polloso_1_test-terraform" {
		org = "test-terraform"
		flow = "flow1"
		email = "gyles.polloso@fairfaxmedia.co.nz"
		message = "Post-Gyles"
	}
`, token, version)
}

func TestAccFlowdock_Invite_Multiple_Resources(t *testing.T) {
	resourceName1 := "flowdock_invitation.gyles-polloso_1_test-terraform"
	resourceName2 := "flowdock_invitation.damian-mackle_1_test-terraform"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFlowdockMultiple(clientVersion, os.Getenv("FLOWDOCK_TOKEN")),
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
	resource "flowdock_invitation" "gyles-polloso_1_test-terraform" {
		org = "test-terraform"
		flow = "flow1"
		email = "gyles.polloso@fairfaxmedia.co.nz"
		message = "Gyles Polloso"
	}

	resource "flowdock_invitation" "damian-mackle_1_test-terraform" {
		org = "test-terraform"
		flow = "flow1"
		email = "damian.mackle@fairfaxmedia.co.nz"
		message = "Damian Mackle"
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
	_, err := client.getUserById(id)
	if err == nil {
		return true
	} else {
		_, error := client.getInvitationByInviteId(org, flow, id)
		if error == nil {
			return true
		}
	}
	return false
}
