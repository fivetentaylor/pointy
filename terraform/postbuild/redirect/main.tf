# Variables for domain configuration
variable "old_root_domain" {
  description = "The root domain being redirected from (e.g., revi.so)"
  type        = string
}

variable "old_web_domain" {
  description = "The www domain being redirected from (e.g., www.revi.so)"
  type        = string
}

variable "old_app_domain" {
  description = "The app domain being redirected from (e.g., app.revi.so)"
  type        = string
}

variable "new_root_domain" {
  description = "The root domain being redirected to (e.g., pointy.ai)"
  type        = string
}

variable "new_web_domain" {
  description = "The www domain being redirected to (e.g., www.pointy.ai)"
  type        = string
}

variable "new_app_domain" {
  description = "The app domain being redirected to (e.g., app.pointy.ai)"
  type        = string
}

# Access the existing Route53 zone for the old domain
data "aws_route53_zone" "old_domain_zone" {
  provider = aws.dns_role
  name     = "${var.old_root_domain}."  # Make sure to include the trailing dot
}

# Request certificates for old domains
resource "aws_acm_certificate" "old_domain_cert" {
  domain_name               = var.old_root_domain
  subject_alternative_names = [var.old_web_domain, var.old_app_domain]
  validation_method         = "DNS"
}

# Validation records for old domain certificate
resource "aws_route53_record" "old_domain_cert_validation" {
  provider = aws.dns_role
  for_each = {
    for dvo in aws_acm_certificate.old_domain_cert.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      type   = dvo.resource_record_type
      value  = dvo.resource_record_value
    }
  }
  name    = each.value.name
  type    = each.value.type
  zone_id = data.aws_route53_zone.old_domain_zone.zone_id
  records = [each.value.value]
  ttl     = 60
}

# Certificate validation for old domain
resource "aws_acm_certificate_validation" "old_domain_cert_validation" {
  certificate_arn         = aws_acm_certificate.old_domain_cert.arn
  validation_record_fqdns = [for record in aws_route53_record.old_domain_cert_validation : record.fqdn]
}

# Create ALB for handling redirects
resource "aws_lb" "redirect_alb" {
  name               = "${replace(var.old_root_domain, ".", "-")}-redirect-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [data.aws_security_group.app_security_group.id]
  subnets            = data.aws_subnets.public.ids
}

# Route 53 records for old domains pointing to the redirect ALB
resource "aws_route53_record" "old_root_domain" {
  provider = aws.dns_role
  zone_id  = data.aws_route53_zone.old_domain_zone.zone_id
  name     = var.old_root_domain
  type     = "A"
  alias {
    name                   = aws_lb.redirect_alb.dns_name
    zone_id                = aws_lb.redirect_alb.zone_id
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "old_web_domain" {
  provider = aws.dns_role
  zone_id  = data.aws_route53_zone.old_domain_zone.zone_id
  name     = var.old_web_domain
  type     = "A"
  alias {
    name                   = aws_lb.redirect_alb.dns_name
    zone_id                = aws_lb.redirect_alb.zone_id
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "old_app_domain" {
  provider = aws.dns_role
  zone_id  = data.aws_route53_zone.old_domain_zone.zone_id
  name     = var.old_app_domain
  type     = "A"
  alias {
    name                   = aws_lb.redirect_alb.dns_name
    zone_id                = aws_lb.redirect_alb.zone_id
    evaluate_target_health = true
  }
}

# HTTPS listener for redirecting old domains to new domains
resource "aws_lb_listener" "https_redirect" {
  load_balancer_arn = aws_lb.redirect_alb.arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = aws_acm_certificate_validation.old_domain_cert_validation.certificate_arn
  
  default_action {
    type = "fixed-response"
    fixed_response {
      content_type = "text/plain"
      message_body = "Not Found"
      status_code  = "404"
    }
  }
}

# Add a listener rule for the root domain redirect
resource "aws_lb_listener_rule" "https_root_domain_redirect" {
  listener_arn = aws_lb_listener.https_redirect.arn
  priority     = 10

  condition {
    host_header {
      values = [var.old_root_domain]
    }
  }

  action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
      host        = var.new_root_domain
      path        = "/#{path}"
      query       = "#{query}"
    }
  }
}

# Add a listener rule for the www domain redirect
resource "aws_lb_listener_rule" "https_www_domain_redirect" {
  listener_arn = aws_lb_listener.https_redirect.arn
  priority     = 20

  condition {
    host_header {
      values = [var.old_web_domain]
    }
  }

  action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
      host        = var.new_web_domain
      path        = "/#{path}"
      query       = "#{query}"
    }
  }
}

# Add a listener rule for the app domain redirect
resource "aws_lb_listener_rule" "https_app_domain_redirect" {
  listener_arn = aws_lb_listener.https_redirect.arn
  priority     = 30

  condition {
    host_header {
      values = [var.old_app_domain]
    }
  }

  action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
      host        = var.new_app_domain
      path        = "/#{path}"
      query       = "#{query}"
    }
  }
}

# HTTP listener for redirecting all HTTP traffic to the appropriate HTTPS domain
resource "aws_lb_listener" "http_redirect" {
  load_balancer_arn = aws_lb.redirect_alb.arn
  port              = 80
  protocol          = "HTTP"
  
  default_action {
    type = "fixed-response"
    fixed_response {
      content_type = "text/plain"
      message_body = "Not Found"
      status_code  = "404"
    }
  }
}

# HTTP listener rule for root domain
resource "aws_lb_listener_rule" "http_root_domain_redirect" {
  listener_arn = aws_lb_listener.http_redirect.arn
  priority     = 10

  condition {
    host_header {
      values = [var.old_root_domain]
    }
  }

  action {
    type = "redirect"
    redirect {
      port        = "80"
      protocol    = "HTTP"
      status_code = "HTTP_301"
      host        = var.new_root_domain
      path        = "/#{path}"
      query       = "#{query}"
    }
  }
}

# HTTP listener rule for www domain
resource "aws_lb_listener_rule" "http_www_domain_redirect" {
  listener_arn = aws_lb_listener.http_redirect.arn
  priority     = 20

  condition {
    host_header {
      values = [var.old_web_domain]
    }
  }

  action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
      host        = var.new_web_domain
      path        = "/#{path}"
      query       = "#{query}"
    }
  }
}

# HTTP listener rule for app domain
resource "aws_lb_listener_rule" "http_app_domain_redirect" {
  listener_arn = aws_lb_listener.http_redirect.arn
  priority     = 30

  condition {
    host_header {
      values = [var.old_app_domain]
    }
  }

  action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
      host        = var.new_app_domain
      path        = "/#{path}"
      query       = "#{query}"
    }
  }
}

# Optional: Create a dummy target group
# ALB requires a target group even if we're only doing redirects
resource "aws_lb_target_group" "dummy_tg" {
  name     = "${replace(var.old_root_domain, ".", "-")}-dummy-tg"
  port     = 80
  protocol = "HTTP"
  vpc_id   = data.aws_vpc.main.id
  
  health_check {
    enabled = false
  }
}
