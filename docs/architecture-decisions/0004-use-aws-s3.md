# 4. Use AWS S3

Date: 2018-07-03

## Status

Accepted

## Context

One of the tenants of Twelve Factor application design is that applications
should be stateless, which includes not relying on local file storage to be
persistent. As such, this project needs a cloud based object store.
See [3. Follow Twelve Factor methodology](0003-follow-twelve-factor-methodology.md)

Amazon Simple Storage Service (S3) is a secure, durable, and scalable object
storage solution to use in the cloud with which we have established an existing
payment relationship.

Amazon provides official SDKs for various programming languages to interact
with S3.

## Decision

We will use Amazon S3 for our object store.

## Consequences

We will have secure, durable, and scalable object storage to use in the cloud
as needed for this project.
