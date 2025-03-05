output "vpc_id" {
  value = aws_vpc.main.id
}
output "repository_url" {
  value = aws_ecr_repository.reviso-server.repository_url
}
output "internet_gateway_id" {
  value = aws_internet_gateway.gw.id
}
output "route_table_id" {
  value = aws_route_table.rt.id
}
output "app_security_group_id" {
  value = aws_security_group.app_sg.id
}
output "internal_security_group_id" {
  value = aws_security_group.internal_sg.id
}
