# SDC End to End Example

## Overview
**This guide will give the steps in order to install and being to use a PowerFlex SDC.**

*It will give examples of the following:* 

1. Deploy the SDC Linux Virtual Machine (vSphere)
2. Install the SDC
3. Create and Attach a Volume to an SDCs
 
## Step 1: VM Deployment

*There are many way in which you can deploy a Virtual Machine. Below is an approch using the **vSphere** terraform provider. If your team use a different hypervisor, or has a differnt process for creating the VMs this step can just be skipped. The examples below show how to use the vSphere terraform provider but for more information you can check out there documentaion: https://registry.terraform.io/providers/hashicorp/vsphere/latest* 

> 1.1 Use the vSphere provider and Setup provider.tf 

### Provider Example:

*Fill in username, password, vsphere_server ip/hostname*

provider.tf
```
terraform {
  required_providers {
    vsphere = {
      source = "hashicorp/vsphere"
      version = "2.8.2"
    }
  }
}

provider "vsphere" {
  user                 = var.vsphere_user
  password             = var.vsphere_password
  vsphere_server       = var.vsphere_server
  allow_unverified_ssl = var.allow_unverified_ssl
  api_timeout          = 10
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

> 1.2 Gather Needed Data

*Using vSphere Terraform Datasources for datacenter, datastore, resource_pool, host and network, we need to gather all of the information needed to create the virtual machine*

Example:

datasources.tf
```
data "vsphere_datacenter" "datacenter" {
  name = var.vsphere_datacenter_name
}

data "vsphere_datastore" "datastore" {
  name = var.vsphere_datastore_name
  datacenter_id = data.vsphere_datacenter.datacenter.id
}

data "vsphere_resource_pool" "pool" {
  name = var.vsphere_resource_pool_name
  datacenter_id = data.vsphere_datacenter.datacenter.id
}

data "vsphere_host" "host" {
  name = var.vsphere_host_name
  datacenter_id = data.vsphere_datacenter.datacenter.id
}

data "vsphere_network" "network" {
  name = var.vsphere_network_name
  datacenter_id = data.vsphere_datacenter.datacenter.id
}


// Variables
variable "vsphere_datacenter_name" {
  type = string
  description = "The name of the datacenter in which to deploy the virtual machines."
}

variable "vsphere_datastore_name" {
  type = string
  description = "The name of the vSphere datastore in which to deploy the virtual machines."
}

variable "vsphere_network_name" {
  type = string
  description = "The name of the vSphere datastore in which to deploy the virtual machines."
}

variable "vsphere_resource_pool_name" {
  type = string
  description = "The name of the vSphere resource pool in which to deploy the virtual machines."
}

variable "vsphere_host_name" {
  type = string
  description = "The name of the vSphere host in which to deploy the virtual machines."
}
```

> 1.3 Deploy the Virtual Machine

*You can deploy one of two ways using a **.ova** image or using a vsphere **template**. Here is an examples of both.*

***Note**: VM should have the following packages **installed fio sshpass unzip yum-utils wget***

### .ova Example:

**Note: the *vapp* properties will need to be filled out based on each usecase as it will be different in each enviorment**
```
resource "vsphere_virtual_machine" "vm-installer" {
  name = var.vm_name
  datacenter_id = data.vsphere_datacenter.datacenter.id
  datastore_id         = data.vsphere_datastore.datastore.id
  host_system_id       = data.vsphere_host.host.id
  resource_pool_id = data.vsphere_resource_pool.pool.id

  num_cpus = var.num_cpus
  memory = var.memory
  wait_for_guest_net_timeout = 0
  wait_for_guest_ip_timeout  = 0

  network_interface {
    network_id = data.vsphere_network.network.id
    adapter_type = var.adapter_type 
  }

  ovf_deploy {
    allow_unverified_ssl_cert = var.allow_unverified_ssl
    remote_ovf_url = var.vm_ova_path
    disk_provisioning = var.disk_provisioning
    ip_protocol = var.ip_protocol
    ip_allocation_policy = var.ip_allocation_policy

    ovf_network_map = {
      network_id = data.vsphere_network.network.id
    }
  }

  // Depends on the image you need to fill out for your specific case/ova image you are using
  vapp {
    properties = {
      "guestinfo.hostname"     = "remote-foo.example.com",
      "guestinfo.ipaddress"    = "172.16.11.101",
      "guestinfo.netmask"      = "255.255.255.0",
      "guestinfo.gateway"      = "172.16.11.1",
      "guestinfo.dns"          = "172.16.11.4",
      "guestinfo.domain"       = "example.com",
      "guestinfo.ntp"          = "ntp.example.com",
      "guestinfo.password"     = "VMware1!",
      "guestinfo.ssh"          = "True"
    }
  }
}


// Variables

variable "vm_name" {
    type = string
    description = "Name of the PFMP installer VM, default to `new_terraform_vm`"
    default = "new_terraform_sdc_vm"
}

variable "disk_provisioning" {
  type = string
  description = "The disk provisioning type for the virtual machine. Options (thin, flat, think, sameAsSource) defaults to `thin`"
  default = "thin"
}

variable "ip_protocol" {
  type = string
  description = "The IP protocol for the virtual machine. Defaults to `IPV4`"
  default = "IPV4"
}

variable "ip_allocation_policy" {
  type = string
  description = "The IP allocation policy for the virtual machine. Defaults to `STATIC_MANUAL`"
  default = "STATIC_MANUAL"
}

variable "num_cpus" {
  type = number
  description = "Number of CPUs for the virtual machine. Defaults to `1`"
  default = 1
}

variable "memory" {
  type = number
  description = "Memory for the virtual machine. Defaults to `4060`"
  default = 4060
}

variable "adapter_type" {
  type = string
  description = "Network Adapter Type for the virtual machine. Defaults to `vmxnet3` Options can be found here: https://docs.vmware.com/en/VMware-vSphere/7.0/com.vmware.vsphere.vm_admin.doc/GUID-AF9E24A8-2CFA-447B-AC83-35D563119667.html"
  default = "vmxnet3"
}
```

### Template Example:

```
data "vsphere_virtual_machine" "template_example" {
  name          = "linux-template"
  datacenter_id = data.vsphere_datacenter.datacenter.id
}

resource "vsphere_virtual_machine" "template" {

  name             = var.vm_name
  resource_pool_id = data.vsphere_resource_pool.pool.id
  datastore_id     = data.vsphere_datastore.primary_datastore.id
  num_cpus         = var.num_cpus
  memory           = var.memory
  folder           = var.vmfolder
  guest_id         = data.vsphere_virtual_machine.template_example.guest_id
  scsi_type        = data.vsphere_virtual_machine.template_example.scsi_type
  scsi_controller_count = 4
  firmware         = "efi"

  network_interface {
    network_id   = data.vsphere_network.VM_Network.id
    adapter_type = data.vsphere_virtual_machine.template_example.network_interface_types[0]
  }

  wait_for_guest_net_timeout  = 10000000000000000
  wait_for_guest_net_routable = false
  wait_for_guest_ip_timeout   = 10000000000000000
  shutdown_wait_timeout       = 10
  migrate_wait_timeout        = 10000000000000000
  force_power_off             = false

  disk {
    label            = "OS"
    size             = 40
    thin_provisioned = data.vsphere_virtual_machine.template_example.disks.0.thin_provisioned
    unit_number      = 0
  }

  clone {
    template_uuid = data.vsphere_virtual_machine.template_example.id
    customize {
      linux_options {
        host_name = var.hostname
        domain    = var.dns[0]
      }
      network_interface {
        ipv4_address = var.mgmt_ip
        ipv4_netmask = 22
      }
      ipv4_gateway = var.gw
      dns_server_list = var.dns
      dns_suffix_list = var.search_domains
    }
  }


// Variables
variable "vm_name" {
    type = string
    description = "Name of the PFMP installer VM, default to `new_terraform_vm`"
    default = "new_terraform_sdc_vm"
}

variable "num_cpus" {
  type = number
  description = "Number of CPUs for the virtual machine. Defaults to `1`"
  default = 1
}

variable "memory" {
  type = number
  description = "Memory for the virtual machine. Defaults to `4060`"
  default = 4060
}

variable "search_domains" {
  type        = list(string)
  description = "List of DNS search domains"
}

variable "dns" {
  type = list(string)
}

variable "gw" {
  type = string
}

variable "vmfolder" {
  description = "Name of Folder to Deploy VM to"
}

variable "mgmt_ip" {
  description = "IP on Management Network"
  type        = string
}
```

## Step 2: SDC Install

*We support installing SDC on Windows EXSi and Linux Hosts. In this example we will use the linux sdc terraform module to do the install on the linux VM you have created above. For more examples including [EXSi](https://dell.github.io/terraform-docs/modules/storage/platforms/powerflex/v1.2.0/product_guide/vsphere_pfmp_installation/) and [Windows](https://dell.github.io/terraform-docs/modules/storage/platforms/powerflex/v1.2.0/product_guide/sdc_host_win/) look [here](https://dell.github.io/terraform-docs/modules/storage/platforms/powerflex/v1.2.0/product_guide).*

> 2.1 Download the SDC Packages

**Download the SDC Packages from [Dell Support](https://www.dell.com/support/product-details/en-us/product/scaleio/drivers). Once downloaded add the sdc package to an http or ftp site for easy access when deploying sdcs**

> 2.2 Install the SDC

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

### Example Module:

```
module "modules_example_sdc_host_linux" {
  source  = "dell/modules/powerflex//examples/sdc_host_linux"
  version = "1.2.0"

  remote_host = {
    user = var.remote_user
    password = var.remote_password
    // If you have a ssh key and certificate.
    // You are able to use those instead of user/pass
    //private_key = ""
    //certificate = ""
  }

  ip = var.ip

  versions = versions={
    pflex =  var.sdc_pflex_version  //"4.5.3000.118"
    kernel = var.linux_kernal_version //"5.15.0-1-generic"
  }
  
  scini = {
    linux_distro = "RHEL9" #"Ubuntu"
    autobuild_scini = true
    // If autobuild is equal to false, give the path to your scini.ko or scini.rpm else this is not needed
    //url = "http://example.com/release/5.15.0-1-generic"
  }

  sdc_pkg = {
    url = var.package_url
    // Do not change this value
    remote_pkg_name = var.remote_pkg_name
    remote_file = var.package_name
    use_remote_path = true
    skip_download_sdc = true
  }
  mdm_ips = var.mdm_ips
}


// Variables
variable "ip" {
  type        = string
  description = "Stores the IP address of the remote Linux host."
}

variable "package_url" {
  type        = string
  description = "URL to where the SDC Package is located for download. This package will be downloaded directly onto the virtual machine. Example: http://example.com/EMC-ScaleIO-sdc-3.6-700.103.Ubuntu.22.04.x86_64.tar, ftp://username:password@ftpserver/path/to/file"
}

variable "remote_pkg_name" {
  type        = string
  description = "The name of the package. Example: emc-sdc-package.tar (for Ubuntu) emc-sdc-package.rpm (for RHEL)"
  default = "emc-sdc-package.tar"
}

variable "package_name" {
  type        = string
  description = "The name of the package. Example: EMC-ScaleIO-sdc-4.5-3000.118.Ubuntu.22.04.x86_64.tar (for Ubuntu) EMC-ScaleIO-sdc-4.5-3000.118.Ubuntu.22.04.x86_64.rpm (for RHEL)"
}

variable "remote_user" {
  type        = string
  description = "Stores the username of the remote Linux host."
}

variable "remote_password" {
  type        = string
  description = "Stores the password of the remote Linux host."
}

variable "sdc_pflex_version" {
  type        = string
  description = "The version of the powerflex sdc you are installing, this can be grabbed from the sdc package name i.e: PowerFlex_4.5.2100.105_Ubuntu18.04_SDC.tar the version would be 4.5.2100.105"
}

variable "linux_kernal_version" {
  type        = string
  description = "The version of linux we are using"
}

variable "mdm_ips" {
  description = "all the mdms (either primary,secondary or virtual ips) in a comma separated list by cluster. If sdc is only being connected to one powerflex leave this unset and it will use the mdms of the cluster set of the provider by default"
  type        = list(string)
  default = []
}
```

## Step 3: Create and Attach Volume

> 3.1 Fill in the provider config

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

> 3.2 Create the Volume

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

> 3.3 Attach the Volume to the SDC

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