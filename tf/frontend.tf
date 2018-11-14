resource "aws_lambda_function" "frontend" {
  role = "${aws_iam_role.lambdalogger.arn}"

  s3_bucket     = "${aws_s3_bucket.zips.id}"
  s3_key        = "${var.project}/frontend_${terraform.workspace == "prod" ?  var.sha : "latest"}.zip"
  function_name = "frontend_${terraform.workspace == "prod" ? "prod" : "dev"}"
  handler       = "frontend.out"
  runtime       = "go1.x"
  memory_size   = 256
  timeout       = 300

  environment {
    variables = {
      ENV           = "${terraform.workspace == "prod" ? "prod" : "dev"}"
      SHA           = "latest"
      LOG_LEVEL     = "debug"
    }
  }

  tags {
    env = "${terraform.workspace == "prod" ? "prod" : "dev"}"
  }
}

resource "aws_cloudwatch_log_subscription_filter" "frontend" {
  name            = "frontend_${terraform.workspace == "prod" ? "prod" : "dev"}"
  log_group_name  = "${aws_cloudwatch_log_group.frontend.name}"
  destination_arn = "${aws_lambda_function.lambdalogger.arn}"
  filter_pattern  = ""
}

resource "aws_cloudwatch_log_group" "frontend" {
  name              = "/aws/lambda/${aws_lambda_function.frontend.function_name}"
  retention_in_days = 30
}