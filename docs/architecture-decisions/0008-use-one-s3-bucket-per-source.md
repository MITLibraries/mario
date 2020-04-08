# 8. Use One S3 Bucket Per Source

Date: 2018-07-05

## Status

Superceded by [15. Use a Single Bucket](0015-use-a-single-bucket.md)

## Context

Each data source will need to upload one or more files to S3 in order to trigger processing. S3 events, which will drive Lambda execution (See [5. Use AWS Lambda](0005-use-aws-lambda.md)), are configured at the bucket level. We may or may not have much control over the environment which is sending data to S3, for example, if it came directly from a vendor. At minimum we must be able to specify a bucket, but we should not assume we will have much more control than this.

Each data source will also need different processing. This implies the need to identify which source a data file came from.

## Decision

Use one S3 bucket per data source.

## Consequences

Each source will need a new S3 bucket. Each bucket must be configured to tie the Lambda function to the object creation event. An advantage to this approach is that it makes it easier to enable and disable sources by simply controlling whether or a not a bucket publishes its events. No changes to application code would be necessary.

The bucket name can be used as the configuration key to identify the source and determine how the data should be processed.
