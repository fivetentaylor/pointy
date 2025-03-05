resource "aws_dynamodb_table" "reviso" {
  name         = var.dynamo_table_name
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "PK"
  range_key    = "SK"

  attribute {
    name = "PK"
    type = "S"
  }

  attribute {
    name = "SK"
    type = "S"
  }

  attribute {
    name = "SK1"
    type = "S"
  }

  attribute {
    name = "SK2"
    type = "S"
  }

  attribute {
    name = "SK3"
    type = "S"
  }

  attribute {
    name = "SK4"
    type = "S"
  }

  attribute {
    name = "SK5"
    type = "S"
  }

  attribute {
    name = "GSI1PK"
    type = "S"
  }

  attribute {
    name = "GSI1SK"
    type = "S"
  }

  local_secondary_index {
    name            = "SK1Index"
    range_key       = "SK1"
    projection_type = "KEYS_ONLY"
  }

  local_secondary_index {
    name            = "SK2Index"
    range_key       = "SK2"
    projection_type = "KEYS_ONLY"
  }

  local_secondary_index {
    name            = "SK3Index"
    range_key       = "SK3"
    projection_type = "KEYS_ONLY"
  }

  local_secondary_index {
    name            = "SK4Index"
    range_key       = "SK4"
    projection_type = "KEYS_ONLY"
  }

  local_secondary_index {
    name            = "SK5Index"
    range_key       = "SK5"
    projection_type = "KEYS_ONLY"
  }

  global_secondary_index {
    name            = "GSI1"
    hash_key        = "GSI1PK"
    range_key       = "GSI1SK"
    projection_type = "ALL"
  }

  tags = {
    Name = var.dynamo_table_name
  }
}

