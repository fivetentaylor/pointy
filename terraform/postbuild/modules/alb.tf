resource "aws_lb_target_group" "server_tg" {
  name     = "${var.preview_prefix}${var.name}-server-tg"
  port     = 9090 // The port your container is listening on
  protocol = "HTTP"
  vpc_id   = data.aws_vpc.main.id

  health_check {
    healthy_threshold   = 3
    unhealthy_threshold = 3
    timeout             = 30
    path                = "/" // Adjust if your service has a specific health check endpoint
    interval            = 60
    matcher             = "200-399"
  }

  target_type = "ip"
}

resource "aws_lb_target_group" "web_tg" {
  name     = "${var.preview_prefix}${var.name}-web-tg"
  port     = 3000
  protocol = "HTTP"
  vpc_id   = data.aws_vpc.main.id

  health_check {
    healthy_threshold   = 3
    unhealthy_threshold = 3
    timeout             = 30
    path                = "/" // Adjust if your service has a specific health check endpoint
    interval            = 60
    matcher             = "200-399"
  }

  target_type = "ip"
}

resource "aws_lb" "app_alb" {
  name               = "${var.preview_prefix}${var.name}-app-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [data.aws_security_group.app_security_group.id]
  subnets            = data.aws_subnets.public.ids
}

resource "aws_lb" "web_alb" {
  name               = "${var.preview_prefix}${var.name}-web-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [data.aws_security_group.app_security_group.id]
  subnets            = data.aws_subnets.public.ids
}

resource "aws_lb_listener" "server_front_end" {
  load_balancer_arn = aws_lb.app_alb.arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = aws_acm_certificate_validation.app_cert_validation.certificate_arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.server_tg.arn
  }
}

resource "aws_lb_listener" "web_front_end" {
  load_balancer_arn = aws_lb.web_alb.arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = aws_acm_certificate_validation.www_cert_validation.certificate_arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.web_tg.arn
  }
}

resource "aws_lb_listener" "server_redirect" {
  load_balancer_arn = aws_lb.app_alb.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type = "redirect"

    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
      host        = "#{host}"
      path        = "/#{path}"
      query       = "#{query}"
    }
  }
}

resource "aws_lb_listener" "www_redirect" {
  load_balancer_arn = aws_lb.web_alb.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type = "redirect"

    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
      host        = "#{host}"
      path        = "/#{path}"
      query       = "#{query}"
    }
  }
}
