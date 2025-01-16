package iis

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rickedb/terraform-provider-iis/iis/agent"
)

func dataSourceApplicationPool() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationPoolRead,
		Schema: map[string]*schema.Schema{
			applicationPoolSchema.Name: {
				Type:     schema.TypeString,
				Required: true,
			},
			applicationPoolSchema.State: {
				Description: "The current state of the application pool",
				Type:        schema.TypeString,
				Computed:    true,
			},
			applicationPoolSchema.StartMode: {
				Description: "Configures application pool to run in On Demand Mode or Always Running Mode",
				Type:        schema.TypeString,
				Computed:    true,
			},
			applicationPoolSchema.PipelineMode: {
				Description: "Configures ASP.NET to run in Classic Mode as an ISAPI extension or in Integrated Mode where manage code is integrated into the request processing pipeline",
				Type:        schema.TypeString,
				Computed:    true,
			},
			applicationPoolSchema.RuntimeVersion: {
				Description: "Configures the application pool to load a specific .NET CLR version. The CLR version chosen should correspond to the appropriate version of the .NET Framework being used by your application",
				Type:        schema.TypeString,
				Computed:    true,
			},
			applicationPoolSchema.Enable32Bit: {
				Description: "If set to true for an application pool on a 64-bit operating system, the worker process(es) serving will be in WOW64 (Windows on Windows64) mode. Processes in WOW64 mode are 32-bit processes that load only 32-bit applications",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			applicationPoolSchema.QueueLength: {
				Description: "Maximum number of requests that HTTP.sys will queue for the application pool. When the queue is full, new requests receive a 503 (Service Unavailable) response",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			applicationPoolSchema.ProcessModelSchema.Key: {
				Description: "Defines the process model settings for the application pool",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						applicationPoolSchema.ProcessModelSchema.IdentityType: {
							Description: "Configures the application pool to run as a built-in account",
							Type:        schema.TypeString,
							Computed:    true,
						},
						applicationPoolSchema.ProcessModelSchema.Username: {
							Description: "Configures the username for the specified identity",
							Type:        schema.TypeString,
							Computed:    true,
						},
						applicationPoolSchema.ProcessModelSchema.Password: {
							Description: "Configures the password for the specified identity",
							Type:        schema.TypeString,
							Computed:    true,
						},
						applicationPoolSchema.ProcessModelSchema.LoadUserProfile: {
							Description: "Specifies whether IIS loads the user profile for an application pool identity",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						applicationPoolSchema.ProcessModelSchema.IdleTimeout: {
							Description: "Amount of trime (in minutes) a worker process will remain idle befor it shuts down. A worker process is idle if it is not processing requests and no new requests are received",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						applicationPoolSchema.ProcessModelSchema.IdleTimeoutAction: {
							Description: "What action to perform when the Idle Time-out duration has been reached",
							Type:        schema.TypeString,
							Computed:    true,
						},
						applicationPoolSchema.ProcessModelSchema.MaxProcesses: {
							Description: "Maximum number of worker processes permitted to service requests for the application pool. If this number is greater than 1, the application pool is a Web Garden. On a NUMA aware system, if this number is 0, IIS will start as many worker processes as tehre are NUMA nodes for optimal performance",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						applicationPoolSchema.ProcessModelSchema.PingingEnabled: {
							Description: "If true, the worker process(es) serving the application pool are pinged periodically to ensure that they are still responsive. This process is called health monitoring",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						applicationPoolSchema.ProcessModelSchema.PingingInterval: {
							Description: "Period of time (in seconds) between health monitoring pings sent to the worker process(es) serving the application",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						applicationPoolSchema.ProcessModelSchema.PingingResponseTime: {
							Description: "Maximum time (in seconds) that a worker process is given to respond to a health monitoring ping. If the worker process does not respond, it is terminated",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						applicationPoolSchema.ProcessModelSchema.StartupTimeLimit: {
							Description: "Period of time (in seconds) a worker process is given to start up and initialize. If the worker process initialization exceeds the startup time limit, it is terminated",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						applicationPoolSchema.ProcessModelSchema.ShutdownTimeLimit: {
							Description: "Period of time (in seconds) a worker process is given to finish processing requests and shut down. If the worker process initialization exceeds the shutdown time limit, it is terminated",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceApplicationPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*agent.Client)
	name := d.Get(applicationPoolSchema.Name).(string)
	appPool, err := client.GetAppPool(name)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = mapAppPoolToResourceData(*appPool, d); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	return nil
}
