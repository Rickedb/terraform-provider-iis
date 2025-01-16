package test

import (
	"testing"

	"github.com/rickedb/terraform-provider-iis/iis/agent"
)

func TestGetWebSite(t *testing.T) {

	client := agent.Client{}

	client.GetWebSite("Default Web S2ite")
}

func TestCreateWebSite(t *testing.T) {

	client := agent.Client{}
	webSite := agent.WebSite{
		Name:                "Test",
		PhysicalPath:        "C:/inetpub/wwwroot/test",
		ApplicationPoolName: "IntegrationTestPool",
		Bindings: []agent.Binding{
			{
				Ip:         "*",
				Protocol:   "http",
				Port:       7272,
				HostHeader: "test",
			},
		},
	}
	client.CreateWebSite(webSite)
}

func TestUpdateWebSite(t *testing.T) {

	client := agent.Client{}
	webSite := agent.WebSite{
		Name:                "Test",
		PhysicalPath:        "C:/inetpub/wwwroot/test",
		ApplicationPoolName: "IntegrationTestPool",
		Bindings: []agent.Binding{
			agent.Binding{
				Ip:         "*",
				Protocol:   "http",
				Port:       7273,
				HostHeader: "test",
			},
			agent.Binding{
				Ip:         "*",
				Protocol:   "http",
				Port:       7271,
				HostHeader: "test",
			},
		},
	}
	client.UpdateWebSite(webSite)
}
