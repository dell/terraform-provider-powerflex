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

# commands to run this tf file : terraform init && terraform apply --auto-approve

# Get all node details present on the cluster
data "powerflex_node" "example1" {
}

// If multiple filter fields are provided then it will show the intersection of all of those fields.
// If there is no intersection between the filters then an empty datasource will be returned
// For more information about how we do our datasource filtering check out our guides: https://dell.github.io/terraform-docs/docs/storage/platforms/powerflex/product_guide/examples/ 
data "powerflex_node" "node" {
  # filter{
  #   ref_id  = ["id1", "id2"]
  #   state = ["READY"]
  #   ip_address = ["ip1", "ip2"]
  #   service_tag = ["serviceTag1", "serviceTag2"]
  #   current_ip_address = ["currentIp1", "currentIp2"]
  #   model = ["model1", "model2"]
  #   device_type = ["deviceType1", "deviceType2"]
  #   discover_device_type = ["discoverDeviceType1", "discoverDeviceType2"]
  #   display_name = ["displayName1", "displayName2"]
  #   managed_state = ["managedState1", "managedState2"]
  #   state = ["state1", "state2"]
  #   in_use = true
  #   custom_firmware = true
  #   needs_attention = true
  #   manufacturer = ["manufacturer1", "manufacturer2"]
  #   system_id = ["systemId1", "systemId2"]
  #   health = ["health1", "health2"]
  #   health_message = ["healthMessage1", "healthMessage2"]
  #   operating_system = ["operatingSystem1", "operatingSystem2"]
  #   number_of_cpus = [3, 2]
  #   nics = [3, 5]
  #   memory_in_gb = [3, 5]
  #   compliance_check_date = ["complianceCheckDate1", "complianceCheckDate2"]
  #   discovered_date = ["discoveredDate1", "discoveredDate2"]
  #   cred_id = ["credId1", "credId2"]
  #   compliance = ["compliance1", "compliance2"]
  #   failures_count = [3, 5]
  #   facts = ["facts1", "facts2"]
  #   puppet_cert_managed = ["puppetCertManaged1", "puppetCertManaged2"]
  #   flex_os_maint_mode=[0]
  #   esxi_maint_mode=[0]
  # }
}

output "node_result" {
  value = data.powerflex_node.example1.node_details
}
