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

# Example for fetching details of the firmware repository using names
data "powerflex_firmware_repository" "test" {
firmware_repository_names = ["PowerFlex 4.5.1.0", "PowerFlex 4.5.2.0"]
}

# Example for fetching details of the firmware repository using id
data "powerflex_firmware_repository" "test" {
	firmware_repository_ids = ["8aaa3fda8f5c2609018f854266e12865", "8aaa3fda8f5c2609018f857b6c0d2ede"]
	}

# Example for fetching all the firmware repository details
data "powerflex_firmware_repository" "test" {
}

