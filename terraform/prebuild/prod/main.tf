provider "aws" {
  region = "us-west-2"
}

terraform {
  backend "s3" {
    bucket         = "reviso-terraform-state"
    key            = "prebuild/terraform.tfstate"
    dynamodb_table = "reviso-terraform-locks"
    encrypt        = true
    region         = "us-west-2"
  }

  required_providers {
    external = {
      source  = "hashicorp/external"
      version = "~> 2.0"
    }
  }
}

module "prebuild" {
  source = "../modules"

  docs_bucket_name   = "reviso-documents"
  images_bucket_name = "reviso-images"
  dynamo_table_name  = "reviso"
  env_name                = "prod"
  email_domain       = "revi.so"
  email_region       = "us-east-1"

  secret_arn = "arn:aws:secretsmanager:us-west-2:998899136269:secret:production-QR5PVQ"

  redis_name          = "reviso-redis"
  redis_node_count    = 1
  redis_instance_type = "cache.t3.medium"

  create_db = false
}

output "vpc_id" {
  value = module.prebuild.vpc_id
}
output "internet_gateway_id" {
  value = module.prebuild.internet_gateway_id
}
output "route_table_id" {
  value = module.prebuild.route_table_id
}
output "repository_url" {
  value = module.prebuild.repository_url
}
output "app_security_group_id" {
  value = module.prebuild.app_security_group_id
}
output "internal_security_group_id" {
  value = module.prebuild.internal_security_group_id
}
