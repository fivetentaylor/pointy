provider "aws" {
  region = "us-west-2"
}

provider "aws" {
  alias  = "dns_role"
  region = "us-west-2"

  assume_role {
    role_arn     = "arn:aws:iam::998899136269:role/StagingTerraformRole"
    session_name = "TerraformSessionRootAccount"
  }
}

terraform {
  backend "s3" {
    bucket         = "staging-reviso-terraform-state"
    key            = "postbuild/staging/terraform.tfstate"
    dynamodb_table = "staging-reviso-terraform-locks"
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
    bucket = "staging-reviso-terraform-state"
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

module "postbuild" {
  source = "../modules"

  providers = {
    aws.dns_role = aws.dns_role
  }

  server_sha        = var.server_sha
  web_sha           = var.web_sha
  slack_webhook_url = var.slack_webhook_url

  desired_web_count    = 1
  desired_server_count = 3

  web_domain         = "www.reviso.biz"
  app_domain         = "app.reviso.biz"
  route53_zone       = "Z0893527DM29OSHJP7NG"
  freeplay_env       = "staging"
  docs_bucket_name   = "stage-reviso-documents"
  images_bucket_name = "stage-reviso-images"
  dynamo_table_name  = "staging-reviso"
  env                = "prod"
  cookie_domain      = "reviso.biz"
  email_domain       = "reviso.biz"
  email_region       = "us-west-2"
  preview_prefix     = ""
  secret_arn         = "arn:aws:secretsmanager:us-west-2:533267310428:secret:staging-4YQM26"

  vpc_id                     = data.terraform_remote_state.prebuild.outputs.vpc_id
  app_security_group_id      = data.terraform_remote_state.prebuild.outputs.app_security_group_id
  internal_security_group_id = data.terraform_remote_state.prebuild.outputs.internal_security_group_id
}

output "app_host" {
  value = module.postbuild.app_host
}

output "ecs_deployment_task_definition" {
  value = module.postbuild.ecs_deployment_task_definition
}
