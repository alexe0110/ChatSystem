terraform {
	required_providers {
		aws = {
			source  = "hashicorp/aws"
			version = "6.17.0"
		}
	}
}
provider "aws" {
	region = var.region
}

# ─── S3 ───────────────────────────────────────
resource "aws_s3_bucket" "files" {
	bucket = "${var.project}-files-${var.environment}"
}

# ─── DynamoDB ─────────────────────────────────
resource "aws_dynamodb_table" "users" {
	name = "users"
	billing_mode = "PAY_PER_REQUEST"
	hash_key = "user_id"

	attribute {
		name = "user_id"
		type = "S"
	}

	attribute {
		name = "login"
		type = "S"
	}

	global_secondary_index {
		hash_key        = "login"
		name            = "login-index"
		projection_type = "ALL"
	}
}

resource "aws_dynamodb_table" "messages" {
	name         = "messages"
	billing_mode = "PAY_PER_REQUEST"
	hash_key     = "chat_id"
	range_key    = "created_at"        # range_key = sort key

	attribute {
		name = "chat_id"
		type = "S"
	}

	attribute {
		name = "created_at"
		type = "S"
	}

	attribute {
		name = "message_id"
		type = "S"
	}

	global_secondary_index {
		hash_key        = "message_id"
		name            = "message-id-index"
		projection_type = "ALL"
	}
}

# ─── ECR ──────────────────────────────────────
resource "aws_ecr_repository" "chat_system" {
	name                 = "chat-system"
	image_tag_mutability = "MUTABLE"
	force_delete         = true        # чтобы terraform destroy удалил даже с образами внутри
}


# ─── ECS ──────────────────────────────────────
resource "aws_ecs_cluster" "main" {
	name = "${var.project}-cluster"
}