provider "aws" {
  region = "us-west-2"
}

provider "aws" {
  alias  = "dns_role"
  region = "us-west-2"
}

terraform {
  backend "s3" {
    bucket         = "reviso-terraform-state"
    key            = "postbuild/terraform.tfstate"
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

data "terraform_remote_state" "prebuild" {
  backend = "s3"
  config = {
    bucket = "reviso-terraform-state"
    key    = "prebuild/terraform.tfstate"
    region = "us-west-2"
  }
}

variable "server_sha" {
  type = string
}

variable "web_sha" {
  type = string
}

variable "slack_webhook_url" {
  type = string
}

module "redirect_reviso" {
  source = "../redirect"

  providers = {
    aws.dns_role = aws.dns_role
  }

  old_root_domain = "revi.so"
  old_web_domain  = "www.revi.so"
  old_app_domain  = "app.revi.so"

  new_root_domain = "pointy.ai"
  new_web_domain  = "www.pointy.ai"
  new_app_domain  = "app.pointy.ai"

  vpc_id                = data.terraform_remote_state.prebuild.outputs.vpc_id
  app_security_group_id = data.terraform_remote_state.prebuild.outputs.app_security_group_id
  old_zone_id           = "Z05640301EQTNY1VO8UTF"
}

module "postbuild_pointy" {
  source = "../modules"

  providers = {
    aws.dns_role = aws.dns_role
  }

  server_sha        = var.server_sha
  web_sha           = var.web_sha
  slack_webhook_url = var.slack_webhook_url

  desired_web_count    = 2
  desired_server_count = 3

  name                       = "pointy"
  web_domain                 = "www.pointy.ai"
  app_domain                 = "app.pointy.ai"
  route53_zone               = "Z036564521TUMO3XKOL1V"
  freeplay_env               = "production"
  docs_bucket_name           = "reviso-documents"
  images_bucket_name         = "reviso-images"
  dynamo_table_name          = "reviso"
  env                        = "prod"
  cookie_domain              = "pointy.ai"
  email_domain               = "pointy.ai"
  email_region               = "us-east-1"
  preview_prefix             = ""
  secret_arn                 = "arn:aws:secretsmanager:us-west-2:998899136269:secret:production-QR5PVQ"
  vpc_id                     = data.terraform_remote_state.prebuild.outputs.vpc_id
  app_security_group_id      = data.terraform_remote_state.prebuild.outputs.app_security_group_id
  internal_security_group_id = data.terraform_remote_state.prebuild.outputs.internal_security_group_id
}

output "app_host" {
  value = module.postbuild_pointy.app_host
}

output "ecs_deployment_task_definition" {
  value = module.postbuild_pointy.ecs_deployment_task_definition
}
