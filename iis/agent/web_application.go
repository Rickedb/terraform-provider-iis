package agent

import (
	"encoding/json"
	"fmt"
	"strings"
)

type WebApplication struct {
	Id                  string
	Name                string `json:"name"`
	Path                string `json:"path"`
	PhysicalPath        string `json:"PhysicalPath"`
	ApplicationPoolName string `json:"applicationPool"`
	Site                string
}

func (client Client) GetWebApplication(site string, name string) (*WebApplication, error) {
	var response WebApplication
	command := fmt.Sprintf("Get-WebApplication -Site '%s' -Name '%s' | ConvertTo-Json -Compress", site, name)
	bytes, err := client.Execute(command)
	if err != nil {
		return nil, err
	}

	if len(*bytes) == 0 {
		return nil, fmt.Errorf("web application '%s/%s' web site could not be found at the host", site, name)
	}

	json.Unmarshal(*bytes, &response)
	response.Site = site
	response.Name = name
	response.Id = fmt.Sprintf("%s_%s", response.Site, response.Name)
	return &response, nil
}

func (client Client) CreateWebApplication(webApplication WebApplication) (*WebApplication, error) {
	physicalPath := strings.ReplaceAll(webApplication.PhysicalPath, "/", `\`)
	command := fmt.Sprintf(`
		$path='%v'
		if (!(Test-Path $path)){
            New-Item -ItemType Directory -Path $path;
			icacls $path /grant "IIS_IUSRS:(OI)(CI)F" /T
        }
		New-WebApplication -Site %q -ApplicationPool %q -Name %q -PhysicalPath $path;
	`, physicalPath,
		webApplication.Site,
		webApplication.ApplicationPoolName,
		webApplication.Name)

	_, err := client.Execute(command)
	if err != nil {
		return nil, err
	}

	return client.GetWebApplication(webApplication.Site, webApplication.Name)
}

func (client Client) UpdateWebApplication(webApplication WebApplication) error {
	existingWebApp, err := client.GetWebApplication(webApplication.Site, webApplication.Name)
	if err != nil {
		return err
	}

	err = client.updateWebApplication(webApplication)
	if err != nil {
		client.updateWebApplication(*existingWebApp)
		return err
	}

	return err
}

func (client Client) DeleteWebApplication(site string, name string) error {
	command := fmt.Sprintf("Remove-WebApplication -Site %q -Name %q", site, name)
	_, err := client.Execute(command)
	if err != nil {
		return err
	}

	return nil
}

func (client Client) updateWebApplication(webApplication WebApplication) error {
	var sb strings.Builder
	physicalPath := strings.ReplaceAll(webApplication.PhysicalPath, "/", `\`)
	sb.WriteString(`Import-Module WebAdministration;`)
	sb.WriteString(fmt.Sprintf(`
		$path='%v'
		if (!(Test-Path $path)){
            New-Item -ItemType Directory -Path $path;
			icacls $path /grant "IIS_IUSRS:(OI)(CI)F" /T
        }
	`, physicalPath))

	setProp := fmt.Sprintf(`Set-ItemProperty -Path 'IIS:\Sites\%s\%s'`, webApplication.Site, webApplication.Name)
	sb.WriteString(fmt.Sprintf(`%s applicationPool %q;`, setProp, webApplication.ApplicationPoolName))
	sb.WriteString(fmt.Sprintf(`%s physicalPath $path;`, setProp))

	command := sb.String()
	_, err := client.Execute(command)
	if err != nil {
		return err
	}

	return nil
}
