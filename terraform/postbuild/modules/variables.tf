locals {
  account_id = data.aws_caller_identity.current.account_id
  region     = data.aws_region.current.name
  partition  = data.aws_partition.current.partition
  task_count = 1
}

variable "name" {
  type = string
}

variable "web_domain" {
  type = string
}

variable "app_domain" {
  type = string
}

variable "route53_zone" {
  type = string
}

variable "freeplay_env" {
  type    = string
  default = "production"
}

variable "preview_prefix" {
  type    = string
  default = ""
}

variable "email_domain" {
  type = string
}

variable "email_region" {
  type    = string
  default = "us-east-1"
}

variable "vpc_id" {
  type = string
}

variable "env" {
  default = "prod"
}

variable "server_sha" {
  type = string
}

variable "web_sha" {
  type = string
}

variable "secret_arn" {
  type = string
}

variable "desired_web_count" {
  type = number
}

variable "desired_server_count" {
  type = number
}

variable "dynamo_table_name" {
  type = string
}

variable "docs_bucket_name" {
  description = "Name of the S3 bucket for reviso docs"
  type        = string
}

variable "images_bucket_name" {
  description = "Name of the S3 bucket for reviso images"
  type        = string
}

variable "cookie_domain" {
  type = string
}

variable "app_security_group_id" {
  type = string
}

variable "internal_security_group_id" {
  type = string
}

variable "slack_webhook_url" {
  type = string
}
