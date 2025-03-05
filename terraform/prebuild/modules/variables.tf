data "aws_caller_identity" "current" {}
data "aws_partition" "current" {}
data "aws_region" "current" {}

locals {
  account_id = data.aws_caller_identity.current.account_id
  region     = data.aws_region.current.name
  partition  = data.aws_partition.current.partition
  task_count = 1
  family     = "reviso"
}

variable "dynamo_table_name" {
  type = string
}

variable "email_domain" {
  type = string
}

variable "email_region" {
  type = string
}

variable "env_name" {
  default = "prod"
}

variable "docs_bucket_name" {
  description = "Name of the S3 bucket for reviso docs"
  type        = string
}

variable "images_bucket_name" {
  description = "Name of the S3 bucket for reviso images"
  type        = string
}

variable "create_db" {
  type    = bool
  default = true # leave as true otherwise terraform tries to destroy the db
}

variable "secret_arn" {
  type = string
}

variable "redis_name" {
  description = "The name of the Redis cluster"
  type        = string
}

variable "redis_instance_type" {
  description = "The instance type of the Redis node"
  type        = string
}

variable "redis_node_count" {
  description = "The number of nodes in the Redis cluster"
  type        = number
}
