output "access_key" {
  value = aws_iam_access_key.terraform_user_key.id
}

output "secret_key" {
  value     = aws_iam_access_key.terraform_user_key.secret
  sensitive = true
}

output "s3_bucket_name" {
  description = "The name of the S3 bucket."
  value       = aws_s3_bucket.terraform_state.bucket
}

output "dynamodb_table_name" {
  description = "The name of the DynamoDB table."
  value       = aws_dynamodb_table.terraform_locks.name
}

