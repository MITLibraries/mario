# 10. Use OpenAPI Specification

Date: 2018-08-10

## Status

Accepted

## Context

By choosing an API documentation standard we make it easier to auto generate developer documentation. There are two existing standards used to document REST APIs--RAML and OpenAPI Specification (OAS). Both seem capable of doing the job. The main difference seems to be that RAML is focused on defining data models while OpenAPI is focused on the nuts and bolts of the API. If we were supporting several APIs, RAML might be more useful for defining reusable types across systems. In this case OAS seems more suited to our task.

## Decision

Use OpenAPI specification to document the API.

## Consequences

Given that Mulesoft uses RAML and Swagger uses OAS, this means we'd be probably be using Swagger as a documentation platform (if we choose to do so). There are tools to convert between RAML and OAS, though, so if we decide we would rather use RAML later it should not be too difficult to switch.
