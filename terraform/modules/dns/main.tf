# Route53のホストゾーンを取得
resource "aws_route53_zone" "main" {
    name = "${var.domain_name}"
}

# Route53のレコードセット
resource "aws_route53_record" "myapp_dns" {
    zone_id = aws_route53_zone.main.zone_id
    name    = "cycle-route.${var.domain_name}"
    type    = "A"
    alias {
        name                   = var.alb_dns_name
        zone_id                = var.alb_zone_id
        evaluate_target_health = true
    }
}
