# ESXi SDC End to End Example

## Overview
**This guide will give the steps in order to install and being to use a PowerFlex SDC in the vSphere context.**

*It will give examples of the following:* 

1. Install the SDC on an ESXi Host
2. Create and Attach a Volume to an SDCs
3. Create VMFS Datastoures based on the Attached volumes

## Step 1: Install the SDC on ESXi Host

> 1.1 Setup Provider.tf

*Fill in username, password, and endpoint ip/hostname*

provider.tf
```
terraform {
   required_providers {
     powerflex = {
       version = ">=1.6.0"
       source  = "registry.terraform.io/dell/powerflex"
     }
   }
 }


provider "powerflex" {
  username = var.username
  password = var.password
  endpoint = var.endpoint
  insecure = true
  timeout  = 120
}


// Variables
variable "username" {
  type        = string
  description = "Stores the username of PowerFlex host."
}

variable "password" {
  type        = string
  description = "Stores the password of PowerFlex host."
}

variable "endpoint" {
  type        = string
  description = "Stores the endpoint of PowerFlex host. eg: https://10.1.1.1:443, here 443 is port where API requests are getting accepted"
}
```
> 1.2 Download the SDC Packages

**Download the SDC Packages from [Dell Support](https://www.dell.com/support/product-details/en-us/product/scaleio/drivers).**

*We support installing SDC on Windows ESXi and Linux Hosts. In this example we will use the ESXi sdc_host terraform resource to do the install of an SDC an ESXi host. For more examples check out (https://dell.github.io/terraform-docs/docs/storage/platforms/powerflex/product_guide/resources/sdc_host/).*

> 1.3 Install the SDC on ESXi Host

```

# # generate a random guid. This is required only for ESXi hosts.
resource "random_uuid" "sdc_guid" {
}

resource "powerflex_sdc_host" "sdc" {
  ip = "10.10.10.10" // The ESXi Host IP
  remote = {
    user = var.esxi-host-username
    password = var.esxi-host-password
  }
  os_family = "esxi"
  esxi = {
    //We generate one using terraform `random_guid` resource above but it should look something like this: "12345678-90AB-CDEF-1234-567890ABCDEF"
    guid = random_uuid.sdc_guid.result
  }
  name         = var.sdc-name
  package_path = var.sdc-path
  // Always just set this to false
  use_remote_path = false
}

variable "sdc-name" {
  type        = string
  description = "Stores the name of the SDC."
  default = "current-user-sdc-esxi"
}

variable "sdc-path" {
  type        = string
  description = "Stores the path to where you downloaded the sdc ESXi package if in this directory then should looks something like './sdc-4.5.0.263-esx8.x.zip'"
}

variable "esxi-host-username" {
  type        = string
  description = "Stores the username of ESXi host."
}

variable "esxi-host-password" {
  type        = string
  description = "Stores the password of the ESXi host."
}
```

## Step 2: Create and Attach Volume

> 2.1 Fill in the provider config

*Fill in username, password, and endpoint ip/hostname*

provider.tf
```
terraform {
   required_providers {
     powerflex = {
       version = ">=1.6.0"
       source  = "registry.terraform.io/dell/powerflex"
     }
   }
 }


provider "powerflex" {
  username = var.username
  password = var.password
  endpoint = var.endpoint
  insecure = true
  timeout  = 120
}


// Variables
variable "username" {
  type        = string
  description = "Stores the username of PowerFlex host."
}

variable "password" {
  type        = string
  description = "Stores the password of PowerFlex host."
}

variable "endpoint" {
  type        = string
  description = "Stores the endpoint of PowerFlex host. eg: https://10.1.1.1:443, here 443 is port where API requests are getting accepted"
}
```

> 2.2 Create the Volume

```
resource "powerflex_volume" "volume-create" {
  name = var.name

  # To create / update, either protection_domain_id or protection_domain_name must be provided
  protection_domain_name = var.protection_domain_name

  # To create / update, either storage_pool_id or storage_pool_name must be provided
  storage_pool_name = var.storage_pool_name

  # The unit of size of the volume is defined by capacity_unit whose default value is "GB".
  size          = var.size
  capacity_unit = var.capacity_unit # GB/TB

  use_rm_cache = var.use_rm_cache # true/false
  volume_type  = var.volume_type # ThickProvisioned/ThinProvisioned volume type
  access_mode  = var.access_mode # ReadWrite/ReadOnly volume access mode
  remove_mode  = var.remove_mode # INCLUDING_DESCENDANTS/ONLY_ME remove mode
}

// Variables

variable "name" {
  type        = string
  description = "name of the Powerflex Volume"
}

variable "protection_domain_name" {
  type        = string
  description = "protection_domain name of the Powerflex Volume"
}

variable "storage_pool_name" {
  type        = string
  description = "storage_pool name of the Powerflex Volume"
}

variable "size" {
  type        = number
  description = "size of the Powerflex Volume"
  default = 8
}

variable "capacity_unit" {
  type        = string
  description = "Capacity Unit of the Powerflex Volume Options:(GB/TB)"
  default = "GB"
}

variable "use_rm_cache" {
  type        = bool
  description = "Sets the rm_cache option of the Powerflex Volume"
  default = true
}

variable "volume_type" {
  type        = string
  description = "Volume Type of the Powerflex Volume Options:(ThickProvisioned/ThinProvisioned)"
  default = "ThickProvisioned"
}

variable "access_mode" {
  type        = string
  description = "Access Mode of the Powerflex Volume Options:(ReadWrite/ReadOnly)"
  default = "ReadWrite"
}

variable "remove_mode" {
  type        = string
  description = "Remove Mode of the Powerflex Volume Options:(INCLUDING_DESCENDANTS/ONLY_ME)"
  default = "INCLUDING_DESCENDANTS"
}
```

> 2.3 Attach the Volume to the SDC

*Grab the SDC ID using the sdc datasource, then map the newly created volume to the SDC* 

```
data "powerflex_sdc" "filtered" {
  filter {
    name = [var.sdc_name]
  }
}

resource "powerflex_sdc_volumes_mapping" "mapping-test" {
  # SDC id
  id = data.powerflex_sdc.filtered.sdcs[0].id
  volume_list = [
    {
      # id of the volume which needs to be mapped. 
      # either volume_id or volume_name can be used.
      volume_id = resource.powerflex_volume.volume-create.id

      # Valid values are 0 or integers greater than 10
      limit_iops = var.limit_iops

      # Default value is 0
      limit_bw_in_mbps = var.limit_bw_in_mbps

      access_mode = var.access_mode # ReadOnly/ReadWrite/NoAccess
    }
  ]
}


// Variables
variable "sdc_name" {
  type        = string
  description = "name of the sdc"
}

variable "limit_iops" {
  type        = number
  description = "limit_iops Valid values are 0 or integers greater than 10"
  default = 140
}

variable "limit_bw_in_mbps" {
  type        = number
  description = "limit_bw_in_mbps Default value is 0"
  default = 0
}

variable "access_mode" {
  type        = string
  description = "access_mode Options(ReadOnly/ReadWrite/NoAccess)"
  default = "ReadOnly"
}
```

## Step 3: Create VMFS Datasource

> 3.1 Setup the vSphere provider

```
terraform {
  required_providers {
    vsphere = {
      source = "hashicorp/vsphere"
      version = "2.11.1"
    }
  }
}

provider "vsphere" {
  user                 = var.vsphere_user
  password             = var.vsphere_password
  vsphere_server       = var.vsphere_server
  allow_unverified_ssl = var.allow_unverified_ssl
  api_timeout          = 100
}

// Variables
variable "vsphere_user" {
  type        = string
  description = "Stores the username of vsphere_user."
}

variable "vsphere_password" {
  type        = string
  description = "Stores the password of vsphere_password."
}

variable "vsphere_server" {
  type        = string
  description = "Stores the host ip/fqdn of the vsphere_server."
}

variable "allow_unverified_ssl" {
  type        = string
  description = "Allow unverified ssl connection"
  default = true
}
```

> 3.2 Rescan and Create the VMFS Datastore

```
data "vsphere_datacenter" "datacenter" {
  name = "Datacenter"
}

data "vsphere_host" "host" {
  name          = var.host_name // Hosts name in vSphere 
  datacenter_id = data.vsphere_datacenter.datacenter.id
}

data "vsphere_vmfs_disks" "available" {
  host_system_id = data.vsphere_host.host.id
  rescan         = true
  filter = var.vmfs_filter // Should look something like eui.bcc545a72b516e0fa45627070000000e
}

resource "vsphere_vmfs_datastore" "datastore" {
  name           = var.vmfs_datastore_name // name of vmfs datastore to be created
  host_system_id = data.vsphere_host.host.id

  disks = "${data.vsphere_vmfs_disks.available.disks}"
}

// Variables
variable "host_name" {
  type        = string
  description = "Stores the name of the vSphere host."
}

variable "vmfs_filter" {
  type        = string
  description = "The filter which is used to grab the disk name to create the VMFS datasource upon. Starts with `eui.`. If you want a specific disk give full eui value, example: `eui.bcc545a72b516e0fa45627070000000e`. Otherwise the default will grab all disks with the eui path on the Powerflex Storage Adaptor"
  default = "eui."
}

variable "vmfs_datastore_name" {
  type        = string
  description = "The name used when creating the VMFS datastore"
  default = "terraform-vmfs-datastore"
}

```