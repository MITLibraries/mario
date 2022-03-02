# Mario

## What is this?

Mario is a metadata processing pipeline that will process data from various
sources and write to Opensearch.

## Installing

The `mario` command can be installed with:

```
$ make install
```

## How to Use This

An OpenSearch index can be started for development purposes by running:

```
$ docker run -p 9200:9200 -p 9600:9600 -e "discovery.type=single-node" \
  -e "plugins.security.disabled=true" \
  opensearchproject/opensearch:1.2.4
```

Alternatively, if you intend to test this with a local instance of TIMDEX
as well, use docker-compose to run both docker and TIMDEX locally using the
instructions in the [TIMDEX README](https://github.com/MITLibraries/timdex/blob/master/README.md#docker-compose-orchestrated-local-environment).

Here are a few sample Mario commands that may be useful for local development:
- `mario ingest -c json -s aspace fixtures/aspace_samples.xml`
  runs the ingest process with ASpace sample files and prints out each record
  as JSON
- `mario ingest -s dspace fixtures/dspace_samples.xml` ingests the
  DSpace sample files into a local OpenSearch instance.
- `mario ingest -s alma --auto fixtures/alma_samples.mrc` ingests the
  Alma sample files into a local OpenSearch instance and promotes the
  index to the timdex-prod alias on completion.
- `mario indexes` list all indexes
- `mario promote -i [index name]` promotes the named index to the
  timdex-prod alias.

## Developing

This project uses modules for dependencies. To upgrade all dependencies to the latest minor/patch version use:

```
$ make update
```

Tests can be run with:

```
$ make test
```

### Adding a new source parser
To add a new source parser:
- (Probably) create a source record struct in `pkg/generator`.
- Add a source parser module in `pkg/generator`.
- Add a tests file that tests ALL fields mapped from the source.
- Update `pkg/ingester/ingester.go` to add a Config.source that uses the new
  generator.
- Update documentation to include the new generator param option (as "type") to
  command options.
- (Probably) don’t need to update the CLI.
- After all of that is completed, tested, and merged, create tasks to harvest
  the source metadata files and ingest them using our [airflow implementation](https://github.com/MITLibraries/workflow).

### Updating the data model
Updating the data model is somewhat complicated because many files need to be
edited across multiple repositories and deployment steps should happen in a
particular order so as not to break production services. Start by updating the data model here in Mario as follows:
- Update `config/es_record_mappings.json` to reflect added/updated/deleted
  fields.
- Update `pkg/record/record.go` to reflect added/updated/deleted fields.
- Update ALL relevant source record definitions and source parser files in
  `pkg/generator`. If a field is edited or deleted, be sure to check every
  source file for usage. If a field is new, add to all relevant sources
  (confirm mapping with metadata folks first).
- Update relevant tests in `pkg/generator`.
- Once the above steps are done, update the data model in TIMDEX following the
  instructions in the [TIMDEX README](https://github.com/MITLibraries/timdex/blob/master/README.md) and test locally with the docker-compose
  orchestrated environment to ensure all changes are properly indexed and
  consumable via the API.

## Config Files
We have several config files that are essential for mapping various metadata
field codes to their human-readable translations, and some of them may need to
be updated from time to time. Most of these config files are pulled from
authoritative sources, with the exception of `marc_rules.json` which we created
and manage ourselves. Sources of the other config files are as follows:

- `dspace_set_list.json` this is harvested from our DSpace repository using our
  OAI-PMH harvester app. The app includes a flag to convert the standard XML
  response to JSON, which just makes it easier to parse.

## System Overview
![alt text](docs/charts/dip_overview.png "Mario system overview chart")

## Architecture Overview
![alt text](docs/charts/dip_architecture.png "Mario system overview chart")

## Architecture Decision Records

This repository contains Architecture Decision Records in the
[docs/architecture-decisions directory](docs/architecture-decisions).

[adr-tools](https://github.com/npryce/adr-tools) should allow easy creation of
additional records with a standardized template.
