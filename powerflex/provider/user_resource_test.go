/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceUser(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create user Test
			{
				Config: ProviderConfigForTesting + UserResourceCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_user.user", "name", "NewUser"),
					resource.TestCheckResourceAttr("powerflex_user.user", "role", "Monitor"),
					resource.TestCheckResourceAttr("powerflex_user.user", "password", "Password123"),
				),
			},
			{
				ResourceName: "powerflex_user.user",
				ImportState:  true,
			},
			{
				ResourceName:  "powerflex_user.user",
				ImportState:   true,
				ImportStateId: "name:NewUser",
			},
			{
				ResourceName:  "powerflex_user.user",
				ImportState:   true,
				ImportStateId: "id:",
				ExpectError:   regexp.MustCompile("Empty import identifier"),
			},
			{
				ResourceName:  "powerflex_user.user",
				ImportState:   true,
				ImportStateId: "name:",
				ExpectError:   regexp.MustCompile("Empty import identifier"),
			},
			{
				ResourceName:  "powerflex_user.user",
				ImportState:   true,
				ImportStateId: "name:invalid",
				ExpectError:   regexp.MustCompile("Could not get user"),
			},
			{
				ResourceName:  "powerflex_user.user",
				ImportState:   true,
				ImportStateId: "id:invalid",
				ExpectError:   regexp.MustCompile("Could not get user"),
			},
			{
				ResourceName:  "powerflex_user.user",
				ImportState:   true,
				ImportStateId: "invalid",
				ExpectError:   regexp.MustCompile("Could not get user"),
			},
			{
				ResourceName:  "powerflex_user.user",
				ImportState:   true,
				ImportStateId: "lastName:dontCare",
				ExpectError:   regexp.MustCompile("Expected import identifier format"),
			},
			// Update user Test
			{
				Config: ProviderConfigForTesting + UserResourceUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_user.user", "name", "NewUser"),
					resource.TestCheckResourceAttr("powerflex_user.user", "role", "SystemAdmin"),
					resource.TestCheckResourceAttr("powerflex_user.user", "password", "Password123"),
				),
			},
		},
	})
}

func TestAccResourceUserCreateNegative(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Dont run with units tests because it will try to create the context")
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = Mock(helper.CreateSsoUser, OptGeneric).Return(nil, fmt.Errorf("Mock error")).Build()
				},
				Config:      ProviderConfigForTesting + UserResourceCreate2,
				ExpectError: regexp.MustCompile(`.*Error creating the user.*`),
			},
			{
				Config:      ProviderConfigForTesting + UserResourceCreate4,
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Value Match.*`),
			},
		},
	})
}

var UserResourceCreate = `
resource "powerflex_user" "user" {
	name = "NewUser"
	role = "Monitor"
	password = "Password123"
}
`

var UserResourceUpdate = `
resource "powerflex_user" "user" {
	name = "NewUser"
	role = "SystemAdmin"
	password = "Password123"
}
`

var UserResourceUpdateName = `
resource "powerflex_user" "user" {
	name = "NewUserRename"
	role = "Monitor"
	password = "Password123"
}
`

var UserResourceUpdatePassword = `
resource "powerflex_user" "user" {
	name = "NewUser"
	role = "Monitor"
	password = "Password123!"
}
`

var UserResourceCreate2 = `
resource "powerflex_user" "user" {
	name = "NewUser"
	role = "Monitor"
	password = "Password123"
}
resource "powerflex_user" "user2" {
	name = "NewUser"
	role = "Monitor"
	password = "Password123"
}
`

var UserResourceCreate3 = `
resource "powerflex_user" "user" {
	name = "NewUser"
	role = "Monitor"
	password = "Password123"
	first_name = "NewUser"
}
`

var UserResourceCreate4 = `
resource "powerflex_user" "user" {
	name = "NewUser"
	role = "NotRealRole"
	password = "Password123"
}
`
