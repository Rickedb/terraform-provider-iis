package iis

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rickedb/terraform-provider-iis/iis/agent"
)

func dataSourceWebSite() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWebSiteRead,
		Schema: map[string]*schema.Schema{
			webSiteSchema.Name: {
				Description: "An unique name for the associated site",
				Type:        schema.TypeString,
				Required:    true,
			},
			webSiteSchema.Id: {
				Description: "An unique numeric identifier for the site. This identifier is used in directory names for log files and trace files",
				Type:        schema.TypeString,
				Computed:    true,
			},
			webSiteSchema.State: {
				Description: "The current state of the site",
				Type:        schema.TypeString,
				Computed:    true,
			},
			webSiteSchema.ApplicationPoolName: {
				Description: "Configures this web site to run in the specified application pool",
				Type:        schema.TypeString,
				Computed:    true,
			},
			webSiteSchema.PhysicalPath: {
				Description: "Physical path to the content for the virtual directory",
				Type:        schema.TypeString,
				Computed:    true,
			},
			webSiteSchema.Username: {
				Description: "Username for the user identity that should be impersonated when accessing the physical path for the virtual directory",
				Type:        schema.TypeString,
				Computed:    true,
			},
			webSiteSchema.Password: {
				Description: "Password for the user identity that should be impersonated when accessing the physical path for the virtual directory",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceWebSiteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*agent.Client)
	name := d.Get(webAppSchema.Name).(string)
	webSite, err := client.GetWebSite(name)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = mapWebSiteToResourceData(*webSite, d); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	return nil
}
