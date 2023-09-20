/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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
# Create, Update, Read, and Delete operations are supported for this resource.
# For this resource, primary_mdm, secondary_mdm, tiebreaker_mdm and cluster_mode are mandatory parameters.
# While specifying primary_mdm, secondary_mdm and tiebreaker_mdm, id or ips is mandatory.
# For adding standby mdm, ips and role parameters are mandatory.
# To check which attributes of the MDM cluster resource can be updated, please refer Product Guide in the documentation

# Example for adding standby MDMs. Before adding standby MDMs, MDM package must be installed on VM with respective role. 
resource "powerflex_mdm_cluster" "test-mdm-cluster" {
  cluster_mode = "ThreeNodes"
  primary_mdm = {
    id = "7f328d0b71711801"
  }
  secondary_mdm = [{
    id = "10912d8a5d412800"
  }]
  tiebreaker_mdm = [{
    id = "0e4f0a2f5978ae02"
  }]
  standby_mdm = [
    {
      ips  = ["10.xxx.xx.xxx"]
      role = "Manager"
    },
    {
      ips  = ["10.yyy.yy.yyy"]
      role = "TieBreaker"
    },
  ]
}

# Example for changing MDM ownership. In this case, the id of the primary and secondary MDM will be swapped.
resource "powerflex_mdm_cluster" "test-mdm-cluster" {
  cluster_mode = "ThreeNodes"
  primary_mdm = {
    id = "10912d8a5d412800"
  }
  secondary_mdm = [{
    id = "7f328d0b71711801"
  }]
  tiebreaker_mdm = [{
    id = "0e4f0a2f5978ae02"
  }]
  standby_mdm = [
    {
      ips  = ["10.xxx.xx.xxx"]
      role = "Manager"
    },
    {
      ips  = ["10.yyy.yy.yyy"]
      role = "TieBreaker"
    },
  ]
}

# Example for switching cluster mode to 5 nodes from 3 nodes. The cluster mode will be FiveNodes. Previously added standby MDMs will be added as Secondary/TieBreaker MDM.
resource "powerflex_mdm_cluster" "test-mdm-cluster" {
  cluster_mode = "FiveNodes"
  primary_mdm = {
    id = "10912d8a5d412800"
  }
  secondary_mdm = [
    {
      id = "7f328d0b71711801"
    },
    {
      ips  = ["10.xxx.xx.xxx"]
      role = "Manager"
    },
  ]
  tiebreaker_mdm = [
    {
      id = "0e4f0a2f5978ae02"
    },
    {
      ips  = ["10.yyy.yy.yyy"]
      role = "TieBreaker"
    },
  ]
  standby_mdm = [
  ]
}

# Example for switching cluster mode to 3 nodes from 5 nodes. The cluster mode will be ThreeNodes. One of the active Secondary/TieBreaker MDM will be moved to standby MDMs.
resource "powerflex_mdm_cluster" "test-mdm-cluster" {
  cluster_mode = "ThreeNodes"
  primary_mdm = {
    id = "10912d8a5d412800"
  }
  secondary_mdm = [
    {
      id = "7f328d0b71711801"
    }
  ]
  tiebreaker_mdm = [
    {
      id = "0e4f0a2f5978ae02"
    }
  ]
  standby_mdm = [
    {
      ips  = ["10.xxx.xx.xxx"]
      role = "Manager"
    },
    {
      ips  = ["10.yyy.yy.yyy"]
      role = "TieBreaker"
    },
  ]
}

# Example for removing standby MDMs. In this case, standby_mdm will be empty list.
resource "powerflex_mdm_cluster" "test-mdm-cluster" {
  cluster_mode = "ThreeNodes"
  primary_mdm = {
    id = "10912d8a5d412800"
  }
  secondary_mdm = [
    {
      id = "7f328d0b71711801"
    }
  ]
  tiebreaker_mdm = [
    {
      id = "0e4f0a2f5978ae02"
    }
  ]
  standby_mdm = []
}