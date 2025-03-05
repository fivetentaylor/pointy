resource "aws_db_instance" "postgres" {
  count                = var.create_db ? 1 : 0
  identifier           = "${var.env_name}-postgres-db"
  engine               = "postgres"
  engine_version       = "16.4"
  instance_class       = "db.t3.micro"
  allocated_storage    = 20
  username             = jsondecode(data.aws_secretsmanager_secret_version.db_secret[0].secret_string)["DB_USERNAME"]
  password             = jsondecode(data.aws_secretsmanager_secret_version.db_secret[0].secret_string)["DB_PASSWORD"]
  parameter_group_name = "default.postgres16"
  skip_final_snapshot  = true

  tags = {
    Name        = "${var.env_name}-postgres-db"
    Environment = var.env_name
  }

  lifecycle {
    prevent_destroy = true
  }
}

data "aws_secretsmanager_secret_version" "db_secret" {
  count     = var.create_db ? 1 : 0
  secret_id = var.secret_arn
}
