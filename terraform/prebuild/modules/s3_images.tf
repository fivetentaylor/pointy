resource "aws_s3_bucket" "images_bucket" {
  bucket = var.images_bucket_name
}

# Configure bucket ownership controls
resource "aws_s3_bucket_ownership_controls" "images_bucket_ownership" {
  bucket = aws_s3_bucket.images_bucket.id

  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

# Configure public access block
resource "aws_s3_bucket_public_access_block" "images_bucket_public_access" {
  bucket = aws_s3_bucket.images_bucket.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# Configure bucket ACL
resource "aws_s3_bucket_acl" "images_bucket_acl" {
  depends_on = [
    aws_s3_bucket_ownership_controls.images_bucket_ownership,
    aws_s3_bucket_public_access_block.images_bucket_public_access,
  ]

  bucket = aws_s3_bucket.images_bucket.id
  acl    = "private"
}

# Configure CORS for the bucket
resource "aws_s3_bucket_cors_configuration" "images_bucket_cors" {
  bucket = aws_s3_bucket.images_bucket.id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "HEAD"]
    allowed_origins = [var.email_domain]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}

# Upload GIF to the bucket
resource "aws_s3_object" "loading_gif" {
  depends_on = [
    aws_s3_bucket.images_bucket,
    aws_s3_bucket_acl.images_bucket_acl,
  ]
  bucket       = aws_s3_bucket.images_bucket.id
  key          = "default/loading.gif"
  source       = "${path.module}/loading.gif"
  content_type = "image/gif"

  tags = {
    Name        = "Loading GIF"
    Environment = var.env_name
  }
}
