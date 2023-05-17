variable "username" {
  type        = string
  description = "Stores the username of PowerFlex host."
  default = "admin"
}

variable "password" {
  type        = string
  description = "Stores the password of PowerFlex host."
  default = "Password123"
}

variable "endpoint" {
  type        = string
  description = "Stores the endpoint of PowerFlex host. eg: https://10.1.1.1:443, here 443 is port where API requests are getting accepted"
  default = "https://10.247.103.159"
}
