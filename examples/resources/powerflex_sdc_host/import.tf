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

# load the private key
data "local_sensitive_file" "ssh_key" {
    filename = "/root/.ssh/linux_rsa"
}

# load the host key
data "local_sensitive_file" "host_key" {
    filename = "linux_host_ecdsa_key.pub"
}

import {
    to = powerflex_sdc_host.import_test_sdc_linux
    id = data.local_sensitive_file.host_key.content[0].id
}

# Example for adding an Linux host as SDC.
resource "powerflex_sdc_host" "import_test_sdc_linux" {
    ip = var.ip
    remote = {
        user = var.user
        # we are not using password auth here, but it can be used as well
        # password = var.password
        private_key = data.local_sensitive_file.ssh_key.content_base64
        host_key    = data.local_sensitive_file.host_key.content_base64
    }
    os_family       = var.os_family
    name            = var.name
    use_remote_path = false
    package_path    =  var.package_path
}