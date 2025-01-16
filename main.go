package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	iis "github.com/rickedb/terraform-provider-iis/iis"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: iis.Provider,
	})
}
