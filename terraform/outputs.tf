output "s3_bucket_name" {
  value = aws_s3_bucket.files.bucket
}

output "ecr_repository_url" {
  value = aws_ecr_repository.chat_system.repository_url
}

output "dynamodb_messages_table" {
  value = aws_dynamodb_table.messages.name
}