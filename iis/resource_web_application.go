package iis

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rickedb/terraform-provider-iis/iis/agent"
)

func resourceWebApplication() *schema.Resource {
	return &schema.Resource{

		CreateContext: resourceWebApplicationCreate,
		ReadContext:   resourceWebApplicationRead,
		UpdateContext: resourceWebApplicationUpdate,
		DeleteContext: resourceWebApplicationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importWebApplicationState,
		},
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
				ForceNew:    true,
			},
			webAppSchema.Site: {
				Description: "Configures this web application to run in the specified application pool",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			webAppSchema.ApplicationPoolName: {
				Description: "Configures this web application to run in the specified application pool",
				Type:        schema.TypeString,
				Required:    true,
			},
			webAppSchema.PhysicalPath: {
				Description:      "Physical path to the content for the virtual directory",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: isValidPath(true),
			},
		},
	}
}

func resourceWebApplicationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceWebApplicationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*agent.Client)

	if d.HasChange(webAppSchema.ApplicationPoolName) {
		appPoolName := d.Get(webAppSchema.ApplicationPoolName).(string)
		if err := validateAppPoolExists(client, appPoolName); err != nil {
			return diag.FromErr(err)
		}
	}

	webApplicationRequest := mapToWebApplication(d)
	webApplication, err := client.CreateWebApplication(webApplicationRequest)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s_%s", webApplication.Site, webApplication.Name))
	return nil
}

func resourceWebApplicationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*agent.Client)

	if d.HasChange(webAppSchema.ApplicationPoolName) {
		appPoolName := d.Get(webAppSchema.ApplicationPoolName).(string)
		if err := validateAppPoolExists(client, appPoolName); err != nil {
			return diag.FromErr(err)
		}
	}

	webApplicationRequest := mapToWebApplication(d)
	err := client.UpdateWebApplication(webApplicationRequest)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	return nil
}

func resourceWebApplicationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*agent.Client)

	site := d.Get(webAppSchema.Site).(string)
	name := d.Get(webAppSchema.Name).(string)
	err := client.DeleteWebApplication(site, name)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	return nil
}

func importWebApplicationState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*agent.Client)
	id := d.Id()
	siteAndName := strings.Split(id, "_")
	if len(siteAndName) < 2 {
		return nil, errors.New("provided id is invalid, please provide the id in the following format: '{site_name}_{web_application_name}")
	}

	webApplication, err := client.GetWebApplication(siteAndName[0], siteAndName[1])
	if err != nil {
		d.SetId("")
		return nil, err
	}

	if err = mapWebApplicationToResourceData(*webApplication, d); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func mapWebApplicationToResourceData(webApplication agent.WebApplication, d *schema.ResourceData) error {
	var err error
	d.SetId(webApplication.Id)
	if err = d.Set(webAppSchema.Name, webApplication.Name); err != nil {
		return err
	}
	if err = d.Set(webAppSchema.Path, webApplication.Path); err != nil {
		return err
	}
	if err = d.Set(webAppSchema.PhysicalPath, webApplication.PhysicalPath); err != nil {
		return err
	}
	if err = d.Set(webAppSchema.Site, webApplication.Site); err != nil {
		return err
	}
	if err = d.Set(webAppSchema.ApplicationPoolName, webApplication.ApplicationPoolName); err != nil {
		return err
	}

	return nil
}

func mapToWebApplication(d *schema.ResourceData) agent.WebApplication {
	return agent.WebApplication{
		Name:                d.Get(webAppSchema.Name).(string),
		Site:                d.Get(webAppSchema.Site).(string),
		ApplicationPoolName: d.Get(webAppSchema.ApplicationPoolName).(string),
		PhysicalPath:        d.Get(webAppSchema.PhysicalPath).(string),
	}
}
