/*
Copyright (c) 2025 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package provider

import (
	"fmt"
	"os"
	"regexp"
	"terraform-provider-powerflex/powerflex/helper"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var exampleSSHKey = `-----BEGIN RSA PRIVATE KEY----- example-ssh-key -----END RSA PRIVATE KEY-----`

var serverCredConfCreate = `
	resource "powerflex_resource_credential" "test" {
	 ## Required values for all credential types
	 name = "test-server-cred-node"
	 type = "Node"
	 password = "example-pass"
	 username = "example-user"
	}
`

var serverCredConfUpdate = `
	resource "powerflex_resource_credential" "test" {
	 ## Required values for all credential types
	 name = "test-server-cred-node"
	 type = "Node"
	 password = "example-pass"
	 username = "example-user"
	 snmp_v3_security_level = "Maximal"
	 snmp_v3_security_name = "example-sec-name"
	 snmp_v3_md5_authentication_password = "md5-password"
	 snmp_v3_des_authentication_password = "des-password"
	}
`

var switchCredConfCreate = `
	resource "powerflex_resource_credential" "test" {
		name = "test-server-cred-switch"
		type = "Switch"
		password = "example-pass"
		username = "example-user"
		ssh_private_key = "` + exampleSSHKey + `"
		key_pair_name = "example_key_pair_name"
	}
`

var switchCredConfUpdateError = `
	resource "powerflex_resource_credential" "test" {
		name = "test-server-cred-switch"
		type = "Switch"
		password = "example-pass"
		username = "example-user"
		key_pair_name = "example_key_pair_name"
	}
`

var switchCredConfUpdateSuccess = `
	resource "powerflex_resource_credential" "test" {
		name = "test-server-cred-switch"
		type = "Switch"
		password = "example-pass"
		username = "example-user"
	}
`

var VCenterCredentialCreate = `
	resource "powerflex_resource_credential" "test" {
		name = "test-server-cred-vcenter"
		type = "vCenter"
		password = "example-pass"
		username = "example-user"
		domain = "example.com"
	}
`

var VCenterCredentialUpdateError = `
	resource "powerflex_resource_credential" "test" {
		name = "test-server-cred-vcenter"
		type = "vCenter"
		password = "example-pass"
		username = "example-user"
	}
`

var VCenterCredentialUpdateSuccess = `
	resource "powerflex_resource_credential" "test" {
		name = "test-server-cred-vcenter"
		type = "vCenter"
		password = "example-pass"
		username = "example-user"
		domain = "example-update.com"
	}
`

var PresentationServerCreate = `
	resource "powerflex_resource_credential" "test" {
		name = "test-server-cred-presentation-server"
		type = "PresentationServer"
		password = "example-pass"
		username = "example-user"
	}
`

var PowerflexGatewayCredentialCreate = `
	resource "powerflex_resource_credential" "test" {
		name = "test-server-cred-powerflex-gateway"
		type = "PowerflexGateway"
		password = "example-pass"
		username = "example-user"
		os_username = "example-os-user"
		os_password = "example-os-pass"
	}
`

var PowerflexGatewayUpdateError = `
	resource "powerflex_resource_credential" "test" {
		name = "test-server-cred-powerflex-gateway"
		type = "PowerflexGateway"
		password = "example-pass"
		username = "example-user"
	}
`

var PowerflexGatewayUpdateSuccess = `
	resource "powerflex_resource_credential" "test" {
		name = "test-server-cred-powerflex-gateway"
		type = "PowerflexGateway"
		password = "example-pass-update"
		username = "example-user"
		os_username = "example-os-user"
		os_password = "example-os-pass"
	}
`

var OsUserCredentialCreate = `
	resource "powerflex_resource_credential" "test" {
		name = "test-server-cred-user"
		type = "OSUser"
		password = "example-pass"
		username = "example-user"
	}
`

var OsUserUpdateSuccess = `
	resource "powerflex_resource_credential" "test" {
		name = "test-server-cred-os-user"
		type = "OSUser"
		password = "example-pass-update"
		username = "example-user"
	}
`

var OsAdminCredentialCreate = `
	resource "powerflex_resource_credential" "test" {
		name = "test-server-cred-os-admin"
		type = "OSAdmin"
		password = "example-pass"
		username = "example-user"
	}
`

var OsAdminUpdateSuccess = `
	resource "powerflex_resource_credential" "test" {
		name = "test-server-cred-os-admin"
		type = "OSAdmin"
		password = "example-pass-update"
		username = "example-user"
	}
`

var ElementManagerCredentialCreate = `
	resource "powerflex_resource_credential" "test" {
		name = "test-server-cred-element-manager"
		type = "ElementManager"
		password = "example-pass"
		username = "example-user"
	}
`

var ElementManagerUpdateSuccess = `
	resource "powerflex_resource_credential" "test" {
		name = "test-server-cred-element-manager"
		type = "ElementManager"
		password = "example-pass-update"
		username = "example-user"
	}
`

// Accptance Tests
func TestAccResourceAcceptanceResourceCredential(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Dont run with units tests, this is an Acceptance test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create Resource Credential
			{
				Config: ProviderConfigForTesting + serverCredConfCreate,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Update Resource Credential
			{
				Config: ProviderConfigForTesting + serverCredConfUpdate,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// Unit Tests
func TestAccResourceResourceCredentialNode(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a Unit test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Get System Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.GetFirstSystem).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + serverCredConfCreate,
				ExpectError: regexp.MustCompile(`.*Unable to Read Powerflex System*.`),
			},
			// Create Resource Credential Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.CreateResourceCredential).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + serverCredConfCreate,
				ExpectError: regexp.MustCompile(`.*Error Creating Resource Credential*.`),
			},
			// Create Resource Credential
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + serverCredConfCreate,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// check that import is working
			{
				ResourceName: "powerflex_resource_credential.test",
				ImportState:  true,
			},
			// Read Resource Credential Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock((*goscaleio.System).GetResourceCredential).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + serverCredConfUpdate,
				ExpectError: regexp.MustCompile(`.*Error Reading Resource Credential*.`),
			},
			// Update Resource Credential Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.ModifyResourceCredential).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfigForTesting + serverCredConfUpdate,
				ExpectError: regexp.MustCompile(`.*Error Modifing Resource Credential*.`),
			},
			// Update Resource Credential
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfigForTesting + serverCredConfUpdate,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Update Resource Credential to different type error
			{
				Config:      ProviderConfigForTesting + switchCredConfCreate,
				ExpectError: regexp.MustCompile(`.*Type cannot be modified*.`),
			},
		},
	})
}

// Unit Tests
func TestAccResourceResourceCredentialSwitch(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a Unit test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create Resource Credential
			{
				Config: ProviderConfigForTesting + switchCredConfCreate,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Update Resource Credential invalid config error
			{
				Config:      ProviderConfigForTesting + switchCredConfUpdateError,
				Check:       resource.ComposeAggregateTestCheckFunc(),
				ExpectError: regexp.MustCompile(`.*ssh_private_key and key_pair_name must either both be set or neither be set.*.`),
			},
			// Update Resource Credential Success
			{
				Config: ProviderConfigForTesting + switchCredConfUpdateSuccess,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// Unit Tests
func TestAccResourceResourceCredentialvCenter(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a Unit test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Resource Credential invalid config error
			{
				Config:      ProviderConfigForTesting + VCenterCredentialUpdateError,
				Check:       resource.ComposeAggregateTestCheckFunc(),
				ExpectError: regexp.MustCompile(`.*domain must be set for type vcenter_credential*.`),
			},
			// Create Resource Credential
			{
				Config: ProviderConfigForTesting + VCenterCredentialCreate,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Update Resource Credential Success
			{
				Config: ProviderConfigForTesting + VCenterCredentialUpdateSuccess,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// Unit Tests
func TestAccResourceResourceCredentialPresentationServer(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a Unit test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{

			// Create Resource Credential
			{
				Config: ProviderConfigForTesting + PresentationServerCreate,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// Unit Tests
func TestAccResourceResourceCredentialPowerflexGateway(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a Unit test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create Resource Credential
			{
				Config: ProviderConfigForTesting + PowerflexGatewayCredentialCreate,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Resource Credential invalid config error
			{
				Config:      ProviderConfigForTesting + PowerflexGatewayUpdateError,
				Check:       resource.ComposeAggregateTestCheckFunc(),
				ExpectError: regexp.MustCompile(`.*os_username and os_password must be set for type powerflex_gateway*.`),
			},
			// Update Resource Credential Success
			{
				Config: ProviderConfigForTesting + PowerflexGatewayUpdateSuccess,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// Unit Tests
func TestAccResourceResourceCredentialOSUser(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a Unit test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create Resource Credential
			{
				Config: ProviderConfigForTesting + OsUserCredentialCreate,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Update Resource Credential Success
			{
				Config: ProviderConfigForTesting + OsUserUpdateSuccess,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// Unit Tests
func TestAccResourceResourceCredentialOSAdmin(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a Unit test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create Resource Credential
			{
				Config: ProviderConfigForTesting + OsAdminCredentialCreate,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Update Resource Credential Success
			{
				Config: ProviderConfigForTesting + OsAdminUpdateSuccess,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

// Unit Tests
func TestAccResourceResourceCredentialElementManager(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		t.Skip("Dont run with acceptance tests, this is a Unit test")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create Resource Credential
			{
				Config: ProviderConfigForTesting + ElementManagerCredentialCreate,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			// Update Resource Credential Success
			{
				Config: ProviderConfigForTesting + ElementManagerUpdateSuccess,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}
