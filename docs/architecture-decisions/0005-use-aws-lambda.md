# 5. Use AWS Lambda

Date: 2018-07-03

## Status

Superceded by [12. Use Lambda and Fargate for Task Execution](0012-use-lambda-and-fargate-for-task-execution.md)

## Context

The bulk of this application will consist of a data processing pipeline that takes metadata from incoming systems and indexes it in Elasticsearch. The processing will only need to be run for relatively short periods of time, usually, when new data arrives. We expect integrations with external systems to be minimal, likely limited only to S3 and Elasticsearch. Given the periodic nature of the application, it seems wasteful and needlessly complex to provision and maintain a VM for providing compute resources.

## Decision

We will use AWS Lambdas as the compute model for the processing pipeline.

## Consequences

Since we are using S3 for storage, the Lambda can be easily configured to run when a new file is placed in an S3 bucket. Lambda supports several different languages. While triggering Lambdas through S3 events will be convenient for most cases, it does potentially add some complexity in situations that require running Lambdas outside of the normal event structure.

Integration with central logging will likely be more complex as Lambda logs currently go directly to CloudWatch.
