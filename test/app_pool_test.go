package test

import (
	"testing"

	"github.com/rickedb/terraform-provider-iis/iis/agent"
)

func TestGetAppPool(t *testing.T) {
	// resource.Test(t, resource.TestCase{
	// 	PreCheck: func() {
	// 		// Perform any setup needed before the test
	// 	},
	// 	Providers: map[string]*schema.Provider{
	// 		"iis": iis.Provider(),
	// 	},
	// 	Steps: []resource.TestStep{
	// 		{
	// 			Config: `
	//             resource "iis_application_pool" "example_pool" {
	//                 name = "test_app_pool"
	//             }
	//             `,
	// 			Check: resource.ComposeTestCheckFunc(
	// 				resource.TestCheckResourceAttr("iis_application_pool.example_pool", "name", "test_app_pool"),
	// 				//resource.TestCheckResourceAttr("iis_application_pool.example_pool", "pipeline_mode", "Classic"),
	// 			),
	// 		},
	// 	},
	// })
	client := agent.Client{}

	client.GetAppPool("Default App Pool")
}

func TestCreateAppPool(t *testing.T) {
	client := agent.Client{}

	pool := agent.ApplicationPool{
		Name:         "abcapp",
		PipelineMode: "Classic",
	}

	res, _ := client.CreateAppPool(pool)
	print(res)
}

func TestUpdateAppPool(t *testing.T) {
	client := agent.Client{}

	pool := agent.ApplicationPool{
		Name:             "IntegrationTestPool",
		Enable32BitWin64: true,
		QueueLength:      20,
		PipelineMode:     "Integrated",
	}
	client.UpdateAppPool(pool)
}

func TestDeleteAppPool(t *testing.T) {
	client := agent.Client{}

	client.DeleteAppPool("NewApp")
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
