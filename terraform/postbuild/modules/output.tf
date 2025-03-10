output "app_host" {
  value = "http://${var.preview_prefix}${var.app_domain}"
}

output "ecs_deployment_task_definition" {
  value = aws_ecs_service.server.task_definition
}
