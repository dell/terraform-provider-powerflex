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

import {
    to = powerflex_sdc_host.import_test_sdc_linux
    id = var.ip
}

# Example for adding an Linux host as SDC.
resource "powerflex_sdc_host" "import_test_sdc_linux" {
    ip = var.ip
    remote = {
      user = var.user
      password = var.password
      port = var.port
    }
    os_family       = var.os_family
    name            = "sdc-host-test-rename"
    performance_profile = "Compact"
    package_path = var.package_path
}