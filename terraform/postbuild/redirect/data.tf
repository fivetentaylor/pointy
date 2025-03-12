data "aws_caller_identity" "current" {}
data "aws_partition" "current" {}
data "aws_region" "current" {}

# Access the existing Route53 zone for the old domain
data "aws_route53_zone" "old_domain_zone" {
  provider = aws.dns_role
  name     = "${var.old_root_domain}."  # Make sure to include the trailing dot
}

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
