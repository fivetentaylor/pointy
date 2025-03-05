resource "aws_s3_bucket" "terraform_state" {
  bucket = "reviso-terraform-state"
}

resource "aws_dynamodb_table" "terraform_locks" {
  name         = "reviso-terraform-locks"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }
}

resource "aws_iam_user" "terraform_user" {
  name = "TerraformUser"
}

resource "aws_iam_user_policy" "terraform_user_inline_policy" {
  name = "TerraformUserInlinePolicy"
  user = aws_iam_user.terraform_user.name
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect   = "Allow"
        Action   = "*"
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_access_key" "terraform_user_key" {
  user = aws_iam_user.terraform_user.name
}

