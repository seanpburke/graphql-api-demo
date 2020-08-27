provider "aws" {
  version = "~> 3.0"
  profile = "default"
  region  = "us-west-2"
}

locals {
  ecr_registry = (var.ecr_registry != "") ? var.ecr_registry : format("%s.dkr.ecr.%s.amazonaws.com", data.aws_caller_identity.current.account_id, data.aws_region.current.id)

  // Define as a local, so that we can encode to JSON in the task definition.
  container_definition = {
    name   = var.ecs_task_name
    image  = "${local.ecr_registry}/${var.ecr_image}:latest"
    cpu    = 256
    memory = 512
    portMappings = [
      { // In Fargate, hostPort must be equal to containerPort.
        containerPort = var.ecs_port
        hostPort      = var.ecs_port
        protocol      = "tcp"
      }
    ]
    logConfiguration = {
      logDriver = "awslogs"
      options = {
        awslogs-group         = "/ecs/graphql-api-demo"
        awslogs-region        = data.aws_region.current.id
        awslogs-stream-prefix = "ecs"
      }
    }
  }
}

resource "aws_ecs_service" "this" {
  name            = var.ecs_service_name
  cluster         = var.ecs_cluster_name
  task_definition = aws_ecs_task_definition.task.arn
  depends_on      = [aws_lb_listener.listener]
  desired_count   = 1
  launch_type     = "FARGATE"
  // Note that FARGATE requires the default scheduling_strategy 'REPLICA'.

  load_balancer {
    target_group_arn = aws_lb_target_group.tg.arn
    container_name   = var.ecs_task_name
    container_port   = var.ecs_port
  }

  // Required when the task definition uses network mode 'awsvpc'.
  network_configuration {
    assign_public_ip = true
    security_groups  = [aws_security_group.sg.id]
    subnets          = aws_subnet.sn[*].id
  }
}

resource "aws_ecs_task_definition" "task" {
  family                   = var.ecs_task_name
  cpu                      = 256
  memory                   = 512
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  task_role_arn            = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/ECS-Task-DynamoDB"
  execution_role_arn       = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/ecsTaskExecutionRole"
  container_definitions    = jsonencode([local.container_definition])
}

resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name        = "ECS-Default-VPC"
    Description = "Created for ECS cluster DEFAULT"
  }
}

resource "aws_subnet" "sn" {
  count             = length(var.subnet_cidr_blocks)
  vpc_id            = aws_vpc.main.id
  availability_zone = data.aws_availability_zones.available.names[(count.index + 2) % length(data.aws_availability_zones.available.names)]
  cidr_block        = var.subnet_cidr_blocks[count.index]

  tags = {
    Name = format("ECS-Default-Subnet-%d", count.index + 1)
  }
}

resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "ECS-Default-InternetGateway"
  }
}

resource "aws_route_table" "rt" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "ECS-Default-RouteTable"
  }
}

resource "aws_route" "public" {
  route_table_id         = aws_route_table.rt.id
  gateway_id             = aws_internet_gateway.gw.id
  destination_cidr_block = "0.0.0.0/0"
}

resource "aws_route_table_association" "a" {
  count          = length(var.subnet_cidr_blocks)
  subnet_id      = aws_subnet.sn[count.index].id
  route_table_id = aws_route_table.rt.id
}

resource "aws_security_group" "sg" {
  description = "ECS Allowed Ports"
  vpc_id      = aws_vpc.main.id

  ingress {
    protocol    = "tcp"
    from_port   = var.ecs_port
    to_port     = var.ecs_port
    cidr_blocks = [var.source_cidr]
  }

  ingress {
    protocol        = "tcp"
    from_port       = 1
    to_port         = 65535
    security_groups = [aws_security_group.elb.id]
  }

  /* By default, AWS creates an ALLOW ALL egress rule when creating a new Security Group inside of
     a VPC. When creating a new Security Group inside a VPC, Terraform will remove this default rule,
     and require you specifically re-create it if you desire that rule. We feel this leads to fewer
     surprises in terms of controlling your egress rules. If you desire this rule to be in place,
     you can use this egress block:
  */
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "ECS-Default-ECS-SecurityGroup"
  }
}

resource "aws_security_group" "elb" {
  description = "ELB Allowed Ports"
  vpc_id      = aws_vpc.main.id

  ingress {
    protocol    = "tcp"
    from_port   = var.elb_port
    to_port     = var.elb_port
    cidr_blocks = [var.source_cidr]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "ECS-Default-ELB-SecurityGroup"
  }
}

resource "aws_lb_target_group" "tg" {
  vpc_id      = aws_vpc.main.id
  port        = var.elb_port
  protocol    = "HTTP"
  target_type = "ip"

  tags = {
    Name = "ECS-Default-ELB-TargetGroup"
  }
}

resource "aws_lb" "elb" {
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.elb.id]
  subnets            = aws_subnet.sn[*].id

  //  access_logs {
  //    bucket  = aws_s3_bucket.lb_logs.bucket
  //    prefix  = "test-lb"
  //    enabled = true
  //  }

  tags = {
    Name = "ECS-Default-ELB"
  }
}

resource "aws_lb_listener" "listener" {
  load_balancer_arn = aws_lb.elb.arn
  port              = var.elb_port
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.tg.arn
  }
}
