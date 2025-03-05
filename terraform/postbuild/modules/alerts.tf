resource "aws_sns_topic" "ecs_alerts" {
  name = "${var.preview_prefix}ecs-cpu-alerts"
}

data "archive_file" "lambda" {
  type        = "zip"
  source_file = "${path.module}/alert_function.py"
  output_path = "${path.module}/alert_function.zip"
}

resource "aws_lambda_function" "ecs_alerts" {
  function_name = "${var.preview_prefix}ECSCPUAlerts"
  filename      = data.archive_file.lambda.output_path
  handler       = "alert_function.lambda_handler" # File name and function name
  runtime       = "python3.8"
  role          = aws_iam_role.lambda_exec.arn

  source_code_hash = filebase64sha256(data.archive_file.lambda.output_path)

  environment {
    variables = {
      SLACK_WEBHOOK_URL = var.slack_webhook_url
    }
  }
}

# IAM role for Lambda function
resource "aws_iam_role" "lambda_exec" {
  name = "${var.preview_prefix}lambda_execution_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      },
    ]
  })
}

# IAM policy to allow publishing to SNS
resource "aws_iam_policy" "lambda_sns_policy" {
  name        = "${var.preview_prefix}lambda_sns_policy"
  path        = "/"
  description = "IAM policy for publishing SNS messages"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Effect   = "Allow"
        Resource = "arn:aws:logs:*:*:*"
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_logs" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = aws_iam_policy.lambda_sns_policy.arn
}

resource "aws_sns_topic_subscription" "lambda_subscription" {
  topic_arn = aws_sns_topic.ecs_alerts.arn
  protocol  = "lambda"
  endpoint  = aws_lambda_function.ecs_alerts.arn
}

resource "aws_lambda_permission" "allow_sns_to_call_lambda" {
  statement_id  = "AllowExecutionFromSNS"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.ecs_alerts.function_name
  principal     = "sns.amazonaws.com"
  source_arn    = aws_sns_topic.ecs_alerts.arn
}

resource "aws_cloudwatch_metric_alarm" "high_cpu_utilization" {
  alarm_name                = "${var.preview_prefix}high-cpu-utilization"
  comparison_operator       = "GreaterThanOrEqualToThreshold"
  evaluation_periods        = "2"
  metric_name               = "CPUUtilization"
  namespace                 = "AWS/ECS"
  period                    = "60"
  statistic                 = "Average"
  threshold                 = "90"
  alarm_description         = "This metric monitors ecs cpu utilization"
  insufficient_data_actions = []

  dimensions = {
    ClusterName = aws_ecs_cluster.reviso.name
    ServiceName = aws_ecs_service.reviso-server.name
  }

  alarm_actions = [aws_sns_topic.ecs_alerts.arn]
}

resource "aws_cloudwatch_metric_alarm" "low_cpu_utilization" {
  alarm_name                = "${var.preview_prefix}low-cpu-utilization"
  comparison_operator       = "LessThanOrEqualToThreshold"
  evaluation_periods        = "2"
  metric_name               = "CPUUtilization"
  namespace                 = "AWS/ECS"
  period                    = "60"
  statistic                 = "Average"
  threshold                 = "0"
  alarm_description         = "This metric monitors ecs cpu utilization"
  insufficient_data_actions = []

  dimensions = {
    ClusterName = aws_ecs_cluster.reviso.name
    ServiceName = aws_ecs_service.reviso-server.name
  }

  alarm_actions = [aws_sns_topic.ecs_alerts.arn]
}

resource "aws_cloudwatch_metric_alarm" "health_check_failures" {
  alarm_name                = "${var.preview_prefix}health-check-failures"
  comparison_operator       = "GreaterThanOrEqualToThreshold"
  evaluation_periods        = "1"
  metric_name               = "HealthCheckFailed"
  namespace                 = "AWS/ECS"
  period                    = "60" # Period in seconds
  statistic                 = "Sum"
  threshold                 = "1"
  alarm_description         = "This metric monitors failed health checks for the ECS service"
  insufficient_data_actions = []

  dimensions = {
    ClusterName = aws_ecs_cluster.reviso.name
    ServiceName = aws_ecs_service.reviso-server.name
  }

  alarm_actions = [aws_sns_topic.ecs_alerts.arn]
}
