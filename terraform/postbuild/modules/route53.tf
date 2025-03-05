data "aws_route53_zone" "dns_zone" {
  provider = aws.dns_role
  zone_id  = var.route53_zone
}

# Request a certificate for the app subdomain
resource "aws_acm_certificate" "app_cert" {
  domain_name       = "${var.preview_prefix}${var.app_domain}"
  validation_method = "DNS"
}

# Request a certificate for the www subdomain
resource "aws_acm_certificate" "www_cert" {
  domain_name       = "${var.preview_prefix}${var.web_domain}"
  validation_method = "DNS"
}

resource "aws_route53_record" "app_cert_validation" {
  provider = aws.dns_role
  for_each = {
    for dvo in toset(aws_acm_certificate.app_cert.domain_validation_options) : dvo.domain_name => {
      name  = dvo.resource_record_name
      type  = dvo.resource_record_type
      value = dvo.resource_record_value
    }
  }

  name    = each.value.name
  type    = each.value.type
  zone_id = data.aws_route53_zone.dns_zone.zone_id
  records = [each.value.value]
  ttl     = 60
}

resource "aws_route53_record" "www_cert_validation" {
  provider = aws.dns_role
  for_each = {
    for dvo in toset(aws_acm_certificate.www_cert.domain_validation_options) : dvo.domain_name => {
      name  = dvo.resource_record_name
      type  = dvo.resource_record_type
      value = dvo.resource_record_value
    }
  }

  name    = each.value.name
  type    = each.value.type
  zone_id = data.aws_route53_zone.dns_zone.zone_id
  records = [each.value.value]
  ttl     = 60
}

resource "aws_acm_certificate_validation" "app_cert_validation" {
  certificate_arn         = aws_acm_certificate.app_cert.arn
  validation_record_fqdns = [for record in values(aws_route53_record.app_cert_validation) : record.fqdn]
}

resource "aws_acm_certificate_validation" "www_cert_validation" {
  certificate_arn         = aws_acm_certificate.www_cert.arn
  validation_record_fqdns = [for record in values(aws_route53_record.www_cert_validation) : record.fqdn]
}

# Route 53 record for the app subdomain
resource "aws_route53_record" "app" {
  provider = aws.dns_role
  zone_id  = data.aws_route53_zone.dns_zone.zone_id
  name     = "${var.preview_prefix}${var.app_domain}"
  type     = "A"

  alias {
    name                   = aws_lb.app_alb.dns_name
    zone_id                = aws_lb.app_alb.zone_id
    evaluate_target_health = true
  }
}

# Route 53 record for the www subdomain
resource "aws_route53_record" "www" {
  provider = aws.dns_role
  zone_id  = data.aws_route53_zone.dns_zone.zone_id
  name     = "${var.preview_prefix}${var.web_domain}"
  type     = "A"

  alias {
    name                   = aws_lb.web_alb.dns_name
    zone_id                = aws_lb.web_alb.zone_id
    evaluate_target_health = true
  }
}
