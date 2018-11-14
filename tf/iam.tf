
resource "aws_iam_role" "lambdalogger" {
    name = "lambdalogger_${terraform.workspace == "default" ? "dev" : "prod"}"

   assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

# this is a POLICY to attach to the above role so we can 
# send logs to cloudwatch, which is necessary for our own output
resource "aws_iam_role_policy" "cloudwatch_logs" {
    name = "allows_cloudwatch_logs_${terraform.workspace == "prod" ? "prod" : "dev"}"
    role = "${aws_iam_role.lambdalogger.id}"

    policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:*"
      ],
      "Resource": "arn:aws:logs:*:*:*"
    }
  ]
}
EOF
}

# S3 access so that we can load our own zip file
# TODO: scope it down to read only
resource "aws_iam_role_policy" "s3_zips" {
    name = "s3_zips_${terraform.workspace == "default" ? "dev" : "prod"}"
    role = "${aws_iam_role.lambdalogger.id}"

    policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:*"
      ],
      "Resource": [
        "arn:aws:s3:::${aws_s3_bucket.zips.id}/*"
      ]
    }
  ]
}
EOF
}