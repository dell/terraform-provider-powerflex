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

# Command to run this tf file : terraform init && terraform plan && terraform apply.

// Resource to manage lifecycle
resource "terraform_data" "always_run" {
  input = timestamp()
}

# Example for setting compatibility management. After successful execution, device will be added to the specified storage pool
resource "powerflex_compatibility_management" "test" {
  # Required Path on your local machine to your gpg file ie(/example/path/secring.gpg)
  repository_path = "/example/path/example.gpg"

  // This will allow terraform create process to trigger each time we run terraform apply.
  lifecycle {
    replace_triggered_by = [
      terraform_data.always_run
    ]
  }
}
