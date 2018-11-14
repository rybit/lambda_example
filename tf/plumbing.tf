provider "aws" {
    version = "1.14.1"
    region  = "us-east-1"
    profile = "personal"
}

# the bucket we are going to drop files in and trigger lambda on
resource "aws_s3_bucket" "zips" {
    bucket = "rybit-lambda-zips"
}

# our current account 
data "aws_caller_identity" "current" {}
