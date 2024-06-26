---
page_title: "Deploying a PowerFlex cluster using ssh key"
title: "Deploying a PowerFlex cluster using ssh key"
linkTitle: "Deploying a PowerFlex cluster using ssh key"
---

<!--
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
-->

# Deploying a PowerFlex cluster using ssh key for passwordless authentication

This guide explains how to use ssh keys for deploying a cluster using passwordless authentication

### Steps for configuring ssh key

- On one of the ScaleIO cluster machine:
    - Go to ~/.ssh
    - ssh-keygen -t rsa
- cat id_rsa.pub > authorized_keys

- If the RSA key was created in OpenSSH key format convert to pem format:
    - ssh-keygen -p -m PEM -f ~/.ssh/id_rsa
- Copy the id_rsa, id_rsa.pub and authorized_keys to ~/.ssh on the rest of the cluster machines
- Replace the key text in block-legacy-gateway-sshkey.yaml (attached below) with the content of id_rsa
- Delete the secret before setting it:
    - kubectl delete -n powerflex secret platform-nodes-credentials
- run on the M&O: 
    - kubectl apply -n powerflex -f <path of block-legacy-gateway-sshkey.yaml>
- delete block-legacy-gateway pod 
    - legacy_gw=$(kubectl get pods -n powerflex|grep -i legacy-gateway |grep -v mds |awk '{print $1}'); kubectl delete pods -n powerflex $legacy_gw
- Remove the value (leave empty) from all the "password" fields in the csv file

### Example of block-legacy-gateway-sshkey.yaml file content

apiVersion: v1
kind: Secret
metadata:
  name: platform-nodes-credentials
type: Opaque
stringData:
  id_rsa: |
    -----BEGIN RSA PRIVATE KEY-----
    PASTE THE RSA KEY HERE
    -----END RSA PRIVATE KEY-----

### Example of cluster resource in case of passwordless authentication

In this case just remove the password attribute from the cluster resource .

```terraform

resource "powerflex_cluster" "test" {
  mdm_password = "Password"
  lia_password = "Password"

  # Advance Security Configuration
  allow_non_secure_communication_with_lia = false
  allow_non_secure_communication_with_mdm = false
  disable_non_mgmt_components_auth        = false

  # Cluster Configuration related fields
  cluster = [
    {
      # MDM Configuration Fields
      username                 = "user",
      operating_system         = "linux",
      is_mdm_or_tb             = "Primary",
      mdm_ips                  = "10.10.10.1",
      mdm_mgmt_ip              = "10.10.10.1",
      mdm_name                 = "mdm1",
      perf_profile_for_mdm     = "HighPerformance",
      is_sds                   = "Yes",
      sds_name                 = "sds1",
      sds_all_ips              = "10.10.10.1",
      protection_domain        = "domain_1",
      sds_storage_device_list  = "/dev/sdb",
      sds_storage_device_names = "sdb",
      storage_pool_list        = "pool1",
      perf_profile_for_sds     = "HighPerformance",
      is_sdc                   = "Yes",
      sdc_name                 = "sdc1",
      perf_profile_for_sdc     = "HighPerformance"

    },
    {
      username                 = "user",
      password                 = "password",
      operating_system         = "linux",
      is_mdm_or_tb             = "Secondary",
      mdm_name                 = "mdm2",
      mdm_ips                  = "10.10.10.2",
      mdm_mgmt_ip              = "10.10.10.2",
      perf_profile_for_mdm     = "HighPerformance",
      is_sds                   = "Yes",
      sds_name                 = "sds2",
      sds_all_ips              = "10.10.10.2",
      sds_storage_device_list  = "/dev/sdb",
      sds_storage_device_names = "sdb",
      protection_domain        = "domain_1",
      storage_pool_list        = "pool1",
      perf_profile_for_sds     = "HighPerformance",
      is_sdc                   = "Yes",
      sdc_name                 = "sdc2",
      perf_profile_for_sdc     = "HighPerformance"
    },
    {
      username                 = "user",
      password                 = "password",
      operating_system         = "linux",
      is_mdm_or_tb             = "TB",
      mdm_name                 = "tb1",
      mdm_ips                  = "10.10.10.3",
      mdm_mgmt_ip              = "10.10.10.3",
      perf_profile_for_mdm     = "HighPerformance",
      is_sds                   = "Yes",
      sds_name                 = "sds3",
      sds_all_ips              = "10.10.10.3",
      sds_storage_device_list  = "/dev/sdb",
      sds_storage_device_names = "sdb",
      protection_domain        = "domain_1",
      storage_pool_list        = "pool1",
      perf_profile_for_sds     = "HighPerformance",
      is_sdc                   = "Yes",
      sdc_name                 = "sdc3",
      perf_profile_for_sdc     = "HighPerformance"
    }
  ]

    storage_pools = [
    {
      media_type        = "SSD"
      protection_domain = "domain_1"
      storage_pool      = "pool1"
      daya_layout       = "MG"
      zero_padding      = "true"
    }
  ]
}

```
This Terraform configuration sets up a PowerFlex cluster without using password.