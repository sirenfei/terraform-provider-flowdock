package flowdock

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"flowdock": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("FLOWDOCK_TOKEN"); v == "" {
		t.Fatal("FLOWDOCK_TOKEN must be set for acceptance tests")

	}
}
func TestAccProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("error: %s", err)
	}
}

func TestAccProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}
