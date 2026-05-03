# CloudWatch Log Group — куда пойдут логи контейнера

resource "aws_cloudwatch_log_group" "app" {
  name              = "/ecs/${var.project}"
  retention_in_days = 7
}

# Task Definition — описание контейнера
resource "aws_ecs_task_definition" "app" {
  family                   = "${var.project}-task"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "256"      # 0.25 vCPU
  memory                   = "512"      # 512 MB
  execution_role_arn       = aws_iam_role.ecs_execution.arn
  task_role_arn            = aws_iam_role.ecs_task.arn

  container_definitions = jsonencode([{
    name      = "chat-system"
    image     = "${aws_ecr_repository.chat_system.repository_url}:latest"
    essential = true

    portMappings = [{
      containerPort = 8080
      protocol      = "tcp"
    }]

    environment = [
      { name = "DB_TYPE",      value = "dynamodb" },
      { name = "STORAGE_TYPE", value = "s3" },
      { name = "AWS_REGION",   value = var.region },
    ]

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = aws_cloudwatch_log_group.app.name
        "awslogs-region"        = var.region
        "awslogs-stream-prefix" = "ecs"
      }
    }
  }])
}

# ECS Service — держит контейнер запущенным
resource "aws_ecs_service" "app" {
  name            = "${var.project}-service"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.app.arn
  desired_count   = 1              # одна копия приложения
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = [
      aws_subnet.public_a.id,
      aws_subnet.public_b.id,
    ]
    security_groups  = [aws_security_group.ecs.id]
    assign_public_ip = true    # нужен для скачивания образа из ECR
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.rest.arn
    container_name   = "chat-system"
    container_port   = 8080
  }

  # Ждём пока listener создастся, иначе ECS не сможет
  # зарегистрироваться в target group
  depends_on = [aws_lb_listener.http]
}