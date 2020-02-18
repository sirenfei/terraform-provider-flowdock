package main

import (
	"terraform-provider-flowdock/flowdock"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: flowdock.Provider,
	})
}
