# 3. Follow Twelve Factor methodology

Date: 2018-07-03

## Status

Accepted

## Context

Designing modern scalable cloud based applications requires intentionally
designing the architecture to take advantage of the cloud.

One leading way to do that is
[The Twelve Factor](https://12factor.net) methodology.

## Decision

We will follow Twelve Factor methodology.

## Consequences

Our application will be deployable in the cloud in a scalable efficient manner.

We will leverage services for some aspects of applications that
previously would have relied on a Virtual Machine, such as storage for files
and logs.
