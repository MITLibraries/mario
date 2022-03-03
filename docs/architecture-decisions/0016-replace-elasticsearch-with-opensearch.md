# 16. Replace Elasticsearch with OpenSearch

Date: 2022-03-02

## Status

Accepted

Supercedes [2. Use Elasticsearch](0002-use-elasticsearch.md)

## Context

Amazon has moved to support OpenSearch over Elasticsearch, and no longer provides Elasticsearch as a managed service past version 7.10. Due to our heavy use of AWS services, it makes sense to use the indexing service that will be supported going forward. See https://aws.amazon.com/blogs/aws/amazon-elasticsearch-service-is-now-amazon-opensearch-service-and-supports-opensearch-10/ for background information.

## Decision

Use the latest version of OpenSearch as our index instead of Elasticsearch.

## Consequences

Currently switching to OpenSearch has very little consequences, as the current version is still nearly identical to the version of Elasticsearch it was forked from (7.10). As OpenSearch evolves and we make adjustments to follow, it may become more difficult to move back to Elasticsearch should we ever desire to.

The documentation for OpenSearch is pretty minimal, however for the same reason as above that isn't very consequential yet since we can still use Elasticsearch documentation, which is very robust.

This doesn't impact other decisions made about our indexing strategy and flow, as they remain the same in OpenSearch.