# 9. Elasticsearch Indexing Strategy

Date: 2018-07-06

## Status

Accepted

## Context

There are a number of different ways we could approach indexing in Elasticsearch. We would like to choose a path that allows us some flexibility to adjust as future needs arise. We also need to think about how to maintain index uptime while modifying the contents of the index.

## Decision

Use an index alias for searching that points to a separate index for each source.

## Consequences

An index alias provides a constant, unchanging endpoint for searches which minimizes the integration impact of modifications to index structure. Changing which indexes the alias points to is an atomic action allowing for smooth transitions to different versions of indexes with no downtime.

Using one index per source allows us to further isolate the impact of bringing new sources online and modifying how different sources are indexed.

The process for indexing (and reindexing) a source would generally follow these steps:

1. Create a new index, using some kind of versioning in the name.
2. Add documents to the new index.
3. Modify the alias to add the new index and remove the old index from its pointers.
4. Delete the old index.
