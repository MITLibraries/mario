# 14. Structure of CLI

Date: 2020-04-08

## Status

Accepted

Supercedes [11. Indexing Commands and Flows](0011-indexing-commands-and-flows.md)

## Context

Our current command line interface is awkward to use and contains a considerable amount of business logic. Rather than listing out every command, this ADR will provide guidelines for adding new commands.

## Decision

The mario command itself will be a collection of subcommands. Only use global options for values that can truly be applied for every subcommand. All other options should be attached to the subcommand. Provide reasonable default values when it makes sense.

A few examples:

```
$ mario ingest --v4 --index aleph-2020-01-01 s3://bucket/key.mrc
$ mario reindex --url http://example.com -s aleph-01 -d aleph-02
```

Additionally, the `main.go` file should be kept small and all business logic should reside elsewhere in the application.

## Consequences

The `main.go` file will need a significant rewrite in order to pull out the business logic which is currently in there. The CLI will also likely change (`--index` is currently a global option, which makes no sense for most subcommands). Changes to the CLI will need to be propagated to the automated ingest workflow processes.
