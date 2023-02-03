variable "username" {
  type = string
  description = "Stores the username of PowerFlex host."
  default = "admin"
}

variable "password" {
  type = string
  description = "Stores the password of PowerFlex host."
  default = "Password123"
}

variable "endpoint" {
    type = string
    description = "Stores the endpoint of PowerFlex host"
    default = "https://10.247.101.69"
}
