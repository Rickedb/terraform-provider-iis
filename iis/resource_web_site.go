package iis

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rickedb/terraform-provider-iis/iis/agent"
)

func resourceWebsite() *schema.Resource {
	return &schema.Resource{
		Description:   "A container for applications and virtual directories which you can access it through one or more unique bindings",
		CreateContext: resourceWebsiteCreate,
		ReadContext:   resourceWebsiteRead,
		UpdateContext: resourceWebsiteUpdate,
		DeleteContext: resourceWebsiteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importWebSiteState,
		},
		Schema: map[string]*schema.Schema{
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
			webSiteSchema.Name: {
				Description: "An unique name for the associated site",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			webSiteSchema.ApplicationPoolName: {
				Description: "Configures this web site to run in the specified application pool",
				Type:        schema.TypeString,
				Required:    true,
			},
			webSiteSchema.PhysicalPath: {
				Description:      "Physical path to the content for the virtual directory",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: isValidPath(true),
			},
			webSiteSchema.Username: {
				Description: "Username for the user identity that should be impersonated when accessing the physical path for the virtual directory",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			webSiteSchema.Password: {
				Description: "Password for the user identity that should be impersonated when accessing the physical path for the virtual directory",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Default:     "",
			},
			webSiteSchema.BindingSchema.Key: {
				Description: "An HTTP binding is a combination of IP address, port and host name (the host name can be a domain name). HTTP.sys listens on the IP/port for incoming requests",
				Type:        schema.TypeSet,
				Optional:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: webSiteBindingsSchema,
				},
			},
		},
	}
}

var webSiteBindingsSchema = map[string]*schema.Schema{
	webSiteSchema.BindingSchema.Protocol: {
		Description:      "Use HTTP if you want the website to have an HTTP binding, or select HTTPS if you want the website to have a Secure Sockets Layer (SSL) binding",
		Type:             schema.TypeString,
		Optional:         true,
		Default:          "http",
		ValidateDiagFunc: validateAllowedValues([]string{"http", "https"}),
	},
	webSiteSchema.BindingSchema.Ip: {
		Description: "An IP address that users can use to access this site",
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "*",
	},
	webSiteSchema.BindingSchema.Port: {
		Description:      "The port on which HTTP.sys must listen for requests made to this website",
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          80,
		ValidateDiagFunc: isInBetweenValues(1, 65535),
	},
	webSiteSchema.BindingSchema.HostHeader: {
		Description: "A host name if you want to assign one or more host names, also known as domain names, to one computer that uses a single IP address. If you specify a host name, clients must use the host name instead of the IP address to access the website",
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
	},
}

func resourceWebsiteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*agent.Client)

	if d.HasChange(webSiteSchema.ApplicationPoolName) {
		appPoolName := d.Get(webSiteSchema.ApplicationPoolName).(string)
		if err := validateAppPoolExists(client, appPoolName); err != nil {
			return diag.FromErr(err)
		}
	}

	webSiteRequest := mapToWebSite(d)
	webSite, err := client.CreateWebSite(webSiteRequest)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(webSite.Id)
	d.Set(webSiteSchema.State, webSite.State)
	return nil
}

func resourceWebsiteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*agent.Client)
	name := d.Get(webSiteSchema.Name).(string)
	webSite, err := client.GetWebSite(name)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	if err = mapWebSiteToResourceData(*webSite, d); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	return nil
}

func resourceWebsiteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*agent.Client)

	if d.HasChange(webSiteSchema.ApplicationPoolName) {
		appPoolName := d.Get(webSiteSchema.ApplicationPoolName).(string)
		if err := validateAppPoolExists(client, appPoolName); err != nil {
			return diag.FromErr(err)
		}
	}

	webSite := mapToWebSite(d)
	err := client.UpdateWebSite(webSite)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	return nil
}

func resourceWebsiteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*agent.Client)

	name := d.Get(webSiteSchema.Name).(string)
	err := client.DeleteWebSite(name)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	return nil
}

func importWebSiteState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*agent.Client)
	webSiteName := d.Id()
	webSite, err := client.GetWebSite(webSiteName)
	if err != nil {
		d.SetId("")
		return nil, err
	}

	if err = mapWebSiteToResourceData(*webSite, d); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func mapWebSiteToResourceData(webSite agent.WebSite, d *schema.ResourceData) error {
	var err error
	d.SetId(webSite.Id)
	if err = d.Set(webSiteSchema.Name, webSite.Name); err != nil {
		return err
	}
	if err = d.Set(webSiteSchema.ApplicationPoolName, webSite.ApplicationPoolName); err != nil {
		return err
	}
	if err = d.Set(webSiteSchema.State, webSite.State); err != nil {
		return err
	}
	if err = d.Set(webSiteSchema.PhysicalPath, webSite.PhysicalPath); err != nil {
		return err
	}
	if err = d.Set(webSiteSchema.Username, webSite.Username); err != nil {
		return err
	}
	if err = d.Set(webSiteSchema.Password, webSite.Password); err != nil {
		return err
	}

	var bindings []map[string]interface{}
	for _, binding := range webSite.Bindings {
		b := map[string]interface{}{
			webSiteSchema.BindingSchema.Ip:         binding.Ip,
			webSiteSchema.BindingSchema.Port:       binding.Port,
			webSiteSchema.BindingSchema.Protocol:   binding.Protocol,
			webSiteSchema.BindingSchema.HostHeader: binding.HostHeader,
		}
		bindings = append(bindings, b)
	}

	err = d.Set(webSiteSchema.BindingSchema.Key, bindings)
	return err
}

func mapToWebSite(d *schema.ResourceData) agent.WebSite {
	bindings := []agent.Binding{}

	bindingsList := d.Get(webSiteSchema.BindingSchema.Key).(*schema.Set).List()
	if len(bindingsList) > 0 {
		for _, binding := range bindingsList {
			bindingResource := binding.(map[string]interface{})
			bindings = append(bindings, agent.Binding{
				Ip:         bindingResource[webSiteSchema.BindingSchema.Ip].(string),
				Port:       bindingResource[webSiteSchema.BindingSchema.Port].(int),
				Protocol:   bindingResource[webSiteSchema.BindingSchema.Protocol].(string),
				HostHeader: bindingResource[webSiteSchema.BindingSchema.HostHeader].(string),
			})
		}
	} else {
		bindings = append(bindings, agent.Binding{
			Ip:         webSiteBindingsSchema[webSiteSchema.BindingSchema.Ip].Default.(string),
			Port:       webSiteBindingsSchema[webSiteSchema.BindingSchema.Port].Default.(int),
			Protocol:   webSiteBindingsSchema[webSiteSchema.BindingSchema.Protocol].Default.(string),
			HostHeader: webSiteBindingsSchema[webSiteSchema.BindingSchema.HostHeader].Default.(string),
		})
	}

	return agent.WebSite{
		Name:                d.Get(webSiteSchema.Name).(string),
		PhysicalPath:        d.Get(webSiteSchema.PhysicalPath).(string),
		Username:            d.Get(webSiteSchema.Username).(string),
		Password:            d.Get(webSiteSchema.Password).(string),
		ApplicationPoolName: d.Get(webSiteSchema.ApplicationPoolName).(string),
		Bindings:            bindings,
	}
}
