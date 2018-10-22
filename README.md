# Mario

## What is this?

Mario is a metadata processing pipeline that will process data from various
sources and write to Elasticsearch.

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

## System Overview
![alt text](docs/charts/dip_overview.png "Mario system overview chart")

## Architecture Overview
![alt text](docs/charts/dip_architecture.png "Mario system overview chart")

## Architecture Decision Records

This repository contains Architecture Decision Records in the
[docs/architecture-decisions directory](docs/architecture-decisions).

[adr-tools](https://github.com/npryce/adr-tools) should allow easy creation of
additional records with a standardized template.
