package iis

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rickedb/terraform-provider-iis/iis/agent"
)

func dataSourceWebApplication() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWebApplicationRead,
		Schema: map[string]*schema.Schema{
			webAppSchema.Id: {
				Description: "An unique numeric identifier for the web application",
				Type:        schema.TypeString,
				Computed:    true,
			},
			webAppSchema.Path: {
				Description: "URL path for the application",
				Type:        schema.TypeString,
				Computed:    true,
			},
			webAppSchema.Name: {
				Description: "An unique name for the web application",
				Type:        schema.TypeString,
				Required:    true,
			},
			webAppSchema.Site: {
				Description: "Configures this web application to run in the specified application pool",
				Type:        schema.TypeString,
				Computed:    true,
			},
			webAppSchema.ApplicationPoolName: {
				Description: "Configures this web application to run in the specified application pool",
				Type:        schema.TypeString,
				Computed:    true,
			},
			webAppSchema.PhysicalPath: {
				Description: "Physical path to the content for the virtual directory",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceWebApplicationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*agent.Client)

	site := d.Get(webAppSchema.Site).(string)
	name := d.Get(webAppSchema.Name).(string)
	webApplication, err := client.GetWebApplication(site, name)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	mapWebApplicationToResourceData(*webApplication, d)
	return nil
}
