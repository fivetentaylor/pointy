data "aws_iam_role" "ecs_task_role" {
  name = "ecs_task_role"
}
data "aws_iam_role" "task_execution_role" {
  name = "task_execution_role"
}

data "aws_ecr_repository" "server" {
  name = "${var.name}-server"
}

data "aws_ecr_repository" "web" {
  name = "${var.name}-web"
}

resource "aws_ecs_cluster" "default" {
  name = "${var.preview_prefix}${var.name}"
}

resource "aws_ecs_task_definition" "server" {
  family                   = "${var.preview_prefix}${var.name}-server"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  execution_role_arn       = data.aws_iam_role.task_execution_role.arn
  task_role_arn            = data.aws_iam_role.ecs_task_role.arn
  cpu                      = "1024"
  memory                   = "2048"

  runtime_platform {
    operating_system_family = "LINUX"
    cpu_architecture        = "X86_64"
  }

  container_definitions = jsonencode([{
    name                 = "main"
    image                = "${data.aws_ecr_repository.server.repository_url}:${var.server_sha}"
    essential            = true
    enableExecuteCommand = true

    portMappings = [
      {
        containerPort = 9090
        hostPort      = 9090
      },
      {
        containerPort = 9091
        hostPort      = 9091
      }
    ]

    environment = [
      {
        name  = "ENV"
        value = var.env
      },
      {
        name  = "NODE_ENV"
        value = "production"
      },
      {
        name  = "EMAIL_DOMAIN"
        value = var.email_domain
      },
      {
        name  = "EMAIL_REGION"
        value = var.email_region
      },
      {
        name  = "ADDR"
        value = ":9090"
      },
      {
        name  = "WORKER_ADDR"
        value = ":9091"
      },
      {
        name  = "WORKER_CONCURRENCY"
        value = "2"
      },
      {
        name  = "ALLOWED_ORIGINS"
        value = "https://${var.preview_prefix}${var.web_domain}"
      },
      {
        name  = "WEB_HOST"
        value = "https://${var.preview_prefix}${var.web_domain}"
      },
      {
        name  = "APP_HOST"
        value = "https://${var.preview_prefix}${var.app_domain}"
      },
      {
        name  = "WS_HOST"
        value = "wss://${var.preview_prefix}${var.app_domain}"
      },
      {
        name  = "AWS_DYNAMODB_TABLE"
        value = var.dynamo_table_name
      },
      {
        name  = "AWS_S3_BUCKET"
        value = var.docs_bucket_name
      },
      {
        name  = "AWS_S3_IMAGES_BUCKET"
        value = var.images_bucket_name
      },
      {
        name  = "COOKIE_DOMAIN"
        value = var.cookie_domain
      },
      {
        name  = "GOOGLE_REDIRECT_URI"
        value = "https://${var.preview_prefix}${var.app_domain}/auth/google/callback"
      },
      {
        name  = "PUBLIC_POSTHOG_HOST"
        value = "https://us.posthog.com"
      },
      {
        name  = "FREEPLAY_PROJECT_ID"
        value = "c6c19eb5-fdf8-4b6c-8664-2344443ec657"
      },
      {
        name  = "FREEPLAY_URL"
        value = "https://reviso.freeplay.ai"
      },
      {
        name  = "FREEPLAY_ENV"
        value = "${var.freeplay_env}"
      },
      {
        name  = "OTEL_SERVICE_NAME"
        value = "reviso-go"
      },
      {
        name  = "OTEL_EXPORTER_OTLP_PROTOCOL"
        value = "http/protobuf"
      },
      {
        name  = "OTEL_EXPORTER_OTLP_ENDPOINT"
        value = "https://api.honeycomb.io:443"
      }
    ]

    secrets = [
      {
        name      = "DATABASE_URL"
        valueFrom = "${var.secret_arn}:DATABASE_URL::"
      },
      {
        name      = "REDIS_URL"
        valueFrom = "${var.secret_arn}:REDIS_URL::"
      },
      {
        name      = "LOOPS_TOKEN"
        valueFrom = "${var.secret_arn}:LOOPS_TOKEN::"
      },
      {
        name      = "GOOGLE_APPLICATION_CREDENTIALS"
        valueFrom = "${var.secret_arn}:GOOGLE_APPLICATION_CREDENTIALS::"
      },
      {
        name      = "GOOGLE_CLOUD_PROJECT"
        valueFrom = "${var.secret_arn}:GOOGLE_CLOUD_PROJECT::"
      },
      {
        name      = "GOOGLE_CLOUD_LOCATION"
        valueFrom = "${var.secret_arn}:GOOGLE_CLOUD_LOCATION::"
      },
      {
        name      = "OPENAI_API_KEY"
        valueFrom = "${var.secret_arn}:OPENAI_API_KEY::"
      },
      {
        name      = "ANTHROPIC_API_KEY"
        valueFrom = "${var.secret_arn}:ANTHROPIC_API_KEY::"
      },
      {
        name      = "GROQ_API_KEY"
        valueFrom = "${var.secret_arn}:GROQ_API_KEY::"
      },
      {
        name      = "EMAIL_ALLOW_LIST"
        valueFrom = "${var.secret_arn}:EMAIL_ALLOW_LIST::"
      },
      {
        name      = "GOOGLE_CLIENT_ID"
        valueFrom = "${var.secret_arn}:GOOGLE_CLIENT_ID::"
      },
      {
        name      = "GOOGLE_CLIENT_SECRET"
        valueFrom = "${var.secret_arn}:GOOGLE_CLIENT_SECRET::"
      },
      {
        name      = "FREEPLAY_API_KEY"
        valueFrom = "${var.secret_arn}:FREEPLAY_API_KEY::"
      },
      {
        name      = "SEGMENT_KEY"
        valueFrom = "${var.secret_arn}:SEGMENT_WRITE_KEY::"
      },
      {
        name      = "PUBLIC_POSTHOG_KEY"
        valueFrom = "${var.secret_arn}:POSTHOG_KEY::"
      },
      {
        name      = "POSTHOG_SERVER_FEATURE_FLAG_KEY"
        valueFrom = "${var.secret_arn}:POSTHOG_SERVER_FEATURE_FLAG_KEY::"
      },
      {
        name      = "OTEL_EXPORTER_OTLP_HEADERS"
        valueFrom = "${var.secret_arn}:OTEL_EXPORTER_OTLP_HEADERS::"
      },
      {
        name      = "JWT_SECRET"
        valueFrom = "${var.secret_arn}:JWT_SECRET::"
      },
      {
        name      = "SCRAPINGBEE_API_KEY"
        valueFrom = "${var.secret_arn}:SCRAPINGBEE_API_KEY::"
      },
      {
        name      = "STRIPE_API_KEY"
        valueFrom = "${var.secret_arn}:STRIPE_API_KEY::"
      },
      {
        name      = "STRIPE_WEBHOOK_SECRET"
        valueFrom = "${var.secret_arn}:STRIPE_WEBHOOK_SECRET::"
      },
    ]

    logConfiguration = {
      logDriver = "awsfirelens"
    }

    healthCheck = {
      command     = ["CMD-SHELL", "curl -f http://localhost:9090/healthcheck || exit 1"]
      interval    = 30
      timeout     = 30
      retries     = 10
      startPeriod = 60
    }
    }, {
    name              = "log_router"
    image             = "betterstack/aws-ecs-fluent-bit:amd64-latest"
    cpu               = 256
    memory            = 512
    memoryReservation = 50
    essential         = true
    secrets = [
      {
        name      = "LOGTAIL_SOURCE_TOKEN"
        valueFrom = "${var.secret_arn}:LOGTAIL_SOURCE_TOKEN::"
      }
    ]
    firelensConfiguration = {
      type = "fluentbit"
      options = {
        "config-file-type"        = "file"
        "config-file-value"       = "/fluent-bit-logtail.conf"
        "enable-ecs-log-metadata" = "true"
      }
    }
  }])
}

resource "aws_ecs_task_definition" "web" {
  family                   = "${var.preview_prefix}${var.name}-web"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  execution_role_arn       = data.aws_iam_role.task_execution_role.arn
  task_role_arn            = data.aws_iam_role.ecs_task_role.arn
  cpu                      = "256"
  memory                   = "1024"

  runtime_platform {
    operating_system_family = "LINUX"
    cpu_architecture        = "X86_64"
  }

  container_definitions = jsonencode([{
    name      = "main"
    image     = "${data.aws_ecr_repository.web.repository_url}:${var.web_sha}"
    essential = true

    portMappings = [
      {
        containerPort = 3000
        hostPort      = 3000
      }
    ]

    environment = [
      {
        name  = "NODE_ENV"
        value = var.env
      },
      {
        name  = "NEXT_PUBLIC_APP_HOST"
        value = "https://${var.preview_prefix}${var.app_domain}"
      },
      {
        name  = "NEXT_PUBLIC_WS_HOST"
        value = "wss://${var.preview_prefix}${var.app_domain}"
      },
      {
        name  = "NEXT_PUBLIC_POSTHOG_HOST"
        value = "https://us.posthog.com"
      },
      {
        name  = "WEB_HOST"
        value = "https://${var.preview_prefix}${var.web_domain}"
      },
      {
        name  = "OTEL_SERVICE_NAME"
        value = "${var.preview_prefix}next-web"
      }
    ]

    secrets = [
      {
        name      = "NEXT_PUBLIC_GOOGLE_CLIENT_ID"
        valueFrom = "${var.secret_arn}:GOOGLE_CLIENT_ID::"
      },
      {
        name      = "NEXT_PUBLIC_POSTHOG_KEY"
        valueFrom = "${var.secret_arn}:POSTHOG_KEY::"
      },
      {
        name      = "NEXT_PUBLIC_SEGMENT_WRITE_KEY"
        valueFrom = "${var.secret_arn}:SEGMENT_WRITE_KEY::"
      },
      {
        name      = "HONEYCOMB_API_KEY"
        valueFrom = "${var.secret_arn}:HONEYCOMB_API_KEY::"
      }
    ]

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-create-group"  = "true"
        "awslogs-group"         = "/ecs/${var.preview_prefix}${var.name}-web/main"
        "awslogs-region"        = local.region
        "awslogs-stream-prefix" = "ecs"
      }
    }
  }])
}

resource "aws_ecs_service" "server" {
  name            = "${var.preview_prefix}${var.name}-server"
  cluster         = aws_ecs_cluster.default.id
  task_definition = aws_ecs_task_definition.server.arn
  desired_count   = var.desired_server_count
  launch_type     = "FARGATE"

  load_balancer {
    target_group_arn = aws_lb_target_group.server_tg.arn
    container_name   = "main"
    container_port   = 9090
  }

  network_configuration {
    subnets          = data.aws_subnets.public.ids
    security_groups  = [data.aws_security_group.app_security_group.id, data.aws_security_group.internal_security_group.id]
    assign_public_ip = true
  }

  deployment_circuit_breaker {
    enable   = true
    rollback = true
  }

  deployment_controller {
    type = "ECS"
  }

  health_check_grace_period_seconds = 120

  deployment_maximum_percent         = 200
  deployment_minimum_healthy_percent = 100

  depends_on = [
    aws_lb_listener.server_front_end,
  ]
}

resource "aws_ecs_service" "web" {
  name            = "${var.preview_prefix}${var.name}-web"
  cluster         = aws_ecs_cluster.default.id
  task_definition = aws_ecs_task_definition.web.arn
  desired_count   = var.desired_web_count
  launch_type     = "FARGATE"

  load_balancer {
    target_group_arn = aws_lb_target_group.web_tg.arn
    container_name   = "main"
    container_port   = 3000
  }

  network_configuration {
    subnets          = data.aws_subnets.public.ids
    security_groups  = [data.aws_security_group.app_security_group.id, data.aws_security_group.internal_security_group.id]
    assign_public_ip = true
  }

  depends_on = [
    aws_lb_listener.web_front_end,
  ]
}
