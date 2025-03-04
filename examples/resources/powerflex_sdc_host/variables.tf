variable "user" {
  type        = string
  description = "username of the sdc host."
  default = ""
}

variable "password" {
  type        = string
  description = "password of the sdc host."
  default = ""
}

variable "port" {
  type        = string
  description = "port of the sdc host."
  default = ""
}

variable "name" {
  type        = string
  description = "name of the sdc host."
  default = ""
}

variable "os_family" {
  type        = string
  description = "OS family of the sdc host."
  default = ""
}
variable "package_path" {
  type        = string
  description = "package path of the sdc host."
  default = ""
}

variable "ip" {
  type        = string
  description = "ip addresse(s) of the sdc host."
  default = ""
}