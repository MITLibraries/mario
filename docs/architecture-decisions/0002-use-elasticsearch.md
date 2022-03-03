# 2. Use Elasticsearch

Date: 2018-07-03

## Status

Superceded by [16. Replace Elasticsearch with OpenSearch](0016-replace-elasticsearch-with-opensearch.md)

## Context

We need to choose between using Solr and Elasticsearch for indexing.

## Decision

We will use Elasticsearch. See https://docs.google.com/document/d/1LX3svZ59f2Ni5TNCPG6jIYb8CnSYOjR0ae0ujPOUN-k/edit for a more detailed description of how this decision was arrived at.

## Consequences

Because Solr and Elasticsearch are so different we won't be able to easily switch from one to the other at a later point. Choosing Elasticsearch means we can't make use of any of the Blacklight family of applications since it is designed to work with Solr.

We expect to have more hosted options available to us with Elasticsearch.
