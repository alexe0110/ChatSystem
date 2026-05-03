# ALB — точка входа из интернета.
# Принимает на 80 порот, раздаёт на контейнеры
resource "aws_lb" "main" {
  name               = "${var.project}-alb"
  internal           = false                # false = публичный
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = [
    aws_subnet.public_a.id,
    aws_subnet.public_b.id,
  ]

  tags = { Name = "${var.project}-alb" }
}

# Target Group — куда ALB отправляет трафик.
resource "aws_lb_target_group" "rest" {
  name        = "${var.project}-tg"
  port        = 8080
  protocol    = "HTTP"
  vpc_id      = aws_vpc.main.id
  target_type = "ip"

  health_check {
    path                = "/health"
    healthy_threshold   = 2
    unhealthy_threshold = 3
    timeout             = 5
    interval            = 30
  }

  tags = { Name = "${var.project}-tg" }
}

# Listener — правило: всё что пришло на порт 80 → отправь в target group
resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.main.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.rest.arn
  }
}