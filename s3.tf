resource "aws_s3_bucket" "feic" {
  bucket = "terraform-bucket-horacio-feic"
  force_destroy = true
}

resource "aws_s3_bucket_public_access_block" "feic" {
  bucket = aws_s3_bucket.feic.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_policy" "feic" {
  bucket = aws_s3_bucket.feic.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = "*"
        Action = [
          "s3:GetObject",
          "s3:PutObject"
        ]
        Resource = "${aws_s3_bucket.feic.arn}/*"
      }
    ]
  })
}

