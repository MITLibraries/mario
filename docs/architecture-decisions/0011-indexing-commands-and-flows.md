# 11. Indexing Commands and Flows

Date: 2018-11-26

## Status

Superceded by [14. Structure of CLI](0014-structure-of-cli.md)

## Context

Our AWS instance of Elasticsearch will be awkward to maintain unless we build the right tools into the CLI. Some of these commands will not be used often, but are anticipatory of being able to intervene and correct issues in production. As such, these are the CLI commands we intend to make available to developers / maintainers of the Discovery Index. Future needs may dictate additional commands or adjustments to these.

## Decision

#### COMMANDS:

    indexes  List Elasticsearch indexes and aliases

    ingest   Parse and ingest the input file
             [note: was parse. Ingest is more reflective of what this does as parsing is just one aspect]

    ping     Request and display general info about the Elasticsearch server

    help, h  Shows a list of commands or help for one command

  Index actions:

    create   Currently available, removed in this proposal.

    delete   Delete an Elasticsearch index

    demote   Remove the given index from the production alias

    promote  Add the given index to production alias

    stats    Stats for provided Index (total records, maybe more?)


In order to ensure we always have a production index available, we want to use a known, semi-hardcoded alias value. The fixed value will be "production".

As we'll have new sources soon, it will be best to have each source maintain their own indexes for ease of use. Elasticsearch will allow multiple indexes to be associated with the alias "production" to allow us to search as many sources as we need. For aleph, the prefix will be `aleph_`. Future sources will declare a prefix that is appropriate as we add them following the `source_` convention.

We currently default to an index name of `timdex`. We should no longer do that and instead default to different values depending on the command. Some commands will not have a default at all and are detailed below.

We currently allow for an independent `create` command. That will be removed and the `ingest` command will handle creating indexes if necessary.

For interactive work, we allow specific index values to be passed in as it may be useful for either intervening in a problem in production, or general local development. A default index value for `ingest` will be a combination of the source, such as `aleph_` and a partial timestamp such as `2018_11_26_1001`.

Once we are confident the new index is ready, we'd then `promote` that index while `demoting` any existing indexes for the source we are working with. `promote` and `demote` are used to signify adding / removing the index to / from the alias "production". `promote` and `demote` will not default to any index value to ensure intent. When in fully automated mode, we'll need to ensure the index value set during the `ingest` process is used to `promote` and the source is used to `demote` as part of the single atomic action. However, that will not be done via the `promote` / `demote` CLI interface and will instead be kicked off as part of the `ingest` process when run in `auto` mode. The `index` argument will be required to inform these command what index to operate on. Additionally, we will ensure we do not allow demoting and index if it will leave the `production` alias with no indexes for the specific source. In other words, we will programmatically ensure that at least one `aleph_` source is always accessible via the `production` alias. Once we add additional sources, they will follow the same requirement of always having one index with the alias of `production`.

We will keep as many old indexes as we deem useful and then go back and `delete` indexes we no longer need. The `delete` process will ensure that the index is not assigned to the alias `production` before proceeding. Additionally, `delete` will not default to any index and one must be supplied to continue.

Those set of commands will allow us to fully manually run the pipeline in production. The following is a proposed automatic flow:

```
mario ingest aleph --promote auto --url esurl --v4 file
```

The source argument (`aleph` above) will allow us to construct the new index automatically when used in conjunction with a current timestamp.

We'll also use the source to check the alias `production` for any indexes that currently have that prefix. `--promote auto` would dd the new index to production alias as well as remove the old index from the production alias. If `--promote auto` is not set, the new index will be created but the old index would remain in place for production use until further interactive steps were taken.

## Consequences

This set of commands and flows will allow us to setup automatic processes while still maintaining manual ability to intervene if necessary for both the current aleph dataset as well as future datasets as we start adding them to the system.
