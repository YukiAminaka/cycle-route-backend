output "alb_dns_name" {
  value = aws_lb.main.dns_name
}

output "alb_security_group_id" {
  value = aws_security_group.alb.id
}

output "target_group_arns" {
  value = {
    frontend = aws_lb_target_group.frontend.arn
    kratos   = aws_lb_target_group.kratos.arn
  }
}
