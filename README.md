# Mario

## What is this?

Mario is a metadata processing pipeline that will process data from various
sources and write to Elasticsearch.

## Installing

The `mario` command can be installed with:

```
$ go get github.com/mitlibraries/mario
```

## How to Use This

An Elasticsearch index can be started for development purposes by running:

```
$ docker run -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" \
    docker.elastic.co/elasticsearch/elasticsearch:6.4.2
```

Create and configure the index with:

```
$ mario create
```

The Mario container can be built and used by running:

```
$ docker build -t mario .
$ docker run --rm -i mario parse -c title - < fixtures/test.mrc
```

## Developing

This project uses modules for dependencies. To upgrade all dependencies to the latest minor/patch version use:

```
$ go get -u ./...
```

Tests can be run with:

```
$ go test -v ./...
```

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
