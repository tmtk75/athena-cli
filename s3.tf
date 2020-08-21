variable aws_profile {}
variable aws_region {}

provider aws {
  profile = var.aws_profile
  region  = var.aws_region
}

resource aws_athena_database main {
  name   = "test_cloudtrail"
  bucket = aws_s3_bucket.main.bucket
}

resource aws_s3_bucket main {
  bucket = "test-cloudtrail-athena"
  acl    = "private"
}

resource aws_athena_workgroup main {
  name = "test_cloudtrail_wg"
  configuration {
    result_configuration {
    }
  }
}

#
# This is an inconvenient feature in AWS. I don't use.
#
#resource aws_athena_named_query lo {
#  name      = "create table for cloudtrail"
#  workgroup = aws_athena_workgroup.main.id
#  database  = aws_athena_database.main.name
#  query     = file("./query.txt")
#}
#
#data aws_caller_identity current {}

