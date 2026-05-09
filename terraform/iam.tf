# Тут две роли
# - Execution Role — для ECS-платформы, чтобы качать образы и тд
# - Task Role — для кода, чтобы он мог рабоать с системами AWS без ключей Access Key/Secret Key

# ECS Task Execution Role — роль для самого ECS-сервиса, чтоб ECS мог скачать образ и писать логи
resource "aws_iam_role" "ecs_execution" {
  name = "${var.project}-ecs-execution"

  # Кто может юзать эту роль
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = { Service = "ecs-tasks.amazonaws.com" }
    }]
  })
}

# AWS политика даёт доступ ecs_execution к ECR и CloudWatch
resource "aws_iam_role_policy_attachment" "ecs_execution" {
  role       = aws_iam_role.ecs_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# Роль для контенера с приложением, чтобы код получил доступ к S3 и тд без ключей
# Как service account в K8s.
resource "aws_iam_role" "ecs_task" {
  name = "${var.project}-ecs-task"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = { Service = "ecs-tasks.amazonaws.com" }
    }]
  })
}

# разрешаем контейнеру ходить в DynamoDB и S3
resource "aws_iam_role_policy" "ecs_task" {
  name = "${var.project}-task-policy"
  role = aws_iam_role.ecs_task.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:PutItem",
          "dynamodb:GetItem",
          "dynamodb:Query",
          "dynamodb:DeleteItem",
          "dynamodb:UpdateItem",
        ]
        # Доступ только к нашим таблицам, не ко всем в аккаунте
        Resource = [
          aws_dynamodb_table.users.arn,
          "${aws_dynamodb_table.users.arn}/index/*",
          aws_dynamodb_table.messages.arn,
          "${aws_dynamodb_table.messages.arn}/index/*",
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "s3:PutObject",
          "s3:GetObject",
        ]
        Resource = "${aws_s3_bucket.files.arn}/*"
      },
    ]
  })
}


# Для работы github actions через OIDC
resource "aws_iam_openid_connect_provider" "github" {
  url             = "https://token.actions.githubusercontent.com"
  client_id_list  = ["sts.amazonaws.com"]
  thumbprint_list = ["1c58a3a8518e8759bf075b76b750d4f2df264fcd"]
}

resource "aws_iam_role" "github_actions" {
  name = "${var.project}-github-actions"

  # https://docs.github.com/ru/actions/how-tos/secure-your-work/security-harden-deployments/oidc-in-aws
  assume_role_policy = jsonencode({
    # Это не дата, просто странное версионирование
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = {
        Federated = aws_iam_openid_connect_provider.github.arn
      }
      Action = "sts:AssumeRoleWithWebIdentity"
      Condition = {
        StringEquals = {
          "token.actions.githubusercontent.com:aud" = "sts.amazonaws.com"
        }
        StringLike = {
          # чтобы только указанный репо мог юзать роль
          "token.actions.githubusercontent.com:sub" = "repo:alexe0110/ChatSystem:*"
        }
      }
    }]
  })
}

# Политика роли - что может делать тот, у кого эта роль
resource "aws_iam_role_policy" "github_actions" {
  name = "${var.project}-github-actions-policy"
  role = aws_iam_role.github_actions.id

  policy = jsonencode({
    Version = "2012-10-17"
    # Массив правил. Каждое правило отвечает на вопрос: "кому разрешить/запретить делать что с чем?"
    Statement = [
      # Можно логиниться в ECR
      {
        Effect = "Allow"
        Action = [
          "ecr:GetAuthorizationToken",
        ]
        Resource = "*"
      },
      # Можно пушить образы
      {
        Effect = "Allow"
        Action = [
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage",
          "ecr:PutImage",
          "ecr:InitiateLayerUpload",
          "ecr:UploadLayerPart",
          "ecr:CompleteLayerUpload",
        ]
        Resource = aws_ecr_repository.chat_system.arn
      },
      # Разрешить обновлять ECS-сервис (передеплоить)
      {
        Effect = "Allow"
        Action = [
          "ecs:UpdateService",
          "ecs:DescribeServices",
        ]
        Resource = "*"
      },
    ]
  })
}