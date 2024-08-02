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

# Command to run this tf file : terraform plan && terraform apply.

// Resource to manage lifecycle
resource "terraform_data" "always_run" {
  input = timestamp()
}

# Example for creating OS repository. After successful execution, OS Repository will be created.
resource "powerflex_os_repository" "test" {
   # Required Fields

   # Name of the OS repository
   name = "TestOsRepo"
   # Source path of the OS image
   source_path = "https://pathtoimage.iso"
   # Supported image types are redhat7, vmware_esxi, sles, windows2016, windows2019
   image_type = "vmware_esxi"
        
    // This will allow terraform create process to trigger each time we run terraform apply.
    lifecycle {
        replace_triggered_by = [
         terraform_data.always_run
        ]
    }
}
