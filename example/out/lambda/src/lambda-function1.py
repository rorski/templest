import boto3


def handler(event, context):
    string = "lorem ipsum"
    encoded_string = string.encode("utf-8")

    bucket_name = "my-bucket-1234abcd"
    s3_path = "path/123456/myfile.txt"

    s3 = boto3.resource("s3")
    s3.Bucket(bucket_name).put_object(Key=s3_path, Body=encoded_string)
