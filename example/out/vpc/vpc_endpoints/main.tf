module "vpc-endpoints-s3" {
  source     = "path-to-module/vpc-endpoints/s3"
  env        = var.env
  vpc_id     = module.my-vpc.vpc.id
  subnet_ids = data.aws_subnet_ids.private.ids
  create_dns = true
}

data "aws_subnet_ids" "private" {
  vpc_id = module.my-vpc.vpc.id

  filter {
    name = "tag:Env"
    value = var.env
  }
}

variable "env" {
  type string
}