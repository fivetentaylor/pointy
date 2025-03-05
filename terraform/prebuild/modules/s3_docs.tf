resource "aws_s3_bucket" "doc_bucket" {
  bucket = var.docs_bucket_name
}

resource "aws_s3_bucket_ownership_controls" "doc_bucket" {
  bucket = aws_s3_bucket.doc_bucket.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "doc_bucket" {
  depends_on = [aws_s3_bucket_ownership_controls.doc_bucket]

  bucket = aws_s3_bucket.doc_bucket.id
  acl    = "private"
}
