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
				Config:      ProviderConfigForGatewayTesting + ParseCSVConfig1,
				ExpectError: regexp.MustCompile(`.*Error During Installation.*`),
			},
			//Create with Packages
			{
				Config: ProviderConfigForGatewayTesting + packageTest + ParseCSVConfig2,
				Check:  resource.TestMatchResourceAttr("powerflex_sdc_expansion.test", "installed_sdc_ips", regexp.MustCompile(GatewayDataPoints.tbIP)),
			},
			//Update
			{
				Config: ProviderConfigForGatewayTesting + packageTest + ParseCSVConfigUpdate,
				Check:  resource.TestMatchResourceAttr("powerflex_sdc_expansion.test", "installed_sdc_ips", regexp.MustCompile(GatewayDataPoints.sdcServerIP)),
			},
		},
	})
}

var packageTest = `
resource "powerflex_package" "upload-test" {
	file_path = ["/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-lia-3.6-700.103.el7.x86_64.rpm",
	"/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-mdm-3.6-700.103.el7.x86_64.rpm",
	"/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-sds-3.6-700.103.el7.x86_64.rpm",
	"/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-sdc-3.6-700.103.el7.x86_64.rpm",
	"/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-sdr-3.6-700.103.el7.x86_64.rpm"]
	}
`
var ParseCSVConfig1 = `
resource "powerflex_sdc_expansion" "test" {
	mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
	lia_password= "` + GatewayDataPoints.liaPassword + `"
	cluster_details = [
		{
			ip = "` + GatewayDataPoints.primaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Primary"
			is_sdc = "No"
		},
		{
			ip = "` + GatewayDataPoints.secondaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Secondary"
			is_sdc = "NO"
		},
		{
			ip = "` + GatewayDataPoints.tbIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "TB"
			is_sdc = "Yes"
	    },
	    {
			ip = "` + GatewayDataPoints.sdcServerIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Standby"
			is_sdc = "No"
   		},
	]
}
`

var ParseCSVConfig2 = `
resource "powerflex_sdc_expansion" "test" {

	depends_on = [
		powerflex_package.upload-test
	]

	mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
	lia_password= "` + GatewayDataPoints.liaPassword + `"
	cluster_details = [
		{
			ip = "` + GatewayDataPoints.primaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Primary"
			is_sdc = "No"
		},
		{
			ip = "` + GatewayDataPoints.secondaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Secondary"
			is_sdc = "No"
		},
		{
			ip = "` + GatewayDataPoints.tbIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "TB"
			is_sdc = "Yes"
	    },
	    {
			ip = "` + GatewayDataPoints.sdcServerIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = ""
			is_sdc = "No"
   		},
	]
}
`

var ParseCSVConfigUpdate = `
resource "powerflex_sdc_expansion" "test" {

	depends_on = [
		powerflex_package.upload-test
	]
	
	mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
	lia_password= "` + GatewayDataPoints.liaPassword + `"
	cluster_details = [
		{
			ip = "` + GatewayDataPoints.primaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Primary"
			is_sdc = "No"
		},
		{
			ip = "` + GatewayDataPoints.secondaryMDMIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "Secondary"
			is_sdc = "No"
		},
		{
			ip = "` + GatewayDataPoints.tbIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
			operating_system = "linux"
			is_mdm_or_tb = "TB"
			is_sdc = "Yes"
	    },
	    {
			ip = "` + GatewayDataPoints.sdcServerIP + `"
			password = "` + GatewayDataPoints.serverPassword + `"
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
		mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
		lia_password= "` + GatewayDataPoints.liaPassword + `"
		cluster_details = [
			{
				ip = "` + GatewayDataPoints.primaryMDMIP + `"
				password = "` + GatewayDataPoints.serverPassword + `"
				operating_system = "linux"
				is_mdm_or_tb = ""
				is_sdc = "Yes"
			},
			{
				ip = "` + GatewayDataPoints.secondaryMDMIP + `"
				password = "` + GatewayDataPoints.serverPassword + `"
				operating_system = "linux"
				is_mdm_or_tb = "Secondary"
				is_sdc = "Yes"
			},
			{
				ip = "` + GatewayDataPoints.tbIP + `"
				password = "` + GatewayDataPoints.serverPassword + `"
				operating_system = "linux"
				is_mdm_or_tb = "TB"
				is_sdc = "Yes"
			},
			{
				ip = "` + GatewayDataPoints.sdcServerIP + `"
				password = "` + GatewayDataPoints.serverPassword + `"
				operating_system = "linux"
				is_mdm_or_tb = ""
				is_sdc = "Yes"
			   },
		]
	}
	`

	var WithoutSecondary = `
	resource "powerflex_sdc_expansion" "test" {
		mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
		lia_password= "` + GatewayDataPoints.liaPassword + `"
		cluster_details = [
			{
				ip = "` + GatewayDataPoints.primaryMDMIP + `"
				password = "` + GatewayDataPoints.serverPassword + `"
				operating_system = "linux"
				is_mdm_or_tb = "Primary"
				is_sdc = "Yes"
			},
			{
				ip = "` + GatewayDataPoints.secondaryMDMIP + `"
				password = "` + GatewayDataPoints.serverPassword + `"
				operating_system = "linux"
				is_mdm_or_tb = ""
				is_sdc = "Yes"
			},
			{
				ip = "` + GatewayDataPoints.tbIP + `"
				password = "` + GatewayDataPoints.serverPassword + `"
				operating_system = "linux"
				is_mdm_or_tb = "TB"
				is_sdc = "Yes"
			},
			{
				ip = "` + GatewayDataPoints.sdcServerIP + `"
				password = "` + GatewayDataPoints.serverPassword + `"
				operating_system = "linux"
				is_mdm_or_tb = ""
				is_sdc = "Yes"
			   },
		]
	}
	`

	var WithoutTB = `
	resource "powerflex_sdc_expansion" "test" {
		mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
		lia_password= "` + GatewayDataPoints.liaPassword + `"
		cluster_details = [
			{
				ip = "` + GatewayDataPoints.primaryMDMIP + `"
				password = "` + GatewayDataPoints.serverPassword + `"
				operating_system = "linux"
				is_mdm_or_tb = "Primary"
				is_sdc = "Yes"
			},
			{
				ip = "` + GatewayDataPoints.secondaryMDMIP + `"
				password = "` + GatewayDataPoints.serverPassword + `"
				operating_system = "linux"
				is_mdm_or_tb = "Secondary"
				is_sdc = "Yes"
			},
			{
				ip = "` + GatewayDataPoints.tbIP + `"
				password = "` + GatewayDataPoints.serverPassword + `"
				operating_system = "linux"
				is_mdm_or_tb = ""
				is_sdc = "Yes"
			},
			{
				ip = "` + GatewayDataPoints.sdcServerIP + `"
				password = "` + GatewayDataPoints.serverPassword + `"
				operating_system = "linux"
				is_mdm_or_tb = ""
				is_sdc = "Yes"
			   },
		]
	}
	`

	var WithoutIP = `
	resource "powerflex_sdc_expansion" "test" {
		mdm_password =  "` + GatewayDataPoints.mdmPassword + `"
		lia_password= "` + GatewayDataPoints.liaPassword + `"
		cluster_details = [
			{
				ip = "` + GatewayDataPoints.primaryMDMIP + `"
				password = "` + GatewayDataPoints.serverPassword + `"
				operating_system = "linux"
				is_mdm_or_tb = "Primary"
				is_sdc = "Yes"
			},
			{
				ip = ""
				password = "` + GatewayDataPoints.serverPassword + `"
				operating_system = "linux"
				is_mdm_or_tb = "Secondary"
				is_sdc = "Yes"
			},
			{
				ip = "` + GatewayDataPoints.tbIP + `"
				password = "` + GatewayDataPoints.serverPassword + `"
				operating_system = "linux"
				is_mdm_or_tb = "TB"
				is_sdc = "Yes"
			},
			{
				ip = "` + GatewayDataPoints.sdcServerIP + `"
				password = "` + GatewayDataPoints.serverPassword + `"
				operating_system = "linux"
				is_mdm_or_tb = ""
				is_sdc = "Yes"
			   },
		]
	}
	`

	var WrongMDMCred = `
	resource "powerflex_sdc_expansion" "test" {
		mdm_password =  "ABCD"
		lia_password= "ABCD"
		cluster_details = [
			{
				ip = "` + GatewayDataPoints.primaryMDMIP + `"
				password = "` + GatewayDataPoints.serverPassword + `"
				operating_system = "linux"
				is_mdm_or_tb = "Primary"
				is_sdc = "Yes"
			},
			{
				ip = "` + GatewayDataPoints.secondaryMDMIP + `"
				password = "` + GatewayDataPoints.serverPassword + `"
				operating_system = "linux"
				is_mdm_or_tb = "Secondary"
				is_sdc = "Yes"
			},
			{
				ip = "` + GatewayDataPoints.tbIP + `"
				password = "` + GatewayDataPoints.serverPassword + `"
				operating_system = "linux"
				is_mdm_or_tb = "TB"
				is_sdc = "Yes"
			},
			{
				ip = "` + GatewayDataPoints.sdcServerIP + `"
				password = "` + GatewayDataPoints.serverPassword + `"
				operating_system = "linux"
				is_mdm_or_tb = ""
				is_sdc = "Yes"
			   },
		]
	}
	`

	var WithoutSDCYes = `
	resource "powerflex_sdc_expansion" "test" {
		mdm_password =  "ABCD"
		lia_password= "ABCD"
		cluster_details = [
			{
				ip = "` + GatewayDataPoints.primaryMDMIP + `"
				password = "` + GatewayDataPoints.serverPassword + `"
				operating_system = "linux"
				is_mdm_or_tb = "Primary"
				is_sdc = "No"
			},
			{
				ip = "` + GatewayDataPoints.secondaryMDMIP + `"
				password = "` + GatewayDataPoints.serverPassword + `"
				operating_system = "linux"
				is_mdm_or_tb = "Secondary"
				is_sdc = "No"
			},
			{
				ip = "` + GatewayDataPoints.tbIP + `"
				password = "` + GatewayDataPoints.serverPassword + `"
				operating_system = "linux"
				is_mdm_or_tb = "TB"
				is_sdc = "NO"
			},
		]
	}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfigForGatewayTesting + WithoutIP,
				ExpectError: regexp.MustCompile(`.*Error while Parsing CSV.*`),
			},
			{
				Config:      ProviderConfigForGatewayTesting + WithoutPrimary,
				ExpectError: regexp.MustCompile(`.*Error while Parsing CSV.*`),
			},
			{
				Config:      ProviderConfigForGatewayTesting + WithoutSecondary,
				ExpectError: regexp.MustCompile(`.*Error while Parsing CSV.*`),
			},
			{
				Config:      ProviderConfigForGatewayTesting + WithoutTB,
				ExpectError: regexp.MustCompile(`.*Error while Parsing CSV.*`),
			},
			{
				Config:      ProviderConfigForGatewayTesting + WrongMDMCred,
				ExpectError: regexp.MustCompile(`.*Error While Validating MDM Credentials.*`),
			},
			{
				Config:      ProviderConfigForGatewayTesting + WithoutSDCYes,
				ExpectError: regexp.MustCompile(`.*No SDC Expansion Details are provided.*`),
			},
		}})
}
