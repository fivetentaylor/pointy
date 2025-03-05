resource "aws_iam_role" "ecs_task_role" {
  name = "ecs_task_role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
        Effect = "Allow"
        Sid    = ""
      }
    ]
  })
}

locals {
  ses_statement = var.email_region != local.region ? [
    {
      Effect = "Allow"
      Action = [
        "ses:SendEmail",
        "ses:SendRawEmail"
      ]
      Resource = [
        "arn:aws:ses:${var.email_region}:${local.account_id}:identity/*",
        "arn:aws:ses:${var.email_region}:${local.account_id}:configuration-set/*",
        "arn:aws:ses:${var.email_region}:${local.account_id}:*/*"
      ]
    }
  ] : []

  policy_statements = concat([
    {
      Effect = "Allow"
      Action = [
        "dynamodb:BatchGetItem",
        "dynamodb:BatchWriteItem",
        "dynamodb:ConditionCheckItem",
        "dynamodb:DeleteItem",
        "dynamodb:DescribeTable",
        "dynamodb:GetItem",
        "dynamodb:ListTables",
        "dynamodb:PutItem",
        "dynamodb:Query",
        "dynamodb:Scan",
        "dynamodb:TagResource",
        "dynamodb:UntagResource",
        "dynamodb:UpdateItem",
        "dynamodb:UpdateTable"
      ]
      Resource = [
        "arn:aws:dynamodb:${local.region}:${local.account_id}:table/${var.dynamo_table_name}",
        "arn:aws:dynamodb:${local.region}:${local.account_id}:table/${var.dynamo_table_name}/*",
        "arn:aws:dynamodb:${local.region}:${local.account_id}:table/${var.dynamo_table_name}/index/*"
      ]
    },
    {
      Effect = "Allow"
      Action = "s3:*"
      Resource = [
        "arn:aws:s3:::${var.docs_bucket_name}",
        "arn:aws:s3:::${var.docs_bucket_name}/*",
        "arn:aws:s3:::${var.images_bucket_name}",
        "arn:aws:s3:::${var.images_bucket_name}/*",
      ]
    },
    {
      Effect = "Allow"
      Action = [
        "ses:SendEmail",
        "ses:SendRawEmail"
      ]
      Resource = [
        "arn:aws:ses:${local.region}:${local.account_id}:identity/*",
        "arn:aws:ses:${local.region}:${local.account_id}:configuration-set/*",
        "arn:aws:ses:${local.region}:${local.account_id}:*/*"
      ]
    },
    {
      Effect = "Allow"
      Action = [
        "ecr:GetDownloadUrlForLayer",
        "ecr:BatchCheckLayerAvailability",
        "ecr:BatchGetImage",
        "ecr:GetAuthorizationToken"
      ]
      Resource = "*"
    },
    {
      Effect = "Allow"
      Action = [
        "ssm:DescribeAssociation",
        "ssm:GetDeployablePatchSnapshotForInstance",
        "ssm:GetDocument",
        "ssm:DescribeDocument",
        "ssm:GetManifest",
        "ssm:GetParameter",
        "ssm:GetParameters",
        "ssm:ListAssociations",
        "ssm:ListInstanceAssociations",
        "ssm:PutInventory",
        "ssm:PutComplianceItems",
        "ssm:PutConfigurePackageResult",
        "ssm:UpdateAssociationStatus",
        "ssm:UpdateInstanceAssociationStatus",
        "ssm:UpdateInstanceInformation"
      ]
      Resource = "*"
    },
    {
      Effect = "Allow"
      Action = [
        "ssmmessages:CreateControlChannel",
        "ssmmessages:CreateDataChannel",
        "ssmmessages:OpenControlChannel",
        "ssmmessages:OpenDataChannel"
      ]
      Resource = "*"
    },
    {
      Effect = "Allow"
      Action = [
        "ec2messages:AcknowledgeMessage",
        "ec2messages:DeleteMessage",
        "ec2messages:FailMessage",
        "ec2messages:GetEndpoint",
        "ec2messages:GetMessages",
        "ec2messages:SendReply"
      ]
      Resource = "*"
    },
    {
      Effect = "Allow"
      Action = [
        "ecs:ExecuteCommand",
        "ssm:SendCommand",
        "ssm:ListCommandInvocations",
        "ssm:GetCommandInvocation",
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ]
      Resource = "*"
    }
  ], local.ses_statement)
}

resource "aws_iam_policy" "ecs_task_role" {
  name        = "ecs_task_role"
  description = "Policy for Reviso ECS tasks"

  policy = jsonencode({
    Version   = "2012-10-17"
    Statement = local.policy_statements
  })
}

resource "aws_iam_role_policy_attachment" "ecs_task_role_s3_access" {
  role       = aws_iam_role.ecs_task_role.name
  policy_arn = aws_iam_policy.ecs_task_role.arn
}

resource "aws_iam_role" "task_execution_role" {
  name = "task_execution_role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = { Service = "ecs-tasks.amazonaws.com" }
    }]
  })
}

resource "aws_iam_policy" "task_execution_policy" {
  name        = "task_execution_policy"
  path        = "/"
  description = "Policy for ecs task execution"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "logs:DescribeLogStreams",

          "secretsmanager:*",

          "ecr:GetAuthorizationToken",
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage",
        ]
        Resource = "*"
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "task_execution_policy_attachment" {
  role       = aws_iam_role.task_execution_role.name
  policy_arn = aws_iam_policy.task_execution_policy.arn
}

