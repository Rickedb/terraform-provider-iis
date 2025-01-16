package iis

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rickedb/terraform-provider-iis/iis/agent"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"hostname": {
				Description: "The remote server which will be hosting the resources",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"username": {
				Description: "The username to be used at credentials when accessing the remote server (it must have administrator permissions)",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"password": {
				Description: "The password to be used at credentials when accessing the remote server",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"iis_application_pool": resourceApplicationPool(),
			"iis_web_site":         resourceWebsite(),
			"iis_web_application":  resourceWebApplication(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"iis_application_pool": dataSourceApplicationPool(),
			"iis_web_site":         dataSourceWebSite(),
			"iis_web_application":  dataSourceWebApplication(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	client := &agent.Client{
		Hostname: d.Get("hostname").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}

	return client, nil
}
