package iis

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rickedb/terraform-provider-iis/iis/agent"
)

func resourceApplicationPool() *schema.Resource {
	return &schema.Resource{
		Description: "Enable you to host multiple web applications on a single server in isolation mode for improved security, availability, and performance",

		CreateContext: resourceApplicationPoolCreate,
		ReadContext:   resourceApplicationPoolRead,
		UpdateContext: resourceApplicationPoolUpdate,
		DeleteContext: resourceApplicationPoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importApplicationPoolState,
		},
		Schema: map[string]*schema.Schema{
			applicationPoolSchema.Name: {
				Description: "The application pool name is the unique identifier for the application pool",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			applicationPoolSchema.State: {
				Description: "The current state of the application pool",
				Type:        schema.TypeString,
				Computed:    true,
			},
			applicationPoolSchema.StartMode: {
				Description:      "Configures application pool to run in On Demand Mode or Always Running Mode",
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "OnDemand",
				ValidateDiagFunc: validateAllowedValues([]string{"OnDemand", "AlwaysRunning"}),
			},
			applicationPoolSchema.PipelineMode: {
				Description:      "Configures ASP.NET to run in Classic Mode as an ISAPI extension or in Integrated Mode where manage code is integrated into the request processing pipeline",
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "Integrated",
				ValidateDiagFunc: validateAllowedValues([]string{"Integrated", "Classic"}),
			},
			applicationPoolSchema.RuntimeVersion: {
				Description:      "Configures the application pool to load a specific .NET CLR version. The CLR version chosen should correspond to the appropriate version of the .NET Framework being used by your application",
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "v4.0",
				ValidateDiagFunc: validateAllowedValues([]string{"v4.0", "v2.0", ""}),
			},
			applicationPoolSchema.Enable32Bit: {
				Description: "If set to true for an application pool on a 64-bit operating system, the worker process(es) serving will be in WOW64 (Windows on Windows64) mode. Processes in WOW64 mode are 32-bit processes that load only 32-bit applications",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			applicationPoolSchema.QueueLength: {
				Description:      "Maximum number of requests that HTTP.sys will queue for the application pool. When the queue is full, new requests receive a 503 (Service Unavailable) response",
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1000,
				ValidateDiagFunc: isInBetweenValues(10, 65535),
			},
			applicationPoolSchema.ProcessModelSchema.Key: {
				Description: "Defines the process model settings for the application pool",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: processModelSchema,
				},
			},
		},
	}
}

var processModelSchema = map[string]*schema.Schema{
	applicationPoolSchema.ProcessModelSchema.IdentityType: {
		Description:      "Configures the application pool to run as a built-in account",
		Type:             schema.TypeString,
		Optional:         true,
		Default:          "ApplicationPoolIdentity",
		ValidateDiagFunc: validateAllowedValues([]string{"LocalSystem", "LocalService", "NetworkService", "SpecificUser", "ApplicationPoolIdentity"}),
	},
	applicationPoolSchema.ProcessModelSchema.Username: {
		Description: "Configures the username for the specified identity",
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
	},
	applicationPoolSchema.ProcessModelSchema.Password: {
		Description: "Configures the password for the specified identity",
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
	},
	applicationPoolSchema.ProcessModelSchema.LoadUserProfile: {
		Description: "Specifies whether IIS loads the user profile for an application pool identity",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
	},
	applicationPoolSchema.ProcessModelSchema.IdleTimeout: {
		Description:      "Amount of trime (in minutes) a worker process will remain idle befor it shuts down. A worker process is idle if it is not processing requests and no new requests are received",
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          20,
		ValidateDiagFunc: isInBetweenValues(0, 43200),
	},
	applicationPoolSchema.ProcessModelSchema.IdleTimeoutAction: {
		Description:      "What action to perform when the Idle Time-out duration has been reached",
		Type:             schema.TypeString,
		Optional:         true,
		Default:          "Terminate",
		ValidateDiagFunc: validateAllowedValues([]string{"Terminate", "Suspend"}),
	},
	applicationPoolSchema.ProcessModelSchema.MaxProcesses: {
		Description:      "Maximum number of worker processes permitted to service requests for the application pool. If this number is greater than 1, the application pool is a Web Garden. On a NUMA aware system, if this number is 0, IIS will start as many worker processes as tehre are NUMA nodes for optimal performance",
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          1,
		ValidateDiagFunc: greaterOrEqualThan(0),
	},
	applicationPoolSchema.ProcessModelSchema.PingingEnabled: {
		Description: "If true, the worker process(es) serving the application pool are pinged periodically to ensure that they are still responsive. This process is called health monitoring",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
	},
	applicationPoolSchema.ProcessModelSchema.PingingInterval: {
		Description:      "Period of time (in seconds) between health monitoring pings sent to the worker process(es) serving the application",
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          30,
		ValidateDiagFunc: isInBetweenValues(1, 4294967),
	},
	applicationPoolSchema.ProcessModelSchema.PingingResponseTime: {
		Description:      "Maximum time (in seconds) that a worker process is given to respond to a health monitoring ping. If the worker process does not respond, it is terminated",
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          90,
		ValidateDiagFunc: isInBetweenValues(1, 4294967),
	},
	applicationPoolSchema.ProcessModelSchema.StartupTimeLimit: {
		Description:      "Period of time (in seconds) a worker process is given to start up and initialize. If the worker process initialization exceeds the startup time limit, it is terminated",
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          90,
		ValidateDiagFunc: isInBetweenValues(1, 4294967),
	},
	applicationPoolSchema.ProcessModelSchema.ShutdownTimeLimit: {
		Description:      "Period of time (in seconds) a worker process is given to finish processing requests and shut down. If the worker process initialization exceeds the shutdown time limit, it is terminated",
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          90,
		ValidateDiagFunc: isInBetweenValues(1, 4294967),
	},
}

func resourceApplicationPoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*agent.Client)

	appPoolRequest := mapToApplicationPool(d)
	appPool, err := client.CreateAppPool(appPoolRequest)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(appPool.Id)
	return nil
}

func resourceApplicationPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*agent.Client)
	name := d.Get(applicationPoolSchema.Name).(string)
	appPool, err := client.GetAppPool(name)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	if err = mapAppPoolToResourceData(*appPool, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceApplicationPoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*agent.Client)

	appPool := mapToApplicationPool(d)
	err := client.UpdateAppPool(appPool)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(appPool.Name)
	return nil
}

func resourceApplicationPoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*agent.Client)
	name := d.Get(applicationPoolSchema.Name).(string)
	err := client.DeleteAppPool(name)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func importApplicationPoolState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*agent.Client)
	appPoolName := d.Id()
	appPool, err := client.GetAppPool(appPoolName)
	if err != nil {
		d.SetId("")
		return nil, err
	}

	if err = mapAppPoolToResourceData(*appPool, d); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func mapToApplicationPool(d *schema.ResourceData) agent.ApplicationPool {
	var processModel agent.ProcessModel
	processModelResourceList := d.Get(applicationPoolSchema.ProcessModelSchema.Key).([]interface{})
	if len(processModelResourceList) > 0 {
		processModelResource := processModelResourceList[0].(map[string]interface{})
		processModel = agent.ProcessModel{
			IdentityType:      processModelResource[applicationPoolSchema.ProcessModelSchema.IdentityType].(string),
			Username:          processModelResource[applicationPoolSchema.ProcessModelSchema.Username].(string),
			LoadUserProfile:   processModelResource[applicationPoolSchema.ProcessModelSchema.LoadUserProfile].(bool),
			IdleTimeout:       processModelResource[applicationPoolSchema.ProcessModelSchema.IdleTimeout].(int),
			IdleTimeoutAction: processModelResource[applicationPoolSchema.ProcessModelSchema.IdleTimeoutAction].(string),
			MaxProcesses:      processModelResource[applicationPoolSchema.ProcessModelSchema.MaxProcesses].(int),
			PingingEnabled:    processModelResource[applicationPoolSchema.ProcessModelSchema.PingingEnabled].(bool),
			PingInterval:      processModelResource[applicationPoolSchema.ProcessModelSchema.PingingInterval].(int),
			PingResponseTime:  processModelResource[applicationPoolSchema.ProcessModelSchema.PingingResponseTime].(int),
			StartupTimeLimit:  processModelResource[applicationPoolSchema.ProcessModelSchema.StartupTimeLimit].(int),
			ShutdownTimeLimit: processModelResource[applicationPoolSchema.ProcessModelSchema.ShutdownTimeLimit].(int),
		}
	} else {
		processModel = agent.ProcessModel{
			IdentityType:      processModelSchema[applicationPoolSchema.ProcessModelSchema.IdentityType].Default.(string),
			Username:          processModelSchema[applicationPoolSchema.ProcessModelSchema.Username].Default.(string),
			LoadUserProfile:   processModelSchema[applicationPoolSchema.ProcessModelSchema.LoadUserProfile].Default.(bool),
			IdleTimeout:       processModelSchema[applicationPoolSchema.ProcessModelSchema.IdleTimeout].Default.(int),
			IdleTimeoutAction: processModelSchema[applicationPoolSchema.ProcessModelSchema.IdleTimeoutAction].Default.(string),
			MaxProcesses:      processModelSchema[applicationPoolSchema.ProcessModelSchema.MaxProcesses].Default.(int),
			PingingEnabled:    processModelSchema[applicationPoolSchema.ProcessModelSchema.PingingEnabled].Default.(bool),
			PingInterval:      processModelSchema[applicationPoolSchema.ProcessModelSchema.PingingInterval].Default.(int),
			PingResponseTime:  processModelSchema[applicationPoolSchema.ProcessModelSchema.PingingResponseTime].Default.(int),
			StartupTimeLimit:  processModelSchema[applicationPoolSchema.ProcessModelSchema.StartupTimeLimit].Default.(int),
			ShutdownTimeLimit: processModelSchema[applicationPoolSchema.ProcessModelSchema.ShutdownTimeLimit].Default.(int),
		}
	}

	appPool := agent.ApplicationPool{
		Name:                  d.Get(applicationPoolSchema.Name).(string),
		StartMode:             d.Get(applicationPoolSchema.StartMode).(string),
		PipelineMode:          d.Get(applicationPoolSchema.PipelineMode).(string),
		ManagedRuntimeVersion: d.Get(applicationPoolSchema.RuntimeVersion).(string),
		Enable32BitWin64:      d.Get(applicationPoolSchema.Enable32Bit).(bool),
		QueueLength:           d.Get(applicationPoolSchema.QueueLength).(int),
		ProcessModel:          processModel,
	}

	return appPool
}

func mapAppPoolToResourceData(appPool agent.ApplicationPool, d *schema.ResourceData) error {
	var err error
	d.SetId(appPool.Id)
	if err = d.Set(applicationPoolSchema.Name, appPool.Name); err != nil {
		return err
	}
	if err = d.Set(applicationPoolSchema.StartMode, appPool.StartMode); err != nil {
		return err
	}
	if err = d.Set(applicationPoolSchema.PipelineMode, appPool.PipelineMode); err != nil {
		return err
	}
	if err = d.Set(applicationPoolSchema.RuntimeVersion, appPool.ManagedRuntimeVersion); err != nil {
		return err
	}
	if err = d.Set(applicationPoolSchema.Enable32Bit, appPool.Enable32BitWin64); err != nil {
		return err
	}
	if err = d.Set(applicationPoolSchema.QueueLength, appPool.QueueLength); err != nil {
		return err
	}

	processModel := map[string]interface{}{
		applicationPoolSchema.ProcessModelSchema.IdentityType:        appPool.ProcessModel.IdentityType,
		applicationPoolSchema.ProcessModelSchema.Username:            appPool.ProcessModel.Username,
		applicationPoolSchema.ProcessModelSchema.LoadUserProfile:     appPool.ProcessModel.LoadUserProfile,
		applicationPoolSchema.ProcessModelSchema.IdleTimeout:         appPool.ProcessModel.IdleTimeout,
		applicationPoolSchema.ProcessModelSchema.IdleTimeoutAction:   appPool.ProcessModel.IdleTimeoutAction,
		applicationPoolSchema.ProcessModelSchema.MaxProcesses:        appPool.ProcessModel.MaxProcesses,
		applicationPoolSchema.ProcessModelSchema.PingingEnabled:      appPool.ProcessModel.PingingEnabled,
		applicationPoolSchema.ProcessModelSchema.PingingInterval:     appPool.ProcessModel.PingInterval,
		applicationPoolSchema.ProcessModelSchema.PingingResponseTime: appPool.ProcessModel.PingResponseTime,
		applicationPoolSchema.ProcessModelSchema.StartupTimeLimit:    appPool.ProcessModel.StartupTimeLimit,
		applicationPoolSchema.ProcessModelSchema.ShutdownTimeLimit:   appPool.ProcessModel.ShutdownTimeLimit,
	}

	err = d.Set(applicationPoolSchema.ProcessModelSchema.Key, []interface{}{processModel})
	return err
}

func validateAppPoolExists(client *agent.Client, appPoolName string) error {
	_, err := client.GetAppPool(appPoolName)
	if err != nil {
		return err
	}

	return nil
}
