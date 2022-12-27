variable "username" {
  type = string
  description = "Stores the username of PowerFlex host."
}

variable "password" {
  type = string
  description = "Stores the password of PowerFlex host."
}

variable "endpoint" {
    type = string
    description = "Stores the endpoint of PowerFlex host"
}
variable "username" {
  type = string
  description = "Stores the username of PowerFlex host."
}

variable "password" {
  type = string
  description = "Stores the password of PowerFlex host."
}

variable "endpoint" {
    type = string
    description = "Stores the endpoint of PowerFlex host"
}

variable "volume_resource_name" {
  type = string
  description = "Stores the name of volume"
}

variable "volume_resource_capacity_unit" {
    type = string 
    description = "Stores the capacity unit of volume"
}

variable "volume_resource_storage_pool_id" {
  type= string
  description = "Stores storage pool id "
}

variable "volume_resource_protection_domain_id" {
  type=string 
  description = "Stores protection domain id"
}

variable "volume_resource_size" {
  type = number
  description = "Stores the size of volume"
}
