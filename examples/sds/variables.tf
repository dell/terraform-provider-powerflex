variable "powerflex_username" {
  type        = string
  default     = "admin"
}

variable "powerflex_password" {
  type        = string
  default     = "Password123"
  sensitive = true
}

variable "powerflex_endpoint" {
  type        = string
  default     = "https://10.247.101.69:/"
}
