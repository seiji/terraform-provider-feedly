package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/seiji/terraform-provider-feedly/feedly"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: feedly.Provider})
}
