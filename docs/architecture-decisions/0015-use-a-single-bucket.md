# 15. Use a Single Bucket

Date: 2020-04-08

## Status

Accepted

Supercedes [8. Use One S3 Bucket Per Source](0008-use-one-s3-bucket-per-source.md)

## Context

There are a couple reasons to switch from multiple buckets to a single bucket. The first is that it simplifies the infrastructure provisioning that needs to be done when a new source is added. The second is that there is a hard limit on the number of buckets an AWS account can have, and our current approach to bucket creation is unsustainable.

## Decision

Use a single namespaced S3 bucket for all source data. The structure of the bucket should be:

```
s3://bucket/<environment>/<source>/<files>
```

Where `environment` would be either `prod` or `stage`, and `source` would be the source identifier. The source identifier used here should also be used as the prefix for the index name. No specific decisions are made here about how the files are structured within a source.

mario should ignore source identifiers that it does not know about.

## Consequences

Changes will be needed both in mario and in the scripts that upload data to the current bucket. It's worth noting that bucket policies around lifecycle and permissions can be applied to objects based on prefix, so any existing needs specific to a source can be supported with this change.
