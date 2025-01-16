package agent

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type WebSite struct {
	Id                         string `json:"id"`
	Name                       string `json:"name"`
	State                      string `json:"state"`
	PhysicalPath               string `json:"physicalPath"`
	Username                   string `json:"username"`
	Password                   string `json:"password"`
	Bindings                   []Binding
	TraceFailedRequestsLogging TraceFailedRequestsLogging `json:"traceFailedRequestsLogging"`
	ApplicationPoolName        string
}

type Binding struct {
	Protocol   string
	HostHeader string
	Ip         string
	Port       int
}

type TraceFailedRequestsLogging struct {
	Enabled     bool   `json:"enabled"`
	Directory   string `json:"directory"`
	MaxLogFiles int    `json:"maxLogFiles"`
}

type Hsts struct {
	Enabled             bool `json:"enabled"`
	MaxAge              int  `json:"max-age"`
	IncludeSubDomains   bool `json:"includeSubDomains"`
	Preload             bool `json:"preload"`
	RedirectHttpToHttps bool `json:"redirectHttpToHttps"`
}

type bindingInformationResponse struct {
	Ip         string
	Port       int
	HostHeader string
}

type websiteResponse struct {
	Id              int             `json:"id"`
	Name            string          `json:"name"`
	ServerAutoStart bool            `json:"serverAutoStart"`
	State           string          `json:"state"`
	PhysicalPath    string          `json:"physicalPath"`
	Username        string          `json:"username"`
	Password        string          `json:"password"`
	Bindings        bindingResponse `json:"bindings"`
	ApplicationPool string          `json:"applicationPool"`

	TraceFailedRequestsLogging TraceFailedRequestsLogging `json:"traceFailedRequestsLogging"`
}

type bindingResponse struct {
	Collection []bindingCollectionItemResponse
}

type bindingCollectionItemResponse struct {
	Protocol           string                     `json:"protocol"`
	BindingInformation bindingInformationResponse `json:"bindingInformation"`
}

func (client Client) GetWebSite(name string) (*WebSite, error) {
	var response websiteResponse
	command := fmt.Sprintf("Get-Website -Name '%s' | ConvertTo-Json -Compress", name)
	bytes, err := client.Execute(command)
	if err != nil {
		return nil, err
	}

	if len(*bytes) == 0 {
		return nil, fmt.Errorf("web site '%s' could not be found at the host", name)
	}

	json.Unmarshal(*bytes, &response)
	webSite := mapWebSite(&response)
	return webSite, nil
}

func (client Client) CreateWebSite(webSite WebSite) (*WebSite, error) {
	physicalPath := strings.ReplaceAll(webSite.PhysicalPath, "/", `\`)
	command := fmt.Sprintf(`
		$path='%v'
		if (!(Test-Path $path)){
            New-Item -ItemType Directory -Path $path;
			icacls $path /grant "IIS_IUSRS:(OI)(CI)F" /T
        }
		New-Website -Name %q -PhysicalPath $path;
	`, physicalPath, webSite.Name)
	_, err := client.Execute(command)
	if err != nil {
		return nil, err
	}

	err = client.updateWebSite(webSite)
	if err != nil {
		return nil, client.DeleteWebSite(webSite.Name)
	}

	return client.GetWebSite(webSite.Name)
}

func (client Client) UpdateWebSite(webSite WebSite) error {
	existingWebSite, err := client.GetWebSite(webSite.Name)
	if err != nil {
		return err
	}

	err = client.updateWebSite(webSite)
	if err != nil {
		client.updateWebSite(*existingWebSite)
		return err
	}

	return nil
}

func (client Client) updateWebSite(webSite WebSite) error {
	var sb strings.Builder
	physicalPath := strings.ReplaceAll(webSite.PhysicalPath, "/", `\`)
	sb.WriteString(`Import-Module WebAdministration;`)
	sb.WriteString(fmt.Sprintf(`
		$path='%v'
		if (!(Test-Path $path)){
            New-Item -ItemType Directory -Path $path;
			icacls $path /grant "IIS_IUSRS:(OI)(CI)F" /T
        }
	`, physicalPath))

	setProp := fmt.Sprintf(`Set-ItemProperty -Path 'IIS:\Sites\%s'`, webSite.Name)
	sb.WriteString(fmt.Sprintf(`%s applicationPool %q;`, setProp, webSite.ApplicationPoolName))
	sb.WriteString(fmt.Sprintf(`%s physicalPath %v;`, setProp, physicalPath))
	sb.WriteString(fmt.Sprintf(`%s userName %q;`, setProp, webSite.Username))
	sb.WriteString(fmt.Sprintf(`%s password %q;`, setProp, webSite.Password))
	sb.WriteString(fmt.Sprintf(` 
		$siteName=%q;
		Get-WebBinding -Name $siteName | ForEach-Object { Remove-WebBinding -Name $siteName -BindingInformation $_.bindingInformation -Protocol $_.protocol; };`, webSite.Name))

	for _, binding := range webSite.Bindings {
		str := fmt.Sprintf("New-WebBinding -Name %q -IPAddress %q -Port %d -HostHeader %q -Protocol %q;", webSite.Name, binding.Ip, binding.Port, binding.HostHeader, binding.Protocol)
		sb.WriteString(str)
	}

	command := sb.String()
	_, err := client.Execute(command)
	if err != nil {
		return err
	}

	return nil
}

func (client Client) DeleteWebSite(webSiteName string) error {
	_, err := client.Execute(fmt.Sprintf("Remove-Website -Name %q", webSiteName))
	if err != nil {
		return err
	}

	return nil
}

func mapWebSite(response *websiteResponse) *WebSite {
	bindings := []Binding{}
	for _, binding := range response.Bindings.Collection {
		bindings = append(bindings, Binding{
			Protocol:   binding.Protocol,
			Ip:         binding.BindingInformation.Ip,
			Port:       binding.BindingInformation.Port,
			HostHeader: binding.BindingInformation.HostHeader,
		})
	}
	return &WebSite{
		Id:                         strconv.Itoa(response.Id),
		Name:                       response.Name,
		State:                      response.State,
		PhysicalPath:               response.PhysicalPath,
		Username:                   response.Username,
		Password:                   response.Password,
		Bindings:                   bindings,
		TraceFailedRequestsLogging: response.TraceFailedRequestsLogging,
		ApplicationPoolName:        response.ApplicationPool,
	}
}

func (binding *bindingInformationResponse) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	splitted := strings.Split(str, ":")
	port, err := strconv.Atoi(splitted[1])
	if err != nil {
		return err
	}

	binding.Ip = splitted[0]
	binding.Port = port
	if len(splitted) > 2 {
		binding.HostHeader = splitted[2]
	}

	return nil
}
