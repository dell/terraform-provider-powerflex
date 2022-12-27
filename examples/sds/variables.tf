variable "powerflex_username" {
  type        = string
  default     = "user"
}

variable "powerflex_password" {
  type        = string
  default     = "pass"
  sensitive = true
}

variable "powerflex_endpoint" {
  type        = string
  default     = "https://localhost:8000/"
}
