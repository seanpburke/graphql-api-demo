variable "ecr_registry" {
  description = "ECR Registry Name, defaults to <id>.dkr.ecr.<region>.amazonaws.com"
  default     = ""
}

variable "ecr_image" {
  description = "ECR Image Name, defaults to graphql-api-demo"
  default     = "graphql-api-demo"
}

variable "ecs_cluster_name" {
  description = "ECS Cluster Name"
  default     = "default"
}

variable "ecs_service_name" {
  description = "ECS Service Name"
  default     = "graphql-api-demo-service"
}

variable "ecs_task_name" {
  description = "ECS Task Family Name"
  default     = "graphql-api-demo-task"
}

variable "ecs_port" {
  description = "Security Group port to open on ECS instances - defaults to port 8080"
  default     = 8080
}

variable "elb_port" {
  description = "Security Group port to open on ELB - port 8080 will be open by default"
  default     = 8080
}

variable "cidr_block" {
  description = "CIDR/IP range for the VPC"
  default     = "10.0.0.0/16"
}

variable "subnet_cidr_blocks" {
  description = "CIDR/IP range for the subnets"
  default     = ["10.0.0.0/24", "10.0.1.0/24"]
}

variable "source_cidr" {
  description = "CIDR/IP range for EcsPort and ElbPort - defaults to 0.0.0.0/0"
  default     = "0.0.0.0/0"
}
