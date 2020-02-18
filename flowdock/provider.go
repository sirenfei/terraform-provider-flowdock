package flowdock

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{

		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("FLOWDOCK_TOKEN", nil),
				Description: "please add your api token from https://www.flowdock.com/account/tokens",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"flowdock_user": DataSourceUser(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"flowdock_invitation":   ResourceInvitation(),
			"flowdock_organization": ResourceOrganization(),
			"flowdock_user":         ResourceUser(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(provider *schema.ResourceData) (interface{}, error) {
	return NewClient(provider.Get("api_token").(string))
}
