package iis

type applicationPoolSchemaKeys struct {
	Id                 string
	Name               string
	State              string
	AutoStart          string
	StartMode          string
	PipelineMode       string
	RuntimeVersion     string
	Enable32Bit        string
	QueueLength        string
	ProcessModelSchema applicationPoolProcessModelSchemaKeys
}

type applicationPoolProcessModelSchemaKeys struct {
	Key                 string
	IdentityType        string
	Username            string
	Password            string
	LoadUserProfile     string
	IdleTimeout         string
	IdleTimeoutAction   string
	MaxProcesses        string
	PingingEnabled      string
	PingingInterval     string
	PingingResponseTime string
	StartupTimeLimit    string
	ShutdownTimeLimit   string
}

var applicationPoolSchema = applicationPoolSchemaKeys{
	Id:             "id",
	Name:           "name",
	State:          "state",
	AutoStart:      "auto_start",
	StartMode:      "start_mode",
	PipelineMode:   "pipeline_mode",
	RuntimeVersion: "runtime_version",
	Enable32Bit:    "enable_32bit",
	QueueLength:    "queue_length",
	ProcessModelSchema: applicationPoolProcessModelSchemaKeys{
		Key:                 "process_model",
		IdentityType:        "identity_type",
		Username:            "username",
		Password:            "password",
		LoadUserProfile:     "load_user_profile",
		IdleTimeout:         "idle_timeout",
		IdleTimeoutAction:   "idle_timeout_action",
		MaxProcesses:        "max_processes",
		PingingEnabled:      "pinging_enabled",
		PingingInterval:     "pinging_interval",
		PingingResponseTime: "pinging_response_time",
		StartupTimeLimit:    "startup_time_limit",
		ShutdownTimeLimit:   "shutdown_time_limit",
	},
}

type webSiteSchemaKeys struct {
	Id                  string
	Name                string
	ApplicationPoolName string
	State               string
	PhysicalPath        string
	Username            string
	Password            string
	BindingSchema       webSiteBindingSchemaKeys
}

type webSiteBindingSchemaKeys struct {
	Key        string
	Protocol   string
	Ip         string
	Port       string
	HostHeader string
}

var webSiteSchema = webSiteSchemaKeys{
	Id:                  "id",
	Name:                "name",
	ApplicationPoolName: "application_pool_name",
	State:               "state",
	PhysicalPath:        "physical_path",
	Username:            "username",
	Password:            "password",
	BindingSchema: webSiteBindingSchemaKeys{
		Key:        "binding",
		Protocol:   "protocol",
		Ip:         "ip",
		Port:       "port",
		HostHeader: "host_header",
	},
}

type webApplicationSchemaKeys struct {
	Id                  string
	Path                string
	Name                string
	PhysicalPath        string
	Site                string
	ApplicationPoolName string
}

var webAppSchema = webApplicationSchemaKeys{
	Id:                  "id",
	Path:                "path",
	Name:                "name",
	PhysicalPath:        "physical_path",
	Site:                "web_site_name",
	ApplicationPoolName: "application_pool_name",
}
