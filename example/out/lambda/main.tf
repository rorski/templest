module "lambda_function" {
  source = "<no value>"

  function_name = "my-lambda1"
  description   = "Some lambda function"
  handler       = "index.handler"
  runtime       = "python3.8"

  source_path = "./src/lambda-function1"

  tags = {
    Name = "my-lambda1"
  }
}
