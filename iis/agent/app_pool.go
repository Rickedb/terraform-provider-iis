package agent

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ApplicationPool struct {
	Id                    string
	Name                  string
	StartMode             string
	PipelineMode          string
	ManagedRuntimeVersion string
	Enable32BitWin64      bool
	QueueLength           int
	CPU                   CPU
	ProcessModel          ProcessModel
	Recycling             Recycling
	RapidFailProtection   Failure
	//ProcessOrphaning      ProcessOrphaning    `json:"Failure"`
}

type ProcessModel struct {
	IdentityType      string
	Username          string
	LoadUserProfile   bool
	IdleTimeout       int
	IdleTimeoutAction string
	MaxProcesses      int
	PingingEnabled    bool
	PingInterval      int
	PingResponseTime  int
	StartupTimeLimit  int
	ShutdownTimeLimit int
}

type applicationPoolResponse struct {
	Name                  string           `json:"Name"`
	State                 AppPoolState     `json:"State"`
	AutoStart             bool             `json:"AutoStart"`
	StartMode             StartMode        `json:"StartMode"`
	PipelineMode          PipelineMode     `json:"ManagedPipelineMode"`
	ManagedRuntimeVersion string           `json:"ManagedRuntimeVersion"`
	Enable32BitWin64      bool             `json:"Enable32BitAppOnWin64"`
	QueueLength           int64            `json:"QueueLength"`
	CPU                   CPU              `json:"Cpu"`
	ProcessModel          JsonProcessModel `json:"ProcessModel"`
	Recycling             Recycling        `json:"Recycling"`
	RapidFailProtection   Failure          `json:"Failure"`
	//ProcessOrphaning      ProcessOrphaning    `json:"Failure"`
}

type AppPoolState string
type StartMode string
type PipelineMode string

type CPU struct {
	Limit int64 `json:"Limit"`
	//LimitInterval            int64  `json:"limit_interval"`
	Action                   string `json:"Action"`
	ProcessorAffinityEnabled bool   `json:"SmpAffinitized"`
	ProcessorAffinityMask32  int64  `json:"SmpProcessorAffinityMask"`
	ProcessorAffinityMask64  int64  `json:"SmpProcessorAffinityMask2"`
}

type JsonProcessModel struct {
	IdentityType      IdentityType      `json:"IdentityType"`
	Username          string            `json:"UserName"`
	LoadUserProfile   bool              `json:"LoadUserProfile"`
	IdleTimeout       DurationMinutes   `json:"IdleTimeout"`
	IdleTimeoutAction IdleTimeoutAction `json:"IdleTimeoutAction"`
	MaxProcesses      int64             `json:"MaxProcesses"`
	PingingEnabled    bool              `json:"PingingEnabled"`
	PingInterval      DurationSeconds   `json:"PingInterval"`
	PingResponseTime  DurationSeconds   `json:"PingResponseTime"`
	StartupTimeLimit  DurationSeconds   `json:"StartupTimeLimit"`
	ShutdownTimeLimit DurationSeconds   `json:"ShutdownTimeLimit"`
}

type IdentityType string
type IdleTimeoutAction string

type Failure struct {
	OrphanWorkerProcessEnabled    bool   `json:"OrphanWorkerProcess"`
	OrphanActionExe               string `json:"OrphanActionExe"`
	OrphanActionParams            string `json:"OrphanActionParams"`
	RapidFailProtectionEnabled    bool   `json:"RapidFailProtection"`
	LoadBalancerCapabilities      int    `json:"LoadBalancerCapabilities"`
	RapidFailProtectionMaxCrashes int64  `json:"RapidFailProtectionMaxCrashes"`
	AutoShutdownExe               string `json:"AutoShutdownExe"`
	AutoShutdownParams            string `json:"AutoShutdownParams"`
	//Interval                 int64  `json:"interval"`
}

type Recycling struct {
	DisableOverlappedRecycle     bool            `json:"DisallowOverlappingRotation"`
	DisableRecycleOnConfigChange bool            `json:"DisallowRotationOnConfigChange"`
	PeriodicRestart              PeriodicRestart `json:"PeriodicRestart"`
}

type LogEvents struct {
	Time           bool `json:"time"`
	Requests       bool `json:"requests"`
	Schedule       bool `json:"schedule"`
	Memory         bool `json:"memory"`
	IsapiUnhealthy bool `json:"isapi_unhealthy"`
	OnDemand       bool `json:"on_demand"`
	ConfigChange   bool `json:"config_change"`
	PrivateMemory  bool `json:"private_memory"`
}

type PeriodicRestart struct {
	//TimeInterval  int64         `json:"time_interval"`
	PrivateMemory int64 `json:"PrivateMemory"`
	RequestLimit  int64 `json:"Requests"`
	VirtualMemory int64 `json:"Memory"`
}

func (client Client) GetAppPool(name string) (*ApplicationPool, error) {
	var response applicationPoolResponse
	command := fmt.Sprintf("Get-IISAppPool -Name '%s' -WarningAction Stop | ConvertTo-Json -Compress", name)
	bytes, err := client.Execute(command)
	if err != nil {
		return nil, err
	}

	if len(*bytes) == 0 {
		return nil, fmt.Errorf("application pool '%s' could not be found at the host", name)
	}

	json.Unmarshal(*bytes, &response)
	appPool := mapToApplicationPool(&response)
	return appPool, nil
}

func (client Client) DeleteAppPool(name string) error {
	command := fmt.Sprintf("Remove-WebAppPool -Name %q", name)
	_, err := client.Execute(command)
	if err != nil {
		return err
	}
	return nil
}

func (client Client) CreateAppPool(appPool ApplicationPool) (*ApplicationPool, error) {
	_, err := client.Execute(fmt.Sprintf(`New-WebAppPool -Name %q;`, appPool.Name))
	if err == nil {
		err = client.UpdateAppPool(appPool)
		if err != nil {
			client.DeleteAppPool(appPool.Name)
			return nil, err
		}
	}

	return client.GetAppPool(appPool.Name)
}

func (client Client) UpdateAppPool(appPool ApplicationPool) error {
	existingAppPool, err := client.GetAppPool(appPool.Name)
	if err != nil {
		return err
	}

	err = client.updateAppPool(appPool)
	if err != nil {
		client.updateAppPool(*existingAppPool)
		return err
	}

	return nil
}

func (client Client) updateAppPool(appPool ApplicationPool) error {
	var sb strings.Builder
	sb.WriteString(`Import-Module WebAdministration;`)
	setProp := fmt.Sprintf(`Set-ItemProperty -Path 'IIS:\AppPools\%s'`, appPool.Name)
	sb.WriteString(fmt.Sprintf(`%s startMode %q;`, setProp, appPool.StartMode))
	sb.WriteString(fmt.Sprintf(`%s managedPipelineMode %q;`, setProp, appPool.PipelineMode))
	sb.WriteString(fmt.Sprintf(`%s managedRuntimeVersion %q;`, setProp, appPool.ManagedRuntimeVersion))
	sb.WriteString(fmt.Sprintf(`%s enable32BitAppOnWin64 %q;`, setProp, toPascalCase(appPool.Enable32BitWin64)))
	sb.WriteString(fmt.Sprintf(`%s queueLength %d;`, setProp, appPool.QueueLength))
	sb.WriteString(fmt.Sprintf(`%s processModel.identityType %q;`, setProp, appPool.ProcessModel.IdentityType))
	sb.WriteString(fmt.Sprintf(`%s processModel.username %q;`, setProp, appPool.ProcessModel.Username))
	sb.WriteString(fmt.Sprintf(`%s processModel.loadUserProfile %q;`, setProp, toPascalCase(appPool.ProcessModel.LoadUserProfile)))
	sb.WriteString(fmt.Sprintf(`%s processModel.idleTimeout %q;`, setProp, DurationMinutes(appPool.ProcessModel.IdleTimeout).toTimeString()))
	sb.WriteString(fmt.Sprintf(`%s processModel.idleTimeoutAction %q;`, setProp, appPool.ProcessModel.IdleTimeoutAction))
	sb.WriteString(fmt.Sprintf(`%s processModel.maxProcesses %d;`, setProp, appPool.ProcessModel.MaxProcesses))
	sb.WriteString(fmt.Sprintf(`%s processModel.pingingEnabled %q;`, setProp, toPascalCase(appPool.ProcessModel.PingingEnabled)))
	sb.WriteString(fmt.Sprintf(`%s processModel.pingInterval %q;`, setProp, DurationSeconds(appPool.ProcessModel.PingInterval).toTimeString()))
	sb.WriteString(fmt.Sprintf(`%s processModel.pingResponseTime %q;`, setProp, DurationSeconds(appPool.ProcessModel.PingResponseTime).toTimeString()))
	sb.WriteString(fmt.Sprintf(`%s processModel.startupTimeLimit %q;`, setProp, DurationSeconds(appPool.ProcessModel.StartupTimeLimit).toTimeString()))
	sb.WriteString(fmt.Sprintf(`%s processModel.shutdownTimeLimit %q;`, setProp, DurationSeconds(appPool.ProcessModel.ShutdownTimeLimit).toTimeString()))
	_, err := client.Execute(sb.String())
	return err
}

func mapToApplicationPool(response *applicationPoolResponse) *ApplicationPool {
	return &ApplicationPool{
		Id:                    response.Name,
		Name:                  response.Name,
		StartMode:             string(response.StartMode),
		PipelineMode:          string(response.PipelineMode),
		ManagedRuntimeVersion: response.ManagedRuntimeVersion,
		Enable32BitWin64:      response.Enable32BitWin64,
		QueueLength:           int(response.QueueLength),
		ProcessModel: ProcessModel{
			IdentityType:      string(response.ProcessModel.IdentityType),
			Username:          response.ProcessModel.Username,
			LoadUserProfile:   response.ProcessModel.LoadUserProfile,
			IdleTimeout:       int(response.ProcessModel.IdleTimeout),
			IdleTimeoutAction: string(response.ProcessModel.IdleTimeoutAction),
			MaxProcesses:      int(response.ProcessModel.MaxProcesses),
			PingingEnabled:    response.ProcessModel.PingingEnabled,
			PingInterval:      int(response.ProcessModel.PingInterval),
			PingResponseTime:  int(response.ProcessModel.PingResponseTime),
			StartupTimeLimit:  int(response.ProcessModel.StartupTimeLimit),
			ShutdownTimeLimit: int(response.ProcessModel.ShutdownTimeLimit),
		},
		CPU: response.CPU,
	}
}

func (state *AppPoolState) UnmarshalJSON(data []byte) error {
	var number int
	if err := json.Unmarshal(data, &number); err != nil {
		return err
	}

	switch number {
	case 0:
		*state = "Starting"
	case 1:
		*state = "Started"
	case 2:
		*state = "Stopping"
	case 3:
		*state = "Stopped"
	default:
		*state = "Unknown"
	}
	return nil
}

func (state *StartMode) UnmarshalJSON(data []byte) error {
	var number int
	if err := json.Unmarshal(data, &number); err != nil {
		return err
	}

	switch number {
	case 0:
		*state = "OnDemand"
	case 1:
		*state = "AlwaysRunning"
	default:
		*state = "Unknown"
	}
	return nil
}

func (state *PipelineMode) UnmarshalJSON(data []byte) error {
	var number int
	if err := json.Unmarshal(data, &number); err != nil {
		return err
	}

	switch number {
	case 0:
		*state = "Integrated"
	case 1:
		*state = "Classic"
	default:
		*state = "Unknown"
	}
	return nil
}

func (state *IdentityType) UnmarshalJSON(data []byte) error {
	var number int
	if err := json.Unmarshal(data, &number); err != nil {
		return err
	}

	switch number {
	case 0:
		*state = "LocalSystem"
	case 1:
		*state = "LocalService"
	case 2:
		*state = "NetworkService"
	case 3:
		*state = "SpecificUser"
	case 4:
		*state = "ApplicationPoolIdentity"
	default:
		*state = "Unknown"
	}
	return nil
}

func (state *IdleTimeoutAction) UnmarshalJSON(data []byte) error {
	var number int
	if err := json.Unmarshal(data, &number); err != nil {
		return err
	}

	switch number {
	case 0:
		*state = "Terminate"
	case 1:
		*state = "Suspend"
	default:
		*state = "Unknown"
	}
	return nil
}
