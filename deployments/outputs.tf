output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "region_name" {
  value = data.aws_region.current.name
}

output "elb" {
  value = aws_lb.elb
}

output "tg" {
  value = aws_lb_target_group.tg
}

output "cd" {
  value = jsonencode(local.container_definition)
}
