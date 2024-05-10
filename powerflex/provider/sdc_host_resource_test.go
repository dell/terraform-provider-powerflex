/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccSDCResource tests the SDC Expansion Operation
func TestAccSDCHostResource(t *testing.T) {
	os.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Create
			{
				Config:      ProviderConfigForTesting + SDCHostConfig1,
				ExpectError: regexp.MustCompile(`.*Error During Installation.*`),
			},
		},
	})
}

var SDCHostConfig1 = `
	resource powerflex_sdc_host sdc {
	ip = "10.247.103.163"
	remote = {
	  user = "root"
	  # we are not using password auth here, but it can be used as well
	  password = "dangerous"
	  #private_key = data.local_sensitive_file.ssh_key.content_base64
	  #host_key = data.local_sensitive_file.host_key.content_base64
	}
	os_family = "linux"
	name = "sdc-linux-rhel"
	package_path = "/root/powerflex_packages/PowerFlex_3.6.700.103_RHEL_OEL7/EMC-ScaleIO-sdc-3.6-700.103.el7.x86_64.rpm"
	#mdm_ips = ["10.10.10.5", "10.10.10.6"]
  }

`
