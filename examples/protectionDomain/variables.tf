variable "powerflex_username" {
  type        = string
  default     = "admin"
}

variable "powerflex_password" {
  type        = string
  default     = "1234"
  sensitive = true
}

variable "powerflex_host" {
  type        = string
  default     = "https://localhost:8000/"
}
