package flowdock

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccFlowdock_Invitation_basic(t *testing.T) {
	resourceName := "flowdock_invitation.mickey_mouse_1_smart-mouse"
	orgName := "smart-mouse"
	flowName := "ops-projects"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: checkItemBasic("1.1.2", os.Getenv("FLOWDOCK_TOKEN")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowdockItemExists(
						resourceName, orgName, flowName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "org", orgName),
					resource.TestCheckResourceAttr(resourceName, "flow", flowName),
					resource.TestCheckResourceAttr(resourceName, "email", "mickey.mouse@gmail.com"),
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
	resource "flowdock_invitation" "mickey_mouse_1_smart-mouse" {
		org = "smart-mouse"
		flow = "ops-projects"
		email = "mickey.mouse@gmail.com"
		message = "gyles"
	}
`, token, version)
}

func TestAccFlowdock_Invitation_Importbasic(t *testing.T) {
	resourceName := "flowdock_invitation.mickey_mouse_1_smart-mouse"
	orgName := "smart-mouse"
	flowName := "ops-projects"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: checkItemImportBasic("1.1.2", os.Getenv("FLOWDOCK_TOKEN")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowdockItemExists(resourceName, orgName, flowName),
				),
			},
			{
				ResourceName:      "flowdock_invitation.mickey_moouse_1_smart-mouse",
				ImportState:       true,
				ImportStateId:     "123456_ops-projects_smart-mouse",
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
	resource "flowdock_invitation" "mickey_mouse_1_smart-mouse" {
		org = "smart-mouse"
		flow = "ops-projects"
		email = "mickey.mouse@gmail.com"
		message = "Mickey Mouse"
	}
`, token, version)
}

func TestAccFlowdock_Invitation_Update(t *testing.T) {
	resourceName := "flowdock_invitation.mickey_mouse_1_smart-mouse"
	orgName := "smart-mouse"
	flowName := "ops-projects"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: checkItemPreUpdate("1.1.2", os.Getenv("FLOWDOCK_TOKEN")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowdockItemExists(resourceName, orgName, flowName),
					resource.TestCheckResourceAttr(resourceName, "org", orgName),
					resource.TestCheckResourceAttr(resourceName, "flow", flowName),
					resource.TestCheckResourceAttr(resourceName, "email", "mickey.mouse@gmail.com"),
					resource.TestCheckResourceAttr(resourceName, "message", "Mickey Mouse"),
				),
			},
			{
				Config: checkItemPostUpdate("1.1.2", os.Getenv("FLOWDOCK_TOKEN")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowdockItemExists(resourceName, orgName, "testing2"),
					resource.TestCheckResourceAttr(resourceName, "flow", "testing2"),
					resource.TestCheckResourceAttr(resourceName, "email", "mickey.mouse@gmail.com"),
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
	resource "flowdock_invitation" "mickey_mouse_1_smart-mouse" {
		org = "smart-mouse"
		flow = "ops-projects"
		email = "mickey.mouse@gmail.com"
		message = "Mickey Mouse"
	}
`, token, version)
}
func checkItemPostUpdate(token string, version string) string {
	return fmt.Sprintf(`
	provider "flowdock" {
		version = "%s"
		api_token = "%s"
	}
	resource "flowdock_invitation" "mickey_mouse_1_smart-mouse" {
		org = "smart-mouse"
		flow = "testing2"
		email = "mickey.mouse@gmail.com"
		message = "Post-Gyles"
	}
`, token, version)
}

func TestAccFlowdock_Multiple(t *testing.T) {
	resourceName1 := "flowdock_invitation.mickey_mouse_1_smart-mouse"
	resourceName2 := "flowdock_invitation.mickey_mouse2_1_smart-mouse"
	orgName := "smart-mouse"
	flowName := "ops-projects"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFlowdockMultiple("1.1.2", os.Getenv("FLOWDOCK_TOKEN")),
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
	resource "flowdock_invitation" "mickey_mouse_1_smart-mouse" {
		org = "smart-mouse"
		flow = "ops-projects"
		email = "mickey.mouse@gmail.com"
		message = "Mickey Mouse"
	}

	resource "flowdock_invitation" "mickey_mouse2_1_smart-mouse" {
		org = "smart-mouse"
		flow = "ops-projects"
		email = "mickey.mouse@gmail.com"
		message = "Mickey Mouse"
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
