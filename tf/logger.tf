resource "aws_lambda_function" "lambdalogger" {
  role = "${aws_iam_role.lambdalogger.arn}"

  s3_bucket     = "${aws_s3_bucket.zips.id}"
  s3_key        = "${var.project}/logger_${terraform.workspace == "prod" ?  var.sha : "latest"}.zip"
  function_name = "logger_${terraform.workspace == "prod" ? "prod" : "dev"}"
  handler       = "logger.out"
  runtime       = "go1.x"
  memory_size   = 256
  timeout       = 300

  environment {
    variables = {
      ENV           = "${terraform.workspace == "prod" ? "prod" : "dev"}"
      SHA           = "latest"
      LOG_LEVEL     = "debug"

      HUMIO_TOKEN      = "XKQuY2DL27gJ8CKlHy63ktinFGAgL0gLidGxjsQBHnvB"
      HUMIO_REPOSITORY = "sandbox"
    }
  }

  tags {
    env = "${terraform.workspace == "prod" ? "prod" : "dev"}"
  }
}

resource "aws_lambda_permission" "allow_cloudwatch" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  principal     = "logs.us-east-1.amazonaws.com"
  function_name = "${aws_lambda_function.lambdalogger.function_name}"
  source_arn    = "arn:aws:logs:us-east-1:${data.aws_caller_identity.current.account_id}:log-group:*:*"
}
