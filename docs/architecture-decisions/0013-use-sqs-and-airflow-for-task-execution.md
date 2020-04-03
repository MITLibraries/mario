# 13. Use SQS and Airflow for Task Execution

Date: 2020-04-03

## Status

Accepted

Supercedes [12. Use Lambda and Fargate for Task Execution](0012-use-lambda-and-fargate-for-task-execution.md)

## Context

The execution model described by ADR 12 was designed before we had Airflow. It works, but we'd like to simplify things by moving it to Airflow to avoid having similar processes handled in different ways.

## Decision

We will change the S3 notification from triggering a Lambda to sending a message to an SQS queue. We will configure a single workflow in Airflow that begins with an SQS sensor.

## Consequences

1. Instead of having a different indexing process for each source, there will be a single indexing workflow in Airflow that gets run on every file uploaded to the S3 bucket. Mario will need logic added to it to handle correctly indexing based only on the name of the bucket and key.

2. The Lambda process was a push process--as soon as the file was added to S3, it was processed. This new process will be a pull process. There will be a delay from when a file is added to S3 to when it is processed. This delay can be controlled by changing how often the workflow is run in Airflow.

3. SQS has a limit of 12 hours between when a message is received to when it can be deleted, meaning a single indexing process can't run for more than 12 hours.
