package powerflex

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccSDCExpansionResource tests the SDC Expansion Operation
func TestAccSDCExpansionResource(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Create
			{
				Config: ProviderConfigForGatewayTesting + ParseCSVConfig1,
				Check:  resource.TestCheckResourceAttr("powerflex_sdc_expansion.test", "installed_sdc_ips", "10.247.103.163,10.247.103.161,10.247.103.162,10.247.103.160"),
			},
			//Update
			{
				Config: ProviderConfigForGatewayTesting + ParseCSVConfigUpdate,
				Check:  resource.TestCheckResourceAttr("powerflex_sdc_expansion.test", "installed_sdc_ips", "10.247.103.163,10.247.103.161,10.247.103.162,10.247.103.160"),
			},
		},
	})
}

var ParseCSVConfig1 = `
resource "powerflex_sdc_expansion" "test" {
	mdm_ip = "10.247.103.160"
	mdm_password = "Password123"
	lia_password="Password123"
	csv_detail = [
		{
			ip = "10.247.103.160"
			password = "dangerous"
			operating_system = "linux"
			is_mdm_or_tb = "Primary"
			is_sdc = "Yes"
		},
		{
			ip = "10.247.103.161"
			password = "dangerous"
			operating_system = "linux"
			is_mdm_or_tb = "Secondary"
			is_sdc = "Yes"
		},
		{
			ip = "10.247.103.162"
			password = "dangerous"
			operating_system = "linux"
			is_mdm_or_tb = "TB"
			is_sdc = "Yes"
	    },
	    {
			ip = "10.247.103.163"
			password = "dangerous"
			operating_system = "linux"
			is_mdm_or_tb = "Standby"
			is_sdc = "Yes"
   		},
	]
}
`

var ParseCSVConfigUpdate = `
resource "powerflex_sdc_expansion" "test" {
	mdm_ip = "10.247.103.160"
	mdm_password = "Password123"
	lia_password="Password123"
	csv_detail = [
		{
			ip = "10.247.103.160"
			password = "dangerous"
			operating_system = "linux"
			is_mdm_or_tb = "Primary"
			is_sdc = "Yes"
		},
		{
			ip = "10.247.103.161"
			password = "dangerous"
			operating_system = "linux"
			is_mdm_or_tb = "Secondary"
			is_sdc = "Yes"
		},
		{
			ip = "10.247.103.162"
			password = "dangerous"
			operating_system = "linux"
			is_mdm_or_tb = "TB"
			is_sdc = "Yes"
	    },
	    {
			ip = "10.247.103.163"
			password = "dangerous"
			operating_system = "linux"
			is_mdm_or_tb = ""
			is_sdc = "Yes"
   		},
	]
}
`

func TestAccSDCExpansionResourceNegative(t *testing.T) {
	var WithoutPrimary = `
		resource "powerflex_sdc_expansion" "test" {
			mdm_ip = "10.247.103.160"
			mdm_password = "Password123"
			lia_password="Password123"
			csv_detail = [
				{
					ip = "10.247.103.160"
					password = "dangerous"
					operating_system = "linux"
					is_mdm_or_tb = ""
					is_sdc = "Yes"
				},
				{
					ip = "10.247.103.161"
					password = "dangerous"
					operating_system = "linux"
					is_mdm_or_tb = "Secondary"
					is_sdc = "Yes"
				},
				{
					ip = "10.247.103.162"
					password = "dangerous"
					operating_system = "linux"
					is_mdm_or_tb = "TB"
					is_sdc = "Yes"
				},
				{
					ip = "10.247.103.163"
					password = "dangerous"
					operating_system = "linux"
					is_mdm_or_tb = "Standby"
					is_sdc = "Yes"
				},
			]
		}
	`

	var WithoutSecondary = `
		resource "powerflex_sdc_expansion" "test" {
			mdm_ip = "10.247.103.160"
			mdm_password = "Password123"
			lia_password="Password123"
			csv_detail = [
				{
					ip = "10.247.103.160"
					password = "dangerous"
					operating_system = "linux"
					is_mdm_or_tb = "Primary"
					is_sdc = "Yes"
				},
				{
					ip = "10.247.103.161"
					password = "dangerous"
					operating_system = "linux"
					is_mdm_or_tb = ""
					is_sdc = "Yes"
				},
				{
					ip = "10.247.103.162"
					password = "dangerous"
					operating_system = "linux"
					is_mdm_or_tb = "TB"
					is_sdc = "Yes"
				},
				{
					ip = "10.247.103.163"
					password = "dangerous"
					operating_system = "linux"
					is_mdm_or_tb = "Standby"
					is_sdc = "Yes"
				},
			]
		}
	`

	var WithoutTB = `
		resource "powerflex_sdc_expansion" "test" {
			mdm_ip = "10.247.103.160"
			mdm_password = "Password123"
			lia_password="Password123"
			csv_detail = [
				{
					ip = "10.247.103.160"
					password = "dangerous"
					operating_system = "linux"
					is_mdm_or_tb = "Primary"
					is_sdc = "Yes"
				},
				{
					ip = "10.247.103.161"
					password = "dangerous"
					operating_system = "linux"
					is_mdm_or_tb = "Secondary"
					is_sdc = "Yes"
				},
				{
					ip = "10.247.103.162"
					password = "dangerous"
					operating_system = "linux"
					is_mdm_or_tb = ""
					is_sdc = "Yes"
				},
				{
					ip = "10.247.103.163"
					password = "dangerous"
					operating_system = "linux"
					is_mdm_or_tb = "Standby"
					is_sdc = "Yes"
				},
			]
		}
	`

	var WithoutIP = `
		resource "powerflex_sdc_expansion" "test" {
			mdm_ip = "10.247.103.160"
			mdm_password = "Password123"
			lia_password="Password123"
			csv_detail = [
				{
					ip = ""
					password = "dangerous"
					operating_system = "linux"
					is_mdm_or_tb = "Primary"
					is_sdc = "Yes"
				},
				{
					ip = "10.247.103.161"
					password = "dangerous"
					operating_system = "linux"
					is_mdm_or_tb = "Secondary"
					is_sdc = "Yes"
				},
				{
					ip = "10.247.103.162"
					password = "dangerous"
					operating_system = "linux"
					is_mdm_or_tb = "TB"
					is_sdc = "Yes"
				},
				{
					ip = "10.247.103.163"
					password = "dangerous"
					operating_system = "linux"
					is_mdm_or_tb = "Standby"
					is_sdc = "Yes"
				},
			]
		}
	`

	var WrongMDMCred = `
		resource "powerflex_sdc_expansion" "test" {
			mdm_ip = "10.247.103.160"
			mdm_password = "ABCD"
			lia_password="ABCD"
			csv_detail = [
				{
					ip = "10.247.103.160"
					password = "dangerous"
					operating_system = "linux"
					is_mdm_or_tb = "Primary"
					is_sdc = "Yes"
				},
				{
					ip = "10.247.103.161"
					password = "dangerous"
					operating_system = "linux"
					is_mdm_or_tb = "Secondary"
					is_sdc = "Yes"
				},
				{
					ip = "10.247.103.162"
					password = "dangerous"
					operating_system = "linux"
					is_mdm_or_tb = "TB"
					is_sdc = "Yes"
				},
				{
					ip = "10.247.103.163"
					password = "dangerous"
					operating_system = "linux"
					is_mdm_or_tb = "Standby"
					is_sdc = "Yes"
				},
			]
		}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForGatewayTesting + WithoutIP,
				ExpectError: regexp.MustCompile(`.*No IPs were provided on line number 2.*`),
			},
			{
				Config:      ProviderConfigForGatewayTesting + WithoutPrimary,
				ExpectError: regexp.MustCompile(`.*Unable to detect a Primary MDM.*`),
			},
			{
				Config:      ProviderConfigForGatewayTesting + WithoutSecondary,
				ExpectError: regexp.MustCompile(`.*Error For Parse CSV.*`),
			},
			{
				Config:      ProviderConfigForGatewayTesting + WithoutTB,
				ExpectError: regexp.MustCompile(`.*Error For Parse CSV.*`),
			},
			{
				Config:      ProviderConfigForGatewayTesting + WrongMDMCred,
				ExpectError: regexp.MustCompile(`.*Error While Validating MDM Credentials.*`),
			},
			//Create
			{
				Config: ProviderConfigForGatewayTesting + ParseCSVConfig1,
				Check:  resource.TestCheckResourceAttr("powerflex_sdc_expansion.test", "installed_sdc_ips", "10.247.103.163,10.247.103.161,10.247.103.162,10.247.103.160"),
			},
			//Negative After Create
			{
				Config:      ProviderConfigForGatewayTesting + WithoutIP,
				ExpectError: regexp.MustCompile(`.*No IPs were provided on line number 2.*`),
			},
			{
				Config:      ProviderConfigForGatewayTesting + WithoutPrimary,
				ExpectError: regexp.MustCompile(`.*Unable to detect a Primary MDM.*`),
			},
			{
				Config:      ProviderConfigForGatewayTesting + WithoutSecondary,
				ExpectError: regexp.MustCompile(`.*Error For Parse CSV.*`),
			},
			{
				Config:      ProviderConfigForGatewayTesting + WithoutTB,
				ExpectError: regexp.MustCompile(`.*Error For Parse CSV.*`),
			},
			{
				Config:      ProviderConfigForGatewayTesting + WrongMDMCred,
				ExpectError: regexp.MustCompile(`.*Error While Validating MDM Credentials.*`),
			},
		}})
}
