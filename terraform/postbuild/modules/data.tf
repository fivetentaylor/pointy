data "aws_caller_identity" "current" {}
data "aws_partition" "current" {}
data "aws_region" "current" {}

data "aws_vpc" "main" {
  id = var.vpc_id
}

data "aws_subnets" "public" {
  filter {
    name   = "tag:Purpose"
    values = ["Public"]
  }
}

data "aws_security_group" "app_security_group" {
  id = var.app_security_group_id
}

data "aws_security_group" "internal_security_group" {
  id = var.internal_security_group_id
}
