terraform {
	required_providers {
		aws = {
			source  = "hashicorp/aws"
			version = "6.17.0"
		}
	}
}
provider "aws" {
	region = "eu-central-1"
}

# ─── S3 ───────────────────────────────────────
resource "aws_s3_bucket" "files" {
	bucket = "chatsystem-files-tf-alexe0110"
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
}

# ─── ECR ──────────────────────────────────────
resource "aws_ecr_repository" "chat_system" {
	name                 = "chat-system"
	image_tag_mutability = "MUTABLE"
	force_delete         = true        # чтобы terraform destroy удалил даже с образами внутри
}