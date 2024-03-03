output "url" {
  description = "The DNS name of the load balancer"
  value       = "http://${module.ecs.dns_name}"
}
