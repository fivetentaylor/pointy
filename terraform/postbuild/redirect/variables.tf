# Variables for domain configuration
variable "vpc_id" {
  description = "The VPC ID where the ALB is deployed"
  type        = string
}

variable "subnets" {
  description = "The subnets where the ALB is deployed"
  type        = list(string)
}

variable "app_security_group_id" {
  description = "The security group ID for the app servers"
  type        = string
}

variable "old_zone_id" {
  description = "The Route 53 zone ID for the old domain"
  type        = string
}


variable "old_root_domain" {
  description = "The root domain being redirected from (e.g., revi.so)"
  type        = string
}

variable "old_web_domain" {
  description = "The www domain being redirected from (e.g., www.revi.so)"
  type        = string
}

variable "old_app_domain" {
  description = "The app domain being redirected from (e.g., app.revi.so)"
  type        = string
}

variable "new_root_domain" {
  description = "The root domain being redirected to (e.g., pointy.ai)"
  type        = string
}

variable "new_web_domain" {
  description = "The www domain being redirected to (e.g., www.pointy.ai)"
  type        = string
}

variable "new_app_domain" {
  description = "The app domain being redirected to (e.g., app.pointy.ai)"
  type        = string
}

