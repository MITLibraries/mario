# 12. Use Lambda and Fargate for Task Execution

Date: 2019-01-31

## Status

Supercedes [5. Use AWS Lambda](0005-use-aws-lambda.md)

Superceded by [13. Use SQS and Airflow for Task Execution](0013-use-sqs-and-airflow-for-task-execution.md)

## Context

Given the limitations of Lambdas we decided to rely on containers to handle the bulk of the processing. Fargate provides a cheap, accessible container runtime.

## Decision

We will use AWS Lambda to trigger a Fargate task for the processing pipeline.

## Consequences

Unfortunately, there's no easy way to trigger a Fargate task from an S3 file upload. The S3 upload will have to trigger a Lambda that can then run the Fargate task, setting the filename as the container's runtime parameters.
