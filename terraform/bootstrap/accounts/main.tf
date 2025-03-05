resource "aws_s3_bucket" "terraform_state" {
  bucket = "${var.environment}-reviso-terraform-state"
}

variable "environment" {
  type    = string
  default = "staging"
}

resource "aws_dynamodb_table" "terraform_locks" {
  name         = "${var.environment}-reviso-terraform-locks"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }
}

resource "aws_iam_role" "terraform_role" {
  name = "TerraformRole"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          AWS = "arn:aws:iam::998899136269:user/TerraformUser"
        }
      },
    ]
  })
}

resource "aws_iam_role_policy" "terraform_role_policy" {
  name = "TerraformRolePolicy"
  role = aws_iam_role.terraform_role.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action   = "*",
        Effect   = "Allow",
        Resource = "*",
      },
    ],
  })
}
